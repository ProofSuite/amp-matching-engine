package dex

import "github.com/ethereum/go-ethereum/common"

type OperatorConfig struct {
	Admin          *Wallet
	Exchange       common.Address
	OperatorParams *OperatorParams
}
