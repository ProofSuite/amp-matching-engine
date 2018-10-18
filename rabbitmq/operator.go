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
					logger.Error(err)
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
			logger.Error(err)
		}

		channels["OPERATOR_SUB"] = nil
	}

	return nil
}

func (c *Connection) UnsubscribeOperator() error {
	ch := c.GetChannel("OPERATOR_SUB")
	q := c.GetQueue(ch, "TX_MESSAGES")

	err := ch.Cancel(q.Name, false)
	if err != nil {
		logger.Error(err)
		return err
	}

	return nil
}

func (c *Connection) PurgeOperatorQueue() error {
	ch := c.GetChannel("OPERATOR_SUB")

	_, err := ch.QueuePurge("TX_MESSAGES", false)
	if err != nil {
		logger.Error(err)
		return err
	}

	return nil
}

// PublishTradeCancelMessage publishes a message when a trade is cancelled
func (c *Connection) PublishTradeCancelMessage(matches *types.Matches) error {
	ch := c.GetChannel("OPERATOR_PUB")
	q := c.GetQueue(ch, "TX_MESSAGES")
	msg := &types.OperatorMessage{
		MessageType: "TRADE_CANCEL",
		Matches:     matches,
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

	logger.Info("PUBLISH TRADE CANCEL MESSAGE")
	return nil
}

// PublishTradeSuccessMessage publishes a message when a trade transaction is successful
func (c *Connection) PublishTradeSuccessMessage(matches *types.Matches) error {
	ch := c.GetChannel("OPERATOR_PUB")
	q := c.GetQueue(ch, "TX_MESSAGES")
	msg := &types.OperatorMessage{
		MessageType: "TRADE_SUCCESS",
		Matches:     matches,
	}

	bytes, err := json.Marshal(msg)
	if err != nil {
		logger.Error(err)
	}

	err = c.Publish(ch, q, bytes)
	if err != nil {
		logger.Error(err)
		return err
	}

	logger.Info("PUBLISH TRADE SUCCESS MESSAGE")
	return nil
}

// PublishTxErrorMessage publishes a messages when a trade execution fails
func (c *Connection) PublishTxErrorMessage(matches *types.Matches, errID int) error {
	ch := c.GetChannel("OPERATOR_PUB")
	q := c.GetQueue(ch, "TX_MESSAGES")
	msg := &types.OperatorMessage{
		MessageType: "TRADE_ERROR",
		Matches:     matches,
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

	logger.Info("PUBLISHED TRADE ERROR MESSAGE. Error ID: %v", errID)
	return nil
}

func (c *Connection) PublishTradeInvalidMessage(matches *types.Matches) error {
	ch := c.GetChannel("OPERATOR_PUB")
	q := c.GetQueue(ch, "TX_MESSAGES")
	msg := &types.OperatorMessage{
		MessageType: "TRADE_INVALID",
		Matches:     matches,
	}

	bytes, err := json.Marshal(msg)
	if err != nil {
		logger.Error(err)
	}

	err = c.Publish(ch, q, bytes)
	if err != nil {
		logger.Error(err)
		return err
	}

	logger.Info("PUBLISHED TRADE INVALID MESSAGE")
	return nil
}

func (c *Connection) PublishTradeSentMessage(matches *types.Matches) error {
	ch := c.GetChannel("OPERATOR_PUB")
	q := c.GetQueue(ch, "TX_MESSAGES")
	msg := &types.OperatorMessage{
		MessageType: "TRADE_PENDING",
		Matches:     matches,
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
