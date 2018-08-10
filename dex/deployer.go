package dex

import (
	"context"
	"fmt"
	"math/big"

	"github.com/Proofsuite/amp-matching-engine/dex/interfaces"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/backends"
	. "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	
)
type Deployer struct {
	Wallet  *Wallet
	Backend bind.ContractBackend
}

// NewDefaultDeployer returns a deployer connected to the local node via HTTP
// (on port 8545) with the first wallet in the configuration
func NewDefaultDeployer() (*Deployer, error) {
	w := config.Wallets[0]
	conn, err := rpc.DialHTTP("http://127.0.0.1:8545")
	if err != nil {
		return nil, err
	}

	backend := ethclient.NewClient(conn)

	return &Deployer{
		Wallet:  w,
		Backend: backend,
	}, nil
}

// NewDeployer returns a deployer connected to the local node via HTTP
// (on port 8545) with the first wallet in the configuration
func NewDeployer(w *Wallet) (*Deployer, error) {
	conn, err := rpc.DialHTTP("http://127.0.0.1:8545")
	if err != nil {
		return nil, err
	}

	backend := ethclient.NewClient(conn)

	return &Deployer{
		Wallet:  w,
		Backend: backend,
	}, nil
}

// NewWebsocketDeployer returns a deployer connected to the local node via websocket
// (on port 8546). The given wallet is used for signing transactions
func NewWebsocketDeployer(w *Wallet) (*Deployer, error) {
	client, err := rpc.DialWebsocket(context.Background(), "ws://127.0.0.1:8546", "")
	if err != nil {
		return nil, err
	}

	backend := ethclient.NewClient(client)

	return &Deployer{
		Wallet:  w,
		Backend: backend,
	}, nil

}

// NewDefaultSimulator returns a simulated deployer useful for unit testing certain functions
// This simulator functions different from a standard deployer. It does not call a blockchain
// and uses a fake backend.
func NewDefaultSimulator() (*Deployer, error) {
	weiBalance := &big.Int{}
	ether := big.NewInt(1e18)
	etherBalance := big.NewInt(1000)
	wallet := config.Wallets[0]

	genesisAlloc := make(core.GenesisAlloc)
	weiBalance.Mul(etherBalance, ether)

	for _, a := range config.Accounts {
		(genesisAlloc)[a] = core.GenesisAccount{Balance: weiBalance}
	}

	simulator := backends.NewSimulatedBackend(genesisAlloc)

	return &Deployer{
		Wallet:  wallet,
		Backend: simulator,
	}, nil
}

// NewSimulator returns a simulated deployer. The given wallet is used for signing transactions.
// Each ethereum address in the list of given accounts is funded with one ether
func NewSimulator(wallet *Wallet, accounts []Address) (*Deployer, error) {
	genesisAlloc := make(core.GenesisAlloc)

	for _, account := range accounts {
		(genesisAlloc)[account] = core.GenesisAccount{Balance: big.NewInt(1e18)}
	}

	simulator := backends.NewSimulatedBackend(genesisAlloc)

	return &Deployer{
		Wallet:  wallet,
		Backend: simulator,
	}, nil
}

// DeployToken deploys a mock ERC20 token. The given `receiver` address receives `amount` tokens. This function
// makes a complete deployment
func (d Deployer) DeployToken(receiver Address, amount *big.Int) (*ERC20Token, *types.Transaction, error) {
	callOptions := &bind.CallOpts{Pending: true}
	txOptions := bind.NewKeyedTransactor(d.Wallet.PrivateKey)

	address, tx, token, err := interfaces.DeployToken(txOptions, d.Backend, receiver, amount)
	if err != nil && err.Error() == "replacement transaction underpriced" {
		txOptions.Nonce = d.GetNonce()
		address, tx, token, err = interfaces.DeployToken(txOptions, d.Backend, receiver, amount)
	} else if err != nil {
		return nil, nil, err
	}

	return &ERC20Token{
		Address:       address,
		Contract:      token,
		CallOptions:   callOptions,
		TxOptions:     txOptions,
		DefaultSender: d.Wallet,
	}, tx, nil
}

// NewTokens returns a mock ERC20 instance from a given address. This does not deploy any new code on the chain
func (d Deployer) NewToken(address Address) (*ERC20Token, error) {
	instance, err := interfaces.NewToken(address, d.Backend)
	if err != nil {
		return nil, err
	}

	callOptions := &bind.CallOpts{Pending: true}
	txOptions := bind.NewKeyedTransactor(d.Wallet.PrivateKey)

	return &ERC20Token{
		Address:       address,
		Contract:      instance,
		CallOptions:   callOptions,
		TxOptions:     txOptions,
		DefaultSender: d.Wallet,
	}, nil
}

// DeployExchange deploys and returns a new instance of the decentralized exchange contract. The
// exchange is deployed with the given fee account which will receive the trading fees.
func (d Deployer) DeployExchange(feeAccount Address) (*Exchange, *types.Transaction, error) {
	callOpts := &bind.CallOpts{Pending: true}
	txOpts := bind.NewKeyedTransactor(d.Wallet.PrivateKey)

	addr, tx, ex, err := interfaces.DeployExchange(txOpts, d.Backend, feeAccount)
	if err != nil && err.Error() == "replacement transaction underpriced" {
		txOpts.Nonce = d.GetNonce()
		addr, tx, ex, err = interfaces.DeployExchange(txOpts, d.Backend, feeAccount)
	} else if err != nil {
		return nil, nil, err
	}

	return &Exchange{
		Address:     addr,
		Contract:    ex,
		CallOptions: callOpts,
		TxOptions:   txOpts,
		Admin:       d.Wallet,
	}, tx, nil
}

// NewExchange returns a new instance of the excha
func (d Deployer) NewExchange(address Address) (*Exchange, error) {
	instance, err := interfaces.NewExchange(address, d.Backend)
	if err != nil {
		return nil, err
	}

	callOptions := &bind.CallOpts{Pending: true}
	txOptions := bind.NewKeyedTransactor(d.Wallet.PrivateKey)

	return &Exchange{
		Address:     address,
		Contract:    instance,
		CallOptions: callOptions,
		TxOptions:   txOptions,
		Admin:       d.Wallet,
	}, nil
}

func (d Deployer) GetNonce() *big.Int {
	ctx := context.Background()
	n, err := d.Backend.PendingNonceAt(ctx, d.Wallet.Address)
	if err != nil {
		fmt.Printf("Error retrieving the account nonce: %v", err)
	}

	return big.NewInt(0).SetUint64(n)
}

func (d Deployer) WaitMined(tx *types.Transaction) (*types.Receipt, error) {
	ctx := context.Background()
	backend := d.Backend.(bind.DeployBackend)

	receipt, err := bind.WaitMined(ctx, backend, tx)
	if err != nil {
		return nil, err
	}

	return receipt, nil
}
