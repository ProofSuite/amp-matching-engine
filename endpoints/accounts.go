package endpoints

import (
	"net/http"

	"github.com/Proofsuite/amp-matching-engine/interfaces"
	"github.com/Proofsuite/amp-matching-engine/types"
	"github.com/Proofsuite/amp-matching-engine/utils/httputils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/gorilla/mux"
)

type accountEndpoint struct {
	accountService interfaces.AccountService
}

func ServeAccountResource(
	r *mux.Router,
	accountService interfaces.AccountService,
) {

	e := &accountEndpoint{accountService}
	r.HandleFunc("/account/create", e.handleCreateAccount).Methods("POST")
	r.HandleFunc("/account/<address>", e.handleGetAccount).Methods("GET")
	r.HandleFunc("/account/{address}/{token}", e.handleGetAccountTokenBalance).Methods("GET")
}

func (e *accountEndpoint) handleCreateAccount(w http.ResponseWriter, r *http.Request) {
	v := r.URL.Query()
	addr := v.Get("address")

	if !common.IsHexAddress(addr) {
		httputils.WriteError(w, http.StatusBadRequest, "Invalid Address")
		return
	}

	a := &types.Account{Address: common.HexToAddress(addr)}
	err := e.accountService.Create(a)
	if err != nil {
		logger.Error(err)
		httputils.WriteError(w, http.StatusInternalServerError, "")
		return
	}

	httputils.WriteJSON(w, http.StatusCreated, a)
}

func (e *accountEndpoint) handleGetAccount(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	addr := vars["address"]
	if !common.IsHexAddress(addr) {
		httputils.WriteError(w, http.StatusBadRequest, "Invalid Address")
		return
	}

	address := common.HexToAddress(addr)
	a, err := e.accountService.GetByAddress(address)
	if err != nil {
		logger.Error(err)
		httputils.WriteError(w, http.StatusInternalServerError, "")
		return
	}

	httputils.WriteJSON(w, http.StatusOK, a)
}

func (e *accountEndpoint) handleGetAccountTokenBalance(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	a := vars["address"]
	if !common.IsHexAddress(a) {
		httputils.WriteError(w, http.StatusBadRequest, "Invalid Address")
	}

	t := vars["token"]
	if !common.IsHexAddress(a) {
		httputils.WriteError(w, http.StatusBadRequest, "Invalid Token Address")
	}

	addr := common.HexToAddress(a)
	tokenAddr := common.HexToAddress(t)

	b, err := e.accountService.GetTokenBalance(addr, tokenAddr)
	if err != nil {
		logger.Error(err)
		httputils.WriteError(w, http.StatusInternalServerError, "")
	}

	httputils.WriteJSON(w, http.StatusOK, b)
}
