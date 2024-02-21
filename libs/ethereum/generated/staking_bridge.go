// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package generated

import (
	"errors"
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
)

// StakingBridgeMetaData contains all meta data concerning the StakingBridge contract.
var StakingBridgeMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"vega_public_key\",\"type\":\"bytes32\"}],\"name\":\"stake\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
}

// StakingBridgeABI is the input ABI used to generate the binding from.
// Deprecated: Use StakingBridgeMetaData.ABI instead.
var StakingBridgeABI = StakingBridgeMetaData.ABI

// StakingBridge is an auto generated Go binding around an Ethereum contract.
type StakingBridge struct {
	StakingBridgeCaller     // Read-only binding to the contract
	StakingBridgeTransactor // Write-only binding to the contract
	StakingBridgeFilterer   // Log filterer for contract events
}

// StakingBridgeCaller is an auto generated read-only Go binding around an Ethereum contract.
type StakingBridgeCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// StakingBridgeTransactor is an auto generated write-only Go binding around an Ethereum contract.
type StakingBridgeTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// StakingBridgeFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type StakingBridgeFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// StakingBridgeSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type StakingBridgeSession struct {
	Contract     *StakingBridge    // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// StakingBridgeCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type StakingBridgeCallerSession struct {
	Contract *StakingBridgeCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts        // Call options to use throughout this session
}

// StakingBridgeTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type StakingBridgeTransactorSession struct {
	Contract     *StakingBridgeTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts        // Transaction auth options to use throughout this session
}

// StakingBridgeRaw is an auto generated low-level Go binding around an Ethereum contract.
type StakingBridgeRaw struct {
	Contract *StakingBridge // Generic contract binding to access the raw methods on
}

// StakingBridgeCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type StakingBridgeCallerRaw struct {
	Contract *StakingBridgeCaller // Generic read-only contract binding to access the raw methods on
}

// StakingBridgeTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type StakingBridgeTransactorRaw struct {
	Contract *StakingBridgeTransactor // Generic write-only contract binding to access the raw methods on
}

// NewStakingBridge creates a new instance of StakingBridge, bound to a specific deployed contract.
func NewStakingBridge(address common.Address, backend bind.ContractBackend) (*StakingBridge, error) {
	contract, err := bindStakingBridge(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &StakingBridge{StakingBridgeCaller: StakingBridgeCaller{contract: contract}, StakingBridgeTransactor: StakingBridgeTransactor{contract: contract}, StakingBridgeFilterer: StakingBridgeFilterer{contract: contract}}, nil
}

// NewStakingBridgeCaller creates a new read-only instance of StakingBridge, bound to a specific deployed contract.
func NewStakingBridgeCaller(address common.Address, caller bind.ContractCaller) (*StakingBridgeCaller, error) {
	contract, err := bindStakingBridge(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &StakingBridgeCaller{contract: contract}, nil
}

// NewStakingBridgeTransactor creates a new write-only instance of StakingBridge, bound to a specific deployed contract.
func NewStakingBridgeTransactor(address common.Address, transactor bind.ContractTransactor) (*StakingBridgeTransactor, error) {
	contract, err := bindStakingBridge(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &StakingBridgeTransactor{contract: contract}, nil
}

// NewStakingBridgeFilterer creates a new log filterer instance of StakingBridge, bound to a specific deployed contract.
func NewStakingBridgeFilterer(address common.Address, filterer bind.ContractFilterer) (*StakingBridgeFilterer, error) {
	contract, err := bindStakingBridge(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &StakingBridgeFilterer{contract: contract}, nil
}

// bindStakingBridge binds a generic wrapper to an already deployed contract.
func bindStakingBridge(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(StakingBridgeABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_StakingBridge *StakingBridgeRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _StakingBridge.Contract.StakingBridgeCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_StakingBridge *StakingBridgeRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _StakingBridge.Contract.StakingBridgeTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_StakingBridge *StakingBridgeRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _StakingBridge.Contract.StakingBridgeTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_StakingBridge *StakingBridgeCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _StakingBridge.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_StakingBridge *StakingBridgeTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _StakingBridge.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_StakingBridge *StakingBridgeTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _StakingBridge.Contract.contract.Transact(opts, method, params...)
}

// Stake is a paid mutator transaction binding the contract method 0x83c592cf.
//
// Solidity: function stake(uint256 amount, bytes32 vega_public_key) returns()
func (_StakingBridge *StakingBridgeTransactor) Stake(opts *bind.TransactOpts, amount *big.Int, vega_public_key [32]byte) (*types.Transaction, error) {
	return _StakingBridge.contract.Transact(opts, "stake", amount, vega_public_key)
}

// Stake is a paid mutator transaction binding the contract method 0x83c592cf.
//
// Solidity: function stake(uint256 amount, bytes32 vega_public_key) returns()
func (_StakingBridge *StakingBridgeSession) Stake(amount *big.Int, vega_public_key [32]byte) (*types.Transaction, error) {
	return _StakingBridge.Contract.Stake(&_StakingBridge.TransactOpts, amount, vega_public_key)
}

// Stake is a paid mutator transaction binding the contract method 0x83c592cf.
//
// Solidity: function stake(uint256 amount, bytes32 vega_public_key) returns()
func (_StakingBridge *StakingBridgeTransactorSession) Stake(amount *big.Int, vega_public_key [32]byte) (*types.Transaction, error) {
	return _StakingBridge.Contract.Stake(&_StakingBridge.TransactOpts, amount, vega_public_key)
}
