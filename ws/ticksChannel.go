package ws

import "github.com/gorilla/websocket"

var tickSubscriptions map[string]map[*websocket.Conn]bool

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
func UnsubscribeTick(channel string, conn *websocket.Conn) {
	if tickSubscriptions == nil {
		tickSubscriptions = make(map[string]map[*websocket.Conn]bool)
	}
	if tickSubscriptions[channel][conn] {
		tickSubscriptions[channel][conn] = false
		delete(tickSubscriptions[channel], conn)
	}
}
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
