package e2e

import (
	"log"
	"math/big"
	"net/http"
	"sync"
	"testing"
	"time"

	"github.com/Proofsuite/amp-matching-engine/app"
	"github.com/Proofsuite/amp-matching-engine/daos"
	"github.com/Proofsuite/amp-matching-engine/ethereum"
	"github.com/Proofsuite/amp-matching-engine/interfaces"
	"github.com/Proofsuite/amp-matching-engine/rabbitmq"
	"github.com/Proofsuite/amp-matching-engine/redis"
	"github.com/Proofsuite/amp-matching-engine/types"
	"github.com/Proofsuite/amp-matching-engine/utils"
	"github.com/Proofsuite/amp-matching-engine/utils/testutils"
	"github.com/Proofsuite/amp-matching-engine/ws"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
)

type OrderTestSetup struct {
	Wallet *types.Wallet
	Client *testutils.Client
}

func SetupTest() (
	*types.Wallet,
	*types.Wallet,
	*testutils.Client,
	*testutils.Client,
	*testutils.OrderFactory,
	*testutils.OrderFactory,
	*types.Pair,
	common.Address,
	common.Address,
	interfaces.OrderDao,
	interfaces.TradeDao,
) {
	err := app.LoadConfig("../config", "test")
	if err != nil {
		panic(err)
	}

	log.SetFlags(log.LstdFlags | log.Llongfile)
	log.SetPrefix("\nLOG: ")

	rabbitmq.InitConnection(app.Config.Rabbitmq)
	ethereum.NewWebsocketProvider()
	redisConn := redis.NewRedisConnection(app.Config.Redis)

	defer redisConn.FlushAll()

	_, err = daos.InitSession(nil)
	if err != nil {
		panic(err)
	}

	pairDao := daos.NewPairDao()
	exchangeAddress := common.HexToAddress(app.Config.Ethereum["exchange_address"])
	pair, err := pairDao.GetByTokenSymbols("ZRX", "WETH")
	if err != nil {
		panic(err)
	}

	orderDao := daos.NewOrderDao()
	orderDao.Drop()

	tradeDao := daos.NewTradeDao()
	tradeDao.Drop()

	ZRX := pair.BaseTokenAddress
	WETH := pair.QuoteTokenAddress
	wallet1 := types.NewWalletFromPrivateKey("7c78c6e2f65d0d84c44ac0f7b53d6e4dd7a82c35f51b251d387c2a69df712660")
	wallet2 := types.NewWalletFromPrivateKey("7c78c6e2f65d0d84c44ac0f7b53d6e4dd7a82c35f51b251d387c2a69df712661")
	NewRouter()

	//setup mock client
	client1 := testutils.NewClient(wallet1, http.HandlerFunc(ws.ConnectionEndpoint))
	client2 := testutils.NewClient(wallet2, http.HandlerFunc(ws.ConnectionEndpoint))
	client1.Start()
	client2.Start()

	factory1, err := testutils.NewOrderFactory(pair, wallet1, exchangeAddress)
	if err != nil {
		panic(err)
	}

	factory2, err := testutils.NewOrderFactory(pair, wallet2, exchangeAddress)
	if err != nil {
		panic(err)
	}

	return wallet1, wallet2, client1, client2, factory1, factory2, pair, ZRX, WETH, orderDao, tradeDao
}

func TestBuyOrder(t *testing.T) {
	_, _, client1, _, factory1, _, _, ZRX, WETH, orderDao, _ := SetupTest()
	m1, o1, err := factory1.NewOrderMessage(ZRX, 1, WETH, 1)
	if err != nil {
		t.Errorf("Could not create new order message: %v", err)
	}

	client1.Requests <- m1
	time.Sleep(time.Second)
	wg := sync.WaitGroup{}
	wg.Add(1)

	go func() {
		for {
			select {
			case l := <-client1.Logs:
				switch l.MessageType {
				case "ORDER_ADDED":
					wg.Done()
				case "ERROR":
					t.Errorf("Received an error")
				}
			}
		}
	}()

	wg.Wait()

	dbo1, _ := orderDao.GetByHash(o1.Hash)

	assert.Equal(t, "1000000", dbo1.PricePoint)
	assert.Equal(t, "BUY", dbo1.Side)
	assert.Equal(t, "OPEN", dbo1.Status)
	assert.Equal(t, "ZRX/WETH", dbo1.PairName)
	assert.Equal(t, ZRX, dbo1.BaseToken)
	assert.Equal(t, WETH, dbo1.QuoteToken)
	assert.Equal(t, big.NewInt(1), dbo1.BuyAmount)
	assert.Equal(t, big.NewInt(1), dbo1.SellAmount)
	assert.Equal(t, o1.Signature, dbo1.Signature)

	utils.PrintJSON(dbo1)
}

