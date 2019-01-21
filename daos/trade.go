package daos

import (
	"time"

	"github.com/Proofsuite/amp-matching-engine/app"
	"github.com/Proofsuite/amp-matching-engine/types"
	"github.com/ethereum/go-ethereum/common"
	mgo "github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
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

	i1 := mgo.Index{
		Key: []string{"baseToken"},
	}

	i2 := mgo.Index{
		Key: []string{"quoteToken"},
	}

	i3 := mgo.Index{
		Key: []string{"createdAt"},
	}

	i4 := mgo.Index{
		Key:    []string{"hash"},
		Sparse: true,
	}

	i5 := mgo.Index{
		Key:    []string{"makerOrderHash"},
		Sparse: true,
	}

	i6 := mgo.Index{
		Key:    []string{"takerOrderHash"},
		Sparse: true,
	}

	i7 := mgo.Index{
		Key: []string{"createdAt", "status", "baseToken", "quoteToken"},
	}

	i8 := mgo.Index{
		Key:       []string{"pricepoint"},
		Collation: &mgo.Collation{NumericOrdering: true, Locale: "en"},
	}

	err := db.Session.DB(dbName).C(collection).EnsureIndex(i1)
	if err != nil {
		panic(err)
	}

	err = db.Session.DB(dbName).C(collection).EnsureIndex(i2)
	if err != nil {
		panic(err)
	}

	err = db.Session.DB(dbName).C(collection).EnsureIndex(i3)
	if err != nil {
		panic(err)
	}

	err = db.Session.DB(dbName).C(collection).EnsureIndex(i4)
	if err != nil {
		panic(err)
	}

	err = db.Session.DB(dbName).C(collection).EnsureIndex(i5)
	if err != nil {
		panic(err)
	}

	err = db.Session.DB(dbName).C(collection).EnsureIndex(i6)
	if err != nil {
		panic(err)
	}

	err = db.Session.DB(dbName).C(collection).EnsureIndex(i7)
	if err != nil {
		panic(err)
	}

	err = db.Session.DB(dbName).C(collection).EnsureIndex(i8)
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
		logger.Error(err)
		return err
	}

	return nil
}

func (dao *TradeDao) DeleteByHashes(hashes ...common.Hash) error {
	err := db.RemoveAll(dao.dbName, dao.collectionName, bson.M{"hash": bson.M{"$in": hashes}})
	if err != nil {
		logger.Error(err)
		return err
	}

	return nil
}

func (dao *TradeDao) Delete(trades ...*types.Trade) error {
	hashes := []common.Hash{}
	for _, t := range trades {
		hashes = append(hashes, t.Hash)
	}

	err := db.RemoveAll(dao.dbName, dao.collectionName, bson.M{"hash": bson.M{"$in": hashes}})
	if err != nil {
		logger.Error(err)
		return err
	}

	return nil
}

func (dao *TradeDao) Update(trade *types.Trade) error {
	trade.UpdatedAt = time.Now()
	err := db.Update(dao.dbName, dao.collectionName, bson.M{"_id": trade.ID}, trade)
	if err != nil {
		logger.Error(err)
		return err
	}

	return nil
}

func (dao *TradeDao) Upsert(id bson.ObjectId, t *types.Trade) error {
	t.UpdatedAt = time.Now()

	err := db.Upsert(dao.dbName, dao.collectionName, bson.M{"_id": id}, t)
	if err != nil {
		logger.Error(err)
		return err
	}

	return nil
}

func (dao *TradeDao) UpsertByHash(h common.Hash, t *types.Trade) error {
	t.UpdatedAt = time.Now()

	err := db.Upsert(dao.dbName, dao.collectionName, bson.M{"hash": h}, t)
	if err != nil {
		logger.Error(err)
		return err
	}

	return nil
}

func (dao *TradeDao) FindAndModify(h common.Hash, t *types.Trade) (*types.Trade, error) {
	t.UpdatedAt = time.Now()
	query := bson.M{"hash": h.Hex()}
	updated := &types.Trade{}
	change := mgo.Change{
		Update:    types.TradeBSONUpdate{t},
		Upsert:    true,
		Remove:    false,
		ReturnNew: true,
	}

	err := db.FindAndModify(dao.dbName, dao.collectionName, query, change, &updated)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	return updated, nil
}

