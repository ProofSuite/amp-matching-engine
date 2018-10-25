package engine

// func (e *Engine) GetRawOrderBook(pair *types.Pair) ([][]types.Order, error) {
// 	code := pair.Code()

// 	ob := e.orderbooks[code]
// 	if ob == nil {
// 		return nil, errors.New("Orderbook error")
// 	}

// 	book := ob.GetRawOrderBook(pair)
// 	// if err != nil {
// 	// 	logger.Error(err)
// 	// 	return nil, err
// 	// }

// 	return book, nil
// }

// func (e *Engine) GetOrderBook(pair *types.Pair) (asks, bids []*map[string]float64, err error) {
// 	code := pair.Code()

// 	ob := e.orderbooks[code]
// 	if ob == nil {
// 		return nil, nil, errors.New("Orderbook error")
// 	}

// 	asks, bids = ob.GetOrderBook(pair)
// 	return asks, bids, nil
// }

// GetRawOrderBook fetches the complete orderbook from redis for the required pair
// func (ob *OrderBook) GetRawOrderBook(p *types.Pair) (book [][]types.Order) {
// 	pattern := p.GetKVPrefix() + "::*::*::orders::*"

// 	length := 100
// 	book = make([][]types.Order, 0)
// 	keys, err := ob.redisConn.Keys(pattern)
// 	if err != nil {
// 		logger.Error(err)
// 		return
// 	}

// 	orders := make([]types.Order, 0)
// 	for start := 0; start < len(keys); start = start + length {
// 		end := start + length
// 		if len(keys) < end {
// 			end = len(keys)
// 		}

// 		res, err := ob.redisConn.MGet(keys[start:end]...)
// 		if err != nil {
// 			logger.Error(err)
// 			return
// 		}

// 		for _, r := range res {
// 			if r == "" {
// 				continue
// 			}

// 			var temp types.Order
// 			if err := json.Unmarshal([]byte(r), &temp); err != nil {
// 				continue
// 			}

// 			orders = append(orders, temp)
// 		}
// 	}

// 	for start := 0; start < len(orders); start = start + length {
// 		end := start + length
// 		if len(keys) < end {
// 			end = len(keys)
// 		}

// 		book = append(book, orders[start:end])
// 	}

// 	return
// }

// // GetOrderBook fetches the complete orderbook from redis for the required pair
// func (ob *OrderBook) GetOrderBook(pair *types.Pair) (asks, bids []*map[string]float64) {
// 	sKey, bKey := pair.GetOrderBookKeys()
// 	res, err := redigo.Int64s(ob.redisConn.Do("SORT", sKey, "GET", sKey+"::book::*", "GET", "#")) // Add price point to order book
// 	if err != nil {
// 		logger.Error(err)
// 	}

// 	for i := 0; i < len(res); i = i + 2 {
// 		temp := &map[string]float64{
// 			"amount": float64(res[i]),
// 			"price":  float64(res[i+1]),
// 		}
// 		asks = append(asks, temp)
// 	}

// 	res, err = redigo.Int64s(ob.redisConn.Do("SORT", bKey, "GET", bKey+"::book::*", "GET", "#", "DESC"))
// 	if err != nil {
// 		logger.Error(err)
// 	}

// 	for i := 0; i < len(res); i = i + 2 {
// 		temp := &map[string]float64{
// 			"amount": float64(res[i]),
// 			"price":  float64(res[i+1]),
// 		}
// 		bids = append(bids, temp)
// 	}

// 	return
// }

// The orderbook currently uses the four following data structures to store engine
// state in redis
// 1. Pricepoints set
// 2. Pricepoints volume set
// 3. Pricepoints hashes set
// 4. Orders map

// 1. The pricepoints set is an ordered set that store all pricepoints.
// Keys: ~ pair addresses + side (BUY or SELL)
// Values: pricepoints set (sorted set but all ranks are actually 0)

// 2. The pricepoints volume set is an order set that store the volume for a given pricepoint
// Keys: pair addresses + side + pricepoint
// Values: volume for corresponding (pair, pricepoint)

// 3. The pricepoints hashes set is an ordered set that stores a set of hashes ranked by creation time for a given pricepoint
// Keys: pair addresses + side + pricepoint
// Values: hashes of orders with corresponding pricepoint

// 4. The orders hashmap is a mapping that stores serialized orders
// Keys: hash
// Values: serialized order

// import (
// 	"encoding/json"
// 	"math/big"
// 	"sync"

// 	"github.com/Proofsuite/amp-matching-engine/rabbitmq"
// 	"github.com/Proofsuite/amp-matching-engine/redis"
// 	"github.com/Proofsuite/amp-matching-engine/types"
// 	"github.com/Proofsuite/amp-matching-engine/utils/math"
// 	"github.com/ethereum/go-ethereum/common"
// 	redigo "github.com/gomodule/redigo/redis"
// )

// type OrderBook struct {
// 	redisConn    *redis.RedisConnection
// 	rabbitMQConn *rabbitmq.Connection
// 	pair         *types.Pair
// 	mutex        *sync.Mutex
// }

