package dex

import (
	"encoding/hex"
	"fmt"
	"log"
	"math/big"

	"github.com/Dvisacker/matching-engine/dex/interfaces"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	. "github.com/ethereum/go-ethereum/common"
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
	Admin       *Wallet
	Address     Address
	Contract    *interfaces.Exchange
	CallOptions *bind.CallOpts
	TxOptions   *bind.TransactOpts
}

// Returns a new exchange interface for a given wallet, contract address and connected backend.
// The exchange contract need to be already deployed at the given address. The given wallet will
// be used by default when sending transactions with this object.
func NewExchange(admin *Wallet, contractAddress Address, backend bind.ContractBackend) (*Exchange, error) {
	instance, err := interfaces.NewExchange(contractAddress, backend)
	if err != nil {
		return nil, err
	}

	callOptions := &bind.CallOpts{Pending: true}
	txOptions := bind.NewKeyedTransactor(admin.PrivateKey)

	return &Exchange{
		Address:     contractAddress,
		Contract:    instance,
		CallOptions: callOptions,
		TxOptions:   txOptions,
		Admin:       admin,
	}, nil
}

// SetDefaultTxOptions resets the transaction value to 0
func (e *Exchange) SetDefaultTxOptions() {
	e.TxOptions = bind.NewKeyedTransactor(e.Admin.PrivateKey)
	e.TxOptions.Value = big.NewInt(0)
}

// SetTxValue sets the transaction ether value
func (e *Exchange) SetTxValue(value *big.Int) {
	e.TxOptions.Value = value
}

// SetCustomSender updates the sender address address to the exchange contract
func (e *Exchange) SetCustomSender(wallet *Wallet) {
	e.TxOptions = bind.NewKeyedTransactor(wallet.PrivateKey)
}

// SetDefaultSender sets the default sender address that will be used when sending a transcation to
// the exchange contract
func (e *Exchange) SetDefaultSender() {
	e.TxOptions = bind.NewKeyedTransactor(e.Admin.PrivateKey)
}

// SetFeeAccount sets the fee account of the exchange contract. The fee account receives
// the trading fees whenever a trade is settled.
func (e *Exchange) SetFeeAccount(account Address) (*types.Transaction, error) {
	tx, err := e.Contract.SetFeeAccount(e.TxOptions, account)
	if err != nil {
		return nil, err
	}

	return tx, nil
}

// SetOperator updates the operator settings of the given address. Only addresses with an
// operator access can execute Withdraw and Trade transactions to the Exchange smart contract
func (e *Exchange) SetOperator(account Address, isOperator bool) (*types.Transaction, error) {
	tx, err := e.Contract.SetOperator(e.TxOptions, account, isOperator)
	if err != nil {
		return nil, err
	}

	return tx, nil
}

// SetWithdrawalSecurityPeriod sets the period after which a non-operator address can send
// a transaction to the exchange smart-contract to withdraw their funds. This acts as security mechanism
// to prevent the operator of the exchange from holding funds
func (e *Exchange) SetWithdrawalSecurityPeriod(p *big.Int) (*types.Transaction, error) {
	tx, err := e.Contract.SetWithdrawalSecurityPeriod(e.TxOptions, p)
	if err != nil {
		return nil, err
	}

	return tx, nil
}

// DepositEther deposits ether into the exchange smart-contract.
func (e *Exchange) DepositEther(value *big.Int) (*types.Transaction, error) {
	e.SetTxValue(value)

	tx, err := e.Contract.DepositEther(e.TxOptions)
	if err != nil {
		return nil, err
	}

	e.SetDefaultTxOptions()
	return tx, nil
}

// DepositEtherFrom deposits ether from a custom address. The transaction sender is reset
// after the transaction is carried out.
func (e *Exchange) DepositEtherFrom(wallet *Wallet, value *big.Int) (*types.Transaction, error) {
	e.SetTxValue(value)
	e.SetCustomSender(wallet)

	tx, err := e.Contract.DepositEther(e.TxOptions)
	if err != nil {
		return nil, err
	}

	e.SetDefaultTxOptions()
	return tx, nil
}

