package services

import (
	"log"
	"math/big"
	"time"

	"github.com/Proofsuite/amp-matching-engine/interfaces"
	"github.com/Proofsuite/amp-matching-engine/types"
	"github.com/Proofsuite/amp-matching-engine/utils"
	"github.com/Proofsuite/amp-matching-engine/utils/math"
	"github.com/globalsign/mgo/bson"
)

type InfoService struct {
	pairDao      interfaces.PairDao
	tokenDao     interfaces.TokenDao
	tradeDao     interfaces.TradeDao
	orderDao     interfaces.OrderDao
	priceService interfaces.PriceService
}

func NewInfoService(
	pairDao interfaces.PairDao,
	tokenDao interfaces.TokenDao,
	tradeDao interfaces.TradeDao,
	orderDao interfaces.OrderDao,
	priceService interfaces.PriceService,
) *InfoService {

	return &InfoService{
		pairDao,
		tokenDao,
		tradeDao,
		orderDao,
		priceService,
	}
}

func (s *InfoService) GetExchangeStats() (*types.ExchangeStats, error) {
	now := time.Now()
	end := time.Unix(now.Unix(), 0)
	start := time.Unix(now.AddDate(0, 0, -7).Unix(), 0)
	one, _ := bson.ParseDecimal128("1")

	tokens, err := s.tokenDao.GetBaseTokens()
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	quoteTokens, err := s.tokenDao.GetQuoteTokens()
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	erroredTradeCount, err := s.tradeDao.GetErroredTradeCount(start, end)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	tokenSymbols := []string{}
	for _, t := range tokens {
		tokenSymbols = append(tokenSymbols, t.Symbol)
	}

	for _, t := range quoteTokens {
		tokenSymbols = append(tokenSymbols, t.Symbol)
	}

	rates, err := s.priceService.GetDollarMarketPrices(tokenSymbols)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	pairs, err := s.pairDao.GetDefaultPairs()
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	tradesQuery := []bson.M{
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
					"pairName":   "$pairName",
					"baseToken":  "$baseToken",
					"quoteToken": "$quoteToken",
				},
				"count":  bson.M{"$sum": one},
				"close":  bson.M{"$last": "$pricepoint"},
				"volume": bson.M{"$sum": bson.M{"$toDecimal": "$amount"}},
			},
		},
	}

	bidsQuery := []bson.M{
		bson.M{
			"$match": bson.M{
				"status": bson.M{"$in": []string{"OPEN", "PARTIAL_FILLED"}},
				"side":   "BUY",
			},
		},
		bson.M{
			"$group": bson.M{
				"_id": bson.M{
					"pairName":   "$pairName",
					"baseToken":  "$baseToken",
					"quoteToken": "$quoteToken",
				},
				"orderCount": bson.M{"$sum": one},
				"orderVolume": bson.M{
					"$sum": bson.M{
						"$subtract": []bson.M{bson.M{"$toDecimal": "$amount"}, bson.M{"$toDecimal": "$filledAmount"}},
					},
				},
				"bestPrice": bson.M{"$max": "$pricepoint"},
			},
		},
	}

	asksQuery := []bson.M{
		bson.M{
			"$match": bson.M{
				"status": bson.M{"$in": []string{"OPEN", "PARTIAL_FILLED"}},
				"side":   "SELL",
			},
		},
		bson.M{
			"$group": bson.M{
				"_id": bson.M{
					"pairName":   "$pairName",
					"baseToken":  "$baseToken",
					"quoteToken": "$quoteToken",
				},
				"orderCount": bson.M{"$sum": one},
				"orderVolume": bson.M{
					"$sum": bson.M{
						"$subtract": []bson.M{bson.M{"$toDecimal": "$amount"}, bson.M{"$toDecimal": "$filledAmount"}},
					},
				},
				"bestPrice": bson.M{"$min": "$pricepoint"},
			},
		},
	}

	tradeData, err := s.tradeDao.Aggregate(tradesQuery)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	bidsData, err := s.orderDao.Aggregate(bidsQuery)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	asksData, err := s.orderDao.Aggregate(asksQuery)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	var totalOrders int
	var totalTrades int
	var totalVolume float64
	var totalBuyOrderAmount float64
	var totalSellOrderAmount float64
	var totalSellOrders int
	var totalBuyOrders int
	pairTradeCounts := map[string]int{}
	tokenTradeCounts := map[string]int{}

	for _, p := range pairs {
		for _, t := range tradeData {
			if t.AddressCode() == p.AddressCode() {
				totalTrades = totalTrades + int(t.Count.Int64())
				totalVolume = totalVolume + t.ConvertedVolume(&p, rates[p.BaseTokenSymbol])

				pairTradeCounts[p.Name()] = int(t.Count.Int64())
				tokenTradeCounts[p.BaseTokenSymbol] = tokenTradeCounts[p.BaseTokenSymbol] + int(t.Count.Int64())
			}
		}

		for _, o := range bidsData {
			if o.AddressCode() == p.AddressCode() {
				// change and replace by equivalent dollar volume instead of order count
				totalBuyOrderAmount = totalBuyOrderAmount + o.ConvertedVolume(&p, rates[p.BaseTokenSymbol])
				totalBuyOrders = totalBuyOrders + int(o.OrderCount.Int64())
			}
		}

		for _, o := range asksData {
			if o.AddressCode() == p.AddressCode() {
				// change and replace by equivalent dollar volume instead of order count
				totalSellOrderAmount = totalSellOrderAmount + o.ConvertedVolume(&p, rates[p.BaseTokenSymbol])
				totalSellOrders = totalSellOrders + int(o.OrderCount.Int64())
			}
		}
	}

	mostTradedToken, _ := utils.MaxIntMap(tokenTradeCounts)
	mostTradedPair, _ := utils.MaxIntMap(pairTradeCounts)
	totalOrders = totalSellOrders + totalBuyOrders
	totalOrderAmount := totalBuyOrderAmount + totalSellOrderAmount

	tradeSuccessRatio := float64(1)
	if totalTrades > 0 {
		tradeSuccessRatio = float64(totalTrades-erroredTradeCount) / float64(totalTrades)
	}

	stats := &types.ExchangeStats{
		TotalOrders:          totalOrders,
		TotalTrades:          totalTrades,
		TotalBuyOrderAmount:  totalBuyOrderAmount,
		TotalSellOrderAmount: totalSellOrderAmount,
		TotalVolume:          totalVolume,
		TotalOrderAmount:     totalOrderAmount,
		TotalBuyOrders:       totalSellOrders,
		TotalSellOrders:      totalBuyOrders,
		MostTradedToken:      mostTradedToken,
		MostTradedPair:       mostTradedPair,
		TradeSuccessRatio:    tradeSuccessRatio,
	}

	log.Printf("%+v\n", stats)

	return stats, nil
}