// // newOrder calls buyOrder/sellOrder based on type of order recieved and
// // publishes the response back to rabbitmq
// func (ob *OrderBook) newOrder(order *types.Order) (err error) {
// 	// Attain lock on engineResource, so that recovery or cancel order function doesn't interfere
// 	ob.mutex.Lock()
// 	defer ob.mutex.Unlock()

// 	resp := &types.EngineResponse{}
// 	if order.Side == "SELL" {
// 		resp, err = ob.sellOrder(order)
// 		if err != nil {
// 			logger.Error(err)
// 			return err
// 		}

// 	} else if order.Side == "BUY" {
// 		resp, err = ob.buyOrder(order)
// 		if err != nil {
// 			logger.Error(err)
// 			return err
// 		}
// 	}

// 	// Note: Plug the option for orders like FOC, Limit here (if needed)
// 	err = ob.rabbitMQConn.PublishEngineResponse(resp)
// 	if err != nil {
// 		logger.Error(err)
// 		return err
// 	}

// 	return nil
// }

// // buyOrder is triggered when a buy order comes in, it fetches the ask list
// // from orderbook. First it checks ths price point list to check whether the order can be matched
// // or not, if there are pricepoints that can satisfy the order then corresponding list of orders
// // are fetched and trade is executed
// func (ob *OrderBook) buyOrder(order *types.Order) (*types.EngineResponse, error) {
// 	res := &types.EngineResponse{
// 		Order:  order,
// 		Status: "NOMATCH",
// 	}

// 	remainingOrder := *order
// 	res.RemainingOrder = &remainingOrder
// 	oskv := order.GetOBMatchKey()

// 	// GET Range of sellOrder between minimum Sell order and order.Price
// 	pps, err := ob.GetMatchingBuyPricePoints(oskv, order.PricePoint.Int64())
// 	if err != nil {
// 		logger.Error(err)
// 		return nil, err
// 	}

// 	if len(pps) == 0 {
// 		order.Status = "OPEN"
// 		res.Status = "NOMATCH"
// 		res.RemainingOrder = nil
// 		ob.addOrder(order)
// 		return res, nil
// 	}

// 	for _, pp := range pps {
// 		entries, err := ob.GetMatchingOrders(oskv, pp)
// 		if err != nil {
// 			logger.Error(err)
// 			return nil, err
// 		}

// 		for _, bookEntry := range entries {
// 			entry := &types.Order{}
// 			err = json.Unmarshal(bookEntry, &entry)
// 			if err != nil {
// 				return nil, err
// 			}

// 			trade, err := ob.execute(order, entry)
// 			if err != nil {
// 				logger.Error(err)
// 				return nil, err
// 			}

// 			match := &types.OrderTradePair{entry, trade}
// 			res.Matches = append(res.Matches, match)
// 			res.RemainingOrder.Amount = math.Sub(res.RemainingOrder.Amount, trade.Amount)

// 			if math.IsZero(res.RemainingOrder.Amount) {
// 				res.Status = "FULL"
// 				res.Order.Status = "FILLED"
// 				res.RemainingOrder = nil
// 				return res, nil
// 			}
// 		}
// 	}

// 	//TODO refactor this in a different function (make above function more clear in general)
// 	res.Order.Status = "PARTIAL_FILLED"
// 	res.Status = "PARTIAL"
// 	res.RemainingOrder.Signature = nil
// 	res.RemainingOrder.Nonce = nil
// 	res.RemainingOrder.Hash = common.HexToHash("")
// 	res.RemainingOrder.BuyAmount = res.RemainingOrder.Amount
// 	res.RemainingOrder.SellAmount = math.Div(
// 		math.Mul(res.RemainingOrder.Amount, res.Order.SellAmount),
// 		res.Order.BuyAmount,
// 	)

// 	return res, nil
// }

// // sellOrder is triggered when a sell order comes in, it fetches the bid list
// // from orderbook. First it checks ths price point list to check whether the order can be matched
// // or not, if there are pricepoints that can satisfy the order then corresponding list of orders
// // are fetched and trade is executed
// func (ob *OrderBook) sellOrder(order *types.Order) (*types.EngineResponse, error) {
// 	res := &types.EngineResponse{
// 		Status: "NOMATCH",
// 		Order:  order,
// 	}

// 	remOrder := *order
// 	res.RemainingOrder = &remOrder
// 	obkv := order.GetOBMatchKey()

// 	pps, err := ob.GetMatchingSellPricePoints(obkv, order.PricePoint.Int64())
// 	if err != nil {
// 		logger.Error(err)
// 		return nil, err
// 	}

// 	if len(pps) == 0 {
// 		res.Status = "NOMATCH"
// 		res.RemainingOrder = nil
// 		ob.addOrder(order)
// 		order.Status = "OPEN"
// 		return res, nil
// 	}

// 	for _, pp := range pps {
// 		entries, err := ob.GetMatchingOrders(obkv, pp)
// 		if err != nil {
// 			logger.Error(err)
// 			return nil, err
// 		}

