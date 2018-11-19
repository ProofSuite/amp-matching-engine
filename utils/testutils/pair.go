package testutils

import (
	"math/big"

	"github.com/Proofsuite/amp-matching-engine/types"
	"github.com/ethereum/go-ethereum/common"
)

func GetZRXWETHTestPair() *types.Pair {
	return &types.Pair{
		BaseTokenSymbol:    "ZRX",
		BaseTokenAddress:   common.HexToAddress("0x2034842261b82651885751fc293bba7ba5398156"),
		BaseTokenDecimals:  18,
		QuoteTokenSymbol:   "WETH",
		PriceMultiplier:    big.NewInt(1e6),
		QuoteTokenAddress:  common.HexToAddress("0x276e16ada4b107332afd776691a7fbbaede168ef"),
		QuoteTokenDecimals: 18,
	}
}
