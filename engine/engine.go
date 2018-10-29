package engine

import (
	"encoding/json"
	"errors"
	"sync"

	"github.com/Proofsuite/amp-matching-engine/interfaces"
	"github.com/Proofsuite/amp-matching-engine/rabbitmq"
	"github.com/Proofsuite/amp-matching-engine/types"
	"github.com/Proofsuite/amp-matching-engine/utils"
)

// Engine contains daos and redis connection required for engine to work
type Engine struct {
	orderbooks   map[string]*OrderBook
	rabbitMQConn *rabbitmq.Connection
}

var logger = utils.EngineLogger

// NewEngine initializes the engine singleton instance
func NewEngine(
	rabbitMQConn *rabbitmq.Connection,
	orderDao interfaces.OrderDao,
	tradeDao interfaces.TradeDao,
	pairDao interfaces.PairDao,
) *Engine {
	pairs, err := pairDao.GetAll()

	if err != nil {
		panic(err)
	}

	obs := map[string]*OrderBook{}
	for _, p := range pairs {
		ob := &OrderBook{
			rabbitMQConn: rabbitMQConn,
			orderDao:     orderDao,
			tradeDao:     tradeDao,
			pair:         &p,
			mutex:        &sync.Mutex{},
		}

		obs[p.Code()] = ob
	}

	engine := &Engine{obs, rabbitMQConn}
	return engine
}

// HandleOrders parses incoming rabbitmq order messages and redirects them to the appropriate
// engine function
func (e *Engine) HandleOrders(msg *rabbitmq.Message) error {
	switch msg.Type {
	case "NEW_ORDER":
		err := e.handleNewOrder(msg.Data)
		if err != nil {
			logger.Error(err)
			return err
		}
	case "ADD_ORDER":
		err := e.handleAddOrder(msg.Data)
		if err != nil {
			logger.Error(err)
			return err
		}
	case "CANCEL_ORDER":
		err := e.handleCancelOrder(msg.Data)
		if err != nil {
			logger.Error(err)
			return err
		}
	case "INVALIDATE_MAKER_ORDERS":
		utils.PrintJSON("receiving invalidate maker orders")
		err := e.handleInvalidateMakerOrders(msg.Data)
		if err != nil {
			logger.Error(err)
			return err
		}
	case "INVALIDATE_TAKER_ORDERS":
		err := e.handleInvalidateTakerOrders(msg.Data)
		if err != nil {
			logger.Error(err)
			return err
		}
	default:
		logger.Error("Unknown message", msg)
	}

	return nil
}

func (e *Engine) handleAddOrder(bytes []byte) error {
	o := &types.Order{}
	err := json.Unmarshal(bytes, o)
	if err != nil {
		logger.Error(err)
		return err
	}

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

func (e *Engine) handleNewOrder(bytes []byte) error {
	o := &types.Order{}
	err := json.Unmarshal(bytes, o)
	if err != nil {
		logger.Error(err)
		return err
	}

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

func (e *Engine) handleCancelOrder(bytes []byte) error {
	o := &types.Order{}
	err := json.Unmarshal(bytes, o)
	if err != nil {
		logger.Error(err)
		return err
	}

	code, err := o.PairCode()
	if err != nil {
		logger.Error(err)
		return err
	}

	ob := e.orderbooks[code]
	if ob == nil {
		return errors.New("Orderbook error")
	}

	err = ob.cancelOrder(o)
	if err != nil {
		logger.Error(err)
		return err
	}

	return nil
}

func (e *Engine) handleInvalidateMakerOrders(bytes []byte) error {
	m := types.Matches{}
	err := json.Unmarshal(bytes, &m)
	if err != nil {
		logger.Error(err)
		return err
	}

	code, err := m.PairCode()
	if err != nil {
		logger.Error(err)
		return err
	}

	ob := e.orderbooks[code]
	if ob == nil {
		return errors.New("Orderbook error")
	}

	err = ob.invalidateMakerOrders(m)
	if err != nil {
		logger.Error(err)
		return err
	}

	return nil
}

func (e *Engine) handleInvalidateTakerOrders(bytes []byte) error {
	m := types.Matches{}
	err := json.Unmarshal(bytes, m)
	if err != nil {
		logger.Error(err)
		return err
	}

	code, err := m.PairCode()
	if err != nil {
		logger.Error(err)
		return err
	}

	ob := e.orderbooks[code]
	if ob == nil {
		logger.Error(err)
		return err
	}

	err = ob.invalidateTakerOrders(m)
	if err != nil {
		logger.Error(err)
		return err
	}

	return nil
}

// func (e *Engine) deleteOrders(bytes []byte) error {
// 	//we assume all the orders correspond to the same pair
// 	orders := []*types.Order{}
// 	code, err := orders[0].PairCode()
// 	if err != nil {
// 		logger.Error(err)
// 		return err
// 	}

// 	ob := e.orderbooks[code]
// 	if ob == nil {
// 		return errors.New("Orderbook error")
// 	}

// 	err = ob.deleteOrders(orders...)
// 	if err != nil {
// 		logger.Error(err)
// 		return err
// 	}

// 	return nil
// }
