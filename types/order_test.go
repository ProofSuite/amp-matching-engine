package types

import (
	"encoding/json"
	"fmt"
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/go-test/deep"
	"github.com/stretchr/testify/assert"
	"gopkg.in/mgo.v2/bson"
)

func TestOrderMarshal(t *testing.T) {

	o := &Order{
		ID:              bson.ObjectIdHex("537f700b537461b70c5f0000"),
		UserAddress:     common.HexToAddress("0x7a9f3cd060ab180f36c17fe6bdf9974f577d77aa"),
		ExchangeAddress: common.HexToAddress("0xae55690d4b079460e6ac28aaa58c9ec7b73a7485"),
		BuyToken:        common.HexToAddress("0xe41d2489571d322189246dafa5ebde1f4699f498"),
		SellToken:       common.HexToAddress("0x12459c951127e0c374ff9105dda097662a027093"),
		BaseToken:       common.HexToAddress("0xe41d2489571d322189246dafa5ebde1f4699f498"),
		QuoteToken:      common.HexToAddress("0x12459c951127e0c374ff9105dda097662a027093"),
		BuyAmount:       big.NewInt(1000),
		SellAmount:      big.NewInt(100),
		Price:           big.NewInt(1000),
		Amount:          big.NewInt(1000),
		FilledAmount:    big.NewInt(100),
		Status:          "NEW",
		Side:            "BUY",
		PairID:          bson.ObjectIdHex("537f700b537461b70c5f0000"),
		PairName:        "ZRX/WETH",
		Expires:         big.NewInt(10000),
		MakeFee:         big.NewInt(50),
		Nonce:           big.NewInt(1000),
		TakeFee:         big.NewInt(50),
		Signature: &Signature{
			V: 28,
			R: common.HexToHash("0x10b30eb0072a4f0a38b6fca0b731cba15eb2e1702845d97c1230b53a839bcb85"),
			S: common.HexToHash("0x6d9ad89548c9e3ce4c97825d027291477f2c44a8caef792095f2cabc978493ff"),
		},
		Hash:      common.HexToHash("0xb9070a2d333403c255ce71ddf6e795053599b2e885321de40353832b96d8880a"),
		CreatedAt: time.Unix(1405544146, 0),
		UpdatedAt: time.Unix(1405544146, 0),
	}

	expected := map[string]interface{}{
		"id":              "537f700b537461b70c5f0000",
		"userAddress":     "0x7a9f3cd060ab180f36c17fe6bdf9974f577d77aa",
		"exchangeAddress": "0xae55690d4b079460e6ac28aaa58c9ec7b73a7485",
		"buyToken":        "0xe41d2489571d322189246dafa5ebde1f4699f498",
		"sellToken":       "0x12459c951127e0c374ff9105dda097662a027093",
		"baseToken":       "0xe41d2489571d322189246dafa5ebde1f4699f498",
		"quoteToken":      "0x12459c951127e0c374ff9105dda097662a027093",
		"buyAmount":       "1000",
		"sellAmount":      "100",
		"price":           "1000",
		"amount":          "1000",
		"filledAmount":    "100",
		"status":          "NEW",
		"side":            "BUY",
		"pairID":          "537f700b537461b70c5f0000",
		"pairName":        "ZRX/WETH",
		"expires":         "10000",
		"makeFee":         "50",
		"takeFee":         "50",
		"nonce":           "1000",
		"signature": map[string]interface{}{
			"V": 28,
			"R": "0x10b30eb0072a4f0a38b6fca0b731cba15eb2e1702845d97c1230b53a839bcb85",
			"S": "0x6d9ad89548c9e3ce4c97825d027291477f2c44a8caef792095f2cabc978493ff",
		},
		"hash":      "0xb9070a2d333403c255ce71ddf6e795053599b2e885321de40353832b96d8880a",
		"createdAt": "2014-07-17 05:55:46 +0900 KST",
		"updatedAt": "2014-07-17 05:55:46 +0900 KST",
	}

	encoded, err := json.Marshal(o)
	if err != nil {
		t.Errorf("Error encoding order: %v", err)
	}

	expectedJSON, err := json.Marshal(expected)
	if err != nil {
		t.Errorf("Error encoding json: %v", err)
	}

	if diff := deep.Equal(encoded, expectedJSON); diff != nil {
		fmt.Printf("Expected: \n%s\nGot: \n%s\n\n", encoded, expectedJSON)
	}

}

