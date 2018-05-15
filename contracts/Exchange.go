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

// ExchangeABI is the input ABI used to generate the binding from.
const ExchangeABI = "[{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"address\"}],\"name\":\"lastTransaction\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"setOwner\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"address\"}],\"name\":\"operators\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"withdrawn\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"address\"},{\"name\":\"\",\"type\":\"address\"}],\"name\":\"tokens\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"feeAccount\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"transferred\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"traded\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"withdrawalSecurityPeriod\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"orderFills\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"address\"}],\"name\":\"protectedFunds\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"name\":\"_feeAccount\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"fallback\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"tokenBuy\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"tokenSell\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"maker\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"taker\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"hsh\",\"type\":\"bytes32\"}],\"name\":\"LogTrade\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"token\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"user\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"balance\",\"type\":\"uint256\"}],\"name\":\"LogDeposit\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"token\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"user\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"balance\",\"type\":\"uint256\"}],\"name\":\"LogWithdraw\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"token\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"user\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"balance\",\"type\":\"uint256\"}],\"name\":\"LogSecurityWithdraw\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"token\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"recipient\",\"type\":\"address\"}],\"name\":\"LogTransfer\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"tokenBuy\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"amountBuy\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"tokenSell\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"amountSell\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"expires\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"nonce\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"maker\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"v\",\"type\":\"uint8\"},{\"indexed\":false,\"name\":\"r\",\"type\":\"bytes32\"},{\"indexed\":false,\"name\":\"s\",\"type\":\"bytes32\"}],\"name\":\"LogCancelOrder\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"orderHash\",\"type\":\"bytes32\"},{\"indexed\":false,\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"tradeNonce\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"taker\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"v\",\"type\":\"uint8\"},{\"indexed\":false,\"name\":\"r\",\"type\":\"bytes32\"},{\"indexed\":false,\"name\":\"s\",\"type\":\"bytes32\"}],\"name\":\"LogCancelTrade\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"errorId\",\"type\":\"uint8\"},{\"indexed\":false,\"name\":\"orderHash\",\"type\":\"bytes32\"}],\"name\":\"LogError\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"SetOwner\",\"type\":\"event\"},{\"constant\":false,\"inputs\":[{\"name\":\"_feeAccount\",\"type\":\"address\"}],\"name\":\"setFeeAccount\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"operator\",\"type\":\"address\"},{\"name\":\"isOperator\",\"type\":\"bool\"}],\"name\":\"setOperator\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_withdrawalSecurityPeriod\",\"type\":\"uint256\"}],\"name\":\"setWithdrawalSecurityPeriod\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"depositEther\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"token\",\"type\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"depositToken\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"trader\",\"type\":\"address\"},{\"name\":\"token\",\"type\":\"address\"}],\"name\":\"tokenBalance\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"trader\",\"type\":\"address\"}],\"name\":\"etherBalance\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"token\",\"type\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"securityWithdraw\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"token\",\"type\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\"},{\"name\":\"trader\",\"type\":\"address\"},{\"name\":\"receiver\",\"type\":\"address\"},{\"name\":\"nonce\",\"type\":\"uint256\"},{\"name\":\"v\",\"type\":\"uint8\"},{\"name\":\"rs\",\"type\":\"bytes32[2]\"},{\"name\":\"feeWithdrawal\",\"type\":\"uint256\"}],\"name\":\"withdraw\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"orderValues\",\"type\":\"uint256[8]\"},{\"name\":\"orderAddresses\",\"type\":\"address[4]\"},{\"name\":\"v\",\"type\":\"uint8[2]\"},{\"name\":\"rs\",\"type\":\"bytes32[4]\"}],\"name\":\"executeTrade\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"orderValues\",\"type\":\"uint256[5]\"},{\"name\":\"orderAddresses\",\"type\":\"address[4]\"},{\"name\":\"v\",\"type\":\"uint8\"},{\"name\":\"r\",\"type\":\"bytes32\"},{\"name\":\"s\",\"type\":\"bytes32\"}],\"name\":\"cancelOrder\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"orderHash\",\"type\":\"bytes32\"},{\"name\":\"amount\",\"type\":\"uint256\"},{\"name\":\"tradeNonce\",\"type\":\"uint256\"},{\"name\":\"taker\",\"type\":\"address\"},{\"name\":\"v\",\"type\":\"uint8\"},{\"name\":\"r\",\"type\":\"bytes32\"},{\"name\":\"s\",\"type\":\"bytes32\"}],\"name\":\"cancelTrade\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"signer\",\"type\":\"address\"},{\"name\":\"hashedData\",\"type\":\"bytes32\"},{\"name\":\"v\",\"type\":\"uint8\"},{\"name\":\"r\",\"type\":\"bytes32\"},{\"name\":\"s\",\"type\":\"bytes32\"}],\"name\":\"isValidSignature\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"}]"

// Exchange is an auto generated Go binding around an Ethereum contract.
type Exchange struct {
	ExchangeCaller     // Read-only binding to the contract
	ExchangeTransactor // Write-only binding to the contract
}

// ExchangeCaller is an auto generated read-only Go binding around an Ethereum contract.
type ExchangeCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ExchangeTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ExchangeTransactor struct {
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
	contract, err := bindExchange(address, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Exchange{ExchangeCaller: ExchangeCaller{contract: contract}, ExchangeTransactor: ExchangeTransactor{contract: contract}}, nil
}

// NewExchangeCaller creates a new read-only instance of Exchange, bound to a specific deployed contract.
func NewExchangeCaller(address common.Address, caller bind.ContractCaller) (*ExchangeCaller, error) {
	contract, err := bindExchange(address, caller, nil)
	if err != nil {
		return nil, err
	}
	return &ExchangeCaller{contract: contract}, nil
}

// NewExchangeTransactor creates a new write-only instance of Exchange, bound to a specific deployed contract.
func NewExchangeTransactor(address common.Address, transactor bind.ContractTransactor) (*ExchangeTransactor, error) {
	contract, err := bindExchange(address, nil, transactor)
	if err != nil {
		return nil, err
	}
	return &ExchangeTransactor{contract: contract}, nil
}

// bindExchange binds a generic wrapper to an already deployed contract.
func bindExchange(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(ExchangeABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor), nil
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
