package services

import (
	"github.com/Proofsuite/amp-matching-engine/interfaces"
	"github.com/ethereum/go-ethereum/common"

	"gopkg.in/mgo.v2/bson"

	"github.com/Proofsuite/amp-matching-engine/types"
)

// PairService struct with daos required, responsible for communicating with daos.
// PairService functions are responsible for interacting with daos and implements business logics.
type PairService struct {
	pairDao      interfaces.PairDao
	tokenDao     interfaces.TokenDao
	eng          interfaces.Engine
	tradeService *TradeService
}

// NewPairService returns a new instance of balance service
func NewPairService(
	pairDao interfaces.PairDao,
	tokenDao interfaces.TokenDao,
	eng interfaces.Engine,
	tradeService *TradeService,
) *PairService {

	return &PairService{pairDao, tokenDao, eng, tradeService}
}

// Create function is responsible for inserting new pair in DB.
// It checks for existence of tokens in DB first
func (s *PairService) Create(pair *types.Pair) error {
	p, err := s.pairDao.GetByBuySellTokenAddress(pair.BaseTokenAddress, pair.QuoteTokenAddress)
	if err != nil {
		return err
	}

	if p != nil {
		return ErrPairExists
	}

	bt, err := s.tokenDao.GetByAddress(pair.BaseTokenAddress)
	if err != nil {
		return err
	}

	if bt == nil {
		return ErrBaseTokenNotFound
	}

	st, err := s.tokenDao.GetByAddress(pair.QuoteTokenAddress)
	if err != nil {
		return err
	}

	if st == nil {
		return ErrQuoteTokenNotFound
	}

	if !st.Quote {
		return ErrQuoteTokenInvalid
	}

	pair.QuoteTokenSymbol = st.Symbol
	pair.QuoteTokenAddress = st.ContractAddress
	pair.QuoteTokenDecimal = st.Decimal
	pair.BaseTokenSymbol = bt.Symbol
	pair.BaseTokenAddress = bt.ContractAddress
	pair.BaseTokenDecimal = bt.Decimal
	err = s.pairDao.Create(pair)
	if err != nil {
		return err
	}

	return nil
}

// GetByID fetches details of a pair using its mongo ID
func (s *PairService) GetByID(id bson.ObjectId) (*types.Pair, error) {
	return s.pairDao.GetByID(id)
}

// GetByTokenAddress fetches details of a pair using contract address of
// its constituting tokens
func (s *PairService) GetByTokenAddress(bt, qt common.Address) (*types.Pair, error) {
	return s.pairDao.GetByTokenAddress(bt, qt)
}

// GetAll is reponsible for fetching all the pairs in the DB
func (s *PairService) GetAll() ([]types.Pair, error) {
	return s.pairDao.GetAll()
}