func TestBuyAndCancelOrder(t *testing.T) {
	_, _, client1, client2, factory1, factory2, _, ZRX, WETH, orderDao, _ := SetupTest()
	m1, o1, err := factory1.NewOrderMessage(ZRX, 1, WETH, 1)
	if err != nil {
		t.Errorf("Error creating order message: %v", err)
	}

	m2, _, err := factory2.NewCancelOrderMessage(o1)
	if err != nil {
		t.Errorf("Error creating cancel order message: %v", err)
	}

	//We put a millisecond delay between both requests to ensure they are
	//received in the same order for each test
	client1.Requests <- m1
	time.Sleep(time.Second)
	client2.Requests <- m2
	time.Sleep(time.Millisecond)

	time.Sleep(time.Second)
	wg := sync.WaitGroup{}
	wg.Add(2)

	go func() {
		for {
			select {
			case l := <-client1.Logs:
				switch l.MessageType {
				case "ORDER_ADDED":
					wg.Done()
				case "ORDER_CANCELLED":
					wg.Done()
				case "ERROR":
					t.Errorf("Received an error")
				}
			}
		}
	}()

	wg.Wait()

	dbo1, _ := orderDao.GetByHash(o1.Hash)
	assert.Equal(t, big.NewInt(1000000), dbo1.PricePoint)
	assert.Equal(t, "BUY", dbo1.Side)
	assert.Equal(t, "CANCELLED", dbo1.Status)
	assert.Equal(t, big.NewInt(0), dbo1.FilledAmount)
}

func TestMatchOrder(t *testing.T) {
	_, _, client1, client2, factory1, factory2, _, ZRX, WETH, orderDao, tradeDao := SetupTest()
	m1, o1, _ := factory1.NewOrderMessage(ZRX, 1e10, WETH, 1e10)
	m2, o2, _ := factory2.NewOrderMessage(WETH, 1e10, ZRX, 1e10)

	//We put a millisecond delay between both requests to ensure they are
	//received in the same order for each test
	client1.Requests <- m1
	time.Sleep(500 * time.Millisecond)
	client2.Requests <- m2
	time.Sleep(time.Millisecond)

	wg := sync.WaitGroup{}
	wg.Add(6)

	go func() {
		for {
			select {
			case l := <-client1.Logs:
				log.Print(l)
				switch l.MessageType {
				case "ORDER_ADDED":
					wg.Done()
				case "REQUEST_SIGNATURE":
					t.Error("Maker should not receive REQUEST_SIGNATURE message")
				case "ORDER_SUCCESS":
					wg.Done()
				case "ORDER_PENDING":
					wg.Done()
				case "ERROR":
					t.Errorf("Received an error")
				}
			}
		}
	}()

	go func() {
		for {
			select {
			case l := <-client2.Logs:
				log.Print(l)
				switch l.MessageType {
				case "ORDER_ADDED":
					wg.Done()
				case "ORDER_PENDING":
					wg.Done()
				case "ORDER_SUCCESS":
					wg.Done()
				case "REQUEST_SIGNATURE":
					wg.Done()
				case "ERROR":
					t.Errorf("Received an error")
				}
			}
		}
	}()

	wg.Wait()

	dbo1, _ := orderDao.GetByHash(o1.Hash)
	dbo2, _ := orderDao.GetByHash(o2.Hash)
	trades, _ := tradeDao.GetAll()

	assert.Equal(t, big.NewInt(1000000), dbo1.PricePoint)
	assert.Equal(t, big.NewInt(1000000), dbo1.PricePoint)
	assert.Equal(t, "BUY", dbo1.Side)
	assert.Equal(t, "SELL", dbo2.Side)
	assert.Equal(t, "FILLED", dbo1.Status)
	assert.Equal(t, "FILLED", dbo2.Status)
	assert.Equal(t, big.NewInt(1e10), dbo1.FilledAmount)
	assert.Equal(t, big.NewInt(1e10), dbo2.FilledAmount)

	assert.Equal(t, 1, len(trades))
	assert.Equal(t, o1.Hash, trades[0].OrderHash)
	assert.Equal(t, o1.UserAddress, trades[0].Maker)
	assert.Equal(t, o2.UserAddress, trades[0].Taker)
	assert.Equal(t, o2.Hash, trades[0].TakerOrderHash)
	assert.Equal(t, big.NewInt(1e10), trades[0].Amount)
	assert.Equal(t, big.NewInt(1000000), trades[0].PricePoint)

	//TODO fix this
	// assert.Equal(t, "SUCCESS", trades[0].Status)
}

