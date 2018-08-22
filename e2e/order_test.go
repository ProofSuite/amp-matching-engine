package e2e

import (
	"log"
	"net/http"
	"sync"
	"testing"
	"time"

	"github.com/Proofsuite/amp-matching-engine/app"
	"github.com/Proofsuite/amp-matching-engine/daos"
	"github.com/Proofsuite/amp-matching-engine/ethereum"
	"github.com/Proofsuite/amp-matching-engine/mocks"
	"github.com/Proofsuite/amp-matching-engine/rabbitmq"
	"github.com/Proofsuite/amp-matching-engine/redis"
	"github.com/Proofsuite/amp-matching-engine/types"
	"github.com/Proofsuite/amp-matching-engine/ws"
	"github.com/ethereum/go-ethereum/common"
)

type OrderTestSetup struct {
	Wallet *types.Wallet
	Client *mocks.Client
}

func SetupTest() (*types.Wallet, *types.Wallet, *mocks.Client, *mocks.Client, *mocks.OrderFactory, *mocks.OrderFactory, *types.Pair, common.Address, common.Address) {
	err := app.LoadConfig("../config")
	if err != nil {
		panic(err)
	}

	log.SetFlags(log.LstdFlags | log.Llongfile)
	log.SetPrefix("\nLOG: ")

	rabbitmq.InitConnection(app.Config.Rabbitmq)
	ethereum.InitConnection(app.Config.Ethereum)
	redis.InitConnection(app.Config.Redis)

	_, err = daos.InitSession()
	if err != nil {
		panic(err)
	}

	pairDao := daos.NewPairDao()
	exchangeAddress := common.HexToAddress("0x")
	pair, err := pairDao.GetByTokenSymbols("ZRX", "WETH")
	if err != nil {
		panic(err)
	}

	ZRX := pair.BaseTokenAddress
	WETH := pair.QuoteTokenAddress
	wallet1 := types.NewWalletFromPrivateKey("7c78c6e2f65d0d84c44ac0f7b53d6e4dd7a82c35f51b251d387c2a69df712660")
	wallet2 := types.NewWalletFromPrivateKey("7c78c6e2f65d0d84c44ac0f7b53d6e4dd7a82c35f51b251d387c2a69df712661")
	NewRouter()

	//setup mock client
	client1 := mocks.NewClient(wallet1, http.HandlerFunc(ws.ConnectionEndpoint))
	client2 := mocks.NewClient(wallet2, http.HandlerFunc(ws.ConnectionEndpoint))
	client1.Start()
	client2.Start()

	factory1, err := mocks.NewOrderFactory(pair, wallet1, exchangeAddress)
	if err != nil {
		panic(err)
	}

	factory2, err := mocks.NewOrderFactory(pair, wallet2, exchangeAddress)
	if err != nil {
		panic(err)
	}

	return wallet1, wallet2, client1, client2, factory1, factory2, pair, ZRX, WETH
}

func TestBuyOrder(t *testing.T) {
	_, _, client1, _, factory1, _, _, ZRX, WETH := SetupTest()
	m1, _, err := factory1.NewOrderMessage(ZRX, 1, WETH, 1)
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
}

func TestBuyAndCancelOrder(t *testing.T) {
	_, _, client1, client2, factory1, factory2, _, ZRX, WETH := SetupTest()
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
}

func TestMatchOrder(t *testing.T) {
	_, _, client1, client2, factory1, factory2, _, ZRX, WETH := SetupTest()
	m1, _, _ := factory1.NewOrderMessage(ZRX, 1e18, WETH, 1e18)
	m2, _, _ := factory2.NewOrderMessage(WETH, 1e18, ZRX, 1e18)

	//We put a millisecond delay between both requests to ensure they are
	//received in the same order for each test
	client1.Requests <- m1
	time.Sleep(time.Millisecond)
	client2.Requests <- m2
	time.Sleep(time.Millisecond)

	wg := sync.WaitGroup{}
	wg.Add(2)

	go func() {
		for {
			select {
			case l := <-client1.Logs:
				switch l.MessageType {
				case "ORDER_ADDED":
					wg.Done()
				case "ORDER_MATCHED":
					wg.Done()
				case "ERROR":
					t.Errorf("Received an error")
				}
			}
		}
	}()

	wg.Wait()
}

