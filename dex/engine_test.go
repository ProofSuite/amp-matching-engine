package dex

import (
	"math/big"
	"reflect"
	"testing"

	"github.com/ethereum/go-ethereum/common"
)

func TestRegisterNewQuoteToken(t *testing.T) {
	quoteToken := testConfig.QuoteTokens["WETH"]

	engine := NewTradingEngine()
	err := engine.RegisterNewQuoteToken(quoteToken)
	if err != nil {
		t.Errorf("Error registering quote token: %v", err)
	}

	if engine.quoteTokens[quoteToken.Address] != quoteToken {
		t.Errorf("Quote token not registered properly")
	}
}

func TestRegisterNewPair(t *testing.T) {
	pair := testConfig.TokenPairs["ZRXWETH"]
	done := make(chan bool)

	engine := NewTradingEngine()
	err := engine.RegisterNewPair(pair, done)
	if err.Error() != "Quote token is not registered" {
		t.Errorf("Expected errror 'Quote token is not registered'")
	}
}

func TestRegisterNewPair2(t *testing.T) {
	quoteToken := testConfig.QuoteTokens["WETH"]
	pair := testConfig.TokenPairs["ZRXWETH"]
	done := make(chan bool)

	engine := NewTradingEngine()
	engine.RegisterNewQuoteToken(quoteToken)
	err := engine.RegisterNewPair(pair, done)
	if err != nil {
		t.Errorf("Could not register token pair: %v", pair)
	}
}

func TestComputeOrderPrice(t *testing.T) {
	quoteToken := testConfig.QuoteTokens["WETH"]
	pair := testConfig.TokenPairs["ZRXWETH"]
	done := make(chan bool)

	engine := NewTradingEngine()
	engine.RegisterNewQuoteToken(quoteToken)
	engine.RegisterNewPair(pair, done)

	order := &Order{
		Id:              0,
		ExchangeAddress: common.HexToAddress("0xae55690d4b079460e6ac28aaa58c9ec7b73a7485"),
		Maker:           common.HexToAddress("0xc9b32e9563fe99612ce3a2695ac2a6404c111dde"),
		TokenBuy:        testConfig.TokenPairs["ZRXWETH"].QuoteToken.Address,
		TokenSell:       testConfig.TokenPairs["ZRXWETH"].BaseToken.Address,
		SymbolBuy:       testConfig.TokenPairs["ZRXWETH"].QuoteToken.Symbol,
		SymbolSell:      testConfig.TokenPairs["ZRXWETH"].BaseToken.Symbol,
		AmountBuy:       big.NewInt(1000),
		AmountSell:      big.NewInt(100),
		Expires:         big.NewInt(10000),
		Nonce:           big.NewInt(1000),
		FeeMake:         big.NewInt(50),
		FeeTake:         big.NewInt(50),
		Signature: &Signature{
			V: 28,
			R: common.HexToHash("0x10b30eb0072a4f0a38b6fca0b731cba15eb2e1702845d97c1230b53a839bcb85"),
			S: common.HexToHash("0x6d9ad89548c9e3ce4c97825d027291477f2c44a8caef792095f2cabc978493ff"),
		},
		PairID: pair.ID,
		Hash:   common.HexToHash("0xb9070a2d333403c255ce71ddf6e795053599b2e885321de40353832b96d8880a"),
	}

	err := engine.ComputeOrderPrice(order)
	if err != nil {
		t.Errorf("Error computing order price: %v", err)
	}
}

