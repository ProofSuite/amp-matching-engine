package engine

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/Proofsuite/amp-matching-engine/ws"

	"github.com/Proofsuite/amp-matching-engine/types"
	"github.com/Proofsuite/amp-matching-engine/utils"
	"github.com/gomodule/redigo/redis"
)

type FillOrder struct {
	Amount int64
	Order  *types.Order
}

func (e *EngineResource) matchOrder(order *types.Order) (err error) {
	engineResponse := &EngineResponse{
		Order: order,
	}
	engineResponse.Trades = make([]*types.Trade, 0)
	engineResponse.RemainingOrder = order
	engineResponse.MatchingOrders = make([]*FillOrder, 0)

	match := &Match{
		Order: order,
	}

	// for match.FillStatus != NO_MATCH {
	if order.Type == types.SELL {
		err = e.sellOrder(order, match)
	} else if order.Type == types.BUY {
		err = e.buyOrder(order, match)
	}
	// Note: Plug the option for orders like FOC, Limit, OnlyFill (If Required)

	// If NO_MATCH add to order book
	if match.FillStatus == NO_MATCH {
		engineResponse.Order.Status = types.OPEN
		e.addOrder(order)
		msg := &types.WsMsg{MsgType: "added_to_orderbook"}
		msg.OrderID = order.ID
		msg.Data = engineResponse
		erab, err := json.Marshal(msg)
		if err != nil {
			log.Fatalf("%s", err)
		}
		ws.Connections[order.ID.Hex()].Conn.WriteMessage(1, erab)

	} else {

		// Execute Trade
		err = e.execute(match, engineResponse)
		if err != nil {
			log.Printf("\nexecute XXXXXXX\n%s\nXXXXXXX execute\n", err)
		}
		msg := &types.WsMsg{MsgType: "trade_remorder_sign"}
		msg.OrderID = order.ID
		msg.Data = engineResponse
		erab, err := json.Marshal(msg)
		if err != nil {
			log.Fatalf("%s", err)
		}

		ws.Connections[order.ID.Hex()].Conn.WriteMessage(1, erab)

		// for {
		t := time.NewTimer(5 * time.Second)
		select {
		case rm := <-ws.Connections[order.ID.Hex()].ReadChannel:
			if rm.MsgType == "trade_remorder_sign" {
				mb, err := json.Marshal(rm.Data)
				if err != nil {
					ws.Connections[order.ID.Hex()].Conn.WriteMessage(1, []byte(err.Error()))
					ws.Connections[order.ID.Hex()].Conn.Close()
				}
				var ersb *EngineResponse
				err = json.Unmarshal(mb, &ersb)
				if err != nil {
					ws.Connections[order.ID.Hex()].Conn.WriteMessage(1, []byte(err.Error()))
					ws.Connections[order.ID.Hex()].Conn.Close()
				}
				if engineResponse.FillStatus == PARTIAL {
					order.OrderBook = &types.OrderSubDoc{Amount: ersb.RemainingOrder.Amount, Signature: ersb.RemainingOrder.Signature}
					e.addOrder(order)
				}
			}
			t.Stop()
			break
		case <-t.C:
			fmt.Printf("\nTimeout\n")
			e.recoverOrders(engineResponse.MatchingOrders)
			engineResponse.FillStatus = ERROR
			engineResponse.Order.Status = types.ERROR
			engineResponse.Trades = nil
			engineResponse.RemainingOrder = nil
			engineResponse.MatchingOrders = nil
			t.Stop()
			break
		}
	}
	e.publishEngineResponse(engineResponse)
	if err != nil {
		log.Printf("\npublishEngineResponse XXXXXXX\n%s\nXXXXXXX publishEngineResponse\n", err)
	}
	return
}

func (e *EngineResource) buyOrder(order *types.Order, match *Match) (err error) {

	if match.MatchingOrders == nil {
		match.MatchingOrders = make([]*FillOrder, 0)
	}

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

	var filledAmount int64
	var orderAmount = order.Amount

	if len(priceRange) == 0 {
		match.FillStatus = NO_MATCH
	} else {
		for _, pr := range priceRange {
			reply, err := redis.ByteSlices(e.redisConn.Do("LRANGE", oskv+"::"+utils.UintToPaddedString(pr), 0, -1)) // "ZREVRANGEBYLEX" key max min
			if err != nil {
				log.Printf("LRANGE: %s\n", err)
				return err
			}

			for _, o := range reply {
				var bookEntry types.Order
				err = json.Unmarshal(o, &bookEntry)
				if err != nil {
					log.Printf("json.Unmarshal: %s\n", err)
					return err
				}

				match.FillStatus = PARTIAL

				// update filledAmount
				beAmtAvailable := bookEntry.Amount - bookEntry.FilledAmount
				if beAmtAvailable > order.Amount {
					match.MatchingOrders = append(match.MatchingOrders, &FillOrder{order.Amount, &bookEntry})
					filledAmount = order.Amount
				} else {
					match.MatchingOrders = append(match.MatchingOrders, &FillOrder{beAmtAvailable, &bookEntry})
					filledAmount = beAmtAvailable
				}

				if filledAmount == orderAmount {
					match.FillStatus = FULL
					// order filled return
					return nil
				}
			}
		}
	}
	return
}

func (e *EngineResource) sellOrder(order *types.Order, match *Match) (err error) {
	if match.MatchingOrders == nil {
		match.MatchingOrders = make([]*FillOrder, 0)
	}

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

	var filledAmount int64
	var orderAmount = order.Amount

	if len(priceRange) == 0 {
		match.FillStatus = NO_MATCH
	} else {
		for _, pr := range priceRange {
			reply, err := redis.ByteSlices(e.redisConn.Do("LRANGE", obkv+"::"+utils.UintToPaddedString(pr), 0, -1)) // "ZREVRANGEBYLEX" key max min
			if err != nil {
				log.Printf("LRANGE: %s\n", err)
				return err
			}

			for _, o := range reply {
				var bookEntry types.Order
				err = json.Unmarshal(o, &bookEntry)
				if err != nil {
					log.Printf("json.Unmarshal: %s\n", err)
					return err
				}

				match.FillStatus = PARTIAL

				// update filledAmount
				beAmtAvailable := bookEntry.Amount - bookEntry.FilledAmount
				if beAmtAvailable > order.Amount {
					match.MatchingOrders = append(match.MatchingOrders, &FillOrder{order.Amount, &bookEntry})
					filledAmount = order.Amount
				} else {
					match.MatchingOrders = append(match.MatchingOrders, &FillOrder{beAmtAvailable, &bookEntry})
					filledAmount = beAmtAvailable
				}

				if filledAmount == orderAmount {
					match.FillStatus = FULL
					// order filled return
					return nil
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
	res, err = e.redisConn.Do("INCRBY", ssKey+"::book::"+utils.UintToPaddedString(order.Price), order.Amount) // Add price point to order book
	if err != nil {
		log.Printf("INCRBY: %s", err)
	}
	fmt.Printf("INCRBY: %s\n", res)

	// Add order to list
	orderAsBytes, err := json.Marshal(order)
	if err != nil {
		log.Printf("ZADD: %s", err)
	}
	res, err = e.redisConn.Do("RPUSH", listKey, orderAsBytes)
	if err != nil {
		log.Printf("RPUSH: %s", err)
	}
	fmt.Printf("RPUSH: %s\n", res)

	return
}

func (e *EngineResource) recoverOrders(orders []*FillOrder) {
	for _, o := range orders {
		e.addOrder(o.Order)
	}
}
