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
	"github.com/Proofsuite/amp-matching-engine/rabbitmq"
	"github.com/Proofsuite/amp-matching-engine/redis"
	"github.com/Proofsuite/amp-matching-engine/types"
	"github.com/Proofsuite/amp-matching-engine/utils/testutils"
	"github.com/Proofsuite/amp-matching-engine/ws"
	"github.com/ethereum/go-ethereum/common"
)

type OrderTestSetup struct {
	Wallet *types.Wallet
	Client *testutils.Client
}

func SetupTest() (*types.Wallet, *types.Wallet, *testutils.Client, *testutils.Client, *testutils.OrderFactory, *testutils.OrderFactory, *types.Pair, common.Address, common.Address) {
	err := app.LoadConfig("../config")
	if err != nil {
		panic(err)
	}

	log.SetFlags(log.LstdFlags | log.Llongfile)
	log.SetPrefix("\nLOG: ")

	rabbitmq.InitConnection(app.Config.Rabbitmq)
	ethereum.InitConnection(app.Config.Ethereum)
	redisConn := redis.NewRedisConnection(app.Config.Redis)
	defer redisConn.FlushAll()

	_, err = daos.InitSession(nil)
	if err != nil {
		panic(err)
	}

	pairDao := daos.NewPairDao()
	exchangeAddress := common.HexToAddress("0x")
	pair, err := pairDao.GetByTokenSymbols("ZRX", "WETH")
	if err != nil {
		panic(err)
	}

	orderDao := daos.NewOrderDao()
	orderDao.Drop()

	ZRX := pair.BaseTokenAddress
	WETH := pair.QuoteTokenAddress
	wallet1 := testutils.GetTestWallet1()
	wallet2 := testutils.GetTestWallet2()
	wallet3 := testutils.GetTestWallet3()
	NewRouter()

	//setup mock client
	factory1, err := testutils.NewOrderFactory(pair, wallet1, exchangeAddress)
	if err != nil {
		panic(err)
	}

	deployer, err := testutils.NewSimulator(walletService, txService, []common.Address{
		wallet1.Address,
		wallet2.Address,
		wallet3.Address,
	})

	if err != nil {
		panic(err)
	}

	op, err := operator.NewOperator(
		walletService,
		txService,
		tradeService,
		ethereumService,
		exchange
	)

	// factory2, err := testutils.NewOrderFactory(pair, wallet2, exchangeAddress)
	// if err != nil {
	// 	panic(err)
	// }

	return wallet1, wallet2, client1, client2, factory1, factory2, pair, ZRX, WETH
}

func TestQueueTrade(t *testing.T) {

}

func TestBuyOrder(t *testing.T) {
	_, _, client1, _, factory1, _, _, ZRX, WETH := SetupTest()
	m1, _, err := factory1.NewOrderMessage(ZRX, 1, WETH, 1)
	if err != nil {
		t.Errorf("Could not create new order message: %v", err)
	}



}
