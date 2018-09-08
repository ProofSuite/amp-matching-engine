package operator

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strconv"

	"github.com/Proofsuite/amp-matching-engine/interfaces"
	"github.com/Proofsuite/amp-matching-engine/rabbitmq"
	"github.com/Proofsuite/amp-matching-engine/types"
	"github.com/Proofsuite/amp-matching-engine/utils"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	eth "github.com/ethereum/go-ethereum/core/types"
	"github.com/streadway/amqp"
)

// Operator manages the transaction queue that will eventually be
// sent to the exchange contract. The Operator Wallet must be equal to the
// account that initially deployed the exchange contract or an address with operator rights
// on the contract
type Operator struct {
	WalletService     interfaces.WalletService
	TradeService      interfaces.TradeService
	OrderService      interfaces.OrderService
	EthereumProvider  interfaces.EthereumProvider
	Exchange          interfaces.Exchange
	TxQueues          []*TxQueue
	QueueAddressIndex map[common.Address]*TxQueue
}

type OperatorInterface interface {
	SubscribeOperatorMessages(fn func(*types.OperatorMessage) error) error
	QueueTrade(o *types.Order, t *types.Trade) error
	GetShortestQueue() (*TxQueue, int, error)
	SetFeeAccount(account common.Address) (*eth.Transaction, error)
	SetOperator(account common.Address, isOperator bool) (*eth.Transaction, error)
	FeeAccount() (common.Address, error)
	Operator(addr common.Address) (bool, error)
}

// NewOperator creates a new operator struct. It creates an exchange contract instance from the
// provided address. The error and trade events are received in the ErrorChannel and TradeChannel.
// Upon receiving errors and trades in their respective channels, event payloads are sent to the
// associated order maker and taker sockets through the through the event channel on the Order and Trade struct.
// In addition, an error event cancels the trade in the trading engine and makes the order available again.
func NewOperator(
	walletService interfaces.WalletService,
	tradeService interfaces.TradeService,
	orderService interfaces.OrderService,
	provider interfaces.EthereumProvider,
	exchange interfaces.Exchange,
) (*Operator, error) {
	txqueues := []*TxQueue{}
	addressIndex := make(map[common.Address]*TxQueue)

	wallets, err := walletService.GetOperatorWallets()
	if err != nil {
		panic(err)
	}

	for i, w := range wallets {
		name := strconv.Itoa(i) + w.Address.Hex()
		ch := rabbitmq.GetChannel("TX_QUEUES:" + name)
		err := rabbitmq.DeclareQueue(ch, "TX_QUEUES:"+name)
		if err != nil {
			panic(err)
		}

		txq := &TxQueue{
			Name:             name,
			Wallet:           w,
			TradeService:     tradeService,
			EthereumProvider: provider,
			Exchange:         exchange,
		}

		txqueues = append(txqueues, txq)
	}

	op := &Operator{
		WalletService:     walletService,
		TradeService:      tradeService,
		OrderService:      orderService,
		EthereumProvider:  provider,
		Exchange:          exchange,
		TxQueues:          txqueues,
		QueueAddressIndex: addressIndex,
	}

	err = op.PurgeQueues()
	if err != nil {
		log.Print(err)
		return nil, err
	}

	go op.HandleEvents()
	return op, nil
}

// SubscribeOperatorMessages
func (op *Operator) SubscribeOperatorMessages(fn func(*types.OperatorMessage) error) error {
	ch := rabbitmq.GetChannel("OPERATOR_SUB")
	q := rabbitmq.GetQueue(ch, "TX_MESSAGES")

	go func() {
		msgs, err := ch.Consume(
			q.Name,
			"",
			true,
			false,
			false,
			false,
			nil,
		)

		if err != nil {
			log.Fatalf("Failed to register a consumer: %s", err)
		}

		forever := make(chan bool)

		go func() {
			for m := range msgs {
				om := &types.OperatorMessage{}
				err := json.Unmarshal(m.Body, &om)
				if err != nil {
					log.Printf("Error: %v", err)
					continue
				}
				go fn(om)
			}
		}()

		<-forever
	}()
	return nil
}

