package services

import (
	"math/big"

	"github.com/Proofsuite/amp-matching-engine/interfaces"
	"github.com/Proofsuite/amp-matching-engine/types"
	"github.com/ethereum/go-ethereum/common"
	"gopkg.in/mgo.v2/bson"
)

type AccountService struct {
	AccountDao interfaces.AccountDao
	TokenDao   interfaces.TokenDao
}

// NewAddressService returns a new instance of accountService
func NewAccountService(
	AccountDao interfaces.AccountDao,
	TokenDao interfaces.TokenDao,
) *AccountService {
	return &AccountService{AccountDao, TokenDao}
}

func (s *AccountService) Create(a *types.Account) error {
	addr := a.Address

	acc, err := s.AccountDao.GetByAddress(addr)
	if err != nil {
		logger.Error(err)
		return err
	}

	if acc != nil {
		return ErrAccountExists
	}

	tokens, err := s.TokenDao.GetAll()
	if err != nil {
		logger.Error(err)
		return err
	}

	a.IsBlocked = false
	a.TokenBalances = make(map[common.Address]*types.TokenBalance)

	// currently by default, the tokens balances are set to 0
	for _, token := range tokens {
		a.TokenBalances[token.ContractAddress] = &types.TokenBalance{
			Address:        token.ContractAddress,
			Symbol:         token.Symbol,
			Balance:        big.NewInt(0),
			Allowance:      big.NewInt(0),
			LockedBalance:  big.NewInt(0),
			PendingBalance: big.NewInt(0),
		}
	}

	if a != nil {
		err = s.AccountDao.Create(a)
		if err != nil {
			logger.Error(err)
			return err
		}
	}

	return nil
}

func (s *AccountService) FindOrCreate(addr common.Address) (*types.Account, error) {
	a, err := s.AccountDao.GetByAddress(addr)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	if a != nil {
		return a, nil
	}

	tokens, err := s.TokenDao.GetAll()
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	a = &types.Account{
		Address:       addr,
		IsBlocked:     false,
		TokenBalances: make(map[common.Address]*types.TokenBalance),
	}

	// currently by default, the tokens balances are set to 0
	for _, t := range tokens {
		a.TokenBalances[t.ContractAddress] = &types.TokenBalance{
			Address:        t.ContractAddress,
			Symbol:         t.Symbol,
			Balance:        big.NewInt(0),
			Allowance:      big.NewInt(0),
			LockedBalance:  big.NewInt(0),
			PendingBalance: big.NewInt(0),
		}
	}

	err = s.AccountDao.Create(a)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	return a, nil
}

func (s *AccountService) GetByID(id bson.ObjectId) (*types.Account, error) {
	return s.AccountDao.GetByID(id)
}

func (s *AccountService) GetAll() ([]types.Account, error) {
	return s.AccountDao.GetAll()
}

func (s *AccountService) GetByAddress(a common.Address) (*types.Account, error) {
	return s.AccountDao.GetByAddress(a)
}

func (s *AccountService) GetTokenBalance(owner common.Address, token common.Address) (*types.TokenBalance, error) {
	return s.AccountDao.GetTokenBalance(owner, token)
}

func (s *AccountService) GetTokenBalances(owner common.Address) (map[common.Address]*types.TokenBalance, error) {
	return s.AccountDao.GetTokenBalances(owner)
}
