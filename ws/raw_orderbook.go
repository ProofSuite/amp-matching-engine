package ws

import (
	"errors"
)

var fullOrderBookSocket *RawOrderBookSocket

// RawOrderBookSocket holds the map of subscribtions subscribed to pair channels
// corresponding to the key/event they have subscribed to.
type RawOrderBookSocket struct {
	subscriptions     map[string]map[*Conn]bool
	subscriptionsList map[*Conn][]string
}

// GetRawOrderBookSocket return singleton instance of PairSockets type struct
func GetRawOrderBookSocket() *RawOrderBookSocket {

	subscriptions := make(map[string]map[*Conn]bool)
	subscriptionsList := make(map[*Conn][]string)

	if fullOrderBookSocket == nil {
		fullOrderBookSocket = &RawOrderBookSocket{subscriptions, subscriptionsList}
	}

	return fullOrderBookSocket
}

// Subscribe handles the subscription of connection to get
// streaming data over the socker for any pair.
// pair := utils.GetPairKey(bt, qt)
func (s *RawOrderBookSocket) Subscribe(channelID string, c *Conn) error {
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
func (s *RawOrderBookSocket) UnsubscribeChannelHandler(channelID string) func(c *Conn) {
	return func(c *Conn) {
		s.UnsubscribeChannel(channelID, c)
	}
}

func (s *RawOrderBookSocket) UnsubscribeHandler() func(c *Conn) {
	return func(c *Conn) {
		s.Unsubscribe(c)
	}
}

// Unsubscribe is used to unsubscribe the connection from listening to the key
// subscribed to. It can be called on unsubscription message from user or due to some other reason by
// system
func (s *RawOrderBookSocket) UnsubscribeChannel(channelID string, c *Conn) {
	if s.subscriptions[channelID][c] {
		s.subscriptions[channelID][c] = false
		delete(s.subscriptions[channelID], c)
	}
}

func (s *RawOrderBookSocket) Unsubscribe(c *Conn) {
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

// SendMessage sends a message on the orderbook channel
func (s *RawOrderBookSocket) SendMessage(c *Conn, msgType string, data interface{}) {
	SendMessage(c, RawOrderBookChannel, msgType, data)
}

// SendErrorMessage sends error message on orderbookchannel
func (s *RawOrderBookSocket) SendOrderMessage(c *Conn, msg string) {
	s.SendMessage(c, "ERROR", map[string]string{"Message": msg})
}

// SendInitMessage sends INIT message on orderbookchannel on subscription event
func (s *RawOrderBookSocket) SendInitMessage(c *Conn, data interface{}) {
	s.SendMessage(c, "INIT", data)
}

// SendUpdateMessage sends UPDATE message on orderbookchannel as new data is created
func (s *RawOrderBookSocket) SendUpdateMessage(c *Conn, data interface{}) {
	s.SendMessage(c, "UPDATE", data)
}

func (s *RawOrderBookSocket) SendErrorMessage(c *Conn, data interface{}) {
	s.SendMessage(c, "ERROR", data)
}
