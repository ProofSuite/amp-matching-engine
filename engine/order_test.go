package engine

import (
	"log"
	"math/big"
	"strconv"
	"testing"
	"time"

	"encoding/json"

	"github.com/Proofsuite/amp-matching-engine/types"
	"github.com/Proofsuite/amp-matching-engine/utils"
	"github.com/Proofsuite/amp-matching-engine/utils/math"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"gopkg.in/mgo.v2/bson"
)

func TestAddOrder(t *testing.T) {
	e := getResource()
	defer flushData(e.redisConn)

	o1 := &types.Order{
		ID:              bson.ObjectIdHex("537f700b537461b70c5f0000"),
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
		PairID:          bson.ObjectIdHex("537f700b537461b70c5f0000"),
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
		Hash:      common.HexToHash("0xb9070a2d333403c255ce71ddf6e795053599b2e885321de40353832b96d8880a"),
		CreatedAt: time.Unix(1405544146, 0),
		UpdatedAt: time.Unix(1405544146, 0),
	}

	bytes, _ := o1.MarshalJSON()

	e.addOrder(o1)
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

	rse, err := getValue(e.redisConn, listKey+"::"+o1.Hash.Hex())
	if err != nil {
		t.Error(err)
	}

	assert.JSONEqf(t, string(bytes), rse, "Expected value for key: %v, was: %s, but got: %v", ssKey, bytes, rse)
	rse, err = getValue(e.redisConn, ssKey+"::book::"+utils.UintToPaddedString(o1.PricePoint.Int64()))
	if err != nil {
		t.Error(err)
	}

	// Currently converting amount to int64. In the future, we need to use strings instead of int64
	assert.Equalf(t, strconv.FormatInt(o1.Amount.Int64(), 10), rse, "Expected value for key: %v, was: %v, but got: %v", ssKey, o1.Amount, rse)

	o2 := &types.Order{
		ID:              bson.ObjectIdHex("537f700b537461b70c5f0000"),
		UserAddress:     common.HexToAddress("0x7a9f3cd060ab180f36c17fe6bdf9974f577d77aa"),
		ExchangeAddress: common.HexToAddress("0xae55690d4b079460e6ac28aaa58c9ec7b73a7485"),
		BuyToken:        common.HexToAddress("0xe41d2489571d322189246dafa5ebde1f4699f498"),
		SellToken:       common.HexToAddress("0x12459c951127e0c374ff9105dda097662a027093"),
		BaseToken:       common.HexToAddress("0xe41d2489571d322189246dafa5ebde1f4699f498"),
		QuoteToken:      common.HexToAddress("0x12459c951127e0c374ff9105dda097662a027093"),
		BuyAmount:       big.NewInt(1000),
		SellAmount:      big.NewInt(100),
		Price:           big.NewInt(6000000000),
		PricePoint:      big.NewInt(6000000000),
		Amount:          big.NewInt(229999999),
		FilledAmount:    big.NewInt(1000000000),
		Status:          "PARTIAL_FILLED",
		Side:            "BUY",
		PairID:          bson.ObjectIdHex("537f700b537461b70c5f0000"),
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
		Hash:      common.HexToHash("0xb9070a2d333403c255ce71ddf6e795053599b2e885321de40353832b96d8880a"),
		CreatedAt: time.Unix(1405544146, 0),
		UpdatedAt: time.Unix(1405544146, 0),
	}

	bytes, _ = o2.MarshalJSON()
	e.addOrder(o2)
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
	rse, err = getValue(e.redisConn, listKey+"::"+o2.Hash.Hex())
	if err != nil {
		t.Error(err)
	}
}

func TestCancelOrder(t *testing.T) {
	e := getResource()
	defer flushData(e.redisConn)

	// Add 1st order (OPEN order)
	o1 := &types.Order{
		ID:              bson.ObjectIdHex("537f700b537461b70c5f0000"),
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
		PairID:          bson.ObjectIdHex("537f700b537461b70c5f0000"),
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
		Hash:      common.HexToHash("0xb9070a2d333403c255ce71ddf6e795053599b2e885321de40353832b96d8880a"),
		CreatedAt: time.Unix(1405544146, 0),
		UpdatedAt: time.Unix(1405544146, 0),
	}

	e.addOrder(o1)

	// Add 2nd order (Partial filled order)
	o2 := &types.Order{
		ID:              bson.ObjectIdHex("537f700b537461b70c5f0000"),
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
		Status:          "PARTIAL_FILLED",
		Side:            "SELL",
		PairID:          bson.ObjectIdHex("537f700b537461b70c5f0000"),
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
		Hash:      common.HexToHash("0xb9070a2d333403c255ce71ddf6e795053599b2e885321de40353832b96d8880b"),
		CreatedAt: time.Unix(1405544146, 0),
		UpdatedAt: time.Unix(1405544146, 0),
	}

	e.addOrder(o2)

	expectedResponse := &Response{
		Order:          o2,
		Trades:         make([]*types.Trade, 0),
		RemainingOrder: &types.Order{},
		FillStatus:     CANCELLED,
		MatchingOrders: make([]*FillOrder, 0),
	}
	expectedResponse.Order.Status = "CANCELLED"

	// cancel o2
	response, err := e.CancelOrder(o2)

	if err != nil {
		t.Errorf("Error while cancelling SampleOrder1: %s", err)
	}

	assert.Equalf(t, expectedResponse, response, "Expected Response: %+v got: %+v", expectedResponse, response)

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
	assert.Equalf(t, strconv.FormatInt(o1.Amount.Int64(), 10), rse, "Expected value for key: %v, was: %v, but got: %v", ss1Key, o2.Amount.Int64()-o2.FilledAmount.Int64(), rse)

	expectedResponse.Order = o1
	expectedResponse.Order.Status = "CANCELLED"

	// cancel o1
	response, err = e.CancelOrder(o1)
	ssKey, listKey := o1.GetOBKeys()

	if err != nil {
		t.Errorf("Error while cancelling SampleOrder: %s", err)
	}

	assert.Equalf(t, expectedResponse, response, "Expected Response from cancelling sampleorder : %+v got: %+v", expectedResponse, response)

	// All keys must have been deleted
	if exists(e.redisConn, ssKey) {
		t.Errorf("Key : %v expected to be deleted but key exists", ssKey)
	}

	if exists(e.redisConn, listKey) {
		t.Errorf("Key : %v expected to be deleted but key exists", ssKey)
	}

	if exists(e.redisConn, listKey+"::"+o1.Hash.Hex()) {
		t.Errorf("Key : %v expected to be deleted but key exists", ssKey)
	}

	if exists(e.redisConn, ssKey+"::book::"+utils.UintToPaddedString(o1.PricePoint.Int64())) {
		t.Errorf("Key : %v expected to be deleted but key exists", ssKey)
	}
}

