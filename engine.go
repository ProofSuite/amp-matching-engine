package main

import "errors"

type TradingEngine struct {
	orderbooks map[string]*OrderBook
}

func NewTradingEngine() *TradingEngine {
	engine := new(TradingEngine)
	engine.orderbooks = make(map[string]*OrderBook)

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

func (engine *TradingEngine) CreateNewOrderBook(symbol string, actions chan<- *Action) {
	orderbook := new(OrderBook)
	orderbook.bid = 0
	orderbook.ask = MAX_PRICE

	for i := range orderbook.prices {
		orderbook.prices[i] = new(PricePoint)
	}

	orderbook.actions = actions
	orderbook.orderIndex = make(map[uint64]*Order)
	engine.orderbooks[symbol] = orderbook
}

func (engine *TradingEngine) AddOrder(order *Order) (bool, error) {
	if orderbook, ok := engine.orderbooks[order.Symbol]; !ok {
		return false, errors.New("Orderbook does not exist")
	} else {
		orderbook.AddOrder(order)
		return true, nil
	}
}

func (engine *TradingEngine) CancelOrder(id uint64, symbol string) (bool, error) {
	if orderbook, ok := engine.orderbooks[symbol]; !ok {
		return false, errors.New("Orderbook does not exist")
	} else {
		orderbook.CancelOrder(id, symbol)
		return true, nil
	}
}

func (engine *TradingEngine) CloseOrderBookChannel(symbol string) (bool, error) {
	if orderbook, ok := engine.orderbooks[symbol]; !ok {
		return false, errors.New("Orderbook does not exist")
	} else {
		orderbook.Done()
		return true, nil
	}
}
