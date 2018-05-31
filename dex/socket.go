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
		fmt.Printf("Receiving message: %v", m)
		switch m.MessageType {
		case PLACE_ORDER:
			s.placeOrder(m.Payload)
		case CANCEL_ORDER:
			s.cancelOrder(m.Payload)
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
		fmt.Printf("Sending message: %v", e)
		switch e.eventType {
		case ORDER_PLACED:
			order := e.payload.(*Order)
			s.sendOrderPlaced(order)
		case ORDER_PARTIALLY_FILLED:
			order := e.payload.(*TradePayload)
			s.sendOrderPartiallyFilled(order)
		case ORDER_FILLED:
			payload := e.payload.(*TradePayload)
			s.sendOrderFilled(payload)
		case ORDER_CANCELED:
			order := e.payload.(*Order)
			s.sendOrderCanceled(order)
		case ORDER_EXECUTED:
			log.Printf("Trade has been executed")
			p := e.payload.(*OrderExecutedPayload)
			s.sendOrderExecuted(p)
		case ORDER_TX_SUCCESS:
			log.Printf("Order TX success")
			p := e.payload.(*TxSuccessPayload)
			s.sendOrderTxSuccess(p)
		case ORDER_TX_ERROR:
			log.Printf("Order TX error")
			p := e.payload.(*TxErrorPayload)
			s.sendOrderTxError(p)
		case TRADE_EXECUTED:
			p := e.payload.(*TradeExecutedPayload)
			s.sendTradeExecuted(p)
		case TRADE_TX_SUCCESS:
			log.Printf("Trade Tx Success")
			p := e.payload.(*TxSuccessPayload)
			s.sendTradeTxSuccess(p)
		case TRADE_TX_ERROR:
			log.Printf("Trade Tx Error")
			p := e.payload.(*TxErrorPayload)
			s.sendTradeTxError(p)
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
	o.Decode(payload)

	o.events = s.events
	if err := s.server.engine.AddOrder(o); err != nil {
		log.Printf("Error: Failed processing order: %v", err)
	}
}

// cancelOrder decodes the message payload and then passes it to the engine object
func (s *Socket) cancelOrder(p Payload) {
	ocp := NewOrderCancelPayload()

	if err := ocp.DecodeOrderCancelPayload(p); err != nil {
		log.Printf("Error: %v", err)
	}

	oc := ocp.OrderCancel

	if err := s.server.engine.CancelOrder(oc); err != nil {
		log.Printf("Error: %v", err)
	}
}

// executeOrder decodes the message payload before passing it to the transaction handler
func (s *Socket) executeOrder(p Payload) {
	log.Printf("Receiving execute order message: %v", p)

	tp := NewTradePayload()

	if err := tp.DecodeTradePayload(p); err != nil {
		log.Printf("Error: %v", err)
	}

	t := tp.Trade
	t.events = s.events

	err := s.server.engine.ExecuteOrder(t)
	if err != nil {
		log.Printf("Error: %v", err)
	}
}

// sendOrderPlacedMessage creates and ORDER_PLACED messages and writes it into the websocket connection
func (s *Socket) sendOrderPlaced(o *Order) error {
	p := &OrderPayload{Order: o}
	m := &Message{MessageType: ORDER_PLACED, Payload: p}

	if err := s.connection.WriteJSON(&m); err != nil {
		return err
	}

	return nil
}

// sendOrderFilledMessage creates an ORDER_FILLED messages and writes it into the websocket connection
func (s *Socket) sendOrderFilled(p *TradePayload) error {
	m := &Message{MessageType: ORDER_FILLED, Payload: p}
	if err := s.connection.WriteJSON(&m); err != nil {
		return err
	}

	return nil
}

// sendOrderPartiallyFilledMessage creates and ORDER_PARTIALLY_FILLED message and writes it into the websocket connection
func (s *Socket) sendOrderPartiallyFilled(p *TradePayload) error {
	fmt.Printf("Send order partially filled message")
	m := &Message{MessageType: ORDER_PARTIALLY_FILLED, Payload: p}
	if err := s.connection.WriteJSON(&m); err != nil {
		return err
	}

	return nil
}

// sendOrderCanceledMessage creates an ORDER_CANCELED message and writes it into the websocket connection
func (s *Socket) sendOrderCanceled(o *Order) error {
	fmt.Printf("Send order canceled message")
	p := &OrderPayload{Order: o}
	m := &Message{MessageType: ORDER_CANCELED, Payload: p}
	if err := s.connection.WriteJSON(&m); err != nil {
		return err
	}

	return nil
}

// sendOrderExecuted creates an ORDER_EXECUTED messages and writes it into the websocket connection
func (s *Socket) sendOrderExecuted(p *OrderExecutedPayload) error {
	m := &Message{MessageType: ORDER_EXECUTED, Payload: p}
	if err := s.connection.WriteJSON(&m); err != nil {
		return err
	}

	return nil
}

// sendTradeExecuted creates an TRADE_EXECUTED messages and writes it into the websocket connection
func (s *Socket) sendTradeExecuted(t *TradeExecutedPayload) error {
	m := &Message{MessageType: TRADE_EXECUTED, Payload: t}

	if err := s.connection.WriteJSON(&m); err != nil {
		return err
	}

	return nil
}

// sendOrderTxSuccess creates an ORDER_TX_SUCCESS message and writes it into the websocket connection
func (s *Socket) sendOrderTxSuccess(p *TxSuccessPayload) error {
	m := &Message{MessageType: ORDER_TX_SUCCESS, Payload: p}

	if err := s.connection.WriteJSON(&m); err != nil {
		return err
	}

	return nil
}

// sendOrderTxError creates an ORDER_TX_ERROR message and writes it into the websocket connection
func (s *Socket) sendOrderTxError(p *TxErrorPayload) error {
	m := &Message{MessageType: ORDER_TX_ERROR, Payload: p}

	if err := s.connection.WriteJSON(&m); err != nil {
		return err
	}

	return nil
}

// sendTradeTxSuccess creates a TRADE_TX_SUCCESS message and writes it into the websocket connection
func (s *Socket) sendTradeTxSuccess(p *TxSuccessPayload) error {
	m := &Message{MessageType: TRADE_TX_SUCCESS, Payload: p}

	if err := s.connection.WriteJSON(&m); err != nil {
		return err
	}

	return nil
}

// sendTradeTxError creates a TRADE_TX_ERROR message and writes it into the websocket connection
func (s *Socket) sendTradeTxError(p *TxErrorPayload) error {
	m := &Message{MessageType: TRADE_TX_ERROR, Payload: p}

	if err := s.connection.WriteJSON(&m); err != nil {
		return err
	}

	return nil
}
