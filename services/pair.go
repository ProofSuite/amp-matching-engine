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
	bt, err := s.tokenDao.GetByAddress(pair.BuyTokenAddress)
	if err != nil {
		return aerrors.InvalidData(map[string]error{"buyToken": errors.New("Token with id " + pair.BuyTokenAddress + " doesn't exists")})
	}
	if bt == nil {
		return aerrors.InvalidData(map[string]error{"buyToken": errors.New("Token with id " + pair.BuyTokenAddress + " doesn't exists")})
	}
	st, err := s.tokenDao.GetByAddress(pair.SellTokenAddress)
	if err != nil {
		return aerrors.InvalidData(map[string]error{"buyToken": errors.New("Token with id " + pair.SellTokenAddress + " doesn't exists")})
	}
	if st == nil {
		return aerrors.InvalidData(map[string]error{"buyToken": errors.New("Token with id " + pair.SellTokenAddress + " doesn't exists")})
	}

	pair.SellTokenSymbol = st.Symbol
	pair.SellToken = st.ID
	pair.SellTokenAddress = st.ContractAddress
	pair.BuyTokenSymbol = bt.Symbol
	pair.BuyToken = bt.ID
	pair.BuyTokenAddress = bt.ContractAddress
	pair.Name = strings.ToUpper(st.Symbol + "-" + bt.Symbol)

	err = s.pairDao.Create(pair)
	return err

}

// GetByID fetches details of a pair using its mongo ID
func (s *PairService) GetByID(id bson.ObjectId) (*types.Pair, error) {
	return s.pairDao.GetByID(id)
}

// GetByTokenAddress fetches details of a pair using contract address of
// its constituting tokens
func (s *PairService) GetByTokenAddress(bt, st string) (*types.Pair, error) {
	return s.pairDao.GetByTokenAddress(bt, st)
}

// GetAll is reponsible for fetching all the pairs in the DB
func (s *PairService) GetAll() ([]types.Pair, error) {
	return s.pairDao.GetAll()
}

// GetOrderBook fetches orderbook from engine/redis and returns it as an map[string]interface
func (s *PairService) GetOrderBook(pairName string) (ob map[string]interface{}, err error) {
	res, err := s.pairDao.GetByName(pairName)
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
func (s *PairService) RegisterForOrderBook(conn *websocket.Conn, pairName string) {

	ob, err := s.GetOrderBook(pairName)
	if err != nil {
		conn.WriteMessage(1, []byte(err.Error()))
	}
	trades, _ := s.tradeService.GetByPairName(pairName)
	ob["trades"] = trades

	if err := ws.GetPairSockets().PairSocketRegister(pairName, conn); err != nil {
		message := map[string]string{
			"Code":    "UNABLE_TO_REGISTER",
			"Message": "UNABLE_TO_REGISTER: " + err.Error(),
		}
		mab, _ := json.Marshal(message)
		conn.WriteMessage(1, mab)
	}
	ws.RegisterConnectionUnsubscribeHandler(conn, ws.GetPairSockets().PairUnsubscribeHandler(pairName))

	rab, _ := json.Marshal(ob)
	conn.WriteMessage(1, rab)
}

// UnRegisterForOrderBook is responsible for handling incoming orderbook unsubscription messages
func (s *PairService) UnRegisterForOrderBook(conn *websocket.Conn, pairName string) {
	ws.GetPairSockets().PairSocketUnregisterConnection(pairName, conn)
}
