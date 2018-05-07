package main

import (
	"fmt"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type AddOrderMessage struct {
	order  *Order
	client *Client
}

type Server struct {
	clients         map[*Client]bool
	addOrderMessage chan *AddOrderMessage
	register        chan *Client
	unregister      chan *Client
	engine          *TradingEngine
}

func NewServer(engine *TradingEngine) *Server {
	return &Server{
		addOrderMessage: make(chan *AddOrderMessage),
		register:        make(chan *Client),
		unregister:      make(chan *Client),
		clients:         make(map[*Client]bool),
		engine:          engine,
	}
}

func (s *Server) start() {
	for {
		select {
		case client := <-s.register:
			s.clients[client] = true
		case client := <-s.unregister:
			if _, ok := s.clients[client]; ok {
				delete(s.clients, client)
			}
		case message := <-s.addOrderMessage:
			order := message.order
			fmt.Printf("Order received: %v", order)
		}
	}
}
