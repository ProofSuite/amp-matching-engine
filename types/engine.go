package types

type OrderTradePair struct {
	Order *Order
	Trade *Trade
}

type EngineResponse struct {
	Status         string            `json:"fillStatus,omitempty"`
	Order          *Order            `json:"order,omitempty"`
	RemainingOrder *Order            `json:"remainingOrder,omitempty"`
	Matches        []*OrderTradePair `json:"matches,omitempty"`
}
