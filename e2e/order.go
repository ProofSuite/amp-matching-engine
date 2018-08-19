package e2e

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/Proofsuite/amp-matching-engine/app"
	"github.com/Proofsuite/amp-matching-engine/daos"
	"github.com/Proofsuite/amp-matching-engine/mocks"
	"github.com/Proofsuite/amp-matching-engine/types"
	"github.com/ethereum/go-ethereum/common"
)

type OrderTestSetup struct {
	Wallet *types.Wallet
	Client *mocks.Client
}

func TestSocketBuyOrder(t *testing.T) {
	wg := sync.WaitGroup{}
	wg.Add(1)

	err := app.LoadConfig("./config")
	if err != nil {
		t.Errorf("Could not load configuration: %v", err)
	}

	_, err = daos.InitSession()
	if err != nil {
		t.Errorf("Could not load db session")
	}

	wallet := types.NewWallet()
	router := buildRouter()
	client := mocks.NewClient(wallet, router)

	accountDao := daos.NewAccountDao()
	pairDao := daos.NewPairDao()
	exchangeAddress := common.HexToAddress("0x")

	client.start()

	pair, err := pairDao.GetByName("ZRX/WETH")
	if err != nil {
		fmt.Printf("Could not retrueve token pair for test")
	}

	factory := mocks.NewOrderFactory(pair, wallet, exchangeAddress)
	m1, o1, _ := factory.NewOrderMessage()
	client.requests <- m1

	time.Sleep(time.Second)

	go func() {
		for {
			select {
			case l := <-makerClient.logs:
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

func TestSocketOrderFill(t *testing.T) {
	wg := sync.WaitGroup{}
	wg.Add(2)

	wallet1 := types.NewWallet()
	wallet2 := types.NewWallet()

	router := buildRouter()

	client1 := mocks.NewClient(wallet1, router)
	client2 := mocks.NewClient(wallet2, router)

	accountDao := daos.NewAccountDao()
	pairDao := daos.NewPairDao()
	exchangeAddress := common.HexToAddress("0x")

	client1.start()
	client2.start()

	pair, err := pairDao.GetByName("ZRX/WETH")
	if err != nil {
		fmt.Printf("Could not retrieve token pair for test")
	}

	ZRX := pair.BaseTokenAddress
	WETH := pair.QuoteTokenAddress

	factory1 := mocks.NewOrderFactory(pair, wallet1, exchangeAddress)
	factory2 := mocks.NewOrderFactory(pair, wallet2, exchangeAddress)

	m1, o1, _ := mocks.NewOrderMessage(ZRX, 1, WETH, 1)
	m2, o2, _ := mocks.NewOrderMessage(WETH, 1, ZRX, 1)

	//We put a millisecond delay between both requests to ensure they are
	//received in the same order for each test
	client1.requests <- m1
	time.Sleep(time.Millisecond)
	client2.requests <- m2
	time.Sleep(time.Millisecond)

	go func() {
		for {
			select {
			case l := <-makerClient.logs:
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
