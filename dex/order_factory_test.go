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
		Nonce:           o.Nonce,
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

func TestNewFactoryBuyOrder(t *testing.T) {
	wallet := testConfig.Wallets[0]
	pair := testConfig.TokenPairs["ZRXWETH"]
	ZRX := pair.BaseToken
	WETH := pair.QuoteToken

	f := NewOrderFactory(&pair, wallet)

	o, err := f.NewBuyOrder(50, 2) //Selling 2ZRX at the price of 2 ZRX = 100 WETH
	if err != nil {
		t.Errorf("Error creating new order: %v", err)
	}

	err = pair.ComputeOrderPrice(o)
	if err != nil {
		t.Errorf("Error getting new order price: %v", err)
	}

	expected := &Order{
		Id:              0,
		OrderType:       BUY,
		ExchangeAddress: testConfig.Exchange,
		Maker:           wallet.Address,
		TokenBuy:        ZRX.Address,
		TokenSell:       WETH.Address,
		SymbolBuy:       ZRX.Symbol,
		SymbolSell:      WETH.Symbol,
		AmountBuy:       big.NewInt(2),
		AmountSell:      big.NewInt(100),
		Expires:         big.NewInt(1e18),
		Nonce:           o.Nonce,
		FeeMake:         big.NewInt(0),
		FeeTake:         big.NewInt(0),
		PairID:          pair.ID,
		Signature:       o.Signature,
		Hash:            o.Hash,
		Price:           50 * 1e3, //1e3 is the current multiplier used in the compute price function
		Amount:          2,
	}

	if !reflect.DeepEqual(o, expected) {
		t.Errorf("Expected order to be equal to %v but got %v instead", expected, o)
	}
}

func TestNewFactorySellOrder(t *testing.T) {
	wallet := testConfig.Wallets[0]
	pair := testConfig.TokenPairs["ZRXWETH"]
	ZRX := pair.BaseToken
	WETH := pair.QuoteToken

	f := NewOrderFactory(&pair, wallet)

	o, err := f.NewSellOrder(100, 1) //Selling 1 ZRX at the price of 1 ZRX = 100 WETH
	if err != nil {
		t.Errorf("Error creating new order: %v", err)
	}

	err = pair.ComputeOrderPrice(o)
	if err != nil {
		t.Errorf("Error getting new order price: %v", err)
	}

	expected := &Order{
		Id:              0,
		OrderType:       SELL,
		ExchangeAddress: testConfig.Exchange,
		Maker:           wallet.Address,
		TokenBuy:        WETH.Address,
		TokenSell:       ZRX.Address,
		SymbolBuy:       WETH.Symbol,
		SymbolSell:      ZRX.Symbol,
		AmountBuy:       big.NewInt(100),
		AmountSell:      big.NewInt(1),
		Expires:         big.NewInt(1e18),
		Nonce:           o.Nonce,
		FeeMake:         big.NewInt(0),
		FeeTake:         big.NewInt(0),
		PairID:          pair.ID,
		Signature:       o.Signature,
		Hash:            o.Hash,
		Price:           100 * 1e3, //1e3 is the current multiplier used in the compute price function
		Amount:          1,
	}

	if !reflect.DeepEqual(o, expected) {
		t.Errorf("Expected order to be equal to %v but got %v instead", expected, o)
	}
}

func TestNewFactorySellOrder2(t *testing.T) {
	wallet := testConfig.Wallets[0]
	pair := testConfig.TokenPairs["ZRXWETH"]
	ZRX := pair.BaseToken
	WETH := pair.QuoteToken

	f := NewOrderFactory(&pair, wallet)

	o, err := f.NewSellOrder(250, 10) //Selling 10 ZRX at the price of 1 ZRX = 250 WETH
	if err != nil {
		t.Errorf("Error creating new order: %v", err)
	}

	err = pair.ComputeOrderPrice(o)
	if err != nil {
		t.Errorf("Error getting new order price: %v", err)
	}

	expected := &Order{
		Id:              0,
		OrderType:       SELL,
		ExchangeAddress: testConfig.Exchange,
		Maker:           wallet.Address,
		TokenBuy:        WETH.Address,
		TokenSell:       ZRX.Address,
		SymbolBuy:       WETH.Symbol,
		SymbolSell:      ZRX.Symbol,
		AmountBuy:       big.NewInt(2500),
		AmountSell:      big.NewInt(10),
		Expires:         big.NewInt(1e18),
		Nonce:           o.Nonce,
		FeeMake:         big.NewInt(0),
		FeeTake:         big.NewInt(0),
		PairID:          pair.ID,
		Signature:       o.Signature,
		Hash:            o.Hash,
		Price:           250 * 1e3, //1e3 is the current multiplier used in the compute price function
		Amount:          10,
	}

	if !reflect.DeepEqual(o, expected) {
		t.Errorf("Expected order to be equal to %v but got %v instead", expected, o)
	}
}
