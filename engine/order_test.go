package engine

import (
	"log"
	"math/big"
	"testing"

	"github.com/Proofsuite/amp-matching-engine/types"
	"github.com/Proofsuite/amp-matching-engine/utils"
	"github.com/Proofsuite/amp-matching-engine/utils/testutils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
)

func setupTest() (
	*Engine,
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

	e := getResource()

	ex := testutils.GetTestAddress1()
	maker := testutils.GetTestWallet1()
	taker := testutils.GetTestWallet2()
	pair := testutils.GetZRXWETHTestPair()
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

	return e, ex, maker, taker, pair, zrx, weth, factory1, factory2
}

func TestAddOrder(t *testing.T) {
	e, _, _, _, _, _, _, factory1, factory2 := setupTest()
	defer e.redisConn.FlushAll()

	o1, _ := factory1.NewSellOrder(1e3, 1e8)
	o2, _ := factory2.NewSellOrder(1e3, 1e8)

	e.addOrder(&o1)

	pricePointSetKey, orderHashListKey := o1.GetOBKeys()

	pricepoints, err := e.redisConn.GetSortedSet(pricePointSetKey)
	if err != nil {
		t.Error("Error getting sorted set")
	}

	pricePointHashes, err := e.redisConn.GetSortedSet(orderHashListKey)
	if err != nil {
		t.Error("Error getting pricepoints order hashes set")
	}

	stored, err := e.GetFromOrderMap(o1.Hash)
	if err != nil {
		t.Error("Error getting sorted set", err)
	}

	volume, err := e.GetPricePointVolume(pricePointSetKey, o1.PricePoint.Int64())
	if err != nil {
		t.Error("Error getting volume set", err)
	}

	assert.Equal(t, 1, len(pricepoints))
	assert.Contains(t, pricepoints, utils.UintToPaddedString(o1.PricePoint.Int64()))
	assert.Equal(t, 1, len(pricePointHashes))
	assert.Contains(t, pricePointHashes, o1.Hash.Hex())
	assert.Equal(t, int64(pricePointHashes[o1.Hash.Hex()]), o1.CreatedAt.Unix())
	assert.Equal(t, "100000000", volume)
	testutils.CompareOrder(t, &o1, stored)

	e.addOrder(&o2)

	pricePointSetKey, orderHashListKey = o2.GetOBKeys()
	pricepoints, err = e.redisConn.GetSortedSet(pricePointSetKey)
	if err != nil {
		t.Error(err)
	}

	pricePointHashes, err = e.redisConn.GetSortedSet(orderHashListKey)
	if err != nil {
		t.Error("Error getting pricepoints order hashes set")
	}

	stored, err = e.GetFromOrderMap(o2.Hash)
	if err != nil {
		t.Error("Error getting order from map", err)
	}

	volume, err = e.GetPricePointVolume(pricePointSetKey, o2.PricePoint.Int64())
	if err != nil {
		t.Error("Error getting volume set", err)
	}

	assert.Equal(t, 1, len(pricepoints))
	assert.Contains(t, pricepoints, utils.UintToPaddedString(o2.PricePoint.Int64()))
	assert.Equal(t, 2, len(pricePointHashes))
	assert.Contains(t, pricePointHashes, o2.Hash.Hex())
	assert.Equal(t, int64(pricePointHashes[o2.Hash.Hex()]), o2.CreatedAt.Unix())
	assert.Equal(t, "200000000", volume)
	testutils.CompareOrder(t, &o2, stored)
}

func TestUpdateOrder(t *testing.T) {
	e, _, _, _, _, _, _, factory1, _ := setupTest()
	defer e.redisConn.FlushAll()

	o1, _ := factory1.NewSellOrder(1e3, 1e8)

	exp1 := o1
	exp1.Status = "PARTIAL_FILLED"
	exp1.FilledAmount = big.NewInt(1000)

	err := e.addOrder(&o1)
	if err != nil {
		t.Error("Could not add order")
	}

	err = e.updateOrder(&o1, big.NewInt(1e3))
	if err != nil {
		t.Error("Could not update order")
	}

	pricePointSetKey, orderHashListKey := o1.GetOBKeys()

	pricepoints, err := e.redisConn.GetSortedSet(pricePointSetKey)
	if err != nil {
		t.Error("Error getting pricepoint set", err)
	}

	pricePointHashes, err := e.redisConn.GetSortedSet(orderHashListKey)
	if err != nil {
		t.Error("Error getting pricepoint hash set", err)
	}

	volume, err := e.GetPricePointVolume(pricePointSetKey, o1.PricePoint.Int64())
	if err != nil {
		t.Error("Error getting pricepoint volume", err)
	}

	stored, err := e.GetFromOrderMap(o1.Hash)
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, 1, len(pricepoints))
	assert.Contains(t, pricepoints, utils.UintToPaddedString(o1.PricePoint.Int64()))
	testutils.Compare(t, &exp1, stored)

	assert.Equal(t, 1, len(pricePointHashes))
	assert.Contains(t, pricePointHashes, o1.Hash.Hex())
	assert.Equal(t, int64(pricePointHashes[o1.Hash.Hex()]), o1.CreatedAt.Unix())
	assert.Equal(t, "99999000", volume)
}

