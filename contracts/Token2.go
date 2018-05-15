// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package main

import (
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// Token2ABI is the input ABI used to generate the binding from.
const Token2ABI = "[{\"constant\":true,\"inputs\":[],\"name\":\"mintingFinished\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"setOwner\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"symbol\",\"outputs\":[{\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"spender\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Approval\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"Mint\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[],\"name\":\"MintFinished\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"SetOwner\",\"type\":\"event\"},{\"constant\":false,\"inputs\":[{\"name\":\"_to\",\"type\":\"address\"},{\"name\":\"_amount\",\"type\":\"uint256\"}],\"name\":\"Token\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"totalSupply\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_to\",\"type\":\"address\"},{\"name\":\"_value\",\"type\":\"uint256\"}],\"name\":\"transfer\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_owner\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_from\",\"type\":\"address\"},{\"name\":\"_to\",\"type\":\"address\"},{\"name\":\"_value\",\"type\":\"uint256\"}],\"name\":\"transferFrom\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_spender\",\"type\":\"address\"},{\"name\":\"_value\",\"type\":\"uint256\"}],\"name\":\"approve\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_owner\",\"type\":\"address\"},{\"name\":\"_spender\",\"type\":\"address\"}],\"name\":\"allowance\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_spender\",\"type\":\"address\"},{\"name\":\"_addedValue\",\"type\":\"uint256\"}],\"name\":\"increaseApproval\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_spender\",\"type\":\"address\"},{\"name\":\"_subtractedValue\",\"type\":\"uint256\"}],\"name\":\"decreaseApproval\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_to\",\"type\":\"address\"},{\"name\":\"_amount\",\"type\":\"uint256\"}],\"name\":\"mint\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"finishMinting\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]"

// Token2 is an auto generated Go binding around an Ethereum contract.
type Token2 struct {
	Token2Caller     // Read-only binding to the contract
	Token2Transactor // Write-only binding to the contract
}

// Token2Caller is an auto generated read-only Go binding around an Ethereum contract.
type Token2Caller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// Token2Transactor is an auto generated write-only Go binding around an Ethereum contract.
type Token2Transactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// Token2Session is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type Token2Session struct {
	Contract     *Token2           // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// Token2CallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type Token2CallerSession struct {
	Contract *Token2Caller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts // Call options to use throughout this session
}

// Token2TransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type Token2TransactorSession struct {
	Contract     *Token2Transactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// Token2Raw is an auto generated low-level Go binding around an Ethereum contract.
type Token2Raw struct {
	Contract *Token2 // Generic contract binding to access the raw methods on
}

// Token2CallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type Token2CallerRaw struct {
	Contract *Token2Caller // Generic read-only contract binding to access the raw methods on
}

// Token2TransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type Token2TransactorRaw struct {
	Contract *Token2Transactor // Generic write-only contract binding to access the raw methods on
}

