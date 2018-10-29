package ws

import (
	"sync"

	"github.com/Proofsuite/amp-matching-engine/types"
	"github.com/Proofsuite/amp-matching-engine/utils"
	"github.com/ethereum/go-ethereum/common"
)

// OrderConn is websocket order connection struct
// It holds the reference to connection and the channel of type OrderMessage
type OrderConnection struct {
	Conn        *Conn
	ReadChannel chan *types.WebsocketEvent
	Active      bool
	Once        sync.Once
}

var orderConnections map[string]*OrderConnection

// GetOrderConn returns the connection associated with an order ID
func GetOrderConnection(a common.Address) (conn *Conn) {
	c := orderConnections[a.Hex()]
	if c == nil {
		logger.Warning("No connection found")
		return nil
	}

	utils.PrintJSON("Connection")
	logger.Info(c)

	return orderConnections[a.Hex()].Conn
}

// GetOrderChannel returns the channel associated with an order ID
func GetOrderChannel(a common.Address) chan *types.WebsocketEvent {
	if orderConnections[a.Hex()] == nil {
		return nil
	}

	if orderConnections[a.Hex()] == nil {
		return nil
	} else if !orderConnections[a.Hex()].Active {
		return nil
	}

	return orderConnections[a.Hex()].ReadChannel
}

// OrderSocketUnsubscribeHandler returns a function of type unsubscribe handler.
func OrderSocketUnsubscribeHandler(a common.Address) func(conn *Conn) {

	return func(conn *Conn) {
		if orderConnections[a.Hex()] != nil {
			orderConnections[a.Hex()] = nil
			delete(orderConnections, a.Hex())
		}
	}
}

// RegisterOrderConnection registers a connection with and orderID.
// It is called whenever a message is recieved over order channel
func RegisterOrderConnection(a common.Address, conn *OrderConnection) {
	if orderConnections == nil {
		orderConnections = make(map[string]*OrderConnection)
	}

	if orderConnections[a.Hex()] == nil {
		conn.Active = true
		orderConnections[a.Hex()] = conn
	}
}

// CloseOrderReadChannel is called whenever an order processing is done
// and no further messages are to be accepted for an hash
func CloseOrderReadChannel(a common.Address) error {
	orderConnections[a.Hex()].Once.Do(func() {
		close(orderConnections[a.Hex()].ReadChannel)
		orderConnections[a.Hex()].Active = false
	})

	return nil
}

func SendOrderMessage(msgType string, a common.Address, h common.Hash, payload interface{}) {
	conn := GetOrderConnection(a)
	if conn == nil {
		return
	}

	SendMessage(conn, OrderChannel, msgType, payload, h)
}
