// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package interfaces

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
const ERC20ABI = "[{\"constant\":true,\"inputs\":[],\"name\":\"name\",\"outputs\":[{\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_spender\",\"type\":\"address\"},{\"name\":\"_value\",\"type\":\"uint256\"}],\"name\":\"approve\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"totalSupply\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_from\",\"type\":\"address\"},{\"name\":\"_to\",\"type\":\"address\"},{\"name\":\"_value\",\"type\":\"uint256\"}],\"name\":\"transferFrom\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"decimals\",\"outputs\":[{\"name\":\"\",\"type\":\"uint8\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"symbol\",\"outputs\":[{\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"allowTransactions\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_to\",\"type\":\"address\"},{\"name\":\"_value\",\"type\":\"uint256\"}],\"name\":\"transfer\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_spender\",\"type\":\"address\"},{\"name\":\"_value\",\"type\":\"uint256\"},{\"name\":\"_extraData\",\"type\":\"bytes\"}],\"name\":\"approveAndCall\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"address\"},{\"name\":\"\",\"type\":\"address\"}],\"name\":\"allowance\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"}]"

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

// AllowTransactions is a free data retrieval call binding the contract method 0xa5488a37.
//
// Solidity: function allowTransactions() constant returns(bool)
func (_ERC20 *ERC20Caller) AllowTransactions(opts *bind.CallOpts) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _ERC20.contract.Call(opts, out, "allowTransactions")
	return *ret0, err
}

// AllowTransactions is a free data retrieval call binding the contract method 0xa5488a37.
//
// Solidity: function allowTransactions() constant returns(bool)
func (_ERC20 *ERC20Session) AllowTransactions() (bool, error) {
	return _ERC20.Contract.AllowTransactions(&_ERC20.CallOpts)
}

