package ws

import "github.com/gorilla/websocket"

var tickSubscriptions map[string]map[*websocket.Conn]bool

// SubscribeTick handles the subscription to ohlcv data streaming
func SubscribeTick(channel string, conn *websocket.Conn) error {
	if tickSubscriptions == nil {
		tickSubscriptions = make(map[string]map[*websocket.Conn]bool)
	}
	if tickSubscriptions[channel] == nil {
		tickSubscriptions[channel] = make(map[*websocket.Conn]bool)
	}
	tickSubscriptions[channel][conn] = true
	return nil
}

// UnsubscribeTick handles the unsubscription from ohlcv data streaming
func UnsubscribeTick(channel string, conn *websocket.Conn) {
	if tickSubscriptions == nil {
		tickSubscriptions = make(map[string]map[*websocket.Conn]bool)
	}
	if tickSubscriptions[channel][conn] {
		tickSubscriptions[channel][conn] = false
		delete(tickSubscriptions[channel], conn)
	}
}

// TickCloseHandler handles the unsubscription from ohlcv data streaming in case of connection close
func TickCloseHandler(channel string) func(conn *websocket.Conn) {
	return func(conn *websocket.Conn) {
		if tickSubscriptions == nil {
			tickSubscriptions = make(map[string]map[*websocket.Conn]bool)
		}
		if tickSubscriptions[channel][conn] {
			tickSubscriptions[channel][conn] = false
			delete(tickSubscriptions[channel], conn)
		}
	}
}

// TickBroadcast broadcasts the ohlcv data to all the subscribed connections
func TickBroadcast(channel string, msg interface{}) {
	go func() {
		for conn, isActive := range tickSubscriptions[channel] {
			var err error
			if isActive {
				err = conn.WriteJSON(msg)
			}
			if err != nil || !isActive {
				conn.Close()
			}
		}
	}()
}