func TestDeleteOrder(t *testing.T) {
	e, _, _, _, _, _, _, factory1, _ := setupTest()
	defer e.redisConn.FlushAll()

	o1, _ := factory1.NewSellOrder(1e3, 1e8)

	e.addOrder(&o1)

	pricePointSetKey, orderHashListKey := o1.GetOBKeys()

	pricepoints, err := e.redisConn.GetSortedSet(pricePointSetKey)
	if err != nil {
		t.Error(err)
	}

	pricePointHashes, err := e.redisConn.GetSortedSet(orderHashListKey)
	if err != nil {
		t.Error(err)
	}

	volume, err := e.GetPricePointVolume(pricePointSetKey, o1.PricePoint.Int64())
	if err != nil {
		t.Error(err)
	}

	stored, err := e.GetFromOrderMap(o1.Hash)
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, 1, len(pricepoints))
	assert.Contains(t, pricepoints, utils.UintToPaddedString(o1.PricePoint.Int64()))
	testutils.CompareOrder(t, &o1, stored)

	assert.Equal(t, 1, len(pricePointHashes))
	assert.Contains(t, pricePointHashes, o1.Hash.Hex())
	assert.Equal(t, int64(pricePointHashes[o1.Hash.Hex()]), o1.CreatedAt.Unix())
	assert.Equal(t, "100000000", volume)

	e.deleteOrder(&o1, big.NewInt(1e8))

	pricePointSetKey, orderHashListKey = o1.GetOBKeys()

	if e.redisConn.Exists(pricePointSetKey) {
		t.Errorf("Key: %v expected to be deleted but exists", pricePointSetKey)
	}

	if e.redisConn.Exists(orderHashListKey) {
		t.Errorf("Key: %v expected to be deleted but key exists", pricePointSetKey)
	}

	if e.redisConn.Exists(orderHashListKey + "::" + o1.Hash.Hex()) {
		t.Errorf("Key: %v expected to be deleted but key exists", pricePointSetKey)
	}

	if e.redisConn.Exists(pricePointSetKey + "::book::" + utils.UintToPaddedString(o1.PricePoint.Int64())) {
		t.Errorf("Key: %v expected to be deleted but key exists", pricePointSetKey)
	}
}

func TestCancelOrder(t *testing.T) {
	e, _, _, _, _, _, _, factory1, _ := setupTest()
	defer e.redisConn.FlushAll()

	o1, _ := factory1.NewSellOrder(1e3, 1e8)
	o2, _ := factory1.NewSellOrder(1e3, 1e8)

	e.addOrder(&o1)
	e.addOrder(&o2)

	expectedOrder := o2
	expectedOrder.Status = "CANCELLED"
	expected := &types.EngineResponse{
		FillStatus:     "CANCELLED",
		Order:          &expectedOrder,
		RemainingOrder: nil,
		MatchingOrders: nil,
		Trades:         nil,
	}

	pricePointSetKey, orderHashListKey := o1.GetOBKeys()

	pricepoints, _ := e.redisConn.GetSortedSet(pricePointSetKey)
	pricePointHashes, _ := e.redisConn.GetSortedSet(orderHashListKey)
	volume, _ := e.GetPricePointVolume(pricePointSetKey, o2.PricePoint.Int64())
	stored1, _ := e.GetFromOrderMap(o1.Hash)
	stored2, _ := e.GetFromOrderMap(o2.Hash)

	assert.Equal(t, 1, len(pricepoints))
	assert.Contains(t, pricepoints, utils.UintToPaddedString(o1.PricePoint.Int64()))
	assert.Equal(t, 2, len(pricePointHashes))
	assert.Contains(t, pricePointHashes, o1.Hash.Hex())
	assert.Contains(t, pricePointHashes, o2.Hash.Hex())
	assert.Equal(t, int64(pricePointHashes[o1.Hash.Hex()]), o1.CreatedAt.Unix())
	assert.Equal(t, int64(pricePointHashes[o2.Hash.Hex()]), o2.CreatedAt.Unix())
	assert.Equal(t, "200000000", volume)
	testutils.Compare(t, &o1, stored1)
	testutils.Compare(t, &o2, stored2)

	res, err := e.CancelOrder(&o2)
	if err != nil {
		t.Error("Error when cancelling order: ", err)
	}

	assert.Equal(t, expected, res)

	pricePointSetKey, orderHashListKey = o1.GetOBKeys()

	pricepoints, _ = e.redisConn.GetSortedSet(pricePointSetKey)
	pricePointHashes, _ = e.redisConn.GetSortedSet(orderHashListKey)
	volume, _ = e.GetPricePointVolume(pricePointSetKey, o2.PricePoint.Int64())
	stored1, _ = e.GetFromOrderMap(o1.Hash)
	stored2, _ = e.GetFromOrderMap(o2.Hash)

	assert.Equal(t, 1, len(pricepoints))
	assert.Contains(t, pricepoints, utils.UintToPaddedString(o1.PricePoint.Int64()))
	assert.Equal(t, 1, len(pricePointHashes))
	assert.Contains(t, pricePointHashes, o1.Hash.Hex())
	assert.NotContains(t, pricePointHashes, o2.Hash.Hex())
	assert.Equal(t, int64(pricePointHashes[o1.Hash.Hex()]), o1.CreatedAt.Unix())
	assert.Equal(t, "100000000", volume)
	testutils.Compare(t, &o1, stored1)
	testutils.Compare(t, nil, stored2)
}

