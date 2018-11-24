package operator_test

import (
	"fmt"
	"log"
	"math/big"
	"reflect"
	"sync"
	"testing"
	"time"

	"github.com/Proofsuite/amp-matching-engine/app"
	"github.com/Proofsuite/amp-matching-engine/contracts"
	"github.com/Proofsuite/amp-matching-engine/ethereum"
	"github.com/Proofsuite/amp-matching-engine/operator"
	"github.com/Proofsuite/amp-matching-engine/rabbitmq"
	"github.com/Proofsuite/amp-matching-engine/services"
	"github.com/Proofsuite/amp-matching-engine/types"
	"github.com/Proofsuite/amp-matching-engine/utils/testutils"
	"github.com/Proofsuite/amp-matching-engine/utils/testutils/mocks"
	"github.com/ethereum/go-ethereum/common"
	eth "github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func SetupTest(t *testing.T) (
	*operator.Operator,
	*contracts.Exchange,
	[]*types.Wallet,
	common.Address,
	common.Address,
	*testutils.OrderFactory,
	*testutils.OrderFactory,
	*ethereum.SimulatedClient,
	*mocks.TradeService,
	*mocks.OrderService,
	*rabbitmq.Connection,
) {

	err := app.LoadConfig("../config", "")
	if err != nil {
		panic(err)
	}

	log.SetFlags(log.LstdFlags | log.Llongfile)
	log.SetPrefix("\nLOG: ")

	rabbitConn := rabbitmq.InitConnection(app.Config.RabbitMQURL)

	wallet1 := testutils.GetTestWallet1()
	wallet2 := testutils.GetTestWallet2()
	wallet3 := testutils.GetTestWallet3()
	wallet4 := testutils.GetTestWallet4()
	wallet5 := testutils.GetTestWallet5()

	tradeService := new(mocks.TradeService)
	orderService := new(mocks.OrderService)

	wallets := []*types.Wallet{wallet1, wallet2, wallet3, wallet4, wallet5}
	admin := wallet1
	maker := wallet4
	taker := wallet5

	//TODO refactor this
	walletDao := new(mocks.WalletDao)
	walletDao.On("GetDefaultAdminWallet").Return(wallet1, nil)
	walletDao.On("GetOperatorWallets").Return([]*types.Wallet{wallet1, wallet2, wallet3}, nil)
	txService := services.NewTxService(walletDao, admin)
	walletService := services.NewWalletService(walletDao)
	//setup mocks

	client := ethereum.NewSimulatedClient([]common.Address{wallet1.Address, wallet2.Address, wallet3.Address, wallet4.Address, wallet5.Address})
	if err != nil {
		panic(err)
	}

	deployer := testutils.NewDeployer(walletService, txService, client)
	provider := ethereum.NewEthereumProvider(client)
	// provider := ethereum.NewEthereumProvider(simulator)

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

	txOpts, err := exchange.DefaultTxOptions()
	if err != nil {
		t.Errorf("Could not retrieve default tx options")
	}

	_, err = exchange.SetOperator(wallet1.Address, true, txOpts)
	if err != nil {
		t.Errorf("Could not set operator: %v", err)
	}

	_, err = exchange.SetOperator(wallet2.Address, true, txOpts)
	if err != nil {
		t.Errorf("Could not set operator: %v", err)
	}

	_, err = exchange.SetOperator(wallet3.Address, true, txOpts)
	if err != nil {
		t.Errorf("Could not set operator: %v", err)
	}

	client.Commit()

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

	client.Commit()

	pair := &types.Pair{
		BaseTokenSymbol:   "ZRX",
		QuoteTokenSymbol:  "WETH",
		BaseTokenAddress:  zrx,
		QuoteTokenAddress: weth,
	}

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

	op, err := operator.NewOperator(
		walletService,
		tradeService,
		orderService,
		provider,
		exchange,
		rabbitConn,
	)

	if err != nil {
		panic(err)
	}

	return op, exchange, wallets, zrx, weth, factory1, factory2, client, tradeService, orderService, rabbitConn
}

