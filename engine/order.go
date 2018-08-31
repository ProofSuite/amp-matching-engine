package engine

import (
	"encoding/json"
	"errors"
	"log"
	"math/big"

	"github.com/Proofsuite/amp-matching-engine/types"
	"github.com/Proofsuite/amp-matching-engine/utils"
	"github.com/Proofsuite/amp-matching-engine/utils/math"
	"github.com/ethereum/go-ethereum/common"
)

// newOrder calls buyOrder/sellOrder based on type of order recieved and
// publishes the response back to rabbitmq
func (e *Engine) newOrder(order *types.Order) (err error) {
	// Attain lock on engineResource, so that recovery or cancel order function doesn't interfere
	e.mutex.Lock()
	defer e.mutex.Unlock()

	resp := &types.EngineResponse{}
	if order.Side == "SELL" {
		resp, err = e.sellOrder(order)
		if err != nil {
			log.Print(err)
			return err
		}

	} else if order.Side == "BUY" {
		resp, err = e.buyOrder(order)
		if err != nil {
			log.Print(err)
			return err
		}
	}

	// Note: Plug the option for orders like FOC, Limit here (if needed)
	err = e.publishEngineResponse(resp)
	if err != nil {
		log.Print(err)
		return err
	}

	return nil
}

// buyOrder is triggered when a buy order comes in, it fetches the ask list
// from orderbook. First it checks ths price point list to check whether the order can be matched
// or not, if there are pricepoints that can satisfy the order then corresponding list of orders
// are fetched and trade is executed
func (e *Engine) buyOrder(order *types.Order) (*types.EngineResponse, error) {
	resp := &types.EngineResponse{
		Order:      order,
		FillStatus: "NOMATCH",
	}

	remainingOrder := *order
	resp.RemainingOrder = &remainingOrder
	resp.Trades = make([]*types.Trade, 0)
	resp.MatchingOrders = make([]*types.FillOrder, 0)
	oskv := order.GetOBMatchKey()

	// GET Range of sellOrder between minimum Sell order and order.Price
	priceRange, err := e.redisConn.ZRangeByLexInt(oskv, "-", "["+utils.UintToPaddedString(order.PricePoint.Int64()))
	if err != nil {
		log.Print(err)
		return nil, err
	}

	if len(priceRange) == 0 {
		resp.FillStatus = "NOMATCH"
		resp.RemainingOrder = nil
		order.Status = "OPEN"
		e.addOrder(order)
		return resp, nil
	}

	for _, pr := range priceRange {
		bookEntries, err := e.redisConn.Sort(oskv+"::"+utils.UintToPaddedString(pr), "", true, false, "*")
		if err != nil {
			log.Print(err)
			return nil, err
		}

		for _, o := range bookEntries {
			bookEntry := &types.Order{}
			err = json.Unmarshal(o, &bookEntry)
			if err != nil {
				log.Print(err)
				return nil, err
			}

			trade, fillOrder, err := e.execute(order, bookEntry)
			if err != nil {
				log.Print(err)
				return nil, err
			}

			resp.Trades = append(resp.Trades, trade)
			resp.MatchingOrders = append(resp.MatchingOrders, fillOrder)
			resp.RemainingOrder.Amount = math.Sub(resp.RemainingOrder.Amount, fillOrder.Amount)

			if math.IsZero(resp.RemainingOrder.Amount) {
				resp.FillStatus = "FULL"
				resp.Order.Status = "FILLED"
				resp.RemainingOrder = nil
				return resp, nil
			}
		}
	}

	resp.Order.Status = "PARTIAL_FILLED"
	resp.FillStatus = "PARTIAL"

	//TODO refactor this in a different function (make above function more clear in general)
	resp.RemainingOrder.Signature = nil
	resp.RemainingOrder.Nonce = nil
	resp.RemainingOrder.Hash = common.HexToHash("")
	resp.RemainingOrder.BuyAmount = resp.RemainingOrder.Amount
	resp.RemainingOrder.SellAmount = math.Div(
		math.Mul(resp.RemainingOrder.Amount, resp.Order.SellAmount),
		resp.Order.BuyAmount,
	)

	return resp, nil
}

