package dex

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/url"
	"sync"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/gorilla/websocket"
	"github.com/posener/wstest"
)

var wg = &sync.WaitGroup{}
var addr = flag.String("addr", "localhost:8080", "http service address")

// Client simulates the client websocket handler that will be used to perform trading.
// requests and responses are respectively the outbound and incoming messages.
// requestLogs and responseLogs are arrays of messages that denote the history of received messages
// wallet is the ethereum account used for orders and trades.
// mutex is used to prevent concurrent writes on the websocket connection
type Client struct {
	connection     *websocket.Conn
	requests       chan *Message
	responses      chan *Message
	requestLogs    []*Message
	responseLogs   []*Message
	ethereumClient *ethclient.Client
	wallet         *Wallet
	mutex          sync.Mutex
	logs           chan *ClientLogMessage
}

// NewClient a default client struct connected to the given server
func NewClient(w *Wallet, s *Server) *Client {
	flag.Parse()
	uri := url.URL{Scheme: "ws", Host: *addr, Path: "/api"}

	rpcClient, err := rpc.DialWebsocket(context.Background(), "ws://127.0.0.1:8546", "")
	if err != nil {
		log.Printf("Could not connect to ethereum client")
	}

	ethClient := ethclient.NewClient(rpcClient)

	d := wstest.NewDialer(s)
	c, _, err := d.Dial(uri.String(), nil)
	if err != nil {
		panic(err)
	}

	reqs := make(chan *Message)
	resps := make(chan *Message)
	logs := make(chan *ClientLogMessage)
	reqLogs := make([]*Message, 0)
	respLogs := make([]*Message, 0)

	return &Client{connection: c,
		wallet:         w,
		requests:       reqs,
		logs:           logs,
		ethereumClient: ethClient,
		responses:      resps,
		requestLogs:    reqLogs,
		responseLogs:   respLogs,
	}
}

// send is used to prevent concurrent writes on the websocket connection
func (c *Client) send(v interface{}) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	return c.connection.WriteJSON(v)
}

// start listening and handling incoming messages
func (c *Client) start() {
	c.handleMessages()
	c.handleIncomingMessages()
}

// handleMessages waits for incoming messages and routes messages to the
// corresponding handler.
// requests are the messages that are written on the client and destined to
// the server. responses are the message that are
func (c *Client) handleMessages() {
	go func() {
		for {
			select {
			case request := <-c.requests:
				// log.Printf("Request is equal to %v", request)
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
				case ORDER_EXECUTED:
					go c.handleOrderExecuted(response.Payload)
				case ORDER_TX_SUCCESS:
					go c.handleOrderTxSuccess(response.Payload)
				case ORDER_TX_ERROR:
					go c.handleOrderTxError(response.Payload)
				case TRADE_EXECUTED:
					go c.handleTradeExecuted(response.Payload)
				case TRADE_TX_SUCCESS:
					go c.handleTradeTxSuccess(response.Payload)
				case TRADE_TX_ERROR:
					go c.handleTradeTxError(response.Payload)
				}
			}
		}
	}()
}

// handleIncomingMessages reads incomings JSON messages from the websocket connection and
// feeds them into the responses channel
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
	err := c.send(request)
	if err != nil {
		fmt.Printf("Error: Could not place order. Payload: %#v\n", request.Payload)
		return
	}

	// fmt.Printf("Log: Place Order Message sent:\n")
	// fmt.Printf("%v\n\n", request.Payload)
}

func (c *Client) sendSignedData(request *Message) {
	err := c.send(request)
	if err != nil {
		fmt.Printf("Error: Could not send signed orders. Payload: %#v", request.Payload)
		return
	}

	// fmt.Printf("Log: Signed Orders Message sent:\n")
	// fmt.Printf("%v\n\n", request)
}

func (c *Client) cancelOrder(request *Message) {
	err := c.send(request)
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

	t := decoded.Trade

	err := c.wallet.SignTrade(t)
	if err != nil {
		fmt.Printf("Error signing trade: %v", err)
	}

	m := &Message{MessageType: SIGNED_DATA, Payload: RequestSignedDataPayload{Trade: t}}
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

func (c *Client) handleOrderExecuted(p Payload) {
	log.Printf("\nLog: Handling Order executed message. Payload: %#v\n\n", p)
}

func (c *Client) handleOrderTxSuccess(p Payload) {
	log.Printf("\nLog: Handling Order Tx Success message. Payload: %#v\n\n", p)
}

func (c *Client) handleOrderTxError(p Payload) {
	log.Printf("\nLog: Handling Order Tx Error message. Payload: %#v\n\n", p)
}

func (c *Client) handleTradeExecuted(p Payload) {
	log.Printf("\nLog: Handling Trade Executed message. Payload: %#v\n\n", p)
}

func (c *Client) handleTradeTxSuccess(p Payload) {
	log.Printf("\nLog: Handling Trade Tx Success message. Payload: %#v\n\n", p)
}

func (c *Client) handleTradeTxError(p Payload) {
	log.Printf("\nLog: Handling Trade Tx Error message. Payload: %#v\n\n", p)
}

func (c *Client) done() {

}
