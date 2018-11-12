package ws

import (
	"github.com/ethereum/go-ethereum/common"
)

// OrderConn is websocket order connection struct
// It holds the reference to connection and the channel of type OrderMessage

type OrderConnection []*Client

var orderConnections map[string]OrderConnection

// GetOrderConn returns the connection associated with an order ID
func GetOrderConnections(a common.Address) OrderConnection {
	c := orderConnections[a.Hex()]
	if c == nil {
		logger.Warning("No connection found")
		return nil
	}

	return orderConnections[a.Hex()]
}

func OrderSocketUnsubscribeHandler(a common.Address) func(client *Client) {
	return func(client *Client) {
		orderConnection := orderConnections[a.Hex()]

		if orderConnection != nil {
			for i, c := range orderConnection {
				if client == c {
					orderConnection = append(orderConnection[:i], orderConnection[:+1]...)
				}
			}
		}

		orderConnections[a.Hex()] = orderConnection
	}
}

// RegisterOrderConnection registers a connection with and orderID.
// It is called whenever a message is recieved over order channel
func RegisterOrderConnection(a common.Address, c *Client) {
	if orderConnections == nil {
		orderConnections = make(map[string]OrderConnection)
	}

	if orderConnections[a.Hex()] == nil {
		orderConnections[a.Hex()] = OrderConnection{c}
		RegisterConnectionUnsubscribeHandler(c, OrderSocketUnsubscribeHandler(a))
	}

	if orderConnections[a.Hex()] != nil {
		if !isClientConnected(a, c) {
			orderConnections[a.Hex()] = append(orderConnections[a.Hex()], c)
			RegisterConnectionUnsubscribeHandler(c, OrderSocketUnsubscribeHandler(a))
		}
	}
}

func isClientConnected(a common.Address, client *Client) bool {
	for _, c := range orderConnections[a.Hex()] {
		if c == client {
			return true
		}
	}

	return false
}

func SendOrderMessage(msgType string, a common.Address, payload interface{}) {
	conn := GetOrderConnections(a)
	if conn == nil {
		return
	}

	for _, c := range conn {
		c.SendMessage(OrderChannel, msgType, payload)
	}
}
