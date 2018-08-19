package mocks

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"sync"

	"github.com/Proofsuite/amp-matching-engine/types"
	"github.com/Proofsuite/go-ethereum/common"
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
	requests       chan *types.WebSocketMessage
	responses      chan *types.WebSocketMessage
	requestLogs    []*types.WebSocketMessage
	responseLogs   []*types.WebSocketMessage
	ethereumClient *ethclient.Client
	wallet         *types.Wallet
	mutex          sync.Mutex
	logs           chan *ClientLogMessage
}

// The client log is mostly used for testing. It optionally takes orders, trade,
// error ids and transaction hashes. All these parameters are optional in order to
// allow the client log message to take in a lot of different types of messages
// An error id of -1 means that there was no error.
type ClientLogMessage struct {
	MessageType string       `json:"messageType"`
	Order       *types.Order `json:"order"`
	Trade       *types.Trade `json:"trade"`
	Tx          common.Hash  `json:"tx"`
	ErrorID     int8         `json:"errorID"`
}

type Server interface {
	ServeHTTP(res http.ResponseWriter, req *http.Request)
}

// NewClient a default client struct connected to the given server
func NewClient(w *types.Wallet, s *Server) *Client {
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

	reqs := make(chan *types.WebSocketMessage)
	resps := make(chan *types.WebSocketMessage)
	logs := make(chan *ClientLogMessage)
	reqLogs := make([]*types.WebSocketMessage, 0)
	respLogs := make([]*types.WebSocketMessage, 0)

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
			case msg := <-c.requests:
				fmt.Printf("Handling request: %v\n", msg)
				c.requestLogs = append(c.requestLogs, msg)

				c.handleOrderChannelMessagesOut(msg)

			case msg := <-c.responses:
				fmt.Printf("Handling response: %v\n")
				c.responseLogs = append(c.responseLogs, msg)

				switch msg.Channel {
				case "orders":
					go c.handleOrderChannelMessagesIn(msg.Payload)
				case "order_book":
					go c.handleOrderBookChannelMessages(msg.Payload)
				case "trades":
					go c.handleTradeChannelMessages(msg.Payload)
				case "ohlcv":
					go c.handleOHLCVMessages(msg.Payload)
				}
			}
		}
	}()
}

// handleChannelMessagesOut
func (c *Client) handleOrderChannelMessagesOut(msg interface{}) {
	msg := &types.Message{}

	err := json.Unmarshal()

	switch msg.Type {
	case "NEW_ORDER":
		c.SendNewOrder(msg)
	case "SUBMIT_SIGNATURE":
		c.SendSubmitNewSignature(msg)
	case "CANCEL_ORDER":
		c.SendCancelOrder(msg)
	case "DONE":
		c.done()
	}
}

// handleChannelMessagesIn
func (c *Client) handleOrderChannelMessagesIn(p types.WebSocketPayload) {
	switch p.Type {
	case "ORDER_ADDED":
		c.handleOrderAdded(p)
	case "ORDER_CANCELED":
		c.handleOrderCanceled(p)
	case "REQUEST_SIGNATURE":
		c.handleSignatureRequested(p)
	case "TRADE_EXECUTED":
		c.handleTradeExecuted(p)
	case "TRADE_TX_SUCCESS":
		c.handleOrderTxSuccess(p)
	case "TRADE_TX_ERROR":
		c.handleOrderTxError(p)
	}
}

func (c *Client) handleOrderBookChannelMessages(p types.WebSocketPayload) {
	switch p.Type {
	case "INIT":
		c.handleOrderBookInit(p)
	case "UPDATE":
		c.handleOrderBookUpdate(p)
	}
}

func (c *Client) handleTradeChannelMessages(p types.WebSocketPayload) {
	switch p.Type {
	case "INIT":
		c.handleTradesInit(p)
	case "UPDATE":
		c.handleTradesUpdate(p)
	}
}

