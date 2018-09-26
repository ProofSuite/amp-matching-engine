package ws

import (
	"sync"

	"github.com/Proofsuite/amp-matching-engine/types"
	"github.com/ethereum/go-ethereum/common"
)

// OrderConn is websocket order connection struct
// It holds the reference to connection and the channel of type OrderMessage
type OrderConnection struct {
	Conn        *Conn
	ReadChannel chan *types.WebSocketPayload
	Active      bool
	Once        sync.Once
}

var orderConnections map[string]*OrderConnection

// GetOrderConn returns the connection associated with an order ID
func GetOrderConnection(hash common.Hash) (conn *Conn) {
	c := orderConnections[hash.Hex()]
	if c == nil {
		logger.Warning("No connection found")
		return nil
	}

	return orderConnections[hash.Hex()].Conn
}

// GetOrderChannel returns the channel associated with an order ID
func GetOrderChannel(h common.Hash) chan *types.WebSocketPayload {
	hash := h.Hex()

	if orderConnections[hash] == nil {
		return nil
	}

	if orderConnections[hash] == nil {
		return nil
	} else if !orderConnections[hash].Active {
		return nil
	}

	return orderConnections[hash].ReadChannel
}

// OrderSocketUnsubscribeHandler returns a function of type unsubscribe handler.
func OrderSocketUnsubscribeHandler(h common.Hash) func(conn *Conn) {
	hash := h.Hex()

	return func(conn *Conn) {
		if orderConnections[hash] != nil {
			orderConnections[hash] = nil
			delete(orderConnections, hash)
		}
	}
}

// RegisterOrderConnection registers a connection with and orderID.
// It is called whenever a message is recieved over order channel
func RegisterOrderConnection(h common.Hash, conn *OrderConnection) {
	hash := h.Hex()

	if orderConnections == nil {
		orderConnections = make(map[string]*OrderConnection)
	}

	if orderConnections[hash] == nil {
		conn.Active = true
		orderConnections[hash] = conn
	}
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

func SendOrderMessage(msgType string, hash common.Hash, data interface{}) {
	conn := GetOrderConnection(hash)
	if conn == nil {
		return
	}

	SendMessage(conn, OrderChannel, msgType, data, hash)
}
