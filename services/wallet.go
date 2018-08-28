package services

import (
	"github.com/Proofsuite/amp-matching-engine/daos"
	"github.com/Proofsuite/amp-matching-engine/types"
	"github.com/ethereum/go-ethereum/common"
)

// WalletService struct with daos required, responsible for communicating with daos
type WalletService struct {
	WalletDao daos.WalletDaoInterface
}

type WalletServiceInterface interface {
	CreateAdminWallet(a common.Address) (*types.Wallet, error)
	GetDefaultAdminWallet() (*types.Wallet, error)
	GetAll() ([]types.Wallet, error)
	GetByAddress(a common.Address) *types.Wallet
}

func NewWalletService(walletDao daos.WalletDaoInterface) *WalletService {
	return &WalletService{walletDao}
}

func (s *WalletService) CreateAdminWallet(a common.Address) (*types.Wallet, error) {
	w := &types.Wallet{
		Address: a,
		Admin:   true,
	}

	err := s.WalletDao.Create(w)
	if err != nil {
		return nil, err
	}

	return w, nil
}

func (s *WalletService) GetDefaultAdminWallet() (*types.Wallet, error) {
	return s.WalletDao.GetDefaultAdminWallet()
}

func (s *WalletService) GetAll() ([]types.Wallet, error) {
	return s.WalletDao.GetAll()
}

func (s *WalletService) GetbyAddress(a common.Address) (*types.Wallet, error) {
	return s.WalletDao.GetByAddress(a)
}
