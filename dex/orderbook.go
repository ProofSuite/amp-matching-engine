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
	prices     [MAX_PRICE]*PricePoint
	actions    chan *Action
	logger     []*Action
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

// Insert adds a new order to a pricepoint in the orderbook.
func (p *PricePoint) Insert(order *Order) {
	if p.orderHead == nil {
		p.orderHead = order
		p.orderTail = order
	} else {
		p.orderTail.next = order
		p.orderTail = order
	}
}

// AddOrder adds a new order to the orderbook. First it checks whether the order is a buy or a sell
// Depending on the order type, the orderbook tries to fill the order against the existing orders of the opposite type in the
// orderbook. If any amount is left, the orderbook opens an additional order with the remaining amount
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

// OpenOrder opens a additional order in the orderbook. It first adds the order
// in the pricepoint mapping corresponding to the given order
// Then, in case of a BUY order for example, the bid price of the orderbook is changed in case the new order
// has a higher price ("best buy order from seller's perspective") than the current bid.
// It works similarly in the case of a SELL order.
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

// CancelOrder removes an order from the orderbook
func (ob *OrderBook) CancelOrder(h Hash) {
	if order, ok := ob.orderIndex[h]; ok {
		order.Amount = 0
		order.status = CANCELLED
		ob.actions <- NewCancelAction(h)
	}
}

// CancelTrade is called when a blockchain transaction execution fails. Then the order needs to
// be re-added to the orderbook.
func (ob *OrderBook) CancelTrade(t *Trade) {
	if order, ok := ob.orderIndex[t.OrderHash]; ok {
		order.Amount = order.Amount + t.Amount.Uint64()
		order.status = OPEN
		ob.actions <- NewCancelTradeAction()
	}
}

// FillBuy tries to fill a BUY order with the existing sell orders.
// In the case, the order price is under the ask price, no orders will be found and the loops thus ends
// Otherwise, FillBuy loops through the different pricepoints by increasing the ask price whenever an
// order is matched
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

// FillSell functions similarly to the FillBuy function
func (ob *OrderBook) FillSell(o *Order) {
	for ob.bid >= o.Price && o.Amount > 0 {
		pricePoint := ob.prices[ob.bid]
		pricePointOrderHead := pricePoint.orderHead
		for pricePointOrderHead != nil {
			ob.fill(o, pricePointOrderHead)
			pricePointOrderHead = pricePointOrderHead.next //perhaps only one of these two lines is necessary ?
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

	amount := big.NewInt(int64(o.Amount))
	trade := NewTrade(pricePointOrderHead, amount, o.Maker)

	pricePointOrderHead.Amount -= o.Amount
	o.Amount = 0
	o.status = FILLED

	if o.events != nil {
		o.events <- o.NewOrderFilledEvent(trade)
	}

	ob.actions <- NewFilledAction(o, pricePointOrderHead, amount.Uint64())
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
