package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math/big"
	"time"

	"github.com/Proofsuite/amp-matching-engine/ws"
	"github.com/ethereum/go-ethereum/common"

	"gopkg.in/mgo.v2/bson"

	"github.com/Proofsuite/amp-matching-engine/daos"
	"github.com/Proofsuite/amp-matching-engine/engine"
	"github.com/Proofsuite/amp-matching-engine/types"
)

// OrderService struct with daos required, responsible for communicating with daos.
// OrderService functions are responsible for interacting with daos and implements business logics.
type OrderService struct {
	orderDao   *daos.OrderDao
	pairDao    *daos.PairDao
	accountDao *daos.AccountDao
	tradeDao   *daos.TradeDao
	engine     *engine.Resource
}

// NewOrderService returns a new instance of orderservice
func NewOrderService(orderDao *daos.OrderDao, pairDao *daos.PairDao, accountDao *daos.AccountDao, tradeDao *daos.TradeDao, engine *engine.Resource) *OrderService {
	return &OrderService{orderDao, pairDao, accountDao, tradeDao, engine}
}

// GetByID fetches the details of an order using order's mongo ID
func (s *OrderService) GetByID(id bson.ObjectId) (*types.Order, error) {
	return s.orderDao.GetByID(id)
}

// GetByUserAddress fetches all the orders placed by passed user address
func (s *OrderService) GetByUserAddress(address common.Address) ([]*types.Order, error) {
	return s.orderDao.GetByUserAddress(address)
}

// Create validates if the passed order is valid or not based on user's available
// funds and order data.
// If valid: Order is inserted in DB with order status as new and order is publiched
// on rabbitmq queue for matching engine to process the order
func (s *OrderService) NewOrder(o *types.Order) error {
	// Validate if the address is not blacklisted
	acc, err := s.accountDao.GetByAddress(o.UserAddress)
	if err != nil {
		return err
	}
	if acc.IsBlocked {
		return fmt.Errorf("Address: %+v isBlocked", acc)
	}

	if err := o.Validate(); err != nil {
		return err
	}

	ok, err := o.VerifySignature()
	if err != nil {
		return err
	}
	if !ok {
		return errors.New("Invalid signature")
	}

	p, err := s.pairDao.GetByBuySellTokenAddress(o.BuyToken, o.SellToken)
	if err != nil {
		return err
	}
	if p == nil {
		return errors.New("Pair not found")
	}

	// Fill token and pair data
	err = o.Process(p)
	if err != nil {
		return err
	}

	// fee balance validation
	wethTokenBalance, err := s.accountDao.GetWethTokenBalance(o.UserAddress)
	if err != nil {
		return err
	}

	if wethTokenBalance.Balance.Cmp(o.MakeFee) == -1 {
		return errors.New("Insufficient WETH Balance")
	}

	if wethTokenBalance.Balance.Cmp(o.TakeFee) == -1 {
		return errors.New("Insufficient WETH Balance")
	}

	if wethTokenBalance.Allowance.Cmp(o.MakeFee) == -1 {
		return errors.New("Insufficient WETH Allowance")
	}

	if wethTokenBalance.Allowance.Cmp(o.TakeFee) == -1 {
		return errors.New("Insufficient WETH Allowance")
	}

	wethTokenBalance.Balance.Sub(wethTokenBalance.Balance, o.MakeFee)
	wethTokenBalance.LockedBalance.Add(wethTokenBalance.LockedBalance, o.TakeFee)

	err = s.accountDao.UpdateTokenBalance(o.UserAddress, o.QuoteToken, wethTokenBalance)
	if err != nil {
		return err
	}

	// balance validation
	sellTokenBalance, err := s.accountDao.GetTokenBalance(o.UserAddress, o.SellToken)
	if err != nil {
		return err
	}

	if sellTokenBalance.Balance.Cmp(o.SellAmount) != 1 {
		return errors.New("Insufficient Balance")
	}

	if sellTokenBalance.Allowance.Cmp(o.SellAmount) != 1 {
		return errors.New("Insufficient Allowance")
	}

	sellTokenBalance.Balance.Sub(sellTokenBalance.Balance, o.SellAmount)
	sellTokenBalance.LockedBalance.Add(sellTokenBalance.Balance, o.SellAmount)
	err = s.accountDao.UpdateTokenBalance(o.UserAddress, o.SellToken, sellTokenBalance)
	if err != nil {
		return err
	}

	if err = s.orderDao.Create(o); err != nil {
		return err
	}

	// Push o to queue
	bytes, _ := json.Marshal(o)

	s.SendMessage("ORDER_ADDED", o.Hash, o)
	s.engine.PublishMessage(&engine.Message{Type: "NEW_ORDER", Data: bytes})
	return nil
}

