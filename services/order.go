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

type OrderService struct {
	orderDao   *daos.OrderDao
	balanceDao *daos.BalanceDao
	pairDao    *daos.PairDao
	tradeDao   *daos.TradeDao
}

func NewOrderService(orderDao *daos.OrderDao, balanceDao *daos.BalanceDao, pairDao *daos.PairDao, tradeDao *daos.TradeDao) *OrderService {
	return &OrderService{orderDao, balanceDao, pairDao, tradeDao}
}

func (s *OrderService) Create(order *types.Order) (err error) {
	
	// Fill token and pair data

	p, err := s.pairDao.GetByTokenAddressPair(order.BuyTokenAddress, order.SellTokenAddress)
	if err != nil {
		return err
	} else if p == nil {
		return errors.New("Pair not found")
	}
	order.PairID = p.ID
	order.BuyToken = p.BuyTokenSymbol
	order.BuyTokenAddress = p.BuyTokenSymbol
	order.SellToken = p.SellTokenSymbol
	order.SellTokenAddress = p.SellTokenAddress

	// Validate if order is valid

	// balance validation
	bal, err := s.balanceDao.GetByAddress(order.UserAddress)
	if err != nil {
		return err
	}
	if order.Type == types.BUY {
		amt := bal.Tokens[order.SellToken]
		if amt.Amount < order.AmountSell+order.Fee {
			return errors.New("Insufficient Balance")
		}
		fmt.Println("Buy : Verified")

		amt.Amount = amt.Amount - (order.AmountSell)             // + order.Fee
		amt.LockedAmount = amt.LockedAmount + (order.AmountSell) // + order.Fee
		err = s.balanceDao.UpdateAmount(order.UserAddress, order.SellToken, &amt)

		if err != nil {
			return err
		}

	} else if order.Type == types.SELL {
		amt := bal.Tokens[order.BuyToken]
		if amt.Amount < order.AmountBuy+order.Fee {
			return errors.New("Insufficient Balance")
		}
		fmt.Println("Sell : Verified")
		amt.Amount = amt.Amount - (order.AmountBuy)             // + order.Fee
		amt.LockedAmount = amt.LockedAmount + (order.AmountBuy) // + order.Fee
		err = s.balanceDao.UpdateAmount(order.UserAddress, order.BuyToken, &amt)

		if err != nil {
			return err
		}
	}

	if err = s.orderDao.Create(order); err != nil {
		return
	}

	// Push order to queue
	engine.Engine.PublishOrder(order)
	return err
}

func (s *OrderService) GetByID(id bson.ObjectId) (*types.Order, error) {
	return s.orderDao.GetByID(id)
}
func (s *OrderService) GetByUserAddress(address string) ([]*types.Order, error) {
	return s.orderDao.GetByUserAddress(address)
}
func (s *OrderService) GetAll() ([]types.Order, error) {
	return s.orderDao.GetAll()
}