func TestMatchPartialOrder1(t *testing.T) {
	_, _, client1, client2, factory1, factory2, _, ZRX, WETH, orderDao, tradeDao := SetupTest()
	m1, o1, _ := factory1.NewOrderMessage(ZRX, 1e10, WETH, 1e10)
	m2, o2, _ := factory2.NewOrderMessage(WETH, 2e10, ZRX, 2e10)

	client1.Requests <- m1
	time.Sleep(200 * time.Millisecond)
	client2.Requests <- m2
	time.Sleep(time.Millisecond)

	wg := sync.WaitGroup{}
	wg.Add(7)

	go func() {
		for {
			select {
			case l := <-client1.Logs:
				log.Print(l)
				switch l.MessageType {
				case "ORDER_ADDED":
					wg.Done()
				case "REQUEST_SIGNATURE":
					t.Error("Maker should not receive REQUEST_SIGNATURE message")
				case "ORDER_SUCCESS":
					wg.Done()
				case "ORDER_PENDING":
					wg.Done()
				case "ERROR":
					t.Errorf("Received an error")
				}
			}
		}
	}()

	go func() {
		for {
			select {
			case l := <-client2.Logs:
				log.Print(l)
				switch l.MessageType {
				case "ORDER_ADDED":
					wg.Done()
				case "ORDER_PENDING":
					wg.Done()
				case "ORDER_SUCCESS":
					wg.Done()
				case "REQUEST_SIGNATURE":
					wg.Done()
				case "ERROR":
					t.Errorf("Received an error")
				}
			}
		}
	}()

	wg.Wait()

	dbo1, _ := orderDao.GetByHash(o1.Hash)
	dbo2, _ := orderDao.GetByHash(o2.Hash)
	trades, _ := tradeDao.GetAll()

	assert.Equal(t, big.NewInt(1000000), dbo1.PricePoint)
	assert.Equal(t, big.NewInt(1000000), dbo1.PricePoint)
	assert.Equal(t, "BUY", dbo1.Side)
	assert.Equal(t, "SELL", dbo2.Side)
	assert.Equal(t, "FILLED", dbo1.Status)
	assert.Equal(t, "REPLACED", dbo2.Status)
	assert.Equal(t, big.NewInt(1e10), dbo1.FilledAmount)
	assert.Equal(t, big.NewInt(1e10), dbo2.FilledAmount)

	assert.Equal(t, 1, len(trades))
	assert.Equal(t, o1.Hash, trades[0].OrderHash)
	assert.Equal(t, o1.UserAddress, trades[0].Maker)
	assert.Equal(t, o2.UserAddress, trades[0].Taker)
	assert.Equal(t, o2.Hash, trades[0].TakerOrderHash)
	assert.Equal(t, big.NewInt(1e10), trades[0].Amount)
	assert.Equal(t, big.NewInt(1000000), trades[0].PricePoint)
}

