package endpoints

import (
	"encoding/json"
	"log"
	"strconv"
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
	rg.Get("/trades/history/<pair>", r.history)
	rg.Get("/trades/ticks", r.ticks)
	rg.Get("/trades/<addr>", r.get)
	ws.RegisterChannel("trade_ticks", r.wsTicks)
}

// history is reponsible for handling pair's trade history requests
func (r *tradeEndpoint) history(c *routing.Context) error {
	pair := c.Param("pair")
	if pair == "" {
		return errors.NewAPIError(400, "INVALID_PAIR_NAME", nil)
	}
	response, err := r.tradeService.GetByPairName(pair)
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
	startTs := time.Date(1970, time.January, 1, 0, 0, 0, 0, time.UTC)
	unit := c.Query("unit", "hour")
	duration := c.Query("duration", "24")
	from := c.Query("from", strconv.FormatInt(startTs.Unix(), 10))
	to := c.Query("to", strconv.FormatInt(time.Now().Unix(), 10))
	pairName := c.Query("pairName", "")
	if pairName == "" {
		return errors.NewAPIError(400, "EMPTY_PAIR_NAME", nil)
	}
	d, err := strconv.ParseInt(duration, 10, 64)
	if err != nil {
		return errors.NewAPIError(400, "INVALID_DURATION", nil)
	}
	fromTs, err := strconv.ParseInt(from, 10, 64)
	if err != nil {
		return errors.NewAPIError(400, "INVALID_FROM_TS", nil)
	}
	toTs, err := strconv.ParseInt(to, 10, 64)
	if err != nil {
		return errors.NewAPIError(400, "INVALID_TO_TS", nil)
	}

	res, err := r.tradeService.GetTicks(pairName, d, unit, fromTs, toTs)
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

	if msg.Key == "" {
		message := map[string]string{
			"Code":    "Invalid_Pair_Name",
			"Message": "Invalid Pair Name passed in query Params",
		}
		mab, _ := json.Marshal(message)
		conn.WriteMessage(1, mab)
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
		r.tradeService.RegisterForTicks(conn, msg.Key, &msg.Params)

	}
	if msg.Event == types.UNSUBSCRIBE {
		r.tradeService.UnregisterForTicks(conn, msg.Key, &msg.Params)
	}
}
