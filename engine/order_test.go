package engine

import (
	"log"
	"math/big"
	"strconv"
	"testing"

	"encoding/json"

	"github.com/Proofsuite/amp-matching-engine/types"
	"github.com/Proofsuite/amp-matching-engine/utils"
	"github.com/Proofsuite/amp-matching-engine/utils/math"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
)

func TestAddOrder(t *testing.T) {
	e := getResource()
	defer flushData(e.redisConn)

	o1 := getSellOrder()

	bytes, _ := o1.MarshalJSON()

	e.addOrder(&o1)
	ssKey, listKey := o1.GetOBKeys()

	rs, err := getSortedSet(e.redisConn, ssKey)
	if err != nil {
		t.Error(err)
	} else {
		for k, v := range rs {
			assert.Equalf(t, utils.UintToPaddedString(o1.PricePoint.Int64()), k, "Expected sorted set value: %v got: %v", utils.UintToPaddedString(o1.PricePoint.Int64()), k)
			assert.Equalf(t, 0.0, v, "Expected sorted set value: %v got: %v", 0, v)
		}
	}

	rs, err = getSortedSet(e.redisConn, listKey)
	if err != nil {
		t.Error(err)
	} else {
		for k, v := range rs {
			assert.Equalf(t, o1.Hash.Hex(), k, "Expected sorted set value: %v got: %v", o1.Hash.Hex(), k)
			assert.Equalf(t, o1.CreatedAt.Unix(), int64(v), "Expected sorted set value: %v got: %v", o1.CreatedAt.Unix(), v)
		}
	}

	rse, err := getValue(e.redisConn, o1.Hash.Hex())
	if err != nil {
		t.Error(err)
	}

	assert.JSONEqf(t, string(bytes), rse, "Expected value for key: %v, was: %s, but got: %v", ssKey, bytes, rse)
	rse, err = getValue(e.redisConn, ssKey+"::book::"+utils.UintToPaddedString(o1.PricePoint.Int64()))
	if err != nil {
		t.Error(err)
	}

	// Currently converting amount to int64. In the future, we need to use strings instead of int64
	assert.Equal(t, o1.Amount.Sub(o1.Amount, o1.FilledAmount).String(), rse)

	o2 := getBuyOrder()

	bytes, _ = o2.MarshalJSON()
	e.addOrder(&o2)
	ssKey, listKey = o2.GetOBKeys()

	rs, err = getSortedSet(e.redisConn, ssKey)
	if err != nil {
		t.Error(err)
	} else {
		var matched = false
		for k, v := range rs {
			if utils.UintToPaddedString(o2.PricePoint.Int64()) == k && v == 0.0 {
				matched = true
			}
		}
		if !matched {
			t.Errorf("Expected sorted set value: %v", utils.UintToPaddedString(o2.PricePoint.Int64()))
		}
	}

	rs, err = getSortedSet(e.redisConn, listKey)
	if err != nil {
		t.Error(err)
	} else {
		var matched = false
		for k, v := range rs {
			if o2.Hash.Hex() == k && o2.CreatedAt.Unix() == int64(v) {
				matched = true
			}
		}

		if !matched {
			t.Errorf("Expected sorted set value: %v ", o2.Hash)
		}
	}
	rse, err = getValue(e.redisConn, o2.Hash.Hex())
	if err != nil {
		t.Error(err)
	}
	assert.JSONEqf(t, string(bytes), rse, "Expected value for key: %v, was: %s, but got: %v", ssKey, bytes, rse)
	rse, err = getValue(e.redisConn, ssKey+"::book::"+utils.UintToPaddedString(o2.PricePoint.Int64()))
	if err != nil {
		t.Error(err)
	}
	rse, err = getValue(e.redisConn, ssKey+"::book::"+utils.UintToPaddedString(o1.PricePoint.Int64()))
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, o2.Amount.Sub(o2.Amount, o2.FilledAmount).String(), rse)

}

func TestUpdateOrder(t *testing.T) {
	e := getResource()
	defer flushData(e.redisConn)
	sampleOrder := getSellOrder()
	e.addOrder(&sampleOrder)

	tAmount := big.NewInt(1000000000)
	e.updateOrder(&sampleOrder, tAmount)
	ssKey, listKey := sampleOrder.GetOBKeys()
	updatedSampleOrder := sampleOrder

	updatedSampleOrder.Status = "PARTIAL_FILLED"
	updatedSampleOrder.FilledAmount = big.NewInt(1000000000)
	updatedSampleOrderBytes, _ := updatedSampleOrder.MarshalJSON()
	rs, err := getSortedSet(e.redisConn, ssKey)
	if err != nil {
		t.Error(err)
	} else {
		for k, v := range rs {
			assert.Equal(t, utils.UintToPaddedString(sampleOrder.Price.Int64()), k)
			assert.Equal(t, 0.0, v)
		}
	}
	rs, err = getSortedSet(e.redisConn, listKey)
	if err != nil {
		t.Error(err)
	} else {
		for k, v := range rs {
			assert.Equal(t, sampleOrder.Hash.Hex(), k)
			assert.Equal(t, sampleOrder.CreatedAt.Unix(), int64(v))
		}
	}
	rse, err := getValue(e.redisConn, sampleOrder.Hash.Hex())
	if err != nil {
		t.Error(err)
	}
	assert.JSONEq(t, string(updatedSampleOrderBytes), rse)
	rse, err = getValue(e.redisConn, ssKey+"::book::"+utils.UintToPaddedString(sampleOrder.PricePoint.Int64()))
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, math.Sub(sampleOrder.Amount, tAmount).String(), rse)
	// Update the order again with negative amount(recover order)
	e.updateOrder(&sampleOrder, math.Mul(tAmount, big.NewInt(-1)))
	rs, err = getSortedSet(e.redisConn, ssKey)
	if err != nil {
		t.Error(err)
	} else {
		for k, v := range rs {
			assert.Equal(t, utils.UintToPaddedString(sampleOrder.PricePoint.Int64()), k)
			assert.Equalf(t, 0.0, v, "Expected sorted set value: %v got: %v", 0, v)
		}
	}
	rs, err = getSortedSet(e.redisConn, listKey)
	if err != nil {
		t.Error(err)
	} else {
		for k, v := range rs {
			assert.Equal(t, sampleOrder.Hash.Hex(), k)
			assert.Equal(t, sampleOrder.CreatedAt.Unix(), int64(v))
		}
	}
	rse, err = getValue(e.redisConn, sampleOrder.Hash.Hex())
	if err != nil {
		t.Error(err)
	}
	o := getSellOrder()
	o.Status = "OPEN"
	orderJSON, _ := o.MarshalJSON()
	assert.JSONEq(t, string(orderJSON), rse)
	rse, err = getValue(e.redisConn, ssKey+"::book::"+utils.UintToPaddedString(sampleOrder.PricePoint.Int64()))
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, sampleOrder.Amount.String(), rse)
}

