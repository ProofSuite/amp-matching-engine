package dex

import (
	"os"
	"reflect"
	"testing"
	"time"
)

func TestSocketBuyOrder(t *testing.T) {
	os.Stdout, _ = os.Open(os.DevNull)

	wallet := testConfig.Wallets[1]
	quotes := testConfig.QuoteTokens
	pairs := testConfig.TokenPairs
	ZRXWETH := pairs["ZRXWETH"]
	ZRX := ZRXWETH.BaseToken
	WETH := ZRXWETH.QuoteToken
	done := make(chan bool)

	factory := NewOrderFactory(&ZRXWETH, wallet)

	server := NewServer()
	server.SetupCurrencies(quotes, pairs, done)

	client := NewClient(wallet, server)
	client.start()

	m1, o1, _ := factory.NewOrderMessage(ZRX, 1, WETH, 1)
	client.requests <- m1

	time.Sleep(time.Second)

	expected := []*Action{
		&Action{actionType: AT_BUY, orderHash: o1.Hash, amount: 1, price: 1e3},
	}

	logs := server.engine.orderbooks[ZRXWETH].logger

	if !reflect.DeepEqual(logs, expected) {
		t.Error("\n\nExpected:\n\n", expected, "\n\nGot:\n\n", logs, "\n\n")
	}
}

func TestSocketOrderFill(t *testing.T) {
	os.Stdout, _ = os.Open(os.DevNull)

	quotes := testConfig.QuoteTokens
	pairs := testConfig.TokenPairs
	ZRXWETH := pairs["ZRXWETH"]
	ZRX := ZRXWETH.BaseToken
	WETH := ZRXWETH.QuoteToken
	done := make(chan bool)

	wallet1 := testConfig.Wallets[1]
	wallet2 := testConfig.Wallets[2]

	factory1 := NewOrderFactory(&ZRXWETH, wallet1)
	factory2 := NewOrderFactory(&ZRXWETH, wallet2)

	server := NewServer()
	server.SetupCurrencies(quotes, pairs, done)

	client1 := NewClient(wallet1, server)
	client1.start()

	client2 := NewClient(wallet2, server)
	client2.start()

	m1, o1, _ := factory1.NewOrderMessage(ZRX, 1, WETH, 1)
	m2, o2, _ := factory2.NewOrderMessage(WETH, 1, ZRX, 1)

	client1.requests <- m1
	time.Sleep(time.Millisecond)
	client2.requests <- m2

	time.Sleep(time.Second)

	logs := server.engine.orderbooks[ZRXWETH].logger

	expected := []*Action{
		&Action{actionType: AT_BUY, orderHash: o1.Hash, amount: 1, price: 1e3},
		&Action{actionType: AT_SELL, orderHash: o2.Hash, amount: 1, price: 1e3},
		&Action{actionType: AT_FILLED, orderHash: o2.Hash, fromOrderHash: o1.Hash, amount: 1, price: 1e3},
	}

	if !reflect.DeepEqual(logs, expected) {
		t.Error("\n\nExpected:\n\n", expected, "\n\nGot:\n\n", logs, "\n\n")
	}
}

func TestSocketOrderCancel(t *testing.T) {
	quotes := testConfig.QuoteTokens
	pairs := testConfig.TokenPairs
	ZRXWETH := pairs["ZRXWETH"]
	ZRX := ZRXWETH.BaseToken
	WETH := ZRXWETH.QuoteToken
	done := make(chan bool)

	wallet1 := testConfig.Wallets[1]
	wallet2 := testConfig.Wallets[2]

	factory1 := NewOrderFactory(&ZRXWETH, wallet1)
	factory2 := NewOrderFactory(&ZRXWETH, wallet2)

	server := NewServer()
	server.SetupCurrencies(quotes, pairs, done)

	client1 := NewClient(wallet1, server)
	client1.start()

	client2 := NewClient(wallet2, server)
	client2.start()

	m1, o1, _ := factory1.NewOrderMessage(ZRX, 1, WETH, 1)
	ocm1, _ := factory1.NewCancelOrderMessage(o1)
	m2, o2, _ := factory2.NewOrderMessage(WETH, 1, ZRX, 1)

	client1.requests <- m1
	time.Sleep(time.Millisecond)
	client1.requests <- ocm1
	time.Sleep(time.Millisecond)
	client2.requests <- m2

	time.Sleep(time.Second)

	logs := server.engine.orderbooks[ZRXWETH].logger

	expected := []*Action{
		&Action{actionType: AT_BUY, amount: 1, price: 1e3, orderHash: o1.Hash},
		&Action{actionType: AT_CANCEL, orderHash: o1.Hash},
		&Action{actionType: AT_SELL, amount: 1, price: 1e3, orderHash: o2.Hash},
	}

	if !reflect.DeepEqual(logs, expected) {
		t.Error("\n\nExpected:\n\n", expected, "\n\nGot:\n\n", logs, "\n\n")
	}
}

func TestSocketExecuteOrder(t *testing.T) {
	testConfig := NewConfiguration()
	server := NewServer()
	done := make(chan bool)

	quotes := testConfig.QuoteTokens
	pairs := testConfig.TokenPairs
	ZRXWETH := pairs["ZRXWETH"]
	ZRX := ZRXWETH.BaseToken
	WETH := ZRXWETH.QuoteToken

	wallet1 := testConfig.Wallets[1]
	wallet2 := testConfig.Wallets[2]

	factory1 := NewOrderFactory(&ZRXWETH, wallet1)
	factory2 := NewOrderFactory(&ZRXWETH, wallet2)

	config := &OperatorConfig{
		Admin:          testConfig.Admin,
		Exchange:       testConfig.Exchange,
		OperatorParams: testConfig.OperatorParams,
	}

	server.SetupTradingEngine(config, quotes, pairs, done)

	client1 := NewClient(wallet1, server)
	client1.start()

	client2 := NewClient(wallet2, server)
	client2.start()

	m1, _, _ := factory1.NewOrderMessage(ZRX, 1, WETH, 1)
	m2, _, _ := factory2.NewOrderMessage(WETH, 1, ZRX, 1)

	client1.requests <- m1
	time.Sleep(time.Millisecond)
	client2.requests <- m2

	time.Sleep(time.Second)

	errChan := server.engine.operator.ErrorChannel
	tradeChan := server.engine.operator.TradeChannel
	cancelOrderChan := server.engine.operator.CancelOrderChannel

	for {
		select {
		case log := <-errChan:
			t.Errorf("Received error log: %v", PrintErrorLog(log))
		case log := <-tradeChan:
			t.Errorf("Received trade log: %v", PrintTradeLog(log))
		case log := <-cancelOrderChan:
			t.Errorf("Receive cancel trade log: %v", PrintCancelOrderLog(log))
		}
	}

	time.Sleep(10 * time.Second)

	// logs := server.engine.orderbooks[ZRXWETH].logger

	// expected := []*Action{
	// 	&Action{actionType: AT_BUY, orderHash: o1.Hash, amount: 1, price: 1e3},
	// 	&Action{actionType: AT_SELL, orderHash: o2.Hash, amount: 1, price: 1e3},
	// 	&Action{actionType: AT_FILLED, orderHash: o2.Hash, fromOrderHash: o1.Hash, amount: 1, price: 1e3},
	// }

	// if !reflect.DeepEqual(logs, expected) {
	// 	t.Error("\n\nExpected:\n\n", expected, "\n\nGot:\n\n", logs, "\n\n")
	// }
}
