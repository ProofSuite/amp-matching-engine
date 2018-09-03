package engine

// import (
// 	"encoding/json"
// 	"math/big"
// 	"testing"

// 	"github.com/Proofsuite/amp-matching-engine/types"
// 	"github.com/Proofsuite/amp-matching-engine/utils/math"
// 	"github.com/stretchr/testify/assert"
// )

// func TestExecute(t *testing.T) {
// 	e := getResource()
// 	defer e.redisConn.FlushAll()
// 	// Test Case1: bookEntry amount is less than order amount
// 	// New Buy Order
// 	bookEntry := getBuyOrder()
// 	bookEntry.FilledAmount = big.NewInt(1000000000)

// 	e.addOrder(&bookEntry)

// 	order := getSellOrder()

// 	expectedAmount := math.Sub(bookEntry.Amount, bookEntry.FilledAmount)

// 	expectedTrade := getTrade(&order, &bookEntry, expectedAmount, big.NewInt(0))

// 	expectedTrade.Hash = expectedTrade.ComputeHash()

// 	etb, _ := json.Marshal(expectedTrade)
// 	expectedBookEntry := bookEntry
// 	expectedBookEntry.Status = "FILLED"
// 	expectedBookEntry.FilledAmount = bookEntry.Amount

// 	expectedFillOrder := &types.FillOrder{
// 		Amount: math.Sub(bookEntry.Amount, bookEntry.FilledAmount),
// 		Order:  &expectedBookEntry,
// 	}
// 	efob, _ := json.Marshal(expectedFillOrder)

// 	trade, fillOrder, err := e.execute(&order, &bookEntry)
// 	if err != nil {
// 		t.Errorf("Error in execute: %s", err)
// 		return
// 	}
// 	tb, _ := json.Marshal(trade)
// 	fob, _ := json.Marshal(fillOrder)
// 	assert.JSONEq(t, string(etb), string(tb))
// 	assert.JSONEq(t, string(efob), string(fob))

// 	// Test Case2: bookEntry amount is equal to order amount
// 	// unmarshal bookentry and order from json string
// 	bookEntry = getBuyOrder()
// 	order = getSellOrder()
// 	bookEntry.FilledAmount = big.NewInt(0)
// 	expectedAmount = math.Sub(bookEntry.Amount, bookEntry.FilledAmount)
// 	expectedTrade = getTrade(&order, &bookEntry, expectedAmount, big.NewInt(0))
// 	expectedTrade.Hash = expectedTrade.ComputeHash()

// 	etb, _ = json.Marshal(expectedTrade)
// 	expectedBookEntry = bookEntry
// 	expectedBookEntry.Status = "FILLED"
// 	expectedBookEntry.FilledAmount = bookEntry.Amount

// 	expectedFillOrder = &types.FillOrder{
// 		Amount: bookEntry.Amount,
// 		Order:  &expectedBookEntry,
// 	}

// 	efob, _ = json.Marshal(expectedFillOrder)

// 	e.addOrder(&bookEntry)

// 	trade, fillOrder, err = e.execute(&order, &bookEntry)
// 	if err != nil {
// 		t.Errorf("Error in execute: %s", err)
// 		return
// 	} else {
// 		tb, _ := json.Marshal(trade)
// 		fob, _ := json.Marshal(fillOrder)
// 		assert.JSONEq(t, string(etb), string(tb))
// 		assert.JSONEq(t, string(efob), string(fob))
// 	}

// 	// Test Case3: bookEntry amount is greater then order amount
// 	// unmarshal bookentry and order from json string
// 	bookEntry = getBuyOrder()
// 	order = getSellOrder()
// 	bookEntry.Amount = math.Add(bookEntry.Amount, big.NewInt(1000000000))
// 	expectedAmount = order.Amount
// 	expectedTrade = getTrade(&order, &bookEntry, expectedAmount, big.NewInt(0))
// 	expectedTrade.Hash = expectedTrade.ComputeHash()

// 	etb, _ = json.Marshal(expectedTrade)
// 	expectedBookEntry = bookEntry
// 	expectedBookEntry.Status = "PARTIAL_FILLED"
// 	expectedBookEntry.FilledAmount = math.Add(expectedBookEntry.FilledAmount, order.Amount)

// 	expectedFillOrder = &types.FillOrder{
// 		Amount: order.Amount,
// 		Order:  &expectedBookEntry,
// 	}

// 	efob, _ = json.Marshal(expectedFillOrder)
// 	e.addOrder(&bookEntry)

// 	trade, fillOrder, err = e.execute(&order, &bookEntry)
// 	if err != nil {
// 		t.Errorf("Error in execute: %s", err)
// 		return
// 	}
// 	tb, _ = json.Marshal(trade)
// 	fob, _ = json.Marshal(fillOrder)
// 	assert.JSONEq(t, string(etb), string(tb))
// 	assert.JSONEq(t, string(efob), string(fob))
// }
