package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math/big"
	"time"

	"github.com/Proofsuite/amp-matching-engine/interfaces"
	"github.com/Proofsuite/amp-matching-engine/utils"
	"github.com/Proofsuite/amp-matching-engine/ws"
	"github.com/ethereum/go-ethereum/common"

	"gopkg.in/mgo.v2/bson"

	"github.com/Proofsuite/amp-matching-engine/rabbitmq"
	"github.com/Proofsuite/amp-matching-engine/types"
)

// OrderService
type OrderService struct {
	orderDao      interfaces.OrderDao
	pairDao       interfaces.PairDao
	accountDao    interfaces.AccountDao
	tradeDao      interfaces.TradeDao
	engine        interfaces.Engine
	validator     interfaces.ValidatorService
	broker        *rabbitmq.Connection
	orderChannels map[string]chan *types.WebsocketEvent
}

// NewOrderService returns a new instance of orderservice
func NewOrderService(
	orderDao interfaces.OrderDao,
	pairDao interfaces.PairDao,
	accountDao interfaces.AccountDao,
	tradeDao interfaces.TradeDao,
	engine interfaces.Engine,
	validator interfaces.ValidatorService,
	broker *rabbitmq.Connection,
) *OrderService {

	orderChannels := make(map[string]chan *types.WebsocketEvent)

	return &OrderService{
		orderDao,
		pairDao,
		accountDao,
		tradeDao,
		engine,
		validator,
		broker,
		orderChannels,
	}
}

// GetByID fetches the details of an order using order's mongo ID
func (s *OrderService) GetByID(id bson.ObjectId) (*types.Order, error) {
	return s.orderDao.GetByID(id)
}

// GetByUserAddress fetches all the orders placed by passed user address
func (s *OrderService) GetByUserAddress(addr common.Address) ([]*types.Order, error) {
	return s.orderDao.GetByUserAddress(addr)
}

// GetByHash fetches all trades corresponding to a trade hash
func (s *OrderService) GetByHash(hash common.Hash) (*types.Order, error) {
	return s.orderDao.GetByHash(hash)
}

// GetCurrentByUserAddress function fetches list of open/partial orders from order collection based on user address.
// Returns array of Order type struct
func (s *OrderService) GetCurrentByUserAddress(addr common.Address) ([]*types.Order, error) {
	return s.orderDao.GetCurrentByUserAddress(addr)
}

// GetHistoryByUserAddress function fetches list of orders which are not in open/partial order status
// from order collection based on user address.
// Returns array of Order type struct
func (s *OrderService) GetHistoryByUserAddress(addr common.Address) ([]*types.Order, error) {
	return s.orderDao.GetHistoryByUserAddress(addr)
}

// NewOrder validates if the passed order is valid or not based on user's available
// funds and order data.
// If valid: Order is inserted in DB with order status as new and order is publiched
// on rabbitmq queue for matching engine to process the order
func (s *OrderService) NewOrder(o *types.Order) error {

	acc, err := s.accountDao.GetByAddress(o.UserAddress)
	if err != nil {
		logger.Error(err)
		return err
	}

	if acc == nil {
		return errors.New("Account not found")
	}

	if acc.IsBlocked {
		return fmt.Errorf("Account %+v is blocked", acc)
	}

	if err := o.Validate(); err != nil {
		logger.Error(err)
		return err
	}

	ok, err := o.VerifySignature()
	if err != nil {
		logger.Error(err)
		return err
	}

	if !ok {
		return errors.New("Invalid signature")
	}

	p, err := s.pairDao.GetByBuySellTokenAddress(o.BuyToken, o.SellToken)
	if err != nil {
		logger.Error(err)
		return err
	}

	if p == nil {
		return errors.New("Pair not found")
	}

	// Fill token and pair data
	err = o.Process(p)
	if err != nil {
		logger.Error(err)
		return err
	}

	err = s.validator.ValidateBalance(o)
	if err != nil {
		logger.Error(err)
		return err
	}

	err = s.broker.PublishNewOrderMessage(o)
	if err != nil {
		logger.Error(err)
		return err
	}

	return nil
}

