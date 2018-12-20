package ws

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/Proofsuite/amp-matching-engine/types"
	"github.com/gorilla/websocket"
)

const (
	writeWait  = 60 * time.Second
	pongWait   = 60 * time.Second
	pingPeriod = (pongWait * 9) / 10
)

var logger = NewWebsocketLogger()

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// ConnectionEndpoint is the the handleFunc function for websocket connections
// It handles incoming websocket messages and routes the message according to
// channel parameter in channelMessage
func ConnectionEndpoint(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		logger.Error(err)
		return
	}

	c := NewClient(conn)
	c.SetCloseHandler(closeHandler(c))

	go readHandler(c)
	go writeHandler(c)
}

func readHandler(c *Client) {
	defer func() {
		logger.Info("Closing connection")
		c.closeConnection()
	}()

	c.SetReadDeadline(time.Now().Add(pongWait))
	c.SetPongHandler(func(string) error {
		c.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		msgType, payload, err := c.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				logger.Error(err)
			}

			return
		}

		if msgType != 1 {
			return
		}

		msg := types.WebsocketMessage{}
		if err := json.Unmarshal(payload, &msg); err != nil {
			logger.Error(err)
			c.SendMessage(msg.Channel, "ERROR", err.Error())
			return
		}

		logger.LogMessageIn(&msg)
		logger.Infof("%v", msg.String())

		if socketChannels[msg.Channel] == nil {
			c.SendMessage(msg.Channel, "ERROR", "INVALID_CHANNEL")
			return
		}

		go socketChannels[msg.Channel](msg.Event, c)
	}
}

func writeHandler(c *Client) {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		logger.Info("Closing connection")
		ticker.Stop()
		c.closeConnection()
	}()

	for {
		select {
		case <-ticker.C:
			c.SetWriteDeadline(time.Now().Add(writeWait))
			err := c.WriteMessage(websocket.PingMessage, nil)
			if err != nil {
				logger.Error(err)
				return
			}

		case m, ok := <-c.send:
			c.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				c.WriteMessage(websocket.CloseMessage, []byte{})
			}

			logger.LogMessageOut(&m)
			logger.Infof("%v", m.String())

			err := c.WriteJSON(m)
			if err != nil {
				logger.Error(err)
				return
			}
		}
	}
}

func closeHandler(c *Client) func(code int, text string) error {
	return func(code int, text string) error {
		c.closeConnection()
		return nil
	}
}

// RegisterConnectionUnsubscribeHandler needs to be called whenever a connection subscribes to
// a new channel.
// At the time of connection closing the ConnectionUnsubscribeHandler handlers associated with
// that connection are triggered.
func RegisterConnectionUnsubscribeHandler(c *Client, fn func(*Client)) {
	subscriptionMutex.Lock()
	defer subscriptionMutex.Unlock()

	logger.Info("Registering a new unsubscribe handler")
	unsubscribeHandlers[c] = append(unsubscribeHandlers[c], fn)
}
