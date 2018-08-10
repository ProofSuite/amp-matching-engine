package services

import (
	"fmt"

	"gopkg.in/mgo.v2/bson"

	"github.com/Proofsuite/amp-matching-engine/daos"
	"github.com/Proofsuite/amp-matching-engine/types"
)

// WalletService struct with daos required, responsible for communicating with daos
type WalletService struct {
	WalletDao *daos.WalletDao
	BalanceDao *daos.BalanceDao
}

func NewWalletService(WalletDao *daos.AddressDao, balanceDao *daos.BalanceDao) *WalletService {
	return &WalletService{WalletDao, BalanceDao}
}

func (s *WalletService) CreateAdminWallet(a common.Address) error {
	w = *types.Wallet{
		address: a,
		admin: true
	}

	err := s.AddressDao.Create(w)
	if err != nil {
		return err
	}
}

func (s *WalletService) GetDefaultAdminWallet() error {
	return s.WalletDao.GetDefaultAdminWallet()
}

func (s *WalletService) GetAll() error {
	return s.WalletDao.GetAll()
}

func (s *WalletService) GetbyAddress(a string) (*types.UserAddress, error) {
	return s.WalletDao.GetByAddress(a)
}

// NewWalletService returns a new instance of addressService
func NewWalletService(AddressDao *daos.AddressDao, balanceDao *daos.BalanceDao, tokenDao *daos.TokenDao) *WalletService {
	return &WalletService{AddressDao, balanceDao, tokenDao}
}