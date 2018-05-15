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

// Token1ABI is the input ABI used to generate the binding from.
const Token1ABI = "[{\"constant\":true,\"inputs\":[],\"name\":\"mintingFinished\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"setOwner\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"symbol\",\"outputs\":[{\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"spender\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Approval\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"Mint\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[],\"name\":\"MintFinished\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"SetOwner\",\"type\":\"event\"},{\"constant\":false,\"inputs\":[{\"name\":\"_to\",\"type\":\"address\"},{\"name\":\"_amount\",\"type\":\"uint256\"}],\"name\":\"Token\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"totalSupply\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_to\",\"type\":\"address\"},{\"name\":\"_value\",\"type\":\"uint256\"}],\"name\":\"transfer\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_owner\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_from\",\"type\":\"address\"},{\"name\":\"_to\",\"type\":\"address\"},{\"name\":\"_value\",\"type\":\"uint256\"}],\"name\":\"transferFrom\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_spender\",\"type\":\"address\"},{\"name\":\"_value\",\"type\":\"uint256\"}],\"name\":\"approve\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_owner\",\"type\":\"address\"},{\"name\":\"_spender\",\"type\":\"address\"}],\"name\":\"allowance\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_spender\",\"type\":\"address\"},{\"name\":\"_addedValue\",\"type\":\"uint256\"}],\"name\":\"increaseApproval\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_spender\",\"type\":\"address\"},{\"name\":\"_subtractedValue\",\"type\":\"uint256\"}],\"name\":\"decreaseApproval\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_to\",\"type\":\"address\"},{\"name\":\"_amount\",\"type\":\"uint256\"}],\"name\":\"mint\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"finishMinting\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]"

// Token1 is an auto generated Go binding around an Ethereum contract.
type Token1 struct {
	Token1Caller     // Read-only binding to the contract
	Token1Transactor // Write-only binding to the contract
}

// Token1Caller is an auto generated read-only Go binding around an Ethereum contract.
type Token1Caller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// Token1Transactor is an auto generated write-only Go binding around an Ethereum contract.
type Token1Transactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// Token1Session is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type Token1Session struct {
	Contract     *Token1           // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// Token1CallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type Token1CallerSession struct {
	Contract *Token1Caller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts // Call options to use throughout this session
}

// Token1TransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type Token1TransactorSession struct {
	Contract     *Token1Transactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// Token1Raw is an auto generated low-level Go binding around an Ethereum contract.
type Token1Raw struct {
	Contract *Token1 // Generic contract binding to access the raw methods on
}

// Token1CallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type Token1CallerRaw struct {
	Contract *Token1Caller // Generic read-only contract binding to access the raw methods on
}

// Token1TransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type Token1TransactorRaw struct {
	Contract *Token1Transactor // Generic write-only contract binding to access the raw methods on
}

