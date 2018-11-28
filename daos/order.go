package daos

import (
	"math/big"
	"time"

	"github.com/Proofsuite/amp-matching-engine/app"
	"github.com/Proofsuite/amp-matching-engine/types"
	"github.com/Proofsuite/amp-matching-engine/utils/math"
	"github.com/ethereum/go-ethereum/common"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// OrderDao contains:
// collectionName: MongoDB collection name
// dbName: name of mongodb to interact with
type OrderDao struct {
	collectionName string
	dbName         string
}

type OrderDaoOption = func(*OrderDao) error

func OrderDaoDBOption(dbName string) func(dao *OrderDao) error {
	return func(dao *OrderDao) error {
		dao.dbName = dbName
		return nil
	}
}

// NewOrderDao returns a new instance of OrderDao
func NewOrderDao(opts ...OrderDaoOption) *OrderDao {
	dao := &OrderDao{}
	dao.collectionName = "orders"
	dao.dbName = app.Config.DBName

	for _, op := range opts {
		err := op(dao)
		if err != nil {
			panic(err)
		}
	}

	index := mgo.Index{
		Key:    []string{"hash"},
		Unique: true,
	}

	i2 := mgo.Index{
		Key: []string{"status"},
	}

	i3 := mgo.Index{
		Key: []string{"baseToken"},
	}

	i4 := mgo.Index{
		Key: []string{"quoteToken"},
	}

	err := db.Session.DB(dao.dbName).C(dao.collectionName).EnsureIndex(index)
	if err != nil {
		panic(err)
	}

	err = db.Session.DB(dao.dbName).C(dao.collectionName).EnsureIndex(i2)
	if err != nil {
		panic(err)
	}

	err = db.Session.DB(dao.dbName).C(dao.collectionName).EnsureIndex(i3)
	if err != nil {
		panic(err)
	}

	err = db.Session.DB(dao.dbName).C(dao.collectionName).EnsureIndex(i4)
	if err != nil {
		panic(err)
	}

	return dao
}

// Create function performs the DB insertion task for Order collection
func (dao *OrderDao) Create(o *types.Order) error {
	o.ID = bson.NewObjectId()
	o.CreatedAt = time.Now()
	o.UpdatedAt = time.Now()

	if o.Status == "" {
		o.Status = "OPEN"
	}

	err := db.Create(dao.dbName, dao.collectionName, o)
	if err != nil {
		logger.Error(err)
		return err
	}

	return nil
}

func (dao *OrderDao) DeleteByHashes(hashes ...common.Hash) error {
	err := db.RemoveAll(dao.dbName, dao.collectionName, bson.M{"hash": bson.M{"$in": hashes}})
	if err != nil {
		logger.Error(err)
		return err
	}

	return nil
}

func (dao *OrderDao) Delete(orders ...*types.Order) error {
	hashes := []common.Hash{}
	for _, o := range orders {
		hashes = append(hashes, o.Hash)
	}

	err := db.RemoveAll(dao.dbName, dao.collectionName, bson.M{"hash": bson.M{"$in": hashes}})
	if err != nil {
		logger.Error(err)
		return err
	}

	return nil
}

// Update function performs the DB updations task for Order collection
// corresponding to a particular order ID
func (dao *OrderDao) Update(id bson.ObjectId, o *types.Order) error {
	o.UpdatedAt = time.Now()

	err := db.Update(dao.dbName, dao.collectionName, bson.M{"_id": id}, o)
	if err != nil {
		logger.Error(err)
		return err
	}

	return nil
}

func (dao *OrderDao) Upsert(id bson.ObjectId, o *types.Order) error {
	o.UpdatedAt = time.Now()

	err := db.Upsert(dao.dbName, dao.collectionName, bson.M{"_id": id}, o)
	if err != nil {
		logger.Error(err)
		return err
	}

	return nil
}

func (dao *OrderDao) UpsertByHash(h common.Hash, o *types.Order) error {
	err := db.Upsert(dao.dbName, dao.collectionName, bson.M{"hash": h.Hex()}, types.OrderBSONUpdate{o})
	if err != nil {
		logger.Error(err)
		return err
	}

	return nil
}

func (dao *OrderDao) UpdateAllByHash(h common.Hash, o *types.Order) error {
	o.UpdatedAt = time.Now()

	err := db.Update(dao.dbName, dao.collectionName, bson.M{"hash": h.Hex()}, o)
	if err != nil {
		logger.Error(err)
		return err
	}

	return nil
}

func (dao *OrderDao) FindAndModify(h common.Hash, o *types.Order) (*types.Order, error) {
	o.UpdatedAt = time.Now()
	query := bson.M{"hash": h.Hex()}
	updated := &types.Order{}
	change := mgo.Change{
		Update:    types.OrderBSONUpdate{o},
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

//UpdateByHash updates fields that are considered updateable for an order.
func (dao *OrderDao) UpdateByHash(h common.Hash, o *types.Order) error {
	o.UpdatedAt = time.Now()
	query := bson.M{"hash": h.Hex()}
	update := bson.M{"$set": bson.M{
		"pricepoint":   o.PricePoint.String(),
		"amount":       o.Amount.String(),
		"status":       o.Status,
		"filledAmount": o.FilledAmount.String(),
		"makeFee":      o.MakeFee.String(),
		"takeFee":      o.TakeFee.String(),
		"updatedAt":    o.UpdatedAt,
	}}

	err := db.Update(dao.dbName, dao.collectionName, query, update)
	if err != nil {
		logger.Error(err)
		return err
	}

	return nil
}

func (dao *OrderDao) UpdateOrderStatus(h common.Hash, status string) error {
	query := bson.M{"hash": h.Hex()}
	update := bson.M{"$set": bson.M{
		"status": status,
	}}

	err := db.Update(dao.dbName, dao.collectionName, query, update)
	if err != nil {
		logger.Error(err)
		return err
	}

	return nil
}

func (dao *OrderDao) UpdateOrderStatusesByHashes(status string, hashes ...common.Hash) ([]*types.Order, error) {
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

	orders := []*types.Order{}
	err = db.Get(dao.dbName, dao.collectionName, query, 0, 0, &orders)
	if err != nil {
		logger.Error(err)
		return nil, nil
	}

	return orders, nil
}

func (dao *OrderDao) UpdateOrderFilledAmount(hash common.Hash, value *big.Int) error {
	q := bson.M{"hash": hash.Hex()}
	res := []types.Order{}
	err := db.Get(dao.dbName, dao.collectionName, q, 0, 1, &res)
	if err != nil {
		logger.Error(err)
		return err
	}

	o := res[0]
	status := ""
	filledAmount := math.Add(o.FilledAmount, value)

	if math.IsEqualOrSmallerThan(filledAmount, big.NewInt(0)) {
		filledAmount = big.NewInt(0)
		status = "OPEN"
	} else if math.IsEqualOrGreaterThan(filledAmount, o.Amount) {
		filledAmount = o.Amount
		status = "FILLED"
	} else {
		status = "PARTIAL_FILLED"
	}

	update := bson.M{"$set": bson.M{
		"status":       status,
		"filledAmount": filledAmount.String(),
	}}

	err = db.Update(dao.dbName, dao.collectionName, q, update)
	if err != nil {
		logger.Error(err)
		return err
	}

	return nil
}

func (dao *OrderDao) UpdateOrderFilledAmounts(hashes []common.Hash, amount []*big.Int) ([]*types.Order, error) {
	hexes := []string{}
	orders := []*types.Order{}
	for i, _ := range hashes {
		hexes = append(hexes, hashes[i].Hex())
	}

	query := bson.M{"hash": bson.M{"$in": hexes}}
	err := db.Get(dao.dbName, dao.collectionName, query, 0, 0, &orders)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	updatedOrders := []*types.Order{}
	for i, o := range orders {
		status := ""
		filledAmount := math.Sub(o.FilledAmount, amount[i])

		if math.IsEqualOrSmallerThan(filledAmount, big.NewInt(0)) {
			filledAmount = big.NewInt(0)
			status = "OPEN"
		} else if math.IsEqualOrGreaterThan(filledAmount, o.Amount) {
			filledAmount = o.Amount
			status = "FILLED"
		} else {
			status = "PARTIAL_FILLED"
		}

		query := bson.M{"hash": o.Hash.Hex()}
		update := bson.M{"$set": bson.M{
			"status":       status,
			"filledAmount": filledAmount.String(),
		}}
		change := mgo.Change{
			Update:    update,
			Upsert:    true,
			Remove:    false,
			ReturnNew: true,
		}

		updated := &types.Order{}
		err := db.FindAndModify(dao.dbName, dao.collectionName, query, change, updated)
		if err != nil {
			logger.Error(err)
			return nil, err
		}

		updatedOrders = append(updatedOrders, updated)
	}

	return updatedOrders, nil
}

// GetByID function fetches a single document from order collection based on mongoDB ID.
// Returns Order type struct
func (dao *OrderDao) GetByID(id bson.ObjectId) (*types.Order, error) {
	var response *types.Order
	err := db.GetByID(dao.dbName, dao.collectionName, id, &response)
	return response, err
}

// GetByHash function fetches a single document from order collection based on mongoDB ID.
// Returns Order type struct
func (dao *OrderDao) GetByHash(hash common.Hash) (*types.Order, error) {
	q := bson.M{"hash": hash.Hex()}
	res := []types.Order{}

	err := db.Get(dao.dbName, dao.collectionName, q, 0, 1, &res)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	if len(res) == 0 {
		return nil, nil
	}

	return &res[0], nil
}

// GetByHashes
func (dao *OrderDao) GetByHashes(hashes []common.Hash) ([]*types.Order, error) {
	hexes := []string{}
	for _, h := range hashes {
		hexes = append(hexes, h.Hex())
	}

	q := bson.M{"hash": bson.M{"$in": hexes}}
	res := []*types.Order{}

	err := db.Get(dao.dbName, dao.collectionName, q, 0, 0, &res)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	return res, nil
}

// GetByUserAddress function fetches list of orders from order collection based on user address.
// Returns array of Order type struct
func (dao *OrderDao) GetByUserAddress(addr common.Address, limit ...int) ([]*types.Order, error) {
	if limit == nil {
		limit = []int{0}
	}

	var res []*types.Order
	q := bson.M{"userAddress": addr.Hex()}

	err := db.Get(dao.dbName, dao.collectionName, q, 0, limit[0], &res)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	if res == nil {
		return []*types.Order{}, nil
	}

	return res, nil
}

// GetCurrentByUserAddress function fetches list of open/partial orders from order collection based on user address.
// Returns array of Order type struct
func (dao *OrderDao) GetCurrentByUserAddress(addr common.Address, limit ...int) ([]*types.Order, error) {
	if limit == nil {
		limit = []int{0}
	}

	var res []*types.Order
	q := bson.M{
		"userAddress": addr.Hex(),
		"status": bson.M{"$in": []string{
			"OPEN",
			"PARTIAL_FILLED",
		},
		},
	}

	err := db.Get(dao.dbName, dao.collectionName, q, 0, limit[0], &res)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	return res, nil
}

// GetHistoryByUserAddress function fetches list of orders which are not in open/partial order status
// from order collection based on user address.
// Returns array of Order type struct
func (dao *OrderDao) GetHistoryByUserAddress(addr common.Address, limit ...int) ([]*types.Order, error) {
	if limit == nil {
		limit = []int{0}
	}

	var res []*types.Order
	q := bson.M{
		"userAddress": addr.Hex(),
		"status": bson.M{"$nin": []string{
			"OPEN",
			"PARTIAL_FILLED",
		},
		},
	}

	err := db.Get(dao.dbName, dao.collectionName, q, 0, limit[0], &res)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	return res, nil
}

func (dao *OrderDao) GetUserLockedBalance(account common.Address, token common.Address) (*big.Int, error) {
	var orders []*types.Order
	q := bson.M{
		"userAddress": account.Hex(),
		"status": bson.M{"$in": []string{
			"OPEN",
			"PARTIAL_FILLED",
		},
		},
		"sellToken": token.Hex(),
	}

	err := db.Get(dao.dbName, dao.collectionName, q, 0, 0, &orders)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	//TODO verify and refactor
	totalLockedBalance := big.NewInt(0)
	for _, o := range orders {
		lockedBalance := big.NewInt(0)
		if o.Side == "BUY" {
			lockedBalance = math.Sub(o.Amount, o.FilledAmount)
		} else if o.Side == "SELL" {
			lockedBalance = math.Mul(math.Sub(o.Amount, o.FilledAmount), o.PricePoint)
		}

		totalLockedBalance = math.Add(totalLockedBalance, lockedBalance)
	}

	return totalLockedBalance, nil
}

func (dao *OrderDao) GetRawOrderBook(p *types.Pair) ([]*types.Order, error) {
	var orders []*types.Order
	q := []bson.M{
		bson.M{
			"$match": bson.M{
				"status":     bson.M{"$in": []string{"OPEN", "PARTIAL_FILLED"}},
				"baseToken":  p.BaseTokenAddress.Hex(),
				"quoteToken": p.QuoteTokenAddress.Hex(),
			},
		},
		bson.M{
			"$addFields": bson.M{
				"priceDecimal": bson.M{"$toDecimal": "$pricepoint"},
			},
		},
		bson.M{
			"$sort": bson.M{
				"priceDecimal": 1,
			},
		},
	}

	err := db.Aggregate(dao.dbName, dao.collectionName, q, &orders)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	return orders, nil
}

func (dao *OrderDao) GetOrderBook(p *types.Pair) ([]map[string]string, []map[string]string, error) {
	bidsQuery := []bson.M{
		bson.M{
			"$match": bson.M{
				"status":     bson.M{"$in": []string{"OPEN", "PARTIAL_FILLED"}},
				"baseToken":  p.BaseTokenAddress.Hex(),
				"quoteToken": p.QuoteTokenAddress.Hex(),
				"side":       "BUY",
			},
		},
		bson.M{
			"$group": bson.M{
				"_id":        bson.M{"$toDecimal": "$pricepoint"},
				"pricepoint": bson.M{"$first": "$pricepoint"},
				"amount": bson.M{
					"$sum": bson.M{
						"$subtract": []bson.M{bson.M{"$toDecimal": "$amount"}, bson.M{"$toDecimal": "$filledAmount"}},
					},
				},
			},
		},
		bson.M{
			"$sort": bson.M{
				"_id": 1,
			},
		},
		bson.M{
			"$project": bson.M{
				"_id":        0,
				"pricepoint": 1,
				"amount":     bson.M{"$toString": "$amount"},
			},
		},
	}

	asksQuery := []bson.M{
		bson.M{
			"$match": bson.M{
				"status":     bson.M{"$in": []string{"OPEN", "PARTIAL_FILLED"}},
				"baseToken":  p.BaseTokenAddress.Hex(),
				"quoteToken": p.QuoteTokenAddress.Hex(),
				"side":       "SELL",
			},
		},
		bson.M{
			"$group": bson.M{
				"_id":        bson.M{"$toDecimal": "$pricepoint"},
				"pricepoint": bson.M{"$first": "$pricepoint"},
				"amount": bson.M{
					"$sum": bson.M{
						"$subtract": []bson.M{bson.M{"$toDecimal": "$amount"}, bson.M{"$toDecimal": "$filledAmount"}},
					},
				},
			},
		},
		bson.M{
			"$sort": bson.M{
				"_id": 1,
			},
		},
		bson.M{
			"$project": bson.M{
				"_id":        0,
				"pricepoint": 1,
				"amount":     bson.M{"$toString": "$amount"},
			},
		},
	}

	bids := []map[string]string{}
	asks := []map[string]string{}
	err := db.Aggregate(dao.dbName, dao.collectionName, bidsQuery, &bids)
	if err != nil {
		logger.Error(err)
		return nil, nil, err
	}

	err = db.Aggregate(dao.dbName, dao.collectionName, asksQuery, &asks)
	if err != nil {
		logger.Error(err)
		return nil, nil, err
	}

	return bids, asks, nil
}

func (dao *OrderDao) GetOrderBookPricePoint(p *types.Pair, pp *big.Int, side string) (*big.Int, error) {
	q := []bson.M{
		bson.M{
			"$match": bson.M{
				"status":     bson.M{"$in": []string{"OPEN", "PARTIAL_FILLED"}},
				"baseToken":  p.BaseTokenAddress.Hex(),
				"quoteToken": p.QuoteTokenAddress.Hex(),
				"pricepoint": pp.String(),
				"side":       side,
			},
		},
		bson.M{
			"$group": bson.M{
				"_id":        bson.M{"$toDecimal": "$pricepoint"},
				"pricepoint": bson.M{"$first": "$pricepoint"},
				"amount": bson.M{
					"$sum": bson.M{
						"$subtract": []bson.M{bson.M{"$toDecimal": "$amount"}, bson.M{"$toDecimal": "$filledAmount"}},
					},
				},
			},
		},
		bson.M{
			"$project": bson.M{
				"_id":        0,
				"pricepoint": 1,
				"amount":     bson.M{"$toString": "$amount"},
			},
		},
	}

	res := []map[string]string{}
	err := db.Aggregate(dao.dbName, dao.collectionName, q, &res)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	if len(res) == 0 {
		return nil, nil
	}

	return math.ToBigInt(res[0]["amount"]), nil
}

func (dao *OrderDao) GetMatchingBuyOrders(o *types.Order) ([]*types.Order, error) {
	var orders []*types.Order
	decimalPricepoint, _ := bson.ParseDecimal128(o.PricePoint.String())

	q := []bson.M{
		bson.M{
			"$match": bson.M{
				"status":     bson.M{"$in": []string{"OPEN", "PARTIAL_FILLED"}},
				"baseToken":  o.BaseToken.Hex(),
				"quoteToken": o.QuoteToken.Hex(),
				"side":       "BUY",
			},
		},
		bson.M{
			"$addFields": bson.M{
				"priceDecimal": bson.M{"$toDecimal": "$pricepoint"},
			},
		},
		bson.M{
			"$match": bson.M{
				"priceDecimal": bson.M{"$gte": decimalPricepoint},
			},
		},
		bson.M{
			"$sort": bson.M{"priceDecimal": -1, "createdAt": 1},
		},
	}

	err := db.Aggregate(dao.dbName, dao.collectionName, q, &orders)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	return orders, nil
}

func (dao *OrderDao) GetMatchingSellOrders(o *types.Order) ([]*types.Order, error) {
	var orders []*types.Order
	decimalPricepoint, _ := bson.ParseDecimal128(o.PricePoint.String())

	q := []bson.M{
		bson.M{
			"$match": bson.M{
				"status":     bson.M{"$in": []string{"OPEN", "PARTIAL_FILLED"}},
				"baseToken":  o.BaseToken.Hex(),
				"quoteToken": o.QuoteToken.Hex(),
				"side":       "SELL",
			},
		},
		bson.M{
			"$addFields": bson.M{
				"priceDecimal": bson.M{"$toDecimal": "$pricepoint"},
			},
		},
		bson.M{
			"$match": bson.M{
				"priceDecimal": bson.M{"$lte": decimalPricepoint},
			},
		},
		bson.M{
			"$sort": bson.M{"priceDecimal": 1, "createdAt": 1},
		},
	}

	err := db.Aggregate(dao.dbName, dao.collectionName, q, &orders)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	return orders, nil
}

// Drop drops all the order documents in the current database
func (dao *OrderDao) Drop() error {
	err := db.DropCollection(dao.dbName, dao.collectionName)
	if err != nil {
		logger.Error(err)
		return err
	}

	return nil
}