func TestRecoverOrders(t *testing.T) {
	e, _, _, _, _, _, _, factory1, _ := setupTest()
	defer e.redisConn.FlushAll()

	o1, _ := factory1.NewSellOrder(1e3, 1e8, 5e7)
	o2, _ := factory1.NewSellOrder(1e3, 1e8, 1e8)
	o3, _ := factory1.NewSellOrder(1e3, 1e8, 1e8)

	e.addOrder(&o1)
	e.addOrder(&o2)
	e.addOrder(&o3)

	pricePointSetKey, orderHashListKey := o1.GetOBKeys()
	pricepoints, _ := e.redisConn.GetSortedSet(pricePointSetKey)
	pricePointHashes, _ := e.redisConn.GetSortedSet(orderHashListKey)
	volume, _ := e.GetPricePointVolume(pricePointSetKey, o2.PricePoint.Int64())
	stored1, _ := e.GetFromOrderMap(o1.Hash)
	stored2, _ := e.GetFromOrderMap(o2.Hash)
	stored3, _ := e.GetFromOrderMap(o3.Hash)

	assert.Equal(t, 1, len(pricepoints))
	assert.Contains(t, pricepoints, utils.UintToPaddedString(o1.PricePoint.Int64()))
	assert.Equal(t, 3, len(pricePointHashes))
	assert.Contains(t, pricePointHashes, o1.Hash.Hex())
	assert.Contains(t, pricePointHashes, o2.Hash.Hex())
	assert.Contains(t, pricePointHashes, o3.Hash.Hex())
	assert.Equal(t, int64(pricePointHashes[o1.Hash.Hex()]), o1.CreatedAt.Unix())
	assert.Equal(t, int64(pricePointHashes[o2.Hash.Hex()]), o2.CreatedAt.Unix())
	assert.Equal(t, int64(pricePointHashes[o3.Hash.Hex()]), o3.CreatedAt.Unix())
	assert.Equal(t, "50000000", volume)

	testutils.Compare(t, &o1, stored1)
	testutils.Compare(t, &o2, stored2)
	testutils.Compare(t, &o3, stored3)

	recoverOrders := []*types.FillOrder{
		&types.FillOrder{big.NewInt(5e7), &o1},
		&types.FillOrder{big.NewInt(5e7), &o2},
		&types.FillOrder{big.NewInt(1e8), &o3},
	}

	err := e.RecoverOrders(recoverOrders)
	if err != nil {
		t.Error("Error when recovering orders", err)
	}

	pricePointSetKey, orderHashListKey = o1.GetOBKeys()

	pricepoints, _ = e.redisConn.GetSortedSet(pricePointSetKey)
	pricePointHashes, _ = e.redisConn.GetSortedSet(orderHashListKey)
	volume, _ = e.GetPricePointVolume(pricePointSetKey, o2.PricePoint.Int64())
	stored1, _ = e.GetFromOrderMap(o1.Hash)
	stored2, _ = e.GetFromOrderMap(o2.Hash)
	stored3, _ = e.GetFromOrderMap(o3.Hash)

	assert.Equal(t, 1, len(pricepoints))
	assert.Contains(t, pricepoints, utils.UintToPaddedString(o1.PricePoint.Int64()))
	assert.Equal(t, 3, len(pricePointHashes))
	assert.Contains(t, pricePointHashes, o1.Hash.Hex())
	assert.Contains(t, pricePointHashes, o2.Hash.Hex())
	assert.Contains(t, pricePointHashes, o3.Hash.Hex())
	assert.Equal(t, int64(pricePointHashes[o1.Hash.Hex()]), o1.CreatedAt.Unix())
	assert.Equal(t, int64(pricePointHashes[o2.Hash.Hex()]), o2.CreatedAt.Unix())
	assert.Equal(t, int64(pricePointHashes[o3.Hash.Hex()]), o3.CreatedAt.Unix())
	assert.Equal(t, "250000000", volume)

	exp1 := o1
	exp1.FilledAmount = big.NewInt(0)
	exp2 := o2
	exp2.FilledAmount = big.NewInt(5e7)
	ex3 := o3
	ex3.FilledAmount = big.NewInt(0)

	testutils.Compare(t, &exp1, stored1)
	testutils.Compare(t, &exp2, stored2)
	testutils.Compare(t, &ex3, stored3)
}