func TestGetShortestQueue(t *testing.T) {
	op, _, _, zrx, weth, factory1, _, _, _, _, _ := SetupTest(t)
	txq1 := op.TxQueues[0]
	txq2 := op.TxQueues[1]
	txq3 := op.TxQueues[2]

	defer txq1.PurgePendingTrades()
	defer txq2.PurgePendingTrades()
	defer txq3.PurgePendingTrades()

	o1, _ := factory1.NewOrder(zrx, 1, weth, 1)
	o2, _ := factory1.NewOrder(zrx, 1, weth, 1)
	o3, _ := factory1.NewOrder(zrx, 1, weth, 1)
	o4, _ := factory1.NewOrder(zrx, 1, weth, 1)
	o5, _ := factory1.NewOrder(zrx, 1, weth, 1)

	t1, _ := factory1.NewTrade(o1, 1)
	t2, _ := factory1.NewTrade(o2, 1)
	t3, _ := factory1.NewTrade(o3, 1)
	t4, _ := factory1.NewTrade(o4, 1)
	t5, _ := factory1.NewTrade(o5, 1)

	txq1.PublishPendingTrade(o1, &t1)
	txq2.PublishPendingTrade(o2, &t2)
	txq3.PublishPendingTrade(o3, &t3)
	txq1.PublishPendingTrade(o4, &t4)
	txq2.PublishPendingTrade(o5, &t5)

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
	op, _, _, zrx, weth, factory1, _, _, _, _, _ := SetupTest(t)

	txq := op.TxQueues[0]
	defer txq.PurgePendingTrades()

	o1, _ := factory1.NewOrder(zrx, 1, weth, 1)
	o2, _ := factory1.NewOrder(zrx, 1, weth, 1)
	o3, _ := factory1.NewOrder(zrx, 1, weth, 1)

	t1, _ := factory1.NewTrade(o1, 1)
	t2, _ := factory1.NewTrade(o2, 1)
	t3, _ := factory1.NewTrade(o3, 1)

	txq.PublishPendingTrade(o1, &t1)
	txq.PublishPendingTrade(o2, &t2)
	txq.PublishPendingTrade(o3, &t3)

	time.Sleep(time.Millisecond)
	assert.Equal(t, 3, txq.Length())

	_, err := txq.PopPendingTrade()
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
	op, _, _, zrx, weth, factory1, _, _, _, _, _ := SetupTest(t)

	txq := op.TxQueues[0]
	defer txq.PurgePendingTrades()

	o1, _ := factory1.NewOrder(zrx, 1, weth, 1)
	o2, _ := factory1.NewOrder(zrx, 1, weth, 1)
	o3, _ := factory1.NewOrder(zrx, 1, weth, 1)

	t1, _ := factory1.NewTrade(o1, 1)
	t2, _ := factory1.NewTrade(o2, 1)
	t3, _ := factory1.NewTrade(o3, 1)

	txq.PublishPendingTrade(o1, &t1)
	txq.PublishPendingTrade(o2, &t2)
	txq.PublishPendingTrade(o3, &t3)

	time.Sleep(time.Millisecond)
	assert.Equal(t, 3, txq.Length())

	_, err := txq.PopPendingTrade()
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
	_, _, wallets, zrx, weth, factory1, _, simulator, _, orderService, rabbitConn := SetupTest(t)

	o1, _ := factory1.NewOrder(zrx, 1, weth, 1)
	t1, _ := factory1.NewTrade(o1, 1)

	provider := ethereum.NewEthereumProvider(simulator)
	tradeService := new(mocks.TradeService)
	exchange := new(mocks.Exchange)
	mockTx := &eth.Transaction{}
	tradeService.On("UpdateTradeTxHash", mock.Anything, mock.Anything).Return(nil)
	exchange.On("Trade", o1, &t1, mock.Anything).Return(mockTx, nil)

	txq, err := operator.NewTxQueue(
		"queue1",
		tradeService,
		provider,
		orderService,
		wallets[0],
		exchange,
		rabbitConn,
	)
	if err != nil {
		t.Errorf("Could not create new queue")
	}
	defer txq.PurgePendingTrades()

	tx, err := txq.ExecuteTrade(o1, &t1)
	if err != nil {
		t.Errorf("Could not execute trade: %v", err)
	}

	tradeService.AssertCalled(t, "UpdateTradeTxHash", mock.Anything, mock.Anything)
	exchange.AssertCalled(t, "Trade", o1, &t1, mock.Anything)

	if !reflect.DeepEqual(tx, mockTx) {
		t.Errorf("Expected: %v\n, Got: %v\n", tx, mockTx)
	}

}

