package ws

import (
	"sync"

	"github.com/Proofsuite/amp-matching-engine/types"
	"github.com/gorilla/websocket"
	"gopkg.in/mgo.v2/bson"
)

type WsOrderConn struct {
	Conn        *websocket.Conn
	ReadChannel chan *types.OrderMessage
	Active      bool
	Once        sync.Once
}

var orderConnections map[string]*WsOrderConn

func OrderSocketCloseHandler(orderId bson.ObjectId) func(conn *websocket.Conn) {
	return func(conn *websocket.Conn) {
		if orderConnections[orderId.Hex()] != nil {
			orderConnections[orderId.Hex()] = nil
			delete(orderConnections, orderId.Hex())
		}
	}
}

func RegisterOrderConnection(orderId bson.ObjectId, conn *WsOrderConn) {
	if orderConnections == nil {
		orderConnections = make(map[string]*WsOrderConn)
	}
	if orderConnections[orderId.Hex()] == nil {
		orderConnections[orderId.Hex()] = conn
	}
}

func GetOrderConn(orderId bson.ObjectId) (conn *websocket.Conn) {
	return orderConnections[orderId.Hex()].Conn
}
func GetOrderChannel(orderId bson.ObjectId) chan *types.OrderMessage {
	if !orderConnections[orderId.Hex()].Active {
		return nil
	}
	return orderConnections[orderId.Hex()].ReadChannel
}
func CloseOrderReadChannel(orderId bson.ObjectId) error {
	orderConnections[orderId.Hex()].Once.Do(func() {
		close(orderConnections[orderId.Hex()].ReadChannel)
		orderConnections[orderId.Hex()].Active = false
	})
	return nil
}
