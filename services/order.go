package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math"
	"time"

	"github.com/Proofsuite/amp-matching-engine/ws"

	"gopkg.in/mgo.v2/bson"

	"github.com/Proofsuite/amp-matching-engine/daos"
	"github.com/Proofsuite/amp-matching-engine/engine"
	"github.com/Proofsuite/amp-matching-engine/types"
)

// OrderService struct with daos required, responsible for communicating with daos.
// OrderService functions are responsible for interacting with daos and implements business logics.
type OrderService struct {
	orderDao   *daos.OrderDao
	balanceDao *daos.BalanceDao
	pairDao    *daos.PairDao
	tradeDao   *daos.TradeDao
	addressDao *daos.AddressDao
	engine     *engine.Resource
}

// NewOrderService returns a new instance of orderservice
func NewOrderService(orderDao *daos.OrderDao, balanceDao *daos.BalanceDao, pairDao *daos.PairDao, tradeDao *daos.TradeDao, addressDao *daos.AddressDao, engine *engine.Resource) *OrderService {
	return &OrderService{orderDao, balanceDao, pairDao, tradeDao, addressDao, engine}
}

// Create validates if the passed order is valid or not based on user's available
// funds and order data.
// If valid: Order is inserted in DB with order status as new and order is publiched
// on rabbitmq queue for matching engine to process the order
func (s *OrderService) Create(order *types.Order) (err error) {

	// Fill token and pair data
	fmt.Printf("\n %v === %v\n", order.BaseToken, order.QuoteToken)

	p, err := s.pairDao.GetByBuySellTokenAddress(order.BaseToken, order.QuoteToken)
	if err != nil {
		return err
	} else if p == nil {
		return errors.New("Pair not found")
	}

	if order.SellToken == p.QuoteTokenAddress {
		order.Side = types.BUY
	} else {
		order.Side = types.SELL
	}

	order.BaseToken = p.BaseTokenAddress
	order.QuoteToken = p.QuoteTokenAddress

	order.PairID = p.ID
	order.PairName = p.Name

	// Validate if order is valid

	addr, err := s.addressDao.GetByAddress(order.UserAddress)
	if err != nil {
		return err
	} else if addr.IsBlocked {
		return fmt.Errorf("Address: %+v isBlocked", addr)
	} else if addr.Nonce != order.Nonce {
		return fmt.Errorf("Order Nonce: %v is not valid expecting nonce: %v", order.Nonce, addr.Nonce)
	}

	// balance validation
	bal, err := s.balanceDao.GetByAddress(order.UserAddress)
	if err != nil {
		return err
	}
	if order.Side == types.BUY {
		amt := bal.Tokens[order.QuoteToken]
		if amt.Amount < order.SellAmount+order.Fee {
			return errors.New("Insufficient Balance")
		}
		fmt.Println("Buy : Verified")

		amt.Amount = amt.Amount - (order.SellAmount)             // + order.Fee
		amt.LockedAmount = amt.LockedAmount + (order.SellAmount) // + order.Fee
		err = s.balanceDao.UpdateAmount(order.UserAddress, order.QuoteToken, &amt)

		if err != nil {
			return err
		}

	} else if order.Side == types.SELL {
		amt := bal.Tokens[order.BaseToken]
		if amt.Amount < order.BuyAmount+order.Fee {
			return errors.New("Insufficient Balance")
		}
		fmt.Println("Sell : Verified")
		amt.Amount = amt.Amount - (order.BuyAmount)             // + order.Fee
		amt.LockedAmount = amt.LockedAmount + (order.BuyAmount) // + order.Fee
		err = s.balanceDao.UpdateAmount(order.UserAddress, order.BaseToken, &amt)

		if err != nil {
			return err
		}
	}

	if err = s.orderDao.Create(order); err != nil {
		return
	}
	if err = s.addressDao.IncrNonce(order.UserAddress); err != nil {
		return
	}

	// Push order to queue
	orderAsBytes, _ := json.Marshal(order)
	s.engine.PublishMessage(&engine.Message{Type: "NEW_ORDER", Data: orderAsBytes})
	return err
}

// GetByID fetches the details of an order using order's mongo ID
func (s *OrderService) GetByID(id bson.ObjectId) (*types.Order, error) {
	return s.orderDao.GetByID(id)
}

// GetByUserAddress fetches all the orders placed by passed user address
func (s *OrderService) GetByUserAddress(address string) ([]*types.Order, error) {
	return s.orderDao.GetByUserAddress(address)
}

