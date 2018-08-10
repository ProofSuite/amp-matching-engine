package endpoints

import (
	"github.com/Proofsuite/amp-matching-engine/errors"
	"github.com/Proofsuite/amp-matching-engine/services"
	"github.com/Proofsuite/amp-matching-engine/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/go-ozzo/ozzo-routing"
)

type addressEndpoint struct {
	addressService *services.AddressService
}

// ServeAddressResource sets up the routing of address endpoints and the corresponding handlers.
func ServeAddressResource(rg *routing.RouteGroup, addressService *services.AddressService) {
	r := &addressEndpoint{addressService}
	rg.Post("/address", r.create)
	rg.Get("/address/<addr>/nonce", r.getNonce)
}

func (r *addressEndpoint) create(c *routing.Context) error {
	var model types.UserAddress
	if err := c.Read(&model); err != nil {
		return errors.NewAPIError(400, "INVALID_DATA", map[string]interface{}{
			"details": err.Error(),
		})
	}
	if err := model.Validate(); err != nil {
		return err
	}
	err := r.addressService.Create(&model)
	if err != nil {
		return err
	}

	return c.Write(model)
}
func (r *addressEndpoint) getNonce(c *routing.Context) error {
	addr := c.Param("addr")
	if !common.IsHexAddress(addr) {
		return errors.NewAPIError(400, "INVALID_ADDRESS", nil)
	}
	nonce, err := r.addressService.GetNonce(addr)
	if err != nil {
		return err
	}

	return c.Write(nonce)
}
