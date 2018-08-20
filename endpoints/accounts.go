package endpoints

import (
	"fmt"

	"github.com/Proofsuite/amp-matching-engine/errors"
	"github.com/Proofsuite/amp-matching-engine/services"
	"github.com/Proofsuite/amp-matching-engine/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/go-ozzo/ozzo-routing"
)

type accountEndpoint struct {
	accountService *services.AccountService
}

func ServeAccountResource(rg *routing.RouteGroup, accountService *services.AccountService) {
	e := &accountEndpoint{accountService}
	rg.Post("/account", e.create)
	rg.Get("/account/<address>", e.get)
}

func (e *accountEndpoint) create(c *routing.Context) error {

	account := &types.Account{}
	if err := c.Read(&account); err != nil {
		return errors.NewAPIError(400, "INVALID_DATA", map[string]interface{}{
			"details": err.Error(),
		})
	}
	if err := account.Validate(); err != nil {
		return errors.NewAPIError(400, "INVALID_ACCOUNT", map[string]interface{}{
			"details": err.Error(),
		})
	}

	if err := e.accountService.Create(account); err != nil {
		fmt.Println(err)
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
		return errors.NewAPIError(400, "ERROR_GETBALANCE", nil)
	}

	return c.Write(balance)
}
