package operator_test

import (
	"encoding/json"
	"errors"
	"log"
	"math/big"
	"testing"

	"github.com/Proofsuite/amp-matching-engine/app"
	"github.com/Proofsuite/amp-matching-engine/contracts"
	"github.com/Proofsuite/amp-matching-engine/operator"
	"github.com/Proofsuite/amp-matching-engine/rabbitmq"
	"github.com/Proofsuite/amp-matching-engine/types"
	"github.com/Proofsuite/amp-matching-engine/utils/testutils"
	"github.com/Proofsuite/amp-matching-engine/utils/testutils/mocks"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
)

func SetupTest(t *testing.T) (*operator.Operator, *contracts.Exchange, common.Address, []*types.Wallet, *types.Token, *types.Token) {
	err := app.LoadConfig("../config", "")
	if err != nil {
		panic(err)
	}

	log.SetFlags(log.LstdFlags | log.Llongfile)
	log.SetPrefix("\nLOG: ")

	rabbitmq.InitConnection(app.Config.Rabbitmq)

	wallet1 := testutils.GetTestWallet1()
	wallet2 := testutils.GetTestWallet2()
	wallet3 := testutils.GetTestWallet3()
	wallet4 := testutils.GetTestWallet4()
	wallet5 := testutils.GetTestWallet5()

	zrx := testutils.GetTestZRXToken()
	weth := testutils.GetTestWETHToken()

	walletService := new(mocks.WalletService)
	txService := new(mocks.TxService)
	tradeService := new(mocks.TradeService)
	ethereumService := new(mocks.EthereumService)

	wallets := []*types.Wallet{wallet1, wallet2, wallet3, wallet4, wallet5}

	//setup mocks
	walletService.On("GetOperatorWallets").Return(wallets, nil)
	txService.On("GetTxSendOptions").Return(bind.NewKeyedTransactor(wallet1.PrivateKey), nil)

	deployer, err := testutils.NewSimulator(walletService, txService, []common.Address{wallet1.Address, wallet2.Address, wallet3.Address})
	if err != nil {
		panic(err)
	}

	factory, err := testutils.NewOrderFactory(pair, wallet1, exchangeAddress)
	if err != nil {
		panic(err)
	}

	feeAccount := common.HexToAddress(app.Config.FeeAccount)
	wethToken := common.HexToAddress(app.Config.WETH)

	exchange, exchangeAddr, _, err := deployer.DeployExchange(feeAccount, wethToken)
	if err != nil {
		t.Errorf("Could not deploy exchange: %v", err)
	}

	op, err := operator.NewOperator(
		walletService,
		txService,
		tradeService,
		ethereumService,
		exchange,
	)

	if err != nil {
		panic(err)
	}

	return op, exchange, exchangeAddr, wallets, zrx, weth
}

func TestNewOperator(t *testing.T) {
	op, _, _, _, _, _ := SetupTest(t)

	log.Print(op.TxQueues)
}

func TestQueueTrade(t *testing.T) {
	op, _, exchangeAddr, wallets, zrx, weth := SetupTest(t)

	maker := wallets[3]
	taker := wallets[4]

	//Maker creates an order that exchanges 'sellAmount' of sellToken for 'buyAmount' of buyToken
	o := &types.Order{
		ExchangeAddress: exchangeAddr,
		BuyAmount:       big.NewInt(1),
		SellAmount:      big.NewInt(1),
		Expires:         big.NewInt(1e8),
		Nonce:           big.NewInt(0),
		MakeFee:         big.NewInt(0),
		TakeFee:         big.NewInt(0),
		BuyToken:        zrx.ContractAddress,
		SellToken:       weth.ContractAddress,
		UserAddress:     maker.Address,
	}

	o.Sign(maker)

	tr := &types.Trade{
		OrderHash:  o.Hash,
		Amount:     big.NewInt(1),
		Taker:      taker.Address,
		TradeNonce: big.NewInt(0),
	}

	tr.Sign(taker)

	err := op.QueueTrade(o, tr)
	if err != nil {
		t.Errorf("Error queuing trade: %v", err)
	}
}

