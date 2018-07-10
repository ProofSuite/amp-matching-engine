package engine

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/Proofsuite/amp-matching-engine/types"
	"github.com/Proofsuite/amp-matching-engine/utils"
	"github.com/gomodule/redigo/redis"
)

type FillOrder struct {
	Amount int64
	Order  *types.Order
}

func (e *EngineResource) matchOrder(order *types.Order) (err error) {

	var engineResponse *EngineResponse
	// for match.FillStatus != NO_MATCH {
	if order.Type == types.SELL {
		engineResponse, err = e.sellOrder(order)
	} else if order.Type == types.BUY {
		engineResponse, err = e.buyOrder(order)
	}
	// Note: Plug the option for orders like FOC, Limit, OnlyFill (If Required)

	e.publishEngineResponse(engineResponse)
	if err != nil {
		log.Printf("\npublishEngineResponse XXXXXXX\n%s\nXXXXXXX publishEngineResponse\n", err)
	}
	return
}

func (e *EngineResource) buyOrder(order *types.Order) (engineResponse *EngineResponse, err error) {

	engineResponse = &EngineResponse{
		Order:      order,
		FillStatus: PARTIAL,
	}
	engineResponse.Trades = make([]*types.Trade, 0)
	engineResponse.RemainingOrder = order
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
		engineResponse.FillStatus = NO_MATCH
		e.addOrder(order)
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
				engineResponse.RemainingOrder.FilledAmount = engineResponse.RemainingOrder.FilledAmount + fillOrder.Amount

				if engineResponse.RemainingOrder.FilledAmount == engineResponse.RemainingOrder.FilledAmount {
					engineResponse.FillStatus = FULL
					return engineResponse, nil
				}
			}
		}
	}
	return
}

func (e *EngineResource) sellOrder(order *types.Order) (engineResponse *EngineResponse, err error) {
	engineResponse = &EngineResponse{
		Order:      order,
		FillStatus: PARTIAL,
	}
	engineResponse.Trades = make([]*types.Trade, 0)
	engineResponse.RemainingOrder = order
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
		engineResponse.FillStatus = NO_MATCH
		e.addOrder(order)
	} else {
		for _, pr := range priceRange {
			bookEntries, err := redis.ByteSlices(e.redisConn.Do("SORT", obkv+"::"+utils.UintToPaddedString(pr), "GET", obkv+"::"+utils.UintToPaddedString(pr)+"::*", "ALPHA")) // "ZREVRANGEBYLEX" key max min
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
				engineResponse.RemainingOrder.FilledAmount = engineResponse.RemainingOrder.FilledAmount + fillOrder.Amount

				if engineResponse.RemainingOrder.FilledAmount == engineResponse.RemainingOrder.FilledAmount {
					engineResponse.FillStatus = FULL
					return engineResponse, nil
				}
			}
		}
	}
	return
}

