package daos

import (
	"errors"
	"time"

	"github.com/Proofsuite/amp-matching-engine/app"
	"github.com/Proofsuite/amp-matching-engine/types"
	"gopkg.in/mgo.v2/bson"
)

// PairDao contains:
// collectionName: MongoDB collection name
// dbName: name of mongodb to interact with
type PairDao struct {
	collectionName string
	dbName         string
}

// NewPairDao returns a new instance of AddressDao
func NewPairDao() *PairDao {
	return &PairDao{"pairs", app.Config.DBName}
}

// Create function performs the DB insertion task for pair collection
func (dao *PairDao) Create(pair *types.Pair) (err error) {

	pair.ID = bson.NewObjectId()
	pair.CreatedAt = time.Now()
	pair.UpdatedAt = time.Now()

	err = db.Create(dao.dbName, dao.collectionName, pair)
	return
}

// GetAll function fetches all the pairs in the pair collection of mongodb.
func (dao *PairDao) GetAll() (response []types.Pair, err error) {
	err = db.Get(dao.dbName, dao.collectionName, bson.M{}, 0, 0, &response)
	return
}

// GetByID function fetches details of a pair using pair's mongo ID.
func (dao *PairDao) GetByID(id bson.ObjectId) (response *types.Pair, err error) {
	err = db.GetByID(dao.dbName, dao.collectionName, id, &response)
	return
}

// GetByName function fetches details of a pair using pair's name.
// It makes CASE INSENSITIVE search query one pair's name
func (dao *PairDao) GetByName(name string) (response *types.Pair, err error) {
	var res []*types.Pair
	q := bson.M{"name": bson.RegEx{
		Pattern: name,
		Options: "i",
	}}
	err = db.Get(dao.dbName, dao.collectionName, q, 0, 1, &res)
	if err != nil {
		return
	} else if len(res) > 0 {
		response = res[0]
	} else {
		err = errors.New("No Pair found")
	}
	return
}

// GetByTokenAddress function fetches pair based on
// CONTRACT ADDRESS of base token and quote token
func (dao *PairDao) GetByTokenAddress(baseToken, quoteToken string) (response *types.Pair, err error) {
	var res []*types.Pair
	q := bson.M{"baseTokenAddress": bson.RegEx{
		Pattern: baseToken,
		Options: "i",
	}, "quoteTokenAddress": bson.RegEx{
		Pattern: quoteToken,
		Options: "i",
	}}
	err = db.Get(dao.dbName, dao.collectionName, q, 0, 1, &res)
	if err != nil {
		return
	} else if len(res) > 0 {
		response = res[0]
	} else {
		err = errors.New("No Pair found")
	}
	return
}

// GetByBuySellTokenAddress function fetches pair based on
// CONTRACT ADDRESS of buy token and sell token
func (dao *PairDao) GetByBuySellTokenAddress(buyToken, sellToken string) (response *types.Pair, err error) {
	var res []*types.Pair
	q := bson.M{
		"$or": []bson.M{
			bson.M{"baseTokenAddress": bson.RegEx{
				Pattern: buyToken,
				Options: "i",
			}, "quoteTokenAddress": bson.RegEx{
				Pattern: sellToken,
				Options: "i",
			},
			},
			bson.M{"baseTokenAddress": bson.RegEx{
				Pattern: sellToken,
				Options: "i",
			}, "quoteTokenAddress": bson.RegEx{
				Pattern: buyToken,
				Options: "i",
			},
			},
		},
	}
	err = db.Get(dao.dbName, dao.collectionName, q, 0, 1, &res)
	if err != nil {
		return
	} else if len(res) > 0 {
		response = res[0]
	} else {
		err = errors.New("No Pair found")
	}
	return
}
