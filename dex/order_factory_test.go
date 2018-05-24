package dex

import (
	"math/big"
	"reflect"
	"testing"
)

func TestNewOrderFromFactory(t *testing.T) {
	wallet := config.Wallets[0]
	pair := config.TokenPairs["ZRXWETH"]
	ZRX := pair.BaseToken
	WETH := pair.QuoteToken

	f := NewOrderFactory(&pair, wallet)

	o, err := f.NewOrder(ZRX, 1, WETH, 1)
	if err != nil {
		t.Errorf("Error creating new order: %v", err)
	}

	expected := &Order{
		Id:              0,
		ExchangeAddress: config.Contracts.exchange,
		Maker:           wallet.Address,
		TokenBuy:        ZRX.Address,
		TokenSell:       WETH.Address,
		SymbolBuy:       ZRX.Symbol,
		SymbolSell:      WETH.Symbol,
		AmountBuy:       big.NewInt(1),
		AmountSell:      big.NewInt(1),
		Expires:         big.NewInt(0),
		Nonce:           big.NewInt(0),
		FeeMake:         big.NewInt(0),
		FeeTake:         big.NewInt(0),
		PairID:          pair.ID,
		Signature:       o.Signature, //TODO Find a better method to test signature/hash
		Hash:            o.Hash,
	}

	if !reflect.DeepEqual(o, expected) {
		t.Errorf("Expected order to be equal to %v but got %v instead", expected, o)
	}

	if f.CurrentOrderID != 1 {
		t.Errorf("Current Order ID should be equal to 1 but got %v instead", f.CurrentOrderID)
	}
}
