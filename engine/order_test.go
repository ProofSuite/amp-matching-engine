package engine

import (
	"encoding/json"
	"strconv"
	"testing"

	"github.com/Proofsuite/amp-matching-engine/types"
	"github.com/Proofsuite/amp-matching-engine/utils"
	"github.com/stretchr/testify/assert"
)

func TestAddOrder(t *testing.T) {
	e, s := getResource()
	defer s.Close()

	sampleOrder := &types.Order{}
	orderJSON := []byte(`{ "id": "5b6ac5297b4457546d64379c", "buyToken": "0x1888a8db0b7db59413ce07150b3373972bf818d3", "sellToken": "0x2034842261b82651885751fc293bba7ba5398156", "baseToken": "0x2034842261b82651885751fc293bba7ba5398156", "quoteToken": "0x1888a8db0b7db59413ce07150b3373972bf818d3", "buyAmount": 6000000000, "sellAmount": 13800000000, "nonce": 0, "userAddress": "0xefD7eB287CeeFCE8256Dd46e25F398acEA7C4b63", "hash": "0xa2d800b77828cb52c83106ca392e465bc0af0d7c319f6956328f739080c19620", "side": "SELL", "amount": 6000000000, "price": 229999999, "filledAmount": 0, "fee": 0, "makeFee": 0, "takeFee": 0, "exchangeAddress": "", "status": "OPEN", "pairID": "5b6ac5117b445753ee755fb8", "pairName": "HPC/AUT", "orderBook": null, "createdAt": "2018-08-08T15:55:45.062141954+05:30", "updatedAt": "2018-08-08T15:55:45.06214208+05:30" }`)
	json.Unmarshal(orderJSON, &sampleOrder)
	e.addOrder(sampleOrder)
	ssKey, listKey := sampleOrder.GetOBKeys()

	rs, err := s.SortedSet(ssKey)
	if err != nil {
		t.Error(err)
	} else {
		for k, v := range rs {
			assert.Equalf(t, utils.UintToPaddedString(sampleOrder.Price), k, "Expected sorted set value: %v got: %v", utils.UintToPaddedString(sampleOrder.Price), k)
			assert.Equalf(t, 0.0, v, "Expected sorted set value: %v got: %v", 0, v)
		}
	}

	rs, err = s.SortedSet(listKey)
	if err != nil {
		t.Error(err)
	} else {
		for k, v := range rs {
			assert.Equalf(t, sampleOrder.Hash, k, "Expected sorted set value: %v got: %v", sampleOrder.Hash, k)
			assert.Equalf(t, sampleOrder.CreatedAt.Unix(), int64(v), "Expected sorted set value: %v got: %v", sampleOrder.CreatedAt.Unix(), v)
		}
	}

	rse, err := s.Get(listKey + "::" + sampleOrder.Hash.Hex())
	if err != nil {
		t.Error(err)
	}
	assert.JSONEqf(t, string(orderJSON), rse, "Expected value for key: %v, was: %s, but got: %v", ssKey, orderJSON, rse)

	rse, err = s.Get(ssKey + "::book::" + utils.UintToPaddedString(sampleOrder.Price))
	if err != nil {
		t.Error(err)
	}
	assert.Equalf(t, strconv.FormatInt(sampleOrder.Amount, 10), rse, "Expected value for key: %v, was: %v, but got: %v", ssKey, sampleOrder.Amount, rse)

	// Partial order addidtion test
	sampleOrder1 := &types.Order{}
	orderJSON = []byte(`{ "id": "5b6ac5297b4457546d64379d", "buyToken": "0x1888a8db0b7db59413ce07150b3373972bf818d3", "sellToken": "0x2034842261b82651885751fc293bba7ba5398156", "baseToken": "0x2034842261b82651885751fc293bba7ba5398156", "quoteToken": "0x1888a8db0b7db59413ce07150b3373972bf818d3", "buyAmount": 6000000000, "sellAmount": 13800000000, "nonce": 0, "userAddress": "0xefD7eB287CeeFCE8256Dd46e25F398acEA7C4b63", "hash": "0xa2d800b77828cb52c83106ca392e465bc0af0d7c319f6956328f739080c19621", "side": "SELL", "amount": 6000000000, "price": 229999999, "filledAmount": 1000000000, "fee": 0, "makeFee": 0, "takeFee": 0, "exchangeAddress": "", "status": "PARTIAL_FILLED", "pairID": "5b6ac5117b445753ee755fb8", "pairName": "HPC/AUT", "orderBook": null, "createdAt": "2018-08-08T15:55:45.062141954+05:30", "updatedAt": "2018-08-08T15:55:45.06214208+05:30" }`)
	json.Unmarshal(orderJSON, &sampleOrder1)
	e.addOrder(sampleOrder1)
	ssKey, listKey = sampleOrder1.GetOBKeys()

	rs, err = s.SortedSet(ssKey)
	if err != nil {
		t.Error(err)
	} else {

		var matched = false
		for k, v := range rs {
			if utils.UintToPaddedString(sampleOrder1.Price) == k && v == 0.0 {
				matched = true
			}
		}
		if !matched {
			t.Errorf("Expected sorted set value: %v", utils.UintToPaddedString(sampleOrder1.Price))
		}
	}

	rs, err = s.SortedSet(listKey)
	if err != nil {
		t.Error(err)
	} else {
		var matched = false
		for k, v := range rs {
			if sampleOrder1.Hash.Hex() == k && sampleOrder1.CreatedAt.Unix() == int64(v) {
				matched = true
			}
		}
		if !matched {
			t.Errorf("Expected sorted set value: %v ", sampleOrder1.Hash)
		}

	}

	rse, err = s.Get(listKey + "::" + sampleOrder1.Hash.Hex())
	if err != nil {
		t.Error(err)
	}
	assert.JSONEqf(t, string(orderJSON), rse, "Expected value for key: %v, was: %s, but got: %v", ssKey, orderJSON, rse)

	rse, err = s.Get(ssKey + "::book::" + utils.UintToPaddedString(sampleOrder1.Price))
	if err != nil {
		t.Error(err)
	}
	assert.Equalf(t, strconv.FormatInt(sampleOrder1.Amount+sampleOrder.Amount-sampleOrder1.FilledAmount, 10), rse, "Expected value for key: %v, was: %v, but got: %v", ssKey, sampleOrder1.Amount+sampleOrder.Amount-sampleOrder1.FilledAmount, rse)

}

