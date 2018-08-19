package engine

import (
	"math/big"
	"strconv"
	"testing"
	"time"

	"github.com/Proofsuite/amp-matching-engine/types"
	"github.com/Proofsuite/amp-matching-engine/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"gopkg.in/mgo.v2/bson"
)

func TestAddOrder(t *testing.T) {
	e, s := getResource()
	defer s.Close()

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
		Price:           229999999,
		Amount:          6000000000,
		FilledAmount:    0,
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

	rs, err := s.SortedSet(ssKey)
	if err != nil {
		t.Error(err)
	} else {
		for k, v := range rs {
			assert.Equalf(t, utils.UintToPaddedString(o1.Price), k, "Expected sorted set value: %v got: %v", utils.UintToPaddedString(o1.Price), k)
			assert.Equalf(t, 0.0, v, "Expected sorted set value: %v got: %v", 0, v)
		}
	}

	rs, err = s.SortedSet(listKey)
	if err != nil {
		t.Error(err)
	} else {
		for k, v := range rs {
			assert.Equalf(t, o1.Hash.Hex(), k, "Expected sorted set value: %v got: %v", o1.Hash.Hex(), k)
			assert.Equalf(t, o1.CreatedAt.Unix(), int64(v), "Expected sorted set value: %v got: %v", o1.CreatedAt.Unix(), v)
		}
	}

	rse, err := s.Get(listKey + "::" + o1.Hash.Hex())
	if err != nil {
		t.Error(err)
	}
	assert.JSONEqf(t, string(bytes), rse, "Expected value for key: %v, was: %s, but got: %v", ssKey, bytes, rse)

	rse, err = s.Get(ssKey + "::book::" + utils.UintToPaddedString(o1.Price))
	if err != nil {
		t.Error(err)
	}

	assert.Equalf(t, strconv.FormatInt(o1.Amount, 10), rse, "Expected value for key: %v, was: %v, but got: %v", ssKey, o1.Amount, rse)

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
		Price:           6000000000,
		Amount:          229999999,
		FilledAmount:    1000000000,
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

	rs, err = s.SortedSet(ssKey)
	if err != nil {
		t.Error(err)
	} else {
		var matched = false
		for k, v := range rs {
			if utils.UintToPaddedString(o2.Price) == k && v == 0.0 {
				matched = true
			}
		}
		if !matched {
			t.Errorf("Expected sorted set value: %v", utils.UintToPaddedString(o2.Price))
		}
	}

	rs, err = s.SortedSet(listKey)
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

	rse, err = s.Get(listKey + "::" + o2.Hash.Hex())
	if err != nil {
		t.Error(err)
	}
}

func TestCancelOrder(t *testing.T) {
	e, s := getResource()
	defer s.Close()

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
		Price:           229999999,
		Amount:          6000000000,
		FilledAmount:    0,
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
		Price:           229999999,
		Amount:          6000000000,
		FilledAmount:    0,
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

	rs, err := s.SortedSet(ss1Key)
	if err != nil {
		t.Error(err)
	} else {
		for k, v := range rs {
			assert.Equalf(t, utils.UintToPaddedString(o2.Price), k, "Expected sorted set value: %v got: %v", utils.UintToPaddedString(o2.Price), k)
			assert.Equalf(t, 0.0, v, "Expected sorted set value: %v got: %v", 0, v)
		}
	}

	rs, err = s.SortedSet(list1Key)
	if err != nil {
		t.Error(err)
	} else {
		for k, v := range rs {
			assert.NotEqualf(t, o2.Hash, k, "Key : %v expected to be deleted but key exists", o2.Hash)
			assert.Equalf(t, o2.CreatedAt.Unix(), int64(v), "Expected sorted set value: %v got: %v", o2.CreatedAt.Unix(), v)
		}
	}

	if s.Exists(list1Key + "::" + o2.Hash.Hex()) {
		t.Errorf("Key : %v expected to be deleted but key exists", list1Key+"::"+o2.Hash.Hex())
	}

	rse, err := s.Get(ss1Key + "::book::" + utils.UintToPaddedString(o2.Price))
	if err != nil {
		t.Error(err)
	}
	assert.Equalf(t, strconv.FormatInt(o1.Amount, 10), rse, "Expected value for key: %v, was: %v, but got: %v", ss1Key, o2.Amount-o2.FilledAmount, rse)

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
	if s.Exists(ssKey) {
		t.Errorf("Key : %v expected to be deleted but key exists", ssKey)
	}

	if s.Exists(listKey) {
		t.Errorf("Key : %v expected to be deleted but key exists", ssKey)
	}

	if s.Exists(listKey + "::" + o1.Hash.Hex()) {
		t.Errorf("Key : %v expected to be deleted but key exists", ssKey)
	}

	if s.Exists(ssKey + "::book::" + utils.UintToPaddedString(o1.Price)) {
		t.Errorf("Key : %v expected to be deleted but key exists", ssKey)
	}
}

