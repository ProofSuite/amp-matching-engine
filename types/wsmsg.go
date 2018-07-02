package types

import (
	"labix.org/v2/mgo/bson"
)

type WsMsg struct {
	MsgType string        `json:"msgType"`
	OrderID bson.ObjectId `json:"orderId"`
	Data    interface{}   `json:"data"`
}