func TestUpdateOrder(t *testing.T) {
	e, s := getResource()
	defer s.Close()

	sampleOrder := &types.Order{}
	orderJSON := []byte(`{ "id": "5b6ac5297b4457546d64379c", "buyToken": "0x1888a8db0b7db59413ce07150b3373972bf818d3", "sellToken": "0x2034842261b82651885751fc293bba7ba5398156", "baseToken": "0x2034842261b82651885751fc293bba7ba5398156", "quoteToken": "0x1888a8db0b7db59413ce07150b3373972bf818d3", "buyAmount": 6000000000, "sellAmount": 13800000000, "nonce": 0, "userAddress": "0xefD7eB287CeeFCE8256Dd46e25F398acEA7C4b63", "hash": "0xa2d800b77828cb52c83106ca392e465bc0af0d7c319f6956328f739080c19620", "side": "SELL", "amount": 6000000000, "price": 229999999, "filledAmount": 0, "fee": 0, "makeFee": 0, "takeFee": 0, "exchangeAddress": "", "status": "OPEN", "pairID": "5b6ac5117b445753ee755fb8", "pairName": "HPC/AUT", "orderBook": null, "createdAt": "2018-08-08T15:55:45.062141954+05:30", "updatedAt": "2018-08-08T15:55:45.06214208+05:30" }`)
	json.Unmarshal(orderJSON, &sampleOrder)
	e.addOrder(sampleOrder)

	var tAmount int64 = 1000000000
	e.updateOrder(sampleOrder, tAmount)
	ssKey, listKey := sampleOrder.GetOBKeys()
	updatedSampleOrder := *sampleOrder
	updatedSampleOrder.FilledAmount = 1000000000
	updatedSampleOrderBytes, _ := json.Marshal(updatedSampleOrder)
	rs, err := s.SortedSet(ssKey)
	if err != nil {
		t.Error(err)
	} else {
		for k, v := range rs {
			assert.Equalf(t, utils.UintToPaddedString(sampleOrder.Price), k, "Expected sorted set value: %v got: %v", utils.UintToPaddedString(sampleOrder.Price-tAmount), k)
			assert.Equalf(t, 0.0, v, "Expected sorted set value: %v got: %v", 0, v)
		}
	}

	rs, err = s.SortedSet(listKey)
	if err != nil {
		t.Error(err)
	} else {
		for k, v := range rs {
			assert.Equalf(t, sampleOrder.Hash, k, "Expected sorted set value: %v got: %v", sampleOrder.Hash, k)
			assert.Equalf(t, sampleOrder.CreatedAt.Unix(), int64(v), "Expected sorted set value: %v got: %v", sampleOrder.CreatedAt.Unix(), v)
		}
	}

	rse, err := s.Get(listKey + "::" + sampleOrder.Hash.Hex())
	if err != nil {
		t.Error(err)
	}
	assert.JSONEqf(t, string(updatedSampleOrderBytes), rse, "Expected value for key: %v, was: %s, but got: %v", ssKey, updatedSampleOrderBytes, rse)

	rse, err = s.Get(ssKey + "::book::" + utils.UintToPaddedString(sampleOrder.Price))
	if err != nil {
		t.Error(err)
	}
	assert.Equalf(t, strconv.FormatInt(sampleOrder.Amount-tAmount, 10), rse, "Expected value for key: %v, was: %v, but got: %v", ssKey, sampleOrder.Amount-tAmount, rse)

	// Update the order again with negative amount(recover order)
	e.updateOrder(sampleOrder, -1*tAmount)

	rs, err = s.SortedSet(ssKey)
	if err != nil {
		t.Error(err)
	} else {
		for k, v := range rs {
			assert.Equalf(t, utils.UintToPaddedString(sampleOrder.Price), k, "Expected sorted set value: %v got: %v", utils.UintToPaddedString(sampleOrder.Price), k)
			assert.Equalf(t, 0.0, v, "Expected sorted set value: %v got: %v", 0, v)
		}
	}

	rs, err = s.SortedSet(listKey)
	if err != nil {
		t.Error(err)
	} else {
		for k, v := range rs {
			assert.Equalf(t, sampleOrder.Hash, k, "Expected sorted set value: %v got: %v", sampleOrder.Hash, k)
			assert.Equalf(t, sampleOrder.CreatedAt.Unix(), int64(v), "Expected sorted set value: %v got: %v", sampleOrder.CreatedAt.Unix(), v)
		}
	}

	rse, err = s.Get(listKey + "::" + sampleOrder.Hash.Hex())
	if err != nil {
		t.Error(err)
	}
	assert.JSONEqf(t, string(orderJSON), rse, "Expected value for key: %v, was: %s, but got: %v", ssKey, orderJSON, rse)

	rse, err = s.Get(ssKey + "::book::" + utils.UintToPaddedString(sampleOrder.Price))
	if err != nil {
		t.Error(err)
	}
	assert.Equalf(t, strconv.FormatInt(sampleOrder.Amount, 10), rse, "Expected value for key: %v, was: %v, but got: %v", ssKey, sampleOrder.Amount, rse)

}

