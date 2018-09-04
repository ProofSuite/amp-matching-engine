package endpoints

import (
	"encoding/json"
	"log"

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

	rg.Get("/orderbook/<baseToken>/<quoteToken>", e.orderBookEndpoint)
	ws.RegisterChannel(ws.LiteOrderBookChannel, e.liteOrderBookWebSocket)
	ws.RegisterChannel(ws.FullOrderBookChannel, e.fullOrderBookWebSocket)
}

// orderBookEndpoint
func (e *OrderBookEndpoint) orderBookEndpoint(c *routing.Context) error {

	bt := c.Param("baseToken")
	if !common.IsHexAddress(bt) {
		return errors.NewAPIError(400, "INVALID_HEX_ADDRESS", nil)
	}

	qt := c.Param("quoteToken")
	if !common.IsHexAddress(qt) {
		return errors.NewAPIError(400, "INVALID_HEX_ADDRESS", nil)
	}

	baseTokenAddress := common.HexToAddress(bt)
	quoteTokenAddress := common.HexToAddress(qt)
	ob, err := e.orderBookService.GetOrderBook(baseTokenAddress, quoteTokenAddress)
	if err != nil {
		return err
	}

	return c.Write(ob)
}

// liteOrderBookWebSocket
func (e *OrderBookEndpoint) fullOrderBookWebSocket(input interface{}, conn *ws.Conn) {
	mab, _ := json.Marshal(input)
	var payload *types.WebSocketPayload
	if err := json.Unmarshal(mab, &payload); err != nil {
		log.Println("unmarshal to wsmsg <==>" + err.Error())
	}

	socket := ws.GetLiteOrderBookSocket()

	if payload.Type != "subscription" {
		log.Println("Payload is not of subscription type")
		socket.SendErrorMessage(conn, "Payload is not of subscription type")
		return
	}
	dab, _ := json.Marshal(payload.Data)
	var msg *types.WebSocketSubscription
	if err := json.Unmarshal(dab, &msg); err != nil {
		log.Println("unmarshal to wsmsg <==>" + err.Error())
	}
	if (msg.Pair.BaseToken == common.Address{}) {
		message := map[string]string{
			"Code":    "Invalid_Pair_BaseToken",
			"Message": "Invalid Pair BaseToken passed in query Params",
		}

		socket.SendErrorMessage(conn, message)
		return
	}

	if (msg.Pair.QuoteToken == common.Address{}) {
		message := map[string]string{
			"Code":    "Invalid_Pair_QuoteToken",
			"Message": "Invalid Pair QuoteToken passed in query Params",
		}

		socket.SendErrorMessage(conn, message)
		return
	}

	if msg.Event == types.SUBSCRIBE {
		e.orderBookService.SubscribeFull(conn, msg.Pair.BaseToken, msg.Pair.QuoteToken)
	}

	if msg.Event == types.UNSUBSCRIBE {
		e.orderBookService.UnsubscribeFull(conn, msg.Pair.BaseToken, msg.Pair.QuoteToken)
	}
}

// liteOrderBookWebSocket
func (e *OrderBookEndpoint) liteOrderBookWebSocket(input interface{}, conn *ws.Conn) {
	mab, _ := json.Marshal(input)
	var payload *types.WebSocketPayload
	if err := json.Unmarshal(mab, &payload); err != nil {
		log.Println("unmarshal to wsmsg <==>" + err.Error())
	}

	socket := ws.GetLiteOrderBookSocket()
	if payload.Type != "subscription" {
		log.Println("Payload is not of subscription type")
		socket.SendErrorMessage(conn, "Payload is not of subscription type")
		return
	}
	dab, _ := json.Marshal(payload.Data)
	var msg *types.WebSocketSubscription
	if err := json.Unmarshal(dab, &msg); err != nil {
		log.Println("unmarshal to wsmsg <==>" + err.Error())
	}
	if (msg.Pair.BaseToken == common.Address{}) {
		message := map[string]string{
			"Code":    "Invalid_Pair_BaseToken",
			"Message": "Invalid Pair BaseToken passed in query Params",
		}

		socket.SendErrorMessage(conn, message)
		return
	}

	if (msg.Pair.QuoteToken == common.Address{}) {
		message := map[string]string{
			"Code":    "Invalid_Pair_QuoteToken",
			"Message": "Invalid Pair QuoteToken passed in query Params",
		}

		socket.SendErrorMessage(conn, message)
		return
	}

	if msg.Event == types.SUBSCRIBE {
		e.orderBookService.SubscribeLite(conn, msg.Pair.BaseToken, msg.Pair.QuoteToken)
	}

	if msg.Event == types.UNSUBSCRIBE {
		e.orderBookService.UnsubscribeLite(conn, msg.Pair.BaseToken, msg.Pair.QuoteToken)
	}
}
