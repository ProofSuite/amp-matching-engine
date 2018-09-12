package rabbitmq

import (
	"encoding/json"
	"log"

	"github.com/Proofsuite/amp-matching-engine/types"
)

func (c *Connection) SubscribeOperator(fn func(*types.OperatorMessage) error) error {
	ch := c.GetChannel("OPERATOR_SUB")
	q := c.GetQueue(ch, "TX_MESSAGES")

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
			log.Fatal("Failed to register a consumer", err)
		}

		forever := make(chan bool)

		go func() {
			for m := range msgs {
				om := &types.OperatorMessage{}
				err := json.Unmarshal(m.Body, &om)
				if err != nil {
					log.Print(err)
					continue
				}

				go fn(om)
			}
		}()

		<-forever
	}()

	return nil
}

func (c *Connection) CloseOperatorChannel() error {
	if channels["OPERATOR_SUB"] != nil {
		ch := c.GetChannel("OPERATOR_SUB")
		err := ch.Close()
		if err != nil {
			log.Print(err)
		}

		channels["OPERATOR_SUB"] = nil
	}

	return nil
}

func (c *Connection) UnSubscribeOperator() error {
	ch := c.GetChannel("OPERATOR_SUB")
	q := c.GetQueue(ch, "TX_MESSAGES")

	err := ch.Cancel(q.Name, false)
	if err != nil {
		log.Print(err)
		return err
	}

	return nil
}

func (c *Connection) PurgeOperatorQueue() error {
	ch := c.GetChannel("OPERATOR_SUB")

	_, err := ch.QueuePurge("TX_MESSAGES", false)
	if err != nil {
		log.Print(err)
		return err
	}

	return nil
}

// PublishTradeCancelMessage publishes a message when a trade is canceled
func (c *Connection) PublishTradeCancelMessage(o *types.Order, tr *types.Trade) error {
	ch := c.GetChannel("OPERATOR_PUB")
	q := c.GetQueue(ch, "TX_MESSAGES")
	msg := &types.OperatorMessage{
		MessageType: "TRADE_CANCEL_MESSAGE",
		Trade:       tr,
	}

	bytes, err := json.Marshal(msg)
	if err != nil {
		logger.Infof("Failed to marshal %s: %s", msg.MessageType, err)
	}

	err = c.Publish(ch, q, bytes)
	if err != nil {
		logger.Error(err)
		return err
	}

	return nil
}

// PublishTradeSuccessMessage publishes a message when a trade transaction is successful
func (c *Connection) PublishTradeSuccessMessage(o *types.Order, tr *types.Trade) error {
	ch := c.GetChannel("OPERATOR_PUB")
	q := c.GetQueue(ch, "TX_MESSAGES")
	msg := &types.OperatorMessage{
		MessageType: "TRADE_SUCCESS_MESSAGE",
		Order:       o,
		Trade:       tr,
	}

	bytes, err := json.Marshal(msg)
	if err != nil {
		logger.Infof("Failed to marshal %s: %s", msg.MessageType, err)
	}

	err = c.Publish(ch, q, bytes)
	if err != nil {
		logger.Error(err)
		return err
	}

	return nil
}

// PublishTxErrorMessage publishes a messages when a trade execution fails
func (c *Connection) PublishTxErrorMessage(tr *types.Trade, errID int) error {
	ch := c.GetChannel("OPERATOR_PUB")
	q := c.GetQueue(ch, "TX_MESSAGES")
	msg := &types.OperatorMessage{
		MessageType: "TRADE_ERROR_MESSAGE",
		Trade:       tr,
		ErrID:       errID,
	}

	bytes, err := json.Marshal(msg)
	if err != nil {
		logger.Infof("Failed to marshal %s: %s", msg.MessageType, err)
	}

	err = c.Publish(ch, q, bytes)
	if err != nil {
		logger.Error(err)
		return err
	}

	return nil
}

func (c *Connection) PublishTradeSentMessage(or *types.Order, tr *types.Trade) error {
	ch := c.GetChannel("OPERATOR_PUB")
	q := c.GetQueue(ch, "TX_MESSAGES")
	msg := &types.OperatorMessage{
		MessageType: "TRADE_SENT_MESSAGE",
		Trade:       tr,
		Order:       or,
	}

	bytes, err := json.Marshal(msg)
	if err != nil {
		logger.Error(err)
		return err
	}

	err = c.Publish(ch, q, bytes)
	if err != nil {
		logger.Error(err)
		return err
	}

	logger.Info("PUBLISHED TRADE SENT MESSAGE")
	return nil
}
