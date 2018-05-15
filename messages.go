package main

import "fmt"

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
	DONE                   = "DONE"
)

type Message struct {
	MessageType MessageType `json:"messageType"`
	Payload     Payload     `json:"payload"`
}

func (m *Message) String() string {
	return fmt.Sprintf("\n Message:{messageType:%v, payload:%v", m.MessageType, m.Payload)
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
type Payload interface{}

type PlaceOrderMesage struct {
	MessageType MessageType  `json:"messageType"`
	Payload     OrderPayload `json:"payload"`
}

type SignedDataMessage struct {
	MessageType MessageType       `json:"messageType"`
	Payload     SignedDataPayload `json:"payload"`
}

type SignedDataPayload struct {
	Order *Order `json:"order"`
}

type CancelOrderMessage struct {
	MessageType MessageType        `json:"messageType"`
	Payload     CancelOrderPayload `json:"payload"`
}

//Messages from server to client
type OrderPlacedMessage struct {
	MessageType MessageType  `json:"messageType"`
	Payload     OrderPayload `json:"payload"`
}

type OrderPayload struct {
	Order *Order `json:"order"`
}

type OrderDataPayload struct {
	Order *OrderData `json:"order"`
}

type OrderIdPayload struct {
	OrderId uint64 `json:"orderId"`
}

type OrderCanceledMessage struct {
	MessageType MessageType    `json:"messageType"`
	Payload     OrderIdPayload `json:"payload"`
}

type RequestSignedDataMessage struct {
	MessageType MessageType              `json:"messageType"`
	Payload     RequestSignedDataPayload `json:"payload"`
}

type RequestSignedDataPayload struct {
	Orders []*Order `json:"order"`
}

type OrderFilledMessage struct {
	MessageType MessageType  `json:"messageType"`
	Payload     OrderPayload `json:"payload"`
}

type OrderPartiallyFilledMessage struct {
	MessageType MessageType  `json:"messageType"`
	Payload     OrderPayload `json:"payload"`
}

type CancelOrderPayload struct {
	OrderId uint64 `json:"orderId"`
	Pair    Pair   `json:"pair"`
}

func NewOrderPlacedMessage(o *Order) *Message {
	return &Message{MessageType: ORDER_PLACED, Payload: OrderPayload{Order: o}}
}

func NewOrderCanceledMessage(o *Order) *Message {
	return &Message{MessageType: ORDER_CANCELED, Payload: OrderIdPayload{OrderId: o.Id}}
}

func NewRequestSignedDataMessage(o *Order) *Message {
	return &Message{MessageType: REQUEST_SIGNED_DATA, Payload: RequestSignedDataPayload{Orders: []*Order{}}}
}

// func (*op OrderPayload) DecodeOrderPayload(data interface{}) {
// 	var orderPayload OrderPayload

// 	orderPayload := data
// }
