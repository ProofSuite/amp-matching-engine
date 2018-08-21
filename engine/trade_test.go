package engine

import (
	"encoding/json"
	"math/big"
	"testing"
	"time"

	"github.com/Proofsuite/amp-matching-engine/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"gopkg.in/mgo.v2/bson"
)

func TestExecute(t *testing.T) {
	e := getResource()
	defer flushData(e.redisConn)
	// Test Case1: bookEntry amount is less than order amount
	// New Buy Order
	bookEntry := &types.Order{
		ID:              bson.ObjectIdHex("5b6ac5297b4457546d64379d"),
		SellToken:       common.HexToAddress("0x1888a8db0b7db59413ce07150b3373972bf818d3"),
		BuyToken:        common.HexToAddress("0x2034842261b82651885751fc293bba7ba5398156"),
		BaseToken:       common.HexToAddress("0x2034842261b82651885751fc293bba7ba5398156"),
		QuoteToken:      common.HexToAddress("0x1888a8db0b7db59413ce07150b3373972bf818d3"),
		SellAmount:      big.NewInt(6000000000),
		BuyAmount:       big.NewInt(13800000000),
		Nonce:           big.NewInt(0),
		Expires:         big.NewInt(0),
		Side:            "BUY",
		Amount:          6000000000,
		Price:           229999999,
		FilledAmount:    1000000000,
		MakeFee:         big.NewInt(0),
		TakeFee:         big.NewInt(0),
		ExchangeAddress: common.HexToAddress("0xae55690d4b079460e6ac28aaa58c9ec7b73a7485"),
		UserAddress:     common.HexToAddress("0x7a9f3cd060ab180f36c17fe6bdf9974f577d77aa"),
		Hash:            common.HexToHash("0xa2d800b77828cb52c83106ca392e465bc0af0d7c319f6956328f739080c19621"),
		PairID:          bson.ObjectIdHex("537f700b537461b70c5f0000"),
		PairName:        "HPC/AUT",
		OrderBook:       nil,
		CreatedAt:       time.Unix(1405544146, 0),
		UpdatedAt:       time.Unix(1405544146, 0),
	}
	bookEntryJSON, _ := json.Marshal(bookEntry)
	// bytes, _ := bookEntry.MarshalJSON()

	// json.Unmarshal(bookEntryJSON, &bookEntry)
	e.addOrder(bookEntry)

	// json.Unmarshal(orderJSON, &order)

	order := &types.Order{
		ID:              bson.ObjectIdHex("5b6ac5297b4457546d64379d"),
		UserAddress:     common.HexToAddress("0x7a9f3cd060ab180f36c17fe6bdf9974f577d77aa"),
		ExchangeAddress: common.HexToAddress("0xae55690d4b079460e6ac28aaa58c9ec7b73a7485"),
		SellToken:       common.HexToAddress("0x2034842261b82651885751fc293bba7ba5398156"),
		BuyToken:        common.HexToAddress("0x1888a8db0b7db59413ce07150b3373972bf818d3"),
		BaseToken:       common.HexToAddress("0x2034842261b82651885751fc293bba7ba5398156"),
		QuoteToken:      common.HexToAddress("0x1888a8db0b7db59413ce07150b3373972bf818d3"),
		Amount:          6000000000,
		Price:           229999999,
		BuyAmount:       big.NewInt(6000000000),
		SellAmount:      big.NewInt(13800000000),
		Side:            "SELL",
		FilledAmount:    0,
		MakeFee:         big.NewInt(0),
		TakeFee:         big.NewInt(0),
		Status:          "NEW",
		PairID:          bson.ObjectIdHex("537f700b537461b70c5f0000"),
		Signature: &types.Signature{
			V: 28,
			R: common.HexToHash("0x10b30eb0072a4f0a38b6fca0b731cba15eb2e1702845d97c1230b53a839bcb85"),
			S: common.HexToHash("0x6d9ad89548c9e3ce4c97825d027291477f2c44a8caef792095f2cabc978493ff"),
		},
		Hash:      common.HexToHash("0xb9070a2d333403c255ce71ddf6e795053599b2e885321de40353832b96d8880a"),
		PairName:  "HPC/AUT",
		OrderBook: nil,
		CreatedAt: time.Unix(1405544146, 0),
		UpdatedAt: time.Unix(1405544146, 0),
	}
	orderJSON, _ := json.Marshal(order)

	// orderBytes, _ := bookEntry.MarshalJSON()
	expectedAmount := bookEntry.Amount - bookEntry.FilledAmount

	expectedTrade := &types.Trade{
		Amount:       big.NewInt(expectedAmount),
		Price:        order.Price,
		BaseToken:    order.BaseToken,
		QuoteToken:   order.QuoteToken,
		OrderHash:    bookEntry.Hash,
		Side:         order.Side,
		Taker:        order.UserAddress,
		PairName:     order.PairName,
		Maker:        bookEntry.UserAddress,
		TradeNonce:   big.NewInt(0),
		TakerOrderID: order.ID,
		MakerOrderID: bookEntry.ID,
		Signature:    &types.Signature{},
	}

	expectedTrade.Hash = expectedTrade.ComputeHash()

	etb, _ := json.Marshal(expectedTrade)
	expectedBookEntry := *bookEntry
	expectedBookEntry.Status = "FILLED"
	expectedBookEntry.FilledAmount = bookEntry.Amount

	expectedFillOrder := &FillOrder{
		Amount: bookEntry.Amount - bookEntry.FilledAmount,
		Order:  &expectedBookEntry,
	}
	efob, _ := json.Marshal(expectedFillOrder)

	trade, fillOrder, err := e.execute(order, bookEntry)
	if err != nil {
		t.Errorf("Error in execute: %s", err)
		return
	} else {
		tb, _ := json.Marshal(trade)
		fob, _ := json.Marshal(fillOrder)
		assert.JSONEq(t, string(etb), string(tb))
		assert.JSONEq(t, string(efob), string(fob))
	}

	// Test Case2: bookEntry amount is equal to order amount
	// unmarshal bookentry and order from json string
	json.Unmarshal(bookEntryJSON, &bookEntry)
	json.Unmarshal(orderJSON, &order)

	bookEntry.FilledAmount = 0
	expectedAmount = bookEntry.Amount - bookEntry.FilledAmount
	expectedTrade = &types.Trade{
		Amount:       big.NewInt(expectedAmount),
		Price:        order.Price,
		BaseToken:    order.BaseToken,
		QuoteToken:   order.QuoteToken,
		OrderHash:    bookEntry.Hash,
		Side:         order.Side,
		Taker:        order.UserAddress,
		PairName:     order.PairName,
		Maker:        bookEntry.UserAddress,
		TakerOrderID: order.ID,
		MakerOrderID: bookEntry.ID,
		TradeNonce:   big.NewInt(0),
		Signature:    &types.Signature{},
	}
	expectedTrade.Hash = expectedTrade.ComputeHash()

	etb, _ = json.Marshal(expectedTrade)
	expectedBookEntry = *bookEntry
	expectedBookEntry.Status = "FILLED"
	expectedBookEntry.FilledAmount = bookEntry.Amount

	expectedFillOrder = &FillOrder{
		Amount: bookEntry.Amount,
		Order:  &expectedBookEntry,
	}

	efob, _ = json.Marshal(expectedFillOrder)

	e.addOrder(bookEntry)

	trade, fillOrder, err = e.execute(order, bookEntry)
	if err != nil {
		t.Errorf("Error in execute: %s", err)
		return
	} else {
		tb, _ := json.Marshal(trade)
		fob, _ := json.Marshal(fillOrder)
		assert.JSONEq(t, string(etb), string(tb))
		assert.JSONEq(t, string(efob), string(fob))
	}

	// Test Case3: bookEntry amount is greater then order amount
	// unmarshal bookentry and order from json string
	json.Unmarshal(bookEntryJSON, &bookEntry)
	json.Unmarshal(orderJSON, &order)
	bookEntry.Amount = bookEntry.Amount + bookEntry.FilledAmount
	bookEntry.FilledAmount = 0
	expectedAmount = order.Amount
	expectedTrade = &types.Trade{
		Amount:       big.NewInt(expectedAmount),
		Price:        order.Price,
		BaseToken:    order.BaseToken,
		QuoteToken:   order.QuoteToken,
		OrderHash:    bookEntry.Hash,
		Side:         order.Side,
		Taker:        order.UserAddress,
		PairName:     order.PairName,
		Maker:        bookEntry.UserAddress,
		TakerOrderID: order.ID,
		MakerOrderID: bookEntry.ID,
		TradeNonce:   big.NewInt(0),
		Signature:    &types.Signature{},
	}
	expectedTrade.Hash = expectedTrade.ComputeHash()

	etb, _ = json.Marshal(expectedTrade)
	expectedBookEntry = *bookEntry
	expectedBookEntry.Status = "PARTIAL_FILLED"
	expectedBookEntry.FilledAmount = expectedBookEntry.FilledAmount + order.Amount

	expectedFillOrder = &FillOrder{
		Amount: order.Amount,
		Order:  &expectedBookEntry,
	}

	efob, _ = json.Marshal(expectedFillOrder)
	e.addOrder(bookEntry)

	trade, fillOrder, err = e.execute(order, bookEntry)
	if err != nil {
		t.Errorf("Error in execute: %s", err)
		return
	} else {
		tb, _ := json.Marshal(trade)
		fob, _ := json.Marshal(fillOrder)
		assert.JSONEq(t, string(etb), string(tb))
		assert.JSONEq(t, string(efob), string(fob))
	}
}