// func TestBuyOrder(t *testing.T) {
// 	m1, _, err := factory1.NewOrderMessage(ZRX, 1, WETH, 1)
// 	if err != nil {
// 		t.Errorf("Could not create new order message: %v", err)
// 	}

// 	client1.Requests <- m1

// 	time.Sleep(time.Second)
// 	wg := sync.WaitGroup{}
// 	wg.Add(1)

// 	go func() {
// 		for {
// 			select {
// 			case l := <-client1.Logs:
// 				switch l.MessageType {
// 				case "ORDER_ADDED":
// 					wg.Done()
// 				case "ERROR":
// 					t.Errorf("Received an error")
// 				}
// 			}
// 		}
// 	}()

// 	wg.Wait()
// }

// func TestSocketBuyOrder(t *testing.T) {
// 	err := app.LoadConfig("../config")
// 	if err != nil {
// 		t.Errorf("Could not load configuration: %v", err)
// 	}

// 	log.SetFlags(log.LstdFlags | log.Llongfile)
// 	log.SetPrefix("\nLOG: ")

// 	rabbitmq.InitConnection(app.Config.Rabbitmq)
// 	ethereum.InitConnection(app.Config.Ethereum)
// 	redis.InitConnection(app.Config.Redis)

// 	_, err = daos.InitSession()
// 	if err != nil {
// 		t.Errorf("Could not load db session")
// 	}

// 	wallet := types.NewWalletFromPrivateKey("7c78c6e2f65d0d84c44ac0f7b53d6e4dd7a82c35f51b251d387c2a69df712660")

// 	NewRouter()
// 	//setup mock client
// 	client := mocks.NewClient(wallet, http.HandlerFunc(ws.ConnectionEndpoint))
// 	client.Start()

// 	pairDao := daos.NewPairDao()
// 	exchangeAddress := common.HexToAddress("0x")
// 	pair, err := pairDao.GetByTokenSymbols("ZRX", "WETH")
// 	if err != nil {
// 		t.Errorf("Could not retrieve token pair: %v", err)
// 	}

// 	ZRX := pair.BaseTokenAddress
// 	WETH := pair.QuoteTokenAddress

// 	factory, err := mocks.NewOrderFactory(pair, wallet, exchangeAddress)
// 	if err != nil {
// 		t.Errorf("Could not create new factory: %v", err)
// 	}

// 	m1, _, err := factory.NewOrderMessage(ZRX, 1, WETH, 1)
// 	if err != nil {
// 		t.Errorf("Could not create new order message: %v", err)
// 	}

// 	client.Requests <- m1

// 	time.Sleep(time.Second)
// 	wg := sync.WaitGroup{}
// 	wg.Add(1)

// 	go func() {
// 		for {
// 			select {
// 			case l := <-client.Logs:
// 				switch l.MessageType {
// 				case "ORDER_ADDED":
// 					wg.Done()
// 				case "ERROR":
// 					t.Errorf("Received an error")
// 				}
// 			}
// 		}
// 	}()

// 	wg.Wait()
// }

// func TestSocketBuyAndCancelOrder(t *testing.T) {
// 	err := app.LoadConfig("../config")
// 	if err != nil {
// 		t.Errorf("Could not load configuration: %v", err)
// 	}

// 	log.SetFlags(log.LstdFlags | log.Llongfile)
// 	log.SetPrefix("\nLOG: ")

// 	rabbitmq.InitConnection(app.Config.Rabbitmq)
// 	ethereum.InitConnection(app.Config.Ethereum)
// 	redis.InitConnection(app.Config.Redis)

// 	_, err = daos.InitSession()
// 	if err != nil {
// 		t.Errorf("Could not load db session")
// 	}

// 	wallet := types.NewWalletFromPrivateKey("7c78c6e2f65d0d84c44ac0f7b53d6e4dd7a82c35f51b251d387c2a69df712660")

// 	NewRouter()
// 	//setup mock client
// 	client := mocks.NewClient(wallet, http.HandlerFunc(ws.ConnectionEndpoint))
// 	client.Start()

// 	pairDao := daos.NewPairDao()
// 	exchangeAddress := common.HexToAddress("0x")
// 	pair, err := pairDao.GetByTokenSymbols("ZRX", "WETH")
// 	if err != nil {
// 		fmt.Printf("Could not retrieve token pair for test")
// 	}

