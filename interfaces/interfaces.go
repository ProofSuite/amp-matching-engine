package interfaces

import (
	"context"
	"math/big"

	"github.com/Proofsuite/amp-matching-engine/contracts/contractsinterfaces"
	"github.com/Proofsuite/amp-matching-engine/rabbitmq"
	"github.com/Proofsuite/amp-matching-engine/types"
	"github.com/Proofsuite/amp-matching-engine/ws"
	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	eth "github.com/ethereum/go-ethereum/core/types"
	"gopkg.in/mgo.v2/bson"
)

type OrderDao interface {
	Create(o *types.Order) error
	Update(id bson.ObjectId, o *types.Order) error
	UpdateAllByHash(hash common.Hash, o *types.Order) error
	UpdateByHash(hash common.Hash, o *types.Order) error
	GetByID(id bson.ObjectId) (*types.Order, error)
	GetByHash(hash common.Hash) (*types.Order, error)
	GetByHashes(hashes []common.Hash) ([]*types.Order, error)
	GetByUserAddress(addr common.Address) ([]*types.Order, error)
	GetCurrentByUserAddress(addr common.Address) ([]*types.Order, error)
	GetHistoryByUserAddress(addr common.Address) ([]*types.Order, error)
	UpdateOrderFilledAmount(hash common.Hash, value *big.Int) error
	GetUserLockedBalance(account common.Address, token common.Address) (*big.Int, error)
	UpdateOrderStatus(hash common.Hash, status string) error
	GetRawOrderBook(*types.Pair) ([]*types.Order, error)
	GetOrderBook(*types.Pair) ([]map[string]string, []map[string]string, error)
	GetOrderBookPricePoint(p *types.Pair, pp *big.Int) (*big.Int, error)
	Drop() error
}

type AccountDao interface {
	Create(account *types.Account) (err error)
	GetAll() (res []types.Account, err error)
	GetByID(id bson.ObjectId) (*types.Account, error)
	GetByAddress(owner common.Address) (response *types.Account, err error)
	GetTokenBalances(owner common.Address) (map[common.Address]*types.TokenBalance, error)
	GetTokenBalance(owner common.Address, token common.Address) (*types.TokenBalance, error)
	UpdateTokenBalance(owner common.Address, token common.Address, tokenBalance *types.TokenBalance) (err error)
	UpdateBalance(owner common.Address, token common.Address, balance *big.Int) (err error)
	UpdateAllowance(owner common.Address, token common.Address, allowance *big.Int) (err error)
	Drop()
}

type WalletDao interface {
	Create(wallet *types.Wallet) error
	GetAll() ([]types.Wallet, error)
	GetByID(id bson.ObjectId) (*types.Wallet, error)
	GetByAddress(addr common.Address) (*types.Wallet, error)
	GetDefaultAdminWallet() (*types.Wallet, error)
	GetOperatorWallets() ([]*types.Wallet, error)
}

type PairDao interface {
	Create(o *types.Pair) error
	GetAll() ([]types.Pair, error)
	GetByID(id bson.ObjectId) (*types.Pair, error)
	GetByName(name string) (*types.Pair, error)
	GetByTokenSymbols(baseTokenSymbol, quoteTokenSymbol string) (*types.Pair, error)
	GetByTokenAddress(baseToken, quoteToken common.Address) (*types.Pair, error)
	GetByBuySellTokenAddress(buyToken, sellToken common.Address) (*types.Pair, error)
}

type TradeDao interface {
	Create(o ...*types.Trade) error
	Update(t *types.Trade) error
	UpdateByHash(hash common.Hash, t *types.Trade) error
	GetAll() ([]types.Trade, error)
	Aggregate(q []bson.M) ([]*types.Tick, error)
	GetByPairName(name string) ([]*types.Trade, error)
	GetByHash(hash common.Hash) (*types.Trade, error)
	GetByOrderHash(hash common.Hash) ([]*types.Trade, error)
	GetByPairAddress(baseToken, quoteToken common.Address) ([]*types.Trade, error)
	GetByUserAddress(addr common.Address) ([]*types.Trade, error)
	UpdateTradeStatus(hash common.Hash, status string) error
	Drop()
}