// sellOrder is triggered when a sell order comes in, it fetches the bid list
// from orderbook. First it checks ths price point list to check whether the order can be matched
// or not, if there are pricepoints that can satisfy the order then corresponding list of orders
// are fetched and trade is executed
func (e *Engine) sellOrder(order *types.Order) (*types.EngineResponse, error) {
	resp := &types.EngineResponse{
		Order:      order,
		FillStatus: "NOMATCH",
	}

	remOrder := *order
	resp.Trades = make([]*types.Trade, 0)
	resp.RemainingOrder = &remOrder
	resp.MatchingOrders = make([]*types.FillOrder, 0)
	obkv := order.GetOBMatchKey()

	// GET Range of sellOrder between minimum Sell order and order.Price
	priceRange, err := e.redisConn.ZRevRangeByLexInt(obkv, "+", "["+utils.UintToPaddedString(order.PricePoint.Int64()))
	if err != nil {
		log.Print(err)
		return nil, err
	}

	if len(priceRange) == 0 {
		resp.FillStatus = "NOMATCH"
		resp.RemainingOrder = nil
		e.addOrder(order)
		order.Status = "OPEN"
		return resp, nil
	}

	for _, pr := range priceRange {
		bookEntries, err := e.redisConn.Sort(obkv+"::"+utils.UintToPaddedString(pr), "", true, false, "*")
		if err != nil {
			log.Print(err)
			return nil, err
		}

		for _, o := range bookEntries {
			bookEntry := &types.Order{}
			err = json.Unmarshal(o, &bookEntry)

			if err != nil {
				log.Print(err)
				return nil, err
			}

			trade, fillOrder, err := e.execute(order, bookEntry)
			if err != nil {
				log.Print(err)
				return nil, err
			}

			order.Status = "PARTIAL_FILLED"
			resp.FillStatus = "PARTIAL"
			resp.Trades = append(resp.Trades, trade)
			resp.MatchingOrders = append(resp.MatchingOrders, fillOrder)
			resp.RemainingOrder.Amount = math.Sub(resp.RemainingOrder.Amount, fillOrder.Amount)

			if math.IsZero(resp.RemainingOrder.Amount) {
				resp.FillStatus = "FULL"
				resp.Order.Status = "FILLED"
				resp.RemainingOrder = nil
				return resp, nil
			}
		}
	}

	resp.Order.Status = "PARTIAL_FILLED"
	resp.FillStatus = "PARTIAL"

	//TODO refactor this in a different function (make above function more clear in general)
	resp.RemainingOrder.Signature = nil
	resp.RemainingOrder.Nonce = nil
	resp.RemainingOrder.Hash = common.HexToHash("")
	resp.RemainingOrder.BuyAmount = resp.RemainingOrder.Amount
	resp.RemainingOrder.SellAmount = math.Div(
		math.Mul(resp.RemainingOrder.Amount, resp.Order.SellAmount),
		resp.Order.BuyAmount,
	)

	return resp, nil
}

// addOrder adds an order to redis
func (e *Engine) addOrder(order *types.Order) error {
	pricePointSetKey, orderHashListKey := order.GetOBKeys()
	if err := e.redisConn.ZAdd(pricePointSetKey, 0, utils.UintToPaddedString(order.PricePoint.Int64())); err != nil {
		log.Print(err)
		return err
	}

	// Currently converting amount to int64. In the future, we need to use strings instead of int64
	amt := math.Sub(order.Amount, order.FilledAmount)
	if _, err := e.redisConn.IncrBy(pricePointSetKey+"::book::"+utils.UintToPaddedString(order.PricePoint.Int64()), amt.Int64()); err != nil {
		log.Print(err)
		return err
	}

	// Add order to list
	orderAsBytes, err := json.Marshal(order)
	if err != nil {
		log.Print(err)
		return err
	}

	decoded := &types.Order{}
	json.Unmarshal(orderAsBytes, decoded)

	if err := e.redisConn.Set(order.Hash.Hex(), string(orderAsBytes)); err != nil {
		log.Print(err)
		return err
	}

	// Add order reference to price sorted set
	if err := e.redisConn.ZAdd(orderHashListKey, order.CreatedAt.Unix(), order.Hash.Hex()); err != nil {
		log.Print(err)
		return err
	}

	return nil
}

