package endpoints

import (
	"encoding/json"
	"net/http"

	"github.com/Proofsuite/amp-matching-engine/interfaces"
	"github.com/Proofsuite/amp-matching-engine/services"
	"github.com/Proofsuite/amp-matching-engine/types"
	"github.com/Proofsuite/amp-matching-engine/utils"
	"github.com/Proofsuite/amp-matching-engine/utils/httputils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/gorilla/mux"
)

type tokenEndpoint struct {
	tokenService interfaces.TokenService
}

// ServeTokenResource sets up the routing of token endpoints and the corresponding handlers.
func ServeTokenResource(
	r *mux.Router,
	tokenService interfaces.TokenService,
) {
	e := &tokenEndpoint{tokenService}
	r.HandleFunc("/tokens/base", e.HandleGetBaseTokens).Methods("GET")
	r.HandleFunc("/tokens/quote", e.HandleGetQuoteTokens).Methods("GET")
	r.HandleFunc("/tokens/{address}", e.HandleGetToken).Methods("GET")
	r.HandleFunc("/tokens", e.HandleGetTokens).Methods("GET")
	r.HandleFunc("/tokens", e.HandleCreateTokens).Methods("POST")
}

func (e *tokenEndpoint) HandleCreateTokens(w http.ResponseWriter, r *http.Request) {
	var t types.Token
	decoder := json.NewDecoder(r.Body)

	err := decoder.Decode(&t)
	if err != nil {
		logger.Error(err)
		httputils.WriteError(w, http.StatusBadRequest, "Invalid payload")
	}

	utils.PrintJSON(t)

	defer r.Body.Close()

	err = e.tokenService.Create(&t)
	if err != nil {
		if err == services.ErrTokenExists {
			httputils.WriteError(w, http.StatusBadRequest, "")
			return
		} else {
			logger.Error(err)
			httputils.WriteError(w, http.StatusInternalServerError, "")
			return
		}
	}

	httputils.WriteJSON(w, http.StatusCreated, t)
}

func (e *tokenEndpoint) HandleGetTokens(w http.ResponseWriter, r *http.Request) {
	res, err := e.tokenService.GetAll()
	if err != nil {
		logger.Error(err)
		httputils.WriteError(w, http.StatusInternalServerError, "")
		return
	}

	httputils.WriteJSON(w, http.StatusOK, res)
}

func (e *tokenEndpoint) HandleGetQuoteTokens(w http.ResponseWriter, r *http.Request) {
	res, err := e.tokenService.GetQuoteTokens()
	if err != nil {
		logger.Error(err)
		httputils.WriteError(w, http.StatusInternalServerError, "")
	}

	httputils.WriteJSON(w, http.StatusOK, res)
}

func (e *tokenEndpoint) HandleGetBaseTokens(w http.ResponseWriter, r *http.Request) {
	res, err := e.tokenService.GetBaseTokens()
	if err != nil {
		logger.Error(err)
		httputils.WriteError(w, http.StatusInternalServerError, "")
		return
	}

	httputils.WriteJSON(w, http.StatusOK, res)
}

func (e *tokenEndpoint) HandleGetToken(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	a := vars["address"]
	if !common.IsHexAddress(a) {
		httputils.WriteError(w, http.StatusBadRequest, "Invalid Address")
	}

	tokenAddress := common.HexToAddress(a)
	res, err := e.tokenService.GetByAddress(tokenAddress)
	if err != nil {
		logger.Error(err)
		httputils.WriteError(w, http.StatusInternalServerError, "")
		return
	}

	httputils.WriteJSON(w, http.StatusOK, res)
}