func TestDeleteOrder(t *testing.T) {
	e := getResource()
	defer flushData(e.redisConn)
	// Add 1st order
	sampleOrder := getSellOrder()
	e.addOrder(&sampleOrder)
	// Add 2nd order
	sampleOrder1 := getSellOrder()
	sampleOrder1.Hash = common.HexToHash("0xa2d800b77828cb52c83106ca392e465bc0af0d7c319f6956328f739080c19621")
	e.addOrder(&sampleOrder1)
	// delete sampleOrder1
	e.deleteOrder(&sampleOrder1, sampleOrder1.Amount)
	ss1Key, list1Key := sampleOrder1.GetOBKeys()
	rs, err := getSortedSet(e.redisConn, ss1Key)
	if err != nil {
		t.Error(err)
	} else {
		for k, v := range rs {
			assert.Equal(t, utils.UintToPaddedString(sampleOrder1.PricePoint.Int64()), k)
			assert.Equal(t, 0.0, v)
		}
	}
	rs, err = getSortedSet(e.redisConn, list1Key)
	if err != nil {
		t.Error(err)
	} else {
		for k, v := range rs {
			assert.NotEqual(t, sampleOrder1.Hash, k)
			assert.Equal(t, sampleOrder1.CreatedAt.Unix(), int64(v))
		}
	}
	if exists(e.redisConn, list1Key+"::"+sampleOrder1.Hash.Hex()) {
		t.Errorf("Key : %v expected to be deleted but key exists", list1Key+"::"+sampleOrder1.Hash.Hex())
	}
	rse, err := getValue(e.redisConn, ss1Key+"::book::"+utils.UintToPaddedString(sampleOrder1.PricePoint.Int64()))
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, sampleOrder.Amount.String(), rse)
	// delete sampleOrder
	e.deleteOrder(&sampleOrder, sampleOrder.Amount)
	ssKey, listKey := sampleOrder.GetOBKeys()
	// All keys must have been deleted
	if exists(e.redisConn, ssKey) {
		t.Errorf("Key : %v expected to be deleted but key exists", ssKey)
	}
	if exists(e.redisConn, listKey) {
		t.Errorf("Key : %v expected to be deleted but key exists", ssKey)
	}
	if exists(e.redisConn, listKey+"::"+sampleOrder.Hash.Hex()) {
		t.Errorf("Key : %v expected to be deleted but key exists", ssKey)
	}
	if exists(e.redisConn, ssKey+"::book::"+utils.UintToPaddedString(sampleOrder.PricePoint.Int64())) {
		t.Errorf("Key : %v expected to be deleted but key exists", ssKey)
	}
}

