package ws

import (
	"errors"

	"github.com/gorilla/websocket"
)

type PairWs map[*websocket.Conn]bool

var pairSockets map[string]PairWs

func PairSocketCloseHandler(pair string) func(conn *websocket.Conn) {
	return func(conn *websocket.Conn) {
		PairSocketUnregisterConnection(pair, conn)
	}
}

func PairSocketUnregisterConnection(pair string, conn *websocket.Conn) {
	if pairSockets[pair][conn] {
		pairSockets[pair][conn] = false
		delete(pairSockets[pair], conn)
	}
}
func PairSocketWriteMessage(pair string, message []byte) error {
	for conn, status := range pairSockets[pair] {
		if status {
			conn.WriteMessage(1, message)
		}
	}
	return nil
}
func PairSocketRegister(pair string, conn *websocket.Conn) error {
	if conn == nil {
		return errors.New("nil not allowed in arguments as *websocket.Conn")
	} else if pairSockets == nil {
		pairSockets = make(map[string]PairWs)
		pairSockets[pair] = make(PairWs)
	} else if pairSockets[pair] == nil {
		pairSockets[pair] = make(PairWs)
	}
	if !pairSockets[pair][conn] {
		pairSockets[pair][conn] = true
	}
	return nil
}
