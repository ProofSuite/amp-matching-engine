package types

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

//To be replaced by WebsocketMessage i think
type ChannelMessage struct {
	Channel string      `json:"channel"`
	Message interface{} `json:"message"`
}

// Message is the model used to send message over socket channel
//To be replaced by WebsocketPayload i think
type Message struct {
	Type string      `json:"type"`
	Hash string      `json:"hash,omitempty"`
	Data interface{} `json:"data"`
}

// Subscription is the model used to send message for subscription to any streaming channel
type Subscription struct {
	Event  SubscriptionEvent `json:"event"`
	Pair   PairSubDoc        `json:"pair"`
	Params `json:"params"`
}

// Params is a sub document used to pass parameters in Subscription messages
type Params struct {
	From     int64  `json:"from"`
	To       int64  `json:"to"`
	Duration int64  `json:"duration"`
	Units    string `json:"units"`
	TickID   string `json:"tickID"`
}

type WebSocketMessage struct {
	Channel string           `json:"channel"`
	Payload WebSocketPayload `json:"payload"`
}

type WebSocketPayload struct {
	Type string      `json:"type"`
	Hash string      `json:"hash,omitempty"`
	Data interface{} `json:"data"`
}

//Data is different for each type of payload

//orders/NEW_ORDER
//Order
//orders/SUBMIT_SIGNATURE
//Order
//Trades
//RemainingOrder
//FillStatus
//MatchingOrders
//orders/CANCEL_ORDER
//Order
//orders/CANCEL_TRADE

//orders/ERROR
//orders/ORDER_ADDED
//orders/ORDER_CANCELED
//orders/REQUEST_SIGNATURE
//orders/TRADE_EXECUTED
//orders/TRADE_TX_SUCCESS
//orders/TRADE_TX_ERROR

//order_book/INIT
//order_book/UPDATE
//trades/INIT
//trades/UPDATE
//ohlcv/INIT
//ohlcv/UPDATE
