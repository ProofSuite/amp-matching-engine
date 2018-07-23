package daos

import (
	"time"

	"github.com/Proofsuite/amp-matching-engine/app"
	"github.com/Proofsuite/amp-matching-engine/types"
	"gopkg.in/mgo.v2/bson"
)

// TokenDao contains:
// collectionName: MongoDB collection name
// dbName: name of mongodb to interact with
type TokenDao struct {
	collectionName string
	dbName         string
}

// NewTokenDao returns a new instance of TokenDao.
func NewTokenDao() *TokenDao {
	return &TokenDao{"tokens", app.Config.DBName}
}

// Create function performs the DB insertion task for token collection
func (dao *TokenDao) Create(token *types.Token) (err error) {

	if err := token.Validate(); err != nil {
		return err
	}

	token.ID = bson.NewObjectId()
	token.CreatedAt = time.Now()
	token.UpdatedAt = time.Now()

	err = db.Create(dao.dbName, dao.collectionName, token)
	return
}

// GetAll function fetches all the tokens in the token collection of mongodb.
func (dao *TokenDao) GetAll() (response []types.Token, err error) {
	err = db.Get(dao.dbName, dao.collectionName, bson.M{}, 0, 0, &response)
	return
}

// GetByID function fetches details of a token based on its mongo id
func (dao *TokenDao) GetByID(id bson.ObjectId) (response *types.Token, err error) {
	err = db.GetByID(dao.dbName, dao.collectionName, id, &response)
	return
}

// GetByAddress function fetches details of a token based on its contract address
func (dao *TokenDao) GetByAddress(addr string) (response *types.Token, err error) {
	q := bson.M{"contractAddress": bson.RegEx{Pattern: addr, Options: "i"}}
	var resp []types.Token
	err = db.Get(dao.dbName, dao.collectionName, q, 0, 1, &resp)
	if err != nil || len(resp) == 0 {
		return
	}
	return &resp[0], nil
}