func TestQueueTrade(t *testing.T) {
	_, _, wallets, zrx, weth, factory1, _, simulator, _, orderService, rabbitConn := SetupTest(t)

	o1, _ := factory1.NewOrder(zrx, 1, weth, 1)
	t1, _ := factory1.NewTrade(o1, 1)

	provider := ethereum.NewEthereumProvider(simulator)
	tradeService := new(mocks.TradeService)
	exchange := new(mocks.Exchange)

	txq, err := operator.NewTxQueue(
		"queue1",
		tradeService,
		provider,
		orderService,
		wallets[0],
		exchange,
		rabbitConn,
	)

	if err != nil {
		t.Error("Could not create new tx queue", err)
	}

	defer txq.PurgePendingTrades()
	defer rabbitConn.PurgeOperatorQueue()

	mockTx := &eth.Transaction{}
	tradeService.On("UpdateTradeTxHash", &t1, mock.Anything).Return(nil)
	exchange.On("Trade", o1, &t1, mock.Anything).Return(mockTx, nil)

	err = txq.QueueTrade(o1, &t1)
	if err != nil {
		t.Errorf("Could not execute trade: %v", err)
	}

	simulator.Commit()

	//we copy the trade in order to copy the nonce generated by the exchange
	tradeService.AssertCalled(t, "UpdateTradeTxHash", &t1, mock.Anything)
	exchange.AssertCalled(t, "Trade", o1, &t1, mock.Anything)
}

func TestHandleEvents1(t *testing.T) {
	_, exchange, wallets, zrx, weth, factory1, factory2, simulator, tradeService, orderService, rabbitConn := SetupTest(t)

	opMessages := make(chan *types.OperatorMessage)
	done := make(chan bool)
	handler := func(msg *types.OperatorMessage) error {
		fmt.Println("RECEIVING MESSAGE IN TEST HANDLE EVENTS 1")
		opMessages <- msg
		return nil
	}

	rabbitConn.SubscribeOperator(handler)

	o1, _ := factory1.NewOrder(zrx, 1, weth, 1)
	t1, _ := factory2.NewTrade(o1, 1)

	provider := ethereum.NewEthereumProvider(simulator)
	tradeService.On("UpdateTradeTxHash", mock.Anything, mock.Anything).Return(nil).Once().Run(func(args mock.Arguments) {
		t1.TxHash = args.Get(1).(common.Hash)
	})
	orderService.On("GetByHash", t1.OrderHash).Return(o1, nil)
	tradeService.On("GetByHash", t1.Hash).Return(&t1, nil)

	txq, err := operator.NewTxQueue(
		"queue1",
		tradeService,
		provider,
		orderService,
		wallets[0],
		exchange,
		rabbitConn,
	)

	if err != nil {
		t.Errorf("Could not create new queue: %v", err)
	}

	err = txq.QueueTrade(o1, &t1)
	if err != nil {
		t.Errorf("Could not execute trade: %v", err)
	}

	simulator.Commit()
	wg := sync.WaitGroup{}
	wg.Add(2)

	go func() {
		for {
			select {
			case msg := <-opMessages:
				fmt.Println(msg)
				switch msg.MessageType {
				case "TRADE_SENT_MESSAGE":
					testutils.Compare(t, msg.Order.Hash, o1.Hash)
					testutils.Compare(t, msg.Trade.Hash, t1.Hash)
					testutils.Compare(t, msg.Trade.OrderHash, o1.Hash)
					go func() {
						time.Sleep(100 * time.Millisecond)
						simulator.Commit()
					}()
					wg.Done()
				case "TRADE_SUCCESS_MESSAGE":
					testutils.Compare(t, msg.Order.Hash, o1.Hash)
					testutils.Compare(t, msg.Trade.Hash, t1.Hash)
					testutils.Compare(t, msg.Trade.OrderHash, o1.Hash)
					wg.Done()
				}
			case <-done:
				return
			}
		}
	}()

	wg.Wait()
	done <- true

	txq.PurgePendingTrades()
	rabbitConn.PurgeOperatorQueue()
	rabbitConn.CloseOperatorChannel()
}

