package operator

import (
	"encoding/json"
	"errors"
	"log"

	"github.com/Proofsuite/amp-matching-engine/interfaces"
	"github.com/Proofsuite/amp-matching-engine/rabbitmq"
	"github.com/Proofsuite/amp-matching-engine/types"
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

	go txq.HandleEvents()
	return txq, nil
}

// GetTxQueue returns the corresponding ampq queue
func GetTxQueue(name string) *amqp.Queue {
	ch := rabbitmq.GetChannel(name)
	q := rabbitmq.GetQueue(ch, name)

	return q
}

// Bug: In certain cases, the trade channel seems to be receiving additional unexpected trades.
// In the case TestSocketExecuteOrder (in file socket_test.go) is run on its own, everything is working correctly.
// However, in the case TestSocketExecuteOrder is run among other tests, some tradeLogs do not correspond to an
// order hash in the ordertrade mapping. I suspect this is because the event listener catches events from previous
// tests. It might be helpful to see how to listen to events from up to a certain block.
func (txq *TxQueue) HandleEvents() error {
	tradeEvents, err := txq.Exchange.ListenToTrades()
	if err != nil {
		log.Print(err)
		return err
	}

	errorEvents, err := txq.Exchange.ListenToErrors()
	if err != nil {
		log.Print(err)
		return err
	}

	for {
		select {
		case event := <-errorEvents:
			tradeHash := event.TradeHash
			errID := int(event.ErrorId)

			tr, err := txq.TradeService.GetByHash(tradeHash)
			if err != nil {
				log.Print(err)
			}

			or, err := txq.OrderService.GetByHash(tr.OrderHash)
			if err != nil {
				log.Print(err)
			}

			err = txq.PublishTxErrorMessage(tr, errID)
			if err != nil {
				log.Print(err)
			}

			err = txq.PublishTradeCancelMessage(or, tr)
			if err != nil {
				log.Print(err)
			}

		case event := <-tradeEvents:
			tr, err := txq.TradeService.GetByHash(event.TradeHash)
			if err != nil {
				log.Print(err)
			}

			or, err := txq.OrderService.GetByHash(tr.OrderHash)
			if err != nil {
				log.Print(err)
			}

			go func() {
				_, err := txq.EthereumService.WaitMined(tr.Tx)
				if err != nil {
					log.Print(err)
				}

				err = txq.PublishTradeSuccessMessage(or, tr)
				if err != nil {
					log.Print(err)
				}

				qname := "TX_QUEUES:" + txq.Name
				ch := rabbitmq.GetChannel(qname)
				q := rabbitmq.GetQueue(ch, qname)

				len := q.Messages
				if len > 0 {
					msg, _, _ := ch.Get(
						q.Name,
						true,
					)

					pding := &types.PendingTradeMessage{}
					err = json.Unmarshal(msg.Body, &pding)
					if err != nil {
						log.Print(err)
					}

					_, err = txq.ExecuteTrade(pding.Order, pding.Trade)
					if err != nil {
						log.Print(err)
					}
				}
			}()
		}
	}
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

// AddTradeToExecutionList adds a new trade to the execution list. If the execution list is empty (= contains 1 element
// after adding the transaction hash), the given order/trade pair gets executed. If the tranasction queue is full,
// we return an error. Ultimately we want to account send the transaction to another queue that is handled by another ethereum account
// func (op *Operator) QueueTrade(o *types.Order, t *types.Trade) error {
func (txq *TxQueue) QueueTrade(o *types.Order, t *types.Trade) error {
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
	tx, err := txq.Exchange.Trade(o, tr)
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

func (txq *TxQueue) PurgePendingTrades() error {
	name := "TX_QUEUES:" + txq.Name
	ch := rabbitmq.GetChannel(name)

	err := rabbitmq.Purge(ch, name)
	if err != nil {
		log.Print(err)
		return err
	}

	return nil
}

// PopPendingTrade
func (txq *TxQueue) PopPendingTrade() (*types.PendingTradeMessage, error) {
	name := "TX_QUEUES:" + txq.Name
	ch := rabbitmq.GetChannel(name)
	q := rabbitmq.GetQueue(ch, name)

	msg, _, _ := ch.Get(
		q.Name,
		true,
	)

	pding := &types.PendingTradeMessage{}
	err := json.Unmarshal(msg.Body, &pding)
	if err != nil {
		log.Print(err)
		return nil, err
	}

	return pding, nil
}

// PublishTradeExecutedMessage publishes a message when a trade is sent
func (txq *TxQueue) PublishTradeSentMessage(or *types.Order, tr *types.Trade) error {
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

	return nil
}

// PublishTxErrorMessage publishes a messages when a trade execution fails
func (txq *TxQueue) PublishTxErrorMessage(tr *types.Trade, errID int) error {
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
func (txq *TxQueue) PublishTradeCancelMessage(o *types.Order, tr *types.Trade) error {
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
func (txq *TxQueue) PublishTradeSuccessMessage(o *types.Order, tr *types.Trade) error {
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
func (txq *TxQueue) Publish(ch *amqp.Channel, q *amqp.Queue, bytes []byte) error {
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
