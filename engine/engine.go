package engine

import (
	"encoding/json"
	"errors"
	"sync"

	"github.com/Proofsuite/amp-matching-engine/rabbitmq"
	"github.com/Proofsuite/amp-matching-engine/redis"
	"github.com/Proofsuite/amp-matching-engine/types"
	"github.com/Proofsuite/amp-matching-engine/utils"
)

// Engine contains daos and redis connection required for engine to work
type Engine struct {
	redisConn    *redis.RedisConnection
	rabbitMQConn *rabbitmq.Connection
	mutex        *sync.Mutex
}

var logger = utils.EngineLogger

// InitEngine initializes the engine singleton instance
func InitEngine(redisConn *redis.RedisConnection, rabbitMQConn *rabbitmq.Connection) (*Engine, error) {
	engine := &Engine{redisConn, rabbitMQConn, &sync.Mutex{}}
	return engine, nil
}

// publishEngineResponse is used by matching engine to publish or send response of matching engine to
// system for further processing
func (e *Engine) publishEngineResponse(res *types.EngineResponse) error {
	ch := e.rabbitMQConn.GetChannel("erPub")
	q := e.rabbitMQConn.GetQueue(ch, "engineResponse")

	bytes, err := json.Marshal(res)
	if err != nil {
		logger.Error("Failed to marshal engine response: ", err)
		return errors.New("Failed to marshal Engine Response: " + err.Error())
	}

	err = e.rabbitMQConn.Publish(ch, q, bytes)
	if err != nil {
		logger.Error("Failed to publish order: ", err)
		return errors.New("Failed to publish order: " + err.Error())
	}

	return nil
}

// SubscribeResponseQueue subscribes to engineResponse queue and triggers the function
// passed as arguments for each message.
func (e *Engine) SubscribeResponseQueue(fn func(*types.EngineResponse) error) error {
	ch := e.rabbitMQConn.GetChannel("erSub")
	q := e.rabbitMQConn.GetQueue(ch, "engineResponse")

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
			logger.Fatal("Failed to register a consumer:", err)
		}

		forever := make(chan bool)

		go func() {
			for d := range msgs {
				var res *types.EngineResponse
				err := json.Unmarshal(d.Body, &res)
				if err != nil {
					logger.Error(err)
					continue
				}
				go fn(res)
			}
		}()

		<-forever
	}()
	return nil
}

// HandleOrders parses incoming rabbitmq order messages and redirects them to the appropriate
// engine function
func (e *Engine) HandleOrders(msg *rabbitmq.Message) error {
	o := &types.Order{}
	err := json.Unmarshal(msg.Data, o)
	if err != nil {
		logger.Error(err)
		return err
	}

	if msg.Type == "NEW_ORDER" {
		e.newOrder(o)
	} else if msg.Type == "ADD_ORDER" {
		e.addOrder(o)
	}

	return nil
}