func TestHandleEvents2(t *testing.T) {
	_, exchange, wallets, zrx, weth, factory1, factory2, simulator, tradeService, orderService, rabbitConn := SetupTest(t)

	opMessages := make(chan *types.OperatorMessage)
	handler := func(msg *types.OperatorMessage) error {
		fmt.Println("RECEIVING MESSAGE IN TEST HANDLE EVENTS 2")
		opMessages <- msg
		return nil
	}

	rabbitConn.SubscribeOperator(handler)
	defer rabbitConn.CloseOperatorChannel()

	o1, _ := factory1.NewOrder(zrx, 1, weth, 1)
	t1, _ := factory2.NewTrade(o1, 1)
	o2, _ := factory1.NewOrder(zrx, 1, weth, 1)
	t2, _ := factory2.NewTrade(o2, 1)

	provider := ethereum.NewEthereumProvider(simulator)
	tradeService.On("UpdateTradeTxHash", mock.Anything, mock.Anything).Return(nil).Once().Run(func(args mock.Arguments) {
		t1.TxHash = args.Get(1).(common.Hash)
	})

	tradeService.On("UpdateTradeTxHash", mock.Anything, mock.Anything).Return(nil).Once().Run(func(args mock.Arguments) {
		t2.TxHash = args.Get(1).(common.Hash)
	})

	orderService.On("GetByHash", t1.OrderHash).Return(o1, nil)
	orderService.On("GetByHash", t2.OrderHash).Return(o2, nil)
	tradeService.On("GetByHash", t1.Hash).Return(&t1, nil)
	tradeService.On("GetByHash", t2.Hash).Return(&t2, nil)

	txq, err := operator.NewTxQueue(
		"queue1",
		tradeService,
		provider,
		orderService,
		wallets[0],
		exchange,
		rabbitConn,
	)

	if err != nil {
		t.Errorf("Could not create tx queue: %v", err)
	}

	defer txq.PurgePendingTrades()
	rabbitConn.PurgeOperatorQueue()
	defer rabbitConn.PurgeOperatorQueue()

	txq.QueueTrade(o1, &t1)
	time.Sleep(10 * time.Millisecond)
	txq.QueueTrade(o2, &t2)

	simulator.Commit()
	wg := sync.WaitGroup{}
	wg.Add(4)

	go func() {
		for {
			msg := <-opMessages
			fmt.Println(msg)
			switch msg.MessageType {
			case "TRADE_SENT_MESSAGE":
				go func() {
					time.Sleep(1000 * time.Millisecond)
					simulator.Commit()
				}()
				wg.Done()
			case "TRADE_SUCCESS_MESSAGE":
				wg.Done()
			case "TRADE_ERROR_MESSAGE":
				t.Errorf("Receive trade error message")
			}
		}
	}()

	wg.Wait()
}

