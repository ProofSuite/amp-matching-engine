package services

import (
	"github.com/Proofsuite/amp-matching-engine/daos"
	"github.com/Proofsuite/amp-matching-engine/types"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
)

// WalletService struct with daos required, responsible for communicating with daos
type TxService struct {
	WalletDao daos.WalletDaoInterface
	Wallet    *types.Wallet
}

func NewTxService(dao daos.WalletDaoInterface, w *types.Wallet) *TxService {
	return &TxService{dao, w}
}

func (s *TxService) GetTxCallOptions() *bind.CallOpts {
	return &bind.CallOpts{Pending: true}
}

func (s *TxService) GetTxDefaultSendOptions() (*bind.TransactOpts, error) {
	wallet, err := s.WalletDao.GetDefaultAdminWallet()
	if err != nil {
		return nil, err
	}

	return bind.NewKeyedTransactor(wallet.PrivateKey), nil
}

func (s *TxService) GetTxSendOptions() (*bind.TransactOpts, error) {
	return bind.NewKeyedTransactor(s.Wallet.PrivateKey), nil
}

func (s *TxService) SetTxSender(w *types.Wallet) {
	s.Wallet = w
}

func (s *TxService) GetCustomTxSendOptions(w *types.Wallet) *bind.TransactOpts {
	return bind.NewKeyedTransactor(w.PrivateKey)
}
