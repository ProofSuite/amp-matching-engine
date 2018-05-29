package dex

import (
	"math/big"
	"reflect"
	"testing"
)

var testConfig = NewDefaultConfiguration()

func TestNewOrderFromFactory(t *testing.T) {
	wallet := testConfig.Wallets[0]
	pair := testConfig.TokenPairs["ZRXWETH"]
	ZRX := pair.BaseToken
	WETH := pair.QuoteToken

	f := NewOrderFactory(&pair, wallet)

	o, err := f.NewOrder(ZRX, 1, WETH, 1)
	if err != nil {
		t.Errorf("Error creating new order: %v", err)
	}

	expected := &Order{
		Id:              0,
		ExchangeAddress: testConfig.Exchange,
		Maker:           wallet.Address,
		TokenBuy:        ZRX.Address,
		TokenSell:       WETH.Address,
		SymbolBuy:       ZRX.Symbol,
		SymbolSell:      WETH.Symbol,
		AmountBuy:       big.NewInt(1),
		AmountSell:      big.NewInt(1),
		Expires:         big.NewInt(1e18),
		Nonce:           big.NewInt(0),
		FeeMake:         big.NewInt(0),
		FeeTake:         big.NewInt(0),
		PairID:          pair.ID,
		Signature:       o.Signature,
		Hash:            o.Hash,
	}

	if !reflect.DeepEqual(o, expected) {
		t.Errorf("Expected order to be equal to %v but got %v instead", expected, o)
	}

	if f.CurrentOrderID != 1 {
		t.Errorf("Current Order ID should be equal to 1 but got %v instead", f.CurrentOrderID)
	}
}