func TestHandleEvents3(t *testing.T) {
	_, exchange, wallets, zrx, weth, factory1, factory2, simulator, tradeService, orderService, rabbitConn := SetupTest(t)

	opMessages := make(chan *types.OperatorMessage)
	rabbitConn.SubscribeOperator(func(msg *types.OperatorMessage) error {
		opMessages <- msg
		return nil
	})

	o1, _ := factory1.NewOrder(zrx, 1, weth, 1)
	t1, _ := factory2.NewTrade(o1, 1)
	o2, _ := factory1.NewOrder(zrx, 2, weth, 2)
	t2, _ := factory2.NewTrade(o2, 2)
	o3, _ := factory1.NewOrder(zrx, 3, weth, 3)
	t3, _ := factory2.NewTrade(o3, 3)
	o4, _ := factory1.NewOrder(zrx, 4, weth, 4)
	t4, _ := factory2.NewTrade(o4, 4)
	o5, _ := factory1.NewOrder(zrx, 5, weth, 5)
	t5, _ := factory2.NewTrade(o5, 5)

	provider := ethereum.NewEthereumProvider(simulator)

	tradeService.On("UpdateTradeTxHash", mock.Anything, mock.Anything).Return(nil).Once().Run(func(args mock.Arguments) {
		t1.TxHash = args.Get(1).(common.Hash)
	})

	tradeService.On("UpdateTradeTxHash", mock.Anything, mock.Anything).Return(nil).Once().Run(func(args mock.Arguments) {
		t2.TxHash = args.Get(1).(common.Hash)
	})

	tradeService.On("UpdateTradeTxHash", mock.Anything, mock.Anything).Return(nil).Once().Run(func(args mock.Arguments) {
		t3.TxHash = args.Get(1).(common.Hash)
	})

	tradeService.On("UpdateTradeTxHash", mock.Anything, mock.Anything).Return(nil).Once().Run(func(args mock.Arguments) {
		t4.TxHash = args.Get(1).(common.Hash)
	})

	tradeService.On("UpdateTradeTxHash", mock.Anything, mock.Anything).Return(nil).Once().Run(func(args mock.Arguments) {
		t5.TxHash = args.Get(1).(common.Hash)
	})

	orderService.On("GetByHash", t1.MakerOrderHash).Return(o1, nil)
	orderService.On("GetByHash", t2.MakerOrderHash).Return(o2, nil)
	orderService.On("GetByHash", t3.MakerOrderHash).Return(o3, nil)
	orderService.On("GetByHash", t4.MakerOrderHash).Return(o4, nil)
	orderService.On("GetByHash", t5.MakerOrderHash).Return(o5, nil)
	tradeService.On("GetByHash", t1.Hash).Return(&t1, nil)
	tradeService.On("GetByHash", t2.Hash).Return(&t2, nil)
	tradeService.On("GetByHash", t3.Hash).Return(&t3, nil)
	tradeService.On("GetByHash", t4.Hash).Return(&t4, nil)
	tradeService.On("GetByHash", t5.Hash).Return(&t5, nil)

	txq, err := operator.NewTxQueue(
		"queue1",
		tradeService,
		provider,
		orderService,
		wallets[0],
		exchange,
		rabbitConn,
	)

	if err != nil {
		t.Errorf("Could not create queue: %v", err)
	}

	defer txq.PurgePendingTrades()

	txq.QueueTrade(o1, &t1)
	time.Sleep(10 * time.Millisecond)
	txq.QueueTrade(o2, &t2)
	time.Sleep(10 * time.Millisecond)
	txq.QueueTrade(o3, &t3)
	time.Sleep(10 * time.Millisecond)
	txq.QueueTrade(o4, &t4)
	time.Sleep(10 * time.Millisecond)
	txq.QueueTrade(o5, &t5)

	go testutils.Mine(simulator)
	wg := sync.WaitGroup{}
	wg.Add(10)

	go func() {
		for {
			msg := <-opMessages
			switch msg.MessageType {
			case "TRADE_SENT_MESSAGE":
				wg.Done()
			case "TRADE_SUCCESS_MESSAGE":
				wg.Done()
			case "TRADE_ERROR_MESSAGE":
				t.Errorf("Received trade error message")
			}
		}
	}()

	wg.Wait()
}

