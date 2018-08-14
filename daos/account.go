package daos

import (
	"math/big"
	"time"

	"github.com/Proofsuite/amp-matching-engine/app"
	"github.com/Proofsuite/amp-matching-engine/types"
	"github.com/ethereum/go-ethereum/common"
	"gopkg.in/mgo.v2/bson"
)

// BalanceDao contains:
// collectionName: MongoDB collection name
// dbName: name of mongodb to interact with
type AccountDao struct {
	collectionName string
	dbName         string
}

// NewBalanceDao returns a new instance of AddressDao
func NewAccountDao() *AccountDao {
	return &AccountDao{"accounts", app.Config.DBName}
}

// Create function performs the DB insertion task for Balance collection
func (dao *AccountDao) Create(account *types.Account) (err error) {
	account.ID = bson.NewObjectId()
	account.CreatedAt = time.Now()
	account.UpdatedAt = time.Now()

	err = db.Create(dao.dbName, dao.collectionName, account)
	return
}

func (dao *AccountDao) GetByAddress(owner common.Address) (response *types.Account, err error) {
	q := bson.M{"address": owner.Hex()}
	err = db.Get(dao.dbName, dao.collectionName, q, 0, 1, &response)
	return
}

func (dao *AccountDao) GetAllTokenBalances(owner common.Address) (map[common.Address]*types.TokenBalance, error) {
	q := bson.M{"address": owner.Hex()}
	response := []types.Account{}
	err := db.Get(dao.dbName, dao.collectionName, q, 0, 1, &response)
	if err != nil {
		return nil, err
	}

	return response[0].TokenBalances, nil
}

func (dao *AccountDao) GetTokenBalance(owner common.Address, token common.Address) (*types.TokenBalance, error) {
	q := []bson.M{
		bson.M{
			"$match": bson.M{
				"address": owner.Hex(),
			},
		},
		bson.M{
			"$project": bson.M{
				"tokenBalances": bson.M{
					"$objectToArray": "$tokenBalances",
				},
				"_id": 0,
			},
		},
		bson.M{
			"$addFields": bson.M{
				"tokenBalances": bson.M{
					"$filter": bson.M{
						"input": "$tokenBalances",
						"as":    "kv",
						"cond": bson.M{
							"$eq": []interface{}{"$$kv.k", token.Hex()},
						},
					},
				},
			},
		},
		bson.M{
			"$addFields": bson.M{
				"tokenBalances": bson.M{
					"$arrayToObject": "$tokenBalances",
				},
			},
		},
	}

	res, err := db.Aggregate(dao.dbName, dao.collectionName, q)
	if err != nil {
		return nil, err
	}

	a := &types.Account{}
	bytes, _ := bson.Marshal(res[0])
	bson.Unmarshal(bytes, &a)

	return a.TokenBalances[token], nil
}

func (dao *AccountDao) UpdateBalance(owner common.Address, token common.Address, balance *big.Int) (err error) {
	q := bson.M{
		"address": owner.Hex(),
	}
	updateQuery := bson.M{
		"$set": bson.M{"tokenBalances." + token.Hex() + ".balance": balance.String()},
	}

	err = db.Update(dao.dbName, dao.collectionName, q, updateQuery)
	return
}

func (dao *AccountDao) UpdateAllowance(owner common.Address, token common.Address, allowance *big.Int) (err error) {
	q := bson.M{
		"address": owner.Hex(),
	}
	updateQuery := bson.M{
		"$set": bson.M{"tokenBalances." + token.Hex() + ".allowance": allowance.String()},
	}

	err = db.Update(dao.dbName, dao.collectionName, q, updateQuery)
	return
}

// func (dao *AccountDao) UpdateAllowance(owner common.Address, token common.Address, allowance *big.Int) (err error) {
// 	q := bson.M{
// 		"address": bson.RegEx{
// 			Pattern: owner.Hex(),
// 			Options: "i",
// 		},
// 	}
// 	updateQuery := bson.M{
// 		"$set": bson.M{
// 			"tokenBalances." + token.Hex() + ".allowance": allowance.String(),
// 		},
// 	}

// 	err = db.Update(dao.dbName, dao.collectionName, q, updateQuery)
// 	return
// }

// func (dao *AccountDao) UpdateLockedBalance(owner common.Address, token common.Address, locked *big.Int) (err error) {
// 	q := bson.M{
// 		"address": bson.RegEx{
// 			Pattern: owner.Hex(),
// 			Options: "i",
// 		},
// 	}
// 	updateQuery := bson.M{
// 		"$set": bson.M{
// 			"tokenBalances." + token.Hex() + ".lockedBalance": locked.String(),
// 		},
// 	}

// 	err = db.Update(dao.dbName, dao.collectionName, q, updateQuery)
// 	return

// }
