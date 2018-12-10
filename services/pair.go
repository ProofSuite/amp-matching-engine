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
	orderDao interfaces.OrderDao
	eng      interfaces.Engine
}

// NewPairService returns a new instance of balance service
func NewPairService(
	pairDao interfaces.PairDao,
	tokenDao interfaces.TokenDao,
	tradeDao interfaces.TradeDao,
	orderDao interfaces.OrderDao,
	eng interfaces.Engine,
) *PairService {

	return &PairService{pairDao, tokenDao, tradeDao, orderDao, eng}
}

// Create function is responsible for inserting new pair in DB.
// It checks for existence of tokens in DB first
func (s *PairService) Create(pair *types.Pair) error {
	p, err := s.pairDao.GetByTokenAddress(pair.BaseTokenAddress, pair.QuoteTokenAddress)
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
	pair.QuoteTokenDecimals = st.Decimals
	pair.BaseTokenSymbol = bt.Symbol
	pair.BaseTokenAddress = bt.ContractAddress
	pair.BaseTokenDecimals = bt.Decimals
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

func (s *PairService) GetAllTokenPairData() ([]*types.PairData, error) {
	now := time.Now()
	end := time.Unix(now.Unix(), 0)
	start := time.Unix(now.AddDate(0, 0, -7).Unix(), 0)
	one, _ := bson.ParseDecimal128("1")

	pairs, err := s.pairDao.GetActivePairs()
	if err != nil {
		return nil, err
	}

	tradeDataQuery := []bson.M{
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

	orderDataQuery := []bson.M{
		bson.M{
			"$match": bson.M{
				"status": bson.M{"$in": []string{"OPEN", "PARTIAL_FILLED"}},
			},
		},
		bson.M{
			"$group": bson.M{
				"_id": bson.M{
					"baseToken":  "$baseToken",
					"quoteToken": "$quoteToken",
					"pairName":   "$pairName",
				},
				"orderCount": bson.M{"$sum": one},
				"orderVolume": bson.M{
					"$sum": bson.M{
						"$subtract": []bson.M{bson.M{"$toDecimal": "$amount"}, bson.M{"$toDecimal": "$filledAmount"}},
					},
				},
			},
		},
	}

	tradeData, err := s.tradeDao.Aggregate(tradeDataQuery)
	if err != nil {
		return nil, err
	}

	orderData := []*types.OrderData{}
	err = s.orderDao.Aggregate(orderDataQuery, orderData)
	if err != nil {
		return nil, err
	}

	pairsData := []*types.PairData{}
	for _, p := range pairs {
		pairData := &types.PairData{Pair: types.PairID{p.Name(), p.BaseTokenAddress, p.QuoteTokenAddress}}

		for _, t := range tradeData {
			if t.AddressCode() == p.AddressCode() {
				pairData.Open = t.Open
				pairData.High = t.High
				pairData.Low = t.Low
				pairData.Volume = t.Volume
				pairData.Close = t.Close
				pairData.Count = t.Count

			}
		}

		for _, o := range orderData {
			if o.AddressCode() == p.AddressCode() {
				pairData.OrderVolume = o.OrderVolume
				pairData.OrderCount = o.OrderCount
			}
		}

		pairsData = append(pairsData, pairData)
	}

	return pairsData, nil
}