//This test verifies whether a transaction queue continues to process transactions after a failing
//transaction. o3/t3 payload is signed with a wrong private key and will be rejected by the smart contracts
//The rest of the transactions are valid and should be sent successfully.
func TestHandleEvents4(t *testing.T) {
	op, exchange, wallets, zrx, weth, factory1, factory2, simulator, tradeService, orderService, rabbitConn := SetupTest(t)

	opMessages := make(chan *types.OperatorMessage)
	rabbitConn.SubscribeOperator(func(msg *types.OperatorMessage) error {
		opMessages <- msg
		return nil
	})

	o1, _ := factory1.NewOrder(zrx, 1, weth, 1)
	t1, _ := factory2.NewTrade(o1, 1)
	o2, _ := factory1.NewOrder(zrx, 2, weth, 2)
	t2, _ := factory2.NewTrade(o2, 2)
	o3, _ := factory1.NewOrder(zrx, 3, weth, 3)
	t3, _ := factory2.NewTrade(o3, 3)
	o4, _ := factory1.NewOrder(zrx, 4, weth, 4)
	t4, _ := factory2.NewTrade(o4, 4)
	o5, _ := factory1.NewOrder(zrx, 5, weth, 5)
	t5, _ := factory2.NewTrade(o5, 5)

	admin := wallets[0]

	// we simulate a failing order with a wrong signature
	t3.Sign(admin)

	tradeService.On("UpdateTradeTxHash", mock.Anything, mock.Anything).Return(nil).Once().Run(func(args mock.Arguments) {
		t1.TxHash = args.Get(1).(common.Hash)
	})

	tradeService.On("UpdateTradeTxHash", mock.Anything, mock.Anything).Return(nil).Once().Run(func(args mock.Arguments) {
		t2.TxHash = args.Get(1).(common.Hash)
	})

	tradeService.On("UpdateTradeTxHash", mock.Anything, mock.Anything).Return(nil).Once().Run(func(args mock.Arguments) {
		t4.TxHash = args.Get(1).(common.Hash)
	})

	tradeService.On("UpdateTradeTxHash", mock.Anything, mock.Anything).Return(nil).Once().Run(func(args mock.Arguments) {
		t5.TxHash = args.Get(1).(common.Hash)
	})

	orderService.On("GetByHash", t1.OrderHash).Return(o1, nil)
	orderService.On("GetByHash", t2.OrderHash).Return(o2, nil)
	orderService.On("GetByHash", t3.OrderHash).Return(o3, nil)
	orderService.On("GetByHash", t4.OrderHash).Return(o4, nil)
	orderService.On("GetByHash", t5.OrderHash).Return(o5, nil)
	tradeService.On("GetByHash", t1.Hash).Return(&t1, nil)
	tradeService.On("GetByHash", t2.Hash).Return(&t2, nil)
	tradeService.On("GetByHash", t3.Hash).Return(&t3, nil)
	tradeService.On("GetByHash", t4.Hash).Return(&t4, nil)
	tradeService.On("GetByHash", t5.Hash).Return(&t5, nil)

	txq, err := operator.NewTxQueue(
		"queue1",
		tradeService,
		op.EthereumProvider,
		orderService,
		wallets[0],
		exchange,
		rabbitConn,
	)

	if err != nil {
		t.Errorf("Could not create queue: %v", err)
	}

	defer txq.PurgePendingTrades()

	time.Sleep(2 * time.Millisecond)
	txq.QueueTrade(o1, &t1)
	time.Sleep(2 * time.Millisecond)
	txq.QueueTrade(o2, &t2)
	time.Sleep(2 * time.Millisecond)
	txq.QueueTrade(o3, &t3)
	time.Sleep(2 * time.Millisecond)
	txq.QueueTrade(o4, &t4)
	time.Sleep(2 * time.Millisecond)
	txq.QueueTrade(o5, &t5)

	go testutils.Mine(simulator)

	wg := sync.WaitGroup{}
	wg.Add(8)

	go func() {
		for {
			msg := <-opMessages
			switch msg.MessageType {
			case "TRADE_SENT_MESSAGE":
				log.Print("TRADE_SENT_MESSAGE")
				wg.Done()
			case "TRADE_SUCCESS_MESSAGE":
				log.Print("TRADE_SUCCESS_MESSAGE")
				wg.Done()
			case "TRADE_ERROR_MESSAGE":
				log.Print("TRADE_ERROR_MESSAGE")
				assert.Equal(t, msg.ErrID, 10)
				assert.Equal(t, msg.Trade.Hash, t3.Hash)
				wg.Done()
			}
		}
	}()

	wg.Wait()
}

func TestHandleEvents5(t *testing.T) {
	op, _, wallets, zrx, weth, factory1, factory2, simulator, tradeService, orderService, rabbitConn := SetupTest(t)

	opMessages := make(chan *types.OperatorMessage)
	rabbitConn.SubscribeOperator(func(msg *types.OperatorMessage) error {
		opMessages <- msg
		return nil
	})

	o1, _ := factory1.NewOrder(zrx, 1, weth, 1)
	t1, _ := factory2.NewTrade(o1, 1)

	//we sign the trade with a random private key
	admin := wallets[0]
	t1.Sign(admin)

	tradeService.On("UpdateTradeTxHash", mock.Anything, mock.Anything).Return(nil).Once().Run(func(args mock.Arguments) {
		t1.TxHash = args.Get(1).(common.Hash)
	})

	orderService.On("GetByHash", t1.OrderHash).Return(o1, nil)
	tradeService.On("GetByHash", t1.Hash).Return(&t1, nil)
	defer op.PurgeQueues()

	txq := op.TxQueues[0]
	err := txq.QueueTrade(o1, &t1)
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
				assert.Equal(t, t1.Hash, msg.Trade.Hash)
				wg.Done()
			case "TRADE_SUCCESS_MESSAGE":
				t.Errorf("Received trade success message")
			case "TRADE_ERROR_MESSAGE":
				// we verify that error ID = 2 (corresponds to invalid signature)
				assert.Equal(t, 2, msg.ErrID)
				assert.Equal(t, t1.Hash, msg.Trade.Hash)
				wg.Done()
			}
		}
	}()

	wg.Wait()
}

