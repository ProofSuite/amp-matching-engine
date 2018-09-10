package testutils

import (
	"context"
	"log"
	"math/big"

	"github.com/Proofsuite/amp-matching-engine/contracts"
	"github.com/Proofsuite/amp-matching-engine/contracts/contractsinterfaces"
	"github.com/Proofsuite/amp-matching-engine/interfaces"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
)

type Deployer struct {
	WalletService interfaces.WalletService
	TxService     interfaces.TxService
	Client        bind.ContractBackend
}

func NewDeployer(
	w interfaces.WalletService,
	tx interfaces.TxService,
	client bind.ContractBackend,
) *Deployer {
	return &Deployer{
		WalletService: w,
		TxService:     tx,
		Client:        client,
	}
}

func NewDefaultDeployer(w interfaces.WalletService, tx interfaces.TxService) (*Deployer, error) {
	conn, err := rpc.DialHTTP("http://127.0.0.1:8545")
	if err != nil {
		return nil, err
	}

	client := ethclient.NewClient(conn)

	return &Deployer{
		WalletService: w,
		TxService:     tx,
		Client:        client,
	}, nil
}

func NewWebSocketDeployer(w interfaces.WalletService, tx interfaces.TxService) (*Deployer, error) {
	conn, err := rpc.DialWebsocket(context.Background(), "ws://127.0.0.1:8546", "")
	if err != nil {
		return nil, err
	}

	client := ethclient.NewClient(conn)

	return &Deployer{
		WalletService: w,
		TxService:     tx,
		Client:        client,
	}, nil
}

// func NewSimulator(accs []common.Address) {
// 	genesisAlloc := make(core.GenesisAlloc)
// }

// // NewDefaultSimulator returns a simulated deployer useful for unit testing certain functions
// // This simulator functions different from a standard deployer. It does not call a blockchain
// // and uses a fake backend.
// func NewSimulator(
// 	w interfaces.WalletService,
// 	tx interfaces.TxService,
// 	accs []common.Address,
// ) (*Deployer, error) {
// 	weiBalance := &big.Int{}
// 	ether := big.NewInt(1e18)
// 	etherBalance := big.NewInt(1000)

// 	genesisAlloc := make(core.GenesisAlloc)
// 	weiBalance.Mul(etherBalance, ether)

// 	for _, a := range accs {
// 		(genesisAlloc)[a] = core.GenesisAccount{Balance: weiBalance}
// 	}

// 	simulator := NewSimulatedBackend(genesisAlloc, 5e6)

// 	return &Deployer{
// 		WalletService: w,
// 		TxService:     tx,
// 		Backend:       simulator,
// 	}, nil
// }

// func NewSimulatedDeployer(
// 	w interfaces.WalletService,
// 	tx interfaces.TxService,
// 	client adfadfasdf
// )

// DeployToken
func (d *Deployer) DeployToken(receiver common.Address, amount *big.Int) (*contracts.Token, common.Address, *ethTypes.Transaction, error) {
	// callOptions := d.TxService.GetTxCallOptions()
	sendOptions, _ := d.TxService.GetTxSendOptions()

	address, tx, tokenInterface, err := contractsinterfaces.DeployToken(sendOptions, d.Client, receiver, amount)
	if err != nil && err.Error() == "replacement transaction underpriced" {
		sendOptions.Nonce, _ = d.GetNonce()
		address, tx, tokenInterface, err = contractsinterfaces.DeployToken(sendOptions, d.Client, receiver, amount)
	} else if err != nil {
		return nil, common.Address{}, nil, err
	}

	return &contracts.Token{
		WalletService: d.WalletService,
		TxService:     d.TxService,
		Interface:     tokenInterface,
	}, address, tx, nil
}

func (d *Deployer) NewToken(addr common.Address) (*contracts.Token, error) {
	tokenInterface, err := contractsinterfaces.NewToken(addr, d.Client)
	if err != nil {
		return nil, err
	}

	return &contracts.Token{
		WalletService: d.WalletService,
		TxService:     d.TxService,
		Interface:     tokenInterface,
	}, nil
}

// DeployExchange
func (d *Deployer) DeployExchange(wethToken common.Address, feeAccount common.Address) (*contracts.Exchange, common.Address, *ethTypes.Transaction, error) {
	sendOptions, _ := d.TxService.GetTxSendOptions()

	address, tx, exchangeInterface, err := contractsinterfaces.DeployExchange(sendOptions, d.Client, wethToken, feeAccount)
	if err != nil && err.Error() == "replacement transaction underpriced" {
		sendOptions.Nonce, _ = d.GetNonce()
		address, tx, exchangeInterface, err = contractsinterfaces.DeployExchange(sendOptions, d.Client, wethToken, feeAccount)
		if err != nil {
			return nil, common.Address{}, nil, err
		}
	} else if err != nil {
		return nil, common.Address{}, nil, err
	}

	return &contracts.Exchange{
		WalletService: d.WalletService,
		Interface:     exchangeInterface,
		Client:        d.Client,
	}, address, tx, err
}

// NewExchange
func (d *Deployer) NewExchange(addr common.Address) (*contracts.Exchange, error) {
	exchangeInterface, err := contractsinterfaces.NewExchange(addr, d.Client)
	if err != nil {
		return nil, err
	}

	return &contracts.Exchange{
		WalletService: d.WalletService,
		Interface:     exchangeInterface,
		Client:        d.Client,
	}, nil
}

// GetNonce
func (d *Deployer) GetNonce() (*big.Int, error) {
	ctx := context.Background()

	wallet, err := d.WalletService.GetDefaultAdminWallet()
	if err != nil {
		log.Print(err)
		return nil, err
	}

	n, err := d.Client.PendingNonceAt(ctx, wallet.Address)
	if err != nil {
		log.Print(err)
		return nil, err
	}

	return big.NewInt(0).SetUint64(n), nil
}

func (d *Deployer) WaitMined(tx *ethTypes.Transaction) (*ethTypes.Receipt, error) {
	ctx := context.Background()
	backend := d.Client.(bind.DeployBackend)

	receipt, err := bind.WaitMined(ctx, backend, tx)
	if err != nil {
		return nil, err
	}

	return receipt, nil
}

// type SimulatedBackend struct {
// 	*ethereum.SimulatedClient
// }

// func (b *SimulatedBackend) PendingBalanceAt(ctx context.Context, account common.Address) (*big.Int, error) {
// 	return nil, errors.New("PendingBalanceAt is not implemented on the simulated backend")
// }

// func NewSimulatedBackend(alloc core.GenesisAlloc, gasLimit uint64) *SimulatedBackend {
// 	return &SimulatedBackend{backends.NewSimulatedBackend(alloc, gasLimit)} //, gasLimit)}
// }
