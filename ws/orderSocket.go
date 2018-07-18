package ws

import (
	"sync"

	"github.com/Proofsuite/amp-matching-engine/types"
	"github.com/gorilla/websocket"
	"gopkg.in/mgo.v2/bson"
)

// OrderConn is websocket order connection struct
// It holds the reference to connection and the channel of type OrderMessage
type OrderConn struct {
	Conn        *websocket.Conn
	ReadChannel chan *types.OrderMessage
	Active      bool
	Once        sync.Once
}

var orderConnections map[string]*OrderConn

// OrderSocketUnsubscribeHandler returns a function of type unsubscribe handler.
func OrderSocketUnsubscribeHandler(orderID bson.ObjectId) func(conn *websocket.Conn) {
	return func(conn *websocket.Conn) {
		if orderConnections[orderID.Hex()] != nil {
			orderConnections[orderID.Hex()] = nil
			delete(orderConnections, orderID.Hex())
		}
	}
}

// RegisterOrderConnection registers a connection with and orderID.
// It is called whenever a message is recieved over order channel
func RegisterOrderConnection(orderID bson.ObjectId, conn *OrderConn) {
	if orderConnections == nil {
		orderConnections = make(map[string]*OrderConn)
	}
	if orderConnections[orderID.Hex()] == nil {
		conn.Active = true
		orderConnections[orderID.Hex()] = conn
	}
}

// GetOrderConn returns the connection associated with an order ID
func GetOrderConn(orderID bson.ObjectId) (conn *websocket.Conn) {
	return orderConnections[orderID.Hex()].Conn
}

// GetOrderChannel returns the channel associated with an order ID
func GetOrderChannel(orderID bson.ObjectId) chan *types.OrderMessage {
	if orderConnections[orderID.Hex()] == nil {
		return nil
	} else if !orderConnections[orderID.Hex()].Active {
		return nil
	}
	return orderConnections[orderID.Hex()].ReadChannel
}

// CloseOrderReadChannel is called whenever an order processing is done
// and no further messages are to be accepted for an orderID
func CloseOrderReadChannel(orderID bson.ObjectId) error {
	orderConnections[orderID.Hex()].Once.Do(func() {
		close(orderConnections[orderID.Hex()].ReadChannel)
		orderConnections[orderID.Hex()].Active = false
	})
	return nil
}