func TestSellOrder(t *testing.T) {
	e, _, _, _, _, _, _, factory1, _ := setupTest()
	defer e.redisConn.FlushAll()

	o1, _ := factory1.NewSellOrder(1e3, 1e8, 0)

	exp1 := o1
	exp1.Status = "OPEN"
	expected := &types.EngineResponse{
		Order:          &exp1,
		Trades:         nil,
		FillStatus:     "NOMATCH",
		RemainingOrder: nil,
		MatchingOrders: nil,
	}

	res, err := e.sellOrder(&o1)
	if err != nil {
		t.Error("Error in sell order: ", err)
	}

	assert.Equal(t, expected, res)
}

func TestBuyOrder(t *testing.T) {
	e, _, _, _, _, _, _, factory1, _ := setupTest()
	defer e.redisConn.FlushAll()

	o1, _ := factory1.NewBuyOrder(1e3, 1e8, 0)

	exp1 := o1
	exp1.Status = "OPEN"
	expected := &types.EngineResponse{
		Order:          &exp1,
		Trades:         nil,
		FillStatus:     "NOMATCH",
		RemainingOrder: nil,
		MatchingOrders: nil,
	}

	res, err := e.buyOrder(&o1)
	if err != nil {
		t.Error("Error in buy order: ", err)
	}

	assert.Equal(t, expected, res)
}

func TestFillOrder1(t *testing.T) {
	e, _, _, _, _, _, _, factory1, factory2 := setupTest()
	defer e.redisConn.FlushAll()

	o1, _ := factory1.NewSellOrder(1e3, 1e8)
	o2, _ := factory2.NewBuyOrder(1e3, 1e8)
	expectedTrade, _ := types.NewUnsignedTrade(&o1, factory2.Wallet.Address, big.NewInt(1e8))

	exp1 := o1
	exp1.Status = "OPEN"
	expectedSellOrderResponse := &types.EngineResponse{
		FillStatus: "NOMATCH",
		Order:      &exp1,
	}

	exp2 := o2
	exp2.Status = "FILLED"
	exp2.FilledAmount = big.NewInt(1e8)

	ex3 := o1
	ex3.Status = "FILLED"
	ex3.FilledAmount = big.NewInt(1e8)

	expectedBuyOrderResponse := &types.EngineResponse{
		FillStatus:     "FULL",
		Order:          &exp2,
		MatchingOrders: []*types.FillOrder{{big.NewInt(1e8), &ex3}},
		Trades:         []*types.Trade{&expectedTrade},
	}

	sellOrderResponse, err := e.sellOrder(&o1)
	if err != nil {
		t.Errorf("Error when calling sell order")
	}
	buyOrderResponse, err := e.buyOrder(&o2)
	if err != nil {
		t.Errorf("Error when calling buy order")
	}

	assert.Equal(t, expectedBuyOrderResponse, buyOrderResponse)
	assert.Equal(t, expectedSellOrderResponse, sellOrderResponse)
}

