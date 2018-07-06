package endpoints

import (
	"encoding/json"
	"log"

	"github.com/Proofsuite/amp-matching-engine/errors"
	"github.com/Proofsuite/amp-matching-engine/services"
	"github.com/Proofsuite/amp-matching-engine/types"
	"github.com/Proofsuite/amp-matching-engine/ws"
	"github.com/go-ozzo/ozzo-routing"
	"github.com/gorilla/websocket"
	"gopkg.in/mgo.v2/bson"
)

type pairEndpoint struct {
	pairService *services.PairService
}

// ServePair sets up the routing of pair endpoints and the corresponding handlers.
func ServePairResource(rg *routing.RouteGroup, pairService *services.PairService) {
	r := &pairEndpoint{pairService}
	rg.Get("/pairs/book/<pairName>", r.orderBookEndpoint)
	rg.Get("/pairs/<id>", r.get)
	rg.Get("/pairs", r.query)
	rg.Post("/pairs", r.create)
	// http.HandleFunc("/pairs/book", r.orderBook)
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
	id := c.Param("id")
	if !bson.IsObjectIdHex(id) {
		return errors.NewAPIError(400, "INVALID_ID", nil)
	}
	response, err := r.pairService.GetByID(bson.ObjectIdHex(id))
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

	if msg.Key == "" {
		message := map[string]string{
			"Code":    "Invalid_Pair_Name",
			"Message": "Invalid Pair Name passed in query Params",
		}
		mab, _ := json.Marshal(message)
		conn.WriteMessage(1, mab)
	}
	if msg.Event == types.SUBSCRIBE {
		r.pairService.RegisterForOrderBook(conn, msg.Key)
	}
	if msg.Event == types.UNSUBSCRIBE {
		r.pairService.UnRegisterForOrderBook(conn, msg.Key)
	}
}
func (r *pairEndpoint) orderBookEndpoint(c *routing.Context) error {
	pName := c.Param("pairName")
	if pName == "" {
		return errors.NewAPIError(401, "EMPTY_PAIR_NAME", map[string]interface{}{})
	}
	ob, err := r.pairService.GetOrderBook(pName)
	if err != nil {
		return err
	}
	return c.Write(ob)
}
