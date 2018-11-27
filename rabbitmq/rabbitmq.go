package rabbitmq

import (
	"crypto/tls"
	"crypto/x509"
	"log"

	"github.com/Proofsuite/amp-matching-engine/app"
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

func InitConnection(address string) *Connection {
	if conn == nil {
		tlsEnabled := app.Config.EnableTLS
		if tlsEnabled {
			newConn := NewTLSConnection()
			conn = &Connection{newConn}
		} else {
			newConn := NewConnection()
			conn = &Connection{newConn}
		}

	}

	return conn
}

func NewConnection() *amqp.Connection {
	uri := amqp.URI{
		Scheme:   "amqp",
		Host:     app.Config.RabbitMQURL,
		Port:     5672,
		Username: "guest",
		Password: "guest",
	}

	conn, err := amqp.Dial(uri.String())
	if err != nil {
		panic(err)
	}

	return conn
}

func NewTLSConnection() *amqp.Connection {
	cfg := &tls.Config{InsecureSkipVerify: true}
	cfg.RootCAs = x509.NewCertPool()

	logger.Info("Connecting RabbitMQ with TLS")

	uri := amqp.URI{
		Scheme:   "amqps",
		Host:     app.Config.RabbitMQURL,
		Port:     5671,
		Username: app.Config.RabbitMQUsername,
		Password: app.Config.RabbitMQPassword,
	}

	// ca, err := ioutil.ReadFile(app.Config.TLSCACertFile)
	// if err != nil {
	// 	panic(err)
	// }

	// logger.Info(ui)

	address := uri.String()

	// address = "amqp://hey:cool@127.0.0.1:5672"

	// cfg.RootCAs.AppendCertsFromPEM(ca)

	cert, err := tls.LoadX509KeyPair(app.Config.RabbitMQCert, app.Config.RabbitMQKey)
	if err != nil {
		panic(err)
	}

	// address = uri.String()
	logger.Info("Connecting RabbitMQ with TLS")
	logger.Info(address)

	cfg.Certificates = append(cfg.Certificates, cert)

	conn, err := amqp.DialTLS(address, cfg)
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
			logger.Error(err)
			return err
		}

		queues[name] = &q
	}

	return nil
}

func (c *Connection) DeclareThrottledQueue(ch *amqp.Channel, name string) error {
	ch.Qos(1, 0, true)

	if queues[name] == nil {
		q, err := ch.QueueDeclare(name, false, false, false, false, nil)
		if err != nil {
			logger.Error(err)
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
		logger.Error(err)
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
		logger.Error(err)
		return nil, err
	}

	return msgs, nil
}

func (c *Connection) ConsumeAfterAck(ch *amqp.Channel, q *amqp.Queue) (<-chan amqp.Delivery, error) {
	msgs, err := ch.Consume(
		q.Name,
		"",
		false,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		logger.Error(err)
		return nil, err
	}

	return msgs, nil
}

func (c *Connection) Purge(ch *amqp.Channel, name string) error {
	_, err := ch.QueueInspect(name)
	if err != nil {
		return nil
	}

	_, err = ch.QueuePurge(name, false)
	if err != nil {
		logger.Error(err)
		return err
	}

	return nil
}