// CancelOrder handles the cancellation order requests.
// Only Orders which are OPEN or NEW i.e. Not yet filled/partially filled
// can be cancelled
func (s *OrderService) CancelOrder(order *types.Order) error {
	dbOrder, err := s.orderDao.GetByHash(order.Hash)
	if err != nil {
		return err
	}

	if dbOrder == nil {
		return fmt.Errorf("No order with this hash present")
	}

	_, err = json.Marshal(dbOrder)
	if err != nil {
		return err
	}

	if dbOrder.Status == "OPEN" || dbOrder.Status == "NEW" {
		res, err := s.engine.CancelOrder(dbOrder)
		if err != nil {
			return err
		}

		s.orderDao.Update(res.Order.ID, res.Order)
		if err := s.cancelOrderUnlockAmount(res.Order); err != nil {
			return err
		}

		s.SendMessage("ORDER_CANCELLED", res.Order.Hash, res)
		s.RelayUpdateOverSocket(res)
		return nil
	}

	return fmt.Errorf("Cannot cancel the order")
}

// HandleEngineResponse listens to messages incoming from the engine and handles websocket
// responses and database updates accordingly
func (s *OrderService) HandleEngineResponse(res *engine.Response) error {
	switch res.FillStatus {
	case engine.ERROR:
		s.handleEngineError(res)
	case engine.NOMATCH:
		s.handleEngineOrderAdded(res)
	case engine.FULL:
	case engine.PARTIAL:
		s.handleEngineOrderMatched(res)
	default:
		s.handleEngineUnknownMessage(res)
	}

	s.RelayUpdateOverSocket(res)
	ws.CloseOrderReadChannel(res.Order.Hash)
	return nil
}

// handleEngineError returns an websocket error message to the client and recovers orders on the
// redis key/value store
func (s *OrderService) handleEngineError(res *engine.Response) {
	s.orderDao.Update(res.Order.ID, res.Order)
	s.cancelOrderUnlockAmount(res.Order)
	ws.OrderSendErrorMessage(ws.GetOrderConn(res.Order.Hash), "Some error", res.Order.Hash)
}

// handleEngineOrderAdded returns a websocket message informing the client that his order has been added
// to the orderbook (but currently not matched)
func (s *OrderService) handleEngineOrderAdded(res *engine.Response) {
	s.SendMessage("ORDER_ADDED", res.Order.Hash, res)
}

