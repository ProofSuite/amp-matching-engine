package dex

import (
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	. "github.com/ethereum/go-ethereum/common"
)

// Operator manages the transaction queue that will eventually be
// sent to the exchange contract. The Operator Wallet must be equal to the
// account that initially deployed the exchange contract or an address with operator rights
// on the contract
type Operator struct {
	Wallet   *Wallet
	Exchange *Exchange
	Chain    *bind.ContractBackend
}

// NewOperator returns a new operator object
func NewOperator(w *Wallet, contractAddr Address, chain bind.ContractBackend) (*Operator, error) {
	o := &Operator{}

	e, err := NewExchange(w, contractAddr, chain)
	if err != nil {
		return nil, err
	}

	o.Wallet = w
	o.Exchange = e
	o.Chain = &chain
	return o, nil
}

// SetDefaultTxOptions resets the transaction value to 0
func (o *Operator) SetDefaultTxOptions() {
	o.Exchange.TxOptions.Value = big.NewInt(0)
}

// SetTxValue sets the transaction ether value
func (o *Operator) SetTxValue(value *big.Int) {
	o.Exchange.TxOptions.Value = value
}

// SetCustomSender updates the sender address address to the exchange contract
func (o *Operator) SetCustomSender(w *Wallet) {
	o.Exchange.TxOptions = bind.NewKeyedTransactor(w.PrivateKey)
}

// SetFeeAccount sets the fee account of the exchange contract. The fee account receives
// the trading fees whenever a trade is settled.
func (o *Operator) SetFeeAccount(account Address) (Transaction, error) {
	tx, err := o.Exchange.SetFeeAccount(account)
	if err != nil {
		return nil, err
	}

	return tx, nil
}

// SetOperator updates the operator settings of the given address. Only addresses with an
// operator access can execute Withdraw and Trade transactions to the Exchange smart contract
func (o *Operator) SetOperator(account Address, isOperator bool) (Transaction, error) {
	tx, err := o.Exchange.SetOperator(account, isOperator)
	if err != nil {
		return nil, err
	}

	return tx, nil
}

// SetWithdrawalSecurityPeriod sets the period after which a non-operator address can send
// a transaction to the exchange smart-contract to withdraw their funds. This acts as security mechanism
// to prevent the operator of the exchange from holding funds
func (o *Operator) SetWithdrawalSecurityPeriod(p *big.Int) (Transaction, error) {
	tx, err := o.Exchange.SetWithdrawalSecurityPeriod(p)
	if err != nil {
		return nil, err
	}

	return tx, nil
}

// DepositEther deposits ether into the exchange smart-contract. A priori this function is not supposed
// to be called by the exchange operator
func (o *Operator) DepositEther(val *big.Int) (Transaction, error) {
	tx, err := o.Exchange.DepositEther(val)
	if err != nil {
		return nil, err
	}

	return tx, nil
}

// DepositToken deposits tokens into the exchange smart-contract. A priori this function is not supposed
// to be called by the exchange operator
func (o *Operator) DepositToken(token Address, amount *big.Int) (Transaction, error) {
	tx, err := o.Exchange.DepositToken(token, amount)
	if err != nil {
		return nil, err
	}

	return tx, err
}

// TokenBalance returns the Exchange token balance of the given token at the given account address.
// Note: This is not the token BalanceOf() function, it's the balance of tokens that have been deposited
// in the exchange smart contract.
func (o *Operator) TokenBalance(account Address, token Address) (*big.Int, error) {
	b, err := o.Exchange.TokenBalance(account, token)
	if err != nil {
		return nil, err
	}

	return b, nil
}

// EtherBalalance returns the Exchange ether balance of the given account address.
// Note: This is not the current ether balance of the given ether address. It's the balance of ether
// that has been deposited in the exchange smart contract.
func (o *Operator) EtherBalance(account Address) (*big.Int, error) {
	b, err := o.Exchange.EtherBalance(account)
	if err != nil {
		return nil, err
	}

	return b, nil
}

// WithdrawalSecurityPeriod is the period after which a non-operator account can withdraw their funds from
// the exchange smart contract.
func (o *Operator) WithdrawalSecurityPeriod() (*big.Int, error) {
	p, err := o.Exchange.WithdrawalSecurityPeriod()
	if err != nil {
		return nil, err
	}

	return p, nil
}

// FeeAccount is the Ethereum towards the exchange trading fees are sent
func (o *Operator) FeeAccount() (Address, error) {
	account, err := o.Exchange.FeeAccount()
	if err != nil {
		return Address{}, err
	}

	return account, nil
}

// Operator returns true if the given address is an operator of the exchange and returns false otherwise
func (o *Operator) Operator(addr Address) (bool, error) {
	isOperator, err := o.Exchange.Operator(addr)
	if err != nil {
		return false, err
	}

	return isOperator, nil
}

// SecurityWithdraw executes a security withdraw transaction. Security withdraw transactions can only be
// executed after the security withdrawal period has ended. A priori, this function should not be called
// by the operator account itself
func (o *Operator) SecurityWithdraw(w *Wallet, token Address, amount *big.Int) (Transaction, error) {
	tx, err := o.Exchange.SecurityWithdraw(w, token, amount)
	if err != nil {
		return nil, err
	}

	return tx, nil
}

// Withdraw executes a normal withdraw transaction. This withdraws tokens or ether from the exchange
// and returns them to the payload Receiver. Only an operator account can send a withdraw
// transaction
func (o *Operator) Withdraw(w *Withdrawal) (Transaction, error) {
	tx, err := o.Exchange.Withdraw(w)
	if err != nil {
		return nil, err
	}

	return tx, nil
}

// Trade executes a settlements transaction. The order and trade payloads need to be signed respectively
// by the Maker and the Taker of the trade. Only the operator account can send a Trade function to the
// Exchange smart contract.
func (o *Operator) Trade(ord *Order, t *Trade) (Transaction, error) {
	tx, err := o.Exchange.Trade(ord, t)
	if err != nil {
		return nil, err
	}

	return tx, nil
}