type TokenDao interface {
	Create(token *types.Token) error
	GetAll() ([]types.Token, error)
	GetByID(id bson.ObjectId) (*types.Token, error)
	GetByAddress(owner common.Address) (*types.Token, error)
	GetQuoteTokens() ([]types.Token, error)
	GetBaseTokens() ([]types.Token, error)
	Drop() error
}

type Exchange interface {
	GetAddress() common.Address
	GetTxCallOptions() *bind.CallOpts
	SetFeeAccount(a common.Address, txOpts *bind.TransactOpts) (*eth.Transaction, error)
	SetOperator(a common.Address, isOperator bool, txOpts *bind.TransactOpts) (*eth.Transaction, error)
	CallTrade(o *types.Order, t *types.Trade, call *ethereum.CallMsg) (uint64, error)
	FeeAccount() (common.Address, error)
	Operator(a common.Address) (bool, error)
	Trade(o *types.Order, t *types.Trade, txOpts *bind.TransactOpts) (*eth.Transaction, error)
	ListenToErrors() (chan *contractsinterfaces.ExchangeLogError, error)
	ListenToTrades() (chan *contractsinterfaces.ExchangeLogTrade, error)
	GetErrorEvents(logs chan *contractsinterfaces.ExchangeLogError) error
	GetTrades(logs chan *contractsinterfaces.ExchangeLogTrade) error
	PrintTrades() error
	PrintErrors() error
}

type Engine interface {
	HandleOrders(msg *rabbitmq.Message) error
	RecoverOrders(orders []*types.OrderTradePair) error
	CancelOrder(order *types.Order) (*types.EngineResponse, error)
	CancelTrades(orders []*types.Order, amount []*big.Int) error
	DeleteOrder(o *types.Order) error
	DeleteOrders(orders ...types.Order) error
}

type WalletService interface {
	CreateAdminWallet(a common.Address) (*types.Wallet, error)
	GetDefaultAdminWallet() (*types.Wallet, error)
	GetOperatorWallets() ([]*types.Wallet, error)
	GetAll() ([]types.Wallet, error)
	GetByAddress(a common.Address) (*types.Wallet, error)
}

type OHLCVService interface {
	Unsubscribe(conn *ws.Conn, bt, qt common.Address, p *types.Params)
	Subscribe(conn *ws.Conn, bt, qt common.Address, p *types.Params)
	GetOHLCV(p []types.PairAddresses, duration int64, unit string, timeInterval ...int64) ([]*types.Tick, error)
}

type EthereumService interface {
	WaitMined(hash common.Hash) (*eth.Receipt, error)
	GetBalanceAt(a common.Address) (*big.Int, error)
	GetPendingNonceAt(a common.Address) (uint64, error)
}

type OrderService interface {
	GetByID(id bson.ObjectId) (*types.Order, error)
	GetByHash(hash common.Hash) (*types.Order, error)
	GetByUserAddress(addr common.Address) ([]*types.Order, error)
	NewOrder(o *types.Order) error
	CancelOrder(oc *types.OrderCancel) error
	CancelTrades(trades []*types.Trade) error
	HandleEngineResponse(res *types.EngineResponse) error
	GetCurrentByUserAddress(addr common.Address) ([]*types.Order, error)
	GetHistoryByUserAddress(addr common.Address) ([]*types.Order, error)
	Rollback(res *types.EngineResponse) *types.EngineResponse
	RollbackOrder(o *types.Order) error
	RollbackTrade(o *types.Order, t *types.Trade) error
}

type OrderBookService interface {
	GetOrderBook(bt, qt common.Address) (map[string]interface{}, error)
	GetRawOrderBook(bt, qt common.Address) ([]*types.Order, error)
	SubscribeOrderBook(conn *ws.Conn, bt, qt common.Address)
	UnSubscribeOrderBook(conn *ws.Conn, bt, qt common.Address)
	SubscribeRawOrderBook(conn *ws.Conn, bt, qt common.Address)
	UnSubscribeRawOrderBook(conn *ws.Conn, bt, qt common.Address)
}

