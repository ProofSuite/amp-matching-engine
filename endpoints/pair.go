package endpoints

import (
	"github.com/ethereum/go-ethereum/common"

	"github.com/Proofsuite/amp-matching-engine/errors"
	"github.com/Proofsuite/amp-matching-engine/services"
	"github.com/Proofsuite/amp-matching-engine/types"
	"github.com/go-ozzo/ozzo-routing"
)

type pairEndpoint struct {
	pairService *services.PairService
}

// ServePairResource sets up the routing of pair endpoints and the corresponding handlers.
func ServePairResource(rg *routing.RouteGroup, pairService *services.PairService) {
	r := &pairEndpoint{pairService}
	rg.Get("/pairs/<baseToken>/<quoteToken>", r.get)
	rg.Get("/pairs", r.query)
	rg.Post("/pairs", r.create)
}

func (r *pairEndpoint) create(c *routing.Context) error {
	var p types.Pair

	if err := c.Read(&p); err != nil {
		return err
	}

	if err := p.Validate(); err != nil {
		return err
	}

	err := r.pairService.Create(&p)
	if err != nil {
		return err
	}

	return c.Write(p)
}

func (r *pairEndpoint) query(c *routing.Context) error {
	res, err := r.pairService.GetAll()
	if err != nil {
		return err
	}

	return c.Write(res)
}

func (r *pairEndpoint) get(c *routing.Context) error {
	baseToken := c.Param("baseToken")
	if !common.IsHexAddress(baseToken) {
		return errors.NewAPIError(400, "INVALID_HEX_ADDRESS", nil)
	}

	quoteToken := c.Param("quoteToken")
	if !common.IsHexAddress(quoteToken) {
		return errors.NewAPIError(400, "INVALID_HEX_ADDRESS", nil)
	}

	baseTokenAddress := common.HexToAddress(baseToken)
	quoteTokenAddress := common.HexToAddress(quoteToken)

	res, err := r.pairService.GetByTokenAddress(baseTokenAddress, quoteTokenAddress)
	if err != nil {
		return err
	}

	return c.Write(res)
}

// func (r *pairEndpoint) orderBook(input interface{}, conn *websocket.Conn) {
// 	mab, _ := json.Marshal(input)
// 	var msg *types.Subscription
// 	if err := json.Unmarshal(mab, &msg); err != nil {
// 		log.Println("unmarshal to wsmsg <==>" + err.Error())
// 	}

// 	if msg.Pair.BaseToken == "" {
// 		message := map[string]string{
// 			"Code":    "Invalid_Pair_BaseToken",
// 			"Message": "Invalid Pair BaseToken passed in query Params",
// 		}
// 		ws.GetPairSockets().SendErrorMessage(conn, message)
// 		return
// 	}

// 	if msg.Pair.QuoteToken == "" {
// 		message := map[string]string{
// 			"Code":    "Invalid_Pair_QuoteToken",
// 			"Message": "Invalid Pair QuoteToken passed in query Params",
// 		}
// 		ws.GetPairSockets().SendErrorMessage(conn, message)
// 		return
// 	}

// 	if msg.Event == types.SUBSCRIBE {
// 		r.pairService.RegisterForOrderBook(conn, msg.Pair.BaseToken, msg.Pair.QuoteToken)
// 	}

// 	if msg.Event == types.UNSUBSCRIBE {
// 		r.pairService.UnRegisterForOrderBook(conn, msg.Pair.BaseToken, msg.Pair.QuoteToken)
// 	}
// }

// func (r *pairEndpoint) orderBookEndpoint(c *routing.Context) error {
// 	pName := c.Param("pairName")
// 	if pName == "" {
// 		return errors.NewAPIError(401, "EMPTY_PAIR_NAME", map[string]interface{}{})
// 	}

// 	bt := c.Param("baseToken")
// 	if !common.IsHexAddress(bt) {
// 		return errors.NewAPIError(400, "INVALID_HEX_ADDRESS", nil)
// 	}

// 	qt := c.Param("quoteToken")
// 	if !common.IsHexAddress(qt) {
// 		return errors.NewAPIError(400, "INVALID_HEX_ADDRESS", nil)
// 	}

// 	baseTokenAddress := common.HexToAddress(bt)
// 	quoteTokenAddress := common.HexToAddress(qt)
// 	ob, err := r.pairService.GetOrderBook(baseTokenAddress, quoteTokenAddress)
// 	if err != nil {
// 		return err
// 	}

// 	return c.Write(ob)
// }
