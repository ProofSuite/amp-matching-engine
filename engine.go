package main

import (
	"errors"
	"fmt"
)

type TradingEngine struct {
	orderbooks map[Pair]*OrderBook
}

func NewTradingEngine() *TradingEngine {
	engine := new(TradingEngine)
	engine.orderbooks = make(map[Pair]*OrderBook)
	return engine
}

func NewOrderBook(actions chan<- *Action) *OrderBook {
	ob := new(OrderBook)
	ob.bid = 0
	ob.ask = MAX_PRICE

	for i := range ob.prices {
		ob.prices[i] = new(PricePoint)
	}

	ob.actions = actions
	ob.orderIndex = make(map[uint64]*Order)
	return ob
}

func (engine *TradingEngine) PrintLogs() {
	for symbol, orderbook := range engine.orderbooks {
		fmt.Printf("Logs for %v:\n", symbol)
		fmt.Printf("%v\n", orderbook.GetLogs())
	}
}

func (engine *TradingEngine) CreateNewOrderBook(pair Pair, done chan<- bool) {
	actions := make(chan *Action)
	logger := make([]*Action, 0)

	orderbook := new(OrderBook)
	orderbook.bid = 0
	orderbook.ask = MAX_PRICE
	orderbook.actions = actions
	orderbook.logger = logger
	orderbook.orderIndex = make(map[uint64]*Order)

	go func() {
		for {
			action := <-actions
			orderbook.logger = append(orderbook.logger, action)
			if action.actionType == AT_DONE {
				done <- true
				return
			}
		}
	}()

	for i := range orderbook.prices {
		orderbook.prices[i] = new(PricePoint)
	}
	engine.orderbooks[pair] = orderbook
}

func (engine *TradingEngine) AddOrder(order *Order) error {
	if orderbook, ok := engine.orderbooks[order.Pair]; !ok {
		return errors.New("Orderbook does not exist")
	} else {
		orderbook.AddOrder(order)
		return nil
	}
}

func (engine *TradingEngine) CancelOrder(id uint64, pair Pair) error {
	if orderbook, ok := engine.orderbooks[pair]; !ok {
		return errors.New("Orderbook does not exist")
	} else {
		orderbook.CancelOrder(id, pair)
		return nil
	}
}

func (engine *TradingEngine) CloseOrderBookChannel(pair Pair) (bool, error) {
	if orderbook, ok := engine.orderbooks[pair]; !ok {
		return false, errors.New("Orderbook does not exist")
	} else {
		orderbook.Done()
		return true, nil
	}
}