func (c *Client) handleOHLCVMessages(p types.WebSocketPayload) {
	switch p.Type {
	case "INIT":
		c.handleOHLCVInit(p)
	case "UPDATE":
		c.handleOHLCVUpdate(p)
	}
}

// handleIncomingMessages reads incomings JSON messages from the websocket connection and
// feeds them into the responses channel
func (c *Client) handleIncomingMessages() {
	message := new(types.WebSocketMessage)
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

func (c *Client) handleOrderAdded(p types.WebSocketPayload) {

	l := &ClientLogMessage{
		MessageType: resp.MessageType,
		Order:       decoded.Order,
		Trade:       decoded.Trade,
		ErrorID:     int8(decoded.ErrorId),
	}
}

func (c *Client) handleOrderCanceled(p types.WebSocketPayload) {

}

func (c *Client) handleSignatureRequested(p types.WebSocketPayload) {

}

func (c *Client) handleTradeExecuted(p types.WebSocketPayload) {

}

func (c *Client) handleOrderTxSuccess(p types.WebSocketPayload) {

}

func (c *Client) handleOrderTxError(p types.WebSocketPayload) {

}

func (c *Client) handleOrderBookInit(p types.WebSocketPayload) {

}

func (c *Client) handleOrderBookUpdate(p types.WebSocketPayload) {

}

func (c *Client) handleTradesInit(p types.WebSocketPayload) {

}

func (c *Client) handleTradesUpdate(p types.WebSocketPayload) {

}

func (c *Client) handleOHLCVInit(p types.WebSocketPayload) {

}

func (c *Client) handleOHLCVUpdate(p types.WebSocketPayload) {

}

// func (c *Client) placeOrder(req *Message) {
// 	err := c.send(req)
// 	if err != nil {
// 		log.Printf("Error: Could not place order. Payload: %#v\n", req.Payload)
// 		return
// 	}
// }

// func (c *Client) sendSignedData(req *Message) {
// 	err := c.send(req)
// 	if err != nil {
// 		log.Printf("Error: Could not send signed orders. Payload: %#v", req.Payload)
// 		return
// 	}
// }

// func (c *Client) cancelOrder(req *Message) {
// 	err := c.send(req)
// 	if err != nil {
// 		log.Printf("Error: Could not cancel order. Payload: %#v", req.Payload)
// 		return
// 	}
// }

// func (c *Client) handleOrderPlaced(resp *Message) {
// 	o := &Order{}
// 	o.DecodeOrderPayload(resp.Payload)

// 	l := &ClientLogMessage{
// 		MessageType: resp.MessageType,
// 		Order:       o,
// 	}

// 	c.logs <- l
// 	fmt.Printf("Log: Handling Order Placed Message:%v\n\n", o)
// }

// func (c *Client) handleOrderCanceled(r *Message) {
// 	fmt.Printf("Log: Handling Order Canceled Message. Payload: %#v", r.Payload)
// }

// func (c *Client) handleOrderFilled(resp *Message) {
// 	decoded := NewTradePayload()
// 	err := decoded.DecodeTradePayload(resp.Payload)
// 	if err != nil {
// 		fmt.Printf("Could not decode trade payload: %v", err)
// 	}

// 	trade := decoded.Trade

// 	err = c.wallet.SignTrade(trade)
// 	if err != nil {
// 		fmt.Printf("Error signing trade: %v", err)
// 	}

// 	l := &ClientLogMessage{
// 		MessageType: resp.MessageType,
// 		Order:       decoded.Order,
// 		Trade:       trade,
// 	}

// 	c.logs <- l
// 	req := &Message{MessageType: SIGNED_DATA, Payload: SignedDataPayload{Trade: trade}}
// 	c.requests <- req
// 	fmt.Printf("Log: Handling Order Filled Message. Payload: %#v", trade)
// }

// func (c *Client) handleOrderPartiallyFilled(resp *Message) {
// 	decoded := NewOrderFilledPayload()
// 	decoded.DecodeOrderFilledPayload(resp.Payload)

// 	trade, err := c.wallet.NewTrade(decoded.MakerOrder, decoded.TakerOrder.Amount)
// 	if err != nil {
// 		fmt.Printf("Error signing trade: %v", err)
// 	}

// 	req := &Message{MessageType: SIGNED_DATA, Payload: SignedDataPayload{Trade: trade}}
// 	c.requests <- req
// }

// func (c *Client) handleOrderExecuted(resp *Message) {
// 	decoded := NewOrderExecutedPayload()
// 	decoded.DecodeOrderExecutedPayload(resp.Payload)

// 	l := &ClientLogMessage{
// 		MessageType: resp.MessageType,
// 		Order:       decoded.Order,
// 		Tx:          decoded.Tx,
// 	}

// 	c.logs <- l
// 	fmt.Printf("\nLog: Handling Order executed message. Tx Hash: %v", decoded.Tx)
// }

// func (c *Client) handleTradeExecuted(resp *Message) {
// 	decoded := NewTradeExecutedPayload()
// 	decoded.DecodeTradeExecutedPayload(resp.Payload)

// 	l := &ClientLogMessage{
// 		MessageType: resp.MessageType,
// 		Trade:       decoded.Trade,
// 		Tx:          decoded.Tx,
// 	}

// 	c.logs <- l
// 	fmt.Printf("\nLog: Handling Trade Executed message. Payload: %x\n\n", decoded.Tx)
// }

// func (c *Client) handleOrderTxSuccess(resp *Message) {
// 	decoded := NewTxSuccessPayload()
// 	decoded.DecodeTxSuccessPayload(resp.Payload)

// 	l := &ClientLogMessage{
// 		MessageType: resp.MessageType,
// 		Order:       decoded.Order,
// 		Trade:       decoded.Trade,
// 		Tx:          decoded.Tx,
// 		ErrorID:     -1,
// 	}

// 	c.logs <- l
// 	fmt.Printf("\nLog: Handling Order Tx Success message. Payload: %x\n\n", decoded.Tx)
// }

// func (c *Client) handleOrderTxError(resp *Message) {
// 	decoded := NewTxErrorPayload()
// 	decoded.DecodeTxErrorPayload(resp.Payload)
// 	errId := decoded.ErrorId

// 	l := &ClientLogMessage{
// 		MessageType: resp.MessageType,
// 		Order:       decoded.Order,
// 		Trade:       decoded.Trade,
// 		ErrorID:     int8(decoded.ErrorId),
// 	}

// 	c.logs <- l
// 	fmt.Printf("\nLog: Handling Tx Error Message. Error ID: %v\n", errId)
// }

// func (c *Client) handleTradeTxSuccess(resp *Message) {
// 	decoded := NewTxSuccessPayload()
// 	decoded.DecodeTxSuccessPayload(resp.Payload)

// 	l := &ClientLogMessage{
// 		MessageType: resp.MessageType,
// 		Order:       decoded.Order,
// 		Trade:       decoded.Trade,
// 		Tx:          decoded.Tx,
// 		ErrorID:     -1,
// 	}

// 	c.logs <- l
// 	fmt.Printf("\nLog: Handling Order Tx Success message. Payload: %x\n\n", decoded.Tx)
// }

// func (c *Client) handleTradeTxError(resp *Message) {
// 	decoded := NewTxErrorPayload()
// 	decoded.DecodeTxErrorPayload(resp.Payload)

// 	l := &ClientLogMessage{
// 		MessageType: resp.MessageType,
// 		Order:       decoded.Order,
// 		Trade:       decoded.Trade,
// 		ErrorID:     int8(decoded.ErrorId),
// 	}

// 	c.logs <- l
// 	fmt.Printf("\nLog: Handling Tx Error Message. Error ID: %v\n", decoded.ErrorId)
// }

// func (c *Client) done() {

// }