// 	ZRX := pair.BaseTokenAddress
// 	WETH := pair.QuoteTokenAddress

// 	factory, err := mocks.NewOrderFactory(pair, wallet, exchangeAddress)
// 	if err != nil {
// 		fmt.Printf("Could not create factory")
// 	}

// 	m1, o1, err := factory.NewOrderMessage(ZRX, 1, WETH, 1)
// 	if err != nil {
// 		t.Errorf("Error creating order message: %v", err)
// 	}

// 	m2, _, err := factory.NewCancelOrderMessage(o1)
// 	if err != nil {
// 		t.Errorf("Error creating cancel order message: %v", err)
// 	}

// 	client.Requests <- m1
// 	time.Sleep(time.Second)
// 	client.Requests <- m2
// 	time.Sleep(time.Millisecond)

// 	time.Sleep(time.Second)
// 	wg := sync.WaitGroup{}
// 	wg.Add(2)

// 	go func() {
// 		for {
// 			select {
// 			case l := <-client.Logs:
// 				switch l.MessageType {
// 				case "ORDER_ADDED":
// 					wg.Done()
// 				case "ORDER_CANCELLED":
// 					wg.Done()
// 				case "ERROR":
// 					t.Errorf("Received an error")
// 				}
// 			}
// 		}
// 	}()

// 	wg.Wait()
// }

// func TestSocketOrderFill(t *testing.T) {
// 	err := app.LoadConfig("../config")
// 	if err != nil {
// 		t.Errorf("Could not load configuration: %v", err)
// 	}

// 	log.SetFlags(log.LstdFlags | log.Llongfile)
// 	log.SetPrefix("\nLOG: ")

// 	rabbitmq.InitConnection(app.Config.Rabbitmq)
// 	ethereum.InitConnection(app.Config.Ethereum)
// 	redis.InitConnection(app.Config.Redis)

// 	_, err = daos.InitSession()
// 	if err != nil {
// 		t.Errorf("Could not load db session")
// 	}

// 	pairDao := daos.NewPairDao()
// 	exchangeAddress := common.HexToAddress("0x")
// 	pair, err := pairDao.GetByTokenSymbols("ZRX", "WETH")
// 	if err != nil {
// 		t.Errorf("Could not retrieve token pair: %v", err)
// 	}

// 	ZRX := pair.BaseTokenAddress
// 	WETH := pair.QuoteTokenAddress

// 	wallet1 := types.NewWalletFromPrivateKey("7c78c6e2f65d0d84c44ac0f7b53d6e4dd7a82c35f51b251d387c2a69df712660")
// 	wallet2 := types.NewWalletFromPrivateKey("7c78c6e2f65d0d84c44ac0f7b53d6e4dd7a82c35f51b251d387c2a69df712661")
// 	NewRouter()
// 	client1 := mocks.NewClient(wallet1, http.HandlerFunc(ws.ConnectionEndpoint))
// 	client2 := mocks.NewClient(wallet2, http.HandlerFunc(ws.ConnectionEndpoint))
// 	client1.Start()
// 	client2.Start()

// 	factory1, err := mocks.NewOrderFactory(pair, wallet1, exchangeAddress)
// 	if err != nil {
// 		fmt.Printf("Could not create order factory: %v", err)
// 	}

// 	factory2, err := mocks.NewOrderFactory(pair, wallet2, exchangeAddress)
// 	if err != nil {
// 		fmt.Printf("Could not create order factory: %v", err)
// 	}

// 	m1, _, _ := factory1.NewOrderMessage(ZRX, 1, WETH, 1)
// 	m2, _, _ := factory2.NewOrderMessage(WETH, 1, ZRX, 1)

// 	//We put a millisecond delay between both requests to ensure they are
// 	//received in the same order for each test
// 	client1.Requests <- m1
// 	time.Sleep(time.Millisecond)
// 	client2.Requests <- m2
// 	time.Sleep(time.Millisecond)

// 	wg := sync.WaitGroup{}
// 	wg.Add(2)

// 	go func() {
// 		for {
// 			select {
// 			case l := <-client1.Logs:
// 				switch l.MessageType {
// 				case "ORDER_ADDED":
// 					wg.Done()
// 				case "ERROR":
// 					t.Errorf("Received an error")
// 				}
// 			}
// 		}
// 	}()

// 	wg.Wait()
// }
