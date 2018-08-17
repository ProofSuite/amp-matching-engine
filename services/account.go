package services

import (
	"errors"
	"math/big"

	"gopkg.in/mgo.v2/bson"

	"github.com/Proofsuite/amp-matching-engine/daos"
	"github.com/Proofsuite/amp-matching-engine/types"
	"github.com/ethereum/go-ethereum/common"
)

type AccountService struct {
	AccountDao *daos.AccountDao
	TokenDao   *daos.TokenDao
}

// NewAddressService returns a new instance of accountService
func NewAccountService(AccountDao *daos.AccountDao, TokenDao *daos.TokenDao) *AccountService {
	return &AccountService{AccountDao, TokenDao}
}

func (s *AccountService) Create(addr common.Address) error {
	acc, err := s.GetByAddress(addr)
	if err != nil {
		return err
	}

	if acc != nil {
		return errors.New("Account already exists")
	}

	tokens, err := s.TokenDao.GetAll()
	if err != nil {
		return err
	}

	acc = &types.Account{}
	acc.Address = addr
	acc.IsBlocked = false

	// currently by default, the tokens balances are set to 0
	for _, token := range tokens {
		acc.TokenBalances[token.ContractAddress] = &types.TokenBalance{
			TokenID:       token.ID,
			TokenAddress:  token.ContractAddress,
			Balance:       big.NewInt(0),
			Allowance:     big.NewInt(0),
			LockedBalance: big.NewInt(0),
		}
	}

	if acc != nil {
		err = s.AccountDao.Create(acc)
		if err != nil {
			return err
		}
	}

	return nil
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
