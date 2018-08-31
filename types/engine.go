package types

import (
	"math/big"
)

// FillStatus is enum used to signify the filled status of order in engineResponse
type FillStatus int

// Response is the structure of message response sent by engine
type EngineResponse struct {
	FillStatus     string       `json:"fillStatus,omitempty"`
	Order          *Order       `json:"order,omitempty"`
	RemainingOrder *Order       `json:"remainingOrder,omitempty"`
	MatchingOrders []*FillOrder `json:"matchingOrders,omitempty"`
	Trades         []*Trade     `json:"trades,omitempty"`
}

// this const block holds the possible valued of FillStatus
// const (
// 	_ FillStatus = iota
// 	NOMATCH
// 	PARTIAL
// 	FULL
// 	ERROR
// 	CANCELLED
// )

// FillOrder is structure holds the matching order and
// the amount that has been filled by taker order
type FillOrder struct {
	Amount *big.Int
	Order  *Order
}