// NewToken2 creates a new instance of Token2, bound to a specific deployed contract.
func NewToken2(address common.Address, backend bind.ContractBackend) (*Token2, error) {
	contract, err := bindToken2(address, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Token2{Token2Caller: Token2Caller{contract: contract}, Token2Transactor: Token2Transactor{contract: contract}}, nil
}

// NewToken2Caller creates a new read-only instance of Token2, bound to a specific deployed contract.
func NewToken2Caller(address common.Address, caller bind.ContractCaller) (*Token2Caller, error) {
	contract, err := bindToken2(address, caller, nil)
	if err != nil {
		return nil, err
	}
	return &Token2Caller{contract: contract}, nil
}

// NewToken2Transactor creates a new write-only instance of Token2, bound to a specific deployed contract.
func NewToken2Transactor(address common.Address, transactor bind.ContractTransactor) (*Token2Transactor, error) {
	contract, err := bindToken2(address, nil, transactor)
	if err != nil {
		return nil, err
	}
	return &Token2Transactor{contract: contract}, nil
}

// bindToken2 binds a generic wrapper to an already deployed contract.
func bindToken2(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(Token2ABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Token2 *Token2Raw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _Token2.Contract.Token2Caller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Token2 *Token2Raw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Token2.Contract.Token2Transactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Token2 *Token2Raw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Token2.Contract.Token2Transactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Token2 *Token2CallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _Token2.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Token2 *Token2TransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Token2.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Token2 *Token2TransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Token2.Contract.contract.Transact(opts, method, params...)
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(_owner address, _spender address) constant returns(uint256)
func (_Token2 *Token2Caller) Allowance(opts *bind.CallOpts, _owner common.Address, _spender common.Address) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Token2.contract.Call(opts, out, "allowance", _owner, _spender)
	return *ret0, err
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(_owner address, _spender address) constant returns(uint256)
func (_Token2 *Token2Session) Allowance(_owner common.Address, _spender common.Address) (*big.Int, error) {
	return _Token2.Contract.Allowance(&_Token2.CallOpts, _owner, _spender)
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(_owner address, _spender address) constant returns(uint256)
func (_Token2 *Token2CallerSession) Allowance(_owner common.Address, _spender common.Address) (*big.Int, error) {
	return _Token2.Contract.Allowance(&_Token2.CallOpts, _owner, _spender)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(_owner address) constant returns(uint256)
func (_Token2 *Token2Caller) BalanceOf(opts *bind.CallOpts, _owner common.Address) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Token2.contract.Call(opts, out, "balanceOf", _owner)
	return *ret0, err
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(_owner address) constant returns(uint256)
func (_Token2 *Token2Session) BalanceOf(_owner common.Address) (*big.Int, error) {
	return _Token2.Contract.BalanceOf(&_Token2.CallOpts, _owner)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(_owner address) constant returns(uint256)
func (_Token2 *Token2CallerSession) BalanceOf(_owner common.Address) (*big.Int, error) {
	return _Token2.Contract.BalanceOf(&_Token2.CallOpts, _owner)
}

// MintingFinished is a free data retrieval call binding the contract method 0x05d2035b.
//
// Solidity: function mintingFinished() constant returns(bool)
func (_Token2 *Token2Caller) MintingFinished(opts *bind.CallOpts) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _Token2.contract.Call(opts, out, "mintingFinished")
	return *ret0, err
}

// MintingFinished is a free data retrieval call binding the contract method 0x05d2035b.
//
// Solidity: function mintingFinished() constant returns(bool)
func (_Token2 *Token2Session) MintingFinished() (bool, error) {
	return _Token2.Contract.MintingFinished(&_Token2.CallOpts)
}

// MintingFinished is a free data retrieval call binding the contract method 0x05d2035b.
//
// Solidity: function mintingFinished() constant returns(bool)
func (_Token2 *Token2CallerSession) MintingFinished() (bool, error) {
	return _Token2.Contract.MintingFinished(&_Token2.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() constant returns(address)
func (_Token2 *Token2Caller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _Token2.contract.Call(opts, out, "owner")
	return *ret0, err
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() constant returns(address)
func (_Token2 *Token2Session) Owner() (common.Address, error) {
	return _Token2.Contract.Owner(&_Token2.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() constant returns(address)
func (_Token2 *Token2CallerSession) Owner() (common.Address, error) {
	return _Token2.Contract.Owner(&_Token2.CallOpts)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() constant returns(string)
func (_Token2 *Token2Caller) Symbol(opts *bind.CallOpts) (string, error) {
	var (
		ret0 = new(string)
	)
	out := ret0
	err := _Token2.contract.Call(opts, out, "symbol")
	return *ret0, err
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() constant returns(string)
func (_Token2 *Token2Session) Symbol() (string, error) {
	return _Token2.Contract.Symbol(&_Token2.CallOpts)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() constant returns(string)
func (_Token2 *Token2CallerSession) Symbol() (string, error) {
	return _Token2.Contract.Symbol(&_Token2.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() constant returns(uint256)
func (_Token2 *Token2Caller) TotalSupply(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Token2.contract.Call(opts, out, "totalSupply")
	return *ret0, err
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() constant returns(uint256)
func (_Token2 *Token2Session) TotalSupply() (*big.Int, error) {
	return _Token2.Contract.TotalSupply(&_Token2.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() constant returns(uint256)
func (_Token2 *Token2CallerSession) TotalSupply() (*big.Int, error) {
	return _Token2.Contract.TotalSupply(&_Token2.CallOpts)
}

// Token is a paid mutator transaction binding the contract method 0xc5944f30.
//
// Solidity: function Token(_to address, _amount uint256) returns()
func (_Token2 *Token2Transactor) Token(opts *bind.TransactOpts, _to common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _Token2.contract.Transact(opts, "Token", _to, _amount)
}

// Token is a paid mutator transaction binding the contract method 0xc5944f30.
//
// Solidity: function Token(_to address, _amount uint256) returns()
func (_Token2 *Token2Session) Token(_to common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _Token2.Contract.Token(&_Token2.TransactOpts, _to, _amount)
}

// Token is a paid mutator transaction binding the contract method 0xc5944f30.
//
// Solidity: function Token(_to address, _amount uint256) returns()
func (_Token2 *Token2TransactorSession) Token(_to common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _Token2.Contract.Token(&_Token2.TransactOpts, _to, _amount)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(_spender address, _value uint256) returns(bool)
func (_Token2 *Token2Transactor) Approve(opts *bind.TransactOpts, _spender common.Address, _value *big.Int) (*types.Transaction, error) {
	return _Token2.contract.Transact(opts, "approve", _spender, _value)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(_spender address, _value uint256) returns(bool)
func (_Token2 *Token2Session) Approve(_spender common.Address, _value *big.Int) (*types.Transaction, error) {
	return _Token2.Contract.Approve(&_Token2.TransactOpts, _spender, _value)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(_spender address, _value uint256) returns(bool)
func (_Token2 *Token2TransactorSession) Approve(_spender common.Address, _value *big.Int) (*types.Transaction, error) {
	return _Token2.Contract.Approve(&_Token2.TransactOpts, _spender, _value)
}

// DecreaseApproval is a paid mutator transaction binding the contract method 0x66188463.
//
// Solidity: function decreaseApproval(_spender address, _subtractedValue uint256) returns(bool)
func (_Token2 *Token2Transactor) DecreaseApproval(opts *bind.TransactOpts, _spender common.Address, _subtractedValue *big.Int) (*types.Transaction, error) {
	return _Token2.contract.Transact(opts, "decreaseApproval", _spender, _subtractedValue)
}

// DecreaseApproval is a paid mutator transaction binding the contract method 0x66188463.
//
// Solidity: function decreaseApproval(_spender address, _subtractedValue uint256) returns(bool)
func (_Token2 *Token2Session) DecreaseApproval(_spender common.Address, _subtractedValue *big.Int) (*types.Transaction, error) {
	return _Token2.Contract.DecreaseApproval(&_Token2.TransactOpts, _spender, _subtractedValue)
}

// DecreaseApproval is a paid mutator transaction binding the contract method 0x66188463.
//
// Solidity: function decreaseApproval(_spender address, _subtractedValue uint256) returns(bool)
func (_Token2 *Token2TransactorSession) DecreaseApproval(_spender common.Address, _subtractedValue *big.Int) (*types.Transaction, error) {
	return _Token2.Contract.DecreaseApproval(&_Token2.TransactOpts, _spender, _subtractedValue)
}

// FinishMinting is a paid mutator transaction binding the contract method 0x7d64bcb4.
//
// Solidity: function finishMinting() returns(bool)
func (_Token2 *Token2Transactor) FinishMinting(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Token2.contract.Transact(opts, "finishMinting")
}

// FinishMinting is a paid mutator transaction binding the contract method 0x7d64bcb4.
//
// Solidity: function finishMinting() returns(bool)
func (_Token2 *Token2Session) FinishMinting() (*types.Transaction, error) {
	return _Token2.Contract.FinishMinting(&_Token2.TransactOpts)
}

// FinishMinting is a paid mutator transaction binding the contract method 0x7d64bcb4.
//
// Solidity: function finishMinting() returns(bool)
func (_Token2 *Token2TransactorSession) FinishMinting() (*types.Transaction, error) {
	return _Token2.Contract.FinishMinting(&_Token2.TransactOpts)
}

// IncreaseApproval is a paid mutator transaction binding the contract method 0xd73dd623.
//
// Solidity: function increaseApproval(_spender address, _addedValue uint256) returns(bool)
func (_Token2 *Token2Transactor) IncreaseApproval(opts *bind.TransactOpts, _spender common.Address, _addedValue *big.Int) (*types.Transaction, error) {
	return _Token2.contract.Transact(opts, "increaseApproval", _spender, _addedValue)
}

// IncreaseApproval is a paid mutator transaction binding the contract method 0xd73dd623.
//
// Solidity: function increaseApproval(_spender address, _addedValue uint256) returns(bool)
func (_Token2 *Token2Session) IncreaseApproval(_spender common.Address, _addedValue *big.Int) (*types.Transaction, error) {
	return _Token2.Contract.IncreaseApproval(&_Token2.TransactOpts, _spender, _addedValue)
}

// IncreaseApproval is a paid mutator transaction binding the contract method 0xd73dd623.
//
// Solidity: function increaseApproval(_spender address, _addedValue uint256) returns(bool)
func (_Token2 *Token2TransactorSession) IncreaseApproval(_spender common.Address, _addedValue *big.Int) (*types.Transaction, error) {
	return _Token2.Contract.IncreaseApproval(&_Token2.TransactOpts, _spender, _addedValue)
}

// Mint is a paid mutator transaction binding the contract method 0x40c10f19.
//
// Solidity: function mint(_to address, _amount uint256) returns(bool)
func (_Token2 *Token2Transactor) Mint(opts *bind.TransactOpts, _to common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _Token2.contract.Transact(opts, "mint", _to, _amount)
}

// Mint is a paid mutator transaction binding the contract method 0x40c10f19.
//
// Solidity: function mint(_to address, _amount uint256) returns(bool)
func (_Token2 *Token2Session) Mint(_to common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _Token2.Contract.Mint(&_Token2.TransactOpts, _to, _amount)
}

// Mint is a paid mutator transaction binding the contract method 0x40c10f19.
//
// Solidity: function mint(_to address, _amount uint256) returns(bool)
func (_Token2 *Token2TransactorSession) Mint(_to common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _Token2.Contract.Mint(&_Token2.TransactOpts, _to, _amount)
}

// SetOwner is a paid mutator transaction binding the contract method 0x13af4035.
//
// Solidity: function setOwner(newOwner address) returns()
func (_Token2 *Token2Transactor) SetOwner(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _Token2.contract.Transact(opts, "setOwner", newOwner)
}

// SetOwner is a paid mutator transaction binding the contract method 0x13af4035.
//
// Solidity: function setOwner(newOwner address) returns()
func (_Token2 *Token2Session) SetOwner(newOwner common.Address) (*types.Transaction, error) {
	return _Token2.Contract.SetOwner(&_Token2.TransactOpts, newOwner)
}

// SetOwner is a paid mutator transaction binding the contract method 0x13af4035.
//
// Solidity: function setOwner(newOwner address) returns()
func (_Token2 *Token2TransactorSession) SetOwner(newOwner common.Address) (*types.Transaction, error) {
	return _Token2.Contract.SetOwner(&_Token2.TransactOpts, newOwner)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(_to address, _value uint256) returns(bool)
func (_Token2 *Token2Transactor) Transfer(opts *bind.TransactOpts, _to common.Address, _value *big.Int) (*types.Transaction, error) {
	return _Token2.contract.Transact(opts, "transfer", _to, _value)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(_to address, _value uint256) returns(bool)
func (_Token2 *Token2Session) Transfer(_to common.Address, _value *big.Int) (*types.Transaction, error) {
	return _Token2.Contract.Transfer(&_Token2.TransactOpts, _to, _value)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(_to address, _value uint256) returns(bool)
func (_Token2 *Token2TransactorSession) Transfer(_to common.Address, _value *big.Int) (*types.Transaction, error) {
	return _Token2.Contract.Transfer(&_Token2.TransactOpts, _to, _value)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(_from address, _to address, _value uint256) returns(bool)
func (_Token2 *Token2Transactor) TransferFrom(opts *bind.TransactOpts, _from common.Address, _to common.Address, _value *big.Int) (*types.Transaction, error) {
	return _Token2.contract.Transact(opts, "transferFrom", _from, _to, _value)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(_from address, _to address, _value uint256) returns(bool)
func (_Token2 *Token2Session) TransferFrom(_from common.Address, _to common.Address, _value *big.Int) (*types.Transaction, error) {
	return _Token2.Contract.TransferFrom(&_Token2.TransactOpts, _from, _to, _value)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(_from address, _to address, _value uint256) returns(bool)
func (_Token2 *Token2TransactorSession) TransferFrom(_from common.Address, _to common.Address, _value *big.Int) (*types.Transaction, error) {
	return _Token2.Contract.TransferFrom(&_Token2.TransactOpts, _from, _to, _value)
}
