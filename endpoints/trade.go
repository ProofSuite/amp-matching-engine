package endpoints

import (
	"strconv"
	"time"

	"github.com/ethereum/go-ethereum/common"

	"github.com/Proofsuite/amp-matching-engine/errors"
	"github.com/Proofsuite/amp-matching-engine/services"
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

	r.tradeService.GetTicks(pairName, d, unit, fromTs, toTs)
	return c.Write("")
}