// CancelOrder handles the cancellation order requests.
// Only Orders which are OPEN or NEW i.e. Not yet filled/partially filled
// can be cancelled
func (s *OrderService) CancelOrder(oc *types.OrderCancel) error {
	o, err := s.orderDao.GetByHash(oc.OrderHash)
	if err != nil {
		logger.Error(err)
		return err
	}

	if o == nil {
		return fmt.Errorf("No order with this hash present")
	}

	if o.Status == "FILLED" || o.Status == "ERROR" {
		return fmt.Errorf("Cannot cancel order (Order status is %v)", o.Status)
	}

	err = s.broker.PublishCancelOrderMessage(o)
	if err != nil {
		logger.Error(err)
		return err
	}

	return nil
}

func (s *OrderService) handleOrderCancelled(res *types.EngineResponse) {
	ws.SendOrderMessage("ORDER_CANCELLED", res.Order.UserAddress, res.Order.Hash, res.Order)
	s.broadcastOrderBookUpdate([]*types.Order{res.Order})
	return
}

// HandleEngineResponse listens to messages incoming from the engine and handles websocket
// responses and database updates accordingly
func (s *OrderService) HandleEngineResponse(res *types.EngineResponse) error {
	switch res.Status {
	case "ERROR":
		s.handleEngineError(res)
	case "ORDER_ADDED":
		s.handleEngineOrderAdded(res)
	case "ORDER_FILLED":
		s.handleEngineOrderMatched(res)
	case "ORDER_PARTIALLY_FILLED":
		s.handleEngineOrderMatched(res)
	case "ORDER_CANCELLED":
		s.handleOrderCancelled(res)
	case "TRADES_CANCELLED":
		s.handleOrdersInvalidated(res)
	default:
		s.handleEngineUnknownMessage(res)
	}

	return nil
}

func (s *OrderService) HandleOperatorMessages(msg *types.OperatorMessage) error {
	switch msg.MessageType {
	case "TRADE_PENDING":
		s.handleOperatorTradePending(msg)
	case "TRADE_SUCCESS":
		s.handleOperatorTradeSuccess(msg)
	case "TRADE_ERROR":
		s.handleOperatorTradeError(msg)
	case "TRADE_INVALID":
		s.handleOperatorTradeError(msg)
	default:
		s.handleOperatorUnknownMessage(msg)
	}

	return nil
}

func (s *OrderService) handleOrdersInvalidated(res *types.EngineResponse) error {
	orders := res.InvalidatedOrders
	trades := res.CancelledTrades

	for _, o := range *orders {
		ws.SendOrderMessage("ORDER_INVALIDATED", o.UserAddress, o.Hash, o)
	}

	if orders != nil && len(*orders) != 0 {
		s.broadcastOrderBookUpdate(*orders)
	}

	if trades != nil && len(*trades) != 0 {
		s.broadcastTradeUpdate(*trades)
	}

	return nil
}

// handleEngineError returns an websocket error message to the client and recovers orders on the
// redis key/value store
func (s *OrderService) handleEngineError(res *types.EngineResponse) {
	o := res.Order
	ws.SendOrderMessage("ERROR", o.UserAddress, o.Hash, nil)
}

// handleEngineOrderAdded returns a websocket message informing the client that his order has been added
// to the orderbook (but currently not matched)
func (s *OrderService) handleEngineOrderAdded(res *types.EngineResponse) {
	o := res.Order
	ws.SendOrderMessage("ORDER_ADDED", o.UserAddress, o.Hash, o)
	s.broadcastOrderBookUpdate([]*types.Order{o})
}

// handleEngineOrderMatched returns a websocket message informing the client that his order has been added.
// The request signature message also signals the client to sign trades.
func (s *OrderService) handleEngineOrderMatched(res *types.EngineResponse) {
	//res.Order is the "taker" order
	o := res.Order

	orders := []*types.Order{o}
	trades := []*types.Trade{}
	matches := types.Matches{}
	invalidMatches := types.Matches{}

	//res.Matches is an array of (order, trade) pairs where each order is an "maker" order that is being matched
	for _, m := range *res.Matches {
		err := s.validator.ValidateBalance(m.Order)
		if err != nil {
			logger.Error(err)
			invalidMatches = append(invalidMatches, m)

		} else {
			orders = append(orders, m.Order)
			trades = append(trades, m.Trade)
			matches = append(matches, m)
		}
	}

	go s.handleSubmitSignatures(res)

	if len(matches) > 0 {
		ws.SendOrderMessage("REQUEST_SIGNATURE",
			o.UserAddress,
			o.Hash,
			types.SignaturePayload{res.Order, res.RemainingOrder, &matches},
		)
	}

	// if there are any invalid matches, the maker orders are at cause (since maker orders have been validated in the
	// newOrder() function. We remove the maker orders from the orderbook)
	if len(invalidMatches) > 0 {
		utils.PrintJSON("one more invalid match")
		err := s.broker.PublishInvalidateMakerOrdersMessage(invalidMatches)
		if err != nil {
			logger.Error(err)
		}
	}

	// we only update the orderbook with the current set of orders if there are no invalid matches.
	// If there are invalid matches, the corresponding maker orders will be removed and the taker order
	// amount filled will be updated as a result, and therefore does not represent the current state of the orderbook
	if len(invalidMatches) == 0 {
		s.broadcastOrderBookUpdate(orders)
	}

}

