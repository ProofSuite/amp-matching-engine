package ws

import (
	"sync"

	"github.com/Proofsuite/amp-matching-engine/types"
	"github.com/gorilla/websocket"
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
func OrderSocketUnsubscribeHandler(hash string) func(conn *websocket.Conn) {
	return func(conn *websocket.Conn) {
		if orderConnections[hash] != nil {
			orderConnections[hash] = nil
			delete(orderConnections, hash)
		}
	}
}

// RegisterOrderConnection registers a connection with and orderID.
// It is called whenever a message is recieved over order channel
func RegisterOrderConnection(hash string, conn *OrderConn) {
	if orderConnections == nil {
		orderConnections = make(map[string]*OrderConn)
	}
	if orderConnections[hash] == nil {
		conn.Active = true
		orderConnections[hash] = conn
	}
}

// GetOrderConn returns the connection associated with an order ID
func GetOrderConn(hash string) (conn *websocket.Conn) {
	return orderConnections[hash].Conn
}

// GetOrderChannel returns the channel associated with an order ID
func GetOrderChannel(hash string) chan *types.OrderMessage {
	if orderConnections[hash] == nil {
		return nil
	} else if !orderConnections[hash].Active {
		return nil
	}
	return orderConnections[hash].ReadChannel
}

// CloseOrderReadChannel is called whenever an order processing is done
// and no further messages are to be accepted for an hash
func CloseOrderReadChannel(hash string) error {
	orderConnections[hash].Once.Do(func() {
		close(orderConnections[hash].ReadChannel)
		orderConnections[hash].Active = false
	})
	return nil
}
