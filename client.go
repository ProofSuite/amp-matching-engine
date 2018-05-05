package main

import "github.com/gorilla/websocket"

type Client struct {
	server     *Server
	connection *websocket.Conn
	send       chan []byte
}