func TestFillOrder2(t *testing.T) {
	e, _, _, _, _, _, _, factory1, factory2 := setupTest()
	defer e.redisConn.FlushAll()

	o1, _ := factory1.NewBuyOrder(1e3, 1e8)
	o2, _ := factory2.NewSellOrder(1e3, 1e8)
	expectedTrade, _ := types.NewUnsignedTrade(&o1, factory2.Wallet.Address, big.NewInt(1e8))

	exp1 := o1
	exp1.Status = "OPEN"
	expectedBuyOrderResponse := &types.EngineResponse{
		FillStatus: "NOMATCH",
		Order:      &exp1,
	}

	exp2 := o2
	exp2.Status = "FILLED"
	exp2.FilledAmount = big.NewInt(1e8)

	ex3 := o1
	ex3.Status = "FILLED"
	ex3.FilledAmount = big.NewInt(1e8)
	expectedSellOrderResponse := &types.EngineResponse{
		FillStatus:     "FULL",
		Order:          &exp2,
		Trades:         []*types.Trade{&expectedTrade},
		MatchingOrders: []*types.FillOrder{{big.NewInt(1e8), &ex3}},
	}

	res1, err := e.buyOrder(&o1)
	if err != nil {
		t.Error("Error when sending buy order")
	}

	res2, err := e.sellOrder(&o2)
	if err != nil {
		t.Error("Error when sending sell order")
	}

	testutils.Compare(t, expectedBuyOrderResponse, res1)
	testutils.Compare(t, expectedSellOrderResponse, res2)
}

func TestMultiMatchOrder1(t *testing.T) {
	e, _, _, _, _, _, _, factory1, factory2 := setupTest()
	defer e.redisConn.FlushAll()

	so1, _ := factory1.NewSellOrder(1e3+1, 1e8)
	so2, _ := factory1.NewSellOrder(1e3+2, 1e8)
	so3, _ := factory1.NewSellOrder(1e3+3, 1e8)
	bo1, _ := factory2.NewBuyOrder(1e3+4, 3e8)

	e.sellOrder(&so1)
	e.sellOrder(&so2)
	e.sellOrder(&so3)

	expso1 := so1
	expso1.Status = "FILLED"
	expso1.FilledAmount = big.NewInt(1e8)
	expso2 := so2
	expso2.Status = "FILLED"
	expso2.FilledAmount = big.NewInt(1e8)
	expso3 := so3
	expso3.Status = "FILLED"
	expso3.FilledAmount = big.NewInt(1e8)
	expbo1 := bo1
	expbo1.Status = "FILLED"
	expbo1.FilledAmount = big.NewInt(3e8)

	trade1, _ := types.NewUnsignedTrade1(&so1, &bo1, big.NewInt(1e8))
	trade2, _ := types.NewUnsignedTrade1(&so2, &bo1, big.NewInt(1e8))
	trade3, _ := types.NewUnsignedTrade1(&so3, &bo1, big.NewInt(1e8))

	expectedResponse := &types.EngineResponse{
		"FULL",
		&expbo1,
		nil,
		[]*types.FillOrder{{big.NewInt(1e8), &expso1}, {big.NewInt(1e8), &expso2}, {big.NewInt(1e8), &expso3}},
		[]*types.Trade{&trade1, &trade2, &trade3},
	}

	response, err := e.buyOrder(&bo1)
	if err != nil {
		t.Errorf("Error in sellOrder: %s", err)
	}

	testutils.Compare(t, expectedResponse, response)
}

// func TestMultiMatchOrder2(t *testing.T) {
// 	e := getResource()
// 	defer e.redisConn.FlushAll()
// 	buyOrder := getBuyOrder()

// 	buyOrder1 := buyOrder
// 	buyOrder1.PricePoint = math.Sub(buyOrder1.Price, big.NewInt(10))
// 	buyOrder1.Nonce = math.Add(buyOrder1.Nonce, big.NewInt(1))
// 	buyOrder1.Hash = buyOrder1.ComputeHash()

// 	buyOrder2 := buyOrder1
// 	buyOrder2.Nonce = math.Add(buyOrder2.Nonce, big.NewInt(1))
// 	buyOrder2.Hash = buyOrder2.ComputeHash()

// 	e.buyOrder(&buyOrder)
// 	e.buyOrder(&buyOrder1)
// 	e.buyOrder(&buyOrder2)

// 	sellOrder := getSellOrder()

// 	sellOrder.PricePoint = math.Sub(sellOrder.Price, big.NewInt(10))
// 	sellOrder.Amount = math.Mul(sellOrder.Amount, big.NewInt(3))

// 	// Test Case2: Send buyOrder first
// 	responseSO := sellOrder
// 	responseSO.FilledAmount = responseSO.Amount
// 	responseSO.Status = "FILLED"

// 	responseBO := buyOrder
// 	responseBO.FilledAmount = responseBO.Amount
// 	responseBO.Status = "FILLED"

// 	responseBO1 := buyOrder1
// 	responseBO1.FilledAmount = responseBO1.Amount
// 	responseBO1.Status = "FILLED"

// 	responseBO2 := buyOrder2
// 	responseBO2.FilledAmount = responseBO2.Amount
// 	responseBO2.Status = "FILLED"

