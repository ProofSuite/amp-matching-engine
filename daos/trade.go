package daos

import (
	"time"

	"github.com/Proofsuite/amp-matching-engine/app"
	"github.com/Proofsuite/amp-matching-engine/types"
	"gopkg.in/mgo.v2/bson"
)

// TradeDao contains:
// collectionName: MongoDB collection name
// dbName: name of mongodb to interact with
type TradeDao struct {
	collectionName string
	dbName         string
}

// NewTradeDao returns a new instance of TradeDao.
func NewTradeDao() *TradeDao {
	return &TradeDao{"trades", app.Config.DBName}
}

// Create function performs the DB insertion task for trade collection
// It accepts 1 or more trades as input.
// All the trades are inserted in one query itself.
func (dao *TradeDao) Create(trades ...*types.Trade) (err error) {
	y := make([]interface{}, len(trades))

	for _, trade := range trades {
		trade.ID = bson.NewObjectId()
		trade.CreatedAt = time.Now()
		trade.UpdatedAt = time.Now()
		y = append(y, trade)
	}

	err = db.Create(dao.dbName, dao.collectionName, y...)
	return
}

// GetAll function fetches all the trades in mongodb
func (dao *TradeDao) GetAll() (response []types.Trade, err error) {
	err = db.Get(dao.dbName, dao.collectionName, bson.M{}, 0, 0, &response)
	return
}

// Aggregate function calls the aggregate pipeline of mongodb
func (dao *TradeDao) Aggregate(q []bson.M) (response []interface{}, err error) {
	return db.Aggregate(dao.dbName, dao.collectionName, q)

}

// GetByPairName fetches all the trades corresponding to a particular pair name.
func (dao *TradeDao) GetByPairName(name string) (response []*types.Trade, err error) {
	q := bson.M{"pairName": bson.RegEx{
		Pattern: name,
		Options: "i",
	}}
	err = db.Get(dao.dbName, dao.collectionName, q, 0, 0, &response)
	if err != nil {
		return
	}
	return
}

// GetByUserAddress fetches all the trades corresponding to a particular user address.
func (dao *TradeDao) GetByUserAddress(addr string) (response []*types.Trade, err error) {
	q := bson.M{"$or": []bson.M{
		bson.M{"maker": bson.RegEx{
			Pattern: addr,
			Options: "i",
		}}, bson.M{"taker": bson.RegEx{
			Pattern: addr,
			Options: "i",
		}},
	}}
	err = db.Get(dao.dbName, dao.collectionName, q, 0, 1, &response)
	if err != nil {
		return
	}
	return
}