func TestDeleteOrder(t *testing.T) {
	e, s := getResource()
	defer s.Close()

	// Add 1st order
	sampleOrder := &types.Order{}
	orderJSON := []byte(`{ "id": "5b6ac5297b4457546d64379c", "buyToken": "0x1888a8db0b7db59413ce07150b3373972bf818d3", "sellToken": "0x2034842261b82651885751fc293bba7ba5398156", "baseToken": "0x2034842261b82651885751fc293bba7ba5398156", "quoteToken": "0x1888a8db0b7db59413ce07150b3373972bf818d3", "buyAmount": 6000000000, "sellAmount": 13800000000, "nonce": 0, "userAddress": "0xefD7eB287CeeFCE8256Dd46e25F398acEA7C4b63", "hash": "0xa2d800b77828cb52c83106ca392e465bc0af0d7c319f6956328f739080c19620", "side": "SELL", "amount": 6000000000, "price": 229999999, "filledAmount": 0, "fee": 0, "makeFee": 0, "takeFee": 0, "exchangeAddress": "", "status": "OPEN", "pairID": "5b6ac5117b445753ee755fb8", "pairName": "HPC/AUT", "orderBook": null, "createdAt": "2018-08-08T15:55:45.062141954+05:30", "updatedAt": "2018-08-08T15:55:45.06214208+05:30" }`)
	json.Unmarshal(orderJSON, &sampleOrder)
	e.addOrder(sampleOrder)

	// Add 2nd order
	sampleOrder1 := &types.Order{}
	orderJSON1 := []byte(`{ "id": "5b6ac5297b4457546d64379d", "buyToken": "0x1888a8db0b7db59413ce07150b3373972bf818d3", "sellToken": "0x2034842261b82651885751fc293bba7ba5398156", "baseToken": "0x2034842261b82651885751fc293bba7ba5398156", "quoteToken": "0x1888a8db0b7db59413ce07150b3373972bf818d3", "buyAmount": 6000000000, "sellAmount": 13800000000, "nonce": 0, "userAddress": "0xefD7eB287CeeFCE8256Dd46e25F398acEA7C4b63", "hash": "0xa2d800b77828cb52c83106ca392e465bc0af0d7c319f6956328f739080c19621", "side": "SELL", "amount": 6000000000, "price": 229999999, "filledAmount": 0, "fee": 0, "makeFee": 0, "takeFee": 0, "exchangeAddress": "", "status": "OPEN", "pairID": "5b6ac5117b445753ee755fb8", "pairName": "HPC/AUT", "orderBook": null, "createdAt": "2018-08-08T15:55:45.062141954+05:30", "updatedAt": "2018-08-08T15:55:45.06214208+05:30" }`)
	json.Unmarshal(orderJSON1, &sampleOrder1)
	e.addOrder(sampleOrder1)

	// delete sampleOrder1
	e.deleteOrder(sampleOrder1, sampleOrder1.Amount)
	ss1Key, list1Key := sampleOrder1.GetOBKeys()

	rs, err := s.SortedSet(ss1Key)
	if err != nil {
		t.Error(err)
	} else {
		for k, v := range rs {
			assert.Equalf(t, utils.UintToPaddedString(sampleOrder1.Price), k, "Expected sorted set value: %v got: %v", utils.UintToPaddedString(sampleOrder1.Price), k)
			assert.Equalf(t, 0.0, v, "Expected sorted set value: %v got: %v", 0, v)
		}
	}

	rs, err = s.SortedSet(list1Key)
	if err != nil {
		t.Error(err)
	} else {
		for k, v := range rs {
			assert.NotEqualf(t, sampleOrder1.Hash, k, "Key : %v expected to be deleted but key exists", sampleOrder1.Hash)
			assert.Equalf(t, sampleOrder1.CreatedAt.Unix(), int64(v), "Expected sorted set value: %v got: %v", sampleOrder1.CreatedAt.Unix(), v)
		}
	}

	if s.Exists(list1Key + "::" + sampleOrder1.Hash.Hex()) {
		t.Errorf("Key : %v expected to be deleted but key exists", list1Key+"::"+sampleOrder1.Hash.Hex())
	}

	rse, err := s.Get(ss1Key + "::book::" + utils.UintToPaddedString(sampleOrder1.Price))
	if err != nil {
		t.Error(err)
	}
	assert.Equalf(t, strconv.FormatInt(sampleOrder.Amount, 10), rse, "Expected value for key: %v, was: %v, but got: %v", ss1Key, sampleOrder.Amount, rse)

	// delete sampleOrder
	e.deleteOrder(sampleOrder, sampleOrder.Amount)
	ssKey, listKey := sampleOrder.GetOBKeys()

	// All keys must have been deleted
	if s.Exists(ssKey) {
		t.Errorf("Key : %v expected to be deleted but key exists", ssKey)
	}
	if s.Exists(listKey) {
		t.Errorf("Key : %v expected to be deleted but key exists", ssKey)
	}
	if s.Exists(listKey + "::" + sampleOrder.Hash.Hex()) {
		t.Errorf("Key : %v expected to be deleted but key exists", ssKey)
	}
	if s.Exists(ssKey + "::book::" + utils.UintToPaddedString(sampleOrder.Price)) {
		t.Errorf("Key : %v expected to be deleted but key exists", ssKey)
	}
}

