package ws

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type channelMessage struct {
	Channel string       `json:"channel"`
	Message *interface{} `json:"message"`
}

var connectionUnsubscribtions map[*websocket.Conn][]func(*websocket.Conn)
var socketChannels map[string]func(*interface{}, *websocket.Conn)

func ConnectionEndpoint(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("==>" + err.Error())
		return
	}
	for {
		_, p, err := conn.ReadMessage()
		if err != nil {
			conn.Close()
		}
		initConnection(conn)
		var msg *channelMessage
		if err := json.Unmarshal(p, &msg); err != nil {
			log.Println("unmarshal to channelMessage <==>" + err.Error())
			conn.WriteJSON(map[string]interface{}{"channelMessage": err.Error()})
		}
		conn.SetCloseHandler(wsCloseHandler(conn))
		if socketChannels[msg.Channel] != nil {
			go socketChannels[msg.Channel](msg.Message, conn)
		} else {
			conn.WriteJSON(map[string]interface{}{"channel": "INVALID_CHANNEL"})
		}
	}
}
func initConnection(conn *websocket.Conn) {
	if connectionUnsubscribtions == nil {
		connectionUnsubscribtions = make(map[*websocket.Conn][]func(*websocket.Conn))
	}
	if connectionUnsubscribtions[conn] == nil {
		connectionUnsubscribtions[conn] = make([]func(*websocket.Conn), 0)
	}
}
func RegisterChannel(channel string, fn func(*interface{}, *websocket.Conn)) error {
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

func getChannelMap() map[string]func(*interface{}, *websocket.Conn) {
	if socketChannels == nil {
		socketChannels = make(map[string]func(*interface{}, *websocket.Conn))
	}
	return socketChannels
}

func RegisterConnectionUnsubscribeHandler(conn *websocket.Conn, fn func(*websocket.Conn)) {
	connectionUnsubscribtions[conn] = append(connectionUnsubscribtions[conn], fn)
}

func wsCloseHandler(conn *websocket.Conn) func(code int, text string) error {
	return func(code int, text string) error {
		for _, unsub := range connectionUnsubscribtions[conn] {
			go unsub(conn)
		}
		return nil
	}
}
