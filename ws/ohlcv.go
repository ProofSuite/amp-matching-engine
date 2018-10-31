package ws

import (
	"errors"
)

var ohlcvSocket *OHLCVSocket

// OHLCVSocket holds the map of subscribtions subscribed to pair channels
// corresponding to the key/event they have subscribed to.
type OHLCVSocket struct {
	subscriptions     map[string]map[*Client]bool
	subscriptionsList map[*Client][]string
}

func NewOHLCVSocket() *OHLCVSocket {
	return &OHLCVSocket{
		subscriptions:     make(map[string]map[*Client]bool),
		subscriptionsList: make(map[*Client][]string),
	}
}

// GetOHLCVSocket return singleton instance of PairSockets type struct
func GetOHLCVSocket() *OHLCVSocket {
	if ohlcvSocket == nil {
		ohlcvSocket = NewOHLCVSocket()
	}

	return ohlcvSocket
}

// Subscribe handles the registration of connection to get
// streaming data over the socker for any pair.
func (s *OHLCVSocket) Subscribe(channelID string, c *Client) error {
	if c == nil {
		return errors.New("No connection found")
	}

	if s.subscriptions[channelID] == nil {
		s.subscriptions[channelID] = make(map[*Client]bool)
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
func (s *OHLCVSocket) UnsubscribeChannelHandler(channelID string) func(c *Client) {
	return func(c *Client) {
		s.UnsubscribeChannel(channelID, c)
	}
}

func (s *OHLCVSocket) UnsubscribeHandler() func(c *Client) {
	return func(c *Client) {
		s.Unsubscribe(c)
	}
}

// Unsubscribe is used to unsubscribe the connection from listening to the key
// subscribed to. It can be called on unsubscription message from user or due to some other reason by
// system
func (s *OHLCVSocket) UnsubscribeChannel(channelID string, c *Client) {
	if s.subscriptions[channelID][c] {
		s.subscriptions[channelID][c] = false
		delete(s.subscriptions[channelID], c)
	}
}

func (s *OHLCVSocket) Unsubscribe(c *Client) {
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
func (s *OHLCVSocket) SendMessage(c *Client, msgType string, p interface{}) {
	c.SendMessage(OHLCVChannel, msgType, p)
}

// SendErrorMessage sends an error message on the trade channel
func (s *OHLCVSocket) SendErrorMessage(c *Client, p interface{}) {
	c.SendMessage(OHLCVChannel, "ERROR", p)
}

// SendInitMessage is responsible for sending message on trade ohlcv channel at subscription
func (s *OHLCVSocket) SendInitMessage(c *Client, p interface{}) {
	c.SendMessage(OHLCVChannel, "INIT", p)
}

// SendUpdateMessage is responsible for sending message on trade ohlcv channel at subscription
func (s *OHLCVSocket) SendUpdateMessage(c *Client, p interface{}) {
	c.SendMessage(OHLCVChannel, "UPDATE", p)
}