func TestMatchPartialOrder2(t *testing.T) {
	_, _, client1, client2, factory1, factory2, pair, ZRX, WETH, _, _ := SetupTest()
	m1, o1, _ := factory1.NewOrderMessage(WETH, 2e18, ZRX, 2e18)
	m2, o2, _ := factory2.NewOrderMessage(ZRX, 1e18, WETH, 1e18)
	m3, o3, _ := factory2.NewOrderMessage(ZRX, 5e17, WETH, 5e17)

	client1.Requests <- m1
	time.Sleep(200 * time.Millisecond)
	client2.Requests <- m2
	time.Sleep(200 * time.Millisecond)
	client2.Requests <- m3

	wg := sync.WaitGroup{}
	wg.Add(4)

	go func() {
		for {
			select {
			case l1 := <-client1.Logs:
				log.Print(l1)
				switch l1.MessageType {
				case "ORDER_ADDED":
					wg.Done()
				case "REQUEST_SIGNATURE":
					wg.Done()
				case "ERROR":
					t.Errorf("Received an error")
				}
			case l2 := <-client2.Logs:
				log.Print(l2)
				switch l2.MessageType {
				case "ORDER_ADDED":
					wg.Done()
				case "REQUEST_SIGNATURE":
					wg.Done()
				case "ERROR":
					t.Errorf("Received an error")
				}
			}
		}
	}()

	wg.Wait()

	t1 := &types.Trade{
		Amount:     big.NewInt(1e18),
		BaseToken:  ZRX,
		QuoteToken: WETH,
		PricePoint: big.NewInt(1e8),
		OrderHash:  o1.Hash,
		Side:       "BUY",
		PairName:   "ZRX/WETH",
		Maker:      factory1.GetAddress(),
		Taker:      factory2.GetAddress(),
		TradeNonce: big.NewInt(0),
	}

	t2 := &types.Trade{
		Amount:     big.NewInt(5e17),
		BaseToken:  ZRX,
		QuoteToken: WETH,
		PricePoint: big.NewInt(1e8),
		OrderHash:  o1.Hash,
		Side:       "BUY",
		PairName:   "ZRX/WETH",
		Maker:      factory1.GetAddress(),
		Taker:      factory2.GetAddress(),
		TradeNonce: big.NewInt(0),
	}

	t1.Hash = t1.ComputeHash()
	t2.Hash = t2.ComputeHash()

	//Responses received by the first client
	expres1 := types.NewOrderAddedWebsocketMessage(o1, pair, 0)
	// Responses received by the second client
	expres2 := types.NewRequestSignaturesWebsocketMessage(o2.Hash, []*types.OrderTradePair{{o1, t1}}, nil)
	expres3 := types.NewRequestSignaturesWebsocketMessage(o3.Hash, []*types.OrderTradePair{{o1, t2}}, nil)

	testutils.Compare(t, expres1, client1.ResponseLogs[0])
	testutils.Compare(t, expres2, client2.ResponseLogs[0])
	testutils.Compare(t, expres3, client2.ResponseLogs[1])
}

