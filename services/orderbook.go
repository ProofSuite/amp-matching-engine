package services

import (
	"errors"

	"github.com/Proofsuite/amp-matching-engine/interfaces"
	"github.com/Proofsuite/amp-matching-engine/types"
	"github.com/Proofsuite/amp-matching-engine/utils"
	"github.com/ethereum/go-ethereum/common"

	"github.com/Proofsuite/amp-matching-engine/ws"
)

// PairService struct with daos required, responsible for communicating with daos.
// PairService functions are responsible for interacting with daos and implements business logics.
type OrderBookService struct {
	pairDao  interfaces.PairDao
	tokenDao interfaces.TokenDao
	orderDao interfaces.OrderDao
	eng      interfaces.Engine
}

// NewPairService returns a new instance of balance service
func NewOrderBookService(
	pairDao interfaces.PairDao,
	tokenDao interfaces.TokenDao,
	orderDao interfaces.OrderDao,
	eng interfaces.Engine,
) *OrderBookService {
	return &OrderBookService{pairDao, tokenDao, orderDao, eng}
}

// GetOrderBook fetches orderbook from engine/redis and returns it as an map[string]interface
func (s *OrderBookService) GetOrderBook(bt, qt common.Address) (map[string]interface{}, error) {
	pair, err := s.pairDao.GetByTokenAddress(bt, qt)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	if pair == nil {
		return nil, errors.New("Pair not found")
	}

	bids, asks, err := s.orderDao.GetOrderBook(pair)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	ob := map[string]interface{}{
		"asks": asks,
		"bids": bids,
	}

	return ob, nil
}

// SubscribeOrderBook is responsible for handling incoming orderbook subscription messages
// It makes an entry of connection in pairSocket corresponding to pair,unit and duration
func (s *OrderBookService) SubscribeOrderBook(c *ws.Client, bt, qt common.Address) {
	socket := ws.GetOrderBookSocket()

	ob, err := s.GetOrderBook(bt, qt)
	if err != nil {
		socket.SendErrorMessage(c, err.Error())
		return
	}

	id := utils.GetOrderBookChannelID(bt, qt)
	err = socket.Subscribe(id, c)
	if err != nil {
		msg := map[string]string{"Message": err.Error()}
		socket.SendErrorMessage(c, msg)
		return
	}

	ws.RegisterConnectionUnsubscribeHandler(c, socket.UnsubscribeHandler(id))
	socket.SendInitMessage(c, ob)
}

// UnsubscribeOrderBook is responsible for handling incoming orderbook unsubscription messages
func (s *OrderBookService) UnsubscribeOrderBook(c *ws.Client) {
	socket := ws.GetOrderBookSocket()
	socket.Unsubscribe(c)
}

func (s *OrderBookService) UnsubscribeOrderBookChannel(c *ws.Client, bt, qt common.Address) {
	socket := ws.GetOrderBookSocket()
	id := utils.GetOrderBookChannelID(bt, qt)
	socket.UnsubscribeChannel(id, c)
}

// GetRawOrderBook fetches complete orderbook from engine/redis
func (s *OrderBookService) GetRawOrderBook(bt, qt common.Address) ([]*types.Order, error) {
	pair, err := s.pairDao.GetByTokenAddress(bt, qt)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	ob, err := s.orderDao.GetRawOrderBook(pair)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	return ob, nil
}

// SubscribeRawOrderBook is responsible for handling incoming orderbook subscription messages
// It makes an entry of connection in pairSocket corresponding to pair,unit and duration
func (s *OrderBookService) SubscribeRawOrderBook(c *ws.Client, bt, qt common.Address) {
	socket := ws.GetRawOrderBookSocket()

	ob, err := s.GetRawOrderBook(bt, qt)
	if err != nil {
		socket.SendErrorMessage(c, err.Error())
		return
	}

	id := utils.GetOrderBookChannelID(bt, qt)
	err = socket.Subscribe(id, c)
	if err != nil {
		msg := map[string]string{"Message": err.Error()}
		socket.SendErrorMessage(c, msg)
		return
	}

	ws.RegisterConnectionUnsubscribeHandler(c, socket.UnsubscribeChannelHandler(id))
	socket.SendInitMessage(c, ob)
}

// UnsubscribeRawOrderBook is responsible for handling incoming orderbook unsubscription messages
func (s *OrderBookService) UnsubscribeRawOrderBook(c *ws.Client) {
	socket := ws.GetRawOrderBookSocket()
	socket.Unsubscribe(c)
}

func (s *OrderBookService) UnsubscribeRawOrderBookChannel(c *ws.Client, bt, qt common.Address) {
	socket := ws.GetRawOrderBookSocket()
	id := utils.GetOrderBookChannelID(bt, qt)
	socket.UnsubscribeChannel(id, c)
}
