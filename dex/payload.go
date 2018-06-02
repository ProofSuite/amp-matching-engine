package dex

import (
	"log"

	"github.com/ethereum/go-ethereum/common"
)

// Payload is the generic that represents websocket messages
type Payload interface{}

// OrderPayload is a simple payload which consists of a single order
type OrderPayload struct {
	Order *Order `json:"order"`
}

// TradePayload contains both an order and a trade object corresponding to that order.
type TradePayload struct {
	Order *Order `json:"order"`
	Trade *Trade `json:"trade"`
}

// TxSuccessPayload contains an order, a trade and the transaction hash of the successful blockchain
// transaction.
type TxSuccessPayload struct {
	Order *Order      `json:"order"`
	Trade *Trade      `json:"trade"`
	Tx    common.Hash `json:"tx"`
}

// TxErrorPayload contains and order, a trade and an error id. The error id corresponds to the error id
//
type TxErrorPayload struct {
	Order   *Order `json:"order"`
	Trade   *Trade `json:"trade"`
	ErrorId uint8  `json:"errorId"`
}

// CancelOrderPayload contains both an OrderId and the pairID of the corresponding orderbook/token pair
type OrderCancelPayload struct {
	OrderCancel *OrderCancel `json:"orderCancel"`
}

// OrderFilledPayload contains a TakerOrder and MakerOrder that corresponds to an order pair that was just matched by
// the matching engine
type OrderFilledPayload struct {
	TakerOrder *Order `json:"takerOrder"`
	MakerOrder *Order `json:"makerOrder"`
}

// OrderExecutedPayload contains the order that was just matched and a transaction hash. The client can use this
// payload to follow the state of the transaction. This payload may change in the future
type OrderExecutedPayload struct {
	Order *Order      `json:"order"`
	Tx    common.Hash `json:"tx"`
}

// TradeExecutedPayload contains a trade that matched and a transaction hash associated with the transaction on that trade.
// The client can use this paylaod to follow the state of the transaction. This payload may change in the future to include
// the order in addition to the trade.
type TradeExecutedPayload struct {
	Trade *Trade      `json:"trade"`
	Tx    common.Hash `json:tx"`
}

// SignedDataPayload contains a trade that was just signed by the client
type SignedDataPayload struct {
	Trade *Trade `json:"trade"`
}

// NewOrderFilledPayload creates a new empty OrderFilledPayload
func NewOrderFilledPayload() *OrderFilledPayload {
	p := &OrderFilledPayload{}
	p.MakerOrder = &Order{}
	p.TakerOrder = &Order{}
	return p
}

// NewTradePayload creates a new empty TradePayload
func NewTradePayload() *TradePayload {
	p := &TradePayload{}
	p.Order = &Order{}
	p.Trade = &Trade{}
	return p
}

// NewOrderCancelPayload creates a new empty OrderCancelPayload
func NewOrderCancelPayload() *OrderCancelPayload {
	p := &OrderCancelPayload{}
	p.OrderCancel = &OrderCancel{}
	return p
}

// NewOrderExecutedPayload creates a new empty OrderExecutedPayload
func NewOrderExecutedPayload() *OrderExecutedPayload {
	p := &OrderExecutedPayload{}
	p.Order = &Order{}
	p.Tx = common.Hash{}
	return p
}

// NewTradeExecutedPayload creates a new empty TradeExecutedPayload
func NewTradeExecutedPayload() *TradeExecutedPayload {
	p := &TradeExecutedPayload{}
	p.Trade = &Trade{}
	p.Tx = common.Hash{}
	return p
}

// NewTxSuccessPayload creates a new empty TxSuccessPayload
func NewTxSuccessPayload() *TxSuccessPayload {
	p := &TxSuccessPayload{}
	p.Order = &Order{}
	p.Trade = &Trade{}
	p.Tx = common.Hash{}
	return p
}

// NewTxErrorPayload creates a new empty TxErrorPayload
func NewTxErrorPayload() *TxErrorPayload {
	p := &TxErrorPayload{}
	p.Order = &Order{}
	p.Trade = &Trade{}
	return p
}

// DecodeOrderPayload takes a payload retrieved from a JSON and decodes it into an Order structure
func (o *Order) DecodeOrderPayload(p Payload) error {
	payload := p.(map[string]interface{})["order"].(map[string]interface{})
	err := o.Decode(payload)
	if err != nil {
		return err
	}

	return nil
}

