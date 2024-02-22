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

// ERC20BridgeMetaData contains all meta data concerning the ERC20Bridge contract.
var ERC20BridgeMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"addresspayable\",\"name\":\"erc20_asset_pool\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"user_address\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"asset_source\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"vega_public_key\",\"type\":\"bytes32\"}],\"name\":\"Asset_Deposited\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"asset_source\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"lifetime_limit\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"withdraw_threshold\",\"type\":\"uint256\"}],\"name\":\"Asset_Limits_Updated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"asset_source\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"vega_asset_id\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"nonce\",\"type\":\"uint256\"}],\"name\":\"Asset_Listed\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"asset_source\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"nonce\",\"type\":\"uint256\"}],\"name\":\"Asset_Removed\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"user_address\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"asset_source\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"nonce\",\"type\":\"uint256\"}],\"name\":\"Asset_Withdrawn\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[],\"name\":\"Bridge_Resumed\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[],\"name\":\"Bridge_Stopped\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"withdraw_delay\",\"type\":\"uint256\"}],\"name\":\"Bridge_Withdraw_Delay_Set\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"depositor\",\"type\":\"address\"}],\"name\":\"Depositor_Exempted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"depositor\",\"type\":\"address\"}],\"name\":\"Depositor_Exemption_Revoked\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"default_withdraw_delay\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"asset_source\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"vega_public_key\",\"type\":\"bytes32\"}],\"name\":\"deposit_asset\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"erc20_asset_pool_address\",\"outputs\":[{\"internalType\":\"addresspayable\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"exempt_depositor\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"asset_source\",\"type\":\"address\"}],\"name\":\"get_asset_deposit_lifetime_limit\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"vega_asset_id\",\"type\":\"bytes32\"}],\"name\":\"get_asset_source\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"get_multisig_control_address\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"asset_source\",\"type\":\"address\"}],\"name\":\"get_vega_asset_id\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"asset_source\",\"type\":\"address\"}],\"name\":\"get_withdraw_threshold\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"nonce\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"signatures\",\"type\":\"bytes\"}],\"name\":\"global_resume\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"nonce\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"signatures\",\"type\":\"bytes\"}],\"name\":\"global_stop\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"asset_source\",\"type\":\"address\"}],\"name\":\"is_asset_listed\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"depositor\",\"type\":\"address\"}],\"name\":\"is_exempt_depositor\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"is_stopped\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"asset_source\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"vega_asset_id\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"lifetime_limit\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"withdraw_threshold\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"nonce\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"signatures\",\"type\":\"bytes\"}],\"name\":\"list_asset\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"asset_source\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"nonce\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"signatures\",\"type\":\"bytes\"}],\"name\":\"remove_asset\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"revoke_exempt_depositor\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"asset_source\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"lifetime_limit\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"threshold\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"nonce\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"signatures\",\"type\":\"bytes\"}],\"name\":\"set_asset_limits\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"delay\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"nonce\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"signatures\",\"type\":\"bytes\"}],\"name\":\"set_withdraw_delay\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"asset_source\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"target\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"creation\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"nonce\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"signatures\",\"type\":\"bytes\"}],\"name\":\"withdraw_asset\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
}

// ERC20BridgeABI is the input ABI used to generate the binding from.
// Deprecated: Use ERC20BridgeMetaData.ABI instead.
var ERC20BridgeABI = ERC20BridgeMetaData.ABI

// ERC20Bridge is an auto generated Go binding around an Ethereum contract.
type ERC20Bridge struct {
	ERC20BridgeCaller     // Read-only binding to the contract
	ERC20BridgeTransactor // Write-only binding to the contract
	ERC20BridgeFilterer   // Log filterer for contract events
}

// ERC20BridgeCaller is an auto generated read-only Go binding around an Ethereum contract.
type ERC20BridgeCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ERC20BridgeTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ERC20BridgeTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ERC20BridgeFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ERC20BridgeFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ERC20BridgeSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ERC20BridgeSession struct {
	Contract     *ERC20Bridge      // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// ERC20BridgeCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ERC20BridgeCallerSession struct {
	Contract *ERC20BridgeCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts      // Call options to use throughout this session
}

// ERC20BridgeTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ERC20BridgeTransactorSession struct {
	Contract     *ERC20BridgeTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts      // Transaction auth options to use throughout this session
}

// ERC20BridgeRaw is an auto generated low-level Go binding around an Ethereum contract.
type ERC20BridgeRaw struct {
	Contract *ERC20Bridge // Generic contract binding to access the raw methods on
}

// ERC20BridgeCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ERC20BridgeCallerRaw struct {
	Contract *ERC20BridgeCaller // Generic read-only contract binding to access the raw methods on
}

// ERC20BridgeTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ERC20BridgeTransactorRaw struct {
	Contract *ERC20BridgeTransactor // Generic write-only contract binding to access the raw methods on
}

