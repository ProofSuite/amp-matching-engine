package engine

import (
	"io/ioutil"
	"log"
	"math/big"
	"testing"

	"github.com/Proofsuite/amp-matching-engine/daos"
	"github.com/Proofsuite/amp-matching-engine/rabbitmq"
	"github.com/Proofsuite/amp-matching-engine/types"
	"github.com/Proofsuite/amp-matching-engine/utils"
	"github.com/Proofsuite/amp-matching-engine/utils/testutils"
	"github.com/Proofsuite/amp-matching-engine/utils/testutils/mocks"
	"github.com/Proofsuite/amp-matching-engine/utils/units"
	"github.com/ethereum/go-ethereum/common"
)

var db *daos.Database

func setupTest() (
	*Engine,
	*OrderBook,
	common.Address,
	*types.Wallet,
	*types.Wallet,
	*types.Pair,
	common.Address,
	common.Address,
	*testutils.OrderFactory,
	*testutils.OrderFactory) {

	log.SetFlags(log.LstdFlags | log.Llongfile)
	log.SetPrefix("\nLOG: ")

	mongoServer := testutils.NewDBTestServer()
	temp, _ := ioutil.TempDir("", "test")
	mongoServer.SetPath(temp)

	session := mongoServer.Session()
	daos.InitSession(session)
	rabbitConn := rabbitmq.InitConnection("amqp://guest:guest@localhost:5672/")

	opts := daos.OrderDaoDBOption("test")
	orderDao := daos.NewOrderDao(opts)
	orderDao.Drop()

	pair := testutils.GetZRXWETHTestPair()
	pairDao := new(mocks.PairDao)
	tradeDao := new(mocks.TradeDao)
	pairDao.On("GetAll").Return([]types.Pair{*pair}, nil)

	eng := NewEngine(rabbitConn, orderDao, tradeDao, pairDao)
	ex := testutils.GetTestAddress1()
	maker := testutils.GetTestWallet1()
	taker := testutils.GetTestWallet2()
	zrx := pair.BaseTokenAddress
	weth := pair.QuoteTokenAddress

	factory1, err := testutils.NewOrderFactory(pair, maker, ex)
	if err != nil {
		panic(err)
	}

	factory2, err := testutils.NewOrderFactory(pair, taker, ex)
	if err != nil {
		panic(err)
	}

	ob := eng.orderbooks[pair.Code()]
	if ob == nil {
		panic("Could not get orderbook")
	}

	return eng, ob, ex, maker, taker, pair, zrx, weth, factory1, factory2
}

func TestSellOrder(t *testing.T) {
	_, ob, _, _, _, _, _, _, factory1, _ := setupTest()

	o1, _ := factory1.NewSellOrder(1e3, 1e8, 0)

	exp1 := o1
	exp1.Status = "OPEN"
	expected := &types.EngineResponse{
		Status:  "ORDER_ADDED",
		Order:   &exp1,
		Matches: nil,
	}

	res, err := ob.sellOrder(&o1)
	if err != nil {
		t.Error("Error in sell order: ", err)
	}

	testutils.CompareEngineResponse(t, expected, res)
}

func TestBuyOrder(t *testing.T) {
	_, ob, _, _, _, _, _, _, factory1, _ := setupTest()

	o1, _ := factory1.NewBuyOrder(1e3, 1e8, 0)

	exp1 := o1
	exp1.Status = "OPEN"
	expected := &types.EngineResponse{
		Status:  "ORDER_ADDED",
		Order:   &exp1,
		Matches: nil,
	}

	res, err := ob.buyOrder(&o1)
	if err != nil {
		t.Error("Error in buy order: ", err)
	}

	testutils.CompareEngineResponse(t, expected, res)
}