// Bug: In certain cases, the trade channel seems to be receiving additional unexpected trades.
// In the case TestSocketExecuteOrder (in file socket_test.go) is run on its own, everything is working correctly.
// However, in the case TestSocketExecuteOrder is run among other tests, some tradeLogs do not correspond to an
// order hash in the ordertrade mapping. I suspect this is because the event listener catches events from previous
// tests. It might be helpful to see how to listen to events from up to a certain block.
func (op *Operator) HandleEvents() error {
	tradeEvents, err := op.Exchange.ListenToTrades()
	if err != nil {
		log.Print(err)
		return err
	}

	errorEvents, err := op.Exchange.ListenToErrors()
	if err != nil {
		log.Print(err)
		return err
	}

	for {
		select {
		case event := <-errorEvents:
			fmt.Println("TRADE_ERROR_EVENT")
			tradeHash := event.TradeHash
			errID := int(event.ErrorId)

			log.Print("The error ID is: ", errID)

			tr, err := op.TradeService.GetByHash(tradeHash)
			if err != nil {
				log.Print(err)
			}

			or, err := op.OrderService.GetByHash(tr.OrderHash)
			if err != nil {
				log.Print(err)
			}

			go func() {
				err = op.PublishTxErrorMessage(tr, errID)
				if err != nil {
					log.Print(err)
				}

				err = op.PublishTradeCancelMessage(or, tr)
				if err != nil {
					log.Print(err)
				}
			}()

		case event := <-tradeEvents:
			tr, err := op.TradeService.GetByHash(event.TradeHash)
			if err != nil {
				log.Print(err)
			}

			fmt.Println("TRADE_SUCCESS_EVENT", tr.Hash.Hex())

			or, err := op.OrderService.GetByHash(tr.OrderHash)
			if err != nil {
				log.Print(err)
			}

			go func() {
				_, err := op.EthereumProvider.WaitMined(tr.Tx)
				if err != nil {
					log.Print(err)
				}

				fmt.Println("TRADE_MINED IN HANDLE EVENTS: ", tr.Hash.Hex())

				err = op.PublishTradeSuccessMessage(or, tr)
				if err != nil {
					log.Print(err)
				}
			}()
		}
	}
}

func (op *Operator) HandleTrades(msg *types.OperatorMessage) error {

	log.Print("Handling trade")
	utils.PrintJSON(msg)

	err := op.QueueTrade(msg.Order, msg.Trade)
	if err != nil {
		log.Print(err)
		return err
	}

	return nil
}

// PublishTxErrorMessage publishes a messages when a trade execution fails
func (op *Operator) PublishTxErrorMessage(tr *types.Trade, errID int) error {
	ch := rabbitmq.GetChannel("OPERATOR_PUB")
	q := rabbitmq.GetQueue(ch, "TX_MESSAGES")
	msg := &types.OperatorMessage{
		MessageType: "TRADE_ERROR_MESSAGE",
		Trade:       tr,
		ErrID:       errID,
	}

	bytes, err := json.Marshal(msg)
	if err != nil {
		log.Printf("Failed to marshal %s: %s", msg.MessageType, err)
	}

	err = rabbitmq.Publish(ch, q, bytes)
	if err != nil {
		log.Print(err)
		return err
	}

	return nil
}

// PublishTradeCancelMessage publishes a message when a trade is canceled
func (op *Operator) PublishTradeCancelMessage(o *types.Order, tr *types.Trade) error {
	ch := rabbitmq.GetChannel("OPERATOR_PUB")
	q := rabbitmq.GetQueue(ch, "TX_MESSAGES")
	msg := &types.OperatorMessage{
		MessageType: "TRADE_CANCEL_MESSAGE",
		Trade:       tr,
	}

	bytes, err := json.Marshal(msg)
	if err != nil {
		log.Printf("Failed to marshal %s: %s", msg.MessageType, err)
	}

	err = rabbitmq.Publish(ch, q, bytes)
	if err != nil {
		log.Print(err)
		return err
	}

	return nil
}

// PublishTradeSuccessMessage publishes a message when a trade transaction is successful
func (op *Operator) PublishTradeSuccessMessage(o *types.Order, tr *types.Trade) error {
	ch := rabbitmq.GetChannel("OPERATOR_PUB")
	q := rabbitmq.GetQueue(ch, "TX_MESSAGES")
	msg := &types.OperatorMessage{
		MessageType: "TRADE_SUCCESS_MESSAGE",
		Order:       o,
		Trade:       tr,
	}

	bytes, err := json.Marshal(msg)
	if err != nil {
		log.Printf("Failed to marshal %s: %s", msg.MessageType, err)
	}

	err = rabbitmq.Publish(ch, q, bytes)
	if err != nil {
		log.Print(err)
		return err
	}

	return nil
}

// Publish
func (op *Operator) Publish(ch *amqp.Channel, q *amqp.Queue, bytes []byte) error {
	err := ch.Publish(
		"",
		q.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/json",
			Body:        bytes,
		},
	)

	if err != nil {
		log.Print(err)
		return err
	}

	return nil
}

// QueueTrade
func (op *Operator) QueueTrade(o *types.Order, t *types.Trade) error {
	txq, len, err := op.GetShortestQueue()
	if err != nil {
		log.Print(err)
		return err
	}

	if len > 10 {
		log.Print("Transaction queue is full")
		return errors.New("Transaction queue is full")
	}

	log.Print("QUEING TRADE", len)
	err = txq.QueueTrade(o, t)
	if err != nil {
		log.Print(err)
		return err
	}

	return nil
}

