package types

import (
	"gopkg.in/mgo.v2/bson"
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
