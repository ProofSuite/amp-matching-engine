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
	ask        uint32
	bid        uint32
	orderIndex map[uint64]*Order
	prices     [MAX_PRICE]*PricePoint
	actions    chan<- *Action
	logger     []*Action
}

func (orderbook *OrderBook) String() string {
	return fmt.Sprintf("Ask:%v, Bid:%v, orderIndex:%v", orderbook.ask, orderbook.bid, orderbook.orderIndex)
}

func (orderbook *OrderBook) GetLogs() []*Action {
	return orderbook.logger
}

func (pricePoint *PricePoint) Insert(order *Order) {
	if pricePoint.orderHead == nil {
		pricePoint.orderHead = order
		pricePoint.orderTail = order
	} else {
		pricePoint.orderTail.next = order
		pricePoint.orderTail = order
	}
}

func (ob *OrderBook) AddOrder(order *Order) {
	fmt.Printf("This is cool")
	if order.OrderType == BUY {
		fmt.Printf("Buying")
		ob.actions <- NewBuyAction(order)
		ob.FillBuy(order)
	} else {
		ob.actions <- NewSellAction(order)
		ob.FillSell(order)
	}
	if order.Amount > 0 {
		ob.openOrder(order)
	}
}

func (ob *OrderBook) openOrder(order *Order) {
	pricePoint := ob.prices[order.Price]
	pricePoint.Insert(order)
	order.status = OPEN
	if order.OrderType == BUY && order.Price > ob.bid {
		ob.bid = order.Price
	} else if order.OrderType == SELL && order.Price < ob.ask {
		ob.ask = order.Price
	}

	ob.orderIndex[order.Id] = order
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

func (ob *OrderBook) fill(order, pricePointOrderHead *Order) {
	if pricePointOrderHead.Amount >= order.Amount {
		ob.actions <- NewFilledAction(order, pricePointOrderHead)
		pricePointOrderHead.Amount -= order.Amount
		order.Amount = 0
		order.status = FILLED
		return
	} else {
		if pricePointOrderHead.Amount > 0 {
			ob.actions <- NewPartialFilledAction(order, pricePointOrderHead)
			order.Amount -= pricePointOrderHead.Amount
			order.status = PARTIAL_FILLED
			pricePointOrderHead.Amount = 0
		}
	}
}

func (ob *OrderBook) Done() {
	ob.actions <- NewDoneAction()
}
