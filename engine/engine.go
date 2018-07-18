package engine

import (
	"encoding/json"
	"errors"
	"log"
	"sync"

	"github.com/gomodule/redigo/redis"

	"github.com/streadway/amqp"

	"github.com/Proofsuite/amp-matching-engine/rabbitmq"
	"github.com/Proofsuite/amp-matching-engine/types"

	"github.com/Proofsuite/amp-matching-engine/daos"
)

// Resource contains daos and redis connection required for engine to work
type Resource struct {
	orderDao  *daos.OrderDao
	redisConn redis.Conn
	mutex     *sync.Mutex
}

// Message is the structure of message that matching engine expects
type Message struct {
	Type string `json:"type"`
	Data []byte `json:"data"`
}

var channels = make(map[string]*amqp.Channel)
var queues = make(map[string]*amqp.Queue)

// Engine is singleton Resource instance
var Engine *Resource

// InitEngine initializes the engine singleton instance
func InitEngine(orderDao *daos.OrderDao, redisConn redis.Conn) (engine *Resource, err error) {
	if Engine == nil {
		if orderDao == nil {
			return nil, errors.New("Need pointer to struct of type daos.OrderDao")
		}
		Engine = &Resource{orderDao, redisConn, &sync.Mutex{}}
		Engine.subscribeMessage()
	}
	engine = Engine
	return
}

// PublishMessage is used to publish order message over the rabbitmq.
func (e *Resource) PublishMessage(order *Message) error {
	ch := getChannel("orderPublish")
	q := getQueue(ch, "order")

	orderAsBytes, err := json.Marshal(order)
	if err != nil {
		log.Fatalf("Failed to marshal order: %s", err)
		return errors.New("Failed to marshal order: " + err.Error())
	}

	err = ch.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "text/json",
			Body:        orderAsBytes,
		})
	if err != nil {
		log.Fatalf("Failed to publish order: %s", err)
		return errors.New("Failed to publish order: " + err.Error())
	}
	return nil
}

// publishEngineResponse is used by matching engine to publish or send response of matching engine to
// system for further processing
func (e *Resource) publishEngineResponse(er *Response) error {
	ch := getChannel("erPub")
	q := getQueue(ch, "engineResponse")

	erAsBytes, err := json.Marshal(er)
	if err != nil {
		log.Fatalf("Failed to marshal Engine Response: %s", err)
		return errors.New("Failed to marshal Engine Response: " + err.Error())
	}

	err = ch.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "text/json",
			Body:        erAsBytes,
		})
	if err != nil {
		log.Fatalf("Failed to publish order: %s", err)
		return errors.New("Failed to publish order: " + err.Error())
	}
	return nil
}

// SubscribeEngineResponse subscribes to engineResponse queue and triggers the function
// passed as arguments for each message.
func (e *Resource) SubscribeEngineResponse(fn func(*Response) error) error {
	ch := getChannel("erSub")
	q := getQueue(ch, "engineResponse")
	go func() {
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
			log.Fatalf("Failed to register a consumer: %s", err)

		}

		forever := make(chan bool)

		go func() {
			for d := range msgs {
				// log.Printf("Received a message: %s", d.Body)
				var er *Response
				err := json.Unmarshal(d.Body, &er)
				if err != nil {
					log.Printf("error: %s", err)
					continue
				}
				go fn(er)
			}
		}()

		<-forever
	}()
	return nil
}

// subscribeMessage is called by matching engine while initializing,
// it subscribes to order message queue and triggers the fn according to message type.
func (e *Resource) subscribeMessage() error {
	ch := getChannel("orderSubscribe")
	q := getQueue(ch, "order")
	go func() {
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
			log.Fatalf("Failed to register a consumer: %s", err)

		}

		forever := make(chan bool)

		go func() {
			for d := range msgs {
				var msg Message
				err := json.Unmarshal(d.Body, &msg)
				if err != nil {
					log.Printf("Message Unmarshal error: %s", err)
					continue
				}
				var order *types.Order
				err = json.Unmarshal(msg.Data, &order)
				if err != nil {
					log.Printf("Order Unmarshal error: %s", err)
					continue
				}
				if msg.Type == "new_order" {
					e.matchOrder(order)
				} else if msg.Type == "remaining_order_add" {
					e.addOrder(order)
				}
			}
		}()

		<-forever
	}()
	return nil
}

func getQueue(ch *amqp.Channel, queue string) *amqp.Queue {
	if queues[queue] == nil {
		q, err := ch.QueueDeclare(queue, false, false, false, false, nil)
		if err != nil {
			log.Fatalf("Failed to declare a queue: %s", err)
		}
		queues[queue] = &q
	}
	return queues[queue]
}

func getChannel(id string) *amqp.Channel {
	if channels[id] == nil {
		ch, err := rabbitmq.Conn.Channel()
		if err != nil {
			log.Fatalf("Failed to open a channel: %s", err)
			panic(err)
		}
		channels[id] = ch
	}
	return channels[id]
}
