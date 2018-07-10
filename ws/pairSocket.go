package ws

import (
	"errors"
	"strings"

	"github.com/gorilla/websocket"
)

type PairSockets struct {
	connections map[string]map[*websocket.Conn]bool
}

var pairSockets *PairSockets

func GetPairSockets() *PairSockets {
	if pairSockets == nil {
		pairSockets = &PairSockets{make(map[string]map[*websocket.Conn]bool)}
	}
	return pairSockets
}
func (ps *PairSockets) PairSocketCloseHandler(pair string) func(conn *websocket.Conn) {
	pair = strings.ToLower(pair)
	return func(conn *websocket.Conn) {
		ps.PairSocketUnregisterConnection(pair, conn)
	}
}

func (ps *PairSockets) PairSocketUnregisterConnection(pair string, conn *websocket.Conn) {
	pair = strings.ToLower(pair)
	if ps.connections[pair][conn] {
		ps.connections[pair][conn] = false
		delete(ps.connections[pair], conn)
	}
}
func (ps *PairSockets) PairSocketWriteMessage(pair string, message interface{}) error {
	pair = strings.ToLower(pair)
	for conn, status := range ps.connections[pair] {

		if status {
			conn.WriteJSON(message)
		}
	}
	return nil
}
func (ps *PairSockets) PairSocketRegister(pair string, conn *websocket.Conn) error {

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