// updateOrder updates the order in redis
func (e *Engine) updateOrder(order *types.Order, tradeAmount *big.Int) error {
	stored := &types.Order{}

	pricePointSetKey, _ := order.GetOBKeys()
	orderString, err := e.redisConn.GetValue(order.Hash.Hex())
	if err != nil {

		log.Print(err)
		return err
	}

	json.Unmarshal([]byte(orderString), &stored)

	stored.FilledAmount = math.Add(stored.FilledAmount, tradeAmount)
	if math.IsZero(stored.FilledAmount) {
		stored.Status = "OPEN"
	} else if math.IsSmallerThan(stored.FilledAmount, stored.Amount) {
		stored.Status = "PARTIAL_FILLED"
	} else {
		stored.Status = "FILLED"
	}

	// Add order to list
	bytes, err := json.Marshal(stored)
	if err != nil {
		log.Print(err)
		return err
	}

	if err := e.redisConn.Set(order.Hash.Hex(), string(bytes)); err != nil {
		log.Print(err)
		return err
	}

	// Currently converting amount to int64. In the future, we need to use strings instead of int64
	if _, err := e.redisConn.IncrBy(pricePointSetKey+"::book::"+utils.UintToPaddedString(order.PricePoint.Int64()), math.Neg(tradeAmount).Int64()); err != nil {
		log.Print(err)
		return err
	}

	return nil
}

// updateOrderAmount is a less general version of update order that updates the order filled amount after
// a trade or recover orders
// EXPERIMENTAL
func (e *Engine) updateOrderAmount(hash common.Hash, amount *big.Int) error {
	stored := &types.Order{}
	orderString, err := e.redisConn.GetValue(hash.Hex())
	if err != nil {
		return err
	}

	err = json.Unmarshal([]byte(orderString), stored)
	if err != nil {
		log.Print(err)
		return err
	}

	stored.FilledAmount = math.Add(stored.FilledAmount, amount)

	if math.IsZero(stored.FilledAmount) {
		stored.Status = "OPEN"
	} else if math.IsSmallerThan(stored.FilledAmount, stored.Amount) {
		stored.Status = "PARTIAL_FILLED"
	} else {
		stored.Status = "FILLED"
	}

	bytes, err := json.Marshal(stored)
	if err != nil {
		log.Print(err)
		return err
	}

	if err := e.redisConn.Set(hash.Hex(), string(bytes)); err != nil {
		log.Print(err)
		return err
	}

	pricePointSetKey, _ := stored.GetOBKeys()

	// Currently converting amount to int64. In the future, we need to use strings instead of int64

	if _, err = e.redisConn.IncrBy(pricePointSetKey+"::book::"+utils.UintToPaddedString(stored.PricePoint.Int64()), math.Neg(amount).Int64()); err != nil {
		log.Print(err)
		return err
	}

	return nil
}

