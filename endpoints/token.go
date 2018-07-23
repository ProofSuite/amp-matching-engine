package endpoints

import (
	"github.com/Proofsuite/amp-matching-engine/errors"
	"github.com/Proofsuite/amp-matching-engine/services"
	"github.com/Proofsuite/amp-matching-engine/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/go-ozzo/ozzo-routing"
)

type tokenEndpoint struct {
	tokenService *services.TokenService
}

// ServeTokenResource sets up the routing of token endpoints and the corresponding handlers.
func ServeTokenResource(rg *routing.RouteGroup, tokenService *services.TokenService) {
	r := &tokenEndpoint{tokenService}
	rg.Get("/tokens/<addr>", r.get)
	rg.Get("/tokens", r.query)
	rg.Post("/tokens", r.create)
}

func (r *tokenEndpoint) create(c *routing.Context) error {
	var model types.Token
	if err := c.Read(&model); err != nil {
		return err
	}
	err := r.tokenService.Create(&model)
	if err != nil {
		return err
	}

	return c.Write(model)
}

func (r *tokenEndpoint) query(c *routing.Context) error {

	response, err := r.tokenService.GetAll()
	if err != nil {
		return err
	}

	return c.Write(response)
}

func (r *tokenEndpoint) get(c *routing.Context) error {
	addr := c.Param("addr")
	if !common.IsHexAddress(addr) {
		return errors.NewAPIError(400, "INVALID_ID", nil)
	}
	response, err := r.tokenService.GetByAddress(addr)
	if err != nil {
		return err
	}

	return c.Write(response)
}
