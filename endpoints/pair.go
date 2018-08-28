package endpoints

import (
	"github.com/ethereum/go-ethereum/common"

	"github.com/Proofsuite/amp-matching-engine/errors"
	"github.com/Proofsuite/amp-matching-engine/services"
	"github.com/Proofsuite/amp-matching-engine/types"
	"github.com/go-ozzo/ozzo-routing"
)

type pairEndpoint struct {
	pairService services.PairServiceInterface
}

// ServePairResource sets up the routing of pair endpoints and the corresponding handlers.
func ServePairResource(
	rg *routing.RouteGroup,
	p services.PairServiceInterface,
) {
	r := &pairEndpoint{p}
	rg.Get("/pairs/<baseToken>/<quoteToken>", r.get)
	rg.Get("/pairs", r.query)
	rg.Post("/pairs", r.create)
}

func (r *pairEndpoint) create(c *routing.Context) error {
	var p types.Pair

	if err := c.Read(&p); err != nil {
		return err
	}

	if err := p.Validate(); err != nil {
		return err
	}

	err := r.pairService.Create(&p)
	if err != nil {
		return err
	}

	return c.Write(p)
}

func (r *pairEndpoint) query(c *routing.Context) error {
	res, err := r.pairService.GetAll()
	if err != nil {
		return err
	}

	return c.Write(res)
}

func (r *pairEndpoint) get(c *routing.Context) error {
	baseToken := c.Param("baseToken")
	if !common.IsHexAddress(baseToken) {
		return errors.NewAPIError(400, "INVALID_HEX_ADDRESS", nil)
	}

	quoteToken := c.Param("quoteToken")
	if !common.IsHexAddress(quoteToken) {
		return errors.NewAPIError(400, "INVALID_HEX_ADDRESS", nil)
	}

	baseTokenAddress := common.HexToAddress(baseToken)
	quoteTokenAddress := common.HexToAddress(quoteToken)

	res, err := r.pairService.GetByTokenAddress(baseTokenAddress, quoteTokenAddress)
	if err != nil {
		return err
	}

	return c.Write(res)
}
