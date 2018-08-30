package operator

import (
	"encoding/json"
	"errors"
	"log"

	"github.com/Proofsuite/amp-matching-engine/contracts"
	"github.com/Proofsuite/amp-matching-engine/rabbitmq"
	"github.com/Proofsuite/amp-matching-engine/services"
	"github.com/Proofsuite/amp-matching-engine/types"
	eth "github.com/ethereum/go-ethereum/core/types"
	"github.com/streadway/amqp"
)

type TxQueue struct {
	Name            string
	Wallet          *types.Wallet
	TradeService    services.TradeServiceInterface
	EthereumService services.EthereumServiceInterface
	Exchange        *contracts.Exchange
}

// NewTxQueue
func NewTxQueue(
	n string,
	tr services.TradeServiceInterface,
	e services.EthereumServiceInterface,
	w *types.Wallet,
	ex *contracts.Exchange,
) (*TxQueue, error) {
	txq := &TxQueue{
		Name:            n,
		TradeService:    tr,
		EthereumService: e,
		Wallet:          w,
		Exchange:        ex,
	}

	go txq.HandleEvents()
	return txq, nil
}

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

			err = txq.PublishTxErrorMessage(tr, errID)
			if err != nil {
				log.Print(err)
			}

			err = txq.PublishTradeCancelMessage(tr)
			if err != nil {
				log.Print(err)
			}

		case event := <-tradeEvents:
			tr, err := txq.TradeService.GetByHash(event.TradeHash)
			if err != nil {
				log.Print(err)
			}

			go func() {
				_, err := txq.EthereumService.WaitMined(tr.Tx)
				if err != nil {
					log.Print(err)
				}

				err = txq.PublishTradeSuccessMessage(tr)
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

					pding := &PendingTradeMessage{}
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
	q := rabbitmq.GetQueue(ch, name)
	return q.Messages
}

// AddTradeToExecutionList adds a new trade to the execution list. If the execution list is empty (= contains 1 element
// after adding the transaction hash), the given order/trade pair gets executed. If the tranasction queue is full,
// we return an error. Ultimately we want to account send the transaction to another queue that is handled by another ethereum account
// func (op *Operator) QueueTrade(o *types.Order, t *types.Trade) error {
func (txq *TxQueue) QueueTrade(o *types.Order, t *types.Trade) error {
	name := "TX_QUEUES:" + txq.Name
	ch := rabbitmq.GetChannel(name)
	q := rabbitmq.GetQueue(ch, name)

	bytes, err := json.Marshal(t)
	if err != nil {
		return errors.New("Failed to marshal trade object")
	}

	err = ch.Publish(
		"",
		q.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/json",
			Body:        bytes,
		},
	)

	length := q.Messages
	if length == 1 {
		txq.ExecuteTrade(o, t)
	}

	return nil
}

// ExecuteTrade send a trade execution order to the smart contract interface. After sending the
// trade message, the trade is updated on the database and is published to the operator subscribers
// (order service)
func (txq *TxQueue) ExecuteTrade(o *types.Order, tr *types.Trade) (*eth.Transaction, error) {
	tx, err := txq.Exchange.Trade(o, tr)
	if err != nil {
		return nil, err
	}

	err = txq.TradeService.UpdateTradeTx(tr, tx)
	if err != nil {
		return nil, errors.New("Could not update trade tx attribute")
	}

	err = txq.PublishTradeSentMessage(tr)
	if err != nil {
		return nil, errors.New("Could not update ")
	}

	return tx, nil
}

// PublishTradeExecutedMessage publishes a message when a trade is sent
func (txq *TxQueue) PublishTradeSentMessage(tr *types.Trade) error {
	msg := &OperatorMessage{
		MessageType: "TRADE_SENT_MESSAGE",
		Trade:       tr,
	}

	err := txq.Publish(msg)
	if err != nil {
		return err
	}

	return nil
}

// PublishTxErrorMessage publishes a messages when a trade execution fails
func (txq *TxQueue) PublishTxErrorMessage(tr *types.Trade, errID int) error {
	msg := &OperatorMessage{
		MessageType: "TRADE_ERROR_MESSAGE",
		Trade:       tr,
		ErrID:       errID,
	}

	err := txq.Publish(msg)
	if err != nil {
		return err
	}

	return nil
}

// PublishTradeCancelMessage publishes a message when a trade is canceled
func (txq *TxQueue) PublishTradeCancelMessage(tr *types.Trade) error {
	msg := &OperatorMessage{
		MessageType: "TRADE_CANCEL_MESSAGE",
		Trade:       tr,
	}

	err := txq.Publish(msg)
	if err != nil {
		return err
	}

	return nil
}

// PublishTradeSuccessMessage publishes a message when a trade transaction
// is successful
func (txq *TxQueue) PublishTradeSuccessMessage(tr *types.Trade) error {
	msg := &OperatorMessage{
		MessageType: "TRADE_SUCCESS_MESSAGE",
		Trade:       tr,
	}

	err := txq.Publish(msg)
	if err != nil {
		return err
	}

	return nil
}

// Publish
func (txq *TxQueue) Publish(msg *OperatorMessage) error {
	ch := rabbitmq.GetChannel("OPERATOR_PUB")
	q := rabbitmq.GetQueue(ch, "TX_MESSAGES")

	bytes, err := json.Marshal(msg)
	if err != nil {
		log.Printf("Failed to marshal %s: %s", msg.MessageType, err)
	}

	err = ch.Publish(
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
		log.Printf("Failed to publish message %s: %s", msg.MessageType, err)
		return err
	}

	return nil
}
