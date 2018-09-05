package operator

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math/big"
	"sync"

	"github.com/Proofsuite/amp-matching-engine/interfaces"
	"github.com/Proofsuite/amp-matching-engine/rabbitmq"
	"github.com/Proofsuite/amp-matching-engine/types"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	eth "github.com/ethereum/go-ethereum/core/types"
	"github.com/streadway/amqp"
)

type TxQueue struct {
	Name            string
	Wallet          *types.Wallet
	TradeService    interfaces.TradeService
	OrderService    interfaces.OrderService
	EthereumService interfaces.EthereumService
	Exchange        interfaces.Exchange
}

// NewTxQueue
func NewTxQueue(
	n string,
	tr interfaces.TradeService,
	e interfaces.EthereumService,
	o interfaces.OrderService,
	w *types.Wallet,
	ex interfaces.Exchange,
) (*TxQueue, error) {

	txq := &TxQueue{
		Name:            n,
		TradeService:    tr,
		OrderService:    o,
		EthereumService: e,
		Wallet:          w,
		Exchange:        ex,
	}

	return txq, nil
}

// GetTxQueue returns the corresponding ampq queue
func GetTxQueue(name string) *amqp.Queue {
	ch := rabbitmq.GetChannel(name)
	q := rabbitmq.GetQueue(ch, name)

	return q
}

func (txq *TxQueue) GetTxSendOptions() *bind.TransactOpts {
	return bind.NewKeyedTransactor(txq.Wallet.PrivateKey)
}

// Length
func (txq *TxQueue) Length() int {
	name := "TX_QUEUES:" + txq.Name
	ch := rabbitmq.GetChannel(name)
	q, err := ch.QueueInspect(name)
	if err != nil {
		log.Print(err)
	}

	return q.Messages
}

var mutex = &sync.Mutex{}

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
	nonce, err := txq.EthereumService.GetPendingNonceAt(txq.Wallet.Address)
	if err != nil {
		log.Print(err)
		return nil, err
	}

	txOpts := txq.GetTxSendOptions()
	txOpts.Nonce = big.NewInt(int64(nonce))

	tx, err := txq.Exchange.Trade(o, tr, txOpts)
	if err != nil {
		log.Print(err)
		return nil, err
	}

	err = txq.TradeService.UpdateTradeTx(tr, tx)
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
		_, err := txq.EthereumService.WaitMined(tx)
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

				// return

				if msg == nil {
					return
				}

				// asdfasdf
			}

			fmt.Println("NEXT_TRADE: ", msg.Trade.Hash.Hex())
			go txq.ExecuteTrade(msg.Order, msg.Trade)
		}
	}()

	return tx, nil
}

func (txq *TxQueue) PublishPendingTrade(o *types.Order, t *types.Trade) error {
	name := "TX_QUEUES:" + txq.Name
	ch := rabbitmq.GetChannel(name)
	q := rabbitmq.GetQueue(ch, name)
	msg := &types.PendingTradeMessage{o, t}

	bytes, err := json.Marshal(msg)
	if err != nil {
		return errors.New("Failed to marshal trade object")
	}

	err = rabbitmq.Publish(ch, q, bytes)
	if err != nil {
		log.Print(err)
		return err
	}

	return nil
}

func (txq *TxQueue) PublishTradeSentMessage(or *types.Order, tr *types.Trade) error {
	fmt.Println("PUBLISHING TRADE SENT MESSAGE")
	ch := rabbitmq.GetChannel("OPERATOR_PUB")
	q := rabbitmq.GetQueue(ch, "TX_MESSAGES")
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

	err = rabbitmq.Publish(ch, q, bytes)
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
	ch := rabbitmq.GetChannel(name)

	err := rabbitmq.Purge(ch, name)
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
	ch := rabbitmq.GetChannel(name)
	q := rabbitmq.GetQueue(ch, name)

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