func TestCancelOrder(t *testing.T) {
	e := getResource()
	defer flushData(e.redisConn)

	// Add 1st order (OPEN order)
	o1 := getSellOrder()

	e.addOrder(&o1)

	// Add 2nd order (Partial filled order)
	o2 := getSellOrder()
	o2.Hash = common.HexToHash("0x10b30eb0072a4f0a38b6fca0b731cba15eb2e1702845d97c1230b53a839bcb85")
	e.addOrder(&o2)

	expectedResponse := getEResponse(&o2, make([]*types.Trade, 0), CANCELLED, make([]*FillOrder, 0), &types.Order{})
	expectedResponse.Order.Status = "CANCELLED"
	// cancel o2
	response, err := e.CancelOrder(&o2)

	if err != nil {
		log.Print("Error while cancelling SampleOrder1: ", err.Error())
	}

	assert.Equal(t, expectedResponse, response)

	ss1Key, list1Key := o2.GetOBKeys()
	rs, err := getSortedSet(e.redisConn, ss1Key)
	if err != nil {
		t.Error(err)
	} else {
		for k, v := range rs {
			assert.Equalf(t, utils.UintToPaddedString(o2.PricePoint.Int64()), k, "Expected sorted set value: %v got: %v", utils.UintToPaddedString(o2.PricePoint.Int64()), k)
			assert.Equalf(t, 0.0, v, "Expected sorted set value: %v got: %v", 0, v)
		}
	}

	rs, err = getSortedSet(e.redisConn, list1Key)
	if err != nil {
		t.Error(err)
	} else {
		for k, v := range rs {
			assert.NotEqualf(t, o2.Hash, k, "Key : %v expected to be deleted but key exists", o2.Hash)
			assert.Equalf(t, o2.CreatedAt.Unix(), int64(v), "Expected sorted set value: %v got: %v", o2.CreatedAt.Unix(), v)
		}
	}

	if exists(e.redisConn, list1Key+"::"+o2.Hash.Hex()) {
		t.Errorf("Key : %v expected to be deleted but key exists", list1Key+"::"+o2.Hash.Hex())
	}

	rse, err := getValue(e.redisConn, ss1Key+"::book::"+utils.UintToPaddedString(o2.PricePoint.Int64()))
	if err != nil {
		t.Error(err)
	}

	// Currently converting amount to int64. In the future, we need to use strings instead of int64
	a := new(big.Int)
	assert.Equal(t, a.Sub(o1.Amount, o1.FilledAmount).String(), rse)

	expectedResponse.Order = &o1
	expectedResponse.Order.Status = "CANCELLED"

	// cancel o1
	response, err = e.CancelOrder(&o1)
	ssKey, listKey := o1.GetOBKeys()

	if err != nil {
		t.Errorf("Error while cancelling SampleOrder: %s", err)
	}

	assert.Equal(t, expectedResponse, response)

	// All keys must have been deleted
	if exists(e.redisConn, ssKey) {
		t.Errorf("Key : %v expected to be deleted but key exists", ssKey)
	}

	if exists(e.redisConn, listKey) {
		t.Errorf("Key : %v expected to be deleted but key exists", ssKey)
	}

	if exists(e.redisConn, o1.Hash.Hex()) {
		t.Errorf("Key : %v expected to be deleted but key exists", ssKey)
	}

	if exists(e.redisConn, ssKey+"::book::"+utils.UintToPaddedString(o1.PricePoint.Int64())) {
		t.Errorf("Key : %v expected to be deleted but key exists", ssKey)
	}
}

func TestRecoverOrders(t *testing.T) {
	e := getResource()
	defer flushData(e.redisConn)

	o1 := getSellOrder()
	o1.Status = "PARTIAL_FILLED"
	o1.FilledAmount = big.NewInt(1000000000)

	expectedOrder1 := o1
	expectedOrder1.Status = "OPEN"
	expectedOrder1.FilledAmount = big.NewInt(0)
	expectedOrder1Json, _ := expectedOrder1.MarshalJSON()

	e.addOrder(&o1)

	o2 := getSellOrder()
	o2.Hash = common.HexToHash("0xa2d800b77828cb52c83106ca392e465bc0af0d7c319f6956328f739080d19621")
	o2.FilledAmount = o2.Amount
	o2.Status = "FILLED"

	expectedOrder2 := o2
	expectedOrder2.Status = "OPEN"
	expectedOrder2.FilledAmount = big.NewInt(0)
	expectedOrder2Json, _ := expectedOrder2.MarshalJSON()

	e.addOrder(&o2)

	o3 := getSellOrder()
	o3.Hash = common.HexToHash("0xa2d800b77828cb52c83106ca392e465bc0af0d7c319f6956328f739180d19621")
	o3.FilledAmount = o2.Amount
	o3.Status = "FILLED"

	expectedOrder3 := o3
	expectedOrder3.Status = "PARTIAL_FILLED"
	expectedOrder3.FilledAmount = big.NewInt(4000000000)
	expectedOrder3Json, _ := expectedOrder3.MarshalJSON()

	e.addOrder(&o3)

	recoverOrders := []*FillOrder{
		&FillOrder{
			Amount: big.NewInt(1000000000),
			Order:  &o1,
		}, &FillOrder{
			Amount: big.NewInt(6000000000),
			Order:  &o2,
		}, &FillOrder{
			Amount: big.NewInt(2000000000),
			Order:  &o3,
		},
	}

	// cancel o2
	err := e.RecoverOrders(recoverOrders)

	if err != nil {
		t.Errorf("Error while recovering Orders: %s", err)
	}

	ssKey, listKey := o1.GetOBKeys()

	rs, err := getSortedSet(e.redisConn, ssKey)
	if err != nil {
		t.Error(err)
	} else {
		for k, v := range rs {
			assert.Equalf(t, utils.UintToPaddedString(o1.PricePoint.Int64()), k, "Expected sorted set value: %v got: %v", utils.UintToPaddedString(o2.PricePoint.Int64()), k)
			assert.Equalf(t, 0.0, v, "Expected sorted set value: %v got: %v", 0, v)
		}
	}

	expectedMap := map[string]float64{
		o1.Hash.Hex(): float64(o1.CreatedAt.Unix()),
		o2.Hash.Hex(): float64(o2.CreatedAt.Unix()),
		o3.Hash.Hex(): float64(o3.CreatedAt.Unix()),
	}

	rs, err = getSortedSet(e.redisConn, listKey)
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, expectedMap, rs)

	rse, err := getValue(e.redisConn, o1.Hash.Hex())
	if err != nil {
		t.Error(err)
	}

	assert.JSONEq(t, string(expectedOrder1Json), rse)

	rse, err = getValue(e.redisConn, o2.Hash.Hex())
	if err != nil {
		t.Error(err)
	}

	assert.JSONEq(t, string(expectedOrder2Json), rse)

	rse, err = getValue(e.redisConn, o3.Hash.Hex())
	if err != nil {
		t.Error(err)
	}

	assert.JSONEq(t, string(expectedOrder3Json), rse)

	rse, err = getValue(e.redisConn, ssKey+"::book::"+utils.UintToPaddedString(o1.PricePoint.Int64()))
	if err != nil {
		t.Error(err)
	}

	expectedAmt := o1.Amount.Int64() + o2.Amount.Int64() + o3.Amount.Int64() - o3.FilledAmount.Int64()
	assert.Equalf(t, strconv.FormatInt(expectedAmt, 10), rse, "Expected value for key: %v, was: %v, but got: %v", ssKey, expectedAmt, rse)
}