// deleteOrder deletes the order in redis
func (e *Engine) deleteOrder(order *types.Order, tradeAmount *big.Int) (err error) {
	pricePointSetKey, orderHashListKey := order.GetOBKeys()
	remVolume, err := e.redisConn.GetValue(pricePointSetKey + "::book::" + utils.UintToPaddedString(order.PricePoint.Int64()))
	if err != nil {
		log.Print(err)
		return
	}

	if math.IsEqual(math.ToBigInt(remVolume), tradeAmount) {
		if err := e.redisConn.ZRem(pricePointSetKey, utils.UintToPaddedString(order.PricePoint.Int64())); err != nil {
			log.Print(err)
			return err
		}
		if err := e.redisConn.Del(pricePointSetKey + "::book::" + utils.UintToPaddedString(order.PricePoint.Int64())); err != nil {
			log.Print(err)
			return err
		}
		if err := e.redisConn.Del(order.Hash.Hex()); err != nil {
			log.Print(err)
			return err
		}
		// Add order reference to price sorted set
		if err := e.redisConn.ZRem(orderHashListKey, order.Hash.Hex()); err != nil {
			log.Print(err)
			return err
		}

	} else {
		if err := e.redisConn.ZAdd(pricePointSetKey, 0, utils.UintToPaddedString(order.PricePoint.Int64())); err != nil {
			log.Print(err)
			return err
		}

		// Currently converting amount to int64. In the future, we need to use strings instead of int64
		if _, err := e.redisConn.IncrBy(pricePointSetKey+"::book::"+utils.UintToPaddedString(order.PricePoint.Int64()), math.Neg(tradeAmount).Int64()); err != nil {
			log.Print(err)
			return err
		}

		if err := e.redisConn.Del(order.Hash.Hex()); err != nil {
			log.Print(err)
			return err
		}
		// Add order reference to price sorted set
		if err := e.redisConn.ZRem(orderHashListKey, order.Hash.Hex()); err != nil {
			log.Print(err)
			return err
		}
	}
	return
}

// RecoverOrders is responsible for recovering the orders that failed to execute after matching
// Orders are updated or added to orderbook based on whether that order exists in orderbook or not.
func (e *Engine) RecoverOrders(orders []*types.FillOrder) error {
	e.mutex.Lock()
	defer e.mutex.Unlock()
	for _, o := range orders {

		// update order's filled amount and status before updating in redis
		o.Order.Status = "PARTIAL_FILLED"
		o.Order.FilledAmount = math.Sub(o.Order.FilledAmount, o.Amount)
		if math.IsZero(o.Order.FilledAmount) {
			o.Order.Status = "OPEN"
		}

		if !e.redisConn.Exists(o.Order.Hash.Hex()) {
			if err := e.addOrder(o.Order); err != nil {
				log.Print(err)
				return err
			}
		} else {
			if err := e.updateOrder(o.Order, math.Neg(o.Amount)); err != nil {
				log.Print(err)
				return err
			}
		}
	}
	return nil
}

// // RecoverOrders2 is an alternative suggestion for RecoverOrders2
// // It would requires an alternative key system for redis where we store only the hash with the listkey prefix
// func (e *Resource) RecoverOrders2(hashes []common.Hash, amounts []*big.Int) error {
// 	for i, _ := range hashes {
// 		err := e.updateOrderAmount(hashes[i], amounts[i])
// 		if err != nil {
// 			log.Print(err)
// 			return err
// 		}
// 	}

// 	return nil
// }

// CancelOrder is used to cancel the order from orderbook
func (e *Engine) CancelOrder(order *types.Order) (*types.EngineResponse, error) {
	e.mutex.Lock()
	defer e.mutex.Unlock()

	res, err := e.redisConn.GetValue(order.Hash.Hex())
	if err != nil {
		log.Print(err)
		return nil, err
	}

	if res == "" {
		return nil, errors.New("Order not found")
	}

	var stored *types.Order
	if err := json.Unmarshal([]byte(res), &stored); err != nil {
		log.Print(err)
		return nil, err
	}

	amt := math.Sub(stored.Amount, stored.FilledAmount)
	if err := e.deleteOrder(order, amt); err != nil {
		log.Print(err)
		return nil, err
	}

	stored.Status = "CANCELLED"

	engineResponse := &types.EngineResponse{
		Order:          stored,
		Trades:         make([]*types.Trade, 0),
		RemainingOrder: &types.Order{},
		FillStatus:     "CANCELLED",
		MatchingOrders: make([]*types.FillOrder, 0),
	}
	return engineResponse, nil
}