// AllowTransactions is a free data retrieval call binding the contract method 0xa5488a37.
//
// Solidity: function allowTransactions() constant returns(bool)
func (_ERC20 *ERC20CallerSession) AllowTransactions() (bool, error) {
	return _ERC20.Contract.AllowTransactions(&_ERC20.CallOpts)
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance( address,  address) constant returns(uint256)
func (_ERC20 *ERC20Caller) Allowance(opts *bind.CallOpts, arg0 common.Address, arg1 common.Address) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _ERC20.contract.Call(opts, out, "allowance", arg0, arg1)
	return *ret0, err
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance( address,  address) constant returns(uint256)
func (_ERC20 *ERC20Session) Allowance(arg0 common.Address, arg1 common.Address) (*big.Int, error) {
	return _ERC20.Contract.Allowance(&_ERC20.CallOpts, arg0, arg1)
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance( address,  address) constant returns(uint256)
func (_ERC20 *ERC20CallerSession) Allowance(arg0 common.Address, arg1 common.Address) (*big.Int, error) {
	return _ERC20.Contract.Allowance(&_ERC20.CallOpts, arg0, arg1)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf( address) constant returns(uint256)
func (_ERC20 *ERC20Caller) BalanceOf(opts *bind.CallOpts, arg0 common.Address) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _ERC20.contract.Call(opts, out, "balanceOf", arg0)
	return *ret0, err
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf( address) constant returns(uint256)
func (_ERC20 *ERC20Session) BalanceOf(arg0 common.Address) (*big.Int, error) {
	return _ERC20.Contract.BalanceOf(&_ERC20.CallOpts, arg0)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf( address) constant returns(uint256)
func (_ERC20 *ERC20CallerSession) BalanceOf(arg0 common.Address) (*big.Int, error) {
	return _ERC20.Contract.BalanceOf(&_ERC20.CallOpts, arg0)
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

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(_spender address, _value uint256) returns(bool)
func (_ERC20 *ERC20Transactor) Approve(opts *bind.TransactOpts, _spender common.Address, _value *big.Int) (*types.Transaction, error) {
	return _ERC20.contract.Transact(opts, "approve", _spender, _value)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(_spender address, _value uint256) returns(bool)
func (_ERC20 *ERC20Session) Approve(_spender common.Address, _value *big.Int) (*types.Transaction, error) {
	return _ERC20.Contract.Approve(&_ERC20.TransactOpts, _spender, _value)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(_spender address, _value uint256) returns(bool)
func (_ERC20 *ERC20TransactorSession) Approve(_spender common.Address, _value *big.Int) (*types.Transaction, error) {
	return _ERC20.Contract.Approve(&_ERC20.TransactOpts, _spender, _value)
}

// ApproveAndCall is a paid mutator transaction binding the contract method 0xcae9ca51.
//
// Solidity: function approveAndCall(_spender address, _value uint256, _extraData bytes) returns(bool)
func (_ERC20 *ERC20Transactor) ApproveAndCall(opts *bind.TransactOpts, _spender common.Address, _value *big.Int, _extraData []byte) (*types.Transaction, error) {
	return _ERC20.contract.Transact(opts, "approveAndCall", _spender, _value, _extraData)
}

// ApproveAndCall is a paid mutator transaction binding the contract method 0xcae9ca51.
//
// Solidity: function approveAndCall(_spender address, _value uint256, _extraData bytes) returns(bool)
func (_ERC20 *ERC20Session) ApproveAndCall(_spender common.Address, _value *big.Int, _extraData []byte) (*types.Transaction, error) {
	return _ERC20.Contract.ApproveAndCall(&_ERC20.TransactOpts, _spender, _value, _extraData)
}

// ApproveAndCall is a paid mutator transaction binding the contract method 0xcae9ca51.
//
// Solidity: function approveAndCall(_spender address, _value uint256, _extraData bytes) returns(bool)
func (_ERC20 *ERC20TransactorSession) ApproveAndCall(_spender common.Address, _value *big.Int, _extraData []byte) (*types.Transaction, error) {
	return _ERC20.Contract.ApproveAndCall(&_ERC20.TransactOpts, _spender, _value, _extraData)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(_to address, _value uint256) returns(bool)
func (_ERC20 *ERC20Transactor) Transfer(opts *bind.TransactOpts, _to common.Address, _value *big.Int) (*types.Transaction, error) {
	return _ERC20.contract.Transact(opts, "transfer", _to, _value)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(_to address, _value uint256) returns(bool)
func (_ERC20 *ERC20Session) Transfer(_to common.Address, _value *big.Int) (*types.Transaction, error) {
	return _ERC20.Contract.Transfer(&_ERC20.TransactOpts, _to, _value)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(_to address, _value uint256) returns(bool)
func (_ERC20 *ERC20TransactorSession) Transfer(_to common.Address, _value *big.Int) (*types.Transaction, error) {
	return _ERC20.Contract.Transfer(&_ERC20.TransactOpts, _to, _value)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(_from address, _to address, _value uint256) returns(bool)
func (_ERC20 *ERC20Transactor) TransferFrom(opts *bind.TransactOpts, _from common.Address, _to common.Address, _value *big.Int) (*types.Transaction, error) {
	return _ERC20.contract.Transact(opts, "transferFrom", _from, _to, _value)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(_from address, _to address, _value uint256) returns(bool)
func (_ERC20 *ERC20Session) TransferFrom(_from common.Address, _to common.Address, _value *big.Int) (*types.Transaction, error) {
	return _ERC20.Contract.TransferFrom(&_ERC20.TransactOpts, _from, _to, _value)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(_from address, _to address, _value uint256) returns(bool)
func (_ERC20 *ERC20TransactorSession) TransferFrom(_from common.Address, _to common.Address, _value *big.Int) (*types.Transaction, error) {
	return _ERC20.Contract.TransferFrom(&_ERC20.TransactOpts, _from, _to, _value)
}

// ExchangeABI is the input ABI used to generate the binding from.
const ExchangeABI = "[{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"address\"}],\"name\":\"lastTransaction\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"trader\",\"type\":\"address\"},{\"name\":\"token\",\"type\":\"address\"}],\"name\":\"tokenBalance\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"setOwner\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"address\"}],\"name\":\"operators\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"orderValues\",\"type\":\"uint256[8]\"},{\"name\":\"orderAddresses\",\"type\":\"address[4]\"},{\"name\":\"v\",\"type\":\"uint8[2]\"},{\"name\":\"rs\",\"type\":\"bytes32[4]\"}],\"name\":\"executeTrade\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"orderValues\",\"type\":\"uint256[5]\"},{\"name\":\"orderAddresses\",\"type\":\"address[4]\"},{\"name\":\"v\",\"type\":\"uint8\"},{\"name\":\"r\",\"type\":\"bytes32\"},{\"name\":\"s\",\"type\":\"bytes32\"}],\"name\":\"cancelOrder\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"token\",\"type\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"depositToken\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"withdrawn\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"orderHash\",\"type\":\"bytes32\"},{\"name\":\"amount\",\"type\":\"uint256\"},{\"name\":\"tradeNonce\",\"type\":\"uint256\"},{\"name\":\"taker\",\"type\":\"address\"},{\"name\":\"v\",\"type\":\"uint8\"},{\"name\":\"r\",\"type\":\"bytes32\"},{\"name\":\"s\",\"type\":\"bytes32\"}],\"name\":\"cancelTrade\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_feeAccount\",\"type\":\"address\"}],\"name\":\"setFeeAccount\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"address\"},{\"name\":\"\",\"type\":\"address\"}],\"name\":\"tokens\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"operator\",\"type\":\"address\"},{\"name\":\"isOperator\",\"type\":\"bool\"}],\"name\":\"setOperator\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"feeAccount\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"token\",\"type\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\"},{\"name\":\"trader\",\"type\":\"address\"},{\"name\":\"receiver\",\"type\":\"address\"},{\"name\":\"nonce\",\"type\":\"uint256\"},{\"name\":\"v\",\"type\":\"uint8\"},{\"name\":\"rs\",\"type\":\"bytes32[2]\"},{\"name\":\"feeWithdrawal\",\"type\":\"uint256\"}],\"name\":\"withdraw\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"signer\",\"type\":\"address\"},{\"name\":\"hashedData\",\"type\":\"bytes32\"},{\"name\":\"v\",\"type\":\"uint8\"},{\"name\":\"r\",\"type\":\"bytes32\"},{\"name\":\"s\",\"type\":\"bytes32\"}],\"name\":\"isValidSignature\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"depositEther\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"token\",\"type\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"securityWithdraw\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"transferred\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"trader\",\"type\":\"address\"}],\"name\":\"etherBalance\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_withdrawalSecurityPeriod\",\"type\":\"uint256\"}],\"name\":\"setWithdrawalSecurityPeriod\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"traded\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"withdrawalSecurityPeriod\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"orderFills\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"address\"}],\"name\":\"protectedFunds\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"name\":\"_feeAccount\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"fallback\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"tokenBuy\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"tokenSell\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"maker\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"taker\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"hsh\",\"type\":\"bytes32\"}],\"name\":\"LogTrade\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"token\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"user\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"balance\",\"type\":\"uint256\"}],\"name\":\"LogDeposit\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"token\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"user\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"balance\",\"type\":\"uint256\"}],\"name\":\"LogWithdraw\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"token\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"user\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"balance\",\"type\":\"uint256\"}],\"name\":\"LogSecurityWithdraw\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"token\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"recipient\",\"type\":\"address\"}],\"name\":\"LogTransfer\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"tokenBuy\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"amountBuy\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"tokenSell\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"amountSell\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"expires\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"nonce\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"maker\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"v\",\"type\":\"uint8\"},{\"indexed\":false,\"name\":\"r\",\"type\":\"bytes32\"},{\"indexed\":false,\"name\":\"s\",\"type\":\"bytes32\"}],\"name\":\"LogCancelOrder\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"orderHash\",\"type\":\"bytes32\"},{\"indexed\":false,\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"tradeNonce\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"taker\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"v\",\"type\":\"uint8\"},{\"indexed\":false,\"name\":\"r\",\"type\":\"bytes32\"},{\"indexed\":false,\"name\":\"s\",\"type\":\"bytes32\"}],\"name\":\"LogCancelTrade\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"errorId\",\"type\":\"uint8\"},{\"indexed\":false,\"name\":\"orderHash\",\"type\":\"bytes32\"}],\"name\":\"LogError\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"SetOwner\",\"type\":\"event\"}]"

// ExchangeBin is the compiled bytecode used for deploying new contracts.
const ExchangeBin = `0x6060604052341561000f57600080fd5b6040516020806120688339810160405280805160008054600160a060020a03338116600160a060020a0319928316179092556001805492909316911617905550506120098061005f6000396000f3006060604052361561012d5763ffffffff60e060020a6000350416630531be92811461013d5780631049334f1461016e57806313af40351461019357806313e7c9d8146101b45780632207148d146101e75780632e813def1461029f578063338b5dea1461030f5780633823d66c14610331578063468ddf2e146103475780634b023cf81461037b578063508493bc1461039a578063558a7297146103bf57806365e17c9d146103e357806374af3ab3146104125780638163681e146104795780638da5cb5b146104a757806398ea5fca146104ba5780639bc2c131146104c2578063afc441e3146104e4578063cd0c5896146104fa578063cf46d3d714610519578063d58133231461052f578063f3198b1614610545578063f7213db614610558578063fe2e2b941461056e575b341561013857600080fd5b600080fd5b341561014857600080fd5b61015c600160a060020a036004351661058d565b60405190815260200160405180910390f35b341561017957600080fd5b61015c600160a060020a036004358116906024351661059f565b341561019e57600080fd5b6101b2600160a060020a03600435166105cb565b005b34156101bf57600080fd5b6101d3600160a060020a0360043516610651565b604051901515815260200160405180910390f35b34156101f257600080fd5b6101d3600461010481600861010060405190810160405291908282610100808284378201915050505050919080608001906004806020026040519081016040529190828260808082843782019150505050509190806040019060028060200260405190810160405280929190826002602002808284378201915050505050919080608001906004806020026040519081016040529190828260808082843750939550610666945050505050565b34156102aa57600080fd5b6101d3600460a481600560a06040519081016040529190828260a08082843782019150505050509190806080019060048060200260405190810160405291908282608080828437509395505050823560ff169260208101359250604001359050610ffd565b341561031a57600080fd5b6101d3600160a060020a0360043516602435611208565b341561033c57600080fd5b6101d36004356112c6565b341561035257600080fd5b6101d3600435602435604435600160a060020a036064351660ff6084351660a43560c4356112db565b341561038657600080fd5b6101d3600160a060020a036004351661141c565b34156103a557600080fd5b61015c600160a060020a036004358116906024351661148d565b34156103ca57600080fd5b6101d3600160a060020a036004351660243515156114aa565b34156103ee57600080fd5b6103f66114f5565b604051600160a060020a03909116815260200160405180910390f35b341561041d57600080fd5b6101d3600160a060020a03600480358216916024359160443582169160643516906084359060ff60a435169061010460c46002604080519081016040528092919082600260200280828437509395505092359250611504915050565b341561048457600080fd5b6101d3600160a060020a036004351660243560ff6044351660643560843561187f565b34156104b257600080fd5b6103f6611947565b6101d3611956565b34156104cd57600080fd5b6101d3600160a060020a0360043516602435611978565b34156104ef57600080fd5b6101d3600435611b97565b341561050557600080fd5b61015c600160a060020a0360043516611bac565b341561052457600080fd5b6101d3600435611be6565b341561053a57600080fd5b6101d3600435611c30565b341561055057600080fd5b61015c611c45565b341561056357600080fd5b61015c600435611c4b565b341561057957600080fd5b61015c600160a060020a0360043516611c5d565b60076020526000908152604090205481565b600160a060020a0380821660009081526009602090815260408083209386168352929052205492915050565b60005433600160a060020a039081169116146105e657600080fd5b600054600160a060020a0380831691167fcbf985117192c8f614a58aaf97226bb80a754772f5f6edf06f87c675f2e6c66360405160405180910390a36000805473ffffffffffffffffffffffffffffffffffffffff1916600160a060020a0392909216919091179055565b60066020526000908152604090205460ff1681565b6000610670611ef9565b610678611f61565b60008054819081908190819033600160a060020a03908116911614806106b65750600160a060020a03331660009081526006602052604090205460ff165b15156106c157600080fd5b610120604051908101604052808d5181526020018d60016020020151815260200160408e0151815260200160608e0151815260200160808e0151815260200160a08e015181526020018c51600160a060020a031681526020018c60016020020151600160a060020a0316815260200160408d0151600160a060020a03169052965060606040519081016040528060c08e0151815260200160e08e0151815260200160608d0151600160a060020a0316905295503060c088015188518960e001518a602001518b604001518c606001518d61010001516040516c01000000000000000000000000600160a060020a03998a168102825297891688026014820152602881019690965293871686026048860152605c850192909252607c840152609c8301529092160260bc82015260d0016040519081900390209450848651876040015188602001516040519384526020840192909252600160a060020a03166c010000000000000000000000000260408084019190915260548301919091526074909101905180910390209350610869876101000151868c518c518d60015b602002015161187f565b15156108a657600080516020611fbe83398151915260095b8660405160ff909216825260208201526040908101905180910390a160009750610fee565b6108c186604001518560208d015160408d01518d600361085f565b15156108fe57600080516020611fbe833981519152600a5b8560405160ff909216825260208201526040908101905180910390a160009750610fee565b438760400151101561092057600080516020611fbe8339815191526004610881565b8551600960008960c00151600160a060020a0316600160a060020a0316815260200190815260200160002060008860400151600160a060020a0316600160a060020a0316815260200190815260200160002054101561098f57600080516020611fbe83398151915260016108d9565b6109b687516109aa88518a602001519063ffffffff611c6f16565b9063ffffffff611c9a16565b600960008960e00151600160a060020a0316600160a060020a031681526020019081526020016000206000896101000151600160a060020a0316600160a060020a03168152602001908152602001600020541015610a2457600080516020611fbe83398151915260006108d9565b60008481526003602052604090205460ff1615610a5157600080516020611fbe83398151915260066108d9565b8651610a7387516000888152600860205260409020549063ffffffff611cb116565b1115610a8f57600080516020611fbe83398151915260076108d9565b6000848152600360205260409020805460ff19166001179055670de0b6b3a7640000610ac7608089015188519063ffffffff611c6f16565b811515610ad057fe5b049250610aee87516109aa88518a602001519063ffffffff611c6f16565b9150670de0b6b3a7640000610b0e838960a001519063ffffffff611c6f16565b811515610b1757fe5b049050610b7883875103600960008a60c00151600160a060020a0316600160a060020a0316815260200190815260200160002060008a6101000151600160a060020a031681526020810191909152604001600020549063ffffffff611cb116565b600960008960c00151600160a060020a0316600160a060020a031681526020019081526020016000206000896101000151600160a060020a03168152602081019190915260400160002055610c1e8651600960008a60c00151600160a060020a0316600160a060020a0316815260200190815260200160002060008960400151600160a060020a031681526020810191909152604001600020549063ffffffff611cc016565b600960008960c00151600160a060020a0316600160a060020a0316815260200190815260200160002060008860400151600160a060020a0316600160a060020a0316815260200190815260200160002081905550610cb883600960008a60c00151600160a060020a03908116825260208083019390935260409182016000908120600154909216815292529020549063ffffffff611cb116565b600960008960c00151600160a060020a0390811682526020808301939093526040918201600090812060015490921681529252902055610d63610d0c88516109aa89518b602001519063ffffffff611c6f16565b600960008a60e00151600160a060020a0316600160a060020a0316815260200190815260200160002060008a6101000151600160a060020a031681526020810191909152604001600020549063ffffffff611cc016565b600960008960e00151600160a060020a0316600160a060020a031681526020019081526020016000206000896101000151600160a060020a03168152602081019190915260400160002055610e17610dc1838363ffffffff611cc016565b600960008a60e00151600160a060020a0316600160a060020a0316815260200190815260200160002060008960400151600160a060020a031681526020810191909152604001600020549063ffffffff611cb116565b600960008960e00151600160a060020a0316600160a060020a0316815260200190815260200160002060008860400151600160a060020a0316600160a060020a0316815260200190815260200160002081905550610eb181600960008a60e00151600160a060020a03908116825260208083019390935260409182016000908120600154909216815292529020549063ffffffff611cb116565b600960008960e00151600160a060020a0390811682526020808301939093526040918201600090812060015490921681529252902055610f0786516000878152600860205260409020549063ffffffff611cb116565b60008681526008602052604081209190915543906007906101008a0151600160a060020a0316600160a060020a031681526020019081526020016000208190555043600760008860400151600160a060020a031681526020810191909152604001600020557facba690473c099075a43376bf85b9c2560e51ee66a208fca92b31922d7d3a36760c08801518860e0015189610100015189604001518a518a604051600160a060020a0396871681529486166020860152928516604080860191909152919094166060840152608083019390935260a082015260c001905180910390a1600197505b50505050505050949350505050565b6000611007611f81565b600060e060405190810160405280895181526020018960016020020151815260200160408a0151815260200160608a015181526020018851600160a060020a031681526020018860016020020151600160a060020a031681526020016040890151600160a060020a03169052915030608083015183518460a001518560200151866040015187606001518860c001516040516c01000000000000000000000000600160a060020a03998a168102825297891688026014820152602881019690965293871686026048860152605c850192909252607c840152609c8301529092160260bc82015260d00160405180910390209050611107338288888861187f565b151561114357600080516020611fbe83398151915260088260405160ff909216825260208201526040908101905180910390a1600092506111fd565b81516000828152600860205260409020557f02b0e3138ecc5c8f35036d06f730effc5f7c9de5835667a26e297ae755ecfd5c608083015183518460a00151856020015186604001518760600151338d8d8d604051600160a060020a039a8b16815260208101999099529689166040808a01919091526060890196909652608088019490945260a087019290925290951660c085015260ff90941660e084015261010083019390935261012082015261014001905180910390a15b505095945050505050565b6000600160a060020a038316158015906112225750600082115b151561122d57600080fd5b6112378383611cd2565b5082600160a060020a03166323b872dd33308560006040516020015260405160e060020a63ffffffff8616028152600160a060020a0393841660048201529190921660248201526044810191909152606401602060405180830381600087803b15156112a257600080fd5b6102c65a03f115156112b357600080fd5b5050506040518051506001949350505050565b60046020526000908152604090205460ff1681565b600080888887896040519384526020840192909252600160a060020a03166c010000000000000000000000000260408084019190915260548301919091526074909101905180910390209050611334338287878761187f565b151561137057600080516020611fbe83398151915260088260405160ff909216825260208201526040908101905180910390a160009150611410565b60008181526003602052604090819020805460ff191660011790557fcbfd142f5f729c33b8c478622822ace3d8474d4eb5c7aaedb914f5d465fb8187908a908a908a908a9033908b908b908b90519788526020880196909652604080880195909552600160a060020a03938416606088015291909216608086015260ff90911660a085015260c084015260e0830191909152610100909101905180910390a15b50979650505050505050565b6000805433600160a060020a03908116911614806114525750600160a060020a03331660009081526006602052604090205460ff165b151561145d57600080fd5b5060018054600160a060020a03831673ffffffffffffffffffffffffffffffffffffffff19909116178155919050565b600960209081526000928352604080842090915290825290205481565b6000805433600160a060020a039081169116146114c657600080fd5b50600160a060020a0382166000908152600660205260409020805482151560ff19909116179055600192915050565b600154600160a060020a031681565b6000805481908190819033600160a060020a03908116911614806115405750600160a060020a03331660009081526006602052604090205460ff165b151561154b57600080fd5b308c8c8c8c8c6040516c01000000000000000000000000600160a060020a039788168102825295871686026014820152602881019490945291851684026048840152909316909102605c8201526070810191909152609001604051809103902092506115bd838d8d8d8d8c8c8c611da8565b15156115cc5760009350611870565b670de0b6b3a76400006115e5868d63ffffffff611c6f16565b8115156115ee57fe5b04915061164c8b600960008f600160a060020a0316600160a060020a0316815260200190815260200160002060008d600160a060020a0316600160a060020a0316815260200190815260200160002054611cc090919063ffffffff16565b600160a060020a038d811660009081526009602090815260408083208f85168452909152808220939093556001549091168152205461168b9083611cb1565b600160a060020a03808e166000908152600960209081526040808320600154909416835292905220556116c48b8363ffffffff611cc016565b9a50600960008d600160a060020a0316600160a060020a0316815260200190815260200160002060008b600160a060020a0316600160a060020a031681526020019081526020016000205490506000600160a060020a03168c600160a060020a0316141561176257600160a060020a0389168b156108fc028c604051600060405180830381858888f19350505050151561175d57600080fd5b6117db565b8b600160a060020a031663a9059cbb8a8d60006040516020015260405160e060020a63ffffffff8516028152600160a060020a0390921660048301526024820152604401602060405180830381600087803b15156117bf57600080fd5b6102c65a03f115156117d057600080fd5b505050604051805150505b6000838152600460209081526040808320805460ff19166001179055600160a060020a038d1683526007909152908190204390557f74217ce088f00bfd283666b763c64f0d1b1c345591dfdd01891dddf52446694e908d908c908e90859051600160a060020a0394851681529290931660208301526040808301919091526060820192909252608001905180910390a1600193505b50505098975050505050505050565b60006001856040517f19457468657265756d205369676e6564204d6573736167653a0a3332000000008152601c810191909152603c0160405180910390208585856040516000815260200160405260006040516020015260405193845260ff90921660208085019190915260408085019290925260608401929092526080909201915160208103908084039060008661646e5a03f1151561191f57600080fd5b505060206040510351600160a060020a031686600160a060020a031614905095945050505050565b600054600160a060020a031681565b60003481901161196557600080fd5b611970600034611cd2565b506001905090565b60025433600160a060020a03811660009081526007602052604081205490928392916119a5904390611cc0565b10156119b057600080fd5b600160a060020a038086166000908152600960209081526040808320339485168452909152902054909250849010156119e857600080fd5b600160a060020a03808616600090815260096020908152604080832093861683529290522054611a1e908563ffffffff611cc016565b600160a060020a038087166000818152600960209081526040808320948816835293905291909120919091551515611a8657600160a060020a03821684156108fc0285604051600060405180830381858888f193505050501515611a8157600080fd5b611b09565b84600160a060020a031663a9059cbb838660006040516020015260405160e060020a63ffffffff8516028152600160a060020a0390921660048301526024820152604401602060405180830381600087803b1515611ae357600080fd5b6102c65a03f11515611af457600080fd5b505050604051805190501515611b0957600080fd5b600160a060020a0380861660009081526009602090815260408083209386168352929052819020547fc1bc7b7b4880798c196e5adcf8dbbe770e901ad0bd5515af4dbc08d66eb18bfd918791859188919051600160a060020a0394851681529290931660208301526040808301919091526060820192909252608001905180910390a1506001949350505050565b60056020526000908152604090205460ff1681565b600160a060020a031660009081527fec8156718a8372b1db44bb411437d0870f3e3790d4a08526d024ce1b0b668f6b602052604090205490565b6000805433600160a060020a0390811691161480611c1c5750600160a060020a03331660009081526006602052604090205460ff165b1515611c2757600080fd5b50600255600190565b60036020526000908152604090205460ff1681565b60025481565b60086020526000908152604090205481565b600a6020526000908152604090205481565b6000828202831580611c8b5750828482811515611c8857fe5b04145b1515611c9357fe5b9392505050565b6000808284811515611ca857fe5b04949350505050565b600082820183811015611c9357fe5b600082821115611ccc57fe5b50900390565b600160a060020a03808316600090815260096020908152604080832033909416835292905290812054611d0b908363ffffffff611cb116565b600160a060020a038481166000908152600960209081526040808320339485168452808352818420958655600783529281902043905591905291547f4e3e4894f24a7c50bcb21d1ef785e34688bee05663c55d822eed7cefc253312392869291869151600160a060020a0394851681529290931660208301526040808301919091526060820192909252608001905180910390a150600192915050565b6000611dba868a86865187600161085f565b1515611df657600080516020611fbe83398151915260085b8a60405160ff909216825260208201526040908101905180910390a1506000611eed565b60008981526004602052604090205460ff1615611e2357600080516020611fbe8339815191526005611dd2565b600160a060020a038089166000908152600960209081526040808320938a168352929052205487901015611eac57600080516020611fbe8339815191526002600160a060020a03808b166000908152600960209081526040808320938c16835292905281902054905160ff909216825260208201526040908101905180910390a1506000611eed565b86821115611ee957600080516020611fbe83398151915260038360405160ff909216825260208201526040908101905180910390a1506000611eed565b5060015b98975050505050505050565b610120604051908101604052806000815260200160008152602001600081526020016000815260200160008152602001600081526020016000600160a060020a031681526020016000600160a060020a031681526020016000600160a060020a031681525090565b606060405190810160409081526000808352602083018190529082015290565b60e06040519081016040908152600080835260208301819052908201819052606082018190526080820181905260a0820181905260c082015290560036d86c59e00bd73dc19ba3adfe068e4b64ac7e92be35546adeddf1b956a87e90a165627a7a72305820825fcb74458483edb58e7e1891651c55844b0253308dc7182e0ad9c6be0e26dd0029`

// DeployExchange deploys a new Ethereum contract, binding an instance of Exchange to it.
func DeployExchange(auth *bind.TransactOpts, backend bind.ContractBackend, _feeAccount common.Address) (common.Address, *types.Transaction, *Exchange, error) {
	parsed, err := abi.JSON(strings.NewReader(ExchangeABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(ExchangeBin), backend, _feeAccount)
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

// EtherBalance is a free data retrieval call binding the contract method 0xcd0c5896.
//
// Solidity: function etherBalance(trader address) constant returns(uint256)
func (_Exchange *ExchangeCaller) EtherBalance(opts *bind.CallOpts, trader common.Address) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Exchange.contract.Call(opts, out, "etherBalance", trader)
	return *ret0, err
}

// EtherBalance is a free data retrieval call binding the contract method 0xcd0c5896.
//
// Solidity: function etherBalance(trader address) constant returns(uint256)
func (_Exchange *ExchangeSession) EtherBalance(trader common.Address) (*big.Int, error) {
	return _Exchange.Contract.EtherBalance(&_Exchange.CallOpts, trader)
}

// EtherBalance is a free data retrieval call binding the contract method 0xcd0c5896.
//
// Solidity: function etherBalance(trader address) constant returns(uint256)
func (_Exchange *ExchangeCallerSession) EtherBalance(trader common.Address) (*big.Int, error) {
	return _Exchange.Contract.EtherBalance(&_Exchange.CallOpts, trader)
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

// IsValidSignature is a free data retrieval call binding the contract method 0x8163681e.
//
// Solidity: function isValidSignature(signer address, hashedData bytes32, v uint8, r bytes32, s bytes32) constant returns(bool)
func (_Exchange *ExchangeCaller) IsValidSignature(opts *bind.CallOpts, signer common.Address, hashedData [32]byte, v uint8, r [32]byte, s [32]byte) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _Exchange.contract.Call(opts, out, "isValidSignature", signer, hashedData, v, r, s)
	return *ret0, err
}

// IsValidSignature is a free data retrieval call binding the contract method 0x8163681e.
//
// Solidity: function isValidSignature(signer address, hashedData bytes32, v uint8, r bytes32, s bytes32) constant returns(bool)
func (_Exchange *ExchangeSession) IsValidSignature(signer common.Address, hashedData [32]byte, v uint8, r [32]byte, s [32]byte) (bool, error) {
	return _Exchange.Contract.IsValidSignature(&_Exchange.CallOpts, signer, hashedData, v, r, s)
}

// IsValidSignature is a free data retrieval call binding the contract method 0x8163681e.
//
// Solidity: function isValidSignature(signer address, hashedData bytes32, v uint8, r bytes32, s bytes32) constant returns(bool)
func (_Exchange *ExchangeCallerSession) IsValidSignature(signer common.Address, hashedData [32]byte, v uint8, r [32]byte, s [32]byte) (bool, error) {
	return _Exchange.Contract.IsValidSignature(&_Exchange.CallOpts, signer, hashedData, v, r, s)
}

// LastTransaction is a free data retrieval call binding the contract method 0x0531be92.
//
// Solidity: function lastTransaction( address) constant returns(uint256)
func (_Exchange *ExchangeCaller) LastTransaction(opts *bind.CallOpts, arg0 common.Address) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Exchange.contract.Call(opts, out, "lastTransaction", arg0)
	return *ret0, err
}

// LastTransaction is a free data retrieval call binding the contract method 0x0531be92.
//
// Solidity: function lastTransaction( address) constant returns(uint256)
func (_Exchange *ExchangeSession) LastTransaction(arg0 common.Address) (*big.Int, error) {
	return _Exchange.Contract.LastTransaction(&_Exchange.CallOpts, arg0)
}

// LastTransaction is a free data retrieval call binding the contract method 0x0531be92.
//
// Solidity: function lastTransaction( address) constant returns(uint256)
func (_Exchange *ExchangeCallerSession) LastTransaction(arg0 common.Address) (*big.Int, error) {
	return _Exchange.Contract.LastTransaction(&_Exchange.CallOpts, arg0)
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

// OrderFills is a free data retrieval call binding the contract method 0xf7213db6.
//
// Solidity: function orderFills( bytes32) constant returns(uint256)
func (_Exchange *ExchangeCaller) OrderFills(opts *bind.CallOpts, arg0 [32]byte) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Exchange.contract.Call(opts, out, "orderFills", arg0)
	return *ret0, err
}

// OrderFills is a free data retrieval call binding the contract method 0xf7213db6.
//
// Solidity: function orderFills( bytes32) constant returns(uint256)
func (_Exchange *ExchangeSession) OrderFills(arg0 [32]byte) (*big.Int, error) {
	return _Exchange.Contract.OrderFills(&_Exchange.CallOpts, arg0)
}

// OrderFills is a free data retrieval call binding the contract method 0xf7213db6.
//
// Solidity: function orderFills( bytes32) constant returns(uint256)
func (_Exchange *ExchangeCallerSession) OrderFills(arg0 [32]byte) (*big.Int, error) {
	return _Exchange.Contract.OrderFills(&_Exchange.CallOpts, arg0)
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

// ProtectedFunds is a free data retrieval call binding the contract method 0xfe2e2b94.
//
// Solidity: function protectedFunds( address) constant returns(uint256)
func (_Exchange *ExchangeCaller) ProtectedFunds(opts *bind.CallOpts, arg0 common.Address) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Exchange.contract.Call(opts, out, "protectedFunds", arg0)
	return *ret0, err
}

// ProtectedFunds is a free data retrieval call binding the contract method 0xfe2e2b94.
//
// Solidity: function protectedFunds( address) constant returns(uint256)
func (_Exchange *ExchangeSession) ProtectedFunds(arg0 common.Address) (*big.Int, error) {
	return _Exchange.Contract.ProtectedFunds(&_Exchange.CallOpts, arg0)
}

// ProtectedFunds is a free data retrieval call binding the contract method 0xfe2e2b94.
//
// Solidity: function protectedFunds( address) constant returns(uint256)
func (_Exchange *ExchangeCallerSession) ProtectedFunds(arg0 common.Address) (*big.Int, error) {
	return _Exchange.Contract.ProtectedFunds(&_Exchange.CallOpts, arg0)
}

// TokenBalance is a free data retrieval call binding the contract method 0x1049334f.
//
// Solidity: function tokenBalance(trader address, token address) constant returns(uint256)
func (_Exchange *ExchangeCaller) TokenBalance(opts *bind.CallOpts, trader common.Address, token common.Address) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Exchange.contract.Call(opts, out, "tokenBalance", trader, token)
	return *ret0, err
}

// TokenBalance is a free data retrieval call binding the contract method 0x1049334f.
//
// Solidity: function tokenBalance(trader address, token address) constant returns(uint256)
func (_Exchange *ExchangeSession) TokenBalance(trader common.Address, token common.Address) (*big.Int, error) {
	return _Exchange.Contract.TokenBalance(&_Exchange.CallOpts, trader, token)
}

// TokenBalance is a free data retrieval call binding the contract method 0x1049334f.
//
// Solidity: function tokenBalance(trader address, token address) constant returns(uint256)
func (_Exchange *ExchangeCallerSession) TokenBalance(trader common.Address, token common.Address) (*big.Int, error) {
	return _Exchange.Contract.TokenBalance(&_Exchange.CallOpts, trader, token)
}

// Tokens is a free data retrieval call binding the contract method 0x508493bc.
//
// Solidity: function tokens( address,  address) constant returns(uint256)
func (_Exchange *ExchangeCaller) Tokens(opts *bind.CallOpts, arg0 common.Address, arg1 common.Address) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Exchange.contract.Call(opts, out, "tokens", arg0, arg1)
	return *ret0, err
}

// Tokens is a free data retrieval call binding the contract method 0x508493bc.
//
// Solidity: function tokens( address,  address) constant returns(uint256)
func (_Exchange *ExchangeSession) Tokens(arg0 common.Address, arg1 common.Address) (*big.Int, error) {
	return _Exchange.Contract.Tokens(&_Exchange.CallOpts, arg0, arg1)
}

// Tokens is a free data retrieval call binding the contract method 0x508493bc.
//
// Solidity: function tokens( address,  address) constant returns(uint256)
func (_Exchange *ExchangeCallerSession) Tokens(arg0 common.Address, arg1 common.Address) (*big.Int, error) {
	return _Exchange.Contract.Tokens(&_Exchange.CallOpts, arg0, arg1)
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

// Transferred is a free data retrieval call binding the contract method 0xafc441e3.
//
// Solidity: function transferred( bytes32) constant returns(bool)
func (_Exchange *ExchangeCaller) Transferred(opts *bind.CallOpts, arg0 [32]byte) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _Exchange.contract.Call(opts, out, "transferred", arg0)
	return *ret0, err
}

// Transferred is a free data retrieval call binding the contract method 0xafc441e3.
//
// Solidity: function transferred( bytes32) constant returns(bool)
func (_Exchange *ExchangeSession) Transferred(arg0 [32]byte) (bool, error) {
	return _Exchange.Contract.Transferred(&_Exchange.CallOpts, arg0)
}

// Transferred is a free data retrieval call binding the contract method 0xafc441e3.
//
// Solidity: function transferred( bytes32) constant returns(bool)
func (_Exchange *ExchangeCallerSession) Transferred(arg0 [32]byte) (bool, error) {
	return _Exchange.Contract.Transferred(&_Exchange.CallOpts, arg0)
}

// WithdrawalSecurityPeriod is a free data retrieval call binding the contract method 0xf3198b16.
//
// Solidity: function withdrawalSecurityPeriod() constant returns(uint256)
func (_Exchange *ExchangeCaller) WithdrawalSecurityPeriod(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Exchange.contract.Call(opts, out, "withdrawalSecurityPeriod")
	return *ret0, err
}

// WithdrawalSecurityPeriod is a free data retrieval call binding the contract method 0xf3198b16.
//
// Solidity: function withdrawalSecurityPeriod() constant returns(uint256)
func (_Exchange *ExchangeSession) WithdrawalSecurityPeriod() (*big.Int, error) {
	return _Exchange.Contract.WithdrawalSecurityPeriod(&_Exchange.CallOpts)
}

// WithdrawalSecurityPeriod is a free data retrieval call binding the contract method 0xf3198b16.
//
// Solidity: function withdrawalSecurityPeriod() constant returns(uint256)
func (_Exchange *ExchangeCallerSession) WithdrawalSecurityPeriod() (*big.Int, error) {
	return _Exchange.Contract.WithdrawalSecurityPeriod(&_Exchange.CallOpts)
}

// Withdrawn is a free data retrieval call binding the contract method 0x3823d66c.
//
// Solidity: function withdrawn( bytes32) constant returns(bool)
func (_Exchange *ExchangeCaller) Withdrawn(opts *bind.CallOpts, arg0 [32]byte) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _Exchange.contract.Call(opts, out, "withdrawn", arg0)
	return *ret0, err
}

// Withdrawn is a free data retrieval call binding the contract method 0x3823d66c.
//
// Solidity: function withdrawn( bytes32) constant returns(bool)
func (_Exchange *ExchangeSession) Withdrawn(arg0 [32]byte) (bool, error) {
	return _Exchange.Contract.Withdrawn(&_Exchange.CallOpts, arg0)
}

// Withdrawn is a free data retrieval call binding the contract method 0x3823d66c.
//
// Solidity: function withdrawn( bytes32) constant returns(bool)
func (_Exchange *ExchangeCallerSession) Withdrawn(arg0 [32]byte) (bool, error) {
	return _Exchange.Contract.Withdrawn(&_Exchange.CallOpts, arg0)
}

// CancelOrder is a paid mutator transaction binding the contract method 0x2e813def.
//
// Solidity: function cancelOrder(orderValues uint256[5], orderAddresses address[4], v uint8, r bytes32, s bytes32) returns(bool)
func (_Exchange *ExchangeTransactor) CancelOrder(opts *bind.TransactOpts, orderValues [5]*big.Int, orderAddresses [4]common.Address, v uint8, r [32]byte, s [32]byte) (*types.Transaction, error) {
	return _Exchange.contract.Transact(opts, "cancelOrder", orderValues, orderAddresses, v, r, s)
}

// CancelOrder is a paid mutator transaction binding the contract method 0x2e813def.
//
// Solidity: function cancelOrder(orderValues uint256[5], orderAddresses address[4], v uint8, r bytes32, s bytes32) returns(bool)
func (_Exchange *ExchangeSession) CancelOrder(orderValues [5]*big.Int, orderAddresses [4]common.Address, v uint8, r [32]byte, s [32]byte) (*types.Transaction, error) {
	return _Exchange.Contract.CancelOrder(&_Exchange.TransactOpts, orderValues, orderAddresses, v, r, s)
}

// CancelOrder is a paid mutator transaction binding the contract method 0x2e813def.
//
// Solidity: function cancelOrder(orderValues uint256[5], orderAddresses address[4], v uint8, r bytes32, s bytes32) returns(bool)
func (_Exchange *ExchangeTransactorSession) CancelOrder(orderValues [5]*big.Int, orderAddresses [4]common.Address, v uint8, r [32]byte, s [32]byte) (*types.Transaction, error) {
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

// DepositEther is a paid mutator transaction binding the contract method 0x98ea5fca.
//
// Solidity: function depositEther() returns(bool)
func (_Exchange *ExchangeTransactor) DepositEther(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Exchange.contract.Transact(opts, "depositEther")
}

// DepositEther is a paid mutator transaction binding the contract method 0x98ea5fca.
//
// Solidity: function depositEther() returns(bool)
func (_Exchange *ExchangeSession) DepositEther() (*types.Transaction, error) {
	return _Exchange.Contract.DepositEther(&_Exchange.TransactOpts)
}

// DepositEther is a paid mutator transaction binding the contract method 0x98ea5fca.
//
// Solidity: function depositEther() returns(bool)
func (_Exchange *ExchangeTransactorSession) DepositEther() (*types.Transaction, error) {
	return _Exchange.Contract.DepositEther(&_Exchange.TransactOpts)
}

// DepositToken is a paid mutator transaction binding the contract method 0x338b5dea.
//
// Solidity: function depositToken(token address, amount uint256) returns(bool)
func (_Exchange *ExchangeTransactor) DepositToken(opts *bind.TransactOpts, token common.Address, amount *big.Int) (*types.Transaction, error) {
	return _Exchange.contract.Transact(opts, "depositToken", token, amount)
}

// DepositToken is a paid mutator transaction binding the contract method 0x338b5dea.
//
// Solidity: function depositToken(token address, amount uint256) returns(bool)
func (_Exchange *ExchangeSession) DepositToken(token common.Address, amount *big.Int) (*types.Transaction, error) {
	return _Exchange.Contract.DepositToken(&_Exchange.TransactOpts, token, amount)
}

// DepositToken is a paid mutator transaction binding the contract method 0x338b5dea.
//
// Solidity: function depositToken(token address, amount uint256) returns(bool)
func (_Exchange *ExchangeTransactorSession) DepositToken(token common.Address, amount *big.Int) (*types.Transaction, error) {
	return _Exchange.Contract.DepositToken(&_Exchange.TransactOpts, token, amount)
}

// ExecuteTrade is a paid mutator transaction binding the contract method 0x2207148d.
//
// Solidity: function executeTrade(orderValues uint256[8], orderAddresses address[4], v uint8[2], rs bytes32[4]) returns(bool)
func (_Exchange *ExchangeTransactor) ExecuteTrade(opts *bind.TransactOpts, orderValues [8]*big.Int, orderAddresses [4]common.Address, v [2]uint8, rs [4][32]byte) (*types.Transaction, error) {
	return _Exchange.contract.Transact(opts, "executeTrade", orderValues, orderAddresses, v, rs)
}

// ExecuteTrade is a paid mutator transaction binding the contract method 0x2207148d.
//
// Solidity: function executeTrade(orderValues uint256[8], orderAddresses address[4], v uint8[2], rs bytes32[4]) returns(bool)
func (_Exchange *ExchangeSession) ExecuteTrade(orderValues [8]*big.Int, orderAddresses [4]common.Address, v [2]uint8, rs [4][32]byte) (*types.Transaction, error) {
	return _Exchange.Contract.ExecuteTrade(&_Exchange.TransactOpts, orderValues, orderAddresses, v, rs)
}

// ExecuteTrade is a paid mutator transaction binding the contract method 0x2207148d.
//
// Solidity: function executeTrade(orderValues uint256[8], orderAddresses address[4], v uint8[2], rs bytes32[4]) returns(bool)
func (_Exchange *ExchangeTransactorSession) ExecuteTrade(orderValues [8]*big.Int, orderAddresses [4]common.Address, v [2]uint8, rs [4][32]byte) (*types.Transaction, error) {
	return _Exchange.Contract.ExecuteTrade(&_Exchange.TransactOpts, orderValues, orderAddresses, v, rs)
}

// SecurityWithdraw is a paid mutator transaction binding the contract method 0x9bc2c131.
//
// Solidity: function securityWithdraw(token address, amount uint256) returns(bool)
func (_Exchange *ExchangeTransactor) SecurityWithdraw(opts *bind.TransactOpts, token common.Address, amount *big.Int) (*types.Transaction, error) {
	return _Exchange.contract.Transact(opts, "securityWithdraw", token, amount)
}

// SecurityWithdraw is a paid mutator transaction binding the contract method 0x9bc2c131.
//
// Solidity: function securityWithdraw(token address, amount uint256) returns(bool)
func (_Exchange *ExchangeSession) SecurityWithdraw(token common.Address, amount *big.Int) (*types.Transaction, error) {
	return _Exchange.Contract.SecurityWithdraw(&_Exchange.TransactOpts, token, amount)
}

// SecurityWithdraw is a paid mutator transaction binding the contract method 0x9bc2c131.
//
// Solidity: function securityWithdraw(token address, amount uint256) returns(bool)
func (_Exchange *ExchangeTransactorSession) SecurityWithdraw(token common.Address, amount *big.Int) (*types.Transaction, error) {
	return _Exchange.Contract.SecurityWithdraw(&_Exchange.TransactOpts, token, amount)
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
// Solidity: function setOperator(operator address, isOperator bool) returns(bool)
func (_Exchange *ExchangeTransactor) SetOperator(opts *bind.TransactOpts, operator common.Address, isOperator bool) (*types.Transaction, error) {
	return _Exchange.contract.Transact(opts, "setOperator", operator, isOperator)
}

// SetOperator is a paid mutator transaction binding the contract method 0x558a7297.
//
// Solidity: function setOperator(operator address, isOperator bool) returns(bool)
func (_Exchange *ExchangeSession) SetOperator(operator common.Address, isOperator bool) (*types.Transaction, error) {
	return _Exchange.Contract.SetOperator(&_Exchange.TransactOpts, operator, isOperator)
}

// SetOperator is a paid mutator transaction binding the contract method 0x558a7297.
//
// Solidity: function setOperator(operator address, isOperator bool) returns(bool)
func (_Exchange *ExchangeTransactorSession) SetOperator(operator common.Address, isOperator bool) (*types.Transaction, error) {
	return _Exchange.Contract.SetOperator(&_Exchange.TransactOpts, operator, isOperator)
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

// SetWithdrawalSecurityPeriod is a paid mutator transaction binding the contract method 0xcf46d3d7.
//
// Solidity: function setWithdrawalSecurityPeriod(_withdrawalSecurityPeriod uint256) returns(bool)
func (_Exchange *ExchangeTransactor) SetWithdrawalSecurityPeriod(opts *bind.TransactOpts, _withdrawalSecurityPeriod *big.Int) (*types.Transaction, error) {
	return _Exchange.contract.Transact(opts, "setWithdrawalSecurityPeriod", _withdrawalSecurityPeriod)
}

// SetWithdrawalSecurityPeriod is a paid mutator transaction binding the contract method 0xcf46d3d7.
//
// Solidity: function setWithdrawalSecurityPeriod(_withdrawalSecurityPeriod uint256) returns(bool)
func (_Exchange *ExchangeSession) SetWithdrawalSecurityPeriod(_withdrawalSecurityPeriod *big.Int) (*types.Transaction, error) {
	return _Exchange.Contract.SetWithdrawalSecurityPeriod(&_Exchange.TransactOpts, _withdrawalSecurityPeriod)
}

// SetWithdrawalSecurityPeriod is a paid mutator transaction binding the contract method 0xcf46d3d7.
//
// Solidity: function setWithdrawalSecurityPeriod(_withdrawalSecurityPeriod uint256) returns(bool)
func (_Exchange *ExchangeTransactorSession) SetWithdrawalSecurityPeriod(_withdrawalSecurityPeriod *big.Int) (*types.Transaction, error) {
	return _Exchange.Contract.SetWithdrawalSecurityPeriod(&_Exchange.TransactOpts, _withdrawalSecurityPeriod)
}

// Withdraw is a paid mutator transaction binding the contract method 0x74af3ab3.
//
// Solidity: function withdraw(token address, amount uint256, trader address, receiver address, nonce uint256, v uint8, rs bytes32[2], feeWithdrawal uint256) returns(bool)
func (_Exchange *ExchangeTransactor) Withdraw(opts *bind.TransactOpts, token common.Address, amount *big.Int, trader common.Address, receiver common.Address, nonce *big.Int, v uint8, rs [2][32]byte, feeWithdrawal *big.Int) (*types.Transaction, error) {
	return _Exchange.contract.Transact(opts, "withdraw", token, amount, trader, receiver, nonce, v, rs, feeWithdrawal)
}

// Withdraw is a paid mutator transaction binding the contract method 0x74af3ab3.
//
// Solidity: function withdraw(token address, amount uint256, trader address, receiver address, nonce uint256, v uint8, rs bytes32[2], feeWithdrawal uint256) returns(bool)
func (_Exchange *ExchangeSession) Withdraw(token common.Address, amount *big.Int, trader common.Address, receiver common.Address, nonce *big.Int, v uint8, rs [2][32]byte, feeWithdrawal *big.Int) (*types.Transaction, error) {
	return _Exchange.Contract.Withdraw(&_Exchange.TransactOpts, token, amount, trader, receiver, nonce, v, rs, feeWithdrawal)
}

// Withdraw is a paid mutator transaction binding the contract method 0x74af3ab3.
//
// Solidity: function withdraw(token address, amount uint256, trader address, receiver address, nonce uint256, v uint8, rs bytes32[2], feeWithdrawal uint256) returns(bool)
func (_Exchange *ExchangeTransactorSession) Withdraw(token common.Address, amount *big.Int, trader common.Address, receiver common.Address, nonce *big.Int, v uint8, rs [2][32]byte, feeWithdrawal *big.Int) (*types.Transaction, error) {
	return _Exchange.Contract.Withdraw(&_Exchange.TransactOpts, token, amount, trader, receiver, nonce, v, rs, feeWithdrawal)
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
	TokenBuy   common.Address
	AmountBuy  *big.Int
	TokenSell  common.Address
	AmountSell *big.Int
	Expires    *big.Int
	Nonce      *big.Int
	Maker      common.Address
	V          uint8
	R          [32]byte
	S          [32]byte
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterLogCancelOrder is a free log retrieval operation binding the contract event 0x02b0e3138ecc5c8f35036d06f730effc5f7c9de5835667a26e297ae755ecfd5c.
//
// Solidity: e LogCancelOrder(tokenBuy address, amountBuy uint256, tokenSell address, amountSell uint256, expires uint256, nonce uint256, maker address, v uint8, r bytes32, s bytes32)
func (_Exchange *ExchangeFilterer) FilterLogCancelOrder(opts *bind.FilterOpts) (*ExchangeLogCancelOrderIterator, error) {

	logs, sub, err := _Exchange.contract.FilterLogs(opts, "LogCancelOrder")
	if err != nil {
		return nil, err
	}
	return &ExchangeLogCancelOrderIterator{contract: _Exchange.contract, event: "LogCancelOrder", logs: logs, sub: sub}, nil
}

// WatchLogCancelOrder is a free log subscription operation binding the contract event 0x02b0e3138ecc5c8f35036d06f730effc5f7c9de5835667a26e297ae755ecfd5c.
//
// Solidity: e LogCancelOrder(tokenBuy address, amountBuy uint256, tokenSell address, amountSell uint256, expires uint256, nonce uint256, maker address, v uint8, r bytes32, s bytes32)
func (_Exchange *ExchangeFilterer) WatchLogCancelOrder(opts *bind.WatchOpts, sink chan<- *ExchangeLogCancelOrder) (event.Subscription, error) {

	logs, sub, err := _Exchange.contract.WatchLogs(opts, "LogCancelOrder")
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
	Sender     common.Address
	V          uint8
	R          [32]byte
	S          [32]byte
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterLogCancelTrade is a free log retrieval operation binding the contract event 0xcbfd142f5f729c33b8c478622822ace3d8474d4eb5c7aaedb914f5d465fb8187.
//
// Solidity: e LogCancelTrade(orderHash bytes32, amount uint256, tradeNonce uint256, taker address, sender address, v uint8, r bytes32, s bytes32)
func (_Exchange *ExchangeFilterer) FilterLogCancelTrade(opts *bind.FilterOpts) (*ExchangeLogCancelTradeIterator, error) {

	logs, sub, err := _Exchange.contract.FilterLogs(opts, "LogCancelTrade")
	if err != nil {
		return nil, err
	}
	return &ExchangeLogCancelTradeIterator{contract: _Exchange.contract, event: "LogCancelTrade", logs: logs, sub: sub}, nil
}

// WatchLogCancelTrade is a free log subscription operation binding the contract event 0xcbfd142f5f729c33b8c478622822ace3d8474d4eb5c7aaedb914f5d465fb8187.
//
// Solidity: e LogCancelTrade(orderHash bytes32, amount uint256, tradeNonce uint256, taker address, sender address, v uint8, r bytes32, s bytes32)
func (_Exchange *ExchangeFilterer) WatchLogCancelTrade(opts *bind.WatchOpts, sink chan<- *ExchangeLogCancelTrade) (event.Subscription, error) {

	logs, sub, err := _Exchange.contract.WatchLogs(opts, "LogCancelTrade")
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

// ExchangeLogDepositIterator is returned from FilterLogDeposit and is used to iterate over the raw logs and unpacked data for LogDeposit events raised by the Exchange contract.
type ExchangeLogDepositIterator struct {
	Event *ExchangeLogDeposit // Event containing the contract specifics and raw log

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
func (it *ExchangeLogDepositIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ExchangeLogDeposit)
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
		it.Event = new(ExchangeLogDeposit)
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
func (it *ExchangeLogDepositIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ExchangeLogDepositIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ExchangeLogDeposit represents a LogDeposit event raised by the Exchange contract.
type ExchangeLogDeposit struct {
	Token   common.Address
	User    common.Address
	Amount  *big.Int
	Balance *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterLogDeposit is a free log retrieval operation binding the contract event 0x4e3e4894f24a7c50bcb21d1ef785e34688bee05663c55d822eed7cefc2533123.
//
// Solidity: e LogDeposit(token address, user address, amount uint256, balance uint256)
func (_Exchange *ExchangeFilterer) FilterLogDeposit(opts *bind.FilterOpts) (*ExchangeLogDepositIterator, error) {

	logs, sub, err := _Exchange.contract.FilterLogs(opts, "LogDeposit")
	if err != nil {
		return nil, err
	}
	return &ExchangeLogDepositIterator{contract: _Exchange.contract, event: "LogDeposit", logs: logs, sub: sub}, nil
}

// WatchLogDeposit is a free log subscription operation binding the contract event 0x4e3e4894f24a7c50bcb21d1ef785e34688bee05663c55d822eed7cefc2533123.
//
// Solidity: e LogDeposit(token address, user address, amount uint256, balance uint256)
func (_Exchange *ExchangeFilterer) WatchLogDeposit(opts *bind.WatchOpts, sink chan<- *ExchangeLogDeposit) (event.Subscription, error) {

	logs, sub, err := _Exchange.contract.WatchLogs(opts, "LogDeposit")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ExchangeLogDeposit)
				if err := _Exchange.contract.UnpackLog(event, "LogDeposit", log); err != nil {
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
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterLogError is a free log retrieval operation binding the contract event 0x36d86c59e00bd73dc19ba3adfe068e4b64ac7e92be35546adeddf1b956a87e90.
//
// Solidity: e LogError(errorId uint8, orderHash bytes32)
func (_Exchange *ExchangeFilterer) FilterLogError(opts *bind.FilterOpts) (*ExchangeLogErrorIterator, error) {

	logs, sub, err := _Exchange.contract.FilterLogs(opts, "LogError")
	if err != nil {
		return nil, err
	}
	return &ExchangeLogErrorIterator{contract: _Exchange.contract, event: "LogError", logs: logs, sub: sub}, nil
}

// WatchLogError is a free log subscription operation binding the contract event 0x36d86c59e00bd73dc19ba3adfe068e4b64ac7e92be35546adeddf1b956a87e90.
//
// Solidity: e LogError(errorId uint8, orderHash bytes32)
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

// ExchangeLogSecurityWithdrawIterator is returned from FilterLogSecurityWithdraw and is used to iterate over the raw logs and unpacked data for LogSecurityWithdraw events raised by the Exchange contract.
type ExchangeLogSecurityWithdrawIterator struct {
	Event *ExchangeLogSecurityWithdraw // Event containing the contract specifics and raw log

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
func (it *ExchangeLogSecurityWithdrawIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ExchangeLogSecurityWithdraw)
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
		it.Event = new(ExchangeLogSecurityWithdraw)
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
func (it *ExchangeLogSecurityWithdrawIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ExchangeLogSecurityWithdrawIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ExchangeLogSecurityWithdraw represents a LogSecurityWithdraw event raised by the Exchange contract.
type ExchangeLogSecurityWithdraw struct {
	Token   common.Address
	User    common.Address
	Amount  *big.Int
	Balance *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterLogSecurityWithdraw is a free log retrieval operation binding the contract event 0xc1bc7b7b4880798c196e5adcf8dbbe770e901ad0bd5515af4dbc08d66eb18bfd.
//
// Solidity: e LogSecurityWithdraw(token address, user address, amount uint256, balance uint256)
func (_Exchange *ExchangeFilterer) FilterLogSecurityWithdraw(opts *bind.FilterOpts) (*ExchangeLogSecurityWithdrawIterator, error) {

	logs, sub, err := _Exchange.contract.FilterLogs(opts, "LogSecurityWithdraw")
	if err != nil {
		return nil, err
	}
	return &ExchangeLogSecurityWithdrawIterator{contract: _Exchange.contract, event: "LogSecurityWithdraw", logs: logs, sub: sub}, nil
}

// WatchLogSecurityWithdraw is a free log subscription operation binding the contract event 0xc1bc7b7b4880798c196e5adcf8dbbe770e901ad0bd5515af4dbc08d66eb18bfd.
//
// Solidity: e LogSecurityWithdraw(token address, user address, amount uint256, balance uint256)
func (_Exchange *ExchangeFilterer) WatchLogSecurityWithdraw(opts *bind.WatchOpts, sink chan<- *ExchangeLogSecurityWithdraw) (event.Subscription, error) {

	logs, sub, err := _Exchange.contract.WatchLogs(opts, "LogSecurityWithdraw")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ExchangeLogSecurityWithdraw)
				if err := _Exchange.contract.UnpackLog(event, "LogSecurityWithdraw", log); err != nil {
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
	TokenBuy  common.Address
	TokenSell common.Address
	Maker     common.Address
	Taker     common.Address
	Amount    *big.Int
	Hsh       [32]byte
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterLogTrade is a free log retrieval operation binding the contract event 0xacba690473c099075a43376bf85b9c2560e51ee66a208fca92b31922d7d3a367.
//
// Solidity: e LogTrade(tokenBuy address, tokenSell address, maker address, taker address, amount uint256, hsh bytes32)
func (_Exchange *ExchangeFilterer) FilterLogTrade(opts *bind.FilterOpts) (*ExchangeLogTradeIterator, error) {

	logs, sub, err := _Exchange.contract.FilterLogs(opts, "LogTrade")
	if err != nil {
		return nil, err
	}
	return &ExchangeLogTradeIterator{contract: _Exchange.contract, event: "LogTrade", logs: logs, sub: sub}, nil
}

// WatchLogTrade is a free log subscription operation binding the contract event 0xacba690473c099075a43376bf85b9c2560e51ee66a208fca92b31922d7d3a367.
//
// Solidity: e LogTrade(tokenBuy address, tokenSell address, maker address, taker address, amount uint256, hsh bytes32)
func (_Exchange *ExchangeFilterer) WatchLogTrade(opts *bind.WatchOpts, sink chan<- *ExchangeLogTrade) (event.Subscription, error) {

	logs, sub, err := _Exchange.contract.WatchLogs(opts, "LogTrade")
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

// ExchangeLogTransferIterator is returned from FilterLogTransfer and is used to iterate over the raw logs and unpacked data for LogTransfer events raised by the Exchange contract.
type ExchangeLogTransferIterator struct {
	Event *ExchangeLogTransfer // Event containing the contract specifics and raw log

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
func (it *ExchangeLogTransferIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ExchangeLogTransfer)
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
		it.Event = new(ExchangeLogTransfer)
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
func (it *ExchangeLogTransferIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ExchangeLogTransferIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ExchangeLogTransfer represents a LogTransfer event raised by the Exchange contract.
type ExchangeLogTransfer struct {
	Token     common.Address
	Recipient common.Address
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterLogTransfer is a free log retrieval operation binding the contract event 0x5d517a7dfb872efa300109cebcb8235ae90602d926b499065254b47383395426.
//
// Solidity: e LogTransfer(token address, recipient address)
func (_Exchange *ExchangeFilterer) FilterLogTransfer(opts *bind.FilterOpts) (*ExchangeLogTransferIterator, error) {

	logs, sub, err := _Exchange.contract.FilterLogs(opts, "LogTransfer")
	if err != nil {
		return nil, err
	}
	return &ExchangeLogTransferIterator{contract: _Exchange.contract, event: "LogTransfer", logs: logs, sub: sub}, nil
}

// WatchLogTransfer is a free log subscription operation binding the contract event 0x5d517a7dfb872efa300109cebcb8235ae90602d926b499065254b47383395426.
//
// Solidity: e LogTransfer(token address, recipient address)
func (_Exchange *ExchangeFilterer) WatchLogTransfer(opts *bind.WatchOpts, sink chan<- *ExchangeLogTransfer) (event.Subscription, error) {

	logs, sub, err := _Exchange.contract.WatchLogs(opts, "LogTransfer")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ExchangeLogTransfer)
				if err := _Exchange.contract.UnpackLog(event, "LogTransfer", log); err != nil {
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

// ExchangeLogWithdrawIterator is returned from FilterLogWithdraw and is used to iterate over the raw logs and unpacked data for LogWithdraw events raised by the Exchange contract.
type ExchangeLogWithdrawIterator struct {
	Event *ExchangeLogWithdraw // Event containing the contract specifics and raw log

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
func (it *ExchangeLogWithdrawIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ExchangeLogWithdraw)
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
		it.Event = new(ExchangeLogWithdraw)
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
func (it *ExchangeLogWithdrawIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ExchangeLogWithdrawIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ExchangeLogWithdraw represents a LogWithdraw event raised by the Exchange contract.
type ExchangeLogWithdraw struct {
	Token   common.Address
	User    common.Address
	Amount  *big.Int
	Balance *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterLogWithdraw is a free log retrieval operation binding the contract event 0x74217ce088f00bfd283666b763c64f0d1b1c345591dfdd01891dddf52446694e.
//
// Solidity: e LogWithdraw(token address, user address, amount uint256, balance uint256)
func (_Exchange *ExchangeFilterer) FilterLogWithdraw(opts *bind.FilterOpts) (*ExchangeLogWithdrawIterator, error) {

	logs, sub, err := _Exchange.contract.FilterLogs(opts, "LogWithdraw")
	if err != nil {
		return nil, err
	}
	return &ExchangeLogWithdrawIterator{contract: _Exchange.contract, event: "LogWithdraw", logs: logs, sub: sub}, nil
}

// WatchLogWithdraw is a free log subscription operation binding the contract event 0x74217ce088f00bfd283666b763c64f0d1b1c345591dfdd01891dddf52446694e.
//
// Solidity: e LogWithdraw(token address, user address, amount uint256, balance uint256)
func (_Exchange *ExchangeFilterer) WatchLogWithdraw(opts *bind.WatchOpts, sink chan<- *ExchangeLogWithdraw) (event.Subscription, error) {

	logs, sub, err := _Exchange.contract.WatchLogs(opts, "LogWithdraw")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ExchangeLogWithdraw)
				if err := _Exchange.contract.UnpackLog(event, "LogWithdraw", log); err != nil {
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
const OwnedBin = `0x6060604052341561000f57600080fd5b60008054600160a060020a033316600160a060020a03199091161790556101588061003b6000396000f300606060405263ffffffff7c010000000000000000000000000000000000000000000000000000000060003504166313af403581146100475780638da5cb5b1461006857600080fd5b341561005257600080fd5b610066600160a060020a0360043516610097565b005b341561007357600080fd5b61007b61011d565b604051600160a060020a03909116815260200160405180910390f35b60005433600160a060020a039081169116146100b257600080fd5b600054600160a060020a0380831691167fcbf985117192c8f614a58aaf97226bb80a754772f5f6edf06f87c675f2e6c66360405160405180910390a36000805473ffffffffffffffffffffffffffffffffffffffff1916600160a060020a0392909216919091179055565b600054600160a060020a0316815600a165627a7a723058204aeeaf4e4dc4759d8bc5a2a939359d5659068f66f4cc519daa70c3ff03ef383b0029`

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
const SafeMathBin = `0x60606040523415600e57600080fd5b603580601b6000396000f3006060604052600080fd00a165627a7a72305820b1664f360b40a53f8cc7b1134106cf8d0e761c5d257318747c569c3ca4670e0b0029`

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
