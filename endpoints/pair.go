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

type pairEndpoint struct {
	pairService *services.PairService
}

// ServePairResource sets up the routing of pair endpoints and the corresponding handlers.
func ServePairResource(rg *routing.RouteGroup, pairService *services.PairService) {
	r := &pairEndpoint{pairService}
	rg.Get("/pairs/book/<bt>/<st>", r.orderBookEndpoint)
	rg.Get("/pairs/<bt>/<st>", r.get)
	rg.Get("/pairs", r.query)
	rg.Post("/pairs", r.create)
	ws.RegisterChannel("order_book", r.orderBook)
}

func (r *pairEndpoint) create(c *routing.Context) error {
	var model types.Pair
	if err := c.Read(&model); err != nil {
		return err
	}
	if err := model.Validate(); err != nil {
		return err
	}

	err := r.pairService.Create(&model)
	if err != nil {
		return err
	}

	return c.Write(model)
}

func (r *pairEndpoint) query(c *routing.Context) error {

	response, err := r.pairService.GetAll()
	if err != nil {
		return err
	}

	return c.Write(response)
}

func (r *pairEndpoint) get(c *routing.Context) error {
	baseToken := c.Param("bt")
	if !common.IsHexAddress(baseToken) {
		return errors.NewAPIError(400, "INVALID_HEX_ADDRESS", nil)
	}
	quoteToken := c.Param("st")
	if !common.IsHexAddress(quoteToken) {
		return errors.NewAPIError(400, "INVALID_HEX_ADDRESS", nil)
	}
	response, err := r.pairService.GetByTokenAddress(baseToken, quoteToken)
	if err != nil {
		return err
	}

	return c.Write(response)
}
func (r *pairEndpoint) orderBook(input *interface{}, conn *websocket.Conn) {
	mab, _ := json.Marshal(input)
	var msg *types.Subscription
	if err := json.Unmarshal(mab, &msg); err != nil {
		log.Println("unmarshal to wsmsg <==>" + err.Error())
	}

	if msg.Pair.BaseToken == "" {
		message := map[string]string{
			"Code":    "Invalid_Pair_BaseToken",
			"Message": "Invalid Pair BaseToken passed in query Params",
		}
		mab, _ := json.Marshal(message)
		conn.WriteMessage(1, mab)
		return

	}
	if msg.Pair.QuoteToken == "" {
		message := map[string]string{
			"Code":    "Invalid_Pair_BaseToken",
			"Message": "Invalid Pair BaseToken passed in query Params",
		}
		mab, _ := json.Marshal(message)
		conn.WriteMessage(1, mab)
		return
	}
	if msg.Event == types.SUBSCRIBE {
		r.pairService.RegisterForOrderBook(conn, msg.Pair.BaseToken, msg.Pair.QuoteToken)
	}
	if msg.Event == types.UNSUBSCRIBE {
		r.pairService.UnRegisterForOrderBook(conn, msg.Pair.BaseToken, msg.Pair.QuoteToken)
	}
}
func (r *pairEndpoint) orderBookEndpoint(c *routing.Context) error {
	pName := c.Param("pairName")
	if pName == "" {
		return errors.NewAPIError(401, "EMPTY_PAIR_NAME", map[string]interface{}{})
	}
	baseToken := c.Param("bt")
	if !common.IsHexAddress(baseToken) {
		return errors.NewAPIError(400, "INVALID_HEX_ADDRESS", nil)
	}
	quoteToken := c.Param("st")
	if !common.IsHexAddress(quoteToken) {
		return errors.NewAPIError(400, "INVALID_HEX_ADDRESS", nil)
	}
	ob, err := r.pairService.GetOrderBook(baseToken, quoteToken)
	if err != nil {
		return err
	}
	return c.Write(ob)
}
