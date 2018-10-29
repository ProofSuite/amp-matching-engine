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
	"math/big"
	"sync"

	"github.com/Proofsuite/amp-matching-engine/interfaces"
	"github.com/Proofsuite/amp-matching-engine/rabbitmq"
	"github.com/Proofsuite/amp-matching-engine/types"
	"github.com/Proofsuite/amp-matching-engine/utils"
	"github.com/Proofsuite/amp-matching-engine/utils/math"
	"github.com/ethereum/go-ethereum/common"
)

type OrderBook struct {
	rabbitMQConn *rabbitmq.Connection
	orderDao     interfaces.OrderDao
	tradeDao     interfaces.TradeDao
	pair         *types.Pair
	mutex        *sync.Mutex
}

// newOrder calls buyOrder/sellOrder based on type of order recieved and
// publishes the response back to rabbitmq
func (ob *OrderBook) newOrder(o *types.Order) (err error) {
	// Attain lock on engineResource, so that recovery or cancel order function doesn't interfere
	ob.mutex.Lock()
	defer ob.mutex.Unlock()

	utils.PrintJSON("In new order")

	res := &types.EngineResponse{}
	if o.Side == "SELL" {
		res, err = ob.sellOrder(o)
		if err != nil {
			logger.Error(err)
			return err
		}

	} else if o.Side == "BUY" {
		res, err = ob.buyOrder(o)
		if err != nil {
			logger.Error(err)
			return err
		}
	}

	// Note: Plug the option for orders like FOC, Limit here (if needed)
	err = ob.rabbitMQConn.PublishEngineResponse(res)
	if err != nil {
		logger.Error(err)
		return err
	}

	return nil
}

// addOrder adds an order to redis
func (ob *OrderBook) addOrder(o *types.Order) error {
	if o.FilledAmount == nil || math.IsZero(o.FilledAmount) {
		o.Status = "OPEN"
	}

	_, err := ob.orderDao.FindAndModify(o.Hash, o)
	if err != nil {
		// we add this condition in the case an order is re-run through the orderbook (in case of invalid counterpart order for example)
		logger.Error(err)
		return err
	}

	return nil
}

// buyOrder is triggered when a buy order comes in, it fetches the ask list
// from orderbook. First it checks ths price point list to check whether the order can be matched
// or not, if there are pricepoints that can satisfy the order then corresponding list of orders
// are fetched and trade is executed
func (ob *OrderBook) buyOrder(o *types.Order) (*types.EngineResponse, error) {
	res := &types.EngineResponse{}
	remainingOrder := *o

	matchingOrders, err := ob.orderDao.GetMatchingSellOrders(o)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	// case where no order is matched
	if len(matchingOrders) == 0 {
		ob.addOrder(o)
		res.Status = "ORDER_ADDED"
		res.Order = o
		return res, nil
	}

	matches := types.Matches{}
	for _, mo := range matchingOrders {
		trade, err := ob.execute(o, mo)
		if err != nil {
			logger.Error(err)
			return nil, err
		}

		matches = append(matches, &types.Match{mo, trade})
		remainingOrder.Amount = math.Sub(remainingOrder.Amount, trade.Amount)

		if math.IsZero(remainingOrder.Amount) {
			o.FilledAmount = o.Amount
			o.Status = "FILLED"

			_, err := ob.orderDao.FindAndModify(o.Hash, o)
			if err != nil {
				logger.Error(err)
				return nil, err
			}

			res.Status = "ORDER_FILLED"
			res.Order = o
			res.RemainingOrder = nil
			res.Matches = &matches
			return res, nil
		}
	}

	//TODO refactor
	o.Status = "REPLACED"
	_, err = ob.orderDao.FindAndModify(o.Hash, o)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	remainingOrder.Signature = nil
	remainingOrder.Nonce = nil
	remainingOrder.Hash = common.HexToHash("")
	remainingOrder.BuyAmount = remainingOrder.Amount
	remainingOrder.SellAmount = math.Div(
		math.Mul(remainingOrder.Amount, o.SellAmount),
		o.BuyAmount,
	)

	res.Status = "ORDER_PARTIALLY_FILLED"
	res.Order = o
	res.RemainingOrder = &remainingOrder
	res.Matches = &matches
	return res, nil
}

