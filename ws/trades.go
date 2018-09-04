package ws

var tradeSocket *TradeSocket

// TradeSocket holds the map of connections subscribed to pair channels
// corresponding to the key/event they have subscribed to.
type TradeSocket struct {
	subscriptions map[string]map[*Conn]bool
}

func GetTradeSocket() *TradeSocket {
	if tradeSocket == nil {
		tradeSocket = &TradeSocket{make(map[string]map[*Conn]bool)}
	}

	return tradeSocket
}

// Subscribe registers a new websocket connections to the trade channel updates
func (s *TradeSocket) Subscribe(channelID string, conn *Conn) error {
	if s.subscriptions[channelID] == nil {
		s.subscriptions[channelID] = make(map[*Conn]bool)
	}

	s.subscriptions[channelID][conn] = true
	return nil
}

// Unsubscribe removes a websocket connection from the trade channel updates
func (s *TradeSocket) Unsubscribe(channelID string, conn *Conn) {
	if s.subscriptions[channelID][conn] {
		s.subscriptions[channelID][conn] = false
		delete(s.subscriptions[channelID], conn)
	}
}

// UnsubscribeHandler unsubscribes a connection from a certain trade channel id
func (s *TradeSocket) UnsubscribeHandler(channelID string) func(conn *Conn) {
	return func(conn *Conn) {
		s.Unsubscribe(channelID, conn)
	}
}

// BroadcastMessage broadcasts trade message to all subscribed sockets
func (s *TradeSocket) BroadcastMessage(channelID string, p interface{}) {
	go func() {
		for conn, active := range tradeSocket.subscriptions[channelID] {
			if active {
				SendTradeUpdateMessage(conn, p)
			}
		}
	}()
}

// SendTradeMessage sends a websocket message on the trade channel
func SendTradeMessage(conn *Conn, msgType string, p interface{}) {
	SendMessage(conn, TradeChannel, msgType, p)
}

// SendTradeErrorMessage sends an error message on the trade channel
func SendTradeErrorMessage(conn *Conn, p interface{}) {
	SendTradeMessage(conn, "ERROR", p)
}

// SendTradeInitMessage is responsible for sending message on trade ohlcv channel at subscription
func SendTradeInitMessage(conn *Conn, p interface{}) {
	SendMessage(conn, TradeChannel, "INIT", p)
}

// SendTradeUpdateMessage is responsible for sending message on trade ohlcv channel at subscription
func SendTradeUpdateMessage(conn *Conn, p interface{}) {
	SendMessage(conn, TradeChannel, "UPDATE", p)
}