func TestGetShortestQueue(t *testing.T) {
	op, exchange, exchangeAddr, wallets, zrx, weth := SetupTest(t)

	opWallet1 := wallets[0]
	opWallet2 := wallets[1]
	opWallet3 := wallets[2]
	maker := wallets[3]
	taker := wallets[4]

	tradeService := new(mocks.TradeService)
	ethService := new(mocks.EthereumService)

	txq1 := &operator.TxQueue{
		Name:            "queue1",
		TradeService:    tradeService,
		EthereumService: ethService,
		Wallet:          opWallet1,
		Exchange:        exchange,
	}

	txq2 := &operator.TxQueue{
		Name:            "queue2",
		TradeService:    tradeService,
		EthereumService: ethService,
		Wallet:          opWallet2,
		Exchange:        exchange,
	}

	txq3 := &operator.TxQueue{
		Name:            "queue3",
		TradeService:    tradeService,
		EthereumService: ethService,
		Wallet:          opWallet3,
		Exchange:        exchange,
	}

	op.TxQueues = []*operator.TxQueue{txq1, txq2, txq3}

	q1 := operator.GetTxQueue("queue1")
	q2 := operator.GetTxQueue("queue2")
	q3 := operator.GetTxQueue("queue3")

	pding1 := factory.NewPendingTradeMessage()
	pding2 := factory.NewPendingTradeMessage()
	pding3 := factory.NewPendingTradeMessage()
	pding4 := factory.NewPendingTradeMessage()

	bytes, err := json.Marshal(t)
	if err != nil {
		return errors.New
	}

}

// func TestQueueTrade(t *testing.T) {
// 	op := SetupTest(t)

// 	order := &types.Order{
// 		ExchangeAddress: exchangeAddr,
// 		BuyAmount:       buyAmount,
// 		SellAmount:      sellAmount,
// 		Expires:         expires,
// 		Nonce:           big.NewInt(0),
// 		MakeFee:         big.NewInt(0),
// 		TakeFee:         big.NewInt(0),
// 		BuyToken:        buyTokenAddr,
// 		SellToken:       sellTokenAddr,
// 		UserAddress:     maker.Address,
// 	}

// 	order.Sign(maker)

// 	trade := &types.Trade{
// 		OrderHash:  order.Hash,
// 		Amount:     amount,
// 		Taker:      taker.Address,
// 		TradeNonce: big.NewInt(0),
// 	}

// 	trade.Sign(taker)
// }

// func TestNewOperator2(t *testing.T) {

// }

// func TestOperator(t *testing.T) {
// 	var wg sync.WaitGroup
// 	wg.Add(2)
// 	testConfig := NewConfiguration()
// 	opParams := testConfig.OperatorParams

// 	ZRX := testConfig.QuoteTokens["ZRX"]
// 	WETH := testConfig.QuoteTokens["WETH"]
// 	ZRXWETH := NewPair(ZRX, WETH)

// 	admin := config.Wallets[0]
// 	maker := config.Wallets[1]
// 	taker := config.Wallets[2]
// 	exchange := config.Exchange

// 	opConfig := &OperatorConfig{
// 		Admin:          admin,
// 		Exchange:       exchange,
// 		OperatorParams: opParams,
// 	}

// 	ZRXAmount := big.NewInt(1e18)
// 	WETHAmount := big.NewInt(1e18)

// 	deployer, err := NewWebsocketDeployer(admin)
// 	if err != nil {
// 		t.Errorf("Could not instantiate deployer: %v", err)
// 	}

// 	ex, err := deployer.NewExchange(exchange)
// 	if err != nil {
// 		t.Errorf("Could not retrieve exchange instance: %v", err)
// 	}

// 	makerFactory := NewOrderFactory(&ZRXWETH, maker)
// 	takerFactory := NewOrderFactory(&ZRXWETH, taker)
// 	makerFactory.SetExchangeAddress(ex.Address)
// 	takerFactory.SetExchangeAddress(ex.Address)

// 	operator, err := NewOperator(opConfig)
// 	if err != nil {
// 		t.Errorf("Could not instantiate operator: %v", err)
// 	}

// 	initialTakerZRXBalance, _ := ex.TokenBalance(taker.Address, ZRX.Address)
// 	initialTakerWETHBalance, _ := ex.TokenBalance(taker.Address, WETH.Address)
// 	initialMakerZRXBalance, _ := ex.TokenBalance(maker.Address, ZRX.Address)
// 	initialMakerWETHBalance, _ := ex.TokenBalance(maker.Address, WETH.Address)

