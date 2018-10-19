package types

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto/sha3"
)

type OrderTradePair struct {
	Order *Order `json:"order"`
	Trade *Trade `json:"trade"`
}

type OrderTradePairs []OrderTradePair

type Matches struct {
	OrderTradePairs []*OrderTradePair
	HashID          common.Hash
}

func NewMatches(orders []*Order, trades []*Trade) *Matches {
	m := &Matches{}
	for i, _ := range orders {
		m.OrderTradePairs = append(m.OrderTradePairs, &OrderTradePair{orders[i], trades[i]})
	}

	return m
}

func (m *Matches) Maker() common.Address {
	return m.OrderTradePairs[0].Trade.Maker
}

func (m *Matches) Taker() common.Address {
	return m.OrderTradePairs[0].Trade.Taker
}

func (m *Matches) Trades() []*Trade {
	trades := []*Trade{}
	for _, m := range m.OrderTradePairs {
		trades = append(trades, m.Trade)
	}

	return trades
}

func (m *Matches) Orders() []*Order {
	orders := []*Order{}
	for _, m := range m.OrderTradePairs {
		orders = append(orders, m.Order)
	}

	return orders
}

func (m *Matches) Validate() error {
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

func (m *Matches) AppendMatch(ot *OrderTradePair) {
	if m.OrderTradePairs == nil {
		m.OrderTradePairs = make([]*OrderTradePair, 0)
	}

	m.OrderTradePairs = append(m.OrderTradePairs, ot)
}

func (m *Matches) AppendMatches(ots []*OrderTradePair) {
	if m.OrderTradePairs == nil {
		m.OrderTradePairs = make([]*OrderTradePair, 0)
	}

	for _, ot := range m.OrderTradePairs {
		m.OrderTradePairs = append(m.OrderTradePairs, ot)
	}
}

func (m *Matches) ComputeHashID() common.Hash {
	sha := sha3.NewKeccak256()
	for _, m := range m.OrderTradePairs {
		sha.Write(m.Order.Hash.Bytes())
		sha.Write(m.Trade.Hash.Bytes())
	}

	m.HashID = common.BytesToHash(sha.Sum(nil))

	return m.HashID
}

type EngineResponse struct {
	Status         string      `json:"fillStatus,omitempty"`
	HashID         common.Hash `json:"hashID, omitempty"`
	Order          *Order      `json:"order,omitempty"`
	RemainingOrder *Order      `json:"remainingOrder,omitempty"`
	Matches        *Matches    `json:"matches,omitempty"`
}

func (r *EngineResponse) AppendMatch(ot *OrderTradePair) {
	if r.Matches == nil {
		r.Matches = &Matches{OrderTradePairs: make([]*OrderTradePair, 0)}
	}

	r.Matches.AppendMatch(ot)
}

func (r *EngineResponse) AppendMatches(ots []*OrderTradePair) {
	r.Matches.AppendMatches(ots)
}

func (r *EngineResponse) ComputeHashID() common.Hash {
	sha := sha3.NewKeccak256()
	for _, m := range r.Matches.OrderTradePairs {
		sha.Write(m.Order.Hash.Bytes())
		sha.Write(m.Trade.Hash.Bytes())
	}

	return common.BytesToHash(sha.Sum(nil))
}
