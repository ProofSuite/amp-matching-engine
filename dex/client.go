package dex

import (
	"flag"
	"fmt"
	"log"
	"net/url"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/posener/wstest"
)

var wg = &sync.WaitGroup{}
var addr = flag.String("addr", "localhost:8080", "http service address")

type Client struct {
	connection   *websocket.Conn
	requests     chan *Message
	responses    chan *Message
	requestLogs  []*Message
	responseLogs []*Message
	wallet       *Wallet
}

func NewClient(w *Wallet, s *Server) *Client {
	flag.Parse()
	uri := url.URL{Scheme: "ws", Host: *addr, Path: "/api"}
	log.Printf("Connecting to %s", uri.String())

	d := wstest.NewDialer(s)
	c, _, err := d.Dial(uri.String(), nil)
	if err != nil {
		panic(err)
	}

	requests := make(chan *Message)
	responses := make(chan *Message)
	requestLogs := make([]*Message, 0)
	responseLogs := make([]*Message, 0)

	return &Client{connection: c,
		wallet:       w,
		requests:     requests,
		responses:    responses,
		requestLogs:  requestLogs,
		responseLogs: responseLogs,
	}
}

func (c *Client) start() {
	c.handleMessages()
	c.handleIncomingMessages()
}

func (c *Client) handleMessages() {
	go func() {
		for {
			select {
			case request := <-c.requests:
				c.requestLogs = append(c.requestLogs, request)
				switch request.MessageType {
				case PLACE_ORDER:
					go c.placeOrder(request)
				case SIGNED_DATA:
					go c.sendSignedData(request)
				case CANCEL_ORDER:
					go c.cancelOrder(request)
				case DONE:
					go c.done()
				}
			case response := <-c.responses:
				c.responseLogs = append(c.responseLogs, response)
				switch response.MessageType {
				case ORDER_PLACED:
					go c.handleOrderPlaced(response.Payload)
				case ORDER_CANCELED:
					go c.handleOrderCanceled(response.Payload)
				case REQUEST_SIGNED_DATA:
					go c.handleRequestSignedData(response.Payload)
				case ORDER_FILLED:
					go c.handleOrderFilled(response.Payload)
				case ORDER_PARTIALLY_FILLED:
					go c.handleOrderPartiallyFilled(response.Payload)
				}
			}
		}
	}()
}

func (c *Client) handleIncomingMessages() {
	message := new(Message)
	go func() {
		for {
			err := c.connection.ReadJSON(&message)
			if err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					log.Printf("Error: %#v", err)
				}
				break
			}

			c.responses <- message
		}
	}()
}

func (c *Client) placeOrder(request *Message) {
	err := c.connection.WriteJSON(request)
	if err != nil {
		fmt.Printf("Error: Could not place order. Payload: %#v\n", request.Payload)
		return
	}

	// fmt.Printf("Log: Place Order Message sent:\n")
	// fmt.Printf("%v\n\n", request.Payload)
}

func (c *Client) sendSignedData(request *Message) {
	err := c.connection.WriteJSON(request)
	if err != nil {
		fmt.Printf("Error: Could not send signed orders. Payload: %#v", request.Payload)
		return
	}

	// fmt.Printf("Log: Signed Orders Message sent:\n")
	// fmt.Printf("%v\n\n", request)
}

func (c *Client) cancelOrder(request *Message) {
	err := c.connection.WriteJSON(request)
	if err != nil {
		fmt.Printf("Error: Could not cancel order. Payload: %#v", request.Payload)
		return
	}

	// fmt.Printf("Log: Cancel Orders Message sent:\n")
	// fmt.Printf("%v\n\n", request)
}

func (c *Client) handleOrderPlaced(p Payload) {
	o := &Order{}
	o.DecodeOrderPayload(p)

	// fmt.Printf("Log: Handling Order Placed Message:\n")
	// fmt.Printf("%v\n\n", o)
}

func (c *Client) handleOrderCanceled(p Payload) {
	fmt.Printf("Log: Handling Order Canceled Message. Payload: %#v", p)
}

func (c *Client) handleOrderFilled(p Payload) {
	decoded := NewTradePayload()
	decoded.DecodeTradePayload(p)

	// order := decoded.Order
	trade := decoded.Trade

	err := c.wallet.SignTrade(trade)
	if err != nil {
		fmt.Printf("Error signing trade: %v", err)
	}

	m := &Message{MessageType: SIGNED_DATA, Payload: RequestSignedDataPayload{Trade: trade}}
	c.requests <- m
	// trade, err := c.wallet.
}

func (c *Client) handleOrderPartiallyFilled(p Payload) {
	decoded := NewOrderFilledPayload()
	decoded.DecodeOrderFilledPayload(p)

	trade, err := c.wallet.NewTrade(decoded.MakerOrder, decoded.TakerOrder.Amount)
	if err != nil {
		fmt.Printf("Error signing trade: %v", err)
	}

	m := &Message{MessageType: SIGNED_DATA, Payload: RequestSignedDataPayload{Trade: trade}}
	c.requests <- m
}

func (c *Client) handleRequestSignedData(p Payload) {
	// fmt.Printf("Log: Handling Request Signed Data. Payload: %#v", p)
	// fakeData := SignedDataPayload{}
	// m := &Message{MessageType: "SIGNED_DATA", Payload: fakeData}
	// c.requests <- m
}

func (c *Client) done() {

}
