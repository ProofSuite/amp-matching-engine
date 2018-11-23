package utils

import (
	"math/big"

	"github.com/Proofsuite/amp-matching-engine/utils/math"
)

func ToDecimal(value *big.Int) float64 {
	bigFloatValue := math.BigIntToBigFloat(value)
	result := math.DivFloat(bigFloatValue, big.NewFloat(1e18))

	floatValue, _ := result.Float64()
	return floatValue
}
