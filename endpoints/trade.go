package endpoints

import (
	"encoding/json"
	"net/http"

	"github.com/Proofsuite/amp-matching-engine/interfaces"
	"github.com/Proofsuite/amp-matching-engine/types"
	"github.com/Proofsuite/amp-matching-engine/utils/httputils"
	"github.com/Proofsuite/amp-matching-engine/ws"
	"github.com/ethereum/go-ethereum/common"
	"github.com/gorilla/mux"
)

type tradeEndpoint struct {
	tradeService interfaces.TradeService
}

// ServeTradeResource sets up the routing of trade endpoints and the corresponding handlers.
// TODO trim down to one single endpoint with the 3 following params: base, quote, address
func ServeTradeResource(
	r *mux.Router,
	tradeService interfaces.TradeService,
) {
	e := &tradeEndpoint{tradeService}
	r.HandleFunc("/trades/pair", e.HandleGetTradeHistory)
	r.HandleFunc("/trades", e.HandleGetTrades)
	ws.RegisterChannel(ws.TradeChannel, e.tradeWebSocket)
}

// history is reponsible for handling pair's trade history requests
func (e *tradeEndpoint) HandleGetTradeHistory(w http.ResponseWriter, r *http.Request) {
	v := r.URL.Query()
	bt := v.Get("baseToken")
	qt := v.Get("quoteToken")

	if bt == "" {
		httputils.WriteError(w, http.StatusBadRequest, "baseToken Parameter missing")
		return
	}

	if qt == "" {
		httputils.WriteError(w, http.StatusBadRequest, "quoteToken Parameter missing")
		return
	}

	if !common.IsHexAddress(bt) {
		httputils.WriteError(w, http.StatusBadRequest, "Invalid base token address")
		return
	}

	if !common.IsHexAddress(qt) {
		httputils.WriteError(w, http.StatusBadRequest, "Invalid quote token address")
		return
	}

	baseToken := common.HexToAddress(bt)
	quoteToken := common.HexToAddress(qt)
	res, err := e.tradeService.GetByPairAddress(baseToken, quoteToken)
	if err != nil {
		logger.Error(err)
		httputils.WriteError(w, http.StatusInternalServerError, "")
		return
	}

	if res == nil {
		httputils.WriteJSON(w, http.StatusOK, []types.Trade{})
		return
	}

	httputils.WriteJSON(w, http.StatusOK, res)
}

// get is reponsible for handling user's trade history requests
func (e *tradeEndpoint) HandleGetTrades(w http.ResponseWriter, r *http.Request) {
	v := r.URL.Query()
	addr := v.Get("address")

	if addr == "" {
		httputils.WriteError(w, http.StatusBadRequest, "address Parameter missing")
		return
	}

	if !common.IsHexAddress(addr) {
		httputils.WriteError(w, http.StatusBadRequest, "Invalid Address")
		return
	}

	address := common.HexToAddress(addr)
	res, err := e.tradeService.GetByUserAddress(address)
	if err != nil {
		logger.Error(err)
		httputils.WriteError(w, http.StatusInternalServerError, "")
		return
	}

	if res == nil {
		httputils.WriteJSON(w, http.StatusOK, []types.Trade{})
		return
	}

	httputils.WriteJSON(w, http.StatusOK, res)
}

func (e *tradeEndpoint) tradeWebSocket(input interface{}, conn *ws.Conn) {
	bytes, _ := json.Marshal(input)
	var payload *types.WebSocketPayload
	if err := json.Unmarshal(bytes, &payload); err != nil {
		logger.Error(err)
	}

	socket := ws.GetTradeSocket()
	if payload.Type != "subscription" {
		err := map[string]string{"Message": "Invalid payload"}
		socket.SendErrorMessage(conn, err)
		return
	}

	bytes, _ = json.Marshal(payload.Data)
	var msg *types.WebSocketSubscription
	err := json.Unmarshal(bytes, &msg)
	if err != nil {
		logger.Error(err)
	}

	if (msg.Pair.BaseToken == common.Address{}) {
		err := map[string]string{"Message": "Invalid base token"}
		socket.SendErrorMessage(conn, err)
		return
	}

	if (msg.Pair.QuoteToken == common.Address{}) {
		err := map[string]string{"Message": "Invalid quote token"}
		socket.SendErrorMessage(conn, err)
		return
	}

	if msg.Event == types.SUBSCRIBE {
		e.tradeService.Subscribe(conn, msg.Pair.BaseToken, msg.Pair.QuoteToken)
	}

	if msg.Event == types.UNSUBSCRIBE {
		e.tradeService.Unsubscribe(conn, msg.Pair.BaseToken, msg.Pair.QuoteToken)
	}
}
