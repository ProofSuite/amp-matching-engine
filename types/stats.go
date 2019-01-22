package types

type ExchangeData struct {
	PairData             []*PairAPIData `json:"pairData"`
	TotalOrders          int            `json:"totalOrders"`
	TotalTrades          int            `json:"totalTrades"`
	TotalSellOrders      int            `json:"totalSellOrders"`
	TotalBuyOrders       int            `json:"totalBuyOrders"`
	TotalBuyOrderAmount  float64        `json:"totalBuyAmount"`
	TotalSellOrderAmount float64        `json:"totalSellAmount"`
	TotalVolume          float64        `json:"totalVolume"`
	TotalOrderAmount     float64        `json:"totalOrderAmount"`
	MostTradedToken      string         `json:"mostTradedToken"`
	MostTradedPair       string         `json:"mostTradedPair"`
	TradeSuccessRatio    float64        `json:"tradeSuccessRatio"`
}

type ExchangeStats struct {
	TotalOrders          int     `json:"totalOrders"`
	TotalTrades          int     `json:"totalTrades"`
	TotalSellOrders      int     `json:"totalSellOrders"`
	TotalBuyOrders       int     `json:"totalBuyOrders"`
	TotalBuyOrderAmount  float64 `json:"totalBuyAmount"`
	TotalSellOrderAmount float64 `json:"totalSellAmount"`
	TotalVolume          float64 `json:"totalVolume"`
	TotalOrderAmount     float64 `json:"totalOrderAmount"`
	MostTradedToken      string  `json:"mostTradedToken"`
	MostTradedPair       string  `json:"mostTradedPair"`
	TradeSuccessRatio    float64 `json:"tradeSuccessRatio"`
}

type PairStats []*PairAPIData