func TestOrderUnmarshal(t *testing.T) {
	expected := Order{
		ID:              bson.ObjectIdHex("537f700b537461b70c5f0000"),
		UserAddress:     common.HexToAddress("0x14d281013d8ee8ccfa0eca87524e5b3cfa6152ba"),
		ExchangeAddress: common.HexToAddress("0xae55690d4b079460e6ac28aaa58c9ec7b73a7485"),
		BuyToken:        common.HexToAddress("0xe41d2489571d322189246dafa5ebde1f4699f498"),
		SellToken:       common.HexToAddress("0x12459c951127e0c374ff9105dda097662a027093"),
		QuoteToken:      common.HexToAddress("0xe41d2489571d322189246dafa5ebde1f4699f498"),
		BaseToken:       common.HexToAddress("0x12459c951127e0c374ff9105dda097662a027093"),
		BuyAmount:       big.NewInt(1000),
		SellAmount:      big.NewInt(100),
		Amount:          big.NewInt(100),
		Price:           big.NewInt(100),
		FilledAmount:    big.NewInt(1000),
		Status:          "NEW",
		Side:            "BUY",
		Expires:         big.NewInt(10000),
		Nonce:           big.NewInt(1000),
		MakeFee:         big.NewInt(50),
		TakeFee:         big.NewInt(50),
		Signature: &Signature{
			V: 28,
			R: common.HexToHash("0x10b30eb0072a4f0a38b6fca0b731cba15eb2e1702845d97c1230b53a839bcb85"),
			S: common.HexToHash("0x6d9ad89548c9e3ce4c97825d027291477f2c44a8caef792095f2cabc978493ff"),
		},
		PairName: "ZRX/WETH",
		PairID:   bson.ObjectIdHex("537f700b537461b70c5f0000"),
		Hash:     common.HexToHash("0xb9070a2d333403c255ce71ddf6e795053599b2e885321de40353832b96d8880a"),
	}

	payload :=
		`{
			"id": "537f700b537461b70c5f0000",
			"userAddress": "0x14d281013d8ee8ccfa0eca87524e5b3cfa6152ba",
			"exchangeAddress": "0xae55690d4b079460e6ac28aaa58c9ec7b73a7485",
			"buyToken":"0xe41d2489571d322189246dafa5ebde1f4699f498",
			"sellToken":"0x12459c951127e0c374ff9105dda097662a027093",
			"quoteToken":"0xe41d2489571d322189246dafa5ebde1f4699f498",
			"baseToken":"0x12459c951127e0c374ff9105dda097662a027093",
			"buyAmount":"1000",
			"sellAmount":"100",
			"amount": "100",
			"price": "100",
			"filledAmount": "1000",
			"status": "NEW",
			"side": "BUY",
			"expires": "10000",
			"nonce":   "1000",
			"makeFee": "50",
			"takeFee": "50",
			"signature": {
				"V": 28,
				"R": "0x10b30eb0072a4f0a38b6fca0b731cba15eb2e1702845d97c1230b53a839bcb85",
				"S": "0x6d9ad89548c9e3ce4c97825d027291477f2c44a8caef792095f2cabc978493ff"
			},
			"pairName": "ZRX/WETH",
			"pairID": "537f700b537461b70c5f0000",
			"hash":"0xb9070a2d333403c255ce71ddf6e795053599b2e885321de40353832b96d8880a"
		}`

	order := Order{}
	err := json.Unmarshal([]byte(payload), &order)
	if err != nil {
		t.Errorf("Could not unmarshal payload: %v", err)
	}

	if diff := deep.Equal(order, expected); diff != nil {
		fmt.Printf("Expected: \n%+v\nGot: \n%+v\n\n", expected, order)
	}
}