// NewERC20Bridge creates a new instance of ERC20Bridge, bound to a specific deployed contract.
func NewERC20Bridge(address common.Address, backend bind.ContractBackend) (*ERC20Bridge, error) {
	contract, err := bindERC20Bridge(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ERC20Bridge{ERC20BridgeCaller: ERC20BridgeCaller{contract: contract}, ERC20BridgeTransactor: ERC20BridgeTransactor{contract: contract}, ERC20BridgeFilterer: ERC20BridgeFilterer{contract: contract}}, nil
}

// NewERC20BridgeCaller creates a new read-only instance of ERC20Bridge, bound to a specific deployed contract.
func NewERC20BridgeCaller(address common.Address, caller bind.ContractCaller) (*ERC20BridgeCaller, error) {
	contract, err := bindERC20Bridge(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ERC20BridgeCaller{contract: contract}, nil
}

// NewERC20BridgeTransactor creates a new write-only instance of ERC20Bridge, bound to a specific deployed contract.
func NewERC20BridgeTransactor(address common.Address, transactor bind.ContractTransactor) (*ERC20BridgeTransactor, error) {
	contract, err := bindERC20Bridge(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ERC20BridgeTransactor{contract: contract}, nil
}

// NewERC20BridgeFilterer creates a new log filterer instance of ERC20Bridge, bound to a specific deployed contract.
func NewERC20BridgeFilterer(address common.Address, filterer bind.ContractFilterer) (*ERC20BridgeFilterer, error) {
	contract, err := bindERC20Bridge(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ERC20BridgeFilterer{contract: contract}, nil
}

// bindERC20Bridge binds a generic wrapper to an already deployed contract.
func bindERC20Bridge(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(ERC20BridgeABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ERC20Bridge *ERC20BridgeRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ERC20Bridge.Contract.ERC20BridgeCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ERC20Bridge *ERC20BridgeRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ERC20Bridge.Contract.ERC20BridgeTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ERC20Bridge *ERC20BridgeRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ERC20Bridge.Contract.ERC20BridgeTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ERC20Bridge *ERC20BridgeCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ERC20Bridge.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ERC20Bridge *ERC20BridgeTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ERC20Bridge.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ERC20Bridge *ERC20BridgeTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ERC20Bridge.Contract.contract.Transact(opts, method, params...)
}

// DefaultWithdrawDelay is a free data retrieval call binding the contract method 0x3f4f199d.
//
// Solidity: function default_withdraw_delay() view returns(uint256)
func (_ERC20Bridge *ERC20BridgeCaller) DefaultWithdrawDelay(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _ERC20Bridge.contract.Call(opts, &out, "default_withdraw_delay")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// DefaultWithdrawDelay is a free data retrieval call binding the contract method 0x3f4f199d.
//
// Solidity: function default_withdraw_delay() view returns(uint256)
func (_ERC20Bridge *ERC20BridgeSession) DefaultWithdrawDelay() (*big.Int, error) {
	return _ERC20Bridge.Contract.DefaultWithdrawDelay(&_ERC20Bridge.CallOpts)
}

// DefaultWithdrawDelay is a free data retrieval call binding the contract method 0x3f4f199d.
//
// Solidity: function default_withdraw_delay() view returns(uint256)
func (_ERC20Bridge *ERC20BridgeCallerSession) DefaultWithdrawDelay() (*big.Int, error) {
	return _ERC20Bridge.Contract.DefaultWithdrawDelay(&_ERC20Bridge.CallOpts)
}

// Erc20AssetPoolAddress is a free data retrieval call binding the contract method 0x9356aab8.
//
// Solidity: function erc20_asset_pool_address() view returns(address)
func (_ERC20Bridge *ERC20BridgeCaller) Erc20AssetPoolAddress(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ERC20Bridge.contract.Call(opts, &out, "erc20_asset_pool_address")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Erc20AssetPoolAddress is a free data retrieval call binding the contract method 0x9356aab8.
//
// Solidity: function erc20_asset_pool_address() view returns(address)
func (_ERC20Bridge *ERC20BridgeSession) Erc20AssetPoolAddress() (common.Address, error) {
	return _ERC20Bridge.Contract.Erc20AssetPoolAddress(&_ERC20Bridge.CallOpts)
}

// Erc20AssetPoolAddress is a free data retrieval call binding the contract method 0x9356aab8.
//
// Solidity: function erc20_asset_pool_address() view returns(address)
func (_ERC20Bridge *ERC20BridgeCallerSession) Erc20AssetPoolAddress() (common.Address, error) {
	return _ERC20Bridge.Contract.Erc20AssetPoolAddress(&_ERC20Bridge.CallOpts)
}

// GetAssetDepositLifetimeLimit is a free data retrieval call binding the contract method 0x354a897a.
//
// Solidity: function get_asset_deposit_lifetime_limit(address asset_source) view returns(uint256)
func (_ERC20Bridge *ERC20BridgeCaller) GetAssetDepositLifetimeLimit(opts *bind.CallOpts, asset_source common.Address) (*big.Int, error) {
	var out []interface{}
	err := _ERC20Bridge.contract.Call(opts, &out, "get_asset_deposit_lifetime_limit", asset_source)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetAssetDepositLifetimeLimit is a free data retrieval call binding the contract method 0x354a897a.
//
// Solidity: function get_asset_deposit_lifetime_limit(address asset_source) view returns(uint256)
func (_ERC20Bridge *ERC20BridgeSession) GetAssetDepositLifetimeLimit(asset_source common.Address) (*big.Int, error) {
	return _ERC20Bridge.Contract.GetAssetDepositLifetimeLimit(&_ERC20Bridge.CallOpts, asset_source)
}

// GetAssetDepositLifetimeLimit is a free data retrieval call binding the contract method 0x354a897a.
//
// Solidity: function get_asset_deposit_lifetime_limit(address asset_source) view returns(uint256)
func (_ERC20Bridge *ERC20BridgeCallerSession) GetAssetDepositLifetimeLimit(asset_source common.Address) (*big.Int, error) {
	return _ERC20Bridge.Contract.GetAssetDepositLifetimeLimit(&_ERC20Bridge.CallOpts, asset_source)
}

// GetAssetSource is a free data retrieval call binding the contract method 0x786b0bc0.
//
// Solidity: function get_asset_source(bytes32 vega_asset_id) view returns(address)
func (_ERC20Bridge *ERC20BridgeCaller) GetAssetSource(opts *bind.CallOpts, vega_asset_id [32]byte) (common.Address, error) {
	var out []interface{}
	err := _ERC20Bridge.contract.Call(opts, &out, "get_asset_source", vega_asset_id)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetAssetSource is a free data retrieval call binding the contract method 0x786b0bc0.
//
// Solidity: function get_asset_source(bytes32 vega_asset_id) view returns(address)
func (_ERC20Bridge *ERC20BridgeSession) GetAssetSource(vega_asset_id [32]byte) (common.Address, error) {
	return _ERC20Bridge.Contract.GetAssetSource(&_ERC20Bridge.CallOpts, vega_asset_id)
}

// GetAssetSource is a free data retrieval call binding the contract method 0x786b0bc0.
//
// Solidity: function get_asset_source(bytes32 vega_asset_id) view returns(address)
func (_ERC20Bridge *ERC20BridgeCallerSession) GetAssetSource(vega_asset_id [32]byte) (common.Address, error) {
	return _ERC20Bridge.Contract.GetAssetSource(&_ERC20Bridge.CallOpts, vega_asset_id)
}

// GetMultisigControlAddress is a free data retrieval call binding the contract method 0xc58dc3b9.
//
// Solidity: function get_multisig_control_address() view returns(address)
func (_ERC20Bridge *ERC20BridgeCaller) GetMultisigControlAddress(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ERC20Bridge.contract.Call(opts, &out, "get_multisig_control_address")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetMultisigControlAddress is a free data retrieval call binding the contract method 0xc58dc3b9.
//
// Solidity: function get_multisig_control_address() view returns(address)
func (_ERC20Bridge *ERC20BridgeSession) GetMultisigControlAddress() (common.Address, error) {
	return _ERC20Bridge.Contract.GetMultisigControlAddress(&_ERC20Bridge.CallOpts)
}

// GetMultisigControlAddress is a free data retrieval call binding the contract method 0xc58dc3b9.
//
// Solidity: function get_multisig_control_address() view returns(address)
func (_ERC20Bridge *ERC20BridgeCallerSession) GetMultisigControlAddress() (common.Address, error) {
	return _ERC20Bridge.Contract.GetMultisigControlAddress(&_ERC20Bridge.CallOpts)
}

// GetVegaAssetId is a free data retrieval call binding the contract method 0xa06b5d39.
//
// Solidity: function get_vega_asset_id(address asset_source) view returns(bytes32)
func (_ERC20Bridge *ERC20BridgeCaller) GetVegaAssetId(opts *bind.CallOpts, asset_source common.Address) ([32]byte, error) {
	var out []interface{}
	err := _ERC20Bridge.contract.Call(opts, &out, "get_vega_asset_id", asset_source)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// GetVegaAssetId is a free data retrieval call binding the contract method 0xa06b5d39.
//
// Solidity: function get_vega_asset_id(address asset_source) view returns(bytes32)
func (_ERC20Bridge *ERC20BridgeSession) GetVegaAssetId(asset_source common.Address) ([32]byte, error) {
	return _ERC20Bridge.Contract.GetVegaAssetId(&_ERC20Bridge.CallOpts, asset_source)
}

// GetVegaAssetId is a free data retrieval call binding the contract method 0xa06b5d39.
//
// Solidity: function get_vega_asset_id(address asset_source) view returns(bytes32)
func (_ERC20Bridge *ERC20BridgeCallerSession) GetVegaAssetId(asset_source common.Address) ([32]byte, error) {
	return _ERC20Bridge.Contract.GetVegaAssetId(&_ERC20Bridge.CallOpts, asset_source)
}

// GetWithdrawThreshold is a free data retrieval call binding the contract method 0xe8a7bce0.
//
// Solidity: function get_withdraw_threshold(address asset_source) view returns(uint256)
func (_ERC20Bridge *ERC20BridgeCaller) GetWithdrawThreshold(opts *bind.CallOpts, asset_source common.Address) (*big.Int, error) {
	var out []interface{}
	err := _ERC20Bridge.contract.Call(opts, &out, "get_withdraw_threshold", asset_source)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetWithdrawThreshold is a free data retrieval call binding the contract method 0xe8a7bce0.
//
// Solidity: function get_withdraw_threshold(address asset_source) view returns(uint256)
func (_ERC20Bridge *ERC20BridgeSession) GetWithdrawThreshold(asset_source common.Address) (*big.Int, error) {
	return _ERC20Bridge.Contract.GetWithdrawThreshold(&_ERC20Bridge.CallOpts, asset_source)
}

// GetWithdrawThreshold is a free data retrieval call binding the contract method 0xe8a7bce0.
//
// Solidity: function get_withdraw_threshold(address asset_source) view returns(uint256)
func (_ERC20Bridge *ERC20BridgeCallerSession) GetWithdrawThreshold(asset_source common.Address) (*big.Int, error) {
	return _ERC20Bridge.Contract.GetWithdrawThreshold(&_ERC20Bridge.CallOpts, asset_source)
}

// IsAssetListed is a free data retrieval call binding the contract method 0x7fd27b7f.
//
// Solidity: function is_asset_listed(address asset_source) view returns(bool)
func (_ERC20Bridge *ERC20BridgeCaller) IsAssetListed(opts *bind.CallOpts, asset_source common.Address) (bool, error) {
	var out []interface{}
	err := _ERC20Bridge.contract.Call(opts, &out, "is_asset_listed", asset_source)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsAssetListed is a free data retrieval call binding the contract method 0x7fd27b7f.
//
// Solidity: function is_asset_listed(address asset_source) view returns(bool)
func (_ERC20Bridge *ERC20BridgeSession) IsAssetListed(asset_source common.Address) (bool, error) {
	return _ERC20Bridge.Contract.IsAssetListed(&_ERC20Bridge.CallOpts, asset_source)
}

// IsAssetListed is a free data retrieval call binding the contract method 0x7fd27b7f.
//
// Solidity: function is_asset_listed(address asset_source) view returns(bool)
func (_ERC20Bridge *ERC20BridgeCallerSession) IsAssetListed(asset_source common.Address) (bool, error) {
	return _ERC20Bridge.Contract.IsAssetListed(&_ERC20Bridge.CallOpts, asset_source)
}

// IsExemptDepositor is a free data retrieval call binding the contract method 0x15c0df9d.
//
// Solidity: function is_exempt_depositor(address depositor) view returns(bool)
func (_ERC20Bridge *ERC20BridgeCaller) IsExemptDepositor(opts *bind.CallOpts, depositor common.Address) (bool, error) {
	var out []interface{}
	err := _ERC20Bridge.contract.Call(opts, &out, "is_exempt_depositor", depositor)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsExemptDepositor is a free data retrieval call binding the contract method 0x15c0df9d.
//
// Solidity: function is_exempt_depositor(address depositor) view returns(bool)
func (_ERC20Bridge *ERC20BridgeSession) IsExemptDepositor(depositor common.Address) (bool, error) {
	return _ERC20Bridge.Contract.IsExemptDepositor(&_ERC20Bridge.CallOpts, depositor)
}

// IsExemptDepositor is a free data retrieval call binding the contract method 0x15c0df9d.
//
// Solidity: function is_exempt_depositor(address depositor) view returns(bool)
func (_ERC20Bridge *ERC20BridgeCallerSession) IsExemptDepositor(depositor common.Address) (bool, error) {
	return _ERC20Bridge.Contract.IsExemptDepositor(&_ERC20Bridge.CallOpts, depositor)
}

// IsStopped is a free data retrieval call binding the contract method 0xe272e9d0.
//
// Solidity: function is_stopped() view returns(bool)
func (_ERC20Bridge *ERC20BridgeCaller) IsStopped(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _ERC20Bridge.contract.Call(opts, &out, "is_stopped")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsStopped is a free data retrieval call binding the contract method 0xe272e9d0.
//
// Solidity: function is_stopped() view returns(bool)
func (_ERC20Bridge *ERC20BridgeSession) IsStopped() (bool, error) {
	return _ERC20Bridge.Contract.IsStopped(&_ERC20Bridge.CallOpts)
}

// IsStopped is a free data retrieval call binding the contract method 0xe272e9d0.
//
// Solidity: function is_stopped() view returns(bool)
func (_ERC20Bridge *ERC20BridgeCallerSession) IsStopped() (bool, error) {
	return _ERC20Bridge.Contract.IsStopped(&_ERC20Bridge.CallOpts)
}

// DepositAsset is a paid mutator transaction binding the contract method 0xf7683932.
//
// Solidity: function deposit_asset(address asset_source, uint256 amount, bytes32 vega_public_key) returns()
func (_ERC20Bridge *ERC20BridgeTransactor) DepositAsset(opts *bind.TransactOpts, asset_source common.Address, amount *big.Int, vega_public_key [32]byte) (*types.Transaction, error) {
	return _ERC20Bridge.contract.Transact(opts, "deposit_asset", asset_source, amount, vega_public_key)
}

// DepositAsset is a paid mutator transaction binding the contract method 0xf7683932.
//
// Solidity: function deposit_asset(address asset_source, uint256 amount, bytes32 vega_public_key) returns()
func (_ERC20Bridge *ERC20BridgeSession) DepositAsset(asset_source common.Address, amount *big.Int, vega_public_key [32]byte) (*types.Transaction, error) {
	return _ERC20Bridge.Contract.DepositAsset(&_ERC20Bridge.TransactOpts, asset_source, amount, vega_public_key)
}

// DepositAsset is a paid mutator transaction binding the contract method 0xf7683932.
//
// Solidity: function deposit_asset(address asset_source, uint256 amount, bytes32 vega_public_key) returns()
func (_ERC20Bridge *ERC20BridgeTransactorSession) DepositAsset(asset_source common.Address, amount *big.Int, vega_public_key [32]byte) (*types.Transaction, error) {
	return _ERC20Bridge.Contract.DepositAsset(&_ERC20Bridge.TransactOpts, asset_source, amount, vega_public_key)
}

// ExemptDepositor is a paid mutator transaction binding the contract method 0xb76fbb75.
//
// Solidity: function exempt_depositor() returns()
func (_ERC20Bridge *ERC20BridgeTransactor) ExemptDepositor(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ERC20Bridge.contract.Transact(opts, "exempt_depositor")
}

// ExemptDepositor is a paid mutator transaction binding the contract method 0xb76fbb75.
//
// Solidity: function exempt_depositor() returns()
func (_ERC20Bridge *ERC20BridgeSession) ExemptDepositor() (*types.Transaction, error) {
	return _ERC20Bridge.Contract.ExemptDepositor(&_ERC20Bridge.TransactOpts)
}

// ExemptDepositor is a paid mutator transaction binding the contract method 0xb76fbb75.
//
// Solidity: function exempt_depositor() returns()
func (_ERC20Bridge *ERC20BridgeTransactorSession) ExemptDepositor() (*types.Transaction, error) {
	return _ERC20Bridge.Contract.ExemptDepositor(&_ERC20Bridge.TransactOpts)
}

// GlobalResume is a paid mutator transaction binding the contract method 0xd72ed529.
//
// Solidity: function global_resume(uint256 nonce, bytes signatures) returns()
func (_ERC20Bridge *ERC20BridgeTransactor) GlobalResume(opts *bind.TransactOpts, nonce *big.Int, signatures []byte) (*types.Transaction, error) {
	return _ERC20Bridge.contract.Transact(opts, "global_resume", nonce, signatures)
}

// GlobalResume is a paid mutator transaction binding the contract method 0xd72ed529.
//
// Solidity: function global_resume(uint256 nonce, bytes signatures) returns()
func (_ERC20Bridge *ERC20BridgeSession) GlobalResume(nonce *big.Int, signatures []byte) (*types.Transaction, error) {
	return _ERC20Bridge.Contract.GlobalResume(&_ERC20Bridge.TransactOpts, nonce, signatures)
}

// GlobalResume is a paid mutator transaction binding the contract method 0xd72ed529.
//
// Solidity: function global_resume(uint256 nonce, bytes signatures) returns()
func (_ERC20Bridge *ERC20BridgeTransactorSession) GlobalResume(nonce *big.Int, signatures []byte) (*types.Transaction, error) {
	return _ERC20Bridge.Contract.GlobalResume(&_ERC20Bridge.TransactOpts, nonce, signatures)
}

// GlobalStop is a paid mutator transaction binding the contract method 0x9dfd3c88.
//
// Solidity: function global_stop(uint256 nonce, bytes signatures) returns()
func (_ERC20Bridge *ERC20BridgeTransactor) GlobalStop(opts *bind.TransactOpts, nonce *big.Int, signatures []byte) (*types.Transaction, error) {
	return _ERC20Bridge.contract.Transact(opts, "global_stop", nonce, signatures)
}

// GlobalStop is a paid mutator transaction binding the contract method 0x9dfd3c88.
//
// Solidity: function global_stop(uint256 nonce, bytes signatures) returns()
func (_ERC20Bridge *ERC20BridgeSession) GlobalStop(nonce *big.Int, signatures []byte) (*types.Transaction, error) {
	return _ERC20Bridge.Contract.GlobalStop(&_ERC20Bridge.TransactOpts, nonce, signatures)
}

// GlobalStop is a paid mutator transaction binding the contract method 0x9dfd3c88.
//
// Solidity: function global_stop(uint256 nonce, bytes signatures) returns()
func (_ERC20Bridge *ERC20BridgeTransactorSession) GlobalStop(nonce *big.Int, signatures []byte) (*types.Transaction, error) {
	return _ERC20Bridge.Contract.GlobalStop(&_ERC20Bridge.TransactOpts, nonce, signatures)
}

// ListAsset is a paid mutator transaction binding the contract method 0x0ff3562c.
//
// Solidity: function list_asset(address asset_source, bytes32 vega_asset_id, uint256 lifetime_limit, uint256 withdraw_threshold, uint256 nonce, bytes signatures) returns()
func (_ERC20Bridge *ERC20BridgeTransactor) ListAsset(opts *bind.TransactOpts, asset_source common.Address, vega_asset_id [32]byte, lifetime_limit *big.Int, withdraw_threshold *big.Int, nonce *big.Int, signatures []byte) (*types.Transaction, error) {
	return _ERC20Bridge.contract.Transact(opts, "list_asset", asset_source, vega_asset_id, lifetime_limit, withdraw_threshold, nonce, signatures)
}

// ListAsset is a paid mutator transaction binding the contract method 0x0ff3562c.
//
// Solidity: function list_asset(address asset_source, bytes32 vega_asset_id, uint256 lifetime_limit, uint256 withdraw_threshold, uint256 nonce, bytes signatures) returns()
func (_ERC20Bridge *ERC20BridgeSession) ListAsset(asset_source common.Address, vega_asset_id [32]byte, lifetime_limit *big.Int, withdraw_threshold *big.Int, nonce *big.Int, signatures []byte) (*types.Transaction, error) {
	return _ERC20Bridge.Contract.ListAsset(&_ERC20Bridge.TransactOpts, asset_source, vega_asset_id, lifetime_limit, withdraw_threshold, nonce, signatures)
}

// ListAsset is a paid mutator transaction binding the contract method 0x0ff3562c.
//
// Solidity: function list_asset(address asset_source, bytes32 vega_asset_id, uint256 lifetime_limit, uint256 withdraw_threshold, uint256 nonce, bytes signatures) returns()
func (_ERC20Bridge *ERC20BridgeTransactorSession) ListAsset(asset_source common.Address, vega_asset_id [32]byte, lifetime_limit *big.Int, withdraw_threshold *big.Int, nonce *big.Int, signatures []byte) (*types.Transaction, error) {
	return _ERC20Bridge.Contract.ListAsset(&_ERC20Bridge.TransactOpts, asset_source, vega_asset_id, lifetime_limit, withdraw_threshold, nonce, signatures)
}

// RemoveAsset is a paid mutator transaction binding the contract method 0xc76de358.
//
// Solidity: function remove_asset(address asset_source, uint256 nonce, bytes signatures) returns()
func (_ERC20Bridge *ERC20BridgeTransactor) RemoveAsset(opts *bind.TransactOpts, asset_source common.Address, nonce *big.Int, signatures []byte) (*types.Transaction, error) {
	return _ERC20Bridge.contract.Transact(opts, "remove_asset", asset_source, nonce, signatures)
}

// RemoveAsset is a paid mutator transaction binding the contract method 0xc76de358.
//
// Solidity: function remove_asset(address asset_source, uint256 nonce, bytes signatures) returns()
func (_ERC20Bridge *ERC20BridgeSession) RemoveAsset(asset_source common.Address, nonce *big.Int, signatures []byte) (*types.Transaction, error) {
	return _ERC20Bridge.Contract.RemoveAsset(&_ERC20Bridge.TransactOpts, asset_source, nonce, signatures)
}

// RemoveAsset is a paid mutator transaction binding the contract method 0xc76de358.
//
// Solidity: function remove_asset(address asset_source, uint256 nonce, bytes signatures) returns()
func (_ERC20Bridge *ERC20BridgeTransactorSession) RemoveAsset(asset_source common.Address, nonce *big.Int, signatures []byte) (*types.Transaction, error) {
	return _ERC20Bridge.Contract.RemoveAsset(&_ERC20Bridge.TransactOpts, asset_source, nonce, signatures)
}

// RevokeExemptDepositor is a paid mutator transaction binding the contract method 0x6a1c6fa4.
//
// Solidity: function revoke_exempt_depositor() returns()
func (_ERC20Bridge *ERC20BridgeTransactor) RevokeExemptDepositor(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ERC20Bridge.contract.Transact(opts, "revoke_exempt_depositor")
}

// RevokeExemptDepositor is a paid mutator transaction binding the contract method 0x6a1c6fa4.
//
// Solidity: function revoke_exempt_depositor() returns()
func (_ERC20Bridge *ERC20BridgeSession) RevokeExemptDepositor() (*types.Transaction, error) {
	return _ERC20Bridge.Contract.RevokeExemptDepositor(&_ERC20Bridge.TransactOpts)
}

// RevokeExemptDepositor is a paid mutator transaction binding the contract method 0x6a1c6fa4.
//
// Solidity: function revoke_exempt_depositor() returns()
func (_ERC20Bridge *ERC20BridgeTransactorSession) RevokeExemptDepositor() (*types.Transaction, error) {
	return _ERC20Bridge.Contract.RevokeExemptDepositor(&_ERC20Bridge.TransactOpts)
}

// SetAssetLimits is a paid mutator transaction binding the contract method 0x41fb776d.
//
// Solidity: function set_asset_limits(address asset_source, uint256 lifetime_limit, uint256 threshold, uint256 nonce, bytes signatures) returns()
func (_ERC20Bridge *ERC20BridgeTransactor) SetAssetLimits(opts *bind.TransactOpts, asset_source common.Address, lifetime_limit *big.Int, threshold *big.Int, nonce *big.Int, signatures []byte) (*types.Transaction, error) {
	return _ERC20Bridge.contract.Transact(opts, "set_asset_limits", asset_source, lifetime_limit, threshold, nonce, signatures)
}

// SetAssetLimits is a paid mutator transaction binding the contract method 0x41fb776d.
//
// Solidity: function set_asset_limits(address asset_source, uint256 lifetime_limit, uint256 threshold, uint256 nonce, bytes signatures) returns()
func (_ERC20Bridge *ERC20BridgeSession) SetAssetLimits(asset_source common.Address, lifetime_limit *big.Int, threshold *big.Int, nonce *big.Int, signatures []byte) (*types.Transaction, error) {
	return _ERC20Bridge.Contract.SetAssetLimits(&_ERC20Bridge.TransactOpts, asset_source, lifetime_limit, threshold, nonce, signatures)
}

// SetAssetLimits is a paid mutator transaction binding the contract method 0x41fb776d.
//
// Solidity: function set_asset_limits(address asset_source, uint256 lifetime_limit, uint256 threshold, uint256 nonce, bytes signatures) returns()
func (_ERC20Bridge *ERC20BridgeTransactorSession) SetAssetLimits(asset_source common.Address, lifetime_limit *big.Int, threshold *big.Int, nonce *big.Int, signatures []byte) (*types.Transaction, error) {
	return _ERC20Bridge.Contract.SetAssetLimits(&_ERC20Bridge.TransactOpts, asset_source, lifetime_limit, threshold, nonce, signatures)
}

// SetWithdrawDelay is a paid mutator transaction binding the contract method 0x5a246728.
//
// Solidity: function set_withdraw_delay(uint256 delay, uint256 nonce, bytes signatures) returns()
func (_ERC20Bridge *ERC20BridgeTransactor) SetWithdrawDelay(opts *bind.TransactOpts, delay *big.Int, nonce *big.Int, signatures []byte) (*types.Transaction, error) {
	return _ERC20Bridge.contract.Transact(opts, "set_withdraw_delay", delay, nonce, signatures)
}

// SetWithdrawDelay is a paid mutator transaction binding the contract method 0x5a246728.
//
// Solidity: function set_withdraw_delay(uint256 delay, uint256 nonce, bytes signatures) returns()
func (_ERC20Bridge *ERC20BridgeSession) SetWithdrawDelay(delay *big.Int, nonce *big.Int, signatures []byte) (*types.Transaction, error) {
	return _ERC20Bridge.Contract.SetWithdrawDelay(&_ERC20Bridge.TransactOpts, delay, nonce, signatures)
}

// SetWithdrawDelay is a paid mutator transaction binding the contract method 0x5a246728.
//
// Solidity: function set_withdraw_delay(uint256 delay, uint256 nonce, bytes signatures) returns()
func (_ERC20Bridge *ERC20BridgeTransactorSession) SetWithdrawDelay(delay *big.Int, nonce *big.Int, signatures []byte) (*types.Transaction, error) {
	return _ERC20Bridge.Contract.SetWithdrawDelay(&_ERC20Bridge.TransactOpts, delay, nonce, signatures)
}

// WithdrawAsset is a paid mutator transaction binding the contract method 0x3ad90635.
//
// Solidity: function withdraw_asset(address asset_source, uint256 amount, address target, uint256 creation, uint256 nonce, bytes signatures) returns()
func (_ERC20Bridge *ERC20BridgeTransactor) WithdrawAsset(opts *bind.TransactOpts, asset_source common.Address, amount *big.Int, target common.Address, creation *big.Int, nonce *big.Int, signatures []byte) (*types.Transaction, error) {
	return _ERC20Bridge.contract.Transact(opts, "withdraw_asset", asset_source, amount, target, creation, nonce, signatures)
}

// WithdrawAsset is a paid mutator transaction binding the contract method 0x3ad90635.
//
// Solidity: function withdraw_asset(address asset_source, uint256 amount, address target, uint256 creation, uint256 nonce, bytes signatures) returns()
func (_ERC20Bridge *ERC20BridgeSession) WithdrawAsset(asset_source common.Address, amount *big.Int, target common.Address, creation *big.Int, nonce *big.Int, signatures []byte) (*types.Transaction, error) {
	return _ERC20Bridge.Contract.WithdrawAsset(&_ERC20Bridge.TransactOpts, asset_source, amount, target, creation, nonce, signatures)
}

// WithdrawAsset is a paid mutator transaction binding the contract method 0x3ad90635.
//
// Solidity: function withdraw_asset(address asset_source, uint256 amount, address target, uint256 creation, uint256 nonce, bytes signatures) returns()
func (_ERC20Bridge *ERC20BridgeTransactorSession) WithdrawAsset(asset_source common.Address, amount *big.Int, target common.Address, creation *big.Int, nonce *big.Int, signatures []byte) (*types.Transaction, error) {
	return _ERC20Bridge.Contract.WithdrawAsset(&_ERC20Bridge.TransactOpts, asset_source, amount, target, creation, nonce, signatures)
}

// ERC20BridgeAssetDepositedIterator is returned from FilterAssetDeposited and is used to iterate over the raw logs and unpacked data for AssetDeposited events raised by the ERC20Bridge contract.
type ERC20BridgeAssetDepositedIterator struct {
	Event *ERC20BridgeAssetDeposited // Event containing the contract specifics and raw log

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
func (it *ERC20BridgeAssetDepositedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ERC20BridgeAssetDeposited)
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
		it.Event = new(ERC20BridgeAssetDeposited)
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
func (it *ERC20BridgeAssetDepositedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ERC20BridgeAssetDepositedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ERC20BridgeAssetDeposited represents a AssetDeposited event raised by the ERC20Bridge contract.
type ERC20BridgeAssetDeposited struct {
	UserAddress   common.Address
	AssetSource   common.Address
	Amount        *big.Int
	VegaPublicKey [32]byte
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterAssetDeposited is a free log retrieval operation binding the contract event 0x3724ff5e82ddc640a08d68b0b782a5991aea0de51a8dd10a59cdbe5b3ec4e6bf.
//
// Solidity: event Asset_Deposited(address indexed user_address, address indexed asset_source, uint256 amount, bytes32 vega_public_key)
func (_ERC20Bridge *ERC20BridgeFilterer) FilterAssetDeposited(opts *bind.FilterOpts, user_address []common.Address, asset_source []common.Address) (*ERC20BridgeAssetDepositedIterator, error) {

	var user_addressRule []interface{}
	for _, user_addressItem := range user_address {
		user_addressRule = append(user_addressRule, user_addressItem)
	}
	var asset_sourceRule []interface{}
	for _, asset_sourceItem := range asset_source {
		asset_sourceRule = append(asset_sourceRule, asset_sourceItem)
	}

	logs, sub, err := _ERC20Bridge.contract.FilterLogs(opts, "Asset_Deposited", user_addressRule, asset_sourceRule)
	if err != nil {
		return nil, err
	}
	return &ERC20BridgeAssetDepositedIterator{contract: _ERC20Bridge.contract, event: "Asset_Deposited", logs: logs, sub: sub}, nil
}

// WatchAssetDeposited is a free log subscription operation binding the contract event 0x3724ff5e82ddc640a08d68b0b782a5991aea0de51a8dd10a59cdbe5b3ec4e6bf.
//
// Solidity: event Asset_Deposited(address indexed user_address, address indexed asset_source, uint256 amount, bytes32 vega_public_key)
func (_ERC20Bridge *ERC20BridgeFilterer) WatchAssetDeposited(opts *bind.WatchOpts, sink chan<- *ERC20BridgeAssetDeposited, user_address []common.Address, asset_source []common.Address) (event.Subscription, error) {

	var user_addressRule []interface{}
	for _, user_addressItem := range user_address {
		user_addressRule = append(user_addressRule, user_addressItem)
	}
	var asset_sourceRule []interface{}
	for _, asset_sourceItem := range asset_source {
		asset_sourceRule = append(asset_sourceRule, asset_sourceItem)
	}

	logs, sub, err := _ERC20Bridge.contract.WatchLogs(opts, "Asset_Deposited", user_addressRule, asset_sourceRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ERC20BridgeAssetDeposited)
				if err := _ERC20Bridge.contract.UnpackLog(event, "Asset_Deposited", log); err != nil {
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

// ParseAssetDeposited is a log parse operation binding the contract event 0x3724ff5e82ddc640a08d68b0b782a5991aea0de51a8dd10a59cdbe5b3ec4e6bf.
//
// Solidity: event Asset_Deposited(address indexed user_address, address indexed asset_source, uint256 amount, bytes32 vega_public_key)
func (_ERC20Bridge *ERC20BridgeFilterer) ParseAssetDeposited(log types.Log) (*ERC20BridgeAssetDeposited, error) {
	event := new(ERC20BridgeAssetDeposited)
	if err := _ERC20Bridge.contract.UnpackLog(event, "Asset_Deposited", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ERC20BridgeAssetLimitsUpdatedIterator is returned from FilterAssetLimitsUpdated and is used to iterate over the raw logs and unpacked data for AssetLimitsUpdated events raised by the ERC20Bridge contract.
type ERC20BridgeAssetLimitsUpdatedIterator struct {
	Event *ERC20BridgeAssetLimitsUpdated // Event containing the contract specifics and raw log

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
func (it *ERC20BridgeAssetLimitsUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ERC20BridgeAssetLimitsUpdated)
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
		it.Event = new(ERC20BridgeAssetLimitsUpdated)
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
func (it *ERC20BridgeAssetLimitsUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ERC20BridgeAssetLimitsUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ERC20BridgeAssetLimitsUpdated represents a AssetLimitsUpdated event raised by the ERC20Bridge contract.
type ERC20BridgeAssetLimitsUpdated struct {
	AssetSource       common.Address
	LifetimeLimit     *big.Int
	WithdrawThreshold *big.Int
	Raw               types.Log // Blockchain specific contextual infos
}

// FilterAssetLimitsUpdated is a free log retrieval operation binding the contract event 0xfc7eab762b8751ad85c101fd1025c763b4e8d48f2093f506629b606618e884fe.
//
// Solidity: event Asset_Limits_Updated(address indexed asset_source, uint256 lifetime_limit, uint256 withdraw_threshold)
func (_ERC20Bridge *ERC20BridgeFilterer) FilterAssetLimitsUpdated(opts *bind.FilterOpts, asset_source []common.Address) (*ERC20BridgeAssetLimitsUpdatedIterator, error) {

	var asset_sourceRule []interface{}
	for _, asset_sourceItem := range asset_source {
		asset_sourceRule = append(asset_sourceRule, asset_sourceItem)
	}

	logs, sub, err := _ERC20Bridge.contract.FilterLogs(opts, "Asset_Limits_Updated", asset_sourceRule)
	if err != nil {
		return nil, err
	}
	return &ERC20BridgeAssetLimitsUpdatedIterator{contract: _ERC20Bridge.contract, event: "Asset_Limits_Updated", logs: logs, sub: sub}, nil
}

// WatchAssetLimitsUpdated is a free log subscription operation binding the contract event 0xfc7eab762b8751ad85c101fd1025c763b4e8d48f2093f506629b606618e884fe.
//
// Solidity: event Asset_Limits_Updated(address indexed asset_source, uint256 lifetime_limit, uint256 withdraw_threshold)
func (_ERC20Bridge *ERC20BridgeFilterer) WatchAssetLimitsUpdated(opts *bind.WatchOpts, sink chan<- *ERC20BridgeAssetLimitsUpdated, asset_source []common.Address) (event.Subscription, error) {

	var asset_sourceRule []interface{}
	for _, asset_sourceItem := range asset_source {
		asset_sourceRule = append(asset_sourceRule, asset_sourceItem)
	}

	logs, sub, err := _ERC20Bridge.contract.WatchLogs(opts, "Asset_Limits_Updated", asset_sourceRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ERC20BridgeAssetLimitsUpdated)
				if err := _ERC20Bridge.contract.UnpackLog(event, "Asset_Limits_Updated", log); err != nil {
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

// ParseAssetLimitsUpdated is a log parse operation binding the contract event 0xfc7eab762b8751ad85c101fd1025c763b4e8d48f2093f506629b606618e884fe.
//
// Solidity: event Asset_Limits_Updated(address indexed asset_source, uint256 lifetime_limit, uint256 withdraw_threshold)
func (_ERC20Bridge *ERC20BridgeFilterer) ParseAssetLimitsUpdated(log types.Log) (*ERC20BridgeAssetLimitsUpdated, error) {
	event := new(ERC20BridgeAssetLimitsUpdated)
	if err := _ERC20Bridge.contract.UnpackLog(event, "Asset_Limits_Updated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ERC20BridgeAssetListedIterator is returned from FilterAssetListed and is used to iterate over the raw logs and unpacked data for AssetListed events raised by the ERC20Bridge contract.
type ERC20BridgeAssetListedIterator struct {
	Event *ERC20BridgeAssetListed // Event containing the contract specifics and raw log

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
func (it *ERC20BridgeAssetListedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ERC20BridgeAssetListed)
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
		it.Event = new(ERC20BridgeAssetListed)
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
func (it *ERC20BridgeAssetListedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ERC20BridgeAssetListedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ERC20BridgeAssetListed represents a AssetListed event raised by the ERC20Bridge contract.
type ERC20BridgeAssetListed struct {
	AssetSource common.Address
	VegaAssetId [32]byte
	Nonce       *big.Int
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterAssetListed is a free log retrieval operation binding the contract event 0x4180d77d05ff0d31650c548c23f2de07a3da3ad42e3dd6edd817b438a150452e.
//
// Solidity: event Asset_Listed(address indexed asset_source, bytes32 indexed vega_asset_id, uint256 nonce)
func (_ERC20Bridge *ERC20BridgeFilterer) FilterAssetListed(opts *bind.FilterOpts, asset_source []common.Address, vega_asset_id [][32]byte) (*ERC20BridgeAssetListedIterator, error) {

	var asset_sourceRule []interface{}
	for _, asset_sourceItem := range asset_source {
		asset_sourceRule = append(asset_sourceRule, asset_sourceItem)
	}
	var vega_asset_idRule []interface{}
	for _, vega_asset_idItem := range vega_asset_id {
		vega_asset_idRule = append(vega_asset_idRule, vega_asset_idItem)
	}

	logs, sub, err := _ERC20Bridge.contract.FilterLogs(opts, "Asset_Listed", asset_sourceRule, vega_asset_idRule)
	if err != nil {
		return nil, err
	}
	return &ERC20BridgeAssetListedIterator{contract: _ERC20Bridge.contract, event: "Asset_Listed", logs: logs, sub: sub}, nil
}

// WatchAssetListed is a free log subscription operation binding the contract event 0x4180d77d05ff0d31650c548c23f2de07a3da3ad42e3dd6edd817b438a150452e.
//
// Solidity: event Asset_Listed(address indexed asset_source, bytes32 indexed vega_asset_id, uint256 nonce)
func (_ERC20Bridge *ERC20BridgeFilterer) WatchAssetListed(opts *bind.WatchOpts, sink chan<- *ERC20BridgeAssetListed, asset_source []common.Address, vega_asset_id [][32]byte) (event.Subscription, error) {

	var asset_sourceRule []interface{}
	for _, asset_sourceItem := range asset_source {
		asset_sourceRule = append(asset_sourceRule, asset_sourceItem)
	}
	var vega_asset_idRule []interface{}
	for _, vega_asset_idItem := range vega_asset_id {
		vega_asset_idRule = append(vega_asset_idRule, vega_asset_idItem)
	}

	logs, sub, err := _ERC20Bridge.contract.WatchLogs(opts, "Asset_Listed", asset_sourceRule, vega_asset_idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ERC20BridgeAssetListed)
				if err := _ERC20Bridge.contract.UnpackLog(event, "Asset_Listed", log); err != nil {
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

// ParseAssetListed is a log parse operation binding the contract event 0x4180d77d05ff0d31650c548c23f2de07a3da3ad42e3dd6edd817b438a150452e.
//
// Solidity: event Asset_Listed(address indexed asset_source, bytes32 indexed vega_asset_id, uint256 nonce)
func (_ERC20Bridge *ERC20BridgeFilterer) ParseAssetListed(log types.Log) (*ERC20BridgeAssetListed, error) {
	event := new(ERC20BridgeAssetListed)
	if err := _ERC20Bridge.contract.UnpackLog(event, "Asset_Listed", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ERC20BridgeAssetRemovedIterator is returned from FilterAssetRemoved and is used to iterate over the raw logs and unpacked data for AssetRemoved events raised by the ERC20Bridge contract.
type ERC20BridgeAssetRemovedIterator struct {
	Event *ERC20BridgeAssetRemoved // Event containing the contract specifics and raw log

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
func (it *ERC20BridgeAssetRemovedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ERC20BridgeAssetRemoved)
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
		it.Event = new(ERC20BridgeAssetRemoved)
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
func (it *ERC20BridgeAssetRemovedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ERC20BridgeAssetRemovedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ERC20BridgeAssetRemoved represents a AssetRemoved event raised by the ERC20Bridge contract.
type ERC20BridgeAssetRemoved struct {
	AssetSource common.Address
	Nonce       *big.Int
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterAssetRemoved is a free log retrieval operation binding the contract event 0x58ad5e799e2df93ab408be0e5c1870d44c80b5bca99dfaf7ddf0dab5e6b155c9.
//
// Solidity: event Asset_Removed(address indexed asset_source, uint256 nonce)
func (_ERC20Bridge *ERC20BridgeFilterer) FilterAssetRemoved(opts *bind.FilterOpts, asset_source []common.Address) (*ERC20BridgeAssetRemovedIterator, error) {

	var asset_sourceRule []interface{}
	for _, asset_sourceItem := range asset_source {
		asset_sourceRule = append(asset_sourceRule, asset_sourceItem)
	}

	logs, sub, err := _ERC20Bridge.contract.FilterLogs(opts, "Asset_Removed", asset_sourceRule)
	if err != nil {
		return nil, err
	}
	return &ERC20BridgeAssetRemovedIterator{contract: _ERC20Bridge.contract, event: "Asset_Removed", logs: logs, sub: sub}, nil
}

// WatchAssetRemoved is a free log subscription operation binding the contract event 0x58ad5e799e2df93ab408be0e5c1870d44c80b5bca99dfaf7ddf0dab5e6b155c9.
//
// Solidity: event Asset_Removed(address indexed asset_source, uint256 nonce)
func (_ERC20Bridge *ERC20BridgeFilterer) WatchAssetRemoved(opts *bind.WatchOpts, sink chan<- *ERC20BridgeAssetRemoved, asset_source []common.Address) (event.Subscription, error) {

	var asset_sourceRule []interface{}
	for _, asset_sourceItem := range asset_source {
		asset_sourceRule = append(asset_sourceRule, asset_sourceItem)
	}

	logs, sub, err := _ERC20Bridge.contract.WatchLogs(opts, "Asset_Removed", asset_sourceRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ERC20BridgeAssetRemoved)
				if err := _ERC20Bridge.contract.UnpackLog(event, "Asset_Removed", log); err != nil {
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

// ParseAssetRemoved is a log parse operation binding the contract event 0x58ad5e799e2df93ab408be0e5c1870d44c80b5bca99dfaf7ddf0dab5e6b155c9.
//
// Solidity: event Asset_Removed(address indexed asset_source, uint256 nonce)
func (_ERC20Bridge *ERC20BridgeFilterer) ParseAssetRemoved(log types.Log) (*ERC20BridgeAssetRemoved, error) {
	event := new(ERC20BridgeAssetRemoved)
	if err := _ERC20Bridge.contract.UnpackLog(event, "Asset_Removed", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ERC20BridgeAssetWithdrawnIterator is returned from FilterAssetWithdrawn and is used to iterate over the raw logs and unpacked data for AssetWithdrawn events raised by the ERC20Bridge contract.
type ERC20BridgeAssetWithdrawnIterator struct {
	Event *ERC20BridgeAssetWithdrawn // Event containing the contract specifics and raw log

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
func (it *ERC20BridgeAssetWithdrawnIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ERC20BridgeAssetWithdrawn)
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
		it.Event = new(ERC20BridgeAssetWithdrawn)
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
func (it *ERC20BridgeAssetWithdrawnIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ERC20BridgeAssetWithdrawnIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ERC20BridgeAssetWithdrawn represents a AssetWithdrawn event raised by the ERC20Bridge contract.
type ERC20BridgeAssetWithdrawn struct {
	UserAddress common.Address
	AssetSource common.Address
	Amount      *big.Int
	Nonce       *big.Int
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterAssetWithdrawn is a free log retrieval operation binding the contract event 0xa79be4f3361e32d396d64c478ecef73732cb40b2a75702c3b3b3226a2c83b5df.
//
// Solidity: event Asset_Withdrawn(address indexed user_address, address indexed asset_source, uint256 amount, uint256 nonce)
func (_ERC20Bridge *ERC20BridgeFilterer) FilterAssetWithdrawn(opts *bind.FilterOpts, user_address []common.Address, asset_source []common.Address) (*ERC20BridgeAssetWithdrawnIterator, error) {

	var user_addressRule []interface{}
	for _, user_addressItem := range user_address {
		user_addressRule = append(user_addressRule, user_addressItem)
	}
	var asset_sourceRule []interface{}
	for _, asset_sourceItem := range asset_source {
		asset_sourceRule = append(asset_sourceRule, asset_sourceItem)
	}

	logs, sub, err := _ERC20Bridge.contract.FilterLogs(opts, "Asset_Withdrawn", user_addressRule, asset_sourceRule)
	if err != nil {
		return nil, err
	}
	return &ERC20BridgeAssetWithdrawnIterator{contract: _ERC20Bridge.contract, event: "Asset_Withdrawn", logs: logs, sub: sub}, nil
}

// WatchAssetWithdrawn is a free log subscription operation binding the contract event 0xa79be4f3361e32d396d64c478ecef73732cb40b2a75702c3b3b3226a2c83b5df.
//
// Solidity: event Asset_Withdrawn(address indexed user_address, address indexed asset_source, uint256 amount, uint256 nonce)
func (_ERC20Bridge *ERC20BridgeFilterer) WatchAssetWithdrawn(opts *bind.WatchOpts, sink chan<- *ERC20BridgeAssetWithdrawn, user_address []common.Address, asset_source []common.Address) (event.Subscription, error) {

	var user_addressRule []interface{}
	for _, user_addressItem := range user_address {
		user_addressRule = append(user_addressRule, user_addressItem)
	}
	var asset_sourceRule []interface{}
	for _, asset_sourceItem := range asset_source {
		asset_sourceRule = append(asset_sourceRule, asset_sourceItem)
	}

	logs, sub, err := _ERC20Bridge.contract.WatchLogs(opts, "Asset_Withdrawn", user_addressRule, asset_sourceRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ERC20BridgeAssetWithdrawn)
				if err := _ERC20Bridge.contract.UnpackLog(event, "Asset_Withdrawn", log); err != nil {
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

// ParseAssetWithdrawn is a log parse operation binding the contract event 0xa79be4f3361e32d396d64c478ecef73732cb40b2a75702c3b3b3226a2c83b5df.
//
// Solidity: event Asset_Withdrawn(address indexed user_address, address indexed asset_source, uint256 amount, uint256 nonce)
func (_ERC20Bridge *ERC20BridgeFilterer) ParseAssetWithdrawn(log types.Log) (*ERC20BridgeAssetWithdrawn, error) {
	event := new(ERC20BridgeAssetWithdrawn)
	if err := _ERC20Bridge.contract.UnpackLog(event, "Asset_Withdrawn", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ERC20BridgeBridgeResumedIterator is returned from FilterBridgeResumed and is used to iterate over the raw logs and unpacked data for BridgeResumed events raised by the ERC20Bridge contract.
type ERC20BridgeBridgeResumedIterator struct {
	Event *ERC20BridgeBridgeResumed // Event containing the contract specifics and raw log

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
func (it *ERC20BridgeBridgeResumedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ERC20BridgeBridgeResumed)
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
		it.Event = new(ERC20BridgeBridgeResumed)
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
func (it *ERC20BridgeBridgeResumedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ERC20BridgeBridgeResumedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ERC20BridgeBridgeResumed represents a BridgeResumed event raised by the ERC20Bridge contract.
type ERC20BridgeBridgeResumed struct {
	Raw types.Log // Blockchain specific contextual infos
}

// FilterBridgeResumed is a free log retrieval operation binding the contract event 0x79c02b0e60e0f00fe0370791204f2f175fe3f06f4816f3506ad4fa1b8e8cde0f.
//
// Solidity: event Bridge_Resumed()
func (_ERC20Bridge *ERC20BridgeFilterer) FilterBridgeResumed(opts *bind.FilterOpts) (*ERC20BridgeBridgeResumedIterator, error) {

	logs, sub, err := _ERC20Bridge.contract.FilterLogs(opts, "Bridge_Resumed")
	if err != nil {
		return nil, err
	}
	return &ERC20BridgeBridgeResumedIterator{contract: _ERC20Bridge.contract, event: "Bridge_Resumed", logs: logs, sub: sub}, nil
}

// WatchBridgeResumed is a free log subscription operation binding the contract event 0x79c02b0e60e0f00fe0370791204f2f175fe3f06f4816f3506ad4fa1b8e8cde0f.
//
// Solidity: event Bridge_Resumed()
func (_ERC20Bridge *ERC20BridgeFilterer) WatchBridgeResumed(opts *bind.WatchOpts, sink chan<- *ERC20BridgeBridgeResumed) (event.Subscription, error) {

	logs, sub, err := _ERC20Bridge.contract.WatchLogs(opts, "Bridge_Resumed")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ERC20BridgeBridgeResumed)
				if err := _ERC20Bridge.contract.UnpackLog(event, "Bridge_Resumed", log); err != nil {
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

// ParseBridgeResumed is a log parse operation binding the contract event 0x79c02b0e60e0f00fe0370791204f2f175fe3f06f4816f3506ad4fa1b8e8cde0f.
//
// Solidity: event Bridge_Resumed()
func (_ERC20Bridge *ERC20BridgeFilterer) ParseBridgeResumed(log types.Log) (*ERC20BridgeBridgeResumed, error) {
	event := new(ERC20BridgeBridgeResumed)
	if err := _ERC20Bridge.contract.UnpackLog(event, "Bridge_Resumed", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ERC20BridgeBridgeStoppedIterator is returned from FilterBridgeStopped and is used to iterate over the raw logs and unpacked data for BridgeStopped events raised by the ERC20Bridge contract.
type ERC20BridgeBridgeStoppedIterator struct {
	Event *ERC20BridgeBridgeStopped // Event containing the contract specifics and raw log

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
func (it *ERC20BridgeBridgeStoppedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ERC20BridgeBridgeStopped)
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
		it.Event = new(ERC20BridgeBridgeStopped)
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
func (it *ERC20BridgeBridgeStoppedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ERC20BridgeBridgeStoppedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ERC20BridgeBridgeStopped represents a BridgeStopped event raised by the ERC20Bridge contract.
type ERC20BridgeBridgeStopped struct {
	Raw types.Log // Blockchain specific contextual infos
}

// FilterBridgeStopped is a free log retrieval operation binding the contract event 0x129d99581c8e70519df1f0733d3212f33d0ed3ea6144adacc336c647f1d36382.
//
// Solidity: event Bridge_Stopped()
func (_ERC20Bridge *ERC20BridgeFilterer) FilterBridgeStopped(opts *bind.FilterOpts) (*ERC20BridgeBridgeStoppedIterator, error) {

	logs, sub, err := _ERC20Bridge.contract.FilterLogs(opts, "Bridge_Stopped")
	if err != nil {
		return nil, err
	}
	return &ERC20BridgeBridgeStoppedIterator{contract: _ERC20Bridge.contract, event: "Bridge_Stopped", logs: logs, sub: sub}, nil
}

// WatchBridgeStopped is a free log subscription operation binding the contract event 0x129d99581c8e70519df1f0733d3212f33d0ed3ea6144adacc336c647f1d36382.
//
// Solidity: event Bridge_Stopped()
func (_ERC20Bridge *ERC20BridgeFilterer) WatchBridgeStopped(opts *bind.WatchOpts, sink chan<- *ERC20BridgeBridgeStopped) (event.Subscription, error) {

	logs, sub, err := _ERC20Bridge.contract.WatchLogs(opts, "Bridge_Stopped")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ERC20BridgeBridgeStopped)
				if err := _ERC20Bridge.contract.UnpackLog(event, "Bridge_Stopped", log); err != nil {
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

// ParseBridgeStopped is a log parse operation binding the contract event 0x129d99581c8e70519df1f0733d3212f33d0ed3ea6144adacc336c647f1d36382.
//
// Solidity: event Bridge_Stopped()
func (_ERC20Bridge *ERC20BridgeFilterer) ParseBridgeStopped(log types.Log) (*ERC20BridgeBridgeStopped, error) {
	event := new(ERC20BridgeBridgeStopped)
	if err := _ERC20Bridge.contract.UnpackLog(event, "Bridge_Stopped", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ERC20BridgeBridgeWithdrawDelaySetIterator is returned from FilterBridgeWithdrawDelaySet and is used to iterate over the raw logs and unpacked data for BridgeWithdrawDelaySet events raised by the ERC20Bridge contract.
type ERC20BridgeBridgeWithdrawDelaySetIterator struct {
	Event *ERC20BridgeBridgeWithdrawDelaySet // Event containing the contract specifics and raw log

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
func (it *ERC20BridgeBridgeWithdrawDelaySetIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ERC20BridgeBridgeWithdrawDelaySet)
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
		it.Event = new(ERC20BridgeBridgeWithdrawDelaySet)
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
func (it *ERC20BridgeBridgeWithdrawDelaySetIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ERC20BridgeBridgeWithdrawDelaySetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ERC20BridgeBridgeWithdrawDelaySet represents a BridgeWithdrawDelaySet event raised by the ERC20Bridge contract.
type ERC20BridgeBridgeWithdrawDelaySet struct {
	WithdrawDelay *big.Int
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterBridgeWithdrawDelaySet is a free log retrieval operation binding the contract event 0x1c7e8f73a01b8af4e18dd34455a42a45ad742bdb79cfda77bbdf50db2391fc88.
//
// Solidity: event Bridge_Withdraw_Delay_Set(uint256 withdraw_delay)
func (_ERC20Bridge *ERC20BridgeFilterer) FilterBridgeWithdrawDelaySet(opts *bind.FilterOpts) (*ERC20BridgeBridgeWithdrawDelaySetIterator, error) {

	logs, sub, err := _ERC20Bridge.contract.FilterLogs(opts, "Bridge_Withdraw_Delay_Set")
	if err != nil {
		return nil, err
	}
	return &ERC20BridgeBridgeWithdrawDelaySetIterator{contract: _ERC20Bridge.contract, event: "Bridge_Withdraw_Delay_Set", logs: logs, sub: sub}, nil
}

// WatchBridgeWithdrawDelaySet is a free log subscription operation binding the contract event 0x1c7e8f73a01b8af4e18dd34455a42a45ad742bdb79cfda77bbdf50db2391fc88.
//
// Solidity: event Bridge_Withdraw_Delay_Set(uint256 withdraw_delay)
func (_ERC20Bridge *ERC20BridgeFilterer) WatchBridgeWithdrawDelaySet(opts *bind.WatchOpts, sink chan<- *ERC20BridgeBridgeWithdrawDelaySet) (event.Subscription, error) {

	logs, sub, err := _ERC20Bridge.contract.WatchLogs(opts, "Bridge_Withdraw_Delay_Set")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ERC20BridgeBridgeWithdrawDelaySet)
				if err := _ERC20Bridge.contract.UnpackLog(event, "Bridge_Withdraw_Delay_Set", log); err != nil {
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

// ParseBridgeWithdrawDelaySet is a log parse operation binding the contract event 0x1c7e8f73a01b8af4e18dd34455a42a45ad742bdb79cfda77bbdf50db2391fc88.
//
// Solidity: event Bridge_Withdraw_Delay_Set(uint256 withdraw_delay)
func (_ERC20Bridge *ERC20BridgeFilterer) ParseBridgeWithdrawDelaySet(log types.Log) (*ERC20BridgeBridgeWithdrawDelaySet, error) {
	event := new(ERC20BridgeBridgeWithdrawDelaySet)
	if err := _ERC20Bridge.contract.UnpackLog(event, "Bridge_Withdraw_Delay_Set", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ERC20BridgeDepositorExemptedIterator is returned from FilterDepositorExempted and is used to iterate over the raw logs and unpacked data for DepositorExempted events raised by the ERC20Bridge contract.
type ERC20BridgeDepositorExemptedIterator struct {
	Event *ERC20BridgeDepositorExempted // Event containing the contract specifics and raw log

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
func (it *ERC20BridgeDepositorExemptedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ERC20BridgeDepositorExempted)
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
		it.Event = new(ERC20BridgeDepositorExempted)
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
func (it *ERC20BridgeDepositorExemptedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ERC20BridgeDepositorExemptedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ERC20BridgeDepositorExempted represents a DepositorExempted event raised by the ERC20Bridge contract.
type ERC20BridgeDepositorExempted struct {
	Depositor common.Address
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterDepositorExempted is a free log retrieval operation binding the contract event 0xf56e0868b913034a60dbca9c89ee79f8b0fa18dadbc5f6665f2f9a2cf3f51cdb.
//
// Solidity: event Depositor_Exempted(address indexed depositor)
func (_ERC20Bridge *ERC20BridgeFilterer) FilterDepositorExempted(opts *bind.FilterOpts, depositor []common.Address) (*ERC20BridgeDepositorExemptedIterator, error) {

	var depositorRule []interface{}
	for _, depositorItem := range depositor {
		depositorRule = append(depositorRule, depositorItem)
	}

	logs, sub, err := _ERC20Bridge.contract.FilterLogs(opts, "Depositor_Exempted", depositorRule)
	if err != nil {
		return nil, err
	}
	return &ERC20BridgeDepositorExemptedIterator{contract: _ERC20Bridge.contract, event: "Depositor_Exempted", logs: logs, sub: sub}, nil
}

// WatchDepositorExempted is a free log subscription operation binding the contract event 0xf56e0868b913034a60dbca9c89ee79f8b0fa18dadbc5f6665f2f9a2cf3f51cdb.
//
// Solidity: event Depositor_Exempted(address indexed depositor)
func (_ERC20Bridge *ERC20BridgeFilterer) WatchDepositorExempted(opts *bind.WatchOpts, sink chan<- *ERC20BridgeDepositorExempted, depositor []common.Address) (event.Subscription, error) {

	var depositorRule []interface{}
	for _, depositorItem := range depositor {
		depositorRule = append(depositorRule, depositorItem)
	}

	logs, sub, err := _ERC20Bridge.contract.WatchLogs(opts, "Depositor_Exempted", depositorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ERC20BridgeDepositorExempted)
				if err := _ERC20Bridge.contract.UnpackLog(event, "Depositor_Exempted", log); err != nil {
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

// ParseDepositorExempted is a log parse operation binding the contract event 0xf56e0868b913034a60dbca9c89ee79f8b0fa18dadbc5f6665f2f9a2cf3f51cdb.
//
// Solidity: event Depositor_Exempted(address indexed depositor)
func (_ERC20Bridge *ERC20BridgeFilterer) ParseDepositorExempted(log types.Log) (*ERC20BridgeDepositorExempted, error) {
	event := new(ERC20BridgeDepositorExempted)
	if err := _ERC20Bridge.contract.UnpackLog(event, "Depositor_Exempted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ERC20BridgeDepositorExemptionRevokedIterator is returned from FilterDepositorExemptionRevoked and is used to iterate over the raw logs and unpacked data for DepositorExemptionRevoked events raised by the ERC20Bridge contract.
type ERC20BridgeDepositorExemptionRevokedIterator struct {
	Event *ERC20BridgeDepositorExemptionRevoked // Event containing the contract specifics and raw log

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
func (it *ERC20BridgeDepositorExemptionRevokedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ERC20BridgeDepositorExemptionRevoked)
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
		it.Event = new(ERC20BridgeDepositorExemptionRevoked)
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
func (it *ERC20BridgeDepositorExemptionRevokedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ERC20BridgeDepositorExemptionRevokedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ERC20BridgeDepositorExemptionRevoked represents a DepositorExemptionRevoked event raised by the ERC20Bridge contract.
type ERC20BridgeDepositorExemptionRevoked struct {
	Depositor common.Address
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterDepositorExemptionRevoked is a free log retrieval operation binding the contract event 0xe74b113dca87276d976f476a9b4b9da3c780a3262eaabad051ee4e98912936a4.
//
// Solidity: event Depositor_Exemption_Revoked(address indexed depositor)
func (_ERC20Bridge *ERC20BridgeFilterer) FilterDepositorExemptionRevoked(opts *bind.FilterOpts, depositor []common.Address) (*ERC20BridgeDepositorExemptionRevokedIterator, error) {

	var depositorRule []interface{}
	for _, depositorItem := range depositor {
		depositorRule = append(depositorRule, depositorItem)
	}

	logs, sub, err := _ERC20Bridge.contract.FilterLogs(opts, "Depositor_Exemption_Revoked", depositorRule)
	if err != nil {
		return nil, err
	}
	return &ERC20BridgeDepositorExemptionRevokedIterator{contract: _ERC20Bridge.contract, event: "Depositor_Exemption_Revoked", logs: logs, sub: sub}, nil
}

// WatchDepositorExemptionRevoked is a free log subscription operation binding the contract event 0xe74b113dca87276d976f476a9b4b9da3c780a3262eaabad051ee4e98912936a4.
//
// Solidity: event Depositor_Exemption_Revoked(address indexed depositor)
func (_ERC20Bridge *ERC20BridgeFilterer) WatchDepositorExemptionRevoked(opts *bind.WatchOpts, sink chan<- *ERC20BridgeDepositorExemptionRevoked, depositor []common.Address) (event.Subscription, error) {

	var depositorRule []interface{}
	for _, depositorItem := range depositor {
		depositorRule = append(depositorRule, depositorItem)
	}

	logs, sub, err := _ERC20Bridge.contract.WatchLogs(opts, "Depositor_Exemption_Revoked", depositorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ERC20BridgeDepositorExemptionRevoked)
				if err := _ERC20Bridge.contract.UnpackLog(event, "Depositor_Exemption_Revoked", log); err != nil {
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

// ParseDepositorExemptionRevoked is a log parse operation binding the contract event 0xe74b113dca87276d976f476a9b4b9da3c780a3262eaabad051ee4e98912936a4.
//
// Solidity: event Depositor_Exemption_Revoked(address indexed depositor)
func (_ERC20Bridge *ERC20BridgeFilterer) ParseDepositorExemptionRevoked(log types.Log) (*ERC20BridgeDepositorExemptionRevoked, error) {
	event := new(ERC20BridgeDepositorExemptionRevoked)
	if err := _ERC20Bridge.contract.UnpackLog(event, "Depositor_Exemption_Revoked", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
