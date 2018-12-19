package units

import (
	"math/big"

	"github.com/Proofsuite/amp-matching-engine/utils/math"
)

func Ethers(value int64) *big.Int {
	return math.Mul(big.NewInt(1e18), big.NewInt(value))
}

func E36() *big.Int {
	return math.Mul(big.NewInt(1e18), big.NewInt(1e18))
}
