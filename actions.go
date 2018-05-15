package main

import "fmt"

type ActionType string

const (
	AT_BUY            = "BUY"
	AT_SELL           = "SELL"
	AT_CANCEL         = "CANCEL"
	AT_CANCELLED      = "CANCELLED"
	AT_PARTIAL_FILLED = "PARTIAL_FILLED"
	AT_FILLED         = "FILLED"
	AT_DONE           = "DONE"
)

type Action struct {
	actionType  ActionType `json:"actionType"`
	pair        Pair       `json:"pair"`
	orderId     uint64     `json:"orderId"`
	fromOrderId uint64     `json:"fromOrderId"`
	amount      uint32     `json:"amount"`
	price       uint32     `json:"price"`
}

func (a *Action) String() string {
	return fmt.Sprintf("\n Action{actionType:%v,pair:%v,orderId:%v,fromOrderId:%v,amount:%v,price:%v}",
		a.actionType, a.pair, a.orderId, a.fromOrderId, a.amount, a.price)
}

func NewBuyAction(o *Order) *Action {
	return &Action{actionType: AT_BUY, pair: o.Pair, orderId: o.Id, amount: o.Amount, price: o.Price}
}

func NewSellAction(o *Order) *Action {
	return &Action{actionType: AT_SELL, pair: o.Pair, orderId: o.Id, amount: o.Amount, price: o.Price}
}

func NewCancelAction(id uint64, pair Pair) *Action {
	return &Action{actionType: AT_CANCEL, pair: pair, orderId: id}
}

func NewCancelledAction(id uint64, pair Pair) *Action {
	return &Action{actionType: AT_CANCELLED, pair: pair, orderId: id}
}

func NewFilledAction(o *Order, fromOrder *Order) *Action {
	return &Action{actionType: AT_FILLED, pair: o.Pair, orderId: o.Id, fromOrderId: fromOrder.Id, amount: o.Amount, price: fromOrder.Price}
}

func NewPartialFilledAction(o *Order, fromOrder *Order) *Action {
	return &Action{actionType: AT_PARTIAL_FILLED, pair: o.Pair, orderId: o.Id, fromOrderId: fromOrder.Id, amount: fromOrder.Amount, price: fromOrder.Price}
}

func NewDoneAction() *Action {
	return &Action{actionType: AT_DONE}
}
