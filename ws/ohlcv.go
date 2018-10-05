package ws

import (
	"errors"
)

var ohlcvSocket *OHLCVSocket

// OHLCVSocket holds the map of subscribtions subscribed to pair channels
// corresponding to the key/event they have subscribed to.
type OHLCVSocket struct {
	subscriptions     map[string]map[*Conn]bool
	subscriptionsList map[*Conn][]string
}

// GetOHLCVSocket return singleton instance of PairSockets type struct
func GetOHLCVSocket() *OHLCVSocket {

	subscriptions := make(map[string]map[*Conn]bool)
	subscriptionsList := make(map[*Conn][]string)

	if ohlcvSocket == nil {
		ohlcvSocket = &OHLCVSocket{subscriptions, subscriptionsList}
	}

	return ohlcvSocket
}

// Subscribe handles the registration of connection to get
// streaming data over the socker for any pair.
func (s *OHLCVSocket) Subscribe(channelID string, c *Conn) error {
	if c == nil {
		return errors.New("No connection found")
	}

	if s.subscriptions[channelID] == nil {
		s.subscriptions[channelID] = make(map[*Conn]bool)
	}

	s.subscriptions[channelID][c] = true

	if s.subscriptionsList[c] == nil {
		s.subscriptionsList[c] = []string{}
	}

	s.subscriptionsList[c] = append(s.subscriptionsList[c], channelID)

	return nil
}

// UnsubscribeHandler returns function of type unsubscribe handler,
// it handles the unsubscription of pair in case of connection closing.
func (s *OHLCVSocket) UnsubscribeChannelHandler(channelID string) func(c *Conn) {
	return func(c *Conn) {
		s.UnsubscribeChannel(channelID, c)
	}
}

func (s *OHLCVSocket) UnsubscribeHandler() func(c *Conn) {
	return func(c *Conn) {
		s.Unsubscribe(c)
	}
}

// Unsubscribe is used to unsubscribe the connection from listening to the key
// subscribed to. It can be called on unsubscription message from user or due to some other reason by
// system
func (s *OHLCVSocket) UnsubscribeChannel(channelID string, c *Conn) {
	if s.subscriptions[channelID][c] {
		s.subscriptions[channelID][c] = false
		delete(s.subscriptions[channelID], c)
	}
}

func (s *OHLCVSocket) Unsubscribe(c *Conn) {
	channelIDs := s.subscriptionsList[c]
	if channelIDs == nil {
		return
	}

	for _, id := range s.subscriptionsList[c] {
		s.UnsubscribeChannel(id, c)
	}
}

// BroadcastOHLCV Message streams message to all the subscribtions subscribed to the pair
func (s *OHLCVSocket) BroadcastOHLCV(channelID string, p interface{}) error {
	for c, status := range s.subscriptions[channelID] {
		if status {
			s.SendUpdateMessage(c, p)
		}
	}

	return nil
}

// SendMessage sends a websocket message on the trade channel
func (s *OHLCVSocket) SendMessage(c *Conn, msgType string, p interface{}) {
	SendMessage(c, OHLCVChannel, msgType, p)
}

// SendErrorMessage sends an error message on the trade channel
func (s *OHLCVSocket) SendErrorMessage(c *Conn, p interface{}) {
	s.SendMessage(c, "ERROR", p)
}

// SendInitMessage is responsible for sending message on trade ohlcv channel at subscription
func (s *OHLCVSocket) SendInitMessage(c *Conn, p interface{}) {
	s.SendMessage(c, "INIT", p)
}

// SendUpdateMessage is responsible for sending message on trade ohlcv channel at subscription
func (s *OHLCVSocket) SendUpdateMessage(c *Conn, p interface{}) {
	s.SendMessage(c, "UPDATE", p)
}
