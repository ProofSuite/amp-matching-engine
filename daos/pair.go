package daos

import (
	"errors"
	"time"

	"github.com/Proofsuite/amp-matching-engine/app"
	"github.com/Proofsuite/amp-matching-engine/types"
	"github.com/ethereum/go-ethereum/common"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// PairDao contains:
// collectionName: MongoDB collection name
// dbName: name of mongodb to interact with
type PairDao struct {
	collectionName string
	dbName         string
}

type PairDaoOption = func(*PairDao) error

func PairDaoDBOption(dbName string) func(dao *PairDao) error {
	return func(dao *PairDao) error {
		dao.dbName = dbName
		return nil
	}
}

// NewPairDao returns a new instance of AddressDao
func NewPairDao(options ...PairDaoOption) *PairDao {
	dao := &PairDao{}
	dao.collectionName = "pairs"
	dao.dbName = app.Config.DBName

	for _, op := range options {
		err := op(dao)
		if err != nil {
			panic(err)
		}
	}

	index := mgo.Index{
		Key:    []string{"baseTokenAddress", "quoteTokenAddress"},
		Unique: true,
	}

	err := db.Session.DB(dao.dbName).C(dao.collectionName).EnsureIndex(index)
	if err != nil {
		panic(err)
	}

	return dao
}

// Create function performs the DB insertion task for pair collection
func (dao *PairDao) Create(pair *types.Pair) error {
	pair.ID = bson.NewObjectId()
	pair.CreatedAt = time.Now()
	pair.UpdatedAt = time.Now()

	err := db.Create(dao.dbName, dao.collectionName, pair)
	return err
}

// GetAll function fetches all the pairs in the pair collection of mongodb.
func (dao *PairDao) GetAll() ([]types.Pair, error) {
	var response []types.Pair
	err := db.Get(dao.dbName, dao.collectionName, bson.M{}, 0, 0, &response)
	return response, err
}

// GetByID function fetches details of a pair using pair's mongo ID.
func (dao *PairDao) GetByID(id bson.ObjectId) (*types.Pair, error) {
	var response *types.Pair
	err := db.GetByID(dao.dbName, dao.collectionName, id, &response)
	return response, err
}

// GetByName function fetches details of a pair using pair's name.
// It makes CASE INSENSITIVE search query one pair's name
func (dao *PairDao) GetByName(name string) (*types.Pair, error) {
	var res []*types.Pair
	q := bson.M{"name": bson.RegEx{
		Pattern: name,
		Options: "i",
	}}

	err := db.Get(dao.dbName, dao.collectionName, q, 0, 1, &res)
	if err != nil {
		return nil, err
	}

	if len(res) == 0 {
		return nil, errors.New("NO_PAIR_FOUND")
	}

	return res[0], nil
}

func (dao *PairDao) GetByTokenSymbols(baseTokenSymbol, quoteTokenSymbol string) (*types.Pair, error) {
	var res []*types.Pair

	q := bson.M{
		"baseTokenSymbol":  baseTokenSymbol,
		"quoteTokenSymbol": quoteTokenSymbol,
	}

	err := db.Get(dao.dbName, dao.collectionName, q, 0, 1, &res)
	if err != nil {
		return nil, err
	}

	if len(res) == 0 {
		return nil, errors.New("No pair found")
	}

	return res[0], nil
}

// GetByTokenAddress function fetches pair based on
// CONTRACT ADDRESS of base token and quote token
func (dao *PairDao) GetByTokenAddress(baseToken, quoteToken common.Address) (*types.Pair, error) {
	var res []*types.Pair

	q := bson.M{
		"baseTokenAddress":  baseToken.Hex(),
		"quoteTokenAddress": quoteToken.Hex(),
	}

	err := db.Get(dao.dbName, dao.collectionName, q, 0, 1, &res)
	if err != nil {
		return nil, err
	}

	if len(res) == 0 {
		return nil, errors.New("NO_PAIR_FOUND")
	}

	return res[0], nil
}

// GetByBuySellTokenAddress function fetches pair based on
// CONTRACT ADDRESS of buy token and sell token
func (dao *PairDao) GetByBuySellTokenAddress(buyToken, sellToken common.Address) (*types.Pair, error) {
	var res []*types.Pair
	q := bson.M{
		"$or": []bson.M{
			bson.M{
				"baseTokenAddress":  buyToken.Hex(),
				"quoteTokenAddress": sellToken.Hex(),
			},
			bson.M{
				"baseTokenAddress":  sellToken.Hex(),
				"quoteTokenAddress": buyToken.Hex(),
			},
		},
	}

	err := db.Get(dao.dbName, dao.collectionName, q, 0, 1, &res)
	if err != nil {
		return nil, err
	}

	if len(res) == 0 {
		return nil, errors.New("NO_PAIR_FOUND")
	}

	return res[0], nil
}