func TestSellOrder(t *testing.T) {
	e := getResource()
	defer flushData(e.redisConn)

	o1 := getSellOrder()

	expectedResponse := getEResponse(&o1, make([]*types.Trade, 0), NOMATCH, make([]*FillOrder, 0), nil)

	expectedResponse.Order.Status = "OPEN"

	response, err := e.sellOrder(&o1)
	if err != nil {
		t.Errorf("Error in sellOrder: %s", err)
	}
	assert.Equal(t, expectedResponse, response)
}

func TestBuyOrder(t *testing.T) {
	e := getResource()
	defer flushData(e.redisConn)

	o1 := getBuyOrder()

	expectedResponse := getEResponse(&o1, make([]*types.Trade, 0), NOMATCH, make([]*FillOrder, 0), nil)

	expectedResponse.Order.Status = "OPEN"

	response, err := e.buyOrder(&o1)
	if err != nil {
		t.Errorf("Error in sellOrder: %s", err)
	}
	assert.Equal(t, expectedResponse, response)
}

// NOTE: Can only be tested with redis server, as SORT command is not supported by miniredis,YET
// TODO: Replace with miniredis as soon as it supports these functions
// TestFillOrder: Test Case1: Send sellOrder first
func TestFillOrder1(t *testing.T) {
	e := getResource()
	defer flushData(e.redisConn)
	sellOrder := getSellOrder()

	buyOrder := getBuyOrder()
	// Test Case1: Send sellOrder first
	expectedResponseBuyOrder := buyOrder
	expectedResponseBuyOrder.FilledAmount = expectedResponseBuyOrder.Amount
	expectedResponseBuyOrder.Status = "FILLED"
	expectedResponseSellOrder := sellOrder
	expectedResponseSellOrder.Status = "OPEN"

	expectedResponse := getEResponse(&expectedResponseSellOrder, make([]*types.Trade, 0), NOMATCH, make([]*FillOrder, 0), nil)

	expBytes, _ := json.Marshal(expectedResponse)
	resp, err := e.sellOrder(&sellOrder)
	if err != nil {
		t.Errorf("Error in sellOrder: %s", err)
	}

	resBytes, _ := json.Marshal(resp)
	assert.JSONEq(t, string(expBytes), string(resBytes))

	expectedResponseSellOrder.FilledAmount = expectedResponseSellOrder.Amount
	expectedResponseSellOrder.Status = "FILLED"
	trade := getTrade(&buyOrder, &sellOrder, expectedResponseBuyOrder.Amount, big.NewInt(0))
	trade.Hash = trade.ComputeHash()

	expectedResponse = getEResponse(&expectedResponseBuyOrder, []*types.Trade{trade}, FULL, []*FillOrder{{buyOrder.Amount, &expectedResponseSellOrder}}, nil)

	expBytes, _ = json.Marshal(expectedResponse)
	resp, err = e.buyOrder(&buyOrder)
	if err != nil {
		t.Errorf("Error in buyOrder: %s", err)
		return
	}
	resBytes, _ = json.Marshal(resp)
	assert.JSONEq(t, string(expBytes), string(resBytes))
}

// TestFillOrder1: Test Case2: Send buyOrder first
func TestFillOrder2(t *testing.T) {
	e := getResource()
	defer flushData(e.redisConn)
	sellOrder := getSellOrder()

	buyOrder := getBuyOrder()

	expectedResponseBuyOrder := buyOrder
	expectedResponseBuyOrder.Status = "OPEN"

	expectedResponse := getEResponse(&expectedResponseBuyOrder, make([]*types.Trade, 0), NOMATCH, make([]*FillOrder, 0), nil)

	erBytes, _ := json.Marshal(expectedResponse)
	response, err := e.buyOrder(&buyOrder)
	if err != nil {
		t.Errorf("Error in sellOrder: %s", err)
	}

	resBytes, _ := json.Marshal(response)
	assert.JSONEq(t, string(erBytes), string(resBytes))

	expectedResponseSellOrder := sellOrder
	expectedResponseSellOrder.FilledAmount = expectedResponseSellOrder.Amount
	expectedResponseSellOrder.Status = "FILLED"

	expectedResponseBuyOrder.FilledAmount = expectedResponseBuyOrder.Amount
	expectedResponseBuyOrder.Status = "FILLED"

	trade := getTrade(&sellOrder, &buyOrder, expectedResponseSellOrder.Amount, big.NewInt(0))

	trade.Hash = trade.ComputeHash()

	expectedResponse = getEResponse(&expectedResponseSellOrder, []*types.Trade{trade}, FULL, []*FillOrder{{buyOrder.Amount, &expectedResponseBuyOrder}}, nil)

	bytes, _ := json.Marshal(expectedResponse)
	response, err = e.sellOrder(&sellOrder)
	if err != nil {
		t.Errorf("Error in sellOrder: %s", err)
	}

	resBytes, _ = json.Marshal(response)
	assert.JSONEq(t, string(bytes), string(resBytes))
}