func TestFillOrder1(t *testing.T) {
	_, ob, _, _, _, _, _, _, factory1, factory2 := setupTest()

	o1, _ := factory1.NewSellOrder(1e3, 1e8)
	o2, _ := factory2.NewBuyOrder(1e3, 1e8)
	expt1 := types.NewTrade(&o1, &o2, units.Ethers(1e8), big.NewInt(1e3))

	expo1 := o1
	expo1.Status = "OPEN"
	expectedSellOrderResponse := &types.EngineResponse{Status: "ORDER_ADDED", Order: &expo1}

	expo2 := o2
	expo2.Status = "FILLED"
	expo2.FilledAmount = units.Ethers(1e8)

	expo3 := o1
	expo3.Status = "FILLED"
	expo3.FilledAmount = units.Ethers(1e8)

	expectedMatches := types.NewMatches(
		[]*types.Order{&expo3},
		&expo2,
		[]*types.Trade{expt1},
	)

	expectedBuyOrderResponse := &types.EngineResponse{
		Status:  "ORDER_FILLED",
		Order:   &expo2,
		Matches: expectedMatches,
	}

	sellOrderResponse, err := ob.sellOrder(&o1)
	if err != nil {
		t.Errorf("Error when calling sell order")
	}

	buyOrderResponse, err := ob.buyOrder(&o2)
	if err != nil {
		t.Errorf("Error when calling buy order")
	}

	testutils.CompareEngineResponse(t, expectedBuyOrderResponse, buyOrderResponse)
	testutils.CompareEngineResponse(t, expectedSellOrderResponse, sellOrderResponse)
}

func TestFillOrder2(t *testing.T) {
	_, ob, _, _, _, _, _, _, factory1, factory2 := setupTest()

	o1, _ := factory1.NewBuyOrder(1e3, 1e8)
	o2, _ := factory2.NewSellOrder(1e3, 1e8)
	expt1 := types.NewTrade(&o1, &o2, units.Ethers(1e8), big.NewInt(1e3))

	expo1 := o1
	expo1.Status = "OPEN"
	expectedBuyOrderResponse := &types.EngineResponse{
		Status: "ORDER_ADDED",
		Order:  &expo1,
	}

	expo2 := o2
	expo2.Status = "FILLED"
	expo2.FilledAmount = utils.Ethers(1e8)

	expo3 := o1
	expo3.Status = "FILLED"
	expo3.FilledAmount = utils.Ethers(1e8)

	expectedMatches := types.NewMatches(
		[]*types.Order{&expo3},
		&expo2,
		[]*types.Trade{expt1},
	)

	expectedSellOrderResponse := &types.EngineResponse{
		Status:  "ORDER_FILLED",
		Order:   &expo2,
		Matches: expectedMatches,
	}

	res1, err := ob.buyOrder(&o1)
	if err != nil {
		t.Error("Error when sending buy order")
	}

	res2, err := ob.sellOrder(&o2)
	if err != nil {
		t.Error("Error when sending sell order")
	}

	testutils.CompareEngineResponse(t, expectedBuyOrderResponse, res1)
	testutils.CompareEngineResponse(t, expectedSellOrderResponse, res2)
}

