package services

import (
	"strings"

	"github.com/Proofsuite/amp-matching-engine/engine"
	"github.com/ethereum/go-ethereum/common"

	"gopkg.in/mgo.v2/bson"

	"github.com/Proofsuite/amp-matching-engine/daos"
	aerrors "github.com/Proofsuite/amp-matching-engine/errors"
	"github.com/Proofsuite/amp-matching-engine/types"
)

// PairService struct with daos required, responsible for communicating with daos.
// PairService functions are responsible for interacting with daos and implements business logics.
type PairService struct {
	pairDao      daos.PairDaoInterface
	tokenDao     daos.TokenDaoInterface
	eng          engine.EngineInterface
	tradeService *TradeService
}

type PairServiceInterface interface {
	Create(pair *types.Pair) error
	GetByID(id bson.ObjectId) (*types.Pair, error)
	GetByTokenAddress(bt, qt common.Address) (*types.Pair, error)
	GetAll() ([]types.Pair, error)
}

// NewPairService returns a new instance of balance service
func NewPairService(
	pairDao daos.PairDaoInterface,
	tokenDao daos.TokenDaoInterface,
	eng engine.EngineInterface,
	tradeService *TradeService,
) *PairService {

	return &PairService{pairDao, tokenDao, eng, tradeService}
}

// Create function is responsible for inserting new pair in DB.
// It checks for existence of tokens in DB first
func (s *PairService) Create(pair *types.Pair) error {
	p, err := s.pairDao.GetByBuySellTokenAddress(pair.BaseTokenAddress, pair.QuoteTokenAddress)
	if err != nil && err.Error() != "NO_PAIR_FOUND" {
		return aerrors.NewAPIError(400, err.Error(), nil)
	} else if p != nil {
		return aerrors.NewAPIError(401, "PAIR_ALREADY_EXISTS", nil)
	}

	bt, err := s.tokenDao.GetByAddress(pair.BaseTokenAddress)
	if err != nil {
		return aerrors.NewAPIError(400, err.Error(), nil)
	}

	if bt == nil {
		return aerrors.NewAPIError(401, "BaseTokenAddress_DOESNT_EXIST", nil)
	}

	st, err := s.tokenDao.GetByAddress(pair.QuoteTokenAddress)
	if err != nil {
		return aerrors.NewAPIError(400, err.Error(), nil)
	}

	if st == nil {
		return aerrors.NewAPIError(401, "QuoteTokenAddress_DOESNT_EXIST", nil)
	}

	if !st.Quote {
		return aerrors.NewAPIError(401, "QuoteTokenAddress_CAN_NOT_BE_USED_AS_QUOTE_TOKEN", nil)
	}

	pair.QuoteTokenSymbol = st.Symbol
	pair.QuoteTokenAddress = st.ContractAddress
	pair.QuoteTokenDecimal = st.Decimal
	pair.BaseTokenSymbol = bt.Symbol
	pair.BaseTokenAddress = bt.ContractAddress
	pair.BaseTokenDecimal = bt.Decimal
	pair.Name = strings.ToUpper(st.Symbol + "/" + bt.Symbol)

	err = s.pairDao.Create(pair)
	return err

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
