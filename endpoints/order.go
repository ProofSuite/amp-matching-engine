package endpoints

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

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
	e.SubscribeEngineResponse(r.engineResponse)
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
	if msg.MsgType == "new_order" {
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
		r.orderService.SendMessage("added_to_orderbook", order.Hash, order)

	} else if msg.MsgType == "cancel_order" {
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

func (r *orderEndpoint) engineResponse(engineResponse *engine.Response) error {
	if engineResponse.FillStatus == engine.NOMATCH {
		r.orderService.SendMessage("added_to_orderbook", engineResponse.Order.Hash, engineResponse)
	} else {
		r.orderService.SendMessage("trade_remaining_order_sign", engineResponse.Order.Hash, engineResponse)

		t := time.NewTimer(10 * time.Second)
		ch := ws.GetOrderChannel(engineResponse.Order.Hash)
		if ch == nil {
			r.orderService.RecoverOrders(engineResponse)
		} else {

			select {
			case rm := <-ch:
				if rm.MsgType == "trade_remaining_order_sign" {
					mb, err := json.Marshal(rm.Data)
					if err != nil {
						fmt.Printf("=== Error while marshaling EngineResponse===")
						r.orderService.RecoverOrders(engineResponse)
						ws.OrderSendErrorMessage(ws.GetOrderConn(engineResponse.Order.Hash), engineResponse.Order.Hash, err.Error())
					}

					var ersb *engine.Response
					err = json.Unmarshal(mb, &ersb)
					if err != nil {
						fmt.Printf("=== Error while unmarshaling EngineResponse===")
						ws.OrderSendErrorMessage(ws.GetOrderConn(engineResponse.Order.Hash), engineResponse.Order.Hash, err.Error())
						r.orderService.RecoverOrders(engineResponse)
					}

					if engineResponse.FillStatus == engine.PARTIAL {
						engineResponse.Order.OrderBook = &types.OrderSubDoc{Amount: ersb.RemainingOrder.Amount, Signature: ersb.RemainingOrder.Signature}
						orderAsBytes, _ := json.Marshal(engineResponse.Order)
						r.engine.PublishMessage(&engine.Message{Type: "remaining_order_add", Data: orderAsBytes})
					}

				}
				t.Stop()
				break

			case <-t.C:
				fmt.Printf("\nTimeout\n")
				r.orderService.RecoverOrders(engineResponse)
				t.Stop()
				break
			}
		}
	}
	r.orderService.UpdateUsingEngineResponse(engineResponse)
	// TODO: send to operator for blockchain execution

	r.orderService.RelayUpdateOverSocket(engineResponse)
	ws.CloseOrderReadChannel(engineResponse.Order.Hash)

	return nil
}
func (r *orderEndpoint) cancelOrder(order *types.Order) error {
	return r.orderService.CancelOrder(order)
}
