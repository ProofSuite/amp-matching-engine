package dex

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
)

type ActionType string

const (
	AT_BUY            = "BUY"
	AT_SELL           = "SELL"
	AT_CANCEL         = "CANCEL"
	AT_CANCELLED      = "CANCELLED"
	AT_PARTIAL_FILLED = "PARTIAL_FILLED"
	AT_FILLED         = "FILLED"
	AT_CANCEL_TRADE   = "AT_CANCEL_TRADE"
	AT_DONE           = "DONE"
)

type Action struct {
	actionType    ActionType  `json:"actionType"`
	hash          common.Hash `json:"hash"`
	pair          TokenPair   `json:"pair"`
	orderId       uint64      `json:"orderId"`
	fromOrderId   uint64      `json:"fromOrderId"`
	orderHash     common.Hash `json:"orderHash"`
	fromOrderHash common.Hash `json:"fromOrderHash"`
	amount        uint64      `json:"amount"`
	price         uint64      `json:"price"`
}

// func (a *Action) String() string {
// 	return fmt.Sprintf("\n Action{actionType:%v,pair:%v,orderId:%v,from")
// }

func (a *Action) String() string {
	return fmt.Sprintf("\n Action{actionType:%v,pair:%v,orderHash:%v,fromOrderHash:%v,amount:%v,price:%v}",
		a.actionType, a.pair, a.orderHash, a.fromOrderHash, a.amount, a.price)
}

func NewBuyAction(o *Order) *Action {
	return &Action{actionType: AT_BUY, orderHash: o.Hash, amount: o.Amount, price: o.Price}
}

func NewSellAction(o *Order) *Action {
	return &Action{actionType: AT_SELL, orderHash: o.Hash, amount: o.Amount, price: o.Price}
}

func NewCancelAction(hash common.Hash) *Action {
	return &Action{actionType: AT_CANCEL, orderHash: hash}
}

func NewCancelTradeAction() *Action {
	return &Action{actionType: AT_CANCEL_TRADE}
}

func NewCancelledAction(hash common.Hash, pair TokenPair) *Action {
	return &Action{actionType: AT_CANCELLED, pair: pair, orderHash: hash}
}

func NewFilledAction(order, fromOrder *Order) *Action {
	return &Action{actionType: AT_FILLED,
		orderHash:     order.Hash,
		fromOrderHash: fromOrder.Hash,
		amount:        fromOrder.Amount,
		price:         fromOrder.Price}
}

func NewPartialFilledAction(order, fromOrder *Order) *Action {
	return &Action{actionType: AT_PARTIAL_FILLED,
		orderHash:     order.Hash,
		fromOrderHash: fromOrder.Hash,
		amount:        fromOrder.Amount,
		price:         fromOrder.Price}
}

func NewDoneAction() *Action {
	return &Action{actionType: AT_DONE}
}
