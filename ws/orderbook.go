package ws

import (
	"errors"

	"github.com/gorilla/websocket"
)

var orderBookSocket *OrderBookSocket

// OrderBookSocket holds the map of subscribtions subscribed to pair channels
// corresponding to the key/event they have subscribed to.
type OrderBookSocket struct {
	subscriptions map[string]map[*websocket.Conn]bool
}

// GetOrderBookSocket return singleton instance of PairSockets type struct
func GetOrderBookSocket() *OrderBookSocket {
	if orderBookSocket == nil {
		orderBookSocket = &OrderBookSocket{make(map[string]map[*websocket.Conn]bool)}
	}

	return orderBookSocket
}

// Subscribe handles the subscription of connection to get
// streaming data over the socker for any pair.
// pair := utils.GetPairKey(bt, qt)
func (s *OrderBookSocket) Subscribe(channelID string, conn *websocket.Conn) error {
	if conn == nil {
		return errors.New("Empty connection object")
	}

	if s.subscriptions[channelID] == nil {
		s.subscriptions[channelID] = make(map[*websocket.Conn]bool)
	}

	s.subscriptions[channelID][conn] = true
	return nil
}

// UnsubscribeHandler returns function of type unsubscribe handler,
// it handles the unsubscription of pair in case of connection closing.
func (s *OrderBookSocket) UnsubscribeHandler(channelID string) func(conn *websocket.Conn) {
	return func(conn *websocket.Conn) {
		s.Unsubscribe(channelID, conn)
	}
}

// Unsubscribe is used to unsubscribe the connection from listening to the key
// subscribed to. It can be called on unsubscription message from user or due to some other reason by
// system
func (s *OrderBookSocket) Unsubscribe(channelID string, conn *websocket.Conn) {
	if s.subscriptions[channelID][conn] {
		s.subscriptions[channelID][conn] = false
		delete(s.subscriptions[channelID], conn)
	}
}

// BroadcastMessage streams message to all the subscribtions subscribed to the pair
func (s *OrderBookSocket) BroadcastMessage(channelID string, p interface{}) error {
	for conn, status := range s.subscriptions[channelID] {
		if status {
			SendOrderBookUpdateMessage(conn, p)
		}
	}

	return nil
}

// SendOrderBookMessage sends a message on the orderbook channel
func SendOrderBookMessage(conn *websocket.Conn, msgType string, data interface{}) {
	SendMessage(conn, OrderBookChannel, msgType, data)
}

// SendOrderBookErrorMessage sends error message on orderbookchannel
func SendOrderBookErrorMessage(conn *websocket.Conn, data interface{}) {
	SendOrderBookMessage(conn, "ERROR", data)
}

// SendOrderBookInitMessage sends INIT message on orderbookchannel on subscription event
func SendOrderBookInitMessage(conn *websocket.Conn, data interface{}) {
	SendOrderBookMessage(conn, "INIT", data)
}

// SendOrderBookUpdateMessage sends UPDATE message on orderbookchannel as new data is created
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
