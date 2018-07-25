package daos

import (
	"time"

	"github.com/Proofsuite/amp-matching-engine/app"
	"github.com/Proofsuite/amp-matching-engine/types"
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
	// index := mgo.Index{
	// 	Key:    []string{"code"},
	// 	Unique: true,
	// }

	// err := db.C("currency").EnsureIndex(index)
	// if err != nil {
	// 	panic(err)
	// }
	return &OrderDao{"orders", app.Config.DBName}
}

// Create function performs the DB insertion task for Order collection
func (dao *OrderDao) Create(order *types.Order) (err error) {

	order.ID = bson.NewObjectId()
	order.Status = types.NEW
	order.CreatedAt = time.Now()
	order.UpdatedAt = time.Now()

	err = db.Create(dao.dbName, dao.collectionName, order)
	return
}

// Update function performs the DB updations task for Order collection
// corresponding to a particular order ID
func (dao *OrderDao) Update(id bson.ObjectId, order *types.Order) (response []types.Order, err error) {
	order.UpdatedAt = time.Now()
	err = db.Update(dao.dbName, dao.collectionName, bson.M{"_id": id}, order)
	return
}

// GetByID function fetches a single document from order collection based on mongoDB ID.
// Returns Order type struct
func (dao *OrderDao) GetByID(id bson.ObjectId) (response *types.Order, err error) {
	err = db.GetByID(dao.dbName, dao.collectionName, id, &response)
	return
}

// GetByHash function fetches a single document from order collection based on mongoDB ID.
// Returns Order type struct
func (dao *OrderDao) GetByHash(hash string) (response *types.Order, err error) {
	q := bson.M{"hash": hash}
	var resp []types.Order
	err = db.Get(dao.dbName, dao.collectionName, q, 0, 1, &resp)
	if err != nil || len(resp) == 0 {
		return
	}
	return &resp[0], nil
}

// GetByUserAddress function fetches list of orders from order collection based on user address.
// Returns array of Order type struct
func (dao *OrderDao) GetByUserAddress(addr string) (response []*types.Order, err error) {
	q := bson.M{"userAddress": addr}
	err = db.Get(dao.dbName, dao.collectionName, q, 0, 0, &response)
	return
}