// 	trade := getTrade(&sellOrder, &buyOrder, buyOrder.Amount, big.NewInt(0))
// 	trade1 := getTrade(&sellOrder, &buyOrder1, buyOrder1.Amount, big.NewInt(0))
// 	trade2 := getTrade(&sellOrder, &buyOrder2, buyOrder2.Amount, big.NewInt(0))

// 	trade.Hash = trade.ComputeHash()
// 	trade1.Hash = trade1.ComputeHash()
// 	trade2.Hash = trade2.ComputeHash()

// 	expectedResponse := getEResponse(&responseSO,
// 		[]*types.Trade{trade, trade1, trade2},
// 		"FULL",
// 		[]*types.FillOrder{
// 			{responseBO.FilledAmount, &responseBO},
// 			{responseBO1.FilledAmount, &responseBO1},
// 			{responseBO2.FilledAmount, &responseBO2}},
// 		nil)

// 	expBytes, _ := json.Marshal(expectedResponse)
// 	response, err := e.sellOrder(&sellOrder)
// 	if err != nil {
// 		t.Errorf("Error in sellOrder: %s", err)
// 	}

// 	resBytes, _ := json.Marshal(response)
// 	assert.JSONEq(t, string(expBytes), string(resBytes))

// }

// func TestPartialMatchOrder1(t *testing.T) {
// 	e := getResource()
// 	defer e.redisConn.FlushAll()

// 	sellOrder := getSellOrder()
// 	buyOrder := getBuyOrder()
// 	sellOrder1 := getSellOrder()

// 	sellOrder1.PricePoint = math.Add(sellOrder1.PricePoint, big.NewInt(10))
// 	sellOrder1.Nonce = math.Add(sellOrder1.Nonce, big.NewInt(1))
// 	sellOrder1.Hash = sellOrder1.ComputeHash()

// 	sellOrder2 := sellOrder1
// 	sellOrder2.Nonce = math.Add(sellOrder2.Nonce, big.NewInt(1))
// 	sellOrder2.Hash = sellOrder2.ComputeHash()

// 	sellOrder3 := sellOrder1
// 	sellOrder3.PricePoint = math.Add(sellOrder3.PricePoint, big.NewInt(10))
// 	sellOrder3.Amount = math.Mul(sellOrder3.Amount, big.NewInt(2))

// 	sellOrder3.Nonce = sellOrder3.Nonce.Add(sellOrder2.Nonce, big.NewInt(1))
// 	sellOrder3.Hash = sellOrder3.ComputeHash()

// 	e.sellOrder(&sellOrder)
// 	e.sellOrder(&sellOrder1)
// 	e.sellOrder(&sellOrder2)
// 	e.sellOrder(&sellOrder3)

// 	buyOrder.PricePoint = math.Add(buyOrder.PricePoint, big.NewInt(20))
// 	buyOrder.Amount = math.Mul(buyOrder.Amount, big.NewInt(4))

// 	// Test Case1: Send sellOrder first
// 	responseBO := buyOrder
// 	responseBO.FilledAmount = buyOrder.Amount
// 	responseBO.Status = "FILLED"
// 	responseSO := sellOrder
// 	responseSO.FilledAmount = responseSO.Amount
// 	responseSO.Status = "FILLED"
// 	responseSO1 := sellOrder1
// 	responseSO1.FilledAmount = responseSO1.Amount
// 	responseSO1.Status = "FILLED"
// 	responseSO2 := sellOrder2
// 	responseSO2.FilledAmount = responseSO2.Amount
// 	responseSO2.Status = "FILLED"
// 	responseSO3 := sellOrder3
// 	responseSO3.FilledAmount = math.Div(responseSO3.Amount, big.NewInt(2))
// 	responseSO3.Status = "PARTIAL_FILLED"

// 	trade := getTrade(&buyOrder, &sellOrder, responseSO.FilledAmount, big.NewInt(0))
// 	trade1 := getTrade(&buyOrder, &sellOrder1, responseSO1.FilledAmount, big.NewInt(0))
// 	trade2 := getTrade(&buyOrder, &sellOrder2, responseSO2.FilledAmount, big.NewInt(0))
// 	trade3 := getTrade(&buyOrder, &sellOrder3, responseSO3.FilledAmount, big.NewInt(0))

// 	trade.Hash = trade.ComputeHash()
// 	trade1.Hash = trade1.ComputeHash()
// 	trade2.Hash = trade2.ComputeHash()
// 	trade3.Hash = trade3.ComputeHash()

// 	expectedResponse := &types.EngineResponse{
// 		Order:      &responseBO,
// 		Trades:     []*types.Trade{trade, trade1, trade2, trade3},
// 		FillStatus: "FULL",
// 		MatchingOrders: []*types.FillOrder{
// 			{responseSO.FilledAmount, &responseSO},
// 			{responseSO1.FilledAmount, &responseSO1},
// 			{responseSO2.FilledAmount, &responseSO2},
// 			{responseSO3.FilledAmount, &responseSO3}},
// 	}

