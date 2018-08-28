package daos

import (
	"errors"
	"time"

	"github.com/Proofsuite/amp-matching-engine/app"
	"github.com/Proofsuite/amp-matching-engine/types"
	"github.com/ethereum/go-ethereum/common"
	"gopkg.in/mgo.v2/bson"
)

// PairDao contains:
// collectionName: MongoDB collection name
// dbName: name of mongodb to interact with
type PairDao struct {
	collectionName string
	dbName         string
}

type PairDaoInterface interface {
	Create(o *types.Pair) error
	GetAll() ([]types.Pair, error)
	GetByID(id bson.ObjectId) (*types.Pair, error)
	GetByName(name string) (*types.Pair, error)
	GetByTokenSymbols(baseTokenSymbol, quoteTokenSymbol string) (*types.Pair, error)
	GetByTokenAddress(baseToken, quoteToken common.Address) (*types.Pair, error)
	GetByBuySellTokenAddress(buyToken, sellToken common.Address) (*types.Pair, error)
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
