package ws

import "errors"

var tradeSocket *TradeSocket

// TradeSocket holds the map of connections subscribed to pair channels
// corresponding to the key/event they have subscribed to.
type TradeSocket struct {
	subscriptions     map[string]map[*Client]bool
	subscriptionsList map[*Client][]string
}

func NewTradeSocket() *TradeSocket {
	return &TradeSocket{
		subscriptions:     make(map[string]map[*Client]bool),
		subscriptionsList: make(map[*Client][]string),
	}
}

func GetTradeSocket() *TradeSocket {
	if tradeSocket == nil {
		tradeSocket = NewTradeSocket()
	}

	return tradeSocket
}

// Subscribe registers a new websocket connections to the trade channel updates
func (s *TradeSocket) Subscribe(channelID string, c *Client) error {
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

// UnsubscribeHandler unsubscribes a connection from a certain trade channel id
func (s *TradeSocket) UnsubscribeChannelHandler(channelID string) func(c *Client) {
	return func(c *Client) {
		s.UnsubscribeChannel(channelID, c)
	}
}

func (s *TradeSocket) UnsubscribeHandler() func(c *Client) {
	return func(c *Client) {
		s.Unsubscribe(c)
	}
}

// Unsubscribe removes a websocket connection from the trade channel updates
func (s *TradeSocket) UnsubscribeChannel(channelID string, c *Client) {
	if s.subscriptions[channelID][c] {
		s.subscriptions[channelID][c] = false
		delete(s.subscriptions[channelID], c)
	}
}

func (s *TradeSocket) Unsubscribe(c *Client) {
	channelIDs := s.subscriptionsList[c]
	if channelIDs == nil {
		return
	}

	for _, id := range s.subscriptionsList[c] {
		s.UnsubscribeChannel(id, c)
	}
}

// BroadcastMessage broadcasts trade message to all subscribed sockets
func (s *TradeSocket) BroadcastMessage(channelID string, p interface{}) {
	go func() {
		for conn, active := range tradeSocket.subscriptions[channelID] {
			if active {
				s.SendUpdateMessage(conn, p)
			}
		}
	}()
}

// SendMessage sends a websocket message on the trade channel
func (s *TradeSocket) SendMessage(c *Client, msgType string, p interface{}) {
	c.SendMessage(TradeChannel, msgType, p)
}

// SendErrorMessage sends an error message on the trade channel
func (s *TradeSocket) SendErrorMessage(c *Client, p interface{}) {
	c.SendMessage(TradeChannel, "ERROR", p)
}

// SendInitMessage is responsible for sending message on trade ohlcv channel at subscription
func (s *TradeSocket) SendInitMessage(c *Client, p interface{}) {
	c.SendMessage(TradeChannel, "INIT", p)
}

// SendUpdateMessage is responsible for sending message on trade ohlcv channel at subscription
func (s *TradeSocket) SendUpdateMessage(c *Client, p interface{}) {
	c.SendMessage(TradeChannel, "UPDATE", p)
}
