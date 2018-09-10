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
	rg *routing.RouteGroup,
	accountService interfaces.AccountService,
) {
	e := &accountEndpoint{accountService}
	rg.Post("/account", e.create)
	rg.Get("/account/<address>", e.get)
}

func (e *accountEndpoint) create(c *routing.Context) error {
	account := &types.Account{}
	err := c.Read(&account)
	if err != nil {
		return errors.NewAPIError(400, "INVALID_DATA", map[string]interface{}{
			"details": err.Error(),
		})
	}

	err = account.Validate()
	if err != nil {
		logger.Error(err)
		return err
	}

	err = e.accountService.Create(account)
	if err != nil {
		logger.Error(err)
		return errors.NewAPIError(400, "CREATE_ACCOUNT_FAIL", map[string]interface{}{
			"details": err.Error(),
		})
	}

	return c.Write(account)
}

func (e *accountEndpoint) get(c *routing.Context) error {
	a := c.Param("address")
	if !common.IsHexAddress(a) {
		return errors.NewAPIError(400, "INVALID_ADDRESS", nil)
	}

	address := common.HexToAddress(a)
	account, err := e.accountService.GetByAddress(address)
	if err != nil {
		logger.Error(err)
		return errors.NewAPIError(400, "ACCOUNT_ERROR", nil)
	}

	return c.Write(account)
}

func (e *accountEndpoint) getBalance(c *routing.Context) error {
	a := c.Param("address")
	if !common.IsHexAddress(a) {
		return errors.NewAPIError(400, "INVALID_ADDRESS", nil)
	}

	t := c.Param("token")
	if !common.IsHexAddress(a) {
		return errors.NewAPIError(400, "INVALID_TOKEN_ADDRESS", nil)
	}

	addr := common.HexToAddress(a)
	tokenAddr := common.HexToAddress(t)

	balance, err := e.accountService.GetTokenBalance(addr, tokenAddr)
	if err != nil {
		logger.Error(err)
		return errors.NewAPIError(400, "ERROR_GETBALANCE", nil)
	}

	return c.Write(balance)
}