// handleEngineOrderMatched returns a websocket message informing the client that his order has been added.
// The request signature message also signals the client to sign trades.
func (s *OrderService) handleEngineOrderMatched(resp *engine.Response) {
	s.SendMessage("REQUEST_SIGNATURE", resp.Order.Hash, resp)
	s.orderDao.Update(resp.Order.ID, resp.Order)
	s.transferAmount(resp.Order, big.NewInt(resp.Order.FilledAmount))

	for _, o := range resp.MatchingOrders {
		s.orderDao.Update(o.Order.ID, resp.Order)
		s.transferAmount(o.Order, big.NewInt(o.Amount))
	}

	if len(resp.Trades) != 0 {
		err := s.tradeDao.Create(resp.Trades...)
		if err != nil {
			log.Fatalf("\n Error saving trades to db: %s\n", err)
		}
	}

	t := time.NewTimer(10 * time.Second)
	ch := ws.GetOrderChannel(resp.Order.Hash)

	if ch == nil {
		s.RecoverOrders(resp)
	} else {
		select {
		case msg := <-ch:
			if msg.Type == "SUBMIT_SIGNATURE" {
				bytes, err := json.Marshal(msg.Data)
				if err != nil {
					s.RecoverOrders(resp)
					ws.OrderSendErrorMessage(ws.GetOrderConn(resp.Order.Hash), err.Error(), resp.Order.Hash)
				}

				clientResponse := &engine.Response{}
				err = json.Unmarshal(bytes, clientResponse)
				if err != nil {
					s.RecoverOrders(resp)
					ws.OrderSendErrorMessage(ws.GetOrderConn(resp.Order.Hash), err.Error(), resp.Order.Hash)
				}

				if clientResponse.FillStatus == engine.PARTIAL {
					resp.Order.OrderBook = &types.OrderSubDoc{Amount: clientResponse.RemainingOrder.Amount, Signature: clientResponse.RemainingOrder.Signature}
					bytes, _ := json.Marshal(resp.Order)
					s.engine.PublishMessage(&engine.Message{Type: "ADD_ORDER", Data: bytes})
				}
			}

			t.Stop()
			break
		case <-t.C:
			s.RecoverOrders(resp)
			t.Stop()
			break
		}
	}
}

// handleEngineUnknownMessage returns a websocket messsage in case the engine response is not recognized
func (s *OrderService) handleEngineUnknownMessage(resp *engine.Response) {
	s.RecoverOrders(resp)
	ws.OrderSendErrorMessage(ws.GetOrderConn(resp.Order.Hash), "UNKNOWN_MESSAGE", resp.Order.Hash)
}

// RecoverOrders recovers orders i.e puts back matched orders to orderbook
// in case of failure of trade signing by the maker
func (s *OrderService) RecoverOrders(resp *engine.Response) {
	//WHY NOT SUBMIT MESSAGE VIA A QUEUE ?
	if err := s.engine.RecoverOrders(resp.MatchingOrders); err != nil {
		panic(err)
	}

	resp.FillStatus = engine.ERROR
	resp.Order.Status = "ERROR"
	resp.Trades = nil
	resp.RemainingOrder = nil
	resp.MatchingOrders = nil
}

// RelayUpdateOverSocket is responsible for notifying listening clients about new order/trade addition/deletion
func (s *OrderService) RelayUpdateOverSocket(resp *engine.Response) {
	if len(resp.Trades) > 0 {
		fmt.Println("Trade relay over socket")
		ws.GetPairSockets().WriteMessage(resp.Order.BaseToken, resp.Order.QuoteToken, "TRADES_ADDED", resp.Trades)
	}

	if resp.RemainingOrder != nil {
		fmt.Println("Order added Relay over socket")
		ws.GetPairSockets().WriteMessage(resp.Order.BaseToken, resp.Order.QuoteToken, "ORDER_ADDED", resp.RemainingOrder)
	}

	if resp.FillStatus == engine.CANCELLED {
		fmt.Println("Order cancelled Relay over socket")
		ws.GetPairSockets().WriteMessage(resp.Order.BaseToken, resp.Order.QuoteToken, "ORDER_CANCELED", resp.Order)
	}
}

// SendMessage is responsible for sending message to socket linked to a particular order
func (s *OrderService) SendMessage(msgType string, hash common.Hash, data interface{}) {
	ws.OrderSendMessage(ws.GetOrderConn(hash), msgType, data, hash)
}

