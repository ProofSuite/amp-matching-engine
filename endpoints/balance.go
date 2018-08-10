package endpoints

import (
	"strconv"

	"github.com/Proofsuite/amp-matching-engine/errors"
	"github.com/Proofsuite/amp-matching-engine/services"
	"github.com/ethereum/go-ethereum/common"
	"github.com/go-ozzo/ozzo-routing"
)

type balanceEndpoint struct {
	balanceService *services.BalanceService
}

// ServeBalanceResource sets up the routing of balance endpoints and the corresponding handlers.
func ServeBalanceResource(rg *routing.RouteGroup, balanceService *services.BalanceService) {
	r := &balanceEndpoint{balanceService}
	rg.Get("/balances/<addr>", r.get)
}

func (r *balanceEndpoint) get(c *routing.Context) error {
	addr := c.Param("addr")
	if !common.IsHexAddress(addr) {
		return errors.NewAPIError(400, "INVALID_ADDRESS", nil)
	}
	nonZero := c.Query("nonZero", "false")
	nonZeroBool, err := strconv.ParseBool(nonZero)
	if err != nil {
		return errors.NewAPIError(400, "INVALID_NONZERO_QUERY_PARAM", nil)
	}
	response, err := r.balanceService.GetByAddress(addr, nonZeroBool)
	if err != nil {
		return err
	}

	return c.Write(response)
}
