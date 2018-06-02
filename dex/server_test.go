package dex

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ethereum/go-ethereum/common"
)

func TestNewQuoteHandler(t *testing.T) {

	q := Token{Symbol: "WETH", Address: common.HexToAddress("0x5d564669ab4cfd96b785d3d05e8c7d66a073daf0")}
	b := new(bytes.Buffer)

	encoder := json.NewEncoder(b)
	encoder.Encode(q)

	s := NewServer()

	req, err := http.NewRequest("POST", "/quote/new", b)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(s.RegisterNewQuoteToken)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned the wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}

func TestNewPairHandler(t *testing.T) {
	WETH := Token{Symbol: "WETH", Address: common.HexToAddress("0x5d564669ab4cfd96b785d3d05e8c7d66a073daf0")}
	ZRX := Token{Symbol: "ZRX", Address: common.HexToAddress("0x9792845456a0075df8a03123e7dac62bb0f69440")}
	p := NewPair(ZRX, WETH)
	b := new(bytes.Buffer)

	encoder := json.NewEncoder(b)
	encoder.Encode(p)

	s := NewServer()
	s.SetupCurrencies(Tokens{"WETH": WETH}, nil, nil)

	req, err := http.NewRequest("POST", "/pair/new", b)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(s.RegisterNewPair)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned the wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}
