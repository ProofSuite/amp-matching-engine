package services

import (
	"fmt"

	"gopkg.in/mgo.v2/bson"

	"github.com/Proofsuite/amp-matching-engine/daos"
	"github.com/Proofsuite/amp-matching-engine/types"
)

// AddressService struct with daos required, responsible for communicating with daos
type AddressService struct {
	AddressDao *daos.AddressDao
	balanceDao *daos.BalanceDao
	tokenDao   *daos.TokenDao
}

// NewAddressService returns a new instance of addressService
func NewAddressService(AddressDao *daos.AddressDao, balanceDao *daos.BalanceDao, tokenDao *daos.TokenDao) *AddressService {
	return &AddressService{AddressDao, balanceDao, tokenDao}
}

// Create validates the address and create wallet for the address
func (s *AddressService) Create(address *types.UserAddress) error {
	ua, err := s.GetByAddress(address.Address)
	if err == nil && ua != nil {
		address.ID = ua.ID
		address.Address = ua.Address
		address.IsBlocked = ua.IsBlocked
		address.CreatedAt = ua.CreatedAt
		address.UpdatedAt = ua.UpdatedAt
		return nil
	}
	err = s.AddressDao.Create(address)
	if err != nil {
		return fmt.Errorf("%s<==", err)
	}
	balService := NewBalanceService(s.balanceDao, s.tokenDao)
	bal := &types.Balance{Address: address.Address}
	err = balService.Create(bal)
	if err != nil {
		return err
	}
	return err

}

// GetByID fetches the address's details based on its mongoID
func (s *AddressService) GetByID(id bson.ObjectId) (*types.UserAddress, error) {
	return s.AddressDao.GetByID(id)
}

// GetAll fetches all the address entries in the db
func (s *AddressService) GetAll() ([]types.UserAddress, error) {
	return s.AddressDao.GetAll()
}

// GetByAddress fetches the address's details
func (s *AddressService) GetByAddress(addr string) (*types.UserAddress, error) {
	return s.AddressDao.GetByAddress(addr)
}