// 		for _, o := range entries {
// 			entry := &types.Order{}
// 			err = json.Unmarshal(o, &entry)

// 			if err != nil {
// 				logger.Error(err)
// 				return nil, err
// 			}

// 			trade, err := ob.execute(order, entry)
// 			if err != nil {
// 				logger.Error(err)
// 				return nil, err
// 			}

// 			order.Status = "PARTIAL_FILLED"
// 			res.Status = "PARTIAL"
// 			match := &types.OrderTradePair{entry, trade}

// 			res.Matches = append(res.Matches, match)
// 			res.RemainingOrder.Amount = math.Sub(res.RemainingOrder.Amount, trade.Amount)

// 			if math.IsZero(res.RemainingOrder.Amount) {
// 				res.Status = "FULL"
// 				res.Order.Status = "FILLED"
// 				res.RemainingOrder = nil
// 				return res, nil
// 			}
// 		}
// 	}

// 	//TODO refactor this in a different function (make above function more clear in general)
// 	res.Order.Status = "PARTIAL_FILLED"
// 	res.Status = "PARTIAL"
// 	res.RemainingOrder.Signature = nil
// 	res.RemainingOrder.Nonce = nil
// 	res.RemainingOrder.Hash = common.HexToHash("")
// 	res.RemainingOrder.BuyAmount = res.RemainingOrder.Amount
// 	res.RemainingOrder.SellAmount = math.Div(
// 		math.Mul(res.RemainingOrder.Amount, res.Order.SellAmount),
// 		res.Order.BuyAmount,
// 	)

// 	return res, nil
// }

// // addOrder adds an order to redis
// func (ob *OrderBook) addOrder(order *types.Order) error {
// 	pricePointSetKey, orderHashListKey := order.GetOBKeys()
// 	err := ob.AddToPricePointSet(pricePointSetKey, order.PricePoint.Int64())
// 	if err != nil {
// 		logger.Error(err)
// 		return err
// 	}

// 	// Currently converting amount to int64. In the future, we need to use strings instead of int64
// 	volume := math.Div(math.Sub(order.Amount, order.FilledAmount), big.NewInt(1e18)).Int64()
// 	err = ob.IncrementPricePointVolume(pricePointSetKey, order.PricePoint.Int64(), volume)
// 	if err != nil {
// 		logger.Error(err)
// 		return err
// 	}

// 	err = ob.AddToOrderMap(order)
// 	if err != nil {
// 		logger.Error(err)
// 		return err
// 	}

// 	err = ob.AddToPricePointHashesSet(orderHashListKey, order.CreatedAt, order.Hash)
// 	if err != nil {
// 		logger.Error(err)
// 		return err
// 	}

// 	return nil
// }

// // updateOrder updates the order in redis
// func (ob *OrderBook) updateOrder(order *types.Order, tradeAmount *big.Int) error {
// 	stored := &types.Order{}

// 	pricePointSetKey, _ := order.GetOBKeys()
// 	stored, err := ob.GetFromOrderMap(order.Hash)
// 	if err != nil {
// 		logger.Error(err)
// 		return err
// 	}

// 	filledAmount := math.Add(stored.FilledAmount, tradeAmount)
// 	if math.IsEqualOrSmallerThan(filledAmount, big.NewInt(0)) {
// 		stored.Status = "OPEN"
// 		stored.FilledAmount = big.NewInt(0)
// 	} else if math.IsEqualOrGreaterThan(filledAmount, stored.Amount) {
// 		stored.Status = "FILLED"
// 		stored.FilledAmount = stored.Amount
// 	} else {
// 		stored.Status = "PARTIAL_FILLED"
// 		stored.FilledAmount = filledAmount
// 	}

// 	err = ob.AddToOrderMap(stored)
// 	if err != nil {
// 		logger.Error(err)
// 		return err
// 	}

// 	volume := math.Div(tradeAmount, big.NewInt(1e18)).Int64()
// 	err = ob.IncrementPricePointVolume(pricePointSetKey, order.PricePoint.Int64(), -volume)
// 	if err != nil {
// 		logger.Error(err)
// 		return err
// 	}

// 	return nil
// }

// // deleteOrder deletes the order in redis
// func (ob *OrderBook) deleteOrder(o *types.Order) (err error) {
// 	//TODO decide to put the mutex on deleteOrder or on cancelOrder
// 	// ob.mutex.Lock()
// 	// defer ob.mutex.Unlock()

// 	pricePointSetKey, orderHashListKey := o.GetOBKeys()
// 	pp := o.PricePoint.Int64()

// 	err = ob.RemoveFromPricePointHashesSet(orderHashListKey, o.Hash)
// 	if err != nil {
// 		logger.Error(err)
// 	}

// 	pplen, err := ob.GetPricePointHashesSetLength(orderHashListKey)
// 	if err != nil {
// 		logger.Error(err)
// 	}

// 	if pplen == 0 {
// 		err = ob.RemoveFromPricePointSet(pricePointSetKey, pp)
// 		if err != nil {
// 			logger.Error(err)
// 		}

