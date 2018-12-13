package contracts

import (
	"errors"
	"fmt"
	"math/big"

	"github.com/Proofsuite/amp-matching-engine/contracts/contractsinterfaces"
	"github.com/Proofsuite/amp-matching-engine/interfaces"
	"github.com/Proofsuite/amp-matching-engine/types"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	eth "github.com/ethereum/go-ethereum/core/types"
)

type Token struct {
	WalletService interfaces.WalletService
	TxService     interfaces.TxService
	Interface     *contractsinterfaces.ERC20
}

func NewToken(
	w interfaces.WalletService,
	tx interfaces.TxService,
	contractAddress common.Address,
	backend bind.ContractBackend,
) (*Token, error) {
	instance, err := contractsinterfaces.NewERC20(contractAddress, backend)
	if err != nil {
		return nil, err
	}

	return &Token{
		WalletService: w,
		TxService:     tx,
		Interface:     instance,
	}, nil
}

func (t *Token) SetTxSender(w *types.Wallet) {
	t.TxService.SetTxSender(w)
}

func (t *Token) GetTxCallOptions() *bind.CallOpts {
	return t.TxService.GetTxCallOptions()
}

func (t *Token) GetTxSendOptions() (*bind.TransactOpts, error) {
	return t.TxService.GetTxSendOptions()
}

func (t *Token) GetCustomTxSendOptions(w *types.Wallet) *bind.TransactOpts {
	return t.TxService.GetCustomTxSendOptions(w)
}

func (t *Token) BalanceOf(owner common.Address) (*big.Int, error) {
	opts := t.GetTxCallOptions()

	b, err := t.Interface.BalanceOf(opts, owner)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func (t *Token) TotalSupply() (*big.Int, error) {
	opts := t.GetTxCallOptions()

	supply, err := t.Interface.TotalSupply(opts)
	if err != nil {
		return nil, errors.New("Could not retrieve total supply of user")
	}

	return supply, nil
}

func (t *Token) Transfer(receiver common.Address, amount *big.Int) (*eth.Transaction, error) {
	opts, _ := t.GetTxSendOptions()

	tx, err := t.Interface.Transfer(opts, receiver, amount)
	if err != nil {
		return nil, errors.New("Error making Transfer() transaction")
	}

	return tx, nil
}

func (t *Token) TransferFromCustomWallet(w *types.Wallet, receiver common.Address, amount *big.Int) (*eth.Transaction, error) {
	opts := t.GetCustomTxSendOptions(w)

	tx, err := t.Interface.Transfer(opts, receiver, amount)
	if err != nil {
		return nil, errors.New("Error making Transfer() transaction")
	}

	return tx, nil
}

func (t *Token) TransferFrom(sender, receiver common.Address, amount *big.Int) (*eth.Transaction, error) {
	opts, _ := t.GetTxSendOptions()
	tx, err := t.Interface.TransferFrom(opts, sender, receiver, amount)
	if err != nil {
		return nil, errors.New("Error making TransferFrom() transaction")
	}

	fmt.Printf("Transfered %v tokens from %v to %v", amount, sender, receiver)
	return tx, nil
}

func (t *Token) Allowance(owner common.Address, spender common.Address) (*big.Int, error) {
	txCallOptions := t.GetTxCallOptions()

	allowance, err := t.Interface.Allowance(txCallOptions, owner, spender)
	if err != nil {
		return nil, errors.New("Error retrieving allowance")
	}

	return allowance, nil
}

func (t *Token) Approve(spender common.Address, amount *big.Int) (*eth.Transaction, error) {
	opts, _ := t.GetTxSendOptions()

	tx, err := t.Interface.Approve(opts, spender, amount)
	if err != nil {
		return nil, errors.New("Error making Approve() transaction")
	}

	return tx, nil
}

func (t *Token) ApproveFrom(w *types.Wallet, spender common.Address, amount *big.Int) (*eth.Transaction, error) {
	opts := t.GetCustomTxSendOptions(w)

	tx, err := t.Interface.Approve(opts, spender, amount)
	if err != nil {
		return nil, errors.New("Error making ApproveFrom() transaction")
	}

	return tx, nil
}

func (t *Token) ListenToTransferEvents() (chan *contractsinterfaces.ERC20Transfer, error) {
	events := make(chan *contractsinterfaces.ERC20Transfer)
	options := &bind.WatchOpts{nil, nil}
	toList := []common.Address{}
	fromList := []common.Address{}

	_, err := t.Interface.WatchTransfer(options, events, fromList, toList)
	if err != nil {
		return nil, err
	}

	return events, nil
}

func (t *Token) PrintTransferEvents() error {
	events := make(chan *contractsinterfaces.ERC20Transfer)
	options := &bind.WatchOpts{nil, nil}
	toList := []common.Address{}
	fromList := []common.Address{}

	_, err := t.Interface.WatchTransfer(options, events, fromList, toList)
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