// DecodeTradePayload takes a payload retrieved from a JSON file and decodes it into a Trade struct
func (tp *TradePayload) DecodeTradePayload(p Payload) error {
	trade := p.(map[string]interface{})["trade"].(map[string]interface{})
	err := tp.Trade.Decode(trade)
	if err != nil {
		return err
	}

	return nil
}

// DecodeOrderFilledPayload takes a payload unmarshalled from a JSON byte string and decodes it into
// an OrderFilledPayload
func (d *OrderFilledPayload) DecodeOrderFilledPayload(p Payload) error {
	makerOrderPayload := p.(map[string]interface{})["makerOrder"].(map[string]interface{})
	takerOrderPayload := p.(map[string]interface{})["takerOrder"].(map[string]interface{})
	err := d.MakerOrder.Decode(makerOrderPayload)
	if err != nil {
		return err
	}
	err = d.TakerOrder.Decode(takerOrderPayload)
	if err != nil {
		return err
	}

	return nil
}

// DecodeOrderCancelPayload takes a payload that was previously unmarshalled from a JSON byte string and decodes it
// into an OrderCancelPayload
func (ocp *OrderCancelPayload) DecodeOrderCancelPayload(p Payload) error {
	orderCancel := p.(map[string]interface{})["orderCancel"].(map[string]interface{})
	err := ocp.OrderCancel.Decode(orderCancel)
	if err != nil {
		return err
	}

	return nil
}

// DecodeOrderExecutedPayload takes a payload that was previously unmarshalled from a JSON byte string and decodes it
// into a OrderExecutedPayload
func (oep *OrderExecutedPayload) DecodeOrderExecutedPayload(p Payload) error {
	log.Printf("The payload is equal to : %v\n", p)
	orderPayload := p.(map[string]interface{})["order"].(map[string]interface{})
	txHash := common.HexToHash(p.(map[string]interface{})["tx"].(string))

	order := &Order{}
	err := order.Decode(orderPayload)
	if err != nil {
		return err
	}

	oep.Order = order
	oep.Tx = txHash
	return nil
}

// DecodeTradeExecutedPayload takes a payload that was previously unmarshalled from a JSON byte string and decodes it
// into a TradeExecutedPayload.
// Developer note: For some reason there is currently a difference between DecodeTradeExecutedPayload and DeodeOrderExecutedPayload
// In the first one, the transaction hash is saved under the "tx" key while it is saved under the "Tx" on the second one.
func (tep *TradeExecutedPayload) DecodeTradeExecutedPayload(p Payload) error {
	log.Printf("The payload is equal to : %v\n", p)
	tradePayload := p.(map[string]interface{})["trade"].(map[string]interface{})
	txHash := common.HexToHash(p.(map[string]interface{})["Tx"].(string))

	// common.HexToHash(p.(string)["tx"])

	trade := &Trade{}
	err := trade.Decode(tradePayload)
	if err != nil {
		return err
	}

	tep.Trade = trade
	tep.Tx = txHash
	return nil
}

// DecodeTxErrorPayload takes a payload that was previously unmarshalled from a JSON byte string and decodes it
// into a TxErrorPayload.
func (d *TxErrorPayload) DecodeTxErrorPayload(p Payload) error {
	orderPayload := p.(map[string]interface{})["order"].(map[string]interface{})
	tradePayload := p.(map[string]interface{})["trade"].(map[string]interface{})
	errId := p.(map[string]interface{})["errorId"].(uint8)

	trade := &Trade{}
	err := trade.Decode(tradePayload)
	if err != nil {
		return err
	}

	order := &Order{}
	err = order.Decode(orderPayload)
	if err != nil {
		return err
	}

	d.Trade = trade
	d.Order = order
	d.ErrorId = errId

	return nil
}

func (d *TxSuccessPayload) DecodeTxSuccessPayload(p Payload) error {
	orderPayload := p.(map[string]interface{})["order"].(map[string]interface{})
	tradePayload := p.(map[string]interface{})["trade"].(map[string]interface{})
	txHash := common.HexToHash(p.(map[string]interface{})["tx"].(string))

	trade := &Trade{}
	err := trade.Decode(tradePayload)
	if err != nil {
		return err
	}

	order := &Order{}
	err = order.Decode(orderPayload)
	if err != nil {
		return err
	}

	d.Trade = trade
	d.Order = order
	d.Tx = txHash

	return nil
}