func TestMultiMatchOrder1(t *testing.T) {
	e := getResource()
	defer flushData(e.redisConn)
	sellOrder := getSellOrder()
	sellOrder.FilledAmount = big.NewInt(0)

	sellOrder1 := sellOrder
	sellOrder1.PricePoint = math.Add(sellOrder.PricePoint, big.NewInt(10))
	sellOrder1.Nonce = math.Add(sellOrder.Nonce, big.NewInt(1))
	sellOrder1.Hash = sellOrder1.ComputeHash()

	sellOrder2 := sellOrder1
	sellOrder2.Nonce = sellOrder2.Nonce.Add(sellOrder2.Nonce, big.NewInt(1))
	sellOrder2.Hash = sellOrder2.ComputeHash()

	e.sellOrder(&sellOrder)
	e.sellOrder(&sellOrder1)
	e.sellOrder(&sellOrder2)

	buyOrder := getBuyOrder()

	buyOrder.PricePoint = math.Add(buyOrder.PricePoint, big.NewInt(10))
	buyOrder.Amount = math.Mul(buyOrder.Amount, big.NewInt(3))

	// Test Case1: Send sellOrder first
	responseBO := buyOrder
	responseBO.FilledAmount = responseBO.Amount
	responseBO.Status = "FILLED"
	responseSO := sellOrder
	responseSO.FilledAmount = responseSO.Amount
	responseSO.Status = "FILLED"
	responseSO1 := sellOrder1
	responseSO1.FilledAmount = responseSO1.Amount
	responseSO1.Status = "FILLED"
	responseSO2 := sellOrder2
	responseSO2.FilledAmount = responseSO2.Amount
	responseSO2.Status = "FILLED"

	trade := getTrade(&buyOrder, &sellOrder, responseSO.Amount, big.NewInt(0))
	trade1 := getTrade(&buyOrder, &sellOrder1, responseSO1.Amount, big.NewInt(0))
	trade2 := getTrade(&buyOrder, &sellOrder2, responseSO2.Amount, big.NewInt(0))

	trade.Hash = trade.ComputeHash()
	trade1.Hash = trade1.ComputeHash()
	trade2.Hash = trade2.ComputeHash()

	expectedResponse := getEResponse(&responseBO,
		[]*types.Trade{trade, trade1, trade2},
		FULL,
		[]*FillOrder{{responseSO.FilledAmount, &responseSO}, {responseSO1.FilledAmount, &responseSO1}, {responseSO2.FilledAmount, &responseSO2}},
		nil)

	expBytes, _ := json.Marshal(expectedResponse)
	response, err := e.buyOrder(&buyOrder)
	if err != nil {
		t.Errorf("Error in sellOrder: %s", err)
	}

	resBytes, _ := json.Marshal(response)
	assert.JSONEq(t, string(expBytes), string(resBytes))
}

func TestMultiMatchOrder2(t *testing.T) {
	e := getResource()
	defer flushData(e.redisConn)
	buyOrder := getBuyOrder()

	buyOrder1 := buyOrder
	buyOrder1.PricePoint = math.Sub(buyOrder1.Price, big.NewInt(10))
	buyOrder1.Nonce = math.Add(buyOrder1.Nonce, big.NewInt(1))
	buyOrder1.Hash = buyOrder1.ComputeHash()

	buyOrder2 := buyOrder1
	buyOrder2.Nonce = math.Add(buyOrder2.Nonce, big.NewInt(1))
	buyOrder2.Hash = buyOrder2.ComputeHash()

	e.buyOrder(&buyOrder)
	e.buyOrder(&buyOrder1)
	e.buyOrder(&buyOrder2)

	sellOrder := getSellOrder()

	sellOrder.PricePoint = math.Sub(sellOrder.Price, big.NewInt(10))
	sellOrder.Amount = math.Mul(sellOrder.Amount, big.NewInt(3))

	// Test Case2: Send buyOrder first
	responseSO := sellOrder
	responseSO.FilledAmount = responseSO.Amount
	responseSO.Status = "FILLED"

	responseBO := buyOrder
	responseBO.FilledAmount = responseBO.Amount
	responseBO.Status = "FILLED"

	responseBO1 := buyOrder1
	responseBO1.FilledAmount = responseBO1.Amount
	responseBO1.Status = "FILLED"

	responseBO2 := buyOrder2
	responseBO2.FilledAmount = responseBO2.Amount
	responseBO2.Status = "FILLED"

	trade := getTrade(&sellOrder, &buyOrder, buyOrder.Amount, big.NewInt(0))
	trade1 := getTrade(&sellOrder, &buyOrder1, buyOrder1.Amount, big.NewInt(0))
	trade2 := getTrade(&sellOrder, &buyOrder2, buyOrder2.Amount, big.NewInt(0))

	trade.Hash = trade.ComputeHash()
	trade1.Hash = trade1.ComputeHash()
	trade2.Hash = trade2.ComputeHash()

	expectedResponse := getEResponse(&responseSO,
		[]*types.Trade{trade, trade1, trade2},
		FULL,
		[]*FillOrder{
			{responseBO.FilledAmount, &responseBO},
			{responseBO1.FilledAmount, &responseBO1},
			{responseBO2.FilledAmount, &responseBO2}},
		nil)

	expBytes, _ := json.Marshal(expectedResponse)
	response, err := e.sellOrder(&sellOrder)
	if err != nil {
		t.Errorf("Error in sellOrder: %s", err)
	}

	resBytes, _ := json.Marshal(response)
	assert.JSONEq(t, string(expBytes), string(resBytes))

}