func TestRecoverOrders(t *testing.T) {
	e, s := getResource()
	defer s.Close()

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
		Price:           229999999,
		Amount:          6000000000,
		FilledAmount:    1000000000,
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
	expectedOrder1.FilledAmount = 0
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
		Price:           229999999,
		Amount:          6000000000,
		FilledAmount:    6000000000,
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
	expectedOrder2.FilledAmount = 0
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
		Price:           229999999,
		Amount:          6000000000,
		FilledAmount:    6000000000,
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
	expectedOrder3.FilledAmount = 4000000000
	expectedOrder3Json, _ := expectedOrder3.MarshalJSON()

	e.addOrder(o3)

	recoverOrders := []*FillOrder{
		&FillOrder{
			Amount: 1000000000,
			Order:  o1,
		}, &FillOrder{
			Amount: 6000000000,
			Order:  o2,
		}, &FillOrder{
			Amount: 2000000000,
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

	rs, err := s.SortedSet(ssKey)
	if err != nil {
		t.Error(err)
	} else {
		for k, v := range rs {
			assert.Equalf(t, utils.UintToPaddedString(o1.Price), k, "Expected sorted set value: %v got: %v", utils.UintToPaddedString(o2.Price), k)
			assert.Equalf(t, 0.0, v, "Expected sorted set value: %v got: %v", 0, v)
		}
	}

	expectedMap := map[string]float64{
		o1.Hash.Hex(): float64(o1.CreatedAt.Unix()),
		o2.Hash.Hex(): float64(o2.CreatedAt.Unix()),
		o3.Hash.Hex(): float64(o3.CreatedAt.Unix()),
	}

	rs, err = s.SortedSet(listKey)
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, expectedMap, rs)

	rse, err := s.Get(listKey + "::" + o1.Hash.Hex())
	if err != nil {
		t.Error(err)
	}

	assert.JSONEqf(t, string(expectedOrder1Json), rse, "Expected value for key: %v, was: %s, but got: %v", listKey+"::"+o1.Hash.Hex(), expectedOrder1Json, rse)

	rse, err = s.Get(list1Key + "::" + o2.Hash.Hex())
	if err != nil {
		t.Error(err)
	}

	assert.JSONEqf(t, string(expectedOrder2Json), rse, "Expected value for key: %v, was: %s, but got: %v", list1Key+"::"+o2.Hash.Hex(), expectedOrder2Json, rse)

	rse, err = s.Get(list2Key + "::" + o3.Hash.Hex())
	if err != nil {
		t.Error(err)
	}

	assert.JSONEqf(t, string(expectedOrder3Json), rse, "Expected value for key: %v, was: %s, but got: %v", list2Key+"::"+o3.Hash.Hex(), expectedOrder3Json, rse)

	rse, err = s.Get(ssKey + "::book::" + utils.UintToPaddedString(o1.Price))
	if err != nil {
		t.Error(err)
	}

	assert.Equalf(t, strconv.FormatInt(o1.Amount+o2.Amount+o3.Amount-o3.FilledAmount, 10), rse, "Expected value for key: %v, was: %v, but got: %v", ssKey, o1.Amount+o2.Amount+o3.Amount, rse)
}

func TestSellOrder(t *testing.T) {
	e, s := getResource()
	defer s.Close()

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
		Price:           229999999,
		Amount:          6000000000,
		FilledAmount:    0,
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
		RemainingOrder: &types.Order{},
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
	e, s := getResource()
	defer s.Close()

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
		Price:           229999999,
		Amount:          6000000000,
		FilledAmount:    0,
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
		RemainingOrder: &types.Order{},
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

/* Can not be tested yet, as SORT command is not supported by miniredis */
// TODO: Find alternative or a way to test SORT function of redis
// func TestFillOrder(t *testing.T) {
// 	e, s := getResource()
// 	defer s.Close()

// 	// New Buy Order
// 	sampleOrder := &types.Order{}
// 	orderJSON := []byte(`{ "id": "5b6ac5297b4457546d64379d", "sellToken": "0x1888a8db0b7db59413ce07150b3373972bf818d3", "buyToken": "0x2034842261b82651885751fc293bba7ba5398156", "baseToken": "0x2034842261b82651885751fc293bba7ba5398156", "quoteToken": "0x1888a8db0b7db59413ce07150b3373972bf818d3", "sellAmount": 6000000000, "buyAmount": 13800000000, "nonce": 0, "userAddress": "0xefD7eB287CeeFCE8256Dd46e25F398acEA7C4b63", "hash": "0xa2d800b77828cb52c83106ca392e465bc0af0d7c319f6956328f739080c19621", "side": "BUY", "amount": 6000000000, "price": 229999999, "filledAmount": 0, "fee": 0, "makeFee": 0, "takeFee": 0, "exchangeAddress": "", "status": "NEW", "pairID": "5b6ac5117b445753ee755fb8", "pairName": "HPC/AUT", "orderBook": null, "createdAt": "2018-08-08T15:55:45.062141954+05:30", "updatedAt": "2018-08-08T15:55:45.06214208+05:30" }`)
// 	json.Unmarshal(orderJSON, &sampleOrder)
// 	expectedResponse := &Response{
// 		Order:          sampleOrder,
// 		RemainingOrder: &types.Order{},
// 		Trades:         make([]*types.Trade, 0),
// 		FillStatus:     NOMATCH,
// 		MatchingOrders: make([]*FillOrder, 0),
// 	}
// 	expectedResponse.Order.Status = types.OPEN
// 	response, err := e.buyOrder(sampleOrder)
// 	if err != nil {
// 		t.Errorf("Error in sellOrder: %s", err)
// 	}
// 	assert.Equal(t, expectedResponse, response)

// 	// New Sell Order
// 	sampleOrder1 := &types.Order{}
// 	orderJSON = []byte(`{ "id": "5b6ac5297b4457546d64379e", "buyToken": "0x1888a8db0b7db59413ce07150b3373972bf818d3", "sellToken": "0x2034842261b82651885751fc293bba7ba5398156", "baseToken": "0x2034842261b82651885751fc293bba7ba5398156", "quoteToken": "0x1888a8db0b7db59413ce07150b3373972bf818d3", "buyAmount": 6000000000, "sellAmount": 13800000000, "nonce": 0, "userAddress": "0xefD7eB287CeeFCE8256Dd46e25F398acEA7C4b63", "hash": "0xa2d800b77828cb52c83106ca392e465bc0af0d7c319f6956328f739080c19622", "side": "SELL", "amount": 6000000000, "price": 229999999, "filledAmount": 0, "fee": 0, "makeFee": 0, "takeFee": 0, "exchangeAddress": "", "status": "NEW", "pairID": "5b6ac5117b445753ee755fb8", "pairName": "HPC/AUT", "orderBook": null, "createdAt": "2018-08-08T15:55:45.062141954+05:30", "updatedAt": "2018-08-08T15:55:45.06214208+05:30" }`)
// 	json.Unmarshal(orderJSON, &sampleOrder1)
// 	expectedResponse = &Response{
// 		Order:          sampleOrder1,
// 		RemainingOrder: &types.Order{},
// 		Trades:         make([]*types.Trade, 0),
// 		FillStatus:     FULL,
// 		MatchingOrders: []*FillOrder{&FillOrder{Amount:6000000000,Order:sampleOrder}},
// 	}
// 	expectedResponse.Order.Status = types.FILLED
// 	expectedResponse.Order.FilledAmount = expectedResponse.Order.Amount

// 	response, err = e.sellOrder(sampleOrder1)
// 	if err != nil {
// 		t.Errorf("Error in sellOrder: %s", err)
// 	}
// 	assert.Equal(t, expectedResponse, response)
// }