func TestCancelOrder(t *testing.T) {
	e, s := getResource()
	defer s.Close()

	// Add 1st order (OPEN order)
	sampleOrder := &types.Order{}
	orderJSON := []byte(`{ "id": "5b6ac5297b4457546d64379c", "buyToken": "0x1888a8db0b7db59413ce07150b3373972bf818d3", "sellToken": "0x2034842261b82651885751fc293bba7ba5398156", "baseToken": "0x2034842261b82651885751fc293bba7ba5398156", "quoteToken": "0x1888a8db0b7db59413ce07150b3373972bf818d3", "buyAmount": 6000000000, "sellAmount": 13800000000, "nonce": 0, "userAddress": "0xefD7eB287CeeFCE8256Dd46e25F398acEA7C4b63", "hash": "0xa2d800b77828cb52c83106ca392e465bc0af0d7c319f6956328f739080c19620", "side": "SELL", "amount": 6000000000, "price": 229999999, "filledAmount": 0, "fee": 0, "makeFee": 0, "takeFee": 0, "exchangeAddress": "", "status": "OPEN", "pairID": "5b6ac5117b445753ee755fb8", "pairName": "HPC/AUT", "orderBook": null, "createdAt": "2018-08-08T15:55:45.062141954+05:30", "updatedAt": "2018-08-08T15:55:45.06214208+05:30" }`)
	json.Unmarshal(orderJSON, &sampleOrder)
	e.addOrder(sampleOrder)

	// Add 2nd order (Partial filled order)
	sampleOrder1 := &types.Order{}
	orderJSON1 := []byte(`{ "id": "5b6ac5297b4457546d64379d", "buyToken": "0x1888a8db0b7db59413ce07150b3373972bf818d3", "sellToken": "0x2034842261b82651885751fc293bba7ba5398156", "baseToken": "0x2034842261b82651885751fc293bba7ba5398156", "quoteToken": "0x1888a8db0b7db59413ce07150b3373972bf818d3", "buyAmount": 6000000000, "sellAmount": 13800000000, "nonce": 0, "userAddress": "0xefD7eB287CeeFCE8256Dd46e25F398acEA7C4b63", "hash": "0xa2d800b77828cb52c83106ca392e465bc0af0d7c319f6956328f739080c19621", "side": "SELL", "amount": 6000000000, "price": 229999999, "filledAmount": 1000000000, "fee": 0, "makeFee": 0, "takeFee": 0, "exchangeAddress": "", "status": "PARTIAL_FILLED", "pairID": "5b6ac5117b445753ee755fb8", "pairName": "HPC/AUT", "orderBook": null, "createdAt": "2018-08-08T15:55:45.062141954+05:30", "updatedAt": "2018-08-08T15:55:45.06214208+05:30" }`)
	json.Unmarshal(orderJSON1, &sampleOrder1)
	e.addOrder(sampleOrder1)

	expectedResponse := &Response{
		Order:          sampleOrder1,
		Trades:         make([]*types.Trade, 0),
		RemainingOrder: &types.Order{},
		FillStatus:     CANCELLED,
		MatchingOrders: make([]*FillOrder, 0),
	}
	expectedResponse.Order.Status = "CANCELLED"

	// cancel sampleOrder1
	response, err := e.CancelOrder(sampleOrder1)

	if err != nil {
		t.Errorf("Error while cancelling SampleOrder1: %s", err)
	}

	assert.Equalf(t, expectedResponse, response, "Expected Response: %+v got: %+v", expectedResponse, response)

	ss1Key, list1Key := sampleOrder1.GetOBKeys()

	rs, err := s.SortedSet(ss1Key)
	if err != nil {
		t.Error(err)
	} else {
		for k, v := range rs {
			assert.Equalf(t, utils.UintToPaddedString(sampleOrder1.Price), k, "Expected sorted set value: %v got: %v", utils.UintToPaddedString(sampleOrder1.Price), k)
			assert.Equalf(t, 0.0, v, "Expected sorted set value: %v got: %v", 0, v)
		}
	}

	rs, err = s.SortedSet(list1Key)
	if err != nil {
		t.Error(err)
	} else {
		for k, v := range rs {
			assert.NotEqualf(t, sampleOrder1.Hash, k, "Key : %v expected to be deleted but key exists", sampleOrder1.Hash)
			assert.Equalf(t, sampleOrder1.CreatedAt.Unix(), int64(v), "Expected sorted set value: %v got: %v", sampleOrder1.CreatedAt.Unix(), v)
		}
	}

	if s.Exists(list1Key + "::" + sampleOrder1.Hash.Hex()) {
		t.Errorf("Key : %v expected to be deleted but key exists", list1Key+"::"+sampleOrder1.Hash.Hex())
	}

	rse, err := s.Get(ss1Key + "::book::" + utils.UintToPaddedString(sampleOrder1.Price))
	if err != nil {
		t.Error(err)
	}
	assert.Equalf(t, strconv.FormatInt(sampleOrder.Amount, 10), rse, "Expected value for key: %v, was: %v, but got: %v", ss1Key, sampleOrder1.Amount-sampleOrder1.FilledAmount, rse)

	expectedResponse.Order = sampleOrder
	expectedResponse.Order.Status = "CANCELLED"

	// cancel sampleOrder
	response, err = e.CancelOrder(sampleOrder)
	ssKey, listKey := sampleOrder.GetOBKeys()

	if err != nil {
		t.Errorf("Error while cancelling SampleOrder: %s", err)
	}

	assert.Equalf(t, expectedResponse, response, "Expected Response from cancelling sampleorder : %+v got: %+v", expectedResponse, response)
	// All keys must have been deleted
	if s.Exists(ssKey) {
		t.Errorf("Key : %v expected to be deleted but key exists", ssKey)
	}
	if s.Exists(listKey) {
		t.Errorf("Key : %v expected to be deleted but key exists", ssKey)
	}
	if s.Exists(listKey + "::" + sampleOrder.Hash.Hex()) {
		t.Errorf("Key : %v expected to be deleted but key exists", ssKey)
	}
	if s.Exists(ssKey + "::book::" + utils.UintToPaddedString(sampleOrder.Price)) {
		t.Errorf("Key : %v expected to be deleted but key exists", ssKey)
	}
}