func TestRecoverOrders(t *testing.T) {
	e := getResource()
	defer flushData(e.redisConn)

	o1 := &types.Order{
		ID:              bson.ObjectIdHex("537f700b537461b70c5f0000"),
		UserAddress:     common.HexToAddress("0x7a9f3cd060ab180f36c17fe6bdf9974f577d77aa"),
		ExchangeAddress: common.HexToAddress("0xae55690d4b079460e6ac28aaa58c9ec7b73a7485"),
		BuyToken:        common.HexToAddress("0x1888a8db0b7db59413ce07150b3373972bf818d3"),
		SellToken:       common.HexToAddress("0x2034842261b82651885751fc293bba7ba5398156"),
		BaseToken:       common.HexToAddress("0x2034842261b82651885751fc293bba7ba5398156"),
		QuoteToken:      common.HexToAddress("0x1888a8db0b7db59413ce07150b3373972bf818d3"),
		BuyAmount:       big.NewInt(6000000000),
		SellAmount:      big.NewInt(13800000000),
		Price:           big.NewInt(229999999),
		PricePoint:      big.NewInt(229999999),
		Amount:          big.NewInt(6000000000),
		FilledAmount:    big.NewInt(1000000000),
		Status:          "PARTIAL_FILLED",
		Side:            "SELL",
		PairID:          bson.ObjectIdHex("537f700b537461b70c5f0000"),
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
		Hash:      common.HexToHash("0xa2d800b77828cb52c83106ca392e465bc0af0d7c319f6956328f739080c19620"),
		CreatedAt: time.Unix(1405544146, 0),
		UpdatedAt: time.Unix(1405544146, 0),
	}

	expectedOrder1 := *o1
	expectedOrder1.Status = "OPEN"
	expectedOrder1.FilledAmount = big.NewInt(0)
	expectedOrder1Json, _ := expectedOrder1.MarshalJSON()

	e.addOrder(o1)

	o2 := &types.Order{
		ID:              bson.ObjectIdHex("537f700b537461b70c5f0000"),
		UserAddress:     common.HexToAddress("0x7a9f3cd060ab180f36c17fe6bdf9974f577d77aa"),
		ExchangeAddress: common.HexToAddress("0xae55690d4b079460e6ac28aaa58c9ec7b73a7485"),
		BuyToken:        common.HexToAddress("0x1888a8db0b7db59413ce07150b3373972bf818d3"),
		SellToken:       common.HexToAddress("0x2034842261b82651885751fc293bba7ba5398156"),
		BaseToken:       common.HexToAddress("0x2034842261b82651885751fc293bba7ba5398156"),
		QuoteToken:      common.HexToAddress("0x1888a8db0b7db59413ce07150b3373972bf818d3"),
		BuyAmount:       big.NewInt(6000000000),
		SellAmount:      big.NewInt(13800000000),
		Price:           big.NewInt(229999999),
		PricePoint:      big.NewInt(229999999),
		Amount:          big.NewInt(6000000000),
		FilledAmount:    big.NewInt(6000000000),
		Status:          "FILLED",
		Side:            "SELL",
		PairID:          bson.ObjectIdHex("537f700b537461b70c5f0000"),
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
		Hash:      common.HexToHash("0xa2d800b77828cb52c83106ca392e465bc0af0d7c319f6956328f739080c19621"),
		CreatedAt: time.Unix(1405544146, 0),
		UpdatedAt: time.Unix(1405544146, 0),
	}

	expectedOrder2 := *o2
	expectedOrder2.Status = "OPEN"
	expectedOrder2.FilledAmount = big.NewInt(0)
	expectedOrder2Json, _ := expectedOrder2.MarshalJSON()

	e.addOrder(o2)

	o3 := &types.Order{
		ID:              bson.ObjectIdHex("537f700b537461b70c5f0000"),
		UserAddress:     common.HexToAddress("0x7a9f3cd060ab180f36c17fe6bdf9974f577d77aa"),
		ExchangeAddress: common.HexToAddress("0xae55690d4b079460e6ac28aaa58c9ec7b73a7485"),
		BuyToken:        common.HexToAddress("0x1888a8db0b7db59413ce07150b3373972bf818d3"),
		SellToken:       common.HexToAddress("0x2034842261b82651885751fc293bba7ba5398156"),
		BaseToken:       common.HexToAddress("0x2034842261b82651885751fc293bba7ba5398156"),
		QuoteToken:      common.HexToAddress("0x1888a8db0b7db59413ce07150b3373972bf818d3"),
		BuyAmount:       big.NewInt(6000000000),
		SellAmount:      big.NewInt(13800000000),
		Price:           big.NewInt(229999999),
		PricePoint:      big.NewInt(229999999),
		Amount:          big.NewInt(6000000000),
		FilledAmount:    big.NewInt(6000000000),
		Status:          "FILLED",
		Side:            "SELL",
		PairID:          bson.ObjectIdHex("537f700b537461b70c5f0000"),
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
		Hash:      common.HexToHash("0xa2d800b77828cb52c83106ca392e465bc0af0d7c319f6956328f739080c19622"),
		CreatedAt: time.Unix(1405544146, 0),
		UpdatedAt: time.Unix(1405544146, 0),
	}

	expectedOrder3 := *o3
	expectedOrder3.Status = "PARTIAL_FILLED"
	expectedOrder3.FilledAmount = big.NewInt(4000000000)
	expectedOrder3Json, _ := expectedOrder3.MarshalJSON()

	e.addOrder(o3)

	recoverOrders := []*FillOrder{
		&FillOrder{
			Amount: big.NewInt(1000000000),
			Order:  o1,
		}, &FillOrder{
			Amount: big.NewInt(6000000000),
			Order:  o2,
		}, &FillOrder{
			Amount: big.NewInt(2000000000),
			Order:  o3,
		},
	}

	// cancel o2
	err := e.RecoverOrders(recoverOrders)

	if err != nil {
		t.Errorf("Error while recovering Orders: %s", err)
	}

	ssKey, listKey := o1.GetOBKeys()
	_, list1Key := o2.GetOBKeys()
	_, list2Key := o3.GetOBKeys()

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

	rse, err := getValue(e.redisConn, listKey+"::"+o1.Hash.Hex())
	if err != nil {
		t.Error(err)
	}

	assert.JSONEqf(t, string(expectedOrder1Json), rse, "Expected value for key: %v, was: %s, but got: %v", listKey+"::"+o1.Hash.Hex(), expectedOrder1Json, rse)

	rse, err = getValue(e.redisConn, list1Key+"::"+o2.Hash.Hex())
	if err != nil {
		t.Error(err)
	}

	assert.JSONEqf(t, string(expectedOrder2Json), rse, "Expected value for key: %v, was: %s, but got: %v", list1Key+"::"+o2.Hash.Hex(), expectedOrder2Json, rse)

	rse, err = getValue(e.redisConn, list2Key+"::"+o3.Hash.Hex())
	if err != nil {
		t.Error(err)
	}

	assert.JSONEqf(t, string(expectedOrder3Json), rse, "Expected value for key: %v, was: %s, but got: %v", list2Key+"::"+o3.Hash.Hex(), expectedOrder3Json, rse)

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

	o1 := &types.Order{
		ID:              bson.ObjectIdHex("537f700b537461b70c5f0000"),
		UserAddress:     common.HexToAddress("0x7a9f3cd060ab180f36c17fe6bdf9974f577d77aa"),
		ExchangeAddress: common.HexToAddress("0xae55690d4b079460e6ac28aaa58c9ec7b73a7485"),
		BuyToken:        common.HexToAddress("0x1888a8db0b7db59413ce07150b3373972bf818d3"),
		SellToken:       common.HexToAddress("0x2034842261b82651885751fc293bba7ba5398156"),
		BaseToken:       common.HexToAddress("0x2034842261b82651885751fc293bba7ba5398156"),
		QuoteToken:      common.HexToAddress("0x1888a8db0b7db59413ce07150b3373972bf818d3"),
		BuyAmount:       big.NewInt(6000000000),
		SellAmount:      big.NewInt(13800000000),
		Price:           big.NewInt(229999999),
		PricePoint:      big.NewInt(229999999),
		Amount:          big.NewInt(6000000000),
		FilledAmount:    big.NewInt(0),
		Status:          "NEW",
		Side:            "SELL",
		PairID:          bson.ObjectIdHex("537f700b537461b70c5f0000"),
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
		Hash:      common.HexToHash("0xa2d800b77828cb52c83106ca392e465bc0af0d7c319f6956328f739080c19621"),
		CreatedAt: time.Unix(1405544146, 0),
		UpdatedAt: time.Unix(1405544146, 0),
	}

	expectedResponse := &Response{
		Order:          o1,
		Trades:         make([]*types.Trade, 0),
		FillStatus:     NOMATCH,
		MatchingOrders: make([]*FillOrder, 0),
	}

	expectedResponse.Order.Status = "OPEN"

	response, err := e.sellOrder(o1)
	if err != nil {
		t.Errorf("Error in sellOrder: %s", err)
	}
	assert.Equal(t, expectedResponse, response)
}