// this function is responsible for unlocking of maker's amount in balance document
// in case maker cancels the order or some error occurs
func (s *OrderService) cancelOrderUnlockAmount(o *types.Order) error {
	// Unlock Amount
	acc, err := s.accountDao.GetByAddress(o.UserAddress)
	if err != nil {
		log.Fatalf("\n%v\n", err)
		return err
	}

	if o.Side == "BUY" {
		tokenBalance := acc.TokenBalances[o.QuoteToken]
		tokenBalance.Balance.Add(tokenBalance.Balance, o.SellAmount)
		tokenBalance.LockedBalance.Sub(tokenBalance.LockedBalance, o.SellAmount)

		err = s.accountDao.UpdateTokenBalance(o.UserAddress, o.QuoteToken, tokenBalance)
		if err != nil {
			log.Fatalf("\n%s\n", err)
		}
	}

	if o.Side == "SELL" {
		tokenBalance := acc.TokenBalances[o.BaseToken]
		tokenBalance.Balance.Add(tokenBalance.Balance, o.SellAmount)
		tokenBalance.LockedBalance.Sub(tokenBalance.LockedBalance, o.SellAmount)

		err = s.accountDao.UpdateTokenBalance(o.UserAddress, o.BaseToken, tokenBalance)
		if err != nil {
			log.Fatalf("\n%v\n", err)
		}
	}

	return nil
}

// transferAmount is used to transfer amount from seller to buyer
// it removes the lockedAmount of one token and adds confirmed amount for another token
// based on the type of order i.e. buy/sell
func (s *OrderService) transferAmount(o *types.Order, filledAmount *big.Int) {
	tokenBalances, err := s.accountDao.GetTokenBalances(o.UserAddress)
	if err != nil {
		log.Fatalf("\n%v\n", err)
	}

	if o.Side == "BUY" {
		sellBalance := tokenBalances[o.QuoteToken]
		sellBalance.LockedBalance = sellBalance.LockedBalance.Sub(sellBalance.LockedBalance, filledAmount)

		err = s.accountDao.UpdateTokenBalance(o.UserAddress, o.QuoteToken, sellBalance)
		if err != nil {
			log.Fatalf("\n%v\n", err)
		}

		buyBalance := tokenBalances[o.BaseToken]
		buyBalance.Balance = buyBalance.Balance.Add(buyBalance.Balance, filledAmount)
		err = s.accountDao.UpdateTokenBalance(o.UserAddress, o.BaseToken, buyBalance)
		if err != nil {
			log.Fatalf("\n%v\n", err)
		}
	}

	if o.Side == "SELL" {
		buyBalance := tokenBalances[o.BaseToken]
		buyBalance.LockedBalance = buyBalance.LockedBalance.Sub(buyBalance.LockedBalance, filledAmount)
		err = s.accountDao.UpdateTokenBalance(o.UserAddress, o.BaseToken, buyBalance)
		if err != nil {
			log.Fatalf("\n%v\n", err)
		}

		sellBalance := tokenBalances[o.QuoteToken]
		sellBalance.Balance = sellBalance.Balance.Add(sellBalance.Balance, filledAmount)
		err = s.accountDao.UpdateTokenBalance(o.UserAddress, o.BaseToken, sellBalance)
		if err != nil {
			log.Fatalf("\n%v\n", err)
		}
	}

	// func (s *OrderService) handleNewTrade(msg *types.Message, res *engine.Response) {
	// 	bytes, err := json.Marshal(msg.Data)
	// 	if err != nil {
	// 		s.RecoverOrders(res)
	// 		ws.OrderSendErrorMessage(ws.GetOrderConn(res.Order.Hash), err.Error(), res.Order.Hash)
	// 	}

	// 	resp := &engine.Response{}
	// 	err = json.Unmarshal(bytes, &resp)
	// 	if err != nil {
	// 		s.RecoverOrders(res)
	// 		ws.OrderSendErrorMessage(ws.GetOrderConn(res.Order.Hash), err.Error(), res.Order.Hash)
	// 	}

	// 	if res.FillStatus == engine.PARTIAL {
	// 		res.Order.OrderBook = &types.asdl
	// 		kfjasd
	// 		ljfk
	// 	}

	// }

	// if o.Side == "SELL" {
	// 	sellBalance, err := s.accountDao.GetTokenBalance(o.UserAddress, o.QuoteToken)
	// 	if err != nil {
	// 		log.Fatalf("\n%v\n")
	// 	}

	// 	sellBalance.LockedBalance = sellBalance.LockedBalance.
	// }

	// if o.Side == "BUY" {
	// 	sbal := res.Tokens[o.QuoteToken]
	// 	sbal.LockedAmount = sbal.LockedAmount - int64((float64(filledAmount)/math.Pow10(8))*float64(o.Price))
	// 	err := s.balanceDao.UpdateAmount(o.UserAddress, o.QuoteToken, &sbal)
	// 	if err != nil {
	// 		log.Fatalf("\n%s\n", err)
	// 	}
	// 	bbal := res.Tokens[o.BaseToken]
	// 	bbal.Amount = bbal.Amount + filledAmount
	// 	err = s.balanceDao.UpdateAmount(o.UserAddress, o.BaseToken, &bbal)
	// 	if err != nil {
	// 		log.Fatalf("\n%s\n", err)
	// 	}
	// 	fmt.Printf("\n Order Buy\n==>sbal: %v \n==>bbal: %v\n==>Unlock Amount: %v\n", sbal, bbal, int64((float64(filledAmount)/math.Pow10(8))*float64(o.Price)))
	// }
	// if o.Side == "SELL" {
	// 	bbal := res.Tokens[o.BaseToken]
	// 	bbal.LockedAmount = bbal.LockedAmount - filledAmount
	// 	err := s.balanceDao.UpdateAmount(o.UserAddress, o.BaseToken, &bbal)
	// 	if err != nil {
	// 		log.Fatalf("\n%s\n", err)
	// 	}

	// 	sbal := res.Tokens[o.QuoteToken]
	// 	sbal.Amount = sbal.Amount + int64((float64(filledAmount)/math.Pow10(8))*float64(o.Price))
	// 	err = s.balanceDao.UpdateAmount(o.UserAddress, o.QuoteToken, &sbal)
	// 	if err != nil {
	// 		log.Fatalf("\n%s\n", err)
	// 	}
	// 	fmt.Printf("\n Order Sell\n==>sbal: %v \n==>bbal: %v\n==>Unlock Amount: %v\n", sbal, bbal, filledAmount)

	// }
}

