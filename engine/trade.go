package engine

import (
	"github.com/Proofsuite/amp-matching-engine/types"
)

// FillStatus is enum used to signify the filled status of order in engineResponse
type FillStatus int

// Response is the structure of message response sent by engine
type Response struct {
	Order          *types.Order
	Trades         []*types.Trade
	RemainingOrder *types.Order

	FillStatus     FillStatus
	MatchingOrders []*FillOrder
}

// this const block holds the possible valued of FillStatus
const (
	_ FillStatus = iota
	NOMATCH
	PARTIAL
	FULL
	ERROR
	CANCELLED
)

// execute function is responsible for executing of matched orders
// i.e it deletes/updates orders in case of order matching and responds 
// with trade instance and fillOrder
func (e *Resource) execute(order *types.Order, bookEntry *types.Order) (trade *types.Trade, fillOrder *FillOrder, err error) {
	fillOrder = &FillOrder{}
	beAmtAvailable := bookEntry.Amount - bookEntry.FilledAmount
	orderUnfilledAmt := order.Amount - order.FilledAmount
	if beAmtAvailable > orderUnfilledAmt {
		fillOrder.Amount = orderUnfilledAmt

		bookEntry.FilledAmount = bookEntry.FilledAmount + orderUnfilledAmt
		bookEntry.Status = types.PARTIALFILLED
		fillOrder.Order = bookEntry

		e.updateOrder(bookEntry, fillOrder.Amount)

	} else {
		fillOrder.Amount = beAmtAvailable

		bookEntry.FilledAmount = bookEntry.FilledAmount + beAmtAvailable
		bookEntry.Status = types.FILLED
		fillOrder.Order = bookEntry

		e.deleteOrder(bookEntry, fillOrder.Amount)

	}
	order.FilledAmount = order.FilledAmount + fillOrder.Amount
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
