package types

// SubscriptionEvent is an enum signifies whether the incoming message is of type Subscribe or unsubscribe
type SubscriptionEvent string

// Enum members for SubscriptionEvent
const (
	SUBSCRIBE   SubscriptionEvent = "subscribe"
	UNSUBSCRIBE SubscriptionEvent = "unsubscribe"
	Fetch       SubscriptionEvent = "fetch"
)

// Message is the model used to send message over socket channel
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
