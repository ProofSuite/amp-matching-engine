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

func NewDefaultDeployer() (*Deployer, error) {
	w := config.Wallets[0]
	conn, err := rpc.DialHTTP("ws://127.0.0.1:8546")
	if err != nil {
		return nil, err
	}

	backend := ethclient.NewClient(conn)

	return &Deployer{
		Wallet:  w,
		Backend: backend,
	}, nil
}

func NewDeployer(w *Wallet) (*Deployer, error) {
	conn, err := rpc.DialHTTP("ws://127.0.0.1:8546")
	if err != nil {
		return nil, err
	}

	backend := ethclient.NewClient(conn)

	return &Deployer{
		Wallet:  w,
		Backend: backend,
	}, nil
}

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
		Address:       address,
		Contract:      token,
		CallOptions:   callOptions,
		TxOptions:     txOptions,
		DefaultSender: d.Wallet,
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
		Address:       address,
		Contract:      instance,
		CallOptions:   callOptions,
		TxOptions:     txOptions,
		DefaultSender: d.Wallet,
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
		Admin:       d.Wallet,
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
		Admin:       d.Wallet,
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
