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
	var msg *types.Message

	bytes, _ := json.Marshal(input)
	if err := json.Unmarshal(bytes, &msg); err != nil {
		log.Println("unmarshal to wsmsg <==>" + err.Error())
	}

	if msg.Type == "NEW_ORDER" {
		r.handleNewOrder(msg * types.Message)
	} else if msg.Type == "CANCEL_ORDER" {
		r.handleCancelOrder(msg * types.Message)
	} else {
		ch := ws.GetOrderChannel(msg.Hash)
		if ch != nil {
			ch <- msg
		}
	}
}

func (r *orderEndpoint) handleNewOrder(msg *types.Message) {
	var or types.OrderRequest
	bytes, err = json.Marshal(msg.Data)

	err := json.Unmarshal(bytes, &or)
	or.Hash = or.ComputeHash()

	if err != nil {
		ws.OrderSendErrorMessage(conn, or.Hash, err.Error())
		return
	}

	if err = or.Validate(); err != nil {
		ws.OrderSendErrorMessage(conn, or.Hash, err.Error())
		return
	}

	ok, err := or.VerifySignature()
	if err != nil {
		ws.OrderSendErrorMessage(conn, or.Hash, err.Error())
		return
	}
	if !ok {
		ws.OrderSendErrorMessage(conn, or.Hash, "Invalid Signature")
		return
	}

	ws.RegisterOrderConnection(or.Hash, &ws.OrderConn{Conn: conn, ReadChannel: ch})
	ws.RegisterConnectionUnsubscribeHandler(conn, ws.OrderSocketUnsubscribeHandler(or.Hash))

	err = r.orderService.Create(or)
	if err != nil {
		ws.OrderSendErrorMessage(conn, or.Hash, err.Error())
		return
	}
}

func (r *orderEndpoint) handleCancelOrder(msg *types.Message) {
	bytes, err := json.Marshal(msg.Data)

	var order *types.Order
	if err := json.Unmarshal(bytes, &order); err != nil {
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
}

func (r *orderEndpoint) cancelOrder(order *types.Order) error {
	return r.orderService.CancelOrder(order)
}