// 	erBytes, _ := json.Marshal(expectedResponse)
// 	response, err := e.buyOrder(&buyOrder)
// 	if err != nil {
// 		t.Errorf("Error in buyOrder: %s", err)
// 	}

// 	resBytes, _ := json.Marshal(response)
// 	assert.JSONEq(t, string(erBytes), string(resBytes))

// 	// Try matching remaining sellOrder with bigger buyOrder amount (partial filled buy Order)
// 	buyOrder = getBuyOrder()
// 	buyOrder.PricePoint = math.Add(buyOrder.PricePoint, big.NewInt(20))
// 	buyOrder.Amount = math.Mul(buyOrder.Amount, big.NewInt(2))

// 	responseBO = buyOrder
// 	responseBO.Status = "PARTIAL_FILLED"
// 	responseBO.FilledAmount = math.Div(buyOrder.Amount, big.NewInt(2))

// 	remOrder := getBuyOrder()
// 	remOrder.PricePoint = math.Add(remOrder.PricePoint, big.NewInt(20))
// 	remOrder.Hash = common.HexToHash("")
// 	remOrder.Signature = nil
// 	remOrder.Nonce = nil

// 	responseSO3.Status = "FILLED"
// 	responseSO3.FilledAmount = responseSO3.Amount
// 	trade4 := &types.Trade{
// 		Amount:     responseBO.FilledAmount,
// 		Price:      buyOrder.PricePoint,
// 		PricePoint: buyOrder.PricePoint,
// 		BaseToken:  buyOrder.BaseToken,
// 		QuoteToken: buyOrder.QuoteToken,
// 		OrderHash:  sellOrder3.Hash,
// 		Side:       buyOrder.Side,
// 		Taker:      buyOrder.UserAddress,
// 		PairName:   buyOrder.PairName,
// 		Maker:      sellOrder.UserAddress,
// 		TradeNonce: big.NewInt(0),
// 	}

// 	trade4.Hash = trade4.ComputeHash()
// 	expectedResponse = &types.EngineResponse{
// 		Order:          &responseBO,
// 		RemainingOrder: &remOrder,
// 		Trades:         []*types.Trade{trade4},
// 		FillStatus:     "PARTIAL",
// 		MatchingOrders: []*types.FillOrder{
// 			{responseBO.FilledAmount, &responseSO3}},
// 	}
// 	erBytes, _ = json.Marshal(expectedResponse)
// 	response, err = e.buyOrder(&buyOrder)
// 	if err != nil {
// 		t.Errorf("Error in buyOrder: %s", err)
// 	}

// 	resBytes, _ = json.Marshal(response)
// 	assert.JSONEq(t, string(erBytes), string(resBytes))
// }

// func TestPartialMatchOrder2(t *testing.T) {
// 	e := getResource()
// 	defer e.redisConn.FlushAll()

// 	sellOrder := getSellOrder()
// 	buyOrder := getBuyOrder()

// 	buyOrder1 := getBuyOrder()
// 	buyOrder1.PricePoint = math.Sub(buyOrder1.PricePoint, big.NewInt(10))
// 	buyOrder1.Nonce = math.Add(buyOrder1.Nonce, big.NewInt(1))
// 	buyOrder1.Hash = buyOrder1.ComputeHash()

// 	buyOrder2 := buyOrder1
// 	buyOrder2.Nonce = math.Add(buyOrder2.Nonce, big.NewInt(1))
// 	buyOrder2.Hash = buyOrder2.ComputeHash()

// 	buyOrder3 := buyOrder1
// 	buyOrder3.PricePoint = math.Sub(buyOrder3.PricePoint, big.NewInt(10))
// 	buyOrder3.Amount = math.Mul(buyOrder3.Amount, big.NewInt(2))
// 	buyOrder3.Nonce = buyOrder3.Nonce.Add(buyOrder2.Nonce, big.NewInt(1))
// 	buyOrder3.Hash = buyOrder3.ComputeHash()

// 	e.buyOrder(&buyOrder)
// 	e.buyOrder(&buyOrder1)
// 	e.buyOrder(&buyOrder2)
// 	e.buyOrder(&buyOrder3)

// 	sellOrder.PricePoint = math.Sub(sellOrder.PricePoint, big.NewInt(20))
// 	sellOrder.Amount = math.Mul(sellOrder.Amount, big.NewInt(4))

