package ws

import (
	"errors"
)

var orderBookSocket *OrderBookSocket

// OrderBookSocket holds the map of subscribtions subscribed to pair channels
// corresponding to the key/event they have subscribed to.
type OrderBookSocket struct {
	subscriptions map[string]map[*Conn]bool
}

// GetOrderBookSocket return singleton instance of PairSockets type struct
func GetOrderBookSocket() *OrderBookSocket {
	if orderBookSocket == nil {
		orderBookSocket = &OrderBookSocket{make(map[string]map[*Conn]bool)}
	}

	return orderBookSocket
}

// Subscribe handles the subscription of connection to get
// streaming data over the socker for any pair.
// pair := utils.GetPairKey(bt, qt)
func (s *OrderBookSocket) Subscribe(channelID string, conn *Conn) error {
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
func (s *OrderBookSocket) UnsubscribeHandler(channelID string) func(conn *Conn) {
	return func(conn *Conn) {
		s.Unsubscribe(channelID, conn)
	}
}

// Unsubscribe is used to unsubscribe the connection from listening to the key
// subscribed to. It can be called on unsubscription message from user or due to some other reason by
// system
func (s *OrderBookSocket) Unsubscribe(channelID string, conn *Conn) {
	if s.subscriptions[channelID][conn] {
		s.subscriptions[channelID][conn] = false
		delete(s.subscriptions[channelID], conn)
	}
}

// BroadcastMessage streams message to all the subscribtions subscribed to the pair
func (s *OrderBookSocket) BroadcastMessage(channelID string, p interface{}) error {
	for conn, status := range s.subscriptions[channelID] {
		if status {
			SendOrderBookUpdateMessage(conn, p)
		}
	}

	return nil
}

// SendOrderBookMessage sends a message on the orderbook channel
func SendOrderBookMessage(conn *Conn, msgType string, data interface{}) {
	SendMessage(conn, OrderBookChannel, msgType, data)
}

// SendOrderBookErrorMessage sends error message on orderbookchannel
func SendOrderBookErrorMessage(conn *Conn, data interface{}) {
	SendOrderBookMessage(conn, "ERROR", data)
}

// SendOrderBookInitMessage sends INIT message on orderbookchannel on subscription event
func SendOrderBookInitMessage(conn *Conn, data interface{}) {
	SendOrderBookMessage(conn, "INIT", data)
}

// SendOrderBookUpdateMessage sends UPDATE message on orderbookchannel as new data is created
func SendOrderBookUpdateMessage(conn *Conn, data interface{}) {
	SendOrderBookMessage(conn, "UPDATE", data)
}
