// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package contractsinterfaces

import (
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// ERC20ABI is the input ABI used to generate the binding from.
const ERC20ABI = "[{\"constant\":true,\"inputs\":[],\"name\":\"name\",\"outputs\":[{\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_spender\",\"type\":\"address\"},{\"name\":\"_value\",\"type\":\"uint256\"}],\"name\":\"approve\",\"outputs\":[{\"name\":\"success\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"totalSupply\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"totalTokenSupply\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_from\",\"type\":\"address\"},{\"name\":\"_to\",\"type\":\"address\"},{\"name\":\"_value\",\"type\":\"uint256\"}],\"name\":\"transferFrom\",\"outputs\":[{\"name\":\"success\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"decimals\",\"outputs\":[{\"name\":\"\",\"type\":\"uint8\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_owner\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"name\":\"balance\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"symbol\",\"outputs\":[{\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_to\",\"type\":\"address\"},{\"name\":\"_value\",\"type\":\"uint256\"}],\"name\":\"transfer\",\"outputs\":[{\"name\":\"success\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_owner\",\"type\":\"address\"},{\"name\":\"_spender\",\"type\":\"address\"}],\"name\":\"allowance\",\"outputs\":[{\"name\":\"remaining\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"_from\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"_to\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"_value\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"_owner\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"_spender\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"_value\",\"type\":\"uint256\"}],\"name\":\"Approval\",\"type\":\"event\"}]"

// ERC20Bin is the compiled bytecode used for deploying new contracts.
const ERC20Bin = `0x`

// DeployERC20 deploys a new Ethereum contract, binding an instance of ERC20 to it.
func DeployERC20(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *ERC20, error) {
	parsed, err := abi.JSON(strings.NewReader(ERC20ABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(ERC20Bin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &ERC20{ERC20Caller: ERC20Caller{contract: contract}, ERC20Transactor: ERC20Transactor{contract: contract}, ERC20Filterer: ERC20Filterer{contract: contract}}, nil
}

// ERC20 is an auto generated Go binding around an Ethereum contract.
type ERC20 struct {
	ERC20Caller     // Read-only binding to the contract
	ERC20Transactor // Write-only binding to the contract
	ERC20Filterer   // Log filterer for contract events
}

// ERC20Caller is an auto generated read-only Go binding around an Ethereum contract.
type ERC20Caller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ERC20Transactor is an auto generated write-only Go binding around an Ethereum contract.
type ERC20Transactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ERC20Filterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ERC20Filterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ERC20Session is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ERC20Session struct {
	Contract     *ERC20            // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// ERC20CallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ERC20CallerSession struct {
	Contract *ERC20Caller  // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts // Call options to use throughout this session
}

// ERC20TransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ERC20TransactorSession struct {
	Contract     *ERC20Transactor  // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// ERC20Raw is an auto generated low-level Go binding around an Ethereum contract.
type ERC20Raw struct {
	Contract *ERC20 // Generic contract binding to access the raw methods on
}

// ERC20CallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ERC20CallerRaw struct {
	Contract *ERC20Caller // Generic read-only contract binding to access the raw methods on
}

// ERC20TransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ERC20TransactorRaw struct {
	Contract *ERC20Transactor // Generic write-only contract binding to access the raw methods on
}

// NewERC20 creates a new instance of ERC20, bound to a specific deployed contract.
func NewERC20(address common.Address, backend bind.ContractBackend) (*ERC20, error) {
	contract, err := bindERC20(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ERC20{ERC20Caller: ERC20Caller{contract: contract}, ERC20Transactor: ERC20Transactor{contract: contract}, ERC20Filterer: ERC20Filterer{contract: contract}}, nil
}

// NewERC20Caller creates a new read-only instance of ERC20, bound to a specific deployed contract.
func NewERC20Caller(address common.Address, caller bind.ContractCaller) (*ERC20Caller, error) {
	contract, err := bindERC20(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ERC20Caller{contract: contract}, nil
}

// NewERC20Transactor creates a new write-only instance of ERC20, bound to a specific deployed contract.
func NewERC20Transactor(address common.Address, transactor bind.ContractTransactor) (*ERC20Transactor, error) {
	contract, err := bindERC20(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ERC20Transactor{contract: contract}, nil
}

// NewERC20Filterer creates a new log filterer instance of ERC20, bound to a specific deployed contract.
func NewERC20Filterer(address common.Address, filterer bind.ContractFilterer) (*ERC20Filterer, error) {
	contract, err := bindERC20(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ERC20Filterer{contract: contract}, nil
}

// bindERC20 binds a generic wrapper to an already deployed contract.
func bindERC20(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(ERC20ABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ERC20 *ERC20Raw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _ERC20.Contract.ERC20Caller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ERC20 *ERC20Raw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ERC20.Contract.ERC20Transactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ERC20 *ERC20Raw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ERC20.Contract.ERC20Transactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ERC20 *ERC20CallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _ERC20.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ERC20 *ERC20TransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ERC20.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ERC20 *ERC20TransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ERC20.Contract.contract.Transact(opts, method, params...)
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(_owner address, _spender address) constant returns(remaining uint256)
func (_ERC20 *ERC20Caller) Allowance(opts *bind.CallOpts, _owner common.Address, _spender common.Address) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _ERC20.contract.Call(opts, out, "allowance", _owner, _spender)
	return *ret0, err
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(_owner address, _spender address) constant returns(remaining uint256)
func (_ERC20 *ERC20Session) Allowance(_owner common.Address, _spender common.Address) (*big.Int, error) {
	return _ERC20.Contract.Allowance(&_ERC20.CallOpts, _owner, _spender)
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(_owner address, _spender address) constant returns(remaining uint256)
func (_ERC20 *ERC20CallerSession) Allowance(_owner common.Address, _spender common.Address) (*big.Int, error) {
	return _ERC20.Contract.Allowance(&_ERC20.CallOpts, _owner, _spender)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(_owner address) constant returns(balance uint256)
func (_ERC20 *ERC20Caller) BalanceOf(opts *bind.CallOpts, _owner common.Address) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _ERC20.contract.Call(opts, out, "balanceOf", _owner)
	return *ret0, err
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(_owner address) constant returns(balance uint256)
func (_ERC20 *ERC20Session) BalanceOf(_owner common.Address) (*big.Int, error) {
	return _ERC20.Contract.BalanceOf(&_ERC20.CallOpts, _owner)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(_owner address) constant returns(balance uint256)
func (_ERC20 *ERC20CallerSession) BalanceOf(_owner common.Address) (*big.Int, error) {
	return _ERC20.Contract.BalanceOf(&_ERC20.CallOpts, _owner)
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() constant returns(uint8)
func (_ERC20 *ERC20Caller) Decimals(opts *bind.CallOpts) (uint8, error) {
	var (
		ret0 = new(uint8)
	)
	out := ret0
	err := _ERC20.contract.Call(opts, out, "decimals")
	return *ret0, err
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() constant returns(uint8)
func (_ERC20 *ERC20Session) Decimals() (uint8, error) {
	return _ERC20.Contract.Decimals(&_ERC20.CallOpts)
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() constant returns(uint8)
func (_ERC20 *ERC20CallerSession) Decimals() (uint8, error) {
	return _ERC20.Contract.Decimals(&_ERC20.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() constant returns(string)
func (_ERC20 *ERC20Caller) Name(opts *bind.CallOpts) (string, error) {
	var (
		ret0 = new(string)
	)
	out := ret0
	err := _ERC20.contract.Call(opts, out, "name")
	return *ret0, err
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() constant returns(string)
func (_ERC20 *ERC20Session) Name() (string, error) {
	return _ERC20.Contract.Name(&_ERC20.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() constant returns(string)
func (_ERC20 *ERC20CallerSession) Name() (string, error) {
	return _ERC20.Contract.Name(&_ERC20.CallOpts)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() constant returns(string)
func (_ERC20 *ERC20Caller) Symbol(opts *bind.CallOpts) (string, error) {
	var (
		ret0 = new(string)
	)
	out := ret0
	err := _ERC20.contract.Call(opts, out, "symbol")
	return *ret0, err
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() constant returns(string)
func (_ERC20 *ERC20Session) Symbol() (string, error) {
	return _ERC20.Contract.Symbol(&_ERC20.CallOpts)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() constant returns(string)
func (_ERC20 *ERC20CallerSession) Symbol() (string, error) {
	return _ERC20.Contract.Symbol(&_ERC20.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() constant returns(uint256)
func (_ERC20 *ERC20Caller) TotalSupply(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _ERC20.contract.Call(opts, out, "totalSupply")
	return *ret0, err
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() constant returns(uint256)
func (_ERC20 *ERC20Session) TotalSupply() (*big.Int, error) {
	return _ERC20.Contract.TotalSupply(&_ERC20.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() constant returns(uint256)
func (_ERC20 *ERC20CallerSession) TotalSupply() (*big.Int, error) {
	return _ERC20.Contract.TotalSupply(&_ERC20.CallOpts)
}

// TotalTokenSupply is a free data retrieval call binding the contract method 0x1ca8b6cb.
//
// Solidity: function totalTokenSupply() constant returns(uint256)
func (_ERC20 *ERC20Caller) TotalTokenSupply(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _ERC20.contract.Call(opts, out, "totalTokenSupply")
	return *ret0, err
}

// TotalTokenSupply is a free data retrieval call binding the contract method 0x1ca8b6cb.
//
// Solidity: function totalTokenSupply() constant returns(uint256)
func (_ERC20 *ERC20Session) TotalTokenSupply() (*big.Int, error) {
	return _ERC20.Contract.TotalTokenSupply(&_ERC20.CallOpts)
}

// TotalTokenSupply is a free data retrieval call binding the contract method 0x1ca8b6cb.
//
// Solidity: function totalTokenSupply() constant returns(uint256)
func (_ERC20 *ERC20CallerSession) TotalTokenSupply() (*big.Int, error) {
	return _ERC20.Contract.TotalTokenSupply(&_ERC20.CallOpts)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(_spender address, _value uint256) returns(success bool)
func (_ERC20 *ERC20Transactor) Approve(opts *bind.TransactOpts, _spender common.Address, _value *big.Int) (*types.Transaction, error) {
	return _ERC20.contract.Transact(opts, "approve", _spender, _value)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(_spender address, _value uint256) returns(success bool)
func (_ERC20 *ERC20Session) Approve(_spender common.Address, _value *big.Int) (*types.Transaction, error) {
	return _ERC20.Contract.Approve(&_ERC20.TransactOpts, _spender, _value)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(_spender address, _value uint256) returns(success bool)
func (_ERC20 *ERC20TransactorSession) Approve(_spender common.Address, _value *big.Int) (*types.Transaction, error) {
	return _ERC20.Contract.Approve(&_ERC20.TransactOpts, _spender, _value)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(_to address, _value uint256) returns(success bool)
func (_ERC20 *ERC20Transactor) Transfer(opts *bind.TransactOpts, _to common.Address, _value *big.Int) (*types.Transaction, error) {
	return _ERC20.contract.Transact(opts, "transfer", _to, _value)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(_to address, _value uint256) returns(success bool)
func (_ERC20 *ERC20Session) Transfer(_to common.Address, _value *big.Int) (*types.Transaction, error) {
	return _ERC20.Contract.Transfer(&_ERC20.TransactOpts, _to, _value)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(_to address, _value uint256) returns(success bool)
func (_ERC20 *ERC20TransactorSession) Transfer(_to common.Address, _value *big.Int) (*types.Transaction, error) {
	return _ERC20.Contract.Transfer(&_ERC20.TransactOpts, _to, _value)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(_from address, _to address, _value uint256) returns(success bool)
func (_ERC20 *ERC20Transactor) TransferFrom(opts *bind.TransactOpts, _from common.Address, _to common.Address, _value *big.Int) (*types.Transaction, error) {
	return _ERC20.contract.Transact(opts, "transferFrom", _from, _to, _value)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(_from address, _to address, _value uint256) returns(success bool)
func (_ERC20 *ERC20Session) TransferFrom(_from common.Address, _to common.Address, _value *big.Int) (*types.Transaction, error) {
	return _ERC20.Contract.TransferFrom(&_ERC20.TransactOpts, _from, _to, _value)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(_from address, _to address, _value uint256) returns(success bool)
func (_ERC20 *ERC20TransactorSession) TransferFrom(_from common.Address, _to common.Address, _value *big.Int) (*types.Transaction, error) {
	return _ERC20.Contract.TransferFrom(&_ERC20.TransactOpts, _from, _to, _value)
}

// ERC20ApprovalIterator is returned from FilterApproval and is used to iterate over the raw logs and unpacked data for Approval events raised by the ERC20 contract.
type ERC20ApprovalIterator struct {
	Event *ERC20Approval // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *ERC20ApprovalIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ERC20Approval)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(ERC20Approval)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *ERC20ApprovalIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ERC20ApprovalIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ERC20Approval represents a Approval event raised by the ERC20 contract.
type ERC20Approval struct {
	Owner   common.Address
	Spender common.Address
	Value   *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterApproval is a free log retrieval operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: e Approval(_owner indexed address, _spender indexed address, _value uint256)
func (_ERC20 *ERC20Filterer) FilterApproval(opts *bind.FilterOpts, _owner []common.Address, _spender []common.Address) (*ERC20ApprovalIterator, error) {

	var _ownerRule []interface{}
	for _, _ownerItem := range _owner {
		_ownerRule = append(_ownerRule, _ownerItem)
	}
	var _spenderRule []interface{}
	for _, _spenderItem := range _spender {
		_spenderRule = append(_spenderRule, _spenderItem)
	}

	logs, sub, err := _ERC20.contract.FilterLogs(opts, "Approval", _ownerRule, _spenderRule)
	if err != nil {
		return nil, err
	}
	return &ERC20ApprovalIterator{contract: _ERC20.contract, event: "Approval", logs: logs, sub: sub}, nil
}

// WatchApproval is a free log subscription operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: e Approval(_owner indexed address, _spender indexed address, _value uint256)
func (_ERC20 *ERC20Filterer) WatchApproval(opts *bind.WatchOpts, sink chan<- *ERC20Approval, _owner []common.Address, _spender []common.Address) (event.Subscription, error) {

	var _ownerRule []interface{}
	for _, _ownerItem := range _owner {
		_ownerRule = append(_ownerRule, _ownerItem)
	}
	var _spenderRule []interface{}
	for _, _spenderItem := range _spender {
		_spenderRule = append(_spenderRule, _spenderItem)
	}

	logs, sub, err := _ERC20.contract.WatchLogs(opts, "Approval", _ownerRule, _spenderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ERC20Approval)
				if err := _ERC20.contract.UnpackLog(event, "Approval", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ERC20TransferIterator is returned from FilterTransfer and is used to iterate over the raw logs and unpacked data for Transfer events raised by the ERC20 contract.
type ERC20TransferIterator struct {
	Event *ERC20Transfer // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *ERC20TransferIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ERC20Transfer)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(ERC20Transfer)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *ERC20TransferIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ERC20TransferIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ERC20Transfer represents a Transfer event raised by the ERC20 contract.
type ERC20Transfer struct {
	From  common.Address
	To    common.Address
	Value *big.Int
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterTransfer is a free log retrieval operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: e Transfer(_from indexed address, _to indexed address, _value uint256)
func (_ERC20 *ERC20Filterer) FilterTransfer(opts *bind.FilterOpts, _from []common.Address, _to []common.Address) (*ERC20TransferIterator, error) {

	var _fromRule []interface{}
	for _, _fromItem := range _from {
		_fromRule = append(_fromRule, _fromItem)
	}
	var _toRule []interface{}
	for _, _toItem := range _to {
		_toRule = append(_toRule, _toItem)
	}

	logs, sub, err := _ERC20.contract.FilterLogs(opts, "Transfer", _fromRule, _toRule)
	if err != nil {
		return nil, err
	}
	return &ERC20TransferIterator{contract: _ERC20.contract, event: "Transfer", logs: logs, sub: sub}, nil
}

// WatchTransfer is a free log subscription operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: e Transfer(_from indexed address, _to indexed address, _value uint256)
func (_ERC20 *ERC20Filterer) WatchTransfer(opts *bind.WatchOpts, sink chan<- *ERC20Transfer, _from []common.Address, _to []common.Address) (event.Subscription, error) {

	var _fromRule []interface{}
	for _, _fromItem := range _from {
		_fromRule = append(_fromRule, _fromItem)
	}
	var _toRule []interface{}
	for _, _toItem := range _to {
		_toRule = append(_toRule, _toItem)
	}

	logs, sub, err := _ERC20.contract.WatchLogs(opts, "Transfer", _fromRule, _toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ERC20Transfer)
				if err := _ERC20.contract.UnpackLog(event, "Transfer", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ExchangeABI is the input ABI used to generate the binding from.
const ExchangeABI = "[{\"constant\":false,\"inputs\":[{\"name\":\"orderValues\",\"type\":\"uint256[8][]\"},{\"name\":\"orderAddresses\",\"type\":\"address[4][]\"},{\"name\":\"v\",\"type\":\"uint8[2][]\"},{\"name\":\"rs\",\"type\":\"bytes32[4][]\"}],\"name\":\"executeBatchTrades\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"setOwner\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"address\"}],\"name\":\"operators\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"numerator\",\"type\":\"uint256\"},{\"name\":\"denominator\",\"type\":\"uint256\"},{\"name\":\"target\",\"type\":\"uint256\"}],\"name\":\"isRoundingError\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"orderHash\",\"type\":\"bytes32[]\"},{\"name\":\"amount\",\"type\":\"uint256[]\"},{\"name\":\"tradeNonce\",\"type\":\"uint256[]\"},{\"name\":\"taker\",\"type\":\"address[]\"},{\"name\":\"v\",\"type\":\"uint8[]\"},{\"name\":\"r\",\"type\":\"bytes32[]\"},{\"name\":\"s\",\"type\":\"bytes32[]\"}],\"name\":\"batchCancelTrades\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"orderValues\",\"type\":\"uint256[8]\"},{\"name\":\"orderAddresses\",\"type\":\"address[4]\"},{\"name\":\"v\",\"type\":\"uint8[2]\"},{\"name\":\"rs\",\"type\":\"bytes32[4]\"}],\"name\":\"executeTrade\",\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\"},{\"name\":\"\",\"type\":\"bytes32\"},{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"filled\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"orderHash\",\"type\":\"bytes32\"},{\"name\":\"amount\",\"type\":\"uint256\"},{\"name\":\"tradeNonce\",\"type\":\"uint256\"},{\"name\":\"taker\",\"type\":\"address\"},{\"name\":\"v\",\"type\":\"uint8\"},{\"name\":\"r\",\"type\":\"bytes32\"},{\"name\":\"s\",\"type\":\"bytes32\"}],\"name\":\"cancelTrade\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_feeAccount\",\"type\":\"address\"}],\"name\":\"setFeeAccount\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"wethToken\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_operator\",\"type\":\"address\"},{\"name\":\"_isOperator\",\"type\":\"bool\"}],\"name\":\"setOperator\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"feeAccount\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"signer\",\"type\":\"address\"},{\"name\":\"hash\",\"type\":\"bytes32\"},{\"name\":\"v\",\"type\":\"uint8\"},{\"name\":\"r\",\"type\":\"bytes32\"},{\"name\":\"s\",\"type\":\"bytes32\"}],\"name\":\"isValidSignature\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_wethToken\",\"type\":\"address\"}],\"name\":\"setWethToken\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"numerator\",\"type\":\"uint256\"},{\"name\":\"denominator\",\"type\":\"uint256\"},{\"name\":\"target\",\"type\":\"uint256\"}],\"name\":\"getPartialAmount\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"traded\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"orderValues\",\"type\":\"uint256[6]\"},{\"name\":\"orderAddresses\",\"type\":\"address[3]\"},{\"name\":\"v\",\"type\":\"uint8\"},{\"name\":\"r\",\"type\":\"bytes32\"},{\"name\":\"s\",\"type\":\"bytes32\"}],\"name\":\"cancelOrder\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"orderValues\",\"type\":\"uint256[6][]\"},{\"name\":\"orderAddresses\",\"type\":\"address[3][]\"},{\"name\":\"v\",\"type\":\"uint8[]\"},{\"name\":\"r\",\"type\":\"bytes32[]\"},{\"name\":\"s\",\"type\":\"bytes32[]\"}],\"name\":\"batchCancelOrders\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"VERSION\",\"outputs\":[{\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"name\":\"_wethToken\",\"type\":\"address\"},{\"name\":\"_feeAccount\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"oldWethToken\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"newWethToken\",\"type\":\"address\"}],\"name\":\"LogWethTokenUpdate\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"oldFeeAccount\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"newFeeAccount\",\"type\":\"address\"}],\"name\":\"LogFeeAccountUpdate\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"operator\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"isOperator\",\"type\":\"bool\"}],\"name\":\"LogOperatorUpdate\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"taker\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"tokenSell\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"tokenBuy\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"orderHashes\",\"type\":\"bytes32[]\"},{\"indexed\":false,\"name\":\"tradeHashes\",\"type\":\"bytes32[]\"},{\"indexed\":false,\"name\":\"filledAmounts\",\"type\":\"uint256[]\"},{\"indexed\":true,\"name\":\"tokenPairHash\",\"type\":\"bytes32\"}],\"name\":\"LogBatchTrades\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"maker\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"taker\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"tokenSell\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"tokenBuy\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"filledAmountSell\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"filledAmountBuy\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"paidFeeMake\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"paidFeeTake\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"orderHash\",\"type\":\"bytes32\"},{\"indexed\":false,\"name\":\"tradeHash\",\"type\":\"bytes32\"},{\"indexed\":true,\"name\":\"tokenPairHash\",\"type\":\"bytes32\"}],\"name\":\"LogTrade\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"errorId\",\"type\":\"uint8\"},{\"indexed\":false,\"name\":\"orderHash\",\"type\":\"bytes32\"},{\"indexed\":false,\"name\":\"tradeHash\",\"type\":\"bytes32\"}],\"name\":\"LogError\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"orderHash\",\"type\":\"bytes32\"},{\"indexed\":false,\"name\":\"tokenBuy\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"amountBuy\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"tokenSell\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"amountSell\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"expires\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"nonce\",\"type\":\"uint256\"},{\"indexed\":true,\"name\":\"maker\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"tokenPairHash\",\"type\":\"bytes32\"}],\"name\":\"LogCancelOrder\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"orderHash\",\"type\":\"bytes32\"},{\"indexed\":false,\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"tradeNonce\",\"type\":\"uint256\"},{\"indexed\":true,\"name\":\"taker\",\"type\":\"address\"}],\"name\":\"LogCancelTrade\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"SetOwner\",\"type\":\"event\"}]"

// ExchangeBin is the compiled bytecode used for deploying new contracts.
const ExchangeBin = `0x608060405234801561001057600080fd5b5060405160408061218183398101604052805160209091015160008054600160a060020a0319908116331790915560018054600160a060020a039485169083161790556002805493909216921691909117905561210f806100726000396000f3006080604052600436106101115763ffffffff7c01000000000000000000000000000000000000000000000000000000006000350416630877f6c1811461011657806313af40351461030057806313e7c9d81461032357806314df96ee14610344578063154e775a146103625780632207148d1461050d578063288cdc91146105df578063468ddf2e146106095780634b023cf81461063f5780634b57b0be14610660578063558a72971461069157806365e17c9d146106b75780638163681e146106cc57806386e09c08146106fc5780638da5cb5b1461071d57806398024a8b14610732578063d581332314610750578063d9a72b5214610768578063e51ad32d146107df578063ffa1ad741461098b575b600080fd5b34801561012257600080fd5b506040805160048035808201356020818102850181019095528084526102ec943694602493909290840191819060009085015b8282101561019257604080516101008181019092529080840287019060089083908390808284375050509183525050600190910190602001610155565b50506040805186358801803560208181028401810190945280835296999897830196919550820193509150819060009085015b828210156102015760408051608081810190925290808402870190600490839083908082843750505091835250506001909101906020016101c5565b50506040805186358801803560208181028401810190945280835296999897830196919550820193509150819060009085015b8282101561026d576040805180820182529080840287019060029083908390808284375050509183525050600190910190602001610234565b50506040805186358801803560208181028401810190945280835296999897830196919550820193509150819060009085015b828210156102dc5760408051608081810190925290808402870190600490839083908082843750505091835250506001909101906020016102a0565b50939650610a1595505050505050565b604080519115158252519081900360200190f35b34801561030c57600080fd5b50610321600160a060020a0360043516610e7c565b005b34801561032f57600080fd5b506102ec600160a060020a0360043516610efb565b34801561035057600080fd5b506102ec600435602435604435610f10565b34801561036e57600080fd5b506040805160206004803580820135838102808601850190965280855261032195369593946024949385019291829185019084908082843750506040805187358901803560208181028481018201909552818452989b9a998901989297509082019550935083925085019084908082843750506040805187358901803560208181028481018201909552818452989b9a998901989297509082019550935083925085019084908082843750506040805187358901803560208181028481018201909552818452989b9a998901989297509082019550935083925085019084908082843750506040805187358901803560208181028481018201909552818452989b9a998901989297509082019550935083925085019084908082843750506040805187358901803560208181028481018201909552818452989b9a998901989297509082019550935083925085019084908082843750506040805187358901803560208181028481018201909552818452989b9a998901989297509082019550935083925085019084908082843750949750610f799650505050505050565b34801561051957600080fd5b50604080516101008181019092526105c19136916004916101049190839060089083908390808284375050604080516080818101909252949796958181019594509250600491508390839080828437505060408051808201825294979695818101959450925060029150839083908082843750506040805160808181019092529497969581810195945092506004915083908390808284375093965061104895505050505050565b60408051938452602084019290925282820152519081900360600190f35b3480156105eb57600080fd5b506105f7600435611680565b60408051918252519081900360200190f35b34801561061557600080fd5b506102ec600435602435604435600160a060020a036064351660ff6084351660a43560c435611692565b34801561064b57600080fd5b506102ec600160a060020a036004351661178f565b34801561066c57600080fd5b50610675611836565b60408051600160a060020a039092168252519081900360200190f35b34801561069d57600080fd5b506102ec600160a060020a03600435166024351515611845565b3480156106c357600080fd5b506106756118e6565b3480156106d857600080fd5b506102ec600160a060020a036004351660243560ff604435166064356084356118f5565b34801561070857600080fd5b506102ec600160a060020a0360043516611a1d565b34801561072957600080fd5b50610675611aad565b34801561073e57600080fd5b506105f7600435602435604435611abc565b34801561075c57600080fd5b506102ec600435611ada565b34801561077457600080fd5b506040805160c08181019092526102ec91369160049160c49190839060069083908390808284375050604080516060818101909252949796958181019594509250600391508390839080828437509396505050823560ff169350505060208101359060400135611aef565b3480156107eb57600080fd5b50604080516004803580820135602081810285018101909552808452610321943694602493909290840191819060009085015b8282101561085a576040805160c0818101909252908084028701906006908390839080828437505050918352505060019091019060200161081e565b50506040805186358801803560208181028401810190945280835296999897830196919550820193509150819060009085015b828210156108c957604080516060818101909252908084028701906003908390839080828437505050918352505060019091019060200161088d565b50505050509192919290803590602001908201803590602001908080602002602001604051908101604052809392919081815260200183836020028082843750506040805187358901803560208181028481018201909552818452989b9a998901989297509082019550935083925085019084908082843750506040805187358901803560208181028481018201909552818452989b9a998901989297509082019550935083925085019084908082843750949750611d3f9650505050505050565b34801561099757600080fd5b506109a0611ddc565b6040805160208082528351818301528351919283929083019185019080838360005b838110156109da5781810151838201526020016109c2565b50505050905090810190601f168015610a075780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b60008054606090819081908490819081908190600160a060020a0316331480610a4d57503360009081526003602052604090205460ff165b1515610a5857600080fd5b60018b5103604051908082528060200260200182016040528015610a86578160200160208202803883390190505b50965060018b5103604051908082528060200260200182016040528015610ab7578160200160208202803883390190505b50955060018b5103604051908082528060200260200182016040528015610ae8578160200160208202803883390190505b509450600093505b8a51841015610bd157610b618c85815181101515610b0a57fe5b906020019060200201518c86815181101515610b2257fe5b906020019060200201518c87815181101515610b3a57fe5b906020019060200201518c88815181101515610b5257fe5b90602001906020020151611048565b9250925092508060001415610b795760009750610e6d565b828785815181101515610b8857fe5b6020908102909101015285518290879086908110610ba257fe5b6020908102909101015284518190869086908110610bbc57fe5b60209081029091010152600190930192610af0565b8a6000815181101515610be057fe5b60209081029190910181015101518b518c906000908110610bfd57fe5b60209081029091010151600260200201516040516020018083600160a060020a0316600160a060020a03166c0100000000000000000000000002815260140182600160a060020a0316600160a060020a03166c01000000000000000000000000028152601401925050506040516020818303038152906040526040518082805190602001908083835b60208310610ca55780518252601f199092019160209182019101610c86565b5181516020939093036101000a600019018019909116921691909117905260405192018290039091208e519093508e9250600091508110610ce257fe5b6020908102909101015160026020020151600160a060020a03167f94bbd7e15e68655b645dc3781ae12285dee432ad5c825161d78d6f69c20c19358d6000815181101515610d2c57fe5b60209081029190910181015101518e518f906000908110610d4957fe5b60209081029091010151600060200201518b8b8b6040518086600160a060020a0316600160a060020a0316815260200185600160a060020a0316600160a060020a03168152602001806020018060200180602001848103845287818151815260200191508051906020019060200280838360005b83811015610dd5578181015183820152602001610dbd565b50505050905001848103835286818151815260200191508051906020019060200280838360005b83811015610e14578181015183820152602001610dfc565b50505050905001848103825285818151815260200191508051906020019060200280838360005b83811015610e53578181015183820152602001610e3b565b505050509050019850505050505050505060405180910390a35b50505050505050949350505050565b600054600160a060020a03163314610e9357600080fd5b60008054604051600160a060020a03808516939216917fcbf985117192c8f614a58aaf97226bb80a754772f5f6edf06f87c675f2e6c66391a36000805473ffffffffffffffffffffffffffffffffffffffff1916600160a060020a0392909216919091179055565b60036020526000908152604090205460ff1681565b600080600084801515610f1f57fe5b8685099150811515610f345760009250610f70565b610f66610f47878663ffffffff611e1316565b610f5a84620f424063ffffffff611e1316565b9063ffffffff611e3e16565b90506103e8811192505b50509392505050565b60005b875181101561103e576110358882815181101515610f9657fe5b906020019060200201518883815181101515610fae57fe5b906020019060200201518884815181101515610fc657fe5b906020019060200201518885815181101515610fde57fe5b906020019060200201518886815181101515610ff657fe5b90602001906020020151888781518110151561100e57fe5b90602001906020020151888881518110151561102657fe5b90602001906020020151611692565b50600101610f7c565b5050505050505050565b6000806000611055612014565b600061105f61207c565b60008054819081908190600160a060020a031633148061108e57503360009081526003602052604090205460ff165b151561109957600080fd5b604080516101208101909152808f6000602090810291909101518252018f6001602090810291909101518252018f6002602090810291909101518252018f6003602090810291909101518252018f6004602090810291909101518252018f6005602090810291909101518252018e600060209081029190910151600160a060020a03168252018e600160209081029190910151600160a060020a03168252018e60026020020151600160a060020a03169052965061115687611e55565b6040805160808101909152818152909650602081018f6006602090810291909101518252018f6007602090810291909101518252018e60036020020151600160a060020a0316905294506111a985611f64565b6101008801518d518d519296506111cb928991908f60015b60200201516118f5565b1515611217576000805160206120a483398151915260015b6040805160ff909216825260208201899052818101879052519081900360600190a18584600080905099509950995061166f565b606085015160208d015160408d0151611235929187918f60036111c1565b1515611251576000805160206120a483398151915260026111e3565b4387604001511015611273576000805160206120a483398151915260036111e3565b60008481526005602052604090205460ff16156112a0576000805160206120a483398151915260046111e3565b8651602080870151600089815260049092526040909120546112c79163ffffffff61200516565b11156112e3576000805160206120a483398151915260056111e3565b6112fa856020015188600001518960200151610f10565b15611315576000805160206120a483398151915260066111e3565b6000848152600560209081526040909120805460ff19166001179055858101518851918901516113459290611abc565b6020808701516000898152600490925260409091205491945061136e919063ffffffff61200516565b60008781526004602081815260408084209490945560e08b01516101008c015160608b015186516000805160206120c48339815191528152600160a060020a0392831695810195909552811660248501526044840189905294519416936323b872dd936064808501948390030190829087803b1580156113ed57600080fd5b505af1158015611401573d6000803e3d6000fd5b505050506040513d602081101561141757600080fd5b5051151561142457600080fd5b60c08701516060860151610100890151602080890151604080516000805160206120c48339815191528152600160a060020a039586166004820152938516602485015260448401919091525192909316926323b872dd926064808401938290030181600087803b15801561149757600080fd5b505af11580156114ab573d6000803e3d6000fd5b505050506040513d60208110156114c157600080fd5b505115156114ce57600080fd5b600087608001511115611598576114f2856020015188600001518960800151611abc565b600154610100890151600254604080516000805160206120c48339815191528152600160a060020a039384166004820152918316602483015260448201859052519395509116916323b872dd916064808201926020929091908290030181600087803b15801561156157600080fd5b505af1158015611575573d6000803e3d6000fd5b505050506040513d602081101561158b57600080fd5b5051151561159857600080fd5b60008760a001511115611661576115bc856020015188600001518960a00151611abc565b6001546060870151600254604080516000805160206120c48339815191528152600160a060020a039384166004820152918316602483015260448201859052519394509116916323b872dd916064808201926020929091908290030181600087803b15801561162a57600080fd5b505af115801561163e573d6000803e3d6000fd5b505050506040513d602081101561165457600080fd5b5051151561166157600080fd5b858486602001519950995099505b505050505050509450945094915050565b60046020526000908152604090205481565b600061169c61207c565b506040805160808101825289815260208101899052908101879052600160a060020a038616606082015260006116d182611f64565b90506116e033828888886118f5565b151561171757604080516000815280820183905290516000805160206120a48339815191529181900360600190a160009250611782565b600081815260056020908152604091829020805460ff1916600117905581518c81529081018b90528082018a90529051600160a060020a038916917f1debd637af55cac936fd656ab3fb0391eb4eb29cb178bf44577ef6cecc10ae25919081900360600190a2600192505b5050979650505050505050565b60008054600160a060020a031633146117a757600080fd5b600160a060020a03821615156117bc57600080fd5b60025460408051600160a060020a039283168152918416602083015280517ff822f5a19627202340985855aeffadb385833332f2f700b3e6287d28547778a99281900390910190a15060028054600160a060020a03831673ffffffffffffffffffffffffffffffffffffffff199091161790556001919050565b600154600160a060020a031681565b60008054600160a060020a0316331461185d57600080fd5b600160a060020a038316151561187257600080fd5b60408051600160a060020a0385168152831515602082015281517f4af650e9ee9ac50b37ec2cd3ddac7e1c69955ffc871bc3e812563775f3bc0e7d929181900390910190a150600160a060020a0382166000908152600360205260409020805482151560ff19909116179055600192915050565b600254600160a060020a031681565b600060018560405160200180807f19457468657265756d205369676e6564204d6573736167653a0a333200000000815250601c0182600019166000191681526020019150506040516020818303038152906040526040518082805190602001908083835b602083106119785780518252601f199092019160209182019101611959565b51815160209384036101000a60001901801990921691161790526040805192909401829003822060008084528383018087529190915260ff8c1683860152606083018b9052608083018a9052935160a08084019750919550601f1981019492819003909101925090865af11580156119f4573d6000803e3d6000fd5b50505060206040510351600160a060020a031686600160a060020a031614905095945050505050565b60008054600160a060020a03163314611a3557600080fd5b60015460408051600160a060020a039283168152918416602083015280517fb8be72b4c168c2f7d3ea469d9f48ccbc62416784a4f6a69ca93ff13f4f36545b9281900390910190a15060018054600160a060020a03831673ffffffffffffffffffffffffffffffffffffffff19909116178155919050565b600054600160a060020a031681565b6000611ad283610f5a868563ffffffff611e1316565b949350505050565b60056020526000908152604090205460ff1681565b6000611af9612014565b506040805161012081018252875181526020808901518183015288830151828401526060808a0151908301526080808a01519083015260a0808a0151908301528751600160a060020a0390811660c084015290880151811660e0830152918701519091166101008201526000611b6e82611e55565b9050611b7d33828888886118f5565b1515611bb55760408051600081526020810183905290516000805160206120a48339815191529181900360600190a160009250611d34565b81516000828152600460209081526040918290209290925560e084015160c085015182516c01000000000000000000000000600160a060020a0393841681028287015292909116909102603482015281516028818303018152604890910191829052805190928291908401908083835b60208310611c445780518252601f199092019160209182019101611c25565b6001836020036101000a038019825116818451168082178552505050505050905001915050604051809103902060001916826101000151600160a060020a03167fbfa78175e8dfd3bfda20dd3ae584843e6ec42822f51f553fc12f5d9f908fdb16838560c0015186600001518760e00151886020015189604001518a6060015160405180886000191660001916815260200187600160a060020a0316600160a060020a0316815260200186815260200185600160a060020a0316600160a060020a0316815260200184815260200183815260200182815260200197505050505050505060405180910390a3600192505b505095945050505050565b60005b8451811015611dd457611dcb8682815181101515611d5c57fe5b906020019060200201518683815181101515611d7457fe5b906020019060200201518684815181101515611d8c57fe5b906020019060200201518685815181101515611da457fe5b906020019060200201518686815181101515611dbc57fe5b90602001906020020151611aef565b50600101611d42565b505050505050565b60408051808201909152600581527f312e302e30000000000000000000000000000000000000000000000000000000602082015281565b6000828202831580611e2f5750828482811515611e2c57fe5b04145b1515611e3757fe5b9392505050565b6000808284811515611e4c57fe5b04949350505050565b61010081015160e082015160c08301516020808501518551608087015160a08801516040808a015160608b015182516c01000000000000000000000000308102828b0152600160a060020a039c8d16810260348301529a8c168b0260488201529a909816909802605c8a01526070890194909452609088019290925260b087015260d086015260f085019390935261011080850192909252825180850390920182526101309093019182905280516000939192918291908401908083835b60208310611f325780518252601f199092019160209182019101611f13565b5181516020939093036101000a6000190180199091169216919091179052604051920182900390912095945050505050565b6000816000015182606001518360200151846040015160405160200180856000191660001916815260200184600160a060020a0316600160a060020a03166c0100000000000000000000000002815260140183815260200182815260200194505050505060405160208183030381529060405260405180828051906020019080838360208310611f325780518252601f199092019160209182019101611f13565b600082820183811015611e3757fe5b610120604051908101604052806000815260200160008152602001600081526020016000815260200160008152602001600081526020016000600160a060020a031681526020016000600160a060020a031681526020016000600160a060020a031681525090565b60408051608081018252600080825260208201819052918101829052606081019190915290560014301341d034ec3c62a1eabc804a79abf3b8c16e6245e82ec572346aa452fabb23b872dd00000000000000000000000000000000000000000000000000000000a165627a7a7230582040880eaa9300426385e3550a08a90a72e72d369db877ac773603a66201c5ba7d0029`

// DeployExchange deploys a new Ethereum contract, binding an instance of Exchange to it.
func DeployExchange(auth *bind.TransactOpts, backend bind.ContractBackend, _wethToken common.Address, _feeAccount common.Address) (common.Address, *types.Transaction, *Exchange, error) {
	parsed, err := abi.JSON(strings.NewReader(ExchangeABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(ExchangeBin), backend, _wethToken, _feeAccount)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Exchange{ExchangeCaller: ExchangeCaller{contract: contract}, ExchangeTransactor: ExchangeTransactor{contract: contract}, ExchangeFilterer: ExchangeFilterer{contract: contract}}, nil
}

// Exchange is an auto generated Go binding around an Ethereum contract.
type Exchange struct {
	ExchangeCaller     // Read-only binding to the contract
	ExchangeTransactor // Write-only binding to the contract
	ExchangeFilterer   // Log filterer for contract events
}

// ExchangeCaller is an auto generated read-only Go binding around an Ethereum contract.
type ExchangeCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ExchangeTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ExchangeTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ExchangeFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ExchangeFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ExchangeSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ExchangeSession struct {
	Contract     *Exchange         // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// ExchangeCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ExchangeCallerSession struct {
	Contract *ExchangeCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts   // Call options to use throughout this session
}

// ExchangeTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ExchangeTransactorSession struct {
	Contract     *ExchangeTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts   // Transaction auth options to use throughout this session
}

// ExchangeRaw is an auto generated low-level Go binding around an Ethereum contract.
type ExchangeRaw struct {
	Contract *Exchange // Generic contract binding to access the raw methods on
}

// ExchangeCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ExchangeCallerRaw struct {
	Contract *ExchangeCaller // Generic read-only contract binding to access the raw methods on
}

// ExchangeTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ExchangeTransactorRaw struct {
	Contract *ExchangeTransactor // Generic write-only contract binding to access the raw methods on
}

// NewExchange creates a new instance of Exchange, bound to a specific deployed contract.
func NewExchange(address common.Address, backend bind.ContractBackend) (*Exchange, error) {
	contract, err := bindExchange(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Exchange{ExchangeCaller: ExchangeCaller{contract: contract}, ExchangeTransactor: ExchangeTransactor{contract: contract}, ExchangeFilterer: ExchangeFilterer{contract: contract}}, nil
}

// NewExchangeCaller creates a new read-only instance of Exchange, bound to a specific deployed contract.
func NewExchangeCaller(address common.Address, caller bind.ContractCaller) (*ExchangeCaller, error) {
	contract, err := bindExchange(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ExchangeCaller{contract: contract}, nil
}

// NewExchangeTransactor creates a new write-only instance of Exchange, bound to a specific deployed contract.
func NewExchangeTransactor(address common.Address, transactor bind.ContractTransactor) (*ExchangeTransactor, error) {
	contract, err := bindExchange(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ExchangeTransactor{contract: contract}, nil
}

// NewExchangeFilterer creates a new log filterer instance of Exchange, bound to a specific deployed contract.
func NewExchangeFilterer(address common.Address, filterer bind.ContractFilterer) (*ExchangeFilterer, error) {
	contract, err := bindExchange(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ExchangeFilterer{contract: contract}, nil
}

// bindExchange binds a generic wrapper to an already deployed contract.
func bindExchange(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(ExchangeABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Exchange *ExchangeRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _Exchange.Contract.ExchangeCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Exchange *ExchangeRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Exchange.Contract.ExchangeTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Exchange *ExchangeRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Exchange.Contract.ExchangeTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Exchange *ExchangeCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _Exchange.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Exchange *ExchangeTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Exchange.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Exchange *ExchangeTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Exchange.Contract.contract.Transact(opts, method, params...)
}

// VERSION is a free data retrieval call binding the contract method 0xffa1ad74.
//
// Solidity: function VERSION() constant returns(string)
func (_Exchange *ExchangeCaller) VERSION(opts *bind.CallOpts) (string, error) {
	var (
		ret0 = new(string)
	)
	out := ret0
	err := _Exchange.contract.Call(opts, out, "VERSION")
	return *ret0, err
}

// VERSION is a free data retrieval call binding the contract method 0xffa1ad74.
//
// Solidity: function VERSION() constant returns(string)
func (_Exchange *ExchangeSession) VERSION() (string, error) {
	return _Exchange.Contract.VERSION(&_Exchange.CallOpts)
}

// VERSION is a free data retrieval call binding the contract method 0xffa1ad74.
//
// Solidity: function VERSION() constant returns(string)
func (_Exchange *ExchangeCallerSession) VERSION() (string, error) {
	return _Exchange.Contract.VERSION(&_Exchange.CallOpts)
}

// FeeAccount is a free data retrieval call binding the contract method 0x65e17c9d.
//
// Solidity: function feeAccount() constant returns(address)
func (_Exchange *ExchangeCaller) FeeAccount(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _Exchange.contract.Call(opts, out, "feeAccount")
	return *ret0, err
}

// FeeAccount is a free data retrieval call binding the contract method 0x65e17c9d.
//
// Solidity: function feeAccount() constant returns(address)
func (_Exchange *ExchangeSession) FeeAccount() (common.Address, error) {
	return _Exchange.Contract.FeeAccount(&_Exchange.CallOpts)
}

// FeeAccount is a free data retrieval call binding the contract method 0x65e17c9d.
//
// Solidity: function feeAccount() constant returns(address)
func (_Exchange *ExchangeCallerSession) FeeAccount() (common.Address, error) {
	return _Exchange.Contract.FeeAccount(&_Exchange.CallOpts)
}

// Filled is a free data retrieval call binding the contract method 0x288cdc91.
//
// Solidity: function filled( bytes32) constant returns(uint256)
func (_Exchange *ExchangeCaller) Filled(opts *bind.CallOpts, arg0 [32]byte) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Exchange.contract.Call(opts, out, "filled", arg0)
	return *ret0, err
}

// Filled is a free data retrieval call binding the contract method 0x288cdc91.
//
// Solidity: function filled( bytes32) constant returns(uint256)
func (_Exchange *ExchangeSession) Filled(arg0 [32]byte) (*big.Int, error) {
	return _Exchange.Contract.Filled(&_Exchange.CallOpts, arg0)
}

// Filled is a free data retrieval call binding the contract method 0x288cdc91.
//
// Solidity: function filled( bytes32) constant returns(uint256)
func (_Exchange *ExchangeCallerSession) Filled(arg0 [32]byte) (*big.Int, error) {
	return _Exchange.Contract.Filled(&_Exchange.CallOpts, arg0)
}

// GetPartialAmount is a free data retrieval call binding the contract method 0x98024a8b.
//
// Solidity: function getPartialAmount(numerator uint256, denominator uint256, target uint256) constant returns(uint256)
func (_Exchange *ExchangeCaller) GetPartialAmount(opts *bind.CallOpts, numerator *big.Int, denominator *big.Int, target *big.Int) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Exchange.contract.Call(opts, out, "getPartialAmount", numerator, denominator, target)
	return *ret0, err
}

// GetPartialAmount is a free data retrieval call binding the contract method 0x98024a8b.
//
// Solidity: function getPartialAmount(numerator uint256, denominator uint256, target uint256) constant returns(uint256)
func (_Exchange *ExchangeSession) GetPartialAmount(numerator *big.Int, denominator *big.Int, target *big.Int) (*big.Int, error) {
	return _Exchange.Contract.GetPartialAmount(&_Exchange.CallOpts, numerator, denominator, target)
}

// GetPartialAmount is a free data retrieval call binding the contract method 0x98024a8b.
//
// Solidity: function getPartialAmount(numerator uint256, denominator uint256, target uint256) constant returns(uint256)
func (_Exchange *ExchangeCallerSession) GetPartialAmount(numerator *big.Int, denominator *big.Int, target *big.Int) (*big.Int, error) {
	return _Exchange.Contract.GetPartialAmount(&_Exchange.CallOpts, numerator, denominator, target)
}

// IsRoundingError is a free data retrieval call binding the contract method 0x14df96ee.
//
// Solidity: function isRoundingError(numerator uint256, denominator uint256, target uint256) constant returns(bool)
func (_Exchange *ExchangeCaller) IsRoundingError(opts *bind.CallOpts, numerator *big.Int, denominator *big.Int, target *big.Int) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _Exchange.contract.Call(opts, out, "isRoundingError", numerator, denominator, target)
	return *ret0, err
}

// IsRoundingError is a free data retrieval call binding the contract method 0x14df96ee.
//
// Solidity: function isRoundingError(numerator uint256, denominator uint256, target uint256) constant returns(bool)
func (_Exchange *ExchangeSession) IsRoundingError(numerator *big.Int, denominator *big.Int, target *big.Int) (bool, error) {
	return _Exchange.Contract.IsRoundingError(&_Exchange.CallOpts, numerator, denominator, target)
}

// IsRoundingError is a free data retrieval call binding the contract method 0x14df96ee.
//
// Solidity: function isRoundingError(numerator uint256, denominator uint256, target uint256) constant returns(bool)
func (_Exchange *ExchangeCallerSession) IsRoundingError(numerator *big.Int, denominator *big.Int, target *big.Int) (bool, error) {
	return _Exchange.Contract.IsRoundingError(&_Exchange.CallOpts, numerator, denominator, target)
}

// IsValidSignature is a free data retrieval call binding the contract method 0x8163681e.
//
// Solidity: function isValidSignature(signer address, hash bytes32, v uint8, r bytes32, s bytes32) constant returns(bool)
func (_Exchange *ExchangeCaller) IsValidSignature(opts *bind.CallOpts, signer common.Address, hash [32]byte, v uint8, r [32]byte, s [32]byte) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _Exchange.contract.Call(opts, out, "isValidSignature", signer, hash, v, r, s)
	return *ret0, err
}

// IsValidSignature is a free data retrieval call binding the contract method 0x8163681e.
//
// Solidity: function isValidSignature(signer address, hash bytes32, v uint8, r bytes32, s bytes32) constant returns(bool)
func (_Exchange *ExchangeSession) IsValidSignature(signer common.Address, hash [32]byte, v uint8, r [32]byte, s [32]byte) (bool, error) {
	return _Exchange.Contract.IsValidSignature(&_Exchange.CallOpts, signer, hash, v, r, s)
}

// IsValidSignature is a free data retrieval call binding the contract method 0x8163681e.
//
// Solidity: function isValidSignature(signer address, hash bytes32, v uint8, r bytes32, s bytes32) constant returns(bool)
func (_Exchange *ExchangeCallerSession) IsValidSignature(signer common.Address, hash [32]byte, v uint8, r [32]byte, s [32]byte) (bool, error) {
	return _Exchange.Contract.IsValidSignature(&_Exchange.CallOpts, signer, hash, v, r, s)
}

// Operators is a free data retrieval call binding the contract method 0x13e7c9d8.
//
// Solidity: function operators( address) constant returns(bool)
func (_Exchange *ExchangeCaller) Operators(opts *bind.CallOpts, arg0 common.Address) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _Exchange.contract.Call(opts, out, "operators", arg0)
	return *ret0, err
}

// Operators is a free data retrieval call binding the contract method 0x13e7c9d8.
//
// Solidity: function operators( address) constant returns(bool)
func (_Exchange *ExchangeSession) Operators(arg0 common.Address) (bool, error) {
	return _Exchange.Contract.Operators(&_Exchange.CallOpts, arg0)
}

// Operators is a free data retrieval call binding the contract method 0x13e7c9d8.
//
// Solidity: function operators( address) constant returns(bool)
func (_Exchange *ExchangeCallerSession) Operators(arg0 common.Address) (bool, error) {
	return _Exchange.Contract.Operators(&_Exchange.CallOpts, arg0)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() constant returns(address)
func (_Exchange *ExchangeCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _Exchange.contract.Call(opts, out, "owner")
	return *ret0, err
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() constant returns(address)
func (_Exchange *ExchangeSession) Owner() (common.Address, error) {
	return _Exchange.Contract.Owner(&_Exchange.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() constant returns(address)
func (_Exchange *ExchangeCallerSession) Owner() (common.Address, error) {
	return _Exchange.Contract.Owner(&_Exchange.CallOpts)
}

// Traded is a free data retrieval call binding the contract method 0xd5813323.
//
// Solidity: function traded( bytes32) constant returns(bool)
func (_Exchange *ExchangeCaller) Traded(opts *bind.CallOpts, arg0 [32]byte) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _Exchange.contract.Call(opts, out, "traded", arg0)
	return *ret0, err
}

// Traded is a free data retrieval call binding the contract method 0xd5813323.
//
// Solidity: function traded( bytes32) constant returns(bool)
func (_Exchange *ExchangeSession) Traded(arg0 [32]byte) (bool, error) {
	return _Exchange.Contract.Traded(&_Exchange.CallOpts, arg0)
}

// Traded is a free data retrieval call binding the contract method 0xd5813323.
//
// Solidity: function traded( bytes32) constant returns(bool)
func (_Exchange *ExchangeCallerSession) Traded(arg0 [32]byte) (bool, error) {
	return _Exchange.Contract.Traded(&_Exchange.CallOpts, arg0)
}

// WethToken is a free data retrieval call binding the contract method 0x4b57b0be.
//
// Solidity: function wethToken() constant returns(address)
func (_Exchange *ExchangeCaller) WethToken(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _Exchange.contract.Call(opts, out, "wethToken")
	return *ret0, err
}

// WethToken is a free data retrieval call binding the contract method 0x4b57b0be.
//
// Solidity: function wethToken() constant returns(address)
func (_Exchange *ExchangeSession) WethToken() (common.Address, error) {
	return _Exchange.Contract.WethToken(&_Exchange.CallOpts)
}

// WethToken is a free data retrieval call binding the contract method 0x4b57b0be.
//
// Solidity: function wethToken() constant returns(address)
func (_Exchange *ExchangeCallerSession) WethToken() (common.Address, error) {
	return _Exchange.Contract.WethToken(&_Exchange.CallOpts)
}

// BatchCancelOrders is a paid mutator transaction binding the contract method 0xe51ad32d.
//
// Solidity: function batchCancelOrders(orderValues uint256[6][], orderAddresses address[3][], v uint8[], r bytes32[], s bytes32[]) returns()
func (_Exchange *ExchangeTransactor) BatchCancelOrders(opts *bind.TransactOpts, orderValues [][6]*big.Int, orderAddresses [][3]common.Address, v []uint8, r [][32]byte, s [][32]byte) (*types.Transaction, error) {
	return _Exchange.contract.Transact(opts, "batchCancelOrders", orderValues, orderAddresses, v, r, s)
}

// BatchCancelOrders is a paid mutator transaction binding the contract method 0xe51ad32d.
//
// Solidity: function batchCancelOrders(orderValues uint256[6][], orderAddresses address[3][], v uint8[], r bytes32[], s bytes32[]) returns()
func (_Exchange *ExchangeSession) BatchCancelOrders(orderValues [][6]*big.Int, orderAddresses [][3]common.Address, v []uint8, r [][32]byte, s [][32]byte) (*types.Transaction, error) {
	return _Exchange.Contract.BatchCancelOrders(&_Exchange.TransactOpts, orderValues, orderAddresses, v, r, s)
}

// BatchCancelOrders is a paid mutator transaction binding the contract method 0xe51ad32d.
//
// Solidity: function batchCancelOrders(orderValues uint256[6][], orderAddresses address[3][], v uint8[], r bytes32[], s bytes32[]) returns()
func (_Exchange *ExchangeTransactorSession) BatchCancelOrders(orderValues [][6]*big.Int, orderAddresses [][3]common.Address, v []uint8, r [][32]byte, s [][32]byte) (*types.Transaction, error) {
	return _Exchange.Contract.BatchCancelOrders(&_Exchange.TransactOpts, orderValues, orderAddresses, v, r, s)
}

// BatchCancelTrades is a paid mutator transaction binding the contract method 0x154e775a.
//
// Solidity: function batchCancelTrades(orderHash bytes32[], amount uint256[], tradeNonce uint256[], taker address[], v uint8[], r bytes32[], s bytes32[]) returns()
func (_Exchange *ExchangeTransactor) BatchCancelTrades(opts *bind.TransactOpts, orderHash [][32]byte, amount []*big.Int, tradeNonce []*big.Int, taker []common.Address, v []uint8, r [][32]byte, s [][32]byte) (*types.Transaction, error) {
	return _Exchange.contract.Transact(opts, "batchCancelTrades", orderHash, amount, tradeNonce, taker, v, r, s)
}

// BatchCancelTrades is a paid mutator transaction binding the contract method 0x154e775a.
//
// Solidity: function batchCancelTrades(orderHash bytes32[], amount uint256[], tradeNonce uint256[], taker address[], v uint8[], r bytes32[], s bytes32[]) returns()
func (_Exchange *ExchangeSession) BatchCancelTrades(orderHash [][32]byte, amount []*big.Int, tradeNonce []*big.Int, taker []common.Address, v []uint8, r [][32]byte, s [][32]byte) (*types.Transaction, error) {
	return _Exchange.Contract.BatchCancelTrades(&_Exchange.TransactOpts, orderHash, amount, tradeNonce, taker, v, r, s)
}

// BatchCancelTrades is a paid mutator transaction binding the contract method 0x154e775a.
//
// Solidity: function batchCancelTrades(orderHash bytes32[], amount uint256[], tradeNonce uint256[], taker address[], v uint8[], r bytes32[], s bytes32[]) returns()
func (_Exchange *ExchangeTransactorSession) BatchCancelTrades(orderHash [][32]byte, amount []*big.Int, tradeNonce []*big.Int, taker []common.Address, v []uint8, r [][32]byte, s [][32]byte) (*types.Transaction, error) {
	return _Exchange.Contract.BatchCancelTrades(&_Exchange.TransactOpts, orderHash, amount, tradeNonce, taker, v, r, s)
}

// CancelOrder is a paid mutator transaction binding the contract method 0xd9a72b52.
//
// Solidity: function cancelOrder(orderValues uint256[6], orderAddresses address[3], v uint8, r bytes32, s bytes32) returns(bool)
func (_Exchange *ExchangeTransactor) CancelOrder(opts *bind.TransactOpts, orderValues [6]*big.Int, orderAddresses [3]common.Address, v uint8, r [32]byte, s [32]byte) (*types.Transaction, error) {
	return _Exchange.contract.Transact(opts, "cancelOrder", orderValues, orderAddresses, v, r, s)
}

// CancelOrder is a paid mutator transaction binding the contract method 0xd9a72b52.
//
// Solidity: function cancelOrder(orderValues uint256[6], orderAddresses address[3], v uint8, r bytes32, s bytes32) returns(bool)
func (_Exchange *ExchangeSession) CancelOrder(orderValues [6]*big.Int, orderAddresses [3]common.Address, v uint8, r [32]byte, s [32]byte) (*types.Transaction, error) {
	return _Exchange.Contract.CancelOrder(&_Exchange.TransactOpts, orderValues, orderAddresses, v, r, s)
}

// CancelOrder is a paid mutator transaction binding the contract method 0xd9a72b52.
//
// Solidity: function cancelOrder(orderValues uint256[6], orderAddresses address[3], v uint8, r bytes32, s bytes32) returns(bool)
func (_Exchange *ExchangeTransactorSession) CancelOrder(orderValues [6]*big.Int, orderAddresses [3]common.Address, v uint8, r [32]byte, s [32]byte) (*types.Transaction, error) {
	return _Exchange.Contract.CancelOrder(&_Exchange.TransactOpts, orderValues, orderAddresses, v, r, s)
}

// CancelTrade is a paid mutator transaction binding the contract method 0x468ddf2e.
//
// Solidity: function cancelTrade(orderHash bytes32, amount uint256, tradeNonce uint256, taker address, v uint8, r bytes32, s bytes32) returns(bool)
func (_Exchange *ExchangeTransactor) CancelTrade(opts *bind.TransactOpts, orderHash [32]byte, amount *big.Int, tradeNonce *big.Int, taker common.Address, v uint8, r [32]byte, s [32]byte) (*types.Transaction, error) {
	return _Exchange.contract.Transact(opts, "cancelTrade", orderHash, amount, tradeNonce, taker, v, r, s)
}

// CancelTrade is a paid mutator transaction binding the contract method 0x468ddf2e.
//
// Solidity: function cancelTrade(orderHash bytes32, amount uint256, tradeNonce uint256, taker address, v uint8, r bytes32, s bytes32) returns(bool)
func (_Exchange *ExchangeSession) CancelTrade(orderHash [32]byte, amount *big.Int, tradeNonce *big.Int, taker common.Address, v uint8, r [32]byte, s [32]byte) (*types.Transaction, error) {
	return _Exchange.Contract.CancelTrade(&_Exchange.TransactOpts, orderHash, amount, tradeNonce, taker, v, r, s)
}

// CancelTrade is a paid mutator transaction binding the contract method 0x468ddf2e.
//
// Solidity: function cancelTrade(orderHash bytes32, amount uint256, tradeNonce uint256, taker address, v uint8, r bytes32, s bytes32) returns(bool)
func (_Exchange *ExchangeTransactorSession) CancelTrade(orderHash [32]byte, amount *big.Int, tradeNonce *big.Int, taker common.Address, v uint8, r [32]byte, s [32]byte) (*types.Transaction, error) {
	return _Exchange.Contract.CancelTrade(&_Exchange.TransactOpts, orderHash, amount, tradeNonce, taker, v, r, s)
}

// ExecuteBatchTrades is a paid mutator transaction binding the contract method 0x0877f6c1.
//
// Solidity: function executeBatchTrades(orderValues uint256[8][], orderAddresses address[4][], v uint8[2][], rs bytes32[4][]) returns(bool)
func (_Exchange *ExchangeTransactor) ExecuteBatchTrades(opts *bind.TransactOpts, orderValues [][8]*big.Int, orderAddresses [][4]common.Address, v [][2]uint8, rs [][4][32]byte) (*types.Transaction, error) {
	return _Exchange.contract.Transact(opts, "executeBatchTrades", orderValues, orderAddresses, v, rs)
}

// ExecuteBatchTrades is a paid mutator transaction binding the contract method 0x0877f6c1.
//
// Solidity: function executeBatchTrades(orderValues uint256[8][], orderAddresses address[4][], v uint8[2][], rs bytes32[4][]) returns(bool)
func (_Exchange *ExchangeSession) ExecuteBatchTrades(orderValues [][8]*big.Int, orderAddresses [][4]common.Address, v [][2]uint8, rs [][4][32]byte) (*types.Transaction, error) {
	return _Exchange.Contract.ExecuteBatchTrades(&_Exchange.TransactOpts, orderValues, orderAddresses, v, rs)
}

// ExecuteBatchTrades is a paid mutator transaction binding the contract method 0x0877f6c1.
//
// Solidity: function executeBatchTrades(orderValues uint256[8][], orderAddresses address[4][], v uint8[2][], rs bytes32[4][]) returns(bool)
func (_Exchange *ExchangeTransactorSession) ExecuteBatchTrades(orderValues [][8]*big.Int, orderAddresses [][4]common.Address, v [][2]uint8, rs [][4][32]byte) (*types.Transaction, error) {
	return _Exchange.Contract.ExecuteBatchTrades(&_Exchange.TransactOpts, orderValues, orderAddresses, v, rs)
}

// ExecuteTrade is a paid mutator transaction binding the contract method 0x2207148d.
//
// Solidity: function executeTrade(orderValues uint256[8], orderAddresses address[4], v uint8[2], rs bytes32[4]) returns(bytes32, bytes32, uint256)
func (_Exchange *ExchangeTransactor) ExecuteTrade(opts *bind.TransactOpts, orderValues [8]*big.Int, orderAddresses [4]common.Address, v [2]uint8, rs [4][32]byte) (*types.Transaction, error) {
	return _Exchange.contract.Transact(opts, "executeTrade", orderValues, orderAddresses, v, rs)
}

// ExecuteTrade is a paid mutator transaction binding the contract method 0x2207148d.
//
// Solidity: function executeTrade(orderValues uint256[8], orderAddresses address[4], v uint8[2], rs bytes32[4]) returns(bytes32, bytes32, uint256)
func (_Exchange *ExchangeSession) ExecuteTrade(orderValues [8]*big.Int, orderAddresses [4]common.Address, v [2]uint8, rs [4][32]byte) (*types.Transaction, error) {
	return _Exchange.Contract.ExecuteTrade(&_Exchange.TransactOpts, orderValues, orderAddresses, v, rs)
}

// ExecuteTrade is a paid mutator transaction binding the contract method 0x2207148d.
//
// Solidity: function executeTrade(orderValues uint256[8], orderAddresses address[4], v uint8[2], rs bytes32[4]) returns(bytes32, bytes32, uint256)
func (_Exchange *ExchangeTransactorSession) ExecuteTrade(orderValues [8]*big.Int, orderAddresses [4]common.Address, v [2]uint8, rs [4][32]byte) (*types.Transaction, error) {
	return _Exchange.Contract.ExecuteTrade(&_Exchange.TransactOpts, orderValues, orderAddresses, v, rs)
}

// SetFeeAccount is a paid mutator transaction binding the contract method 0x4b023cf8.
//
// Solidity: function setFeeAccount(_feeAccount address) returns(bool)
func (_Exchange *ExchangeTransactor) SetFeeAccount(opts *bind.TransactOpts, _feeAccount common.Address) (*types.Transaction, error) {
	return _Exchange.contract.Transact(opts, "setFeeAccount", _feeAccount)
}

// SetFeeAccount is a paid mutator transaction binding the contract method 0x4b023cf8.
//
// Solidity: function setFeeAccount(_feeAccount address) returns(bool)
func (_Exchange *ExchangeSession) SetFeeAccount(_feeAccount common.Address) (*types.Transaction, error) {
	return _Exchange.Contract.SetFeeAccount(&_Exchange.TransactOpts, _feeAccount)
}

// SetFeeAccount is a paid mutator transaction binding the contract method 0x4b023cf8.
//
// Solidity: function setFeeAccount(_feeAccount address) returns(bool)
func (_Exchange *ExchangeTransactorSession) SetFeeAccount(_feeAccount common.Address) (*types.Transaction, error) {
	return _Exchange.Contract.SetFeeAccount(&_Exchange.TransactOpts, _feeAccount)
}

// SetOperator is a paid mutator transaction binding the contract method 0x558a7297.
//
// Solidity: function setOperator(_operator address, _isOperator bool) returns(bool)
func (_Exchange *ExchangeTransactor) SetOperator(opts *bind.TransactOpts, _operator common.Address, _isOperator bool) (*types.Transaction, error) {
	return _Exchange.contract.Transact(opts, "setOperator", _operator, _isOperator)
}

// SetOperator is a paid mutator transaction binding the contract method 0x558a7297.
//
// Solidity: function setOperator(_operator address, _isOperator bool) returns(bool)
func (_Exchange *ExchangeSession) SetOperator(_operator common.Address, _isOperator bool) (*types.Transaction, error) {
	return _Exchange.Contract.SetOperator(&_Exchange.TransactOpts, _operator, _isOperator)
}

// SetOperator is a paid mutator transaction binding the contract method 0x558a7297.
//
// Solidity: function setOperator(_operator address, _isOperator bool) returns(bool)
func (_Exchange *ExchangeTransactorSession) SetOperator(_operator common.Address, _isOperator bool) (*types.Transaction, error) {
	return _Exchange.Contract.SetOperator(&_Exchange.TransactOpts, _operator, _isOperator)
}

// SetOwner is a paid mutator transaction binding the contract method 0x13af4035.
//
// Solidity: function setOwner(newOwner address) returns()
func (_Exchange *ExchangeTransactor) SetOwner(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _Exchange.contract.Transact(opts, "setOwner", newOwner)
}

// SetOwner is a paid mutator transaction binding the contract method 0x13af4035.
//
// Solidity: function setOwner(newOwner address) returns()
func (_Exchange *ExchangeSession) SetOwner(newOwner common.Address) (*types.Transaction, error) {
	return _Exchange.Contract.SetOwner(&_Exchange.TransactOpts, newOwner)
}

// SetOwner is a paid mutator transaction binding the contract method 0x13af4035.
//
// Solidity: function setOwner(newOwner address) returns()
func (_Exchange *ExchangeTransactorSession) SetOwner(newOwner common.Address) (*types.Transaction, error) {
	return _Exchange.Contract.SetOwner(&_Exchange.TransactOpts, newOwner)
}

// SetWethToken is a paid mutator transaction binding the contract method 0x86e09c08.
//
// Solidity: function setWethToken(_wethToken address) returns(bool)
func (_Exchange *ExchangeTransactor) SetWethToken(opts *bind.TransactOpts, _wethToken common.Address) (*types.Transaction, error) {
	return _Exchange.contract.Transact(opts, "setWethToken", _wethToken)
}

// SetWethToken is a paid mutator transaction binding the contract method 0x86e09c08.
//
// Solidity: function setWethToken(_wethToken address) returns(bool)
func (_Exchange *ExchangeSession) SetWethToken(_wethToken common.Address) (*types.Transaction, error) {
	return _Exchange.Contract.SetWethToken(&_Exchange.TransactOpts, _wethToken)
}

// SetWethToken is a paid mutator transaction binding the contract method 0x86e09c08.
//
// Solidity: function setWethToken(_wethToken address) returns(bool)
func (_Exchange *ExchangeTransactorSession) SetWethToken(_wethToken common.Address) (*types.Transaction, error) {
	return _Exchange.Contract.SetWethToken(&_Exchange.TransactOpts, _wethToken)
}

// ExchangeLogBatchTradesIterator is returned from FilterLogBatchTrades and is used to iterate over the raw logs and unpacked data for LogBatchTrades events raised by the Exchange contract.
type ExchangeLogBatchTradesIterator struct {
	Event *ExchangeLogBatchTrades // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *ExchangeLogBatchTradesIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ExchangeLogBatchTrades)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(ExchangeLogBatchTrades)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *ExchangeLogBatchTradesIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ExchangeLogBatchTradesIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ExchangeLogBatchTrades represents a LogBatchTrades event raised by the Exchange contract.
type ExchangeLogBatchTrades struct {
	Taker         common.Address
	TokenSell     common.Address
	TokenBuy      common.Address
	OrderHashes   [][32]byte
	TradeHashes   [][32]byte
	FilledAmounts []*big.Int
	TokenPairHash [32]byte
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterLogBatchTrades is a free log retrieval operation binding the contract event 0x94bbd7e15e68655b645dc3781ae12285dee432ad5c825161d78d6f69c20c1935.
//
// Solidity: e LogBatchTrades(taker indexed address, tokenSell address, tokenBuy address, orderHashes bytes32[], tradeHashes bytes32[], filledAmounts uint256[], tokenPairHash indexed bytes32)
func (_Exchange *ExchangeFilterer) FilterLogBatchTrades(opts *bind.FilterOpts, taker []common.Address, tokenPairHash [][32]byte) (*ExchangeLogBatchTradesIterator, error) {

	var takerRule []interface{}
	for _, takerItem := range taker {
		takerRule = append(takerRule, takerItem)
	}

	var tokenPairHashRule []interface{}
	for _, tokenPairHashItem := range tokenPairHash {
		tokenPairHashRule = append(tokenPairHashRule, tokenPairHashItem)
	}

	logs, sub, err := _Exchange.contract.FilterLogs(opts, "LogBatchTrades", takerRule, tokenPairHashRule)
	if err != nil {
		return nil, err
	}
	return &ExchangeLogBatchTradesIterator{contract: _Exchange.contract, event: "LogBatchTrades", logs: logs, sub: sub}, nil
}

// WatchLogBatchTrades is a free log subscription operation binding the contract event 0x94bbd7e15e68655b645dc3781ae12285dee432ad5c825161d78d6f69c20c1935.
//
// Solidity: e LogBatchTrades(taker indexed address, tokenSell address, tokenBuy address, orderHashes bytes32[], tradeHashes bytes32[], filledAmounts uint256[], tokenPairHash indexed bytes32)
func (_Exchange *ExchangeFilterer) WatchLogBatchTrades(opts *bind.WatchOpts, sink chan<- *ExchangeLogBatchTrades, taker []common.Address, tokenPairHash [][32]byte) (event.Subscription, error) {

	var takerRule []interface{}
	for _, takerItem := range taker {
		takerRule = append(takerRule, takerItem)
	}

	var tokenPairHashRule []interface{}
	for _, tokenPairHashItem := range tokenPairHash {
		tokenPairHashRule = append(tokenPairHashRule, tokenPairHashItem)
	}

	logs, sub, err := _Exchange.contract.WatchLogs(opts, "LogBatchTrades", takerRule, tokenPairHashRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ExchangeLogBatchTrades)
				if err := _Exchange.contract.UnpackLog(event, "LogBatchTrades", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ExchangeLogCancelOrderIterator is returned from FilterLogCancelOrder and is used to iterate over the raw logs and unpacked data for LogCancelOrder events raised by the Exchange contract.
type ExchangeLogCancelOrderIterator struct {
	Event *ExchangeLogCancelOrder // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *ExchangeLogCancelOrderIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ExchangeLogCancelOrder)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(ExchangeLogCancelOrder)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *ExchangeLogCancelOrderIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ExchangeLogCancelOrderIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ExchangeLogCancelOrder represents a LogCancelOrder event raised by the Exchange contract.
type ExchangeLogCancelOrder struct {
	OrderHash     [32]byte
	TokenBuy      common.Address
	AmountBuy     *big.Int
	TokenSell     common.Address
	AmountSell    *big.Int
	Expires       *big.Int
	Nonce         *big.Int
	Maker         common.Address
	TokenPairHash [32]byte
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterLogCancelOrder is a free log retrieval operation binding the contract event 0xbfa78175e8dfd3bfda20dd3ae584843e6ec42822f51f553fc12f5d9f908fdb16.
//
// Solidity: e LogCancelOrder(orderHash bytes32, tokenBuy address, amountBuy uint256, tokenSell address, amountSell uint256, expires uint256, nonce uint256, maker indexed address, tokenPairHash indexed bytes32)
func (_Exchange *ExchangeFilterer) FilterLogCancelOrder(opts *bind.FilterOpts, maker []common.Address, tokenPairHash [][32]byte) (*ExchangeLogCancelOrderIterator, error) {

	var makerRule []interface{}
	for _, makerItem := range maker {
		makerRule = append(makerRule, makerItem)
	}
	var tokenPairHashRule []interface{}
	for _, tokenPairHashItem := range tokenPairHash {
		tokenPairHashRule = append(tokenPairHashRule, tokenPairHashItem)
	}

	logs, sub, err := _Exchange.contract.FilterLogs(opts, "LogCancelOrder", makerRule, tokenPairHashRule)
	if err != nil {
		return nil, err
	}
	return &ExchangeLogCancelOrderIterator{contract: _Exchange.contract, event: "LogCancelOrder", logs: logs, sub: sub}, nil
}

// WatchLogCancelOrder is a free log subscription operation binding the contract event 0xbfa78175e8dfd3bfda20dd3ae584843e6ec42822f51f553fc12f5d9f908fdb16.
//
// Solidity: e LogCancelOrder(orderHash bytes32, tokenBuy address, amountBuy uint256, tokenSell address, amountSell uint256, expires uint256, nonce uint256, maker indexed address, tokenPairHash indexed bytes32)
func (_Exchange *ExchangeFilterer) WatchLogCancelOrder(opts *bind.WatchOpts, sink chan<- *ExchangeLogCancelOrder, maker []common.Address, tokenPairHash [][32]byte) (event.Subscription, error) {

	var makerRule []interface{}
	for _, makerItem := range maker {
		makerRule = append(makerRule, makerItem)
	}
	var tokenPairHashRule []interface{}
	for _, tokenPairHashItem := range tokenPairHash {
		tokenPairHashRule = append(tokenPairHashRule, tokenPairHashItem)
	}

	logs, sub, err := _Exchange.contract.WatchLogs(opts, "LogCancelOrder", makerRule, tokenPairHashRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ExchangeLogCancelOrder)
				if err := _Exchange.contract.UnpackLog(event, "LogCancelOrder", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ExchangeLogCancelTradeIterator is returned from FilterLogCancelTrade and is used to iterate over the raw logs and unpacked data for LogCancelTrade events raised by the Exchange contract.
type ExchangeLogCancelTradeIterator struct {
	Event *ExchangeLogCancelTrade // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *ExchangeLogCancelTradeIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ExchangeLogCancelTrade)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(ExchangeLogCancelTrade)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *ExchangeLogCancelTradeIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ExchangeLogCancelTradeIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ExchangeLogCancelTrade represents a LogCancelTrade event raised by the Exchange contract.
type ExchangeLogCancelTrade struct {
	OrderHash  [32]byte
	Amount     *big.Int
	TradeNonce *big.Int
	Taker      common.Address
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterLogCancelTrade is a free log retrieval operation binding the contract event 0x1debd637af55cac936fd656ab3fb0391eb4eb29cb178bf44577ef6cecc10ae25.
//
// Solidity: e LogCancelTrade(orderHash bytes32, amount uint256, tradeNonce uint256, taker indexed address)
func (_Exchange *ExchangeFilterer) FilterLogCancelTrade(opts *bind.FilterOpts, taker []common.Address) (*ExchangeLogCancelTradeIterator, error) {

	var takerRule []interface{}
	for _, takerItem := range taker {
		takerRule = append(takerRule, takerItem)
	}

	logs, sub, err := _Exchange.contract.FilterLogs(opts, "LogCancelTrade", takerRule)
	if err != nil {
		return nil, err
	}
	return &ExchangeLogCancelTradeIterator{contract: _Exchange.contract, event: "LogCancelTrade", logs: logs, sub: sub}, nil
}

// WatchLogCancelTrade is a free log subscription operation binding the contract event 0x1debd637af55cac936fd656ab3fb0391eb4eb29cb178bf44577ef6cecc10ae25.
//
// Solidity: e LogCancelTrade(orderHash bytes32, amount uint256, tradeNonce uint256, taker indexed address)
func (_Exchange *ExchangeFilterer) WatchLogCancelTrade(opts *bind.WatchOpts, sink chan<- *ExchangeLogCancelTrade, taker []common.Address) (event.Subscription, error) {

	var takerRule []interface{}
	for _, takerItem := range taker {
		takerRule = append(takerRule, takerItem)
	}

	logs, sub, err := _Exchange.contract.WatchLogs(opts, "LogCancelTrade", takerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ExchangeLogCancelTrade)
				if err := _Exchange.contract.UnpackLog(event, "LogCancelTrade", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ExchangeLogErrorIterator is returned from FilterLogError and is used to iterate over the raw logs and unpacked data for LogError events raised by the Exchange contract.
type ExchangeLogErrorIterator struct {
	Event *ExchangeLogError // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *ExchangeLogErrorIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ExchangeLogError)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(ExchangeLogError)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *ExchangeLogErrorIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ExchangeLogErrorIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ExchangeLogError represents a LogError event raised by the Exchange contract.
type ExchangeLogError struct {
	ErrorId   uint8
	OrderHash [32]byte
	TradeHash [32]byte
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterLogError is a free log retrieval operation binding the contract event 0x14301341d034ec3c62a1eabc804a79abf3b8c16e6245e82ec572346aa452fabb.
//
// Solidity: e LogError(errorId uint8, orderHash bytes32, tradeHash bytes32)
func (_Exchange *ExchangeFilterer) FilterLogError(opts *bind.FilterOpts) (*ExchangeLogErrorIterator, error) {

	logs, sub, err := _Exchange.contract.FilterLogs(opts, "LogError")
	if err != nil {
		return nil, err
	}
	return &ExchangeLogErrorIterator{contract: _Exchange.contract, event: "LogError", logs: logs, sub: sub}, nil
}

// WatchLogError is a free log subscription operation binding the contract event 0x14301341d034ec3c62a1eabc804a79abf3b8c16e6245e82ec572346aa452fabb.
//
// Solidity: e LogError(errorId uint8, orderHash bytes32, tradeHash bytes32)
func (_Exchange *ExchangeFilterer) WatchLogError(opts *bind.WatchOpts, sink chan<- *ExchangeLogError) (event.Subscription, error) {

	logs, sub, err := _Exchange.contract.WatchLogs(opts, "LogError")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ExchangeLogError)
				if err := _Exchange.contract.UnpackLog(event, "LogError", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ExchangeLogFeeAccountUpdateIterator is returned from FilterLogFeeAccountUpdate and is used to iterate over the raw logs and unpacked data for LogFeeAccountUpdate events raised by the Exchange contract.
type ExchangeLogFeeAccountUpdateIterator struct {
	Event *ExchangeLogFeeAccountUpdate // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *ExchangeLogFeeAccountUpdateIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ExchangeLogFeeAccountUpdate)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(ExchangeLogFeeAccountUpdate)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *ExchangeLogFeeAccountUpdateIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ExchangeLogFeeAccountUpdateIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ExchangeLogFeeAccountUpdate represents a LogFeeAccountUpdate event raised by the Exchange contract.
type ExchangeLogFeeAccountUpdate struct {
	OldFeeAccount common.Address
	NewFeeAccount common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterLogFeeAccountUpdate is a free log retrieval operation binding the contract event 0xf822f5a19627202340985855aeffadb385833332f2f700b3e6287d28547778a9.
//
// Solidity: e LogFeeAccountUpdate(oldFeeAccount address, newFeeAccount address)
func (_Exchange *ExchangeFilterer) FilterLogFeeAccountUpdate(opts *bind.FilterOpts) (*ExchangeLogFeeAccountUpdateIterator, error) {

	logs, sub, err := _Exchange.contract.FilterLogs(opts, "LogFeeAccountUpdate")
	if err != nil {
		return nil, err
	}
	return &ExchangeLogFeeAccountUpdateIterator{contract: _Exchange.contract, event: "LogFeeAccountUpdate", logs: logs, sub: sub}, nil
}

// WatchLogFeeAccountUpdate is a free log subscription operation binding the contract event 0xf822f5a19627202340985855aeffadb385833332f2f700b3e6287d28547778a9.
//
// Solidity: e LogFeeAccountUpdate(oldFeeAccount address, newFeeAccount address)
func (_Exchange *ExchangeFilterer) WatchLogFeeAccountUpdate(opts *bind.WatchOpts, sink chan<- *ExchangeLogFeeAccountUpdate) (event.Subscription, error) {

	logs, sub, err := _Exchange.contract.WatchLogs(opts, "LogFeeAccountUpdate")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ExchangeLogFeeAccountUpdate)
				if err := _Exchange.contract.UnpackLog(event, "LogFeeAccountUpdate", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ExchangeLogOperatorUpdateIterator is returned from FilterLogOperatorUpdate and is used to iterate over the raw logs and unpacked data for LogOperatorUpdate events raised by the Exchange contract.
type ExchangeLogOperatorUpdateIterator struct {
	Event *ExchangeLogOperatorUpdate // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *ExchangeLogOperatorUpdateIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ExchangeLogOperatorUpdate)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(ExchangeLogOperatorUpdate)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *ExchangeLogOperatorUpdateIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ExchangeLogOperatorUpdateIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ExchangeLogOperatorUpdate represents a LogOperatorUpdate event raised by the Exchange contract.
type ExchangeLogOperatorUpdate struct {
	Operator   common.Address
	IsOperator bool
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterLogOperatorUpdate is a free log retrieval operation binding the contract event 0x4af650e9ee9ac50b37ec2cd3ddac7e1c69955ffc871bc3e812563775f3bc0e7d.
//
// Solidity: e LogOperatorUpdate(operator address, isOperator bool)
func (_Exchange *ExchangeFilterer) FilterLogOperatorUpdate(opts *bind.FilterOpts) (*ExchangeLogOperatorUpdateIterator, error) {

	logs, sub, err := _Exchange.contract.FilterLogs(opts, "LogOperatorUpdate")
	if err != nil {
		return nil, err
	}
	return &ExchangeLogOperatorUpdateIterator{contract: _Exchange.contract, event: "LogOperatorUpdate", logs: logs, sub: sub}, nil
}

// WatchLogOperatorUpdate is a free log subscription operation binding the contract event 0x4af650e9ee9ac50b37ec2cd3ddac7e1c69955ffc871bc3e812563775f3bc0e7d.
//
// Solidity: e LogOperatorUpdate(operator address, isOperator bool)
func (_Exchange *ExchangeFilterer) WatchLogOperatorUpdate(opts *bind.WatchOpts, sink chan<- *ExchangeLogOperatorUpdate) (event.Subscription, error) {

	logs, sub, err := _Exchange.contract.WatchLogs(opts, "LogOperatorUpdate")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ExchangeLogOperatorUpdate)
				if err := _Exchange.contract.UnpackLog(event, "LogOperatorUpdate", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ExchangeLogTradeIterator is returned from FilterLogTrade and is used to iterate over the raw logs and unpacked data for LogTrade events raised by the Exchange contract.
type ExchangeLogTradeIterator struct {
	Event *ExchangeLogTrade // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *ExchangeLogTradeIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ExchangeLogTrade)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(ExchangeLogTrade)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *ExchangeLogTradeIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ExchangeLogTradeIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ExchangeLogTrade represents a LogTrade event raised by the Exchange contract.
type ExchangeLogTrade struct {
	Maker            common.Address
	Taker            common.Address
	TokenSell        common.Address
	TokenBuy         common.Address
	FilledAmountSell *big.Int
	FilledAmountBuy  *big.Int
	PaidFeeMake      *big.Int
	PaidFeeTake      *big.Int
	OrderHash        [32]byte
	TradeHash        [32]byte
	TokenPairHash    [32]byte
	Raw              types.Log // Blockchain specific contextual infos
}

// FilterLogTrade is a free log retrieval operation binding the contract event 0x174a42d8fdc3a48bf80a4e95ac4b280ef69189e4603105caac770bf9771357fc.
//
// Solidity: e LogTrade(maker indexed address, taker indexed address, tokenSell address, tokenBuy address, filledAmountSell uint256, filledAmountBuy uint256, paidFeeMake uint256, paidFeeTake uint256, orderHash bytes32, tradeHash bytes32, tokenPairHash indexed bytes32)
func (_Exchange *ExchangeFilterer) FilterLogTrade(opts *bind.FilterOpts, maker []common.Address, taker []common.Address, tokenPairHash [][32]byte) (*ExchangeLogTradeIterator, error) {

	var makerRule []interface{}
	for _, makerItem := range maker {
		makerRule = append(makerRule, makerItem)
	}
	var takerRule []interface{}
	for _, takerItem := range taker {
		takerRule = append(takerRule, takerItem)
	}

	var tokenPairHashRule []interface{}
	for _, tokenPairHashItem := range tokenPairHash {
		tokenPairHashRule = append(tokenPairHashRule, tokenPairHashItem)
	}

	logs, sub, err := _Exchange.contract.FilterLogs(opts, "LogTrade", makerRule, takerRule, tokenPairHashRule)
	if err != nil {
		return nil, err
	}
	return &ExchangeLogTradeIterator{contract: _Exchange.contract, event: "LogTrade", logs: logs, sub: sub}, nil
}

// WatchLogTrade is a free log subscription operation binding the contract event 0x174a42d8fdc3a48bf80a4e95ac4b280ef69189e4603105caac770bf9771357fc.
//
// Solidity: e LogTrade(maker indexed address, taker indexed address, tokenSell address, tokenBuy address, filledAmountSell uint256, filledAmountBuy uint256, paidFeeMake uint256, paidFeeTake uint256, orderHash bytes32, tradeHash bytes32, tokenPairHash indexed bytes32)
func (_Exchange *ExchangeFilterer) WatchLogTrade(opts *bind.WatchOpts, sink chan<- *ExchangeLogTrade, maker []common.Address, taker []common.Address, tokenPairHash [][32]byte) (event.Subscription, error) {

	var makerRule []interface{}
	for _, makerItem := range maker {
		makerRule = append(makerRule, makerItem)
	}
	var takerRule []interface{}
	for _, takerItem := range taker {
		takerRule = append(takerRule, takerItem)
	}

	var tokenPairHashRule []interface{}
	for _, tokenPairHashItem := range tokenPairHash {
		tokenPairHashRule = append(tokenPairHashRule, tokenPairHashItem)
	}

	logs, sub, err := _Exchange.contract.WatchLogs(opts, "LogTrade", makerRule, takerRule, tokenPairHashRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ExchangeLogTrade)
				if err := _Exchange.contract.UnpackLog(event, "LogTrade", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ExchangeLogWethTokenUpdateIterator is returned from FilterLogWethTokenUpdate and is used to iterate over the raw logs and unpacked data for LogWethTokenUpdate events raised by the Exchange contract.
type ExchangeLogWethTokenUpdateIterator struct {
	Event *ExchangeLogWethTokenUpdate // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *ExchangeLogWethTokenUpdateIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ExchangeLogWethTokenUpdate)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(ExchangeLogWethTokenUpdate)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *ExchangeLogWethTokenUpdateIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ExchangeLogWethTokenUpdateIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ExchangeLogWethTokenUpdate represents a LogWethTokenUpdate event raised by the Exchange contract.
type ExchangeLogWethTokenUpdate struct {
	OldWethToken common.Address
	NewWethToken common.Address
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterLogWethTokenUpdate is a free log retrieval operation binding the contract event 0xb8be72b4c168c2f7d3ea469d9f48ccbc62416784a4f6a69ca93ff13f4f36545b.
//
// Solidity: e LogWethTokenUpdate(oldWethToken address, newWethToken address)
func (_Exchange *ExchangeFilterer) FilterLogWethTokenUpdate(opts *bind.FilterOpts) (*ExchangeLogWethTokenUpdateIterator, error) {

	logs, sub, err := _Exchange.contract.FilterLogs(opts, "LogWethTokenUpdate")
	if err != nil {
		return nil, err
	}
	return &ExchangeLogWethTokenUpdateIterator{contract: _Exchange.contract, event: "LogWethTokenUpdate", logs: logs, sub: sub}, nil
}

// WatchLogWethTokenUpdate is a free log subscription operation binding the contract event 0xb8be72b4c168c2f7d3ea469d9f48ccbc62416784a4f6a69ca93ff13f4f36545b.
//
// Solidity: e LogWethTokenUpdate(oldWethToken address, newWethToken address)
func (_Exchange *ExchangeFilterer) WatchLogWethTokenUpdate(opts *bind.WatchOpts, sink chan<- *ExchangeLogWethTokenUpdate) (event.Subscription, error) {

	logs, sub, err := _Exchange.contract.WatchLogs(opts, "LogWethTokenUpdate")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ExchangeLogWethTokenUpdate)
				if err := _Exchange.contract.UnpackLog(event, "LogWethTokenUpdate", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ExchangeSetOwnerIterator is returned from FilterSetOwner and is used to iterate over the raw logs and unpacked data for SetOwner events raised by the Exchange contract.
type ExchangeSetOwnerIterator struct {
	Event *ExchangeSetOwner // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *ExchangeSetOwnerIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ExchangeSetOwner)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(ExchangeSetOwner)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *ExchangeSetOwnerIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ExchangeSetOwnerIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ExchangeSetOwner represents a SetOwner event raised by the Exchange contract.
type ExchangeSetOwner struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterSetOwner is a free log retrieval operation binding the contract event 0xcbf985117192c8f614a58aaf97226bb80a754772f5f6edf06f87c675f2e6c663.
//
// Solidity: e SetOwner(previousOwner indexed address, newOwner indexed address)
func (_Exchange *ExchangeFilterer) FilterSetOwner(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*ExchangeSetOwnerIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _Exchange.contract.FilterLogs(opts, "SetOwner", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &ExchangeSetOwnerIterator{contract: _Exchange.contract, event: "SetOwner", logs: logs, sub: sub}, nil
}

// WatchSetOwner is a free log subscription operation binding the contract event 0xcbf985117192c8f614a58aaf97226bb80a754772f5f6edf06f87c675f2e6c663.
//
// Solidity: e SetOwner(previousOwner indexed address, newOwner indexed address)
func (_Exchange *ExchangeFilterer) WatchSetOwner(opts *bind.WatchOpts, sink chan<- *ExchangeSetOwner, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _Exchange.contract.WatchLogs(opts, "SetOwner", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ExchangeSetOwner)
				if err := _Exchange.contract.UnpackLog(event, "SetOwner", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// OwnedABI is the input ABI used to generate the binding from.
const OwnedABI = "[{\"constant\":false,\"inputs\":[{\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"setOwner\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"SetOwner\",\"type\":\"event\"}]"

// OwnedBin is the compiled bytecode used for deploying new contracts.
const OwnedBin = `0x608060405234801561001057600080fd5b5060008054600160a060020a031916331790556101ac806100326000396000f30060806040526004361061004b5763ffffffff7c010000000000000000000000000000000000000000000000000000000060003504166313af403581146100505780638da5cb5b14610080575b600080fd5b34801561005c57600080fd5b5061007e73ffffffffffffffffffffffffffffffffffffffff600435166100be565b005b34801561008c57600080fd5b50610095610164565b6040805173ffffffffffffffffffffffffffffffffffffffff9092168252519081900360200190f35b60005473ffffffffffffffffffffffffffffffffffffffff1633146100e257600080fd5b6000805460405173ffffffffffffffffffffffffffffffffffffffff808516939216917fcbf985117192c8f614a58aaf97226bb80a754772f5f6edf06f87c675f2e6c66391a36000805473ffffffffffffffffffffffffffffffffffffffff191673ffffffffffffffffffffffffffffffffffffffff92909216919091179055565b60005473ffffffffffffffffffffffffffffffffffffffff16815600a165627a7a7230582076e36060821a6167abb283e930822baede8db6b867d5ad2ae62a2e530ed0aad10029`

// DeployOwned deploys a new Ethereum contract, binding an instance of Owned to it.
func DeployOwned(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *Owned, error) {
	parsed, err := abi.JSON(strings.NewReader(OwnedABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(OwnedBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Owned{OwnedCaller: OwnedCaller{contract: contract}, OwnedTransactor: OwnedTransactor{contract: contract}, OwnedFilterer: OwnedFilterer{contract: contract}}, nil
}

// Owned is an auto generated Go binding around an Ethereum contract.
type Owned struct {
	OwnedCaller     // Read-only binding to the contract
	OwnedTransactor // Write-only binding to the contract
	OwnedFilterer   // Log filterer for contract events
}

// OwnedCaller is an auto generated read-only Go binding around an Ethereum contract.
type OwnedCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// OwnedTransactor is an auto generated write-only Go binding around an Ethereum contract.
type OwnedTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// OwnedFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type OwnedFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// OwnedSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type OwnedSession struct {
	Contract     *Owned            // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// OwnedCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type OwnedCallerSession struct {
	Contract *OwnedCaller  // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts // Call options to use throughout this session
}

// OwnedTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type OwnedTransactorSession struct {
	Contract     *OwnedTransactor  // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// OwnedRaw is an auto generated low-level Go binding around an Ethereum contract.
type OwnedRaw struct {
	Contract *Owned // Generic contract binding to access the raw methods on
}

// OwnedCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type OwnedCallerRaw struct {
	Contract *OwnedCaller // Generic read-only contract binding to access the raw methods on
}

// OwnedTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type OwnedTransactorRaw struct {
	Contract *OwnedTransactor // Generic write-only contract binding to access the raw methods on
}

// NewOwned creates a new instance of Owned, bound to a specific deployed contract.
func NewOwned(address common.Address, backend bind.ContractBackend) (*Owned, error) {
	contract, err := bindOwned(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Owned{OwnedCaller: OwnedCaller{contract: contract}, OwnedTransactor: OwnedTransactor{contract: contract}, OwnedFilterer: OwnedFilterer{contract: contract}}, nil
}

// NewOwnedCaller creates a new read-only instance of Owned, bound to a specific deployed contract.
func NewOwnedCaller(address common.Address, caller bind.ContractCaller) (*OwnedCaller, error) {
	contract, err := bindOwned(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &OwnedCaller{contract: contract}, nil
}

// NewOwnedTransactor creates a new write-only instance of Owned, bound to a specific deployed contract.
func NewOwnedTransactor(address common.Address, transactor bind.ContractTransactor) (*OwnedTransactor, error) {
	contract, err := bindOwned(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &OwnedTransactor{contract: contract}, nil
}

// NewOwnedFilterer creates a new log filterer instance of Owned, bound to a specific deployed contract.
func NewOwnedFilterer(address common.Address, filterer bind.ContractFilterer) (*OwnedFilterer, error) {
	contract, err := bindOwned(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &OwnedFilterer{contract: contract}, nil
}

// bindOwned binds a generic wrapper to an already deployed contract.
func bindOwned(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(OwnedABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Owned *OwnedRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _Owned.Contract.OwnedCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Owned *OwnedRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Owned.Contract.OwnedTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Owned *OwnedRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Owned.Contract.OwnedTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Owned *OwnedCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _Owned.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Owned *OwnedTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Owned.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Owned *OwnedTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Owned.Contract.contract.Transact(opts, method, params...)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() constant returns(address)
func (_Owned *OwnedCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _Owned.contract.Call(opts, out, "owner")
	return *ret0, err
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() constant returns(address)
func (_Owned *OwnedSession) Owner() (common.Address, error) {
	return _Owned.Contract.Owner(&_Owned.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() constant returns(address)
func (_Owned *OwnedCallerSession) Owner() (common.Address, error) {
	return _Owned.Contract.Owner(&_Owned.CallOpts)
}

// SetOwner is a paid mutator transaction binding the contract method 0x13af4035.
//
// Solidity: function setOwner(newOwner address) returns()
func (_Owned *OwnedTransactor) SetOwner(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _Owned.contract.Transact(opts, "setOwner", newOwner)
}

// SetOwner is a paid mutator transaction binding the contract method 0x13af4035.
//
// Solidity: function setOwner(newOwner address) returns()
func (_Owned *OwnedSession) SetOwner(newOwner common.Address) (*types.Transaction, error) {
	return _Owned.Contract.SetOwner(&_Owned.TransactOpts, newOwner)
}

// SetOwner is a paid mutator transaction binding the contract method 0x13af4035.
//
// Solidity: function setOwner(newOwner address) returns()
func (_Owned *OwnedTransactorSession) SetOwner(newOwner common.Address) (*types.Transaction, error) {
	return _Owned.Contract.SetOwner(&_Owned.TransactOpts, newOwner)
}

// OwnedSetOwnerIterator is returned from FilterSetOwner and is used to iterate over the raw logs and unpacked data for SetOwner events raised by the Owned contract.
type OwnedSetOwnerIterator struct {
	Event *OwnedSetOwner // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *OwnedSetOwnerIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OwnedSetOwner)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(OwnedSetOwner)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *OwnedSetOwnerIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *OwnedSetOwnerIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// OwnedSetOwner represents a SetOwner event raised by the Owned contract.
type OwnedSetOwner struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterSetOwner is a free log retrieval operation binding the contract event 0xcbf985117192c8f614a58aaf97226bb80a754772f5f6edf06f87c675f2e6c663.
//
// Solidity: e SetOwner(previousOwner indexed address, newOwner indexed address)
func (_Owned *OwnedFilterer) FilterSetOwner(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*OwnedSetOwnerIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _Owned.contract.FilterLogs(opts, "SetOwner", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &OwnedSetOwnerIterator{contract: _Owned.contract, event: "SetOwner", logs: logs, sub: sub}, nil
}

// WatchSetOwner is a free log subscription operation binding the contract event 0xcbf985117192c8f614a58aaf97226bb80a754772f5f6edf06f87c675f2e6c663.
//
// Solidity: e SetOwner(previousOwner indexed address, newOwner indexed address)
func (_Owned *OwnedFilterer) WatchSetOwner(opts *bind.WatchOpts, sink chan<- *OwnedSetOwner, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _Owned.contract.WatchLogs(opts, "SetOwner", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(OwnedSetOwner)
				if err := _Owned.contract.UnpackLog(event, "SetOwner", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// SafeMathABI is the input ABI used to generate the binding from.
const SafeMathABI = "[]"

// SafeMathBin is the compiled bytecode used for deploying new contracts.
const SafeMathBin = `0x604c602c600b82828239805160001a60731460008114601c57601e565bfe5b5030600052607381538281f30073000000000000000000000000000000000000000030146080604052600080fd00a165627a7a723058201d05756c063faa48d970328dfbd42d1d8e0bab5ab14a3e37f25c5fedefc384b90029`

// DeploySafeMath deploys a new Ethereum contract, binding an instance of SafeMath to it.
func DeploySafeMath(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *SafeMath, error) {
	parsed, err := abi.JSON(strings.NewReader(SafeMathABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(SafeMathBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &SafeMath{SafeMathCaller: SafeMathCaller{contract: contract}, SafeMathTransactor: SafeMathTransactor{contract: contract}, SafeMathFilterer: SafeMathFilterer{contract: contract}}, nil
}

// SafeMath is an auto generated Go binding around an Ethereum contract.
type SafeMath struct {
	SafeMathCaller     // Read-only binding to the contract
	SafeMathTransactor // Write-only binding to the contract
	SafeMathFilterer   // Log filterer for contract events
}

// SafeMathCaller is an auto generated read-only Go binding around an Ethereum contract.
type SafeMathCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SafeMathTransactor is an auto generated write-only Go binding around an Ethereum contract.
type SafeMathTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SafeMathFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type SafeMathFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SafeMathSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type SafeMathSession struct {
	Contract     *SafeMath         // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// SafeMathCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type SafeMathCallerSession struct {
	Contract *SafeMathCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts   // Call options to use throughout this session
}

// SafeMathTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type SafeMathTransactorSession struct {
	Contract     *SafeMathTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts   // Transaction auth options to use throughout this session
}

// SafeMathRaw is an auto generated low-level Go binding around an Ethereum contract.
type SafeMathRaw struct {
	Contract *SafeMath // Generic contract binding to access the raw methods on
}

// SafeMathCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type SafeMathCallerRaw struct {
	Contract *SafeMathCaller // Generic read-only contract binding to access the raw methods on
}

// SafeMathTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type SafeMathTransactorRaw struct {
	Contract *SafeMathTransactor // Generic write-only contract binding to access the raw methods on
}

// NewSafeMath creates a new instance of SafeMath, bound to a specific deployed contract.
func NewSafeMath(address common.Address, backend bind.ContractBackend) (*SafeMath, error) {
	contract, err := bindSafeMath(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &SafeMath{SafeMathCaller: SafeMathCaller{contract: contract}, SafeMathTransactor: SafeMathTransactor{contract: contract}, SafeMathFilterer: SafeMathFilterer{contract: contract}}, nil
}

// NewSafeMathCaller creates a new read-only instance of SafeMath, bound to a specific deployed contract.
func NewSafeMathCaller(address common.Address, caller bind.ContractCaller) (*SafeMathCaller, error) {
	contract, err := bindSafeMath(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &SafeMathCaller{contract: contract}, nil
}

// NewSafeMathTransactor creates a new write-only instance of SafeMath, bound to a specific deployed contract.
func NewSafeMathTransactor(address common.Address, transactor bind.ContractTransactor) (*SafeMathTransactor, error) {
	contract, err := bindSafeMath(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &SafeMathTransactor{contract: contract}, nil
}

// NewSafeMathFilterer creates a new log filterer instance of SafeMath, bound to a specific deployed contract.
func NewSafeMathFilterer(address common.Address, filterer bind.ContractFilterer) (*SafeMathFilterer, error) {
	contract, err := bindSafeMath(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &SafeMathFilterer{contract: contract}, nil
}

// bindSafeMath binds a generic wrapper to an already deployed contract.
func bindSafeMath(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(SafeMathABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_SafeMath *SafeMathRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _SafeMath.Contract.SafeMathCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_SafeMath *SafeMathRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SafeMath.Contract.SafeMathTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_SafeMath *SafeMathRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _SafeMath.Contract.SafeMathTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_SafeMath *SafeMathCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _SafeMath.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_SafeMath *SafeMathTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SafeMath.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_SafeMath *SafeMathTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _SafeMath.Contract.contract.Transact(opts, method, params...)
}
