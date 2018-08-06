package daos

import (
	"errors"
	"time"

	"github.com/Proofsuite/amp-matching-engine/app"
	"github.com/Proofsuite/amp-matching-engine/types"
	"gopkg.in/mgo.v2/bson"
)

// AddressDao contains:
// collectionName: MongoDB collection name
// dbName: name of mongodb to interact with
type AddressDao struct {
	collectionName string
	dbName         string
}

// NewAddressDao returns a new instance of AddressDao
func NewAddressDao() *AddressDao {
	return &AddressDao{"addresses", app.Config.DBName}
}

// Create function performs the DB insertion task for Address collection
func (dao *AddressDao) Create(address *types.UserAddress) (err error) {

	address.ID = bson.NewObjectId()
	address.CreatedAt = time.Now()
	address.UpdatedAt = time.Now()

	err = db.Create(dao.dbName, dao.collectionName, address)
	return
}

// GetAll function fetches all the documents in address collection.
// Returns array of UserAddress type structs
func (dao *AddressDao) GetAll() (response []types.UserAddress, err error) {
	err = db.Get(dao.dbName, dao.collectionName, bson.M{}, 0, 0, &response)
	return
}

// GetByID function fetches a single document from address collection based on mongoDB ID.
// Returns UserAddress type struct
func (dao *AddressDao) GetByID(id bson.ObjectId) (response *types.UserAddress, err error) {
	err = db.GetByID(dao.dbName, dao.collectionName, id, &response)
	return
}

// GetByAddress function fetches document from address collection based on user address.
// Returns UserAddress type struct
func (dao *AddressDao) GetByAddress(addr string) (response *types.UserAddress, err error) {
	var res []*types.UserAddress
	q := bson.M{"address": addr}
	err = db.Get(dao.dbName, dao.collectionName, q, 0, 1, &res)
	if err != nil {
		return
	} else if len(res) > 0 {
		response = res[0]
	} else {
		err = errors.New("Address not registered")
	}
	return
}

// GetNonce function fetches document from address collection based on user address.
// Returns nonce int64
func (dao *AddressDao) GetNonce(addr string) (nonce int64, err error) {
	res, err := dao.GetByAddress(addr)
	if err != nil {
		return 0, err
	}
	return res.Nonce, nil
}

// IncrNonce function increaments the nonce of the address
func (dao *AddressDao) IncrNonce(addr string) (err error) {
	return db.Update(dao.dbName, dao.collectionName, bson.M{"address": addr}, bson.M{"$inc": bson.M{"nonce": 1}})
}