// 		err = ob.DeletePricePointVolume(pricePointSetKey, pp)
// 		if err != nil {
// 			logger.Error(err)
// 		}
// 	} else {
// 		amt := math.Div(math.Sub(o.Amount, o.FilledAmount), big.NewInt(1e18)).Int64()
// 		err = ob.DecrementPricePointVolume(pricePointSetKey, pp, amt)
// 		if err != nil {
// 			logger.Error(err)
// 		}
// 	}

// 	err = ob.RemoveFromOrderMap(o.Hash)
// 	if err != nil {
// 		logger.Error(err)
// 	}

// 	return err
// }

// func (ob *OrderBook) deleteOrders(orders ...types.Order) error {
// 	for _, o := range orders {
// 		err := ob.deleteOrder(&o)
// 		if err != nil {
// 			logger.Error(err)
// 			return err
// 		}
// 	}

// 	return nil
// }

// // RecoverOrders is responsible for recovering the orders that failed to execute after matching
// // Orders are updated or added to orderbook based on whether that order exists in orderbook or not.
// func (ob *OrderBook) RecoverOrders(matches []*types.OrderTradePair) error {
// 	ob.mutex.Lock()
// 	defer ob.mutex.Unlock()

// 	for _, m := range matches {
// 		t := m.Trade
// 		o := m.Order

// 		o.Status = "PARTIAL_FILLED"
// 		filledAmount := math.Sub(o.FilledAmount, t.Amount)
// 		if math.IsEqualOrSmallerThan(filledAmount, big.NewInt(0)) {
// 			o.Status = "OPEN"
// 			o.FilledAmount = big.NewInt(0)
// 		}

// 		_, obListKey := o.GetOBKeys()
// 		if !ob.redisConn.Exists(obListKey + "::orders::" + o.Hash.Hex()) {
// 			if err := ob.addOrder(o); err != nil {
// 				logger.Error(err)
// 				return err
// 			}
// 		} else {
// 			if err := ob.updateOrder(o, math.Neg(t.Amount)); err != nil {
// 				logger.Error(err)
// 				return err
// 			}
// 		}
// 	}

// 	return nil
// }

// func (ob *OrderBook) CancelTrades(orders []*types.Order, amounts []*big.Int) error {
// 	ob.mutex.Lock()
// 	defer ob.mutex.Unlock()

// 	for i, o := range orders {
// 		o.Status = "PARTIAL_FILLED"
// 		o.FilledAmount = math.Sub(o.FilledAmount, amounts[i])
// 		if math.IsZero(o.FilledAmount) {
// 			o.Status = "OPEN"
// 		}
// 		_, obListKey := o.GetOBKeys()
// 		if !ob.redisConn.Exists(obListKey + "::orders::" + o.Hash.Hex()) {
// 			if err := ob.addOrder(o); err != nil {
// 				logger.Error(err)
// 				return err
// 			}
// 		} else {
// 			err := ob.updateOrder(o, math.Neg(o.Amount))
// 			if err != nil {
// 				logger.Error(err)
// 				return err
// 			}
// 		}
// 	}

// 	return nil
// }

// // CancelOrder is used to cancel the order from orderbook
// func (ob *OrderBook) CancelOrder(o *types.Order) (*types.EngineResponse, error) {
// 	ob.mutex.Lock()
// 	defer ob.mutex.Unlock()

// 	stored, err := ob.GetFromOrderMap(o.Hash)
// 	if err != nil {
// 		logger.Error(err)
// 		return nil, err
// 	}

// 	if err := ob.deleteOrder(o); err != nil {
// 		logger.Error(err)
// 		return nil, err
// 	}

// 	stored.Status = "CANCELLED"
// 	res := &types.EngineResponse{
// 		Status:         "CANCELLED",
// 		Order:          stored,
// 		RemainingOrder: nil,
// 		Matches:        nil,
// 	}

// 	return res, nil
// }

// // execute function is responsible for executing of matched orders
// // i.e it deletes/updates orders in case of order matching and responds
// // with trade instance and fillOrder
// func (ob *OrderBook) execute(order *types.Order, bookEntry *types.Order) (*types.Trade, error) {
// 	trade := &types.Trade{}
// 	tradeAmount := big.NewInt(0)
// 	bookEntryAvailableAmount := math.Sub(bookEntry.Amount, bookEntry.FilledAmount)
// 	orderAvailableAmount := math.Sub(order.Amount, order.FilledAmount)

// 	if math.IsGreaterThan(bookEntryAvailableAmount, orderAvailableAmount) {
// 		tradeAmount = orderAvailableAmount
// 		bookEntry.FilledAmount = math.Add(bookEntry.FilledAmount, orderAvailableAmount)
// 		bookEntry.Status = "PARTIAL_FILLED"

// 		err := ob.updateOrder(bookEntry, tradeAmount)
// 		if err != nil {
// 			logger.Error(err)
// 			return nil, err
// 		}

// 	} else {
// 		err := ob.deleteOrder(bookEntry)
// 		if err != nil {
// 			logger.Error(err)
// 			return nil, err
// 		}

