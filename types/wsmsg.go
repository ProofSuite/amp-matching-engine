package types

import (
	"gopkg.in/mgo.v2/bson"
)

// SubsciptionEvent is an enum signifies whether the incoming message is of type Subscribe or unsubscribe
type SubsciptionEvent string

// Enum members for SubsciptionEvent
const (
	SUBSCRIBE   SubsciptionEvent = "subscribe"
	UNSUBSCRIBE SubsciptionEvent = "unsubscribe"
	Fetch       SubsciptionEvent = "fetch"
)

// OrderMessage is the model used to send message over order_channel
type OrderMessage struct {
	MsgType string        `json:"msgType"`
	OrderID bson.ObjectId `json:"orderId"`
	Data    interface{}   `json:"data"`
}

// Subscription is the model used to send message for subscription to any streaming channel
type Subscription struct {
	Event  SubsciptionEvent `json:"event"`
	Key    string           `json:"key"`
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
