package operator

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math/big"

	"github.com/Proofsuite/amp-matching-engine/interfaces"
	"github.com/Proofsuite/amp-matching-engine/rabbitmq"
	"github.com/Proofsuite/amp-matching-engine/types"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	eth "github.com/ethereum/go-ethereum/core/types"
)

type TxQueue struct {
	Name             string
	Wallet           *types.Wallet
	TradeService     interfaces.TradeService
	OrderService     interfaces.OrderService
	EthereumProvider interfaces.EthereumProvider
	Exchange         interfaces.Exchange
	RabbitMQConn     *rabbitmq.Connection
}

// NewTxQueue
func NewTxQueue(
	n string,
	tr interfaces.TradeService,
	p interfaces.EthereumProvider,
	o interfaces.OrderService,
	w *types.Wallet,
	ex interfaces.Exchange,
) (*TxQueue, error) {

	txq := &TxQueue{
		Name:             n,
		TradeService:     tr,
		OrderService:     o,
		EthereumProvider: p,
		Wallet:           w,
		Exchange:         ex,
	}

	return txq, nil
}

func (txq *TxQueue) GetTxSendOptions() *bind.TransactOpts {
	return bind.NewKeyedTransactor(txq.Wallet.PrivateKey)
}

// Length
func (txq *TxQueue) Length() int {
	name := "TX_QUEUES:" + txq.Name
	ch := txq.RabbitMQConn.GetChannel(name)
	q, err := ch.QueueInspect(name)
	if err != nil {
		log.Print(err)
	}

	return q.Messages
}

// AddTradeToExecutionList adds a new trade to the execution list. If the execution list is empty (= contains 1 element
// after adding the transaction hash), the given order/trade pair gets executed. If the tranasction queue is full,
// we return an error. Ultimately we want to account send the transaction to another queue that is handled by another ethereum account
// func (op *Operator) QueueTrade(o *types.Order, t *types.Trade) error {
// TODO: Currently doesn't seem thread safe and fails unless called with a sleep time between each call.
func (txq *TxQueue) QueueTrade(o *types.Order, t *types.Trade) error {

	fmt.Println("Length of the queue is ", txq.Length())
	if txq.Length() == 0 {
		_, err := txq.ExecuteTrade(o, t)
		if err != nil {
			log.Print(err)
			return err
		}
	}

	err := txq.PublishPendingTrade(o, t)
	if err != nil {
		log.Print(err)
		return err
	}

	return nil
}

// ExecuteTrade send a trade execution order to the smart contract interface. After sending the
// trade message, the trade is updated on the database and is published to the operator subscribers
// (order service)
func (txq *TxQueue) ExecuteTrade(o *types.Order, tr *types.Trade) (*eth.Transaction, error) {
	fmt.Println("EXECUTE_TRADE: ", tr.Hash.Hex())
	nonce, err := txq.EthereumProvider.GetPendingNonceAt(txq.Wallet.Address)
	if err != nil {
		log.Print(err)
		return nil, err
	}

	txOpts := txq.GetTxSendOptions()
	txOpts.Nonce = big.NewInt(int64(nonce))

	log.Print("NONCE IS EQUAL TO", txOpts.Nonce)
	log.Print("QUEUE IS ", txq.Name)

	tx, err := txq.Exchange.Trade(o, tr, txOpts)
	if err != nil {
		log.Print(err)
		return nil, err
	}

	err = txq.TradeService.UpdateTradeTxHash(tr, tx.Hash())
	if err != nil {
		log.Print(err)
		return nil, errors.New("Could not update trade tx attribute")
	}

	err = txq.PublishTradeSentMessage(o, tr)
	if err != nil {
		log.Print(err)
		return nil, errors.New("Could not update")
	}

	go func() {
		fmt.Println("MINING TRADE")
		_, err := txq.EthereumProvider.WaitMined(tx.Hash())
		if err != nil {
			log.Print(err)
		}

		fmt.Println("TRADE_MINED IN EXECUTE TRADE: ", tr.Hash.Hex())

		len := txq.Length()
		fmt.Println("LENGTH of the queue is ", len)
		if len > 0 {
			msg, err := txq.PopPendingTrade()
			if err != nil {
				log.Print(err)
				return
			}

			// very hacky
			if msg.Trade.Hash == tr.Hash {
				fmt.Println("HACKY POP PENDING TRADE: ", tr.Hash.Hex())
				msg, err = txq.PopPendingTrade()
				if err != nil {
					log.Print(err)
					return
				}

				if msg == nil {
					return
				}
			}

			fmt.Println("NEXT_TRADE: ", msg.Trade.Hash.Hex())
			go txq.ExecuteTrade(msg.Order, msg.Trade)
		}
	}()

	return tx, nil
}

func (txq *TxQueue) PublishPendingTrade(o *types.Order, t *types.Trade) error {
	name := "TX_QUEUES:" + txq.Name
	ch := txq.RabbitMQConn.GetChannel(name)
	q := txq.RabbitMQConn.GetQueue(ch, name)
	msg := &types.PendingTradeMessage{o, t}

	bytes, err := json.Marshal(msg)
	if err != nil {
		return errors.New("Failed to marshal trade object")
	}

	err = txq.RabbitMQConn.Publish(ch, q, bytes)
	if err != nil {
		log.Print(err)
		return err
	}

	return nil
}

func (txq *TxQueue) PublishTradeSentMessage(or *types.Order, tr *types.Trade) error {
	fmt.Println("PUBLISHING TRADE SENT MESSAGE")

	ch := txq.RabbitMQConn.GetChannel("OPERATOR_PUB")
	q := txq.RabbitMQConn.GetQueue(ch, "TX_MESSAGES")
	msg := &types.OperatorMessage{
		MessageType: "TRADE_SENT_MESSAGE",
		Trade:       tr,
		Order:       or,
	}

	bytes, err := json.Marshal(msg)
	if err != nil {
		log.Print(err)
		return err
	}

	err = txq.RabbitMQConn.Publish(ch, q, bytes)
	if err != nil {
		log.Print(err)
		return err
	}

	fmt.Println("PUBLISHED TRADE SENT MESSAGE")
	return nil
}

func (txq *TxQueue) PurgePendingTrades() error {
	fmt.Println("PURGING PENDING TRADES")
	name := "TX_QUEUES:" + txq.Name
	ch := txq.RabbitMQConn.GetChannel(name)

	err := txq.RabbitMQConn.Purge(ch, name)
	if err != nil {
		log.Print(err)
		return err
	}

	fmt.Println("PURGED PENDING TRADES")
	return nil
}

// PopPendingTrade
func (txq *TxQueue) PopPendingTrade() (*types.PendingTradeMessage, error) {
	fmt.Println("POPPING PENDING TRADE")
	name := "TX_QUEUES:" + txq.Name
	ch := txq.RabbitMQConn.GetChannel(name)
	q := txq.RabbitMQConn.GetQueue(ch, name)

	msg, _, _ := ch.Get(
		q.Name,
		true,
	)

	if len(msg.Body) == 0 {
		return nil, nil
	}

	pding := &types.PendingTradeMessage{}
	err := json.Unmarshal(msg.Body, &pding)
	if err != nil {
		log.Print(err)
		return nil, err
	}

	fmt.Println("POPPED PENDING TRADE", pding.Trade.Hash.Hex())
	return pding, nil
}