func TestRecoverOrders(t *testing.T) {
	e, s := getResource()
	defer s.Close()

	// Add Partial filled order
	sampleOrder := &types.Order{}
	orderJSON := []byte(`{ "id": "5b6ac5297b4457546d64379d", "buyToken": "0x1888a8db0b7db59413ce07150b3373972bf818d3", "sellToken": "0x2034842261b82651885751fc293bba7ba5398156", "baseToken": "0x2034842261b82651885751fc293bba7ba5398156", "quoteToken": "0x1888a8db0b7db59413ce07150b3373972bf818d3", "buyAmount": 6000000000, "sellAmount": 13800000000, "nonce": 0, "userAddress": "0xefD7eB287CeeFCE8256Dd46e25F398acEA7C4b63", "hash": "0xa2d800b77828cb52c83106ca392e465bc0af0d7c319f6956328f739080c19621", "side": "SELL", "amount": 6000000000, "price": 229999999, "filledAmount": 1000000000, "fee": 0, "makeFee": 0, "takeFee": 0, "exchangeAddress": "", "status": "PARTIAL_FILLED", "pairID": "5b6ac5117b445753ee755fb8", "pairName": "HPC/AUT", "orderBook": null, "createdAt": "2018-08-08T15:55:45.062141954+05:30", "updatedAt": "2018-08-08T15:55:45.06214208+05:30" }`)
	expectedOrderJSON := []byte(`{ "id": "5b6ac5297b4457546d64379d", "buyToken": "0x1888a8db0b7db59413ce07150b3373972bf818d3", "sellToken": "0x2034842261b82651885751fc293bba7ba5398156", "baseToken": "0x2034842261b82651885751fc293bba7ba5398156", "quoteToken": "0x1888a8db0b7db59413ce07150b3373972bf818d3", "buyAmount": 6000000000, "sellAmount": 13800000000, "nonce": 0, "userAddress": "0xefD7eB287CeeFCE8256Dd46e25F398acEA7C4b63", "hash": "0xa2d800b77828cb52c83106ca392e465bc0af0d7c319f6956328f739080c19621", "side": "SELL", "amount": 6000000000, "price": 229999999, "filledAmount": 0, "fee": 0, "makeFee": 0, "takeFee": 0, "exchangeAddress": "", "status": "OPEN", "pairID": "5b6ac5117b445753ee755fb8", "pairName": "HPC/AUT", "orderBook": null, "createdAt": "2018-08-08T15:55:45.062141954+05:30", "updatedAt": "2018-08-08T15:55:45.06214208+05:30" }`)
	json.Unmarshal(orderJSON, &sampleOrder)
	e.addOrder(sampleOrder)

	// Order to be recovered(not added) to test recovery of completely filled matched orders
	sampleOrder1 := &types.Order{}
	orderJSON1 := []byte(`{ "id": "5b6ac5297b4457546d64379c", "buyToken": "0x1888a8db0b7db59413ce07150b3373972bf818d3", "sellToken": "0x2034842261b82651885751fc293bba7ba5398156", "baseToken": "0x2034842261b82651885751fc293bba7ba5398156", "quoteToken": "0x1888a8db0b7db59413ce07150b3373972bf818d3", "buyAmount": 6000000000, "sellAmount": 13800000000, "nonce": 0, "userAddress": "0xefD7eB287CeeFCE8256Dd46e25F398acEA7C4b63", "hash": "0xa2d800b77828cb52c83106ca392e465bc0af0d7c319f6956328f739080c19620", "side": "SELL", "amount": 6000000000, "price": 229999999, "filledAmount": 6000000000, "fee": 0, "makeFee": 0, "takeFee": 0, "exchangeAddress": "", "status": "FILLED", "pairID": "5b6ac5117b445753ee755fb8", "pairName": "HPC/AUT", "orderBook": null, "createdAt": "2018-08-08T15:55:45.062141954+05:30", "updatedAt": "2018-08-08T15:55:45.06214208+05:30" }`)
	expectedOrderJSON1 := []byte(`{ "id": "5b6ac5297b4457546d64379c", "buyToken": "0x1888a8db0b7db59413ce07150b3373972bf818d3", "sellToken": "0x2034842261b82651885751fc293bba7ba5398156", "baseToken": "0x2034842261b82651885751fc293bba7ba5398156", "quoteToken": "0x1888a8db0b7db59413ce07150b3373972bf818d3", "buyAmount": 6000000000, "sellAmount": 13800000000, "nonce": 0, "userAddress": "0xefD7eB287CeeFCE8256Dd46e25F398acEA7C4b63", "hash": "0xa2d800b77828cb52c83106ca392e465bc0af0d7c319f6956328f739080c19620", "side": "SELL", "amount": 6000000000, "price": 229999999, "filledAmount": 0, "fee": 0, "makeFee": 0, "takeFee": 0, "exchangeAddress": "", "status": "OPEN", "pairID": "5b6ac5117b445753ee755fb8", "pairName": "HPC/AUT", "orderBook": null, "createdAt": "2018-08-08T15:55:45.062141954+05:30", "updatedAt": "2018-08-08T15:55:45.06214208+05:30" }`)
	json.Unmarshal(orderJSON1, &sampleOrder1)

	// Order to be recovered(not added) to test recovery of completely filled order which was matched partially with failed order
	sampleOrder2 := &types.Order{}
	orderJSON2 := []byte(`{ "id": "5b6ac5297b4457546d64379e", "buyToken": "0x1888a8db0b7db59413ce07150b3373972bf818d3", "sellToken": "0x2034842261b82651885751fc293bba7ba5398156", "baseToken": "0x2034842261b82651885751fc293bba7ba5398156", "quoteToken": "0x1888a8db0b7db59413ce07150b3373972bf818d3", "buyAmount": 6000000000, "sellAmount": 13800000000, "nonce": 0, "userAddress": "0xefD7eB287CeeFCE8256Dd46e25F398acEA7C4b63", "hash": "0xa2d800b77828cb52c83106ca392e465bc0af0d7c319f6956328f739080c19622", "side": "SELL", "amount": 6000000000, "price": 229999999, "filledAmount": 6000000000, "fee": 0, "makeFee": 0, "takeFee": 0, "exchangeAddress": "", "status": "FILLED", "pairID": "5b6ac5117b445753ee755fb8", "pairName": "HPC/AUT", "orderBook": null, "createdAt": "2018-08-08T15:55:45.062141954+05:30", "updatedAt": "2018-08-08T15:55:45.06214208+05:30" }`)
	expectedOrderJSON2 := []byte(`{ "id": "5b6ac5297b4457546d64379e", "buyToken": "0x1888a8db0b7db59413ce07150b3373972bf818d3", "sellToken": "0x2034842261b82651885751fc293bba7ba5398156", "baseToken": "0x2034842261b82651885751fc293bba7ba5398156", "quoteToken": "0x1888a8db0b7db59413ce07150b3373972bf818d3", "buyAmount": 6000000000, "sellAmount": 13800000000, "nonce": 0, "userAddress": "0xefD7eB287CeeFCE8256Dd46e25F398acEA7C4b63", "hash": "0xa2d800b77828cb52c83106ca392e465bc0af0d7c319f6956328f739080c19622", "side": "SELL", "amount": 6000000000, "price": 229999999, "filledAmount": 4000000000, "fee": 0, "makeFee": 0, "takeFee": 0, "exchangeAddress": "", "status": "PARTIAL_FILLED", "pairID": "5b6ac5117b445753ee755fb8", "pairName": "HPC/AUT", "orderBook": null, "createdAt": "2018-08-08T15:55:45.062141954+05:30", "updatedAt": "2018-08-08T15:55:45.06214208+05:30" }`)
	json.Unmarshal(orderJSON2, &sampleOrder2)

	reoverOrders := []*FillOrder{
		&FillOrder{
			Amount: 1000000000,
			Order:  sampleOrder,
		}, &FillOrder{
			Amount: 6000000000,
			Order:  sampleOrder1,
		}, &FillOrder{
			Amount: 2000000000,
			Order:  sampleOrder2,
		},
	}

	// cancel sampleOrder1
	err := e.RecoverOrders(reoverOrders)

	if err != nil {
		t.Errorf("Error while recovering Orders: %s", err)
	}

	ssKey, listKey := sampleOrder.GetOBKeys()
	_, list1Key := sampleOrder1.GetOBKeys()
	_, list2Key := sampleOrder2.GetOBKeys()

	rs, err := s.SortedSet(ssKey)
	if err != nil {
		t.Error(err)
	} else {
		for k, v := range rs {
			assert.Equalf(t, utils.UintToPaddedString(sampleOrder.Price), k, "Expected sorted set value: %v got: %v", utils.UintToPaddedString(sampleOrder1.Price), k)
			assert.Equalf(t, 0.0, v, "Expected sorted set value: %v got: %v", 0, v)
		}
	}

	expectedMap := map[string]float64{
		sampleOrder.Hash.Hex():  float64(sampleOrder.CreatedAt.Unix()),
		sampleOrder1.Hash.Hex(): float64(sampleOrder1.CreatedAt.Unix()),
		sampleOrder2.Hash.Hex(): float64(sampleOrder2.CreatedAt.Unix()),
	}
	rs, err = s.SortedSet(listKey)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, expectedMap, rs)

	rse, err := s.Get(listKey + "::" + sampleOrder.Hash.Hex())
	if err != nil {
		t.Error(err)
	}
	assert.JSONEqf(t, string(expectedOrderJSON), rse, "Expected value for key: %v, was: %s, but got: %v", listKey+"::"+sampleOrder.Hash.Hex(), expectedOrderJSON, rse)

	rse, err = s.Get(list1Key + "::" + sampleOrder1.Hash.Hex())
	if err != nil {
		t.Error(err)
	}
	assert.JSONEqf(t, string(expectedOrderJSON1), rse, "Expected value for key: %v, was: %s, but got: %v", list1Key+"::"+sampleOrder1.Hash.Hex(), expectedOrderJSON1, rse)

	rse, err = s.Get(list2Key + "::" + sampleOrder2.Hash.Hex())
	if err != nil {
		t.Error(err)
	}
	assert.JSONEqf(t, string(expectedOrderJSON2), rse, "Expected value for key: %v, was: %s, but got: %v", list2Key+"::"+sampleOrder2.Hash.Hex(), expectedOrderJSON2, rse)

	rse, err = s.Get(ssKey + "::book::" + utils.UintToPaddedString(sampleOrder.Price))
	if err != nil {
		t.Error(err)
	}
	assert.Equalf(t, strconv.FormatInt(sampleOrder.Amount+sampleOrder1.Amount+sampleOrder2.Amount-sampleOrder2.FilledAmount, 10), rse, "Expected value for key: %v, was: %v, but got: %v", ssKey, sampleOrder.Amount+sampleOrder1.Amount+sampleOrder2.Amount, rse)

}

