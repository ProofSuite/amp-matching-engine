package contracts

import (
	"context"
	"encoding/hex"
	"fmt"
	"log"
	"math/big"

	"github.com/Proofsuite/amp-matching-engine/contracts/contractsinterfaces"
	"github.com/Proofsuite/amp-matching-engine/interfaces"
	"github.com/Proofsuite/amp-matching-engine/types"
	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	eth "github.com/ethereum/go-ethereum/core/types"
)

type ethereumClientInterface interface {
	CodeAt(ctx context.Context, contract common.Address, blockNumber *big.Int) ([]byte, error)
	CallContract(ctx context.Context, call ethereum.CallMsg, blockNumber *big.Int) ([]byte, error)
	PendingCodeAt(ctx context.Context, account common.Address) ([]byte, error)
	PendingNonceAt(ctx context.Context, account common.Address) (uint64, error)
	SuggestGasPrice(ctx context.Context) (*big.Int, error)
	EstimateGas(ctx context.Context, call ethereum.CallMsg) (gas uint64, err error)
	SendTransaction(ctx context.Context, tx *eth.Transaction) error
	FilterLogs(ctx context.Context, query ethereum.FilterQuery) ([]eth.Log, error)
	SubscribeFilterLogs(ctx context.Context, query ethereum.FilterQuery, ch chan<- eth.Log) (ethereum.Subscription, error)
}

// Exchange is an augmented interface to the Exchange.sol smart-contract. It uses the
// smart-contract bindings generated with abigen and adds additional functionality and
// simplifications to these bindings.
// Address is the Ethereum address of the exchange contract
// Contract is the original abigen bindings
// CallOptions are options for making read calls to the connected backend
// TxOptions are options for making write txs to the connected backend
type Exchange struct {
	WalletService interfaces.WalletService
	TxService     interfaces.TxService
	Interface     *contractsinterfaces.Exchange
	Client        ethereumClientInterface
}

// Returns a new exchange interface for a given wallet, contract address and connected backend.
// The exchange contract need to be already deployed at the given address. The given wallet will
// be used by default when sending transactions with this object.
func NewExchange(
	w interfaces.WalletService,
	tx interfaces.TxService,
	contractAddress common.Address,
	backend ethereumClientInterface,
) (*Exchange, error) {
	instance, err := contractsinterfaces.NewExchange(contractAddress, backend)
	if err != nil {
		return nil, err
	}

	return &Exchange{
		WalletService: w,
		TxService:     tx,
		Interface:     instance,
		Client:        backend,
	}, nil
}

func (e *Exchange) SetTxSender(w *types.Wallet) {
	e.TxService.SetTxSender(w)
}

func (e *Exchange) GetTxCallOptions() *bind.CallOpts {
	return e.TxService.GetTxCallOptions()
}

func (e *Exchange) GetTxSendOptions() (*bind.TransactOpts, error) {
	return e.TxService.GetTxSendOptions()
}

func (e *Exchange) GetCustomTxSendOptions(w *types.Wallet) *bind.TransactOpts {
	return e.TxService.GetCustomTxSendOptions(w)
}

// SetFeeAccount sets the fee account of the exchange contract. The fee account receives
// the trading fees whenever a trade is settled.
func (e *Exchange) SetFeeAccount(a common.Address) (*eth.Transaction, error) {
	txOptions, _ := e.GetTxSendOptions()

	tx, err := e.Interface.SetFeeAccount(txOptions, a)
	if err != nil {
		return nil, err
	}

	return tx, nil
}

// SetOperator updates the operator settings of the given address. Only addresses with an
// operator access can execute Withdraw and Trade transactions to the Exchange smart contract
func (e *Exchange) SetOperator(a common.Address, isOperator bool) (*eth.Transaction, error) {
	txOptions, _ := e.GetTxSendOptions()

	tx, err := e.Interface.SetOperator(txOptions, a, isOperator)
	if err != nil {
		return nil, err
	}

	return tx, nil
}

// FeeAccount is the Ethereum towards the exchange trading fees are sent
func (e *Exchange) FeeAccount() (common.Address, error) {
	callOptions := e.GetTxCallOptions()

	account, err := e.Interface.FeeAccount(callOptions)
	if err != nil {
		return common.Address{}, err
	}

	return account, nil
}

func (e *Exchange) Operator(a common.Address) (bool, error) {
	callOptions := e.GetTxCallOptions()
	// Operator returns true if the given address is an operator of the exchange and returns false otherwise
	isOperator, err := e.Interface.Operators(callOptions, a)
	if err != nil {
		return false, err
	}

	return isOperator, nil
}