// NewToken1 creates a new instance of Token1, bound to a specific deployed contract.
func NewToken1(address common.Address, backend bind.ContractBackend) (*Token1, error) {
	contract, err := bindToken1(address, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Token1{Token1Caller: Token1Caller{contract: contract}, Token1Transactor: Token1Transactor{contract: contract}}, nil
}

// NewToken1Caller creates a new read-only instance of Token1, bound to a specific deployed contract.
func NewToken1Caller(address common.Address, caller bind.ContractCaller) (*Token1Caller, error) {
	contract, err := bindToken1(address, caller, nil)
	if err != nil {
		return nil, err
	}
	return &Token1Caller{contract: contract}, nil
}

// NewToken1Transactor creates a new write-only instance of Token1, bound to a specific deployed contract.
func NewToken1Transactor(address common.Address, transactor bind.ContractTransactor) (*Token1Transactor, error) {
	contract, err := bindToken1(address, nil, transactor)
	if err != nil {
		return nil, err
	}
	return &Token1Transactor{contract: contract}, nil
}

// bindToken1 binds a generic wrapper to an already deployed contract.
func bindToken1(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(Token1ABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Token1 *Token1Raw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _Token1.Contract.Token1Caller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Token1 *Token1Raw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Token1.Contract.Token1Transactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Token1 *Token1Raw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Token1.Contract.Token1Transactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Token1 *Token1CallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _Token1.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Token1 *Token1TransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Token1.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Token1 *Token1TransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Token1.Contract.contract.Transact(opts, method, params...)
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(_owner address, _spender address) constant returns(uint256)
func (_Token1 *Token1Caller) Allowance(opts *bind.CallOpts, _owner common.Address, _spender common.Address) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Token1.contract.Call(opts, out, "allowance", _owner, _spender)
	return *ret0, err
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(_owner address, _spender address) constant returns(uint256)
func (_Token1 *Token1Session) Allowance(_owner common.Address, _spender common.Address) (*big.Int, error) {
	return _Token1.Contract.Allowance(&_Token1.CallOpts, _owner, _spender)
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(_owner address, _spender address) constant returns(uint256)
func (_Token1 *Token1CallerSession) Allowance(_owner common.Address, _spender common.Address) (*big.Int, error) {
	return _Token1.Contract.Allowance(&_Token1.CallOpts, _owner, _spender)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(_owner address) constant returns(uint256)
func (_Token1 *Token1Caller) BalanceOf(opts *bind.CallOpts, _owner common.Address) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Token1.contract.Call(opts, out, "balanceOf", _owner)
	return *ret0, err
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(_owner address) constant returns(uint256)
func (_Token1 *Token1Session) BalanceOf(_owner common.Address) (*big.Int, error) {
	return _Token1.Contract.BalanceOf(&_Token1.CallOpts, _owner)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(_owner address) constant returns(uint256)
func (_Token1 *Token1CallerSession) BalanceOf(_owner common.Address) (*big.Int, error) {
	return _Token1.Contract.BalanceOf(&_Token1.CallOpts, _owner)
}

// MintingFinished is a free data retrieval call binding the contract method 0x05d2035b.
//
// Solidity: function mintingFinished() constant returns(bool)
func (_Token1 *Token1Caller) MintingFinished(opts *bind.CallOpts) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _Token1.contract.Call(opts, out, "mintingFinished")
	return *ret0, err
}

// MintingFinished is a free data retrieval call binding the contract method 0x05d2035b.
//
// Solidity: function mintingFinished() constant returns(bool)
func (_Token1 *Token1Session) MintingFinished() (bool, error) {
	return _Token1.Contract.MintingFinished(&_Token1.CallOpts)
}

// MintingFinished is a free data retrieval call binding the contract method 0x05d2035b.
//
// Solidity: function mintingFinished() constant returns(bool)
func (_Token1 *Token1CallerSession) MintingFinished() (bool, error) {
	return _Token1.Contract.MintingFinished(&_Token1.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() constant returns(address)
func (_Token1 *Token1Caller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _Token1.contract.Call(opts, out, "owner")
	return *ret0, err
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() constant returns(address)
func (_Token1 *Token1Session) Owner() (common.Address, error) {
	return _Token1.Contract.Owner(&_Token1.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() constant returns(address)
func (_Token1 *Token1CallerSession) Owner() (common.Address, error) {
	return _Token1.Contract.Owner(&_Token1.CallOpts)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() constant returns(string)
func (_Token1 *Token1Caller) Symbol(opts *bind.CallOpts) (string, error) {
	var (
		ret0 = new(string)
	)
	out := ret0
	err := _Token1.contract.Call(opts, out, "symbol")
	return *ret0, err
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() constant returns(string)
func (_Token1 *Token1Session) Symbol() (string, error) {
	return _Token1.Contract.Symbol(&_Token1.CallOpts)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() constant returns(string)
func (_Token1 *Token1CallerSession) Symbol() (string, error) {
	return _Token1.Contract.Symbol(&_Token1.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() constant returns(uint256)
func (_Token1 *Token1Caller) TotalSupply(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Token1.contract.Call(opts, out, "totalSupply")
	return *ret0, err
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() constant returns(uint256)
func (_Token1 *Token1Session) TotalSupply() (*big.Int, error) {
	return _Token1.Contract.TotalSupply(&_Token1.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() constant returns(uint256)
func (_Token1 *Token1CallerSession) TotalSupply() (*big.Int, error) {
	return _Token1.Contract.TotalSupply(&_Token1.CallOpts)
}

// Token is a paid mutator transaction binding the contract method 0xc5944f30.
//
// Solidity: function Token(_to address, _amount uint256) returns()
func (_Token1 *Token1Transactor) Token(opts *bind.TransactOpts, _to common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _Token1.contract.Transact(opts, "Token", _to, _amount)
}

// Token is a paid mutator transaction binding the contract method 0xc5944f30.
//
// Solidity: function Token(_to address, _amount uint256) returns()
func (_Token1 *Token1Session) Token(_to common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _Token1.Contract.Token(&_Token1.TransactOpts, _to, _amount)
}

// Token is a paid mutator transaction binding the contract method 0xc5944f30.
//
// Solidity: function Token(_to address, _amount uint256) returns()
func (_Token1 *Token1TransactorSession) Token(_to common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _Token1.Contract.Token(&_Token1.TransactOpts, _to, _amount)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(_spender address, _value uint256) returns(bool)
func (_Token1 *Token1Transactor) Approve(opts *bind.TransactOpts, _spender common.Address, _value *big.Int) (*types.Transaction, error) {
	return _Token1.contract.Transact(opts, "approve", _spender, _value)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(_spender address, _value uint256) returns(bool)
func (_Token1 *Token1Session) Approve(_spender common.Address, _value *big.Int) (*types.Transaction, error) {
	return _Token1.Contract.Approve(&_Token1.TransactOpts, _spender, _value)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(_spender address, _value uint256) returns(bool)
func (_Token1 *Token1TransactorSession) Approve(_spender common.Address, _value *big.Int) (*types.Transaction, error) {
	return _Token1.Contract.Approve(&_Token1.TransactOpts, _spender, _value)
}

// DecreaseApproval is a paid mutator transaction binding the contract method 0x66188463.
//
// Solidity: function decreaseApproval(_spender address, _subtractedValue uint256) returns(bool)
func (_Token1 *Token1Transactor) DecreaseApproval(opts *bind.TransactOpts, _spender common.Address, _subtractedValue *big.Int) (*types.Transaction, error) {
	return _Token1.contract.Transact(opts, "decreaseApproval", _spender, _subtractedValue)
}

// DecreaseApproval is a paid mutator transaction binding the contract method 0x66188463.
//
// Solidity: function decreaseApproval(_spender address, _subtractedValue uint256) returns(bool)
func (_Token1 *Token1Session) DecreaseApproval(_spender common.Address, _subtractedValue *big.Int) (*types.Transaction, error) {
	return _Token1.Contract.DecreaseApproval(&_Token1.TransactOpts, _spender, _subtractedValue)
}

// DecreaseApproval is a paid mutator transaction binding the contract method 0x66188463.
//
// Solidity: function decreaseApproval(_spender address, _subtractedValue uint256) returns(bool)
func (_Token1 *Token1TransactorSession) DecreaseApproval(_spender common.Address, _subtractedValue *big.Int) (*types.Transaction, error) {
	return _Token1.Contract.DecreaseApproval(&_Token1.TransactOpts, _spender, _subtractedValue)
}

// FinishMinting is a paid mutator transaction binding the contract method 0x7d64bcb4.
//
// Solidity: function finishMinting() returns(bool)
func (_Token1 *Token1Transactor) FinishMinting(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Token1.contract.Transact(opts, "finishMinting")
}

// FinishMinting is a paid mutator transaction binding the contract method 0x7d64bcb4.
//
// Solidity: function finishMinting() returns(bool)
func (_Token1 *Token1Session) FinishMinting() (*types.Transaction, error) {
	return _Token1.Contract.FinishMinting(&_Token1.TransactOpts)
}

// FinishMinting is a paid mutator transaction binding the contract method 0x7d64bcb4.
//
// Solidity: function finishMinting() returns(bool)
func (_Token1 *Token1TransactorSession) FinishMinting() (*types.Transaction, error) {
	return _Token1.Contract.FinishMinting(&_Token1.TransactOpts)
}

// IncreaseApproval is a paid mutator transaction binding the contract method 0xd73dd623.
//
// Solidity: function increaseApproval(_spender address, _addedValue uint256) returns(bool)
func (_Token1 *Token1Transactor) IncreaseApproval(opts *bind.TransactOpts, _spender common.Address, _addedValue *big.Int) (*types.Transaction, error) {
	return _Token1.contract.Transact(opts, "increaseApproval", _spender, _addedValue)
}

// IncreaseApproval is a paid mutator transaction binding the contract method 0xd73dd623.
//
// Solidity: function increaseApproval(_spender address, _addedValue uint256) returns(bool)
func (_Token1 *Token1Session) IncreaseApproval(_spender common.Address, _addedValue *big.Int) (*types.Transaction, error) {
	return _Token1.Contract.IncreaseApproval(&_Token1.TransactOpts, _spender, _addedValue)
}

// IncreaseApproval is a paid mutator transaction binding the contract method 0xd73dd623.
//
// Solidity: function increaseApproval(_spender address, _addedValue uint256) returns(bool)
func (_Token1 *Token1TransactorSession) IncreaseApproval(_spender common.Address, _addedValue *big.Int) (*types.Transaction, error) {
	return _Token1.Contract.IncreaseApproval(&_Token1.TransactOpts, _spender, _addedValue)
}

// Mint is a paid mutator transaction binding the contract method 0x40c10f19.
//
// Solidity: function mint(_to address, _amount uint256) returns(bool)
func (_Token1 *Token1Transactor) Mint(opts *bind.TransactOpts, _to common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _Token1.contract.Transact(opts, "mint", _to, _amount)
}

// Mint is a paid mutator transaction binding the contract method 0x40c10f19.
//
// Solidity: function mint(_to address, _amount uint256) returns(bool)
func (_Token1 *Token1Session) Mint(_to common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _Token1.Contract.Mint(&_Token1.TransactOpts, _to, _amount)
}

// Mint is a paid mutator transaction binding the contract method 0x40c10f19.
//
// Solidity: function mint(_to address, _amount uint256) returns(bool)
func (_Token1 *Token1TransactorSession) Mint(_to common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _Token1.Contract.Mint(&_Token1.TransactOpts, _to, _amount)
}

// SetOwner is a paid mutator transaction binding the contract method 0x13af4035.
//
// Solidity: function setOwner(newOwner address) returns()
func (_Token1 *Token1Transactor) SetOwner(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _Token1.contract.Transact(opts, "setOwner", newOwner)
}

// SetOwner is a paid mutator transaction binding the contract method 0x13af4035.
//
// Solidity: function setOwner(newOwner address) returns()
func (_Token1 *Token1Session) SetOwner(newOwner common.Address) (*types.Transaction, error) {
	return _Token1.Contract.SetOwner(&_Token1.TransactOpts, newOwner)
}

// SetOwner is a paid mutator transaction binding the contract method 0x13af4035.
//
// Solidity: function setOwner(newOwner address) returns()
func (_Token1 *Token1TransactorSession) SetOwner(newOwner common.Address) (*types.Transaction, error) {
	return _Token1.Contract.SetOwner(&_Token1.TransactOpts, newOwner)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(_to address, _value uint256) returns(bool)
func (_Token1 *Token1Transactor) Transfer(opts *bind.TransactOpts, _to common.Address, _value *big.Int) (*types.Transaction, error) {
	return _Token1.contract.Transact(opts, "transfer", _to, _value)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(_to address, _value uint256) returns(bool)
func (_Token1 *Token1Session) Transfer(_to common.Address, _value *big.Int) (*types.Transaction, error) {
	return _Token1.Contract.Transfer(&_Token1.TransactOpts, _to, _value)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(_to address, _value uint256) returns(bool)
func (_Token1 *Token1TransactorSession) Transfer(_to common.Address, _value *big.Int) (*types.Transaction, error) {
	return _Token1.Contract.Transfer(&_Token1.TransactOpts, _to, _value)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(_from address, _to address, _value uint256) returns(bool)
func (_Token1 *Token1Transactor) TransferFrom(opts *bind.TransactOpts, _from common.Address, _to common.Address, _value *big.Int) (*types.Transaction, error) {
	return _Token1.contract.Transact(opts, "transferFrom", _from, _to, _value)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(_from address, _to address, _value uint256) returns(bool)
func (_Token1 *Token1Session) TransferFrom(_from common.Address, _to common.Address, _value *big.Int) (*types.Transaction, error) {
	return _Token1.Contract.TransferFrom(&_Token1.TransactOpts, _from, _to, _value)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(_from address, _to address, _value uint256) returns(bool)
func (_Token1 *Token1TransactorSession) TransferFrom(_from common.Address, _to common.Address, _value *big.Int) (*types.Transaction, error) {
	return _Token1.Contract.TransferFrom(&_Token1.TransactOpts, _from, _to, _value)
}
