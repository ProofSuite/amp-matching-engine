package types

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

// SubscriptionEvent is an enum signifies whether the incoming message is of type Subscribe or unsubscribe
type SubscriptionEvent string

// Enum members for SubscriptionEvent
const (
	SUBSCRIBE   SubscriptionEvent = "subscribe"
	UNSUBSCRIBE SubscriptionEvent = "unsubscribe"
	Fetch       SubscriptionEvent = "fetch"
)

const TradeChannel = "trades"
const OrderbookChannel = "order_book"
const OrderChannel = "orders"
const OHLCVChannel = "ohlcv"

type WebSocketMessage struct {
	Channel string           `json:"channel"`
	Payload WebSocketPayload `json:"payload"`
}

type WebSocketPayload struct {
	Type string      `json:"type"`
	Hash string      `json:"hash,omitempty"`
	Data interface{} `json:"data"`
}

type WebSocketSubscription struct {
	Event  SubscriptionEvent `json:"event"`
	Pair   PairAddresses     `json:"pair"`
	Params `json:"params"`
}

// Params is a sub document used to pass parameters in Subscription messages
type Params struct {
	From     int64  `json:"from"`
	To       int64  `json:"to"`
	Duration int64  `json:"duration"`
	Units    string `json:"units"`
	PairID   string `json:"pair"`
}

type SignaturePayload struct {
	Order   *Order            `json:"order"`
	Matches []*OrderTradePair `json:"matches"`
}

func NewOrderWebsocketMessage(o *Order) *WebSocketMessage {
	return &WebSocketMessage{
		Channel: "orders",
		Payload: WebSocketPayload{
			Type: "NEW_ORDER",
			Hash: o.Hash.Hex(),
			Data: o,
		},
	}
}

func NewOrderAddedWebsocketMessage(o *Order, p *Pair, filled int64) *WebSocketMessage {
	o.Process(p)
	o.FilledAmount = big.NewInt(filled)
	o.Status = "OPEN"
	return &WebSocketMessage{
		Channel: "orders",
		Payload: WebSocketPayload{
			Type: "ORDER_ADDED",
			Hash: o.Hash.Hex(),
			Data: o,
		},
	}
}

func NewOrderCancelWebsocketMessage(oc *OrderCancel) *WebSocketMessage {
	return &WebSocketMessage{
		Channel: "orders",
		Payload: WebSocketPayload{
			Type: "CANCEL_ORDER",
			Hash: oc.Hash.Hex(),
			Data: oc,
		},
	}
}

func NewRequestSignaturesWebsocketMessage(hash common.Hash, m []*OrderTradePair, o *Order) *WebSocketMessage {
	return &WebSocketMessage{
		Channel: "orders",
		Payload: WebSocketPayload{
			Type: "REQUEST_SIGNATURE",
			Hash: hash.Hex(),
			Data: SignaturePayload{o, m},
		},
	}
}

func NewSubmitSignatureWebsocketMessage(hash string, m []*OrderTradePair, o *Order) *WebSocketMessage {
	return &WebSocketMessage{
		Channel: "orders",
		Payload: WebSocketPayload{
			Type: "SUBMIT_SIGNATURE",
			Hash: hash,
			Data: SignaturePayload{o, m},
		},
	}
}
