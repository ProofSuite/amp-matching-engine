package dex

import (
	"math/big"
	"sync"
	"testing"
)

func TestOperator(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(2)
	testConfig := NewConfiguration()
	opParams := testConfig.OperatorParams

	ZRX := testConfig.QuoteTokens["ZRX"]
	WETH := testConfig.QuoteTokens["WETH"]
	ZRXWETH := NewPair(ZRX, WETH)

	admin := config.Wallets[0]
	maker := config.Wallets[1]
	taker := config.Wallets[2]
	exchange := config.Exchange

	opConfig := &OperatorConfig{
		Admin:          admin,
		Exchange:       exchange,
		OperatorParams: opParams,
	}

	ZRXAmount := big.NewInt(1e18)
	WETHAmount := big.NewInt(1e18)

	deployer, err := NewWebsocketDeployer(admin)
	if err != nil {
		t.Errorf("Could not instantiate deployer: %v", err)
	}

	ex, err := deployer.NewExchange(exchange)
	if err != nil {
		t.Errorf("Could not retrieve exchange instance: %v", err)
	}

	makerFactory := NewOrderFactory(&ZRXWETH, maker)
	takerFactory := NewOrderFactory(&ZRXWETH, taker)
	makerFactory.SetExchangeAddress(ex.Address)
	takerFactory.SetExchangeAddress(ex.Address)

	operator, err := NewOperator(opConfig)
	if err != nil {
		t.Errorf("Could not instantiate operator: %v", err)
	}

	initialTakerZRXBalance, _ := ex.TokenBalance(taker.Address, ZRX.Address)
	initialTakerWETHBalance, _ := ex.TokenBalance(taker.Address, WETH.Address)
	initialMakerZRXBalance, _ := ex.TokenBalance(maker.Address, ZRX.Address)
	initialMakerWETHBalance, _ := ex.TokenBalance(maker.Address, WETH.Address)

	order, _ := makerFactory.NewOrderWithEvents(WETH, WETHAmount.Int64(), ZRX, ZRXAmount.Int64())
	trade, _ := takerFactory.NewTradeWithEvents(order, 1)

	err = operator.AddTradeToExecutionList(order, trade)
	if err != nil {
		t.Errorf("Could not execute trade: %v", err)
	}

	//Each received success event sends a done signal to the wait group
	//In total, the wait group (wg) needs to receive 2 different signals
	//corresponding to ORDER_TX_SUCCESS and TRADE_TX_SUCCESS
	go func() {
		for {
			select {
			case e := <-order.events:
				switch e.eventType {
				case ORDER_TX_SUCCESS:
					wg.Done()
				case ORDER_TX_ERROR:
					p := e.payload.(*TxErrorPayload)
					t.Errorf("Received Tx error payload: %v", p)
				}
			case e := <-trade.events:
				switch e.eventType {
				case TRADE_TX_SUCCESS:
					wg.Done()
				case TRADE_TX_ERROR:
					p := e.payload.(*TxErrorPayload)
					t.Errorf("Received Tx error payload: %v", p)
				}
			}
		}
	}()

	wg.Wait()

	TakerZRXBalance, _ := ex.TokenBalance(taker.Address, ZRX.Address)
	TakerWETHBalance, _ := ex.TokenBalance(taker.Address, WETH.Address)
	MakerZRXBalance, _ := ex.TokenBalance(maker.Address, ZRX.Address)
	MakerWETHBalance, _ := ex.TokenBalance(maker.Address, WETH.Address)

	TakerZRXIncrement := big.NewInt(0)
	TakerWETHIncrement := big.NewInt(0)
	MakerZRXIncrement := big.NewInt(0)
	MakerWETHIncrement := big.NewInt(0)

	MakerZRXIncrement.Sub(MakerZRXBalance, initialMakerZRXBalance)
	MakerWETHIncrement.Sub(MakerWETHBalance, initialMakerWETHBalance)
	TakerZRXIncrement.Sub(TakerZRXBalance, initialTakerZRXBalance)
	TakerWETHIncrement.Sub(TakerWETHBalance, initialTakerWETHBalance)

	if MakerZRXIncrement.Cmp(big.NewInt(-1)) != 0 {
		t.Errorf("Expected Maker Balance to be equal to be %v but got %v instead", -1, MakerZRXIncrement)
	}

	if MakerWETHIncrement.Cmp(big.NewInt(1)) != 0 {
		t.Errorf("Expected Taker Balance to be equal to be %v but got %v instead", 1, MakerWETHIncrement)
	}

	if TakerWETHIncrement.Cmp(big.NewInt(-1)) != 0 {
		t.Errorf("Expected Maker Balance to be equal to %v but got %v instead", -1, TakerWETHIncrement)
	}

	if TakerZRXIncrement.Cmp(big.NewInt(1)) != 0 {
		t.Errorf("Expected Maker Balance to be equal to %v but got %v instead", 1, TakerZRXIncrement)
	}
}

