package engine

import (
	"math/big"

	"github.com/Proofsuite/amp-matching-engine/types"
	"github.com/Proofsuite/amp-matching-engine/utils/math"
)

// execute function is responsible for executing of matched orders
// i.e it deletes/updates orders in case of order matching and responds
// with trade instance and fillOrder
func (e *Engine) execute(order *types.Order, bookEntry *types.Order) (*types.Trade, error) {
	// fillOrder := &types.FillOrder{}
	trade := &types.Trade{}
	tradeAmount := big.NewInt(0)
	bookEntryAvailableAmount := math.Sub(bookEntry.Amount, bookEntry.FilledAmount)
	orderAvailableAmount := math.Sub(order.Amount, order.FilledAmount)

	if math.IsGreaterThan(bookEntryAvailableAmount, orderAvailableAmount) {
		tradeAmount = orderAvailableAmount
		bookEntry.FilledAmount = math.Add(bookEntry.FilledAmount, orderAvailableAmount)
		bookEntry.Status = "PARTIAL_FILLED"

		err := e.updateOrder(bookEntry, tradeAmount)
		if err != nil {
			logger.Error(err)
			return nil, err
		}

	} else {
		tradeAmount = bookEntryAvailableAmount
		bookEntry.FilledAmount = math.Add(bookEntry.FilledAmount, bookEntryAvailableAmount)
		bookEntry.Status = "FILLED"

		err := e.deleteOrder(bookEntry, tradeAmount)
		if err != nil {
			logger.Error(err)
			return nil, err
		}
	}

	order.FilledAmount = math.Add(order.FilledAmount, tradeAmount)
	trade = &types.Trade{
		Amount:     tradeAmount,
		Price:      order.PricePoint,
		PricePoint: order.PricePoint,
		BaseToken:  order.BaseToken,
		QuoteToken: order.QuoteToken,
		OrderHash:  bookEntry.Hash,
		Side:       order.Side,
		Taker:      order.UserAddress,
		PairName:   order.PairName,
		Maker:      bookEntry.UserAddress,
	}

	return trade, nil
}
