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
	Event  SubsciptionEvent `json:"event"`
	Key    string           `json:"key"`
	Params `json:"params"`
}

type Params struct {
	From     int64  `json:"from"`
	To       int64  `json:"to"`
	Duration int64  `json:"duration"`
	Units    string `json:"units"`
	TickID   string `json:"tickID"`
}
