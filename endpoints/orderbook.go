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

type OrderBookEndpoint struct {
	orderBookService interfaces.OrderBookService
}

// ServePairResource sets up the routing of pair endpoints and the corresponding handlers.
func ServeOrderBookResource(
	r *mux.Router,
	orderBookService interfaces.OrderBookService,
) {
	e := &OrderBookEndpoint{orderBookService}
	r.HandleFunc("/orderbook/raw", e.handleGetRawOrderBook)
	r.HandleFunc("/orderbook", e.handleGetOrderBook)
	ws.RegisterChannel(ws.LiteOrderBookChannel, e.orderBookWebSocket)
	ws.RegisterChannel(ws.RawOrderBookChannel, e.rawOrderBookWebSocket)
}

// orderBookEndpoint
func (e *OrderBookEndpoint) handleGetOrderBook(w http.ResponseWriter, r *http.Request) {
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
		httputils.WriteError(w, http.StatusBadRequest, "Invalid Base Token Address")
		return
	}

	if !common.IsHexAddress(qt) {
		httputils.WriteError(w, http.StatusBadRequest, "Invalid Quote Token Address")
		return
	}

	baseTokenAddress := common.HexToAddress(bt)
	quoteTokenAddress := common.HexToAddress(qt)
	ob, err := e.orderBookService.GetOrderBook(baseTokenAddress, quoteTokenAddress)
	if err != nil {
		logger.Error(err)
		httputils.WriteError(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	httputils.WriteJSON(w, http.StatusOK, ob)
}

// orderBookEndpoint
func (e *OrderBookEndpoint) handleGetRawOrderBook(w http.ResponseWriter, r *http.Request) {
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
		httputils.WriteError(w, http.StatusBadRequest, "Invalid Base Token Address")
		return
	}

	if !common.IsHexAddress(qt) {
		httputils.WriteError(w, http.StatusBadRequest, "Invalid Quote Token Address")
		return
	}

	baseTokenAddress := common.HexToAddress(bt)
	quoteTokenAddress := common.HexToAddress(qt)
	ob, err := e.orderBookService.GetRawOrderBook(baseTokenAddress, quoteTokenAddress)
	if err != nil {
		httputils.WriteError(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	httputils.WriteJSON(w, http.StatusOK, ob)
}

// liteOrderBookWebSocket
func (e *OrderBookEndpoint) rawOrderBookWebSocket(input interface{}, conn *ws.Conn) {
	mab, _ := json.Marshal(input)
	var payload *types.WebSocketPayload

	err := json.Unmarshal(mab, &payload)
	if err != nil {
		logger.Error(err)
		return
	}

	socket := ws.GetRawOrderBookSocket()

	if payload.Type != "subscription" {
		logger.Error("Payload is not of subscription type")
		socket.SendErrorMessage(conn, "Payload is not of subscription type")
		return
	}

	b, _ := json.Marshal(payload.Data)
	var msg *types.WebSocketSubscription

	err = json.Unmarshal(b, &msg)
	if err != nil {
		logger.Error(err)
	}

	if (msg.Pair.BaseToken == common.Address{}) {
		message := map[string]string{"Message": "Invalid Base Token"}
		socket.SendErrorMessage(conn, message)
		return
	}

	if (msg.Pair.QuoteToken == common.Address{}) {
		message := map[string]string{"Message": "Invalid Quote Token"}
		socket.SendErrorMessage(conn, message)
		return
	}

	if msg.Event == types.SUBSCRIBE {
		e.orderBookService.SubscribeRawOrderBook(conn, msg.Pair.BaseToken, msg.Pair.QuoteToken)
	}

	if msg.Event == types.UNSUBSCRIBE {
		e.orderBookService.UnSubscribeRawOrderBook(conn, msg.Pair.BaseToken, msg.Pair.QuoteToken)
	}
}

// liteOrderBookWebSocket
func (e *OrderBookEndpoint) orderBookWebSocket(input interface{}, conn *ws.Conn) {
	bytes, _ := json.Marshal(input)
	var payload *types.WebSocketPayload
	err := json.Unmarshal(bytes, &payload)
	if err != nil {
		logger.Error(err)
	}

	socket := ws.GetOrderBookSocket()
	if payload.Type != "subscription" {
		message := map[string]string{"Message": "Invalid subscription payload"}
		socket.SendErrorMessage(conn, message)
		return
	}

	bytes, _ = json.Marshal(payload.Data)
	var msg *types.WebSocketSubscription

	err = json.Unmarshal(bytes, &msg)
	if err != nil {
		logger.Error(err)
		message := map[string]string{"Message": "Internal server error"}
		socket.SendErrorMessage(conn, message)
	}

	if (msg.Pair.BaseToken == common.Address{}) {
		message := map[string]string{"Message": "Invalid base token"}
		socket.SendErrorMessage(conn, message)
		return
	}

	if (msg.Pair.QuoteToken == common.Address{}) {
		message := map[string]string{"Message": "Invalid quote token"}
		socket.SendErrorMessage(conn, message)
		return
	}

	if msg.Event == types.SUBSCRIBE {
		e.orderBookService.SubscribeOrderBook(conn, msg.Pair.BaseToken, msg.Pair.QuoteToken)
	}

	if msg.Event == types.UNSUBSCRIBE {
		e.orderBookService.UnSubscribeOrderBook(conn, msg.Pair.BaseToken, msg.Pair.QuoteToken)
	}
}