func (s *InfoService) GetExchangeData() (*types.ExchangeData, error) {
	now := time.Now()
	end := time.Unix(now.Unix(), 0)
	start := time.Unix(now.AddDate(0, 0, -7).Unix(), 0)
	one, _ := bson.ParseDecimal128("1")

	tokens, err := s.tokenDao.GetBaseTokens()
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	quoteTokens, err := s.tokenDao.GetQuoteTokens()
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	erroredTradeCount, err := s.tradeDao.GetErroredTradeCount(start, end)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	tokenSymbols := []string{}
	for _, t := range tokens {
		tokenSymbols = append(tokenSymbols, t.Symbol)
	}

	for _, t := range quoteTokens {
		tokenSymbols = append(tokenSymbols, t.Symbol)
	}

	rates, err := s.priceService.GetDollarMarketPrices(tokenSymbols)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	pairs, err := s.pairDao.GetDefaultPairs()
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	tradesQuery := []bson.M{
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
					"pairName":   "$pairName",
					"baseToken":  "$baseToken",
					"quoteToken": "$quoteToken",
				},
				"count":  bson.M{"$sum": one},
				"close":  bson.M{"$last": "$pricepoint"},
				"volume": bson.M{"$sum": bson.M{"$toDecimal": "$amount"}},
			},
		},
	}

	bidsQuery := []bson.M{
		bson.M{
			"$match": bson.M{
				"status": bson.M{"$in": []string{"OPEN", "PARTIAL_FILLED"}},
				"side":   "BUY",
			},
		},
		bson.M{
			"$group": bson.M{
				"_id": bson.M{
					"pairName":   "$pairName",
					"baseToken":  "$baseToken",
					"quoteToken": "$quoteToken",
				},
				"orderCount": bson.M{"$sum": one},
				"orderVolume": bson.M{
					"$sum": bson.M{
						"$subtract": []bson.M{bson.M{"$toDecimal": "$amount"}, bson.M{"$toDecimal": "$filledAmount"}},
					},
				},
				"bestPrice": bson.M{"$max": "$pricepoint"},
			},
		},
	}

	asksQuery := []bson.M{
		bson.M{
			"$match": bson.M{
				"status": bson.M{"$in": []string{"OPEN", "PARTIAL_FILLED"}},
				"side":   "SELL",
			},
		},
		bson.M{
			"$group": bson.M{
				"_id": bson.M{
					"pairName":   "$pairName",
					"baseToken":  "$baseToken",
					"quoteToken": "$quoteToken",
				},
				"orderCount": bson.M{"$sum": one},
				"orderVolume": bson.M{
					"$sum": bson.M{
						"$subtract": []bson.M{bson.M{"$toDecimal": "$amount"}, bson.M{"$toDecimal": "$filledAmount"}},
					},
				},
				"bestPrice": bson.M{"$min": "$pricepoint"},
			},
		},
	}

	tradeData, err := s.tradeDao.Aggregate(tradesQuery)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	bidsData, err := s.orderDao.Aggregate(bidsQuery)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	asksData, err := s.orderDao.Aggregate(asksQuery)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	var pairAPIData []*types.PairAPIData
	var totalOrders int
	var totalTrades int
	var totalVolume float64
	var totalBuyOrderAmount float64
	var totalSellOrderAmount float64
	var totalSellOrders int
	var totalBuyOrders int
	pairTradeCounts := map[string]int{}
	tokenTradeCounts := map[string]int{}

	// //total orderbook volume per quote token
	// totalOrderBookVolume := map[string]int{}

	for _, p := range pairs {
		pairData := &types.PairData{
			Pair:               types.PairID{p.Name(), p.BaseTokenAddress, p.QuoteTokenAddress},
			Volume:             big.NewInt(0),
			Close:              big.NewInt(0),
			Count:              big.NewInt(0),
			OrderVolume:        big.NewInt(0),
			OrderCount:         big.NewInt(0),
			BidPrice:           big.NewInt(0),
			AskPrice:           big.NewInt(0),
			Price:              big.NewInt(0),
			AverageOrderAmount: big.NewInt(0),
			AverageTradeAmount: big.NewInt(0),
		}

		for _, t := range tradeData {
			if t.AddressCode() == p.AddressCode() {
				pairData.Volume = t.Volume
				pairData.Close = t.Close
				pairData.Count = t.Count
				pairData.AverageTradeAmount = math.Div(t.Volume, t.Count)

				totalTrades = totalTrades + int(t.Count.Int64())
				totalVolume = totalVolume + t.ConvertedVolume(&p, rates[p.BaseTokenSymbol])

				pairTradeCounts[p.Name()] = int(t.Count.Int64())
				tokenTradeCounts[p.BaseTokenSymbol] = tokenTradeCounts[p.BaseTokenSymbol] + int(t.Count.Int64())
			}
		}

		for _, o := range bidsData {
			if o.AddressCode() == p.AddressCode() {
				pairData.OrderVolume = o.OrderVolume
				pairData.OrderCount = o.OrderCount
				pairData.BidPrice = o.BestPrice
				pairData.AverageOrderAmount = math.Div(pairData.OrderVolume, pairData.OrderCount)

				// change and replace by equivalent dollar volume instead of order count
				totalBuyOrderAmount = totalBuyOrderAmount + o.ConvertedVolume(&p, rates[p.BaseTokenSymbol])
				totalBuyOrders = totalBuyOrders + int(o.OrderCount.Int64())
			}
		}

		for _, o := range asksData {
			if o.AddressCode() == p.AddressCode() {
				pairData.OrderVolume = math.Add(pairData.OrderVolume, o.OrderVolume)
				pairData.OrderCount = math.Add(pairData.OrderCount, o.OrderCount)
				pairData.AskPrice = o.BestPrice
				pairData.AverageOrderAmount = math.Div(pairData.OrderVolume, pairData.OrderCount)

				// change and replace by equivalent dollar volume instead of order count
				totalSellOrderAmount = totalSellOrderAmount + o.ConvertedVolume(&p, rates[p.BaseTokenSymbol])
				totalSellOrders = totalSellOrders + int(o.OrderCount.Int64())

				//TODO change price into orderbook price
				if math.IsNotEqual(pairData.BidPrice, big.NewInt(0)) && math.IsNotEqual(pairData.AskPrice, big.NewInt(0)) {
					pairData.Price = math.Avg(pairData.BidPrice, pairData.AskPrice)
				} else {
					pairData.Price = big.NewInt(0)
				}
			}
		}

		pairAPIData = append(pairAPIData, pairData.ToAPIData(&p))
	}

	mostTradedToken, _ := utils.MaxIntMap(tokenTradeCounts)
	mostTradedPair, _ := utils.MaxIntMap(pairTradeCounts)
	totalOrders = totalSellOrders + totalBuyOrders
	totalOrderAmount := totalBuyOrderAmount + totalSellOrderAmount
	tradeSuccessRatio := float64(totalTrades-erroredTradeCount) / float64(totalTrades)

	exchangeData := &types.ExchangeData{
		PairData:             pairAPIData,
		TotalOrders:          totalOrders,
		TotalTrades:          totalTrades,
		TotalBuyOrderAmount:  totalBuyOrderAmount,
		TotalSellOrderAmount: totalSellOrderAmount,
		TotalVolume:          totalVolume,
		TotalOrderAmount:     totalOrderAmount,
		TotalBuyOrders:       totalSellOrders,
		TotalSellOrders:      totalBuyOrders,
		MostTradedToken:      mostTradedToken,
		MostTradedPair:       mostTradedPair,
		TradeSuccessRatio:    tradeSuccessRatio,
	}

	return exchangeData, nil
}

