package testutils

import (
	"encoding/json"
	"testing"

	"github.com/Proofsuite/amp-matching-engine/types"
	"github.com/go-test/deep"
	"github.com/stretchr/testify/assert"
	"gopkg.in/mgo.v2/dbtest"
)

func CompareEngineResponse(t *testing.T, a, b *types.EngineResponse) {
	assert.Equal(t, a.Status, b.Status)

	if a.Order != nil && b.Order != nil {
		assert.NotNil(t, a.Order)
		assert.NotNil(t, b.Order)
		CompareOrder(t, a.Order, b.Order)
	}

	if a.Matches != nil && b.Matches != nil {
		assert.NotNil(t, a.Matches)
		assert.NotNil(t, b.Matches)
		CompareMatches(t, a.Matches, b.Matches)
	}
}

func CompareMatches(t *testing.T, a, b *types.Matches) {
	if a != nil && b != nil {
		assert.NotNil(t, a)
		assert.NotNil(t, b)

		for i, _ := range a.MakerOrders {
			ComparePublicOrder(t, a.MakerOrders[i], b.MakerOrders[i])
		}

		ComparePublicOrder(t, a.TakerOrder, b.TakerOrder)

		for i, _ := range a.Trades {
			ComparePublicTrade(t, a.Trades[i], b.Trades[i])
		}
	}
}

func ComparePublicOrder(t *testing.T, a, b *types.Order) {
	assert.Equal(t, a.UserAddress, b.UserAddress)
	assert.Equal(t, a.ExchangeAddress, b.ExchangeAddress)
	assert.Equal(t, a.PricePoint, b.PricePoint)
	assert.Equal(t, a.Amount, b.Amount)
	assert.Equal(t, a.FilledAmount, b.FilledAmount)
	assert.Equal(t, a.Status, b.Status)
	assert.Equal(t, a.Side, b.Side)
	assert.Equal(t, a.PairName, b.PairName)
	assert.Equal(t, a.MakeFee, b.MakeFee)
	assert.Equal(t, a.Nonce, b.Nonce)
	assert.Equal(t, a.TakeFee, b.TakeFee)
	assert.Equal(t, a.Signature, b.Signature)
	assert.Equal(t, a.Hash, b.Hash)
}

func CompareOrder(t *testing.T, a, b *types.Order) {
	assert.Equal(t, a.ID, b.ID)
	assert.Equal(t, a.UserAddress, b.UserAddress)
	assert.Equal(t, a.ExchangeAddress, b.ExchangeAddress)
	assert.Equal(t, a.PricePoint, b.PricePoint)
	assert.Equal(t, a.Amount, b.Amount)
	assert.Equal(t, a.FilledAmount, b.FilledAmount)
	assert.Equal(t, a.Status, b.Status)
	assert.Equal(t, a.Side, b.Side)
	assert.Equal(t, a.PairName, b.PairName)
	assert.Equal(t, a.MakeFee, b.MakeFee)
	assert.Equal(t, a.Nonce, b.Nonce)
	assert.Equal(t, a.TakeFee, b.TakeFee)
	assert.Equal(t, a.Signature, b.Signature)
	assert.Equal(t, a.Hash, b.Hash)
}

func ComparePair(t *testing.T, a, b *types.Pair) {
	assert.Equal(t, a.ID, b.ID)
	assert.Equal(t, a.BaseTokenSymbol, b.BaseTokenSymbol)
	assert.Equal(t, a.BaseTokenAddress, b.BaseTokenAddress)
	assert.Equal(t, a.QuoteTokenSymbol, b.QuoteTokenSymbol)
	assert.Equal(t, a.QuoteTokenAddress, b.QuoteTokenAddress)
	assert.Equal(t, a.Active, b.Active)
	assert.Equal(t, a.MakeFee, b.MakeFee)
	assert.Equal(t, a.TakeFee, b.TakeFee)
}

func CompareToken(t *testing.T, a, b *types.Token) {
	assert.Equal(t, a.Symbol, b.Symbol)
	assert.Equal(t, a.Address, b.Address)
	assert.Equal(t, a.Active, b.Active)
	assert.Equal(t, a.Quote, b.Quote)
	assert.Equal(t, a.ID, b.ID)
}

func CompareTrade(t *testing.T, a, b *types.Trade) {
	assert.Equal(t, a.ID, b.ID)
	assert.Equal(t, a.Maker, b.Maker)
	assert.Equal(t, a.Taker, b.Taker)
	assert.Equal(t, a.BaseToken, b.BaseToken)
	assert.Equal(t, a.QuoteToken, b.QuoteToken)
	assert.Equal(t, a.MakerOrderHash, b.MakerOrderHash)
	assert.Equal(t, a.TakerOrderHash, b.TakerOrderHash)
	assert.Equal(t, a.Hash, b.Hash)
	assert.Equal(t, a.PricePoint, b.PricePoint)
	assert.Equal(t, a.PairName, b.PairName)
	assert.Equal(t, a.TxHash, b.TxHash)
	assert.Equal(t, a.Amount, b.Amount)
}

func ComparePublicTrade(t *testing.T, a, b *types.Trade) {
	assert.Equal(t, a.Maker, b.Maker)
	assert.Equal(t, a.Taker, b.Taker)
	assert.Equal(t, a.BaseToken, b.BaseToken)
	assert.Equal(t, a.QuoteToken, b.QuoteToken)
	assert.Equal(t, a.MakerOrderHash, b.MakerOrderHash)
	assert.Equal(t, a.TakerOrderHash, b.TakerOrderHash)
	assert.Equal(t, a.Hash, b.Hash)
	assert.Equal(t, a.PricePoint, b.PricePoint)
	assert.Equal(t, a.PairName, b.PairName)
	assert.Equal(t, a.TxHash, b.TxHash)
	assert.Equal(t, a.Amount, b.Amount)
}

func CompareAccount(t *testing.T, a, b *types.Account) {
	assert.Equal(t, a.Address, b.Address)
	assert.Equal(t, a.TokenBalances, b.TokenBalances)
	assert.Equal(t, a.IsBlocked, b.IsBlocked)
}

func CompareAccountStrict(t *testing.T, a, b *types.Account) {
	assert.Equal(t, a.ID, b.ID)
	assert.Equal(t, a.Address, b.Address)
	assert.Equal(t, a.TokenBalances, b.TokenBalances)
	assert.Equal(t, a.IsBlocked, b.IsBlocked)
	assert.Equal(t, a.UpdatedAt, b.UpdatedAt)
	assert.Equal(t, a.CreatedAt, b.CreatedAt)
}

func Compare(t *testing.T, expected interface{}, value interface{}) {
	expectedBytes, _ := json.Marshal(expected)
	bytes, _ := json.Marshal(value)

	assert.JSONEqf(t, string(expectedBytes), string(bytes), "")
}

func CompareStructs(t *testing.T, expected interface{}, order interface{}) {
	diff := deep.Equal(expected, order)
	if diff != nil {
		t.Errorf("\n%+v\nGot: \n%+v\n\n", expected, order)
	}
}

func NewDBTestServer() *dbtest.DBServer {
	return &dbtest.DBServer{}
}
