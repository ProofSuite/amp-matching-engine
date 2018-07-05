package daos

import (
	"errors"
	"time"

	"github.com/Proofsuite/amp-matching-engine/app"
	"github.com/Proofsuite/amp-matching-engine/types"
	"gopkg.in/mgo.v2/bson"
)

type BalanceDao struct {
	collectionName string
	dbName         string
}

func NewBalanceDao() *BalanceDao {
	return &BalanceDao{"balances", app.Config.DBName}
}

func (dao *BalanceDao) Create(balance *types.Balance) (err error) {

	balance.ID = bson.NewObjectId()
	balance.CreatedAt = time.Now()
	balance.UpdatedAt = time.Now()

	err = DB.Create(dao.dbName, dao.collectionName, balance)
	return
}

func (dao *BalanceDao) GetAll() (response []types.Balance, err error) {
	err = DB.Get(dao.dbName, dao.collectionName, bson.M{}, 0, 0, &response)
	return
}

func (dao *BalanceDao) GetByID(id bson.ObjectId) (response *types.Balance, err error) {
	err = DB.GetByID(dao.dbName, dao.collectionName, id, &response)
	return
}
func (dao *BalanceDao) GetByAddress(addr string) (response *types.Balance, err error) {
	var res []*types.Balance
	q := bson.M{"address": addr}
	err = DB.Get(dao.dbName, dao.collectionName, q, 0, 1, &res)
	if err != nil {
		return
	} else if len(res) > 0 {
		response = res[0]
	} else {
		err = errors.New("No Wallet for address found")
	}
	return
}

func (dao *BalanceDao) UpdateAmount(address string, token string, amount *types.TokenBalance) (err error) {
	q := bson.M{"address": address}
	updateQuery := bson.M{
		"$set": bson.M{
			"tokens." + token: amount,
		},
	}

	err = DB.Update(dao.dbName, dao.collectionName, q, updateQuery)
	return
}

// func (dao *BalanceDao) UnlockFunds(address string, token string, amount *types.TokenBalance) (err error) {
// 	q := bson.M{"address": address}
// 	updateQuery := bson.M{
// 		"$set": bson.M{
// 			"tokens." + token: amount,
// 		},
// 	}

// 	err = DB.Update(dao.dbName, dao.collectionName, q, updateQuery)
// 	return
// }
