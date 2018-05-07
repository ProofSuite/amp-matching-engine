package main

import (
	"encoding/json"
	"errors"
	"fmt"
)

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
type Order struct {
	Id        uint64 `json:"id"`
	Symbol    string `json:"symbol"`
	Price     uint32 `json:"price"`
	Amount    uint32 `json:"amount"`
	OrderType OrderType
	status    OrderStatus
	next      *Order
}

func (order *Order) String() string {
	return fmt.Sprintf("\nOrder{id:%d,symbol:%s,orderType:%v,price:%d,amount:%d}", order.Id, order.Symbol, order.OrderType, order.Price, order.Amount)
}

func NewOrder(id uint64, symbol string, orderType OrderType, price uint32, amount uint32) *Order {
	return &Order{Id: id, Symbol: symbol, OrderType: orderType, Price: price, Amount: amount, status: NEW}
}
