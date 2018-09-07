package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math/big"
	"time"

	"github.com/streadway/amqp"

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
	amqp             *rabbitmq.Connection
}

// NewOrderService returns a new instance of orderservice
func NewOrderService(
	orderDao interfaces.OrderDao,
	pairDao interfaces.PairDao,
	accountDao interfaces.AccountDao,
	tradeDao interfaces.TradeDao,
	engine interfaces.Engine,
	ethereumProvider interfaces.EthereumProvider,
	amqp *rabbitmq.Connection,
) *OrderService {
	return &OrderService{
		orderDao,
		pairDao,
		accountDao,
		tradeDao,
		engine,
		ethereumProvider,
		amqp,
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
		log.Print(err)
		return err
	}

	if p == nil {
		return errors.New("Pair not found")
	}

	// Fill token and pair data
	err = o.Process(p)
	if err != nil {
		log.Print(err)
		return err
	}

	// fee balance validation
	wethAddress := common.HexToAddress("0x2EB24432177e82907dE24b7c5a6E0a5c03226135")
	wethTokenBalance, err := s.accountDao.GetTokenBalance(o.UserAddress, wethAddress)
	if err != nil {
		log.Printf("Error retrieving WETH balance: %v", err.Error())
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

	err = s.accountDao.UpdateTokenBalance(o.UserAddress, wethAddress, wethTokenBalance)
	if err != nil {
		log.Print(err)
		return err
	}

	// balance validation
	sellTokenBalance, err := s.accountDao.GetTokenBalance(o.UserAddress, o.SellToken)
	if err != nil {
		log.Print(err)
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
		log.Print(err)
		return err
	}

	if err = s.orderDao.Create(o); err != nil {
		log.Print(err)
		return err
	}

	// Push o to queue
	bytes, _ := json.Marshal(o)
	s.PublishOrder(&rabbitmq.Message{Type: "NEW_ORDER", Data: bytes})
	return nil
}

// CancelOrder handles the cancellation order requests.
// Only Orders which are OPEN or NEW i.e. Not yet filled/partially filled
// can be cancelled
func (s *OrderService) CancelOrder(oc *types.OrderCancel) error {
	dbOrder, err := s.orderDao.GetByHash(oc.OrderHash)
	if err != nil {
		log.Print(err)
		return err
	}

	if dbOrder == nil {
		return fmt.Errorf("No order with this hash present")
	}

	_, err = json.Marshal(dbOrder)
	if err != nil {
		log.Print(err)
		return err
	}

	if dbOrder.Status == "OPEN" || dbOrder.Status == "NEW" {
		res, err := s.engine.CancelOrder(dbOrder)
		if err != nil {
			log.Print(err)
			return err
		}

		s.orderDao.UpdateByHash(res.Order.Hash, res.Order)
		if err := s.cancelOrderUnlockAmount(res.Order); err != nil {
			log.Print(err)
			return err
		}

		s.SendMessage("ORDER_CANCELLED", res.Order.Hash, res.Order)
		s.RelayUpdateOverSocket(res)
		return nil
	}

	return fmt.Errorf("Cannot cancel the order")
}

// HandleEngineResponse listens to messages incoming from the engine and handles websocket
// responses and database updates accordingly
func (s *OrderService) HandleEngineResponse(res *types.EngineResponse) error {
	switch res.FillStatus {
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

// handleEngineError returns an websocket error message to the client and recovers orders on the
// redis key/value store
func (s *OrderService) handleEngineError(res *types.EngineResponse) {
	s.orderDao.UpdateByHash(res.Order.Hash, res.Order)
	s.cancelOrderUnlockAmount(res.Order)
	ws.SendOrderErrorMessage(ws.GetOrderConnection(res.Order.Hash), "Some error", res.Order.Hash)
}

// handleEngineOrderAdded returns a websocket message informing the client that his order has been added
// to the orderbook (but currently not matched)
func (s *OrderService) handleEngineOrderAdded(res *types.EngineResponse) {
	s.SendMessage("ORDER_ADDED", res.Order.Hash, res.Order)
}

// handleEngineOrderMatched returns a websocket message informing the client that his order has been added.
// The request signature message also signals the client to sign trades.
func (s *OrderService) handleEngineOrderMatched(res *types.EngineResponse) {
	err := s.orderDao.UpdateByHash(res.Order.Hash, res.Order)
	if err != nil {
		log.Print(err)
	}

	s.transferAmount(res.Order, res.Order.FilledAmount)

	for _, mo := range res.MatchingOrders {
		err := s.orderDao.UpdateByHash(mo.Order.Hash, mo.Order)
		if err != nil {
			log.Print(err)
		}

		s.transferAmount(mo.Order, mo.Amount)
	}

	utils.PrintJSON(res)

	if len(res.Trades) != 0 {
		err := s.tradeDao.Create(res.Trades...)
		if err != nil {
			log.Fatal(err)
		}
	}

	go s.handleSubmitSignatures(res)
	s.SendMessage("REQUEST_SIGNATURE", res.Order.Hash, types.SignaturePayload{res.RemainingOrder, res.Trades})
}

// handleSubmitSignatures wait for a submit signature message that provides the matching engine with orders
// that can be broadcast to the exchange smart contrct
func (s *OrderService) handleSubmitSignatures(res *types.EngineResponse) {
	ch := ws.GetOrderChannel(res.Order.Hash)
	t := time.NewTimer(30 * time.Second)

	select {
	case msg := <-ch:
		if msg != nil && msg.Type == "SUBMIT_SIGNATURE" {

			log.Print("SUBMITTING SIGNATURES")

			bytes, err := json.Marshal(msg.Data)
			if err != nil {
				s.RecoverOrders(res)
				ws.SendOrderErrorMessage(ws.GetOrderConnection(res.Order.Hash), err.Error(), res.Order.Hash)
			}

			data := &types.SignaturePayload{}
			err = json.Unmarshal(bytes, data)
			if err != nil {
				log.Print(err)
				s.RecoverOrders(res)
				ws.SendOrderErrorMessage(ws.GetOrderConnection(res.Order.Hash), err.Error(), res.Order.Hash)
			}

			if data.Order != nil {
				bytes, err := json.Marshal(res.Order)
				if err != nil {
					log.Print(err)
					ws.SendOrderErrorMessage(ws.GetOrderConnection(res.Order.Hash), err.Error(), res.Order.Hash)
				}

				log.Print("ADDING NEW ORDER")
				s.PublishOrder(&rabbitmq.Message{Type: "ADD_ORDER", Data: bytes})
			}

			if data.Trades != nil {
				_, err := json.Marshal(res.Order)
				if err != nil {
					log.Print(err)
					ws.SendOrderErrorMessage(ws.GetOrderConnection(res.Order.Hash), err.Error(), res.Order.Hash)
				}
			}
		}

	case <-t.C:
		s.RecoverOrders(res)
		t.Stop()
		break
	}
}

// handleEngineUnknownMessage returns a websocket messsage in case the engine resonse is not recognized
func (s *OrderService) handleEngineUnknownMessage(res *types.EngineResponse) {
	s.RecoverOrders(res)
	ws.SendOrderErrorMessage(ws.GetOrderConnection(res.Order.Hash), "UNKNOWN_MESSAGE", res.Order.Hash)
}

// RecoverOrders recovers orders i.e puts back matched orders to orderbook
// in case of failure of trade signing by the maker
func (s *OrderService) RecoverOrders(res *types.EngineResponse) {
	if err := s.engine.RecoverOrders(res.MatchingOrders); err != nil {
		panic(err)
	}

	res.FillStatus = "ERROR"
	res.Order.Status = "ERROR"
	res.Trades = nil
	res.RemainingOrder = nil
	res.MatchingOrders = nil
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
		log.Print(err)
		return err
	}

	err = s.engine.CancelTrades(orders, amounts)
	if err != nil {
		log.Print(err)
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
	go broadcastLiteOBUpdate(res.Order.BaseToken, res.Order.QuoteToken, getLightOBPayload(res))
	go broadcastFullOBUpdate(res.Order)
}

// RelayTradeUpdate is resonsible for notifying listening clients about new trades
func (s *OrderService) RelayTradeUpdate(res *types.EngineResponse) {

	if len(res.Trades) == 0 {
		return
	}
	// broadcast trades
	go broadcastTradeUpdate(res.Trades)
}

func broadcastLiteOBUpdate(baseToken, quoteToken common.Address, data interface{}) {
	cid := utils.GetOrderBookChannelID(baseToken, quoteToken)
	ws.GetLiteOrderBookSocket().BroadcastMessage(cid, data)
}

func broadcastFullOBUpdate(order *types.Order) {
	cid := utils.GetOrderBookChannelID(order.BaseToken, order.QuoteToken)
	ws.GetFullOrderBookSocket().BroadcastMessage(cid, order)
}

func broadcastTradeUpdate(trades []*types.Trade) {
	cid := utils.GetTradeChannelID(trades[0].BaseToken, trades[0].QuoteToken)
	ws.GetTradeSocket().BroadcastMessage(cid, trades)
}

func getLightOBPayload(res *types.EngineResponse) interface{} {

	orderSide := make(map[string]string)
	matchSide := make([]map[string]string, 0)

	matchSideMap := make(map[string]*big.Int)

	if math.Sub(res.Order.Amount, res.Order.FilledAmount).Cmp(big.NewInt(0)) != 0 {
		orderSide["price"] = res.Order.PricePoint.String()
		orderSide["amount"] = math.Sub(res.Order.Amount, res.Order.FilledAmount).String()
	}

	if len(res.MatchingOrders) > 0 {
		for _, mo := range res.MatchingOrders {
			pp := mo.Order.PricePoint.String()
			if matchSideMap[pp] == nil {
				matchSideMap[pp] = big.NewInt(0)
			}
			matchSideMap[pp] = math.Add(matchSideMap[pp], mo.Amount)
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

// SendMessage is resonsible for sending message to socket linked to a particular order
func (s *OrderService) SendMessage(msgType string, hash common.Hash, data interface{}) {
	ws.SendOrderMessage(ws.GetOrderConnection(hash), msgType, data, hash)
}

// this function is resonsible for unlocking of maker's amount in balance document
// in case maker cancels the order or some error occurs
func (s *OrderService) cancelOrderUnlockAmount(o *types.Order) error {
	// Unlock Amount
	acc, err := s.accountDao.GetByAddress(o.UserAddress)
	if err != nil {
		log.Fatal(err)
		return err
	}

	if o.Side == "BUY" {
		tokenBalance := acc.TokenBalances[o.QuoteToken]
		tokenBalance.Balance.Add(tokenBalance.Balance, o.SellAmount)
		tokenBalance.LockedBalance.Sub(tokenBalance.LockedBalance, o.SellAmount)

		err = s.accountDao.UpdateTokenBalance(o.UserAddress, o.QuoteToken, tokenBalance)
		if err != nil {
			log.Fatal(err)
		}
	}

	if o.Side == "SELL" {
		tokenBalance := acc.TokenBalances[o.BaseToken]
		tokenBalance.Balance.Add(tokenBalance.Balance, o.SellAmount)
		tokenBalance.LockedBalance.Sub(tokenBalance.LockedBalance, o.SellAmount)

		err = s.accountDao.UpdateTokenBalance(o.UserAddress, o.BaseToken, tokenBalance)
		if err != nil {
			log.Fatal(err)
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
		log.Fatal(err)
	}

	if o.Side == "BUY" {
		sellBalance := tokenBalances[o.QuoteToken]
		sellBalance.LockedBalance = sellBalance.LockedBalance.Sub(sellBalance.LockedBalance, filledAmount)

		err = s.accountDao.UpdateTokenBalance(o.UserAddress, o.QuoteToken, sellBalance)
		if err != nil {
			log.Fatal(err)
		}

		buyBalance := tokenBalances[o.BaseToken]
		buyBalance.Balance = buyBalance.Balance.Add(buyBalance.Balance, filledAmount)
		err = s.accountDao.UpdateTokenBalance(o.UserAddress, o.BaseToken, buyBalance)
		if err != nil {
			log.Fatal(err)
		}
	}

	if o.Side == "SELL" {
		buyBalance := tokenBalances[o.BaseToken]
		buyBalance.LockedBalance = buyBalance.LockedBalance.Sub(buyBalance.LockedBalance, filledAmount)
		err = s.accountDao.UpdateTokenBalance(o.UserAddress, o.BaseToken, buyBalance)
		if err != nil {
			log.Fatal(err)
		}

		sellBalance := tokenBalances[o.QuoteToken]
		sellBalance.Balance = sellBalance.Balance.Add(sellBalance.Balance, filledAmount)
		err = s.accountDao.UpdateTokenBalance(o.UserAddress, o.BaseToken, sellBalance)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func (s *OrderService) SubscribeQueue(fn func(*rabbitmq.Message) error) error {
	ch := s.amqp.GetChannel("orderSubscribe")
	q := s.amqp.GetQueue(ch, "order")

	go func() {
		msgs, err := ch.Consume(
			q.Name, // queue
			"",     // consumer
			true,   // auto-ack
			false,  // exclusive
			false,  // no-local
			false,  // no-wait
			nil,    // args
		)

		if err != nil {
			log.Print(err)
		}

		forever := make(chan bool)

		go func() {
			for d := range msgs {
				msg := &rabbitmq.Message{}
				err := json.Unmarshal(d.Body, msg)
				if err != nil {
					log.Print(err)
					continue
				}

				go fn(msg)
			}
		}()

		<-forever
	}()
	return nil
}

func (s *OrderService) PublishOrder(order *rabbitmq.Message) error {
	ch := s.amqp.GetChannel("orderPublish")
	q := s.amqp.GetQueue(ch, "order")

	bytes, err := json.Marshal(order)
	if err != nil {
		log.Fatal("Failed to marshal order: ", err)
		return errors.New("Failed to marshal order: " + err.Error())
	}

	err = ch.Publish(
		"",
		q.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/json",
			Body:        bytes,
		})

	if err != nil {
		log.Fatal("Failed to publish order: ", err)
		return errors.New("Failed to publish order: " + err.Error())
	}

	return nil
}