func (s *InfoService) GetPairStats() (*types.PairStats, error) {
	now := time.Now()
	end := time.Unix(now.Unix(), 0)
	start := time.Unix(now.AddDate(0, 0, -7).Unix(), 0)
	one, _ := bson.ParseDecimal128("1")

	tokens, err := s.tokenDao.GetBaseTokens()
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	quoteTokens, err := s.tokenDao.GetQuoteTokens()
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	tokenSymbols := []string{}
	for _, t := range tokens {
		tokenSymbols = append(tokenSymbols, t.Symbol)
	}

	for _, t := range quoteTokens {
		tokenSymbols = append(tokenSymbols, t.Symbol)
	}

	pairs, err := s.pairDao.GetDefaultPairs()
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	tradesQuery := []bson.M{
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
					"pairName":   "$pairName",
					"baseToken":  "$baseToken",
					"quoteToken": "$quoteToken",
				},
				"count":  bson.M{"$sum": one},
				"close":  bson.M{"$last": "$pricepoint"},
				"volume": bson.M{"$sum": bson.M{"$toDecimal": "$amount"}},
			},
		},
	}

	bidsQuery := []bson.M{
		bson.M{
			"$match": bson.M{
				"status": bson.M{"$in": []string{"OPEN", "PARTIAL_FILLED"}},
				"side":   "BUY",
			},
		},
		bson.M{
			"$group": bson.M{
				"_id": bson.M{
					"pairName":   "$pairName",
					"baseToken":  "$baseToken",
					"quoteToken": "$quoteToken",
				},
				"orderCount": bson.M{"$sum": one},
				"orderVolume": bson.M{
					"$sum": bson.M{
						"$subtract": []bson.M{bson.M{"$toDecimal": "$amount"}, bson.M{"$toDecimal": "$filledAmount"}},
					},
				},
				"bestPrice": bson.M{"$max": "$pricepoint"},
			},
		},
	}

	asksQuery := []bson.M{
		bson.M{
			"$match": bson.M{
				"status": bson.M{"$in": []string{"OPEN", "PARTIAL_FILLED"}},
				"side":   "SELL",
			},
		},
		bson.M{
			"$group": bson.M{
				"_id": bson.M{
					"pairName":   "$pairName",
					"baseToken":  "$baseToken",
					"quoteToken": "$quoteToken",
				},
				"orderCount": bson.M{"$sum": one},
				"orderVolume": bson.M{
					"$sum": bson.M{
						"$subtract": []bson.M{bson.M{"$toDecimal": "$amount"}, bson.M{"$toDecimal": "$filledAmount"}},
					},
				},
				"bestPrice": bson.M{"$min": "$pricepoint"},
			},
		},
	}

	tradeData, err := s.tradeDao.Aggregate(tradesQuery)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	bidsData, err := s.orderDao.Aggregate(bidsQuery)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	asksData, err := s.orderDao.Aggregate(asksQuery)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	var pairStatistics types.PairStats

	for _, p := range pairs {
		pairData := &types.PairData{
			Pair:               types.PairID{p.Name(), p.BaseTokenAddress, p.QuoteTokenAddress},
			Volume:             big.NewInt(0),
			Close:              big.NewInt(0),
			Count:              big.NewInt(0),
			OrderVolume:        big.NewInt(0),
			OrderCount:         big.NewInt(0),
			BidPrice:           big.NewInt(0),
			AskPrice:           big.NewInt(0),
			Price:              big.NewInt(0),
			AverageOrderAmount: big.NewInt(0),
			AverageTradeAmount: big.NewInt(0),
		}

		for _, t := range tradeData {
			if t.AddressCode() == p.AddressCode() {
				pairData.Volume = t.Volume
				pairData.Close = t.Close
				pairData.Count = t.Count
				pairData.AverageTradeAmount = math.Div(t.Volume, t.Count)
			}
		}

		for _, o := range bidsData {
			if o.AddressCode() == p.AddressCode() {
				pairData.OrderVolume = o.OrderVolume
				pairData.OrderCount = o.OrderCount
				pairData.BidPrice = o.BestPrice
				pairData.AverageOrderAmount = math.Div(pairData.OrderVolume, pairData.OrderCount)
			}
		}

		for _, o := range asksData {
			if o.AddressCode() == p.AddressCode() {
				pairData.OrderVolume = math.Add(pairData.OrderVolume, o.OrderVolume)
				pairData.OrderCount = math.Add(pairData.OrderCount, o.OrderCount)
				pairData.AskPrice = o.BestPrice
				pairData.AverageOrderAmount = math.Div(pairData.OrderVolume, pairData.OrderCount)

				//TODO change price into orderbook price
				if math.IsNotEqual(pairData.BidPrice, big.NewInt(0)) && math.IsNotEqual(pairData.AskPrice, big.NewInt(0)) {
					pairData.Price = math.Avg(pairData.BidPrice, pairData.AskPrice)
				} else {
					pairData.Price = big.NewInt(0)
				}
			}
		}

		pairStatistics = append(pairStatistics, pairData.ToAPIData(&p))
	}

	return &pairStatistics, nil
}
