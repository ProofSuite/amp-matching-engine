package daos

import (
	"io/ioutil"
	"math/big"
	"testing"
	"time"

	"github.com/Proofsuite/amp-matching-engine/types"
	"github.com/Proofsuite/amp-matching-engine/utils/testutils"
	"github.com/ethereum/go-ethereum/common"
	"gopkg.in/mgo.v2/bson"
)

func init() {
	temp, _ := ioutil.TempDir("", "test")
	server.SetPath(temp)

	session := server.Session()
	db = &Database{session}
}

func TestUpdateOrderByHash(t *testing.T) {
	o := &types.Order{
		ID:              bson.ObjectIdHex("537f700b537461b70c5f0001"),
		UserAddress:     common.HexToAddress("0x7a9f3cd060ab180f36c17fe6bdf9974f577d77aa"),
		ExchangeAddress: common.HexToAddress("0xae55690d4b079460e6ac28aaa58c9ec7b73a7485"),
		BuyToken:        common.HexToAddress("0xe41d2489571d322189246dafa5ebde1f4699f498"),
		SellToken:       common.HexToAddress("0x12459c951127e0c374ff9105dda097662a027093"),
		BaseToken:       common.HexToAddress("0xe41d2489571d322189246dafa5ebde1f4699f498"),
		QuoteToken:      common.HexToAddress("0x12459c951127e0c374ff9105dda097662a027093"),
		BuyAmount:       big.NewInt(1000),
		SellAmount:      big.NewInt(100),
		Price:           big.NewInt(1000),
		PricePoint:      big.NewInt(1000),
		Amount:          big.NewInt(1000),
		FilledAmount:    big.NewInt(100),
		Status:          "NEW",
		Side:            "BUY",
		PairID:          bson.ObjectIdHex("537f700b537461b70c5f0001"),
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

	dao := NewOrderDao()

	err := dao.Create(o)
	if err != nil {
		t.Errorf("Could not create order object")
	}

	updated := &types.Order{
		ID:              o.ID,
		UserAddress:     o.UserAddress,
		ExchangeAddress: o.ExchangeAddress,
		BuyToken:        o.BuyToken,
		SellToken:       o.SellToken,
		BaseToken:       o.BaseToken,
		QuoteToken:      o.QuoteToken,
		BuyAmount:       big.NewInt(1000),
		SellAmount:      big.NewInt(100),
		Price:           big.NewInt(4000),
		Amount:          big.NewInt(4000),
		FilledAmount:    big.NewInt(200),
		Status:          "FILLED",
		Side:            "BUY",
		PairID:          o.PairID,
		PairName:        "ZRX/WETH",
		Expires:         big.NewInt(10000),
		MakeFee:         big.NewInt(50),
		Nonce:           big.NewInt(1000),
		TakeFee:         big.NewInt(50),
		Signature:       o.Signature,
		Hash:            o.Hash,
		CreatedAt:       o.CreatedAt,
		UpdatedAt:       o.UpdatedAt,
	}

	err = dao.UpdateByHash(
		o.Hash,
		updated,
	)

	if err != nil {
		t.Errorf("Could not updated order from hash %v", err)
	}

	queried, err := dao.GetByHash(o.Hash)
	if err != nil {
		t.Errorf("Could not get order by hash")
	}

	testutils.CompareOrder(t, queried, updated)
}

func TestOrderUpdate(t *testing.T) {
	dao := NewOrderDao()
	err := dao.Drop()
	if err != nil {
		t.Errorf("Could not drop previous order state")
	}

	o := &types.Order{
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
		Signature: &types.Signature{
			V: 28,
			R: common.HexToHash("0x10b30eb0072a4f0a38b6fca0b731cba15eb2e1702845d97c1230b53a839bcb85"),
			S: common.HexToHash("0x6d9ad89548c9e3ce4c97825d027291477f2c44a8caef792095f2cabc978493ff"),
		},
		Hash:      common.HexToHash("0xb9070a2d333403c255ce71ddf6e795053599b2e885321de40353832b96d8880a"),
		CreatedAt: time.Unix(1405544146, 0),
		UpdatedAt: time.Unix(1405544146, 0),
	}

	err = dao.Create(o)
	if err != nil {
		t.Errorf("Could not create order object")
	}

	updated := &types.Order{
		ID:              o.ID,
		UserAddress:     o.UserAddress,
		ExchangeAddress: o.ExchangeAddress,
		BuyToken:        o.BuyToken,
		SellToken:       o.SellToken,
		BaseToken:       o.BaseToken,
		QuoteToken:      o.QuoteToken,
		BuyAmount:       big.NewInt(1000),
		SellAmount:      big.NewInt(100),
		Price:           big.NewInt(4000),
		Amount:          big.NewInt(4000),
		FilledAmount:    big.NewInt(200),
		Status:          "FILLED",
		Side:            "BUY",
		PairID:          o.PairID,
		PairName:        "ZRX/WETH",
		Expires:         big.NewInt(10000),
		MakeFee:         big.NewInt(50),
		Nonce:           big.NewInt(1000),
		TakeFee:         big.NewInt(50),
		Signature:       o.Signature,
		Hash:            o.Hash,
		CreatedAt:       o.CreatedAt,
		UpdatedAt:       o.UpdatedAt,
	}

	err = dao.Update(
		o.ID,
		updated,
	)

	if err != nil {
		t.Errorf("Could not updated order from hash %v", err)
	}

	queried, err := dao.GetByHash(o.Hash)
	if err != nil {
		t.Errorf("Could not get order by hash")
	}

	testutils.CompareOrder(t, queried, updated)
}

func TestOrderDao(t *testing.T) {
	dao := NewOrderDao()
	err := dao.Drop()
	if err != nil {
		t.Errorf("Could not drop previous order state")
	}

	o := &types.Order{
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
		Signature: &types.Signature{
			V: 28,
			R: common.HexToHash("0x10b30eb0072a4f0a38b6fca0b731cba15eb2e1702845d97c1230b53a839bcb85"),
			S: common.HexToHash("0x6d9ad89548c9e3ce4c97825d027291477f2c44a8caef792095f2cabc978493ff"),
		},
		Hash:      common.HexToHash("0xb9070a2d333403c255ce71ddf6e795053599b2e885321de40353832b96d8880a"),
		CreatedAt: time.Unix(1405544146, 0),
		UpdatedAt: time.Unix(1405544146, 0),
	}

	err = dao.Create(o)
	if err != nil {
		t.Errorf("Could not create order object")
	}

	o1, err := dao.GetByHash(common.HexToHash("0xb9070a2d333403c255ce71ddf6e795053599b2e885321de40353832b96d8880a"))
	if err != nil {
		t.Errorf("Could not get order by hash")
	}

	testutils.CompareOrder(t, o, o1)

	o2, err := dao.GetByUserAddress(common.HexToAddress("0x7a9f3cd060ab180f36c17fe6bdf9974f577d77aa"))
	if err != nil {
		t.Errorf("Could not get order by user address")
	}

	testutils.CompareOrder(t, o, o2[0])
}
