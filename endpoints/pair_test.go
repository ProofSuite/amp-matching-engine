package endpoints

import (
	"bytes"
	"encoding/json"
	"math/big"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Proofsuite/amp-matching-engine/types"
	"github.com/Proofsuite/amp-matching-engine/utils/testutils"
	"github.com/Proofsuite/amp-matching-engine/utils/testutils/mocks"
	"github.com/ethereum/go-ethereum/common"
	"github.com/gorilla/mux"
)

func SetupPairEndpointTest() (*mux.Router, *mocks.PairService) {
	r := mux.NewRouter()
	pairService := new(mocks.PairService)

	ServePairResource(r, pairService)

	return r, pairService
}

func TestHandleCreatePair(t *testing.T) {
	router, pairService := SetupPairEndpointTest()

	pair := types.Pair{
		BaseTokenSymbol:   "ZRX",
		BaseTokenAddress:  common.HexToAddress("0x1"),
		QuoteTokenSymbol:  "WETH",
		QuoteTokenAddress: common.HexToAddress("0x2"),
		MakeFee:           big.NewInt(1e4),
		TakeFee:           big.NewInt(1e4),
	}

	pairService.On("Create", &pair).Return(nil)

	b, _ := json.Marshal(pair)
	req, err := http.NewRequest("POST", "/pairs", bytes.NewBuffer(b))
	if err != nil {
		t.Error(err)
	}

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusCreated {
		t.Errorf("Handler return wrong status. Got %v want %v", rr.Code, http.StatusCreated)
	}

	created := types.Pair{}
	json.NewDecoder(rr.Body).Decode(&created)

	pairService.AssertCalled(t, "Create", &pair)
	testutils.ComparePair(t, &pair, &created)
}

func TestHandleCreateInvalidPair(t *testing.T) {
	router, pairService := SetupPairEndpointTest()

	pair := types.Pair{
		BaseTokenSymbol:   "ZRX",
		BaseTokenAddress:  common.HexToAddress("0x1"),
		QuoteTokenAddress: common.HexToAddress("0x2"),
		MakeFee:           big.NewInt(1e4),
		TakeFee:           big.NewInt(1e4),
	}

	pairService.On("Create", &pair).Return(nil)

	b, _ := json.Marshal(pair)
	req, err := http.NewRequest("POST", "/pairs", bytes.NewBuffer(b))
	if err != nil {
		t.Error(err)
	}

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("Handler return wrong status. Got %v want %v", rr.Code, http.StatusBadRequest)
	}
}

func TestHandleGetAllPairs(t *testing.T) {
	router, pairService := SetupPairEndpointTest()

	p1 := types.Pair{
		BaseTokenSymbol:   "ZRX",
		BaseTokenAddress:  common.HexToAddress("0x1"),
		QuoteTokenAddress: common.HexToAddress("0x2"),
		MakeFee:           big.NewInt(1e4),
		TakeFee:           big.NewInt(1e4),
	}

	p2 := types.Pair{
		BaseTokenSymbol:   "WETH",
		BaseTokenAddress:  common.HexToAddress("0x3"),
		QuoteTokenAddress: common.HexToAddress("0x4"),
		MakeFee:           big.NewInt(1e4),
		TakeFee:           big.NewInt(1e4),
	}

	pairService.On("GetAll").Return([]types.Pair{p1, p2}, nil)

	req, err := http.NewRequest("GET", "/pairs", nil)
	if err != nil {
		t.Error(err)
	}

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Handler return wrong status. Got %v want %v", rr.Code, http.StatusOK)
	}

	result := []types.Pair{}
	json.NewDecoder(rr.Body).Decode(&result)

	pairService.AssertCalled(t, "GetAll")
	testutils.ComparePair(t, &p1, &result[0])
	testutils.ComparePair(t, &p2, &result[1])
}

func TestHandleGetPair(t *testing.T) {
	router, pairService := SetupPairEndpointTest()

	base := common.HexToAddress("0x1")
	quote := common.HexToAddress("0x2")

	p1 := types.Pair{
		BaseTokenSymbol:   "ZRX",
		QuoteTokenSymbol:  "WETH",
		BaseTokenAddress:  base,
		QuoteTokenAddress: quote,
		MakeFee:           big.NewInt(1e4),
		TakeFee:           big.NewInt(1e4),
	}

	pairService.On("GetByTokenAddress", base, quote).Return(&p1, nil)

	url := "/pairs/" + base.Hex() + "/" + quote.Hex()
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		t.Error(err)
	}

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Handler return wrong status. Got %v want %v", rr.Code, http.StatusOK)
	}

	result := types.Pair{}
	json.NewDecoder(rr.Body).Decode(&result)

	pairService.AssertCalled(t, "GetByTokenAddress", base, quote)
	testutils.ComparePair(t, &p1, &result)
}
