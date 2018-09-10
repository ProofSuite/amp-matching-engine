package rabbitmq

import (
	"encoding/json"
	"log"

	"github.com/Proofsuite/amp-matching-engine/types"
	"github.com/streadway/amqp"
)

// Conn is singleton rabbitmq connection
var conn *Connection
var channels = make(map[string]*amqp.Channel)
var queues = make(map[string]*amqp.Queue)

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
