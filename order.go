package main

import (
	"encoding/json"
	"errors"
	"fmt"
)

type Pair struct {
	BaseToken  string `json:"baseToken"`
	QuoteToken string `json:"quoteToken"`
}

func NewPair(baseToken string, quoteToken string) Pair {
	return Pair{BaseToken: baseToken, QuoteToken: quoteToken}
}

type OrderStatus int

const (
	NEW OrderStatus = iota
	OPEN
	PARTIAL_FILLED
	FILLED
	CANCELLED
	PENDING
)

func (orderStatus *OrderStatus) UnmarshalJSON(data []byte) error {
	var s string
	err := json.Unmarshal(data, &s)
	if err != nil {
		return err
	}

	value, ok := map[string]OrderStatus{
		"NEW":            NEW,
		"OPEN":           OPEN,
		"PARTIAL_FILLED": PARTIAL_FILLED,
		"FILLED":         FILLED,
		"CANCELLED":      CANCELLED,
		"PENDING":        PENDING}[s]
	if !ok {
		return errors.New("Invalid Enum Status Value")
	}
	*orderStatus = value
	return nil
}

func (orderStatus *OrderStatus) MarshalJSON() ([]byte, error) {
	value, ok := map[OrderStatus]string{
		NEW:            "NEW",
		OPEN:           "OPEN",
		PARTIAL_FILLED: "PARTIAL_FILLED",
		FILLED:         "FILLED",
		CANCELLED:      "CANCELLED",
		PENDING:        "PENDING"}[*orderStatus]
	if !ok {
		return nil, errors.New("Invalid Enum Type")
	}
	return json.Marshal(value)
}

type OrderType int

const (
	BUY OrderType = iota
	SELL
)

func (orderType *OrderType) UnmarshalJSON(data []byte) error {
	var s string
	err := json.Unmarshal(data, &s)
	if err != nil {
		return err
	}

	value, ok := map[string]OrderType{"BUY": BUY, "SELL": SELL}[s]
	if !ok {
		return errors.New("Invalid Enum Type Value")
	}
	*orderType = value
	return nil
}

func (orderType *OrderType) MarshalJSON() ([]byte, error) {
	value, ok := map[OrderType]string{BUY: "BUY", SELL: "SELL"}[*orderType]
	if !ok {
		return nil, errors.New("Invalid Enum Type")
	}
	return json.Marshal(value)
}

//Each order is linked to the next order at the same price point
type OrderData struct {
	Id            uint64 `json:"id"`
	Price         uint32 `json:"price"`
	BaseToken     string `json:"baseToken"`
	QuoteToken    string `json:"quoteToken"`
	Amount        uint32 `json:"amount"`
	InitialAmount uint32 `json:"initialAmount"`
	OrderType     OrderType
	status        OrderStatus
}

func (o *OrderData) String() string {
	return fmt.Sprintf("Order{id:%d,baseToken:%v,quoteToken:%v,orderType:%v,price:%d,amount:%d}", o.Id, o.BaseToken, o.QuoteToken, o.OrderType, o.Price, o.Amount)
}

type Order struct {
	Id            uint64 `json:"id"`
	Pair          Pair   `json:"pair"`
	Price         uint32 `json:"price"`
	Amount        uint32 `json:"amount"`
	InitialAmount uint32 `json:"initialAmount"`
	OrderType     OrderType
	status        OrderStatus
	next          *Order
	client        *ClientInterface
}

func (order *Order) String() string {
	return fmt.Sprintf("Order{id:%d,pair:%v,orderType:%v,price:%d,amount:%d}", order.Id, order.Pair, order.OrderType, order.Price, order.Amount)
}

func NewOrder(id uint64, pair Pair, orderType OrderType, price uint32, amount uint32) *Order {
	return &Order{Id: id, Pair: pair, OrderType: orderType, Price: price, Amount: amount, InitialAmount: amount, status: NEW}
}

func (o *Order) signalOrderCreated() {
	o.client.signals <- NewOrderPlacedSignal(o)
}

// func (o *Order) cancelOrderSignal() {
// 	o.client.actions <- NewCancelOrderSignal(o)
// }

// func (o *Order) cancelledOrderSignal() {
// 	o.client.actions <- NewCancelledOrderSignal(o)
// }

// func (o *Order) cancelTradeSignal() {
// 	o.client.actions <- NewCancelTradeSignal(o)
// }

func (o *Order) signalOrderMatched() {
	o.client.signals <- NewOrderMatchedSignal(o)
}

func (o *Order) signalOrderPartiallyFilled() {
	o.client.signals <- NewOrderPartiallyFilledSignal(o)
}

func (o *Order) signalOrderFilled() {
	o.client.signals <- NewOrderFilledSignal(o)
}
