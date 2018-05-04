package main

import "fmt"

const MAX_PRICE = 10000000

//Pricepoint contains pointers to the first and the last order entered at that price
type PricePoint struct {
	orderHead *Order
	orderTail *Order
}

type OrderStatus int

const (
	NEW OrderStatus = iota
	OPEN
	PARTIAL_FILLED
	FILLED
	CANCELLED
)

type OrderType int

const (
	BUY OrderType = iota
	SELL
)

//Each order is linked to the next order at the same price point
type Order struct {
	id        uint64
	symbol    string
	price     uint32
	amount    uint32
	orderType OrderType
	status    OrderStatus
	next      *Order
}

func (order *Order) String() string {
	return fmt.Sprintf("\nOrder{id:%d,symbol:%s,orderType:%v,price:%d,amount:%d}", order.id, order.symbol, order.orderType, order.price, order.amount)
}

func NewOrder(id uint64, symbol string, orderType OrderType, price uint32, amount uint32) *Order {
	return &Order{id: id, symbol: symbol, orderType: orderType, price: price, amount: amount, status: NEW}
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
}

func (orderbook *OrderBook) String() string {
	return fmt.Sprintf("Ask:%v, Bid:%v, orderIndex:%v", orderbook.ask, orderbook.bid, orderbook.orderIndex)
}

func (pricePoint *PricePoint) Insert(order *Order) {
	fmt.Printf("\nInserting a new price point at: %v", order.price)
	fmt.Printf("\n*************")
	if pricePoint.orderHead == nil {
		pricePoint.orderHead = order
		pricePoint.orderTail = order
	} else {
		pricePoint.orderTail.next = order
		pricePoint.orderTail = order
	}
	fmt.Printf("\nThe price point tail and head for this price point are:")
	fmt.Printf("\nPrice point order tail: %v", pricePoint.orderTail)
	fmt.Printf("\nPrice point order head: %v", pricePoint.orderHead)
	fmt.Printf("\n**************\n")
}

func (ob *OrderBook) AddOrder(order *Order) {
	if order.orderType == BUY {
		ob.actions <- NewBuyAction(order)
		ob.FillBuy(order)
	} else {
		ob.actions <- NewSellAction(order)
		ob.FillSell(order)
	}
	if order.amount > 0 {
		ob.openOrder(order)
	}
}

func (ob *OrderBook) openOrder(order *Order) {
	fmt.Printf("\nOpening a new order")
	fmt.Printf("\n*********************")
	fmt.Printf("\nOrder: %v", order)
	fmt.Printf("\n*********************\n")
	pricePoint := ob.prices[order.price]
	pricePoint.Insert(order)
	order.status = OPEN
	if order.orderType == BUY && order.price > ob.bid {
		fmt.Printf("\nThe order type is BUY and order price is higher than the bid")
		ob.bid = order.price
	} else if order.orderType == SELL && order.price < ob.ask {
		fmt.Printf("\nThe order type is SELL and order price is lower than the ask")
		ob.ask = order.price
	}

	ob.orderIndex[order.id] = order
}

func (ob *OrderBook) CancelOrder(id uint64, symbol string) {
	ob.actions <- NewCancelAction(id, symbol)
	if order, ok := ob.orderIndex[id]; ok {
		order.amount = 0
		order.status = CANCELLED
	}

	ob.actions <- NewCancelledAction(id, symbol)
}

func (ob *OrderBook) FillBuy(order *Order) {

	for ob.ask < order.price && order.amount > 0 {
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
	for ob.bid >= order.price && order.amount > 0 {
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
	if pricePointOrderHead.amount >= order.amount {
		ob.actions <- NewFilledAction(order, pricePointOrderHead)
		pricePointOrderHead.amount -= order.amount
		order.amount = 0
		order.status = FILLED
		return
	} else {
		if pricePointOrderHead.amount > 0 {
			ob.actions <- NewPartialFilledAction(order, pricePointOrderHead)
			order.amount -= pricePointOrderHead.amount
			order.status = PARTIAL_FILLED
			pricePointOrderHead.amount = 0
		}
	}
}

func (ob *OrderBook) Done() {
	ob.actions <- NewDoneAction()
}
