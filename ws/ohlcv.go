package ws

import (
	"errors"

	"github.com/gorilla/websocket"
)

var ohlcvSocket *OHLCVSocket

// OHLCVSocket holds the map of subscribtions subscribed to pair channels
// corresponding to the key/event they have subscribed to.
type OHLCVSocket struct {
	subscriptions map[string]map[*websocket.Conn]bool
}

// GetOHLCVSocket return singleton instance of PairSockets type struct
func GetOHLCVSocket() *OHLCVSocket {
	if ohlcvSocket == nil {
		ohlcvSocket = &OHLCVSocket{make(map[string]map[*websocket.Conn]bool)}
	}

	return ohlcvSocket
}

// Register handles the registration of connection to get
// streaming data over the socker for any pair.
func (s *OHLCVSocket) Subscribe(channelId string, conn *websocket.Conn) error {
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
func (s *OHLCVSocket) UnsubscribeHandler(channelId string) func(conn *websocket.Conn) {
	return func(conn *websocket.Conn) {
		s.Unsubscribe(channelId, conn)
	}
}

// UnregisterConnection is used to unsubscribe the connection from listening to the key
// subscribed to. It can be called on unsubscription message from user or due to some other reason by
// system
func (s *OHLCVSocket) Unsubscribe(channelId string, conn *websocket.Conn) {
	if s.subscriptions[channelId][conn] {
		s.subscriptions[channelId][conn] = false
		delete(s.subscriptions[channelId], conn)
	}
}

// Broadcast Message streams message to all the subscribtions subscribed to the pair
func (s *OHLCVSocket) BroadcastMessage(channelId string, msgType string, p interface{}) error {
	for conn, status := range s.subscriptions[channelId] {
		if status {
			SendOHLCVMessage(conn, msgType, p)
		}
	}

	return nil
}

// SendMessage sends a message on the orderbook channel
func SendOHLCVMessage(conn *websocket.Conn, msgType string, p interface{}) {
	SendMessage(conn, OHLCVChannel, msgType, p)
}

// SendErrorMessage sends
func SendOHLCVErrorMessage(conn *websocket.Conn, p interface{}) {
	SendOHLCVMessage(conn, "ERROR", p)
}

func SendOHLCVInitMesssage(conn *websocket.Conn, p interface{}) {
	SendOHLCVMessage(conn, "INIT", p)
}

func SendOHLCVUpdateMessage(conn *websocket.Conn, p interface{}) {
	SendOHLCVMessage(conn, "UPDATE", p)
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

// func (s *OHLCVSocket) Subscribe(bt, qt common.Address, conn *websocket.Conn) error {
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
