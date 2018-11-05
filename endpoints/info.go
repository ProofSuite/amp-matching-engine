package endpoints

import (
	"net/http"

	"github.com/Proofsuite/amp-matching-engine/app"
	"github.com/Proofsuite/amp-matching-engine/interfaces"
	"github.com/Proofsuite/amp-matching-engine/utils/httputils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/gorilla/mux"
)

type infoEndpoint struct {
	walletService interfaces.WalletService
}

func ServeInfoResource(
	r *mux.Router,
	walletService interfaces.WalletService,
) {

	e := &infoEndpoint{walletService}
	r.HandleFunc("/info", e.handleGetInfo)
	r.HandleFunc("/info/exchange", e.handleGetExchangeInfo)
	r.HandleFunc("/info/operators", e.handleGetOperatorsInfo)
	r.HandleFunc("/info/fees", e.handleGetFeeInfo)
}

func (e *infoEndpoint) handleGetInfo(w http.ResponseWriter, r *http.Request) {
	ex := common.HexToAddress(app.Config.Ethereum["exchange_address"])

	operators, err := e.walletService.GetOperatorAddresses()
	if err != nil {
		logger.Error(err)
		httputils.WriteJSON(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	res := map[string]interface{}{
		"exchangeAddress": ex.Hex(),
		"makeFee":         "0",
		"takeFee":         "0",
		"operators":       operators,
	}

	httputils.WriteJSON(w, http.StatusOK, res)
}

func (e *infoEndpoint) handleGetExchangeInfo(w http.ResponseWriter, r *http.Request) {
	ex := common.HexToAddress(app.Config.Ethereum["exchange_address"])

	res := map[string]string{"exchangeAddress": ex.Hex()}

	httputils.WriteJSON(w, http.StatusOK, res)
}

func (e *infoEndpoint) handleGetOperatorsInfo(w http.ResponseWriter, r *http.Request) {
	addresses, err := e.walletService.GetOperatorAddresses()
	if err != nil {
		logger.Error(err)
		httputils.WriteJSON(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	res := map[string][]common.Address{"operators": addresses}

	httputils.WriteJSON(w, http.StatusOK, res)
}

func (e *infoEndpoint) handleGetFeeInfo(w http.ResponseWriter, r *http.Request) {
	fees := map[string]string{
		"makeFee": "0",
		"takeFee": "0",
	}

	httputils.WriteJSON(w, http.StatusOK, fees)
}
