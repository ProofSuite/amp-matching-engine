package ws

import (
	"errors"

	"github.com/Proofsuite/amp-matching-engine/utils"
	"github.com/gorilla/websocket"
)

// PairSockets holds the map of connections subscribed to pair channels
// corresponding to the key/event they have subscribed to.
type PairSockets struct {
	connections map[string]map[*websocket.Conn]bool
}

var pairSockets *PairSockets

// GetPairSockets return singleton instance of PairSockets type struct
func GetPairSockets() *PairSockets {
	if pairSockets == nil {
		pairSockets = &PairSockets{make(map[string]map[*websocket.Conn]bool)}
	}
	return pairSockets
}

// UnsubscribeHandler returns function of type unsubscribe handler,
// it handles the unsubscription of pair in case of connection closing.
func (ps *PairSockets) UnsubscribeHandler(bt, qt string) func(conn *websocket.Conn) {
	return func(conn *websocket.Conn) {
		ps.UnregisterConnection(bt, qt, conn)
	}
}

// UnregisterConnection is used to unsubscribe the connection from listening to the key
// subscribed to. It can be called on unsubscription message from user or due to some other reason by
// system
func (ps *PairSockets) UnregisterConnection(bt, qt string, conn *websocket.Conn) {
	pair := utils.GetPairKey(bt, qt)
	if ps.connections[pair][conn] {
		ps.connections[pair][conn] = false
		delete(ps.connections[pair], conn)
	}
}

// WriteMessage streams message to all the connections subscribed to the pair
func (ps *PairSockets) WriteMessage(bt, qt string, msgType string, message interface{}) error {
	pair := utils.GetPairKey(bt, qt)

	for conn, status := range ps.connections[pair] {
		if status {
			ps.SendMessage(conn, msgType, message)
		}
	}

	return nil
}

// Register handles the registration of connection to get
// streaming data over the socker for any pair.
func (ps *PairSockets) Register(bt, qt string, conn *websocket.Conn) error {
	pair := utils.GetPairKey(bt, qt)
	if conn == nil {
		return errors.New("nil not allowed in arguments as *websocket.Conn")
	} else if ps.connections == nil {
		ps.connections = make(map[string]map[*websocket.Conn]bool)
		ps.connections[pair] = make(map[*websocket.Conn]bool)
	} else if ps.connections[pair] == nil {
		ps.connections[pair] = make(map[*websocket.Conn]bool)
	}
	if !ps.connections[pair][conn] {
		ps.connections[pair][conn] = true
	}
	return nil
}

// SendMessage is responsible for sending message on orderbook channel
func (ps *PairSockets) SendMessage(conn *websocket.Conn, msgType string, msg interface{}) {
	SendMessage(conn, OrderbookChannel, msgType, msg)
}

// SendErrorMessage is responsible for sending error messages on orderbook channel
func (ps *PairSockets) SendErrorMessage(conn *websocket.Conn, msg interface{}) {
	ps.SendMessage(conn, "ERROR", msg)
}

// SendBookMessage is responsible for sending complete order book on subscription request
func (ps *PairSockets) SendBookMessage(conn *websocket.Conn, msg interface{}) {
	ps.SendMessage(conn, "INIT", msg)
}
