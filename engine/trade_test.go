package engine

import (
	"testing"
)

func TestExecute(t *testing.T) {
	// e, s := getResource()
	// defer s.Close()

	// // Test Case1: bookEntry amount is less than order amount
	// // New Buy Order
	// // bookEntryJSON := []byte(`{ "id": "5b6ac5297b4457546d64379d", "sellToken": "0x1888a8db0b7db59413ce07150b3373972bf818d3", "buyToken": "0x2034842261b82651885751fc293bba7ba5398156", "baseToken": "0x2034842261b82651885751fc293bba7ba5398156", "quoteToken": "0x1888a8db0b7db59413ce07150b3373972bf818d3", "sellAmount": 6000000000, "buyAmount": 13800000000, "nonce": 0, "userAddress": "0xefD7eB287CeeFCE8256Dd46e25F398acEA7C4b63", "hash": "0xa2d800b77828cb52c83106ca392e465bc0af0d7c319f6956328f739080c19621", "side": "BUY", "amount": 6000000000, "price": 229999999, "filledAmount": 1000000000, "fee": 0, "makeFee": 0, "takeFee": 0, "exchangeAddress": "", "status": "OPEN", "pairID": "5b6ac5117b445753ee755fb8", "pairName": "HPC/AUT", "orderBook": null, "createdAt": "2018-08-08T15:55:45.062141954+05:30", "updatedAt": "2018-08-08T15:55:45.06214208+05:30" }`)
	// bookEntry := &types.Order{
	// 	ID:        bson.ObjectIdHex("5b6ac5297b4457546d64379d"),
	// 	SellToken: common.HexToAddress("0x1888a8db0b7db59413ce07150b3373972bf818d3"),
	// 	BuyToken:  common.HexToAddress("0x2034842261b82651885751fc293bba7ba5398156"),
	// 	PairName:  "HPC/AUT",
	// 	OrderBook: nil,
	// 	CreatedAt: time.Unix(1405544146, 0),
	// 	UpdatedAt: time.Unix(1405544146, 0),
	// }

	// bytes, _ := bookEntry.MarshalJSON()
	// // "2018-08-08T15:55:45.06214208+05:30",

	// // json.Unmarshal(bookEntryJSON, &bookEntry)
	// e.addOrder(bookEntry)

	// // New Sell Order
	// // order := &types.Order{}
	// // orderJSON := []byte(`{ "id": "5b6ac5297b4457546d64379e",
	// // 	"buyToken": "0x1888a8db0b7db59413ce07150b3373972bf818d3",
	// // 	"sellToken": "0x2034842261b82651885751fc293bba7ba5398156",
	// // 	"baseToken": "0x2034842261b82651885751fc293bba7ba5398156",
	// // 	 "quoteToken": "0x1888a8db0b7db59413ce07150b3373972bf818d3",
	// // 		"buyAmount": 6000000000,
	// // 		 "sellAmount": 13800000000,
	// // 			"nonce": 0,
	// // 			 "userAddress": "0xefD7eB287CeeFCE8256Dd46e25F398acEA7C4b63",
	// // 			 "hash": "0xa2d800b77828cb52c83106ca392e465bc0af0d7c319f6956328f739080c19622",
	// // 			 "side": "SELL",
	// // 			 "amount": 6000000000,
	// // 				"price": 229999999,
	// // 				 "filledAmount": 0,
	// // 					"fee": 0,
	// // 					 "makeFee": 0,
	// // 					 "takeFee": 0,
	// // 						"exchangeAddress": "",
	// // 						"status": "NEW",
	// // 						"pairID": "5b6ac5117b445753ee755fb8",
	// // 						"pairName": "HPC/AUT",
	// // 						"orderBook": null,
	// // 						"createdAt": "2018-08-08T15:55:45.062141954+05:30",
	// // 						 "updatedAt": "2018-08-08T15:55:45.06214208+05:30" }`)

	// // json.Unmarshal(orderJSON, &order)

	// order := &types.Order{
	// 	ID:              bson.ObjectIdHex("5b6ac5297b4457546d64379d"),
	// 	UserAddress:     common.HexToAddress("0x7a9f3cd060ab180f36c17fe6bdf9974f577d77aa"),
	// 	ExchangeAddress: common.HexToAddress("0xae55690d4b079460e6ac28aaa58c9ec7b73a7485"),
	// 	SellToken:       common.HexToAddress("0x2034842261b82651885751fc293bba7ba5398156"),
	// 	BuyToken:        common.HexToAddress("0x1888a8db0b7db59413ce07150b3373972bf818d3"),
	// 	BaseToken:       common.HexToAddress("0x2034842261b82651885751fc293bba7ba5398156"),
	// 	QuoteToken:      common.HexToAddress("0x1888a8db0b7db59413ce07150b3373972bf818d3"),
	// 	Amount:          6000000000,
	// 	Price:           229999999,
	// 	BuyAmount:       big.NewInt(6000000000),
	// 	SellAmount:      big.NewInt(13800000000),
	// 	FilledAmount:    0,
	// 	MakeFee:         big.NewInt(0),
	// 	TakeFee:         big.NewInt(0),
	// 	Status:          "NEW",
	// 	PairID:          bson.ObjectIdHex("537f700b537461b70c5f0000"),
	// 	Signature: &types.Signature{
	// 		V: 28,
	// 		R: common.HexToHash("0x10b30eb0072a4f0a38b6fca0b731cba15eb2e1702845d97c1230b53a839bcb85"),
	// 		S: common.HexToHash("0x6d9ad89548c9e3ce4c97825d027291477f2c44a8caef792095f2cabc978493ff"),
	// 	},
	// 	Hash:      common.HexToHash("0xb9070a2d333403c255ce71ddf6e795053599b2e885321de40353832b96d8880a"),
	// 	PairName:  "HPC/AUT",
	// 	OrderBook: nil,
	// 	CreatedAt: time.Unix(1405544146, 0),
	// 	UpdatedAt: time.Unix(1405544146, 0),
	// }

	// orderBytes, _ := bookEntry.MarshalJSON()
	// expectedAmount := bookEntry.Amount.Minus(bookEntry.FilledAmount)

	// expectedTrade := &types.Trade{
	// 	Amount:       expectedAmount,
	// 	Price:        order.Price,
	// 	BaseToken:    order.BaseToken,
	// 	QuoteToken:   order.QuoteToken,
	// 	OrderHash:    bookEntry.Hash,
	// 	Side:         order.Side,
	// 	Taker:        order.UserAddress,
	// 	PairName:     order.PairName,
	// 	Maker:        bookEntry.UserAddress,
	// 	TakerOrderID: order.ID,
	// 	MakerOrderID: bookEntry.ID,
	// }
	// expectedTrade.Hash = expectedTrade.ComputeHash()

	// etb, _ := json.Marshal(expectedTrade)
	// expectedBookEntry := *bookEntry
	// expectedBookEntry.Status = "FILLED"
	// expectedBookEntry.FilledAmount = bookEntry.Amount

	// expectedFillOrder := &FillOrder{
	// 	Amount: bookEntry.Amount - bookEntry.FilledAmount,
	// 	Order:  &expectedBookEntry,
	// }
	// efob, _ := json.Marshal(expectedFillOrder)

	// trade, fillOrder, err := e.execute(order, bookEntry)
	// if err != nil {
	// 	t.Errorf("Error in execute: %s", err)
	// 	return
	// } else {
	// 	tb, _ := json.Marshal(trade)
	// 	fob, _ := json.Marshal(fillOrder)
	// 	fmt.Println(expectedFillOrder.Order.Status == fillOrder.Order.Status)
	// 	assert.JSONEq(t, string(etb), string(tb))
	// 	assert.JSONEq(t, string(efob), string(fob))
	// }

	// // Test Case2: bookEntry amount is equal to order amount
	// // unmarshal bookentry and order from json string
	// json.Unmarshal(bookEntryJSON, &bookEntry)
	// json.Unmarshal(orderJSON, &order)
	// bookEntry.FilledAmount = 0
	// expectedTrade = &types.Trade{
	// 	Amount:       bookEntry.Amount,
	// 	Price:        order.Price,
	// 	BaseToken:    order.BaseToken,
	// 	QuoteToken:   order.QuoteToken,
	// 	OrderHash:    bookEntry.Hash,
	// 	Side:         order.Side,
	// 	Taker:        order.UserAddress,
	// 	PairName:     order.PairName,
	// 	Maker:        bookEntry.UserAddress,
	// 	TakerOrderID: order.ID,
	// 	MakerOrderID: bookEntry.ID,
	// }
	// expectedTrade.Hash = expectedTrade.ComputeHash()

	// etb, _ = json.Marshal(expectedTrade)
	// expectedBookEntry = *bookEntry
	// expectedBookEntry.Status = "FILLED"
	// expectedBookEntry.FilledAmount = bookEntry.Amount

	// expectedFillOrder = &FillOrder{
	// 	Amount: bookEntry.Amount,
	// 	Order:  &expectedBookEntry,
	// }
	// efob, _ = json.Marshal(expectedFillOrder)

	// e.addOrder(bookEntry)

	// trade, fillOrder, err = e.execute(order, bookEntry)
	// if err != nil {
	// 	t.Errorf("Error in execute: %s", err)
	// 	return
	// } else {
	// 	tb, _ := json.Marshal(trade)
	// 	fob, _ := json.Marshal(fillOrder)
	// 	fmt.Println(expectedFillOrder.Order.Status == fillOrder.Order.Status)
	// 	assert.JSONEq(t, string(etb), string(tb))
	// 	assert.JSONEq(t, string(efob), string(fob))
	// }

	// // Test Case3: bookEntry amount is greater then order amount
	// // unmarshal bookentry and order from json string
	// json.Unmarshal(bookEntryJSON, &bookEntry)
	// json.Unmarshal(orderJSON, &order)
	// bookEntry.Amount = bookEntry.Amount + bookEntry.FilledAmount
	// bookEntry.FilledAmount = 0
	// expectedTrade = &types.Trade{
	// 	Amount:       order.Amount,
	// 	Price:        order.Price,
	// 	BaseToken:    order.BaseToken,
	// 	QuoteToken:   order.QuoteToken,
	// 	OrderHash:    bookEntry.Hash,
	// 	Side:         order.Side,
	// 	Taker:        order.UserAddress,
	// 	PairName:     order.PairName,
	// 	Maker:        bookEntry.UserAddress,
	// 	TakerOrderID: order.ID,
	// 	MakerOrderID: bookEntry.ID,
	// }
	// expectedTrade.Hash = expectedTrade.ComputeHash()

	// etb, _ = json.Marshal(expectedTrade)
	// expectedBookEntry = *bookEntry
	// expectedBookEntry.Status = types.PARTIALFILLED
	// expectedBookEntry.FilledAmount = expectedBookEntry.FilledAmount + order.Amount

	// expectedFillOrder = &FillOrder{
	// 	Amount: order.Amount,
	// 	Order:  &expectedBookEntry,
	// }

	// efob, _ = json.Marshal(expectedFillOrder)
	// e.addOrder(bookEntry)

	// trade, fillOrder, err = e.execute(order, bookEntry)
	// if err != nil {
	// 	t.Errorf("Error in execute: %s", err)
	// 	return
	// } else {
	// 	tb, _ := json.Marshal(trade)
	// 	fob, _ := json.Marshal(fillOrder)
	// 	fmt.Println(expectedFillOrder.Order.Status == fillOrder.Order.Status)
	// 	assert.JSONEq(t, string(etb), string(tb))
	// 	assert.JSONEq(t, string(efob), string(fob))
	// }
}
