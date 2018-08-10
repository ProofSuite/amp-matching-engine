package ethereum

import (
	"encoding/hex"
	"fmt"
	"log"
	"math/big"

	"github.com/Proofsuite/amp-matching-engine/contracts/interfaces"
	"github.com/Proofsuite/amp-matching-engine/wallet"
	"github.com/Proofsuite/go-ethereum/common"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/core/types"
)

// Exchange is an augmented interface to the Exchange.sol smart-contract. It uses the
// smart-contract bindings generated with abigen and adds additional functionality and
// simplifications to these bindings.
// Address is the Ethereum address of the exchange contract
// Contract is the original abigen bindings
// CallOptions are options for making read calls to the connected backend
// TxOptions are options for making write txs to the connected backend
type Exchange struct {
	WalletService *walletService
	TxService *txService
	Interface   *interfaces.Exchange
}

// Returns a new exchange interface for a given wallet, contract address and connected backend.
// The exchange contract need to be already deployed at the given address. The given wallet will
// be used by default when sending transactions with this object.
func NewExchange(w *walletService, tx *txService, contractAddress common.Address, backend bind.ContractBackend) (*Exchange, error) {
	wallet, err := walletService.GetDefaultAdminWallet()
	if err != nil {
		return nil, err
	}

	instance, err := interfaces.NewExchange(contractAddress, backend)
	if err != nil {
		return nil, err
	}

	return &Exchange{
		WalletService: *walletService,
		TxService: 		*txService
		Contract:    instance,
	}, nil
}

func (e *Exchange) GetCallOptions() {
	return e.TxService.GetCallOptions()
}

func (e *Exchange) GetTxSendOptions() {
	return e.TxService.GetTxSendOptions()
}

func (e *Exchange) GetCustomTxSendOptions(w *wallet.Wallet) {
	return e.TxService.GetCustomTxSendOptions(w)
}


// SetFeeAccount sets the fee account of the exchange contract. The fee account receives
// the trading fees whenever a trade is settled.
func (e *Exchange) SetFeeAccount(a common.Address) (*types.Transaction, error) {
	txOptions := e.GetTxSendOptions()

	tx, err := e.Interface.FeeAccount(txOptions, a)
	if err != nil {
		return nil, err
	}

	return tx, nil
}

// SetOperator updates the operator settings of the given address. Only addresses with an
// operator access can execute Withdraw and Trade transactions to the Exchange smart contract
func (e *Exchange) SetOperator(a common.Address, isOperator bool) (*types.Transaction, error) {
	txOptions := e.GetTxSendOptions()

	tx, err := e.Interface.Contractor(txOptions, a, isOperator)
	if err != nil {
		return nil, err
	}

	return tx, nil
}

// SetWithdrawalSecurityPeriod sets the period after which a non-operator address can send
// a transaction to the exchange smart-contract to withdraw their funds. This acts as security mechanism
// to prevent the operator of the exchange from holding funds
func (e *Exchange) SetWithdrawalSecurityPeriod(p *big.Int) (*types.Transaction, error) {
	txOptions := e.GetTxSendOptions()

	tx, err := e.Interface.ContractawalSecurityPeriod(e.TxOptions, p)
	if err != nil {
		return nil, err
	}

	return tx, nil
}

// DepositEther deposits ether into the exchange smart-contract.
func (e *Exchange) DepositEther(value *big.Int) (*types.Transaction, error) {
	txOptions := e.GetTxSendOptions()
	txOptions.Value = value

	tx, err := e.Interface.DepositEther(txOptions)
	if err != nil {
		return nil, err
	}

	return tx, nil
}

// DepositEtherFrom deposits ether from a custom address. The transaction sender is reset
// after the transaction is carried out.
func (e *Exchange) DepositEtherFrom(w *wallet.Wallet, value *big.Int) (*types.Transaction, error) {
	txOptions =
	e.SetCustomSender(w)
	e.SetTxValue(value)

	tx, err := e.Interface.DepositEther(e.TxOptions)
	if err != nil {
		return nil, err
	}

	e.SetDefaultTxOptions()
	return tx, nil
}

