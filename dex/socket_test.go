package dex

import (
	"os"
	"reflect"
	"testing"
	"time"
)

func TestSocket(t *testing.T) {
	os.Stdout, _ = os.Open(os.DevNull)

	wallet := config.Wallets[1]
	quotes := config.QuoteTokens
	pairs := config.TokenPairs
	ZRXWETH := pairs["ZRXWETH"]
	ZRX := ZRXWETH.BaseToken
	WETH := ZRXWETH.QuoteToken
	done := make(chan bool)

	factory := NewOrderFactory(&ZRXWETH, wallet)

	server := NewServer()
	server.Setup(quotes, pairs, done)

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

func TestSocket2(t *testing.T) {
	os.Stdout, _ = os.Open(os.DevNull)

	quotes := config.QuoteTokens
	pairs := config.TokenPairs
	ZRXWETH := pairs["ZRXWETH"]
	ZRX := ZRXWETH.BaseToken
	WETH := ZRXWETH.QuoteToken
	done := make(chan bool)

	wallet1 := config.Wallets[1]
	wallet2 := config.Wallets[2]

	factory1 := NewOrderFactory(&ZRXWETH, wallet1)
	factory2 := NewOrderFactory(&ZRXWETH, wallet2)

	server := NewServer()
	server.Setup(quotes, pairs, done)

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