// func (s *OrderService) HandleClientResponse() error {
// 	t := time.NewTimer(10 * time.Second)
// 	ch := ws.GetOrderChannel(res.Order.Hash)
// 	if ch == nil {
// 		s.RecoverOrders(res)
// 	} else {

// 		select {
// 		case msg := <-ch:
// 			if msg.Type == "REQUEST_SIGNATURE" {
// 				bytes, err := json.Marshal(msg.Data)
// 				if err != nil {
// 					fmt.Printf("=== Error while marshaling EngineResponse ===")
// 					s.RecoverOrders(res)
// 					ws.OrderSendErrorMessage(ws.GetOrderConn(res.Order.Hash), res.Order.Hash, err.Error())
// 				}

// 				var ersb *engine.Response
// 				err = json.Unmarshal(bytes, &ersb)
// 				if err != nil {
// 					fmt.Printf("=== Error while unmarshaling EngineResponse ===")
// 					ws.OrderSendErrorMessage(ws.GetOrderConn(res.Order.Hash), res.Order.Hash, err.Error())
// 					s.RecoverOrders(res)
// 				}

// 				if res.FillStatus == engine.PARTIAL {
// 					res.Order.OrderBook = &types.OrderSubDoc{Amount: ersb.RemainingOrder.Amount, Signature: ersb.RemainingOrder.Signature}
// 					orderAsBytes, _ := json.Marshal(res.Order)
// 					s.engine.PublishMessage(&engine.Message{Type: "remaining_order_add", Data: orderAsBytes})
// 				}

// 			}
// 			t.Stop()
// 			break

// 		case <-t.C:
// 			fmt.Printf("\nTimeout\n")
// 			s.RecoverOrders(res)
// 			t.Stop()
// 			break
// 		}
// 	}
// }

