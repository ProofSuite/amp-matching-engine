package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/Proofsuite/amp-matching-engine/app"
	"github.com/Proofsuite/amp-matching-engine/interfaces"
	"github.com/Proofsuite/amp-matching-engine/utils"
	"github.com/Proofsuite/amp-matching-engine/utils/math"
	"github.com/Proofsuite/amp-matching-engine/ws"
	"github.com/ethereum/go-ethereum/common"

	"gopkg.in/mgo.v2/bson"

	"github.com/Proofsuite/amp-matching-engine/rabbitmq"
	"github.com/Proofsuite/amp-matching-engine/types"
)

// OrderService struct with daos required, responsible for communicating with daos.
// OrderService functions are responsible for interacting with daos and implements business logics.
type OrderService struct {
	orderDao         interfaces.OrderDao
	pairDao          interfaces.PairDao
	accountDao       interfaces.AccountDao
	tradeDao         interfaces.TradeDao
	engine           interfaces.Engine
	ethereumProvider interfaces.EthereumProvider
	broker           *rabbitmq.Connection
}

// NewOrderService returns a new instance of orderservice
func NewOrderService(
	orderDao interfaces.OrderDao,
	pairDao interfaces.PairDao,
	accountDao interfaces.AccountDao,
	tradeDao interfaces.TradeDao,
	engine interfaces.Engine,
	ethereumProvider interfaces.EthereumProvider,
	broker *rabbitmq.Connection,
) *OrderService {
	return &OrderService{
		orderDao,
		pairDao,
		accountDao,
		tradeDao,
		engine,
		ethereumProvider,
		broker,
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
	// Validate if the address is not blacklisted
	acc, err := s.accountDao.GetByAddress(o.UserAddress)
	if err != nil {
		logger.Error(err)
		return err
	}

	if acc.IsBlocked {
		return fmt.Errorf("Address: %+v isBlocked", acc)
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

	// fee balance validation
	wethAddress := common.HexToAddress(app.Config.Ethereum["weth_address"])
	exchangeAddress := common.HexToAddress(app.Config.Ethereum["exchange_address"])
	balanceRecord, err := s.accountDao.GetTokenBalances(o.UserAddress)
	if err != nil {
		logger.Error(err)
		return err
	}

	wethBalance, err := s.ethereumProvider.BalanceOf(o.UserAddress, wethAddress)
	if err != nil {
		logger.Error(err)
		return err
	}

	wethAllowance, err := s.ethereumProvider.Allowance(o.UserAddress, exchangeAddress, wethAddress)
	if err != nil {
		logger.Error(err)
		return err
	}

	wethLockedBalance, err := s.orderDao.GetUserLockedBalance(o.UserAddress, wethAddress)
	if err != nil {
		logger.Error(err)
		return err
	}

	sellTokenBalance, err := s.ethereumProvider.BalanceOf(o.UserAddress, o.SellToken)
	if err != nil {
		logger.Error(err)
		return err
	}

	sellTokenAllowance, err := s.ethereumProvider.Allowance(o.UserAddress, exchangeAddress, o.SellToken)
	if err != nil {
		logger.Error(err)
		return err
	}

	sellTokenLockedBalance, err := s.orderDao.GetUserLockedBalance(o.UserAddress, o.SellToken)
	if err != nil {
		logger.Error(err)
		return err
	}

	fee := math.Max(o.MakeFee, o.TakeFee)
	availableWethBalance := math.Sub(wethBalance, wethLockedBalance)
	availableSellTokenBalance := math.Sub(sellTokenBalance, sellTokenLockedBalance)

	if availableWethBalance.Cmp(fee) == -1 {
		return errors.New("Insufficient WETH Balance")
	}

	if wethAllowance.Cmp(fee) == -1 {
		return errors.New("Insufficient WETH Balance")
	}

	if availableSellTokenBalance.Cmp(o.SellAmount) != 1 {
		return errors.New("Insufficient Balance")
	}

	if sellTokenAllowance.Cmp(o.SellAmount) != 1 {
		return errors.New("Insufficient Allowance")
	}

	sellTokenBalanceRecord := balanceRecord[o.SellToken]
	sellTokenBalanceRecord.Balance.Set(sellTokenBalance)
	sellTokenBalanceRecord.Allowance.Set(sellTokenAllowance)
	wethTokenBalanceRecord := balanceRecord[wethAddress]
	wethTokenBalanceRecord.Balance.Set(wethBalance)
	wethTokenBalanceRecord.Allowance.Set(wethAllowance)

	err = s.accountDao.UpdateTokenBalance(o.UserAddress, wethAddress, wethTokenBalanceRecord)
	if err != nil {
		logger.Error(err)
		return err
	}

	err = s.accountDao.UpdateTokenBalance(o.UserAddress, o.SellToken, sellTokenBalanceRecord)
	if err != nil {
		logger.Error(err)
		return err
	}

	if err = s.orderDao.Create(o); err != nil {
		logger.Error(err)
		return err
	}

	bytes, err := json.Marshal(o)
	if err != nil {
		logger.Error(err)
		return err
	}

	s.broker.PublishOrder(&rabbitmq.Message{Type: "NEW_ORDER", Data: bytes})
	return nil
}

// CancelOrder handles the cancellation order requests.
// Only Orders which are OPEN or NEW i.e. Not yet filled/partially filled
// can be cancelled
func (s *OrderService) CancelOrder(oc *types.OrderCancel) error {
	dbOrder, err := s.orderDao.GetByHash(oc.OrderHash)
	if err != nil {
		logger.Error(err)
		return err
	}

	if dbOrder == nil {
		return fmt.Errorf("No order with this hash present")
	}

	_, err = json.Marshal(dbOrder)
	if err != nil {
		logger.Error(err)
		return err
	}

	if dbOrder.Status == "OPEN" || dbOrder.Status == "NEW" {
		res, err := s.engine.CancelOrder(dbOrder)
		if err != nil {
			logger.Error(err)
			return err
		}

		err = s.orderDao.UpdateOrderStatus(res.Order.Hash, "CANCELLED")
		if err != nil {
			logger.Error(err)
		}

		// err = s.orderDao.UpdateByHash(res.Order.Hash, res.Order)
		// if err != nil {
		// 	logger.Error(err)
		// 	return err
		// }

		// err = s.ReverseBalances(res.Order)
		// if err != nil {
		// 	logger.Error(err)
		// 	return err
		// }

		ws.SendOrderMessage("ORDER_CANCELLED", res.Order.Hash, res.Order)
		s.RelayUpdateOverSocket(res)
		return nil
	}

	return fmt.Errorf("Cannot cancel the order")
}

// HandleEngineResponse listens to messages incoming from the engine and handles websocket
// responses and database updates accordingly
func (s *OrderService) HandleEngineResponse(res *types.EngineResponse) error {
	switch res.Status {
	case "ERROR":
		s.handleEngineError(res)
	case "NOMATCH":
		s.handleEngineOrderAdded(res)
	case "FULL":
		s.handleEngineOrderMatched(res)
	case "PARTIAL":
		s.handleEngineOrderMatched(res)
	default:
		s.handleEngineUnknownMessage(res)
	}

	s.RelayUpdateOverSocket(res)
	// ws.CloseOrderReadChannel(res.Order.Hash)
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

// handleEngineError returns an websocket error message to the client and recovers orders on the
// redis key/value store
func (s *OrderService) handleEngineError(res *types.EngineResponse) {
	// err := s.orderDao.UpdateByHash(res.Order.Hash, res.Order)
	// if err != nil {
	// 	logger.Error(err)
	// }

	err := s.orderDao.UpdateOrderStatus(res.Order.Hash, "ERROR")
	if err != nil {
		logger.Error(err)
	}

	ws.SendOrderMessage("ERROR", res.Order.Hash, nil)
}

// handleEngineOrderAdded returns a websocket message informing the client that his order has been added
// to the orderbook (but currently not matched)
func (s *OrderService) handleEngineOrderAdded(res *types.EngineResponse) {
	ws.SendOrderMessage("ORDER_ADDED", res.Order.Hash, res.Order)
}

// handleEngineOrderMatched returns a websocket message informing the client that his order has been added.
// The request signature message also signals the client to sign trades.
func (s *OrderService) handleEngineOrderMatched(res *types.EngineResponse) {
	err := s.orderDao.UpdateByHash(res.Order.Hash, res.Order)
	if err != nil {
		logger.Error(err)
	}

	for _, m := range res.Matches {
		err := s.orderDao.UpdateByHash(m.Order.Hash, m.Order)
		if err != nil {
			logger.Error(err)
		}
	}

	go s.handleSubmitSignatures(res)
	ws.SendOrderMessage("REQUEST_SIGNATURE", res.Order.Hash, types.SignaturePayload{res.RemainingOrder, res.Matches})
}

// handleSubmitSignatures wait for a submit signature message that provides the matching engine with orders
// that can be broadcast to the exchange smart contrct
func (s *OrderService) handleSubmitSignatures(res *types.EngineResponse) {
	ch := ws.GetOrderChannel(res.Order.Hash)
	t := time.NewTimer(30 * time.Second)

	select {
	case msg := <-ch:
		if msg != nil && msg.Type == "SUBMIT_SIGNATURE" {
			bytes, err := json.Marshal(msg.Data)
			if err != nil {
				logger.Error(err)
				s.Rollback(res)
				ws.SendOrderMessage("ERROR", res.Order.Hash, err)
			}

			data := &types.SignaturePayload{}
			err = json.Unmarshal(bytes, data)
			if err != nil {
				logger.Error(err)
				s.Rollback(res)
				ws.SendOrderMessage("ERROR", res.Order.Hash, err)
			}

			// remaining order
			if data.Order != nil {
				err := s.orderDao.Create(data.Order)
				if err != nil {
					//TODO consider if we should going on with execution or not
					logger.Error(err)
					s.Rollback(res)
					ws.SendOrderMessage("ERROR", res.Order.Hash, err)
				}

				bytes, err := json.Marshal(data.Order)
				if err != nil {
					logger.Error(err)
					//TODO not sure whether rolling back is good here
					s.Rollback(res)
					ws.SendOrderMessage("ERROR", res.Order.Hash, err)
				}

				s.broker.PublishOrder(&rabbitmq.Message{Type: "ADD_ORDER", Data: bytes})
			}

			if data.Matches != nil {
				trades := []*types.Trade{}
				for _, m := range data.Matches {
					trades = append(trades, m.Trade)
				}

				//TODO include this in the handleOrderMatched step
				err := s.tradeDao.Create(trades...)
				if err != nil {
					logger.Error(err)
				}

				_, err = json.Marshal(res.Order)
				if err != nil {
					logger.Error(err)
					s.Rollback(res)
					ws.SendOrderMessage("ERROR", res.Order.Hash, err)
				}

				for _, m := range data.Matches {
					err := s.broker.PublishTrade(m.Order, m.Trade)
					if err != nil {
						logger.Error(err)
						s.Rollback(res)
						ws.SendOrderMessage("ERROR", res.Order.Hash, err)
					}
				}
			}
		}
	case <-t.C:
		s.Rollback(res)
		t.Stop()
		break
	}
}

// handleEngineUnknownMessage returns a websocket messsage in case the engine resonse is not recognized
func (s *OrderService) handleEngineUnknownMessage(res *types.EngineResponse) {
	s.Rollback(res)
	ws.SendOrderMessage("ERROR", res.Order.Hash, nil)
}

func (s *OrderService) handleOperatorUnknownMessage(msg *types.OperatorMessage) {
	// s.Rollback(res)
	ws.SendOrderMessage("ERROR", msg.Order.Hash, nil)
}

func (s *OrderService) handleOperatorTradePending(msg *types.OperatorMessage) {
	t := msg.Trade

	err := s.tradeDao.UpdateTradeStatus(t.Hash, "ORDER_PENDING")
	if err != nil {
		logger.Error(err)
	}

	ws.SendOrderMessage("ORDER_PENDING", t.OrderHash, t)
	ws.SendOrderMessage("ORDER_PENDING", t.TakerOrderHash, t)
}

// handleTradeMakerInvalid handles the case where a "MAKER_INVALID" message is received from the
// operator. It reinclues the TAKER order in the db and in the redis orderbook and invalidates the
// MAKER order
func (s *OrderService) handleTradeMakerInvalid(msg *types.OperatorMessage) {
	t := msg.Trade

	err := s.tradeDao.UpdateTradeStatus(t.Hash, "INVALID")
	if err != nil {
		logger.Error(err)
	}

	err = s.tradeDao.UpdateTradeStatus(t.TakerOrderHash, "INVALID")
	if err != nil {
		logger.Error(err)
	}

	err = s.orderDao.UpdateOrderFilledAmount(t.OrderHash, math.Neg(t.Amount))
	if err != nil {
		logger.Error(err)
	}

	takerOrder, err := s.orderDao.GetByHash(t.TakerOrderHash)
	if err != nil {
		logger.Error(err)
	}

	op := &types.OrderTradePair{takerOrder, t}
	err = s.engine.RecoverOrders([]*types.OrderTradePair{op})
	if err != nil {
		logger.Error(err)
	}

	//TODO decide whether we should also send a message to the taker
	ws.SendOrderMessage("ORDER_INVALID", t.OrderHash, t)
}

// handleTradeMakerInvalid handles the case where a "TAKER_INVALID" message is received from the
// operator. It reinclues the MAKER order in the db and in the redis orderbook and invalidates the
// TAKER order
func (s *OrderService) handleTradeTakerInvalid(msg *types.OperatorMessage) {
	t := msg.Trade

	err := s.tradeDao.UpdateTradeStatus(t.Hash, "INVALID")
	if err != nil {
		logger.Error(err)
	}

	err = s.orderDao.UpdateOrderStatus(t.TakerOrderHash, "INVALID")
	if err != nil {
		logger.Error(err)
	}

	//we reinclude the amount "lost" of the MAKER ORDER due to this failed trade back in the mongo record
	err = s.orderDao.UpdateOrderFilledAmount(t.OrderHash, math.Neg(t.Amount))
	if err != nil {
		logger.Error(err)
	}

	makerOrder, err := s.orderDao.GetByHash(t.OrderHash)
	if err != nil {
		logger.Error(err)
	}

	// we recover and include the maker order in the redis orderbook again
	//TODO only the trade amount should be needed and not the full trade
	op := &types.OrderTradePair{makerOrder, t}
	err = s.engine.RecoverOrders([]*types.OrderTradePair{op})
	if err != nil {
		logger.Error(err)
	}

	ws.SendOrderMessage("ORDER_INVALID", t.TakerOrderHash, t)
	//TODO decide whether we should also send a message to the maker. In
	//TODO theory might has well not take the trouble since this will not happen
	//TODO often and he will likely not know
}

// handleOperatorTradeSuccess handles successfull trade messages from the orderbook. It updates
// the trade status in the database and
func (s *OrderService) handleOperatorTradeSuccess(msg *types.OperatorMessage) {
	t := msg.Trade
	err := s.tradeDao.UpdateTradeStatus(t.Hash, "SUCCESS")
	if err != nil {
		logger.Error(err)
	}

	ws.SendOrderMessage("ORDER_SUCCESS", t.OrderHash, t)
	ws.SendOrderMessage("ORDER_SUCCESS", t.TakerOrderHash, t)
}

// handleOperatorTradeError handles error messages from the operator (case where the blockchain tx was made
// but ended up failing. It updates the trade status in the db. None of the orders are reincluded in the redis
// orderbook.
func (s *OrderService) handleOperatorTradeError(msg *types.OperatorMessage) {
	t := msg.Trade
	ws.SendOrderMessage("ORDER_ERROR", t.OrderHash, t)
	ws.SendOrderMessage("ORDER_ERROR", t.TakerOrderHash, t)

	err := s.tradeDao.UpdateTradeStatus(t.Hash, "ERROR")
	if err != nil {
		logger.Error(err)
	}
}

func (s *OrderService) Rollback(res *types.EngineResponse) *types.EngineResponse {
	if res.RemainingOrder != nil {
		err := s.orderDao.UpdateOrderStatus(res.RemainingOrder.Hash, "ERROR")
		if err != nil {
			logger.Error(err)
		}
	}

	//TODO what do we do with remaining order ?
	if len(res.Matches) > 0 {
		for _, ot := range res.Matches {
			t := ot.Trade

			err := s.orderDao.UpdateOrderFilledAmount(t.OrderHash, math.Neg(t.Amount))
			if err != nil {
				logger.Error(err)
			}

			err = s.orderDao.UpdateOrderStatus(t.TakerOrderHash, "ERROR")
			if err != nil {
				logger.Error(err)
			}

			err = s.tradeDao.UpdateTradeStatus(t.Hash, "ERROR")
			if err != nil {
				logger.Error(err)
			}
		}

		//TODO should we simply delete the orders from the orderbook
		err := s.engine.RecoverOrders(res.Matches)
		if err != nil {
			logger.Error(err)
		}
	}

	res.Status = "ERROR"
	res.Order.Status = "ERROR"
	res.Matches = nil
	res.RemainingOrder = nil
	return res
}

func (s *OrderService) RollbackOrder(o *types.Order) (err error) {
	err = s.orderDao.UpdateOrderStatus(o.Hash, "ERROR")
	if err != nil {
		logger.Error(err)
	}

	err = s.engine.DeleteOrder(o)
	if err != nil {
		logger.Error(err)
	}

	return err
}

func (s *OrderService) RollbackTrade(o *types.Order, t *types.Trade) (err error) {
	err = s.tradeDao.UpdateTradeStatus(t.Hash, "ERROR")
	if err != nil {
		logger.Error(err)
	}

	//TODO check that this also updates the order status
	err = s.orderDao.UpdateOrderFilledAmount(t.OrderHash, math.Neg(t.Amount))
	if err != nil {
		logger.Error(err)
	}

	err = s.engine.RecoverOrders([]*types.OrderTradePair{{o, t}})
	if err != nil {
		logger.Error(err)
	}

	return err
}

func (s *OrderService) CancelTrades(trades []*types.Trade) error {
	orderHashes := []common.Hash{}
	amounts := []*big.Int{}

	for _, t := range trades {
		orderHashes = append(orderHashes, t.OrderHash)
		amounts = append(amounts, t.Amount)
	}

	orders, err := s.orderDao.GetByHashes(orderHashes)
	if err != nil {
		logger.Error(err)
		return err
	}

	err = s.engine.CancelTrades(orders, amounts)
	if err != nil {
		logger.Error(err)
		return err
	}

	return nil
}

// RelayUpdateOverSocket is resonsible for notifying listening clients about new order/trade addition/deletion
func (s *OrderService) RelayUpdateOverSocket(res *types.EngineResponse) {
	// broadcast order's latest state
	s.RelayOrderUpdate(res)
	s.RelayTradeUpdate(res)
}

// RelayOrderUpdate is resonsible for notifying listening clients about new order addition/deletion
func (s *OrderService) RelayOrderUpdate(res *types.EngineResponse) {
	// broadcast order's latest state
	go broadcastOrderBookUpdate(res.Order.BaseToken, res.Order.QuoteToken, getOrderBookPayload(res))
	go broadcastRawOrderBookUpdate(res.Order)
}

// RelayTradeUpdate is resonsible for notifying listening clients about new trades
func (s *OrderService) RelayTradeUpdate(res *types.EngineResponse) {
	if len(res.Matches) == 0 {
		return
	}

	trades := []*types.Trade{}
	for _, m := range res.Matches {
		trades = append(trades, m.Trade)
	}

	// broadcast trades
	go broadcastTradeUpdate(trades)
}

func broadcastOrderBookUpdate(baseToken, quoteToken common.Address, data interface{}) {
	cid := utils.GetOrderBookChannelID(baseToken, quoteToken)
	ws.GetLiteOrderBookSocket().BroadcastMessage(cid, data)
}

func broadcastRawOrderBookUpdate(order *types.Order) {
	cid := utils.GetOrderBookChannelID(order.BaseToken, order.QuoteToken)
	ws.GetRawOrderBookSocket().BroadcastMessage(cid, order)
}

func broadcastTradeUpdate(trades []*types.Trade) {
	cid := utils.GetTradeChannelID(trades[0].BaseToken, trades[0].QuoteToken)
	ws.GetTradeSocket().BroadcastMessage(cid, trades)
}

func getOrderBookPayload(res *types.EngineResponse) interface{} {
	orderSide := make(map[string]string)
	matchSide := make([]map[string]string, 0)

	matchSideMap := make(map[string]*big.Int)

	if math.Sub(res.Order.Amount, res.Order.FilledAmount).Cmp(big.NewInt(0)) != 0 {
		orderSide["price"] = res.Order.PricePoint.String()
		orderSide["amount"] = math.Sub(res.Order.Amount, res.Order.FilledAmount).String()
	}

	if len(res.Matches) > 0 {
		for _, mo := range res.Matches {
			pp := mo.Order.PricePoint.String()
			if matchSideMap[pp] == nil {
				matchSideMap[pp] = big.NewInt(0)
			}

			matchSideMap[pp] = math.Add(matchSideMap[pp], mo.Trade.Amount)
		}
	}

	for price, amount := range matchSideMap {
		temp := map[string]string{
			"price":  price,
			"amount": math.Neg(amount).String(),
		}
		matchSide = append(matchSide, temp)
	}

	var response map[string]interface{}
	if res.Order.Side == "SELL" {
		response = map[string]interface{}{
			"asks": []map[string]string{orderSide},
			"bids": matchSide,
		}
	} else {
		response = map[string]interface{}{
			"asks": matchSide,
			"bids": []map[string]string{orderSide},
		}
	}

	return response
}

// transferAmount is used to transfer amount from seller to buyer
// it removes the lockedAmount of one token and adds confirmed amount for another token
// based on the type of order i.e. buy/sell
// func (s *OrderService) transferAmount(o *types.Order, filledAmount *big.Int) {
// 	tokenBalances, err := s.accountDao.GetTokenBalances(o.UserAddress)
// 	if err != nil {
// 		logger.Error(err)
// 	}

// 	if o.Side == "BUY" {
// 		sellBalance := tokenBalances[o.QuoteToken]
// 		sellBalance.LockedBalance = sellBalance.LockedBalance.Sub(sellBalance.LockedBalance, filledAmount)

// 		err = s.accountDao.UpdateTokenBalance(o.UserAddress, o.QuoteToken, sellBalance)
// 		if err != nil {
// 			logger.Error(err)
// 		}

// 		buyBalance := tokenBalances[o.BaseToken]
// 		buyBalance.Balance = buyBalance.Balance.Add(buyBalance.Balance, filledAmount)
// 		err = s.accountDao.UpdateTokenBalance(o.UserAddress, o.BaseToken, buyBalance)
// 		if err != nil {
// 			logger.Error(err)
// 		}
// 	}

// 	if o.Side == "SELL" {
// 		buyBalance := tokenBalances[o.BaseToken]
// 		buyBalance.LockedBalance = buyBalance.LockedBalance.Sub(buyBalance.LockedBalance, filledAmount)
// 		err = s.accountDao.UpdateTokenBalance(o.UserAddress, o.BaseToken, buyBalance)
// 		if err != nil {
// 			logger.Error(err)
// 		}

// 		sellBalance := tokenBalances[o.QuoteToken]
// 		sellBalance.Balance = sellBalance.Balance.Add(sellBalance.Balance, filledAmount)
// 		err = s.accountDao.UpdateTokenBalance(o.UserAddress, o.BaseToken, sellBalance)
// 		if err != nil {
// 			logger.Error(err)
// 		}
// 	}
// }

// UnlockBalances unlocks balances in case an order is canceled for examples. It unlocks both
// the token sell balance as well as the weth balance for fees.
// func (s *OrderService) ReverseBalances(orders ...*types.Order) error {
// 	for _, o := range orders {
// 		balanceRecord, err := s.accountDao.GetTokenBalances(o.UserAddress)
// 		if err != nil {
// 			logger.Error(err)
// 			return err
// 		}

// 		wethAddress := common.HexToAddress(app.Config.Ethereum["weth_address"])

// 		fee := math.Max(o.MakeFee, o.TakeFee)
// 		sellTokenBalanceRecord := balanceRecord[o.SellToken]
// 		sellTokenBalanceRecord.LockedBalance = math.Sub(sellTokenBalanceRecord.LockedBalance, o.SellAmount)
// 		wethTokenBalanceRecord := balanceRecord[wethAddress]
// 		wethTokenBalanceRecord.LockedBalance = math.Sub(wethTokenBalanceRecord.LockedBalance, fee)

// 		err = s.accountDao.UpdateTokenBalance(o.UserAddress, wethAddress, wethTokenBalanceRecord)
// 		if err != nil {
// 			logger.Error(err)
// 			return err
// 		}

// 		err = s.accountDao.UpdateTokenBalance(o.UserAddress, o.SellToken, sellTokenBalanceRecord)
// 		if err != nil {
// 			logger.Error(err)
// 			return err
// 		}
// 	}

// 	return nil
// }

// // this function is resonsible for unlocking of maker's amount in balance document
// // in case maker cancels the order or some error occurs
// func (s *OrderService) cancelOrderUnlockAmount(o *types.Order) error {
// 	// Unlock Amount
// 	acc, err := s.accountDao.GetByAddress(o.UserAddress)
// 	if err != nil {
// 		logger.Error(err)
// 		return err
// 	}

// 	if o.Side == "BUY" {
// 		tokenBalance := acc.TokenBalances[o.QuoteToken]
// 		tokenBalance.Balance.Add(tokenBalance.Balance, o.SellAmount)
// 		tokenBalance.LockedBalance.Sub(tokenBalance.LockedBalance, o.SellAmount)

// 		err = s.accountDao.UpdateTokenBalance(o.UserAddress, o.QuoteToken, tokenBalance)
// 		if err != nil {
// 			logger.Error(err)
// 			return err
// 		}
// 	}

// 	if o.Side == "SELL" {
// 		tokenBalance := acc.TokenBalances[o.BaseToken]
// 		tokenBalance.Balance.Add(tokenBalance.Balance, o.SellAmount)
// 		tokenBalance.LockedBalance.Sub(tokenBalance.LockedBalance, o.SellAmount)

// 		err = s.accountDao.UpdateTokenBalance(o.UserAddress, o.BaseToken, tokenBalance)
// 		if err != nil {
// 			logger.Error(err)
// 			return err
// 		}
// 	}

// 	return nil
// }