func TestMatchPartialOrder3(t *testing.T) {
	_, _, client1, client2, factory1, factory2, pair, ZRX, WETH, _, _ := SetupTest()
	m1, o1, _ := factory1.NewOrderMessage(WETH, 1e18, ZRX, 1e18)
	m2, o2, _ := factory2.NewOrderMessage(ZRX, 2e18, WETH, 2e18)

	client1.Requests <- m1
	time.Sleep(200 * time.Millisecond)
	client2.Requests <- m2

	wg := sync.WaitGroup{}
	wg.Add(3)

	go func() {
		for {
			select {
			case l1 := <-client1.Logs:
				switch l1.MessageType {
				case "ORDER_ADDED":
					wg.Done()
				case "REQUEST_SIGNATURE":
					wg.Done()
				case "ERROR":
					t.Errorf("Received an error")
				}
			case l2 := <-client2.Logs:
				switch l2.MessageType {
				case "ORDER_ADDED":
					wg.Done()
				case "REQUEST_SIGNATURE":
					wg.Done()
				case "ERROR":
					t.Errorf("Received an error")
				}
			}
		}
	}()

	wg.Wait()

	t1 := &types.Trade{
		Amount:     big.NewInt(1e18),
		BaseToken:  ZRX,
		QuoteToken: WETH,
		PricePoint: big.NewInt(1e8),
		OrderHash:  o1.Hash,
		Side:       "BUY",
		PairName:   "ZRX/WETH",
		Maker:      factory1.GetAddress(),
		Taker:      factory2.GetAddress(),
		TradeNonce: big.NewInt(0),
	}

	ro1 := &types.Order{
		Amount:          big.NewInt(1e18),
		BaseToken:       ZRX,
		QuoteToken:      WETH,
		BuyToken:        ZRX,
		SellToken:       WETH,
		BuyAmount:       big.NewInt(1e18),
		SellAmount:      big.NewInt(1e18),
		FilledAmount:    big.NewInt(0),
		ExchangeAddress: factory2.GetExchangeAddress(),
		UserAddress:     factory2.GetAddress(),
		PricePoint:      big.NewInt(1e8),
		Side:            "BUY",
		PairName:        "ZRX/WETH",
		Status:          "OPEN",
		TakeFee:         big.NewInt(0),
		MakeFee:         big.NewInt(0),
		Expires:         o1.Expires,
	}

	t1.Hash = t1.ComputeHash()

	//Responses received by the first client
	res1 := types.NewOrderAddedWebsocketMessage(o1, pair, 0)
	// Responses received by the second client
	res2 := types.NewRequestSignaturesWebsocketMessage(o2.Hash, []*types.OrderTradePair{{o1, t1}}, ro1)

	testutils.Compare(t, res1, client1.ResponseLogs[0])
	testutils.Compare(t, res2, client2.ResponseLogs[0])
}

func TestMatchPartialOrder4(t *testing.T) {
	_, _, client1, client2, factory1, factory2, pair, ZRX, WETH, _, _ := SetupTest()
	m1, o1, _ := factory1.NewOrderMessage(WETH, 1e18, ZRX, 1e18)
	m2, o2, _ := factory1.NewOrderMessage(WETH, 1e18, ZRX, 1e18)
	m3, o3, _ := factory2.NewOrderMessage(ZRX, 3e18, WETH, 3e18)

	client1.Requests <- m1
	time.Sleep(1000 * time.Millisecond)
	client1.Requests <- m2
	time.Sleep(1000 * time.Millisecond)
	client2.Requests <- m3

	wg := sync.WaitGroup{}
	wg.Add(4)

	go func() {
		for {
			select {
			case l1 := <-client1.Logs:
				log.Print(l1)
				switch l1.MessageType {
				case "ORDER_ADDED":
					wg.Done()
				case "REQUEST_SIGNATURE":
					wg.Done()
				case "ERROR":
					t.Errorf("Received an error")
				}
			case l2 := <-client2.Logs:
				log.Print(l2)
				switch l2.MessageType {
				case "ORDER_ADDED":
					wg.Done()
				case "REQUEST_SIGNATURE":
					wg.Done()
				case "ERROR":
					t.Errorf("Received an error")
				}
			}
		}
	}()

	wg.Wait()

	t1 := &types.Trade{
		Amount:     big.NewInt(1e18),
		BaseToken:  ZRX,
		QuoteToken: WETH,
		PricePoint: big.NewInt(1e8),
		OrderHash:  o1.Hash,
		Side:       "BUY",
		PairName:   "ZRX/WETH",
		Maker:      factory1.GetAddress(),
		Taker:      factory2.GetAddress(),
		TradeNonce: big.NewInt(0),
	}

	t2 := &types.Trade{
		Amount:     big.NewInt(1e18),
		BaseToken:  ZRX,
		QuoteToken: WETH,
		PricePoint: big.NewInt(1e8),
		OrderHash:  o2.Hash,
		Side:       "BUY",
		PairName:   "ZRX/WETH",
		Maker:      factory1.GetAddress(),
		Taker:      factory2.GetAddress(),
		TradeNonce: big.NewInt(0),
	}

	//Remaining order
	ro1 := &types.Order{
		Amount:          big.NewInt(1e18),
		BaseToken:       ZRX,
		QuoteToken:      WETH,
		BuyToken:        ZRX,
		SellToken:       WETH,
		BuyAmount:       big.NewInt(1e18),
		SellAmount:      big.NewInt(1e18),
		FilledAmount:    big.NewInt(0),
		ExchangeAddress: factory2.GetExchangeAddress(),
		UserAddress:     factory2.GetAddress(),
		PricePoint:      big.NewInt(1e8),
		Side:            "BUY",
		PairName:        "ZRX/WETH",
		Status:          "OPEN",
		TakeFee:         big.NewInt(0),
		MakeFee:         big.NewInt(0),
		Expires:         o1.Expires,
	}

	t1.Hash = t1.ComputeHash()
	t2.Hash = t2.ComputeHash()

	res1 := types.NewOrderAddedWebsocketMessage(o1, pair, 0)
	res2 := types.NewOrderAddedWebsocketMessage(o2, pair, 0)
	res3 := types.NewRequestSignaturesWebsocketMessage(o3.Hash, []*types.OrderTradePair{{o1, t1}, {o2, t2}}, ro1)
	testutils.Compare(t, res1, client1.ResponseLogs[0])
	testutils.Compare(t, res2, client1.ResponseLogs[1])
	testutils.Compare(t, res3, client2.ResponseLogs[0])
}

