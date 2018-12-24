// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import bson "github.com/globalsign/mgo/bson"
import common "github.com/ethereum/go-ethereum/common"

import mock "github.com/stretchr/testify/mock"
import types "github.com/Proofsuite/amp-matching-engine/types"

// TradeDao is an autogenerated mock type for the TradeDao type
type TradeDao struct {
	mock.Mock
}

// Aggregate provides a mock function with given fields: q
func (_m *TradeDao) Aggregate(q []bson.M) ([]*types.Tick, error) {
	ret := _m.Called(q)

	var r0 []*types.Tick
	if rf, ok := ret.Get(0).(func([]bson.M) []*types.Tick); ok {
		r0 = rf(q)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*types.Tick)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func([]bson.M) error); ok {
		r1 = rf(q)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Create provides a mock function with given fields: o
func (_m *TradeDao) Create(o ...*types.Trade) error {
	_va := make([]interface{}, len(o))
	for _i := range o {
		_va[_i] = o[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 error
	if rf, ok := ret.Get(0).(func(...*types.Trade) error); ok {
		r0 = rf(o...)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Drop provides a mock function with given fields:
func (_m *TradeDao) Drop() {
	_m.Called()
}

// FindAndModify provides a mock function with given fields: h, t
func (_m *TradeDao) FindAndModify(h common.Hash, t *types.Trade) (*types.Trade, error) {
	ret := _m.Called(h, t)

	var r0 *types.Trade
	if rf, ok := ret.Get(0).(func(common.Hash, *types.Trade) *types.Trade); ok {
		r0 = rf(h, t)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*types.Trade)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(common.Hash, *types.Trade) error); ok {
		r1 = rf(h, t)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetAll provides a mock function with given fields:
func (_m *TradeDao) GetAll() ([]types.Trade, error) {
	ret := _m.Called()

	var r0 []types.Trade
	if rf, ok := ret.Get(0).(func() []types.Trade); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]types.Trade)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetAllTradesByPairAddress provides a mock function with given fields: bt, qt
func (_m *TradeDao) GetAllTradesByPairAddress(bt common.Address, qt common.Address) ([]*types.Trade, error) {
	ret := _m.Called(bt, qt)

	var r0 []*types.Trade
	if rf, ok := ret.Get(0).(func(common.Address, common.Address) []*types.Trade); ok {
		r0 = rf(bt, qt)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*types.Trade)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(common.Address, common.Address) error); ok {
		r1 = rf(bt, qt)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetByHash provides a mock function with given fields: h
func (_m *TradeDao) GetByHash(h common.Hash) (*types.Trade, error) {
	ret := _m.Called(h)

	var r0 *types.Trade
	if rf, ok := ret.Get(0).(func(common.Hash) *types.Trade); ok {
		r0 = rf(h)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*types.Trade)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(common.Hash) error); ok {
		r1 = rf(h)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetByMakerOrderHash provides a mock function with given fields: h
func (_m *TradeDao) GetByMakerOrderHash(h common.Hash) ([]*types.Trade, error) {
	ret := _m.Called(h)

	var r0 []*types.Trade
	if rf, ok := ret.Get(0).(func(common.Hash) []*types.Trade); ok {
		r0 = rf(h)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*types.Trade)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(common.Hash) error); ok {
		r1 = rf(h)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetByOrderHashes provides a mock function with given fields: hashes
func (_m *TradeDao) GetByOrderHashes(hashes []common.Hash) ([]*types.Trade, error) {
	ret := _m.Called(hashes)

	var r0 []*types.Trade
	if rf, ok := ret.Get(0).(func([]common.Hash) []*types.Trade); ok {
		r0 = rf(hashes)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*types.Trade)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func([]common.Hash) error); ok {
		r1 = rf(hashes)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetByPairName provides a mock function with given fields: name
func (_m *TradeDao) GetByPairName(name string) ([]*types.Trade, error) {
	ret := _m.Called(name)

	var r0 []*types.Trade
	if rf, ok := ret.Get(0).(func(string) []*types.Trade); ok {
		r0 = rf(name)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*types.Trade)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(name)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetByTakerOrderHash provides a mock function with given fields: h
func (_m *TradeDao) GetByTakerOrderHash(h common.Hash) ([]*types.Trade, error) {
	ret := _m.Called(h)

	var r0 []*types.Trade
	if rf, ok := ret.Get(0).(func(common.Hash) []*types.Trade); ok {
		r0 = rf(h)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*types.Trade)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(common.Hash) error); ok {
		r1 = rf(h)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetByUserAddress provides a mock function with given fields: a
func (_m *TradeDao) GetByUserAddress(a common.Address) ([]*types.Trade, error) {
	ret := _m.Called(a)

	var r0 []*types.Trade
	if rf, ok := ret.Get(0).(func(common.Address) []*types.Trade); ok {
		r0 = rf(a)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*types.Trade)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(common.Address) error); ok {
		r1 = rf(a)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetNTradesByPairAddress provides a mock function with given fields: bt, qt, n
func (_m *TradeDao) GetNTradesByPairAddress(bt common.Address, qt common.Address, n int) ([]*types.Trade, error) {
	ret := _m.Called(bt, qt, n)

	var r0 []*types.Trade
	if rf, ok := ret.Get(0).(func(common.Address, common.Address, int) []*types.Trade); ok {
		r0 = rf(bt, qt, n)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*types.Trade)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(common.Address, common.Address, int) error); ok {
		r1 = rf(bt, qt, n)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetSortedTrades provides a mock function with given fields: bt, qt, n
func (_m *TradeDao) GetSortedTrades(bt common.Address, qt common.Address, n int) ([]*types.Trade, error) {
	ret := _m.Called(bt, qt, n)

	var r0 []*types.Trade
	if rf, ok := ret.Get(0).(func(common.Address, common.Address, int) []*types.Trade); ok {
		r0 = rf(bt, qt, n)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*types.Trade)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(common.Address, common.Address, int) error); ok {
		r1 = rf(bt, qt, n)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetSortedTradesByUserAddress provides a mock function with given fields: a, limit
func (_m *TradeDao) GetSortedTradesByUserAddress(a common.Address, limit ...int) ([]*types.Trade, error) {
	_va := make([]interface{}, len(limit))
	for _i := range limit {
		_va[_i] = limit[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, a)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 []*types.Trade
	if rf, ok := ret.Get(0).(func(common.Address, ...int) []*types.Trade); ok {
		r0 = rf(a, limit...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*types.Trade)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(common.Address, ...int) error); ok {
		r1 = rf(a, limit...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetTradesByPairAddress provides a mock function with given fields: bt, qt, n
func (_m *TradeDao) GetTradesByPairAddress(bt common.Address, qt common.Address, n int) ([]*types.Trade, error) {
	ret := _m.Called(bt, qt, n)

	var r0 []*types.Trade
	if rf, ok := ret.Get(0).(func(common.Address, common.Address, int) []*types.Trade); ok {
		r0 = rf(bt, qt, n)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*types.Trade)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(common.Address, common.Address, int) error); ok {
		r1 = rf(bt, qt, n)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Update provides a mock function with given fields: t
func (_m *TradeDao) Update(t *types.Trade) error {
	ret := _m.Called(t)

	var r0 error
	if rf, ok := ret.Get(0).(func(*types.Trade) error); ok {
		r0 = rf(t)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// UpdateByHash provides a mock function with given fields: h, t
func (_m *TradeDao) UpdateByHash(h common.Hash, t *types.Trade) error {
	ret := _m.Called(h, t)

	var r0 error
	if rf, ok := ret.Get(0).(func(common.Hash, *types.Trade) error); ok {
		r0 = rf(h, t)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// UpdateTradeStatus provides a mock function with given fields: h, status
func (_m *TradeDao) UpdateTradeStatus(h common.Hash, status string) error {
	ret := _m.Called(h, status)

	var r0 error
	if rf, ok := ret.Get(0).(func(common.Hash, string) error); ok {
		r0 = rf(h, status)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// UpdateTradeStatuses provides a mock function with given fields: status, hashes
func (_m *TradeDao) UpdateTradeStatuses(status string, hashes ...common.Hash) ([]*types.Trade, error) {
	_va := make([]interface{}, len(hashes))
	for _i := range hashes {
		_va[_i] = hashes[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, status)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 []*types.Trade
	if rf, ok := ret.Get(0).(func(string, ...common.Hash) []*types.Trade); ok {
		r0 = rf(status, hashes...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*types.Trade)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, ...common.Hash) error); ok {
		r1 = rf(status, hashes...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdateTradeStatusesByOrderHashes provides a mock function with given fields: status, hashes
func (_m *TradeDao) UpdateTradeStatusesByOrderHashes(status string, hashes ...common.Hash) ([]*types.Trade, error) {
	_va := make([]interface{}, len(hashes))
	for _i := range hashes {
		_va[_i] = hashes[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, status)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 []*types.Trade
	if rf, ok := ret.Get(0).(func(string, ...common.Hash) []*types.Trade); ok {
		r0 = rf(status, hashes...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*types.Trade)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, ...common.Hash) error); ok {
		r1 = rf(status, hashes...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
