package contracts

import (
	"encoding/hex"
	"fmt"
	"math/big"

	"github.com/Dvisacker/proofsuite-orderbook/dex"
	"github.com/Dvisacker/proofsuite-orderbook/dex/contracts/interfaces"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	. "github.com/ethereum/go-ethereum/common"
)

var config = dex.NewDefaultConfiguration()

// type Exchange struct {
// 	contract *interfaces.ExchangeSession
// }

type Exchange struct {
	Address     Address
	Contract    *interfaces.Exchange
	CallOptions *bind.CallOpts
	TxOptions   *bind.TransactOpts
}

func NewExchange(wallet *dex.Wallet, contractAddress Address, backend bind.ContractBackend) (*Exchange, error) {
	instance, err := interfaces.NewExchange(contractAddress, backend)
	if err != nil {
		return nil, err
	}

	callOptions := &bind.CallOpts{Pending: true}
	txOptions := bind.NewKeyedTransactor(wallet.PrivateKey)

	return &Exchange{
		Address:     contractAddress,
		Contract:    instance,
		CallOptions: callOptions,
		TxOptions:   txOptions,
	}, nil
}

//TODO add more default options
func (e *Exchange) SetDefaultTxOptions() {
	e.TxOptions.Value = big.NewInt(0)
	// e.Contract.TransactOpts.Value = big.NewInt(0)
}

func (e *Exchange) SetTxValue(value *big.Int) {
	e.TxOptions.Value = value
}

func (e *Exchange) SetCustomSender(wallet *dex.Wallet) {
	txOptions := bind.NewKeyedTransactor(wallet.PrivateKey)
	e.TxOptions = txOptions
}

func (e *Exchange) SetDefaultSender() {
	txOptions := bind.NewKeyedTransactor(config.Wallets[0].PrivateKey)
	e.TxOptions = txOptions
}

func (e *Exchange) SetFeeAccount(account Address) (dex.Transaction, error) {
	tx, err := e.Contract.SetFeeAccount(e.TxOptions, account)
	if err != nil {
		return nil, err
	}

	return tx, nil
}

func (e *Exchange) SetOperator(account Address, isOperator bool) (dex.Transaction, error) {
	tx, err := e.Contract.SetOperator(e.TxOptions, account, isOperator)
	if err != nil {
		return nil, err
	}

	return tx, nil
}

func (e *Exchange) SetWithdrawalSecurityPeriod(p *big.Int) (dex.Transaction, error) {
	tx, err := e.Contract.SetWithdrawalSecurityPeriod(e.TxOptions, p)
	if err != nil {
		return nil, err
	}

	return tx, nil
}

func (e *Exchange) DepositEther(value *big.Int) (dex.Transaction, error) {
	e.SetTxValue(value)

	tx, err := e.Contract.DepositEther(e.TxOptions)
	if err != nil {
		return nil, err
	}

	e.SetDefaultTxOptions()
	return tx, nil
}

func (e *Exchange) DepositToken(token Address, amount *big.Int) (dex.Transaction, error) {
	// e.SetDefaultTxOptions()

	tx, err := e.Contract.DepositToken(e.TxOptions, token, amount)
	if err != nil {
		return nil, err
	}

	return tx, err
}

func (e *Exchange) TokenBalance(trader Address, token Address) (*big.Int, error) {
	balance, err := e.Contract.TokenBalance(e.CallOptions, trader, token)
	if err != nil {
		return nil, err
	}

	return balance, nil
}

func (e *Exchange) EtherBalance(trader Address) (*big.Int, error) {
	balance, err := e.Contract.EtherBalance(e.CallOptions, trader)
	if err != nil {
		return nil, err
	}

	return balance, nil
}

func (e *Exchange) WithdrawalSecurityPeriod() (*big.Int, error) {
	period, err := e.Contract.WithdrawalSecurityPeriod(e.CallOptions)
	if err != nil {
		return nil, err
	}

	return period, nil
}

func (e *Exchange) FeeAccount() (Address, error) {
	account, err := e.Contract.FeeAccount(e.CallOptions)
	if err != nil {
		return Address{}, err
	}

	return account, nil
}

func (e *Exchange) Operator(address Address) (bool, error) {
	isOperator, err := e.Contract.Operators(e.CallOptions, address)
	if err != nil {
		return false, err
	}

	return isOperator, nil
}

func (e *Exchange) SecurityWithdraw(wallet *dex.Wallet, token Address, amount *big.Int) (dex.Transaction, error) {
	e.SetDefaultTxOptions()
	e.SetCustomSender(wallet)

	tx, err := e.Contract.SecurityWithdraw(e.TxOptions, token, amount)
	if err != nil {
		return nil, err
	}

	return tx, nil
}

func (e *Exchange) Withdraw(w *dex.Withdrawal) (dex.Transaction, error) {
	e.SetDefaultTxOptions()

	s := w.Signature
	tx, err := e.Contract.Withdraw(e.TxOptions, w.Token, w.Amount, w.Trader, w.Receiver, w.Nonce, s.V, [2][32]byte{s.R, s.S}, w.Fee)
	if err != nil {
		return nil, err
	}

	return tx, nil
}

func (e *Exchange) Trade(o *dex.Order, t *dex.Trade) (dex.Transaction, error) {
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

func (e *Exchange) ListenToErrorEvents() (chan *interfaces.ExchangeLogError, error) {
	events := make(chan *interfaces.ExchangeLogError)
	options := &bind.WatchOpts{nil, nil}

	_, err := e.Contract.WatchLogError(options, events)
	if err != nil {
		return nil, err
	}

	return events, nil
}

func (e *Exchange) ListenToTrades() (chan *interfaces.ExchangeLogTrade, error) {
	events := make(chan *interfaces.ExchangeLogTrade)
	options := &bind.WatchOpts{nil, nil}

	_, err := e.Contract.WatchLogTrade(options, events)
	if err != nil {
		return nil, err
	}

	return events, nil
}

func (e *Exchange) ListenToDeposits() (chan *interfaces.ExchangeLogDeposit, error) {
	events := make(chan *interfaces.ExchangeLogDeposit)
	options := &bind.WatchOpts{nil, nil}

	_, err := e.Contract.WatchLogDeposit(options, events)
	if err != nil {
		return nil, err
	}

	return events, nil
}

func (e *Exchange) PrintTrades() error {
	events := make(chan *interfaces.ExchangeLogTrade)
	options := &bind.WatchOpts{nil, nil}

	_, err := e.Contract.WatchLogTrade(options, events)
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
	options := &bind.WatchOpts{nil, nil}

	_, err := e.Contract.WatchLogError(options, events)
	if err != nil {
		return err
	}

	go func() {
		for {
			event := <-events
			fmt.Printf("New Error Event. Id: %v, Hash: %v\n\n", event.ErrorId, hex.EncodeToString(event.OrderHash[:]))
		}
	}()

	return nil
}