func (s *OrderService) UpdateUsingEngineResponse(er *engine.EngineResponse) {
	if er.FillStatus == engine.ERROR {
		fmt.Println("Error")
		s.orderDao.Update(er.Order.ID, er.Order)
		res, err := s.balanceDao.GetByAddress(er.Order.UserAddress)
		if err != nil {
			log.Fatalf("\n%s\n", err)
		}

		// Unlock Amount
		if er.Order.Type == types.BUY {
			bal := res.Tokens[er.Order.SellToken]
			fmt.Println("===> buy bal")
			fmt.Println(bal)
			bal.Amount = bal.Amount + (er.Order.AmountSell)
			bal.LockedAmount = bal.LockedAmount - (er.Order.AmountSell)

			fmt.Println("===> updated bal")
			fmt.Println(bal)
			err = s.balanceDao.UpdateAmount(er.Order.UserAddress, er.Order.SellToken, &bal)
			if err != nil {
				log.Fatalf("\n%s\n", err)
			}
		}
		if er.Order.Type == types.SELL {
			bal := res.Tokens[er.Order.BuyToken]
			fmt.Println("===> sell bal")
			fmt.Println(bal)
			bal.Amount = bal.Amount + (er.Order.AmountBuy)
			bal.LockedAmount = bal.LockedAmount - (er.Order.AmountBuy)
			fmt.Println("===> updated bal")
			fmt.Println(bal)
			err = s.balanceDao.UpdateAmount(er.Order.UserAddress, er.Order.BuyToken, &bal)
			if err != nil {
				log.Fatalf("\n%s\n", err)
			}
		}
	} else if er.FillStatus == engine.NO_MATCH {
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
		res, _ := s.balanceDao.GetByAddress(er.Order.UserAddress)

		// TODO: Move this code block to different function

		if er.Order.Type == types.BUY {
			sbal := res.Tokens[er.Order.SellToken]
			sbal.LockedAmount = sbal.LockedAmount - int64((float64(er.Order.FilledAmount)/math.Pow10(8))*float64(er.Order.Price))
			err := s.balanceDao.UpdateAmount(er.Order.UserAddress, er.Order.SellToken, &sbal)
			if err != nil {
				log.Fatalf("\n%s\n", err)
			}
			bbal := res.Tokens[er.Order.BuyToken]
			bbal.Amount = bbal.Amount + er.Order.FilledAmount
			err = s.balanceDao.UpdateAmount(er.Order.UserAddress, er.Order.BuyToken, &bbal)
			if err != nil {
				log.Fatalf("\n%s\n", err)
			}
			fmt.Printf("\n Order Buy\n==>sbal: %s \n==>bbal: %s\n==>Unlock Amount: %d\n", sbal, bbal, int64((float64(er.Order.FilledAmount)/math.Pow10(8))*float64(er.Order.Price)))
		}
		if er.Order.Type == types.SELL {
			bbal := res.Tokens[er.Order.BuyToken]
			bbal.LockedAmount = bbal.LockedAmount - er.Order.FilledAmount
			err := s.balanceDao.UpdateAmount(er.Order.UserAddress, er.Order.BuyToken, &bbal)
			if err != nil {
				log.Fatalf("\n%s\n", err)
			}

			sbal := res.Tokens[er.Order.SellToken]
			sbal.Amount = sbal.Amount + int64((float64(er.Order.FilledAmount)/math.Pow10(8))*float64(er.Order.Price))
			err = s.balanceDao.UpdateAmount(er.Order.UserAddress, er.Order.SellToken, &sbal)
			if err != nil {
				log.Fatalf("\n%s\n", err)
			}
			fmt.Printf("\n Order Sell\n==>sbal: %s \n==>bbal: %s\n==>Unlock Amount: %d\n", sbal, bbal, er.Order.FilledAmount)

		}
		for _, mo := range er.MatchingOrders {
			// fmt.Println(mo.Order)
			s.orderDao.Update(mo.Order.ID, mo.Order)
			res, _ := s.balanceDao.GetByAddress(mo.Order.UserAddress)
			// Unlock Amount
			if mo.Order.Type == types.BUY {
				sbal := res.Tokens[mo.Order.SellToken]
				sbal.LockedAmount = sbal.LockedAmount - int64((float64(mo.Amount)/math.Pow10(8))*float64(er.Order.Price))
				err := s.balanceDao.UpdateAmount(er.Order.UserAddress, er.Order.SellToken, &sbal)
				if err != nil {
					log.Fatalf("\n%s\n", err)
				}
				bbal := res.Tokens[er.Order.BuyToken]
				bbal.Amount = bbal.Amount + mo.Amount
				err = s.balanceDao.UpdateAmount(er.Order.UserAddress, er.Order.BuyToken, &bbal)
				if err != nil {
					log.Fatalf("\n%s\n", err)
				}
				fmt.Printf("\n Match Buy\n==>sbal: %s \n==>bbal: %s\n==>Unlock Amount: %d\n", sbal, bbal, int64((float64(mo.Amount)/math.Pow10(8))*float64(er.Order.Price)))

			}
			if mo.Order.Type == types.SELL {
				bbal := res.Tokens[mo.Order.BuyToken]
				bbal.LockedAmount = bbal.LockedAmount - mo.Amount
				err := s.balanceDao.UpdateAmount(er.Order.UserAddress, er.Order.BuyToken, &bbal)
				if err != nil {
					log.Fatalf("\n%s\n", err)
				}

				sbal := res.Tokens[er.Order.SellToken]
				sbal.Amount = sbal.Amount + int64((float64(mo.Amount)/math.Pow10(8))*float64(er.Order.Price))
				err = s.balanceDao.UpdateAmount(er.Order.UserAddress, er.Order.SellToken, &sbal)
				if err != nil {
					log.Fatalf("\n%s\n", err)
				}
				fmt.Printf("\n Match Sell\n==>sbal: %s \n==>bbal: %s\n==>Unlock Amount: %d\n", sbal, bbal, mo.Amount)

			}
		}
		err := s.tradeDao.Create(er.Trades...)
		if err != nil {
			log.Fatalf("\n Error adding trades to db: %s\n", err)
		}
	}
}

// RelayUpdateOverSocket is responsible for notifying listening clients about new order/trade addition/deletion
func (s *OrderService) RelayUpdateOverSocket(er *engine.EngineResponse) {
	if len(er.Trades) > 0 {
		message := &types.OrderMessage{
			MsgType: "trades_added",
			Data:    er.Trades,
		}
		mab, _ := json.Marshal(message)
		ws.PairSocketWriteMessage(er.Order.PairName, mab)
	}
	if er.RemainingOrder != nil {
		message := &types.OrderMessage{
			MsgType: "order_added",
			Data:    er.RemainingOrder,
		}
		mab, _ := json.Marshal(message)
		ws.PairSocketWriteMessage(er.Order.PairName, mab)
	}
}
