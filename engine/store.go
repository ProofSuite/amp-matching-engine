package engine

import (
	"encoding/json"
	"fmt"
	"math/big"
	"time"

	"github.com/Proofsuite/amp-matching-engine/types"
	"github.com/Proofsuite/amp-matching-engine/utils"
	"github.com/Proofsuite/amp-matching-engine/utils/math"
	"github.com/ethereum/go-ethereum/common"
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

func (ob *OrderBook) GetFromOrderMap(hash common.Hash) (*types.Order, error) {
	o := &types.Order{}
	keys, err := ob.redisConn.Keys("*::" + hash.Hex())
	if err != nil {
		return nil, err
	}

	if len(keys) == 0 {
		return nil, fmt.Errorf("Key doesn't exist")
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
// func (ob *OrderBook) AddToPricePointSet(pricePointSetKey string, pricePoint int64) error {
// 	err := ob.redisConn.ZAdd(pricePointSetKey, 0, utils.UintToPaddedString(pricePoint))
// 	if err != nil {
// 		logger.Error(err)
// 		return err
// 	}

// 	return nil
// }

func (ob *OrderBook) AddToPricePointSet(pricePointSetKey string, pricepoint *big.Int, amount *big.Int) error {
	volume := math.Div(amount, big.NewInt(1e18)).Int64()
	_, err := ob.redisConn.ZIncrBy(pricePointSetKey, volume, utils.UintToPaddedString(pricepoint.Int64()))
	if err != nil {
		logger.Error(err)
		return err
	}

	return nil
}

func (ob *OrderBook) RemoveFromPricePointSet(pricePointSetKey string, pricepoint *big.Int) error {
	err := ob.redisConn.ZRem(pricePointSetKey, utils.UintToPaddedString(pricepoint.Int64()))
	if err != nil {
		logger.Error(err)
		return err
	}

	return nil
}

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

func (ob *OrderBook) GetPricePointSetLength(pricePointSetKey string) (int64, error) {
	count, err := ob.redisConn.ZCount(pricePointSetKey)
	if err != nil {
		logger.Error(err)
		return 0, err
	}

	return count, nil
}

// AddPricePointHashesSet
func (ob *OrderBook) AddToPricePointHashesSet(orderHashListKey string, createdAt time.Time, hash common.Hash) error {
	err := ob.redisConn.ZAdd(orderHashListKey, createdAt.Unix(), hash.Hex())
	if err != nil {
		logger.Error(err)
		return err
	}

	return nil
}

// RemoveFromPricePointHashesSet
func (ob *OrderBook) RemoveFromPricePointHashesSet(orderHashListKey string, hash common.Hash) error {
	err := ob.redisConn.ZRem(orderHashListKey, hash.Hex())
	if err != nil {
		logger.Error(err)
		return err
	}

	return nil
}

func (ob *OrderBook) GetPricePointHashesSetLength(orderHashListKey string) (int64, error) {
	count, err := ob.redisConn.ZCount(orderHashListKey)
	if err != nil {
		logger.Error(err)
		return 0, err
	}

	return count, nil
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
