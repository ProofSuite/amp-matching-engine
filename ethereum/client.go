package ethereum

import (
	"context"
	"errors"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind/backends"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
)

type SimulatedClient struct {
	*backends.SimulatedBackend
}

func (b *SimulatedClient) PendingBalanceAt(ctx context.Context, acc common.Address) (*big.Int, error) {
	return nil, errors.New("PendingBalanceAt is not implemented on the simulated backend")
}

func NewSimulatedClient(accs []common.Address) *SimulatedClient {
	weiBalance := &big.Int{}
	ether := big.NewInt(1e18)
	etherBalance := big.NewInt(1000)

	alloc := make(core.GenesisAlloc)
	weiBalance.Mul(etherBalance, ether)

	for _, a := range accs {
		(alloc)[a] = core.GenesisAccount{Balance: weiBalance}
	}

	client := backends.NewSimulatedBackend(alloc, 5e6)
	return &SimulatedClient{client}
}