func TestPartialMatchOrder1(t *testing.T) {
	e := getResource()
	defer flushData(e.redisConn)

	sellOrder := getSellOrder()
	buyOrder := getBuyOrder()
	sellOrder1 := getSellOrder()

	sellOrder1.PricePoint = math.Add(sellOrder1.PricePoint, big.NewInt(10))
	sellOrder1.Nonce = math.Add(sellOrder1.Nonce, big.NewInt(1))
	sellOrder1.Hash = sellOrder1.ComputeHash()

	sellOrder2 := sellOrder1
	sellOrder2.Nonce = math.Add(sellOrder2.Nonce, big.NewInt(1))
	sellOrder2.Hash = sellOrder2.ComputeHash()

	sellOrder3 := sellOrder1
	sellOrder3.PricePoint = math.Add(sellOrder3.PricePoint, big.NewInt(10))
	sellOrder3.Amount = math.Mul(sellOrder3.Amount, big.NewInt(2))

	sellOrder3.Nonce = sellOrder3.Nonce.Add(sellOrder2.Nonce, big.NewInt(1))
	sellOrder3.Hash = sellOrder3.ComputeHash()

	e.sellOrder(&sellOrder)
	e.sellOrder(&sellOrder1)
	e.sellOrder(&sellOrder2)
	e.sellOrder(&sellOrder3)

	buyOrder.PricePoint = math.Add(buyOrder.PricePoint, big.NewInt(20))
	buyOrder.Amount = math.Mul(buyOrder.Amount, big.NewInt(4))

	// Test Case1: Send sellOrder first
	responseBO := buyOrder
	responseBO.FilledAmount = buyOrder.Amount
	responseBO.Status = "FILLED"
	responseSO := sellOrder
	responseSO.FilledAmount = responseSO.Amount
	responseSO.Status = "FILLED"
	responseSO1 := sellOrder1
	responseSO1.FilledAmount = responseSO1.Amount
	responseSO1.Status = "FILLED"
	responseSO2 := sellOrder2
	responseSO2.FilledAmount = responseSO2.Amount
	responseSO2.Status = "FILLED"
	responseSO3 := sellOrder3
	responseSO3.FilledAmount = math.Div(responseSO3.Amount, big.NewInt(2))
	responseSO3.Status = "PARTIAL_FILLED"

	trade := getTrade(&buyOrder, &sellOrder, responseSO.FilledAmount, big.NewInt(0))
	trade1 := getTrade(&buyOrder, &sellOrder1, responseSO1.FilledAmount, big.NewInt(0))
	trade2 := getTrade(&buyOrder, &sellOrder2, responseSO2.FilledAmount, big.NewInt(0))
	trade3 := getTrade(&buyOrder, &sellOrder3, responseSO3.FilledAmount, big.NewInt(0))

	trade.Hash = trade.ComputeHash()
	trade1.Hash = trade1.ComputeHash()
	trade2.Hash = trade2.ComputeHash()
	trade3.Hash = trade3.ComputeHash()

	expectedResponse := &Response{
		Order:      &responseBO,
		Trades:     []*types.Trade{trade, trade1, trade2, trade3},
		FillStatus: FULL,
		MatchingOrders: []*FillOrder{
			{responseSO.FilledAmount, &responseSO},
			{responseSO1.FilledAmount, &responseSO1},
			{responseSO2.FilledAmount, &responseSO2},
			{responseSO3.FilledAmount, &responseSO3}},
	}

	erBytes, _ := json.Marshal(expectedResponse)
	response, err := e.buyOrder(&buyOrder)
	if err != nil {
		t.Errorf("Error in buyOrder: %s", err)
	}

	resBytes, _ := json.Marshal(response)
	assert.JSONEq(t, string(erBytes), string(resBytes))

	// Try matching remaining sellOrder with bigger buyOrder amount (partial filled buy Order)
	buyOrder = getBuyOrder()
	buyOrder.PricePoint = math.Add(buyOrder.PricePoint, big.NewInt(20))
	buyOrder.Amount = math.Mul(buyOrder.Amount, big.NewInt(2))

	responseBO = buyOrder
	responseBO.Status = "PARTIAL_FILLED"
	responseBO.FilledAmount = math.Div(buyOrder.Amount, big.NewInt(2))

	remOrder := getBuyOrder()
	remOrder.PricePoint = math.Add(remOrder.PricePoint, big.NewInt(20))
	remOrder.Hash = common.HexToHash("")
	remOrder.Signature = nil
	remOrder.Nonce = nil

	responseSO3.Status = "FILLED"
	responseSO3.FilledAmount = responseSO3.Amount
	trade4 := &types.Trade{
		Amount:     responseBO.FilledAmount,
		Price:      buyOrder.PricePoint,
		PricePoint: buyOrder.PricePoint,
		BaseToken:  buyOrder.BaseToken,
		QuoteToken: buyOrder.QuoteToken,
		OrderHash:  sellOrder3.Hash,
		Side:       buyOrder.Side,
		Taker:      buyOrder.UserAddress,
		PairName:   buyOrder.PairName,
		Maker:      sellOrder.UserAddress,
		TradeNonce: big.NewInt(0),
	}

	trade4.Hash = trade4.ComputeHash()
	expectedResponse = &Response{
		Order:          &responseBO,
		RemainingOrder: &remOrder,
		Trades:         []*types.Trade{trade4},
		FillStatus:     PARTIAL,
		MatchingOrders: []*FillOrder{
			{responseBO.FilledAmount, &responseSO3}},
	}
	erBytes, _ = json.Marshal(expectedResponse)
	response, err = e.buyOrder(&buyOrder)
	if err != nil {
		t.Errorf("Error in buyOrder: %s", err)
	}

	resBytes, _ = json.Marshal(response)
	assert.JSONEq(t, string(erBytes), string(resBytes))
}

