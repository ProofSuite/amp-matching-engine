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
			case req := <-c.requests:
				fmt.Printf("Handling Request: \n%v\n", req)
				c.requestLogs = append(c.requestLogs, req)
				switch req.MessageType {
				case PLACE_ORDER:
					go c.placeOrder(req)
				case SIGNED_DATA:
					go c.sendSignedData(req)
				case CANCEL_ORDER:
					go c.cancelOrder(req)
				case DONE:
					go c.done()
				}
			case resp := <-c.responses:
				c.responseLogs = append(c.responseLogs, resp)
				switch resp.MessageType {
				case ORDER_PLACED:
					go c.handleOrderPlaced(resp)
				case ORDER_CANCELED:
					go c.handleOrderCanceled(resp)
				case ORDER_FILLED:
					go c.handleOrderFilled(resp)
				case ORDER_EXECUTED:
					go c.handleOrderExecuted(resp)
				case ORDER_TX_SUCCESS:
					go c.handleOrderTxSuccess(resp)
				case ORDER_TX_ERROR:
					go c.handleOrderTxError(resp)
				case TRADE_EXECUTED:
					go c.handleTradeExecuted(resp)
				case TRADE_TX_SUCCESS:
					go c.handleTradeTxSuccess(resp)
				case TRADE_TX_ERROR:
					go c.handleTradeTxError(resp)
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

func (c *Client) placeOrder(req *Message) {
	err := c.send(req)
	if err != nil {
		log.Printf("Error: Could not place order. Payload: %#v\n", req.Payload)
		return
	}
}

func (c *Client) sendSignedData(req *Message) {
	err := c.send(req)
	if err != nil {
		log.Printf("Error: Could not send signed orders. Payload: %#v", req.Payload)
		return
	}
}

func (c *Client) cancelOrder(req *Message) {
	err := c.send(req)
	if err != nil {
		log.Printf("Error: Could not cancel order. Payload: %#v", req.Payload)
		return
	}
}

func (c *Client) handleOrderPlaced(resp *Message) {
	o := &Order{}
	o.DecodeOrderPayload(resp.Payload)

	l := &ClientLogMessage{
		MessageType: resp.MessageType,
		Order:       o,
	}

	c.logs <- l
	fmt.Printf("Log: Handling Order Placed Message:%v\n\n", o)
}

func (c *Client) handleOrderCanceled(r *Message) {
	fmt.Printf("Log: Handling Order Canceled Message. Payload: %#v", r.Payload)
}

func (c *Client) handleOrderFilled(resp *Message) {
	decoded := NewTradePayload()
	err := decoded.DecodeTradePayload(resp.Payload)
	if err != nil {
		fmt.Printf("Could not decode trade payload: %v", err)
	}

	trade := decoded.Trade

	err = c.wallet.SignTrade(trade)
	if err != nil {
		fmt.Printf("Error signing trade: %v", err)
	}

	l := &ClientLogMessage{
		MessageType: resp.MessageType,
		Order:       decoded.Order,
		Trade:       trade,
	}

	c.logs <- l
	req := &Message{MessageType: SIGNED_DATA, Payload: SignedDataPayload{Trade: trade}}
	c.requests <- req
	fmt.Printf("Log: Handling Order Filled Message. Payload: %#v", trade)
}

func (c *Client) handleOrderPartiallyFilled(resp *Message) {
	decoded := NewOrderFilledPayload()
	decoded.DecodeOrderFilledPayload(resp.Payload)

	trade, err := c.wallet.NewTrade(decoded.MakerOrder, decoded.TakerOrder.Amount)
	if err != nil {
		fmt.Printf("Error signing trade: %v", err)
	}

	req := &Message{MessageType: SIGNED_DATA, Payload: SignedDataPayload{Trade: trade}}
	c.requests <- req
}

func (c *Client) handleOrderExecuted(resp *Message) {
	decoded := NewOrderExecutedPayload()
	decoded.DecodeOrderExecutedPayload(resp.Payload)

	l := &ClientLogMessage{
		MessageType: resp.MessageType,
		Order:       decoded.Order,
		Tx:          decoded.Tx,
	}

	c.logs <- l
	fmt.Printf("\nLog: Handling Order executed message. Tx Hash: %v", decoded.Tx)
}

func (c *Client) handleTradeExecuted(resp *Message) {
	decoded := NewTradeExecutedPayload()
	decoded.DecodeTradeExecutedPayload(resp.Payload)

	l := &ClientLogMessage{
		MessageType: resp.MessageType,
		Trade:       decoded.Trade,
		Tx:          decoded.Tx,
	}

	c.logs <- l
	fmt.Printf("\nLog: Handling Trade Executed message. Payload: %x\n\n", decoded.Tx)
}

func (c *Client) handleOrderTxSuccess(resp *Message) {
	decoded := NewTxSuccessPayload()
	decoded.DecodeTxSuccessPayload(resp.Payload)

	l := &ClientLogMessage{
		MessageType: resp.MessageType,
		Order:       decoded.Order,
		Trade:       decoded.Trade,
		Tx:          decoded.Tx,
		ErrorID:     -1,
	}

	c.logs <- l
	fmt.Printf("\nLog: Handling Order Tx Success message. Payload: %x\n\n", decoded.Tx)
}

func (c *Client) handleOrderTxError(resp *Message) {
	decoded := NewTxErrorPayload()
	decoded.DecodeTxErrorPayload(resp.Payload)
	errId := decoded.ErrorId

	l := &ClientLogMessage{
		MessageType: resp.MessageType,
		Order:       decoded.Order,
		Trade:       decoded.Trade,
		ErrorID:     int8(decoded.ErrorId),
	}

	c.logs <- l
	fmt.Printf("\nLog: Handling Tx Error Message. Error ID: %v\n", errId)
}

func (c *Client) handleTradeTxSuccess(resp *Message) {
	decoded := NewTxSuccessPayload()
	decoded.DecodeTxSuccessPayload(resp.Payload)

	l := &ClientLogMessage{
		MessageType: resp.MessageType,
		Order:       decoded.Order,
		Trade:       decoded.Trade,
		Tx:          decoded.Tx,
		ErrorID:     -1,
	}

	c.logs <- l
	fmt.Printf("\nLog: Handling Order Tx Success message. Payload: %x\n\n", decoded.Tx)
}

func (c *Client) handleTradeTxError(resp *Message) {
	decoded := NewTxErrorPayload()
	decoded.DecodeTxErrorPayload(resp.Payload)

	l := &ClientLogMessage{
		MessageType: resp.MessageType,
		Order:       decoded.Order,
		Trade:       decoded.Trade,
		ErrorID:     int8(decoded.ErrorId),
	}

	c.logs <- l
	fmt.Printf("\nLog: Handling Tx Error Message. Error ID: %v\n", decoded.ErrorId)
}

func (c *Client) done() {

}
