package daos

import (
	"log"
	"time"

	"github.com/Proofsuite/amp-matching-engine/app"
	"github.com/Proofsuite/amp-matching-engine/types"
	"github.com/ethereum/go-ethereum/common"
	mgo "gopkg.in/mgo.v2"
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
	dbName := app.Config.DBName
	collection := "trades"
	index := mgo.Index{
		Key:    []string{"hash"},
		Sparse: true,
	}

	err := db.Session.DB(dbName).C(collection).EnsureIndex(index)
	if err != nil {
		panic(err)
	}
	return &TradeDao{collection, dbName}
}

// Create function performs the DB insertion task for trade collection
// It accepts 1 or more trades as input.
// All the trades are inserted in one query itself.
func (dao *TradeDao) Create(trades ...*types.Trade) error {
	y := make([]interface{}, len(trades))

	for _, trade := range trades {
		trade.ID = bson.NewObjectId()
		trade.CreatedAt = time.Now()
		trade.UpdatedAt = time.Now()
		y = append(y, trade)
	}

	err := db.Create(dao.dbName, dao.collectionName, y...)
	if err != nil {
		log.Print(err)
		return err
	}

	return nil
}

func (dao *TradeDao) Update(trade *types.Trade) error {
	trade.UpdatedAt = time.Now()
	err := db.Update(dao.dbName, dao.collectionName, bson.M{"_id": trade.ID}, trade)
	return err
}

// UpdateByHash updates the fields that can be normally updated in a structure. For a
// complete update, use the Update or UpdateAllByHash function
func (dao *TradeDao) UpdateByHash(hash common.Hash, t *types.Trade) error {
	t.UpdatedAt = time.Now()
	query := bson.M{"hash": hash.Hex()}
	update := bson.M{"$set": bson.M{
		"price":      t.Price.String(),
		"pricepoint": t.PricePoint.String(),
		"tradeNonce": t.TradeNonce.String(),
		"txHash":     t.TxHash.String(),
		"signature": &types.SignatureRecord{
			V: t.Signature.V,
			R: t.Signature.R.Hex(),
			S: t.Signature.S.Hex(),
		},
		"updatedAt": t.UpdatedAt,
	}}

	err := db.Update(dao.dbName, dao.collectionName, query, update)
	if err != nil {
		return err
	}

	return nil
}

// GetAll function fetches all the trades in mongodb
func (dao *TradeDao) GetAll() ([]types.Trade, error) {
	var response []types.Trade
	err := db.Get(dao.dbName, dao.collectionName, bson.M{}, 0, 0, &response)
	return response, err
}

// Aggregate function calls the aggregate pipeline of mongodb
func (dao *TradeDao) Aggregate(q []bson.M) ([]*types.Tick, error) {
	var response []*types.Tick
	err := db.Aggregate(dao.dbName, dao.collectionName, q, &response)
	return response, err
}

// GetByPairName fetches all the trades corresponding to a particular pair name.
func (dao *TradeDao) GetByPairName(name string) ([]*types.Trade, error) {
	var response []*types.Trade
	q := bson.M{"pairName": bson.RegEx{
		Pattern: name,
		Options: "i",
	}}

	err := db.Get(dao.dbName, dao.collectionName, q, 0, 0, &response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// GetByHash fetches the first record that matches a certain hash
func (dao *TradeDao) GetByHash(hash common.Hash) (*types.Trade, error) {
	q := bson.M{"hash": hash.Hex()}

	response := []*types.Trade{}
	err := db.Get(dao.dbName, dao.collectionName, q, 0, 0, &response)
	if err != nil {
		return nil, err
	}

	return response[0], nil
}

// GetByOrderHash fetches the first trade record which matches a certain order hash
func (dao *TradeDao) GetByOrderHash(hash common.Hash) ([]*types.Trade, error) {
	q := bson.M{"orderHash": hash.Hex()}

	response := []*types.Trade{}
	err := db.Get(dao.dbName, dao.collectionName, q, 0, 0, &response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// GetByPairAddress fetches all the trades corresponding to a particular pair token address.
func (dao *TradeDao) GetByPairAddress(baseToken, quoteToken common.Address) ([]*types.Trade, error) {
	var response []*types.Trade

	q := bson.M{"baseToken": baseToken.Hex(), "quoteToken": quoteToken.Hex()}
	err := db.Get(dao.dbName, dao.collectionName, q, 0, 0, &response)
	return response, err
}

// GetByUserAddress fetches all the trades corresponding to a particular user address.
func (dao *TradeDao) GetByUserAddress(addr common.Address) ([]*types.Trade, error) {
	var response []*types.Trade
	q := bson.M{"$or": []bson.M{
		{"maker": addr.Hex()}, {"taker": addr.Hex()},
	}}
	err := db.Get(dao.dbName, dao.collectionName, q, 0, 1, &response)
	return response, err
}

// Drop drops all the order documents in the current database
func (dao *TradeDao) Drop() {
	db.DropCollection(dao.dbName, dao.collectionName)
}