// handleSubmitSignatures wait for a submit signature message that provides the matching engine with orders
// that can be broadcast to the exchange smart contrct
func (s *OrderService) handleSubmitSignatures(res *types.EngineResponse) {
	addr := res.Order.UserAddress
	hashID := res.Order.Hash

	ch := s.CreateOrderChannel(hashID)
	defer s.DeleteOrderChannel(hashID)

	t := time.NewTimer(30 * time.Second)

	select {
	case msg := <-ch:
		if msg != nil && msg.Type == "SUBMIT_SIGNATURE" {
			bytes, err := json.Marshal(msg.Payload)
			if err != nil {
				logger.Error(err)
				ws.SendOrderMessage("ERROR", addr, hashID, err.Error())
			}

			p := &types.SignaturePayload{}
			err = json.Unmarshal(bytes, p)
			if err != nil {
				logger.Error(err)
				ws.SendOrderMessage("ERROR", addr, hashID, err.Error())
			}

			// handle remaining order ()
			ro := p.RemainingOrder
			if ro != nil {
				err := ro.ValidateComplete()
				if err != nil {
					ws.SendOrderMessage("ERROR", addr, hashID, err.Error())
					return
				}

				b, err := json.Marshal(ro)
				if err != nil {
					logger.Error(err)
					ws.SendOrderMessage("ERROR", addr, hashID, err.Error())
					return
				}

				//TODO do we need the Hash ID ?
				s.broker.PublishOrder(&rabbitmq.Message{Type: "NEW_ORDER", Data: b})
			}

			// trades
			if p.Matches != nil {
				matches := p.Matches
				trades := matches.Trades()
				taker := matches.Taker()

				err := matches.Validate()
				if err != nil {
					logger.Error(err)
					ws.SendOrderMessage("ERROR", taker, hashID, err)
					return
				}

				//TODO include this in the handleOrderMatched step
				err = s.tradeDao.Create(trades...)
				if err != nil {
					logger.Error(err)
					ws.SendOrderMessage("ERROR", taker, hashID, err)
					return
				}

				logger.Info("PUBLISHING TRADES")
				err = s.broker.PublishTrades(matches)
				if err != nil {
					logger.Error(err)
					ws.SendOrderMessage("ERROR", taker, hashID, err)
					return
				}
			}
		}
	case <-t.C:
		t.Stop()
		break
	}
}

// handleEngineUnknownMessage returns a websocket messsage in case the engine resonse is not recognized
func (s *OrderService) handleEngineUnknownMessage(res *types.EngineResponse) {
	log.Print("Receiving unknown engine message")
	utils.PrintJSON(res)
}

func (s *OrderService) handleOperatorUnknownMessage(msg *types.OperatorMessage) {
	log.Print("Receiving unknown message")
	utils.PrintJSON(msg)
}

func (s *OrderService) handleOperatorTradePending(msg *types.OperatorMessage) {
	matches := msg.Matches
	trades := matches.Trades()
	orders := matches.Orders()

	for _, t := range trades {
		err := s.tradeDao.UpdateTradeStatus(t.Hash, "PENDING")
		if err != nil {
			logger.Error(err)
		}

		t.Status = "PENDING"
	}

	taker := trades[0].Taker
	takerOrderHash := trades[0].TakerOrderHash
	ws.SendOrderMessage("ORDER_PENDING", taker, takerOrderHash, types.OrderPendingPayload{matches})

	for _, o := range orders {
		maker := o.UserAddress
		orderHash := o.Hash
		ws.SendOrderMessage("ORDER_PENDING", maker, orderHash, types.OrderPendingPayload{matches})
	}

	//TODO separate in different function or find more idiomatic code
	s.broadcastTradeUpdate(trades)
}

