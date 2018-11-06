package ws

import (
	"sync"

	"github.com/Proofsuite/amp-matching-engine/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/gorilla/websocket"
)

type Client struct {
	*websocket.Conn
	mu   sync.Mutex
	send chan types.WebsocketMessage
}

var unsubscribeHandlers map[*Client][]func(*Client)

func NewClient(c *websocket.Conn) *Client {
	conn := &Client{Conn: c, mu: sync.Mutex{}, send: make(chan types.WebsocketMessage)}

	if unsubscribeHandlers == nil {
		unsubscribeHandlers = make(map[*Client][]func(*Client))
	}

	if unsubscribeHandlers[conn] == nil {
		unsubscribeHandlers[conn] = make([]func(*Client), 0)
	}

	return conn
}

// SendMessage constructs the message with proper structure to be sent over websocket
func (c *Client) SendMessage(channel string, msgType string, payload interface{}, h ...common.Hash) {
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
	c.send <- m
}

func (c *Client) closeConnection() {
	for _, unsub := range unsubscribeHandlers[c] {
		go unsub(c)
	}

	c.Close()
}
