package main

import (
	"fmt"
	"log"

	"github.com/gorilla/websocket"
)

type ClientInterface struct {
	server         *Server
	connection     *websocket.Conn
	send           chan []byte
	actions        chan *Action
	signals        chan *Signal
	clientMessages chan *Message
}

func (c *ClientInterface) sendOrderPlacedMessage(o *Order) {
	p := &OrderPayload{Order: o}
	m := &Message{MessageType: "ORDER_PLACED", Payload: p}
	if err := c.connection.WriteJSON(&m); err != nil {
		log.Println("Write:", err)
		return
	}

	fmt.Printf("LOG. Sending Order Placed Message:\n")
	fmt.Printf("%v\n\n", o)
}

func (c *ClientInterface) requestSignedDataMessage(o *Order) {
	p := &RequestSignedDataPayload{Orders: []*Order{o}}
	m := &Message{MessageType: "REQUEST_SIGNED_DATA", Payload: p}
	fmt.Printf("Requesting Signatures")
	if err := c.connection.WriteJSON(&m); err != nil {
		log.Println("Write:", err)
		return
	}
}

func (c *ClientInterface) sendOrderFilledMessage(o *Order) {
	fmt.Printf("Send order filled message")
	p := &OrderPayload{Order: o}
	m := &Message{MessageType: ORDER_FILLED, Payload: p}
	if err := c.connection.WriteJSON(&m); err != nil {
		log.Println("Write:", err)
		return
	}
}

func (c *ClientInterface) sendOrderPartiallyFilledMessage(o *Order) {
	fmt.Printf("Send order partially filled message")
	p := &OrderPayload{Order: o}
	m := &Message{MessageType: ORDER_PARTIALLY_FILLED, Payload: p}
	if err := c.connection.WriteJSON(&m); err != nil {
		log.Println("Write:", err)
		return
	}
}

func (c *ClientInterface) sendOrderCanceledMessage(o *Order) {
	fmt.Printf("Send order canceled message")
	p := &OrderPayload{Order: o}
	m := &Message{MessageType: ORDER_CANCELED, Payload: p}
	if err := c.connection.WriteJSON(&m); err != nil {
		log.Println("Write: ", err)
		return
	}
}
