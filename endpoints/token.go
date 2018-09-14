package endpoints

import (
	"github.com/Proofsuite/amp-matching-engine/errors"
	"github.com/Proofsuite/amp-matching-engine/interfaces"
	"github.com/Proofsuite/amp-matching-engine/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/go-ozzo/ozzo-routing"
)

type tokenEndpoint struct {
	tokenService interfaces.TokenService
}

// ServeTokenResource sets up the routing of token endpoints and the corresponding handlers.
func ServeTokenResource(rg *routing.RouteGroup, tokenService interfaces.TokenService) {
	r := &tokenEndpoint{tokenService}
	rg.Get("/tokens/base", r.queryBase)
	rg.Get("/tokens/quote", r.queryQuote)
	rg.Get("/tokens/<address>", r.get)
	rg.Get("/tokens", r.query)
	rg.Post("/tokens", r.create)
}

func (r *tokenEndpoint) create(c *routing.Context) error {
	var model types.Token
	if err := c.Read(&model); err != nil {
		logger.Error(err)
		return errors.NewHTTPError(400, "Invalid Payload", nil)
	}

	err := r.tokenService.Create(&model)
	if err != nil {
		logger.Error(err)
		return errors.NewHTTPError(500, "Internal Server Error", nil)
	}

	return c.Write(model)
}

func (r *tokenEndpoint) query(c *routing.Context) error {
	response, err := r.tokenService.GetAll()
	if err != nil {
		return errors.NewHTTPError(500, "Internal Server Error", nil)
	}

	return c.Write(response)
}

func (r *tokenEndpoint) queryQuote(c *routing.Context) error {
	response, err := r.tokenService.GetQuote()
	if err != nil {
		logger.Error(err)
		return errors.NewHTTPError(500, "Internal Server Error", nil)
	}

	return c.Write(response)
}

func (r *tokenEndpoint) queryBase(c *routing.Context) error {
	response, err := r.tokenService.GetBase()
	if err != nil {
		logger.Error(err)
		return errors.NewHTTPError(500, "Internal Server Error", nil)
	}

	return c.Write(response)
}

func (r *tokenEndpoint) get(c *routing.Context) error {
	a := c.Param("address")
	if !common.IsHexAddress(a) {
		return errors.NewHTTPError(400, "Invalid Address", nil)
	}

	tokenAddress := common.HexToAddress(a)
	response, err := r.tokenService.GetByAddress(tokenAddress)
	if err != nil {
		logger.Error(err)
		return errors.NewHTTPError(500, "Internal Server Error", nil)
	}

	return c.Write(response)
}
