package endpoints

import (
	"encoding/json"
	"log"

	"github.com/Proofsuite/amp-matching-engine/errors"
	"github.com/Proofsuite/amp-matching-engine/services"
	"github.com/Proofsuite/amp-matching-engine/types"
	"github.com/Proofsuite/amp-matching-engine/ws"
	"github.com/ethereum/go-ethereum/common"
	"github.com/go-ozzo/ozzo-routing"
	"github.com/gorilla/websocket"
)

type tradeEndpoint struct {
	tradeService services.TradeServiceInterface
}

// ServeTradeResource sets up the routing of trade endpoints and the corresponding handlers.
func ServeTradeResource(
	rg *routing.RouteGroup,
	tradeService services.TradeServiceInterface,
) {
	e := &tradeEndpoint{tradeService}
	rg.Get("/trades/history/<bt>/<qt>", e.history)
	rg.Get("/trades/<addr>", e.get)

	ws.RegisterChannel(ws.TradeChannel, e.tradeWebSocket)
}

// history is reponsible for handling pair's trade history requests
func (r *tradeEndpoint) history(c *routing.Context) error {
	bt := c.Param("bt")
	if !common.IsHexAddress(bt) {
		return errors.NewAPIError(400, "INVALID_HEX_ADDRESS", nil)
	}

	qt := c.Param("qt")
	if !common.IsHexAddress(qt) {
		return errors.NewAPIError(400, "INVALID_HEX_ADDRESS", nil)
	}

	baseToken := common.HexToAddress(bt)
	quoteToken := common.HexToAddress(qt)
	response, err := r.tradeService.GetByPairAddress(baseToken, quoteToken)
	if err != nil {
		return err
	}

	return c.Write(response)
}

// get is reponsible for handling user's trade history requests
func (r *tradeEndpoint) get(c *routing.Context) error {
	addr := c.Param("addr")
	if !common.IsHexAddress(addr) {
		return errors.NewAPIError(400, "INVALID_ADDRESS", nil)
	}

	address := common.HexToAddress(addr)
	response, err := r.tradeService.GetByUserAddress(address)
	if err != nil {
		return err
	}

	return c.Write(response)
}

func (e *tradeEndpoint) tradeWebSocket(input interface{}, conn *websocket.Conn) {
	mab, _ := json.Marshal(input)
	var msg *types.WebSocketSubscription
	if err := json.Unmarshal(mab, &msg); err != nil {
		log.Println("unmarshal to wsmsg <==>" + err.Error())
	}

	if (msg.Pair.BaseToken == common.Address{}) {
		message := map[string]string{
			"Code":    "Invalid_Pair_BaseToken",
			"Message": "Invalid Pair BaseToken passed in Params",
		}
		ws.SendTradeErrorMessage(conn, message)
		return
	}

	if (msg.Pair.QuoteToken == common.Address{}) {
		message := map[string]string{
			"Code":    "Invalid_Pair_BaseToken",
			"Message": "Invalid Pair BaseToken passed in Params",
		}
		ws.SendTradeErrorMessage(conn, message)
		return
	}

	if msg.Event == types.SUBSCRIBE {
		e.tradeService.Subscribe(conn, msg.Pair.BaseToken, msg.Pair.QuoteToken)
	}

	if msg.Event == types.UNSUBSCRIBE {
		e.tradeService.Unsubscribe(conn, msg.Pair.BaseToken, msg.Pair.QuoteToken)
	}
}
