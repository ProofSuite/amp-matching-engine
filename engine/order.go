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
	"github.com/gomodule/redigo/redis"
)

// FillOrder is structure holds the matching order and
// the amount that has been filled by taker order
type FillOrder struct {
	Amount *big.Int
	Order  *types.Order
}

// newOrder calls buyOrder/sellOrder based on type of order recieved and
// publishes the response back to rabbitmq
func (e *Resource) newOrder(order *types.Order) (err error) {
	// Attain lock on engineResource, so that recovery or cancel order function doesn't interfere
	e.mutex.Lock()
	defer e.mutex.Unlock()

	resp := &Response{}
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
func (e *Resource) buyOrder(order *types.Order) (*Response, error) {
	resp := &Response{
		Order:      order,
		FillStatus: NOMATCH,
	}

	remainingOrder := *order
	resp.RemainingOrder = &remainingOrder
	resp.Trades = make([]*types.Trade, 0)
	resp.MatchingOrders = make([]*FillOrder, 0)
	oskv := order.GetOBMatchKey()

	// GET Range of sellOrder between minimum Sell order and order.Price
	orders, err := redis.Values(e.redisConn.Do("ZRANGEBYLEX", oskv, "-", "["+utils.UintToPaddedString(order.PricePoint.Int64()))) // "ZRANGEBYLEX" key min max
	if err != nil {
		log.Print(err)
		return nil, err
	}

	priceRange := make([]int64, 0)
	if err := redis.ScanSlice(orders, &priceRange); err != nil {
		log.Print(err)
		return nil, err
	}

	if len(priceRange) == 0 {
		resp.FillStatus = NOMATCH
		resp.RemainingOrder = nil
		order.Status = "OPEN"
		e.addOrder(order)
		return resp, nil
	}

	for _, pr := range priceRange {
		bookEntries, err := redis.ByteSlices(e.redisConn.Do("SORT", oskv+"::"+utils.UintToPaddedString(pr), "GET", oskv+"::"+utils.UintToPaddedString(pr)+"::*", "ALPHA")) // "ZREVRANGEBYLEX" key max min
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
				resp.FillStatus = FULL
				resp.Order.Status = "FILLED"
				resp.RemainingOrder = nil
				return resp, nil
			}
		}
	}

	resp.Order.Status = "PARTIAL_FILLED"
	resp.FillStatus = PARTIAL

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
func (e *Resource) sellOrder(order *types.Order) (*Response, error) {
	resp := &Response{
		Order:      order,
		FillStatus: NOMATCH,
	}

	remOrder := *order
	resp.Trades = make([]*types.Trade, 0)
	resp.RemainingOrder = &remOrder
	resp.MatchingOrders = make([]*FillOrder, 0)
	obkv := order.GetOBMatchKey()

	// GET Range of sellOrder between minimum Sell order and order.Price
	orders, err := redis.Values(e.redisConn.Do("ZREVRANGEBYLEX", obkv, "+", "["+utils.UintToPaddedString(order.PricePoint.Int64()))) // "ZREVRANGEBYLEX" key max min
	if err != nil {
		log.Print(err)
		return nil, err
	}

	priceRange := make([]int64, 0)
	if err := redis.ScanSlice(orders, &priceRange); err != nil {
		log.Print(err)
		return nil, err
	}

	if len(priceRange) == 0 {
		resp.FillStatus = NOMATCH
		resp.RemainingOrder = nil
		e.addOrder(order)
		order.Status = "OPEN"
		return resp, nil
	}

	for _, pr := range priceRange {
		bookEntries, err := redis.ByteSlices(e.redisConn.Do("SORT", obkv+"::"+utils.UintToPaddedString(pr), "GET", obkv+"::"+utils.UintToPaddedString(pr)+"::*", "ALPHA")) // "ZREVRANGEBYLEX" key max min
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
			resp.FillStatus = PARTIAL
			resp.Trades = append(resp.Trades, trade)
			resp.MatchingOrders = append(resp.MatchingOrders, fillOrder)
			resp.RemainingOrder.Amount = math.Sub(resp.RemainingOrder.Amount, fillOrder.Amount)

			if math.IsZero(resp.RemainingOrder.Amount) {
				resp.FillStatus = FULL
				resp.Order.Status = "FILLED"
				resp.RemainingOrder = nil
				return resp, nil
			}
		}
	}

	resp.Order.Status = "PARTIAL_FILLED"
	resp.FillStatus = PARTIAL

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
func (e *Resource) addOrder(order *types.Order) error {
	ssKey, listKey := order.GetOBKeys()
	_, err := e.redisConn.Do("ZADD", ssKey, "NX", 0, utils.UintToPaddedString(order.PricePoint.Int64())) // Add price point to order book
	if err != nil {
		log.Print(err)
		return err
	}

	// Currently converting amount to int64. In the future, we need to use strings instead of int64
	amt := math.Sub(order.Amount, order.FilledAmount)
	_, err = e.redisConn.Do("INCRBY", ssKey+"::book::"+utils.UintToPaddedString(order.PricePoint.Int64()), amt.Int64()) // Add price point to order book
	if err != nil {
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

	_, err = e.redisConn.Do("SET", listKey+"::"+order.Hash.Hex(), string(orderAsBytes))
	if err != nil {
		log.Print(err)
		return err
	}

	// Add order reference to price sorted set
	_, err = e.redisConn.Do("ZADD", listKey, "NX", order.CreatedAt.Unix(), order.Hash.Hex())
	if err != nil {
		log.Print(err)
		return err
	}

	return nil
}

// updateOrder updates the order in redis
func (e *Resource) updateOrder(order *types.Order, tradeAmount *big.Int) error {
	stored := &types.Order{}

	ssKey, listKey := order.GetOBKeys()
	bytes, err := redis.Bytes(e.redisConn.Do("GET", listKey+"::"+order.Hash.Hex()))
	if err != nil {
		log.Print(err)
		return err
	}

	json.Unmarshal(bytes, &stored)

	stored.FilledAmount = math.Add(stored.FilledAmount, tradeAmount)
	if math.IsZero(stored.FilledAmount) {
		stored.Status = "OPEN"
	} else if math.IsSmallerThan(stored.FilledAmount, stored.Amount) {
		stored.Status = "PARTIAL_FILLED"
	} else {
		stored.Status = "FILLED"
	}

	// Add order to list
	bytes, err = json.Marshal(stored)
	if err != nil {
		log.Print(err)
		return err
	}

	_, err = e.redisConn.Do("SET", listKey+"::"+order.Hash.Hex(), string(bytes))
	if err != nil {
		log.Print(err)
		return err
	}

	// Currently converting amount to int64. In the future, we need to use strings instead of int64
	_, err = e.redisConn.Do("INCRBY", ssKey+"::book::"+utils.UintToPaddedString(order.PricePoint.Int64()), math.Neg(tradeAmount))
	if err != nil {
		log.Print(err)
		return err
	}

	return nil
}

// updateOrderAmount is a less general version of update order that updates the order filled amount after
// a trade or recover orders
// EXPERIMENTAL
func (e *Resource) updateOrderAmount(hash common.Hash, amount *big.Int) error {
	stored := &types.Order{}
	bytes, err := redis.Bytes(e.redisConn.Do("GET", hash.Hex()))
	if err != nil {
		return err
	}

	err = json.Unmarshal(bytes, stored)
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

	bytes, err = json.Marshal(stored)
	if err != nil {
		log.Print(err)
		return err
	}

	_, err = e.redisConn.Do("SET", hash.Hex(), string(bytes))
	if err != nil {
		log.Print(err)
		return err
	}

	ssKey, _ := stored.GetOBKeys()

	// Currently converting amount to int64. In the future, we need to use strings instead of int64
	_, err = e.redisConn.Do("INCRBY", ssKey+"::book::"+utils.UintToPaddedString(stored.PricePoint.Int64()), math.Neg(amount))
	if err != nil {
		log.Print(err)
		return err
	}

	return nil
}

// deleteOrder deletes the order in redis
func (e *Resource) deleteOrder(order *types.Order, tradeAmount *big.Int) (err error) {
	ssKey, listKey := order.GetOBKeys()
	remVolume, err := redis.String(e.redisConn.Do("GET", ssKey+"::book::"+utils.UintToPaddedString(order.PricePoint.Int64())))
	if err != nil {
		log.Print(err)
		return
	}

	if math.IsEqual(math.ToBigInt(remVolume), tradeAmount) {
		_, err := e.redisConn.Do("ZREM", ssKey, "NX", 0, utils.UintToPaddedString(order.PricePoint.Int64()))
		if err != nil {
			log.Print(err)
			return err
		}
		// fmt.Printf("ZREM: %s\n", res)
		_, err = e.redisConn.Do("DEL", ssKey+"::book::"+utils.UintToPaddedString(order.PricePoint.Int64()))
		if err != nil {
			log.Print(err)
			return err
		}
		// fmt.Printf("DEL: %s\n", res)

		_, err = e.redisConn.Do("DEL", listKey+"::"+order.Hash.Hex())
		if err != nil {
			log.Print(err)
			return err
		}
		// Add order reference to price sorted set
		_, err = e.redisConn.Do("ZREM", listKey, order.Hash.Hex())
		if err != nil {
			log.Print(err)
			return err
		}

	} else {
		_, err := e.redisConn.Do("ZADD", ssKey, "NX", 0, utils.UintToPaddedString(order.PricePoint.Int64()))
		if err != nil {
			log.Print(err)
			return err
		}

		// Currently converting amount to int64. In the future, we need to use strings instead of int64
		_, err = e.redisConn.Do("INCRBY", ssKey+"::book::"+utils.UintToPaddedString(order.PricePoint.Int64()), math.Neg(tradeAmount))
		if err != nil {
			log.Print(err)
			return err
		}

		_, err = e.redisConn.Do("DEL", listKey+"::"+order.Hash.Hex())
		if err != nil {
			log.Print(err)
			return err
		}
		// Add order reference to price sorted set
		_, err = e.redisConn.Do("ZREM", listKey, order.Hash.Hex())
		if err != nil {
			log.Print(err)
			return err
		}
	}

	return
}

// RecoverOrders is responsible for recovering the orders that failed to execute after matching
// Orders are updated or added to orderbook based on whether that order exists in orderbook or not.
func (e *Resource) RecoverOrders(orders []*FillOrder) error {
	e.mutex.Lock()
	defer e.mutex.Unlock()
	for _, o := range orders {

		// update order's filled amount and status before updating in redis
		o.Order.Status = "PARTIAL_FILLED"
		o.Order.FilledAmount = math.Sub(o.Order.FilledAmount, o.Amount)
		if math.IsZero(o.Order.FilledAmount) {
			o.Order.Status = "OPEN"
		}

		_, listKey := o.Order.GetOBKeys()
		res, _ := redis.Bytes(e.redisConn.Do("GET", listKey+"::"+o.Order.Hash.Hex()))
		if res == nil {
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

// RecoverOrders2 is an alternative suggestion for RecoverOrders2
// It would requires an alternative key system for redis where we store only the hash with the listkey prefix
func (e *Resource) RecoverOrders2(hashes []common.Hash, amounts []*big.Int) error {
	for i, _ := range hashes {
		err := e.updateOrderAmount(hashes[i], amounts[i])
		if err != nil {
			log.Print(err)
			return err
		}
	}

	return nil
}

// CancelOrder is used to cancel the order from orderbook
func (e *Resource) CancelOrder(order *types.Order) (*Response, error) {
	e.mutex.Lock()
	defer e.mutex.Unlock()

	_, listKey := order.GetOBKeys()
	res, err := redis.Bytes(e.redisConn.Do("GET", listKey+"::"+order.Hash.Hex()))
	if err != nil {
		log.Print(err)
		return nil, err
	}

	if res == nil {
		return nil, errors.New("Order not found")
	}

	var stored *types.Order
	if err := json.Unmarshal(res, &stored); err != nil {
		log.Print(err)
		return nil, err
	}

	amt := math.Sub(stored.Amount, stored.FilledAmount)
	if err := e.deleteOrder(order, amt); err != nil {
		log.Print(err)
		return nil, err
	}

	stored.Status = "CANCELLED"

	engineResponse := &Response{
		Order:          stored,
		Trades:         make([]*types.Trade, 0),
		RemainingOrder: &types.Order{},
		FillStatus:     CANCELLED,
		MatchingOrders: make([]*FillOrder, 0),
	}
	return engineResponse, nil
}
