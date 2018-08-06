package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math"

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

	p, err := s.pairDao.GetByTokenAddress(order.BaseTokenAddress, order.QuoteTokenAddress)
	if err != nil {
		return err
	} else if p == nil {
		return errors.New("Pair not found")
	}
	order.PairID = p.ID
	order.PairName = p.Name
	order.BaseToken = p.BaseTokenSymbol
	order.BaseTokenAddress = p.BaseTokenAddress
	order.QuoteToken = p.QuoteTokenSymbol
	order.QuoteTokenAddress = p.QuoteTokenAddress

	// Validate if order is valid

	addr, err := s.addressDao.GetByAddress(order.UserAddress)
	if err != nil {
		return err
	} else if addr.IsBlocked {
		return fmt.Errorf("Address: %s isBlocked", addr)
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
		if amt.Amount < order.AmountSell+order.Fee {
			return errors.New("Insufficient Balance")
		}
		fmt.Println("Buy : Verified")

		amt.Amount = amt.Amount - (order.AmountSell)             // + order.Fee
		amt.LockedAmount = amt.LockedAmount + (order.AmountSell) // + order.Fee
		err = s.balanceDao.UpdateAmount(order.UserAddress, order.QuoteToken, &amt)

		if err != nil {
			return err
		}

	} else if order.Side == types.SELL {
		amt := bal.Tokens[order.BaseToken]
		if amt.Amount < order.AmountBuy+order.Fee {
			return errors.New("Insufficient Balance")
		}
		fmt.Println("Sell : Verified")
		amt.Amount = amt.Amount - (order.AmountBuy)             // + order.Fee
		amt.LockedAmount = amt.LockedAmount + (order.AmountBuy) // + order.Fee
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
	s.engine.PublishMessage(&engine.Message{Type: "new_order", Data: orderAsBytes})
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
		s.SendMessage("cancel_order", engineResponse.Order.Hash, engineResponse)
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
		message := &types.Message{
			MsgType: "trades_added",
			Data:    er.Trades,
		}
		ws.GetPairSockets().PairSocketWriteMessage(er.Order.BaseTokenAddress, er.Order.QuoteTokenAddress, message)
	}
	if er.RemainingOrder != nil {
		fmt.Println("Order added Relay over socket")
		message := &types.Message{
			MsgType: "order_added",
			Data:    er.RemainingOrder,
		}
		ws.GetPairSockets().PairSocketWriteMessage(er.Order.BaseTokenAddress, er.Order.QuoteTokenAddress, message)
	}
	if er.FillStatus == engine.CANCELLED {
		fmt.Println("Order cancelled Relay over socket")
		message := &types.Message{
			MsgType: "order_cancelled",
			Data:    er.Order,
		}
		ws.GetPairSockets().PairSocketWriteMessage(er.Order.BaseTokenAddress, er.Order.QuoteTokenAddress, message)
	}
}

// SendMessage is responsible for sending message to socket linked to a particular order
func (s *OrderService) SendMessage(msgType string, hash string, data interface{}) {
	msg := &types.Message{MsgType: msgType}
	msg.Hash = hash
	msg.Data = data
	ws.GetOrderConn(hash).WriteJSON(msg)
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
		bal.Amount = bal.Amount + (er.Order.AmountSell)
		bal.LockedAmount = bal.LockedAmount - (er.Order.AmountSell)

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
		bal.Amount = bal.Amount + (er.Order.AmountBuy)
		bal.LockedAmount = bal.LockedAmount - (er.Order.AmountBuy)
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
