package ws

import (
	"errors"
)

var liteOrderBook *LiteOrderBookSocket

// LiteOrderBookSocket holds the map of subscribtions subscribed to pair channels
// corresponding to the key/event they have subscribed to.
type LiteOrderBookSocket struct {
	subscriptions map[string]map[*Conn]bool
}

// GetLiteOrderBookSocket return singleton instance of PairSockets type struct
func GetLiteOrderBookSocket() *LiteOrderBookSocket {
	if liteOrderBook == nil {
		liteOrderBook = &LiteOrderBookSocket{make(map[string]map[*Conn]bool)}
	}

	return liteOrderBook
}

// Subscribe handles the subscription of connection to get
// streaming data over the socker for any pair.
// pair := utils.GetPairKey(bt, qt)
func (s *LiteOrderBookSocket) Subscribe(channelID string, conn *Conn) error {
	if conn == nil {
		return errors.New("Empty connection object")
	}

	if s.subscriptions[channelID] == nil {
		s.subscriptions[channelID] = make(map[*Conn]bool)
	}

	s.subscriptions[channelID][conn] = true
	return nil
}

// UnsubscribeHandler returns function of type unsubscribe handler,
// it handles the unsubscription of pair in case of connection closing.
func (s *LiteOrderBookSocket) UnsubscribeHandler(channelID string) func(conn *Conn) {
	return func(conn *Conn) {
		s.Unsubscribe(channelID, conn)
	}
}

// Unsubscribe is used to unsubscribe the connection from listening to the key
// subscribed to. It can be called on unsubscription message from user or due to some other reason by
// system
func (s *LiteOrderBookSocket) Unsubscribe(channelID string, conn *Conn) {
	if s.subscriptions[channelID][conn] {
		s.subscriptions[channelID][conn] = false
		delete(s.subscriptions[channelID], conn)
	}
}

// BroadcastMessage streams message to all the subscribtions subscribed to the pair
func (s *LiteOrderBookSocket) BroadcastMessage(channelID string, p interface{}) error {
	for conn, status := range s.subscriptions[channelID] {
		if status {
			s.SendUpdateMessage(conn, p)
		}
	}

	return nil
}

// SendMessage sends a message on the orderbook channel
func (s *LiteOrderBookSocket) SendMessage(conn *Conn, msgType string, data interface{}) {
	SendMessage(conn, LiteOrderBookChannel, msgType, data)
}

// SendErrorMessage sends error message on orderbookchannel
func (s *LiteOrderBookSocket) SendErrorMessage(conn *Conn, data interface{}) {
	s.SendMessage(conn, "ERROR", data)
}

// SendInitMessage sends INIT message on orderbookchannel on subscription event
func (s *LiteOrderBookSocket) SendInitMessage(conn *Conn, data interface{}) {
	s.SendMessage(conn, "INIT", data)
}

// SendUpdateMessage sends UPDATE message on orderbookchannel as new data is created
func (s *LiteOrderBookSocket) SendUpdateMessage(conn *Conn, data interface{}) {
	s.SendMessage(conn, "UPDATE", data)
}