// 	order, _ := makerFactory.NewOrderWithEvents(WETH, WETHAmount.Int64(), ZRX, ZRXAmount.Int64())
// 	trade, _ := takerFactory.NewTradeWithEvents(order, 1)

// 	err = operator.AddTradeToExecutionList(order, trade)
// 	if err != nil {
// 		t.Errorf("Could not execute trade: %v", err)
// 	}

// 	//Each received success event sends a done signal to the wait group
// 	//In total, the wait group (wg) needs to receive 2 different signals
// 	//corresponding to ORDER_TX_SUCCESS and TRADE_TX_SUCCESS
// 	go func() {
// 		for {
// 			select {
// 			case e := <-order.events:
// 				switch e.eventType {
// 				case ORDER_TX_SUCCESS:
// 					wg.Done()
// 				case ORDER_TX_ERROR:
// 					p := e.payload.(*TxErrorPayload)
// 					t.Errorf("Received Tx error payload: %v", p)
// 				}
// 			case e := <-trade.events:
// 				switch e.eventType {
// 				case TRADE_TX_SUCCESS:
// 					wg.Done()
// 				case TRADE_TX_ERROR:
// 					p := e.payload.(*TxErrorPayload)
// 					t.Errorf("Received Tx error payload: %v", p)
// 				}
// 			}
// 		}
// 	}()

// 	wg.Wait()

// 	TakerZRXBalance, _ := ex.TokenBalance(taker.Address, ZRX.Address)
// 	TakerWETHBalance, _ := ex.TokenBalance(taker.Address, WETH.Address)
// 	MakerZRXBalance, _ := ex.TokenBalance(maker.Address, ZRX.Address)
// 	MakerWETHBalance, _ := ex.TokenBalance(maker.Address, WETH.Address)

// 	TakerZRXIncrement := big.NewInt(0)
// 	TakerWETHIncrement := big.NewInt(0)
// 	MakerZRXIncrement := big.NewInt(0)
// 	MakerWETHIncrement := big.NewInt(0)

// 	MakerZRXIncrement.Sub(MakerZRXBalance, initialMakerZRXBalance)
// 	MakerWETHIncrement.Sub(MakerWETHBalance, initialMakerWETHBalance)
// 	TakerZRXIncrement.Sub(TakerZRXBalance, initialTakerZRXBalance)
// 	TakerWETHIncrement.Sub(TakerWETHBalance, initialTakerWETHBalance)

// 	if MakerZRXIncrement.Cmp(big.NewInt(-1)) != 0 {
// 		t.Errorf("Expected Maker Balance to be equal to be %v but got %v instead", -1, MakerZRXIncrement)
// 	}

// 	if MakerWETHIncrement.Cmp(big.NewInt(1)) != 0 {
// 		t.Errorf("Expected Taker Balance to be equal to be %v but got %v instead", 1, MakerWETHIncrement)
// 	}

// 	if TakerWETHIncrement.Cmp(big.NewInt(-1)) != 0 {
// 		t.Errorf("Expected Maker Balance to be equal to %v but got %v instead", -1, TakerWETHIncrement)
// 	}

// 	if TakerZRXIncrement.Cmp(big.NewInt(1)) != 0 {
// 		t.Errorf("Expected Maker Balance to be equal to %v but got %v instead", 1, TakerZRXIncrement)
// 	}
// }

// func TestOperator2(t *testing.T) {
// 	var wg sync.WaitGroup
// 	wg.Add(4)
// 	testConfig := NewConfiguration()
// 	opParams := testConfig.OperatorParams

// 	ZRX := testConfig.QuoteTokens["ZRX"]
// 	WETH := testConfig.QuoteTokens["WETH"]
// 	ZRXWETH := NewPair(ZRX, WETH)

// 	admin := config.Wallets[0]
// 	maker1 := config.Wallets[1]
// 	taker1 := config.Wallets[2]
// 	maker2 := config.Wallets[3]
// 	taker2 := config.Wallets[4]

// 	exchange := config.Exchange

// 	opConfig := &OperatorConfig{
// 		Admin:          admin,
// 		Exchange:       exchange,
// 		OperatorParams: opParams,
// 	}

