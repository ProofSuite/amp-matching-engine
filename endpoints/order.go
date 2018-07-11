package endpoints

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
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
	engine       *engine.EngineResource
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// ServeOrder sets up the routing of order endpoints and the corresponding handlers.
func ServeOrderResource(rg *routing.RouteGroup, orderService *services.OrderService, e *engine.EngineResource) {
	r := &orderEndpoint{orderService, e}
	rg.Get("/orders/<addr>", r.get)
	// http.HandleFunc("/orders/ws", r.ws)
	// http.HandleFunc("/orders/book/<pair>", r.ws)
	ws.RegisterChannel("order_channel", r.ws)
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

func (r *orderEndpoint) ws(input *interface{}, conn *websocket.Conn) {

	ch := make(chan *types.OrderMessage)
	mab, _ := json.Marshal(input)
	var msg *types.OrderMessage
	if err := json.Unmarshal(mab, &msg); err != nil {
		log.Println("unmarshal to wsmsg <==>" + err.Error())
	}
	messageType := 1
	if msg.MsgType == "new_order" {
		oab, err := json.Marshal(msg.Data)

		var model types.OrderRequest
		if err := json.Unmarshal(oab, &model); err != nil {
			conn.WriteMessage(messageType, []byte(err.Error()))

			return
		}
		if err := model.Validate(); err != nil {
			conn.WriteMessage(messageType, []byte(err.Error()))

			return
		}
		if ok, err := model.VerifySignature(); err != nil {
			conn.WriteMessage(messageType, []byte(err.Error()))
			return
		} else if !ok {
			conn.WriteMessage(messageType, []byte("Invalid Signature"))
			return
		}
		order, err := model.ToOrder()
		if err != nil {
			conn.WriteMessage(messageType, []byte(err.Error()))
			conn.Close()

			return
		}
		err = r.orderService.Create(order)
		if err != nil {
			conn.WriteMessage(messageType, []byte(err.Error()))
			return
		}
		oab, _ = json.Marshal(order)
		conn.WriteMessage(messageType, oab)
		ws.RegisterOrderConnection(order.ID, &ws.WsOrderConn{Conn: conn, ReadChannel: ch})
		ws.RegisterConnectionUnsubscribeHandler(conn, ws.OrderSocketCloseHandler(order.ID))
	} else if msg.MsgType == "order_cancel" {
		oab, err := json.Marshal(msg.Data)

		var model *types.Order
		if err := json.Unmarshal(oab, &model); err != nil {
			conn.WriteMessage(messageType, []byte(err.Error()))

			return
		}
		err = r.cancelOrder(model)
		if err != nil {
			conn.WriteMessage(messageType, []byte(err.Error()))
			return
		}
	} else {
		ch := ws.GetOrderChannel(msg.OrderID)
		if ch != nil {
			ch <- msg
		}
	}
}

func (r *orderEndpoint) engineResponse(engineResponse *engine.EngineResponse) error {
	if engineResponse.FillStatus == engine.NO_MATCH {
		r.orderService.SendMessage("added_to_orderbook", engineResponse.Order.ID, engineResponse)
	} else {
		r.orderService.SendMessage("trade_remorder_sign", engineResponse.Order.ID, engineResponse)

		t := time.NewTimer(10 * time.Second)
		ch := ws.GetOrderChannel(engineResponse.Order.ID)
		fmt.Println(ch)
		if ch == nil {
			r.orderService.RecoverOrders(engineResponse)
		} else {

			select {
			case rm := <-ch:
				if rm.MsgType == "trade_remorder_sign" {
					mb, err := json.Marshal(rm.Data)
					if err != nil {
						r.orderService.RecoverOrders(engineResponse)
						ws.GetOrderConn(engineResponse.Order.ID).WriteMessage(1, []byte(err.Error()))
					}

					var ersb *engine.EngineResponse
					err = json.Unmarshal(mb, &ersb)
					if err != nil {
						ws.GetOrderConn(engineResponse.Order.ID).WriteMessage(1, []byte(err.Error()))
						r.orderService.RecoverOrders(engineResponse)
					}

					if engineResponse.FillStatus == engine.PARTIAL {
						engineResponse.Order.OrderBook = &types.OrderSubDoc{Amount: ersb.RemainingOrder.Amount, Signature: ersb.RemainingOrder.Signature}
						orderAsBytes, _ := json.Marshal(engineResponse.Order)
						r.engine.PublishMessage(&engine.Message{Type: "remaining_order_add", Data: orderAsBytes})
					}
					r.orderService.UpdateUsingEngineResponse(engineResponse)
					// TODO: send to operator for blockchain execution

					r.orderService.RelayUpdateOverSocket(engineResponse)
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
	ws.CloseOrderReadChannel(engineResponse.Order.ID)

	return nil
}
func (r *orderEndpoint) cancelOrder(order *types.Order) error {
	return r.orderService.CancelOrder(order)
}
