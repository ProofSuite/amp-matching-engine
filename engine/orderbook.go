package engine

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

import (
	"encoding/json"
	"math/big"
	"sync"

	"github.com/Proofsuite/amp-matching-engine/rabbitmq"
	"github.com/Proofsuite/amp-matching-engine/redis"
	"github.com/Proofsuite/amp-matching-engine/types"
	"github.com/Proofsuite/amp-matching-engine/utils/math"
	"github.com/ethereum/go-ethereum/common"
)

type OrderBook struct {
	redisConn    *redis.RedisConnection
	rabbitMQConn *rabbitmq.Connection
	pair         *types.Pair
	mutex        *sync.Mutex
}

// newOrder calls buyOrder/sellOrder based on type of order recieved and
// publishes the response back to rabbitmq
func (ob *OrderBook) newOrder(order *types.Order, hashID common.Hash) (err error) {
	// Attain lock on engineResource, so that recovery or cancel order function doesn't interfere
	ob.mutex.Lock()
	defer ob.mutex.Unlock()

	resp := &types.EngineResponse{}
	if order.Side == "SELL" {
		resp, err = ob.sellOrder(order)
		if err != nil {
			logger.Error(err)
			return err
		}

	} else if order.Side == "BUY" {
		resp, err = ob.buyOrder(order)
		if err != nil {
			logger.Error(err)
			return err
		}
	}

	// Note: Plug the option for orders like FOC, Limit here (if needed)
	resp.HashID = hashID
	err = ob.rabbitMQConn.PublishEngineResponse(resp)
	if err != nil {
		logger.Error(err)
		return err
	}

	return nil
}

// buyOrder is triggered when a buy order comes in, it fetches the ask list
// from orderbook. First it checks ths price point list to check whether the order can be matched
// or not, if there are pricepoints that can satisfy the order then corresponding list of orders
// are fetched and trade is executed
func (ob *OrderBook) buyOrder(order *types.Order) (*types.EngineResponse, error) {
	res := &types.EngineResponse{
		Order:  order,
		Status: "NOMATCH",
	}

	remainingOrder := *order
	res.RemainingOrder = &remainingOrder
	oskv := order.GetOBMatchKey()

	// GET Range of sellOrder between minimum Sell order and order.Price
	pps, err := ob.GetMatchingBuyPricePoints(oskv, order.PricePoint.Int64())
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	if len(pps) == 0 {
		order.Status = "OPEN"
		res.Status = "NOMATCH"
		res.RemainingOrder = nil
		ob.addOrder(order)
		return res, nil
	}

	for _, pp := range pps {
		entries, err := ob.GetMatchingOrders(oskv, pp)
		if err != nil {
			logger.Error(err)
			return nil, err
		}

		for _, bookEntry := range entries {
			entry := &types.Order{}
			err = json.Unmarshal(bookEntry, &entry)
			if err != nil {
				return nil, err
			}

			trade, err := ob.execute(order, entry)
			if err != nil {
				logger.Error(err)
				return nil, err
			}

			match := &types.OrderTradePair{entry, trade}
			res.Matches = append(res.Matches, match)
			res.RemainingOrder.Amount = math.Sub(res.RemainingOrder.Amount, trade.Amount)

			if math.IsZero(res.RemainingOrder.Amount) {
				res.Status = "FULL"
				res.Order.Status = "FILLED"
				res.RemainingOrder = nil
				return res, nil
			}
		}
	}

	//TODO refactor this in a different function (make above function more clear in general)
	res.Order.Status = "REPLACED"
	res.Status = "PARTIAL"
	res.RemainingOrder.Signature = nil
	res.RemainingOrder.Nonce = nil
	res.RemainingOrder.Hash = common.HexToHash("")
	res.RemainingOrder.BuyAmount = res.RemainingOrder.Amount
	res.RemainingOrder.SellAmount = math.Div(
		math.Mul(res.RemainingOrder.Amount, res.Order.SellAmount),
		res.Order.BuyAmount,
	)

	return res, nil
}

