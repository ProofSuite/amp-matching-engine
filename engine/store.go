package engine

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/Proofsuite/amp-matching-engine/types"
	"github.com/Proofsuite/amp-matching-engine/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/gomodule/redigo/redis"
)

// GetPricePoints returns the pricepoints matching a certain (pair, pricepoint)
func (ob *OrderBook) GetMatchingBuyPricePoints(obKey string, pricePoint int64) ([]int64, error) {
	pps, err := ob.redisConn.ZRangeByLexInt(obKey, "-", "["+utils.UintToPaddedString(pricePoint))
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	return pps, nil
}

func (ob *OrderBook) GetMatchingSellPricePoints(obkv string, pricePoint int64) ([]int64, error) {
	pps, err := ob.redisConn.ZRevRangeByLexInt(obkv, "+", "["+utils.UintToPaddedString(pricePoint))
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	return pps, nil
}

func (ob *OrderBook) GetPricePointVolume(pricePointSetKey string, pricePoint int64) (int64, error) {
	val, err := ob.redisConn.GetValue(pricePointSetKey + "::book::" + utils.UintToPaddedString(pricePoint))
	if err != nil {
		logger.Error(err)
		return 0, err
	}

	vol, err := strconv.ParseInt(val, 10, 64)
	if err != nil {
		logger.Error(err)
		return 0, err
	}

	return vol, nil
}

func (ob *OrderBook) GetFromOrderMap(hash common.Hash) (*types.Order, error) {
	o := &types.Order{}
	keys, err := ob.redisConn.Keys("*::" + hash.Hex())
	if err != nil {
		return nil, err
	}

	if len(keys) == 0 {
		return nil, fmt.Errorf("Key doesn't exists")
	}

	serialized, err := ob.redisConn.GetValue(keys[0])
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	err = json.Unmarshal([]byte(serialized), &o)
	if err != nil {
		return nil, err
	}

	return o, nil
}

// GetPricePointOrders returns the orders hashes for a (pair, pricepoint)
func (ob *OrderBook) GetMatchingOrders(obKey string, pricePoint int64) ([][]byte, error) {
	k := obKey + "::" + utils.UintToPaddedString(pricePoint)
	orders, err := ob.redisConn.Sort(k, "", true, false, k+"::orders::*")
	if err != nil {
		return nil, err
	}

	return orders, nil
}

// AddPricePointToSet
func (ob *OrderBook) AddToPricePointSet(pricePointSetKey string, pricePoint int64) error {
	err := ob.redisConn.ZAdd(pricePointSetKey, 0, utils.UintToPaddedString(pricePoint))
	if err != nil {
		logger.Error(err)
		return err
	}

	return nil
}

// RemoveFromPricePointSet
func (ob *OrderBook) RemoveFromPricePointSet(pricePointSetKey string, pricePoint int64) error {
	err := ob.redisConn.ZRem(pricePointSetKey, utils.UintToPaddedString(pricePoint))
	if err != nil {
		logger.Error(err)
		return err
	}

	return nil
}

// AddPricePointHashesSet
func (ob *OrderBook) AddToPricePointHashesSet(orderHashListKey string, createdAt time.Time, hash common.Hash) error {
	err := ob.redisConn.ZAdd(orderHashListKey, createdAt.Unix(), hash.Hex())
	if err != nil {
		return err
	}

	return nil
}

// RemoveFromPricePointHashesSet
func (ob *OrderBook) RemoveFromPricePointHashesSet(orderHashListKey string, hash common.Hash) error {
	err := ob.redisConn.ZRem(orderHashListKey, hash.Hex())
	if err != nil {
		return err
	}

	return nil
}

// IncrementPricePointVolume increases the value of a certain pricepoint at a certain volume
func (ob *OrderBook) IncrementPricePointVolume(pricePointSetKey string, pricePoint int64, amount int64) error {
	_, err := ob.redisConn.IncrBy(pricePointSetKey+"::book::"+utils.UintToPaddedString(pricePoint), amount)
	if err != nil {
		logger.Error(err)
		return err
	}

	return nil
}

