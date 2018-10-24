package engine

import (
	"io/ioutil"
	"log"
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

	pair := testutils.GetZRXWETHTestPair()
	pairDao := new(mocks.PairDao)
	pairDao.On("GetAll").Return([]types.Pair{*pair}, nil)

	eng := NewEngine(rabbitConn, orderDao, pairDao)
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
		Status:         "NOMATCH",
		Order:          &exp1,
		RemainingOrder: nil,
		Matches:        nil,
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
		Status:         "NOMATCH",
		Order:          &exp1,
		RemainingOrder: nil,
		Matches:        nil,
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
	expt1, _ := types.NewUnsignedTrade1(&o1, &o2, units.Ethers(1e8))

	expo1 := o1
	expo1.Status = "OPEN"
	expectedSellOrderResponse := &types.EngineResponse{
		Status: "NOMATCH",
		Order:  &expo1,
	}

	expo2 := o2
	expo2.Status = "FILLED"
	expo2.FilledAmount = units.Ethers(1e8)

	expo3 := o1
	expo3.Status = "FILLED"
	expo3.FilledAmount = units.Ethers(1e8)

	expectedMatches := types.NewMatches([]*types.Order{&expo3}, []*types.Trade{&expt1})

	expectedBuyOrderResponse := &types.EngineResponse{
		Status:  "FULL",
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
	expt1, _ := types.NewUnsignedTrade1(&o1, &o2, utils.Ethers(1e8))

	expo1 := o1
	expo1.Status = "OPEN"
	expectedBuyOrderResponse := &types.EngineResponse{
		Status: "NOMATCH",
		Order:  &expo1,
	}

	expo2 := o2
	expo2.Status = "FILLED"
	expo2.FilledAmount = utils.Ethers(1e8)

	expo3 := o1
	expo3.Status = "FILLED"
	expo3.FilledAmount = utils.Ethers(1e8)

	expectedMatches := types.NewMatches([]*types.Order{&expo3}, []*types.Trade{&expt1})

	expectedSellOrderResponse := &types.EngineResponse{
		Status:  "FULL",
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

	expt1, _ := types.NewUnsignedTrade1(&so1, &bo1, utils.Ethers(1e8))
	expt2, _ := types.NewUnsignedTrade1(&so2, &bo1, utils.Ethers(1e8))
	expt3, _ := types.NewUnsignedTrade1(&so3, &bo1, utils.Ethers(1e8))

	expectedMatches := types.NewMatches([]*types.Order{&expso1, &expso2, &expso3}, []*types.Trade{&expt1, &expt2, &expt3})

	expectedResponse := &types.EngineResponse{
		"FULL",
		common.Hash{},
		&expbo1,
		nil,
		expectedMatches,
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

	expt1, _ := types.NewUnsignedTrade1(&bo1, &so1, units.Ethers(1e8))
	expt2, _ := types.NewUnsignedTrade1(&bo2, &so1, units.Ethers(1e8))
	expt3, _ := types.NewUnsignedTrade1(&bo3, &so1, units.Ethers(1e8))
	expectedMatches := types.NewMatches([]*types.Order{&expbo3, &expbo2, &expbo1}, []*types.Trade{&expt3, &expt2, &expt1})

	expectedResponse := &types.EngineResponse{
		"FULL",
		common.Hash{},
		&expso1,
		nil,
		expectedMatches,
	}

	res, err := ob.sellOrder(&so1)
	if err != nil {
		t.Errorf("Error in sell order: %s", err)
	}

	testutils.CompareMatches(t, expectedResponse.Matches, res.Matches)
}

func TestPartialMatchOrder1(t *testing.T) {
	_, ob, _, _, _, _, _, _, factory1, factory2 := setupTest()
	// defer e.redisConn.FlushAll()

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

	expt1, _ := types.NewUnsignedTrade1(&so1, &bo1, units.Ethers(1e8))
	expt2, _ := types.NewUnsignedTrade1(&so2, &bo1, units.Ethers(1e8))
	expt3, _ := types.NewUnsignedTrade1(&so3, &bo1, units.Ethers(1e8))
	expt4, _ := types.NewUnsignedTrade1(&so4, &bo1, units.Ethers(1e8))

	ob.sellOrder(&so1)
	ob.sellOrder(&so2)
	ob.sellOrder(&so3)
	ob.sellOrder(&so4)

	res, err := ob.buyOrder(&bo1)
	if err != nil {
		t.Errorf("Error when buying order")
	}

	expectedMatches := types.NewMatches([]*types.Order{&expso1, &expso2, &expso3, &expso4}, []*types.Trade{&expt1, &expt2, &expt3, &expt4})

	expectedResponse := &types.EngineResponse{
		"FULL",
		common.Hash{},
		&expbo1,
		nil,
		expectedMatches,
	}

	testutils.CompareEngineResponse(t, expectedResponse, res)
}

func TestPartialMatchOrder2(t *testing.T) {
	_, ob, _, _, _, _, _, _, factory1, factory2 := setupTest()
	// defer e.redisConn.FlushAll()

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

	expt1, _ := types.NewUnsignedTrade1(&bo1, &so1, utils.Ethers(1e8))
	expt2, _ := types.NewUnsignedTrade1(&bo2, &so1, utils.Ethers(1e8))
	expt3, _ := types.NewUnsignedTrade1(&bo3, &so1, utils.Ethers(1e8))
	expt4, _ := types.NewUnsignedTrade1(&bo4, &so1, utils.Ethers(1e8))

	ob.buyOrder(&bo1)
	ob.buyOrder(&bo2)
	ob.buyOrder(&bo3)
	ob.buyOrder(&bo4)

	res, err := ob.sellOrder(&so1)
	if err != nil {
		t.Errorf("Error when buying order")
	}

	expectedMatches := types.NewMatches([]*types.Order{&expbo1, &expbo2, &expbo3, &expbo4}, []*types.Trade{&expt1, &expt2, &expt3, &expt4})

	expectedResponse := &types.EngineResponse{
		"FULL",
		common.Hash{},
		&expso1,
		nil,
		expectedMatches,
	}

	testutils.CompareEngineResponse(t, expectedResponse, res)
}

// func TestRecoverOrders(t *testing.T) {
// 	e, _, _, _, _, _, _, factory1, factory2 := setupTest()
// 	defer e.redisConn.FlushAll()

// 	o1, _ := factory1.NewSellOrder(1e3, 1e8, 5e7)
// 	o2, _ := factory1.NewSellOrder(1e3, 1e8, 1e8)
// 	o3, _ := factory1.NewSellOrder(1e3, 1e8, 1e8)

// 	t1, _ := factory2.NewTrade(o1, 5e7)
// 	t2, _ := factory2.NewTrade(o2, 5e7)
// 	t3, _ := factory2.NewTrade(o3, 1e8)

// 	e.addOrder(&o1)
// 	e.addOrder(&o2)
// 	e.addOrder(&o3)

// 	pricePointSetKey, orderHashListKey := o1.GetOBKeys()
// 	pricepoints, _ := e.redisConn.GetSortedSet(pricePointSetKey)
// 	pricePointHashes, _ := e.redisConn.GetSortedSet(orderHashListKey)
// 	volume, _ := e.GetPricePointVolume(pricePointSetKey, o2.PricePoint.Int64())
// 	stored1, _ := e.GetFromOrderMap(o1.Hash)
// 	stored2, _ := e.GetFromOrderMap(o2.Hash)
// 	stored3, _ := e.GetFromOrderMap(o3.Hash)

// 	assert.Equal(t, 1, len(pricepoints))
// 	assert.Contains(t, pricepoints, utils.UintToPaddedString(o1.PricePoint.Int64()))
// 	assert.Equal(t, 3, len(pricePointHashes))
// 	assert.Contains(t, pricePointHashes, o1.Hash.Hex())
// 	assert.Contains(t, pricePointHashes, o2.Hash.Hex())
// 	assert.Contains(t, pricePointHashes, o3.Hash.Hex())
// 	assert.Equal(t, int64(pricePointHashes[o1.Hash.Hex()]), o1.CreatedAt.Unix())
// 	assert.Equal(t, int64(pricePointHashes[o2.Hash.Hex()]), o2.CreatedAt.Unix())
// 	assert.Equal(t, int64(pricePointHashes[o3.Hash.Hex()]), o3.CreatedAt.Unix())
// 	assert.Equal(t, "50000000", volume)

// 	testutils.Compare(t, &o1, stored1)
// 	testutils.Compare(t, &o2, stored2)
// 	testutils.Compare(t, &o3, stored3)

// 	matches := []*types.OrderTradePair{
// 		{&o1, &t1},
// 		{&o2, &t2},
// 		{&o3,
// 	}

// 	recoverOrders := []*types.FillOrder{
// 		&types.FillOrder{big.NewInt(5e7), &o1},
// 		&types.FillOrder{big.NewInt(5e7), &o2},
// 		&types.FillOrder{big.NewInt(1e8), &o3},
// 	}

// 	err := e.RecoverOrders(recoverOrders)
// 	if err != nil {
// 		t.Error("Error when recovering orders", err)
// 	}

// 	pricePointSetKey, orderHashListKey = o1.GetOBKeys()

// 	pricepoints, _ = e.redisConn.GetSortedSet(pricePointSetKey)
// 	pricePointHashes, _ = e.redisConn.GetSortedSet(orderHashListKey)
// 	volume, _ = e.GetPricePointVolume(pricePointSetKey, o2.PricePoint.Int64())
// 	stored1, _ = e.GetFromOrderMap(o1.Hash)
// 	stored2, _ = e.GetFromOrderMap(o2.Hash)
// 	stored3, _ = e.GetFromOrderMap(o3.Hash)

// 	assert.Equal(t, 1, len(pricepoints))
// 	assert.Contains(t, pricepoints, utils.UintToPaddedString(o1.PricePoint.Int64()))
// 	assert.Equal(t, 3, len(pricePointHashes))
// 	assert.Contains(t, pricePointHashes, o1.Hash.Hex())
// 	assert.Contains(t, pricePointHashes, o2.Hash.Hex())
// 	assert.Contains(t, pricePointHashes, o3.Hash.Hex())
// 	assert.Equal(t, int64(pricePointHashes[o1.Hash.Hex()]), o1.CreatedAt.Unix())
// 	assert.Equal(t, int64(pricePointHashes[o2.Hash.Hex()]), o2.CreatedAt.Unix())
// 	assert.Equal(t, int64(pricePointHashes[o3.Hash.Hex()]), o3.CreatedAt.Unix())
// 	assert.Equal(t, "250000000", volume)

// 	exp1 := o1
// 	exp1.FilledAmount = big.NewInt(0)
// 	exp2 := o2
// 	exp2.FilledAmount = big.NewInt(5e7)
// 	ex3 := o3
// 	ex3.FilledAmount = big.NewInt(0)

// 	testutils.Compare(t, &exp1, stored1)
// 	testutils.Compare(t, &exp2, stored2)
// 	testutils.Compare(t, &ex3, stored3)
// }

// import (
// 	"encoding/json"
// 	"math/big"
// 	"testing"

// 	"github.com/Proofsuite/amp-matching-engine/types"
// 	"github.com/Proofsuite/amp-matching-engine/utils/math"
// 	"github.com/stretchr/testify/assert"
// )

// func TestExecute(t *testing.T) {
// 	e := getResource()
// 	defer e.redisConn.FlushAll()
// 	// Test Case1: bookEntry amount is less than order amount
// 	// New Buy Order
// 	bookEntry := getBuyOrder()
// 	bookEntry.FilledAmount = big.NewInt(1000000000)

// 	e.addOrder(&bookEntry)

// 	order := getSellOrder()

// 	expectedAmount := math.Sub(bookEntry.Amount, bookEntry.FilledAmount)

// 	expectedTrade := getTrade(&order, &bookEntry, expectedAmount, big.NewInt(0))

// 	expectedTrade.Hash = expectedTrade.ComputeHash()

// 	etb, _ := json.Marshal(expectedTrade)
// 	expectedBookEntry := bookEntry
// 	expectedBookEntry.Status = "FILLED"
// 	expectedBookEntry.FilledAmount = bookEntry.Amount

// 	expectedFillOrder := &types.FillOrder{
// 		Amount: math.Sub(bookEntry.Amount, bookEntry.FilledAmount),
// 		Order:  &expectedBookEntry,
// 	}
// 	efob, _ := json.Marshal(expectedFillOrder)

// 	trade, fillOrder, err := e.execute(&order, &bookEntry)
// 	if err != nil {
// 		t.Errorf("Error in execute: %s", err)
// 		return
// 	}
// 	tb, _ := json.Marshal(trade)
// 	fob, _ := json.Marshal(fillOrder)
// 	assert.JSONEq(t, string(etb), string(tb))
// 	assert.JSONEq(t, string(efob), string(fob))

// 	// Test Case2: bookEntry amount is equal to order amount
// 	// unmarshal bookentry and order from json string
// 	bookEntry = getBuyOrder()
// 	order = getSellOrder()
// 	bookEntry.FilledAmount = big.NewInt(0)
// 	expectedAmount = math.Sub(bookEntry.Amount, bookEntry.FilledAmount)
// 	expectedTrade = getTrade(&order, &bookEntry, expectedAmount, big.NewInt(0))
// 	expectedTrade.Hash = expectedTrade.ComputeHash()

// 	etb, _ = json.Marshal(expectedTrade)
// 	expectedBookEntry = bookEntry
// 	expectedBookEntry.Status = "FILLED"
// 	expectedBookEntry.FilledAmount = bookEntry.Amount

// 	expectedFillOrder = &types.FillOrder{
// 		Amount: bookEntry.Amount,
// 		Order:  &expectedBookEntry,
// 	}

// 	efob, _ = json.Marshal(expectedFillOrder)

// 	e.addOrder(&bookEntry)

// 	trade, fillOrder, err = e.execute(&order, &bookEntry)
// 	if err != nil {
// 		t.Errorf("Error in execute: %s", err)
// 		return
// 	} else {
// 		tb, _ := json.Marshal(trade)
// 		fob, _ := json.Marshal(fillOrder)
// 		assert.JSONEq(t, string(etb), string(tb))
// 		assert.JSONEq(t, string(efob), string(fob))
// 	}

// 	// Test Case3: bookEntry amount is greater then order amount
// 	// unmarshal bookentry and order from json string
// 	bookEntry = getBuyOrder()
// 	order = getSellOrder()
// 	bookEntry.Amount = math.Add(bookEntry.Amount, big.NewInt(1000000000))
// 	expectedAmount = order.Amount
// 	expectedTrade = getTrade(&order, &bookEntry, expectedAmount, big.NewInt(0))
// 	expectedTrade.Hash = expectedTrade.ComputeHash()

// 	etb, _ = json.Marshal(expectedTrade)
// 	expectedBookEntry = bookEntry
// 	expectedBookEntry.Status = "PARTIAL_FILLED"
// 	expectedBookEntry.FilledAmount = math.Add(expectedBookEntry.FilledAmount, order.Amount)

// 	expectedFillOrder = &types.FillOrder{
// 		Amount: order.Amount,
// 		Order:  &expectedBookEntry,
// 	}

// 	efob, _ = json.Marshal(expectedFillOrder)
// 	e.addOrder(&bookEntry)

// 	trade, fillOrder, err = e.execute(&order, &bookEntry)
// 	if err != nil {
// 		t.Errorf("Error in execute: %s", err)
// 		return
// 	}
// 	tb, _ = json.Marshal(trade)
// 	fob, _ = json.Marshal(fillOrder)
// 	assert.JSONEq(t, string(etb), string(tb))
// 	assert.JSONEq(t, string(efob), string(fob))
// }
