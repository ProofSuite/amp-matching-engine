package services

import (
	"math/big"
	"testing"

	"github.com/Proofsuite/amp-matching-engine/types"
	"github.com/Proofsuite/amp-matching-engine/utils/testutils"
	"github.com/Proofsuite/amp-matching-engine/utils/testutils/mocks"
	"github.com/ethereum/go-ethereum/common"
)

func TestCancelTrades(t *testing.T) {
	orderDao := new(mocks.OrderDao)
	pairDao := new(mocks.PairDao)
	accountDao := new(mocks.AccountDao)
	tradeDao := new(mocks.TradeDao)
	engine := new(mocks.Engine)

	orderService := NewOrderService(
		orderDao,
		pairDao,
		accountDao,
		tradeDao,
		engine,
	)

	t1 := testutils.GetTestTrade1()
	t2 := testutils.GetTestTrade2()
	o1 := testutils.GetTestOrder1()
	o2 := testutils.GetTestOrder2()

	trades := []*types.Trade{&t1, &t2}
	hashes := []common.Hash{t1.OrderHash, t2.OrderHash}
	amounts := []*big.Int{t1.Amount, t2.Amount}
	orders := []*types.Order{&o1, &o2}

	orderDao.On("GetByHashes", hashes).Return(orders, nil)
	engine.On("CancelTrades", orders, amounts).Return(nil)

	err := orderService.CancelTrades(trades)
	if err != nil {
		t.Error("Could not cancel trades", err)
	}

	orderDao.AssertCalled(t, "GetByHashes", hashes)
	engine.AssertCalled(t, "CancelTrades", orders, amounts)
}