func TestPartialMatchOrder2(t *testing.T) {
	e := getResource()
	defer flushData(e.redisConn)

	sellOrder := getSellOrder()
	buyOrder := getBuyOrder()

	buyOrder1 := getBuyOrder()
	buyOrder1.PricePoint = math.Sub(buyOrder1.PricePoint, big.NewInt(10))
	buyOrder1.Nonce = math.Add(buyOrder1.Nonce, big.NewInt(1))
	buyOrder1.Hash = buyOrder1.ComputeHash()

	buyOrder2 := buyOrder1
	buyOrder2.Nonce = math.Add(buyOrder2.Nonce, big.NewInt(1))
	buyOrder2.Hash = buyOrder2.ComputeHash()

	buyOrder3 := buyOrder1
	buyOrder3.PricePoint = math.Sub(buyOrder3.PricePoint, big.NewInt(10))
	buyOrder3.Amount = math.Mul(buyOrder3.Amount, big.NewInt(2))
	buyOrder3.Nonce = buyOrder3.Nonce.Add(buyOrder2.Nonce, big.NewInt(1))
	buyOrder3.Hash = buyOrder3.ComputeHash()

	e.buyOrder(&buyOrder)
	e.buyOrder(&buyOrder1)
	e.buyOrder(&buyOrder2)
	e.buyOrder(&buyOrder3)

	sellOrder.PricePoint = math.Sub(sellOrder.PricePoint, big.NewInt(20))
	sellOrder.Amount = math.Mul(sellOrder.Amount, big.NewInt(4))

	// Test Case1: Send sellOrder first
	responseSO := sellOrder
	responseSO.FilledAmount = sellOrder.Amount
	responseSO.Status = "FILLED"
	responseBO := buyOrder
	responseBO.FilledAmount = responseBO.Amount
	responseBO.Status = "FILLED"
	responseBO1 := buyOrder1
	responseBO1.FilledAmount = responseBO1.Amount
	responseBO1.Status = "FILLED"
	responseBO2 := buyOrder2
	responseBO2.FilledAmount = responseBO2.Amount
	responseBO2.Status = "FILLED"
	responseBO3 := buyOrder3
	responseBO3.FilledAmount = math.Div(responseBO3.Amount, big.NewInt(2))
	responseBO3.Status = "PARTIAL_FILLED"

	trade := getTrade(&sellOrder, &buyOrder, responseBO.FilledAmount, big.NewInt(0))
	trade1 := getTrade(&sellOrder, &buyOrder1, responseBO1.FilledAmount, big.NewInt(0))
	trade2 := getTrade(&sellOrder, &buyOrder2, responseBO2.FilledAmount, big.NewInt(0))
	trade3 := getTrade(&sellOrder, &buyOrder3, responseBO3.FilledAmount, big.NewInt(0))

	trade.Hash = trade.ComputeHash()
	trade1.Hash = trade1.ComputeHash()
	trade2.Hash = trade2.ComputeHash()
	trade3.Hash = trade3.ComputeHash()

	expectedResponse := &Response{
		Order:      &responseSO,
		Trades:     []*types.Trade{trade, trade1, trade2, trade3},
		FillStatus: FULL,
		MatchingOrders: []*FillOrder{
			{responseBO.FilledAmount, &responseBO},
			{responseBO1.FilledAmount, &responseBO1},
			{responseBO2.FilledAmount, &responseBO2},
			{responseBO3.FilledAmount, &responseBO3}},
	}

	erBytes, _ := json.Marshal(expectedResponse)
	response, err := e.sellOrder(&sellOrder)
	if err != nil {
		t.Errorf("Error in buyOrder: %s", err)
	}

	resBytes, _ := json.Marshal(response)
	assert.JSONEq(t, string(erBytes), string(resBytes))

	// Try matching remaining buyOrder with bigger sellOrder amount (partial filled sell Order)
	sellOrder = getSellOrder()
	sellOrder.PricePoint = math.Sub(sellOrder.PricePoint, big.NewInt(20))
	sellOrder.Amount = math.Mul(sellOrder.Amount, big.NewInt(2))

	responseSO = sellOrder
	responseSO.Status = "PARTIAL_FILLED"
	responseSO.FilledAmount = math.Div(sellOrder.Amount, big.NewInt(2))

	remOrder := getSellOrder()
	remOrder.PricePoint = math.Sub(remOrder.PricePoint, big.NewInt(20))
	remOrder.Hash = common.HexToHash("")
	remOrder.Signature = nil
	remOrder.Nonce = nil
	responseBO3.Status = "FILLED"
	responseBO3.FilledAmount = responseBO3.Amount

	trade4 := &types.Trade{
		Amount:     responseSO.FilledAmount,
		Price:      sellOrder.PricePoint,
		PricePoint: sellOrder.PricePoint,
		BaseToken:  sellOrder.BaseToken,
		QuoteToken: sellOrder.QuoteToken,
		OrderHash:  buyOrder3.Hash,
		Side:       sellOrder.Side,
		Taker:      sellOrder.UserAddress,
		PairName:   sellOrder.PairName,
		Maker:      buyOrder.UserAddress,
		TradeNonce: big.NewInt(0),
	}

	trade4.Hash = trade4.ComputeHash()
	expectedResponse = &Response{
		Order:          &responseSO,
		RemainingOrder: &remOrder,
		Trades:         []*types.Trade{trade4},
		FillStatus:     PARTIAL,
		MatchingOrders: []*FillOrder{
			{responseSO.FilledAmount, &responseBO3}},
	}

	erBytes, _ = json.Marshal(expectedResponse)
	response, err = e.sellOrder(&sellOrder)
	if err != nil {
		t.Errorf("Error in buyOrder: %s", err)
	}

	resBytes, _ = json.Marshal(response)
	assert.JSONEq(t, string(erBytes), string(resBytes))
}

