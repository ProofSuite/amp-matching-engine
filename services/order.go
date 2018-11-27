package services

import (
	"errors"
	"fmt"
	"log"
	"math/big"

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
func (s *OrderService) GetByUserAddress(addr common.Address, limit ...int) ([]*types.Order, error) {
	return s.orderDao.GetByUserAddress(addr, limit...)
}

// GetByHash fetches all trades corresponding to a trade hash
func (s *OrderService) GetByHash(hash common.Hash) (*types.Order, error) {
	return s.orderDao.GetByHash(hash)
}

func (s *OrderService) GetByHashes(hashes []common.Hash) ([]*types.Order, error) {
	return s.orderDao.GetByHashes(hashes)
}

// GetCurrentByUserAddress function fetches list of open/partial orders from order collection based on user address.
// Returns array of Order type struct
func (s *OrderService) GetCurrentByUserAddress(addr common.Address, limit ...int) ([]*types.Order, error) {
	return s.orderDao.GetCurrentByUserAddress(addr, limit...)
}

// GetHistoryByUserAddress function fetches list of orders which are not in open/partial order status
// from order collection based on user address.
// Returns array of Order type struct
func (s *OrderService) GetHistoryByUserAddress(addr common.Address, limit ...int) ([]*types.Order, error) {
	return s.orderDao.GetHistoryByUserAddress(addr, limit...)
}

// NewOrder validates if the passed order is valid or not based on user's available
// funds and order data.
// If valid: Order is inserted in DB with order status as new and order is publiched
// on rabbitmq queue for matching engine to process the order
func (s *OrderService) NewOrder(o *types.Order) error {
	if err := o.Validate(); err != nil {
		logger.Error(err)
		return err
	}

	ok, err := o.VerifySignature()
	if err != nil {
		logger.Error(err)
	}

	if !ok {
		return errors.New("Invalid Signature")
	}

	p, err := s.pairDao.GetByTokenAddress(o.BaseToken, o.QuoteToken)
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
		return errors.New("No order with corresponding hash")
	}

	if o.Status == "FILLED" || o.Status == "ERROR" || o.Status == "CANCEL" {
		return fmt.Errorf("Cannot cancel order. Status is %v", o.Status)
	}

	err = s.broker.PublishCancelOrderMessage(o)
	if err != nil {
		logger.Error(err)
		return err
	}

	return nil
}

func (s *OrderService) handleOrderCancelled(res *types.EngineResponse) {
	ws.SendOrderMessage("ORDER_CANCELLED", res.Order.UserAddress, res.Order)
	s.broadcastOrderBookUpdate([]*types.Order{res.Order})
	s.broadcastRawOrderBookUpdate([]*types.Order{res.Order})
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
		ws.SendOrderMessage("ORDER_INVALIDATED", o.UserAddress, o)
	}

	if orders != nil && len(*orders) != 0 {
		s.broadcastOrderBookUpdate(*orders)
	}

	if orders != nil && len(*orders) != 0 {
		s.broadcastRawOrderBookUpdate(*orders)
	}

	if trades != nil && len(*trades) != 0 {
		s.broadcastTradeUpdate(*trades)
	}

	return nil
}

// handleEngineError returns an websocket error message to the client
func (s *OrderService) handleEngineError(res *types.EngineResponse) {
	o := res.Order
	ws.SendOrderMessage("ERROR", o.UserAddress, nil)
}

// handleEngineOrderAdded returns a websocket message informing the client that his order has been added
// to the orderbook (but currently not matched)
func (s *OrderService) handleEngineOrderAdded(res *types.EngineResponse) {
	o := res.Order
	ws.SendOrderMessage("ORDER_ADDED", o.UserAddress, o)

	s.broadcastOrderBookUpdate([]*types.Order{o})
	s.broadcastRawOrderBookUpdate([]*types.Order{o})
}

