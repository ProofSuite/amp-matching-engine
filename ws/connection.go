package ws

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"sync"

	"github.com/Proofsuite/amp-matching-engine/types"
	"github.com/Proofsuite/amp-matching-engine/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/gorilla/websocket"
)

const (
	TradeChannel        = "trades"
	RawOrderBookChannel = "raw_orderbook"
	OrderBookChannel    = "orderbook"
	OrderChannel        = "orders"
	OHLCVChannel        = "ohlcv"
)

var logger = utils.Logger

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Conn struct {
	*websocket.Conn
	mu sync.Mutex
}

var connectionUnsubscribtions map[*Conn][]func(*Conn)
var socketChannels map[string]func(interface{}, *Conn)

// ConnectionEndpoint is the the handleFunc function for websocket connections
// It handles incoming websocket messages and routes the message according to
// channel parameter in channelMessage
func ConnectionEndpoint(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		logger.Error(err)
		return
	}

	conn := &Conn{c, sync.Mutex{}}
	initConnection(conn)

	go func() {
		// Recover in case of any panic in websocket. So that the app doesn't crash ===
		defer func() {
			if r := recover(); r != nil {
				err, ok := r.(error)
				if err != nil {
					logger.Error(err)
				}

				if !ok {
					logger.Error("Failed attempt at recovering websocket panic")
				}
			}
		}()

		for {
			messageType, p, err := conn.ReadMessage()
			if err != nil {
				logger.Error(err)
				conn.Close()
			}

			if messageType != 1 {
				return
			}

			msg := types.WebsocketMessage{}
			if err := json.Unmarshal(p, &msg); err != nil {
				logger.Error(err)
				SendMessage(conn, msg.Channel, "ERROR", err.Error())
				return
			}

			conn.SetCloseHandler(wsCloseHandler(conn))

			if socketChannels[msg.Channel] == nil {
				SendMessage(conn, msg.Channel, "ERROR", "INVALID_CHANNEL")
				return
			}

			go socketChannels[msg.Channel](msg.Event, conn)
		}
	}()
}

func NewConnection(c *websocket.Conn) *Conn {
	return &Conn{c, sync.Mutex{}}
}

// initConnection initializes connection in connectionUnsubscribtions map
func initConnection(conn *Conn) {
	if connectionUnsubscribtions == nil {
		connectionUnsubscribtions = make(map[*Conn][]func(*Conn))
	}

	if connectionUnsubscribtions[conn] == nil {
		connectionUnsubscribtions[conn] = make([]func(*Conn), 0)
	}
}

// RegisterChannel function needs to be called whenever the system is interested in listening to
// a new channel. A channel needs function which will handle the incoming messages for that channel.
//
// channelMessage handler function receives message from channelMessage and pointer to connection
func RegisterChannel(channel string, fn func(interface{}, *Conn)) error {
	if channel == "" {
		return errors.New("Channel can not be empty string")
	}

	if fn == nil {
		logger.Error("Handler should not be nil")
		return errors.New("Handler should not be nil")
	}

	ch := getChannelMap()
	if ch[channel] != nil {
		logger.Error("Channel already registered")
		return fmt.Errorf("Channel already registered")
	}

	ch[channel] = fn
	return nil
}

// getChannelMap returns singleton map of channels with there handler functions
func getChannelMap() map[string]func(interface{}, *Conn) {
	if socketChannels == nil {
		socketChannels = make(map[string]func(interface{}, *Conn))
	}
	return socketChannels
}

// RegisterConnectionUnsubscribeHandler needs to be called whenever a connection subscribes to
// a new channel.
// At the time of connection closing the ConnectionUnsubscribeHandler handlers associated with
// that connection are triggered.
func RegisterConnectionUnsubscribeHandler(c *Conn, fn func(*Conn)) {
	connectionUnsubscribtions[c] = append(connectionUnsubscribtions[c], fn)
}

// wsCloseHandler handles the closing of connection.
// it triggers all the UnsubscribeHandler associated with the closing
// connection in a separate go routine
func wsCloseHandler(c *Conn) func(code int, text string) error {
	return func(code int, text string) error {
		for _, unsub := range connectionUnsubscribtions[c] {
			go unsub(c)
		}
		return nil
	}
}

// SendMessage constructs the message with proper structure to be sent over websocket
func SendMessage(c *Conn, channel string, msgType string, payload interface{}, h ...common.Hash) {
	e := types.WebsocketEvent{
		Type:    msgType,
		Payload: payload,
	}

	if len(h) > 0 {
		e.Hash = h[0].Hex()
	}

	m := types.WebsocketMessage{
		Channel: channel,
		Event:   e,
	}

	c.mu.Lock()
	defer c.mu.Unlock()
	err := c.WriteJSON(m)
	if err != nil {
		logger.Error(err)
		c.Close()
	}
}
