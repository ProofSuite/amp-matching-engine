package dex

import "github.com/ethereum/go-ethereum/core/types"

type OrderTradePair struct {
	order *Order             `json:"order"`
	trade *Trade             `json:"trade"`
	tx    *types.Transaction `json:"tx"`
}
