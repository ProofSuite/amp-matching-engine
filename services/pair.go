package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/Proofsuite/amp-matching-engine/engine"

	"github.com/Proofsuite/amp-matching-engine/ws"

	"github.com/gorilla/websocket"

	"gopkg.in/mgo.v2/bson"

	"github.com/Proofsuite/amp-matching-engine/daos"
	aerrors "github.com/Proofsuite/amp-matching-engine/errors"
	"github.com/Proofsuite/amp-matching-engine/types"
)

// PairService struct with daos required, responsible for communicating with daos.
// PairService functions are responsible for interacting with daos and implements business logics.
type PairService struct {
	pairDao      *daos.PairDao
	tokenDao     *daos.TokenDao
	eng          *engine.Resource
	tradeService *TradeService
}

// NewPairService returns a new instance of balance service
func NewPairService(pairDao *daos.PairDao, tokenDao *daos.TokenDao, eng *engine.Resource, tradeService *TradeService) *PairService {

	return &PairService{pairDao, tokenDao, eng, tradeService}
}

// Create function is responsible for inserting new pair in DB.
// It checks for existence of tokens in DB first
func (s *PairService) Create(pair *types.Pair) error {
	p, err := s.GetByTokenAddress(pair.BaseTokenAddress, pair.QuoteTokenAddress)
	if err != nil && err.Error() != "No Pair found" {
		return aerrors.NewAPIError(400, err.Error(), nil)
	} else if p != nil {
		return aerrors.NewAPIError(401, "PAIR_ALREADY_EXISTS", nil)
	}

	p, err = s.GetByTokenAddress(pair.QuoteTokenAddress, pair.BaseTokenAddress)
	if err != nil && err.Error() != "No Pair found" {
		return aerrors.NewAPIError(400, err.Error(), nil)
	} else if p != nil {
		return aerrors.NewAPIError(401, "PAIR_ALREADY_EXISTS", nil)
	}

	bt, err := s.tokenDao.GetByAddress(pair.BaseTokenAddress)
	if err != nil {
		return aerrors.NewAPIError(400, err.Error(), nil)
	}
	if bt == nil {
		return aerrors.NewAPIError(401, "BaseTokenAddress_DOESNT_EXISTS", nil)
	}

	st, err := s.tokenDao.GetByAddress(pair.QuoteTokenAddress)
	if err != nil {
		return aerrors.NewAPIError(400, err.Error(), nil)
	}
	if st == nil {
		return aerrors.NewAPIError(401, "QuoteTokenAddress_DOESNT_EXISTS", nil)
	}
	if !st.Quote {
		return aerrors.NewAPIError(401, "QuoteTokenAddress_CAN_NOT_BE_USED_AS_QUOTE_TOKEN", nil)
	}

	pair.QuoteTokenSymbol = st.Symbol
	pair.QuoteToken = st.ID
	pair.QuoteTokenAddress = st.ContractAddress
	pair.BaseTokenSymbol = bt.Symbol
	pair.BaseToken = bt.ID
	pair.BaseTokenAddress = bt.ContractAddress
	pair.Name = strings.ToUpper(st.Symbol + "/" + bt.Symbol)

	err = s.pairDao.Create(pair)
	return err

}

// GetByID fetches details of a pair using its mongo ID
func (s *PairService) GetByID(id bson.ObjectId) (*types.Pair, error) {
	return s.pairDao.GetByID(id)
}

// GetByTokenAddress fetches details of a pair using contract address of
// its constituting tokens
func (s *PairService) GetByTokenAddress(bt, qt string) (*types.Pair, error) {
	return s.pairDao.GetByTokenAddress(bt, qt)
}

// GetAll is reponsible for fetching all the pairs in the DB
func (s *PairService) GetAll() ([]types.Pair, error) {
	return s.pairDao.GetAll()
}

// GetOrderBook fetches orderbook from engine/redis and returns it as an map[string]interface
func (s *PairService) GetOrderBook(bt, qt string) (ob map[string]interface{}, err error) {
	res, err := s.GetByTokenAddress(bt, qt)
	if err != nil {
		message := map[string]string{
			"Code":    "Invalid_Pair_Name",
			"Message": "Invalid Pair Name " + err.Error(),
		}
		mab, _ := json.Marshal(message)
		return nil, errors.New(string(mab))
	}
	sKey, bKey := res.GetOrderBookKeys()
	fmt.Printf("\n Sell Key: %s \n Buy Key: %s \n", sKey, bKey)

	sBook, bBook := s.eng.GetOrderBook(res)
	ob = map[string]interface{}{
		"sell": sBook,
		"buy":  bBook,
	}
	return
}

// RegisterForOrderBook is responsible for handling incoming orderbook subscription messages
// It makes an entry of connection in pairSocket corresponding to pair,unit and duration
func (s *PairService) RegisterForOrderBook(conn *websocket.Conn, bt, qt string) {

	ob, err := s.GetOrderBook(bt, qt)
	if err != nil {
		conn.WriteMessage(1, []byte(err.Error()))
	}
	trades, _ := s.tradeService.GetByPairAddress(bt, qt)
	ob["trades"] = trades

	if err := ws.GetPairSockets().PairSocketRegister(bt, qt, conn); err != nil {
		message := map[string]string{
			"Code":    "UNABLE_TO_REGISTER",
			"Message": "UNABLE_TO_REGISTER: " + err.Error(),
		}
		mab, _ := json.Marshal(message)
		conn.WriteJSON(mab)
	}
	ws.RegisterConnectionUnsubscribeHandler(conn, ws.GetPairSockets().PairUnsubscribeHandler(bt, qt))

	response := types.Message{
		MsgType: "order_book",
		Data:    ob,
	}
	rab, _ := json.Marshal(response)
	conn.WriteJSON(rab)
}

// UnRegisterForOrderBook is responsible for handling incoming orderbook unsubscription messages
func (s *PairService) UnRegisterForOrderBook(conn *websocket.Conn, bt, qt string) {
	ws.GetPairSockets().PairSocketUnregisterConnection(bt, qt, conn)
}
