package services

//ORDERBOOK.GO

// GetOrderBook fetches orderbook from engine/redis and returns it as an map[string]interface
// func (s *OrderBookService) GetOrderBook(bt, qt common.Address) (map[string]interface{}, error) {
// 	res, err := s.pairDao.GetByTokenAddress(bt, qt)
// 	if err != nil {
// 		message := map[string]string{
// 			"Code":    "Invalid Pair",
// 			"Message": err.Error(),
// 		}
// 		bytes, _ := json.Marshal(message)
// 		return nil, errors.New(string(bytes))
// 	}

// 	bids, asks, err := s.eng.GetOrderBook(res)
// 	if err != nil {
// 		message := map[string]string{
// 			"Code":    "Invalid Pair",
// 			"Message": err.Error(),
// 		}

// 		bytes, _ := json.Marshal(message)
// 		return nil, errors.New(string(bytes))
// 	}

// 	ob := map[string]interface{}{
// 		"asks": asks,
// 		"bids": bids,
// 	}
// 	return ob, nil
// }

// // SubscribeOrderBook is responsible for handling incoming orderbook subscription messages
// // It makes an entry of connection in pairSocket corresponding to pair,unit and duration
// func (s *OrderBookService) SubscribeOrderBook(conn *ws.Conn, bt, qt common.Address) {
// 	socket := ws.GetOrderBookSocket()

// 	ob, err := s.GetOrderBook(bt, qt)
// 	if err != nil {
// 		socket.SendErrorMessage(conn, err.Error())
// 		return
// 	}

// 	id := utils.GetOrderBookChannelID(bt, qt)
// 	err = socket.Subscribe(id, conn)
// 	if err != nil {
// 		message := map[string]string{
// 			"Code":    "Internal Server Error",
// 			"Message": err.Error(),
// 		}

// 		socket.SendErrorMessage(conn, message)
// 		return
// 	}

// 	ws.RegisterConnectionUnsubscribeHandler(conn, socket.UnsubscribeHandler(id))
// 	socket.SendInitMessage(conn, ob)
// }

// // UnsubscribeOrderBook is responsible for handling incoming orderbook unsubscription messages
// func (s *OrderBookService) UnsubscribeOrderBook(conn *ws.Conn, bt, qt common.Address) {
// 	socket := ws.GetOrderBookSocket()

// 	id := utils.GetOrderBookChannelID(bt, qt)
// 	socket.Unsubscribe(id, conn)
// }

// // GetRawOrderBook fetches complete orderbook from engine/redis
// func (s *OrderBookService) GetRawOrderBook(bt, qt common.Address) ([][]types.Order, error) {
// 	res, err := s.pairDao.GetByTokenAddress(bt, qt)
// 	if err != nil {
// 		message := map[string]string{
// 			"Code":    "Invalid Pair",
// 			"Message": err.Error(),
// 		}
// 		bytes, _ := json.Marshal(message)
// 		return nil, errors.New(string(bytes))
// 	}

// 	book, err := s.eng.GetRawOrderBook(res)
// 	if err != nil {
// 		message := map[string]string{
// 			"Code":    "Internal Server Error",
// 			"Message": err.Error(),
// 		}
// 		bytes, _ := json.Marshal(message)
// 		return nil, errors.New(string(bytes))
// 	}

// 	return book, nil
// }

// // SubscribeRawOrderBook is responsible for handling incoming orderbook subscription messages
// // It makes an entry of connection in pairSocket corresponding to pair,unit and duration
// func (s *OrderBookService) SubscribeRawOrderBook(conn *ws.Conn, bt, qt common.Address) {
// 	socket := ws.GetRawOrderBookSocket()

// 	ob, err := s.GetRawOrderBook(bt, qt)
// 	if err != nil {
// 		socket.SendErrorMessage(conn, err.Error())
// 		return
// 	}

// 	id := utils.GetOrderBookChannelID(bt, qt)
// 	err = socket.Subscribe(id, conn)
// 	if err != nil {
// 		message := map[string]string{
// 			"Code":    "Internal Server Error",
// 			"Message": err.Error(),
// 		}

// 		socket.SendErrorMessage(conn, message)
// 		return
// 	}

// 	ws.RegisterConnectionUnsubscribeHandler(conn, socket.UnsubscribeHandler(id))
// 	socket.SendInitMessage(conn, ob)
// }

// // UnsubscribeRawOrderBook is responsible for handling incoming orderbook unsubscription messages
// func (s *OrderBookService) UnsubscribeRawOrderBook(conn *ws.Conn, bt, qt common.Address) {
// 	socket := ws.GetRawOrderBookSocket()

// 	id := utils.GetOrderBookChannelID(bt, qt)
// 	socket.Unsubscribe(id, conn)
// }
// func getOrderBookPayload(res *types.EngineResponse) interface{} {
// 	orderSide := make(map[string]string)
// 	matchSide := make([]map[string]string, 0)
// 	matchSideMap := make(map[string]*big.Int)

// 	if math.Sub(res.Order.Amount, res.Order.FilledAmount).Cmp(big.NewInt(0)) != 0 {
// 		orderSide["price"] = res.Order.PricePoint.String()
// 		orderSide["amount"] = math.Sub(res.Order.Amount, res.Order.FilledAmount).String()
// 	}

// 	if len(res.Matches) > 0 {
// 		for _, mo := range res.Matches {
// 			pp := mo.Order.PricePoint.String()
// 			if matchSideMap[pp] == nil {
// 				matchSideMap[pp] = big.NewInt(0)
// 			}

// 			matchSideMap[pp] = math.Add(matchSideMap[pp], mo.Trade.Amount)
// 		}
// 	}

// 	for price, amount := range matchSideMap {
// 		temp := map[string]string{
// 			"price":  price,
// 			"amount": math.Neg(amount).String(),
// 		}
// 		matchSide = append(matchSide, temp)
// 	}

// 	var response map[string]interface{}
// 	if res.Order.Side == "SELL" {
// 		response = map[string]interface{}{
// 			"asks": []map[string]string{orderSide},
// 			"bids": matchSide,
// 		}
// 	} else {
// 		response = map[string]interface{}{
// 			"asks": matchSide,
// 			"bids": []map[string]string{orderSide},
// 		}
// 	}

// 	return response
// }
