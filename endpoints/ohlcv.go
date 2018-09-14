package endpoints

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/Proofsuite/amp-matching-engine/interfaces"
	"github.com/Proofsuite/amp-matching-engine/types"
	"github.com/Proofsuite/amp-matching-engine/utils/httputils"
	"github.com/Proofsuite/amp-matching-engine/ws"
	"github.com/ethereum/go-ethereum/common"
	"github.com/gorilla/mux"
)

type OHLCVEndpoint struct {
	ohlcvService interfaces.OHLCVService
}

func ServeOHLCVResource(
	r *mux.Router,
	ohlcvService interfaces.OHLCVService,
) {
	e := &OHLCVEndpoint{ohlcvService}
	r.HandleFunc("/ohlcv", e.handleGetOHLCV).Methods("POST")
	ws.RegisterChannel(ws.OHLCVChannel, e.ohlcvWebSocket)
}

func (e *OHLCVEndpoint) handleGetOHLCV(w http.ResponseWriter, r *http.Request) {
	var model types.TickRequest

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(model)
	if err != nil {
		logger.Error(err)
		httputils.WriteError(w, http.StatusBadRequest, "Invalid payload")
		return
	}

	if model.Units == "" {
		model.Units = "hour"
	}

	if model.Duration == 0 {
		model.Duration = 24
	}

	if model.To == 0 {
		model.To = time.Now().Unix()
	}

	res, err := e.ohlcvService.GetOHLCV(model.Pair, model.Duration, model.Units, model.From, model.To)
	if err != nil {
		logger.Error(err)
		httputils.WriteError(w, http.StatusInternalServerError, "")
		return
	}

	httputils.WriteJSON(w, http.StatusOK, res)
}

func (e *OHLCVEndpoint) ohlcvWebSocket(input interface{}, conn *ws.Conn) {
	startTs := time.Date(1970, time.January, 1, 0, 0, 0, 0, time.UTC)

	mab, _ := json.Marshal(input)
	var payload *types.WebSocketPayload

	err := json.Unmarshal(mab, &payload)
	if err != nil {
		logger.Error(err)
	}

	socket := ws.GetOHLCVSocket()

	if payload.Type != "subscription" {
		socket.SendErrorMessage(conn, "Invalid payload")
		return
	}

	dab, _ := json.Marshal(payload.Data)
	var msg *types.WebSocketSubscription

	err = json.Unmarshal(dab, &msg)
	if err != nil {
		logger.Error(err)
	}

	if (msg.Pair.BaseToken == common.Address{}) {
		socket.SendErrorMessage(conn, "Invalid base token")
		return
	}

	if (msg.Pair.QuoteToken == common.Address{}) {
		socket.SendErrorMessage(conn, "Invalid Quote Token")
		return
	}

	if msg.Params.From == 0 {
		msg.Params.From = startTs.Unix()
	}

	if msg.Params.To == 0 {
		msg.Params.To = time.Now().Unix()
	}

	if msg.Params.Duration == 0 {
		msg.Params.Duration = 24
	}

	if msg.Params.Units == "" {
		msg.Params.Units = "hour"
	}

	if msg.Event == types.SUBSCRIBE {
		e.ohlcvService.Subscribe(conn, msg.Pair.BaseToken, msg.Pair.QuoteToken, &msg.Params)
	}

	if msg.Event == types.UNSUBSCRIBE {
		e.ohlcvService.Unsubscribe(conn, msg.Pair.BaseToken, msg.Pair.QuoteToken, &msg.Params)
	}
}
