package services

import (
	"encoding/json"
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
	eng      interfaces.Engine
}

// NewPairService returns a new instance of balance service
func NewOrderBookService(
	pairDao interfaces.PairDao,
	tokenDao interfaces.TokenDao,
	eng interfaces.Engine,
) *OrderBookService {
	return &OrderBookService{pairDao, tokenDao, eng}
}

// GetOrderBook fetches orderbook from engine/redis and returns it as an map[string]interface
func (s *OrderBookService) GetOrderBook(bt, qt common.Address) (map[string]interface{}, error) {
	res, err := s.pairDao.GetByTokenAddress(bt, qt)
	if err != nil {
		message := map[string]string{
			"Code":    "Invalid Pair",
			"Message": err.Error(),
		}
		bytes, _ := json.Marshal(message)
		return nil, errors.New(string(bytes))
	}

	bids, asks, err := s.eng.GetOrderBook(res)
	if err != nil {
		message := map[string]string{
			"Code":    "Invalid Pair",
			"Message": err.Error(),
		}

		bytes, _ := json.Marshal(message)
		return nil, errors.New(string(bytes))
	}

	ob := map[string]interface{}{
		"asks": asks,
		"bids": bids,
	}
	return ob, nil
}

// SubscribeOrderBook is responsible for handling incoming orderbook subscription messages
// It makes an entry of connection in pairSocket corresponding to pair,unit and duration
func (s *OrderBookService) SubscribeOrderBook(conn *ws.Conn, bt, qt common.Address) {
	socket := ws.GetLiteOrderBookSocket()

	ob, err := s.GetOrderBook(bt, qt)
	if err != nil {
		socket.SendErrorMessage(conn, err.Error())
		return
	}

	id := utils.GetOrderBookChannelID(bt, qt)
	err = socket.Subscribe(id, conn)
	if err != nil {
		message := map[string]string{
			"Code":    "Internal Server Error",
			"Message": err.Error(),
		}

		socket.SendErrorMessage(conn, message)
		return
	}

	ws.RegisterConnectionUnsubscribeHandler(conn, socket.UnsubscribeHandler(id))
	socket.SendInitMessage(conn, ob)
}

// UnSubscribeOrderBook is responsible for handling incoming orderbook unsubscription messages
func (s *OrderBookService) UnSubscribeOrderBook(conn *ws.Conn, bt, qt common.Address) {
	socket := ws.GetLiteOrderBookSocket()

	id := utils.GetOrderBookChannelID(bt, qt)
	socket.Unsubscribe(id, conn)
}

// GetRawOrderBook fetches complete orderbook from engine/redis
func (s *OrderBookService) GetRawOrderBook(bt, qt common.Address) ([][]types.Order, error) {
	res, err := s.pairDao.GetByTokenAddress(bt, qt)
	if err != nil {
		message := map[string]string{
			"Code":    "Invalid Pair",
			"Message": err.Error(),
		}
		bytes, _ := json.Marshal(message)
		return nil, errors.New(string(bytes))
	}

	book, err := s.eng.GetRawOrderBook(res)
	if err != nil {
		message := map[string]string{
			"Code":    "Internal Server Error",
			"Message": err.Error(),
		}
		bytes, _ := json.Marshal(message)
		return nil, errors.New(string(bytes))
	}

	return book, nil
}

// SubscribeRawOrderBook is responsible for handling incoming orderbook subscription messages
// It makes an entry of connection in pairSocket corresponding to pair,unit and duration
func (s *OrderBookService) SubscribeRawOrderBook(conn *ws.Conn, bt, qt common.Address) {
	socket := ws.GetRawOrderBookSocket()

	ob, err := s.GetRawOrderBook(bt, qt)
	if err != nil {
		socket.SendErrorMessage(conn, err.Error())
		return
	}

	id := utils.GetOrderBookChannelID(bt, qt)
	err = socket.Subscribe(id, conn)
	if err != nil {
		message := map[string]string{
			"Code":    "Internal Server Error",
			"Message": err.Error(),
		}

		socket.SendErrorMessage(conn, message)
		return
	}

	ws.RegisterConnectionUnsubscribeHandler(conn, socket.UnsubscribeHandler(id))
	socket.SendInitMessage(conn, ob)
}

// UnSubscribeRawOrderBook is responsible for handling incoming orderbook unsubscription messages
func (s *OrderBookService) UnSubscribeRawOrderBook(conn *ws.Conn, bt, qt common.Address) {
	socket := ws.GetRawOrderBookSocket()

	id := utils.GetOrderBookChannelID(bt, qt)
	socket.Unsubscribe(id, conn)
}
