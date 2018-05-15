package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/mitchellh/mapstructure"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type Server struct {
	clients    map[*ClientInterface]bool
	register   chan *ClientInterface
	unregister chan *ClientInterface
	engine     *TradingEngine
}

func NewServer() *Server {
	return &Server{
		register:   make(chan *ClientInterface),
		unregister: make(chan *ClientInterface),
		clients:    make(map[*ClientInterface]bool),
		engine:     nil,
	}
}

func (s *Server) start(pairs []Pair, done chan bool) {
	fmt.Printf("Starting server ....\n\n\n")
	s.engine = NewTradingEngine()

	for _, pair := range pairs {
		s.engine.CreateNewOrderBook(pair, done)
	}

	for {
		select {
		case client := <-s.register:
			s.clients[client] = true
		case client := <-s.unregister:
			if _, ok := s.clients[client]; ok {
				delete(s.clients, client)
			}
		}
	}
}

func (s *Server) openWebsocketConnection(w http.ResponseWriter, r *http.Request) {
	signals := make(chan *Signal)
	clientMessages := make(chan *Message)
	connection, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	fmt.Printf("Opening new connection ...\n\n")

	// defer connection.Close()
	clientInterface := &ClientInterface{server: s, connection: connection, signals: signals, clientMessages: clientMessages}
	clientInterface.server.register <- clientInterface

	// Listen to messages from server to client
	go func() {
		for {
			signal := <-clientInterface.signals
			switch signal.signalType {
			case ORDER_PLACED:
				acceptedOrder := &Order{Id: signal.orderId, Pair: signal.pair, Amount: signal.amount, Price: signal.price}
				clientInterface.sendOrderPlacedMessage(acceptedOrder)
			case ORDER_MATCHED:
				acceptedOrder := &Order{Id: signal.orderId, Pair: signal.pair, Amount: signal.amount, Price: signal.price}
				clientInterface.requestSignedDataMessage(acceptedOrder)
			case ORDER_PARTIALLY_FILLED:
				acceptedOrder := &Order{Id: signal.orderId, Pair: signal.pair, Amount: signal.amount, Price: signal.price}
				clientInterface.sendOrderPartiallyFilledMessage(acceptedOrder)
			case ORDER_FILLED:
				acceptedOrder := &Order{Id: signal.orderId, Pair: signal.pair, Amount: signal.amount, Price: signal.price}
				clientInterface.sendOrderFilledMessage(acceptedOrder)
			case ORDER_CANCELED:
				acceptedOrder := &Order{Id: signal.orderId, Pair: signal.pair, Amount: signal.amount, Price: signal.price}
				clientInterface.sendOrderCanceledMessage(acceptedOrder)
			case DONE:
			default:
				panic("Unknown action type")
			}
		}
	}()

	// Listen to messages from client to server
	go func() {
		for {
			m := <-clientInterface.clientMessages
			switch m.MessageType {
			case PLACE_ORDER:
				o := &Order{}
				payload := m.Payload.(map[string]interface{})["order"]
				mapstructure.Decode(payload, o)
				s.placeOrder(o, clientInterface)
			case CANCEL_ORDER:
				p := &CancelOrderPayload{}
				cancelOrderPayload := m.Payload
				mapstructure.Decode(cancelOrderPayload, p)
				s.cancelOrder(p)
			case SIGNED_DATA:
				p := &SignedDataPayload{}
				signedDataPayload := m.Payload
				mapstructure.Decode(signedDataPayload, p)
				s.executeOrder(p, clientInterface)
			case DONE:
			default:
				panic("Unknown message type")
			}
		}
	}()

	for {
		message := new(Message)
		err := clientInterface.connection.ReadJSON(&message)

		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("Error: %v", err)
			}
			break
		}

		clientInterface.clientMessages <- message
	}
}

func (s *Server) placeOrder(o *Order, c *ClientInterface) {
	fmt.Printf("LOG. Placing Order:\n%v\n\n", o)
	o.client = c
	if err := s.engine.AddOrder(o); err != nil {
		log.Println("Error: Failed processing order: %v", o)
	}
}

func (s *Server) cancelOrder(p *CancelOrderPayload) {
	fmt.Printf("LOG. Canceling Order:\n%v\n\n", p)
	if err := s.engine.CancelOrder(p.OrderId, p.Pair); err != nil {
		log.Println("Error: Failed canceling order")
	}
}

func (s *Server) executeOrder(p *SignedDataPayload, c *ClientInterface) {
	fmt.Printf("LOG: Executing order. Payload:\n%v\n\n", p)
	if err := s.engine.ExecuteOrder(p.OrderId, p.Pair); err != nil {
		log.Println("ERROR: Could not cancel order")
	}
}