func TestMatchPartialOrder5(t *testing.T) {
	_, _, client1, client2, factory1, factory2, pair, ZRX, WETH, _, _ := SetupTest()
	m1, o1, _ := factory1.NewBuyOrderMessage(50, 10) // buy 1e18 ZRX at 1ZRX = 50WETH
	m2, o2, _ := factory2.NewSellOrderMessage(50, 10)

	client1.Requests <- m1
	time.Sleep(200 * time.Millisecond)
	client2.Requests <- m2
	time.Sleep(200 * time.Millisecond)

	wg := sync.WaitGroup{}
	wg.Add(4)

	go func() {
		for {
			select {
			case l1 := <-client1.Logs:
				log.Print(l1)
				switch l1.MessageType {
				case "ORDER_ADDED":
					wg.Done()
				case "REQUEST_SIGNATURE":
					wg.Done()
				case "ERROR":
					t.Errorf("Received an error")
				}
			case l2 := <-client2.Logs:
				log.Print(l2)
				switch l2.MessageType {
				case "ORDER_ADDED":
					wg.Done()
				case "REQUEST_SIGNATURE":
					wg.Done()
				case "ERROR":
					t.Errorf("Received an error")
				}
			}
		}
	}()

	wg.Wait()

	res1 := types.NewOrderAddedWebsocketMessage(o1, pair, 0)

	t1 := &types.Trade{
		Amount:     big.NewInt(1e18),
		BaseToken:  ZRX,
		QuoteToken: WETH,
		PricePoint: big.NewInt(50e8),
		OrderHash:  o1.Hash,
		Side:       "SELL",
		PairName:   "ZRX/WETH",
		Maker:      factory1.GetAddress(),
		Taker:      factory2.GetAddress(),
		TradeNonce: big.NewInt(0),
	}

	t1.Hash = t1.ComputeHash()
	res2 := types.NewRequestSignaturesWebsocketMessage(o2.Hash, []*types.OrderTradePair{{o1, t1}}, nil)

	testutils.Compare(t, res1, client1.ResponseLogs[0])
	testutils.Compare(t, res2, client2.ResponseLogs[0])

}

