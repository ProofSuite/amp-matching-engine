package services

import (
	"github.com/Proofsuite/amp-matching-engine/interfaces"
	"github.com/Proofsuite/amp-matching-engine/types"
	"github.com/Proofsuite/amp-matching-engine/utils"
	"github.com/Proofsuite/amp-matching-engine/ws"

	"github.com/ethereum/go-ethereum/common"
	eth "github.com/ethereum/go-ethereum/core/types"
)

// TradeService struct with daos required, responsible for communicating with daos.
// TradeService functions are responsible for interacting with daos and implements business logics.
type TradeService struct {
	tradeDao interfaces.TradeDao
}

// NewTradeService returns a new instance of TradeService
func NewTradeService(TradeDao interfaces.TradeDao) *TradeService {
	return &TradeService{TradeDao}
}

// GetByPairName fetches all the trades corresponding to a pair using pair's name
func (t *TradeService) GetByPairName(pairName string) ([]*types.Trade, error) {
	return t.tradeDao.GetByPairName(pairName)
}

// GetTrades is currently not implemented correctly
func (t *TradeService) GetTrades(bt, qt common.Address) ([]types.Trade, error) {
	return t.tradeDao.GetAll()
}

// GetByPairAddress fetches all the trades corresponding to a pair using pair's token address
func (t *TradeService) GetByPairAddress(bt, qt common.Address) ([]*types.Trade, error) {
	return t.tradeDao.GetByPairAddress(bt, qt)
}

// GetByUserAddress fetches all the trades corresponding to a user address
func (t *TradeService) GetByUserAddress(addr common.Address) ([]*types.Trade, error) {
	return t.tradeDao.GetByUserAddress(addr)
}

// GetByHash fetches all trades corresponding to a trade hash
func (t *TradeService) GetByHash(hash common.Hash) (*types.Trade, error) {
	return t.tradeDao.GetByHash(hash)
}

// GetByOrderHash fetches all trades corresponding to an order hash
func (t *TradeService) GetByOrderHash(hash common.Hash) ([]*types.Trade, error) {
	return t.tradeDao.GetByOrderHash(hash)
}

func (t *TradeService) UpdateTradeTx(tr *types.Trade, tx *eth.Transaction) error {
	tr.Tx = tx

	err := t.tradeDao.Update(tr)
	if err != nil {
		return err
	}

	return nil
}

// Subscribe
func (t *TradeService) Subscribe(conn *ws.Conn, bt, qt common.Address) {
	socket := ws.GetTradeSocket()

	trades, err := t.GetTrades(bt, qt)
	if err != nil {
		socket.SendErrorMessage(conn, err.Error())
		return
	}

	id := utils.GetTradeChannelID(bt, qt)
	err = socket.Subscribe(id, conn)
	if err != nil {
		message := map[string]string{
			"Code":    "UNABLE_TO_REGISTER",
			"Message": "UNABLE_TO_REGISTER " + err.Error(),
		}

		socket.SendErrorMessage(conn, message)
		return
	}

	ws.RegisterConnectionUnsubscribeHandler(conn, socket.UnsubscribeHandler(id))
	socket.SendInitMessage(conn, trades)
}

// Unsubscribe
func (t *TradeService) Unsubscribe(conn *ws.Conn, bt, qt common.Address) {
	socket := ws.GetTradeSocket()

	id := utils.GetTradeChannelID(bt, qt)
	socket.Unsubscribe(id, conn)
}