func TestOperator2(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(4)
	testConfig := NewConfiguration()
	opParams := testConfig.OperatorParams

	ZRX := testConfig.QuoteTokens["ZRX"]
	WETH := testConfig.QuoteTokens["WETH"]
	ZRXWETH := NewPair(ZRX, WETH)

	admin := config.Wallets[0]

	maker1 := config.Wallets[1]
	taker1 := config.Wallets[2]
	maker2 := config.Wallets[3]
	taker2 := config.Wallets[4]

	exchange := config.Exchange

	opConfig := &OperatorConfig{
		Admin:          admin,
		Exchange:       exchange,
		OperatorParams: opParams,
	}

	ZRXAmount := big.NewInt(1e18)
	WETHAmount := big.NewInt(1e18)

	deployer, err := NewWebsocketDeployer(admin)
	if err != nil {
		t.Errorf("Could not instantiate deployer: %v", err)
	}

	ex, err := deployer.NewExchange(exchange)
	if err != nil {
		t.Errorf("Could not retrieve exchange instance: %v", err)
	}

	maker1Factory := NewOrderFactory(&ZRXWETH, maker1)
	taker1Factory := NewOrderFactory(&ZRXWETH, taker1)
	maker2Factory := NewOrderFactory(&ZRXWETH, maker2)
	taker2Factory := NewOrderFactory(&ZRXWETH, taker2)

	maker1Factory.SetExchangeAddress(ex.Address)
	taker1Factory.SetExchangeAddress(ex.Address)
	maker2Factory.SetExchangeAddress(ex.Address)
	taker2Factory.SetExchangeAddress(ex.Address)

	operator, err := NewOperator(opConfig)
	if err != nil {
		t.Errorf("Could not instantiate operator: %v", err)
	}

	initialTaker1ZRXBalance, _ := ex.TokenBalance(taker1.Address, ZRX.Address)
	initialTaker1WETHBalance, _ := ex.TokenBalance(taker1.Address, WETH.Address)
	initialMaker1ZRXBalance, _ := ex.TokenBalance(maker1.Address, ZRX.Address)
	initialMaker1WETHBalance, _ := ex.TokenBalance(maker1.Address, WETH.Address)

	initialTaker2ZRXBalance, _ := ex.TokenBalance(taker2.Address, ZRX.Address)
	initialTaker2WETHBalance, _ := ex.TokenBalance(taker2.Address, WETH.Address)
	initialMaker2ZRXBalance, _ := ex.TokenBalance(maker2.Address, ZRX.Address)
	initialMaker2WETHBalance, _ := ex.TokenBalance(maker2.Address, WETH.Address)

	o1, _ := maker1Factory.NewOrderWithEvents(WETH, WETHAmount.Int64(), ZRX, ZRXAmount.Int64())
	t1, _ := taker1Factory.NewTradeWithEvents(o1, 1)

	o2, _ := maker2Factory.NewOrderWithEvents(WETH, WETHAmount.Int64(), ZRX, ZRXAmount.Int64())
	t2, _ := taker2Factory.NewTradeWithEvents(o2, 1)

	err = operator.AddTradeToExecutionList(o1, t1)
	err = operator.AddTradeToExecutionList(o2, t2)

	//Each received success event sends a done signal to the wait group
	//In total, the wait group (wg) needs to receive 4 different signals
	//corresponding to 2 ORDER_TX_SUCCESS (o1 and o2) and 2 TRADE_TX_SUCCESS (t1 and t2)
	go func() {
		for {
			select {
			case e := <-o1.events:
				switch e.eventType {
				case ORDER_TX_SUCCESS:
					wg.Done()
				case ORDER_TX_ERROR:
					p := e.payload.(*TxErrorPayload)
					t.Errorf("Received Tx error payload: %v", p)
				}
			case e := <-o2.events:
				switch e.eventType {
				case ORDER_TX_SUCCESS:
					wg.Done()
				case ORDER_TX_ERROR:
					p := e.payload.(*TxErrorPayload)
					t.Errorf("Received Tx error payload: %v", p)
				}
			case e := <-t1.events:
				switch e.eventType {
				case TRADE_TX_SUCCESS:
					wg.Done()
				case TRADE_TX_ERROR:
					p := e.payload.(*TxErrorPayload)
					t.Errorf("Received Tx error payload: %v", p)
				}
			case e := <-t2.events:
				switch e.eventType {
				case TRADE_TX_SUCCESS:
					wg.Done()
				case TRADE_TX_ERROR:
					p := e.payload.(*TxErrorPayload)
					t.Errorf("Received Tx error payload: %v", p)
				}
			}
		}
	}()

	wg.Wait()

	Taker1ZRXBalance, _ := ex.TokenBalance(taker1.Address, ZRX.Address)
	Taker1WETHBalance, _ := ex.TokenBalance(taker1.Address, WETH.Address)
	Maker1ZRXBalance, _ := ex.TokenBalance(maker1.Address, ZRX.Address)
	Maker1WETHBalance, _ := ex.TokenBalance(maker1.Address, WETH.Address)

	Taker2ZRXBalance, _ := ex.TokenBalance(taker2.Address, ZRX.Address)
	Taker2WETHBalance, _ := ex.TokenBalance(taker2.Address, WETH.Address)
	Maker2ZRXBalance, _ := ex.TokenBalance(maker2.Address, ZRX.Address)
	Maker2WETHBalance, _ := ex.TokenBalance(maker2.Address, WETH.Address)

	Taker1ZRXIncrement := big.NewInt(0)
	Taker2ZRXIncrement := big.NewInt(0)
	Taker1WETHIncrement := big.NewInt(0)
	Taker2WETHIncrement := big.NewInt(0)
	Maker1ZRXIncrement := big.NewInt(0)
	Maker2ZRXIncrement := big.NewInt(0)
	Maker1WETHIncrement := big.NewInt(0)
	Maker2WETHIncrement := big.NewInt(0)

	Maker1ZRXIncrement.Sub(Maker1ZRXBalance, initialMaker1ZRXBalance)
	Maker1WETHIncrement.Sub(Maker1WETHBalance, initialMaker1WETHBalance)
	Taker1ZRXIncrement.Sub(Taker1ZRXBalance, initialTaker1ZRXBalance)
	Taker1WETHIncrement.Sub(Taker1WETHBalance, initialTaker1WETHBalance)
	Maker2ZRXIncrement.Sub(Maker2ZRXBalance, initialMaker2ZRXBalance)
	Maker2WETHIncrement.Sub(Maker2WETHBalance, initialMaker2WETHBalance)
	Taker2ZRXIncrement.Sub(Taker2ZRXBalance, initialTaker2ZRXBalance)
	Taker2WETHIncrement.Sub(Taker2WETHBalance, initialTaker2WETHBalance)

	if Maker1ZRXIncrement.Cmp(big.NewInt(-1)) != 0 {
		t.Errorf("Expected Maker Balance to be equal to be %v but got %v instead", -1, Maker1ZRXIncrement)
	}

	if Maker1WETHIncrement.Cmp(big.NewInt(1)) != 0 {
		t.Errorf("Expected Taker Balance to be equal to be %v but got %v instead", 1, Maker1WETHIncrement)
	}

	if Taker1WETHIncrement.Cmp(big.NewInt(-1)) != 0 {
		t.Errorf("Expected Maker Balance to be equal to %v but got %v instead", -1, Taker1WETHIncrement)
	}

	if Taker1ZRXIncrement.Cmp(big.NewInt(1)) != 0 {
		t.Errorf("Expected Maker Balance to be equal to %v but got %v instead", 1, Taker1ZRXIncrement)
	}

	if Maker2ZRXIncrement.Cmp(big.NewInt(-1)) != 0 {
		t.Errorf("Expected Maker Balance to be equal to be %v but got %v instead", -1, Maker2ZRXIncrement)
	}

	if Maker2WETHIncrement.Cmp(big.NewInt(1)) != 0 {
		t.Errorf("Expected Taker Balance to be equal to be %v but got %v instead", 1, Maker2WETHIncrement)
	}

	if Taker2WETHIncrement.Cmp(big.NewInt(-1)) != 0 {
		t.Errorf("Expected Maker Balance to be equal to %v but got %v instead", -1, Taker2WETHIncrement)
	}

	if Taker2ZRXIncrement.Cmp(big.NewInt(1)) != 0 {
		t.Errorf("Expected Maker Balance to be equal to %v but got %v instead", 1, Taker2ZRXIncrement)
	}
}