// 	// Test Case1: Send sellOrder first
// 	responseSO := sellOrder
// 	responseSO.FilledAmount = sellOrder.Amount
// 	responseSO.Status = "FILLED"
// 	responseBO := buyOrder
// 	responseBO.FilledAmount = responseBO.Amount
// 	responseBO.Status = "FILLED"
// 	responseBO1 := buyOrder1
// 	responseBO1.FilledAmount = responseBO1.Amount
// 	responseBO1.Status = "FILLED"
// 	responseBO2 := buyOrder2
// 	responseBO2.FilledAmount = responseBO2.Amount
// 	responseBO2.Status = "FILLED"
// 	responseBO3 := buyOrder3
// 	responseBO3.FilledAmount = math.Div(responseBO3.Amount, big.NewInt(2))
// 	responseBO3.Status = "PARTIAL_FILLED"

// 	trade := getTrade(&sellOrder, &buyOrder, responseBO.FilledAmount, big.NewInt(0))
// 	trade1 := getTrade(&sellOrder, &buyOrder1, responseBO1.FilledAmount, big.NewInt(0))
// 	trade2 := getTrade(&sellOrder, &buyOrder2, responseBO2.FilledAmount, big.NewInt(0))
// 	trade3 := getTrade(&sellOrder, &buyOrder3, responseBO3.FilledAmount, big.NewInt(0))

// 	trade.Hash = trade.ComputeHash()
// 	trade1.Hash = trade1.ComputeHash()
// 	trade2.Hash = trade2.ComputeHash()
// 	trade3.Hash = trade3.ComputeHash()

// 	expectedResponse := &types.EngineResponse{
// 		Order:      &responseSO,
// 		Trades:     []*types.Trade{trade, trade1, trade2, trade3},
// 		FillStatus: "FULL",
// 		MatchingOrders: []*types.FillOrder{
// 			{responseBO.FilledAmount, &responseBO},
// 			{responseBO1.FilledAmount, &responseBO1},
// 			{responseBO2.FilledAmount, &responseBO2},
// 			{responseBO3.FilledAmount, &responseBO3}},
// 	}

// 	erBytes, _ := json.Marshal(expectedResponse)
// 	response, err := e.sellOrder(&sellOrder)
// 	if err != nil {
// 		t.Errorf("Error in buyOrder: %s", err)
// 	}

// 	resBytes, _ := json.Marshal(response)
// 	assert.JSONEq(t, string(erBytes), string(resBytes))

// 	// Try matching remaining buyOrder with bigger sellOrder amount (partial filled sell Order)
// 	sellOrder = getSellOrder()
// 	sellOrder.PricePoint = math.Sub(sellOrder.PricePoint, big.NewInt(20))
// 	sellOrder.Amount = math.Mul(sellOrder.Amount, big.NewInt(2))

// 	responseSO = sellOrder
// 	responseSO.Status = "PARTIAL_FILLED"
// 	responseSO.FilledAmount = math.Div(sellOrder.Amount, big.NewInt(2))

// 	remOrder := getSellOrder()
// 	remOrder.PricePoint = math.Sub(remOrder.PricePoint, big.NewInt(20))
// 	remOrder.Hash = common.HexToHash("")
// 	remOrder.Signature = nil
// 	remOrder.Nonce = nil
// 	responseBO3.Status = "FILLED"
// 	responseBO3.FilledAmount = responseBO3.Amount

// 	trade4 := &types.Trade{
// 		Amount:     responseSO.FilledAmount,
// 		Price:      sellOrder.PricePoint,
// 		PricePoint: sellOrder.PricePoint,
// 		BaseToken:  sellOrder.BaseToken,
// 		QuoteToken: sellOrder.QuoteToken,
// 		OrderHash:  buyOrder3.Hash,
// 		Side:       sellOrder.Side,
// 		Taker:      sellOrder.UserAddress,
// 		PairName:   sellOrder.PairName,
// 		Maker:      buyOrder.UserAddress,
// 		TradeNonce: big.NewInt(0),
// 	}

// 	trade4.Hash = trade4.ComputeHash()
// 	expectedResponse = &types.EngineResponse{
// 		Order:          &responseSO,
// 		RemainingOrder: &remOrder,
// 		Trades:         []*types.Trade{trade4},
// 		FillStatus:     "PARTIAL",
// 		MatchingOrders: []*types.FillOrder{
// 			{responseSO.FilledAmount, &responseBO3}},
// 	}

// 	erBytes, _ = json.Marshal(expectedResponse)
// 	response, err = e.sellOrder(&sellOrder)
// 	if err != nil {
// 		t.Errorf("Error in buyOrder: %s", err)
// 	}

// 	resBytes, _ = json.Marshal(response)
// 	assert.JSONEq(t, string(erBytes), string(resBytes))
// }
