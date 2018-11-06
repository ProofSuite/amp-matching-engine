package types

import (
	"fmt"
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

func (ev *WebsocketMessage) String() string {
	return fmt.Sprintf("%v/%v", ev.Channel, ev.Event.String())
}

type WebsocketEvent struct {
	Type    string      `json:"type"`
	Hash    string      `json:"hash,omitempty"`
	Payload interface{} `json:"payload"`
}

func (ev *WebsocketEvent) String() string {
	return fmt.Sprintf("%v", ev.Type)
}

// Params is a sub document used to pass parameters in Subscription messages
type Params struct {
	From     int64  `json:"from"`
	To       int64  `json:"to"`
	Duration int64  `json:"duration"`
	Units    string `json:"units"`
	PairID   string `json:"pair"`
}

type OrderPendingPayload struct {
	Matches *Matches `json:"matches"`
}

// type OrderSuccessPayload struct {
// 	Order *Order `json:"order"`
// 	Trade *Trade `json:"trade"`
// }

type OrderSuccessPayload struct {
	Matches *Matches `json:"matches"`
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
