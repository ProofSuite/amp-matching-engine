package services

import (
	"gopkg.in/mgo.v2/bson"

	"github.com/Proofsuite/amp-matching-engine/daos"
	"github.com/Proofsuite/amp-matching-engine/types"
)

// TokenService struct with daos required, responsible for communicating with daos.
// TokenService functions are responsible for interacting with daos and implements business logics.
type TokenService struct {
	tokenDao *daos.TokenDao
}

// NewTokenService returns a new instance of TokenService
func NewTokenService(tokenDao *daos.TokenDao) *TokenService {
	return &TokenService{tokenDao}
}

// Create inserts a new token into the database
func (s *TokenService) Create(token *types.Token) error {
	return s.tokenDao.Create(token)

}

// GetByID fetches the detailed document of a token using its mongo ID
func (s *TokenService) GetByID(id bson.ObjectId) (*types.Token, error) {
	return s.tokenDao.GetByID(id)
}

// GetByAddress fetches the detailed document of a token using its contract address
func (s *TokenService) GetByAddress(addr string) (*types.Token, error) {
	return s.tokenDao.GetByAddress(addr)
}

// GetAll fetches all the tokens from db
func (s *TokenService) GetAll() ([]types.Token, error) {
	return s.tokenDao.GetAll()
}
