package ws

import (
	"errors"

	"github.com/gorilla/websocket"
)

type PairWs map[*websocket.Conn]bool

var PairSockets map[string]PairWs

func PairSocketCloseHandler(pair string, conn *websocket.Conn) func(code int, text string) error {
	return func(code int, text string) error {
		if PairSockets[pair][conn] {
			PairSockets[pair][conn] = false
			delete(PairSockets[pair], conn)
		}
		return nil
	}
}

func PairSocketWriteMessage(pair string, message []byte) error {
	for conn, status := range PairSockets[pair] {
		if status {
			conn.WriteMessage(1, message)
		}
	}
	return nil
}
func PairSocketRegister(pair string, conn *websocket.Conn) error {
	if conn == nil {
		return errors.New("nil not allowed in arguments as *websocket.Conn")
	} else if PairSockets == nil {
		PairSockets = make(map[string]PairWs)
		PairSockets[pair] = make(PairWs)
	} else if PairSockets[pair] == nil {
		PairSockets[pair] = make(PairWs)
	}
	if !PairSockets[pair][conn] {
		PairSockets[pair][conn] = true
	}
	return nil
}
