package dex

import (
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	. "github.com/ethereum/go-ethereum/common"
)

type Operator struct {
	Wallet   *Wallet
	Exchange *Exchange
	Chain    *bind.ContractBackend
}

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

func (o *Operator) SetDefaultTxOptions() {
	o.Exchange.TxOptions.Value = big.NewInt(0)
}

func (o *Operator) SetTxValue(value *big.Int) {
	o.Exchange.TxOptions.Value = value
}

func (o *Operator) SetCustomSender(w *Wallet) {
	o.Exchange.TxOptions = bind.NewKeyedTransactor(w.PrivateKey)
}

func (o *Operator) SetFeeAccount(account Address) (Transaction, error) {
	tx, err := o.Exchange.SetFeeAccount(account)
	if err != nil {
		return nil, err
	}

	return tx, nil
}

func (o *Operator) SetOperator(account Address, isOperator bool) (Transaction, error) {
	tx, err := o.Exchange.SetOperator(account, isOperator)
	if err != nil {
		return nil, err
	}

	return tx, nil
}

func (o *Operator) SetWithdrawalSecurityPeriod(p *big.Int) (Transaction, error) {
	tx, err := o.Exchange.SetWithdrawalSecurityPeriod(p)
	if err != nil {
		return nil, err
	}

	return tx, nil
}

func (o *Operator) DepositEther(val *big.Int) (Transaction, error) {
	tx, err := o.Exchange.DepositEther(val)
	if err != nil {
		return nil, err
	}

	return tx, nil
}

func (o *Operator) DepositToken(token Address, amount *big.Int) (Transaction, error) {
	tx, err := o.Exchange.DepositToken(token, amount)
	if err != nil {
		return nil, err
	}

	return tx, err
}

func (o *Operator) TokenBalance(account Address, token Address) (*big.Int, error) {
	b, err := o.Exchange.TokenBalance(account, token)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func (o *Operator) EtherBalance(account Address) (*big.Int, error) {
	b, err := o.Exchange.EtherBalance(account)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func (o *Operator) WithdrawalSecurityPeriod() (*big.Int, error) {
	p, err := o.Exchange.WithdrawalSecurityPeriod()
	if err != nil {
		return nil, err
	}

	return p, nil
}

func (o *Operator) FeeAccount() (Address, error) {
	account, err := o.Exchange.FeeAccount()
	if err != nil {
		return Address{}, err
	}

	return account, nil
}

func (o *Operator) Operator(addr Address) (bool, error) {
	isOperator, err := o.Exchange.Operator(addr)
	if err != nil {
		return false, err
	}

	return isOperator, nil
}

func (o *Operator) SecurityWithdraw(w *Wallet, token Address, amount *big.Int) (Transaction, error) {
	tx, err := o.Exchange.SecurityWithdraw(w, token, amount)
	if err != nil {
		return nil, err
	}

	return tx, nil
}

func (o *Operator) Withdraw(w *Withdrawal) (Transaction, error) {
	tx, err := o.Exchange.Withdraw(w)
	if err != nil {
		return nil, err
	}

	return tx, nil
}

func (o *Operator) Trade(ord *Order, t *Trade) (Transaction, error) {
	tx, err := o.Exchange.Trade(ord, t)
	if err != nil {
		return nil, err
	}

	return tx, nil
}
