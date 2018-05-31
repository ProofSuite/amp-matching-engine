package dex

import (
	"encoding/json"
	"fmt"

	"github.com/kr/pretty"
)

type MessageType string

const (
	PLACE_ORDER            = "PLACE_ORDER"
	CANCEL_ORDER           = "CANCEL_ORDER"
	ORDER_PLACED           = "ORDER_PLACED"
	SIGNED_DATA            = "SIGNED_DATA"
	REQUEST_SIGNED_DATA    = "REQUEST_SIGNED_DATA"
	ORDER_PARTIALLY_FILLED = "ORDER_PARTIALLY_FILLED"
	ORDER_FILLED           = "ORDER_FILLED"
	ORDER_CANCELED         = "ORDER_CANCELED"
	ORDER_MATCHED          = "ORDER_MATCHED"
	ORDER_EXECUTED         = "ORDER_EXECUTED"
	ORDER_TX_SUCCESS       = "ORDER_TX_SUCCESS"
	ORDER_TX_ERROR         = "ORDER_TX_ERROR"
	TRADE_EXECUTED         = "TRADE_EXECUTED"
	TRADE_TX_SUCCESS       = "TRADE_TX_SUCCESS"
	TRADE_TX_ERROR         = "TRADE_TX_ERROR"
	DONE                   = "DONE"
)

type Message struct {
	MessageType MessageType `json:"messageType"`
	Payload     Payload     `json:"payload"`
}

func (m *Message) UnmarshalJSON(b []byte) error {
	message := map[string]interface{}{}

	err := json.Unmarshal(b, &message)
	if err != nil {
		return err
	}

	m.MessageType = MessageType(message["messageType"].(string))
	m.Payload = message["payload"]

	return nil
}

func (m *Message) String() string {
	return fmt.Sprintf("\nMessage:\nMessageType: %v\nPayload:\n%v\n", pretty.Formatter(m.MessageType), pretty.Formatter(m.Payload))
}

//Messages from client to server
// 1. Place order message
// 2. Signed data message
// 3. Cancel order message

//Messages from server to client
// 1. Order created message
// 2. Order canceled message
// 3. Request Signed Data message
// 4. Order partially filled message
// 5. Order filled message

//Messages from client to server
type PlaceOrderMesage struct {
	MessageType MessageType  `json:"messageType"`
	Payload     OrderPayload `json:"payload"`
}

type SignedDataMessage struct {
	MessageType MessageType       `json:"messageType"`
	Payload     SignedDataPayload `json:"payload"`
}

type CancelOrderMessage struct {
	MessageType MessageType        `json:"messageType"`
	Payload     OrderCancelPayload `json:"payload"`
}

//Messages from server to client
type OrderPlacedMessage struct {
	MessageType MessageType  `json:"messageType"`
	Payload     OrderPayload `json:"payload"`
}

type OrderCanceledMessage struct {
	MessageType MessageType    `json:"messageType"`
	Payload     OrderIdPayload `json:"payload"`
}

type RequestSignedDataMessage struct {
	MessageType MessageType              `json:"messageType"`
	Payload     RequestSignedDataPayload `json:"payload"`
}

type OrderFilledMessage struct {
	MessageType MessageType  `json:"messageType"`
	Payload     OrderPayload `json:"payload"`
}

type OrderPartiallyFilledMessage struct {
	MessageType MessageType  `json:"messageType"`
	Payload     OrderPayload `json:"payload"`
}

type OrderExecutedMessage struct {
	MessageType MessageType          `json:"messageType"`
	Payload     OrderExecutedPayload `json:"payload"`
}

type TradeExecutedMessage struct {
	MessageType MessageType          `json:"messageType"`
	Payload     TradeExecutedPayload `json:"payload"`
}

// The client log is mostly used for testing. It optionally takes orders, trade,
// error ids and transaction hashes. All these parameters are optional in order to
// allow the client log message to take in a lot of different types of messages
// An error id of -1 means that there was no error.
type ClientLogMessage struct {
	MessageType MessageType `json:"messageType"`
	Order       *Order      `json:"order"`
	Trade       *Trade      `json:"trade"`
	Tx          common.Hash `json:"tx"`
	ErrorID     int8        `json:"errorID"`
}

func (m *ClientLogMessage) String() string {
	return fmt.Sprintf("\nMessageType: %v\nOrder: %v\nTrade: %v\nTx: %v\nErrorID: %v\n\n",
		m.MessageType, m.Order, m.Trade, m.Tx, m.ErrorID)
}
