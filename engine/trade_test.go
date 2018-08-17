package engine

import (
	"github.com/Proofsuite/amp-matching-engine/types"
	"github.com/stretchr/testify/assert"
	"testing"
	"encoding/json"
	"fmt"
)

func TestExecute(t *testing.T) {
	e, s := getResource()
	defer s.Close()

	// Test Case1: bookEntry amount is less than order amount
	// New Buy Order
	bookEntry := &types.Order{}
	bookEntryJSON := []byte(`{ "id": "5b6ac5297b4457546d64379d", "sellToken": "0x1888a8db0b7db59413ce07150b3373972bf818d3", "buyToken": "0x2034842261b82651885751fc293bba7ba5398156", "baseToken": "0x2034842261b82651885751fc293bba7ba5398156", "quoteToken": "0x1888a8db0b7db59413ce07150b3373972bf818d3", "sellAmount": 6000000000, "buyAmount": 13800000000, "nonce": 0, "userAddress": "0xefD7eB287CeeFCE8256Dd46e25F398acEA7C4b63", "hash": "0xa2d800b77828cb52c83106ca392e465bc0af0d7c319f6956328f739080c19621", "side": "BUY", "amount": 6000000000, "price": 229999999, "filledAmount": 1000000000, "fee": 0, "makeFee": 0, "takeFee": 0, "exchangeAddress": "", "status": "OPEN", "pairID": "5b6ac5117b445753ee755fb8", "pairName": "HPC/AUT", "orderBook": null, "createdAt": "2018-08-08T15:55:45.062141954+05:30", "updatedAt": "2018-08-08T15:55:45.06214208+05:30" }`)
	json.Unmarshal(bookEntryJSON, &bookEntry)
	e.addOrder(bookEntry)

	// New Sell Order
	order := &types.Order{}
	orderJSON := []byte(`{ "id": "5b6ac5297b4457546d64379e", "buyToken": "0x1888a8db0b7db59413ce07150b3373972bf818d3", "sellToken": "0x2034842261b82651885751fc293bba7ba5398156", "baseToken": "0x2034842261b82651885751fc293bba7ba5398156", "quoteToken": "0x1888a8db0b7db59413ce07150b3373972bf818d3", "buyAmount": 6000000000, "sellAmount": 13800000000, "nonce": 0, "userAddress": "0xefD7eB287CeeFCE8256Dd46e25F398acEA7C4b63", "hash": "0xa2d800b77828cb52c83106ca392e465bc0af0d7c319f6956328f739080c19622", "side": "SELL", "amount": 6000000000, "price": 229999999, "filledAmount": 0, "fee": 0, "makeFee": 0, "takeFee": 0, "exchangeAddress": "", "status": "NEW", "pairID": "5b6ac5117b445753ee755fb8", "pairName": "HPC/AUT", "orderBook": null, "createdAt": "2018-08-08T15:55:45.062141954+05:30", "updatedAt": "2018-08-08T15:55:45.06214208+05:30" }`)
	json.Unmarshal(orderJSON, &order)

	expectedTrade := &types.Trade{
		Amount:       bookEntry.Amount-bookEntry.FilledAmount,
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
	}
	expectedTrade.Hash = expectedTrade.ComputeHash()

	etb,_:=json.Marshal(expectedTrade)
	expectedBookEntry:=*bookEntry
	expectedBookEntry.Status=types.FILLED
	expectedBookEntry.FilledAmount=bookEntry.Amount

	expectedFillOrder := &FillOrder{
		Amount: bookEntry.Amount-bookEntry.FilledAmount,
		Order:  &expectedBookEntry,
	}
	efob,_:=json.Marshal(expectedFillOrder)

	trade, fillOrder, err := e.execute(order, bookEntry)
	if err != nil {
		t.Errorf("Error in execute: %s", err)
		return
	} else{
		tb,_:=json.Marshal(trade)
		fob,_:=json.Marshal(fillOrder)
		fmt.Println(expectedFillOrder.Order.Status == fillOrder.Order.Status)
		assert.JSONEq(t, string(etb),  string(tb))
		assert.JSONEq(t, string(efob),  string(fob))
	}

	// Test Case2: bookEntry amount is equal to order amount
	// unmarshal bookentry and order from json string
	json.Unmarshal(bookEntryJSON, &bookEntry)
	json.Unmarshal(orderJSON, &order)
	bookEntry.FilledAmount=0
	expectedTrade = &types.Trade{
		Amount:       bookEntry.Amount,
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
	}
	expectedTrade.Hash = expectedTrade.ComputeHash()

	etb,_=json.Marshal(expectedTrade)
	expectedBookEntry=*bookEntry
	expectedBookEntry.Status=types.FILLED
	expectedBookEntry.FilledAmount=bookEntry.Amount

	expectedFillOrder = &FillOrder{
		Amount: bookEntry.Amount,
		Order:  &expectedBookEntry,
	}
	efob,_=json.Marshal(expectedFillOrder)

	e.addOrder(bookEntry)

	trade, fillOrder, err = e.execute(order, bookEntry)
	if err != nil {
		t.Errorf("Error in execute: %s", err)
		return
	} else{
		tb,_:=json.Marshal(trade)
		fob,_:=json.Marshal(fillOrder)
		fmt.Println(expectedFillOrder.Order.Status == fillOrder.Order.Status)
		assert.JSONEq(t, string(etb),  string(tb))
		assert.JSONEq(t, string(efob),  string(fob))
	}

	// Test Case3: bookEntry amount is greater then order amount
	// unmarshal bookentry and order from json string
	json.Unmarshal(bookEntryJSON, &bookEntry)
	json.Unmarshal(orderJSON, &order)
	bookEntry.Amount=bookEntry.Amount+bookEntry.FilledAmount
	bookEntry.FilledAmount=0
	expectedTrade = &types.Trade{
		Amount:       order.Amount,
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
	}
	expectedTrade.Hash = expectedTrade.ComputeHash()

	etb,_=json.Marshal(expectedTrade)
	expectedBookEntry=*bookEntry
	expectedBookEntry.Status=types.PARTIALFILLED
	expectedBookEntry.FilledAmount=expectedBookEntry.FilledAmount+order.Amount

	expectedFillOrder = &FillOrder{
		Amount: order.Amount,
		Order:  &expectedBookEntry,
	}

	efob,_=json.Marshal(expectedFillOrder)
	e.addOrder(bookEntry)

	trade, fillOrder, err = e.execute(order, bookEntry)
	if err != nil {
		t.Errorf("Error in execute: %s", err)
		return
	} else{
		tb,_:=json.Marshal(trade)
		fob,_:=json.Marshal(fillOrder)
		fmt.Println(expectedFillOrder.Order.Status == fillOrder.Order.Status)
		assert.JSONEq(t, string(etb),  string(tb))
		assert.JSONEq(t, string(efob),  string(fob))
	}
}