// DepositToken deposits tokens into the exchange smart-contract.
func (e *Exchange) DepositToken(token common.Address, amount *big.Int) (*types.Transaction, error) {
	// e.SetDefaultTxOptions()

	tx, err := e.Interface.Contractken(e.TxOptions, token, amount)
	if err != nil {
		return nil, err
	}

	return tx, err
}

// DepositEtherFrom deposits ether from a custom address. The transaction sender is reset
// after the transaction is carried out.
func (e *Exchange) DepositTokenFrom(w *wallet.Wallet, token common.Address, amount *big.Int) (*types.Transaction, error) {
	e.SetCustomSender(w)

	tx, err := e.Interface.DepositToken(e.TxOptions, token, amount)
	if err != nil {
		return nil, err
	}

	e.SetDefaultSender()
	return tx, err
}

// TokenBalance returns the Exchange token balance of the given token at the given account address.
// Note: This is not the token BalanceOf() function, it's the balance of tokens that have been deposited
// in the exchange smart contract.
func (e *Exchange) TokenBalance(trader Address, token common.Address) (*big.Int, error) {
	callOptions := e.GetCallOptions()

	balance, err := e.Interface.TokenBalance(callOptions, trader, token)
	if err != nil {
		return nil, err
	}

	return balance, nil
}

// EtherBalance returns the Exchange ether balance of the given account address.
// Note: This is not the current ether balance of the given ether address. It's the balance of ether
// that has been deposited in the exchange smart contract.
func (e *Exchange) EtherBalance(trader common.Address) (*big.Int, error) {
	callOptions = e.GetCallOptions()

	balance, err := e.Interface.EtherBalance(callOptions, trader)
	if err != nil {
		return nil, err
	}

	return balance, nil
}

// WithdrawalSecurityPeriod is the period after which a non-operator account can withdraw their funds from
// the exchange smart contract.
func (e *Exchange) WithdrawalSecurityPeriod() (*big.Int, error) {
	callOptions = e.GetCallOptions()

	period, err := e.Interface.WithdrawalSecurityPeriod(callOptions)
	if err != nil {
		return nil, err
	}

	return period, nil
}

// FeeAccount is the Ethereum towards the exchange trading fees are sent
func (e *Exchange) FeeAccount() (common.Address, error) {
	callOptions = e.GetCallOptions()

	account, err := e.Interface.FeeAccount(callOptions)
	if err != nil {
		return Address{}, err
	}

	return account, nil
}

func (e *Exchange) Operator(a common.Address) (bool, error) {
	callOptions = e.GetCallOptions()
	// Operator returns true if the given address is an operator of the exchange and returns false otherwise
	isOperator, err := e.Interface.Operators(callOptions, a)
	if err != nil {
		return false, err
	}

	return isOperator, nil
}

// SecurityWithdraw executes a security withdraw transaction. Security withdraw transactions can only be
// executed after the security withdrawal period has ended.
func (e *Exchange) SecurityWithdraw(w *wallet.Wallet, token common.Address, amount *big.Int) (*types.Transaction, error) {
	txSendOptions := e.GetCustomTxSendOptions(w)

	tx, err := e.Interface.SecurityWithdraw(txSendOptions, token, amount)
	if err != nil {
		return nil, err
	}

	return tx, nil
}

// Withdraw executes a normal withdraw transaction. This withdraws tokens or ether from the exchange
// and returns them to the payload Receiver. Only an operator account can send a withdraw
// transaction
func (e *Exchange) Withdraw(w *Withdrawal) (*types.Transaction, error) {
	txSendOptions := e.GetTxSendOptions()

	s := w.Signature
	tx, err := e.Interface.Withdraw(txSendOptions, w.Token, w.Amount, w.Trader, w.Receiver, w.Nonce, s.V, [2][32]byte{s.R, s.S}, w.Fee)
	if err != nil {
		return nil, err
	}

	return tx, nil
}