func TestEngine(t *testing.T) {
	wallet := testConfig.Wallets[0]
	pair := testConfig.TokenPairs["ZRXWETH"]
	ZRX := pair.BaseToken
	WETH := pair.QuoteToken
	done := make(chan bool)

	engine := NewTradingEngine()
	engine.RegisterNewQuoteToken(WETH)
	engine.RegisterNewPair(pair, done)

	factory := NewOrderFactory(&pair, wallet)

	o1, _ := factory.NewOrder(ZRX, 1, WETH, 1)
	o2, _ := factory.NewOrder(WETH, 1, ZRX, 1)

	engine.AddOrder(o1)
	engine.AddOrder(o2)

	logs := engine.orderbooks[pair].logger

	expected := []*Action{
		&Action{actionType: AT_BUY, amount: 1, price: 1e3, orderHash: o1.Hash},
		&Action{actionType: AT_SELL, amount: 1, price: 1e3, orderHash: o2.Hash},
		&Action{actionType: AT_FILLED, amount: 1, price: 1e3, orderHash: o2.Hash, fromOrderHash: o1.Hash},
	}

	if !reflect.DeepEqual(logs, expected) {
		t.Error("\n\nExpected:\n\n", expected, "\n\nGot:\n\n", logs, "\n\n")
	}
}

func TestEngine2(t *testing.T) {
	wallet := testConfig.Wallets[1]
	pair := testConfig.TokenPairs["ZRXWETH"]
	ZRX := pair.BaseToken
	WETH := pair.QuoteToken
	done := make(chan bool)

	engine := NewTradingEngine()
	engine.RegisterNewQuoteToken(WETH)
	engine.RegisterNewPair(pair, done)

	factory := NewOrderFactory(&pair, wallet)
	o1, _ := factory.NewOrder(ZRX, 1, WETH, 1)
	oc1, _ := factory.NewOrderCancel(o1)
	o2, _ := factory.NewOrder(WETH, 1, ZRX, 1)

	engine.AddOrder(o1)
	engine.CancelOrder(oc1)
	engine.AddOrder(o2)
	engine.CloseOrderBook(pair.ID)

	<-done

	expected := []*Action{
		&Action{actionType: AT_BUY, amount: 1, price: 1e3, orderHash: o1.Hash},
		&Action{actionType: AT_CANCEL, orderHash: o1.Hash},
		&Action{actionType: AT_SELL, amount: 1, price: 1e3, orderHash: o2.Hash},
		&Action{actionType: AT_DONE},
	}

	logs := engine.orderbooks[pair].logger

	if !reflect.DeepEqual(logs, expected) {
		t.Error("\n\nExpected:\n\n", expected, "\n\nGot:\n\n", logs, "\n\n")
	}
}

func TestEngine3(t *testing.T) {
	wallet := testConfig.Wallets[0]
	pair := testConfig.TokenPairs["ZRXWETH"]
	ZRX := pair.BaseToken
	WETH := pair.QuoteToken
	done := make(chan bool)

	engine := NewTradingEngine()
	engine.RegisterNewQuoteToken(WETH)
	engine.RegisterNewPair(pair, done)

	factory := NewOrderFactory(&pair, wallet)
	o1, _ := factory.NewOrder(ZRX, 1, WETH, 10)
	o2, _ := factory.NewOrder(WETH, 11, ZRX, 1)

	engine.AddOrder(o1)
	engine.AddOrder(o2)
	engine.CloseOrderBook(pair.ID)

	<-done

	expected := []*Action{
		&Action{actionType: AT_BUY, amount: 1, price: 10 * 1e3, orderHash: o1.Hash},
		&Action{actionType: AT_SELL, amount: 1, price: 11 * 1e3, orderHash: o2.Hash},
		&Action{actionType: AT_DONE},
	}

	logs := engine.orderbooks[pair].logger

	if !reflect.DeepEqual(logs, expected) {
		t.Error("\n\nExpected:\n\n", expected, "\n\nGot:\n\n", logs, "\n\n")
	}
}