func TestBuyOrder(t *testing.T) {
	e := getResource()
	defer flushData(e.redisConn)

	o1 := &types.Order{
		ID:              bson.ObjectIdHex("537f700b537461b70c5f0000"),
		UserAddress:     common.HexToAddress("0x7a9f3cd060ab180f36c17fe6bdf9974f577d77aa"),
		ExchangeAddress: common.HexToAddress("0xae55690d4b079460e6ac28aaa58c9ec7b73a7485"),
		BuyToken:        common.HexToAddress("0x1888a8db0b7db59413ce07150b3373972bf818d3"),
		SellToken:       common.HexToAddress("0x2034842261b82651885751fc293bba7ba5398156"),
		BaseToken:       common.HexToAddress("0x2034842261b82651885751fc293bba7ba5398156"),
		QuoteToken:      common.HexToAddress("0x1888a8db0b7db59413ce07150b3373972bf818d3"),
		BuyAmount:       big.NewInt(6000000000),
		SellAmount:      big.NewInt(13800000000),
		Price:           big.NewInt(229999999),
		PricePoint:      big.NewInt(229999999),
		Amount:          big.NewInt(6000000000),
		FilledAmount:    big.NewInt(0),
		Status:          "NEW",
		Side:            "BUY",
		PairID:          bson.ObjectIdHex("537f700b537461b70c5f0000"),
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
		Hash:      common.HexToHash("0xa2d800b77828cb52c83106ca392e465bc0af0d7c319f6956328f739080c19621"),
		CreatedAt: time.Unix(1405544146, 0),
		UpdatedAt: time.Unix(1405544146, 0),
	}

	expectedResponse := &Response{
		Order:          o1,
		Trades:         make([]*types.Trade, 0),
		FillStatus:     NOMATCH,
		MatchingOrders: make([]*FillOrder, 0),
	}

	expectedResponse.Order.Status = "OPEN"

	response, err := e.buyOrder(o1)
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
	sellOrder := &types.Order{
		ID:              bson.ObjectIdHex("537f700b537461b70c5f0000"),
		UserAddress:     common.HexToAddress("0x7a9f3cd060ab180f36c17fe6bdf9974f577d77aa"),
		ExchangeAddress: common.HexToAddress("0xae55690d4b079460e6ac28aaa58c9ec7b73a7485"),
		BaseToken:       common.HexToAddress("0x2034842261b82651885751fc293bba7ba5398156"),
		QuoteToken:      common.HexToAddress("0x1888a8db0b7db59413ce07150b3373972bf818d3"),
		BuyAmount:       big.NewInt(6000000000),
		SellAmount:      big.NewInt(13800000000),
		Price:           big.NewInt(229999999),
		PricePoint:      big.NewInt(229999999),
		Amount:          big.NewInt(6000000000),
		FilledAmount:    big.NewInt(0),
		Status:          "NEW",
		Side:            "SELL",
		PairID:          bson.ObjectIdHex("537f700b537461b70c5f0000"),
		PairName:        "ZRX/WETH",
		Expires:         big.NewInt(0),
		MakeFee:         big.NewInt(0),
		Nonce:           big.NewInt(0),
		TakeFee:         big.NewInt(0),
		Signature: &types.Signature{
			V: 28,
			R: common.HexToHash("0x10b30eb0072a4f0a38b6fca0b731cba15eb2e1702845d97c1230b53a839bcb85"),
			S: common.HexToHash("0x6d9ad89548c9e3ce4c97825d027291477f2c44a8caef792095f2cabc978493ff"),
		},
		Hash:      common.HexToHash("0xa2d800b77828cb52c83106ca392e465bc0af0d7c319f6956328f739080c19621"),
		CreatedAt: time.Unix(1405544146, 0),
		UpdatedAt: time.Unix(1405544146, 0),
	}

	buyOrder := &types.Order{
		ID:              bson.ObjectIdHex("537f700b537461b70c5f0000"),
		UserAddress:     common.HexToAddress("0x7a9f3cd060ab180f36c17fe6bdf9974f577d77aa"),
		ExchangeAddress: common.HexToAddress("0xae55690d4b079460e6ac28aaa58c9ec7b73a7485"),
		BaseToken:       common.HexToAddress("0x2034842261b82651885751fc293bba7ba5398156"),
		QuoteToken:      common.HexToAddress("0x1888a8db0b7db59413ce07150b3373972bf818d3"),
		BuyAmount:       big.NewInt(6000000000),
		SellAmount:      big.NewInt(13800000000),
		Price:           big.NewInt(229999999),
		PricePoint:      big.NewInt(229999999),
		Amount:          big.NewInt(6000000000),
		FilledAmount:    big.NewInt(0),
		Status:          "NEW",
		Side:            "BUY",
		PairID:          bson.ObjectIdHex("537f700b537461b70c5f0000"),
		PairName:        "ZRX/WETH",
		Expires:         big.NewInt(0),
		MakeFee:         big.NewInt(0),
		Nonce:           big.NewInt(0),
		TakeFee:         big.NewInt(0),
		Signature: &types.Signature{
			V: 28,
			R: common.HexToHash("0x10b30eb0072a4f0a38b6fca0b731cba15eb2e1702845d97c1230b53a839bcb85"),
			S: common.HexToHash("0x6d9ad89548c9e3ce4c97825d027291477f2c44a8caef792095f2cabc978493ff"),
		},
		Hash:      common.HexToHash("0xa2d800b77828cb52c83106ca392e465bc0af0d7c319f6956328f739080c19621"),
		CreatedAt: time.Unix(1405544146, 0),
		UpdatedAt: time.Unix(1405544146, 0),
	}

	// Test Case1: Send sellOrder first
	expectedResponseBuyOrder := *buyOrder
	expectedResponseBuyOrder.FilledAmount = expectedResponseBuyOrder.Amount
	expectedResponseBuyOrder.Status = "FILLED"
	expectedResponseSellOrder := *sellOrder
	expectedResponseSellOrder.Status = "OPEN"

	expectedResponse := &Response{
		FillStatus:     NOMATCH,
		Order:          &expectedResponseSellOrder,
		Trades:         make([]*types.Trade, 0),
		MatchingOrders: make([]*FillOrder, 0),
	}

	expBytes, _ := json.Marshal(expectedResponse)
	resp, err := e.sellOrder(sellOrder)
	if err != nil {
		t.Errorf("Error in sellOrder: %s", err)
	}

	resBytes, _ := json.Marshal(resp)
	assert.JSONEq(t, string(expBytes), string(resBytes))

	expectedResponseSellOrder.FilledAmount = expectedResponseSellOrder.Amount
	expectedResponseSellOrder.Status = "FILLED"
	trade := &types.Trade{
		Amount:       expectedResponseBuyOrder.Amount,
		Price:        buyOrder.PricePoint,
		PricePoint:   buyOrder.PricePoint,
		BaseToken:    buyOrder.BaseToken,
		QuoteToken:   buyOrder.QuoteToken,
		OrderHash:    sellOrder.Hash,
		Side:         buyOrder.Side,
		Taker:        buyOrder.UserAddress,
		PairName:     buyOrder.PairName,
		Maker:        sellOrder.UserAddress,
		TakerOrderID: buyOrder.ID,
		MakerOrderID: sellOrder.ID,
		TradeNonce:   big.NewInt(0),
	}
	trade.Hash = trade.ComputeHash()

	expectedResponse = &Response{
		Order:          &expectedResponseBuyOrder,
		Trades:         []*types.Trade{trade},
		FillStatus:     FULL,
		MatchingOrders: []*FillOrder{{buyOrder.Amount, &expectedResponseSellOrder}},
	}

	expBytes, _ = json.Marshal(expectedResponse)
	resp, err = e.buyOrder(buyOrder)
	if err != nil {
		t.Errorf("Error in sellOrder: %s", err)
	}
	resBytes, _ = json.Marshal(resp)
	assert.JSONEq(t, string(expBytes), string(resBytes))
}

