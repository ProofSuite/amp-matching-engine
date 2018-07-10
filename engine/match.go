package engine

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"github.com/Proofsuite/amp-matching-engine/utils"

	"github.com/Proofsuite/amp-matching-engine/types"
	"github.com/gomodule/redigo/redis"
)

type Match struct {
	Order          *types.Order
	FillStatus     FillStatus
	MatchingOrders []*FillOrder
}
type FillStatus int

type EngineResponse struct {
	Order          *types.Order
	Trades         []*types.Trade
	RemainingOrder *types.Order

	FillStatus     FillStatus
	MatchingOrders []*FillOrder
}

const (
	_ FillStatus = iota
	NO_MATCH
	PARTIAL
	FULL
	ERROR
)

func (e *EngineResource) execute(m *Match, er *EngineResponse) (err error) {

	if m == nil {
		err = errors.New("No match passed")
		return
	}

	var filledAmount int64

	order := er.Order
	matchedOrders := m.MatchingOrders

	for i, o := range matchedOrders {
		mo := o.Order
		ss, list := mo.GetOBKeys()
		// POP the order from the top of list
		reply, err := redis.Bytes(e.redisConn.Do("LPOP", list)) // "ZREVRANGEBYLEX" key max min
		if err != nil {
			log.Printf("LPOP: %s\n", err)
			return err
		}

		var bookEntry types.Order
		err = json.Unmarshal(reply, &bookEntry)
		if err != nil {
			log.Printf("json.Unmarshal: %s\n", err)
			return err
		}

		if bookEntry.ID != mo.ID {
			log.Fatal("Invalid matching order passed: ", bookEntry.ID, mo.ID, list)
			return errors.New("Invalid matching order passed")
		}
		bookEntry.FilledAmount = bookEntry.FilledAmount + o.Amount

		filledAmount = filledAmount + o.Amount

		// Create trade object to be passed to the system for further processing
		t := &types.Trade{
			Amount:       o.Amount,
			Price:        order.Price,
			OrderHash:    mo.Hash,
			Type:         order.Type,
			TradeNonce:   int64(i),
			Taker:        order.UserAddress,
			PairName:     order.PairName,
			Maker:        mo.UserAddress,
			TakerOrderID: order.ID,
			MakerOrderID: mo.ID,
		}
		t.Hash = t.ComputeHash()

		er.Trades = append(er.Trades, t)

		// If book entry order is not filled completely then update the filledAmount and push it back to the head of list

		if (bookEntry.Amount - bookEntry.FilledAmount) > 0 {
			bookEntryAsBytes, err := json.Marshal(bookEntry)
			if err != nil {
				log.Printf("json.Marshal: %s", err)
			}
			res, err := e.redisConn.Do("LPUSH", list, bookEntryAsBytes)
			if err != nil {
				log.Printf("LPUSH: %s", err)
			}
			fmt.Println("LPUSH: ", res)
		}

		if bookEntry.FilledAmount == bookEntry.Amount {
			bookEntry.Status = types.FILLED
		} else {
			bookEntry.Status = types.PARTIAL_FILLED
		}

		er.MatchingOrders = append(er.MatchingOrders, &FillOrder{o.Amount, &bookEntry})

		// Update OrderBook
		res, err := e.redisConn.Do("INCRBY", ss+"::book::"+utils.UintToPaddedString(order.Price), -1*order.Amount) // Add price point to order book
		if err != nil {
			log.Printf("DEL BookEntry: %s", err)
		}
		fmt.Printf("DEL BookEntry: %s\n", res)
		// Get length of remaining orders in the list

		l, err := redis.Int64(e.redisConn.Do("LLEN", list))
		if err != nil {
			log.Printf("LLEN: %s", err)
		} else if l == 0 {
			// If list is empty: remove the list and remove the price point from sorted set
			_, err := e.redisConn.Do("del", list)
			if err != nil {
				log.Printf("del: %s", err)
				return err
			}
			// fmt.Printf("del: %s", res)
			res, err := e.redisConn.Do("DEL", ss+"::book::"+utils.UintToPaddedString(order.Price)) // Add price point to order book
			if err != nil {
				log.Printf("DEL BookEntry: %s", err)
			}
			fmt.Printf("DEL BookEntry: %s\n", res)
			_, err = e.redisConn.Do("ZREM", ss, utils.UintToPaddedString(mo.Price))
			if err != nil {
				log.Printf("ZREM: %s", err)
				return err
			}
			// fmt.Printf("ZREM: %s", res)
		}
	}

	order.FilledAmount = filledAmount

	if order.Amount != order.FilledAmount {
		er.Order.Status = types.PARTIAL_FILLED
		er.FillStatus = PARTIAL
		remOrder := *order
		remOrder.Amount = order.Amount - order.FilledAmount
		remOrder.FilledAmount = 0
		er.RemainingOrder = &remOrder
	} else {
		er.Order.Status = types.FILLED
		er.FillStatus = FULL
		er.RemainingOrder = nil
	}
	return
}