func TestEngine4(t *testing.T) {
	wallet := testConfig.Wallets[0]
	pair := testConfig.TokenPairs["ZRXWETH"]
	ZRX := pair.BaseToken
	WETH := pair.QuoteToken
	done := make(chan bool)

	engine := NewTradingEngine()
	engine.RegisterNewQuoteToken(WETH)
	engine.RegisterNewPair(pair, done)

	factory := NewOrderFactory(&pair, wallet)
	o1, _ := factory.NewOrder(ZRX, 1, WETH, 11)
	o2, _ := factory.NewOrder(WETH, 10, ZRX, 1)

	engine.AddOrder(o1)
	engine.AddOrder(o2)
	engine.CloseOrderBook(pair.ID)

	<-done

	expected := []*Action{
		&Action{actionType: AT_BUY, amount: 1, price: 11 * 1e3, orderHash: o1.Hash},
		&Action{actionType: AT_SELL, amount: 1, price: 10 * 1e3, orderHash: o2.Hash},
		&Action{actionType: AT_FILLED, amount: 1, price: 11 * 1e3, orderHash: o2.Hash, fromOrderHash: o1.Hash},
		&Action{actionType: AT_DONE},
	}

	logs := engine.orderbooks[pair].logger

	if !reflect.DeepEqual(logs, expected) {
		t.Error("\n\nExpected:\n\n", expected, "\n\nGot:\n\n", logs, "\n\n")
	}
}

func TestEngine5(t *testing.T) {
	wallet := testConfig.Wallets[0]
	pair := testConfig.TokenPairs["ZRXWETH"]
	WETH := pair.QuoteToken // ZRX is the base token and WETH is the quote token
	done := make(chan bool)

	engine := NewTradingEngine()
	engine.RegisterNewQuoteToken(WETH)
	engine.RegisterNewPair(pair, done)

	factory := NewOrderFactory(&pair, wallet)
	o1, _ := factory.NewSellOrder(50, 50)
	o2, _ := factory.NewSellOrder(45, 25)
	o3, _ := factory.NewSellOrder(45, 25)
	o4, _ := factory.NewBuyOrder(55, 75)
	oc5, _ := factory.NewOrderCancel(o1)
	o6, _ := factory.NewBuyOrder(55, 20)
	o7, _ := factory.NewBuyOrder(50, 15)
	o8, _ := factory.NewSellOrder(45, 25)

	engine.AddOrder(o1)
	engine.AddOrder(o2)
	engine.AddOrder(o3)
	engine.AddOrder(o4)
	engine.CancelOrder(oc5)
	engine.AddOrder(o6)
	engine.AddOrder(o7)
	engine.AddOrder(o8)
	engine.CloseOrderBook(pair.ID)

	<-done

	expected := []*Action{
		&Action{actionType: AT_SELL, price: 50 * 1e3, amount: 50, orderHash: o1.Hash},
		&Action{actionType: AT_SELL, price: 45 * 1e3, amount: 25, orderHash: o2.Hash},
		&Action{actionType: AT_SELL, price: 45 * 1e3, amount: 25, orderHash: o3.Hash},
		&Action{actionType: AT_BUY, price: 55 * 1e3, amount: 75, orderHash: o4.Hash},
		&Action{actionType: AT_PARTIAL_FILLED, price: 45 * 1e3, amount: 25, orderHash: o4.Hash, fromOrderHash: o2.Hash},
		&Action{actionType: AT_PARTIAL_FILLED, price: 45 * 1e3, amount: 25, orderHash: o4.Hash, fromOrderHash: o3.Hash},
		&Action{actionType: AT_FILLED, price: 50 * 1e3, amount: 25, orderHash: o4.Hash, fromOrderHash: o1.Hash},
		&Action{actionType: AT_CANCEL, orderHash: o1.Hash},
		&Action{actionType: AT_BUY, price: 55 * 1e3, amount: 20, orderHash: o6.Hash},
		&Action{actionType: AT_BUY, price: 50 * 1e3, amount: 15, orderHash: o7.Hash},
		&Action{actionType: AT_SELL, price: 45 * 1e3, amount: 25, orderHash: o8.Hash},
		&Action{actionType: AT_PARTIAL_FILLED, price: 55 * 1e3, amount: 20, orderHash: o8.Hash, fromOrderHash: o6.Hash},
		&Action{actionType: AT_FILLED, price: 50 * 1e3, amount: 5, orderHash: o8.Hash, fromOrderHash: o7.Hash},
		&Action{actionType: AT_DONE},
	}

	logs := engine.orderbooks[pair].logger

	if !reflect.DeepEqual(logs, expected) {
		t.Error("\n\nExpected:\n\n", expected, "\n\nGot:\n\n", logs, "\n\n")
	}
}
