package services

import (
	"github.com/Proofsuite/amp-matching-engine/errors"
	"github.com/Proofsuite/amp-matching-engine/interfaces"
	"github.com/ethereum/go-ethereum/common"
	"gopkg.in/mgo.v2/bson"

	"github.com/Proofsuite/amp-matching-engine/types"
)

// TokenService struct with daos required, responsible for communicating with daos.
// TokenService functions are responsible for interacting with daos and implements business logics.
type TokenService struct {
	tokenDao interfaces.TokenDao
}

// NewTokenService returns a new instance of TokenService
func NewTokenService(tokenDao interfaces.TokenDao) *TokenService {
	return &TokenService{tokenDao}
}

// Create inserts a new token into the database
func (s *TokenService) Create(token *types.Token) error {
	t, err := s.tokenDao.GetByAddress(token.ContractAddress)
	if err != nil {
		return err
	}

	if t != nil {
		return errors.NewAPIError(401, "TOKEN_ALREADY_EXISTS", nil)
	}

	return s.tokenDao.Create(token)
}

// GetByID fetches the detailed document of a token using its mongo ID
func (s *TokenService) GetByID(id bson.ObjectId) (*types.Token, error) {
	return s.tokenDao.GetByID(id)
}

// GetByAddress fetches the detailed document of a token using its contract address
func (s *TokenService) GetByAddress(addr common.Address) (*types.Token, error) {
	return s.tokenDao.GetByAddress(addr)
}

// GetAll fetches all the tokens from db
func (s *TokenService) GetAll() ([]types.Token, error) {
	return s.tokenDao.GetAll()
}