func TestMatchPartialOrder6(t *testing.T) {
	_, _, client1, client2, factory1, factory2, pair, _, _, _, _ := SetupTest()
	m1, o1, _ := factory1.NewSellOrderMessage(51, 1000) // buy 1e18 ZRX at 1ZRX = 49WETH
	m2, o2, _ := factory2.NewBuyOrderMessage(49, 1000)  // sell 1e18 ZRX at 1ZRX = 51WETH

	client1.Requests <- m1
	time.Sleep(200 * time.Millisecond)
	client2.Requests <- m2
	time.Sleep(200 * time.Millisecond)

	wg := sync.WaitGroup{}
	wg.Add(4)

	go func() {
		for {
			select {
			case l1 := <-client1.Logs:
				switch l1.MessageType {
				case "ORDER_ADDED":
					wg.Done()
				case "REQUEST_SIGNATURE":
					wg.Done()
				case "ERROR":
					t.Errorf("Received an error")
				}
			case l2 := <-client2.Logs:
				switch l2.MessageType {
				case "ORDER_ADDED":
					wg.Done()
				case "REQUEST_SIGNATURE":
					wg.Done()
				case "ERROR":
					t.Errorf("Received an error")
				}
			}
		}
	}()

	wg.Wait()

	res1 := types.NewOrderAddedWebsocketMessage(o1, pair, 0)
	res2 := types.NewOrderAddedWebsocketMessage(o2, pair, 0)
	testutils.Compare(t, res1, client1.ResponseLogs[0])
	testutils.Compare(t, res2, client2.ResponseLogs[0])
}

func TestOrders1(t *testing.T) {
	_, wallet2, client1, client2, factory1, factory2, pair, ZRX, WETH, _, _ := SetupTest()
	m1, o1, _ := factory1.NewOrderMessage(WETH, 1e18, ZRX, 1e18)
	m2, o2, _ := factory1.NewOrderMessage(WETH, 1e18, ZRX, 1e18)
	m3, o3, _ := factory2.NewOrderMessage(ZRX, 3e18, WETH, 3e18)

	// we simulated the order has an invalid signature
	o2.Sign(wallet2)
	m2 = types.NewOrderWebsocketMessage(o2)

	client1.Requests <- m1
	time.Sleep(1000 * time.Millisecond)
	client1.Requests <- m2
	time.Sleep(1000 * time.Millisecond)
	client2.Requests <- m3

	wg := sync.WaitGroup{}
	wg.Add(4)

	go func() {
		for {
			select {
			case l1 := <-client1.Logs:
				log.Print(l1)
				switch l1.MessageType {
				case "ORDER_ADDED":
					wg.Done()
				case "REQUEST_SIGNATURE":
					wg.Done()
				case "ERROR":
					t.Errorf("Received an error")
				}
			case l2 := <-client2.Logs:
				log.Print(l2)
				switch l2.MessageType {
				case "ORDER_ADDED":
					wg.Done()
				case "REQUEST_SIGNATURE":
					wg.Done()
				case "ERROR":
					t.Errorf("Received an error")
				}
			}
		}
	}()

	wg.Wait()

	t1 := &types.Trade{
		Amount:     big.NewInt(1e18),
		BaseToken:  ZRX,
		QuoteToken: WETH,
		PricePoint: big.NewInt(1e8),
		OrderHash:  o1.Hash,
		Side:       "BUY",
		PairName:   "ZRX/WETH",
		Maker:      factory1.GetAddress(),
		Taker:      factory2.GetAddress(),
		TradeNonce: big.NewInt(0),
	}

	t2 := &types.Trade{
		Amount:     big.NewInt(1e18),
		BaseToken:  ZRX,
		QuoteToken: WETH,
		PricePoint: big.NewInt(1e8),
		OrderHash:  o2.Hash,
		Side:       "BUY",
		PairName:   "ZRX/WETH",
		Maker:      factory1.GetAddress(),
		Taker:      factory2.GetAddress(),
		TradeNonce: big.NewInt(0),
	}

	//Remaining order
	ro1 := &types.Order{
		Amount:          big.NewInt(1e18),
		BaseToken:       ZRX,
		QuoteToken:      WETH,
		BuyToken:        ZRX,
		SellToken:       WETH,
		BuyAmount:       big.NewInt(1e18),
		SellAmount:      big.NewInt(1e18),
		FilledAmount:    big.NewInt(0),
		ExchangeAddress: factory2.GetExchangeAddress(),
		UserAddress:     factory2.GetAddress(),
		PricePoint:      big.NewInt(1e8),
		Side:            "BUY",
		PairName:        "ZRX/WETH",
		Status:          "OPEN",
		TakeFee:         big.NewInt(0),
		MakeFee:         big.NewInt(0),
		Expires:         o1.Expires,
	}

	t1.Hash = t1.ComputeHash()
	t2.Hash = t2.ComputeHash()

	res1 := types.NewOrderAddedWebsocketMessage(o1, pair, 0)
	res2 := types.NewOrderAddedWebsocketMessage(o2, pair, 0)
	res3 := types.NewRequestSignaturesWebsocketMessage(o3.Hash, []*types.OrderTradePair{{o1, t1}, {o2, t2}}, ro1)
	testutils.Compare(t, res1, client1.ResponseLogs[0])
	testutils.Compare(t, res2, client1.ResponseLogs[1])
	testutils.Compare(t, res3, client2.ResponseLogs[0])
}

