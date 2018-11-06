package ws

import (
	"github.com/Proofsuite/amp-matching-engine/types"
	"github.com/ethereum/go-ethereum/common"
)

// OrderConn is websocket order connection struct
// It holds the reference to connection and the channel of type OrderMessage
type OrderConnection struct {
	Client      *Client
	ReadChannel chan *types.WebsocketEvent
}

var orderConnections map[string]*OrderConnection

// GetOrderConn returns the connection associated with an order ID
func GetOrderConnection(a common.Address) *Client {
	c := orderConnections[a.Hex()]
	if c == nil {
		logger.Warning("No connection found")
		return nil
	}

	return orderConnections[a.Hex()].Client
}

// GetOrderChannel returns the channel associated with an order ID
func GetOrderChannel(a common.Address) chan *types.WebsocketEvent {
	if orderConnections[a.Hex()] == nil {
		return nil
	}

	return orderConnections[a.Hex()].ReadChannel
}

// OrderSocketUnsubscribeHandler returns a function of type unsubscribe handler.
func OrderSocketUnsubscribeHandler(a common.Address) func(client *Client) {
	return func(client *Client) {
		if orderConnections[a.Hex()] != nil {
			logger.Info("Unsubscribing order connection")
			orderConnections[a.Hex()] = nil
			delete(orderConnections, a.Hex())
		}
	}
}

// RegisterOrderConnection registers a connection with and orderID.
// It is called whenever a message is recieved over order channel
func RegisterOrderConnection(a common.Address, c *OrderConnection) {
	if orderConnections == nil {
		orderConnections = make(map[string]*OrderConnection)
	}

	if orderConnections[a.Hex()] == nil {
		logger.Info("Registering a new order connection")
		orderConnections[a.Hex()] = c

		RegisterConnectionUnsubscribeHandler(c.Client, OrderSocketUnsubscribeHandler(a))
	}
}

func SendOrderMessage(msgType string, a common.Address, payload interface{}) {
	c := GetOrderConnection(a)
	if c == nil {
		return
	}

	c.SendMessage(OrderChannel, msgType, payload)
}