func TestSellOrder(t *testing.T) {
	e, s := getResource()
	defer s.Close()

	sampleOrder := &types.Order{}
	orderJSON := []byte(`{ "id": "5b6ac5297b4457546d64379d", "buyToken": "0x1888a8db0b7db59413ce07150b3373972bf818d3", "sellToken": "0x2034842261b82651885751fc293bba7ba5398156", "baseToken": "0x2034842261b82651885751fc293bba7ba5398156", "quoteToken": "0x1888a8db0b7db59413ce07150b3373972bf818d3", "buyAmount": 6000000000, "sellAmount": 13800000000, "nonce": 0, "userAddress": "0xefD7eB287CeeFCE8256Dd46e25F398acEA7C4b63", "hash": "0xa2d800b77828cb52c83106ca392e465bc0af0d7c319f6956328f739080c19621", "side": "SELL", "amount": 6000000000, "price": 229999999, "filledAmount": 0, "fee": 0, "makeFee": 0, "takeFee": 0, "exchangeAddress": "", "status": "NEW", "pairID": "5b6ac5117b445753ee755fb8", "pairName": "HPC/AUT", "orderBook": null, "createdAt": "2018-08-08T15:55:45.062141954+05:30", "updatedAt": "2018-08-08T15:55:45.06214208+05:30" }`)
	json.Unmarshal(orderJSON, &sampleOrder)
	expectedResponse := &Response{
		Order:          sampleOrder,
		RemainingOrder: &types.Order{},
		Trades:         make([]*types.Trade, 0),
		FillStatus:     NOMATCH,
		MatchingOrders: make([]*FillOrder, 0),
	}
	expectedResponse.Order.Status = types.OPEN
	response, err := e.sellOrder(sampleOrder)
	if err != nil {
		t.Errorf("Error in sellOrder: %s", err)
	}
	assert.Equal(t, expectedResponse, response)
}

