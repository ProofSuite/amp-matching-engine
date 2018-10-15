package types

import "github.com/ethereum/go-ethereum/common"

type OrderTradePair struct {
	Order *Order `json:"order"`
	Trade *Trade `json:"trade"`
}

type EngineResponse struct {
	Status         string            `json:"fillStatus,omitempty"`
	HashID         common.Hash       `json:"hashID, omitempty"`
	Order          *Order            `json:"order,omitempty"`
	RemainingOrder *Order            `json:"remainingOrder,omitempty"`
	Matches        []*OrderTradePair `json:"matches,omitempty"`
}
