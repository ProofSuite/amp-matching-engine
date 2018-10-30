package types

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"gopkg.in/mgo.v2/bson"
)

func ComparePair(t *testing.T, a, b *Pair) {
	assert.Equal(t, a.ID, b.ID)
	assert.Equal(t, a.Name, b.Name)
	assert.Equal(t, a.BaseTokenSymbol, b.BaseTokenSymbol)
	assert.Equal(t, a.BaseTokenAddress, b.BaseTokenAddress)
	assert.Equal(t, a.QuoteTokenSymbol, b.QuoteTokenSymbol)
	assert.Equal(t, a.QuoteTokenAddress, b.QuoteTokenAddress)
	assert.Equal(t, a.Active, b.Active)
	assert.Equal(t, a.MakeFee, b.MakeFee)
	assert.Equal(t, a.TakeFee, b.TakeFee)
}

func TestPairBSON(t *testing.T) {
	pair := &Pair{
		ID:                bson.NewObjectId(),
		BaseTokenSymbol:   "REQ",
		BaseTokenAddress:  common.HexToAddress("0xcf7389dc6c63637598402907d5431160ec8972a5"),
		QuoteTokenSymbol:  "WETH",
		QuoteTokenAddress: common.HexToAddress("0x7a9f3cd060ab180f36c17fe6bdf9974f577d77aa"),
		Active:            true,
		MakeFee:           big.NewInt(10000),
		TakeFee:           big.NewInt(10000),
	}

	data, err := bson.Marshal(pair)
	if err != nil {
		t.Errorf("%+v", err)
	}

	decoded := &Pair{}
	if err := bson.Unmarshal(data, decoded); err != nil {
		t.Error(err)
	}

	ComparePair(t, pair, decoded)
}