// 	ZRXAmount := big.NewInt(1e18)
// 	WETHAmount := big.NewInt(1e18)

// 	deployer, err := NewWebsocketDeployer(admin)
// 	if err != nil {
// 		t.Errorf("Could not instantiate deployer: %v", err)
// 	}

// 	ex, err := deployer.NewExchange(exchange)
// 	if err != nil {
// 		t.Errorf("Could not retrieve exchange instance: %v", err)
// 	}

// 	maker1Factory := NewOrderFactory(&ZRXWETH, maker1)
// 	taker1Factory := NewOrderFactory(&ZRXWETH, taker1)
// 	maker2Factory := NewOrderFactory(&ZRXWETH, maker2)
// 	taker2Factory := NewOrderFactory(&ZRXWETH, taker2)

// 	maker1Factory.SetExchangeAddress(ex.Address)
// 	taker1Factory.SetExchangeAddress(ex.Address)
// 	maker2Factory.SetExchangeAddress(ex.Address)
// 	taker2Factory.SetExchangeAddress(ex.Address)

// 	operator, err := NewOperator(opConfig)
// 	if err != nil {
// 		t.Errorf("Could not instantiate operator: %v", err)
// 	}

// 	initialTaker1ZRXBalance, _ := ex.TokenBalance(taker1.Address, ZRX.Address)
// 	initialTaker1WETHBalance, _ := ex.TokenBalance(taker1.Address, WETH.Address)
// 	initialMaker1ZRXBalance, _ := ex.TokenBalance(maker1.Address, ZRX.Address)
// 	initialMaker1WETHBalance, _ := ex.TokenBalance(maker1.Address, WETH.Address)

// 	initialTaker2ZRXBalance, _ := ex.TokenBalance(taker2.Address, ZRX.Address)
// 	initialTaker2WETHBalance, _ := ex.TokenBalance(taker2.Address, WETH.Address)
// 	initialMaker2ZRXBalance, _ := ex.TokenBalance(maker2.Address, ZRX.Address)
// 	initialMaker2WETHBalance, _ := ex.TokenBalance(maker2.Address, WETH.Address)

// 	o1, _ := maker1Factory.NewOrderWithEvents(WETH, WETHAmount.Int64(), ZRX, ZRXAmount.Int64())
// 	t1, _ := taker1Factory.NewTradeWithEvents(o1, 1)

// 	o2, _ := maker2Factory.NewOrderWithEvents(WETH, WETHAmount.Int64(), ZRX, ZRXAmount.Int64())
// 	t2, _ := taker2Factory.NewTradeWithEvents(o2, 1)

// 	err = operator.AddTradeToExecutionList(o1, t1)
// 	err = operator.AddTradeToExecutionList(o2, t2)

// 	//Each received success event sends a done signal to the wait group
// 	//In total, the wait group (wg) needs to receive 4 different signals
// 	//corresponding to 2 ORDER_TX_SUCCESS (o1 and o2) and 2 TRADE_TX_SUCCESS (t1 and t2)
// 	go func() {
// 		for {
// 			select {
// 			case e := <-o1.events:
// 				switch e.eventType {
// 				case ORDER_TX_SUCCESS:
// 					wg.Done()
// 				case ORDER_TX_ERROR:
// 					p := e.payload.(*TxErrorPayload)
// 					t.Errorf("Received Tx error payload: %v", p)
// 				}
// 			case e := <-o2.events:
// 				switch e.eventType {
// 				case ORDER_TX_SUCCESS:
// 					wg.Done()
// 				case ORDER_TX_ERROR:
// 					p := e.payload.(*TxErrorPayload)
// 					t.Errorf("Received Tx error payload: %v", p)
// 				}
// 			case e := <-t1.events:
// 				switch e.eventType {
// 				case TRADE_TX_SUCCESS:
// 					wg.Done()
// 				case TRADE_TX_ERROR:
// 					p := e.payload.(*TxErrorPayload)
// 					t.Errorf("Received Tx error payload: %v", p)
// 				}
// 			case e := <-t2.events:
// 				switch e.eventType {
// 				case TRADE_TX_SUCCESS:
// 					wg.Done()
// 				case TRADE_TX_ERROR:
// 					p := e.payload.(*TxErrorPayload)
// 					t.Errorf("Received Tx error payload: %v", p)
// 				}
// 			}
// 		}
// 	}()