type PairService interface {
	Create(pair *types.Pair) error
	GetByID(id bson.ObjectId) (*types.Pair, error)
	GetByTokenAddress(bt, qt common.Address) (*types.Pair, error)
	GetAll() ([]types.Pair, error)
}

type TokenService interface {
	Create(token *types.Token) error
	GetByID(id bson.ObjectId) (*types.Token, error)
	GetByAddress(addr common.Address) (*types.Token, error)
	GetAll() ([]types.Token, error)
	GetQuoteTokens() ([]types.Token, error)
	GetBaseTokens() ([]types.Token, error)
}

type TradeService interface {
	GetByPairName(p string) ([]*types.Trade, error)
	GetTrades(bt, qt common.Address) ([]types.Trade, error)
	GetByPairAddress(bt, qt common.Address) ([]*types.Trade, error)
	GetByUserAddress(addr common.Address) ([]*types.Trade, error)
	GetByHash(hash common.Hash) (*types.Trade, error)
	GetByOrderHash(hash common.Hash) ([]*types.Trade, error)
	UpdateTradeTxHash(tr *types.Trade, txHash common.Hash) error
	Subscribe(conn *ws.Conn, bt, qt common.Address)
	Unsubscribe(conn *ws.Conn, bt, qt common.Address)
}

type TxService interface {
	GetTxCallOptions() *bind.CallOpts
	GetTxSendOptions() (*bind.TransactOpts, error)
	GetTxDefaultSendOptions() (*bind.TransactOpts, error)
	SetTxSender(w *types.Wallet)
	GetCustomTxSendOptions(w *types.Wallet) *bind.TransactOpts
}

type AccountService interface {
	Create(account *types.Account) error
	GetByID(id bson.ObjectId) (*types.Account, error)
	GetAll() ([]types.Account, error)
	GetByAddress(a common.Address) (*types.Account, error)
	GetTokenBalance(owner common.Address, token common.Address) (*types.TokenBalance, error)
	GetTokenBalances(owner common.Address) (map[common.Address]*types.TokenBalance, error)
}

type EthereumConfig interface {
	GetURL() string
	ExchangeAddress() common.Address
	WethAddress() common.Address
}

type EthereumClient interface {
	CodeAt(ctx context.Context, contract common.Address, blockNumber *big.Int) ([]byte, error)
	CallContract(ctx context.Context, call ethereum.CallMsg, blockNumber *big.Int) ([]byte, error)
	PendingCodeAt(ctx context.Context, account common.Address) ([]byte, error)
	TransactionReceipt(ctx context.Context, txHash common.Hash) (*eth.Receipt, error)
	EstimateGas(ctx context.Context, call ethereum.CallMsg) (gas uint64, err error)
	SendTransaction(ctx context.Context, tx *eth.Transaction) error
	PendingNonceAt(ctx context.Context, account common.Address) (uint64, error)
	BalanceAt(ctx context.Context, contract common.Address, blockNumber *big.Int) (*big.Int, error)
	FilterLogs(ctx context.Context, query ethereum.FilterQuery) ([]eth.Log, error)
	SubscribeFilterLogs(ctx context.Context, query ethereum.FilterQuery, ch chan<- eth.Log) (ethereum.Subscription, error)
	SuggestGasPrice(ctx context.Context) (*big.Int, error)
}

type EthereumProvider interface {
	WaitMined(hash common.Hash) (*eth.Receipt, error)
	GetBalanceAt(a common.Address) (*big.Int, error)
	GetPendingNonceAt(a common.Address) (uint64, error)
	BalanceOf(owner common.Address, token common.Address) (*big.Int, error)
	Allowance(owner, spender, token common.Address) (*big.Int, error)
	ExchangeAllowance(owner, token common.Address) (*big.Int, error)
}
