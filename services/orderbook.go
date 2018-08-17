package services

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/Proofsuite/amp-matching-engine/engine"
	"github.com/ethereum/go-ethereum/common"

	"github.com/Proofsuite/amp-matching-engine/ws"

	"github.com/gorilla/websocket"

	"github.com/Proofsuite/amp-matching-engine/daos"
)

// PairService struct with daos required, responsible for communicating with daos.
// PairService functions are responsible for interacting with daos and implements business logics.
type OrderBookService struct {
	pairDao  *daos.PairDao
	tokenDao *daos.TokenDao
	eng      *engine.Resource
}

// NewPairService returns a new instance of balance service
func NewOrderBookService(pairDao *daos.PairDao, tokenDao *daos.TokenDao, eng *engine.Resource) *OrderBookService {
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

	sKey, bKey := res.GetOrderBookKeys()
	fmt.Printf("\n Sell Key: %s \n Buy Key: %s \n", sKey, bKey)

	bids, asks := s.eng.GetOrderBook(res)
	ob = map[string]interface{}{
		"asks": asks,
		"bids": bids,
	}
	return
}

// RegisterForOrderBook is responsible for handling incoming orderbook subscription messages
// It makes an entry of connection in pairSocket corresponding to pair,unit and duration
func (s *OrderBookService) SubscribeOrderBook(conn *websocket.Conn, bt, qt common.Address) {
	ob, err := s.GetOrderBook(bt, qt)
	if err != nil {
		ws.GetPairSockets().SendErrorMessage(conn, err.Error())
		return
	}

	if err := ws.GetPairSockets().Register(bt, qt, conn); err != nil {
		message := map[string]string{
			"Code":    "UNABLE_TO_REGISTER",
			"Message": "UNABLE_TO_REGISTER: " + err.Error(),
		}
		ws.GetPairSockets().SendErrorMessage(conn, message)
		return
	}

	ws.RegisterConnectionUnsubscribeHandler(conn, ws.GetPairSockets().UnsubscribeHandler(bt, qt))
	ws.GetPairSockets().SendBookMessage(conn, ob)
}

// UnRegisterForOrderBook is responsible for handling incoming orderbook unsubscription messages
func (s *OrderBookService) UnSubscribeOrderBook(conn *websocket.Conn, bt, qt common.Address) {
	ws.GetPairSockets().UnregisterConnection(bt, qt, conn)
}
