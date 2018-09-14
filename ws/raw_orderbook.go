package ws

import (
	"errors"
)

var fullOrderBookSocket *RawOrderBookSocket

// RawOrderBookSocket holds the map of subscribtions subscribed to pair channels
// corresponding to the key/event they have subscribed to.
type RawOrderBookSocket struct {
	subscriptions map[string]map[*Conn]bool
}

// GetRawOrderBookSocket return singleton instance of PairSockets type struct
func GetRawOrderBookSocket() *RawOrderBookSocket {
	if fullOrderBookSocket == nil {
		fullOrderBookSocket = &RawOrderBookSocket{make(map[string]map[*Conn]bool)}
	}

	return fullOrderBookSocket
}

// Subscribe handles the subscription of connection to get
// streaming data over the socker for any pair.
// pair := utils.GetPairKey(bt, qt)
func (s *RawOrderBookSocket) Subscribe(channelID string, conn *Conn) error {
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
func (s *RawOrderBookSocket) UnsubscribeHandler(channelID string) func(conn *Conn) {
	return func(conn *Conn) {
		s.Unsubscribe(channelID, conn)
	}
}

// Unsubscribe is used to unsubscribe the connection from listening to the key
// subscribed to. It can be called on unsubscription message from user or due to some other reason by
// system
func (s *RawOrderBookSocket) Unsubscribe(channelID string, conn *Conn) {
	if s.subscriptions[channelID][conn] {
		s.subscriptions[channelID][conn] = false
		delete(s.subscriptions[channelID], conn)
	}
}

// BroadcastMessage streams message to all the subscribtions subscribed to the pair
func (s *RawOrderBookSocket) BroadcastMessage(channelID string, p interface{}) error {
	for conn, status := range s.subscriptions[channelID] {
		if status {
			s.SendUpdateMessage(conn, p)
		}
	}

	return nil
}

// SendMessage sends a message on the orderbook channel
func (s *RawOrderBookSocket) SendMessage(conn *Conn, msgType string, data interface{}) {
	SendMessage(conn, RawOrderBookChannel, msgType, data)
}

// SendErrorMessage sends error message on orderbookchannel
func (s *RawOrderBookSocket) SendOrderMessage(conn *Conn, msg string) {
	s.SendMessage(conn, "ERROR", map[string]string{"Message": msg})
}

// SendInitMessage sends INIT message on orderbookchannel on subscription event
func (s *RawOrderBookSocket) SendInitMessage(conn *Conn, data interface{}) {
	s.SendMessage(conn, "INIT", data)
}

// SendUpdateMessage sends UPDATE message on orderbookchannel as new data is created
func (s *RawOrderBookSocket) SendUpdateMessage(conn *Conn, data interface{}) {
	s.SendMessage(conn, "UPDATE", data)
}
