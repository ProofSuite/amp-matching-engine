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
	r := &orderEndpoint{orderService}
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
	} else {
		ch := ws.GetOrderChannel(msg.OrderID)
		if ch != nil {
			ch <- msg
		}
	}
}

func (r *orderEndpoint) engineResponse(engineResponse *engine.EngineResponse) error {
	// b, _ := json.Marshal(er)
	// fmt.Printf("\n======> \n%s\n <======\n", b)
	// If NO_MATCH add to order book
	if engineResponse.FillStatus == engine.NO_MATCH {

		msg := &types.OrderMessage{MsgType: "added_to_orderbook"}
		msg.OrderID = engineResponse.Order.ID
		msg.Data = engineResponse
		ws.GetOrderConn(engineResponse.Order.ID).WriteJSON(msg)

	} else {
		msg := &types.OrderMessage{MsgType: "trade_remorder_sign"}
		msg.OrderID = engineResponse.Order.ID
		msg.Data = engineResponse

		ws.GetOrderConn(engineResponse.Order.ID).WriteJSON(msg)

		// for {
		t := time.NewTimer(10 * time.Second)
		ch := ws.GetOrderChannel(engineResponse.Order.ID)
		if ch == nil {
			// e.recoverOrders(engineResponse.MatchingOrders)
		} else {
			select {
			case rm := <-ch:
				if rm.MsgType == "trade_remorder_sign" {
					mb, err := json.Marshal(rm.Data)
					if err != nil {
						ws.GetOrderConn(engineResponse.Order.ID).WriteMessage(1, []byte(err.Error()))
					}
					var ersb *engine.EngineResponse
					err = json.Unmarshal(mb, &ersb)
					if err != nil {
						ws.GetOrderConn(engineResponse.Order.ID).WriteMessage(1, []byte(err.Error()))
					}
					if engineResponse.FillStatus == engine.PARTIAL {
						engineResponse.Order.OrderBook = &types.OrderSubDoc{Amount: ersb.RemainingOrder.Amount, Signature: ersb.RemainingOrder.Signature}
						// e.addOrder(order)
					}
				}
				t.Stop()
				break
			case <-t.C:
				fmt.Printf("\nTimeout\n")
				// e.recoverOrders(engineResponse.MatchingOrders)
				engineResponse.FillStatus = engine.ERROR
				engineResponse.Order.Status = types.ERROR
				engineResponse.Trades = nil
				engineResponse.RemainingOrder = nil
				engineResponse.MatchingOrders = nil
				t.Stop()
				break
			}
		}
	}
	ws.CloseOrderReadChannel(engineResponse.Order.ID)
	r.orderService.UpdateUsingEngineResponse(engineResponse)
	// TODO: send to operator for blockchain execution

	r.orderService.RelayUpdateOverSocket(engineResponse)
	return nil
}
