package endpoints

import (
	"encoding/json"
	"log"

	"github.com/Proofsuite/amp-matching-engine/errors"
	"github.com/ethereum/go-ethereum/common"

	"github.com/Proofsuite/amp-matching-engine/engine"
	"github.com/Proofsuite/amp-matching-engine/services"
	"github.com/Proofsuite/amp-matching-engine/types"
	"github.com/Proofsuite/amp-matching-engine/ws"
	"github.com/go-ozzo/ozzo-routing"
	"github.com/gorilla/websocket"
)

type orderEndpoint struct {
	orderService *services.OrderService
	engine       *engine.Resource
}

// ServeOrderResource sets up the routing of order endpoints and the corresponding handlers.
func ServeOrderResource(rg *routing.RouteGroup, orderService *services.OrderService, e *engine.Resource) {
	r := &orderEndpoint{orderService, e}
	rg.Get("/orders/<addr>", r.get)
	ws.RegisterChannel(ws.OrderChannel, r.ws)
	e.SubscribeEngineResponse(r.orderService.EngineResponse)
}

func (r *orderEndpoint) get(c *routing.Context) error {
	address := c.Param("addr")
	if !common.IsHexAddress(address) {
		return errors.NewAPIError(400, "Invalid Adrress", map[string]interface{}{})
	}
	orders, err := r.orderService.GetByUserAddress(address)
	if err != nil {
		return errors.NewAPIError(400, "Fetch Error", map[string]interface{}{})
	}
	return c.Write(orders)
}

func (r *orderEndpoint) ws(input interface{}, conn *websocket.Conn) {

	ch := make(chan *types.Message)
	mab, _ := json.Marshal(input)
	var msg *types.Message
	if err := json.Unmarshal(mab, &msg); err != nil {
		log.Println("unmarshal to wsmsg <==>" + err.Error())
	}
	if msg.MsgType == "NEW_ORDER" {
		oab, err := json.Marshal(msg.Data)

		var model types.OrderRequest
		if err := json.Unmarshal(oab, &model); err != nil {
			ws.OrderSendErrorMessage(conn, model.ComputeHash(), err.Error())
			return
		}
		if err := model.Validate(); err != nil {
			ws.OrderSendErrorMessage(conn, model.ComputeHash(), err.Error())
			return
		}
		if ok, err := model.VerifySignature(); err != nil {
			ws.OrderSendErrorMessage(conn, model.ComputeHash(), err.Error())
			return
		} else if !ok {
			ws.OrderSendErrorMessage(conn, model.ComputeHash(), "Invalid Signature")
			return
		}
		order, err := model.ToOrder()
		if err != nil {
			ws.OrderSendErrorMessage(conn, model.ComputeHash(), err.Error())
			return
		}
		ws.RegisterOrderConnection(order.Hash, &ws.OrderConn{Conn: conn, ReadChannel: ch})
		ws.RegisterConnectionUnsubscribeHandler(conn, ws.OrderSocketUnsubscribeHandler(order.Hash))

		err = r.orderService.Create(order)
		if err != nil {
			ws.OrderSendErrorMessage(conn, model.ComputeHash(), err.Error())
			return
		}
		r.orderService.SendMessage("ORDER_ADDED", order.Hash, order)

	} else if msg.MsgType == "CANCEL_ORDER" {
		oab, err := json.Marshal(msg.Data)

		var order *types.Order
		if err := json.Unmarshal(oab, &order); err != nil {
			ws.OrderSendErrorMessage(conn, order.Hash, err.Error())
			return
		}

		ws.RegisterOrderConnection(order.Hash, &ws.OrderConn{Conn: conn, Active: true})
		ws.RegisterConnectionUnsubscribeHandler(conn, ws.OrderSocketUnsubscribeHandler(order.Hash))

		err = r.cancelOrder(order)
		if err != nil {
			ws.OrderSendErrorMessage(conn, order.Hash, err.Error())
			return
		}
	} else {
		ch := ws.GetOrderChannel(msg.Hash)
		if ch != nil {
			ch <- msg
		}
	}
}

func (r *orderEndpoint) cancelOrder(order *types.Order) error {
	return r.orderService.CancelOrder(order)
}
