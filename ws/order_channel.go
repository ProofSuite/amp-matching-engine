package ws

import (
	"sync"

	"github.com/Proofsuite/amp-matching-engine/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/gorilla/websocket"
)

// OrderConn is websocket order connection struct
// It holds the reference to connection and the channel of type OrderMessage
type OrderConn struct {
	Conn        *websocket.Conn
	ReadChannel chan *types.Message
	Active      bool
	Once        sync.Once
}

var orderConnections map[string]*OrderConn

// OrderSocketUnsubscribeHandler returns a function of type unsubscribe handler.
func OrderSocketUnsubscribeHandler(h common.Hash) func(conn *websocket.Conn) {
	hash := h.Hex()

	return func(conn *websocket.Conn) {
		if orderConnections[hash] != nil {
			orderConnections[hash] = nil
			delete(orderConnections, hash)
		}
	}
}

// RegisterOrderConnection registers a connection with and orderID.
// It is called whenever a message is recieved over order channel
func RegisterOrderConnection(h common.Hash, conn *OrderConn) {
	hash := h.Hex()

	if orderConnections == nil {
		orderConnections = make(map[string]*OrderConn)
	}
	if orderConnections[hash] == nil {
		conn.Active = true
		orderConnections[hash] = conn
	}
}

// GetOrderConn returns the connection associated with an order ID
func GetOrderConn(hash common.Hash) (conn *websocket.Conn) {
	return orderConnections[hash.Hex()].Conn
}

// GetOrderChannel returns the channel associated with an order ID
func GetOrderChannel(h common.Hash) chan *types.Message {
	hash := h.Hex()

	if orderConnections[hash] == nil {
		return nil
	} else if !orderConnections[hash].Active {
		return nil
	}
	return orderConnections[hash].ReadChannel
}

// CloseOrderReadChannel is called whenever an order processing is done
// and no further messages are to be accepted for an hash
func CloseOrderReadChannel(h common.Hash) error {
	hash := h.Hex()

	orderConnections[hash].Once.Do(func() {
		close(orderConnections[hash].ReadChannel)
		orderConnections[hash].Active = false
	})
	return nil
}

// OrderSendMessage is responsible for sending message on order channel
func OrderSendMessage(conn *websocket.Conn, msgType string, msg interface{}, hash ...common.Hash) {
	SendMessage(conn, OrderChannel, msgType, msg, hash...)
}

// OrderSendErrorMessage is responsible for sending error message on order channel
func OrderSendErrorMessage(conn *websocket.Conn, msg interface{}, hash ...common.Hash) {
	OrderSendMessage(conn, "ERROR", msg, hash...)
}