// GetShortestQueue
func (op *Operator) GetShortestQueue() (*TxQueue, int, error) {
	shortest := &TxQueue{}
	min := 1000

	for _, txq := range op.TxQueues {
		if shortest == nil {
			shortest = txq
			min = txq.Length()
		}

		ln := txq.Length()
		if ln < min {
			shortest = txq
			min = ln
		}
	}

	return shortest, min, nil
}

// SetFeeAccount sets the fee account of the exchange contract. The fee account receives
// the trading fees whenever a trade is settled.
func (op *Operator) SetFeeAccount(account common.Address) (*eth.Transaction, error) {
	txOpts, err := op.GetTxSendOptions()
	if err != nil {
		return nil, err
	}

	tx, err := op.Exchange.SetFeeAccount(account, txOpts)
	if err != nil {
		return nil, err
	}

	return tx, nil
}

// SetOperator updates the operator settings of the given address. Only addresses with an
// operator access can execute Withdraw and Trade transactions to the Exchange smart contract
func (op *Operator) SetOperator(account common.Address, isOperator bool) (*eth.Transaction, error) {
	txOpts, err := op.GetTxSendOptions()
	if err != nil {
		return nil, err
	}

	tx, err := op.Exchange.SetOperator(account, isOperator, txOpts)
	if err != nil {
		return nil, err
	}

	return tx, nil
}

// FeeAccount is the Ethereum towards the exchange trading fees are sent
func (op *Operator) FeeAccount() (common.Address, error) {
	account, err := op.Exchange.FeeAccount()
	if err != nil {
		return common.Address{}, err
	}

	return account, nil
}

// Operator returns true if the given address is an operator of the exchange and returns false otherwise
func (op *Operator) Operator(addr common.Address) (bool, error) {
	isOperator, err := op.Exchange.Operator(addr)
	if err != nil {
		return false, err
	}

	return isOperator, nil
}

func (op *Operator) PurgeQueues() error {
	for _, txq := range op.TxQueues {
		err := txq.PurgePendingTrades()
		if err != nil {
			return err
		}
	}

	return nil
}

func (op *Operator) GetTxSendOptions() (*bind.TransactOpts, error) {
	wallet, err := op.WalletService.GetDefaultAdminWallet()
	if err != nil {
		return nil, err
	}

	return bind.NewKeyedTransactor(wallet.PrivateKey), nil
}

// addressindex := op.QueueAddressIndex
// maker := o.UserAddress
// taker := t.Taker

// txq, ok := addressindex[maker]
// if ok {
// 	if txq.Length() < 10 {
// 		err := txq.QueueTrade(o, t)
// 		if err != nil {
// 			log.Print(err)
// 			return err
// 		}
// 		return nil
// 	} else {
// 		return errors.New("User transaction queue full")
// 	}
// }

// txq, ok = addressindex[taker]
// if ok {
// 	if txq.Length() < 10 {
// 		err := txq.QueueTrade(o, t)
// 		if err != nil {
// 			log.Print(err)
// 			return err
// 		}
// 		return nil
// 	} else {
// 		log.Print("Transaction queue is full")
// 		return errors.New("User transaction queue full")
// 	}
// }

// func NewDefaultOperator(
// 	walletService interfaces.WalletService,
// 	tradeService interfaces.TradeService,
// 	orderService interfaces.OrderService,
// 	provider interfaces.EthereumProvider,
// ) (*Operator, error) {
// 	exchangeAddress := common.HexToAddress(app.Config.Ethereum["exchange_address"])
// 	txqueues := []*TxQueue{}
// 	addressIndex := make(map[common.Address]*TxQueue)

// 	wallets, err := walletService.GetOperatorWallets()
// 	if err != nil {
// 		panic(err)
// 	}

// 	exchange, err := contractsinterfaces.NewExchange(exchangeAddress, provider)
// 	if err != nil {
// 		return nil, err
// 	}

// 	for i, w := range wallets {
// 		name := strconv.Itoa(i) + w.Address.Hex()
// 		txq := &TxQueue{
// 			Name:             name,
// 			Wallet:           w,
// 			TradeService:     tradeService,
// 			EthereumProvider: provider,
// 			Exchange:         exchange,
// 		}

// 		txqueues = append(txqueues, txq)
// 	}

// 	op := &Operator{
// 		WalletService:     walletService,
// 		TradeService:      tradeService,
// 		OrderService:      orderService,
// 		EthereumProvider:  provider,
// 		Exchange:          exchange,
// 		TxQueues:          txqueues,
// 		QueueAddressIndex: addressIndex,
// 	}

// 	go op.HandleEvents()

// 	return op, nil
// }