// 		tradeAmount = bookEntryAvailableAmount
// 		bookEntry.FilledAmount = bookEntry.Amount
// 		bookEntry.Status = "FILLED"
// 	}

// 	order.FilledAmount = math.Add(order.FilledAmount, tradeAmount)
// 	trade = &types.Trade{
// 		Amount:         tradeAmount,
// 		PricePoint:     order.PricePoint,
// 		BaseToken:      order.BaseToken,
// 		QuoteToken:     order.QuoteToken,
// 		OrderHash:      bookEntry.Hash,
// 		TakerOrderHash: order.Hash,
// 		Side:           order.Side,
// 		Taker:          order.UserAddress,
// 		PairName:       order.PairName,
// 		Maker:          bookEntry.UserAddress,
// 	}

// 	return trade, nil
// }

// func (ob *OrderBook) GetPricePointVolume(pricePointSetKey string, pricePoint int64) (int64, error) {
// 	val, err := ob.redisConn.GetValue(pricePointSetKey + "::book::" + utils.UintToPaddedString(pricePoint))
// 	if err != nil {
// 		logger.Error(err)
// 		return 0, err
// 	}

// 	vol, err := strconv.ParseInt(val, 10, 64)
// 	if err != nil {
// 		logger.Error(err)
// 		return 0, err
// 	}

// 	return vol, nil
// }

// // IncrementPricePointVolume increases the value of a certain pricepoint at a certain volume
// func (ob *OrderBook) IncrementPricePointVolume(pricePointSetKey string, pricePoint int64, amount int64) error {
// 	_, err := ob.redisConn.IncrBy(pricePointSetKey+"::book::"+utils.UintToPaddedString(pricePoint), amount)
// 	if err != nil {
// 		logger.Error(err)
// 		return err
// 	}

// 	return nil
// }

// // DecrementPricePoint
// func (ob *OrderBook) DecrementPricePointVolume(pricePointSetKey string, pricePoint int64, amount int64) error {
// 	_, err := ob.redisConn.IncrBy(pricePointSetKey+"::book::"+utils.UintToPaddedString(pricePoint), -amount)
// 	if err != nil {
// 		logger.Error(err)
// 		return err
// 	}

// 	return nil
// }

// // DeletePricePoint
// func (ob *OrderBook) DeletePricePointVolume(pricePointSetKey string, pricePoint int64) error {
// 	err := ob.redisConn.Del(pricePointSetKey + "::book::" + utils.UintToPaddedString(pricePoint))
// 	if err != nil {
// 		logger.Error(err)
// 		return err
// 	}

// 	return nil
// }

// // GetPricePoints returns the pricepoints matching a certain (pair, pricepoint)
// func (ob *OrderBook) GetMatchingBuyPricePoints(obKey string, pricePoint int64) ([]int64, error) {
// 	pps, err := ob.redisConn.ZRangeByLexInt(obKey, "-", "["+utils.UintToPaddedString(pricePoint))
// 	if err != nil {
// 		logger.Error(err)
// 		return nil, err
// 	}

// 	return pps, nil
// }

// func (ob *OrderBook) GetMatchingSellPricePoints(obkv string, pricePoint int64) ([]int64, error) {
// 	pps, err := ob.redisConn.ZRevRangeByLexInt(obkv, "+", "["+utils.UintToPaddedString(pricePoint))
// 	if err != nil {
// 		logger.Error(err)
// 		return nil, err
// 	}

// 	return pps, nil
// }

// func (ob *OrderBook) GetFromOrderMap(hash common.Hash) (*types.Order, error) {
// 	o := &types.Order{}
// 	keys, err := ob.redisConn.Keys("*::" + hash.Hex())
// 	if err != nil {
// 		return nil, err
// 	}

// 	if len(keys) == 0 {
// 		return nil, fmt.Errorf("Key doesn't exist")
// 	}

// 	serialized, err := ob.redisConn.GetValue(keys[0])
// 	if err != nil {
// 		logger.Error(err)
// 		return nil, err
// 	}

// 	err = json.Unmarshal([]byte(serialized), &o)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return o, nil
// }

// // GetPricePointOrders returns the orders hashes for a (pair, pricepoint)
// func (ob *OrderBook) GetMatchingOrders(obKey string, pricePoint int64) ([][]byte, error) {
// 	k := obKey + "::" + utils.UintToPaddedString(pricePoint)
// 	orders, err := ob.redisConn.Sort(k, "", true, false, k+"::orders::*")
// 	if err != nil {
// 		return nil, err
// 	}

// 	return orders, nil
// }

// AddPricePointToSet
// func (ob *OrderBook) AddToPricePointSet(pricePointSetKey string, pricePoint int64) error {
// 	err := ob.redisConn.ZAdd(pricePointSetKey, 0, utils.UintToPaddedString(pricePoint))
// 	if err != nil {
// 		logger.Error(err)
// 		return err
// 	}

// 	return nil
// }

// func (ob *OrderBook) AddToPricePointSet(pricePointSetKey string, pricepoint *big.Int, amount *big.Int) error {
// 	volume := math.Div(amount, big.NewInt(1e18)).Int64()

