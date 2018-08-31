package operator_test

import (
	"log"
	"math/big"
	"reflect"
	"sync"
	"testing"
	"time"

	"github.com/Proofsuite/amp-matching-engine/app"
	"github.com/Proofsuite/amp-matching-engine/contracts"
	"github.com/Proofsuite/amp-matching-engine/operator"
	"github.com/Proofsuite/amp-matching-engine/rabbitmq"
	"github.com/Proofsuite/amp-matching-engine/services"
	"github.com/Proofsuite/amp-matching-engine/types"
	"github.com/Proofsuite/amp-matching-engine/utils/testutils"
	"github.com/Proofsuite/amp-matching-engine/utils/testutils/mocks"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/backends"
	"github.com/ethereum/go-ethereum/common"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func SetupTest(t *testing.T) (
	*operator.Operator,
	*contracts.Exchange,
	common.Address,
	[]*types.Wallet,
	common.Address,
	common.Address,
	*testutils.OrderFactory,
	*testutils.OrderFactory,
	*backends.SimulatedBackend,
) {

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

	tradeService := new(mocks.TradeService)
	ethereumService := new(mocks.EthereumService)

	wallets := []*types.Wallet{wallet1, wallet2, wallet3, wallet4, wallet5}
	admin := wallet1
	maker := wallet4
	taker := wallet5

	walletDao := new(mocks.WalletDao)
	walletDao.On("GetDefaultAdminWallet").Return(wallet1, nil)
	txService := services.NewTxService(walletDao, admin)
	walletService := services.NewWalletService(walletDao)

	//setup mocks
	deployer, err := testutils.NewSimulator(walletService, txService, []common.Address{wallet1.Address, wallet2.Address, wallet3.Address, wallet4.Address, wallet5.Address})
	if err != nil {
		panic(err)
	}

	//Initially Maker owns 1e18 units of sellToken and Taker owns 1e18 units buyToken
	wethToken, weth, _, err := deployer.DeployToken(maker.Address, big.NewInt(1e18))
	if err != nil {
		t.Errorf("Error deploying token 1: %v", err)
	}

	zrxToken, zrx, _, err := deployer.DeployToken(taker.Address, big.NewInt(1e18))
	if err != nil {
		t.Errorf("Error deploying token 2: %v", err)
	}

	exchange, exchangeAddr, _, err := deployer.DeployExchange(admin.Address, weth)
	if err != nil {
		t.Errorf("Could not deploy exchange: %v", err)
	}

	_, err = exchange.SetOperator(admin.Address, true)
	if err != nil {
		t.Errorf("Could not set operator: %v", err)
	}

	simulator := deployer.Backend.(*backends.SimulatedBackend)
	simulator.Commit()

	exchange.PrintErrors()

	wethToken.SetTxSender(maker)
	_, err = wethToken.Approve(exchangeAddr, big.NewInt(1e18))
	if err != nil {
		t.Errorf("Could not approve sellToken: %v", err)
	}

	zrxToken.SetTxSender(taker)
	_, err = zrxToken.Approve(exchangeAddr, big.NewInt(1e18))
	if err != nil {
		t.Errorf("Could not approve buyToken: %v", err)
	}

	simulator.Commit()

	pair := &types.Pair{
		BaseTokenSymbol:   "ZRX",
		QuoteTokenSymbol:  "WETH",
		BaseTokenAddress:  zrx,
		QuoteTokenAddress: weth,
	}

	op, err := operator.NewOperator(
		walletService,
		txService,
		tradeService,
		ethereumService,
		exchange,
	)

	factory1, err := testutils.NewOrderFactory(pair, maker, exchangeAddr)
	if err != nil {
		panic(err)
	}

	factory2, err := testutils.NewOrderFactory(pair, taker, exchangeAddr)
	if err != nil {
		panic(err)
	}

	if err != nil {
		panic(err)
	}

	return op, exchange, exchangeAddr, wallets, zrx, weth, factory1, factory2, simulator
}

func TestNewOperator(t *testing.T) {
	op, _, _, _, _, _, _, _, _ := SetupTest(t)

	log.Print(op.TxQueues)
}

func TestGetShortestQueue(t *testing.T) {
	op, exchange, _, wallets, zrx, weth, factory, _, _ := SetupTest(t)

	opWallet1 := wallets[0]
	opWallet2 := wallets[1]
	opWallet3 := wallets[2]

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

	//Refactor this into better helper function or something shorter
	err := txq1.PurgePendingTrades()
	if err != nil {
		t.Errorf("Could not purge pending trades")
	}

	err = txq2.PurgePendingTrades()
	if err != nil {
		t.Errorf("Could not purge pending trades")
	}

	err = txq3.PurgePendingTrades()
	if err != nil {
		t.Errorf("Could not purge pending trades")
	}

	op.TxQueues = []*operator.TxQueue{txq1, txq2, txq3}

	o1, _ := factory.NewOrder(zrx, 1, weth, 1)
	o2, _ := factory.NewOrder(zrx, 1, weth, 1)
	o3, _ := factory.NewOrder(zrx, 1, weth, 1)
	o4, _ := factory.NewOrder(zrx, 1, weth, 1)
	o5, _ := factory.NewOrder(zrx, 1, weth, 1)

	t1, _ := factory.NewTrade(o1, 1)
	t2, _ := factory.NewTrade(o2, 1)
	t3, _ := factory.NewTrade(o3, 1)
	t4, _ := factory.NewTrade(o4, 1)
	t5, _ := factory.NewTrade(o5, 1)

	txq1.PublishPendingTrade(o1, t1)
	txq2.PublishPendingTrade(o2, t2)
	txq3.PublishPendingTrade(o3, t3)
	txq1.PublishPendingTrade(o4, t4)
	txq2.PublishPendingTrade(o5, t5)

	time.Sleep(time.Second)

	shortest, ln, err := op.GetShortestQueue()
	if err != nil {
		t.Errorf("Could not get shortest queue: %v", err)
	}

	assert.Equal(t, 1, ln)

	if !reflect.DeepEqual(shortest, txq3) {
		t.Errorf("Wrong shortest queue:\n Expected: %v\n, Got: %v\n", shortest, txq3)
	}
}

func TestPublishPendingTrade(t *testing.T) {
	op, exchange, _, wallets, zrx, weth, factory, _, _ := SetupTest(t)

	opWallet1 := wallets[0]

	tradeService := new(mocks.TradeService)
	ethService := new(mocks.EthereumService)

	txq := &operator.TxQueue{
		Name:            "queue1",
		TradeService:    tradeService,
		EthereumService: ethService,
		Wallet:          opWallet1,
		Exchange:        exchange,
	}

	err := txq.PurgePendingTrades()
	if err != nil {
		t.Errorf("Could not purge pending trades")
	}

	op.TxQueues = []*operator.TxQueue{txq}

	o1, _ := factory.NewOrder(zrx, 1, weth, 1)
	o2, _ := factory.NewOrder(zrx, 1, weth, 1)
	o3, _ := factory.NewOrder(zrx, 1, weth, 1)

	t1, _ := factory.NewTrade(o1, 1)
	t2, _ := factory.NewTrade(o2, 1)
	t3, _ := factory.NewTrade(o3, 1)

	txq.PublishPendingTrade(o1, t1)
	txq.PublishPendingTrade(o2, t2)
	txq.PublishPendingTrade(o3, t3)

	time.Sleep(time.Millisecond)
	assert.Equal(t, 3, txq.Length())

	_, err = txq.PopPendingTrade()
	if err != nil {
		t.Errorf("Could not pop pending trade")
	}

	time.Sleep(time.Millisecond)
	assert.Equal(t, 2, txq.Length())

	_, err = txq.PopPendingTrade()
	if err != nil {
		t.Errorf("Could not pop pending trade")
	}

	_, err = txq.PopPendingTrade()
	if err != nil {
		t.Errorf("Could not pop pending trade")
	}

	time.Sleep(time.Millisecond)
	assert.Equal(t, 0, txq.Length())
}

func TestPopPendingTrade(t *testing.T) {
	op, exchange, _, wallets, zrx, weth, factory, _, _ := SetupTest(t)

	opWallet1 := wallets[0]

	tradeService := new(mocks.TradeService)
	ethService := new(mocks.EthereumService)

	txq := &operator.TxQueue{
		Name:            "queue1",
		TradeService:    tradeService,
		EthereumService: ethService,
		Wallet:          opWallet1,
		Exchange:        exchange,
	}

	err := txq.PurgePendingTrades()
	if err != nil {
		t.Errorf("Could not purge pending trades")
	}

	op.TxQueues = []*operator.TxQueue{txq}

	o1, _ := factory.NewOrder(zrx, 1, weth, 1)
	o2, _ := factory.NewOrder(zrx, 1, weth, 1)
	o3, _ := factory.NewOrder(zrx, 1, weth, 1)

	t1, _ := factory.NewTrade(o1, 1)
	t2, _ := factory.NewTrade(o2, 1)
	t3, _ := factory.NewTrade(o3, 1)

	txq.PublishPendingTrade(o1, t1)
	txq.PublishPendingTrade(o2, t2)
	txq.PublishPendingTrade(o3, t3)

	time.Sleep(time.Millisecond)
	assert.Equal(t, 3, txq.Length())

	_, err = txq.PopPendingTrade()
	if err != nil {
		t.Errorf("Could not pop pending trade")
	}

	time.Sleep(time.Millisecond)
	assert.Equal(t, 2, txq.Length())

	_, err = txq.PopPendingTrade()
	if err != nil {
		t.Errorf("Could not pop pending trade")
	}

	_, err = txq.PopPendingTrade()
	if err != nil {
		t.Errorf("Could not pop pending trade")
	}

	time.Sleep(time.Millisecond)
	assert.Equal(t, 0, txq.Length())
}

func TestExecuteTrade(t *testing.T) {
	_, _, _, wallets, zrx, weth, factory, _, _ := SetupTest(t)

	o1, _ := factory.NewOrder(zrx, 1, weth, 1)
	t1, _ := factory.NewTrade(o1, 1)

	mockTx := &ethTypes.Transaction{}
	tradeService := new(mocks.TradeService)
	ethService := new(mocks.EthereumService)
	exchange := new(mocks.Exchange)
	tradeService.On("UpdateTradeTx", t1, mockTx).Return(nil)
	exchange.On("Trade", o1, t1).Return(mockTx, nil)

	txq := &operator.TxQueue{
		Name:            "queue1",
		TradeService:    tradeService,
		EthereumService: ethService,
		Wallet:          wallets[0],
		Exchange:        exchange,
	}

	err := txq.PurgePendingTrades()
	if err != nil {
		t.Errorf("Could not purge pending trades")
	}

	tx, err := txq.ExecuteTrade(o1, t1)
	if err != nil {
		t.Errorf("Could not execute trade: %v", err)
	}

	tradeService.AssertCalled(t, "UpdateTradeTx", t1, mockTx)
	exchange.AssertCalled(t, "Trade", o1, t1)

	if !reflect.DeepEqual(tx, mockTx) {
		t.Errorf("Expected: %v\n, Got: %v\n", tx, mockTx)
	}

}

func TestQueueTrade(t *testing.T) {
	_, _, _, wallets, zrx, weth, factory, _, _ := SetupTest(t)

	o1, _ := factory.NewOrder(zrx, 1, weth, 1)
	t1, _ := factory.NewTrade(o1, 1)

	mockTx := &ethTypes.Transaction{}
	tradeService := new(mocks.TradeService)
	ethService := new(mocks.EthereumService)
	exchange := new(mocks.Exchange)
	tradeService.On("UpdateTradeTx", t1, mockTx).Return(nil)
	exchange.On("Trade", o1, t1).Return(mockTx, nil)

	opWallet1 := wallets[0]

	txq := &operator.TxQueue{
		Name:            "queue1",
		TradeService:    tradeService,
		EthereumService: ethService,
		Wallet:          opWallet1,
		Exchange:        exchange,
	}

	err := txq.PurgePendingTrades()
	if err != nil {
		t.Errorf("Could not purge pending trades")
	}

	err = txq.QueueTrade(o1, t1)
	if err != nil {
		t.Errorf("Could not execute trade: %v", err)
	}

	tradeService.AssertCalled(t, "UpdateTradeTx", t1, mockTx)
	exchange.AssertCalled(t, "Trade", o1, t1)
}

func TestHandleEvents(t *testing.T) {
	_, exchange, _, wallets, zrx, weth, factory1, factory2, simulator := SetupTest(t)

	opMessages := make(chan *types.OperatorMessage)
	rabbitmq.SubscribeOperator(func(msg *types.OperatorMessage) error {
		opMessages <- msg
		return nil
	})

	o1, _ := factory1.NewOrder(zrx, 1, weth, 1)
	t1, _ := factory2.NewTrade(o1, 1)

	admin := wallets[0]

	tradeService := new(mocks.TradeService)
	orderService := new(mocks.OrderService)
	ethService := new(mocks.EthereumService)
	tradeService.On("UpdateTradeTx", t1, mock.Anything).Return(nil)
	orderService.On("GetByHash", t1.OrderHash).Return(o1, nil)
	tradeService.On("GetByHash", t1.Hash).Return(t1, nil)
	ethService.On("WaitMined", mock.Anything).Return(nil, nil)

	txq := &operator.TxQueue{
		Name:            "queue1",
		TradeService:    tradeService,
		OrderService:    orderService,
		EthereumService: ethService,
		Wallet:          wallets[0],
		Exchange:        exchange,
	}

	err := txq.PurgePendingTrades()
	if err != nil {
		t.Errorf("Could not purge pending trades")
	}

	time.Sleep(500 * time.Millisecond)
	go txq.HandleEvents()

	exchange.SetTxSender(admin)
	err = txq.QueueTrade(o1, t1)
	if err != nil {
		t.Errorf("Could not execute trade: %v", err)
	}

	simulator.Commit()

	wg := sync.WaitGroup{}
	wg.Add(2)

	go func() {
		for {
			msg := <-opMessages
			switch msg.MessageType {
			case "TRADE_SENT_MESSAGE":
				testutils.Compare(t, msg.Order.Hash, o1.Hash)
				testutils.Compare(t, msg.Trade.Hash, t1.Hash)
				testutils.Compare(t, msg.Trade.OrderHash, o1.Hash)
				wg.Done()
			case "TRADE_SUCCESS_MESSAGE":
				testutils.Compare(t, msg.Order.Hash, o1.Hash)
				testutils.Compare(t, msg.Trade.Hash, t1.Hash)
				testutils.Compare(t, msg.Trade.OrderHash, o1.Hash)
				wg.Done()
			}
		}
	}()

	wg.Wait()
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
