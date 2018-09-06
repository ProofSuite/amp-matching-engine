package ethereum

import (
	"context"
	"log"
	"math/big"

	"github.com/Proofsuite/amp-matching-engine/app"
	"github.com/Proofsuite/amp-matching-engine/contracts/contractsinterfaces"
	"github.com/Proofsuite/amp-matching-engine/interfaces"
	"github.com/Proofsuite/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	eth "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rpc"
)

type EthereumProvider struct {
	Client interfaces.EthereumClient
	Config interfaces.EthereumConfig
}

func NewEthereumProvider(c interfaces.EthereumClient) *EthereumProvider {
	err := app.LoadConfig("../config", "")
	if err != nil {
		panic(err)
	}

	url := app.Config.Ethereum["URL"]
	exchange := common.HexToAddress(app.Config.Ethereum["exchange_address"])
	weth := common.HexToAddress(app.Config.Ethereum["weth_address"])
	config := NewEthereumConfig(url, exchange, weth)

	return &EthereumProvider{
		Client: c,
		Config: config,
	}
}

func NewDefaultEthereumProvider() *EthereumProvider {
	err := app.LoadConfig("./config", "")
	if err != nil {
		panic(err)
	}

	url := app.Config.Ethereum["URL"]
	exchange := common.HexToAddress(app.Config.Ethereum["exchange_address"])
	weth := common.HexToAddress(app.Config.Ethereum["weth_address"])

	conn, err := rpc.DialHTTP(app.Config.Ethereum["URL"])
	if err != nil {
		panic(err)
	}

	client := ethclient.NewClient(conn)
	config := NewEthereumConfig(url, exchange, weth)

	return &EthereumProvider{
		Client: client,
		Config: config,
	}
}

func NewSimulatedEthereumProvider(accs []common.Address) *EthereumProvider {
	url := app.Config.Ethereum["URL"]
	exchange := common.HexToAddress(app.Config.Ethereum["exchange_address"])
	weth := common.HexToAddress(app.Config.Ethereum["weth_address"])

	err := app.LoadConfig("../config", "")
	if err != nil {
		panic(err)
	}

	config := NewEthereumConfig(url, exchange, weth)
	client := NewSimulatedClient(accs)

	return &EthereumProvider{
		Client: client,
		Config: config,
	}
}

func (e *EthereumProvider) WaitMined(tx *eth.Transaction) (*eth.Receipt, error) {
	ctx := context.Background()
	receipt, err := bind.WaitMined(ctx, e.Client, tx)
	if err != nil {
		log.Print(err)
		return &eth.Receipt{}, err
	}

	return receipt, nil
}

func (e *EthereumProvider) GetBalanceAt(a common.Address) (*big.Int, error) {
	ctx := context.Background()
	nonce, err := e.Client.BalanceAt(ctx, a, nil)
	if err != nil {
		log.Print(err)
		return big.NewInt(0), err
	}

	return nonce, nil
}

func (e *EthereumProvider) GetPendingNonceAt(a common.Address) (uint64, error) {
	ctx := context.Background()
	nonce, err := e.Client.PendingNonceAt(ctx, a)
	if err != nil {
		log.Print(err)
		return 0, err
	}

	return nonce, nil
}

func (e *EthereumProvider) BalanceOf(owner common.Address, token common.Address) (*big.Int, error) {
	tokenInterface, err := contractsinterfaces.NewToken(token, e.Client)
	if err != nil {
		return nil, err
	}

	opts := &bind.CallOpts{Pending: true}
	b, err := tokenInterface.BalanceOf(opts, owner)
	if err != nil {
		log.Print(err)
		return nil, err
	}

	return b, nil
}

func (e *EthereumProvider) Allowance(owner, spender, token common.Address) (*big.Int, error) {
	tokenInterface, err := contractsinterfaces.NewToken(token, e.Client)
	if err != nil {
		return nil, err
	}

	opts := &bind.CallOpts{Pending: true}
	a, err := tokenInterface.Allowance(opts, owner, spender)
	if err != nil {
		log.Print(err)
		return nil, err
	}

	return a, nil
}

func (e *EthereumProvider) ExchangeAllowance(owner, token common.Address) (*big.Int, error) {
	tokenInterface, err := contractsinterfaces.NewToken(token, e.Client)
	if err != nil {
		log.Print(err)
		return nil, err
	}

	exchange := e.Config.ExchangeAddress()
	opts := &bind.CallOpts{Pending: true}
	a, err := tokenInterface.Allowance(opts, owner, exchange)
	if err != nil {
		log.Print(err)
		return nil, err
	}

	return a, nil
}

// func NewEthereumWebSocketConnection(config app.Config) *Ethereum {
// 	conn, err := rpc.DialWebsocket(context.Background(), config.EthereumURL)
// 	if err != nil {
// 		panic(err)
// 	}

// 	client = ethclient.NewClient(conn)
// 	config := NewEthereumConfig(config.EthereumURL, config.ExchangeAddress, config.WethAddress)

// 	return &Ethereum{
// 		Client: client,
// 		Config: config
// 	}
// }
