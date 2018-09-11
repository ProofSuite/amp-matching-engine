package testutils

import (
	"github.com/Proofsuite/amp-matching-engine/types"
	"github.com/ethereum/go-ethereum/common"
)

func GetZRXWETHTestPair() *types.Pair {
	return &types.Pair{
		BaseTokenSymbol:   "ZRX",
		BaseTokenAddress:  common.HexToAddress("0x2034842261b82651885751fc293bba7ba5398156"),
		BaseTokenDecimal:  18,
		QuoteTokenSymbol:  "WETH",
		QuoteTokenAddress: common.HexToAddress("0x276e16ada4b107332afd776691a7fbbaede168ef"),
		QuoteTokenDecimal: 18,
	}
}