// DepositToken deposits tokens into the exchange smart-contract.
func (e *Exchange) DepositToken(token Address, amount *big.Int) (*types.Transaction, error) {
	// e.SetDefaultTxOptions()

	tx, err := e.Contract.DepositToken(e.TxOptions, token, amount)
	if err != nil {
		return nil, err
	}

	return tx, err
}

// DepositEtherFrom deposits ether from a custom address. The transaction sender is reset
// after the transaction is carried out.
func (e *Exchange) DepositTokenFrom(wallet *Wallet, token Address, amount *big.Int) (*types.Transaction, error) {
	e.SetCustomSender(wallet)

	tx, err := e.Contract.DepositToken(e.TxOptions, token, amount)
	if err != nil {
		return nil, err
	}

	e.SetDefaultSender()
	return tx, err
}

// TokenBalance returns the Exchange token balance of the given token at the given account address.
// Note: This is not the token BalanceOf() function, it's the balance of tokens that have been deposited
// in the exchange smart contract.
func (e *Exchange) TokenBalance(trader Address, token Address) (*big.Int, error) {
	balance, err := e.Contract.TokenBalance(e.CallOptions, trader, token)
	if err != nil {
		return nil, err
	}

	return balance, nil
}

// EtherBalance returns the Exchange ether balance of the given account address.
// Note: This is not the current ether balance of the given ether address. It's the balance of ether
// that has been deposited in the exchange smart contract.
func (e *Exchange) EtherBalance(trader Address) (*big.Int, error) {
	balance, err := e.Contract.EtherBalance(e.CallOptions, trader)
	if err != nil {
		return nil, err
	}

	return balance, nil
}

// WithdrawalSecurityPeriod is the period after which a non-operator account can withdraw their funds from
// the exchange smart contract.
func (e *Exchange) WithdrawalSecurityPeriod() (*big.Int, error) {
	period, err := e.Contract.WithdrawalSecurityPeriod(e.CallOptions)
	if err != nil {
		return nil, err
	}

	return period, nil
}

// FeeAccount is the Ethereum towards the exchange trading fees are sent
func (e *Exchange) FeeAccount() (Address, error) {
	account, err := e.Contract.FeeAccount(e.CallOptions)
	if err != nil {
		return Address{}, err
	}

	return account, nil
}

// Operator returns true if the given address is an operator of the exchange and returns false otherwise
func (e *Exchange) Operator(address Address) (bool, error) {
	isOperator, err := e.Contract.Operators(e.CallOptions, address)
	if err != nil {
		return false, err
	}

	return isOperator, nil
}

// SecurityWithdraw executes a security withdraw transaction. Security withdraw transactions can only be
// executed after the security withdrawal period has ended.
func (e *Exchange) SecurityWithdraw(wallet *Wallet, token Address, amount *big.Int) (*types.Transaction, error) {
	e.SetDefaultTxOptions()
	e.SetCustomSender(wallet)

	tx, err := e.Contract.SecurityWithdraw(e.TxOptions, token, amount)
	if err != nil {
		return nil, err
	}

	return tx, nil
}

// Withdraw executes a normal withdraw transaction. This withdraws tokens or ether from the exchange
// and returns them to the payload Receiver. Only an operator account can send a withdraw
// transaction
func (e *Exchange) Withdraw(w *Withdrawal) (*types.Transaction, error) {
	e.SetDefaultTxOptions()

	s := w.Signature
	tx, err := e.Contract.Withdraw(e.TxOptions, w.Token, w.Amount, w.Trader, w.Receiver, w.Nonce, s.V, [2][32]byte{s.R, s.S}, w.Fee)
	if err != nil {
		return nil, err
	}

	return tx, nil
}

