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

// TODO: refactor into non-global variables
var unsubscribeHandlers map[*Client][]func(*Client)
var subscriptionMutex sync.Mutex

func NewClient(c *websocket.Conn) *Client {
	subscriptionMutex.Lock()
	defer subscriptionMutex.Unlock()
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
	subscriptionMutex.Lock()
	defer subscriptionMutex.Unlock()

	for _, unsub := range unsubscribeHandlers[c] {
		go unsub(c)
	}

	c.Close()
}

func (c *Client) SendOrderErrorMessage(err error, h common.Hash) {
	p := map[string]interface{}{
		"message": err.Error(),
		"hash":    h.Hex(),
	}

	e := types.WebsocketEvent{
		Type:    "ERROR",
		Payload: p,
	}

	m := types.WebsocketMessage{
		Channel: OrderChannel,
		Event:   e,
	}

	c.mu.Lock()
	defer c.mu.Unlock()
	c.send <- m
}
