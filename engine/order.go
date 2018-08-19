package engine

import (
	"encoding/json"
	"errors"
	"log"

	"github.com/Proofsuite/amp-matching-engine/types"
	"github.com/Proofsuite/amp-matching-engine/utils"
	"github.com/Proofsuite/go-ethereum/common"
	"github.com/gomodule/redigo/redis"
)

// FillOrder is structure holds the matching order and
// the amount that has been filled by taker order
type FillOrder struct {
	Amount int64
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
			log.Printf("\n%s\n", err)
			return err
		}
	} else if order.Side == "BUY" {
		resp, err = e.buyOrder(order)
		if err != nil {
			log.Printf("\n%s\n", err)
			return err
		}
	}

	// Note: Plug the option for orders like FOC, Limit
	err = e.publishEngineResponse(resp)
	if err != nil {
		log.Printf("\npublishEngineResponse XXXXXXX\n%s\nXXXXXXX publishEngineResponse\n", err)
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

	remOrder := *order
	resp.Trades = make([]*types.Trade, 0)
	resp.RemainingOrder = &remOrder
	resp.MatchingOrders = make([]*FillOrder, 0)
	oskv := order.GetOBMatchKey()

	// GET Range of sellOrder between minimum Sell order and order.Price
	orders, err := redis.Values(e.redisConn.Do("ZRANGEBYLEX", oskv, "-", "["+utils.UintToPaddedString(order.Price))) // "ZRANGEBYLEX" key min max
	if err != nil {
		log.Printf("ZRANGEBYLEX: %s\n", err)
		return nil, err
	}

	priceRange := make([]int64, 0)
	if err := redis.ScanSlice(orders, &priceRange); err != nil {
		log.Printf("Scan %s\n", err)
		return nil, err
	}

	if len(priceRange) == 0 {
		resp.FillStatus = NOMATCH
		order.Status = "OPEN"
		e.addOrder(order)
		return resp, nil
	}

	for _, pr := range priceRange {
		bookEntries, err := redis.ByteSlices(e.redisConn.Do("SORT", oskv+"::"+utils.UintToPaddedString(pr), "GET", oskv+"::"+utils.UintToPaddedString(pr)+"::*", "ALPHA")) // "ZREVRANGEBYLEX" key max min
		if err != nil {
			log.Printf("LRANGE: %s\n", err)
			return nil, err
		}

		for _, o := range bookEntries {
			var bookEntry *types.Order
			err = json.Unmarshal(o, &bookEntry)
			if err != nil {
				log.Printf("json.Unmarshal: %s\n", err)
				return nil, err
			}

			trade, fillOrder, err := e.execute(order, bookEntry)
			if err != nil {
				log.Printf("Error Executing Order: %s\n", err)
				return nil, err
			}

			resp.Trades = append(resp.Trades, trade)
			resp.MatchingOrders = append(resp.MatchingOrders, fillOrder)
			resp.RemainingOrder.Amount = resp.RemainingOrder.Amount - fillOrder.Amount

			if resp.RemainingOrder.Amount == 0 {
				resp.FillStatus = FULL
				resp.Order.Status = "FILLED"
				resp.RemainingOrder = nil
				return resp, nil
			}

			resp.Order.Status = "PARTIAL_FILLED"
		}
	}

	return resp, nil
}

// sellOrder is triggered when a sell order comes in, it fetches the bid list
// from orderbook. First it checks ths price point list to check whether the order can be matched
// or not, if there are pricepoints that can satisfy the order then corresponding list of orders
// are fetched and trade is executed
func (e *Resource) sellOrder(order *types.Order) (resp *Response, err error) {
	resp = &Response{
		Order:      order,
		FillStatus: NOMATCH,
	}

	remOrder := *order
	resp.Trades = make([]*types.Trade, 0)
	resp.RemainingOrder = &remOrder
	resp.MatchingOrders = make([]*FillOrder, 0)

	obkv := order.GetOBMatchKey()
	// GET Range of sellOrder between minimum Sell order and order.Price
	orders, err := redis.Values(e.redisConn.Do("ZREVRANGEBYLEX", obkv, "+", "["+utils.UintToPaddedString(order.Price))) // "ZREVRANGEBYLEX" key max min
	if err != nil {
		log.Printf("ZREVRANGEBYLEX: %s\n", err)
		return nil, err
	}

	priceRange := make([]int64, 0)
	if err := redis.ScanSlice(orders, &priceRange); err != nil {
		log.Printf("Scan %s\n", err)
		return nil, err
	}

	if len(priceRange) == 0 {
		resp.FillStatus = NOMATCH
		resp.RemainingOrder = &types.Order{}
		e.addOrder(order)
		order.Status = "OPEN"
		return resp, nil
	}

	for _, pr := range priceRange {
		bookEntries, err := redis.ByteSlices(e.redisConn.Do("SORT", obkv+"::"+utils.UintToPaddedString(pr), "GET", obkv+"::"+utils.UintToPaddedString(pr)+"::*", "ALPHA")) // "ZREVRANGEBYLEX" key max min
		if err != nil {
			log.Printf("SORT: %s\n", err)
			return nil, err
		}
		for _, o := range bookEntries {
			var bookEntry *types.Order
			err = json.Unmarshal(o, &bookEntry)
			if err != nil {
				log.Printf("json.Unmarshal: %s\n", err)
				return nil, err
			}

			trade, fillOrder, err := e.execute(order, bookEntry)
			if err != nil {
				log.Printf("Error Executing Order: %s\n", err)
				return nil, err
			}

			resp.Trades = append(resp.Trades, trade)
			resp.MatchingOrders = append(resp.MatchingOrders, fillOrder)
			resp.RemainingOrder.Amount = resp.RemainingOrder.Amount - fillOrder.Amount

			if resp.RemainingOrder.Amount == 0 {
				resp.FillStatus = FULL
				resp.Order.Status = "FILLED"
				resp.RemainingOrder = nil
				return resp, nil
			}
			resp.Order.Status = "PARTIAL_FILLED"

		}
	}
	return
}