// 	wg.Wait()

// 	Taker1ZRXBalance, _ := ex.TokenBalance(taker1.Address, ZRX.Address)
// 	Taker1WETHBalance, _ := ex.TokenBalance(taker1.Address, WETH.Address)
// 	Maker1ZRXBalance, _ := ex.TokenBalance(maker1.Address, ZRX.Address)
// 	Maker1WETHBalance, _ := ex.TokenBalance(maker1.Address, WETH.Address)

// 	Taker2ZRXBalance, _ := ex.TokenBalance(taker2.Address, ZRX.Address)
// 	Taker2WETHBalance, _ := ex.TokenBalance(taker2.Address, WETH.Address)
// 	Maker2ZRXBalance, _ := ex.TokenBalance(maker2.Address, ZRX.Address)
// 	Maker2WETHBalance, _ := ex.TokenBalance(maker2.Address, WETH.Address)

// 	Taker1ZRXIncrement := big.NewInt(0)
// 	Taker2ZRXIncrement := big.NewInt(0)
// 	Taker1WETHIncrement := big.NewInt(0)
// 	Taker2WETHIncrement := big.NewInt(0)
// 	Maker1ZRXIncrement := big.NewInt(0)
// 	Maker2ZRXIncrement := big.NewInt(0)
// 	Maker1WETHIncrement := big.NewInt(0)
// 	Maker2WETHIncrement := big.NewInt(0)

// 	Maker1ZRXIncrement.Sub(Maker1ZRXBalance, initialMaker1ZRXBalance)
// 	Maker1WETHIncrement.Sub(Maker1WETHBalance, initialMaker1WETHBalance)
// 	Taker1ZRXIncrement.Sub(Taker1ZRXBalance, initialTaker1ZRXBalance)
// 	Taker1WETHIncrement.Sub(Taker1WETHBalance, initialTaker1WETHBalance)
// 	Maker2ZRXIncrement.Sub(Maker2ZRXBalance, initialMaker2ZRXBalance)
// 	Maker2WETHIncrement.Sub(Maker2WETHBalance, initialMaker2WETHBalance)
// 	Taker2ZRXIncrement.Sub(Taker2ZRXBalance, initialTaker2ZRXBalance)
// 	Taker2WETHIncrement.Sub(Taker2WETHBalance, initialTaker2WETHBalance)

// 	if Maker1ZRXIncrement.Cmp(big.NewInt(-1)) != 0 {
// 		t.Errorf("Expected Maker Balance to be equal to be %v but got %v instead", -1, Maker1ZRXIncrement)
// 	}

// 	if Maker1WETHIncrement.Cmp(big.NewInt(1)) != 0 {
// 		t.Errorf("Expected Taker Balance to be equal to be %v but got %v instead", 1, Maker1WETHIncrement)
// 	}

// 	if Taker1WETHIncrement.Cmp(big.NewInt(-1)) != 0 {
// 		t.Errorf("Expected Maker Balance to be equal to %v but got %v instead", -1, Taker1WETHIncrement)
// 	}

// 	if Taker1ZRXIncrement.Cmp(big.NewInt(1)) != 0 {
// 		t.Errorf("Expected Maker Balance to be equal to %v but got %v instead", 1, Taker1ZRXIncrement)
// 	}

// 	if Maker2ZRXIncrement.Cmp(big.NewInt(-1)) != 0 {
// 		t.Errorf("Expected Maker Balance to be equal to be %v but got %v instead", -1, Maker2ZRXIncrement)
// 	}

// 	if Maker2WETHIncrement.Cmp(big.NewInt(1)) != 0 {
// 		t.Errorf("Expected Taker Balance to be equal to be %v but got %v instead", 1, Maker2WETHIncrement)
// 	}

// 	if Taker2WETHIncrement.Cmp(big.NewInt(-1)) != 0 {
// 		t.Errorf("Expected Maker Balance to be equal to %v but got %v instead", -1, Taker2WETHIncrement)
// 	}

// 	if Taker2ZRXIncrement.Cmp(big.NewInt(1)) != 0 {
// 		t.Errorf("Expected Maker Balance to be equal to %v but got %v instead", 1, Taker2ZRXIncrement)
// 	}
// }
