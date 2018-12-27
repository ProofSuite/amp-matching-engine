package ethereum

import (
	"context"
	"math/big"
	"time"

	"github.com/Proofsuite/amp-matching-engine/app"
	"github.com/Proofsuite/amp-matching-engine/contracts/contractsinterfaces"
	"github.com/Proofsuite/amp-matching-engine/interfaces"
	"github.com/Proofsuite/amp-matching-engine/utils"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	eth "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
)

type EthereumProvider struct {
	Client interfaces.EthereumClient
	Config interfaces.EthereumConfig
}

func NewEthereumProvider(c interfaces.EthereumClient) *EthereumProvider {
	url := app.Config.Ethereum["http_url"]
	exchange := common.HexToAddress(app.Config.Ethereum["exchange_address"])
	config := NewEthereumConfig(url, exchange)

	return &EthereumProvider{
		Client: c,
		Config: config,
	}
}

func NewDefaultEthereumProvider() *EthereumProvider {
	url := app.Config.Ethereum["http_url"]
	exchange := common.HexToAddress(app.Config.Ethereum["exchange_address"])

	conn, err := rpc.DialHTTP(app.Config.Ethereum["http_url"])
	if err != nil {
		panic(err)
	}

	client := ethclient.NewClient(conn)
	config := NewEthereumConfig(url, exchange)

	return &EthereumProvider{
		Client: client,
		Config: config,
	}
}

func NewWebsocketProvider() *EthereumProvider {
	url := app.Config.Ethereum["ws_url"]
	exchange := common.HexToAddress(app.Config.Ethereum["exchange_address"])

	conn, err := rpc.DialWebsocket(context.Background(), url, "")
	if err != nil {
		panic(err)
	}

	client := ethclient.NewClient(conn)
	config := NewEthereumConfig(url, exchange)

	return &EthereumProvider{
		Client: client,
		Config: config,
	}
}

func NewSimulatedEthereumProvider(accs []common.Address) *EthereumProvider {
	url := app.Config.Ethereum["http_url"]
	exchange := common.HexToAddress(app.Config.Ethereum["exchange_address"])

	config := NewEthereumConfig(url, exchange)
	client := NewSimulatedClient(accs)

	return &EthereumProvider{
		Client: client,
		Config: config,
	}
}

func (e *EthereumProvider) WaitMined(hash common.Hash) (*eth.Receipt, error) {
	ctx := context.Background()
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		receipt, _ := e.Client.TransactionReceipt(ctx, hash)
		if receipt != nil {
			return receipt, nil
		}

		// if err != nil {
		// 	logger.Error(err)
		// 	// return nil, err
		// }

		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-ticker.C:
		}
	}
}

func (e *EthereumProvider) GetBalanceAt(a common.Address) (*big.Int, error) {
	ctx := context.Background()
	nonce, err := e.Client.BalanceAt(ctx, a, nil)
	if err != nil {
		logger.Error(err)
		return big.NewInt(0), err
	}

	return nonce, nil
}

func (e *EthereumProvider) GetPendingNonceAt(a common.Address) (uint64, error) {
	ctx := context.Background()
	nonce, err := e.Client.PendingNonceAt(ctx, a)
	if err != nil {
		logger.Error(err)
		return 0, err
	}

	return nonce, nil
}

func (e *EthereumProvider) Decimals(token common.Address) (uint8, error) {
	var tokenInterface *contractsinterfaces.ERC20
	var err error

	// retry in case the connection with the ethereum client is asleep
	err = utils.Retry(3, func() error {
		tokenInterface, err = contractsinterfaces.NewERC20(token, e.Client)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		logger.Error(err)
		return 0, err
	}

	opts := &bind.CallOpts{Pending: true}
	decimals, err := tokenInterface.Decimals(opts)
	if err != nil {
		logger.Error(err)
		return 0, err
	}

	return decimals, nil
}

func (e *EthereumProvider) Symbol(token common.Address) (string, error) {
	// retry in case the connection with the ethereum client is asleep
	var tokenInterface *contractsinterfaces.ERC20
	var err error

	err = utils.Retry(3, func() error {
		tokenInterface, err = contractsinterfaces.NewERC20(token, e.Client)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		logger.Error(err)
		return "", err
	}

	opts := &bind.CallOpts{Pending: true}
	symbol, err := tokenInterface.Symbol(opts)
	if err != nil {
		logger.Error(err)
		return "", err
	}

	return symbol, nil
}

func (e *EthereumProvider) BalanceOf(owner common.Address, token common.Address) (*big.Int, error) {
	tokenInterface, err := contractsinterfaces.NewERC20(token, e.Client)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	opts := &bind.CallOpts{Pending: true}
	b, err := tokenInterface.BalanceOf(opts, owner)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	return b, nil
}

func (e *EthereumProvider) Allowance(owner, spender, token common.Address) (*big.Int, error) {
	tokenInterface, err := contractsinterfaces.NewERC20(token, e.Client)
	if err != nil {
		return nil, err
	}

	opts := &bind.CallOpts{Pending: true}
	a, err := tokenInterface.Allowance(opts, owner, spender)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	return a, nil
}

func (e *EthereumProvider) ExchangeAllowance(owner, token common.Address) (*big.Int, error) {
	tokenInterface, err := contractsinterfaces.NewERC20(token, e.Client)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	exchange := common.HexToAddress(app.Config.Ethereum["exchange_address"])
	opts := &bind.CallOpts{Pending: true}
	a, err := tokenInterface.Allowance(opts, owner, exchange)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	return a, nil
}

// func (e *EthereumProvider) NewERC20Instance(
// 	w interfaces.WalletService,
// 	tx interfaces.TxService,
// 	token common.Address,
// ) (*contractsinterfaces.Token, error) {
// 	tokenInterface, err := contractsinterfaces.NewToken(token, e.Client)
// 	if err != nil {
// 		logger.Error(err)
// 		return nil, err
// 	}

// 	return &contracts.Token{
// 		WalletService: w,
// 		TxService:     tx,
// 		Interface:     tokenInterface,
// 	}, nil
// }

// func (e *EthereumProvider) NewExchangeInstance(w interfaces.WalletService, tx interfaces.TxService) (*contracts.Exchange, error) {
// 	exchangeAddress := app.Config.Ethereum["exchange_address"]
// 	if exchangeAddress == "" {
// 		return nil, errors.New("Exchange address configuration not found")
// 	}

// 	exchangeInterface, err := contractsinterfaces.NewExchange(exchangeAddress, e.Client)
// 	if err != nil {
// 		logger.Error(err)
// 		return nil, err
// 	}

// 	return &contracts.Exchange{
// 		WalletService: w,
// 		TxService:     tx,
// 		Interface:     exchangeInterface,
// 		Client:        e.Client,
// 	}, nil
// }

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
