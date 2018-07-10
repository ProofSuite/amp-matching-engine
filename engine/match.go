package engine

import (
	"github.com/Proofsuite/amp-matching-engine/types"
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

func (e *EngineResource) execute(order *types.Order, bookEntry *types.Order) (trade *types.Trade, fillOrder *FillOrder, err error) {
	fillOrder = &FillOrder{
		Order: bookEntry,
	}
	beAmtAvailable := bookEntry.Amount - bookEntry.FilledAmount
	orderUnfilledAmt := order.Amount - order.FilledAmount
	if beAmtAvailable > orderUnfilledAmt {
		fillOrder.Amount = orderUnfilledAmt

		bookEntry.FilledAmount = bookEntry.FilledAmount + orderUnfilledAmt
		bookEntry.Status = types.PARTIAL_FILLED

		e.updateOrder(bookEntry, fillOrder.Amount)

	} else {
		fillOrder.Amount = beAmtAvailable

		bookEntry.FilledAmount = bookEntry.FilledAmount + beAmtAvailable
		bookEntry.Status = types.FILLED

		e.deleteOrder(bookEntry, fillOrder.Amount)

	}

	// Create trade object to be passed to the system for further processing
	trade = &types.Trade{
		Amount:       fillOrder.Amount,
		Price:        order.Price,
		OrderHash:    bookEntry.Hash,
		Type:         order.Type,
		Taker:        order.UserAddress,
		PairName:     order.PairName,
		Maker:        bookEntry.UserAddress,
		TakerOrderID: order.ID,
		MakerOrderID: bookEntry.ID,
	}
	trade.Hash = trade.ComputeHash()
	return
}