func TestMultiMatchOrder1(t *testing.T) {
	_, ob, _, _, _, _, _, _, factory1, factory2 := setupTest()

	so1, _ := factory1.NewSellOrder(1e3+1, 1e8)
	so2, _ := factory1.NewSellOrder(1e3+2, 1e8)
	so3, _ := factory1.NewSellOrder(1e3+3, 1e8)
	bo1, _ := factory2.NewBuyOrder(1e3+4, 3e8)

	ob.sellOrder(&so1)
	ob.sellOrder(&so2)
	ob.sellOrder(&so3)

	expso1 := so1
	expso1.Status = "FILLED"
	expso1.FilledAmount = utils.Ethers(1e8)
	expso2 := so2
	expso2.Status = "FILLED"
	expso2.FilledAmount = utils.Ethers(1e8)
	expso3 := so3
	expso3.Status = "FILLED"
	expso3.FilledAmount = utils.Ethers(1e8)
	expbo1 := bo1
	expbo1.Status = "FILLED"
	expbo1.FilledAmount = utils.Ethers(3e8)

	expt1 := types.NewTrade(&so1, &bo1, utils.Ethers(1e8), big.NewInt(1e3+4))
	expt2 := types.NewTrade(&so2, &bo1, utils.Ethers(1e8), big.NewInt(1e3+4))
	expt3 := types.NewTrade(&so3, &bo1, utils.Ethers(1e8), big.NewInt(1e3+4))

	expectedMatches := types.NewMatches(
		[]*types.Order{&expso1, &expso2, &expso3},
		&bo1,
		[]*types.Trade{expt1, expt2, expt3},
	)

	expectedResponse := &types.EngineResponse{
		Status:  "ORDER_FILLED",
		Order:   &bo1,
		Matches: expectedMatches,
	}

	response, err := ob.buyOrder(&bo1)
	if err != nil {
		t.Errorf("Error in sellOrder: %s", err)
	}

	testutils.CompareEngineResponse(t, expectedResponse, response)
}

func TestMultiMatchOrder2(t *testing.T) {
	_, ob, _, _, _, _, _, _, factory1, factory2 := setupTest()

	bo1, _ := factory1.NewBuyOrder(1e3+1, 1e8)
	bo2, _ := factory1.NewBuyOrder(1e3+2, 1e8)
	bo3, _ := factory1.NewBuyOrder(1e3+3, 1e8)
	so1, _ := factory2.NewSellOrder(1e3, 3e8)

	expbo1 := bo1
	expbo1.Status = "FILLED"
	expbo1.FilledAmount = units.Ethers(1e8)
	expbo2 := bo2
	expbo2.Status = "FILLED"
	expbo2.FilledAmount = units.Ethers(1e8)
	expbo3 := bo3
	expbo3.Status = "FILLED"
	expbo3.FilledAmount = units.Ethers(1e8)
	expso1 := so1
	expso1.Status = "FILLED"
	expso1.FilledAmount = utils.Ethers(3e8)

	ob.buyOrder(&bo1)
	ob.buyOrder(&bo2)
	ob.buyOrder(&bo3)

	expt1 := types.NewTrade(&bo1, &so1, units.Ethers(1e8), big.NewInt(1e3))
	expt2 := types.NewTrade(&bo2, &so1, units.Ethers(1e8), big.NewInt(1e3))
	expt3 := types.NewTrade(&bo3, &so1, units.Ethers(1e8), big.NewInt(1e3))

	expectedMatches := types.NewMatches(
		[]*types.Order{&expbo3, &expbo2, &expbo1},
		&so1,
		[]*types.Trade{expt3, expt2, expt1},
	)

	expectedResponse := &types.EngineResponse{
		Status:  "ORDER_FILLED",
		Order:   &so1,
		Matches: expectedMatches,
	}

	res, err := ob.sellOrder(&so1)
	if err != nil {
		t.Errorf("Error in sell order: %s", err)
	}

	testutils.CompareMatches(t, expectedResponse.Matches, res.Matches)
}

