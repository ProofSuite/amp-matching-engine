package endpoints

import (
	"encoding/json"
	"log"

	"github.com/ethereum/go-ethereum/common"

	"github.com/Proofsuite/amp-matching-engine/errors"
	"github.com/Proofsuite/amp-matching-engine/services"
	"github.com/Proofsuite/amp-matching-engine/types"
	"github.com/Proofsuite/amp-matching-engine/ws"
	"github.com/go-ozzo/ozzo-routing"
	"github.com/gorilla/websocket"
)

type OrderBookEndpoint struct {
	orderBookService *services.OrderBookService
}

// ServePairResource sets up the routing of pair endpoints and the corresponding handlers.
func ServeOrderBookResource(rg *routing.RouteGroup, orderBookService *services.OrderBookService) {
	e := &OrderBookEndpoint{orderBookService}

	rg.Get("/orderbook/<baseToken>/<quoteToken>", e.orderBookEndpoint)
	ws.RegisterChannel(ws.OrderBookChannel, e.orderBookWebSocket)
}

func (e *OrderBookEndpoint) orderBookEndpoint(c *routing.Context) error {
	p := c.Param("pair")
	if p == "" {
		return errors.NewAPIError(401, "EMPTY_PAIR_NAME", map[string]interface{}{})
	}

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

func (e *OrderBookEndpoint) orderBookWebSocket(input interface{}, conn *websocket.Conn) {
	mab, _ := json.Marshal(input)
	var msg *types.Subscription
	if err := json.Unmarshal(mab, &msg); err != nil {
		log.Println("unmarshal to wsmsg <==>" + err.Error())
	}

	//NOTE: Should be BaseToken rather be a string or a common address ?
	if (msg.Pair.BaseToken == common.Address{}) {
		message := map[string]string{
			"Code":    "Invalid_Pair_BaseToken",
			"Message": "Invalid Pair BaseToken passed in query Params",
		}

		ws.GetPairSockets().SendErrorMessage(conn, message)
		return
	}

	//NOTE: Should QuoteToken rather be a string or a common address ?
	if (msg.Pair.QuoteToken == common.Address{}) {
		message := map[string]string{
			"Code":    "Invalid_Pair_QuoteToken",
			"Message": "Invalid Pair QuoteToken passed in query Params",
		}

		ws.GetPairSockets().SendErrorMessage(conn, message)
		return
	}

	if msg.Event == types.SUBSCRIBE {
		e.orderBookService.SubscribeOrderBook(conn, msg.Pair.BaseToken, msg.Pair.QuoteToken)
	}

	if msg.Event == types.UNSUBSCRIBE {
		e.orderBookService.UnSubscribeOrderBook(conn, msg.Pair.BaseToken, msg.Pair.QuoteToken)
	}
}