// DEPRECATED
// bal, err := s.balanceDao.GetByAddress(order.UserAddress)
// if err != nil {
// 	return err
// }
// if order.Side == "BUY" {
// 	amt := bal.Tokens[order.QuoteToken]
// 	if amt.Amount < order.SellAmount+order.Fee {
// 		return errors.New("Insufficient Balance")
// 	}
// 	fmt.Println("Buy : Verified")

// 	amt.Amount = amt.Amount - (order.SellAmount)             // + order.Fee
// 	amt.LockedAmount = amt.LockedAmount + (order.SellAmount) // + order.Fee
// 	err = s.balanceDao.UpdateAmount(order.UserAddress, order.QuoteToken, &amt)

// 	if err != nil {
// 		return err
// 	}

// } else if order.Side == "SELL" {
// 	amt := bal.Tokens[order.BaseToken]
// 	if amt.Amount < order.BuyAmount+order.Fee {
// 		return errors.New("Insufficient Balance")
// 	}
// 	fmt.Println("Sell : Verified")
// 	amt.Amount = amt.Amount - (order.BuyAmount)             // + order.Fee
// 	amt.LockedAmount = amt.LockedAmount + (order.BuyAmount) // + order.Fee
// 	err = s.balanceDao.UpdateAmount(order.UserAddress, order.BaseToken, &amt)

// 	if res.FillStatus == engine.ERROR {
// 		fmt.Println("Error")
// 		s.orderDao.Update(res.Ordres.ID, res.Order)
// 		s.cancelOrderUnlockAmount(res.Order)
// 	} else if res.FillStatus == engine.NOMATCH {
// 		fmt.Println("No Match")
// 		s.orderDao.Update(res.Ordres.ID, res.Order)

// 		// TODO: Update locked amount (change taker fee to maker fee)
// 		// res, err := s.balanceDao.GetByAddress(res.Ordres.UserAddress)
// 		// if err != nil {
// 		// 	log.Fatalf("\n%s\n", err)
// 		// }

// 	} else if res.FillStatus == engine.FULL || res.FillStatus == engine.PARTIAL {
// 		fmt.Printf("\nPartial Or filled: %d\n", res.FillStatus)

// 		s.orderDao.Update(res.Order.ID, res.Order)
// 		// Unlock and transfer Amount
// 		s.transferAmount(res.Order, res.Order.FilledAmount)

// 		for _, mo := range res.MatchingOrders {
// 			s.orderDao.Update(mo.Order.ID, mo.Order)

// 			// Unlock and transfer Amount
// 			s.transferAmount(mo.Order, mo.Amount)
// 		}

// 		if len(res.Trades) != 0 {
// 			err := s.tradeDao.Create(res.Trades...)
// 			if err != nil {
// 				log.Fatalf("\n Error adding trades to db: %s\n", err)
// 			}
// 		}

// // 	}
// // }

// // UpdateUsingEngineResponse is responsible for updating order status of maker
// // and taker orders and transfer/unlock amount based on the response sent by the
// // matching engine
// func (s *OrderService) UpdateUsingEngineResponse(res *engine.Response) {
// 	switch res.FillStatus {
// 	case engine.ERROR:
// 		s.orderDao.Update(res.Order.ID, res.Order)
// 		s.cancelOrderUnlockAmount(res.Order)

// 	case engine.NOMATCH:
// 		s.orderDao.Update(res.Order.ID, res.Order)

// 	case engine.FULL:
// 	case engine.PARTIAL:
// 		s.orderDao.Update(res.Order.ID, res.Order)
// 		s.transferAmount(res.Order, big.NewInt(res.Order.FilledAmount))

// 		for _, mo := range res.MatchingOrders {
// 			s.orderDao.Update(mo.Order.ID, mo.Order)
// 			s.transferAmount(mo.Order, big.NewInt(mo.Amount))
// 		}

// 		if len(res.Trades) != 0 {
// 			err := s.tradeDao.Create(res.Trades...)
// 			if err != nil {
// 				log.Fatalf("\n Error saving trades to db: %s\n", err)
// 			}
// 		}
// 	}
