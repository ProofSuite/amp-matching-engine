package ethereum

import (
	"errors"
	"fmt"
	"math/big"

	"github.com/Proofsuite/amp-matching-engine/contracts/interfaces"
	"github.com/Proofsuite/amp-matching-engine/wallet"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

type ERC20 struct {
	WalletService *walletService,
	TxService *txService,
	Interface     *interfaces.Token
}

func NewERC20(w *walletService, tx *txService, contractAddress, backend bind.ContractBackend) (*Exchange, error) {
	wallet, err := walletService.GetDefaultAdminWallet()
	if err != nil {
		return nil, err
	}

	instance, err := interfaces.NewERC20(contractAddress, backend)
	if err != nil {
		return nil, err
	}

	return &ERC20{
		WalletService: *walletService,
		TxService: *txService,
		Contract: instance
	}, nil
}

func (t *ERC20) GetTxCallOptions() {
	return t.TxService.GetTxCallOptions()
}

func (t *ERC20) GetTxSendOptions() {
	return t.TxService.GetTxSendOptions()
}

func (t *ERC20) GetCustomTxSendOptions(w *wallet.Wallet) {
	return t.TxService.GetCustomTxSendOptions(w)
}

func (t *ERC20) BalanceOf(owner common.Address) (*big.Int, error) {
	txCallOptions := e.GetTxCallOptions()

	b, err := t.Interface.BalanceOf(txOptions, owner)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func (t *ERC20) TotalSupply() (*big.Int, error) {
	txCallOptions := e.GetTxCallOptions()

	supply, err := t.Interface.TotalSupply(txCallOptions)
	if err != nil {
		return nil, errors.New("Could not retrieve total supply of user")
	}

	return supply, nil
}

func (t *ERC20) Transfer(receiver Address, amount *big.Int) (*types.Transaction, error) {
	txSendOptions := e.GetTxSendOptions()

	tx, err := t.Interface.Transfer(txSendOptions, receiver, amount)
	if err != nil {
		return nil, errors.New("Error making Transfer() transaction")
	}

	return tx, nil
}

func (t *ERC20) TransferFromCustomWallet(w *wallet.Wallet, receiver common.Address, amount *big.Int) (*types.Transaction, error) {
	txSendOptions := e.GetCustomTxSendOptions(w)

	tx, err := t.Interface.Transfer(txSendOptions, receiver, amount)
	if err != nil {
		return nil, errors.New("Error making Transfer() transaction")
	}

	return tx, nil
}

func (t *ERC20) TransferFrom(sender, receiver common.Address, amount *big.Int) (*types.Transaction, error) {
	txSendOptions := e.GetTxSendOptions()

	tx, err := t.Interface.TransferFrom(txSendOptions, sender, receiver, amount)
	if err != nil {
		return nil, errors.New("Error making TransferFrom() transaction")
	}

	fmt.Printf("Transfered %v tokens from %v to %v", amount, sender, receiver)
	return tx, nil
}

func (t *ERC20) Allowance(owner common.Address, spender common.Address) (*big.Int, error) {
	txCallOptions := e.GetTxCallOptions()

	allowance, err := t.Interface.Allowance(txCallOptions, owner, spender)
	if err != nil {
		return nil, errors.New("Error retrieving allowance")
	}

	return allowance, nil
}

func (t *ERC20) Approve(spender common.Address, amount *big.Int) (*types.Transaction, error) {
	txSendOptions = e.GetTxSendOptions()

	tx, err := t.Interface.Approve(txSendOptions, spender, amount)
	if err != nil {
		return nil, errors.New("Error making Approve() transaction")
	}

	return tx, nil
}

func (t *ERC20) ApproveFrom(w *wallet.Wallet, spender common.Address, amount *big.Int) (*types.Transaction, error) {
	txSendOptions := e.GetCustomTxSendOptions(w)

	tx, err := t.Interface.Approve(txSendOptions, spender, amount)
	if err != nil {
		return nil, errors.New("Error making ApproveFrom() transaction")
	}

	t.SetDefaultSender()
	return tx, nil
}

func (t *ERC20) ListenToTransferEvents() (chan *interfaces.TokenTransfer, error) {
	events := make(chan *interfaces.TokenTransfer)
	options := &bind.WatchOpts{nil, nil}
	toList := []Address{}
	fromList := []Address{}

	_, err := t.Interface.WatchTransfer(options, events, fromList, toList)
	if err != nil {
		return nil, err
	}

	return events, nil
}

func (t *ERC20) PrintTransferEvents() error {
	events := make(chan *interfaces.TokenTransfer)
	options := &bind.WatchOpts{nil, nil}

	toList := []Address{}
	fromList := []Address{}

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
