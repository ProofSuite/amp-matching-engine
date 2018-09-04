package services

import (
	"encoding/json"
	"errors"

	"github.com/Proofsuite/amp-matching-engine/interfaces"
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
func (s *OrderBookService) GetOrderBook(bt, qt common.Address) (ob map[string]interface{}, err error) {
	res, err := s.pairDao.GetByTokenAddress(bt, qt)
	if err != nil {
		message := map[string]string{
			"Code":    "Invalid_Pair",
			"Message": "Invalid Pair " + err.Error(),
		}
		bytes, _ := json.Marshal(message)
		return nil, errors.New(string(bytes))
	}

	// sKey, bKey := res.GetOrderBookKeys()

	bids, asks := s.eng.GetOrderBook(res)
	ob = map[string]interface{}{
		"asks": asks,
		"bids": bids,
	}
	return
}

// SubscribeLite is responsible for handling incoming orderbook subscription messages
// It makes an entry of connection in pairSocket corresponding to pair,unit and duration
func (s *OrderBookService) SubscribeLite(conn *ws.Conn, bt, qt common.Address) {
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
			"Code":    "UNABLE_TO_REGISTER",
			"Message": "UNABLE_TO_REGISTER " + err.Error(),
		}

		socket.SendErrorMessage(conn, message)
		return
	}

	ws.RegisterConnectionUnsubscribeHandler(conn, socket.UnsubscribeHandler(id))
	socket.SendInitMessage(conn, ob)
}

// UnsubscribeLite is responsible for handling incoming orderbook unsubscription messages
func (s *OrderBookService) UnsubscribeLite(conn *ws.Conn, bt, qt common.Address) {
	socket := ws.GetLiteOrderBookSocket()

	id := utils.GetOrderBookChannelID(bt, qt)
	socket.Unsubscribe(id, conn)
}

// GetFullOrderBook fetches complete orderbook from engine/redis
func (s *OrderBookService) GetFullOrderBook(bt, qt common.Address) (ob [][]string, err error) {
	res, err := s.pairDao.GetByTokenAddress(bt, qt)
	if err != nil {
		message := map[string]string{
			"Code":    "Invalid_Pair",
			"Message": "Invalid Pair " + err.Error(),
		}
		bytes, _ := json.Marshal(message)
		return nil, errors.New(string(bytes))
	}

	return s.eng.GetFullOrderBook(res), nil
}

// SubscribeFull is responsible for handling incoming orderbook subscription messages
// It makes an entry of connection in pairSocket corresponding to pair,unit and duration
func (s *OrderBookService) SubscribeFull(conn *ws.Conn, bt, qt common.Address) {
	socket := ws.GetFullOrderBookSocket()

	ob, err := s.GetFullOrderBook(bt, qt)
	if err != nil {
		socket.SendErrorMessage(conn, err.Error())
		return
	}

	id := utils.GetOrderBookChannelID(bt, qt)
	err = socket.Subscribe(id, conn)
	if err != nil {
		message := map[string]string{
			"Code":    "UNABLE_TO_REGISTER",
			"Message": "UNABLE_TO_REGISTER " + err.Error(),
		}

		socket.SendErrorMessage(conn, message)
		return
	}

	ws.RegisterConnectionUnsubscribeHandler(conn, socket.UnsubscribeHandler(id))
	socket.SendInitMessage(conn, ob)
}

// UnsubscribeFull is responsible for handling incoming orderbook unsubscription messages
func (s *OrderBookService) UnsubscribeFull(conn *ws.Conn, bt, qt common.Address) {
	socket := ws.GetFullOrderBookSocket()

	id := utils.GetOrderBookChannelID(bt, qt)
	socket.Unsubscribe(id, conn)
}
