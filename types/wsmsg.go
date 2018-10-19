package types

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

// SubscriptionEvent is an enum signifies whether the incoming message is of type Subscribe or unsubscribe
type SubscriptionEvent string

// Enum members for SubscriptionEvent
const (
	SUBSCRIBE   SubscriptionEvent = "SUBSCRIBE"
	UNSUBSCRIBE SubscriptionEvent = "UNSUBSCRIBE"
	Fetch       SubscriptionEvent = "fetch"
)

const TradeChannel = "trades"
const OrderbookChannel = "order_book"
const OrderChannel = "orders"
const OHLCVChannel = "ohlcv"

type WebsocketMessage struct {
	Channel string         `json:"channel"`
	Event   WebsocketEvent `json:"event"`
}

type WebsocketEvent struct {
	Type    string      `json:"type"`
	Hash    string      `json:"hash,omitempty"`
	Payload interface{} `json:"payload"`
}

// Params is a sub document used to pass parameters in Subscription messages
type Params struct {
	From     int64  `json:"from"`
	To       int64  `json:"to"`
	Duration int64  `json:"duration"`
	Units    string `json:"units"`
	PairID   string `json:"pair"`
}

// Signature payload is a websocket message payload struct for the "REQUEST_SIGNATURE" message
//
type SignaturePayload struct {
	Order          *Order            `json:"order"`
	RemainingOrder *Order            `json:"remainingOrder,omitempty"`
	Matches        []*OrderTradePair `json:"matches"`
}

// type OrderPendingPayload struct {
// 	Order *Order `json:"order"`
// 	Trade *Trade `json:"trade"`
// }

type OrderPendingPayload struct {
	Matches []*OrderTradePair `json:"matches"`
}

// type OrderSuccessPayload struct {
// 	Order *Order `json:"order"`
// 	Trade *Trade `json:"trade"`
// }

type OrderSuccessPayload struct {
	Matches []*OrderTradePair `json:"matches"`
}

type SubscriptionPayload struct {
	PairName   string         `json:"pairName,omitempty"`
	QuoteToken common.Address `json:"quoteToken,omitempty"`
	BaseToken  common.Address `json:"baseToken,omitempty"`
	From       int64          `json"from"`
	To         int64          `json:"to"`
	Duration   int64          `json:"duration"`
	Units      string         `json:"units"`
}

func NewOrderWebsocketMessage(o *Order) *WebsocketMessage {
	return &WebsocketMessage{
		Channel: "orders",
		Event: WebsocketEvent{
			Type:    "NEW_ORDER",
			Hash:    o.Hash.Hex(),
			Payload: o,
		},
	}
}

func NewOrderAddedWebsocketMessage(o *Order, p *Pair, filled int64) *WebsocketMessage {
	o.Process(p)
	o.FilledAmount = big.NewInt(filled)
	o.Status = "OPEN"
	return &WebsocketMessage{
		Channel: "orders",
		Event: WebsocketEvent{
			Type:    "ORDER_ADDED",
			Hash:    o.Hash.Hex(),
			Payload: o,
		},
	}
}

func NewOrderCancelWebsocketMessage(oc *OrderCancel) *WebsocketMessage {
	return &WebsocketMessage{
		Channel: "orders",
		Event: WebsocketEvent{
			Type:    "CANCEL_ORDER",
			Hash:    oc.Hash.Hex(),
			Payload: oc,
		},
	}
}

func NewRequestSignaturesWebsocketMessage(h common.Hash, order *Order, remainingOrder *Order, matches []*OrderTradePair) *WebsocketMessage {
	return &WebsocketMessage{
		Channel: "orders",
		Event: WebsocketEvent{
			Type: "REQUEST_SIGNATURE",
			Hash: h.Hex(),
			Payload: SignaturePayload{
				order,
				remainingOrder,
				matches,
			},
		},
	}
}

func NewSubmitSignatureWebsocketMessage(hash string, order *Order, remainingOrder *Order, matches []*OrderTradePair) *WebsocketMessage {
	return &WebsocketMessage{
		Channel: "orders",
		Event: WebsocketEvent{
			Type: "SUBMIT_SIGNATURE",
			Hash: hash,
			Payload: SignaturePayload{
				order,
				remainingOrder,
				matches,
			},
		},
	}
}
