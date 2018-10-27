package types

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto/sha3"
)

type Match struct {
	Order *Order `json:"order"`
	Trade *Trade `json:"trade"`
}

type Matches []*Match

func NewMatches(orders []*Order, trades []*Trade) *Matches {
	m := Matches{}
	for i, _ := range orders {
		m = append(m, &Match{orders[i], trades[i]})
	}

	return &m
}

func (m Matches) Maker() common.Address {
	return m[0].Trade.Maker
}

func (m Matches) Taker() common.Address {
	return m[0].Trade.Taker
}

func (m Matches) Trades() []*Trade {
	trades := []*Trade{}
	for _, match := range m {
		trades = append(trades, match.Trade)
	}

	return trades
}

func (m Matches) PairCode() (string, error) {
	return m[0].Order.PairCode()
}

func (m Matches) Orders() []*Order {
	orders := []*Order{}
	for _, match := range m {
		orders = append(orders, match.Order)
	}

	return orders
}

func (m Matches) TradeAmounts() []*big.Int {
	amounts := []*big.Int{}
	for _, match := range m {
		amounts = append(amounts, match.Trade.Amount)
	}

	return amounts
}

func (m Matches) Validate() error {
	trades := m.Trades()
	orders := m.Orders()

	for _, t := range trades {
		err := t.ValidateComplete()
		if err != nil {
			logger.Error(err)
			return err
		}
	}

	for _, o := range orders {
		err := o.ValidateComplete()
		if err != nil {
			logger.Error(err)
			return err
		}
	}

	return nil
}

func (m Matches) HashID() common.Hash {
	sha := sha3.NewKeccak256()
	for _, match := range m {
		sha.Write(match.Order.Hash.Bytes())
		sha.Write(match.Trade.Hash.Bytes())
	}

	return common.BytesToHash(sha.Sum(nil))
}

type EngineResponse struct {
	Status            string    `json:"fillStatus,omitempty"`
	Order             *Order    `json:"order,omitempty"`
	RemainingOrder    *Order    `json:"remainingOrder,omitempty"`
	Matches           *Matches  `json:"matches,omitempty"`
	RecoveredOrders   *[]*Order `json:"recoveredOrders,omitempty"`
	InvalidatedOrders *[]*Order `json:"invalidatedOrders,omitempty"`
	CancelledTrades   *[]*Trade `json:"cancelledTrades,omitempty"`
}

func (r *EngineResponse) AppendMatch(m *Match) {
	if r.Matches == nil {
		r.Matches = &Matches{m}
		return
	}

	matches := append(*r.Matches, m)
	r.Matches = &matches
}

func (r *EngineResponse) AppendMatches(m []*Match) {
	if r.Matches == nil {
		matches := Matches(m)
		r.Matches = &matches
		return
	}

	matches := append(*r.Matches, m...)
	r.Matches = &matches
}

func (r *EngineResponse) HashID() common.Hash {
	if r.Status == "NOMATCH" {
		return r.Order.Hash
	}

	sha := sha3.NewKeccak256()
	for _, m := range *r.Matches {
		sha.Write(m.Order.Hash.Bytes())
		sha.Write(m.Trade.Hash.Bytes())
	}

	return common.BytesToHash(sha.Sum(nil))
}
