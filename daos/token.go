package daos

import (
	"time"

	"github.com/Proofsuite/amp-matching-engine/app"
	"github.com/Proofsuite/amp-matching-engine/types"
	"labix.org/v2/mgo/bson"
)

type TokenDao struct {
	collectionName string
	dbName         string
}

func NewTokenDao() *TokenDao {
	return &TokenDao{"tokens", app.Config.DBName}
}

func (dao *TokenDao) Create(token *types.Token) (err error) {

	if err := token.Validate(); err != nil {
		return err
	}

	token.ID = bson.NewObjectId()
	token.CreatedAt = time.Now()
	token.UpdatedAt = time.Now()

	err = DB.Create(dao.dbName, dao.collectionName, token)
	return
}

func (dao *TokenDao) GetAll() (response []types.Token, err error) {
	err = DB.Get(dao.dbName, dao.collectionName, bson.M{}, 0, 0, &response)
	return
}

func (dao *TokenDao) GetByID(id bson.ObjectId) (response *types.Token, err error) {
	err = DB.GetByID(dao.dbName, dao.collectionName, id, &response)
	return
}
