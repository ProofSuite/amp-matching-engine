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
			if isActive {
				TradeSendTickMessage(conn, msg)
			}
		}
	}()
}

// TradeSendMessage is responsible for sending message on trade ohlcv channel
func TradeSendMessage(conn *websocket.Conn, msgType string, msg interface{}) {
	SendMessage(conn, TradeChannel, msgType, msg)
}

// TradeSendErrorMessage is responsible for sending error messages on trade channel
func TradeSendErrorMessage(conn *websocket.Conn, msg interface{}) {
	TradeSendMessage(conn, "Error", msg)
}

// TradeSendTicksMessage is responsible for sending message on trade ohlcv channel at subscription
func TradeSendTicksMessage(conn *websocket.Conn, msg interface{}) {
	SendMessage(conn, TradeChannel, "trade_ticks", msg)
}

// TradeSendTickMessage is responsible for sending message on trade ohlcv channel at subscription
func TradeSendTickMessage(conn *websocket.Conn, msg interface{}) {
	SendMessage(conn, TradeChannel, "trade_tick", msg)
}
