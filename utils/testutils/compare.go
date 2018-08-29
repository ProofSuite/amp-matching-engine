package testutils

import (
	"encoding/json"
	"testing"

	"github.com/Proofsuite/amp-matching-engine/types"
	"github.com/go-test/deep"
	"github.com/stretchr/testify/assert"
	"gopkg.in/mgo.v2/dbtest"
)

func CompareOrder(t *testing.T, a, b *types.Order) {
	assert.Equal(t, a.ID, b.ID)
	assert.Equal(t, a.UserAddress, b.UserAddress)
	assert.Equal(t, a.ExchangeAddress, b.ExchangeAddress)
	assert.Equal(t, a.BuyToken, b.BuyToken)
	assert.Equal(t, a.SellToken, b.SellToken)
	assert.Equal(t, a.BaseToken, b.BaseToken)
	assert.Equal(t, a.BuyAmount, b.BuyAmount)
	assert.Equal(t, a.SellAmount, b.SellAmount)
	assert.Equal(t, a.Price, b.Price)
	assert.Equal(t, a.Amount, b.Amount)
	assert.Equal(t, a.FilledAmount, b.FilledAmount)
	assert.Equal(t, a.Status, b.Status)
	assert.Equal(t, a.Side, b.Side)
	assert.Equal(t, a.PairName, b.PairName)
	assert.Equal(t, a.Expires, b.Expires)
	assert.Equal(t, a.MakeFee, b.MakeFee)
	assert.Equal(t, a.Nonce, b.Nonce)
	assert.Equal(t, a.TakeFee, b.TakeFee)
	assert.Equal(t, a.Signature, b.Signature)
	assert.Equal(t, a.Hash, b.Hash)
}

func ComparePair(t *testing.T, a, b *types.Pair) {
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

func CompareToken(t *testing.T, a, b *types.Token) {
	assert.Equal(t, a.Name, b.Name)
	assert.Equal(t, a.Symbol, b.Symbol)
	assert.Equal(t, a.ContractAddress, b.ContractAddress)
	assert.Equal(t, a.Decimal, b.Decimal)
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
	assert.Equal(t, a.OrderHash, b.OrderHash)
	assert.Equal(t, a.Hash, b.Hash)
	assert.Equal(t, a.PairName, b.PairName)
	assert.Equal(t, a.TradeNonce, b.TradeNonce)
	assert.Equal(t, a.Signature, b.Signature)
	assert.Equal(t, a.Tx, b.Tx)

	assert.Equal(t, a.Price, b.Price)
	assert.Equal(t, a.Side, b.Side)
	assert.Equal(t, a.Amount, b.Amount)
}

func CompareAccount(t *testing.T, a, b *types.Account) {
	assert.Equal(t, a.ID, b.ID)
	assert.Equal(t, a.Address, b.Address)
	assert.Equal(t, a.TokenBalances, b.TokenBalances)
	assert.Equal(t, a.IsBlocked, b.IsBlocked)
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