// UpdateByHash updates the fields that can be normally updated in a structure. For a
// complete update, use the Update or UpdateAllByHash function
func (dao *TradeDao) UpdateByHash(h common.Hash, t *types.Trade) error {
	t.UpdatedAt = time.Now()
	query := bson.M{"hash": h.Hex()}
	update := bson.M{"$set": bson.M{
		"pricepoint":     t.PricePoint.String(),
		"amount":         t.Amount.String(),
		"txHash":         t.TxHash.String(),
		"takerOrderHash": t.TakerOrderHash.String(),
		"makerOrderHash": t.MakerOrderHash.String(),
		"updatedAt":      t.UpdatedAt,
	}}

	err := db.Update(dao.dbName, dao.collectionName, query, update)
	if err != nil {
		logger.Error(err)
		return err
	}

	return nil
}

// GetAll function fetches all the trades in mongodb
func (dao *TradeDao) GetAll() ([]types.Trade, error) {
	var response []types.Trade
	err := db.Get(dao.dbName, dao.collectionName, bson.M{}, 0, 0, &response)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	return response, nil
}

func (dao *TradeDao) GetErroredTradeCount(start, end time.Time) (int, error) {
	q := bson.M{
		"status": bson.M{"$in": []string{"ERROR"}},
		"createdAt": bson.M{
			"$gte": start,
			"$lt":  end,
		},
	}

	n, err := db.Count(dao.dbName, dao.collectionName, q)
	if err != nil {
		logger.Error(err)
		return 0, err
	}

	return n, nil
}

// Aggregate function calls the aggregate pipeline of mongodb
func (dao *TradeDao) Aggregate(q []bson.M) ([]*types.Tick, error) {
	var res []*types.Tick

	err := db.Aggregate(dao.dbName, dao.collectionName, q, &res)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	return res, nil
}

// GetByPairName fetches all the trades corresponding to a particular pair name.
func (dao *TradeDao) GetByPairName(name string) ([]*types.Trade, error) {
	var res []*types.Trade
	q := bson.M{"pairName": bson.RegEx{
		Pattern: name,
		Options: "i",
	}}

	err := db.Get(dao.dbName, dao.collectionName, q, 0, 0, &res)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	return res, nil
}

// GetByHash fetches the first record that matches a certain hash
func (dao *TradeDao) GetByHash(h common.Hash) (*types.Trade, error) {
	q := bson.M{"hash": h.Hex()}

	res := []*types.Trade{}
	err := db.Get(dao.dbName, dao.collectionName, q, 0, 0, &res)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	return res[0], nil
}

// GetByOrderHash fetches the first trade record which matches a certain order hash
func (dao *TradeDao) GetByMakerOrderHash(h common.Hash) ([]*types.Trade, error) {
	q := bson.M{"makerOrderHash": h.Hex()}

	res := []*types.Trade{}
	err := db.Get(dao.dbName, dao.collectionName, q, 0, 0, &res)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	return res, nil
}

func (dao *TradeDao) GetByTakerOrderHash(h common.Hash) ([]*types.Trade, error) {
	q := bson.M{"takerOrderHash": h.Hex()}

	res := []*types.Trade{}
	err := db.Get(dao.dbName, dao.collectionName, q, 0, 0, &res)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	return res, nil
}

func (dao *TradeDao) GetByOrderHashes(hashes []common.Hash) ([]*types.Trade, error) {
	hexes := []string{}
	for _, h := range hashes {
		hexes = append(hexes, h.Hex())
	}

	q := bson.M{"makerOrderHash": bson.M{"$in": hexes}}
	res := []*types.Trade{}

	err := db.Get(dao.dbName, dao.collectionName, q, 0, 0, &res)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	return res, nil
}

func (dao *TradeDao) GetSortedTrades(bt, qt common.Address, n int) ([]*types.Trade, error) {
	res := []*types.Trade{}

	q := bson.M{"baseToken": bt.Hex(), "quoteToken": qt.Hex()}
	sort := []string{"-createdAt"}
	err := db.GetAndSort(dao.dbName, dao.collectionName, q, sort, 0, n, &res)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	return res, nil
}

func (dao *TradeDao) GetNTradesByPairAddress(bt, qt common.Address, n int) ([]*types.Trade, error) {
	res, err := dao.GetTradesByPairAddress(bt, qt, n)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	return res, nil
}

func (dao *TradeDao) GetAllTradesByPairAddress(bt, qt common.Address) ([]*types.Trade, error) {
	res, err := dao.GetTradesByPairAddress(bt, qt, 0)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	return res, nil
}