func TestOperator1(t *testing.T) {
	op, _, _, zrx, weth, factory1, factory2, simulator, tradeService, orderService, rabbitConn := SetupTest(t)

	opMessages := make(chan *types.OperatorMessage)
	rabbitConn.SubscribeOperator(func(msg *types.OperatorMessage) error {
		opMessages <- msg
		return nil
	})

	o1, _ := factory1.NewOrder(zrx, 1, weth, 1)
	t1, _ := factory2.NewTrade(o1, 1)
	o2, _ := factory1.NewOrder(zrx, 2, weth, 2)
	t2, _ := factory2.NewTrade(o2, 2)
	o3, _ := factory1.NewOrder(zrx, 3, weth, 3)
	t3, _ := factory2.NewTrade(o3, 3)
	o4, _ := factory1.NewOrder(zrx, 4, weth, 4)
	t4, _ := factory2.NewTrade(o4, 4)
	o5, _ := factory1.NewOrder(zrx, 5, weth, 5)
	t5, _ := factory2.NewTrade(o5, 5)

	tradeService.On("UpdateTradeTxHash", mock.Anything, mock.Anything).Return(nil).Once().Run(func(args mock.Arguments) {
		t1.TxHash = args.Get(1).(common.Hash)
	})

	tradeService.On("UpdateTradeTxHash", mock.Anything, mock.Anything).Return(nil).Once().Run(func(args mock.Arguments) {
		t2.TxHash = args.Get(1).(common.Hash)
	})

	tradeService.On("UpdateTradeTxHash", mock.Anything, mock.Anything).Return(nil).Once().Run(func(args mock.Arguments) {
		t3.TxHash = args.Get(1).(common.Hash)
	})

	tradeService.On("UpdateTradeTxHash", mock.Anything, mock.Anything).Return(nil).Once().Run(func(args mock.Arguments) {
		t4.TxHash = args.Get(1).(common.Hash)
	})

	tradeService.On("UpdateTradeTxHash", mock.Anything, mock.Anything).Return(nil).Once().Run(func(args mock.Arguments) {
		t4.TxHash = args.Get(1).(common.Hash)
	})

	orderService.On("GetByHash", t1.OrderHash).Return(o1, nil)
	orderService.On("GetByHash", t2.OrderHash).Return(o2, nil)
	orderService.On("GetByHash", t3.OrderHash).Return(o3, nil)
	orderService.On("GetByHash", t4.OrderHash).Return(o4, nil)
	orderService.On("GetByHash", t5.OrderHash).Return(o5, nil)
	tradeService.On("GetByHash", t1.Hash).Return(&t1, nil)
	tradeService.On("GetByHash", t2.Hash).Return(&t2, nil)
	tradeService.On("GetByHash", t3.Hash).Return(&t3, nil)
	tradeService.On("GetByHash", t4.Hash).Return(&t4, nil)
	tradeService.On("GetByHash", t5.Hash).Return(&t5, nil)

	defer op.PurgeQueues()

	time.Sleep(2 * time.Millisecond)
	op.QueueTrade(o1, &t1)
	time.Sleep(2 * time.Millisecond)
	op.QueueTrade(o2, &t2)
	time.Sleep(2 * time.Millisecond)
	op.QueueTrade(o3, &t3)
	time.Sleep(2 * time.Millisecond)
	op.QueueTrade(o4, &t4)
	time.Sleep(2 * time.Millisecond)
	op.QueueTrade(o5, &t5)

	wg := sync.WaitGroup{}
	wg.Add(10)

	go func() {
		for {
			msg := <-opMessages
			switch msg.MessageType {
			case "TRADE_SENT_MESSAGE":
				// we simulate that transactions take 100 milliseconds
				go func() {
					time.Sleep(20 * time.Millisecond)
					simulator.Commit()
				}()
				wg.Done()
			case "TRADE_SUCCESS_MESSAGE":
				wg.Done()
			case "TRADE_ERROR_MESSAGE":
				t.Errorf("Received trade error message")
			}
		}
	}()

	wg.Wait()
}

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
