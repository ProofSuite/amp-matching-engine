package engine

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/Proofsuite/amp-matching-engine/types"
	"github.com/Proofsuite/amp-matching-engine/utils"
	"github.com/gomodule/redigo/redis"
)

// FillOrder is structure holds the matching order and
// the amount that has been filled by taker order
type FillOrder struct {
	Amount int64
	Order  *types.Order
}

// matchOrder calls buyOrder/sellOrder based on type of order recieved and
// publishes the response back to rabbitmq
func (e *Resource) matchOrder(order *types.Order) (err error) {

	// Attain lock on engineResource, so that recovery or cancel order function doesn't interfere

	e.mutex.Lock()
	defer e.mutex.Unlock()
	var engineResponse *Response
	if order.Side == types.SELL {
		engineResponse, err = e.sellOrder(order)
	} else if order.Side == types.BUY {
		engineResponse, err = e.buyOrder(order)
	}
	if err != nil {
		log.Printf("\n%s\n", err)
		return err
	}

	// Note: Plug the option for orders like FOC, Limit

	e.publishEngineResponse(engineResponse)
	if err != nil {
		log.Printf("\npublishEngineResponse XXXXXXX\n%s\nXXXXXXX publishEngineResponse\n", err)
	}
	return
}

// buyOrder is triggered when a buy order comes in, it fetches the ask list
// from orderbook. First it checks ths price point list to check whether the order can be matched
// or not, if there are pricepoints that can satisfy the order then corresponding list of orders
// are fetched and trade is executed
func (e *Resource) buyOrder(order *types.Order) (engineResponse *Response, err error) {

	engineResponse = &Response{
		Order:      order,
		FillStatus: NOMATCH,
	}
	engineResponse.Trades = make([]*types.Trade, 0)
	remOrder := *order
	engineResponse.RemainingOrder = &remOrder
	engineResponse.MatchingOrders = make([]*FillOrder, 0)

	oskv := order.GetOBMatchKey()

	// GET Range of sellOrder between minimum Sell order and order.Price
	orders, err := redis.Values(e.redisConn.Do("ZRANGEBYLEX", oskv, "-", "["+utils.UintToPaddedString(order.Price))) // "ZRANGEBYLEX" key min max
	if err != nil {
		log.Printf("ZRANGEBYLEX: %s\n", err)
		return
	}
	priceRange := make([]int64, 0)
	if err := redis.ScanSlice(orders, &priceRange); err != nil {
		log.Printf("Scan %s\n", err)
	}

	if len(priceRange) == 0 {
		engineResponse.FillStatus = NOMATCH
		e.addOrder(order)
		engineResponse.RemainingOrder = &types.Order{}
		order.Status = types.OPEN
	} else {
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
				engineResponse.Trades = append(engineResponse.Trades, trade)
				engineResponse.MatchingOrders = append(engineResponse.MatchingOrders, fillOrder)
				engineResponse.RemainingOrder.Amount = engineResponse.RemainingOrder.Amount - fillOrder.Amount

				if engineResponse.RemainingOrder.Amount == 0 {
					engineResponse.FillStatus = FULL
					engineResponse.Order.Status = types.FILLED
					engineResponse.RemainingOrder = nil
					return engineResponse, nil
				}
				engineResponse.Order.Status = types.PARTIALFILLED
			}
		}
	}
	return
}

