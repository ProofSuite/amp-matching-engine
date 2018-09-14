package endpoints

import (
	"github.com/Proofsuite/amp-matching-engine/errors"
	"github.com/Proofsuite/amp-matching-engine/interfaces"
	"github.com/Proofsuite/amp-matching-engine/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/go-ozzo/ozzo-routing"
)

type accountEndpoint struct {
	accountService interfaces.AccountService
}

func ServeAccountResource(
	r *routing.RouteGroup,
	accountService interfaces.AccountService,
) {

	e := &accountEndpoint{accountService}
	r.Post("/account", e.create)
	r.Get("/account/<address>", e.get)
}

func (e *accountEndpoint) create(c *routing.Context) error {
	account := &types.Account{}
	err := c.Read(&account)
	if err != nil {
		return errors.NewHTTPError(400, "Invalid payload", nil)
	}

	err = account.Validate()
	if err != nil {
		logger.Error(err)
		return err
	}

	err = e.accountService.Create(account)
	if err != nil {
		logger.Error(err)
		return errors.NewHTTPError(400, "Internal server error", nil)
	}

	return c.Write(account)
}

func (e *accountEndpoint) get(c *routing.Context) error {
	a := c.Param("address")
	if !common.IsHexAddress(a) {
		return errors.NewHTTPError(400, "Invalid Address", nil)
	}

	address := common.HexToAddress(a)
	account, err := e.accountService.GetByAddress(address)
	if err != nil {
		logger.Error(err)
		return errors.NewHTTPError(400, "Internal Server Error", nil)
	}

	return c.Write(account)
}

func (e *accountEndpoint) getBalance(c *routing.Context) error {
	a := c.Param("address")
	if !common.IsHexAddress(a) {
		return errors.NewHTTPError(400, "Invalid Address", nil)
	}

	t := c.Param("token")
	if !common.IsHexAddress(a) {
		return errors.NewHTTPError(400, "Invalid Token Address", nil)
	}

	addr := common.HexToAddress(a)
	tokenAddr := common.HexToAddress(t)

	balance, err := e.accountService.GetTokenBalance(addr, tokenAddr)
	if err != nil {
		logger.Error(err)
		return errors.NewHTTPError(400, "Internal Server Error", nil)
	}

	return c.Write(balance)
}
