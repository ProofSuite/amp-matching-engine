package daos

import (
	"errors"
	"time"

	"github.com/Proofsuite/amp-matching-engine/app"
	"github.com/Proofsuite/amp-matching-engine/types"
	"gopkg.in/mgo.v2/bson"
)

// BalanceDao contains:
// collectionName: MongoDB collection name
// dbName: name of mongodb to interact with
type BalanceDao struct {
	collectionName string
	dbName         string
}

// NewBalanceDao returns a new instance of AddressDao
func NewBalanceDao() *BalanceDao {
	return &BalanceDao{"balances", app.Config.DBName}
}

// Create function performs the DB insertion task for Balance collection
func (dao *BalanceDao) Create(balance *types.Balance) (err error) {

	balance.ID = bson.NewObjectId()
	balance.CreatedAt = time.Now()
	balance.UpdatedAt = time.Now()

	err = db.Create(dao.dbName, dao.collectionName, balance)
	return
}

// GetByAddress function fetches document from db collection based on user address.
// Returns Balance type struct
func (dao *BalanceDao) GetByAddress(addr string) (response *types.Balance, err error) {
	var res []*types.Balance
	q := bson.M{"address": addr}
	err = db.Get(dao.dbName, dao.collectionName, q, 0, 1, &res)
	if err != nil {
		return
	} else if len(res) > 0 {
		response = res[0]
	} else {
		err = errors.New("No Wallet for address found")
	}
	return
}

// UpdateAmount updates amount corresponding to a particular token for a given user
func (dao *BalanceDao) UpdateAmount(address string, token string, amount *types.TokenBalance) (err error) {
	q := bson.M{"address": address}
	updateQuery := bson.M{
		"$set": bson.M{
			"tokens." + token: amount,
		},
	}

	err = db.Update(dao.dbName, dao.collectionName, q, updateQuery)
	return
}
