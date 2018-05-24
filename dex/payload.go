package dex

// Payload is the generic that represents websocket messages
type Payload interface{}

type SignedDataPayload struct {
	Order *Order `json:"order"`
}

// OrderPayload is a simple payload which consists of a single order
type OrderPayload struct {
	Order *Order `json:"order"`
}

// TradePayload contains both an order and a trade object corresponding to that order.
type TradePayload struct {
	Order *Order `json:"order"`
	Trade *Trade `json:"trade"`
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

type OrderIdPayload struct {
	OrderId uint64 `json:"orderId"`
}

type RequestSignedDataPayload struct {
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

func NewOrderCancelPayload() *OrderCancelPayload {
	p := &OrderCancelPayload{}
	p.OrderCancel = &OrderCancel{}
	return p
}

// DecodeOrderPayload takes a payload retrieved from a JSON and decodes it into an Order structure
func (o *Order) DecodeOrderPayload(p Payload) {
	payload := p.(map[string]interface{})["order"].(map[string]interface{})
	o.DecodeOrder(payload)
}

// DecodeOrderFilledPayload takes a payload unmarshalled from a JSON byte string and decodes it into
// an OrderFilledPayload
func (d *OrderFilledPayload) DecodeOrderFilledPayload(p Payload) {
	makerOrderPayload := p.(map[string]interface{})["makerOrder"].(map[string]interface{})
	takerOrderPayload := p.(map[string]interface{})["takerOrder"].(map[string]interface{})
	d.MakerOrder.DecodeOrder(makerOrderPayload)
	d.TakerOrder.DecodeOrder(takerOrderPayload)
}

// DecodeTradePayload takes a payload unmarshalled from a JSON byte string and decodes it into a DecodeTradePayload
func (d *TradePayload) DecodeTradePayload(p Payload) {
	trade := p.(map[string]interface{})["trade"].(map[string]interface{})
	order := p.(map[string]interface{})["order"].(map[string]interface{})
	d.Order.DecodeOrder(order)
	d.Trade.DecodeTrade(trade)
}

func (d *OrderCancelPayload) DecodeOrderCancelPayload(p Payload) {
	orderCancel := p.(map[string]interface{})["orderCancel"].(map[string]interface{})
	d.OrderCancel.Decode(orderCancel)
}
