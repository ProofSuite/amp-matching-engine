package dex

import (
	"errors"
	"fmt"
	"math/big"

	"github.com/Dvisacker/matching-engine/dex/interfaces"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	. "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

type ERC20Token struct {
	Address       Address
	Contract      *interfaces.Token
	CallOptions   *bind.CallOpts
	TxOptions     *bind.TransactOpts
	DefaultSender *Wallet
}

func (t *ERC20Token) SetCustomSender(w *Wallet) {
	t.TxOptions = bind.NewKeyedTransactor(w.PrivateKey)
}

func (t *ERC20Token) SetDefaultSender() {
	t.TxOptions = bind.NewKeyedTransactor(t.DefaultSender.PrivateKey)
}

func (t *ERC20Token) BalanceOf(owner Address) (*big.Int, error) {
	b, err := t.Contract.BalanceOf(t.CallOptions, owner)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func (t *ERC20Token) TotalSupply() (*big.Int, error) {
	supply, err := t.Contract.TotalSupply(t.CallOptions)
	if err != nil {
		return nil, errors.New("Could not retrieve total supply of user")
	}

	return supply, nil
}

func (t *ERC20Token) Transfer(receiver Address, amount *big.Int) (*types.Transaction, error) {
	tx, err := t.Contract.Transfer(t.TxOptions, receiver, amount)
	if err != nil {
		return nil, errors.New("Error making Transfer() transaction")
	}

	return tx, nil
}

func (t *ERC20Token) TransferFromCustomWallet(wallet *Wallet, receiver Address, amount *big.Int) (*types.Transaction, error) {
	t.SetCustomSender(wallet)

	tx, err := t.Contract.Transfer(t.TxOptions, receiver, amount)
	if err != nil {
		return nil, errors.New("Error making Transfer() transaction")
	}

	return tx, nil
}

func (t *ERC20Token) TransferFrom(sender, receiver Address, amount *big.Int) (*types.Transaction, error) {
	tx, err := t.Contract.TransferFrom(t.TxOptions, sender, receiver, amount)
	if err != nil {
		return nil, errors.New("Error making TransferFrom() transaction")
	}

	fmt.Printf("Transfered %v tokens from %v to %v", amount, sender, receiver)
	return tx, nil
}

func (t *ERC20Token) Allowance(owner Address, spender Address) (*big.Int, error) {
	allowance, err := t.Contract.Allowance(t.CallOptions, owner, spender)
	if err != nil {
		return nil, errors.New("Error retrieving allowance")
	}

	return allowance, nil
}

func (t *ERC20Token) Approve(spender Address, amount *big.Int) (*types.Transaction, error) {
	tx, err := t.Contract.Approve(t.TxOptions, spender, amount)
	if err != nil {
		return nil, errors.New("Error making Approve() transaction")
	}

	return tx, nil
}

func (t *ERC20Token) ApproveFrom(wallet *Wallet, spender Address, amount *big.Int) (*types.Transaction, error) {
	t.SetCustomSender(wallet)

	tx, err := t.Contract.Approve(t.TxOptions, spender, amount)
	if err != nil {
		return nil, errors.New("Error making ApproveFrom() transaction")
	}

	t.SetDefaultSender()
	return tx, nil
}

func (t *ERC20Token) ListenToTransferEvents() (chan *interfaces.TokenTransfer, error) {
	events := make(chan *interfaces.TokenTransfer)
	options := &bind.WatchOpts{nil, nil}
	toList := []Address{}
	fromList := []Address{}

	_, err := t.Contract.WatchTransfer(options, events, fromList, toList)
	if err != nil {
		return nil, err
	}

	return events, nil
}

func (t *ERC20Token) PrintTransferEvents() error {
	events := make(chan *interfaces.TokenTransfer)
	options := &bind.WatchOpts{nil, nil}

	toList := []Address{}
	fromList := []Address{}

	_, err := t.Contract.WatchTransfer(options, events, fromList, toList)
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
