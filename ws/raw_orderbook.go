package ws

import (
	"errors"
)

var rawOrderBookSocket *RawOrderBookSocket

// RawOrderBookSocket holds the map of subscribtions subscribed to pair channels
// corresponding to the key/event they have subscribed to.
type RawOrderBookSocket struct {
	subscriptions     map[string]map[*Client]bool
	subscriptionsList map[*Client][]string
}

func NewRawOrderBookSocket() *RawOrderBookSocket {
	return &RawOrderBookSocket{
		subscriptions:     make(map[string]map[*Client]bool),
		subscriptionsList: make(map[*Client][]string),
	}
}

// GetRawOrderBookSocket return singleton instance of PairSockets type struct
func GetRawOrderBookSocket() *RawOrderBookSocket {
	if rawOrderBookSocket == nil {
		rawOrderBookSocket = NewRawOrderBookSocket()
	}

	return rawOrderBookSocket
}

// Subscribe handles the subscription of connection to get
// streaming data over the socker for any pair.
// pair := utils.GetPairKey(bt, qt)
func (s *RawOrderBookSocket) Subscribe(channelID string, c *Client) error {
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
func (s *RawOrderBookSocket) UnsubscribeChannelHandler(channelID string) func(c *Client) {
	return func(c *Client) {
		s.UnsubscribeChannel(channelID, c)
	}
}

func (s *RawOrderBookSocket) UnsubscribeHandler() func(c *Client) {
	return func(c *Client) {
		s.Unsubscribe(c)
	}
}

// Unsubscribe is used to unsubscribe the connection from listening to the key
// subscribed to. It can be called on unsubscription message from user or due to some other reason by
// system
func (s *RawOrderBookSocket) UnsubscribeChannel(channelID string, c *Client) {
	if s.subscriptions[channelID][c] {
		s.subscriptions[channelID][c] = false
		delete(s.subscriptions[channelID], c)
	}
}

func (s *RawOrderBookSocket) Unsubscribe(c *Client) {
	channelIDs := s.subscriptionsList[c]
	if channelIDs == nil {
		return
	}

	for _, id := range s.subscriptionsList[c] {
		s.UnsubscribeChannel(id, c)
	}
}

// BroadcastMessage streams message to all the subscribtions subscribed to the pair
func (s *RawOrderBookSocket) BroadcastMessage(channelID string, p interface{}) error {
	for c, status := range s.subscriptions[channelID] {
		if status {
			s.SendUpdateMessage(c, p)
		}
	}

	return nil
}

// SendInitMessage sends INIT message on orderbookchannel on subscription event
func (s *RawOrderBookSocket) SendInitMessage(c *Client, data interface{}) {
	c.SendMessage(RawOrderBookChannel, "INIT", data)
}

// SendUpdateMessage sends UPDATE message on orderbookchannel as new data is created
func (s *RawOrderBookSocket) SendUpdateMessage(c *Client, data interface{}) {
	c.SendMessage(RawOrderBookChannel, "UPDATE", data)
}

func (s *RawOrderBookSocket) SendErrorMessage(c *Client, data interface{}) {
	c.SendMessage(RawOrderBookChannel, "ERROR", data)
}
