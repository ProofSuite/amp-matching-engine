package endpoints

import (
	"encoding/json"
	"log"
	"time"

	"github.com/Proofsuite/amp-matching-engine/services"
	"github.com/Proofsuite/amp-matching-engine/types"
	"github.com/Proofsuite/amp-matching-engine/ws"
	"github.com/ethereum/go-ethereum/common"
	"github.com/go-ozzo/ozzo-routing"
	"github.com/gorilla/websocket"
)

type OHLCVEndpoint struct {
	ohlcvService services.OHLCVServiceInterface
}

func ServeOHLCVResource(
	rg *routing.RouteGroup,
	ohlcvService services.OHLCVServiceInterface,
) {
	e := &OHLCVEndpoint{ohlcvService}
	rg.Post("/ohlcv", e.ohlcv)
	ws.RegisterChannel(ws.OHLCVChannel, e.ohlcvWebSocket)
}

func (e *OHLCVEndpoint) ohlcv(c *routing.Context) error {
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

	res, err := e.ohlcvService.GetOHLCV(model.Pair, model.Duration, model.Units, model.From, model.To)
	if err != nil {
		return err
	}

	return c.Write(res)
}

func (e *OHLCVEndpoint) ohlcvWebSocket(input interface{}, conn *websocket.Conn) {
	startTs := time.Date(1970, time.January, 1, 0, 0, 0, 0, time.UTC)

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
		ws.SendOHLCVErrorMessage(conn, message)
		return
	}

	if (msg.Pair.QuoteToken == common.Address{}) {
		message := map[string]string{
			"Code":    "Invalid_Pair_BaseToken",
			"Message": "Invalid Pair BaseToken passed in Params",
		}
		ws.SendOHLCVErrorMessage(conn, message)
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
