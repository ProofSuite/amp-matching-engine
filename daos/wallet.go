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
type WalletDao struct {
	collectionName string
	dbName         string
}

type WalletRecord struct {
	ID
	Address string
	PrivateKey string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewWalletDao() *WalletDao {
	return &WalletDao{"wallet", app.Config.DBName}
}

func (r *WalletRecord) Marshal(w *types.Wallet) {

}

func (r *WalletRecord) Unmarshal()

func (dao *WalletDao) Create(wallet *types.Wallet) (err error) {
	if err := wallet.Validate(); err != nil {
		return err
	}

	wallet.ID = bson.NewObjectId()

	r := &WalletRecord{
		ID: wallet.ID,
		Address: wallet.Address.Hex(),
		PrivateKey: wallet.PrivateKey.Hex(),
		CreatedAt: time.Now(),
		UpdateAt: time.Now()
	}

	err = db.Create(dao.dbName, dao.collectionName, wallet)
	return
}

func (dao *WalletDao) GetAll() (response []types.Wallet, err error) {
	err = db.Get(dao.dbName, dao.collectionName, bson.M{}, 0, 0, &response)
	return
}

// GetByID function fetches details of a token based on its mongo id
func (dao *WalletDao) GetByID(id bson.ObjectId) (response *types.Wallet, err error) {
	err = db.GetByID(dao.dbName, dao.collectionName, id, &response)
	return
}

// GetByAddress function fetches details of a token based on its contract address
func (dao *WalletDao) GetByAddress(addr string) (response *types.Wallet, err error) {
	q := bson.M{"address": bson.RegEx{Pattern: addr, Options: "i"}}
	var resp []types.Wallet
	err = db.Get(dao.dbName, dao.collectionName, q, 0, 1, &resp)
	if err != nil || len(resp) == 0 {
		return
	}
	return &resp[0], nil
}

func (dao *WalletDao) GetDefaultAdminWallet() (err error) {
	q := bson.M{"admin": true}
	var resp []types.Wallet
	err = db.Get(dao.dbName, dao.collectionName, q, 0, 1, &resp)
	if err != nil || len(resp) == 0 {
		return
	}

	return &resp[0], nil
}
