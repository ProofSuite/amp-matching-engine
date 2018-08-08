package services

import (
	"math"

	"github.com/Proofsuite/amp-matching-engine/daos"
	"github.com/Proofsuite/amp-matching-engine/types"
)

// BalanceService struct with daos required, responsible for communicating with daos.
// BalanceService functions are responsible for interacting with daos and implements business logics.
type BalanceService struct {
	balanceDao *daos.BalanceDao
	tokenDao   *daos.TokenDao
}

// NewBalanceService returns a new instance of balance service
func NewBalanceService(balanceDao *daos.BalanceDao, tokenDao *daos.TokenDao) *BalanceService {
	return &BalanceService{balanceDao, tokenDao}
}

// Create function is responsible for creating balance entry corresponding
// to a user address.
func (s *BalanceService) Create(balance *types.Balance) error {
	tb := make(map[string]types.TokenBalance)
	tokens, err := s.tokenDao.GetAll()
	if err != nil {
		return err
	}
	for _, token := range tokens {
		tb[token.ContractAddress] = types.TokenBalance{
			Amount:       int64(10000 * math.Pow10(8)),
			LockedAmount: 0,
			TokenID:      token.ID,
			TokenAddress: token.ContractAddress,
			TokenSymbol:  token.Symbol,
		}
	}
	balance.Tokens = tb
	err = s.balanceDao.Create(balance)
	return err

}

// GetByAddress is responsible for fetching the balance details of a user address
func (s *BalanceService) GetByAddress(addr string) (*types.Balance, error) {
	return s.balanceDao.GetByAddress(addr)
}
