package operator

import (
	"encoding/json"
	"errors"
	"log"
	"strconv"

	"github.com/Proofsuite/amp-matching-engine/contracts"
	"github.com/Proofsuite/amp-matching-engine/rabbitmq"
	"github.com/Proofsuite/amp-matching-engine/services"
	"github.com/Proofsuite/amp-matching-engine/types"
	"github.com/ethereum/go-ethereum/common"
	eth "github.com/ethereum/go-ethereum/core/types"
)

// Operator manages the transaction queue that will eventually be
// sent to the exchange contract. The Operator Wallet must be equal to the
// account that initially deployed the exchange contract or an address with operator rights
// on the contract
type Operator struct {
	WalletService     services.WalletServiceInterface
	TxService         services.TxServiceInterface
	TradeService      services.TradeServiceInterface
	EthereumService   services.EthereumServiceInterface
	Exchange          *contracts.Exchange
	TxQueues          []*TxQueue
	QueueAddressIndex map[common.Address]*TxQueue
}

type OperatorInterface interface {
	SubscribeOperatorMessages(fn func(*OperatorMessage) error) error
	QueueTrade(o *types.Order, t *types.Trade) error
	GetShortestQueue() (*TxQueue, int, error)
	SetFeeAccount(account common.Address) (*eth.Transaction, error)
	SetOperator(account common.Address, isOperator bool) (*eth.Transaction, error)
	FeeAccount() (common.Address, error)
	Operator(addr common.Address) (bool, error)
}

type OperatorMessage struct {
	MessageType string
	Order       *types.Order
	Trade       *types.Trade
	ErrID       int
}

type PendingTradeMessage struct {
	Order *types.Order
	Trade *types.Trade
}

// NewOperator creates a new operator struct. It creates an exchange contract instance from the
// provided address. The error and trade events are received in the ErrorChannel and TradeChannel.
// Upon receiving errors and trades in their respective channels, event payloads are sent to the
// associated order maker and taker sockets through the through the event channel on the Order and Trade struct.
// In addition, an error event cancels the trade in the trading engine and makes the order available again.
func NewOperator(
	walletService services.WalletServiceInterface,
	txService services.TxServiceInterface,
	tradeService services.TradeServiceInterface,
	ethereumService services.EthereumServiceInterface,
	exchange *contracts.Exchange,
) (*Operator, error) {

	txqueues := []*TxQueue{}
	addressIndex := make(map[common.Address]*TxQueue)

	wallets, err := walletService.GetOperatorWallets()
	if err != nil {
		panic(err)
	}

	for i, w := range wallets {
		name := strconv.Itoa(i) + w.Address.Hex()
		txq := &TxQueue{
			Name:            name,
			Wallet:          w,
			TradeService:    tradeService,
			EthereumService: ethereumService,
			Exchange:        exchange,
		}

		txqueues = append(txqueues, txq)
	}

	op := &Operator{
		WalletService:     walletService,
		TxService:         txService,
		TradeService:      tradeService,
		EthereumService:   ethereumService,
		Exchange:          exchange,
		TxQueues:          txqueues,
		QueueAddressIndex: addressIndex,
	}

	return op, nil
}

