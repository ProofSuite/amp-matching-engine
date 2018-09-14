package endpoints

import (
	"encoding/json"
	"log"

	"github.com/Proofsuite/amp-matching-engine/errors"
	"github.com/Proofsuite/amp-matching-engine/interfaces"
	"github.com/ethereum/go-ethereum/common"

	"github.com/Proofsuite/amp-matching-engine/types"
	"github.com/Proofsuite/amp-matching-engine/ws"
	"github.com/go-ozzo/ozzo-routing"
)

type orderEndpoint struct {
	orderService interfaces.OrderService
	engine       interfaces.Engine
}

// ServeOrderResource sets up the routing of order endpoints and the corresponding handlers.
func ServeOrderResource(
	rg *routing.RouteGroup,
	orderService interfaces.OrderService,
	engine interfaces.Engine,
) {
	e := &orderEndpoint{orderService, engine}
	rg.Get("/orders/<address>/history", e.queryHistory)
	rg.Get("/orders/<address>/current", e.queryCurrent)
	rg.Get("/orders/<address>", e.query)
	ws.RegisterChannel(ws.OrderChannel, e.ws)
}

func (e *orderEndpoint) query(c *routing.Context) error {
	addr := c.Param("address")
	if !common.IsHexAddress(addr) {
		return errors.NewAPIError(400, "Invalid Adrress", map[string]interface{}{})
	}

	address := common.HexToAddress(addr)
	orders, err := e.orderService.GetByUserAddress(address)
	if err != nil {
		return errors.NewAPIError(400, "Fetch Error", map[string]interface{}{})
	}

	return c.Write(orders)
}

func (e *orderEndpoint) queryCurrent(c *routing.Context) error {
	addr := c.Param("address")
	if !common.IsHexAddress(addr) {
		return errors.NewAPIError(400, "Invalid Adrress", map[string]interface{}{})
	}

	address := common.HexToAddress(addr)
	orders, err := e.orderService.GetCurrentByUserAddress(address)
	if err != nil {
		return errors.NewAPIError(400, "Fetch Error", map[string]interface{}{})
	}

	return c.Write(orders)
}

func (e *orderEndpoint) queryHistory(c *routing.Context) error {
	addr := c.Param("address")
	if !common.IsHexAddress(addr) {
		return errors.NewAPIError(400, "Invalid Adrress", map[string]interface{}{})
	}

	address := common.HexToAddress(addr)
	orders, err := e.orderService.GetHistoryByUserAddress(address)
	if err != nil {
		return errors.NewAPIError(400, "Fetch Error", map[string]interface{}{})
	}

	return c.Write(orders)
}

// ws function handles incoming websocket messages on the order channel
func (e *orderEndpoint) ws(input interface{}, conn *ws.Conn) {
	msg := &types.WebSocketPayload{}

	bytes, _ := json.Marshal(input)
	if err := json.Unmarshal(bytes, &msg); err != nil {
		logger.Error(err)
		ws.SendMessage(conn, ws.OrderChannel, "ERROR", err.Error())
	}

	switch msg.Type {
	case "NEW_ORDER":
		e.handleNewOrder(msg, conn)
	case "CANCEL_ORDER":
		e.handleCancelOrder(msg, conn)
	case "SUBMIT_SIGNATURE":
		e.handleSubmitSignatures(msg, conn)
	default:
		log.Print("Response with error")
	}
}

// handleSubmitSignatures handles NewTrade messages. New trade messages are transmitted to the corresponding order channel
// and received in the handleClientResponse.
func (e *orderEndpoint) handleSubmitSignatures(p *types.WebSocketPayload, conn *ws.Conn) {
	hash := common.HexToHash(p.Hash)
	ch := ws.GetOrderChannel(hash)

	if ch != nil {
		ch <- p
	}
}

// handleNewOrder handles NewOrder message. New order messages are transmitted to the order service after being unmarshalled
func (e *orderEndpoint) handleNewOrder(msg *types.WebSocketPayload, conn *ws.Conn) {
	ch := make(chan *types.WebSocketPayload)
	o := &types.Order{}

	bytes, err := json.Marshal(msg.Data)
	if err != nil {
		logger.Error(err)
		ws.SendMessage(conn, ws.OrderChannel, "ERROR", err.Error())
		return
	}

	err = json.Unmarshal(bytes, &o)
	if err != nil {
		logger.Error(err)
		ws.SendMessage(conn, ws.OrderChannel, "ERROR", err.Error())
		return
	}

	o.Hash = o.ComputeHash()
	ws.RegisterOrderConnection(o.Hash, &ws.OrderConnection{Conn: conn, ReadChannel: ch})
	ws.RegisterConnectionUnsubscribeHandler(conn, ws.OrderSocketUnsubscribeHandler(o.Hash))

	err = e.orderService.NewOrder(o)
	if err != nil {
		logger.Error(err)
		ws.SendMessage(conn, ws.OrderChannel, "ERROR", err.Error())
		return
	}
}

// handleCancelOrder handles CancelOrder message.
func (e *orderEndpoint) handleCancelOrder(p *types.WebSocketPayload, conn *ws.Conn) {
	bytes, err := json.Marshal(p.Data)
	oc := &types.OrderCancel{}

	err = oc.UnmarshalJSON(bytes)
	if err != nil {
		logger.Error(err)
		ws.SendMessage(conn, ws.OrderChannel, "ERROR", err.Error())
	}

	ws.RegisterOrderConnection(oc.Hash, &ws.OrderConnection{Conn: conn, Active: true})
	ws.RegisterConnectionUnsubscribeHandler(
		conn,
		ws.OrderSocketUnsubscribeHandler(oc.Hash),
	)

	err = e.orderService.CancelOrder(oc)
	if err != nil {
		logger.Error(err)
		ws.SendMessage(conn, ws.OrderChannel, "ERROR", err.Error())
		return
	}
}