func TestOrderBSON(t *testing.T) {
	order := &Order{
		ID:              bson.ObjectIdHex("537f700b537461b70c5f0000"),
		UserAddress:     common.HexToAddress("0x7a9f3cd060ab180f36c17fe6bdf9974f577d77aa"),
		ExchangeAddress: common.HexToAddress("0xae55690d4b079460e6ac28aaa58c9ec7b73a7485"),
		BuyToken:        common.HexToAddress("0xe41d2489571d322189246dafa5ebde1f4699f498"),
		SellToken:       common.HexToAddress("0x12459c951127e0c374ff9105dda097662a027093"),
		BaseToken:       common.HexToAddress("0xe41d2489571d322189246dafa5ebde1f4699f498"),
		QuoteToken:      common.HexToAddress("0x12459c951127e0c374ff9105dda097662a027093"),
		BuyAmount:       big.NewInt(1000),
		SellAmount:      big.NewInt(100),
		Price:           big.NewInt(1000),
		Amount:          big.NewInt(1000),
		FilledAmount:    big.NewInt(100),
		Status:          "NEW",
		Side:            "BUY",
		PairID:          bson.ObjectIdHex("537f700b537461b70c5f0000"),
		PairName:        "ZRX/WETH",
		Expires:         big.NewInt(10000),
		MakeFee:         big.NewInt(50),
		Nonce:           big.NewInt(1000),
		TakeFee:         big.NewInt(50),
		Signature: &Signature{
			V: 28,
			R: common.HexToHash("0x10b30eb0072a4f0a38b6fca0b731cba15eb2e1702845d97c1230b53a839bcb85"),
			S: common.HexToHash("0x6d9ad89548c9e3ce4c97825d027291477f2c44a8caef792095f2cabc978493ff"),
		},
		Hash:      common.HexToHash("0xb9070a2d333403c255ce71ddf6e795053599b2e885321de40353832b96d8880a"),
		CreatedAt: time.Unix(1405544146, 0),
		UpdatedAt: time.Unix(1405544146, 0),
	}

	data, err := bson.Marshal(order)
	if err != nil {
		t.Error(err)
	}

	decoded := &Order{}
	if err := bson.Unmarshal(data, decoded); err != nil {
		t.Error(err)
	}

	assert.Equal(t, decoded, order)
}

// func TestAccountBSON(t *testing.T) {
// 	assert := assert.New(t)

// 	address := NewAddressFromString("0xe8e84ee367bc63ddb38d3d01bccef106c194dc47")
// 	tokenAddress1 := NewAddressFromString("0xcf7389dc6c63637598402907d5431160ec8972a5")
// 	tokenAddress2 := NewAddressFromString("0x7a9f3cd060ab180f36c17fe6bdf9974f577d77aa")

// 	tokenBalance1 := &TokenBalance{
// 		ID:       bson.NewObjectId(),
// 		Address:  tokenAddress1,
// 		Symbol:   "EOS",
// 		Balance:       NewBigInt("10000"),
// 		Allowance:     NewBigInt("10000"),
// 		LockedBalance: NewBigInt("5000"),
// 	}

// 	tokenBalance2 := &TokenBalance{
// 		ID:       bson.NewObjectId(),
// 		Address:  tokenAddress2,
// 		Symbol:   "ZRX",
// 		Balance:       NewBigInt("10000"),
// 		Allowance:     NewBigInt("10000"),
// 		LockedBalance: NewBigInt("5000"),
// 	}

// 	account := &Account{
// 		ID:      bson.NewObjectId(),
// 		Address: address,
// 		TokenBalances: map[common.Address]*TokenBalance{
// 			tokenAddress1: tokenBalance1,
// 			tokenAddress2: tokenBalance2,
// 		},
// 		IsBlocked: false,
// 	}

// 	data, err := bson.Marshal(account)
// 	if err != nil {
// 		t.Error(err)
// 	}

// 	decoded := &Account{}
// 	if err := bson.Unmarshal(data, decoded); err != nil {
// 		t.Error(err)
// 	}

// 	assert.Equal(decoded, account)
// }