func TestBuyOrder(t *testing.T) {
	e, s := getResource()
	defer s.Close()

	sampleOrder := &types.Order{}
	orderJSON := []byte(`{ "id": "5b6ac5297b4457546d64379d", "sellToken": "0x1888a8db0b7db59413ce07150b3373972bf818d3", "buyToken": "0x2034842261b82651885751fc293bba7ba5398156", "baseToken": "0x2034842261b82651885751fc293bba7ba5398156", "quoteToken": "0x1888a8db0b7db59413ce07150b3373972bf818d3", "buyAmount": 6000000000, "sellAmount": 13800000000, "nonce": 0, "userAddress": "0xefD7eB287CeeFCE8256Dd46e25F398acEA7C4b63", "hash": "0xa2d800b77828cb52c83106ca392e465bc0af0d7c319f6956328f739080c19621", "side": "BUY", "amount": 6000000000, "price": 229999999, "filledAmount": 0, "fee": 0, "makeFee": 0, "takeFee": 0, "exchangeAddress": "", "status": "NEW", "pairID": "5b6ac5117b445753ee755fb8", "pairName": "HPC/AUT", "orderBook": null, "createdAt": "2018-08-08T15:55:45.062141954+05:30", "updatedAt": "2018-08-08T15:55:45.06214208+05:30" }`)
	json.Unmarshal(orderJSON, &sampleOrder)
	expectedResponse := &Response{
		Order:          sampleOrder,
		RemainingOrder: &types.Order{},
		Trades:         make([]*types.Trade, 0),
		FillStatus:     NOMATCH,
		MatchingOrders: make([]*FillOrder, 0),
	}
	expectedResponse.Order.Status = types.OPEN
	response, err := e.buyOrder(sampleOrder)
	if err != nil {
		t.Errorf("Error in sellOrder: %s", err)
	}
	assert.Equal(t, expectedResponse, response)
}

