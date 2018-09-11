package engine

import (
	"encoding/json"
	"errors"
	"sync"

	"github.com/Proofsuite/amp-matching-engine/interfaces"
	"github.com/Proofsuite/amp-matching-engine/rabbitmq"
	"github.com/Proofsuite/amp-matching-engine/redis"
	"github.com/Proofsuite/amp-matching-engine/types"
	"github.com/Proofsuite/amp-matching-engine/utils"
)

// Engine contains daos and redis connection required for engine to work
type Engine struct {
	orderbooks   map[string]*OrderBook
	redisConn    *redis.RedisConnection
	rabbitMQConn *rabbitmq.Connection
}

var logger = utils.EngineLogger

// NewEngine initializes the engine singleton instance
func NewEngine(
	redisConn *redis.RedisConnection,
	rabbitMQConn *rabbitmq.Connection,
	pairDao interfaces.PairDao,
) *Engine {

	pairs, err := pairDao.GetAll()
	if err != nil {
		panic(err)
	}

	obs := map[string]*OrderBook{}
	for _, p := range pairs {
		ob := &OrderBook{
			redisConn:    redisConn,
			rabbitMQConn: rabbitMQConn,
			pair:         &p,
			mutex:        &sync.Mutex{},
		}

		obs[p.Code()] = ob
	}

	engine := &Engine{obs, redisConn, rabbitMQConn}
	return engine
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
		err := e.newOrder(o)
		if err != nil {
			logger.Error(err)
			return err
		}
	} else if msg.Type == "ADD_ORDER" {
		err := e.addOrder(o)
		if err != nil {
			logger.Error(err)
			return err
		}
	}

	return nil
}

func (e *Engine) addOrder(o *types.Order) error {
	code, err := o.PairCode()
	if err != nil {
		logger.Error(err)
		return err
	}

	ob := e.orderbooks[code]
	if ob == nil {
		return errors.New("Orderbook error")
	}

	err = ob.addOrder(o)
	if err != nil {
		logger.Error(err)
		return err
	}

	return nil
}

func (e *Engine) newOrder(o *types.Order) error {
	code, err := o.PairCode()
	if err != nil {
		logger.Error(err)
		return err
	}

	ob := e.orderbooks[code]
	if ob == nil {
		return errors.New("Orderbook error")
	}

	err = ob.newOrder(o)
	if err != nil {
		logger.Error(err)
		return err
	}

	return nil
}
