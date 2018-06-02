package dex

import (
	"math/big"
	"os"
	"reflect"
	"sync"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
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
	var wg sync.WaitGroup
	wg.Add(2)
	testConfig := NewConfiguration()
	server := NewServer()
	done := make(chan bool)

	quotes := testConfig.QuoteTokens
	pairs := testConfig.TokenPairs
	ZRXWETH := pairs["ZRXWETH"]
	ZRX := ZRXWETH.BaseToken
	WETH := ZRXWETH.QuoteToken

	maker := testConfig.Wallets[1]
	taker := testConfig.Wallets[2]

	makerFactory := NewOrderFactory(&ZRXWETH, maker)
	takerFactory := NewOrderFactory(&ZRXWETH, taker)

	config := &OperatorConfig{
		Admin:          testConfig.Admin,
		Exchange:       testConfig.Exchange,
		OperatorParams: testConfig.OperatorParams,
	}

	err := server.SetupTradingEngine(config, quotes, pairs, done)
	if err != nil {
		t.Errorf("Error during trading engine setup : %v", err)
	}

	makerClient := NewClient(maker, server)
	makerClient.start()

	takerClient := NewClient(taker, server)
	takerClient.start()

	m1, _, _ := makerFactory.NewOrderMessage(ZRX, 1, WETH, 1)
	m2, _, _ := takerFactory.NewOrderMessage(WETH, 1, ZRX, 1)

	initialTakerZRXBalance, _ := server.engine.TokenBalance(taker.Address, ZRX.Address)
	initialTakerWETHBalance, _ := server.engine.TokenBalance(taker.Address, WETH.Address)
	initialMakerZRXBalance, _ := server.engine.TokenBalance(maker.Address, ZRX.Address)
	initialMakerWETHBalance, _ := server.engine.TokenBalance(maker.Address, WETH.Address)

	makerClient.requests <- m1
	time.Sleep(time.Millisecond)
	takerClient.requests <- m2
	time.Sleep(time.Millisecond)

	var takerTx, makerTx common.Hash

	go func() {
		for {
			select {
			case l := <-makerClient.logs:
				switch l.MessageType {
				case ORDER_EXECUTED:
					takerTx = l.Tx
				case ORDER_TX_SUCCESS:
					if l.Tx != takerTx {
						t.Errorf("Tx hash is not matching")
					}
					wg.Done()
				case ORDER_TX_ERROR:
					t.Errorf("Transaction failed: The error id is %v", l.ErrorID)
				}
			case l := <-takerClient.logs:
				switch l.MessageType {
				case TRADE_EXECUTED:
					makerTx = l.Tx
				case TRADE_TX_SUCCESS:
					if l.Tx != makerTx {
						t.Errorf("Tx hash is not matching")
					}
					wg.Done()
				case ORDER_TX_ERROR:
					t.Errorf("Transaction failed: The error id is %v", l.ErrorID)
				}
			}
		}
	}()

	wg.Wait()

	if makerTx != takerTx {
		t.Errorf("The maker transaction hash and the taker transaction hash are different")
	}

	TakerZRXBalance, _ := server.engine.TokenBalance(taker.Address, ZRX.Address)
	TakerWETHBalance, _ := server.engine.TokenBalance(taker.Address, WETH.Address)
	MakerZRXBalance, _ := server.engine.TokenBalance(maker.Address, ZRX.Address)
	MakerWETHBalance, _ := server.engine.TokenBalance(maker.Address, WETH.Address)

	TakerZRXIncrement := big.NewInt(0)
	TakerWETHIncrement := big.NewInt(0)
	MakerZRXIncrement := big.NewInt(0)
	MakerWETHIncrement := big.NewInt(0)

	MakerZRXIncrement.Sub(MakerZRXBalance, initialMakerZRXBalance)
	MakerWETHIncrement.Sub(MakerWETHBalance, initialMakerWETHBalance)
	TakerZRXIncrement.Sub(TakerZRXBalance, initialTakerZRXBalance)
	TakerWETHIncrement.Sub(TakerWETHBalance, initialTakerWETHBalance)

	if MakerZRXIncrement.Cmp(big.NewInt(1)) != 0 {
		t.Errorf("Expected Maker Balance to be equal to be %v but got %v instead", -1, MakerZRXIncrement)
	}

	if MakerWETHIncrement.Cmp(big.NewInt(-1)) != 0 {
		t.Errorf("Expected Taker Balance to be equal to be %v but got %v instead", 1, MakerWETHIncrement)
	}

	if TakerWETHIncrement.Cmp(big.NewInt(1)) != 0 {
		t.Errorf("Expected Maker Balance to be equal to %v but got %v instead", -1, TakerWETHIncrement)
	}

	if TakerZRXIncrement.Cmp(big.NewInt(-1)) != 0 {
		t.Errorf("Expected Maker Balance to be equal to %v but got %v instead", 1, TakerZRXIncrement)
	}
}