// TestFillOrder1: Test Case2: Send buyOrder first
func TestFillOrder2(t *testing.T) {
	e := getResource()
	defer flushData(e.redisConn)
	sellOrder := &types.Order{
		ID:              bson.ObjectIdHex("537f700b537461b70c5f0000"),
		UserAddress:     common.HexToAddress("0x7a9f3cd060ab180f36c17fe6bdf9974f577d77aa"),
		ExchangeAddress: common.HexToAddress("0xae55690d4b079460e6ac28aaa58c9ec7b73a7485"),
		BaseToken:       common.HexToAddress("0x2034842261b82651885751fc293bba7ba5398156"),
		QuoteToken:      common.HexToAddress("0x1888a8db0b7db59413ce07150b3373972bf818d3"),
		BuyAmount:       big.NewInt(6000000000),
		SellAmount:      big.NewInt(13800000000),
		Price:           big.NewInt(229999999),
		PricePoint:      big.NewInt(229999999),
		Amount:          big.NewInt(6000000000),
		FilledAmount:    big.NewInt(0),
		Status:          "NEW",
		Side:            "SELL",
		PairID:          bson.ObjectIdHex("537f700b537461b70c5f0000"),
		PairName:        "ZRX/WETH",
		Expires:         big.NewInt(0),
		MakeFee:         big.NewInt(0),
		Nonce:           big.NewInt(0),
		TakeFee:         big.NewInt(0),
		Signature: &types.Signature{
			V: 28,
			R: common.HexToHash("0x10b30eb0072a4f0a38b6fca0b731cba15eb2e1702845d97c1230b53a839bcb85"),
			S: common.HexToHash("0x6d9ad89548c9e3ce4c97825d027291477f2c44a8caef792095f2cabc978493ff"),
		},
		Hash:      common.HexToHash("0xa2d800b77828cb52c83106ca392e465bc0af0d7c319f6956328f739080c19621"),
		CreatedAt: time.Unix(1405544146, 0),
		UpdatedAt: time.Unix(1405544146, 0),
	}

	buyOrder := &types.Order{
		ID:              bson.ObjectIdHex("537f700b537461b70c5f0000"),
		UserAddress:     common.HexToAddress("0x7a9f3cd060ab180f36c17fe6bdf9974f577d77aa"),
		ExchangeAddress: common.HexToAddress("0xae55690d4b079460e6ac28aaa58c9ec7b73a7485"),
		BaseToken:       common.HexToAddress("0x2034842261b82651885751fc293bba7ba5398156"),
		QuoteToken:      common.HexToAddress("0x1888a8db0b7db59413ce07150b3373972bf818d3"),
		BuyAmount:       big.NewInt(6000000000),
		SellAmount:      big.NewInt(13800000000),
		Price:           big.NewInt(229999999),
		PricePoint:      big.NewInt(229999999),
		Amount:          big.NewInt(6000000000),
		FilledAmount:    big.NewInt(0),
		Status:          "NEW",
		Side:            "BUY",
		PairID:          bson.ObjectIdHex("537f700b537461b70c5f0000"),
		PairName:        "ZRX/WETH",
		Expires:         big.NewInt(0),
		MakeFee:         big.NewInt(0),
		Nonce:           big.NewInt(0),
		TakeFee:         big.NewInt(0),
		Signature: &types.Signature{
			V: 28,
			R: common.HexToHash("0x10b30eb0072a4f0a38b6fca0b731cba15eb2e1702845d97c1230b53a839bcb85"),
			S: common.HexToHash("0x6d9ad89548c9e3ce4c97825d027291477f2c44a8caef792095f2cabc978493ff"),
		},
		Hash:      common.HexToHash("0xa2d800b77828cb52c83106ca392e465bc0af0d7c319f6956328f739080c19621"),
		CreatedAt: time.Unix(1405544146, 0),
		UpdatedAt: time.Unix(1405544146, 0),
	}

	expectedResponseBuyOrder := *buyOrder
	expectedResponseBuyOrder.Status = "OPEN"

	expectedResponse := &Response{
		Order:          &expectedResponseBuyOrder,
		Trades:         make([]*types.Trade, 0),
		FillStatus:     NOMATCH,
		MatchingOrders: make([]*FillOrder, 0),
	}

	erBytes, _ := json.Marshal(expectedResponse)
	response, err := e.buyOrder(buyOrder)
	if err != nil {
		t.Errorf("Error in sellOrder: %s", err)
	}

	resBytes, _ := json.Marshal(response)
	assert.JSONEq(t, string(erBytes), string(resBytes))

	expectedResponseSellOrder := *sellOrder
	expectedResponseSellOrder.FilledAmount = expectedResponseSellOrder.Amount
	expectedResponseSellOrder.Status = "FILLED"

	expectedResponseBuyOrder.FilledAmount = expectedResponseBuyOrder.Amount
	expectedResponseBuyOrder.Status = "FILLED"

	trade := &types.Trade{
		Amount:       expectedResponseSellOrder.Amount,
		Price:        sellOrder.Price,
		PricePoint:   sellOrder.PricePoint,
		BaseToken:    sellOrder.BaseToken,
		QuoteToken:   sellOrder.QuoteToken,
		OrderHash:    buyOrder.Hash,
		Side:         sellOrder.Side,
		Taker:        sellOrder.UserAddress,
		PairName:     sellOrder.PairName,
		Maker:        buyOrder.UserAddress,
		TakerOrderID: sellOrder.ID,
		MakerOrderID: buyOrder.ID,
		TradeNonce:   big.NewInt(0),
	}

	trade.Hash = trade.ComputeHash()

	expectedResponse = &Response{
		Order:          &expectedResponseSellOrder,
		Trades:         []*types.Trade{trade},
		FillStatus:     FULL,
		MatchingOrders: []*FillOrder{{buyOrder.Amount, &expectedResponseBuyOrder}},
	}

	bytes, _ := json.Marshal(expectedResponse)
	response, err = e.sellOrder(sellOrder)
	if err != nil {
		t.Errorf("Error in sellOrder: %s", err)
	}

	resBytes, _ = json.Marshal(response)
	assert.JSONEq(t, string(bytes), string(resBytes))
}

