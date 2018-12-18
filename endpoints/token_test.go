package endpoints

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Proofsuite/amp-matching-engine/types"
	"github.com/Proofsuite/amp-matching-engine/utils/testutils"
	"github.com/Proofsuite/amp-matching-engine/utils/testutils/mocks"
	"github.com/ethereum/go-ethereum/common"
	"github.com/gorilla/mux"
)

func SetupTest() (*mux.Router, *mocks.TokenService) {
	r := mux.NewRouter()
	tokenService := new(mocks.TokenService)

	ServeTokenResource(r, tokenService)

	return r, tokenService
}

func TestHandleCreateTokens(t *testing.T) {
	router, tokenService := SetupTest()

	token := types.Token{
		Name:            "ZRX",
		Symbol:          "ZRX",
		Decimal:         18,
		Quote:           false,
		Address: common.HexToAddress("0x1"),
	}

	tokenService.On("Create", &token).Return(nil)

	b, _ := json.Marshal(token)
	req, err := http.NewRequest("POST", "/tokens", bytes.NewBuffer(b))
	if err != nil {
		t.Error(err)
	}

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusCreated {
		t.Errorf("Handler return wrong status. Got %v want %v", rr.Code, http.StatusCreated)
	}

	created := types.Token{}
	json.NewDecoder(rr.Body).Decode(&created)

	tokenService.AssertCalled(t, "Create", &token)
	testutils.CompareToken(t, &token, &created)
}

func TestHandleGetTokens(t *testing.T) {
	router, tokenService := SetupTest()

	t1 := testutils.GetTestZRXToken()
	t2 := testutils.GetTestWETHToken()

	tokenService.On("GetAll").Return([]types.Token{t1, t2}, nil)

	req, err := http.NewRequest("GET", "/tokens", nil)
	if err != nil {
		t.Error(err)
	}

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Handler return wrong status. Got %v want %v", rr.Code, http.StatusOK)
	}

	result := []types.Token{}
	json.NewDecoder(rr.Body).Decode(&result)

	tokenService.AssertCalled(t, "GetAll")
	testutils.CompareToken(t, &t1, &result[0])
	testutils.CompareToken(t, &t2, &result[1])
}

func TestHandleGetQuoteTokens(t *testing.T) {
	router, tokenService := SetupTest()

	t1 := types.Token{
		Name:            "WETH",
		Symbol:          "WETH",
		Decimal:         18,
		Quote:           true,
		Address: common.HexToAddress("0x1"),
	}

	t2 := types.Token{
		Name:            "DAI",
		Symbol:          "DAI",
		Decimal:         18,
		Quote:           true,
		Address: common.HexToAddress("0x2"),
	}

	tokenService.On("GetQuoteTokens").Return([]types.Token{t1, t2}, nil)

	req, err := http.NewRequest("GET", "/tokens/quote", nil)
	if err != nil {
		t.Error(err)
	}

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Handler return wrong status. Got %v want %v", rr.Code, http.StatusOK)
	}

	result := []types.Token{}
	json.NewDecoder(rr.Body).Decode(&result)

	tokenService.AssertCalled(t, "GetQuoteTokens")
	testutils.CompareToken(t, &t1, &result[0])
	testutils.CompareToken(t, &t2, &result[1])
}

func TestHandleGetBaseTokens(t *testing.T) {
	router, tokenService := SetupTest()

	t1 := types.Token{
		Name:            "WETH",
		Symbol:          "WETH",
		Decimal:         18,
		Quote:           false,
		Address: common.HexToAddress("0x1"),
	}

	t2 := types.Token{
		Name:            "DAI",
		Symbol:          "DAI",
		Decimal:         18,
		Quote:           false,
		Address: common.HexToAddress("0x2"),
	}

	tokenService.On("GetBaseTokens").Return([]types.Token{t1, t2}, nil)

	req, err := http.NewRequest("GET", "/tokens/base", nil)
	if err != nil {
		t.Error(err)
	}

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Handler return wrong status. Got %v want %v", rr.Code, http.StatusOK)
	}

	result := []types.Token{}
	json.NewDecoder(rr.Body).Decode(&result)

	tokenService.AssertCalled(t, "GetBaseTokens")
	testutils.CompareToken(t, &t1, &result[0])
	testutils.CompareToken(t, &t2, &result[1])
}

func TestHandleGetToken(t *testing.T) {
	router, tokenService := SetupTest()

	addr := common.HexToAddress("0x1")

	t1 := types.Token{
		Name:            "DAI",
		Symbol:          "DAI",
		Decimal:         18,
		Quote:           false,
		Address: addr,
	}

	tokenService.On("GetByAddress", addr).Return(&t1, nil)

	url := "/tokens/" + addr.Hex()
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		t.Error(err)
	}

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Handler return wrong status. Got %v want %v", rr.Code, http.StatusOK)
	}

	result := types.Token{}
	json.NewDecoder(rr.Body).Decode(&result)

	tokenService.AssertCalled(t, "GetByAddress", addr)
	testutils.Compare(t, &t1, &result)
}

// func TestHandleGetTokens(t *testing.T) {
// 	router, tokenService := SetupTest()

// }

// func TestHandleGetQuoteTokens(t *testing.T) {
// 	router, tokenService := SetupTest()

// }

// func TestHandleGetBaseTokens(t *testing.T) {
// 	router, tokenService := SetupTest()
// }

// func TestHandleGetToken(t *testing.T) {
// 	router, tokenService := SetupTest()

// }

// var resp interface{}
// 			if err := json.Unmarshal(res.Body.Bytes(), &resp); err != nil {
// 				fmt.Printf("%v", err)
// 			}

// 			bytes.NewBufferString