// sellOrder is triggered when a sell order comes in, it fetches the bid list
// from orderbook. First it checks ths price point list to check whether the order can be matched
// or not, if there are pricepoints that can satisfy the order then corresponding list of orders
// are fetched and trade is executed
func (ob *OrderBook) sellOrder(order *types.Order) (*types.EngineResponse, error) {
	res := &types.EngineResponse{
		Status: "NOMATCH",
		Order:  order,
	}

	remOrder := *order
	res.RemainingOrder = &remOrder
	obkv := order.GetOBMatchKey()

	pps, err := ob.GetMatchingSellPricePoints(obkv, order.PricePoint.Int64())
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	if len(pps) == 0 {
		res.Status = "NOMATCH"
		res.RemainingOrder = nil
		ob.addOrder(order)
		order.Status = "OPEN"
		return res, nil
	}

	for _, pp := range pps {
		entries, err := ob.GetMatchingOrders(obkv, pp)
		if err != nil {
			logger.Error(err)
			return nil, err
		}

		for _, o := range entries {
			entry := &types.Order{}
			err = json.Unmarshal(o, &entry)

			if err != nil {
				logger.Error(err)
				return nil, err
			}

			trade, err := ob.execute(order, entry)
			if err != nil {
				logger.Error(err)
				return nil, err
			}

			match := &types.OrderTradePair{entry, trade}
			res.Matches = append(res.Matches, match)
			res.RemainingOrder.Amount = math.Sub(res.RemainingOrder.Amount, trade.Amount)

			if math.IsZero(res.RemainingOrder.Amount) {
				res.Status = "FULL"
				res.Order.Status = "FILLED"
				res.RemainingOrder = nil
				return res, nil
			}
		}
	}

	//TODO refactor this in a different function (make above function more clear in general)
	res.Order.Status = "REPLACED"
	res.Status = "PARTIAL"
	res.RemainingOrder.Signature = nil
	res.RemainingOrder.Nonce = nil
	res.RemainingOrder.Hash = common.HexToHash("")
	res.RemainingOrder.BuyAmount = res.RemainingOrder.Amount
	res.RemainingOrder.SellAmount = math.Div(
		math.Mul(res.RemainingOrder.Amount, res.Order.SellAmount),
		res.Order.BuyAmount,
	)

	return res, nil
}

// addOrder adds an order to redis
func (ob *OrderBook) addOrder(order *types.Order) error {
	order.Status = "OPEN"
	pricePointSetKey, orderHashListKey := order.GetOBKeys()
	err := ob.AddToPricePointSet(pricePointSetKey, order.PricePoint.Int64())
	if err != nil {
		logger.Error(err)
		return err
	}

	err = ob.AddToOrderMap(order)
	if err != nil {
		logger.Error(err)
		return err
	}

	err = ob.AddToPricePointHashesSet(orderHashListKey, order.CreatedAt, order.Hash)
	if err != nil {
		logger.Error(err)
		return err
	}

	return nil
}

// updateOrder updates the order in redis
func (ob *OrderBook) updateOrder(order *types.Order, tradeAmount *big.Int) error {
	stored := &types.Order{}

	stored, err := ob.GetFromOrderMap(order.Hash)
	if err != nil {
		logger.Error(err)
		return err
	}

	filledAmount := math.Add(stored.FilledAmount, tradeAmount)
	if math.IsEqualOrSmallerThan(filledAmount, big.NewInt(0)) {
		stored.Status = "OPEN"
		stored.FilledAmount = big.NewInt(0)
	} else if math.IsEqualOrGreaterThan(filledAmount, stored.Amount) {
		stored.Status = "FILLED"
		stored.FilledAmount = stored.Amount
	} else {
		stored.Status = "PARTIAL_FILLED"
		stored.FilledAmount = filledAmount
	}

	err = ob.AddToOrderMap(stored)
	if err != nil {
		logger.Error(err)
		return err
	}

	return nil
}

// deleteOrder deletes the order in redis
func (ob *OrderBook) deleteOrder(o *types.Order) (err error) {
	//TODO decide to put the mutex on deleteOrder or on cancelOrder
	// ob.mutex.Lock()
	// defer ob.mutex.Unlock()

	pricePointSetKey, orderHashListKey := o.GetOBKeys()
	pp := o.PricePoint.Int64()

	err = ob.RemoveFromPricePointHashesSet(orderHashListKey, o.Hash)
	if err != nil {
		logger.Error(err)
	}

	pplen, err := ob.GetPricePointHashesSetLength(orderHashListKey)
	if err != nil {
		logger.Error(err)
	}

	if pplen == 0 {
		err = ob.RemoveFromPricePointSet(pricePointSetKey, pp)
		if err != nil {
			logger.Error(err)
		}
	}

	err = ob.RemoveFromOrderMap(o.Hash)
	if err != nil {
		logger.Error(err)
	}

	return err
}

