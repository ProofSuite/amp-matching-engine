package services

import (
	"gopkg.in/mgo.v2/bson"

	"github.com/Proofsuite/amp-matching-engine/daos"
	"github.com/Proofsuite/amp-matching-engine/types"
)

type TokenService struct {
	tokenDao *daos.TokenDao
}

func NewTokenService(tokenDao *daos.TokenDao) *TokenService {
	return &TokenService{tokenDao}
}

func (s *TokenService) Create(token *types.Token) error {
	return s.tokenDao.Create(token)

}

func (s *TokenService) GetByID(id bson.ObjectId) (*types.Token, error) {
	return s.tokenDao.GetByID(id)
}

func (s *TokenService) GetAll() ([]types.Token, error) {
	return s.tokenDao.GetAll()
}