// handleEngineOrderMatched returns a websocket message informing the client that his order has been added.
// The request signature message also signals the client to sign trades.
func (s *OrderService) handleEngineOrderMatched(res *types.EngineResponse) {
	o := res.Order //res.Order is the "taker" order
	taker := o.UserAddress
	matches := *res.Matches
	ws.SendOrderMessage("ORDER_MATCHED", taker, types.OrderMatchedPayload{&matches})

	orders := []*types.Order{o}
	validMatches := types.Matches{TakerOrder: o}
	invalidMatches := types.Matches{TakerOrder: o}

	//res.Matches is an array of (order, trade) pairs where each order is an "maker" order that is being matched
	for i, _ := range matches.Trades {
		err := s.validator.ValidateBalance(matches.MakerOrders[i])
		if err != nil {
			logger.Error(err)
			invalidMatches.AppendMatch(matches.MakerOrders[i], matches.Trades[i])

		} else {
			validMatches.AppendMatch(matches.MakerOrders[i], matches.Trades[i])
			orders = append(orders, matches.MakerOrders[i])
		}
	}

	// if there are any invalid matches, the maker orders are at cause (since maker orders have been validated in the
	// newOrder() function. We remove the maker orders from the orderbook)
	if invalidMatches.Length() > 0 {
		err := s.broker.PublishInvalidateMakerOrdersMessage(invalidMatches)
		if err != nil {
			logger.Error(err)
		}
	}

	err := s.tradeDao.Create(validMatches.Trades...)
	if err != nil {
		logger.Error(err)
		ws.SendOrderMessage("ERROR", taker, err)
		return
	}

	err = s.broker.PublishTrades(&validMatches)
	if err != nil {
		logger.Error(err)
		ws.SendOrderMessage("ERROR", taker, err)
		return
	}

	// we only update the orderbook with the current set of orders if there are no invalid matches.
	// If there are invalid matches, the corresponding maker orders will be removed and the taker order
	// amount filled will be updated as a result, and therefore does not represent the current state of the orderbook
	if invalidMatches.Length() == 0 {
		s.broadcastOrderBookUpdate(orders)
		s.broadcastRawOrderBookUpdate(orders)
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
	trades := matches.Trades
	orders := matches.MakerOrders

	taker := trades[0].Taker
	ws.SendOrderMessage("ORDER_PENDING", taker, types.OrderPendingPayload{matches})

	for _, o := range orders {
		maker := o.UserAddress
		ws.SendOrderMessage("ORDER_PENDING", maker, types.OrderPendingPayload{matches})
	}

	s.broadcastTradeUpdate(trades)
}

// handleOperatorTradeSuccess handles successfull trade messages from the orderbook. It updates
// the trade status in the database and
func (s *OrderService) handleOperatorTradeSuccess(msg *types.OperatorMessage) {
	matches := msg.Matches
	hashes := []common.Hash{}
	trades := matches.Trades

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

	// Send ORDER_SUCCESS message to order takers
	taker := trades[0].Taker
	ws.SendOrderMessage("ORDER_SUCCESS", taker, types.OrderSuccessPayload{matches})

	// Send ORDER_SUCCESS message to order makers
	for i, _ := range trades {
		match := matches.NthMatch(i)
		maker := match.MakerOrders[0].UserAddress
		ws.SendOrderMessage("ORDER_SUCCESS", maker, types.OrderSuccessPayload{match})
	}

	s.broadcastTradeUpdate(trades)
}

// handleOperatorTradeError handles error messages from the operator (case where the blockchain tx was made
// but ended up failing. It updates the trade status in the db. None of the orders are reincluded in the engine
// orderbook.
func (s *OrderService) handleOperatorTradeError(msg *types.OperatorMessage) {
	matches := msg.Matches
	trades := matches.Trades
	orders := matches.MakerOrders

	errType := msg.ErrorType
	if errType != "" {
		logger.Error("")
	}

	for _, t := range trades {
		err := s.tradeDao.UpdateTradeStatus(t.Hash, "ERROR")
		if err != nil {
			logger.Error(err)
		}

		t.Status = "ERROR"
	}

	taker := trades[0].Taker
	ws.SendOrderMessage("ORDER_ERROR", taker, matches)

	for _, o := range orders {
		maker := o.UserAddress
		ws.SendOrderMessage("ORDER_ERROR", maker, o)
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

	logger.Info("BIDS", utils.JSON(bids))
	logger.Info("ASKS", utils.JSON(asks))

	id := utils.GetOrderBookChannelID(p.BaseTokenAddress, p.QuoteTokenAddress)
	ws.GetOrderBookSocket().BroadcastMessage(id, map[string]interface{}{
		"pairName": orders[0].PairName,
		"bids":     bids,
		"asks":     asks,
	})
}

func (s *OrderService) broadcastRawOrderBookUpdate(orders []*types.Order) {
	p, err := orders[0].Pair()
	if err != nil {
		logger.Error(err)
		return
	}

	id := utils.GetOrderBookChannelID(p.BaseTokenAddress, p.QuoteTokenAddress)
	ws.GetRawOrderBookSocket().BroadcastMessage(id, orders)
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
