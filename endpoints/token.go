package endpoints

import (
	"log"

	"github.com/Proofsuite/amp-matching-engine/errors"
	"github.com/Proofsuite/amp-matching-engine/services"
	"github.com/Proofsuite/amp-matching-engine/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/go-ozzo/ozzo-routing"
)

type tokenEndpoint struct {
	tokenService services.TokenServiceInterface
}

// ServeTokenResource sets up the routing of token endpoints and the corresponding handlers.
func ServeTokenResource(rg *routing.RouteGroup, tokenService services.TokenServiceInterface) {
	r := &tokenEndpoint{tokenService}
	rg.Get("/tokens/<address>", r.get)
	rg.Get("/tokens", r.query)
	rg.Post("/tokens", r.create)
}

func (r *tokenEndpoint) create(c *routing.Context) error {
	var model types.Token
	if err := c.Read(&model); err != nil {
		log.Print(err)
		return err
	}

	err := r.tokenService.Create(&model)
	if err != nil {
		log.Print(err)
		return err
	}

	return c.Write(model)
}

func (r *tokenEndpoint) query(c *routing.Context) error {
	response, err := r.tokenService.GetAll()
	if err != nil {
		log.Print(err)
		return err
	}

	return c.Write(response)
}

func (r *tokenEndpoint) get(c *routing.Context) error {
	a := c.Param("address")
	if !common.IsHexAddress(a) {
		return errors.NewAPIError(400, "INVALID_ID", nil)
	}

	tokenAddress := common.HexToAddress(a)
	response, err := r.tokenService.GetByAddress(tokenAddress)
	if err != nil {
		log.Print(err)
		return err
	}

	return c.Write(response)
}