func TestPartialMatchOrder1(t *testing.T) {
	_, ob, _, _, _, _, _, _, factory1, factory2 := setupTest()

	so1, _ := factory1.NewSellOrder(1e3+1, 1e8)
	so2, _ := factory1.NewSellOrder(1e3+2, 1e8)
	so3, _ := factory1.NewSellOrder(1e3+3, 1e8)
	so4, _ := factory1.NewSellOrder(1e3+4, 2e8)
	bo1, _ := factory2.NewBuyOrder(1e3+5, 4e8)

	expso1 := so1
	expso1.FilledAmount = units.Ethers(1e8)
	expso1.Status = "FILLED"
	expso2 := so2
	expso2.FilledAmount = units.Ethers(1e8)
	expso2.Status = "FILLED"
	expso3 := so3
	expso3.FilledAmount = units.Ethers(1e8)
	expso3.Status = "FILLED"
	expso4 := so4
	expso4.FilledAmount = units.Ethers(1e8)
	expso4.Status = "PARTIAL_FILLED"
	expbo1 := bo1
	expbo1.FilledAmount = units.Ethers(4e8)
	expbo1.Status = "FILLED"

	expt1 := types.NewTrade(&so1, &bo1, units.Ethers(1e8), big.NewInt(1e3+5))
	expt2 := types.NewTrade(&so2, &bo1, units.Ethers(1e8), big.NewInt(1e3+5))
	expt3 := types.NewTrade(&so3, &bo1, units.Ethers(1e8), big.NewInt(1e3+5))
	expt4 := types.NewTrade(&so4, &bo1, units.Ethers(1e8), big.NewInt(1e3+5))

	ob.sellOrder(&so1)
	ob.sellOrder(&so2)
	ob.sellOrder(&so3)
	ob.sellOrder(&so4)

	res, err := ob.buyOrder(&bo1)
	if err != nil {
		t.Errorf("Error when buying order")
	}

	expectedMatches := types.NewMatches(
		[]*types.Order{&expso1, &expso2, &expso3, &expso4},
		&bo1,
		[]*types.Trade{expt1, expt2, expt3, expt4},
	)

	expectedResponse := &types.EngineResponse{
		Status:  "ORDER_FILLED",
		Order:   &expbo1,
		Matches: expectedMatches,
	}

	testutils.CompareEngineResponse(t, expectedResponse, res)
}

func TestPartialMatchOrder2(t *testing.T) {
	_, ob, _, _, _, _, _, _, factory1, factory2 := setupTest()

	bo1, _ := factory1.NewBuyOrder(1e3+5, 1e8)
	bo2, _ := factory1.NewBuyOrder(1e3+4, 1e8)
	bo3, _ := factory1.NewBuyOrder(1e3+3, 1e8)
	bo4, _ := factory1.NewBuyOrder(1e3+2, 2e8)
	so1, _ := factory2.NewSellOrder(1e3+1, 4e8)

	expbo1 := bo1
	expbo1.FilledAmount = utils.Ethers(1e8)
	expbo1.Status = "FILLED"
	expbo2 := bo2
	expbo2.FilledAmount = utils.Ethers(1e8)
	expbo2.Status = "FILLED"
	expbo3 := bo3
	expbo3.FilledAmount = utils.Ethers(1e8)
	expbo3.Status = "FILLED"
	expbo4 := bo4
	expbo4.FilledAmount = utils.Ethers(1e8)
	expbo4.Status = "PARTIAL_FILLED"

	expso1 := so1
	expso1.FilledAmount = utils.Ethers(4e8)
	expso1.Status = "FILLED"

	expt1 := types.NewTrade(&bo1, &so1, utils.Ethers(1e8), big.NewInt(1e3+1))
	expt2 := types.NewTrade(&bo2, &so1, utils.Ethers(1e8), big.NewInt(1e3+1))
	expt3 := types.NewTrade(&bo3, &so1, utils.Ethers(1e8), big.NewInt(1e3+1))
	expt4 := types.NewTrade(&bo4, &so1, utils.Ethers(1e8), big.NewInt(1e3+1))

	ob.buyOrder(&bo1)
	ob.buyOrder(&bo2)
	ob.buyOrder(&bo3)
	ob.buyOrder(&bo4)

	res, err := ob.sellOrder(&so1)
	if err != nil {
		t.Errorf("Error when buying order")
	}

	expectedMatches := types.NewMatches(
		[]*types.Order{&expbo1, &expbo2, &expbo3, &expbo4},
		&so1,
		[]*types.Trade{expt1, expt2, expt3, expt4},
	)

	expectedResponse := &types.EngineResponse{
		Status:  "ORDER_FILLED",
		Order:   &expso1,
		Matches: expectedMatches,
	}

	testutils.CompareEngineResponse(t, expectedResponse, res)
}