// addOrder adds an order to redis
func (e *Resource) addOrder(order *types.Order) error {
	ssKey, listKey := order.GetOBKeys()
	_, err := e.redisConn.Do("ZADD", ssKey, "NX", 0, utils.UintToPaddedString(order.Price)) // Add price point to order book
	if err != nil {
		log.Printf("ZADD: %s", err)
		return err
	}

	// fmt.Printf("ZADD: %s\n", res)
	_, err = e.redisConn.Do("INCRBY", ssKey+"::book::"+utils.UintToPaddedString(order.Price), order.Amount-order.FilledAmount) // Add price point to order book
	if err != nil {
		log.Printf("INCRBY: %s", err)
		return err
	}

	// Add order to list
	orderAsBytes, err := json.Marshal(order)
	if err != nil {
		log.Printf("orderAsBytes: %s", err)
		return err
	}

	_, err = e.redisConn.Do("SET", listKey+"::"+order.Hash.Hex(), string(orderAsBytes))
	if err != nil {
		log.Printf("SET: %s", err)
		return err
	}

	// Add order reference to price sorted set
	_, err = e.redisConn.Do("ZADD", listKey, "NX", order.CreatedAt.Unix(), order.Hash.Hex())
	if err != nil {
		log.Printf("ZADD: %s", err)
		return err
	}

	// fmt.Printf("ZADD: %s\n", res)
	return nil
}

// updateOrder updates the order in redis
func (e *Resource) updateOrder(order *types.Order, tradeAmount int64) error {
	ssKey, listKey := order.GetOBKeys()
	var storedOrder *types.Order
	storedOrderAsBytes, err := redis.Bytes(e.redisConn.Do("GET", listKey+"::"+order.Hash.Hex()))
	if err != nil {
		log.Printf("orderAsBytes: %s", err)
		return err
	}

	json.Unmarshal(storedOrderAsBytes, &storedOrder)

	storedOrder.FilledAmount = storedOrder.FilledAmount + tradeAmount
	if storedOrder.FilledAmount == 0 {
		storedOrder.Status = "OPEN"
	} else if storedOrder.FilledAmount < storedOrder.Amount {
		storedOrder.Status = "PARTIAL_FILLED"
	} else {
		storedOrder.Status = "FILLED"
	}

	// Add order to list
	orderAsBytes, err := json.Marshal(storedOrder)
	if err != nil {
		log.Printf("orderAsBytes: %s", err)
		return err
	}

	_, err = e.redisConn.Do("SET", listKey+"::"+order.Hash.Hex(), string(orderAsBytes))
	if err != nil {
		log.Printf("SET: %s", err)
		return err
	}

	_, err = e.redisConn.Do("INCRBY", ssKey+"::book::"+utils.UintToPaddedString(order.Price), -1*tradeAmount) // Add price point to order book
	if err != nil {
		log.Printf("INCRBY: %s", err)
		return err
	}

	return nil
}

// updateOrderAmount is a less general version of update order that updates the order filled amount after
// a trade or recover orders
// EXPERIMENTAL
func (e *Resource) updateOrderAmount(hash common.Hash, amount int64) error {
	stored := &types.Order{}
	bytes, err := redis.Bytes(e.redisConn.Do("GET", hash.Hex()))
	if err != nil {
		return err
	}

	err = json.Unmarshal(bytes, stored)
	if err != nil {
		return err
	}

	stored.FilledAmount += amount
	if stored.FilledAmount == 0 {
		stored.Status = "OPEN"
	} else if stored.FilledAmount < stored.Amount {
		stored.Status = "PARTIAL_FILLED"
	} else {
		stored.Status = "FILLED"
	}

	bytes, err = json.Marshal(stored)
	if err != nil {
		return err
	}

	_, err = e.redisConn.Do("SET", hash.Hex(), string(bytes))
	if err != nil {
		return err
	}

	ssKey, _ := stored.GetOBKeys()
	_, err = e.redisConn.Do("INCRBY", ssKey+"::book::"+utils.UintToPaddedString(stored.Price), -amount)
	if err != nil {
		return err
	}

	return nil
}