// handleOperatorTradeSuccess handles successfull trade messages from the orderbook. It updates
// the trade status in the database and
func (s *OrderService) handleOperatorTradeSuccess(msg *types.OperatorMessage) {
	matches := msg.Matches
	hashes := []common.Hash{}
	trades := matches.Trades()

	for _, t := range trades {
		hashes = append(hashes, t.Hash)
	}

	if len(hashes) == 0 {
		return
	}

	trades, err := s.tradeDao.UpdateTradeStatuses("SUCCESS", hashes...)
	if err != nil {
		logger.Error(err)
	}

	//TODO VERIFY that the status of the trades in the matches are modified to "SUCCESS"
	taker := trades[0].Taker
	takerHash := trades[0].TakerOrderHash
	ws.SendOrderMessage("ORDER_SUCCESS", taker, takerHash, types.OrderSuccessPayload{matches})

	for _, m := range *matches {
		maker := m.Order.UserAddress
		orderHash := m.Order.Hash
		//TODO Only send the corresponding order in the payload
		payload := types.Matches{m}
		ws.SendOrderMessage("ORDER_SUCCESS", maker, orderHash, types.OrderSuccessPayload{&payload})
	}

	s.broadcastTradeUpdate(trades)
}

// handleOperatorTradeError handles error messages from the operator (case where the blockchain tx was made
// but ended up failing. It updates the trade status in the db. None of the orders are reincluded in the redis
// orderbook.
func (s *OrderService) handleOperatorTradeError(msg *types.OperatorMessage) {
	matches := msg.Matches
	trades := matches.Trades()
	orders := matches.Orders()

	for _, t := range trades {
		err := s.tradeDao.UpdateTradeStatus(t.Hash, "ERROR")
		if err != nil {
			logger.Error(err)
		}

		t.Status = "ERROR"
	}

	taker := trades[0].Taker
	takerHash := trades[0].TakerOrderHash
	ws.SendOrderMessage("ORDER_ERROR", taker, takerHash, matches)

	for _, o := range orders {
		maker := o.UserAddress
		orderHash := o.Hash
		ws.SendOrderMessage("ORDER_ERROR", maker, orderHash, o)
	}

	s.broadcastTradeUpdate(trades)
}

func (s *OrderService) broadcastOrderBookUpdate(orders []*types.Order) {
	bids := []map[string]string{}
	asks := []map[string]string{}

	p, err := orders[0].Pair()
	if err != nil {
		logger.Error()
		return
	}

	for _, o := range orders {
		pp := o.PricePoint
		side := o.Side

		amount, err := s.orderDao.GetOrderBookPricePoint(p, pp, side)
		if err != nil {
			logger.Error(err)
		}

		// case where the amount at the pricepoint is equal to 0
		if amount == nil {
			amount = big.NewInt(0)
		}

		update := map[string]string{
			"pricepoint": pp.String(),
			"amount":     amount.String(),
		}

		if side == "BUY" {
			bids = append(bids, update)
		} else {
			asks = append(asks, update)
		}
	}

	id := utils.GetOrderBookChannelID(p.BaseTokenAddress, p.QuoteTokenAddress)
	ws.GetOrderBookSocket().BroadcastMessage(id, map[string][]map[string]string{
		"bids": bids,
		"asks": asks,
	})
}

func (s *OrderService) broadcastTradeUpdate(trades []*types.Trade) {
	p, err := trades[0].Pair()
	if err != nil {
		logger.Error(err)
		return
	}

	id := utils.GetTradeChannelID(p.BaseTokenAddress, p.QuoteTokenAddress)
	ws.GetTradeSocket().BroadcastMessage(id, trades)
}

func (s *OrderService) CreateOrderChannel(h common.Hash) chan *types.WebsocketEvent {
	if s.orderChannels == nil {
		s.orderChannels = make(map[string]chan *types.WebsocketEvent)
	}

	ch := make(chan *types.WebsocketEvent)
	if s.orderChannels[h.Hex()] == nil {
		s.orderChannels[h.Hex()] = ch
	}

	return ch
}

func (s *OrderService) GetOrderChannel(h common.Hash) chan *types.WebsocketEvent {
	if s.orderChannels[h.Hex()] == nil {
		return nil
	}

	if s.orderChannels[h.Hex()] == nil {
		return nil
	}

	return s.orderChannels[h.Hex()]
}

func (s *OrderService) DeleteOrderChannel(h common.Hash) {
	delete(s.orderChannels, h.Hex())
}
