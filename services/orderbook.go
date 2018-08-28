package services

import (
	"encoding/json"
	"errors"

	"github.com/Proofsuite/amp-matching-engine/engine"
	"github.com/Proofsuite/amp-matching-engine/utils"
	"github.com/ethereum/go-ethereum/common"

	"github.com/Proofsuite/amp-matching-engine/ws"

	"github.com/gorilla/websocket"

	"github.com/Proofsuite/amp-matching-engine/daos"
)

// PairService struct with daos required, responsible for communicating with daos.
// PairService functions are responsible for interacting with daos and implements business logics.
type OrderBookService struct {
	pairDao  daos.PairDaoInterface
	tokenDao daos.TokenDaoInterface
	eng      engine.EngineInterface
}

type OrderBookServiceInterface interface {
	GetOrderBook(bt, qt common.Address) (ob map[string]interface{}, err error)
	Subscribe(conn *websocket.Conn, bt, qt common.Address)
	Unsubscribe(conn *websocket.Conn, bt, qt common.Address)
}

// NewPairService returns a new instance of balance service
func NewOrderBookService(
	pairDao daos.PairDaoInterface,
	tokenDao daos.TokenDaoInterface,
	eng engine.EngineInterface,
) *OrderBookService {
	return &OrderBookService{pairDao, tokenDao, eng}
}

// Get fetches orderbook from engine/redis and returns it as an map[string]interface
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

// RegisterForOrderBook is responsible for handling incoming orderbook subscription messages
// It makes an entry of connection in pairSocket corresponding to pair,unit and duration
func (s *OrderBookService) Subscribe(conn *websocket.Conn, bt, qt common.Address) {
	socket := ws.GetOrderBookSocket()

	ob, err := s.GetOrderBook(bt, qt)
	if err != nil {
		ws.SendOrderBookErrorMessage(conn, err.Error())
		return
	}

	id := utils.GetOrderBookChannelID(bt, qt)
	err = socket.Subscribe(id, conn)
	if err != nil {
		message := map[string]string{
			"Code":    "UNABLE_TO_REGISTER",
			"Message": "UNABLE_TO_REGISTER " + err.Error(),
		}

		ws.SendOrderBookErrorMessage(conn, message)
		return
	}

	ws.RegisterConnectionUnsubscribeHandler(conn, socket.UnsubscribeHandler(id))
	ws.SendOrderBookInitMessage(conn, ob)
}

// UnRegisterForOrderBook is responsible for handling incoming orderbook unsubscription messages
func (s *OrderBookService) Unsubscribe(conn *websocket.Conn, bt, qt common.Address) {
	socket := ws.GetOrderBookSocket()

	id := utils.GetOrderBookChannelID(bt, qt)
	socket.Unsubscribe(id, conn)
}
