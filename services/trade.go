package services

import (
	"github.com/Proofsuite/amp-matching-engine/daos"
	"github.com/Proofsuite/amp-matching-engine/types"
)

type TradeService struct {
	tradeDao *daos.TradeDao
}

func NewTradeService(TradeDao *daos.TradeDao) *TradeService {
	return &TradeService{TradeDao}
}

func (t *TradeService) GetByPairName(pairName string) ([]*types.Trade, error) {
	return t.tradeDao.GetByPairName(pairName)
}
