package endpoints

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/Proofsuite/matching-engine/errors"
	"github.com/ethereum/go-ethereum/common"

	"github.com/Proofsuite/matching-engine/engine"
	"github.com/Proofsuite/matching-engine/services"
	"github.com/Proofsuite/matching-engine/types"
	"github.com/Proofsuite/matching-engine/ws"
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
	http.HandleFunc("/orders/ws", r.ws)
	http.HandleFunc("/orders/book/<pair>", r.ws)
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

func (r *orderEndpoint) ws(w http.ResponseWriter, req *http.Request) {
	conn, err := upgrader.Upgrade(w, req, nil)
	if err != nil {
		log.Println("==>" + err.Error())
		return
	}
	// go func() {
	ch := make(chan *types.WsMsg)
	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			log.Println("<==>" + err.Error())
			conn.Close()
			return
		}
		var msg *types.WsMsg
		if err := json.Unmarshal(p, &msg); err != nil {
			log.Println("unmarshal to wsmsg <==>" + err.Error())
			conn.Close()

			return
		}
		if msg.MsgType == "new_order" {
			oab, err := json.Marshal(msg.Data)

			var model types.OrderRequest
			if err := json.Unmarshal(oab, &model); err != nil {
				conn.WriteMessage(messageType, []byte(err.Error()))
				conn.Close()

				return
			}
			if err := model.Validate(); err != nil {
				conn.WriteMessage(messageType, []byte(err.Error()))
				conn.Close()

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
				conn.Close()
				return
			}
			oab, _ = json.Marshal(order)
			conn.WriteMessage(messageType, oab)
			if ws.Connections == nil {
				ws.Connections = make(map[string]*ws.Ws)
			}
			conn.SetCloseHandler(ws.OrderSocketCloseHandler(order.ID))
			ws.Connections[order.ID.Hex()] = &ws.Ws{Conn: conn, ReadChannel: ch}
		} else {
			ws.Connections[msg.OrderID.Hex()].ReadChannel <- msg
		}
	}
}

func (r *orderEndpoint) engineResponse(er *engine.EngineResponse) error {
	// b, _ := json.Marshal(er)
	// fmt.Printf("\n======> \n%s\n <======\n", b)
	r.orderService.UpdateUsingEngineResponse(er)
	// TODO: send to operator for blockchain execution

	r.orderService.RelayUpdateOverSocket(er)
	return nil
}