// sellOrder is triggered when a sell order comes in, it fetches the bid list
// from orderbook. First it checks ths price point list to check whether the order can be matched
// or not, if there are pricepoints that can satisfy the order then corresponding list of orders
// are fetched and trade is executed
func (e *Resource) sellOrder(order *types.Order) (engineResponse *Response, err error) {
	engineResponse = &Response{
		Order:      order,
		FillStatus: NOMATCH,
	}
	engineResponse.Trades = make([]*types.Trade, 0)
	remOrder := *order
	engineResponse.RemainingOrder = &remOrder
	engineResponse.MatchingOrders = make([]*FillOrder, 0)

	obkv := order.GetOBMatchKey()
	// GET Range of sellOrder between minimum Sell order and order.Price
	orders, err := redis.Values(e.redisConn.Do("ZREVRANGEBYLEX", obkv, "+", "["+utils.UintToPaddedString(order.Price))) // "ZREVRANGEBYLEX" key max min
	if err != nil {
		log.Printf("ZREVRANGEBYLEX: %s\n", err)
		return
	}

	priceRange := make([]int64, 0)
	if err := redis.ScanSlice(orders, &priceRange); err != nil {
		log.Printf("Scan %s\n", err)
	}

	if len(priceRange) == 0 {
		engineResponse.FillStatus = NOMATCH
		engineResponse.RemainingOrder = &types.Order{}
		e.addOrder(order)
		order.Status = types.OPEN
	} else {
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
				engineResponse.Trades = append(engineResponse.Trades, trade)
				engineResponse.MatchingOrders = append(engineResponse.MatchingOrders, fillOrder)
				engineResponse.RemainingOrder.Amount = engineResponse.RemainingOrder.Amount - fillOrder.Amount

				if engineResponse.RemainingOrder.Amount == 0 {
					engineResponse.FillStatus = FULL
					engineResponse.Order.Status = types.FILLED
					engineResponse.RemainingOrder = nil
					return engineResponse, nil
				}
				engineResponse.Order.Status = types.PARTIALFILLED

			}
		}
	}
	return
}

// addOrder adds an order to redis
func (e *Resource) addOrder(order *types.Order) error {

	ssKey, listKey := order.GetOBKeys()
	res, err := e.redisConn.Do("ZADD", ssKey, "NX", 0, utils.UintToPaddedString(order.Price)) // Add price point to order book
	if err != nil {
		log.Printf("ZADD: %s", err)
		return err
	}
	fmt.Printf("ZADD: %s\n", res)
	res, err = e.redisConn.Do("INCRBY", ssKey+"::book::"+utils.UintToPaddedString(order.Price), order.Amount-order.FilledAmount) // Add price point to order book
	if err != nil {
		log.Printf("INCRBY: %s", err)
		return err
	}
	fmt.Printf("INCRBY: %s\n", res)

	// Add order to list
	orderAsBytes, err := json.Marshal(order)
	if err != nil {
		log.Printf("orderAsBytes: %s", err)
		return err
	}
	res, err = e.redisConn.Do("SET", listKey+"::"+order.Hash, string(orderAsBytes))
	if err != nil {
		log.Printf("SET: %s", err)
		return err
	}
	// Add order reference to price sorted set
	res, err = e.redisConn.Do("ZADD", listKey, "NX", order.CreatedAt.Unix(), order.Hash)
	if err != nil {
		log.Printf("ZADD: %s", err)
		return err
	}

	fmt.Printf("ZADD: %s\n", res)

	return nil
}

