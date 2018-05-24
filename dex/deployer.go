package dex

import (
	"context"
	"fmt"
	"math/big"

	"github.com/Dvisacker/matching-engine/dex/interfaces"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/backends"
	. "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
)

type Deployer struct {
	Wallet  *Wallet
	Backend bind.ContractBackend
}

func NewDefaultLocalDeployer() (*Deployer, error) {
	Wallet := config.Wallets[0]
	connection, err := rpc.DialHTTP("http://127.0.0.1:8545")
	if err != nil {
		return nil, err
	}

	backend := ethclient.NewClient(connection)

	return &Deployer{
		Wallet:  Wallet,
		Backend: backend,
	}, nil
}

func NewDefaultSimulator() (*Deployer, error) {
	Wallet := config.Wallets[0]
	genesisAllocation := make(core.GenesisAlloc)

	for _, account := range config.Accounts {
		(genesisAllocation)[account] = core.GenesisAccount{Balance: big.NewInt(1e18)}
	}

	simulator := backends.NewSimulatedBackend(genesisAllocation)

	return &Deployer{
		Wallet:  Wallet,
		Backend: simulator,
	}, nil
}

func (d Deployer) DeployToken(receiver Address, amount *big.Int) (*ERC20Token, error) {
	callOptions := &bind.CallOpts{Pending: true}
	txOptions := bind.NewKeyedTransactor(d.Wallet.PrivateKey)

	address, _, token, err := interfaces.DeployToken(txOptions, d.Backend, receiver, amount)
	if err != nil && err.Error() == "replacement transaction underpriced" {
		txOptions.Nonce = d.GetNonce()
		address, _, token, err = interfaces.DeployToken(txOptions, d.Backend, receiver, amount)
	} else if err != nil {
		return nil, err
	}

	return &ERC20Token{
		Address:     address,
		Contract:    token,
		CallOptions: callOptions,
		TxOptions:   txOptions,
	}, nil
}

func (d Deployer) NewToken(address Address) (*ERC20Token, error) {
	instance, err := interfaces.NewToken(address, d.Backend)
	if err != nil {
		return nil, err
	}

	callOptions := &bind.CallOpts{Pending: true}
	txOptions := bind.NewKeyedTransactor(d.Wallet.PrivateKey)

	return &ERC20Token{
		Address:     address,
		Contract:    instance,
		CallOptions: callOptions,
		TxOptions:   txOptions,
	}, nil
}

func (d Deployer) DeployExchange(feeAccount Address) (*Exchange, error) {
	callOptions := &bind.CallOpts{Pending: true}
	txOptions := bind.NewKeyedTransactor(d.Wallet.PrivateKey)

	address, _, exchange, err := interfaces.DeployExchange(txOptions, d.Backend, feeAccount)
	if err != nil && err.Error() == "replacement transaction underpriced" {
		txOptions.Nonce = d.GetNonce()
		address, _, exchange, err = interfaces.DeployExchange(txOptions, d.Backend, feeAccount)
	} else if err != nil {
		return nil, err
	}

	return &Exchange{
		Address:     address,
		Contract:    exchange,
		CallOptions: callOptions,
		TxOptions:   txOptions,
	}, nil
}

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
	}, nil
}

func (d Deployer) GetNonce() *big.Int {
	context := context.Background()
	nonce, err := d.Backend.PendingNonceAt(context, d.Wallet.Address)
	if err != nil {
		fmt.Printf("Error retrieving the account nonce: %v", err)
	}

	return big.NewInt(0).SetUint64(nonce)
}
