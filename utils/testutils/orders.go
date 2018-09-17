package testutils

import (
	"math/big"
	"time"

	"github.com/Proofsuite/amp-matching-engine/types"
	"github.com/ethereum/go-ethereum/common"
	"gopkg.in/mgo.v2/bson"
)

func GetTestOrder1() types.Order {
	return types.Order{
		ID:              bson.ObjectIdHex("537f700b537461b70c5f0000"),
		UserAddress:     common.HexToAddress("0x7a9f3cd060ab180f36c17fe6bdf9974f577d77aa"),
		ExchangeAddress: common.HexToAddress("0xae55690d4b079460e6ac28aaa58c9ec7b73a7485"),
		BuyToken:        common.HexToAddress("0xe41d2489571d322189246dafa5ebde1f4699f498"),
		SellToken:       common.HexToAddress("0x12459c951127e0c374ff9105dda097662a027093"),
		BaseToken:       common.HexToAddress("0xe41d2489571d322189246dafa5ebde1f4699f498"),
		QuoteToken:      common.HexToAddress("0x12459c951127e0c374ff9105dda097662a027093"),
		BuyAmount:       big.NewInt(1000),
		SellAmount:      big.NewInt(100),
		PricePoint:      big.NewInt(1000),
		Amount:          big.NewInt(1000),
		FilledAmount:    big.NewInt(100),
		Status:          "OPEN",
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
		Hash:      common.HexToHash("0xb9070a2d333403c255ce71ddf6e795053599b2e885321de40353832b96d8880a"),
		CreatedAt: time.Unix(1405544146, 0),
		UpdatedAt: time.Unix(1405544146, 0),
	}
}

func GetTestOrder2() types.Order {
	return types.Order{
		ID:              bson.ObjectIdHex("537f700b537461b70c5f0000"),
		UserAddress:     common.HexToAddress("0x7a9f3cd060ab180f36c17fe6bdf9974f577d77aa"),
		ExchangeAddress: common.HexToAddress("0xae55690d4b079460e6ac28aaa58c9ec7b73a7485"),
		BuyToken:        common.HexToAddress("0x4bc89ac6f1c55ea645294f3fed949813a768ac6d"),
		SellToken:       common.HexToAddress("0xd27a76b12bc4a870c1045c86844161337393d9fa"),
		BaseToken:       common.HexToAddress("0x4bc89ac6f1c55ea645294f3fed949813a768ac6d"),
		QuoteToken:      common.HexToAddress("0xd27a76b12bc4a870c1045c86844161337393d9fa"),
		BuyAmount:       big.NewInt(1000),
		SellAmount:      big.NewInt(100),
		PricePoint:      big.NewInt(1200),
		Amount:          big.NewInt(1000),
		FilledAmount:    big.NewInt(100),
		Status:          "OPEN",
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
		Hash:      common.HexToHash("0xecf27444c5ce65a88f73db628687fb9b4ac2686b5577df405958d47bee8eaa53"),
		CreatedAt: time.Unix(1405544146, 0),
		UpdatedAt: time.Unix(1405544146, 0),
	}
}

func GetTestOrder3() types.Order {
	return types.Order{
		ID:              bson.ObjectIdHex("537f700b537461b70c5f0000"),
		UserAddress:     common.HexToAddress("0x7a9f3cd060ab180f36c17fe6bdf9974f577d77aa"),
		ExchangeAddress: common.HexToAddress("0xae55690d4b079460e6ac28aaa58c9ec7b73a7485"),
		BuyToken:        common.HexToAddress("0x4bc89ac6f1c55ea645294f3fed949813a768ac6d"),
		SellToken:       common.HexToAddress("0xd27a76b12bc4a870c1045c86844161337393d9fa"),
		BaseToken:       common.HexToAddress("0x4bc89ac6f1c55ea645294f3fed949813a768ac6d"),
		QuoteToken:      common.HexToAddress("0xd27a76b12bc4a870c1045c86844161337393d9fa"),
		BuyAmount:       big.NewInt(1000),
		SellAmount:      big.NewInt(100),
		PricePoint:      big.NewInt(1200),
		Amount:          big.NewInt(1000),
		FilledAmount:    big.NewInt(100),
		Status:          "OPEN",
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
		Hash:      common.HexToHash("0x400558b2f5a7b20dd06241c2313c08f652b297e819926b5a51a5abbc60f451e6"),
		CreatedAt: time.Unix(1405544146, 0),
		UpdatedAt: time.Unix(1405544146, 0),
	}
}
