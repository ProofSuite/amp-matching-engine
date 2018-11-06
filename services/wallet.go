package services

import (
	"github.com/Proofsuite/amp-matching-engine/interfaces"
	"github.com/Proofsuite/amp-matching-engine/types"
	"github.com/ethereum/go-ethereum/common"
)

// WalletService struct with daos required, responsible for communicating with daos
type WalletService struct {
	WalletDao interfaces.WalletDao
}

func NewWalletService(walletDao interfaces.WalletDao) *WalletService {
	return &WalletService{walletDao}
}

func (s *WalletService) CreateAdminWallet(a common.Address) (*types.Wallet, error) {
	w := &types.Wallet{
		Address: a,
		Admin:   true,
	}

	err := s.WalletDao.Create(w)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	return w, nil
}

func (s *WalletService) GetDefaultAdminWallet() (*types.Wallet, error) {
	return s.WalletDao.GetDefaultAdminWallet()
}

func (s *WalletService) GetOperatorAddresses() ([]common.Address, error) {
	wallets, err := s.WalletDao.GetOperatorWallets()
	if err != nil {
		return nil, err
	}

	if wallets == nil {
		return []common.Address{}, nil
	}

	addresses := []common.Address{}
	for _, w := range wallets {
		addresses = append(addresses, w.Address)
	}

	return addresses, nil
}

func (s *WalletService) GetOperatorWallets() ([]*types.Wallet, error) {
	return s.WalletDao.GetOperatorWallets()
}

func (s *WalletService) GetAll() ([]types.Wallet, error) {
	return s.WalletDao.GetAll()
}

func (s *WalletService) GetByAddress(a common.Address) (*types.Wallet, error) {
	return s.WalletDao.GetByAddress(a)
}
