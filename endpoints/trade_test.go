package endpoints

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Proofsuite/amp-matching-engine/types"
	"github.com/Proofsuite/amp-matching-engine/utils/testutils"
	"github.com/Proofsuite/amp-matching-engine/utils/testutils/mocks"
	"github.com/gorilla/mux"
)

func SetupTradeTest() (*mux.Router, *mocks.TradeService) {
	r := mux.NewRouter()
	tradeService := new(mocks.TradeService)

	ServeTradeResource(r, tradeService)

	return r, tradeService
}

func TestHandleGetTradeHistory(t *testing.T) {
	router, tradeService := SetupTest()

	t1 := testutils.GetTestZRXToken()
	t2 := testutils.GetTestWETHToken()

	tr1 := types.Trade{}
	tr2 := types.Trade{}
	trs := []types.Trade{tr1, tr2}

	tradeService.On("GetByPairAddress", t1.Address, t2.Address).Returns(trs)

	req, err := http.NewRequest("GET", "/trades/history/{")

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Handler return wrong status. Got %v want %v", rr.Code, http.StatusOK)
	}

	json.NewDecoder()

}