// Trade executes a settlements transaction. The order and trade payloads need to be signed respectively
// by the Maker and the Taker of the trade. Only the operator account can send a Trade function to the
// Exchange smart contract.
func (e *Exchange) Trade(o *types.Order, t *types.Trade, txOpts *bind.TransactOpts) (*eth.Transaction, error) {
	orderValues := [8]*big.Int{o.BuyAmount, o.SellAmount, o.Expires, o.Nonce, o.MakeFee, o.TakeFee, t.Amount, t.TradeNonce}
	orderAddresses := [4]common.Address{o.BuyToken, o.SellToken, o.UserAddress, t.Taker}
	vValues := [2]uint8{o.Signature.V, t.Signature.V}
	rsValues := [4][32]byte{o.Signature.R, o.Signature.S, t.Signature.R, t.Signature.S}

	tx, err := e.Interface.ExecuteTrade(txOpts, orderValues, orderAddresses, vValues, rsValues)
	if err != nil {
		return nil, err
	}

	return tx, nil
}

// ListenToErrorEvents returns a channel that receives errors logs (events) from the exchange smart contract.
// The error IDs correspond to the following codes:
// 1. MAKER_INSUFFICIENT_BALANCE,
// 2. TAKER_INSUFFICIENT_BALANCE,
// 3. WITHDRAW_INSUFFICIENT_BALANCE,
// 4. WITHDRAW_FEE_TO_HIGH,
// 5. ORDER_EXPIRED,
// 6. WITHDRAW_ALREADY_COMPLETED,
// 7. TRADE_ALREADY_COMPLETED,
// 8. TRADE_AMOUNT_TOO_BIG,
// 9. SIGNATURE_INVALID,
// 10. MAKER_SIGNATURE_INVALID,
// 11. TAKER_SIGNATURE_INVALID
func (e *Exchange) ListenToErrors() (chan *contractsinterfaces.ExchangeLogError, error) {
	events := make(chan *contractsinterfaces.ExchangeLogError)
	opts := &bind.WatchOpts{nil, nil}

	_, err := e.Interface.WatchLogError(opts, events)
	if err != nil {
		return nil, err
	}

	return events, nil
}

// ListenToTrades returns a channel that receivs trade logs (events) from the underlying exchange smart contract
func (e *Exchange) ListenToTrades() (chan *contractsinterfaces.ExchangeLogTrade, error) {
	events := make(chan *contractsinterfaces.ExchangeLogTrade)
	opts := &bind.WatchOpts{nil, nil}

	_, err := e.Interface.WatchLogTrade(opts, events, nil, nil, nil)
	if err != nil {
		return nil, err
	}

	return events, nil
}

func (e *Exchange) GetErrorEvents(logs chan *contractsinterfaces.ExchangeLogError) error {
	opts := &bind.WatchOpts{nil, nil}

	_, err := e.Interface.WatchLogError(opts, logs)
	if err != nil {
		return err
	}

	return nil
}

func (e *Exchange) GetTrades(logs chan *contractsinterfaces.ExchangeLogTrade) error {
	opts := &bind.WatchOpts{nil, nil}

	_, err := e.Interface.WatchLogTrade(opts, logs, nil, nil, nil)
	if err != nil {
		return err
	}

	return nil
}

func (e *Exchange) PrintTrades() error {
	events := make(chan *contractsinterfaces.ExchangeLogTrade)
	opts := &bind.WatchOpts{nil, nil}

	_, err := e.Interface.WatchLogTrade(opts, events, nil, nil, nil)
	if err != nil {
		return err
	}

	go func() {
		for {
			event := <-events
			fmt.Printf("New event: %v", event)
		}
	}()

	return nil
}

func (e *Exchange) PrintErrors() error {
	events := make(chan *contractsinterfaces.ExchangeLogError)
	opts := &bind.WatchOpts{nil, nil}

	_, err := e.Interface.WatchLogError(opts, events)
	if err != nil {
		return err
	}

	go func() {
		for {
			event := <-events
			log.Printf("New Error Event. Id: %v, Hash: %v\n\n", event.ErrorId, hex.EncodeToString(event.TradeHash[:]))
		}
	}()

	return nil
}

func PrintErrorLog(log *contractsinterfaces.ExchangeLogError) string {
	return fmt.Sprintf("Error:\nErrorID: %v\nOrderHash: %v\n\n", log.ErrorId, log.OrderHash)
}

func PrintTradeLog(log *contractsinterfaces.ExchangeLogTrade) string {
	return fmt.Sprintf("Error:\nMaker: %v\nTaker: %v\nTokenBuy: %v\nTokenSell: %v\nOrderHash: %v\nTradeHash: %v\n\n",
		log.Maker, log.Taker, log.TokenBuy, log.TokenSell, log.OrderHash, log.TradeHash)
}

func PrintCancelOrderLog(log *contractsinterfaces.ExchangeLogCancelOrder) string {
	return fmt.Sprintf("Error:\nSender: %v\nOrderHash: %v\n\n", log.Maker, log.OrderHash)
}

func PrintCancelTradeLog(log *contractsinterfaces.ExchangeLogCancelTrade) string {
	return fmt.Sprintf("Error:\nSender: %v\nTradeHash: %v\n\n", log.Taker, log.OrderHash)
}
