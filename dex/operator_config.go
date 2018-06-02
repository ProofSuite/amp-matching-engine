package dex

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

type OperatorConfig struct {
	Admin          *Wallet         `json:"admin"`
	Exchange       common.Address  `json:"exchange"`
	OperatorParams *OperatorParams `json:"operatorParams"`
}

func (c *OperatorConfig) String() string {
	return fmt.Sprintf(
		"Operator Configuration:\n"+
			"Admin Address: %x\nExchange Address: %x\nGas Price: %v\nMax Gas: %v\nMinimum Balance: %v\nRPC URL: %v\n",
		c.Admin.Address, c.Exchange, c.OperatorParams.gasPrice, c.OperatorParams.maxGas, c.OperatorParams.minBalance, c.OperatorParams.rpcURL)
}

// OperatorParams contains numerical values that define how the operator should send transcations.
// rpcURL is the url of the ethereum node. By default it is ws://localhost:8546. For non-websocket
// (not supported a priori) http://localhost:8545
type OperatorParams struct {
	gasPrice   *big.Int `json:"gasPrice"`
	maxGas     uint64   `json"maxGas"`
	minBalance *big.Int `json:"minBalance"`
	rpcURL     string   `json:"rpcURL"`
}
