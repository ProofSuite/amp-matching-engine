package rabbitmq

import (
	"log"

	"github.com/streadway/amqp"
)

// Conn is singleton rabbitmq connection
var Conn *amqp.Connection
var channels = make(map[string]*amqp.Channel)
var queues = make(map[string]*amqp.Queue)

type Message struct {
	Type string `json:"type"`
	Data []byte `json:"data"`
}

// InitConnection Initializes single rabbitmq connection for whole system
func InitConnection(address string) {
	if Conn == nil {
		conn, err := amqp.Dial(address)
		if err != nil {
			panic(err)
		}

		Conn = conn
	}
}

func GetQueue(ch *amqp.Channel, queue string) *amqp.Queue {
	if queues[queue] == nil {
		q, err := ch.QueueDeclare(queue, false, false, false, false, nil)
		if err != nil {
			log.Fatalf("Failed to declare a queue: %s", err)
		}
		queues[queue] = &q
	}
	return queues[queue]
}

func GetChannel(id string) *amqp.Channel {
	if channels[id] == nil {
		ch, err := Conn.Channel()
		if err != nil {
			log.Fatalf("Failed to open a channel: %s", err)
			panic(err)
		}
		channels[id] = ch
	}
	return channels[id]
}
