package ws

import (
	"github.com/Proofsuite/amp-matching-engine/types"
	"github.com/gorilla/websocket"
)

var tradeSocket *TradeSocket

// TradeSocket holds the map of connections subscribed to pair channels
// corresponding to the key/event they have subscribed to.
type TradeSocket struct {
	subscriptions map[string]map[*websocket.Conn]bool
}

func GetTradeSocket() *TradeSocket {
	if tradeSocket == nil {
		tradeSocket = &TradeSocket{make(map[string]map[*websocket.Conn]bool)}
	}

	return tradeSocket
}

// Subscribe registers a new websocket connections to the trade channel updates
func (s *TradeSocket) Subscribe(channelId string, conn *websocket.Conn) error {
	if s.subscriptions[channelId] == nil {
		s.subscriptions[channelId] = make(map[*websocket.Conn]bool)
	}

	s.subscriptions[channelId][conn] = true
	return nil
}

// Unsubscribe removes a websocket connection from the trade channel updates
func (s *TradeSocket) Unsubscribe(channelId string, conn *websocket.Conn) {
	if s.subscriptions[channelId][conn] {
		s.subscriptions[channelId][conn] = false
		delete(s.subscriptions[channelId], conn)
	}
}

// TradeUnSubscribeHandler unsubscribes a connection from a certain trade channel id
func (s *TradeSocket) UnsubscribeHandler(channelId string) func(conn *websocket.Conn) {
	return func(conn *websocket.Conn) {
		s.Unsubscribe(channelId, conn)
	}
}

func (s *TradeSocket) BroadcastMessage(channelId string, msgType string, p *types.WebSocketPayload) {
	go func() {
		for conn, active := range tradeSocket.subscriptions[channelId] {
			if active {
				SendTradeMessage(conn, msgType, p)
			}
		}
	}()
}

// SendTradeMesage sends a websocket message on the trade channel
func SendTradeMessage(conn *websocket.Conn, msgType string, p interface{}) {
	SendMessage(conn, TradeChannel, msgType, p)
}

// SendTradeErrorMessage sends an error message on the trade channel
func SendTradeErrorMessage(conn *websocket.Conn, p interface{}) {
	SendTradeMessage(conn, "ERROR", p)
}

// SendTradeTradessMessage is responsible for sending message on trade ohlcv channel at subscription
func SendTradeInitMessage(conn *websocket.Conn, p interface{}) {
	SendMessage(conn, TradeChannel, "INIT", p)
}

// TradeSendTradesMessage is responsible for sending message on trade ohlcv channel at subscription
func SendTradeUpdateMessage(conn *websocket.Conn, p interface{}) {
	SendMessage(conn, TradeChannel, "UPDATE", p)
}

// // UnsubscribeTrades unsubscribes a websocket connection from trades streaming
// func (sadf) UnsubscribeTrades(channelId string, conn *websocket.Conn) {
// 	if tradeSocket == nil {
// 		tradeSocket = &TradeSocket{
// 			subscriptions: make(map[string]map[*websocket.Conn]bool)
// 		}
// 	}

// 	if tradeSocket[channelId][conn] {
// 		tradeSocket[channelId][conn] = false
// 		delete(tradeSocket[channelId], conn)
// 	}
// }

// // TradesCloseHandler handles the unsubscription from ohlcv data streaming in case of connection close
// func TradeUnsubscribeHandler(channelId string) func(conn *websocket.Conn) {
// 	return func(conn *websocket.Conn) {
// 		if tradeSocket == nil {
// 			tradeSocket = make(map[string]map[*websocket.Conn]bool)
// 		}

// 		if tradeSocket[channelId][conn] {
// 			tradeSocket[channelId][conn] = false
// 			delete(tradeSocket[channelId], conn)
// 		}
// 	}
// }

// // BroadcastTrades send an OHLCV udpate message to all subscribed sockets
// func BroadcastTrades(channel string, p *types.WebSocketPayload) {
// 	go func() {
// 		for conn, isActive := range tradeSocket[channel] {
// 			if isActive {
// 				SendTradeUpdateMessage(conn, p)
// 			}
// 		}
// 	}()
// }

// //
