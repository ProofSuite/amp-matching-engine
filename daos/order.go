package daos

import (
	"time"

	"github.com/Proofsuite/amp-matching-engine/app"
	"github.com/Proofsuite/amp-matching-engine/types"
	"labix.org/v2/mgo/bson"
)

type OrderDao struct {
	collectionName string
	dbName         string
}

func NewOrderDao() *OrderDao {
	return &OrderDao{"orders", app.Config.DBName}
}

func (dao *OrderDao) Create(order *types.Order) (err error) {

	order.ID = bson.NewObjectId()
	order.Status = types.NEW
	order.CreatedAt = time.Now()
	order.UpdatedAt = time.Now()

	err = DB.Create(dao.dbName, dao.collectionName, order)
	return
}

func (dao *OrderDao) GetAll() (response []types.Order, err error) {
	err = DB.Get(dao.dbName, dao.collectionName, bson.M{}, 0, 0, &response)
	return
}
func (dao *OrderDao) Update(id bson.ObjectId, order *types.Order) (response []types.Order, err error) {
	order.UpdatedAt = time.Now()
	err = DB.Update(dao.dbName, dao.collectionName, bson.M{"_id": id}, order)
	return
}
func (dao *OrderDao) GetByID(id bson.ObjectId) (response *types.Order, err error) {
	err = DB.GetByID(dao.dbName, dao.collectionName, id, &response)
	return
}
func (dao *OrderDao) GetByUserAddress(addr string) (response []*types.Order, err error) {
	q := bson.M{"userAddress": addr}
	err = DB.Get(dao.dbName, dao.collectionName, q, 0, 0, &response)
	return
}
