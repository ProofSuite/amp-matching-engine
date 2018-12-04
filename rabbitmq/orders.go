package rabbitmq

import (
	"encoding/json"
	"errors"
	"log"

	"github.com/Proofsuite/amp-matching-engine/types"
)

func (c *Connection) SubscribeOrders(fn func(*Message) error) error {
	ch := c.GetChannel("orderSubscribe")
	q := c.GetQueue(ch, "order")

	go func() {
		msgs, err := c.Consume(ch, q)
		if err != nil {
			logger.Error(err)
		}

		forever := make(chan bool)

		go func() {
			for d := range msgs {
				msg := &Message{}
				err := json.Unmarshal(d.Body, msg)
				if err != nil {
					logger.Error(err)
					continue
				}

				go fn(msg)
			}
		}()

		<-forever
	}()
	return nil
}

func (c *Connection) SubscribeTrades(fn func(*types.OperatorMessage) error) error {
	ch := c.GetChannel("tradeSubscribe")
	q := c.GetQueue(ch, "trades")

	go func() {
		msgs, err := c.Consume(ch, q)
		if err != nil {
			logger.Error(err)
		}

		forever := make(chan bool)

		go func() {
			for d := range msgs {
				msg := &types.OperatorMessage{}
				err := json.Unmarshal(d.Body, msg)
				if err != nil {
					logger.Error(err)
					continue
				}

				go fn(msg)
			}
		}()

		<-forever
	}()
	return nil
}

func (c *Connection) PublishTrades(matches *types.Matches) error {
	ch := c.GetChannel("tradePublish")
	q := c.GetQueue(ch, "trades")

	msg := &types.OperatorMessage{
		MessageType: "NEW_ORDER",
		Matches:     matches,
	}

	logger.Info("operator/", msg.String())

	b, err := json.Marshal(msg)
	if err != nil {
		logger.Error(err)
		return err
	}

	err = c.Publish(ch, q, b)
	if err != nil {
		logger.Error(err)
		return err
	}

	return nil
}

func (c *Connection) PublishNewOrderMessage(o *types.Order) error {
	b, err := json.Marshal(o)
	if err != nil {
		logger.Error(err)
		return err
	}

	err = c.PublishOrder(&Message{
		Type: "NEW_ORDER",
		Data: b,
	})

	if err != nil {
		logger.Error(err)
		return err
	}

	return nil
}

func (c *Connection) PublishCancelOrderMessage(o *types.Order) error {
	b, err := json.Marshal(o)
	if err != nil {
		logger.Error(err)
		return err
	}

	err = c.PublishOrder(&Message{
		Type: "CANCEL_ORDER",
		Data: b,
	})

	if err != nil {
		logger.Error(err)
		return err
	}

	return nil
}

func (c *Connection) PublishInvalidateMakerOrdersMessage(m types.Matches) error {
	b, err := json.Marshal(m)
	if err != nil {
		logger.Error(err)
		return err
	}

	err = c.PublishOrder(&Message{
		Type: "INVALIDATE_MAKER_ORDERS",
		Data: b,
	})

	if err != nil {
		logger.Error(err)
		return err
	}

	return nil
}

func (c *Connection) PublishInvalidateTakerOrdersMessage(m types.Matches) error {
	b, err := json.Marshal(m)
	if err != nil {
		logger.Error(err)
		return err
	}

	err = c.PublishOrder(&Message{
		Type: "INVALIDATE_TAKER_ORDERS",
		Data: b,
	})

	if err != nil {
		logger.Error(err)
		return err
	}

	return nil
}

func (c *Connection) PublishOrder(order *Message) error {
	ch := c.GetChannel("orderPublish")
	q := c.GetQueue(ch, "order")

	bytes, err := json.Marshal(order)
	if err != nil {
		log.Fatal("Failed to marshal order: ", err)
		return errors.New("Failed to marshal order: " + err.Error())
	}

	err = c.Publish(ch, q, bytes)
	if err != nil {
		logger.Error(err)
		return err
	}

	return nil
}