/* Can not be tested yet, as SORT command is not supported by miniredis */
// TODO: Find alternative or a way to test SORT function of redis
//func TestFillOrder(t *testing.T) {
//	e, s := getResource()
//	defer s.Close()
//
//	// New Buy Order
//	sampleOrder := &types.Order{}
//	orderJSON := []byte(`{ "id": "5b6ac5297b4457546d64379d", "sellToken": "0x1888a8db0b7db59413ce07150b3373972bf818d3", "buyToken": "0x2034842261b82651885751fc293bba7ba5398156", "baseToken": "0x2034842261b82651885751fc293bba7ba5398156", "quoteToken": "0x1888a8db0b7db59413ce07150b3373972bf818d3", "sellAmount": 6000000000, "buyAmount": 13800000000, "nonce": 0, "userAddress": "0xefD7eB287CeeFCE8256Dd46e25F398acEA7C4b63", "hash": "0xa2d800b77828cb52c83106ca392e465bc0af0d7c319f6956328f739080c19621", "side": "BUY", "amount": 6000000000, "price": 229999999, "filledAmount": 0, "fee": 0, "makeFee": 0, "takeFee": 0, "exchangeAddress": "", "status": "NEW", "pairID": "5b6ac5117b445753ee755fb8", "pairName": "HPC/AUT", "orderBook": null, "createdAt": "2018-08-08T15:55:45.062141954+05:30", "updatedAt": "2018-08-08T15:55:45.06214208+05:30" }`)
//	json.Unmarshal(orderJSON, &sampleOrder)
//	expectedResponse := &Response{
//		Order:          sampleOrder,
//		RemainingOrder: &types.Order{},
//		Trades:         make([]*types.Trade, 0),
//		FillStatus:     NOMATCH,
//		MatchingOrders: make([]*FillOrder, 0),
//	}
//	expectedResponse.Order.Status = types.OPEN
//	response, err := e.buyOrder(sampleOrder)
//	if err != nil {
//		t.Errorf("Error in sellOrder: %s", err)
//	}
//	assert.Equal(t, expectedResponse, response)
//
//	// New Sell Order
//	sampleOrder1 := &types.Order{}
//	orderJSON = []byte(`{ "id": "5b6ac5297b4457546d64379e", "buyToken": "0x1888a8db0b7db59413ce07150b3373972bf818d3", "sellToken": "0x2034842261b82651885751fc293bba7ba5398156", "baseToken": "0x2034842261b82651885751fc293bba7ba5398156", "quoteToken": "0x1888a8db0b7db59413ce07150b3373972bf818d3", "buyAmount": 6000000000, "sellAmount": 13800000000, "nonce": 0, "userAddress": "0xefD7eB287CeeFCE8256Dd46e25F398acEA7C4b63", "hash": "0xa2d800b77828cb52c83106ca392e465bc0af0d7c319f6956328f739080c19622", "side": "SELL", "amount": 6000000000, "price": 229999999, "filledAmount": 0, "fee": 0, "makeFee": 0, "takeFee": 0, "exchangeAddress": "", "status": "NEW", "pairID": "5b6ac5117b445753ee755fb8", "pairName": "HPC/AUT", "orderBook": null, "createdAt": "2018-08-08T15:55:45.062141954+05:30", "updatedAt": "2018-08-08T15:55:45.06214208+05:30" }`)
//	json.Unmarshal(orderJSON, &sampleOrder1)
//	expectedResponse = &Response{
//		Order:          sampleOrder1,
//		RemainingOrder: &types.Order{},
//		Trades:         make([]*types.Trade, 0),
//		FillStatus:     FULL,
//		MatchingOrders: []*FillOrder{&FillOrder{Amount:6000000000,Order:sampleOrder}},
//	}
//	expectedResponse.Order.Status = types.FILLED
//	expectedResponse.Order.FilledAmount = expectedResponse.Order.Amount
//
//	response, err = e.sellOrder(sampleOrder1)
//	if err != nil {
//		t.Errorf("Error in sellOrder: %s", err)
//	}
//	assert.Equal(t, expectedResponse, response)
//}