// RecoverOrders recovers orders i.e puts back matched orders to orderbook
// in case of failure of trade signing by the maker
func (s *OrderService) RecoverOrders(engineResponse *engine.Response) {
	if err := s.engine.RecoverOrders(engineResponse.MatchingOrders); err != nil {
		panic(err)
	}
	engineResponse.FillStatus = engine.ERROR
	engineResponse.Order.Status = types.ERROR
	engineResponse.Trades = nil
	engineResponse.RemainingOrder = nil
	engineResponse.MatchingOrders = nil
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
	dab, _ := json.Marshal(dbOrder)
	fmt.Printf("%s", dab)
	if dbOrder.Status == types.OPEN || dbOrder.Status == types.NEW {
		engineResponse, err := s.engine.CancelOrder(dbOrder)
		if err != nil {
			return err
		}
		s.orderDao.Update(engineResponse.Order.ID, engineResponse.Order)
		if err := s.cancelOrderUnlockAmount(engineResponse); err != nil {
			return err
		}
		s.SendMessage("CANCEL_ORDER", engineResponse.Order.Hash, engineResponse)
		s.RelayUpdateOverSocket(engineResponse)
		return nil
	}
	return fmt.Errorf("Cannot cancel the order")
}

// UpdateUsingEngineResponse is responsible for updating order status of maker
// and taker orders and transfer/unlock amount based on the response sent by the
// matching engine
func (s *OrderService) UpdateUsingEngineResponse(er *engine.Response) {
	if er.FillStatus == engine.ERROR {
		fmt.Println("Error")
		s.orderDao.Update(er.Order.ID, er.Order)
		s.cancelOrderUnlockAmount(er)
	} else if er.FillStatus == engine.NOMATCH {
		fmt.Println("No Match")
		s.orderDao.Update(er.Order.ID, er.Order)

		// TODO: Update locked amount (change taker fee to maker fee)
		// res, err := s.balanceDao.GetByAddress(er.Order.UserAddress)
		// if err != nil {
		// 	log.Fatalf("\n%s\n", err)
		// }

	} else if er.FillStatus == engine.FULL || er.FillStatus == engine.PARTIAL {
		fmt.Printf("\nPartial Or filled: %d\n", er.FillStatus)

		s.orderDao.Update(er.Order.ID, er.Order)
		// Unlock and transfer Amount
		s.transferAmount(er.Order, er.Order.FilledAmount)

		for _, mo := range er.MatchingOrders {
			s.orderDao.Update(mo.Order.ID, mo.Order)

			// Unlock and transfer Amount
			s.transferAmount(mo.Order, mo.Amount)
		}
		if len(er.Trades) != 0 {
			err := s.tradeDao.Create(er.Trades...)
			if err != nil {
				log.Fatalf("\n Error adding trades to db: %s\n", err)
			}
		}
	}
}

// RelayUpdateOverSocket is responsible for notifying listening clients about new order/trade addition/deletion
func (s *OrderService) RelayUpdateOverSocket(er *engine.Response) {
	if len(er.Trades) > 0 {
		fmt.Println("Trade relay over socket")
		ws.GetPairSockets().WriteMessage(er.Order.BaseToken, er.Order.QuoteToken, "TRADES_ADDED", er.Trades)
	}
	if er.RemainingOrder != nil {
		fmt.Println("Order added Relay over socket")
		ws.GetPairSockets().WriteMessage(er.Order.BaseToken, er.Order.QuoteToken, "ORDER_ADDED", er.RemainingOrder)
	}
	if er.FillStatus == engine.CANCELLED {
		fmt.Println("Order cancelled Relay over socket")
		ws.GetPairSockets().WriteMessage(er.Order.BaseToken, er.Order.QuoteToken, "ORDER_CANCELED", er.Order)
	}
}

// SendMessage is responsible for sending message to socket linked to a particular order
func (s *OrderService) SendMessage(msgType string, hash string, data interface{}) {
	ws.OrderSendMessage(ws.GetOrderConn(hash), msgType, data, hash)
}

// this function is responsible for unlocking of maker's amount in balance document
// in case maker cancels the order or some error occurs
func (s *OrderService) cancelOrderUnlockAmount(er *engine.Response) error {
	// Unlock Amount
	res, err := s.balanceDao.GetByAddress(er.Order.UserAddress)
	if err != nil {
		log.Fatalf("\n%s\n", err)
		return err
	}
	if er.Order.Side == types.BUY {
		bal := res.Tokens[er.Order.QuoteToken]
		fmt.Println("===> buy bal")
		fmt.Println(bal)
		bal.Amount = bal.Amount + (er.Order.SellAmount)
		bal.LockedAmount = bal.LockedAmount - (er.Order.SellAmount)

		fmt.Println("===> updated bal")
		fmt.Println(bal)
		err = s.balanceDao.UpdateAmount(er.Order.UserAddress, er.Order.QuoteToken, &bal)
		if err != nil {
			log.Fatalf("\n%s\n", err)
			return err
		}
	}
	if er.Order.Side == types.SELL {
		bal := res.Tokens[er.Order.BaseToken]
		fmt.Println("===> sell bal")
		fmt.Println(bal)
		bal.Amount = bal.Amount + (er.Order.BuyAmount)
		bal.LockedAmount = bal.LockedAmount - (er.Order.BuyAmount)
		fmt.Println("===> updated bal")
		fmt.Println(bal)
		err = s.balanceDao.UpdateAmount(er.Order.UserAddress, er.Order.BaseToken, &bal)
		if err != nil {
			log.Fatalf("\n%s\n", err)
			return err
		}
	}
	return nil
}