// 	_, err := ob.redisConn.ZIncrBy(pricePointSetKey, pricepoint.Int64(), utils.UintToPaddedString())
// 	_, err := ob.redisConn.ZIncrBy(pricePointSetKey, volume, utils.UintToPaddedString(pricepoint.Int64()))
// 	if err != nil {
// 		logger.Error(err)
// 		return err
// 	}

// 	return nil
// }

// func (ob *OrderBook) RemoveFromPricePointSet(pricePointSetKey string, pricepoint *big.Int) error {
// 	err := ob.redisConn.ZRem(pricePointSetKey, utils.UintToPaddedString(pricepoint.Int64()))
// 	if err != nil {
// 		logger.Error(err)
// 		return err
// 	}

// 	return nil
// }

// func (ob *OrderBook) UpdatePricePointSet(pricePointSetKey string, pricepoint *big.Int, amount *big.Int) error {
// 	volume := math.Div(amount, big.NewInt(1e18))
// 	err := ob.redisConn.ZAdd(pricePointSetKey, )
// }

// // RemoveFromPricePointSet
// func (ob *OrderBook) RemoveFromPricePointSet(pricePointSetKey string, pricePoint int64) error {
// 	err := ob.redisConn.ZRem(pricePointSetKey, utils.UintToPaddedString(pricePoint))
// 	if err != nil {
// 		logger.Error(err)
// 		return err
// 	}

// 	return nil
// }

// func (ob *OrderBook) GetPricePointSetLength(pricePointSetKey string) (int64, error) {
// 	count, err := ob.redisConn.ZCount(pricePointSetKey)
// 	if err != nil {
// 		logger.Error(err)
// 		return 0, err
// 	}

// 	return count, nil
// }

// // AddPricePointHashesSet
// func (ob *OrderBook) AddToPricePointHashesSet(orderHashListKey string, createdAt time.Time, hash common.Hash) error {
// 	err := ob.redisConn.ZAdd(orderHashListKey, createdAt.Unix(), hash.Hex())
// 	if err != nil {
// 		logger.Error(err)
// 		return err
// 	}

// 	return nil
// }

// // RemoveFromPricePointHashesSet
// func (ob *OrderBook) RemoveFromPricePointHashesSet(orderHashListKey string, hash common.Hash) error {
// 	err := ob.redisConn.ZRem(orderHashListKey, hash.Hex())
// 	if err != nil {
// 		logger.Error(err)
// 		return err
// 	}

// 	return nil
// }

// func (ob *OrderBook) GetPricePointHashesSetLength(orderHashListKey string) (int64, error) {
// 	count, err := ob.redisConn.ZCount(orderHashListKey)
// 	if err != nil {
// 		logger.Error(err)
// 		return 0, err
// 	}

// 	return count, nil
// }

// // AddToOrderMap
// func (ob *OrderBook) AddToOrderMap(o *types.Order) error {
// 	bytes, err := json.Marshal(o)
// 	if err != nil {
// 		logger.Error(err)
// 		return err
// 	}

// 	decoded := &types.Order{}
// 	json.Unmarshal(bytes, decoded)
// 	_, orderHashListKey := o.GetOBKeys()
// 	err = ob.redisConn.Set(orderHashListKey+"::orders::"+o.Hash.Hex(), string(bytes))
// 	if err != nil {
// 		logger.Error(err)
// 		return err
// 	}

// 	return nil
// }

// // RemoveFromOrderMap
// func (ob *OrderBook) RemoveFromOrderMap(hash common.Hash) error {
// 	keys, _ := ob.redisConn.Keys("*::" + hash.Hex())
// 	err := ob.redisConn.Del(keys[0])
// 	if err != nil {
// 		logger.Error(err)
// 		return err
// 	}

// 	return nil
// }

// func TestAddOrder(t *testing.T) {
// 	e, ob, _, _, _, _, _, _, factory1, factory2 := setupTest()
// 	defer e.redisConn.FlushAll()

// 	o1, _ := factory1.NewSellOrder(1e3, 1e8)
// 	o2, _ := factory2.NewSellOrder(1e3, 1e8)

// 	e.addOrder(&o1)

// 	pricePointSetKey, orderHashListKey := o1.GetOBKeys()

// 	pricepoints, err := e.redisConn.GetSortedSet(pricePointSetKey)
// 	if err != nil {
// 		t.Error("Error getting sorted set")
// 	}

// 	pricePointHashes, err := e.redisConn.GetSortedSet(orderHashListKey)
// 	if err != nil {
// 		t.Error("Error getting pricepoints order hashes set")
// 	}

// 	stored, err := ob.GetFromOrderMap(o1.Hash)
// 	if err != nil {
// 		t.Error("Error getting sorted set", err)
// 	}

// 	// volume, err := ob.GetPricePointVolume(pricePointSetKey, o1.PricePoint.Int64())
// 	// if err != nil {
// 	// 	t.Error("Error getting volume set", err)
// 	// }

