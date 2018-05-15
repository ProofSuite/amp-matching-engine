package main

import "fmt"

type SignalType string

type Signal struct {
	signalType  SignalType `json:"signalType"`
	pair        Pair       `json:"pair"`
	orderId     uint64     `json:"orderId"`
	fromOrderId uint64     `json:"fromOrderId"`
	amount      uint32     `json:"amount"`
	price       uint32     `json:"price"`
}

func (s *Signal) String() string {
	return fmt.Sprintf("\n Action{actionType:%v,pair:%v,orderId:%v,fromOrderId:%v,amount:%v,price:%v}",
		s.signalType, s.pair, s.orderId, s.fromOrderId, s.amount, s.price)
}

func NewOrderPlacedSignal(o *Order) *Signal {
	return &Signal{signalType: ORDER_PLACED, pair: o.Pair, orderId: o.Id, amount: o.Amount, price: o.Price}
}

func NewOrderMatchedSignal(o *Order) *Signal {
	return &Signal{signalType: ORDER_MATCHED, pair: o.Pair, orderId: o.Id, amount: o.Amount, price: o.Price}
}

func NewOrderPartiallyFilledSignal(o *Order) *Signal {
	return &Signal{signalType: ORDER_PARTIALLY_FILLED, pair: o.Pair, orderId: o.Id, amount: o.Amount, price: o.Price}
}

func NewOrderFilledSignal(o *Order) *Signal {
	return &Signal{signalType: ORDER_FILLED, pair: o.Pair, orderId: o.Id, amount: o.Amount, price: o.Price}
}

func NewOrderCanceledSignal(o *Order) *Signal {
	return &Signal{signalType: ORDER_CANCELED, pair: o.Pair, orderId: o.Id, amount: o.Amount, price: o.Price}
}

func NewDoneSignal() *Signal {
	return &Signal{signalType: DONE}
}

// func NewCancelOrderSignal(o *Order) *Signal {
// 	return &Signal{signalType: CANCEL_ORDER, pair: o.Pair, orderId: id}
// }

// func NewCancelledOrderSignal(o *Order) *Signal {
// 	return &Signal{signalType: CANCELLED_ORDER, pair: o.Pair, orderId: id}
// }

// func NewCancelTradeSignal(o *Order) *Signal {
// 	return &Signal{signalType: CANCEL_TRADE, pair: o.Pair, orderId: id}
// }

// func NewCancelledTradeSignal(o *Order) *Signal {
// 	return &Signal{signalType: CANCELLED_TRADE, pair: o.Pair, orderId: id}
// }
