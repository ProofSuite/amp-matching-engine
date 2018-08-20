package ws

import (
	"errors"

	"github.com/Proofsuite/amp-matching-engine/types"
	"github.com/gorilla/websocket"
)

var orderBookSocket *OrderBookSocket

// OrderBookSocket holds the map of subscribtions subscribed to pair channels
// corresponding to the key/event they have subscribed to.
type OrderBookSocket struct {
	subscriptions map[string]map[*websocket.Conn]bool
}

// GetPairSockets return singleton instance of PairSockets type struct
func GetOrderBookSocket() *OrderBookSocket {
	if orderBookSocket == nil {
		orderBookSocket = &OrderBookSocket{make(map[string]map[*websocket.Conn]bool)}
	}

	return orderBookSocket
}

// Register handles the registration of connection to get
// streaming data over the socker for any pair.
// pair := utils.GetPairKey(bt, qt)
func (s *OrderBookSocket) Subscribe(channelId string, conn *websocket.Conn) error {
	if conn == nil {
		return errors.New("Empty connection object")
	}

	if s.subscriptions[channelId] == nil {
		s.subscriptions[channelId] = make(map[*websocket.Conn]bool)
	}

	s.subscriptions[channelId][conn] = true
	return nil
}

// UnsubscribeHandler returns function of type unsubscribe handler,
// it handles the unsubscription of pair in case of connection closing.
func (s *OrderBookSocket) UnsubscribeHandler(channelId string) func(conn *websocket.Conn) {
	return func(conn *websocket.Conn) {
		s.Unsubscribe(channelId, conn)
	}
}

// UnregisterConnection is used to unsubscribe the connection from listening to the key
// subscribed to. It can be called on unsubscription message from user or due to some other reason by
// system
func (s *OrderBookSocket) Unsubscribe(channelId string, conn *websocket.Conn) {
	if s.subscriptions[channelId][conn] {
		s.subscriptions[channelId][conn] = false
		delete(s.subscriptions[channelId], conn)
	}
}

// Broadcast Message streams message to all the subscribtions subscribed to the pair
func (s *OrderBookSocket) BroadcastMessage(channelId string, msgType string, p *types.WebSocketPayload) error {
	for conn, status := range s.subscriptions[channelId] {
		if status {
			SendOrderBookMessage(conn, msgType, p)
		}
	}

	return nil
}

// SendMessage sends a message on the orderbook channel
func SendOrderBookMessage(conn *websocket.Conn, msgType string, data interface{}) {
	SendMessage(conn, OrderBookChannel, msgType, data)
}

// SendErrorMessage sends
func SendOrderBookErrorMessage(conn *websocket.Conn, data interface{}) {
	SendOrderBookMessage(conn, "ERROR", data)
}

func SendOrderBookInitMessage(conn *websocket.Conn, data interface{}) {
	SendOrderBookMessage(conn, "INIT", data)
}

func SendOrderBookUpdateMessage(conn *websocket.Conn, data interface{}) {
	SendOrderBookMessage(conn, "UPDATE", data)
}

// // SendErrorMessage is responsible for sending error messages on orderbook channel
// func (ps *PairSockets) SendErrorMessage(conn *websocket.Conn, msg interface{}) {
// 	ps.SendMessage(conn, "ERROR", msg)
// }

// // SendBookMessage is responsible for sending complete order book on subscription request
// func (ps *PairSockets) SendBookMessage(conn *websocket.Conn, msg interface{}) {
// 	ps.SendMessage(conn, "INIT", msg)
// }

// func (ps *PairSockets) SendMessage(conn *websocket.Conn, msgType string, msg interface{}) {
// 	SendMessage(conn, OrderBookChannel, msgType, msg)
// }

// // SendErrorMessage is responsible for sending error messages on orderbook channel
// func (ps *PairSockets) SendErrorMessage(conn *websocket.Conn, msg interface{}) {
// 	ps.SendMessage(conn, "ERROR", msg)
// }

// // SendBookMessage is responsible for sending complete order book on subscription request
// func (ps *PairSockets) SendBookMessage(conn *websocket.Conn, msg interface{}) {
// 	ps.SendMessage(conn, "INIT", msg)
// }

// func (s *OrderBookSocket) Subscribe(bt, qt common.Address, conn *websocket.Conn) error {
// 	pair := utils.GetPairKey(bt, qt)
// 	if conn == nil {
// 		return errors.New("nil not allowed in arguments as *websocket.Conn")
// 	} else if s.subscribtions == nil {
// 		s.subscribtions = make(map[string]map[*websocket.Conn]bool)
// 		s.subscribtions[pair] = make(map[*websocket.Conn]bool)
// 	} else if s.subscribtions[pair] == nil {
// 		s.subscribtions[pair] = make(map[*websocket.Conn]bool)
// 	}

// 	if !s.subscribtions[pair][conn] {
// 		s.subscribtions[pair][conn] = true
// 	}

// 	return nil
// }
