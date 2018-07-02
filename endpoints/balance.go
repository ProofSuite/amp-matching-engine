package endpoints

import (
	"github.com/Proofsuite/amp-matching-engine/errors"
	"github.com/Proofsuite/amp-matching-engine/services"
	"github.com/ethereum/go-ethereum/common"
	"github.com/go-ozzo/ozzo-routing"
)

type balanceEndpoint struct {
	balanceService *services.BalanceService
}

// ServeBalance sets up the routing of balance endpoints and the corresponding handlers.
func ServeBalanceResource(rg *routing.RouteGroup, balanceService *services.BalanceService) {
	r := &balanceEndpoint{balanceService}
	rg.Get("/balances/<addr>", r.get)
	// rg.Post("/balances", r.create)
}

// func (r *balanceEndpoint) create(c *routing.Context) error {
// 	var model types.Balance
// 	if err := c.Read(&model); err != nil {
// 		return err
// 	}
// 	if err := model.Validate(); err != nil {
// 		return err
// 	}
// 	err := r.balanceService.Create(&model)
// 	if err != nil {
// 		return err
// 	}

// 	return c.Write(model)
// }

func (r *balanceEndpoint) query(c *routing.Context) error {

	response, err := r.balanceService.GetAll()
	if err != nil {
		return err
	}

	return c.Write(response)
}

func (r *balanceEndpoint) get(c *routing.Context) error {
	addr := c.Param("addr")
	if !common.IsHexAddress(addr) {
		return errors.NewAPIError(400, "INVALID_ADDRESS", nil)
	}
	response, err := r.balanceService.GetByAddress(addr)
	if err != nil {
		return err
	}

	return c.Write(response)
}
