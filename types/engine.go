package types

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto/sha3"
)

type Matches struct {
	MakerOrders []*Order `json:"makerOrders"`
	TakerOrder  *Order   `json:"takerOrder"`
	Trades      []*Trade `json:"trades"`
}

func NewMatches(makerOrders []*Order, takerOrder *Order, trades []*Trade) *Matches {
	return &Matches{
		MakerOrders: makerOrders,
		TakerOrder:  takerOrder,
		Trades:      trades,
	}
}

func (m Matches) NthMatch(i int) *Matches {
	return NewMatches(
		[]*Order{m.MakerOrders[i]},
		m.TakerOrder,
		[]*Trade{m.Trades[i]},
	)
}

func (m Matches) Taker() common.Address {
	return m.TakerOrder.UserAddress
}

func (m Matches) TakerOrderHash() common.Hash {
	return m.TakerOrder.Hash
}

func (m Matches) PairCode() (string, error) {
	return m.TakerOrder.PairCode()
}

func (m Matches) TradeAmounts() []*big.Int {
	amounts := []*big.Int{}
	for _, t := range m.Trades {
		amounts = append(amounts, t.Amount)
	}

	return amounts
}

func (m Matches) Length() int {
	return len(m.Trades)
}

func (m Matches) AppendMatch(mo *Order, t *Trade) {
	m.MakerOrders = append(m.MakerOrders, mo)
	m.Trades = append(m.Trades, t)
}

func (m Matches) Validate() error {
	for _, t := range m.Trades {
		err := t.ValidateComplete()
		if err != nil {
			logger.Error(err)
			return err
		}
	}

	for _, mo := range m.MakerOrders {
		err := mo.ValidateComplete()
		if err != nil {
			logger.Error(err)
			return err
		}
	}

	err := m.TakerOrder.ValidateComplete()
	if err != nil {
		logger.Error(err)
		return err
	}

	return nil
}

func (m Matches) HashID() common.Hash {
	sha := sha3.NewKeccak256()
	for _, t := range m.Trades {
		sha.Write(t.Hash.Bytes())
	}

	return common.BytesToHash(sha.Sum(nil))
}

type EngineResponse struct {
	Status            string    `json:"fillStatus,omitempty"`
	Order             *Order    `json:"order,omitempty"`
	Matches           *Matches  `json:"matches,omitempty"`
	RecoveredOrders   *[]*Order `json:"recoveredOrders,omitempty"`
	InvalidatedOrders *[]*Order `json:"invalidatedOrders,omitempty"`
	CancelledTrades   *[]*Trade `json:"cancelledTrades,omitempty"`
}

func (r *EngineResponse) AppendMatch(mo *Order, t *Trade) {
	if r.Matches == nil {
		r.Matches = &Matches{}
	}

	r.Matches.MakerOrders = append(r.Matches.MakerOrders, mo)
	r.Matches.Trades = append(r.Matches.Trades, t)
}

func (r *EngineResponse) AppendMatches(mo []*Order, t []*Trade) {
	if r.Matches == nil {
		r.Matches = &Matches{}
	}

	r.Matches.MakerOrders = append(r.Matches.MakerOrders, mo...)
	r.Matches.Trades = append(r.Matches.Trades, t...)

}

func (r *EngineResponse) HashID() common.Hash {
	if r.Status == "NOMATCH" {
		return r.Order.Hash
	}

	sha := sha3.NewKeccak256()
	for _, t := range r.Matches.Trades {
		sha.Write(t.Hash.Bytes())
	}

	return common.BytesToHash(sha.Sum(nil))
}