// updateOrder updates the order in redis
func (e *Resource) updateOrder(order *types.Order, tradeAmount int64) error {

	ssKey, listKey := order.GetOBKeys()
	var storedOrder *types.Order
	storedOrderAsBytes, err := redis.Bytes(e.redisConn.Do("GET", listKey+"::"+order.Hash))
	if err != nil {
		log.Printf("orderAsBytes: %s", err)
		return err
	}
	json.Unmarshal(storedOrderAsBytes, &storedOrder)

	storedOrder.FilledAmount = storedOrder.FilledAmount + tradeAmount
	if storedOrder.FilledAmount == 0 {
		storedOrder.Status = types.OPEN
	}

	// Add order to list
	orderAsBytes, err := json.Marshal(storedOrder)
	if err != nil {
		log.Printf("orderAsBytes: %s", err)
		return err
	}

	res, err := e.redisConn.Do("SET", listKey+"::"+order.Hash, string(orderAsBytes))
	if err != nil {
		log.Printf("SET: %s", err)
		return err
	}

	res, err = e.redisConn.Do("INCRBY", ssKey+"::book::"+utils.UintToPaddedString(order.Price), -1*tradeAmount) // Add price point to order book
	if err != nil {
		log.Printf("INCRBY: %s", err)
		return err
	}
	fmt.Printf("INCRBY: %s\n", res)

	fmt.Printf("ZADD: %s\n", res)

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
		res, err := e.redisConn.Do("ZREM", ssKey, "NX", 0, utils.UintToPaddedString(order.Price)) // Add price point to order book
		if err != nil {
			log.Printf("ZREM: %s", err)
			return err
		}
		fmt.Printf("ZREM: %s\n", res)
		res, err = e.redisConn.Do("DEL", ssKey+"::book::"+utils.UintToPaddedString(order.Price)) // Add price point to order book
		if err != nil {
			log.Printf("DEL: %s", err)
			return err
		}
		fmt.Printf("DEL: %s\n", res)

		res, err = e.redisConn.Do("DEL", listKey+"::"+order.Hash)
		if err != nil {
			log.Printf("DEL: %s", err)
			return err
		}
		// Add order reference to price sorted set
		res, err = e.redisConn.Do("ZREM", listKey, order.Hash)
		if err != nil {
			log.Printf("ZREM: %s", err)
			return err
		}

		fmt.Printf("ZREM: %s\n", res)
	} else {
		res, err := e.redisConn.Do("ZADD", ssKey, "NX", 0, utils.UintToPaddedString(order.Price)) // Add price point to order book
		if err != nil {
			log.Printf("ZADD: %s", err)
			return err
		}
		fmt.Printf("ZADD: %s\n", res)
		res, err = e.redisConn.Do("INCRBY", ssKey+"::book::"+utils.UintToPaddedString(order.Price), -1*tradeAmount) // Add price point to order book
		if err != nil {
			log.Printf("INCRBY: %s", err)
			return err
		}
		fmt.Printf("INCRBY: %s\n", res)

		res, err = e.redisConn.Do("DEL", listKey+"::"+order.Hash)
		if err != nil {
			log.Printf("DEL: %s", err)
			return err
		}
		// Add order reference to price sorted set
		res, err = e.redisConn.Do("ZREM", listKey, order.Hash)
		if err != nil {
			log.Printf("ZREM: %s", err)
			return err
		}

		fmt.Printf("ZREM: %s\n", res)
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
		o.Order.Status = types.PARTIALFILLED
		o.Order.FilledAmount = o.Order.FilledAmount - o.Amount
		if o.Order.FilledAmount == 0 {
			o.Order.Status = types.OPEN
		}

		_, listKey := o.Order.GetOBKeys()
		res, _ := redis.Bytes(e.redisConn.Do("GET", listKey+"::"+o.Order.Hash))
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

// CancelOrder is used to cancel the order from orderbook
func (e *Resource) CancelOrder(order *types.Order) (engineResponse *Response, err error) {
	e.mutex.Lock()
	defer e.mutex.Unlock()
	_, listKey := order.GetOBKeys()
	res, err := redis.Bytes(e.redisConn.Do("GET", listKey+"::"+order.Hash))
	if err != nil {
		log.Printf("GET: %s", err)
		return
	}
	if res == nil {
		return
	}
	var storedOrder *types.Order

	if err := json.Unmarshal(res, &storedOrder); err != nil {
		log.Printf("GET: %s", err)
		return nil, err
	}
	if err := e.deleteOrder(order, storedOrder.Amount-storedOrder.FilledAmount); err != nil {
		log.Printf("\n%s\n", err)
		return nil, err
	}
	storedOrder.Status = types.CANCELLED
	engineResponse = &Response{
		Order:          storedOrder,
		Trades:         make([]*types.Trade, 0),
		RemainingOrder: &types.Order{},
		FillStatus:     CANCELLED,
		MatchingOrders: make([]*FillOrder, 0),
	}
	return
}