// Trade executes a settlements transaction. The order and trade payloads need to be signed respectively
// by the Maker and the Taker of the trade. Only the operator account can send a Trade function to the
// Exchange smart contract.
func (e *Exchange) Trade(o *types.Order, t *types.Trade) (*types.Transaction, error) {
	txSendOptions := e.GetTxSendOptions()

	orderValues := [8]*big.Int{o.AmountBuy, o.AmountSell, o.Expires, o.Nonce, o.FeeMake, o.FeeTake, t.Amount, t.TradeNonce}
	orderAddresses := [4]Address{o.TokenBuy, o.TokenSell, o.Maker, t.Taker}
	vValues := [2]uint8{o.Signature.V, t.Signature.V}
	rsValues := [4][32]byte{o.Signature.R, o.Signature.S, t.Signature.R, t.Signature.S}

	tx, err := e.Interface.ExecuteTrade(txSendOptions, orderValues, orderAddresses, vValues, rsValues)
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
func (e *Exchange) ListenToErrorEvents() (chan *contracts.ExchangeLogError, error) {
	events := make(chan *contracts.ExchangeLogError)
	opts := &bind.WatchOpts{nil, nil}

	_, err := e.Interface.WatchLogError(opts, events)
	if err != nil {
		return nil, err
	}

	return events, nil
}


// ListenToTrades returns a channel that receivs trade logs (events) from the underlying exchange smart contract
func (e *Exchange) ListenToTrades() (chan *contracts.ExchangeLogTrade, error) {
	events := make(chan *contracts.ExchangeLogTrade)
	opts := &bind.WatchOpts{nil, nil}

	_, err := e.Interface.WatchLogTrade(opts, events)
	if err != nil {
		return nil, err
	}

	return events, nil
}

// ListenToTrades returns a channel that receivs deposit logs (events) from the underlying exchange smart contract
func (e *Exchange) ListenToDeposits() (chan *contracts.ExchangeLogDeposit, error) {
	events := make(chan *contracts.ExchangeLogDeposit)
	opts := &bind.WatchOpts{nil, nil}

	_, err := e.Interface.WatchLogDeposit(opts, events)
	if err != nil {
		return nil, err
	}

	return events, nil
}


func (e *Exchange) GetErrorEvents(logs chan *contracts.ExchangeLogError) error {
	opts := &bind.WatchOpts{nil, nil}

	_, err := e.Interface.WatchLogError(opts, logs)
	if err != nil {
		return err
	}

	return nil
}

func (e *Exchange) GetTrades(logs chan *contracts.ExchangeLogTrade) error {
	opts := &bind.WatchOpts{nil, nil}

	_, err := e.Interface.WatchLogTrade(opts, logs)
	if err != nil {
		return err
	}

	return nil
}

func (e *Exchange) PrintTrades() error {
	events := make(chan *contracts.ExchangeLogTrade)
	opts := &bind.WatchOpts{nil, nil}

	_, err := e.Interface.WatchLogTrade(opts, events)
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
	events := make(chan *contracts.ExchangeLogError)
	opts := &bind.WatchOpts{nil, nil}

	_, err := e.Interface.WatchLogError(opts, events)
	if err != nil {
		return err
	}

	go func() {
		for {
			event := <-events
			log.Printf("New Error Event. Id: %v, Hash: %v\n\n", event.ErrorId, hex.EncodeToString(event.OrderHash[:]))
		}
	}()

	return nil
}

func PrintErrorLog(log *contracts.ExchangeLogError) string {
	return fmt.Sprintf("Error:\nErrorID: %v\nOrderHash: %v\n\n", log.ErrorId, log.OrderHash)
}

func PrintTradeLog(log *contracts.ExchangeLogTrade) string {
	return fmt.Sprintf("Error:\nAmount: %v\nMaker: %v\nTaker: %v\nTokenBuy: %v\nTokenSell: %v\nOrderHash: %v\nTradeHash: %v\n\n",
		log.Amount, log.Maker, log.Taker, log.TokenBuy, log.TokenSell, log.OrderHash, log.TradeHash)
}

func PrintCancelOrderLog(log *contracts.ExchangeLogCancelOrder) string {
	return fmt.Sprintf("Error:\nSender: %v\nOrderHash: %v\n\n", log.Sender, log.OrderHash)
}

func PrintCancelTradeLog(log *contracts.ExchangeLogCancelTrade) string {
	return fmt.Sprintf("Error:\nSender: %v\nTradeHash: %v\n\n", log.Sender, log.TradeHash)
}

func PrintWithdrawalErrorLog(log *contracts.ExchangeLogWithdrawalError) string {
	return fmt.Sprintf("Error:\nError ID: %v\n, WithdrawalHash: %v\n\n", log.ErrorId, log.WithdrawalHash)
}