func (e *EngineResource) addOrder(order *types.Order) {

	ssKey, listKey := order.GetOBKeys()
	res, err := e.redisConn.Do("ZADD", ssKey, "NX", 0, utils.UintToPaddedString(order.Price)) // Add price point to order book
	if err != nil {
		log.Printf("ZADD: %s", err)
	}
	fmt.Printf("ZADD: %s\n", res)
	res, err = e.redisConn.Do("INCRBY", ssKey+"::book::"+utils.UintToPaddedString(order.Price), order.Amount-order.FilledAmount) // Add price point to order book
	if err != nil {
		log.Printf("INCRBY: %s", err)
	}
	fmt.Printf("INCRBY: %s\n", res)

	// Add order to list
	orderAsBytes, err := json.Marshal(order)
	if err != nil {
		log.Printf("orderAsBytes: %s", err)
	}
	res, err = e.redisConn.Do("SET", listKey+"::"+order.ID.Hex(), string(orderAsBytes))
	if err != nil {
		log.Printf("SET: %s", err)
	}
	// Add order reference to price sorted set
	res, err = e.redisConn.Do("ZADD", listKey, "NX", order.CreatedAt.Unix(), order.ID.Hex())
	if err != nil {
		log.Printf("ZADD: %s", err)
	}

	fmt.Printf("ZADD: %s\n", res)

	return
}
func (e *EngineResource) updateOrder(order *types.Order, tradeAmount int64) {

	ssKey, listKey := order.GetOBKeys()
	res, err := e.redisConn.Do("ZADD", ssKey, "NX", 0, utils.UintToPaddedString(order.Price)) // Add price point to order book
	if err != nil {
		log.Printf("ZADD: %s", err)
	}
	fmt.Printf("ZADD: %s\n", res)
	res, err = e.redisConn.Do("INCRBY", ssKey+"::book::"+utils.UintToPaddedString(order.Price), -1*tradeAmount) // Add price point to order book
	if err != nil {
		log.Printf("INCRBY: %s", err)
	}
	fmt.Printf("INCRBY: %s\n", res)

	// Add order to list
	orderAsBytes, err := json.Marshal(order)
	if err != nil {
		log.Printf("orderAsBytes: %s", err)
	}
	res, err = e.redisConn.Do("SET", listKey+"::"+order.ID.Hex(), string(orderAsBytes))
	if err != nil {
		log.Printf("SET: %s", err)
	}
	// Add order reference to price sorted set
	res, err = e.redisConn.Do("ZADD", listKey, "NX", order.CreatedAt.Unix(), order.ID.Hex())
	if err != nil {
		log.Printf("ZADD: %s", err)
	}

	fmt.Printf("ZADD: %s\n", res)

	return
}

func (e *EngineResource) deleteOrder(order *types.Order, tradeAmount int64) {

	ssKey, listKey := order.GetOBKeys()
	remVolume, err := redis.Int64(e.redisConn.Do("GET", ssKey+"::book::"+utils.UintToPaddedString(order.Price)))
	if err != nil {
		log.Printf("GET remVolume: %s", err)
	}
	if remVolume == tradeAmount {
		res, err := e.redisConn.Do("ZREM", ssKey, "NX", 0, utils.UintToPaddedString(order.Price)) // Add price point to order book
		if err != nil {
			log.Printf("ZREM: %s", err)
		}
		fmt.Printf("ZREM: %s\n", res)
		res, err = e.redisConn.Do("DEL", ssKey+"::book::"+utils.UintToPaddedString(order.Price)) // Add price point to order book
		if err != nil {
			log.Printf("DEL: %s", err)
		}
		fmt.Printf("DEL: %s\n", res)

		res, err = e.redisConn.Do("DEL", listKey+"::"+order.ID.Hex())
		if err != nil {
			log.Printf("DEL: %s", err)
		}
		// Add order reference to price sorted set
		res, err = e.redisConn.Do("ZREM", listKey, order.ID.Hex())
		if err != nil {
			log.Printf("ZREM: %s", err)
		}

		fmt.Printf("ZREM: %s\n", res)
	} else {
		res, err := e.redisConn.Do("ZADD", ssKey, "NX", 0, utils.UintToPaddedString(order.Price)) // Add price point to order book
		if err != nil {
			log.Printf("ZADD: %s", err)
		}
		fmt.Printf("ZADD: %s\n", res)
		res, err = e.redisConn.Do("INCRBY", ssKey+"::book::"+utils.UintToPaddedString(order.Price), -1*tradeAmount) // Add price point to order book
		if err != nil {
			log.Printf("INCRBY: %s", err)
		}
		fmt.Printf("INCRBY: %s\n", res)

		res, err = e.redisConn.Do("DEL", listKey+"::"+order.ID.Hex())
		if err != nil {
			log.Printf("DEL: %s", err)
		}
		// Add order reference to price sorted set
		res, err = e.redisConn.Do("ZREM", listKey, order.ID.Hex())
		if err != nil {
			log.Printf("ZREM: %s", err)
		}

		fmt.Printf("ZREM: %s\n", res)
	}

	return
}
func (e *EngineResource) recoverOrders(orders []*FillOrder) {
	for _, o := range orders {
		e.addOrder(o.Order)
	}
}
