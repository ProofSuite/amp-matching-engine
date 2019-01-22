package endpoints

import (
	"net/http"

	"github.com/Proofsuite/amp-matching-engine/app"
	"github.com/Proofsuite/amp-matching-engine/interfaces"
	"github.com/Proofsuite/amp-matching-engine/types"
	"github.com/Proofsuite/amp-matching-engine/utils/httputils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/gorilla/mux"
)

type infoEndpoint struct {
	walletService interfaces.WalletService
	tokenService  interfaces.TokenService
	infoService   interfaces.InfoService
}

func ServeInfoResource(
	r *mux.Router,
	walletService interfaces.WalletService,
	tokenService interfaces.TokenService,
	infoService interfaces.InfoService,
) {

	e := &infoEndpoint{walletService, tokenService, infoService}
	r.HandleFunc("/info", e.handleGetInfo)
	r.HandleFunc("/info/exchange", e.handleGetExchangeInfo)
	r.HandleFunc("/info/operators", e.handleGetOperatorsInfo)
	r.HandleFunc("/info/fees", e.handleGetFeeInfo)
	r.HandleFunc("/stats/trading", e.handleGetTradingStats)
	// r.HandleFunc("/stats/all", e.handleGetStats)
	// r.HandleFunc("/stats/pairs", e.handleGetPairStats)
}

func (e *infoEndpoint) handleGetInfo(w http.ResponseWriter, r *http.Request) {
	ex := common.HexToAddress(app.Config.Ethereum["exchange_address"])

	operators, err := e.walletService.GetOperatorAddresses()
	if err != nil {
		logger.Error(err)
		httputils.WriteJSON(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	quotes, err := e.tokenService.GetQuoteTokens()
	if err != nil {
		logger.Error(err)
	}

	fees := []map[string]string{}
	for _, q := range quotes {
		fees = append(fees, map[string]string{
			"quote":   q.Symbol,
			"makeFee": q.MakeFee.String(),
			"takeFee": q.TakeFee.String(),
		})
	}

	res := map[string]interface{}{
		"exchangeAddress": ex.Hex(),
		"fees":            fees,
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
	quotes, err := e.tokenService.GetQuoteTokens()
	if err != nil {
		logger.Error(err)
	}

	fees := []map[string]string{}
	for _, q := range quotes {
		fees = append(fees, map[string]string{
			"quote":   q.Symbol,
			"makeFee": q.MakeFee.String(),
			"takeFee": q.TakeFee.String(),
		})
	}

	httputils.WriteJSON(w, http.StatusOK, fees)
}

func (e *infoEndpoint) handleGetTradingStats(w http.ResponseWriter, r *http.Request) {
	res, err := e.infoService.GetExchangeStats()
	if err != nil {
		logger.Error(err)
		httputils.WriteError(w, http.StatusInternalServerError, "")
		return
	}

	if res == nil {
		httputils.WriteJSON(w, http.StatusOK, []types.Pair{})
		return
	}

	httputils.WriteJSON(w, http.StatusOK, res)
}

func (e *infoEndpoint) handleGetPairStats(w http.ResponseWriter, r *http.Request) {
	res, err := e.infoService.GetPairStats()
	if err != nil {
		logger.Error(err)
		httputils.WriteError(w, http.StatusInternalServerError, "")
		return
	}

	if res == nil {
		httputils.WriteJSON(w, http.StatusOK, []types.Pair{})
		return
	}

	httputils.WriteJSON(w, http.StatusOK, res)
}

func (e *infoEndpoint) handleGetStats(w http.ResponseWriter, r *http.Request) {
	res, err := e.infoService.GetExchangeData()
	if err != nil {
		logger.Error(err)
		httputils.WriteError(w, http.StatusInternalServerError, "")
		return
	}

	if res == nil {
		httputils.WriteJSON(w, http.StatusOK, []types.Pair{})
		return
	}

	httputils.WriteJSON(w, http.StatusOK, res)
}
