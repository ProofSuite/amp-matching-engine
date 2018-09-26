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

	err := db.Session.DB(dao.dbName).C(dao.collectionName).EnsureIndex(index)
	if err != nil {
		panic(err)
	}

	return dao
}

// Create function performs the DB insertion task for Order collection
func (dao *OrderDao) Create(order *types.Order) error {
	order.ID = bson.NewObjectId()
	order.CreatedAt = time.Now()
	order.UpdatedAt = time.Now()

	if order.Status == "" {
		order.Status = "OPEN"
	}

	err := db.Create(dao.dbName, dao.collectionName, order)
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

func (dao *OrderDao) UpdateAllByHash(hash common.Hash, o *types.Order) error {
	o.UpdatedAt = time.Now()

	err := db.Update(dao.dbName, dao.collectionName, bson.M{"hash": hash.Hex()}, o)
	if err != nil {
		logger.Error(err)
		return err
	}

	return nil
}

//UpdateByHash updates fields that are considered updateable for an order.
func (dao *OrderDao) UpdateByHash(hash common.Hash, o *types.Order) error {
	o.UpdatedAt = time.Now()
	query := bson.M{"hash": hash.Hex()}
	update := bson.M{"$set": bson.M{
		"buyAmount":    o.BuyAmount.String(),
		"sellAmount":   o.SellAmount.String(),
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

func (dao *OrderDao) UpdateOrderStatus(hash common.Hash, status string) error {
	query := bson.M{"hash": hash.Hex()}
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
		status = "PARTIALLY_FILLED"
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
func (dao *OrderDao) GetByUserAddress(addr common.Address) ([]*types.Order, error) {
	var res []*types.Order
	q := bson.M{"userAddress": addr.Hex()}

	err := db.Get(dao.dbName, dao.collectionName, q, 0, 0, &res)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	return res, nil
}

// GetCurrentByUserAddress function fetches list of open/partial orders from order collection based on user address.
// Returns array of Order type struct
func (dao *OrderDao) GetCurrentByUserAddress(addr common.Address) ([]*types.Order, error) {
	var res []*types.Order
	q := bson.M{
		"userAddress": addr.Hex(),
		"status": bson.M{"$in": []string{
			"OPEN",
			"PARTIALLY_FILLED",
		},
		},
	}

	err := db.Get(dao.dbName, dao.collectionName, q, 0, 0, &res)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	return res, nil
}

// GetHistoryByUserAddress function fetches list of orders which are not in open/partial order status
// from order collection based on user address.
// Returns array of Order type struct
func (dao *OrderDao) GetHistoryByUserAddress(addr common.Address) ([]*types.Order, error) {
	var res []*types.Order
	q := bson.M{
		"userAddress": addr.Hex(),
		"status": bson.M{"$nin": []string{
			"OPEN",
			"PARTIALLY_FILLED",
		},
		},
	}

	err := db.Get(dao.dbName, dao.collectionName, q, 0, 0, &res)
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
			"PARTIALLY_FILLED",
		},
		},
		"sellToken": token.Hex(),
	}

	err := db.Get(dao.dbName, dao.collectionName, q, 0, 0, &orders)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	totalLockedBalance := big.NewInt(0)
	for _, o := range orders {
		filledSellAmount := math.Div(math.Mul(o.FilledAmount, o.SellAmount), o.BuyAmount)
		lockedBalance := math.Sub(o.SellAmount, filledSellAmount)
		totalLockedBalance = math.Add(totalLockedBalance, lockedBalance)
	}

	return totalLockedBalance, nil
}

func (dao *OrderDao) GetRawOrderBook(p *types.Pair) ([]*types.Order, error) {
	var orders []*types.Order
	q := bson.M{
		"status":     bson.M{"$in": []string{"OPEN", "PARTIALLY_FILLED"}},
		"baseToken":  p.BaseTokenAddress.Hex(),
		"quoteToken": p.QuoteTokenAddress.Hex(),
	}

	sort := []string{"pricepoint"}
	err := db.GetAndSort(dao.dbName, dao.collectionName, q, sort, 0, 0, &orders)
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
				"status":     bson.M{"$in": []string{"OPEN", "PARTIALLY_FILLED"}},
				"baseToken":  p.BaseTokenAddress.Hex(),
				"quoteToken": p.QuoteTokenAddress.Hex(),
				"side":       "BUY",
			},
		},
		bson.M{
			"$group": bson.M{
				"_id": "$pricepoint",
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
				"pricepoint": "$_id",
				"amount":     bson.M{"$toString": "$amount"},
			},
		},
	}

	asksQuery := []bson.M{
		bson.M{
			"$match": bson.M{
				"status":     bson.M{"$in": []string{"OPEN", "PARTIALLY_FILLED"}},
				"baseToken":  p.BaseTokenAddress.Hex(),
				"quoteToken": p.QuoteTokenAddress.Hex(),
				"side":       "SELL",
			},
		},
		bson.M{
			"$group": bson.M{
				"_id": "$pricepoint",
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
				"pricepoint": "$_id",
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

func (dao *OrderDao) GetOrderBookPricePoint(p *types.Pair, pp *big.Int) (*big.Int, error) {
	q := []bson.M{
		bson.M{
			"$match": bson.M{
				"status":     bson.M{"$in": []string{"OPEN", "PARTIALLY_FILLED"}},
				"baseToken":  p.BaseTokenAddress.Hex(),
				"quoteToken": p.QuoteTokenAddress.Hex(),
				"pricepoint": pp.String(),
			},
		},
		bson.M{
			"$project": bson.M{
				"_id":        0,
				"pricepoint": "$pricepoint",
				"amount": bson.M{
					"$toString": bson.M{
						"$subtract": []bson.M{bson.M{"$toDecimal": "$amount"}, bson.M{"$toDecimal": "$filledAmount"}},
					},
				},
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

// Drop drops all the order documents in the current database
func (dao *OrderDao) Drop() error {
	err := db.DropCollection(dao.dbName, dao.collectionName)
	if err != nil {
		logger.Error(err)
		return err
	}

	return nil
}