func TestMultiMatchOrder1(t *testing.T) {
	e := getResource()
	defer flushData(e.redisConn)
	sellOrder := types.Order{
		ID:              bson.ObjectIdHex("537f700b537461b70c5f0000"),
		UserAddress:     common.HexToAddress("0x7a9f3cd060ab180f36c17fe6bdf9974f577d77aa"),
		ExchangeAddress: common.HexToAddress("0xae55690d4b079460e6ac28aaa58c9ec7b73a7485"),
		BaseToken:       common.HexToAddress("0x2034842261b82651885751fc293bba7ba5398156"),
		QuoteToken:      common.HexToAddress("0x1888a8db0b7db59413ce07150b3373972bf818d3"),
		BuyAmount:       big.NewInt(6000000000),
		SellAmount:      big.NewInt(13800000000),
		Price:           big.NewInt(229999999),
		PricePoint:      big.NewInt(229999999),
		Amount:          big.NewInt(6000000000),
		FilledAmount:    big.NewInt(0),
		Status:          "NEW",
		Side:            "SELL",
		PairID:          bson.ObjectIdHex("537f700b537461b70c5f0000"),
		PairName:        "ZRX/WETH",
		Expires:         big.NewInt(0),
		MakeFee:         big.NewInt(0),
		Nonce:           big.NewInt(0),
		TakeFee:         big.NewInt(0),
		Signature: &types.Signature{
			V: 28,
			R: common.HexToHash("0x10b30eb0072a4f0a38b6fca0b731cba15eb2e1702845d97c1230b53a839bcb85"),
			S: common.HexToHash("0x6d9ad89548c9e3ce4c97825d027291477f2c44a8caef792095f2cabc978493ff"),
		},
		Hash:      common.HexToHash("0xa2d800b77828cb52c83106ca392e465bc0af0d7c319f6956328f739080c19621"),
		CreatedAt: time.Unix(1405544146, 0),
		UpdatedAt: time.Unix(1405544146, 0),
	}

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

	buyOrder := types.Order{
		ID:              bson.ObjectIdHex("537f700b537461b70c5f0000"),
		UserAddress:     common.HexToAddress("0x7a9f3cd060ab180f36c17fe6bdf9974f577d77aa"),
		ExchangeAddress: common.HexToAddress("0xae55690d4b079460e6ac28aaa58c9ec7b73a7485"),
		BaseToken:       common.HexToAddress("0x2034842261b82651885751fc293bba7ba5398156"),
		QuoteToken:      common.HexToAddress("0x1888a8db0b7db59413ce07150b3373972bf818d3"),
		BuyAmount:       big.NewInt(6000000000),
		SellAmount:      big.NewInt(13800000000),
		Price:           big.NewInt(229999999),
		PricePoint:      big.NewInt(229999999),
		Amount:          big.NewInt(6000000000),
		FilledAmount:    big.NewInt(0),
		Status:          "NEW",
		Side:            "BUY",
		PairID:          bson.ObjectIdHex("537f700b537461b70c5f0000"),
		PairName:        "ZRX/WETH",
		Expires:         big.NewInt(0),
		MakeFee:         big.NewInt(0),
		Nonce:           big.NewInt(0),
		TakeFee:         big.NewInt(0),
		Signature: &types.Signature{
			V: 28,
			R: common.HexToHash("0x10b30eb0072a4f0a38b6fca0b731cba15eb2e1702845d97c1230b53a839bcb85"),
			S: common.HexToHash("0x6d9ad89548c9e3ce4c97825d027291477f2c44a8caef792095f2cabc978493ff"),
		},
		Hash:      common.HexToHash("0xa2d800b77828cb52c83106ca392e465bc0af0d7c319f6956328f739080c19621"),
		CreatedAt: time.Unix(1405544146, 0),
		UpdatedAt: time.Unix(1405544146, 0),
	}

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

	trade := &types.Trade{
		Amount:       responseSO.Amount,
		Price:        buyOrder.PricePoint,
		PricePoint:   buyOrder.PricePoint,
		BaseToken:    buyOrder.BaseToken,
		QuoteToken:   buyOrder.QuoteToken,
		OrderHash:    sellOrder.Hash,
		Side:         buyOrder.Side,
		Taker:        buyOrder.UserAddress,
		PairName:     buyOrder.PairName,
		Maker:        sellOrder.UserAddress,
		TakerOrderID: buyOrder.ID,
		MakerOrderID: sellOrder.ID,
		TradeNonce:   big.NewInt(0),
	}

	trade.Hash = trade.ComputeHash()

	trade1 := &types.Trade{
		Amount:       responseSO1.Amount,
		Price:        buyOrder.PricePoint,
		PricePoint:   buyOrder.PricePoint,
		BaseToken:    buyOrder.BaseToken,
		QuoteToken:   buyOrder.QuoteToken,
		OrderHash:    sellOrder1.Hash,
		Side:         buyOrder.Side,
		Taker:        buyOrder.UserAddress,
		PairName:     buyOrder.PairName,
		Maker:        sellOrder.UserAddress,
		TakerOrderID: buyOrder.ID,
		MakerOrderID: sellOrder.ID,
		TradeNonce:   big.NewInt(0),
	}

	trade1.Hash = trade1.ComputeHash()

	trade2 := &types.Trade{
		Amount:       responseSO2.Amount,
		Price:        buyOrder.PricePoint,
		PricePoint:   buyOrder.PricePoint,
		BaseToken:    buyOrder.BaseToken,
		QuoteToken:   buyOrder.QuoteToken,
		OrderHash:    sellOrder2.Hash,
		Side:         buyOrder.Side,
		Taker:        buyOrder.UserAddress,
		PairName:     buyOrder.PairName,
		Maker:        sellOrder.UserAddress,
		TakerOrderID: buyOrder.ID,
		MakerOrderID: sellOrder.ID,
		TradeNonce:   big.NewInt(0),
	}

	trade2.Hash = trade2.ComputeHash()

	expectedResponse := &Response{
		Order:          &responseBO,
		Trades:         []*types.Trade{trade, trade1, trade2},
		FillStatus:     FULL,
		MatchingOrders: []*FillOrder{{responseSO.FilledAmount, &responseSO}, {responseSO1.FilledAmount, &responseSO1}, {responseSO2.FilledAmount, &responseSO2}},
	}

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
	buyOrder := types.Order{
		ID:              bson.ObjectIdHex("537f700b537461b70c5f0000"),
		UserAddress:     common.HexToAddress("0x7a9f3cd060ab180f36c17fe6bdf9974f577d77aa"),
		ExchangeAddress: common.HexToAddress("0xae55690d4b079460e6ac28aaa58c9ec7b73a7485"),
		BaseToken:       common.HexToAddress("0x2034842261b82651885751fc293bba7ba5398156"),
		QuoteToken:      common.HexToAddress("0x1888a8db0b7db59413ce07150b3373972bf818d3"),
		BuyAmount:       big.NewInt(6000000000),
		SellAmount:      big.NewInt(13800000000),
		Price:           big.NewInt(229999999),
		PricePoint:      big.NewInt(229999999),
		Amount:          big.NewInt(6000000000),
		FilledAmount:    big.NewInt(0),
		Status:          "NEW",
		Side:            "BUY",
		PairID:          bson.ObjectIdHex("537f700b537461b70c5f0000"),
		PairName:        "ZRX/WETH",
		Expires:         big.NewInt(0),
		MakeFee:         big.NewInt(0),
		Nonce:           big.NewInt(0),
		TakeFee:         big.NewInt(0),
		Signature: &types.Signature{
			V: 28,
			R: common.HexToHash("0x10b30eb0072a4f0a38b6fca0b731cba15eb2e1702845d97c1230b53a839bcb85"),
			S: common.HexToHash("0x6d9ad89548c9e3ce4c97825d027291477f2c44a8caef792095f2cabc978493ff"),
		},
		Hash:      common.HexToHash("0xa2d800b77828cb52c83106ca392e465bc0af0d7c319f6956328f739080c19621"),
		CreatedAt: time.Unix(1405544146, 0),
		UpdatedAt: time.Unix(1405544146, 0),
	}

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

	sellOrder := types.Order{
		ID:              bson.ObjectIdHex("537f700b537461b70c5f0000"),
		UserAddress:     common.HexToAddress("0x7a9f3cd060ab180f36c17fe6bdf9974f577d77aa"),
		ExchangeAddress: common.HexToAddress("0xae55690d4b079460e6ac28aaa58c9ec7b73a7485"),
		BaseToken:       common.HexToAddress("0x2034842261b82651885751fc293bba7ba5398156"),
		QuoteToken:      common.HexToAddress("0x1888a8db0b7db59413ce07150b3373972bf818d3"),
		BuyAmount:       big.NewInt(6000000000),
		SellAmount:      big.NewInt(13800000000),
		Price:           big.NewInt(229999999),
		PricePoint:      big.NewInt(229999999),
		Amount:          big.NewInt(6000000000),
		FilledAmount:    big.NewInt(0),
		Status:          "NEW",
		Side:            "SELL",
		PairID:          bson.ObjectIdHex("537f700b537461b70c5f0000"),
		PairName:        "ZRX/WETH",
		Expires:         big.NewInt(0),
		MakeFee:         big.NewInt(0),
		Nonce:           big.NewInt(0),
		TakeFee:         big.NewInt(0),
		Signature: &types.Signature{
			V: 28,
			R: common.HexToHash("0x10b30eb0072a4f0a38b6fca0b731cba15eb2e1702845d97c1230b53a839bcb85"),
			S: common.HexToHash("0x6d9ad89548c9e3ce4c97825d027291477f2c44a8caef792095f2cabc978493ff"),
		},
		Hash:      common.HexToHash("0xa2d800b77828cb52c83106ca392e465bc0af0d7c319f6956328f739080c19621"),
		CreatedAt: time.Unix(1405544146, 0),
		UpdatedAt: time.Unix(1405544146, 0),
	}

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

	trade := &types.Trade{
		Amount:       buyOrder.Amount,
		Price:        sellOrder.PricePoint,
		PricePoint:   sellOrder.PricePoint,
		BaseToken:    sellOrder.BaseToken,
		QuoteToken:   sellOrder.QuoteToken,
		OrderHash:    buyOrder.Hash,
		Side:         sellOrder.Side,
		Taker:        sellOrder.UserAddress,
		PairName:     sellOrder.PairName,
		Maker:        buyOrder.UserAddress,
		TakerOrderID: sellOrder.ID,
		MakerOrderID: buyOrder.ID,
		TradeNonce:   big.NewInt(0),
	}

	trade.Hash = trade.ComputeHash()

	trade1 := &types.Trade{
		Amount:       buyOrder1.Amount,
		Price:        sellOrder.PricePoint,
		PricePoint:   sellOrder.PricePoint,
		BaseToken:    sellOrder.BaseToken,
		QuoteToken:   sellOrder.QuoteToken,
		OrderHash:    buyOrder1.Hash,
		Side:         sellOrder.Side,
		Taker:        sellOrder.UserAddress,
		PairName:     sellOrder.PairName,
		Maker:        buyOrder.UserAddress,
		TakerOrderID: sellOrder.ID,
		MakerOrderID: buyOrder.ID,
		TradeNonce:   big.NewInt(0),
	}

	trade1.Hash = trade1.ComputeHash()

	trade2 := &types.Trade{
		Amount:       buyOrder2.Amount,
		Price:        sellOrder.PricePoint,
		PricePoint:   sellOrder.PricePoint,
		BaseToken:    sellOrder.BaseToken,
		QuoteToken:   sellOrder.QuoteToken,
		OrderHash:    buyOrder2.Hash,
		Side:         sellOrder.Side,
		Taker:        sellOrder.UserAddress,
		PairName:     sellOrder.PairName,
		Maker:        buyOrder.UserAddress,
		TakerOrderID: sellOrder.ID,
		MakerOrderID: buyOrder.ID,
		TradeNonce:   big.NewInt(0),
	}
	trade2.Hash = trade2.ComputeHash()

	expectedResponse := &Response{
		Order:      &responseSO,
		Trades:     []*types.Trade{trade, trade1, trade2},
		FillStatus: FULL,
		MatchingOrders: []*FillOrder{
			{responseBO.FilledAmount, &responseBO},
			{responseBO1.FilledAmount, &responseBO1},
			{responseBO2.FilledAmount, &responseBO2}},
	}

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

	sampleSellOrder := types.Order{
		ID:              bson.ObjectIdHex("537f700b537461b70c5f0000"),
		UserAddress:     common.HexToAddress("0x7a9f3cd060ab180f36c17fe6bdf9974f577d77aa"),
		ExchangeAddress: common.HexToAddress("0xae55690d4b079460e6ac28aaa58c9ec7b73a7485"),
		BaseToken:       common.HexToAddress("0x2034842261b82651885751fc293bba7ba5398156"),
		QuoteToken:      common.HexToAddress("0x1888a8db0b7db59413ce07150b3373972bf818d3"),
		BuyAmount:       big.NewInt(6000000000),
		SellAmount:      big.NewInt(13800000000),
		Price:           big.NewInt(229999999),
		PricePoint:      big.NewInt(229999999),
		Amount:          big.NewInt(6000000000),
		FilledAmount:    big.NewInt(0),
		Status:          "NEW",
		Side:            "SELL",
		PairID:          bson.ObjectIdHex("537f700b537461b70c5f0000"),
		PairName:        "ZRX/WETH",
		Expires:         big.NewInt(0),
		MakeFee:         big.NewInt(0),
		Nonce:           big.NewInt(0),
		TakeFee:         big.NewInt(0),
		Signature: &types.Signature{
			V: 28,
			R: common.HexToHash("0x10b30eb0072a4f0a38b6fca0b731cba15eb2e1702845d97c1230b53a839bcb85"),
			S: common.HexToHash("0x6d9ad89548c9e3ce4c97825d027291477f2c44a8caef792095f2cabc978493ff"),
		},
		Hash:      common.HexToHash("0xa2d800b77828cb52c83106ca392e465bc0af0d7c319f6956328f739080c19621"),
		CreatedAt: time.Unix(1405544146, 0),
		UpdatedAt: time.Unix(1405544146, 0),
	}

	sampleBuyOrder := types.Order{
		ID:              bson.ObjectIdHex("537f700b537461b70c5f0000"),
		UserAddress:     common.HexToAddress("0x7a9f3cd060ab180f36c17fe6bdf9974f577d77aa"),
		ExchangeAddress: common.HexToAddress("0xae55690d4b079460e6ac28aaa58c9ec7b73a7485"),
		BaseToken:       common.HexToAddress("0x2034842261b82651885751fc293bba7ba5398156"),
		QuoteToken:      common.HexToAddress("0x1888a8db0b7db59413ce07150b3373972bf818d3"),
		BuyAmount:       big.NewInt(6000000000),
		SellAmount:      big.NewInt(13800000000),
		Price:           big.NewInt(229999999),
		PricePoint:      big.NewInt(229999999),
		Amount:          big.NewInt(6000000000),
		FilledAmount:    big.NewInt(0),
		Status:          "NEW",
		Side:            "BUY",
		PairID:          bson.ObjectIdHex("537f700b537461b70c5f0000"),
		PairName:        "ZRX/WETH",
		Expires:         big.NewInt(0),
		MakeFee:         big.NewInt(0),
		Nonce:           big.NewInt(0),
		TakeFee:         big.NewInt(0),
		Signature: &types.Signature{
			V: 28,
			R: common.HexToHash("0x10b30eb0072a4f0a38b6fca0b731cba15eb2e1702845d97c1230b53a839bcb85"),
			S: common.HexToHash("0x6d9ad89548c9e3ce4c97825d027291477f2c44a8caef792095f2cabc978493ff"),
		},
		Hash:      common.HexToHash("0xa2d800b77828cb52c83106ca392e465bc0af0d7c319f6956328f739080c19671"),
		CreatedAt: time.Unix(1405544146, 0),
		UpdatedAt: time.Unix(1405544146, 0),
	}

	sellOrder := sampleSellOrder
	buyOrder := sampleBuyOrder
	sellOrder1 := sampleSellOrder

	sellOrder1.PricePoint = math.Add(sellOrder1.PricePoint, big.NewInt(10))
	sellOrder1.Nonce = math.Add(sellOrder1.Nonce, big.NewInt(1))
	sellOrder1.Hash = sellOrder1.ComputeHash()

	sellOrder2 := sellOrder1
	sellOrder2.Nonce = math.Add(sellOrder2.Nonce, big.NewInt(1))
	sellOrder2.Hash = sellOrder2.ComputeHash()

	sellOrder3 := sellOrder1
	sellOrder3.PricePoint = math.Add(sellOrder3.PricePoint, big.NewInt(10))
	sellOrder3.Amount = math.Mul(sellOrder3.Amount, big.NewInt(2))

	sellOrder3.Nonce = sellOrder3.Nonce.Add(sellOrder3.Nonce, big.NewInt(1))
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

	trade := &types.Trade{
		Amount:       responseSO.FilledAmount,
		Price:        buyOrder.PricePoint,
		PricePoint:   buyOrder.PricePoint,
		BaseToken:    buyOrder.BaseToken,
		QuoteToken:   buyOrder.QuoteToken,
		OrderHash:    sellOrder.Hash,
		Side:         buyOrder.Side,
		Taker:        buyOrder.UserAddress,
		PairName:     buyOrder.PairName,
		Maker:        sellOrder.UserAddress,
		TakerOrderID: buyOrder.ID,
		MakerOrderID: sellOrder.ID,
		TradeNonce:   big.NewInt(0),
	}

	trade.Hash = trade.ComputeHash()

	trade1 := &types.Trade{
		Amount:       responseSO1.FilledAmount,
		Price:        buyOrder.PricePoint,
		PricePoint:   buyOrder.PricePoint,
		BaseToken:    buyOrder.BaseToken,
		QuoteToken:   buyOrder.QuoteToken,
		OrderHash:    sellOrder1.Hash,
		Side:         buyOrder.Side,
		Taker:        buyOrder.UserAddress,
		PairName:     buyOrder.PairName,
		Maker:        sellOrder.UserAddress,
		TakerOrderID: buyOrder.ID,
		MakerOrderID: sellOrder.ID,
		TradeNonce:   big.NewInt(0),
	}

	trade1.Hash = trade1.ComputeHash()

	trade2 := &types.Trade{
		Amount:       responseSO2.FilledAmount,
		Price:        buyOrder.PricePoint,
		PricePoint:   buyOrder.PricePoint,
		BaseToken:    buyOrder.BaseToken,
		QuoteToken:   buyOrder.QuoteToken,
		OrderHash:    sellOrder2.Hash,
		Side:         buyOrder.Side,
		Taker:        buyOrder.UserAddress,
		PairName:     buyOrder.PairName,
		Maker:        sellOrder.UserAddress,
		TakerOrderID: buyOrder.ID,
		MakerOrderID: sellOrder.ID,
		TradeNonce:   big.NewInt(0),
	}

	trade2.Hash = trade2.ComputeHash()

	trade3 := &types.Trade{
		Amount:       responseSO3.FilledAmount,
		Price:        buyOrder.PricePoint,
		PricePoint:   buyOrder.PricePoint,
		BaseToken:    buyOrder.BaseToken,
		QuoteToken:   buyOrder.QuoteToken,
		OrderHash:    sellOrder3.Hash,
		Side:         buyOrder.Side,
		Taker:        buyOrder.UserAddress,
		PairName:     buyOrder.PairName,
		Maker:        sellOrder.UserAddress,
		TakerOrderID: buyOrder.ID,
		MakerOrderID: sellOrder.ID,
		TradeNonce:   big.NewInt(0),
	}

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
	buyOrder = sampleBuyOrder
	buyOrder.PricePoint = math.Add(buyOrder.PricePoint, big.NewInt(20))
	buyOrder.Amount = math.Mul(buyOrder.Amount, big.NewInt(2))

	responseBO = buyOrder
	responseBO.Status = "PARTIAL_FILLED"
	responseBO.FilledAmount = math.Div(buyOrder.Amount, big.NewInt(2))

	remOrder := sampleBuyOrder
	remOrder.PricePoint = math.Add(remOrder.PricePoint, big.NewInt(20))

	responseSO3.Status = "FILLED"
	responseSO3.FilledAmount = responseSO3.Amount
	trade4 := &types.Trade{
		Amount:       responseBO.FilledAmount,
		Price:        buyOrder.PricePoint,
		PricePoint:   buyOrder.PricePoint,
		BaseToken:    buyOrder.BaseToken,
		QuoteToken:   buyOrder.QuoteToken,
		OrderHash:    sellOrder3.Hash,
		Side:         buyOrder.Side,
		Taker:        buyOrder.UserAddress,
		PairName:     buyOrder.PairName,
		Maker:        sellOrder.UserAddress,
		TakerOrderID: buyOrder.ID,
		MakerOrderID: sellOrder.ID,
		TradeNonce:   big.NewInt(0),
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
	sampleSellOrder := types.Order{
		ID:              bson.ObjectIdHex("537f700b537461b70c5f0000"),
		UserAddress:     common.HexToAddress("0x7a9f3cd060ab180f36c17fe6bdf9974f577d77aa"),
		ExchangeAddress: common.HexToAddress("0xae55690d4b079460e6ac28aaa58c9ec7b73a7485"),
		BaseToken:       common.HexToAddress("0x2034842261b82651885751fc293bba7ba5398156"),
		QuoteToken:      common.HexToAddress("0x1888a8db0b7db59413ce07150b3373972bf818d3"),
		BuyAmount:       big.NewInt(6000000000),
		SellAmount:      big.NewInt(13800000000),
		Price:           big.NewInt(229999999),
		PricePoint:      big.NewInt(229999999),
		Amount:          big.NewInt(6000000000),
		FilledAmount:    big.NewInt(0),
		Status:          "NEW",
		Side:            "SELL",
		PairID:          bson.ObjectIdHex("537f700b537461b70c5f0000"),
		PairName:        "ZRX/WETH",
		Expires:         big.NewInt(0),
		MakeFee:         big.NewInt(0),
		Nonce:           big.NewInt(0),
		TakeFee:         big.NewInt(0),
		Signature: &types.Signature{
			V: 28,
			R: common.HexToHash("0x10b30eb0072a4f0a38b6fca0b731cba15eb2e1702845d97c1230b53a839bcb85"),
			S: common.HexToHash("0x6d9ad89548c9e3ce4c97825d027291477f2c44a8caef792095f2cabc978493ff"),
		},
		Hash:      common.HexToHash("0xa2d800b77828cb52c83106ca392e465bc0af0d7c319f6956328f739080c19621"),
		CreatedAt: time.Unix(1405544146, 0),
		UpdatedAt: time.Unix(1405544146, 0),
	}

	sampleBuyOrder := types.Order{
		ID:              bson.ObjectIdHex("537f700b537461b70c5f0000"),
		UserAddress:     common.HexToAddress("0x7a9f3cd060ab180f36c17fe6bdf9974f577d77aa"),
		ExchangeAddress: common.HexToAddress("0xae55690d4b079460e6ac28aaa58c9ec7b73a7485"),
		BaseToken:       common.HexToAddress("0x2034842261b82651885751fc293bba7ba5398156"),
		QuoteToken:      common.HexToAddress("0x1888a8db0b7db59413ce07150b3373972bf818d3"),
		BuyAmount:       big.NewInt(6000000000),
		SellAmount:      big.NewInt(13800000000),
		PricePoint:      big.NewInt(229999999),
		Price:           big.NewInt(229999999),
		Amount:          big.NewInt(6000000000),
		FilledAmount:    big.NewInt(0),
		Status:          "NEW",
		Side:            "BUY",
		PairID:          bson.ObjectIdHex("537f700b537461b70c5f0000"),
		PairName:        "ZRX/WETH",
		Expires:         big.NewInt(0),
		MakeFee:         big.NewInt(0),
		Nonce:           big.NewInt(0),
		TakeFee:         big.NewInt(0),
		Signature: &types.Signature{
			V: 28,
			R: common.HexToHash("0x10b30eb0072a4f0a38b6fca0b731cba15eb2e1702845d97c1230b53a839bcb85"),
			S: common.HexToHash("0x6d9ad89548c9e3ce4c97825d027291477f2c44a8caef792095f2cabc978493ff"),
		},
		Hash:      common.HexToHash("0xa2d800b77828cb52c83106ca392e465bc0af0d7c319f6956328f739080c19671"),
		CreatedAt: time.Unix(1405544146, 0),
		UpdatedAt: time.Unix(1405544146, 0),
	}

	sellOrder := sampleSellOrder
	buyOrder := sampleBuyOrder

	buyOrder1 := sampleBuyOrder
	buyOrder1.PricePoint = math.Sub(buyOrder1.PricePoint, big.NewInt(10))
	buyOrder1.Nonce = math.Add(buyOrder1.Nonce, big.NewInt(1))
	buyOrder1.Hash = buyOrder1.ComputeHash()

	buyOrder2 := buyOrder1
	buyOrder2.Nonce = math.Add(buyOrder2.Nonce, big.NewInt(1))
	buyOrder2.Hash = buyOrder2.ComputeHash()

	buyOrder3 := buyOrder1
	buyOrder3.PricePoint = math.Sub(buyOrder3.PricePoint, big.NewInt(10))
	buyOrder3.Amount = math.Mul(buyOrder3.Amount, big.NewInt(2))
	buyOrder3.Nonce = buyOrder3.Nonce.Add(buyOrder3.Nonce, big.NewInt(1))
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

	trade := &types.Trade{
		Amount:       responseBO.FilledAmount,
		Price:        sellOrder.PricePoint,
		PricePoint:   sellOrder.PricePoint,
		BaseToken:    sellOrder.BaseToken,
		QuoteToken:   sellOrder.QuoteToken,
		OrderHash:    buyOrder.Hash,
		Side:         sellOrder.Side,
		Taker:        sellOrder.UserAddress,
		PairName:     sellOrder.PairName,
		Maker:        buyOrder.UserAddress,
		TakerOrderID: sellOrder.ID,
		MakerOrderID: buyOrder.ID,
		TradeNonce:   big.NewInt(0),
	}

	trade.Hash = trade.ComputeHash()

	trade1 := &types.Trade{
		Amount:       responseBO1.FilledAmount,
		Price:        sellOrder.PricePoint,
		PricePoint:   sellOrder.PricePoint,
		BaseToken:    sellOrder.BaseToken,
		QuoteToken:   sellOrder.QuoteToken,
		OrderHash:    buyOrder1.Hash,
		Side:         sellOrder.Side,
		Taker:        sellOrder.UserAddress,
		PairName:     sellOrder.PairName,
		Maker:        buyOrder.UserAddress,
		TakerOrderID: sellOrder.ID,
		MakerOrderID: buyOrder.ID,
		TradeNonce:   big.NewInt(0),
	}

	trade1.Hash = trade1.ComputeHash()

	trade2 := &types.Trade{
		Amount:       responseBO2.FilledAmount,
		Price:        sellOrder.PricePoint,
		PricePoint:   sellOrder.PricePoint,
		BaseToken:    sellOrder.BaseToken,
		QuoteToken:   sellOrder.QuoteToken,
		OrderHash:    buyOrder2.Hash,
		Side:         sellOrder.Side,
		Taker:        sellOrder.UserAddress,
		PairName:     sellOrder.PairName,
		Maker:        buyOrder.UserAddress,
		TakerOrderID: sellOrder.ID,
		MakerOrderID: buyOrder.ID,
		TradeNonce:   big.NewInt(0),
	}

	trade2.Hash = trade2.ComputeHash()

	trade3 := &types.Trade{
		Amount:       responseBO3.FilledAmount,
		Price:        sellOrder.PricePoint,
		PricePoint:   sellOrder.PricePoint,
		BaseToken:    sellOrder.BaseToken,
		QuoteToken:   sellOrder.QuoteToken,
		OrderHash:    buyOrder3.Hash,
		Side:         sellOrder.Side,
		Taker:        sellOrder.UserAddress,
		PairName:     sellOrder.PairName,
		Maker:        buyOrder.UserAddress,
		TakerOrderID: sellOrder.ID,
		MakerOrderID: buyOrder.ID,
		TradeNonce:   big.NewInt(0),
	}

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

	log.Print("Lets asdfkasdjf;asldkf")
	response.Print()

	resBytes, _ := json.Marshal(response)
	assert.JSONEq(t, string(erBytes), string(resBytes))

	// Try matching remaining buyOrder with bigger sellOrder amount (partial filled sell Order)
	sellOrder = sampleSellOrder
	sellOrder.PricePoint = math.Sub(sellOrder.PricePoint, big.NewInt(20))
	sellOrder.Amount = math.Mul(sellOrder.Amount, big.NewInt(2))

	responseSO = sellOrder
	responseSO.Status = "PARTIAL_FILLED"
	responseSO.FilledAmount = math.Div(sellOrder.Amount, big.NewInt(2))

	remOrder := sampleSellOrder
	remOrder.PricePoint = math.Sub(remOrder.PricePoint, big.NewInt(20))
	remOrder.Hash = remOrder.ComputeHash()
	remOrder.Signature = nil

	responseBO3.Status = "FILLED"
	responseBO3.FilledAmount = responseBO3.Amount

	trade4 := &types.Trade{
		Amount:       responseSO.FilledAmount,
		Price:        sellOrder.PricePoint,
		PricePoint:   sellOrder.PricePoint,
		BaseToken:    sellOrder.BaseToken,
		QuoteToken:   sellOrder.QuoteToken,
		OrderHash:    buyOrder3.Hash,
		Side:         sellOrder.Side,
		Taker:        sellOrder.UserAddress,
		PairName:     sellOrder.PairName,
		Maker:        buyOrder.UserAddress,
		TakerOrderID: sellOrder.ID,
		MakerOrderID: buyOrder.ID,
		TradeNonce:   big.NewInt(0),
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
