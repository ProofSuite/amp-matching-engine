package rabbitmq

import (
	"encoding/json"
	"log"

	"github.com/Proofsuite/amp-matching-engine/types"
	"github.com/Proofsuite/amp-matching-engine/utils"
	"github.com/streadway/amqp"
)

// Conn is singleton rabbitmq connection
var conn *Connection
var channels = make(map[string]*amqp.Channel)
var queues = make(map[string]*amqp.Queue)

var logger = utils.RabbitLogger

type Connection struct {
	Conn *amqp.Connection
}
type Message struct {
	Type string `json:"type"`
	Data []byte `json:"data"`
}

// InitConnection Initializes single rabbitmq connection for whole system
func InitConnection(address string) *Connection {
	if conn == nil {
		newConn, err := amqp.Dial(address)
		if err != nil {
			panic(err)
		}
		conn = &Connection{newConn}
	}
	return conn
}

func (c *Connection) NewConnection(address string) *amqp.Connection {
	conn, err := amqp.Dial(address)
	if err != nil {
		panic(err)
	}

	return conn
}

func (c *Connection) GetQueue(ch *amqp.Channel, queue string) *amqp.Queue {
	if queues[queue] == nil {
		q, err := ch.QueueDeclare(queue, false, false, false, false, nil)
		if err != nil {
			log.Fatalf("Failed to declare a queue: %s", err)
		}

		queues[queue] = &q
	}

	return queues[queue]
}

func (c *Connection) DeclareQueue(ch *amqp.Channel, name string) error {
	if queues[name] == nil {
		q, err := ch.QueueDeclare(name, false, false, false, false, nil)
		if err != nil {
			log.Print(err)
			return err
		}

		queues[name] = &q
	}

	return nil
}

func (c *Connection) GetChannel(id string) *amqp.Channel {
	if channels[id] == nil {
		ch, err := c.Conn.Channel()
		if err != nil {
			log.Fatalf("Failed to open a channel: %s", err)
			panic(err)
		}

		channels[id] = ch
	}

	return channels[id]
}

// Publish
func (c *Connection) Publish(ch *amqp.Channel, q *amqp.Queue, bytes []byte) error {
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

func (c *Connection) Consume(ch *amqp.Channel, q *amqp.Queue) (<-chan amqp.Delivery, error) {
	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)

	if err != nil {
		log.Print(err)
		return nil, err
	}

	return msgs, nil
}

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

func (c *Connection) Purge(ch *amqp.Channel, name string) error {
	_, err := ch.QueueInspect(name)
	if err != nil {
		return nil
	}

	_, err = ch.QueuePurge(name, false)
	if err != nil {
		log.Print(err)
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