func (ob *OrderBook) deleteOrders(orders ...types.Order) error {
	for _, o := range orders {
		err := ob.deleteOrder(&o)
		if err != nil {
			logger.Error(err)
			return err
		}
	}

	return nil
}

// RecoverOrders is responsible for recovering the orders that failed to execute after matching
// Orders are updated or added to orderbook based on whether that order exists in orderbook or not.
func (ob *OrderBook) RecoverOrders(matches []*types.OrderTradePair) error {
	ob.mutex.Lock()
	defer ob.mutex.Unlock()

	for _, m := range matches {
		t := m.Trade
		o := m.Order

		o.Status = "PARTIAL_FILLED"
		filledAmount := math.Sub(o.FilledAmount, t.Amount)
		if math.IsEqualOrSmallerThan(filledAmount, big.NewInt(0)) {
			o.Status = "OPEN"
			o.FilledAmount = big.NewInt(0)
		}

		_, obListKey := o.GetOBKeys()
		if !ob.redisConn.Exists(obListKey + "::orders::" + o.Hash.Hex()) {
			if err := ob.addOrder(o); err != nil {
				logger.Error(err)
				return err
			}
		} else {
			if err := ob.updateOrder(o, math.Neg(t.Amount)); err != nil {
				logger.Error(err)
				return err
			}
		}
	}

	return nil
}

func (ob *OrderBook) CancelTrades(orders []*types.Order, amounts []*big.Int) error {
	ob.mutex.Lock()
	defer ob.mutex.Unlock()

	for i, o := range orders {
		o.Status = "PARTIAL_FILLED"
		o.FilledAmount = math.Sub(o.FilledAmount, amounts[i])
		if math.IsZero(o.FilledAmount) {
			o.Status = "OPEN"
		}
		_, obListKey := o.GetOBKeys()
		if !ob.redisConn.Exists(obListKey + "::orders::" + o.Hash.Hex()) {
			if err := ob.addOrder(o); err != nil {
				logger.Error(err)
				return err
			}
		} else {
			err := ob.updateOrder(o, math.Neg(o.Amount))
			if err != nil {
				logger.Error(err)
				return err
			}
		}
	}

	return nil
}

// CancelOrder is used to cancel the order from orderbook
func (ob *OrderBook) CancelOrder(o *types.Order) (*types.EngineResponse, error) {
	ob.mutex.Lock()
	defer ob.mutex.Unlock()

	stored, err := ob.GetFromOrderMap(o.Hash)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	if err := ob.deleteOrder(o); err != nil {
		logger.Error(err)
		return nil, err
	}

	stored.Status = "CANCELLED"
	res := &types.EngineResponse{
		Status:         "CANCELLED",
		Order:          stored,
		RemainingOrder: nil,
		Matches:        nil,
	}

	return res, nil
}

// execute function is responsible for executing of matched orders
// i.e it deletes/updates orders in case of order matching and responds
// with trade instance and fillOrder
func (ob *OrderBook) execute(order *types.Order, bookEntry *types.Order) (*types.Trade, error) {
	trade := &types.Trade{}
	tradeAmount := big.NewInt(0)
	bookEntryAvailableAmount := math.Sub(bookEntry.Amount, bookEntry.FilledAmount)
	orderAvailableAmount := math.Sub(order.Amount, order.FilledAmount)

	if math.IsGreaterThan(bookEntryAvailableAmount, orderAvailableAmount) {
		tradeAmount = orderAvailableAmount
		bookEntry.FilledAmount = math.Add(bookEntry.FilledAmount, orderAvailableAmount)
		bookEntry.Status = "PARTIAL_FILLED"

		err := ob.updateOrder(bookEntry, tradeAmount)
		if err != nil {
			logger.Error(err)
			return nil, err
		}

	} else {
		err := ob.deleteOrder(bookEntry)
		if err != nil {
			logger.Error(err)
			return nil, err
		}

		tradeAmount = bookEntryAvailableAmount
		bookEntry.FilledAmount = bookEntry.Amount
		bookEntry.Status = "FILLED"
	}

	order.FilledAmount = math.Add(order.FilledAmount, tradeAmount)
	trade = &types.Trade{
		Amount:         tradeAmount,
		PricePoint:     order.PricePoint,
		BaseToken:      order.BaseToken,
		QuoteToken:     order.QuoteToken,
		OrderHash:      bookEntry.Hash,
		TakerOrderHash: order.Hash,
		Side:           order.Side,
		Taker:          order.UserAddress,
		PairName:       order.PairName,
		Maker:          bookEntry.UserAddress,
	}

	return trade, nil
}
