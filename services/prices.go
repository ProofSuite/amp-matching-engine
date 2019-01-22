package services

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
)

type PriceService struct{}

func NewPriceService() *PriceService {
	return &PriceService{}
}

func (s *PriceService) GetDollarMarketPrices(baseCurrencies []string) (map[string]float64, error) {
	bases := strings.Join(baseCurrencies[:], ",")
	url := "https://min-api.cryptocompare.com/data/price?fsym=USD&tsyms=" + bases

	res, err := http.Get(url)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	result := map[string]float64{}

	err = json.Unmarshal(b, &result)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	return result, nil
}

func (s *PriceService) GetMultipleMarketPrices(baseCurrencies []string, quoteCurrencies []string) (map[string]map[string]float64, error) {
	base := strings.Join(baseCurrencies[:], ",")
	quotes := strings.Join(quoteCurrencies[:], ",")

	url := "https://min-api.cryptocompare.com/data/pricemulti?fsyms=" + base + "&tsyms=" + quotes

	res, err := http.Get(url)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	// the API result
	result := map[string]map[string]float64{}
	err = json.Unmarshal(b, &result)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	return result, nil
}
