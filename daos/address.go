package daos

import (
	"errors"
	"time"

	"github.com/Proofsuite/matching-engine/app"
	"github.com/Proofsuite/matching-engine/types"
	"labix.org/v2/mgo/bson"
)

type AddressDao struct {
	collectionName string
	dbName         string
}

func NewAddressDao() *AddressDao {
	return &AddressDao{"addresses", app.Config.DBName}
}

func (dao *AddressDao) Create(address *types.UserAddress) (err error) {

	address.ID = bson.NewObjectId()
	address.CreatedAt = time.Now()
	address.UpdatedAt = time.Now()

	err = DB.Create(dao.dbName, dao.collectionName, address)
	return
}

func (dao *AddressDao) GetAll() (response []types.UserAddress, err error) {
	err = DB.Get(dao.dbName, dao.collectionName, bson.M{}, 0, 0, &response)
	return
}

func (dao *AddressDao) GetByID(id bson.ObjectId) (response *types.UserAddress, err error) {
	err = DB.GetByID(dao.dbName, dao.collectionName, id, &response)
	return
}
func (dao *AddressDao) GetByAddress(addr string) (response *types.UserAddress, err error) {
	var res []*types.UserAddress
	q := bson.M{"address": addr}
	err = DB.Get(dao.dbName, dao.collectionName, q, 0, 1, &res)
	if err != nil {
		return
	} else if len(res) > 0 {
		response = res[0]
	} else {
		err = errors.New("Address not registered")
	}
	return
}
