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
	"github.com/gorilla/websocket"
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
	ws.RegisterChannel(ws.OrderBookChannel, e.orderBookWebSocket)
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

// orderBookWebSocket
func (e *OrderBookEndpoint) orderBookWebSocket(input interface{}, conn *websocket.Conn) {
	mab, _ := json.Marshal(input)
	var msg *types.WebSocketSubscription
	if err := json.Unmarshal(mab, &msg); err != nil {
		log.Println("unmarshal to wsmsg <==>" + err.Error())
	}

	if (msg.Pair.BaseToken == common.Address{}) {
		message := map[string]string{
			"Code":    "Invalid_Pair_BaseToken",
			"Message": "Invalid Pair BaseToken passed in query Params",
		}

		ws.SendOrderBookErrorMessage(conn, message)
		return
	}

	if (msg.Pair.QuoteToken == common.Address{}) {
		message := map[string]string{
			"Code":    "Invalid_Pair_QuoteToken",
			"Message": "Invalid Pair QuoteToken passed in query Params",
		}

		ws.SendOrderBookErrorMessage(conn, message)
		return
	}

	if msg.Event == types.SUBSCRIBE {
		e.orderBookService.Subscribe(conn, msg.Pair.BaseToken, msg.Pair.QuoteToken)
	}

	if msg.Event == types.UNSUBSCRIBE {
		e.orderBookService.Unsubscribe(conn, msg.Pair.BaseToken, msg.Pair.QuoteToken)
	}
}
