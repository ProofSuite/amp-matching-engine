package daos

import (
	"github.com/Proofsuite/amp-matching-engine/app"
	"github.com/Proofsuite/amp-matching-engine/types"
	"github.com/ethereum/go-ethereum/common"
	"gopkg.in/mgo.v2/bson"
)

// TokenDao contains:
// collectionName: MongoDB collection name
// dbName: name of mongodb to interact with
type WalletDao struct {
	collectionName string
	dbName         string
}

func NewWalletDao() *WalletDao {
	return &WalletDao{"wallets", app.Config.DBName}
}

func (dao *WalletDao) Create(wallet *types.Wallet) error {
	err := wallet.Validate()
	if err != nil {
		logger.Error(err)
		return err
	}

	wallet.ID = bson.NewObjectId()
	err = db.Create(dao.dbName, dao.collectionName, wallet)
	if err != nil {
		logger.Error(err)
		return err
	}

	return nil
}

func (dao *WalletDao) GetAll() ([]types.Wallet, error) {
	var response []types.Wallet

	err := db.Get(dao.dbName, dao.collectionName, bson.M{}, 0, 0, &response)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	return response, nil
}

// GetByID function fetches details of a token based on its mongo id
func (dao *WalletDao) GetByID(id bson.ObjectId) (*types.Wallet, error) {
	var response *types.Wallet

	err := db.GetByID(dao.dbName, dao.collectionName, id, &response)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	return response, nil
}

// GetByAddress function fetches details of a token based on its contract address
func (dao *WalletDao) GetByAddress(a common.Address) (*types.Wallet, error) {
	q := bson.M{"address": a.Hex()}
	var resp []types.Wallet

	err := db.Get(dao.dbName, dao.collectionName, q, 0, 1, &resp)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	if len(resp) == 0 {
		logger.Info("No wallets found")
		return nil, nil
	}

	return &resp[0], nil
}

func (dao *WalletDao) GetDefaultAdminWallet() (*types.Wallet, error) {
	q := bson.M{"admin": true}
	var resp []types.Wallet

	err := db.Get(dao.dbName, dao.collectionName, q, 0, 1, &resp)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	if len(resp) == 0 {
		logger.Info("No default admin wallet")
		return nil, nil
	}

	return &resp[0], nil
}

func (dao *WalletDao) GetOperatorWallets() ([]*types.Wallet, error) {
	q := bson.M{"operator": true}
	res := []*types.Wallet{}

	err := db.Get(dao.dbName, dao.collectionName, q, 0, 1, &res)
	if err != nil || len(res) == 0 {
		logger.Error(err)
		return nil, err
	}

	return res, nil
}
