package daos

import (
	"errors"
	"time"

	"github.com/Proofsuite/amp-matching-engine/app"
	"github.com/Proofsuite/amp-matching-engine/types"
	"labix.org/v2/mgo/bson"
)

type PairDao struct {
	collectionName string
	dbName         string
}

func NewPairDao() *PairDao {
	return &PairDao{"pairs", app.Config.DBName}
}

func (dao *PairDao) Create(pair *types.Pair) (err error) {

	pair.ID = bson.NewObjectId()
	pair.CreatedAt = time.Now()
	pair.UpdatedAt = time.Now()

	err = DB.Create(dao.dbName, dao.collectionName, pair)
	return
}

func (dao *PairDao) GetAll() (response []types.Pair, err error) {
	err = DB.Get(dao.dbName, dao.collectionName, bson.M{}, 0, 0, &response)
	return
}

func (dao *PairDao) GetByID(id bson.ObjectId) (response *types.Pair, err error) {
	err = DB.GetByID(dao.dbName, dao.collectionName, id, &response)
	return
}
func (dao *PairDao) GetByName(name string) (response *types.Pair, err error) {
	var res []*types.Pair
	q := bson.M{"name": bson.RegEx{name, "i"}}
	err = DB.Get(dao.dbName, dao.collectionName, q, 0, 1, &res)
	if err != nil {
		return
	} else if len(res) > 0 {
		response = res[0]
	} else {
		err = errors.New("No Pair found")
	}
	return
}
