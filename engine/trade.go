package engine

import (
	"log"
	"math/big"

	"github.com/Proofsuite/amp-matching-engine/types"
	"github.com/Proofsuite/amp-matching-engine/utils/math"
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
	bookEntryAvailableAmount := math.Sub(bookEntry.Amount, bookEntry.FilledAmount)
	orderAvailableAmount := math.Sub(order.Amount, order.FilledAmount)

	if math.IsGreaterThan(bookEntryAvailableAmount, orderAvailableAmount) {
		fillOrder.Amount = orderAvailableAmount
		bookEntry.FilledAmount = math.Add(bookEntry.FilledAmount, orderAvailableAmount)
		bookEntry.Status = "PARTIAL_FILLED"
		fillOrder.Order = bookEntry

		err := e.updateOrder(bookEntry, fillOrder.Amount)
		if err != nil {
			log.Print(err)
			return nil, nil, err
		}

	} else {
		fillOrder.Amount = bookEntryAvailableAmount
		bookEntry.FilledAmount = math.Add(bookEntry.FilledAmount, bookEntryAvailableAmount)
		bookEntry.Status = "FILLED"
		fillOrder.Order = bookEntry

		err := e.deleteOrder(bookEntry, fillOrder.Amount)
		if err != nil {
			log.Print(err)
			return nil, nil, err
		}
	}

	order.FilledAmount = math.Add(order.FilledAmount, fillOrder.Amount)
	// Create trade object to be passed to the system for further processing
	trade = &types.Trade{
		Amount:       fillOrder.Amount,
		Price:        order.PricePoint,
		BaseToken:    order.BaseToken,
		QuoteToken:   order.QuoteToken,
		OrderHash:    bookEntry.Hash,
		Side:         order.Side,
		Taker:        order.UserAddress,
		PairName:     order.PairName,
		Maker:        bookEntry.UserAddress,
		TakerOrderID: order.ID,
		MakerOrderID: bookEntry.ID,
		TradeNonce:   big.NewInt(0),
		Signature:    &types.Signature{},
	}

	trade.Hash = trade.ComputeHash()
	return
}
