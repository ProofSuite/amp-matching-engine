package main

import (
	"fmt"
)

const MAX_PRICE = 10000000

//Pricepoint contains pointers to the first and the last order entered at that price
type PricePoint struct {
	orderHead *Order
	orderTail *Order
}

//The orderbook keeps track of the maximum bid and minimum ask
//The orderIndex is a mapping of OrderIDs to pointers so that we can easily cancel outstanding orders
//prices is an array of all possible pricepoints
//actions is a channels to report some action to a handler as they occur
type OrderBook struct {
	ask             uint32
	bid             uint32
	orderIndex      map[uint64]*Order
	prices          [MAX_PRICE]*PricePoint
	actions         chan *Action
	outboundActions chan *Action
	logger          []*Action
}

func (orderbook *OrderBook) String() string {
	return fmt.Sprintf("Ask:%v, Bid:%v, orderIndex:%v", orderbook.ask, orderbook.bid, orderbook.orderIndex)
}

func (orderbook *OrderBook) GetLogs() []*Action {
	return orderbook.logger
}

func (p *PricePoint) Insert(order *Order) {
	if p.orderHead == nil {
		p.orderHead = order
		p.orderTail = order
	} else {
		p.orderTail.next = order
		p.orderTail = order
	}
}

func (ob *OrderBook) AddOrder(o *Order) {

	if o.OrderType == BUY {
		ob.actions <- NewBuyAction(o)
		ob.FillBuy(o)
	} else {
		ob.actions <- NewSellAction(o)
		ob.FillSell(o)
	}
	if o.Amount > 0 {
		ob.openOrder(o)
	}
}

func (ob *OrderBook) openOrder(o *Order) {
	o.signalOrderCreated()
	pricePoint := ob.prices[o.Price]
	pricePoint.Insert(o)
	o.status = OPEN
	if o.OrderType == BUY && o.Price > ob.bid {
		ob.bid = o.Price
	} else if o.OrderType == SELL && o.Price < ob.ask {
		ob.ask = o.Price
	}

	ob.orderIndex[o.Id] = o
}

func (ob *OrderBook) CancelOrder(id uint64, pair Pair) {
	ob.actions <- NewCancelAction(id, pair)
	if order, ok := ob.orderIndex[id]; ok {
		order.Amount = 0
		order.status = CANCELLED
	}

	ob.actions <- NewCancelledAction(id, pair)
}

func (ob *OrderBook) FillBuy(order *Order) {

	for ob.ask < order.Price && order.Amount > 0 {
		pricePoint := ob.prices[ob.ask]
		pricePointOrderHead := pricePoint.orderHead
		for pricePointOrderHead != nil {
			ob.fill(order, pricePointOrderHead)
			pricePointOrderHead = pricePointOrderHead.next
			pricePoint.orderHead = pricePointOrderHead
		}
		ob.ask++
	}
}

func (ob *OrderBook) FillSell(order *Order) {
	for ob.bid >= order.Price && order.Amount > 0 {
		pricePoint := ob.prices[ob.bid]
		pricePointOrderHead := pricePoint.orderHead
		for pricePointOrderHead != nil {
			ob.fill(order, pricePointOrderHead)
			pricePointOrderHead = pricePointOrderHead.next
			pricePoint.orderHead = pricePointOrderHead
		}
		ob.bid--
	}
}

func (ob *OrderBook) fill(o, pricePointOrderHead *Order) {
	if pricePointOrderHead.Amount >= o.Amount {
		ob.fillCompletely(o, pricePointOrderHead)
		return
	} else {
		if pricePointOrderHead.Amount > 0 {
			ob.fillPartially(o, pricePointOrderHead)
			return
		}
	}
}

func (ob *OrderBook) fillCompletely(o, pricePointOrderHead *Order) {
	ob.actions <- NewFilledAction(o, pricePointOrderHead)
	pricePointOrderHead.Amount -= o.Amount
	o.Amount = 0
	o.status = FILLED
	o.signalOrderFilled()
	return
}

func (ob *OrderBook) fillPartially(o, pricePointOrderHead *Order) {
	ob.actions <- NewPartialFilledAction(o, pricePointOrderHead)
	o.Amount -= pricePointOrderHead.Amount
	o.status = PARTIAL_FILLED
	o.signalOrderPartiallyFilled()
	return
}

func (ob *OrderBook) Done() {
	ob.actions <- NewDoneAction()
}
