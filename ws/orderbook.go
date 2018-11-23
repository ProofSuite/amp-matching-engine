package ws

import (
	"errors"
)

var orderbook *OrderBookSocket

// OrderBookSocket holds the map of subscribtions subscribed to pair channels
// corresponding to the key/event they have subscribed to.
type OrderBookSocket struct {
	subscriptions     map[string]map[*Client]bool
	subscriptionsList map[*Client][]string
}

func NewOrderBookSocket() *OrderBookSocket {
	return &OrderBookSocket{
		subscriptions:     make(map[string]map[*Client]bool),
		subscriptionsList: make(map[*Client][]string),
	}
}

// GetOrderBookSocket return singleton instance of PairSockets type struct
func GetOrderBookSocket() *OrderBookSocket {
	if orderbook == nil {
		orderbook = NewOrderBookSocket()
	}

	return orderbook
}

// Subscribe handles the subscription of connection to get
// streaming data over the socker for any pair.
// pair := utils.GetPairKey(bt, qt)
func (s *OrderBookSocket) Subscribe(channelID string, c *Client) error {
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
func (s *OrderBookSocket) UnsubscribeHandler(channelID string) func(c *Client) {
	return func(c *Client) {
		s.UnsubscribeChannel(channelID, c)
	}
}

// Unsubscribe is used to unsubscribe the connection from listening to the key
// subscribed to. It can be called on unsubscription message from user or due to some other reason by
// system
func (s *OrderBookSocket) UnsubscribeChannel(channelID string, c *Client) {
	if s.subscriptions[channelID][c] {
		s.subscriptions[channelID][c] = false
		delete(s.subscriptions[channelID], c)
	}
}

func (s *OrderBookSocket) Unsubscribe(c *Client) {
	channelIDs := s.subscriptionsList[c]
	if channelIDs == nil {
		return
	}

	for _, id := range s.subscriptionsList[c] {
		s.UnsubscribeChannel(id, c)
	}
}

// BroadcastMessage streams message to all the subscribtions subscribed to the pair
func (s *OrderBookSocket) BroadcastMessage(channelID string, p interface{}) error {

	for c, status := range s.subscriptions[channelID] {
		if status {
			s.SendUpdateMessage(c, p)
		}
	}

	return nil
}

// SendErrorMessage sends error message on orderbookchannel
func (s *OrderBookSocket) SendErrorMessage(c *Client, data interface{}) {
	c.SendMessage(OrderBookChannel, "ERROR", data)
}

// SendInitMessage sends INIT message on orderbookchannel on subscription event
func (s *OrderBookSocket) SendInitMessage(c *Client, data interface{}) {
	c.SendMessage(OrderBookChannel, "INIT", data)
}

// SendUpdateMessage sends UPDATE message on orderbookchannel as new data is created
func (s *OrderBookSocket) SendUpdateMessage(c *Client, data interface{}) {
	c.SendMessage(OrderBookChannel, "UPDATE", data)
}
