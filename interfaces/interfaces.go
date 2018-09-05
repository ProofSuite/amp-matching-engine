package interfaces

import (
	"math/big"

	"github.com/Proofsuite/amp-matching-engine/contracts/contractsinterfaces"
	"github.com/Proofsuite/amp-matching-engine/rabbitmq"
	"github.com/Proofsuite/amp-matching-engine/types"
	"github.com/Proofsuite/amp-matching-engine/ws"
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
	Drop() error
}

type AccountDao interface {
	Create(account *types.Account) (err error)
	GetAll() (res []types.Account, err error)
	GetByID(id bson.ObjectId) (*types.Account, error)
	GetByAddress(owner common.Address) (response *types.Account, err error)
	GetTokenBalances(owner common.Address) (map[common.Address]*types.TokenBalance, error)
	GetWethTokenBalance(owner common.Address) (*types.TokenBalance, error)
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
	Drop()
}

type TokenDao interface {
	Create(token *types.Token) error
	GetAll() ([]types.Token, error)
	GetByID(id bson.ObjectId) (*types.Token, error)
	GetByAddress(owner common.Address) (*types.Token, error)
	Drop() error
}

type Exchange interface {
	SetTxSender(w *types.Wallet)
	GetTxCallOptions() *bind.CallOpts
	GetTxSendOptions() (*bind.TransactOpts, error)
	GetCustomTxSendOptions(w *types.Wallet) *bind.TransactOpts
	SetFeeAccount(a common.Address) (*eth.Transaction, error)
	SetOperator(a common.Address, isOperator bool) (*eth.Transaction, error)
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
	SubscribeResponseQueue(fn func(*types.EngineResponse) error) error
	RecoverOrders(orders []*types.FillOrder) error
	CancelOrder(order *types.Order) (*types.EngineResponse, error)
	GetOrderBook(pair *types.Pair) (asks, bids []*map[string]float64)
	CancelTrades(orders []*types.Order, amount []*big.Int) error
	GetFullOrderBook(pair *types.Pair) [][]types.Order
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
	GetOHLCV(p []types.PairSubDoc, duration int64, unit string, timeInterval ...int64) ([]*types.Tick, error)
}

type EthereumService interface {
	WaitMined(tx *eth.Transaction) (*eth.Receipt, error)
	GetPendingNonceAt(a common.Address) (uint64, error)
	// GetPendingBalanceAt(a common.Address) (*big.Int, error)
}

type OrderService interface {
	GetByID(id bson.ObjectId) (*types.Order, error)
	GetByHash(hash common.Hash) (*types.Order, error)
	GetByUserAddress(addr common.Address) ([]*types.Order, error)
	NewOrder(o *types.Order) error
	CancelOrder(oc *types.OrderCancel) error
	HandleEngineResponse(res *types.EngineResponse) error
	RecoverOrders(res *types.EngineResponse)
	RelayUpdateOverSocket(res *types.EngineResponse)
	SendMessage(msgType string, hash common.Hash, data interface{})
	SubscribeQueue(fn func(*rabbitmq.Message) error) error
	PublishOrder(order *rabbitmq.Message) error
}

type OrderBookService interface {
	GetOrderBook(bt, qt common.Address) (ob map[string]interface{}, err error)
	GetFullOrderBook(bt, qt common.Address) (ob [][]types.Order, err error)
	SubscribeLite(conn *ws.Conn, bt, qt common.Address)
	UnsubscribeLite(conn *ws.Conn, bt, qt common.Address)
	SubscribeFull(conn *ws.Conn, bt, qt common.Address)
	UnsubscribeFull(conn *ws.Conn, bt, qt common.Address)
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
}

type TradeService interface {
	GetByPairName(p string) ([]*types.Trade, error)
	GetTrades(bt, qt common.Address) ([]types.Trade, error)
	GetByPairAddress(bt, qt common.Address) ([]*types.Trade, error)
	GetByUserAddress(addr common.Address) ([]*types.Trade, error)
	GetByHash(hash common.Hash) (*types.Trade, error)
	GetByOrderHash(hash common.Hash) ([]*types.Trade, error)
	UpdateTradeTx(tr *types.Trade, tx *eth.Transaction) error
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
