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
	"log"
	"math/big"
	"time"

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
	res := &types.EngineResponse{
		Order:      order,
		FillStatus: "NOMATCH",
	}

	remainingOrder := *order
	res.RemainingOrder = &remainingOrder
	oskv := order.GetOBMatchKey()

	// GET Range of sellOrder between minimum Sell order and order.Price
	pps, err := e.GetMatchingBuyPricePoints(oskv, order.PricePoint.Int64())
	if err != nil {
		log.Print(err)
		return nil, err
	}

	if len(pps) == 0 {
		res.FillStatus = "NOMATCH"
		res.RemainingOrder = nil
		order.Status = "OPEN"
		e.addOrder(order)
		return res, nil
	}

	for _, pp := range pps {
		entries, err := e.GetMatchingOrders(oskv, pp)
		if err != nil {
			log.Print(err)
			return nil, err
		}

		for _, bookEntry := range entries {
			entry := &types.Order{}
			err = json.Unmarshal(bookEntry, &entry)
			if err != nil {
				log.Print(err)
				return nil, err
			}

			trade, fillOrder, err := e.execute(order, entry)
			if err != nil {
				log.Print(err)
				return nil, err
			}

			res.Trades = append(res.Trades, trade)
			res.MatchingOrders = append(res.MatchingOrders, fillOrder)
			res.RemainingOrder.Amount = math.Sub(res.RemainingOrder.Amount, fillOrder.Amount)

			if math.IsZero(res.RemainingOrder.Amount) {
				res.FillStatus = "FULL"
				res.Order.Status = "FILLED"
				res.RemainingOrder = nil
				return res, nil
			}
		}
	}

	//TODO refactor this in a different function (make above function more clear in general)
	res.Order.Status = "PARTIAL_FILLED"
	res.FillStatus = "PARTIAL"
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
func (e *Engine) sellOrder(order *types.Order) (*types.EngineResponse, error) {
	res := &types.EngineResponse{
		Order:      order,
		FillStatus: "NOMATCH",
	}

	remOrder := *order
	res.RemainingOrder = &remOrder
	obkv := order.GetOBMatchKey()

	// // GET Range of sellOrder between minimum Sell order and order.Price
	pps, err := e.GetMatchingSellPricePoints(obkv, order.PricePoint.Int64())
	if err != nil {
		log.Print(err)
		return nil, err
	}

	if len(pps) == 0 {
		res.FillStatus = "NOMATCH"
		res.RemainingOrder = nil
		e.addOrder(order)
		order.Status = "OPEN"
		return res, nil
	}

	for _, pp := range pps {
		entries, err := e.GetMatchingOrders(obkv, pp)
		if err != nil {
			log.Print(err)
			return nil, err
		}

		for _, o := range entries {
			entry := &types.Order{}
			err = json.Unmarshal(o, &entry)

			if err != nil {
				log.Print(err)
				return nil, err
			}

			trade, fillOrder, err := e.execute(order, entry)
			if err != nil {
				log.Print(err)
				return nil, err
			}

			order.Status = "PARTIAL_FILLED"
			res.FillStatus = "PARTIAL"
			res.Trades = append(res.Trades, trade)
			res.MatchingOrders = append(res.MatchingOrders, fillOrder)
			res.RemainingOrder.Amount = math.Sub(res.RemainingOrder.Amount, fillOrder.Amount)

			if math.IsZero(res.RemainingOrder.Amount) {
				res.FillStatus = "FULL"
				res.Order.Status = "FILLED"
				res.RemainingOrder = nil
				return res, nil
			}
		}
	}

	//TODO refactor this in a different function (make above function more clear in general)
	res.Order.Status = "PARTIAL_FILLED"
	res.FillStatus = "PARTIAL"
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
func (e *Engine) addOrder(order *types.Order) error {
	pricePointSetKey, orderHashListKey := order.GetOBKeys()
	err := e.AddToPricePointSet(pricePointSetKey, order.PricePoint.Int64())
	if err != nil {
		log.Print(err)
		return err
	}

	// Currently converting amount to int64. In the future, we need to use strings instead of int64
	amt := math.Sub(order.Amount, order.FilledAmount)
	err = e.IncrementPricePointVolume(pricePointSetKey, order.PricePoint.Int64(), amt.Int64())
	if err != nil {
		log.Print(err)
		return err
	}

	err = e.AddToOrderMap(order)
	if err != nil {
		log.Print(err)
		return err
	}

	err = e.AddToPricePointHashesSet(orderHashListKey, order.CreatedAt, order.Hash)
	if err != nil {
		log.Print(err)
		return err
	}

	return nil
}

// updateOrder updates the order in redis
func (e *Engine) updateOrder(order *types.Order, tradeAmount *big.Int) error {
	stored := &types.Order{}

	pricePointSetKey, _ := order.GetOBKeys()
	stored, err := e.GetFromOrderMap(order.Hash)
	if err != nil {
		log.Print(err)
		return err
	}

	stored.FilledAmount = math.Add(stored.FilledAmount, tradeAmount)
	if math.IsZero(stored.FilledAmount) {
		stored.Status = "OPEN"
	} else if math.IsSmallerThan(stored.FilledAmount, stored.Amount) {
		stored.Status = "PARTIAL_FILLED"
	} else {
		stored.Status = "FILLED"
	}

	err = e.AddToOrderMap(stored)
	if err != nil {
		log.Print(err)
		return err
	}

	// Currently converting amount to int64. In the future, we need to use strings instead of int64
	err = e.IncrementPricePointVolume(pricePointSetKey, order.PricePoint.Int64(), math.Neg(tradeAmount).Int64())
	if err != nil {
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

	stored, err := e.GetFromOrderMap(hash)
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

	err = e.AddToOrderMap(stored)
	if err != nil {
		log.Print(err)
		return err
	}

	pricePointSetKey, _ := stored.GetOBKeys()

	// Currently converting amount to int64. In the future, we need to use strings instead of int64
	err = e.IncrementPricePointVolume(pricePointSetKey, stored.PricePoint.Int64(), math.Neg(amount).Int64())
	if err != nil {
		log.Print(err)
		return err
	}

	return nil
}

// deleteOrder deletes the order in redis
func (e *Engine) deleteOrder(order *types.Order, tradeAmount *big.Int) (err error) {
	pricePointSetKey, orderHashListKey := order.GetOBKeys()

	vol, err := e.GetPricePointVolume(pricePointSetKey, order.PricePoint.Int64())
	if err != nil {
		log.Print(err)
		return
	}

	if math.IsEqual(math.ToBigInt(vol), tradeAmount) {
		err := e.RemoveFromPricePointSet(pricePointSetKey, order.PricePoint.Int64())
		if err != nil {
			log.Print(err)
			return err
		}

		err = e.RemoveFromPricePointHashesSet(orderHashListKey, order.Hash)
		if err != nil {
			log.Print(err)
			return err
		}

		err = e.DeletePricePointVolume(pricePointSetKey, order.PricePoint.Int64())
		if err != nil {
			log.Print(err)
			return err
		}

		err = e.RemoveFromOrderMap(order.Hash)
		if err != nil {
			log.Print(err)
			return err
		}

	} else {
		err := e.AddToPricePointSet(pricePointSetKey, order.PricePoint.Int64())
		if err != nil {
			log.Print(err)
			return err
		}

		// Currently converting amount to int64. In the future, we need to use strings instead of int64
		err = e.DecrementPricePointVolume(pricePointSetKey, order.PricePoint.Int64(), tradeAmount.Int64())
		if err != nil {
			log.Print(err)
			return err
		}

		err = e.RemoveFromOrderMap(order.Hash)
		if err != nil {
			log.Print(err)
			return err
		}

		err = e.RemoveFromPricePointHashesSet(orderHashListKey, order.Hash)
		if err != nil {
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

func (e *Engine) CancelTrades(orders []*types.Order, amount []*big.Int) error {
	e.mutex.Lock()
	defer e.mutex.Unlock()

	for _, o := range orders {
		o.Status = "PARTIAL_FILLED"
		o.FilledAmount = math.Sub(o.FilledAmount, o.Amount)
		if math.IsZero(o.FilledAmount) {
			o.Status = "OPEN"
		}

		if !e.redisConn.Exists(o.Hash.Hex()) {
			if err := e.addOrder(o); err != nil {
				log.Print(err)
				return err
			}
		} else {
			err := e.updateOrder(o, math.Neg(o.Amount))
			if err != nil {
				log.Print(err)
				return err
			}
		}
	}

	return nil
}

// CancelOrder is used to cancel the order from orderbook
func (e *Engine) CancelOrder(order *types.Order) (*types.EngineResponse, error) {
	e.mutex.Lock()
	defer e.mutex.Unlock()

	stored, err := e.GetFromOrderMap(order.Hash)
	if err != nil {
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
		FillStatus:     "CANCELLED",
		Order:          stored,
		RemainingOrder: nil,
		Trades:         nil,
		MatchingOrders: nil,
	}

	return engineResponse, nil
}

// GetPricePoints returns the pricepoints matching a certain (pair, pricepoint)
func (e *Engine) GetMatchingBuyPricePoints(obKey string, pricePoint int64) ([]int64, error) {
	pps, err := e.redisConn.ZRangeByLexInt(obKey, "-", "["+utils.UintToPaddedString(pricePoint))
	if err != nil {
		log.Print(err)
		return nil, err
	}

	return pps, nil
}

func (e *Engine) GetMatchingSellPricePoints(obkv string, pricePoint int64) ([]int64, error) {
	pps, err := e.redisConn.ZRevRangeByLexInt(obkv, "+", "["+utils.UintToPaddedString(pricePoint))
	if err != nil {
		log.Print(err)
		return nil, err
	}

	return pps, nil
}

func (e *Engine) GetPricePointVolume(pricePointSetKey string, pricePoint int64) (string, error) {
	vol, err := e.redisConn.GetValue(pricePointSetKey + "::book::" + utils.UintToPaddedString(pricePoint))
	if err != nil {
		log.Print(err)
		return "", err
	}

	return vol, nil
}

func (e *Engine) GetFromOrderMap(hash common.Hash) (*types.Order, error) {
	o := &types.Order{}

	serialized, err := e.redisConn.GetValue(hash.Hex())
	if err != nil {
		log.Print(err)
		return nil, err
	}

	err = json.Unmarshal([]byte(serialized), &o)
	if err != nil {
		return nil, err
	}

	return o, nil
}

// GetPricePointOrders returns the orders hashes for a (pair, pricepoint)
func (e *Engine) GetMatchingOrders(obKey string, pricePoint int64) ([][]byte, error) {
	orders, err := e.redisConn.Sort(obKey+"::"+utils.UintToPaddedString(pricePoint), "", true, false, "*")
	if err != nil {
		return nil, err
	}

	return orders, nil
}

// AddPricePointToSet
func (e *Engine) AddToPricePointSet(pricePointSetKey string, pricePoint int64) error {
	err := e.redisConn.ZAdd(pricePointSetKey, 0, utils.UintToPaddedString(pricePoint))
	if err != nil {
		log.Print(err)
		return err
	}

	return nil
}

// RemoveFromPricePointSet
func (e *Engine) RemoveFromPricePointSet(pricePointSetKey string, pricePoint int64) error {
	err := e.redisConn.ZRem(pricePointSetKey, utils.UintToPaddedString(pricePoint))
	if err != nil {
		log.Print(err)
		return err
	}

	return nil
}

// AddPricePointHashesSet
func (e *Engine) AddToPricePointHashesSet(orderHashListKey string, createdAt time.Time, hash common.Hash) error {
	err := e.redisConn.ZAdd(orderHashListKey, createdAt.Unix(), hash.Hex())
	if err != nil {
		return err
	}

	return nil
}

// RemoveFromPricePointHashesSet
func (e *Engine) RemoveFromPricePointHashesSet(orderHashListKey string, hash common.Hash) error {
	err := e.redisConn.ZRem(orderHashListKey, hash.Hex())
	if err != nil {
		return err
	}

	return nil
}

// IncrementPricePointVolume increases the value of a certain pricepoint at a certain volume
func (e *Engine) IncrementPricePointVolume(pricePointSetKey string, pricePoint int64, amount int64) error {
	_, err := e.redisConn.IncrBy(pricePointSetKey+"::book::"+utils.UintToPaddedString(pricePoint), amount)
	if err != nil {
		log.Print(err)
		return err
	}

	return nil
}

// DecrementPricePoint
func (e *Engine) DecrementPricePointVolume(pricePointSetKey string, pricePoint int64, amount int64) error {
	_, err := e.redisConn.IncrBy(pricePointSetKey+"::book::"+utils.UintToPaddedString(pricePoint), -amount)
	if err != nil {
		log.Print(err)
		return err
	}

	return nil
}

// DeletePricePoint
func (e *Engine) DeletePricePointVolume(pricePointSetKey string, pricePoint int64) error {
	err := e.redisConn.Del(pricePointSetKey + "::book::" + utils.UintToPaddedString(pricePoint))
	if err != nil {
		log.Print(err)
		return err
	}

	return nil
}

// AddToOrderMap
func (e *Engine) AddToOrderMap(o *types.Order) error {
	bytes, err := json.Marshal(o)
	if err != nil {
		log.Print(err)
		return err
	}

	decoded := &types.Order{}
	json.Unmarshal(bytes, decoded)

	err = e.redisConn.Set(o.Hash.Hex(), string(bytes))
	if err != nil {
		log.Print(err)
		return err
	}

	return nil
}

// RemoveFromOrderMap
func (e *Engine) RemoveFromOrderMap(hash common.Hash) error {
	err := e.redisConn.Del(hash.Hex())
	if err != nil {
		log.Print(err)
		return err
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
