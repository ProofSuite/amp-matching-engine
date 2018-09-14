package endpoints

import (
	"encoding/json"

	"github.com/Proofsuite/amp-matching-engine/errors"
	"github.com/Proofsuite/amp-matching-engine/interfaces"
	"github.com/Proofsuite/amp-matching-engine/types"
	"github.com/Proofsuite/amp-matching-engine/ws"
	"github.com/ethereum/go-ethereum/common"
	"github.com/go-ozzo/ozzo-routing"
)

type OrderBookEndpoint struct {
	orderBookService interfaces.OrderBookService
}

// ServePairResource sets up the routing of pair endpoints and the corresponding handlers.
func ServeOrderBookResource(
	rg *routing.RouteGroup,
	orderBookService interfaces.OrderBookService,
) {
	e := &OrderBookEndpoint{orderBookService}

	rg.Get("/orderbook/<baseToken>/<quoteToken>/full", e.fullOrderBookEndpoint)
	rg.Get("/orderbook/<baseToken>/<quoteToken>", e.orderBookEndpoint)
	ws.RegisterChannel(ws.LiteOrderBookChannel, e.liteOrderBookWebSocket)
	ws.RegisterChannel(ws.RawOrderBookChannel, e.fullOrderBookWebSocket)
}

// orderBookEndpoint
func (e *OrderBookEndpoint) orderBookEndpoint(c *routing.Context) error {

	bt := c.Param("baseToken")
	if !common.IsHexAddress(bt) {
		return errors.NewHTTPError(400, "Invalid Hex Address", nil)
	}

	qt := c.Param("quoteToken")
	if !common.IsHexAddress(qt) {
		return errors.NewHTTPError(400, "Invalid Hex Address", nil)
	}

	baseTokenAddress := common.HexToAddress(bt)
	quoteTokenAddress := common.HexToAddress(qt)
	ob, err := e.orderBookService.GetOrderBook(baseTokenAddress, quoteTokenAddress)
	if err != nil {
		return errors.NewHTTPError(500, "Internal Server Error", nil)
	}

	return c.Write(ob)
}

// orderBookEndpoint
func (e *OrderBookEndpoint) fullOrderBookEndpoint(c *routing.Context) error {

	bt := c.Param("baseToken")
	if !common.IsHexAddress(bt) {
		return errors.NewHTTPError(400, "Invalid Hex Address", nil)
	}

	qt := c.Param("quoteToken")
	if !common.IsHexAddress(qt) {
		return errors.NewHTTPError(400, "Invalid Hex Address", nil)
	}

	baseTokenAddress := common.HexToAddress(bt)
	quoteTokenAddress := common.HexToAddress(qt)
	ob, err := e.orderBookService.GetRawOrderBook(baseTokenAddress, quoteTokenAddress)
	if err != nil {
		return errors.NewHTTPError(500, "Internal Server Error", nil)
	}

	return c.Write(ob)
}

// liteOrderBookWebSocket
func (e *OrderBookEndpoint) fullOrderBookWebSocket(input interface{}, conn *ws.Conn) {
	mab, _ := json.Marshal(input)
	var payload *types.WebSocketPayload

	err := json.Unmarshal(mab, &payload)
	if err != nil {
		logger.Error(err)
		return
	}

	socket := ws.GetLiteOrderBookSocket()

	if payload.Type != "subscription" {
		logger.Error("Payload is not of subscription type")
		socket.SendErrorMessage(conn, "Payload is not of subscription type")
		return
	}

	dab, _ := json.Marshal(payload.Data)
	var msg *types.WebSocketSubscription

	err = json.Unmarshal(dab, &msg)
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
func (e *OrderBookEndpoint) liteOrderBookWebSocket(input interface{}, conn *ws.Conn) {
	bytes, _ := json.Marshal(input)
	var payload *types.WebSocketPayload
	err := json.Unmarshal(bytes, &payload)
	if err != nil {
		logger.Error(err)
	}

	socket := ws.GetLiteOrderBookSocket()
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