// sellOrder is triggered when a sell order comes in, it fetches the bid list
// from orderbook. First it checks ths price point list to check whether the order can be matched
// or not, if there are pricepoints that can satisfy the order then corresponding list of orders
// are fetched and trade is executed
func (ob *OrderBook) sellOrder(o *types.Order) (*types.EngineResponse, error) {
	res := &types.EngineResponse{}
	remainingOrder := *o

	matchingOrders, err := ob.orderDao.GetMatchingBuyOrders(o)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	if len(matchingOrders) == 0 {
		o.Status = "OPEN"
		ob.addOrder(o)

		res.Status = "ORDER_ADDED"
		res.Order = o
		return res, nil
	}

	matches := types.Matches{}
	for _, mo := range matchingOrders {
		trade, err := ob.execute(o, mo)
		if err != nil {
			logger.Error(err)
			return nil, err
		}

		matches = append(matches, &types.Match{mo, trade})
		remainingOrder.Amount = math.Sub(remainingOrder.Amount, trade.Amount)

		if math.IsZero(remainingOrder.Amount) {
			o.FilledAmount = o.Amount
			o.Status = "FILLED"

			_, err := ob.orderDao.FindAndModify(o.Hash, o)
			if err != nil {
				logger.Error(err)
				return nil, err
			}

			res.Status = "ORDER_FILLED"
			res.Order = o
			res.RemainingOrder = nil
			res.Matches = &matches
			return res, nil
		}
	}

	//TODO refactor
	o.Status = "REPLACED"
	_, err = ob.orderDao.FindAndModify(o.Hash, o)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	remainingOrder.Signature = nil
	remainingOrder.Nonce = nil
	remainingOrder.Hash = common.HexToHash("")
	remainingOrder.BuyAmount = remainingOrder.Amount
	remainingOrder.SellAmount = math.Div(
		math.Mul(res.RemainingOrder.Amount, o.SellAmount),
		o.BuyAmount,
	)

	res.Status = "ORDER_PARTIALLY_FILLED"
	res.Order = o
	res.RemainingOrder = &remainingOrder
	res.Matches = &matches
	return res, nil
}

// execute function is responsible for executing of matched orders
// i.e it deletes/updates orders in case of order matching and responds
// with trade instance and fillOrder
func (ob *OrderBook) execute(takerOrder *types.Order, makerOrder *types.Order) (*types.Trade, error) {
	trade := &types.Trade{}
	tradeAmount := big.NewInt(0)

	//TODO changes 'strictly greater than' condition. The orders that are almost completely filled
	//TODO should be removed/skipped
	if math.IsStrictlyGreaterThan(makerOrder.RemainingAmount(), takerOrder.RemainingAmount()) {
		tradeAmount = takerOrder.RemainingAmount()
		makerOrder.FilledAmount = math.Add(makerOrder.FilledAmount, tradeAmount)
		makerOrder.Status = "PARTIAL_FILLED"

		_, err := ob.orderDao.FindAndModify(makerOrder.Hash, makerOrder)
		if err != nil {
			logger.Error(err)
			return nil, err
		}
	} else {
		tradeAmount = makerOrder.RemainingAmount()
		makerOrder.FilledAmount = makerOrder.Amount
		makerOrder.Status = "FILLED"

		_, err := ob.orderDao.FindAndModify(makerOrder.Hash, makerOrder)
		if err != nil {
			logger.Error(err)
			return nil, err
		}
	}

	takerOrder.FilledAmount = math.Add(takerOrder.FilledAmount, tradeAmount)
	trade = &types.Trade{
		Amount:         tradeAmount,
		PricePoint:     takerOrder.PricePoint,
		BaseToken:      takerOrder.BaseToken,
		QuoteToken:     takerOrder.QuoteToken,
		OrderHash:      makerOrder.Hash,
		TakerOrderHash: takerOrder.Hash,
		Side:           takerOrder.Side,
		Taker:          takerOrder.UserAddress,
		PairName:       takerOrder.PairName,
		Maker:          makerOrder.UserAddress,
		Status:         "PENDING",
	}

	return trade, nil
}

// CancelOrder is used to cancel the order from orderbook
func (ob *OrderBook) cancelOrder(o *types.Order) error {
	ob.mutex.Lock()
	defer ob.mutex.Unlock()

	o.Status = "CANCELLED"
	err := ob.orderDao.UpdateOrderStatus(o.Hash, "CANCELLED")
	if err != nil {
		logger.Error(err)
		return err
	}

	res := &types.EngineResponse{
		Status:         "ORDER_CANCELLED",
		Order:          o,
		RemainingOrder: nil,
		Matches:        nil,
	}

	err = ob.rabbitMQConn.PublishEngineResponse(res)
	if err != nil {
		logger.Error(err)
		return err
	}

	return nil
}

