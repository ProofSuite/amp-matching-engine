package endpoints

import (
	"encoding/json"
	"net/http"
	"strconv"

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
	ws.RegisterChannel(ws.TradeChannel, e.tradeWebsocket)
}

// history is reponsible for handling pair's trade history requests
func (e *tradeEndpoint) HandleGetTradeHistory(w http.ResponseWriter, r *http.Request) {
	v := r.URL.Query()
	bt := v.Get("baseToken")
	qt := v.Get("quoteToken")
	l := v.Get("length")

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

	length := 20
	if l != "" {
		length, _ = strconv.Atoi(l)
	}

	baseToken := common.HexToAddress(bt)
	quoteToken := common.HexToAddress(qt)
	res, err := e.tradeService.GetSortedTradesByDate(baseToken, quoteToken, length)
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

func (e *tradeEndpoint) tradeWebsocket(input interface{}, c *ws.Conn) {
	b, _ := json.Marshal(input)
	var ev *types.WebsocketEvent
	if err := json.Unmarshal(b, &ev); err != nil {
		logger.Error(err)
		return
	}

	socket := ws.GetTradeSocket()
	if ev.Type != "SUBSCRIBE" && ev.Type != "UNSUBSCRIBE" {
		logger.Info("Event Type", ev.Type)
		err := map[string]string{"Message": "Invalid payload"}
		socket.SendErrorMessage(c, err)
		return
	}

	b, _ = json.Marshal(ev.Payload)
	var p *types.SubscriptionPayload
	err := json.Unmarshal(b, &p)
	if err != nil {
		logger.Error(err)
		return
	}

	if ev.Type == "SUBSCRIBE" {
		if (p.BaseToken == common.Address{}) {
			err := map[string]string{"Message": "Invalid base token"}
			socket.SendErrorMessage(c, err)
			return
		}

		if (p.QuoteToken == common.Address{}) {
			err := map[string]string{"Message": "Invalid quote token"}
			socket.SendErrorMessage(c, err)
			return
		}

		e.tradeService.Subscribe(c, p.BaseToken, p.QuoteToken)
	}

	if ev.Type == "UNSUBSCRIBE" {
		if p == nil {
			e.tradeService.Unsubscribe(c)
			return
		}

		e.tradeService.UnsubscribeChannel(c, p.BaseToken, p.QuoteToken)
	}
}