// transferAmount is used to transfer amount from seller to buyer
// it removes the lockedAmount of one token and adds confirmed amount for another token
// based on the type of order i.e. buy/sell
func (s *OrderService) transferAmount(order *types.Order, filledAmount int64) {

	res, _ := s.balanceDao.GetByAddress(order.UserAddress)

	if order.Side == types.BUY {
		sbal := res.Tokens[order.QuoteToken]
		sbal.LockedAmount = sbal.LockedAmount - int64((float64(filledAmount)/math.Pow10(8))*float64(order.Price))
		err := s.balanceDao.UpdateAmount(order.UserAddress, order.QuoteToken, &sbal)
		if err != nil {
			log.Fatalf("\n%s\n", err)
		}
		bbal := res.Tokens[order.BaseToken]
		bbal.Amount = bbal.Amount + filledAmount
		err = s.balanceDao.UpdateAmount(order.UserAddress, order.BaseToken, &bbal)
		if err != nil {
			log.Fatalf("\n%s\n", err)
		}
		fmt.Printf("\n Order Buy\n==>sbal: %v \n==>bbal: %v\n==>Unlock Amount: %v\n", sbal, bbal, int64((float64(filledAmount)/math.Pow10(8))*float64(order.Price)))
	}
	if order.Side == types.SELL {
		bbal := res.Tokens[order.BaseToken]
		bbal.LockedAmount = bbal.LockedAmount - filledAmount
		err := s.balanceDao.UpdateAmount(order.UserAddress, order.BaseToken, &bbal)
		if err != nil {
			log.Fatalf("\n%s\n", err)
		}

		sbal := res.Tokens[order.QuoteToken]
		sbal.Amount = sbal.Amount + int64((float64(filledAmount)/math.Pow10(8))*float64(order.Price))
		err = s.balanceDao.UpdateAmount(order.UserAddress, order.QuoteToken, &sbal)
		if err != nil {
			log.Fatalf("\n%s\n", err)
		}
		fmt.Printf("\n Order Sell\n==>sbal: %v \n==>bbal: %v\n==>Unlock Amount: %v\n", sbal, bbal, filledAmount)

	}
}

func (s *OrderService) EngineResponse(engineResponse *engine.Response) error {
	if engineResponse.FillStatus == engine.NOMATCH {
		s.SendMessage("ORDER_ADDED", engineResponse.Order.Hash, engineResponse)
	} else {
		s.SendMessage("REQUEST_SIGNATURE", engineResponse.Order.Hash, engineResponse)

		t := time.NewTimer(10 * time.Second)
		ch := ws.GetOrderChannel(engineResponse.Order.Hash)
		if ch == nil {
			s.RecoverOrders(engineResponse)
		} else {

			select {
			case rm := <-ch:
				if rm.MsgType == "REQUEST_SIGNATURE" {
					mb, err := json.Marshal(rm.Data)
					if err != nil {
						fmt.Printf("=== Error while marshaling EngineResponse===")
						s.RecoverOrders(engineResponse)
						ws.OrderSendErrorMessage(ws.GetOrderConn(engineResponse.Order.Hash), engineResponse.Order.Hash, err.Error())
					}

					var ersb *engine.Response
					err = json.Unmarshal(mb, &ersb)
					if err != nil {
						fmt.Printf("=== Error while unmarshaling EngineResponse===")
						ws.OrderSendErrorMessage(ws.GetOrderConn(engineResponse.Order.Hash), engineResponse.Order.Hash, err.Error())
						s.RecoverOrders(engineResponse)
					}

					if engineResponse.FillStatus == engine.PARTIAL {
						engineResponse.Order.OrderBook = &types.OrderSubDoc{Amount: ersb.RemainingOrder.Amount, Signature: ersb.RemainingOrder.Signature}
						orderAsBytes, _ := json.Marshal(engineResponse.Order)
						s.engine.PublishMessage(&engine.Message{Type: "remaining_order_add", Data: orderAsBytes})
					}

				}
				t.Stop()
				break

			case <-t.C:
				fmt.Printf("\nTimeout\n")
				s.RecoverOrders(engineResponse)
				t.Stop()
				break
			}
		}
	}
	s.UpdateUsingEngineResponse(engineResponse)
	// TODO: send to operator for blockchain execution

	s.RelayUpdateOverSocket(engineResponse)
	ws.CloseOrderReadChannel(engineResponse.Order.Hash)

	return nil
}