// GetByPairAddress fetches all the trades corresponding to a particular pair token address.
func (dao *TradeDao) GetTradesByPairAddress(bt, qt common.Address, n int) ([]*types.Trade, error) {
	var res []*types.Trade

	q := bson.M{"baseToken": bt.Hex(), "quoteToken": qt.Hex()}
	err := db.Get(dao.dbName, dao.collectionName, q, 0, n, &res)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	return res, nil
}

func (dao *TradeDao) GetSortedTradesByUserAddress(a common.Address, limit ...int) ([]*types.Trade, error) {
	if limit == nil {
		limit = []int{0}
	}

	var res []*types.Trade
	q := bson.M{"$or": []bson.M{{"maker": a.Hex()}, {"taker": a.Hex()}}}
	sort := []string{"-createdAt"}

	err := db.GetAndSort(dao.dbName, dao.collectionName, q, sort, 0, limit[0], &res)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	return res, nil
}

// GetByUserAddress fetches all the trades corresponding to a particular user address.
func (dao *TradeDao) GetByUserAddress(a common.Address) ([]*types.Trade, error) {
	var res []*types.Trade
	q := bson.M{"$or": []bson.M{{"maker": a.Hex()}, {"taker": a.Hex()}}}

	err := db.Get(dao.dbName, dao.collectionName, q, 0, 0, &res)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	return res, nil
}

func (dao *TradeDao) UpdateTradeStatus(h common.Hash, status string) error {
	query := bson.M{"hash": h.Hex()}
	update := bson.M{"$set": bson.M{"status": status}}

	err := db.Update(dao.dbName, dao.collectionName, query, update)
	if err != nil {
		logger.Error(err)
		return err
	}

	return nil
}

func (dao *TradeDao) UpdateTradeStatuses(status string, hashes ...common.Hash) ([]*types.Trade, error) {
	hexes := []string{}
	for _, h := range hashes {
		hexes = append(hexes, h.Hex())
	}

	query := bson.M{"hash": bson.M{"$in": hexes}}
	update := bson.M{
		"$set": bson.M{
			"updatedAt": time.Now(),
			"status":    status,
		},
	}

	err := db.UpdateAll(dao.dbName, dao.collectionName, query, update)
	if err != nil {
		logger.Error(err)
		return nil, nil
	}

	trades := []*types.Trade{}
	err = db.Get(dao.dbName, dao.collectionName, query, 0, 0, &trades)
	if err != nil {
		logger.Error(err)
		return nil, nil
	}

	return trades, nil
}

func (dao *TradeDao) UpdateTradeStatusesByOrderHashes(status string, hashes ...common.Hash) ([]*types.Trade, error) {
	hexes := []string{}
	for _, h := range hashes {
		hexes = append(hexes, h.Hex())
	}

	query := bson.M{"makerOrderHash": bson.M{"$in": hexes}}
	update := bson.M{
		"$set": bson.M{
			"updatedAt": time.Now(),
			"status":    status,
		},
	}

	err := db.UpdateAll(dao.dbName, dao.collectionName, query, update)
	if err != nil {
		logger.Error(err)
		return nil, nil
	}

	trades := []*types.Trade{}
	err = db.Get(dao.dbName, dao.collectionName, query, 0, 0, &trades)
	if err != nil {
		logger.Error(err)
		return nil, nil
	}

	return trades, nil
}

func (dao *TradeDao) UpdateTradeStatusesByHashes(status string, hashes ...common.Hash) ([]*types.Trade, error) {
	hexes := []string{}
	for _, h := range hashes {
		hexes = append(hexes, h.Hex())
	}

	query := bson.M{"hash": bson.M{"$in": hexes}}
	update := bson.M{
		"$set": bson.M{
			"updatedAt": time.Now(),
			"status":    status,
		},
	}

	err := db.UpdateAll(dao.dbName, dao.collectionName, query, update)
	if err != nil {
		logger.Error(err)
		return nil, nil
	}

	trades := []*types.Trade{}
	err = db.Get(dao.dbName, dao.collectionName, query, 0, 0, &trades)
	if err != nil {
		logger.Error(err)
		return nil, nil
	}

	return trades, nil
}

// Drop drops all the order documents in the current database
func (dao *TradeDao) Drop() {
	db.DropCollection(dao.dbName, dao.collectionName)
}