// SubscribeOperatorMessages
func (op *Operator) SubscribeOperatorMessages(fn func(*OperatorMessage) error) error {
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
				var om *OperatorMessage
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

// QueueTrade
func (op *Operator) QueueTrade(o *types.Order, t *types.Trade) error {
	addressindex := op.QueueAddressIndex
	maker := o.UserAddress
	taker := t.Taker

	txq, ok := addressindex[maker]
	if ok {
		if txq.Length() < 10 {
			err := txq.QueueTrade(o, t)
			if err != nil {
				log.Print(err)
				return err
			}
			return nil
		} else {
			return errors.New("User transaction queue full")
		}
	}

	txq, ok = addressindex[taker]
	if ok {
		if txq.Length() < 10 {
			err := txq.QueueTrade(o, t)
			if err != nil {
				log.Print(err)
				return err
			}
			return nil
		} else {
			return errors.New("User transaction queue full")
		}
	}

	txq, len, err := op.GetShortestQueue()
	if err != nil {
		log.Print(err)
		return err
	}

	if len > 10 {
		log.Print("Transaction queue is full")
		return errors.New("Transaction queue is full")
	}

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
	min := op.TxQueues[0].Length()

	for _, txq := range op.TxQueues {
		if shortest == nil {
			shortest = txq
			min = txq.Length()
		}

		ln := txq.Length()
		if ln < min {
			min = ln
		}
	}

	return shortest, min, nil
}

// SetFeeAccount sets the fee account of the exchange contract. The fee account receives
// the trading fees whenever a trade is settled.
func (op *Operator) SetFeeAccount(account common.Address) (*eth.Transaction, error) {
	tx, err := op.Exchange.SetFeeAccount(account)
	if err != nil {
		return nil, err
	}

	return tx, nil
}

// SetOperator updates the operator settings of the given address. Only addresses with an
// operator access can execute Withdraw and Trade transactions to the Exchange smart contract
func (op *Operator) SetOperator(account common.Address, isOperator bool) (*eth.Transaction, error) {
	tx, err := op.Exchange.SetOperator(account, isOperator)
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

// 	// err := t.Validate()
// 	// if err != nil {
// 	// 	return err
// 	// }

// 	ch := getChannel("tradeTxs")
// 	q := getQueue(ch, "tradeTxs")

// 	bytes, err := json.Marshal(t)
// 	if err != nil {
// 		return errors.New("Failed to marshal trade object")
// 	}

// 	err = ch.Publish(
// 		"",
// 		q.Name,
// 		false,
// 		false,
// 		amqp.Publishing{
// 			ContentType: "text/json",
// 			Body:        bytes,
// 		})

// 	length := q.Messages
// 	if length == 1 {
// 		op.ExecuteTrade(o, t)
// 	}

// 	if length == 10 {
// 		return errors.New("Transaction queue is full")
// 	}

// 	return nil
// }

// Trade executes a settlements transaction. The order and trade payloads need to be signed respectively
// // by the Maker and the Taker of the trade. Only the operator account can send a Trade function to the
// // Exchange smart contract.
// func (op *Operator) ExecuteTrade(o *types.Order, tr *types.Trade) (*eth.Transaction, error) {
// 	tx, err := op.Exchange.Trade(o, tr)
// 	if err != nil {
// 		return nil, err
// 	}

// 	err = op.TradeService.UpdateTradeTx(tr, tx)
// 	if err != nil {
// 		return nil, errors.New("Could not update trade tx attribute")
// 	}

// 	err = op.PublishTradeExecutedMessage(tr)
// 	if err != nil {
// 		return nil, errors.New("Could not publish trade executed message")
// 	}

// 	return tx, nil
// }

// // Validate checks that the operator configuration is sufficient.
// func (op *Operator) Validate() error {
// 	// wallet, err := op.WalletService.GetDefaultAdminWallet()
// 	// if err != nil {
// 	// 	return err
// 	// }

// 	// balance, err := op.EthereumService.GetPendingBalanceAt(wallet.Address)
// 	// if err != nil {
// 	// 	return err
// 	// }

// 	return nil
// }

// SetDefaultTxOptions resets the transaction value to 0
// func (op *Operator) SetDefaultTxOptions() {
// 	op.Exchange.TxOptions.Value = big.NewInt(0)
// }

// // SetTxValue sets the transaction ether value
// func (op *Operator) SetTxValue(value *big.Int) {
// 	op.Exchange.TxOptions.Value = value
// }

// // SetCustomSender updates the sender address address to the exchange contract
// func (op *Operator) SetCustomSender(w *types.Wallet) {
// 	op.Exchange.TxOptions = bind.NewKeyedTransactor(w.PrivateKey)
// }

// func getQueue(ch *amqp.Channel, queue string) *amqp.Queue {
// 	if queues[queue] == nil {
// 		q, err := ch.QueueDeclare(queue, false, false, false, false, nil)
// 		if err != nil {
// 			log.Fatalf("Failed to declare a queue: %s", err)
// 		}
// 		queues[queue] = &q
// 	}
// 	return queues[queue]
// }

// func getChannel(id string) *amqp.Channel {
// 	if channels[id] == nil {
// 		ch, err := rabbitmq.Conn.Channel()
// 		if err != nil {
// 			log.Fatalf("Failed to open a channel: %s", err)
// 			panic(err)
// 		}
// 		channels[id] = ch
// 	}
// 	return channels[id]
// }

// go func() {
// 	for {
// 		select {
// 		case event := <-errorEvents:
// 			tradeHash := event.TradeHash
// 			errID := int(event.ErrorId)
// 			//TODO add this function in the trade service
// 			tr, err := op.TradeService.GetByHash(tradeHash)
// 			if err != nil {
// 				log.Printf("Could not retrieve hash")
// 				return
// 			}

// 			err = op.PublishTxErrorMessage(tr, errID)
// 			if err != nil {
// 				log.Printf("Could not publish tx error message")
// 			}

// 			err = op.PublishTradeCancelMessage(tr)
// 			if err != nil {
// 				log.Printf("Could not publish cancel trade message")
// 			}

// 		case event := <-tradeEvents:
// 			//TODO add this function in the trade service
// 			tr, err := tradeService.GetByHash(event.TradeHash)
// 			if err != nil {
// 				log.Printf("Could not retrieve initial hash")
// 				return
// 			}

// 			// only execute the next transaction in the queue when this transaction is mined
// 			go func() {
// 				_, err := op.EthereumService.WaitMined(tr.Tx)
// 				if err != nil {
// 					log.Printf("Could not execute trade: %v\n", err)
// 				}

// 				err = op.PublishTradeSuccessMessage(tr)
// 				if err != nil {
// 					log.Printf("Could not publish order success message")
// 				}

// 				ch := getChannel("PENDING_TRADES")
// 				q := getQueue(ch, "PENDING_TRADES")

// 				length := q.Messages
// 				if length > 0 {
// 					msg, _, _ := ch.Get(
// 						q.Name,
// 						true,
// 					)

// 					var pendingTrade PendingTradeMessage
// 					err = json.Unmarshal(msg.Body, &pendingTrade)
// 					if err != nil {
// 						log.Printf("Could not executed trade: %v\n", err)
// 					}

// 					_, err = op.ExecuteTrade(pendingTrade.Order, pendingTrade.Trade)
// 					if err != nil {
// 						log.Printf("Could not execute trade: %v", err)
// 					}
// 				}
// 			}()
// 		}
// 	}
// }()