// deleteOrder deletes the order in redis
func (e *Resource) deleteOrder(order *types.Order, tradeAmount int64) (err error) {

	ssKey, listKey := order.GetOBKeys()
	remVolume, err := redis.Int64(e.redisConn.Do("GET", ssKey+"::book::"+utils.UintToPaddedString(order.Price)))
	if err != nil {
		log.Printf("GET remVolume: %s", err)
		return
	}

	if remVolume == tradeAmount {
		_, err := e.redisConn.Do("ZREM", ssKey, "NX", 0, utils.UintToPaddedString(order.Price)) // Add price point to order book
		if err != nil {
			log.Printf("ZREM: %s", err)
			return err
		}
		// fmt.Printf("ZREM: %s\n", res)
		_, err = e.redisConn.Do("DEL", ssKey+"::book::"+utils.UintToPaddedString(order.Price)) // Add price point to order book
		if err != nil {
			log.Printf("DEL: %s", err)
			return err
		}
		// fmt.Printf("DEL: %s\n", res)

		_, err = e.redisConn.Do("DEL", listKey+"::"+order.Hash.Hex())
		if err != nil {
			log.Printf("DEL: %s", err)
			return err
		}
		// Add order reference to price sorted set
		_, err = e.redisConn.Do("ZREM", listKey, order.Hash.Hex())
		if err != nil {
			log.Printf("ZREM: %s", err)
			return err
		}

		// fmt.Printf("ZREM: %s\n", res)
	} else {
		_, err := e.redisConn.Do("ZADD", ssKey, "NX", 0, utils.UintToPaddedString(order.Price)) // Add price point to order book
		if err != nil {
			log.Printf("ZADD: %s", err)
			return err
		}

		// fmt.Printf("ZADD: %s\n", res)
		_, err = e.redisConn.Do("INCRBY", ssKey+"::book::"+utils.UintToPaddedString(order.Price), -1*tradeAmount) // Add price point to order book
		if err != nil {
			log.Printf("INCRBY: %s", err)
			return err
		}
		// fmt.Printf("INCRBY: %s\n", res)

		_, err = e.redisConn.Do("DEL", listKey+"::"+order.Hash.Hex())
		if err != nil {
			log.Printf("DEL: %s", err)
			return err
		}
		// Add order reference to price sorted set
		_, err = e.redisConn.Do("ZREM", listKey, order.Hash.Hex())
		if err != nil {
			log.Printf("ZREM: %s", err)
			return err
		}

		// fmt.Printf("ZREM: %s\n", res)
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
		o.Order.FilledAmount = o.Order.FilledAmount - o.Amount
		if o.Order.FilledAmount == 0 {
			o.Order.Status = "OPEN"
		}

		_, listKey := o.Order.GetOBKeys()
		res, _ := redis.Bytes(e.redisConn.Do("GET", listKey+"::"+o.Order.Hash.Hex()))
		if res == nil {
			if err := e.addOrder(o.Order); err != nil {
				return err
			}
		} else {
			if err := e.updateOrder(o.Order, -1*o.Amount); err != nil {
				return err
			}
		}
	}
	return nil
}

// RecoverOrders2 is an alternative suggestion for RecoverOrders2
// It would requires an alternative key system for redis where we store only the hash with the listkey prefix
func (e *Resource) RecoverOrders2(hashes []common.Hash, amounts []int64) error {
	for i, _ := range hashes {
		err := e.updateOrderAmount(hashes[i], amounts[i])
		if err != nil {
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
		log.Printf("GET: %s", err)
		return nil, err
	}
	if res == nil {
		return nil, errors.New("Order not found")
	}

	var storedOrder *types.Order
	if err := json.Unmarshal(res, &storedOrder); err != nil {
		log.Printf("G2ET: %s", err)
		return nil, err
	}
	if err := e.deleteOrder(order, storedOrder.Amount-storedOrder.FilledAmount); err != nil {
		log.Printf("\n%s\n", err)
		return nil, err
	}

	storedOrder.Status = "CANCELLED"

	engineResponse := &Response{
		Order:          storedOrder,
		Trades:         make([]*types.Trade, 0),
		RemainingOrder: &types.Order{},
		FillStatus:     CANCELLED,
		MatchingOrders: make([]*FillOrder, 0),
	}
	return engineResponse, nil
}
