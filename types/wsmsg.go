package types

import (
	"labix.org/v2/mgo/bson"
)

type SubsciptionEvent string

const (
	SUBSCRIBE   SubsciptionEvent = "subscribe"
	UNSUBSCRIBE SubsciptionEvent = "unsubscribe"
)

type OrderMessage struct {
	MsgType string        `json:"msgType"`
	OrderID bson.ObjectId `json:"orderId"`
	Data    interface{}   `json:"data"`
}

type Subscription struct {
	Event SubsciptionEvent `json:"event"`
	Key   string           `json:"key"`
}
