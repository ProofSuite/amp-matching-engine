package engine

import (
	"encoding/json"
	"errors"
	"log"

	"github.com/gomodule/redigo/redis"

	"github.com/streadway/amqp"

	"github.com/Proofsuite/amp-matching-engine/rabbitmq"
	"github.com/Proofsuite/amp-matching-engine/types"

	"github.com/Proofsuite/amp-matching-engine/daos"
)

type EngineResource struct {
	orderDao  *daos.OrderDao
	redisConn redis.Conn
}

var channels = make(map[string]*amqp.Channel)
var queues = make(map[string]*amqp.Queue)
var Engine *EngineResource

func InitEngine(orderDao *daos.OrderDao, redisConn redis.Conn) (engine *EngineResource, err error) {
	if Engine == nil {
		if orderDao == nil {
			return nil, errors.New("Need pointer to struct of type daos.OrderDao")
		}
		Engine = &EngineResource{orderDao, redisConn}
		Engine.subscribeOrder()
	}
	engine = Engine
	return
}

func (e *EngineResource) PublishOrder(order *types.Order) error {
	ch := getChannel("orderPublish")
	q := getOrderQueue(ch, "order")

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
func (e *EngineResource) publishEngineResponse(er *EngineResponse) error {
	ch := getChannel("erPub")
	q := getOrderQueue(ch, "engineResponse")

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
func (e *EngineResource) SubscribeEngineResponse(fn func(*EngineResponse) error) error {
	ch := getChannel("erSub")
	q := getOrderQueue(ch, "engineResponse")
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
				var er *EngineResponse
				err := json.Unmarshal(d.Body, &er)
				if err != nil {
					log.Printf("error: %s", err)
				}
				fn(er)
			}
		}()

		log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
		<-forever
	}()
	return nil
}
func (e *EngineResource) subscribeOrder() error {
	ch := getChannel("orderSubscribe")
	q := getOrderQueue(ch, "order")
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
				var order *types.Order
				err := json.Unmarshal(d.Body, &order)
				if err != nil {
					log.Printf("error: %s", err)
				}
				e.matchOrder(order)
			}
		}()

		log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
		<-forever
	}()
	return nil
}

func getOrderQueue(ch *amqp.Channel, queue string) *amqp.Queue {
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
