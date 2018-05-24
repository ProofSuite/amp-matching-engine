package dex

import (
	"fmt"
	"log"

	"github.com/gorilla/websocket"
)

// Socket acts as a hub that handles messages from the client application and responses
// from the server matching engine.
type Socket struct {
	server      *Server
	connection  *websocket.Conn
	actions     chan *Action
	messagesIn  chan *Message
	messagesOut chan *Message
	events      chan *Event
}

// listenToMessagesIn reads incoming messages from the websocket connection
// and sends these messages into the messageIn channel
func (s *Socket) listenToMessagesIn() {
	for {
		message := new(Message)
		err := s.connection.ReadJSON(&message)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("Error: %v", err)
			}
			break
		}
		s.messagesIn <- message
	}
}

// handleMessagesIN listens on the messageIn channel and routes them to the appropriate
// handler based on their MessageType
func (s *Socket) handleMessagesIn() {
	for {
		m := <-s.messagesIn
		switch m.MessageType {
		case PLACE_ORDER:
			s.placeOrder(m.Payload)
		case CANCEL_ORDER:
			s.cancelOrder(&m.Payload)
		case SIGNED_DATA:
			s.executeOrder(m.Payload)
		case DONE:
		default:
			panic("Unknown message type")
		}
	}
}

// handleMessagesOut listens on the event channel (events sent from the matching engine) and routes them
// to the appropriate handler based on their event type
func (s *Socket) handleMessagesOut() {
	for {
		e := <-s.events
		switch e.eventType {
		case ORDER_PLACED:
			order := e.payload.(*Order)
			s.sendOrderPlacedMessage(order)
		case ORDER_PARTIALLY_FILLED:
			order := e.payload.(*TradePayload)
			s.sendOrderPartiallyFilledMessage(order)
		case ORDER_FILLED:
			payload := e.payload.(*TradePayload)
			s.sendOrderFilledMessage(payload)
		case ORDER_CANCELED:
			order := e.payload.(*Order)
			s.sendOrderCanceledMessage(order)
		case DONE:
		default:
			panic("Unknown action type")
		}
	}
}

// placeOrder decodes orders and then passes it to the engine object
func (s *Socket) placeOrder(p Payload) {
	payload := p.(map[string]interface{})["order"].(map[string]interface{})
	o := &Order{}
	o.DecodeOrder(payload)

	o.events = s.events
	if err := s.server.engine.AddOrder(o); err != nil {
		log.Printf("Error: Failed processing order: %v", err)
	}
}

// cancelOrder decodes the message payload and then passes it to the engine object
func (s *Socket) cancelOrder(p Payload) {
	decoded := NewOrderCancelPayload()
	decoded.DecodeOrderCancelPayload(p)

	oc := decoded.OrderCancel
	// fmt.Printf("\nLOG. Canceling Order:\n%v\n\n", payload)
	if err := s.server.engine.CancelOrder(oc); err != nil {
		log.Println("Error: Failing canceling order")
	}
}

// executeOrder decodes the message payload before passing it to the transaction handler
func (s *Socket) executeOrder(p Payload) {
	payload := p.(map[string]interface{})["trade"].(map[string]interface{})
	t := &Trade{}
	t.DecodeTrade(payload)

	fmt.Printf("\nLOG: Executing order. Payload:\n%v\n\n", t)
}

// sendOrderPlacedMessage creates and ORDER_PLACED messages and writes it into the websocket connection
func (s *Socket) sendOrderPlacedMessage(o *Order) error {
	p := &OrderPayload{Order: o}
	m := &Message{MessageType: ORDER_PLACED, Payload: p}

	if err := s.connection.WriteJSON(&m); err != nil {
		return err
	}

	// fmt.Printf("\nLOG. Sending Order Placed Message:\n%v\n\n", o)
	return nil
}

// sendOrderFilledMessage creates an ORDER_FILLED messages and writes it into the websocket connection
func (s *Socket) sendOrderFilledMessage(p *TradePayload) error {
	m := &Message{MessageType: ORDER_FILLED, Payload: p}
	if err := s.connection.WriteJSON(&m); err != nil {
		return err
	}

	// fmt.Printf("\nLOG. Sending Order Filled Message:\n%v\n\n", p)
	return nil
}

// sendOrderPartiallyFilledMessage creates and ORDER_PARTIALLY_FILLED message and writes it into the websocket connection
func (s *Socket) sendOrderPartiallyFilledMessage(p *TradePayload) error {
	fmt.Printf("Send order partially filled message")
	m := &Message{MessageType: ORDER_PARTIALLY_FILLED, Payload: p}
	if err := s.connection.WriteJSON(&m); err != nil {
		return err
	}

	// fmt.Printf("\nLOG. Sending Partially Filled Message:\n%v\n\n", p)
	return nil
}

// sendOrderCanceledMessage creates an ORDER_CANCELED message and writes it into the websocket connection
func (s *Socket) sendOrderCanceledMessage(o *Order) error {
	fmt.Printf("Send order canceled message")
	p := &OrderPayload{Order: o}
	m := &Message{MessageType: ORDER_CANCELED, Payload: p}
	if err := s.connection.WriteJSON(&m); err != nil {
		return err
	}

	// fmt.Printf("\nLOG: Sending Order Canceled Message:\n%v\n\n", o)
	return nil
}
