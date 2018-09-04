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
			SendOHLCVUpdateMessage(conn, p)
		}
	}

	return nil
}

// SendOHLCVMessage sends a message on the ohlcv channel
func SendOHLCVMessage(conn *Conn, msgType string, p interface{}) {
	SendMessage(conn, OHLCVChannel, msgType, p)
}

// SendOHLCVErrorMessage is responsible for sending error messages on ohlcv channel
func SendOHLCVErrorMessage(conn *Conn, p interface{}) {
	SendOHLCVMessage(conn, "ERROR", p)
}

// SendOHLCVInitMesssage is responsible for sending complete order book on subscription request
func SendOHLCVInitMesssage(conn *Conn, p interface{}) {
	SendOHLCVMessage(conn, "INIT", p)
}

// SendOHLCVUpdateMessage is responsible for update ohlcv tick messages
func SendOHLCVUpdateMessage(conn *Conn, p interface{}) {
	SendOHLCVMessage(conn, "UPDATE", p)
}
