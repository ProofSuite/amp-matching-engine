package types

// func Compare(t *testing.T, expected interface{}, value interface{}) {
// 	expectedBytes, _ := json.Marshal(expected)
// 	bytes, _ := json.Marshal(value)

// 	assert.JSONEqf(t, string(expectedBytes), string(bytes), "")
// }

// func CompareStructs(t *testing.T, expected interface{}, order interface{}) {
// 	diff := deep.Equal(expected, order)
// 	if diff != nil {
// 		t.Errorf("\n%+v\nGot: \n%+v\n\n", expected, order)
// 	}
// }

// func TestWebSocketMessageJSON(t *testing.T) {
// 	order := &Order{
// 		ID:              bson.ObjectIdHex("537f700b537461b70c5f0000"),
// 		UserAddress:     common.HexToAddress("0x7a9f3cd060ab180f36c17fe6bdf9974f577d77aa"),
// 		ExchangeAddress: common.HexToAddress("0xae55690d4b079460e6ac28aaa58c9ec7b73a7485"),
// 		BuyToken:        common.HexToAddress("0xe41d2489571d322189246dafa5ebde1f4699f498"),
// 		SellToken:       common.HexToAddress("0x12459c951127e0c374ff9105dda097662a027093"),
// 		BaseToken:       common.HexToAddress("0xe41d2489571d322189246dafa5ebde1f4699f498"),
// 		QuoteToken:      common.HexToAddress("0x12459c951127e0c374ff9105dda097662a027093"),
// 		BuyAmount:       big.NewInt(1000),
// 		SellAmount:      big.NewInt(100),
// 		Amount:          big.NewInt(1000),
// 		FilledAmount:    big.NewInt(100),
// 		Status:          "OPEN",
// 		Side:            "BUY",
// 		PairName:        "ZRX/WETH",
// 		MakeFee:         big.NewInt(50),
// 		Nonce:           big.NewInt(1000),
// 		TakeFee:         big.NewInt(50),
// 		Signature: &Signature{
// 			V: 28,
// 			R: common.HexToHash("0x10b30eb0072a4f0a38b6fca0b731cba15eb2e1702845d97c1230b53a839bcb85"),
// 			S: common.HexToHash("0x6d9ad89548c9e3ce4c97825d027291477f2c44a8caef792095f2cabc978493ff"),
// 		},
// 		Hash:      common.HexToHash("0xb9070a2d333403c255ce71ddf6e795053599b2e885321de40353832b96d8880a"),
// 		CreatedAt: time.Unix(1405544146, 0),
// 		UpdatedAt: time.Unix(1405544146, 0),
// 	}

// 	msg := &WebSocketMessage{
// 		Channel: "orders",
// 		Payload: WebSocketPayload{
// 			Type: "NEW_ORDER",
// 			Hash: "0xb9070a2d333403c255ce71ddf6e795053599b2e885321de40353832b96d8880a",
// 			Data: order,
// 		},
// 	}

// 	encoded, err := json.Marshal(msg)
// 	if err != nil {
// 		t.Errorf("Error encoding order: %v", err)
// 	}

// 	decoded := &WebSocketMessage{}
// 	err = json.Unmarshal([]byte(encoded), &decoded)
// 	if err != nil {
// 		t.Errorf("Could not unmarshal payload: %v", err)
// 	}

// 	//we re-encode the order because the order is unmarshalled as a map[string]interface
// 	encodedOrder, err := json.Marshal(decoded.Payload.Data)
// 	if err != nil {
// 		t.Errorf("Error encoding order: %v", err)
// 	}

// 	decodedOrder := &Order{}
// 	err = json.Unmarshal(encodedOrder, decodedOrder)
// 	if err != nil {
// 		t.Errorf("Error decoding order; %v", err)
// 	}

// 	testutils.Compare(t, msg, decoded)
// 	testutils.CompareStructs(t, order, decodedOrder)
// }

// func TestOrderCancelWebSocketMessageJSON(t *testing.T) {
// 	oc := &OrderCancel{
// 		Signature: &Signature{
// 			V: 28,
// 			R: common.HexToHash("0x10b30eb0072a4f0a38b6fca0b731cba15eb2e1702845d97c1230b53a839bcb85"),
// 			S: common.HexToHash("0x6d9ad89548c9e3ce4c97825d027291477f2c44a8caef792095f2cabc978493ff"),
// 		},
// 		Hash:      common.HexToHash("0xb9070a2d333403c255ce71ddf6e795053599b2e885321de40353832b96d8880a"),
// 		OrderHash: common.HexToHash("0xb9070a2d333403c255ce71ddf6e795053599b2e885321de40353832b96d8880a"),
// 	}

// 	msg := &WebSocketMessage{
// 		Channel: "orders",
// 		Payload: WebSocketPayload{
// 			Type: "CANCEL",
// 			Hash: "0xb9070a2d333403c255ce71ddf6e795053599b2e885321de40353832b96d8880a",
// 			Data: oc,
// 		},
// 	}

// 	encoded, err := json.Marshal(msg)
// 	if err != nil {
// 		t.Errorf("Error encoding order: %v", err)
// 	}

// 	decoded := &WebSocketMessage{}
// 	err = json.Unmarshal([]byte(encoded), &decoded)
// 	if err != nil {
// 		t.Errorf("Could not unmarshal payload: %v", err)
// 	}

