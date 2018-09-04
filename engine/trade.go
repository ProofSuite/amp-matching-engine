package engine

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/Proofsuite/amp-matching-engine/types"
	"github.com/Proofsuite/amp-matching-engine/utils/math"
)

// FillStatus is enum used to signify the filled status of order in engineResponse
type FillStatus int

// Response is the structure of message response sent by engine
type Response struct {
	FillStatus     types.FillStatus   `json:"fillStatus,omitempty"`
	Order          *types.Order       `json:"order,omitempty"`
	RemainingOrder *types.Order       `json:"remainingOrder,omitempty"`
	MatchingOrders []*types.FillOrder `json:"matchingOrders,omitempty"`
	Trades         []*types.Trade     `json:"trades,omitempty"`
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
func (e *Engine) execute(order *types.Order, bookEntry *types.Order) (*types.Trade, *types.FillOrder, error) {
	fillOrder := &types.FillOrder{}
	trade := &types.Trade{}

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

	trade = &types.Trade{
		Amount:     fillOrder.Amount,
		Price:      order.PricePoint,
		PricePoint: order.PricePoint,
		//NOTE: I don't think these are publicly needed but leaving this until confirmation
		BaseToken:  order.BaseToken,
		QuoteToken: order.QuoteToken,
		OrderHash:  bookEntry.Hash,
		Side:       order.Side,
		Taker:      order.UserAddress,
		PairName:   order.PairName,
		Maker:      bookEntry.UserAddress,
	}

	return trade, fillOrder, nil
}

func (resp *Response) Print() {
	b, err := json.MarshalIndent(resp, "", "  ")
	if err != nil {
		log.Print(err)
	}

	fmt.Print("\n", string(b))
}
