package daos

import (
	"time"

	"github.com/Proofsuite/amp-matching-engine/app"
	"github.com/Proofsuite/amp-matching-engine/types"
	"labix.org/v2/mgo/bson"
)

type TradeDao struct {
	collectionName string
	dbName         string
}

func NewTradeDao() *TradeDao {
	return &TradeDao{"trades", app.Config.DBName}
}

func (dao *TradeDao) Create(trades ...*types.Trade) (err error) {
	y := make([]interface{}, len(trades))

	for _, trade := range trades {
		trade.ID = bson.NewObjectId()
		trade.CreatedAt = time.Now()
		trade.UpdatedAt = time.Now()
		y = append(y, trade)
	}

	err = DB.Create(dao.dbName, dao.collectionName, y...)
	return
}

func (dao *TradeDao) GetAll() (response []types.Trade, err error) {
	err = DB.Get(dao.dbName, dao.collectionName, bson.M{}, 0, 0, &response)
	return
}
func (dao *TradeDao) Update(id bson.ObjectId, trade *types.Trade) (response []types.Trade, err error) {
	trade.UpdatedAt = time.Now()
	err = DB.Update(dao.dbName, dao.collectionName, bson.M{"_id": id}, trade)
	return
}
func (dao *TradeDao) GetByID(id bson.ObjectId) (response *types.Trade, err error) {
	err = DB.GetByID(dao.dbName, dao.collectionName, id, &response)
	return
}
func (dao *TradeDao) GetByPairName(name string) (response []*types.Trade, err error) {
	q := bson.M{"pairName": bson.RegEx{
		Pattern: name,
		Options: "i",
	}}
	err = DB.Get(dao.dbName, dao.collectionName, q, 0, 1, &response)
	if err != nil {
		return
	}
	return
}
