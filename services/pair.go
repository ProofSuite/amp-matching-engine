package services

import (
	"time"

	"github.com/Proofsuite/amp-matching-engine/interfaces"
	"github.com/Proofsuite/amp-matching-engine/types"
	"github.com/ethereum/go-ethereum/common"
	"gopkg.in/mgo.v2/bson"
)

// PairService struct with daos required, responsible for communicating with daos.
// PairService functions are responsible for interacting with daos and implements business logics.
type PairService struct {
	pairDao  interfaces.PairDao
	tokenDao interfaces.TokenDao
	tradeDao interfaces.TradeDao
	eng      interfaces.Engine
}

// NewPairService returns a new instance of balance service
func NewPairService(
	pairDao interfaces.PairDao,
	tokenDao interfaces.TokenDao,
	tradeDao interfaces.TradeDao,
	eng interfaces.Engine,
) *PairService {

	return &PairService{pairDao, tokenDao, tradeDao, eng}
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

func (s *PairService) GetTokenPairData(bt, qt common.Address) ([]*types.Tick, error) {
	now := time.Now()
	end := time.Unix(now.Unix(), 0)
	start := time.Unix(now.AddDate(0, 0, -7).Unix(), 0)
	one, _ := bson.ParseDecimal128("1")

	q := []bson.M{
		bson.M{
			"$match": bson.M{
				"createdAt": bson.M{
					"$gte": start,
					"$lt":  end,
				},
				"status":     bson.M{"$in": []string{"SUCCESS"}},
				"baseToken":  bt.Hex(),
				"quoteToken": qt.Hex(),
			},
		},
		bson.M{
			"$group": bson.M{
				"_id": bson.M{
					"baseToken":  "$baseToken",
					"pairName":   "$pairName",
					"quoteToken": "$quoteToken",
				},
				"count":  bson.M{"$sum": one},
				"open":   bson.M{"$first": bson.M{"$toDecimal": "$pricepoint"}},
				"high":   bson.M{"$max": bson.M{"$toDecimal": "$pricepoint"}},
				"low":    bson.M{"$min": bson.M{"$toDecimal": "$pricepoint"}},
				"close":  bson.M{"$last": bson.M{"$toDecimal": "$pricepoint"}},
				"volume": bson.M{"$sum": bson.M{"$toDecimal": "$amount"}},
			},
		},
	}

	res, err := s.tradeDao.Aggregate(q)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (s *PairService) GetAllTokenPairData() ([]*types.Tick, error) {
	now := time.Now()
	end := time.Unix(now.Unix(), 0)
	start := time.Unix(now.AddDate(0, 0, -7).Unix(), 0)
	one, _ := bson.ParseDecimal128("1")

	q := []bson.M{
		bson.M{
			"$match": bson.M{
				"createdAt": bson.M{
					"$gte": start,
					"$lt":  end,
				},
				"status": bson.M{"$in": []string{"SUCCESS"}},
			},
		},
		bson.M{
			"$group": bson.M{
				"_id": bson.M{
					"baseToken":  "$baseToken",
					"pairName":   "$pairName",
					"quoteToken": "$quoteToken",
				},
				"count":  bson.M{"$sum": one},
				"open":   bson.M{"$first": bson.M{"$toDecimal": "$pricepoint"}},
				"high":   bson.M{"$max": bson.M{"$toDecimal": "$pricepoint"}},
				"low":    bson.M{"$min": bson.M{"$toDecimal": "$pricepoint"}},
				"close":  bson.M{"$last": bson.M{"$toDecimal": "$pricepoint"}},
				"volume": bson.M{"$sum": bson.M{"$toDecimal": "$amount"}},
			},
		},
	}

	res, err := s.tradeDao.Aggregate(q)
	if err != nil {
		return nil, err
	}

	return res, nil
}
