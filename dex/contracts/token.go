package contracts

import (
	"errors"
	"fmt"
	"math/big"

	"github.com/Dvisacker/proofsuite-orderbook/dex"
	"github.com/Dvisacker/proofsuite-orderbook/dex/contracts/interfaces"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	. "github.com/ethereum/go-ethereum/common"
)

type Token struct {
	Address     Address
	Contract    *interfaces.Token
	CallOptions *bind.CallOpts
	TxOptions   *bind.TransactOpts
}

func (t *Token) SetCustomSender(wallet *dex.Wallet) {
	txOptions := bind.NewKeyedTransactor(wallet.PrivateKey)
	t.TxOptions = txOptions
}

func (t *Token) SetDefaultSender() {
	txOptions := bind.NewKeyedTransactor(config.Wallets[0].PrivateKey)
	t.TxOptions = txOptions
}

func (t *Token) BalanceOf(owner Address) (*big.Int, error) {
	balance, err := t.Contract.BalanceOf(t.CallOptions, owner)
	if err != nil {
		return nil, err
	}

	return balance, nil
}

func (t *Token) TotalSupply() (*big.Int, error) {
	supply, err := t.Contract.TotalSupply(t.CallOptions)
	if err != nil {
		return nil, errors.New("Could not retrieve total supply of user")
	}

	return supply, nil
}

func (t *Token) Transfer(receiver Address, amount *big.Int) (Hash, error) {
	tx, err := t.Contract.Transfer(t.TxOptions, receiver, amount)
	if err != nil {
		return Hash{}, errors.New("Error making Transfer() transaction")
	}

	return tx.Hash(), nil
}

func (t *Token) TransferFrom(sender Address, receiver Address, amount *big.Int) (Hash, error) {
	tx, err := t.Contract.TransferFrom(t.TxOptions, sender, receiver, amount)
	if err != nil {
		return Hash{}, errors.New("Error making TransferFrom() transaction")
	}

	fmt.Printf("Transfered %v tokens from %v to %v", amount, sender, receiver)
	return tx.Hash(), nil
}

func (t *Token) Allowance(owner Address, spender Address) (*big.Int, error) {
	allowance, err := t.Contract.Allowance(t.CallOptions, owner, spender)
	if err != nil {
		return nil, errors.New("Error retrieving allowance")
	}

	return allowance, nil
}

func (t *Token) Approve(spender Address, amount *big.Int) (Hash, error) {
	tx, err := t.Contract.Approve(t.TxOptions, spender, amount)
	if err != nil {
		return Hash{}, errors.New("Error making Approve() transaction")
	}

	return tx.Hash(), nil
}

func (t *Token) ListenToTransferEvents() (chan *interfaces.TokenTransfer, error) {
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

func (t *Token) PrintTransferEvents() error {
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
