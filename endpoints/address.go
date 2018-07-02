package endpoints

import (
	"github.com/Proofsuite/amp-matching-engine/services"
	"github.com/Proofsuite/amp-matching-engine/types"
	"github.com/go-ozzo/ozzo-routing"
)

type addressEndpoint struct {
	addressService *services.AddressService
}

// ServeAddress sets up the routing of address endpoints and the corresponding handlers.
func ServeAddressResource(rg *routing.RouteGroup, addressService *services.AddressService) {
	r := &addressEndpoint{addressService}
	// rg.Get("/addresses/<id>", r.get)
	// rg.Get("/addresses", r.query)
	rg.Post("/address", r.create)
}

func (r *addressEndpoint) create(c *routing.Context) error {
	var model types.UserAddress
	if err := c.Read(&model); err != nil {
		return err
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
