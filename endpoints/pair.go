package endpoints

import (
	"encoding/json"
	"net/http"

	"github.com/ethereum/go-ethereum/common"

	"github.com/Proofsuite/amp-matching-engine/interfaces"
	"github.com/Proofsuite/amp-matching-engine/services"
	"github.com/Proofsuite/amp-matching-engine/types"
	"github.com/Proofsuite/amp-matching-engine/utils/httputils"
	"github.com/gorilla/mux"
)

type pairEndpoint struct {
	pairService interfaces.PairService
}

// ServePairResource sets up the routing of pair endpoints and the corresponding handlers.
func ServePairResource(
	r *mux.Router,
	p interfaces.PairService,
) {
	e := &pairEndpoint{p}
	r.HandleFunc("/pairs", e.HandleCreatePair).Methods("POST")
	r.HandleFunc("/pairs/{baseToken}/{quoteToken}", e.HandleGetPair).Methods("GET")
	r.HandleFunc("/pairs", e.HandleGetAllPairs).Methods("GET")
}

func (e *pairEndpoint) HandleCreatePair(w http.ResponseWriter, r *http.Request) {
	p := &types.Pair{}

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(p)
	if err != nil {
		httputils.WriteError(w, http.StatusBadRequest, "Invalid payload")
		return
	}

	defer r.Body.Close()

	err = p.Validate()
	if err != nil {
		httputils.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	err = e.pairService.Create(p)
	if err != nil {
		switch err {
		case services.ErrPairExists:
			httputils.WriteError(w, http.StatusBadRequest, "Pair exists")
			return
		case services.ErrBaseTokenNotFound:
			httputils.WriteError(w, http.StatusBadRequest, "Base token not found")
			return
		case services.ErrQuoteTokenNotFound:
			httputils.WriteError(w, http.StatusBadRequest, "Quote token not found")
			return
		case services.ErrQuoteTokenInvalid:
			httputils.WriteError(w, http.StatusBadRequest, "Quote token invalid (token is not registered as quote")
			return
		default:
			logger.Error(err)
			httputils.WriteError(w, http.StatusInternalServerError, "")
			return
		}
	}

	httputils.WriteJSON(w, http.StatusCreated, p)
}

func (e *pairEndpoint) HandleGetAllPairs(w http.ResponseWriter, r *http.Request) {
	res, err := e.pairService.GetAll()
	if err != nil {
		logger.Error(err)
		httputils.WriteError(w, http.StatusInternalServerError, "")
		return
	}

	httputils.WriteJSON(w, http.StatusOK, res)
}

func (e *pairEndpoint) HandleGetPair(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	baseToken := vars["baseToken"]
	quoteToken := vars["quoteToken"]

	if !common.IsHexAddress(baseToken) {
		httputils.WriteError(w, http.StatusBadRequest, "Invalid Address")
	}

	if !common.IsHexAddress(quoteToken) {
		httputils.WriteError(w, http.StatusBadRequest, "Invalid Address")
	}

	baseTokenAddress := common.HexToAddress(baseToken)
	quoteTokenAddress := common.HexToAddress(quoteToken)
	res, err := e.pairService.GetByTokenAddress(baseTokenAddress, quoteTokenAddress)
	if err != nil {
		logger.Error(err)
		httputils.WriteError(w, http.StatusInternalServerError, "")
		return
	}

	httputils.WriteJSON(w, http.StatusOK, res)
}