// 	assert.Equal(t, 1, len(pricepoints))
// 	assert.Contains(t, pricepoints, utils.UintToPaddedString(o1.PricePoint.Int64()))
// 	assert.Equal(t, 1, len(pricePointHashes))
// 	assert.Contains(t, pricePointHashes, o1.Hash.Hex())
// 	assert.Equal(t, int64(pricePointHashes[o1.Hash.Hex()]), o1.CreatedAt.Unix())
// 	// assert.Equal(t, int64(1e8), volume)
// 	testutils.CompareOrder(t, &o1, stored)

// 	e.addOrder(&o2)

// 	pricePointSetKey, orderHashListKey = o2.GetOBKeys()
// 	pricepoints, err = e.redisConn.GetSortedSet(pricePointSetKey)
// 	if err != nil {
// 		t.Error(err)
// 	}

// 	pricePointHashes, err = e.redisConn.GetSortedSet(orderHashListKey)
// 	if err != nil {
// 		t.Error("Error getting pricepoints order hashes set")
// 	}

// 	stored, err = ob.GetFromOrderMap(o2.Hash)
// 	if err != nil {
// 		t.Error("Error getting order from map", err)
// 	}

// 	// volume, err = ob.GetPricePointVolume(pricePointSetKey, o2.PricePoint.Int64())
// 	// if err != nil {
// 	// 	t.Error("Error getting volume set", err)
// 	// }

// 	assert.Equal(t, 1, len(pricepoints))
// 	assert.Contains(t, pricepoints, utils.UintToPaddedString(o2.PricePoint.Int64()))
// 	assert.Equal(t, 2, len(pricePointHashes))
// 	assert.Contains(t, pricePointHashes, o2.Hash.Hex())
// 	assert.Equal(t, int64(pricePointHashes[o2.Hash.Hex()]), o2.CreatedAt.Unix())
// 	// assert.Equal(t, int64(2e8), volume)
// 	testutils.CompareOrder(t, &o2, stored)
// }

// func TestUpdateOrder(t *testing.T) {
// 	e, ob, _, _, _, _, _, _, factory1, _ := setupTest()
// 	defer e.redisConn.FlushAll()

// 	o1, _ := factory1.NewSellOrder(1e3, 1e8)

// 	exp1 := o1
// 	exp1.Status = "PARTIAL_FILLED"
// 	exp1.FilledAmount = units.Ethers(1e3)

// 	err := ob.addOrder(&o1)
// 	if err != nil {
// 		t.Error("Could not add order")
// 	}

// 	err = ob.updateOrder(&o1, units.Ethers(1e3))
// 	if err != nil {
// 		t.Error("Could not update order")
// 	}

// 	pricePointSetKey, orderHashListKey := o1.GetOBKeys()

// 	pricepoints, err := ob.redisConn.GetSortedSet(pricePointSetKey)
// 	if err != nil {
// 		t.Error("Error getting pricepoint set", err)
// 	}

// 	pricePointHashes, err := ob.redisConn.GetSortedSet(orderHashListKey)
// 	if err != nil {
// 		t.Error("Error getting pricepoint hash set", err)
// 	}

// 	// volume, err := ob.GetPricePointVolume(pricePointSetKey, o1.PricePoint.Int64())
// 	// if err != nil {
// 	// 	t.Error("Error getting pricepoint volume", err)
// 	// }

// 	stored, err := ob.GetFromOrderMap(o1.Hash)
// 	if err != nil {
// 		t.Error(err)
// 	}

// 	assert.Equal(t, 1, len(pricepoints))
// 	assert.Contains(t, pricepoints, utils.UintToPaddedString(o1.PricePoint.Int64()))
// 	testutils.Compare(t, &exp1, stored)

// 	assert.Equal(t, 1, len(pricePointHashes))
// 	assert.Contains(t, pricePointHashes, o1.Hash.Hex())
// 	assert.Equal(t, int64(pricePointHashes[o1.Hash.Hex()]), o1.CreatedAt.Unix())
// 	// assert.Equal(t, int64(99999000), volume)
// }

// func TestDeleteOrder(t *testing.T) {
// 	e, ob, _, _, _, _, _, _, factory1, _ := setupTest()
// 	defer e.redisConn.FlushAll()

// 	o1, _ := factory1.NewSellOrder(1e3, 1e8)

// 	e.addOrder(&o1)

// 	pricePointSetKey, orderHashListKey := o1.GetOBKeys()

// 	pricepoints, err := e.redisConn.GetSortedSet(pricePointSetKey)
// 	if err != nil {
// 		t.Error(err)
// 	}

// 	pricePointHashes, err := e.redisConn.GetSortedSet(orderHashListKey)
// 	if err != nil {
// 		t.Error(err)
// 	}

// 	// volume, err := ob.GetPricePointVolume(pricePointSetKey, o1.PricePoint.Int64())
// 	// if err != nil {
// 	// 	t.Error(err)
// 	// }

// 	stored, err := ob.GetFromOrderMap(o1.Hash)
// 	if err != nil {
// 		t.Error(err)
// 	}

