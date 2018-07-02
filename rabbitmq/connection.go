package rabbitmq

import (
	"log"

	"github.com/streadway/amqp"
)

// Conn is singleton rabbitmq connection
var Conn *amqp.Connection

// InitConnection Initializes single rabbitmq connection for whole system
func InitConnection(address string) {
	if Conn == nil {
		conn, err := amqp.Dial(address)
		if err != nil {
			log.Fatalf("failed to open a connection: %s", err)
			panic(err)
		}
		Conn = conn
	}
}
