package dex

import (
	"fmt"
	"math/big"

	. "github.com/ethereum/go-ethereum/common"
)

const MAX_PRICE = 10000000

//Pricepoint contains pointers to the first and the last order entered at that price
type PricePoint struct {
	orderHead *Order
	orderTail *Order
}

// The orderbook keeps track of the maximum bid and minimum ask
//The orderIndex is a mapping of OrderIDs to pointers so that we can easily cancel outstanding orders
//prices is an array of all possible pricepoints
//actions is a channels to report some action to a handler as they occur
type OrderBook struct {
	ask        uint64
	bid        uint64
	orderIndex map[Hash]*Order
	// orderIndex map[uint64]*Order
	prices  [MAX_PRICE]*PricePoint
	actions chan *Action
	logger  []*Action
}

// NewOrderbook returns a default orderbook struct
func NewOrderBook(actions chan *Action) *OrderBook {
	ob := new(OrderBook)
	ob.bid = 0
	ob.ask = MAX_PRICE
	ob.actions = actions
	ob.orderIndex = make(map[Hash]*Order)

	for i := range ob.prices {
		ob.prices[i] = new(PricePoint)
	}

	return ob
}

func (ob *OrderBook) String() string {
	return fmt.Sprintf("Ask:%v, Bid:%v, orderIndex:%v", ob.ask, ob.bid, ob.orderIndex)
}

func (ob *OrderBook) GetLogs() []*Action {
	return ob.logger
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
	if o.events != nil {
		o.events <- o.NewOrderPlacedEvent()
	}
	pricePoint := ob.prices[o.Price]
	pricePoint.Insert(o)
	o.status = OPEN
	if o.OrderType == BUY && o.Price > ob.bid {
		ob.bid = o.Price
	} else if o.OrderType == SELL && o.Price < ob.ask {
		ob.ask = o.Price
	}

	ob.orderIndex[o.Hash] = o
}

func (ob *OrderBook) CancelOrder(h Hash) {
	if order, ok := ob.orderIndex[h]; ok {
		order.Amount = 0
		order.status = CANCELLED
		ob.actions <- NewCancelAction(h)
	}
}

func (ob *OrderBook) FillBuy(o *Order) {
	for ob.ask <= o.Price && o.Amount > 0 {
		pricePoint := ob.prices[ob.ask]
		pricePointOrderHead := pricePoint.orderHead
		for pricePointOrderHead != nil {
			ob.fill(o, pricePointOrderHead)
			pricePointOrderHead = pricePointOrderHead.next
			pricePoint.orderHead = pricePointOrderHead
		}
		ob.ask++
	}
}

func (ob *OrderBook) FillSell(o *Order) {
	for ob.bid >= o.Price && o.Amount > 0 {
		pricePoint := ob.prices[ob.bid]
		pricePointOrderHead := pricePoint.orderHead
		for pricePointOrderHead != nil {
			ob.fill(o, pricePointOrderHead)
			pricePointOrderHead = pricePointOrderHead.next //only of these two lines is necessary
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

	amount := big.NewInt(int64(o.Amount))
	trade := NewTrade(pricePointOrderHead, amount, o.Maker)

	pricePointOrderHead.Amount -= o.Amount
	o.Amount = 0
	o.status = FILLED

	if o.events != nil {
		o.events <- o.NewOrderFilledEvent(trade)
	}
	return
}

func (ob *OrderBook) fillPartially(o, pricePointOrderHead *Order) {
	ob.actions <- NewPartialFilledAction(o, pricePointOrderHead)
	o.Amount -= pricePointOrderHead.Amount
	o.status = PARTIAL_FILLED

	if o.events != nil {
		o.events <- o.NewOrderPartiallyFilledEvent()
	}
	return
}

func (ob *OrderBook) Done() {
	ob.actions <- NewDoneAction()
}