// 	//we re-encode the order because the order is unmarshalled as a map[string]interface
// 	encodedOrderCancel, err := json.Marshal(decoded.Payload.Data)
// 	if err != nil {
// 		t.Errorf("Error encoding order: %v", err)
// 	}

// 	decodedOrderCancel := &OrderCancel{}
// 	err = json.Unmarshal(encodedOrderCancel, decodedOrderCancel)
// 	if err != nil {
// 		t.Errorf("Error decoding order; %v", err)
// 	}

// 	Compare(t, msg, decoded)
// 	CompareStructs(t, oc, decodedOrderCancel)
// }

// func TestNewOrderCancelWebSocketMessage(t *testing.T) {
// 	oc := &OrderCancel{
// 		Signature: &Signature{
// 			V: 28,
// 			R: common.HexToHash("0x10b30eb0072a4f0a38b6fca0b731cba15eb2e1702845d97c1230b53a839bcb85"),
// 			S: common.HexToHash("0x6d9ad89548c9e3ce4c97825d027291477f2c44a8caef792095f2cabc978493ff"),
// 		},
// 		Hash:      common.HexToHash("0xb9070a2d333403c255ce71ddf6e795053599b2e885321de40353832b96d8880a"),
// 		OrderHash: common.HexToHash("0xb9070a2d333403c255ce71ddf6e795053599b2e885321de40353832b96d8880a"),
// 	}

// 	msg := NewOrderCancelWebsocketMessage(oc)

// 	expected := &WebSocketMessage{
// 		Channel: "orders",
// 		Payload: WebSocketPayload{
// 			Type: "CANCEL_ORDER",
// 			Hash: "0xb9070a2d333403c255ce71ddf6e795053599b2e885321de40353832b96d8880a",
// 			Data: &OrderCancel{
// 				Signature: &Signature{
// 					V: 28,
// 					R: common.HexToHash("0x10b30eb0072a4f0a38b6fca0b731cba15eb2e1702845d97c1230b53a839bcb85"),
// 					S: common.HexToHash("0x6d9ad89548c9e3ce4c97825d027291477f2c44a8caef792095f2cabc978493ff"),
// 				},
// 				Hash:      common.HexToHash("0xb9070a2d333403c255ce71ddf6e795053599b2e885321de40353832b96d8880a"),
// 				OrderHash: common.HexToHash("0xb9070a2d333403c255ce71ddf6e795053599b2e885321de40353832b96d8880a"),
// 			},
// 		},
// 	}

// 	testutils.Compare(t, expected, msg)
// 	testutils.CompareStructs(t, expected, msg)
// }

// func TestNewWebsocketMessage(t *testing.T) {
// 	o := &Order{
// 		ID:              bson.ObjectIdHex("537f700b537461b70c5f0000"),
// 		UserAddress:     common.HexToAddress("0x7a9f3cd060ab180f36c17fe6bdf9974f577d77aa"),
// 		ExchangeAddress: common.HexToAddress("0xae55690d4b079460e6ac28aaa58c9ec7b73a7485"),
// 		BuyToken:        common.HexToAddress("0xe41d2489571d322189246dafa5ebde1f4699f498"),
// 		SellToken:       common.HexToAddress("0x12459c951127e0c374ff9105dda097662a027093"),
// 		BaseToken:       common.HexToAddress("0xe41d2489571d322189246dafa5ebde1f4699f498"),
// 		QuoteToken:      common.HexToAddress("0x12459c951127e0c374ff9105dda097662a027093"),
// 		BuyAmount:       big.NewInt(1000),
// 		SellAmount:      big.NewInt(100),
// 		PricePoint:      big.NewInt(1000000),
// 		Amount:          big.NewInt(1000),
// 		FilledAmount:    big.NewInt(100),
// 		Status:          "OPEN",
// 		Side:            "BUY",
// 		PairName:        "ZRX/WETH",
// 		MakeFee:         big.NewInt(50),
// 		Nonce:           big.NewInt(1000),
// 		TakeFee:         big.NewInt(50),
// 		Signature: &Signature{
// 			V: 28,
// 			R: common.HexToHash("0x10b30eb0072a4f0a38b6fca0b731cba15eb2e1702845d97c1230b53a839bcb85"),
// 			S: common.HexToHash("0x6d9ad89548c9e3ce4c97825d027291477f2c44a8caef792095f2cabc978493ff"),
// 		},
// 		Hash:      common.HexToHash("0xb9070a2d333403c255ce71ddf6e795053599b2e885321de40353832b96d8880a"),
// 		CreatedAt: time.Unix(1405544146, 0),
// 		UpdatedAt: time.Unix(1405544146, 0),
// 	}

// 	msg := NewOrderWebsocketMessage(o)

// 	expected := &WebSocketMessage{
// 		Channel: "orders",
// 		Payload: WebSocketPayload{
// 			Type: "NEW_ORDER",
// 			Hash: "0xb9070a2d333403c255ce71ddf6e795053599b2e885321de40353832b96d8880a",
// 			Data: o,
// 		},
// 	}

// 	testutils.Compare(t, expected, msg)
// 	testutils.CompareStructs(t, expected, msg)
// }