// 	assert.Equal(t, 1, len(pricepoints))
// 	assert.Contains(t, pricepoints, utils.UintToPaddedString(o1.PricePoint.Int64()))
// 	testutils.CompareOrder(t, &o1, stored)

// 	assert.Equal(t, 1, len(pricePointHashes))
// 	assert.Contains(t, pricePointHashes, o1.Hash.Hex())
// 	assert.Equal(t, int64(pricePointHashes[o1.Hash.Hex()]), o1.CreatedAt.Unix())
// 	// assert.Equal(t, int64(100000000), volume)

// 	err = ob.deleteOrder(&o1)
// 	if err != nil {
// 		t.Error(err)
// 	}

// 	pricePointSetKey, orderHashListKey = o1.GetOBKeys()

// 	if e.redisConn.Exists(pricePointSetKey) {
// 		t.Errorf("Key: %v expected to be deleted but exists", pricePointSetKey)
// 	}

// 	if e.redisConn.Exists(orderHashListKey) {
// 		t.Errorf("Key: %v expected to be deleted but key exists", pricePointSetKey)
// 	}

// 	if e.redisConn.Exists(orderHashListKey + "::" + o1.Hash.Hex()) {
// 		t.Errorf("Key: %v expected to be deleted but key exists", pricePointSetKey)
// 	}

// 	if e.redisConn.Exists(pricePointSetKey + "::book::" + utils.UintToPaddedString(o1.PricePoint.Int64())) {
// 		t.Errorf("Key: %v expected to be deleted but key exists", pricePointSetKey)
// 	}
// }

// func TestCancelOrder(t *testing.T) {
// 	e, ob, _, _, _, _, _, _, factory1, _ := setupTest()
// 	defer e.redisConn.FlushAll()

// 	o1, _ := factory1.NewSellOrder(1e3, 1e8)
// 	o2, _ := factory1.NewSellOrder(1e3, 1e8)

// 	e.addOrder(&o1)
// 	e.addOrder(&o2)

// 	expectedOrder := o2
// 	expectedOrder.Status = "CANCELLED"
// 	expected := &types.EngineResponse{
// 		Status:         "CANCELLED",
// 		HashID:         o2.Hash,
// 		Order:          &expectedOrder,
// 		RemainingOrder: nil,
// 		Matches:        nil,
// 	}

// 	pricePointSetKey, orderHashListKey := o1.GetOBKeys()

// 	pricepoints, _ := e.redisConn.GetSortedSet(pricePointSetKey)
// 	pricePointHashes, _ := e.redisConn.GetSortedSet(orderHashListKey)
// 	// volume, _ := ob.GetPricePointVolume(pricePointSetKey, o2.PricePoint.Int64())
// 	stored1, _ := ob.GetFromOrderMap(o1.Hash)
// 	stored2, _ := ob.GetFromOrderMap(o2.Hash)

// 	assert.Equal(t, 1, len(pricepoints))
// 	assert.Contains(t, pricepoints, utils.UintToPaddedString(o1.PricePoint.Int64()))
// 	assert.Equal(t, 2, len(pricePointHashes))
// 	assert.Contains(t, pricePointHashes, o1.Hash.Hex())
// 	assert.Contains(t, pricePointHashes, o2.Hash.Hex())
// 	assert.Equal(t, int64(pricePointHashes[o1.Hash.Hex()]), o1.CreatedAt.Unix())
// 	assert.Equal(t, int64(pricePointHashes[o2.Hash.Hex()]), o2.CreatedAt.Unix())
// 	// assert.Equal(t, int64(200000000), volume)
// 	testutils.Compare(t, &o1, stored1)
// 	testutils.Compare(t, &o2, stored2)

// 	res, err := ob.CancelOrder(&o2)
// 	if err != nil {
// 		t.Error("Error when cancelling order: ", err)
// 	}

// 	testutils.Compare(t, expected, res)

// 	pricePointSetKey, orderHashListKey = o1.GetOBKeys()

// 	pricepoints, _ = e.redisConn.GetSortedSet(pricePointSetKey)
// 	pricePointHashes, _ = e.redisConn.GetSortedSet(orderHashListKey)
// 	// volume, _ = ob.GetPricePointVolume(pricePointSetKey, o2.PricePoint.Int64())
// 	stored1, _ = ob.GetFromOrderMap(o1.Hash)
// 	stored2, _ = ob.GetFromOrderMap(o2.Hash)

// 	assert.Equal(t, 1, len(pricepoints))
// 	assert.Contains(t, pricepoints, utils.UintToPaddedString(o1.PricePoint.Int64()))
// 	assert.Equal(t, 1, len(pricePointHashes))
// 	assert.Contains(t, pricePointHashes, o1.Hash.Hex())
// 	assert.NotContains(t, pricePointHashes, o2.Hash.Hex())
// 	assert.Equal(t, int64(pricePointHashes[o1.Hash.Hex()]), o1.CreatedAt.Unix())
// 	// assert.Equal(t, int64(100000000), volume)
// 	testutils.Compare(t, &o1, stored1)
// 	testutils.Compare(t, nil, stored2)
// }
