package ws

import (
	"github.com/Proofsuite/amp-matching-engine/types"
	"github.com/gorilla/websocket"
	"labix.org/v2/mgo/bson"
)

type Ws struct {
	Conn        *websocket.Conn
	ReadChannel chan *types.OrderMessage
}

var Connections map[string]*Ws

func OrderSocketCloseHandler(orderId bson.ObjectId) func(code int, text string) error {
	return func(code int, text string) error {
		if Connections[orderId.Hex()] != nil {
			Connections[orderId.Hex()] = nil
			delete(Connections, orderId.Hex())
		}
		return nil
	}
}