// DecrementPricePoint
func (ob *OrderBook) DecrementPricePointVolume(pricePointSetKey string, pricePoint int64, amount int64) error {
	_, err := ob.redisConn.IncrBy(pricePointSetKey+"::book::"+utils.UintToPaddedString(pricePoint), -amount)
	if err != nil {
		logger.Error(err)
		return err
	}

	return nil
}

// DeletePricePoint
func (ob *OrderBook) DeletePricePointVolume(pricePointSetKey string, pricePoint int64) error {
	err := ob.redisConn.Del(pricePointSetKey + "::book::" + utils.UintToPaddedString(pricePoint))
	if err != nil {
		logger.Error(err)
		return err
	}

	return nil
}

// AddToOrderMap
func (ob *OrderBook) AddToOrderMap(o *types.Order) error {
	bytes, err := json.Marshal(o)
	if err != nil {
		logger.Error(err)
		return err
	}

	decoded := &types.Order{}
	json.Unmarshal(bytes, decoded)
	_, orderHashListKey := o.GetOBKeys()
	err = ob.redisConn.Set(orderHashListKey+"::orders::"+o.Hash.Hex(), string(bytes))
	if err != nil {
		logger.Error(err)
		return err
	}

	return nil
}

// RemoveFromOrderMap
func (ob *OrderBook) RemoveFromOrderMap(hash common.Hash) error {
	keys, _ := ob.redisConn.Keys("*::" + hash.Hex())
	err := ob.redisConn.Del(keys[0])
	if err != nil {
		logger.Error(err)
		return err
	}

	return nil
}

// GetOrderBook fetches the complete orderbook from redis for the required pair
func (ob *OrderBook) GetOrderBook(pair *types.Pair) (sellBook, buyBook []*map[string]float64) {
	sKey, bKey := pair.GetOrderBookKeys()
	res, err := redis.Int64s(ob.redisConn.Do("SORT", sKey, "GET", sKey+"::book::*", "GET", "#")) // Add price point to order book
	if err != nil {
		logger.Error(err)
	}

	for i := 0; i < len(res); i = i + 2 {
		temp := &map[string]float64{
			"amount": float64(res[i]),
			"price":  float64(res[i+1]),
		}
		sellBook = append(sellBook, temp)
	}

	res, err = redis.Int64s(ob.redisConn.Do("SORT", bKey, "GET", bKey+"::book::*", "GET", "#", "DESC"))
	if err != nil {
		logger.Error(err)
	}

	for i := 0; i < len(res); i = i + 2 {
		temp := &map[string]float64{
			"amount": float64(res[i]),
			"price":  float64(res[i+1]),
		}
		buyBook = append(buyBook, temp)
	}

	return
}

// GetFullOrderBook fetches the complete orderbook from redis for the required pair
func (ob *OrderBook) GetFullOrderBook(pair *types.Pair) (book [][]types.Order) {
	pattern := pair.GetKVPrefix() + "::*::*::orders::*"

	groupInt := 100

	book = make([][]types.Order, 0)
	keys, err := ob.redisConn.Keys(pattern)
	if err != nil {
		logger.Error(err)
		return
	}

	orders := make([]types.Order, 0)

	for start := 0; start < len(keys); start = start + groupInt {
		end := start + groupInt
		if len(keys) < end {
			end = len(keys)
		}
		res, err := ob.redisConn.MGet(keys[start:end]...)
		if err != nil {
			logger.Error(err)
			return
		}
		for _, r := range res {
			if r == "" {
				continue
			}
			var temp types.Order
			if err := json.Unmarshal([]byte(r), &temp); err != nil {
				continue
			}
			orders = append(orders, temp)
		}
	}

	for start := 0; start < len(orders); start = start + groupInt {
		end := start + groupInt
		if len(keys) < end {
			end = len(keys)
		}
		book = append(book, orders[start:end])
	}

	return
}