func TestInvalidPairOrder(t *testing.T) {
	_, _, client1, _, factory1, _, _, ZRX, _, _, _ := SetupTest()
	m1, _, err := factory1.NewOrderMessage(ZRX, 1, ZRX, 1)
	if err != nil {
		t.Errorf("Could not create new order message: %v", err)
	}

	client1.Requests <- m1

	time.Sleep(time.Second)
	wg := sync.WaitGroup{}
	wg.Add(1)

	go func() {
		for {
			select {
			case l := <-client1.Logs:
				switch l.MessageType {
				case "ORDER_ADDED":
					wg.Done()
				case "ERROR":
					utils.PrintJSON(l)
					t.Errorf("Received an error")
				}
			}
		}
	}()
	wg.Wait()
}

func TestInvalidAmountOrder(t *testing.T) {
	_, _, client1, _, factory1, _, _, ZRX, _, _, _ := SetupTest()
	m1, _, err := factory1.NewOrderMessage(ZRX, 1, ZRX, -1)
	if err != nil {
		t.Errorf("Could not create new order message: %v", err)
	}

	client1.Requests <- m1

	time.Sleep(time.Second)
	wg := sync.WaitGroup{}
	wg.Add(1)

	go func() {
		for {
			select {
			case l := <-client1.Logs:
				switch l.MessageType {
				case "ORDER_ADDED":
					wg.Done()
				case "ERROR":
					utils.PrintJSON(l)
					t.Errorf("Received an error")
				}
			}
		}
	}()
	wg.Wait()
}

func TestInvalidNonceOrder(t *testing.T) {
	_, _, client1, _, factory1, _, _, ZRX, _, _, _ := SetupTest()
	m1, o1, err := factory1.NewOrderMessage(ZRX, 1, ZRX, 1)
	if err != nil {
		t.Errorf("Could not create new order message: %v", err)
	}

	o1.Nonce = big.NewInt(-1)
	m1.Payload.Data = o1

	client1.Requests <- m1

	time.Sleep(time.Second)
	wg := sync.WaitGroup{}
	wg.Add(1)

	go func() {
		for {
			select {
			case l := <-client1.Logs:
				switch l.MessageType {
				case "ORDER_ADDED":
					wg.Done()
				case "ERROR":
					utils.PrintJSON(l)
					t.Errorf("Received an error")
				}
			}
		}
	}()
	wg.Wait()
}

func TestInvalidExchangeAddress(t *testing.T) {
	_, _, client1, _, factory1, _, _, ZRX, _, _, _ := SetupTest()
	m1, o1, err := factory1.NewOrderMessage(ZRX, 1, ZRX, 1)
	if err != nil {
		t.Errorf("Could not create new order message: %v", err)
	}

	o1.ExchangeAddress = common.HexToAddress("0x1")
	m1.Payload.Data = o1

	client1.Requests <- m1

	time.Sleep(time.Second)
	wg := sync.WaitGroup{}
	wg.Add(1)

	go func() {
		for {
			select {
			case l := <-client1.Logs:
				switch l.MessageType {
				case "ORDER_ADDED":
					wg.Done()
				case "ERROR":
					utils.PrintJSON(l)
					t.Errorf("Received an error")
				}
			}
		}
	}()
	wg.Wait()
}