// Trade executes a settlements transaction. The order and trade payloads need to be signed respectively
// by the Maker and the Taker of the trade. Only the operator account can send a Trade function to the
// Exchange smart contract.
func (e *Exchange) Trade(o *Order, t *Trade) (*types.Transaction, error) {
	e.SetDefaultTxOptions()

	orderValues := [8]*big.Int{o.AmountBuy, o.AmountSell, o.Expires, o.Nonce, o.FeeMake, o.FeeTake, t.Amount, t.TradeNonce}
	orderAddresses := [4]Address{o.TokenBuy, o.TokenSell, o.Maker, t.Taker}
	vValues := [2]uint8{o.Signature.V, t.Signature.V}
	rsValues := [4][32]byte{o.Signature.R, o.Signature.S, t.Signature.R, t.Signature.S}

	tx, err := e.Contract.ExecuteTrade(e.TxOptions, orderValues, orderAddresses, vValues, rsValues)
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
func (e *Exchange) ListenToErrorEvents() (chan *interfaces.ExchangeLogError, error) {
	events := make(chan *interfaces.ExchangeLogError)
	opts := &bind.WatchOpts{nil, nil}

	_, err := e.Contract.WatchLogError(opts, events)
	if err != nil {
		return nil, err
	}

	return events, nil
}

func (e *Exchange) GetErrorEvents(logs chan *interfaces.ExchangeLogError) error {
	opts := &bind.WatchOpts{nil, nil}

	_, err := e.Contract.WatchLogError(opts, logs)
	if err != nil {
		return err
	}

	return nil
}

// ListenToTrades returns a channel that receivs trade logs (events) from the underlying exchange smart contract
func (e *Exchange) ListenToTrades() (chan *interfaces.ExchangeLogTrade, error) {
	events := make(chan *interfaces.ExchangeLogTrade)
	opts := &bind.WatchOpts{nil, nil}

	_, err := e.Contract.WatchLogTrade(opts, events)
	if err != nil {
		return nil, err
	}

	return events, nil
}

func (e *Exchange) GetTrades(logs chan *interfaces.ExchangeLogTrade) error {
	opts := &bind.WatchOpts{nil, nil}

	_, err := e.Contract.WatchLogTrade(opts, logs)
	if err != nil {
		return err
	}

	return nil
}

// ListenToTrades returns a channel that receivs deposit logs (events) from the underlying exchange smart contract
func (e *Exchange) ListenToDeposits() (chan *interfaces.ExchangeLogDeposit, error) {
	events := make(chan *interfaces.ExchangeLogDeposit)
	opts := &bind.WatchOpts{nil, nil}

	_, err := e.Contract.WatchLogDeposit(opts, events)
	if err != nil {
		return nil, err
	}

	return events, nil
}

func (e *Exchange) PrintTrades() error {
	events := make(chan *interfaces.ExchangeLogTrade)
	opts := &bind.WatchOpts{nil, nil}

	_, err := e.Contract.WatchLogTrade(opts, events)
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
	events := make(chan *interfaces.ExchangeLogError)
	opts := &bind.WatchOpts{nil, nil}

	_, err := e.Contract.WatchLogError(opts, events)
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

func PrintErrorLog(log *interfaces.ExchangeLogError) string {
	return fmt.Sprintf("Error:\nErrorID: %v\nOrderHash: %v\n\n", log.ErrorId, log.OrderHash)
}

func PrintTradeLog(log *interfaces.ExchangeLogTrade) string {
	return fmt.Sprintf("Error:\nAmount: %v\nMaker: %v\nTaker: %v\nTokenBuy: %v\nTokenSell: %v\nOrderHash: %v\nTradeHash: %v\n\n",
		log.Amount, log.Maker, log.Taker, log.TokenBuy, log.TokenSell, log.OrderHash, log.TradeHash)
}

func PrintCancelOrderLog(log *interfaces.ExchangeLogCancelOrder) string {
	return fmt.Sprintf("Error:\nSender: %v\nOrderHash: %v\n\n", log.Sender, log.OrderHash)
}

func PrintCancelTradeLog(log *interfaces.ExchangeLogCancelTrade) string {
	return fmt.Sprintf("Error:\nSender: %v\nTradeHash: %v\n\n", log.Sender, log.TradeHash)
}

func PrintWithdrawalErrorLog(log *interfaces.ExchangeLogWithdrawalError) string {
	return fmt.Sprintf("Error:\nError ID: %v\n, WithdrawalHash: %v\n\n", log.ErrorId, log.WithdrawalHash)
}

// func Print(log *interfaces.ExchangeLogError) string {
// 	return fmt.Sprintf("Error:\nErrorID: %v\nOrderHash: %v\n\n", log.)
// }
