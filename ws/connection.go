package ws

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/Proofsuite/amp-matching-engine/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/gorilla/websocket"
)

// websocket channel's string
const (
	TradeChannel     = "trades"
	OrderBookChannel = "order_book"
	OrderChannel     = "orders"
	OHLCVChannel     = "ohlcv"
)

// gorilla websocket upgrader instance with configuration
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var connectionUnsubscribtions map[*websocket.Conn][]func(*websocket.Conn)
var socketChannels map[string]func(interface{}, *websocket.Conn)

// ConnectionEndpoint is the the handleFunc function for websocket connections
// It handles incoming websocket messages and routes the message according to
// channel parameter in channelMessage
func ConnectionEndpoint(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("==>" + err.Error())
		return
	}

	initConnection(conn)
	go func() {
		// Recover in case of any panic in websocket. So that the app doesn't crash ===
		defer func() {
			if r := recover(); r != nil {
				var ok bool
				err, ok = r.(error)
				if !ok {
					err = fmt.Errorf("Panic in websocket: %v", r)
				}
				log.Fatal(err)
			}
		}()

		for {
			messageType, p, err := conn.ReadMessage()
			if err != nil {
				conn.Close()
			}

			if messageType != 1 {
				return
			}

			msg := types.WebSocketMessage{}
			if err := json.Unmarshal(p, &msg); err != nil {
				log.Println("unmarshal to channelMessage <==>" + err.Error())
				SendMessage(conn, msg.Channel, "ERROR", err.Error())
				return
			}

			conn.SetCloseHandler(wsCloseHandler(conn))
			if socketChannels[msg.Channel] != nil {
				go socketChannels[msg.Channel](msg.Payload, conn)
			} else {
				SendMessage(conn, msg.Channel, "ERROR", "INVALID_CHANNEL")
			}
		}
	}()
}

// initConnection initializes connection in connectionUnsubscribtions map
func initConnection(conn *websocket.Conn) {
	if connectionUnsubscribtions == nil {
		connectionUnsubscribtions = make(map[*websocket.Conn][]func(*websocket.Conn))
	}

	if connectionUnsubscribtions[conn] == nil {
		connectionUnsubscribtions[conn] = make([]func(*websocket.Conn), 0)
	}
}

// RegisterChannel function needs to be called whenever the system is interested in listening to
// a new channel. A channel needs function which will handle the incoming messages for that channel.
//
// channelMessage handler function receives message from channelMessage and pointer to connection
func RegisterChannel(channel string, fn func(interface{}, *websocket.Conn)) error {
	if channel == "" {
		return errors.New("Channel can not be empty string")
	}

	if fn == nil {
		return errors.New("fn can not be nil")
	}

	ch := getChannelMap()
	if ch[channel] != nil {
		return fmt.Errorf("channel %s already registered", channel)
	}

	ch[channel] = fn
	return nil
}

// getChannelMap returns singleton map of channels with there handler functions
func getChannelMap() map[string]func(interface{}, *websocket.Conn) {
	if socketChannels == nil {
		socketChannels = make(map[string]func(interface{}, *websocket.Conn))
	}
	return socketChannels
}

// RegisterConnectionUnsubscribeHandler needs to be called whenever a connection subscribes to
// a new channel.
// At the time of connection closing the ConnectionUnsubscribeHandler handlers associated with
// that connection are triggered.
func RegisterConnectionUnsubscribeHandler(conn *websocket.Conn, fn func(*websocket.Conn)) {
	connectionUnsubscribtions[conn] = append(connectionUnsubscribtions[conn], fn)
}

// wsCloseHandler handles the closing of connection.
// it triggers all the UnsubscribeHandler associated with the closing
// connection in a separate go routine
func wsCloseHandler(conn *websocket.Conn) func(code int, text string) error {
	return func(code int, text string) error {
		for _, unsub := range connectionUnsubscribtions[conn] {
			go unsub(conn)
		}
		return nil
	}
}

// SendMessage constructs the message with proper structure to be sent over websocket
func SendMessage(conn *websocket.Conn, channel string, msgType string, data interface{}, hash ...common.Hash) {
	payload := types.WebSocketPayload{
		Type: msgType,
		Data: data,
	}

	if len(hash) > 0 {
		payload.Hash = hash[0].Hex()
	}

	message := types.WebSocketMessage{
		Channel: channel,
		Payload: payload,
	}
	err := conn.WriteJSON(message)
	if err != nil {
		conn.Close()
	}
}
