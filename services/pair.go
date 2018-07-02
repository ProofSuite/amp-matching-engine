package services

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/Proofsuite/amp-matching-engine/engine"

	"github.com/Proofsuite/amp-matching-engine/ws"

	"github.com/gorilla/websocket"

	"labix.org/v2/mgo/bson"

	"github.com/Proofsuite/amp-matching-engine/daos"
	aerrors "github.com/Proofsuite/amp-matching-engine/errors"
	"github.com/Proofsuite/amp-matching-engine/types"
)

type PairService struct {
	pairDao      *daos.PairDao
	tokenDao     *daos.TokenDao
	eng          *engine.EngineResource
	tradeService *TradeService
}

func NewPairService(pairDao *daos.PairDao, tokenDao *daos.TokenDao, eng *engine.EngineResource, tradeService *TradeService) *PairService {

	return &PairService{pairDao, tokenDao, eng, tradeService}
}

func (s *PairService) Create(pair *types.Pair) error {
	bt, err := s.tokenDao.GetByID(pair.BuyToken)
	if err != nil {
		return aerrors.InvalidData(map[string]error{"buyToken": errors.New("Token with id " + pair.BuyToken.Hex() + " doesn't exists")})
	}
	st, err := s.tokenDao.GetByID(pair.SellToken)
	if err != nil {
		return aerrors.InvalidData(map[string]error{"buyToken": errors.New("Token with id " + pair.SellToken.Hex() + " doesn't exists")})
	}
	pair.SellTokenSymbol = st.Symbol
	pair.SellTokenAddress = st.ContractAddress
	pair.BuyTokenSymbol = bt.Symbol
	pair.BuyTokenAddress = bt.ContractAddress

	err = s.pairDao.Create(pair)
	return err

}

func (s *PairService) GetByID(id bson.ObjectId) (*types.Pair, error) {
	return s.pairDao.GetByID(id)
}

func (s *PairService) GetAll() ([]types.Pair, error) {
	return s.pairDao.GetAll()
}

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
func (s *PairService) RegisterForOrderBook(conn *websocket.Conn, pairName string) {

	ob, err := s.GetOrderBook(pairName)
	if err != nil {
		conn.WriteMessage(1, []byte(err.Error()))
		conn.Close()
	}
	trades, _ := s.tradeService.GetByPairName(pairName)
	ob["trades"] = trades

	if err := ws.PairSocketRegister(pairName, conn); err != nil {
		message := map[string]string{
			"Code":    "UNABLE_TO_REGISTER",
			"Message": "UNABLE_TO_REGISTER: " + err.Error(),
		}
		mab, _ := json.Marshal(message)
		conn.WriteMessage(1, mab)
		conn.Close()
	}
	conn.SetCloseHandler(ws.PairSocketCloseHandler(pairName, conn))

	go func() {
		for {
			if _, _, err := conn.NextReader(); err != nil {
				conn.Close()
				break
			}
		}
	}()

	rab, _ := json.Marshal(ob)
	conn.WriteMessage(1, rab)
}