// cancelTrades revertTrades and reintroduces the taker orders in the orderbook
func (ob *OrderBook) invalidateMakerOrders(matches types.Matches) error {
	ob.mutex.Lock()
	defer ob.mutex.Unlock()

	logger.Info("In invalidate maker orders")

	orders := matches.Orders()
	trades := matches.Trades()
	tradeAmounts := matches.TradeAmounts()
	makerOrderHashes := []common.Hash{}
	takerOrderHashes := []common.Hash{}

	for i, _ := range orders {
		makerOrderHashes = append(makerOrderHashes, trades[i].OrderHash)
		takerOrderHashes = append(takerOrderHashes, trades[i].TakerOrderHash)
	}

	takerOrders, err := ob.orderDao.UpdateOrderFilledAmounts(takerOrderHashes, tradeAmounts)
	if err != nil {
		logger.Error(err)
		return err
	}

	makerOrders, err := ob.orderDao.UpdateOrderStatusesByHashes("INVALIDATED", makerOrderHashes...)
	if err != nil {
		logger.Error(err)
		return err
	}

	//TODO in the case the trades are not in the database they should not be created.
	cancelledTrades, err := ob.tradeDao.UpdateTradeStatusesByOrderHashes("CANCELLED", takerOrderHashes...)
	if err != nil {
		logger.Error(err)
		return err
	}

	res := &types.EngineResponse{
		Status:            "TRADES_CANCELLED",
		InvalidatedOrders: &makerOrders,
		CancelledTrades:   &cancelledTrades,
	}

	err = ob.rabbitMQConn.PublishEngineResponse(res)
	if err != nil {
		logger.Error(err)
	}

	for _, o := range takerOrders {
		err := ob.rabbitMQConn.PublishNewOrderMessage(o)
		if err != nil {
			logger.Error(err)
		}
	}

	return nil
}

func (ob *OrderBook) invalidateTakerOrders(matches types.Matches) error {
	ob.mutex.Lock()
	defer ob.mutex.Unlock()

	orders := matches.Orders()
	trades := matches.Trades()
	tradeAmounts := matches.TradeAmounts()
	makerOrderHashes := []common.Hash{}
	takerOrderHashes := []common.Hash{}

	for i, _ := range orders {
		makerOrderHashes = append(makerOrderHashes, trades[i].OrderHash)
		takerOrderHashes = append(takerOrderHashes, trades[i].TakerOrderHash)
	}

	makerOrders, err := ob.orderDao.UpdateOrderFilledAmounts(makerOrderHashes, tradeAmounts)
	if err != nil {
		logger.Error(err)
		return err
	}

	takerOrders, err := ob.orderDao.UpdateOrderStatusesByHashes("INVALIDATED", takerOrderHashes...)
	if err != nil {
		logger.Error(err)
		return err
	}

	cancelledTrades, err := ob.tradeDao.UpdateTradeStatusesByOrderHashes("CANCELLED", makerOrderHashes...)
	if err != nil {
		logger.Error(err)
		return err
	}

	res := &types.EngineResponse{
		Status:            "TRADES_CANCELLED",
		InvalidatedOrders: &takerOrders,
		CancelledTrades:   &cancelledTrades,
	}

	err = ob.rabbitMQConn.PublishEngineResponse(res)
	if err != nil {
		logger.Error(err)
		return err
	}

	for _, o := range makerOrders {
		err := ob.rabbitMQConn.PublishNewOrderMessage(o)
		if err != nil {
			logger.Error(err)
		}
	}

	return nil
}

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

// // RecoverOrders is responsible for recovering the orders that failed to execute after matching
// // Orders are updated or added to orderbook based on whether that order exists in orderbook or not.
// func (ob *OrderBook) recoverOrders(matches types.Matches) error {
// 	ob.mutex.Lock()
// 	defer ob.mutex.Unlock()

// 	for _, m := range matches {
// 		t := m.Trade
// 		o := m.Order

// 		filledAmount := math.Sub(o.FilledAmount, t.Amount)
// 		if math.IsEqualOrSmallerThan(filledAmount, big.NewInt(0)) {
// 			o.Status = "OPEN"
// 			o.FilledAmount = big.NewInt(0)
// 			_, err := ob.orderDao.FindAndModify(o.Hash, o)
// 			if err != nil {
// 				logger.Error(err)
// 				return err
// 			}

// 		} else {
// 			o.Status = "PARTIAL_FILLED"
// 			o.FilledAmount = math.Sub(o.FilledAmount, t.Amount)
// 			_, err := ob.orderDao.FindAndModify(o.Hash, o)
// 			if err != nil {
// 				logger.Error(err)
// 				return err
// 			}
// 		}
// 	}

// 	return nil
// }

// func (ob *OrderBook) deleteOrders(orders ...*types.Order) error {
// 	ob.mutex.Lock()
// 	defer ob.mutex.Unlock()

// 	err := ob.orderDao.Delete(orders...)
// 	if err != nil {
// 		logger.Error(err)
// 		return err
// 	}

// 	res := &types.EngineResponse{
// 		Status:        "DELETED",
// 		DeletedOrders: &orders,
// 	}

// 	// Note: Plug the option for orders like FOC, Limit here (if needed)
// 	err = ob.rabbitMQConn.PublishEngineResponse(res)
// 	if err != nil {
// 		logger.Error(err)
// 		return err
// 	}

// 	return nil
// }

func (ob *OrderBook) InvalidateOrder(o *types.Order) (*types.EngineResponse, error) {
	ob.mutex.Lock()
	defer ob.mutex.Unlock()

	o.Status = "ERROR"
	err := ob.orderDao.UpdateOrderStatus(o.Hash, "ERROR")
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	res := &types.EngineResponse{
		Status:         "INVALIDATED",
		Order:          o,
		RemainingOrder: nil,
		Matches:        nil,
	}

	return res, nil
}
