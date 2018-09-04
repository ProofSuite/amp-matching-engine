package ws

import (
	"errors"
)

var ohlcvSocket *OHLCVSocket

// OHLCVSocket holds the map of subscribtions subscribed to pair channels
// corresponding to the key/event they have subscribed to.
type OHLCVSocket struct {
	subscriptions map[string]map[*Conn]bool
}

// GetOHLCVSocket return singleton instance of PairSockets type struct
func GetOHLCVSocket() *OHLCVSocket {
	if ohlcvSocket == nil {
		ohlcvSocket = &OHLCVSocket{make(map[string]map[*Conn]bool)}
	}

	return ohlcvSocket
}

// Subscribe handles the registration of connection to get
// streaming data over the socker for any pair.
func (s *OHLCVSocket) Subscribe(channelID string, conn *Conn) error {
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
func (s *OHLCVSocket) UnsubscribeHandler(channelID string) func(conn *Conn) {
	return func(conn *Conn) {
		s.Unsubscribe(channelID, conn)
	}
}

// Unsubscribe is used to unsubscribe the connection from listening to the key
// subscribed to. It can be called on unsubscription message from user or due to some other reason by
// system
func (s *OHLCVSocket) Unsubscribe(channelID string, conn *Conn) {
	if s.subscriptions[channelID][conn] {
		s.subscriptions[channelID][conn] = false
		delete(s.subscriptions[channelID], conn)
	}
}

// BroadcastOHLCV Message streams message to all the subscribtions subscribed to the pair
func (s *OHLCVSocket) BroadcastOHLCV(channelID string, p interface{}) error {
	for conn, status := range s.subscriptions[channelID] {
		if status {
			s.SendUpdateMessage(conn, p)
		}
	}

	return nil
}

// SendMessage sends a websocket message on the trade channel
func (s *OHLCVSocket) SendMessage(conn *Conn, msgType string, p interface{}) {
	SendMessage(conn, OHLCVChannel, msgType, p)
}

// SendErrorMessage sends an error message on the trade channel
func (s *OHLCVSocket) SendErrorMessage(conn *Conn, p interface{}) {
	s.SendMessage(conn, "ERROR", p)
}

// SendInitMessage is responsible for sending message on trade ohlcv channel at subscription
func (s *OHLCVSocket) SendInitMessage(conn *Conn, p interface{}) {
	s.SendMessage(conn, "INIT", p)
}

// SendUpdateMessage is responsible for sending message on trade ohlcv channel at subscription
func (s *OHLCVSocket) SendUpdateMessage(conn *Conn, p interface{}) {
	s.SendMessage(conn, "UPDATE", p)
}
