package endpoints

import (
	"encoding/json"
	"log"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/gorilla/websocket"

	"github.com/Proofsuite/amp-matching-engine/errors"
	"github.com/Proofsuite/amp-matching-engine/services"
	"github.com/Proofsuite/amp-matching-engine/types"
	"github.com/Proofsuite/amp-matching-engine/ws"
	"github.com/go-ozzo/ozzo-routing"
)

type tradeEndpoint struct {
	tradeService *services.TradeService
}

// ServeTradeResource sets up the routing of trade endpoints and the corresponding handlers.
func ServeTradeResource(rg *routing.RouteGroup, tradeService *services.TradeService) {
	r := &tradeEndpoint{tradeService}
	rg.Get("/trades/history/<bt>/<qt>", r.history)
	rg.Post("/trades/ticks", r.ticks)
	rg.Get("/trades/<addr>", r.get)
	ws.RegisterChannel("trades", r.wsTicks)
}

// history is reponsible for handling pair's trade history requests
func (r *tradeEndpoint) history(c *routing.Context) error {
	pair := c.Param("pair")
	if pair == "" {
		return errors.NewAPIError(400, "INVALID_PAIR_NAME", nil)
	}
	baseToken := c.Param("bt")
	if !common.IsHexAddress(baseToken) {
		return errors.NewAPIError(400, "INVALID_HEX_ADDRESS", nil)
	}
	quoteToken := c.Param("qt")
	if !common.IsHexAddress(quoteToken) {
		return errors.NewAPIError(400, "INVALID_HEX_ADDRESS", nil)
	}
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
	response, err := r.tradeService.GetByUserAddress(addr)
	if err != nil {
		return err
	}

	return c.Write(response)
}

func (r *tradeEndpoint) ticks(c *routing.Context) error {
	var model types.TickRequest
	if err := c.Read(&model); err != nil {
		return err
	}

	startTs := time.Date(1970, time.January, 1, 0, 0, 0, 0, time.UTC)
	if model.Units == "" {
		model.Units = "hour"
	}
	if model.Duration == 0 {
		model.Duration = 24
	}
	if model.From == 0 {
		model.From = startTs.Unix()
	}
	if model.To == 0 {
		model.To = time.Now().Unix()
	}
	res, err := r.tradeService.GetTicks(model.Pair, model.Duration, model.Units, model.From, model.To)
	if err != nil {
		return err
	}
	return c.Write(res)
}
func (r *tradeEndpoint) wsTicks(input *interface{}, conn *websocket.Conn) {
	startTs := time.Date(1970, time.January, 1, 0, 0, 0, 0, time.UTC)

	mab, _ := json.Marshal(input)
	var msg *types.Subscription
	if err := json.Unmarshal(mab, &msg); err != nil {
		log.Println("unmarshal to wsmsg <==>" + err.Error())
	}

	if msg.Pair.BaseToken == "" {
		message := map[string]string{
			"Code":    "Invalid_Pair_BaseToken",
			"Message": "Invalid Pair BaseToken passed in Params",
		}
		mab, _ := json.Marshal(message)
		conn.WriteMessage(1, mab)
		return

	}
	if msg.Pair.QuoteToken == "" {
		message := map[string]string{
			"Code":    "Invalid_Pair_BaseToken",
			"Message": "Invalid Pair BaseToken passed in Params",
		}
		mab, _ := json.Marshal(message)
		conn.WriteMessage(1, mab)
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
		r.tradeService.RegisterForTicks(conn, msg.Pair.BaseToken, msg.Pair.QuoteToken, &msg.Params)

	}
	if msg.Event == types.UNSUBSCRIBE {
		r.tradeService.UnregisterForTicks(conn, msg.Pair.BaseToken, msg.Pair.QuoteToken, &msg.Params)
	}
}