func getBuyOrder() types.Order {
	return types.Order{
		UserAddress:     common.HexToAddress("0x7a9f3cd060ab180f36c17fe6bdf9974f577d77aa"),
		ExchangeAddress: common.HexToAddress("0xae55690d4b079460e6ac28aaa58c9ec7b73a7485"),
		BuyToken:        common.HexToAddress("0xe41d2489571d322189246dafa5ebde1f4699f498"),
		SellToken:       common.HexToAddress("0x12459c951127e0c374ff9105dda097662a027093"),
		BaseToken:       common.HexToAddress("0xe41d2489571d322189246dafa5ebde1f4699f498"),
		QuoteToken:      common.HexToAddress("0x12459c951127e0c374ff9105dda097662a027093"),
		BuyAmount:       big.NewInt(6000000000),
		SellAmount:      big.NewInt(13800000000),
		Price:           big.NewInt(229999999),
		PricePoint:      big.NewInt(229999999),
		Amount:          big.NewInt(6000000000),
		FilledAmount:    big.NewInt(0),
		Status:          "NEW",
		Side:            "BUY",
		PairName:        "ZRX/WETH",
		Expires:         big.NewInt(10000),
		MakeFee:         big.NewInt(50),
		Nonce:           big.NewInt(1000),
		TakeFee:         big.NewInt(50),
		Signature: &types.Signature{
			V: 28,
			R: common.HexToHash("0x10b30eb0072a4f0a38b6fca0b731cba15eb2e1702845d97c1230b53a839bcb85"),
			S: common.HexToHash("0x6d9ad89548c9e3ce4c97825d027291477f2c44a8caef792095f2cabc978493ff"),
		},
		Hash: common.HexToHash("0xb9070a2d333403c255ce71ddf6e795053599b2e885321de40353832b96d8880a"),
	}
}
func getSellOrder() types.Order {
	return types.Order{
		UserAddress:     common.HexToAddress("0x7a9f3cd060ab180f36c17fe6bdf9974f577d77aa"),
		ExchangeAddress: common.HexToAddress("0xae55690d4b079460e6ac28aaa58c9ec7b73a7485"),
		BuyToken:        common.HexToAddress("0xe41d2489571d322189246dafa5ebde1f4699f498"),
		SellToken:       common.HexToAddress("0x12459c951127e0c374ff9105dda097662a027093"),
		BaseToken:       common.HexToAddress("0xe41d2489571d322189246dafa5ebde1f4699f498"),
		QuoteToken:      common.HexToAddress("0x12459c951127e0c374ff9105dda097662a027093"),
		BuyAmount:       big.NewInt(6000000000),
		SellAmount:      big.NewInt(13800000000),
		Price:           big.NewInt(229999999),
		PricePoint:      big.NewInt(229999999),
		Amount:          big.NewInt(6000000000),
		FilledAmount:    big.NewInt(0),
		Status:          "NEW",
		Side:            "SELL",
		PairName:        "ZRX/WETH",
		Expires:         big.NewInt(10000),
		MakeFee:         big.NewInt(50),
		Nonce:           big.NewInt(1000),
		TakeFee:         big.NewInt(50),
		Signature: &types.Signature{
			V: 28,
			R: common.HexToHash("0x10b30eb0072a4f0a38b6fca0b731cba15eb2e1702845d97c1230b53a839bcb85"),
			S: common.HexToHash("0x6d9ad89548c9e3ce4c97825d027291477f2c44a8caef792095f2cabc978493ff"),
		},
		Hash: common.HexToHash("0xb9070a2d333403c255ce71ddf6e795053599b2e885321de40353832b96d8880a"),
	}
}

func getEResponse(order *types.Order, trades []*types.Trade, status FillStatus, matchingOrders []*FillOrder, remOrder *types.Order) *Response {
	return &Response{
		Order:          order,
		FillStatus:     status,
		Trades:         trades,
		MatchingOrders: matchingOrders,
		RemainingOrder: remOrder,
	}
}

func getTrade(takerOrder *types.Order, makerOrder *types.Order, amount, nonce *big.Int) *types.Trade {
	return &types.Trade{
		Amount:     amount,
		Price:      takerOrder.PricePoint,
		PricePoint: takerOrder.PricePoint,
		BaseToken:  takerOrder.BaseToken,
		QuoteToken: takerOrder.QuoteToken,
		OrderHash:  makerOrder.Hash,
		Side:       takerOrder.Side,
		Taker:      takerOrder.UserAddress,
		PairName:   takerOrder.PairName,
		Maker:      makerOrder.UserAddress,
		TradeNonce: nonce,
	}
}
