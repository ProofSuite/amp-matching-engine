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

// OrderDao contains:
// collectionName: MongoDB collection name
// dbName: name of mongodb to interact with
type OrderDao struct {
	collectionName string
	dbName         string
}

// NewOrderDao returns a new instance of OrderDao
func NewOrderDao() *OrderDao {
	dbName := app.Config.DBName
	collection := "orders"
	index := mgo.Index{
		Key:    []string{"hash"},
		Unique: true,
	}

	err := db.Session.DB(dbName).C(collection).EnsureIndex(index)
	if err != nil {
		panic(err)
	}
	return &OrderDao{collection, dbName}
}

// Create function performs the DB insertion task for Order collection
func (dao *OrderDao) Create(order *types.Order) error {
	order.ID = bson.NewObjectId()
	order.Status = "NEW"
	order.CreatedAt = time.Now()
	order.UpdatedAt = time.Now()

	err := db.Create(dao.dbName, dao.collectionName, order)
	if err != nil {
		log.Print(err)
		return err
	}

	return nil
}

// q := bson.M{
// 	"address": owner.Hex(),
// }
// updateQuery := bson.M{
// 	"$set": bson.M{
// 		"tokenBalances." + token.Hex() + ".balance":       tokenBalance.Balance.String(),
// 		"tokenBalances." + token.Hex() + ".allowance":     tokenBalance.Allowance.String(),
// 		"tokenBalances." + token.Hex() + ".lockedBalance": tokenBalance.LockedBalance.String(),
// 	},
// }

// Update function performs the DB updations task for Order collection
// corresponding to a particular order ID
func (dao *OrderDao) Update(id bson.ObjectId, o *types.Order) error {
	o.UpdatedAt = time.Now()

	err := db.Update(dao.dbName, dao.collectionName, bson.M{"_id": id}, o)
	if err != nil {
		log.Print(err)
		return err
	}

	return nil

	// q := bson.M{
	// 	"_id": id,
	// }

	// updateQuery := bson.M{
	// 	"$set": bson.M{
	// 		"baseToken": o.BaseToken,
	// 		"hash":      o.Hash,
	// 	},
	// }

	// log.Println(id)

	// log.Print(response)

	// log.Print(id)

	// err := db.Update(dao.dbName, dao.collectionName, bson.M{"_id": id}, bson.M{
	// 	"baseToken":    o.BaseToken,
	// 	"quoteToken":   o.QuoteToken,
	// 	"buyAmount":    o.BuyAmount,
	// 	"sellAmount":   o.SellAmount,
	// 	"status":       o.Status,
	// 	"side":         o.Side,
	// 	"hash":         o.Hash,
	// 	"price":        o.Price,
	// 	"pricepoint":   o.PricePoint,
	// 	"amount":       o.Amount,
	// 	"filledAmount": o.FilledAmount,
	// 	"nonce":        o.Nonce,
	// 	"expires":      o.Expires,
	// 	"makeFee":      o.MakeFee,
	// 	"takeFee":      o.TakeFee,
	// 	"signature":    o.Signature,
	// 	"orderBook":    o.OrderBook,
	// })

}

func (dao *OrderDao) UpdateByHash(hash common.Hash, o *types.Order) error {
	o.UpdatedAt = time.Now()

	err := db.Update(dao.dbName, dao.collectionName, bson.M{"hash": hash.Hex()}, o)
	if err != nil {
		log.Print(err)
		return err
	}

	return nil
}

// GetByID function fetches a single document from order collection based on mongoDB ID.
// Returns Order type struct
func (dao *OrderDao) GetByID(id bson.ObjectId) (response *types.Order, err error) {
	err = db.GetByID(dao.dbName, dao.collectionName, id, &response)
	return
}

// GetByHash function fetches a single document from order collection based on mongoDB ID.
// Returns Order type struct
func (dao *OrderDao) GetByHash(hash common.Hash) (response *types.Order, err error) {
	q := bson.M{"hash": hash.Hex()}
	var resp []types.Order
	err = db.Get(dao.dbName, dao.collectionName, q, 0, 1, &resp)
	if err != nil || len(resp) == 0 {
		return
	}
	return &resp[0], nil
}

// GetByUserAddress function fetches list of orders from order collection based on user address.
// Returns array of Order type struct
func (dao *OrderDao) GetByUserAddress(addr common.Address) (response []*types.Order, err error) {
	q := bson.M{"userAddress": addr.Hex()}
	err = db.Get(dao.dbName, dao.collectionName, q, 0, 0, &response)
	return
}

// Drop drops all the order documents in the current database
func (dao *OrderDao) Drop() error {
	err := db.DropCollection(dao.dbName, dao.collectionName)
	if err != nil {
		log.Print(err)
		return err
	}

	return nil
}
