// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package blockchain

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

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = abi.U256
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
)

// ProxeusFSABI is the input ABI used to generate the binding from.
const ProxeusFSABI = "[{\"constant\":false,\"inputs\":[{\"name\":\"_hash\",\"type\":\"bytes32\"},{\"name\":\"_data\",\"type\":\"bytes32\"}],\"name\":\"registerFile\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_hash\",\"type\":\"bytes32\"}],\"name\":\"getFileSigners\",\"outputs\":[{\"name\":\"\",\"type\":\"address[]\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_hash\",\"type\":\"bytes32\"}],\"name\":\"verifyFile\",\"outputs\":[{\"name\":\"exists\",\"type\":\"bool\"},{\"name\":\"issuer\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_hash\",\"type\":\"bytes32\"}],\"name\":\"signFile\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"files\",\"outputs\":[{\"name\":\"issuer\",\"type\":\"address\"},{\"name\":\"data\",\"type\":\"bytes32\"},{\"name\":\"exists\",\"type\":\"bool\"},{\"name\":\"signersCount\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"name\":\"_owner\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"hash\",\"type\":\"bytes32\"}],\"name\":\"UpdatedEvent\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"hash\",\"type\":\"bytes32\"},{\"indexed\":true,\"name\":\"signer\",\"type\":\"address\"}],\"name\":\"FileSignedEvent\",\"type\":\"event\"}]"

// ProxeusFS is an auto generated Go binding around an Ethereum contract.
type ProxeusFS struct {
	ProxeusFSCaller     // Read-only binding to the contract
	ProxeusFSTransactor // Write-only binding to the contract
	ProxeusFSFilterer   // Log filterer for contract events
}

// ProxeusFSCaller is an auto generated read-only Go binding around an Ethereum contract.
type ProxeusFSCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ProxeusFSTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ProxeusFSTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ProxeusFSFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ProxeusFSFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ProxeusFSSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ProxeusFSSession struct {
	Contract     *ProxeusFS        // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// ProxeusFSCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ProxeusFSCallerSession struct {
	Contract *ProxeusFSCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts    // Call options to use throughout this session
}

// ProxeusFSTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ProxeusFSTransactorSession struct {
	Contract     *ProxeusFSTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts    // Transaction auth options to use throughout this session
}

// ProxeusFSRaw is an auto generated low-level Go binding around an Ethereum contract.
type ProxeusFSRaw struct {
	Contract *ProxeusFS // Generic contract binding to access the raw methods on
}

// ProxeusFSCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ProxeusFSCallerRaw struct {
	Contract *ProxeusFSCaller // Generic read-only contract binding to access the raw methods on
}

// ProxeusFSTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ProxeusFSTransactorRaw struct {
	Contract *ProxeusFSTransactor // Generic write-only contract binding to access the raw methods on
}

// NewProxeusFS creates a new instance of ProxeusFS, bound to a specific deployed contract.
func NewProxeusFS(address common.Address, backend bind.ContractBackend) (*ProxeusFS, error) {
	contract, err := bindProxeusFS(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ProxeusFS{ProxeusFSCaller: ProxeusFSCaller{contract: contract}, ProxeusFSTransactor: ProxeusFSTransactor{contract: contract}, ProxeusFSFilterer: ProxeusFSFilterer{contract: contract}}, nil
}

// NewProxeusFSCaller creates a new read-only instance of ProxeusFS, bound to a specific deployed contract.
func NewProxeusFSCaller(address common.Address, caller bind.ContractCaller) (*ProxeusFSCaller, error) {
	contract, err := bindProxeusFS(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ProxeusFSCaller{contract: contract}, nil
}

// NewProxeusFSTransactor creates a new write-only instance of ProxeusFS, bound to a specific deployed contract.
func NewProxeusFSTransactor(address common.Address, transactor bind.ContractTransactor) (*ProxeusFSTransactor, error) {
	contract, err := bindProxeusFS(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ProxeusFSTransactor{contract: contract}, nil
}

// NewProxeusFSFilterer creates a new log filterer instance of ProxeusFS, bound to a specific deployed contract.
func NewProxeusFSFilterer(address common.Address, filterer bind.ContractFilterer) (*ProxeusFSFilterer, error) {
	contract, err := bindProxeusFS(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ProxeusFSFilterer{contract: contract}, nil
}

// bindProxeusFS binds a generic wrapper to an already deployed contract.
func bindProxeusFS(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(ProxeusFSABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ProxeusFS *ProxeusFSRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _ProxeusFS.Contract.ProxeusFSCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ProxeusFS *ProxeusFSRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ProxeusFS.Contract.ProxeusFSTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ProxeusFS *ProxeusFSRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ProxeusFS.Contract.ProxeusFSTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ProxeusFS *ProxeusFSCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _ProxeusFS.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ProxeusFS *ProxeusFSTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ProxeusFS.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ProxeusFS *ProxeusFSTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ProxeusFS.Contract.contract.Transact(opts, method, params...)
}

// Files is a free data retrieval call binding the contract method 0x98c9adff.
//
// Solidity: function files(bytes32 ) constant returns(address issuer, bytes32 data, bool exists, uint256 signersCount)
func (_ProxeusFS *ProxeusFSCaller) Files(opts *bind.CallOpts, arg0 [32]byte) (struct {
	Issuer       common.Address
	Data         [32]byte
	Exists       bool
	SignersCount *big.Int
}, error) {
	ret := new(struct {
		Issuer       common.Address
		Data         [32]byte
		Exists       bool
		SignersCount *big.Int
	})
	out := ret
	err := _ProxeusFS.contract.Call(opts, out, "files", arg0)
	return *ret, err
}

// Files is a free data retrieval call binding the contract method 0x98c9adff.
//
// Solidity: function files(bytes32 ) constant returns(address issuer, bytes32 data, bool exists, uint256 signersCount)
func (_ProxeusFS *ProxeusFSSession) Files(arg0 [32]byte) (struct {
	Issuer       common.Address
	Data         [32]byte
	Exists       bool
	SignersCount *big.Int
}, error) {
	return _ProxeusFS.Contract.Files(&_ProxeusFS.CallOpts, arg0)
}

// Files is a free data retrieval call binding the contract method 0x98c9adff.
//
// Solidity: function files(bytes32 ) constant returns(address issuer, bytes32 data, bool exists, uint256 signersCount)
func (_ProxeusFS *ProxeusFSCallerSession) Files(arg0 [32]byte) (struct {
	Issuer       common.Address
	Data         [32]byte
	Exists       bool
	SignersCount *big.Int
}, error) {
	return _ProxeusFS.Contract.Files(&_ProxeusFS.CallOpts, arg0)
}

// GetFileSigners is a free data retrieval call binding the contract method 0x24f225f5.
//
// Solidity: function getFileSigners(bytes32 _hash) constant returns(address[])
func (_ProxeusFS *ProxeusFSCaller) GetFileSigners(opts *bind.CallOpts, _hash [32]byte) ([]common.Address, error) {
	var (
		ret0 = new([]common.Address)
	)
	out := ret0
	err := _ProxeusFS.contract.Call(opts, out, "getFileSigners", _hash)
	return *ret0, err
}

// GetFileSigners is a free data retrieval call binding the contract method 0x24f225f5.
//
// Solidity: function getFileSigners(bytes32 _hash) constant returns(address[])
func (_ProxeusFS *ProxeusFSSession) GetFileSigners(_hash [32]byte) ([]common.Address, error) {
	return _ProxeusFS.Contract.GetFileSigners(&_ProxeusFS.CallOpts, _hash)
}

// GetFileSigners is a free data retrieval call binding the contract method 0x24f225f5.
//
// Solidity: function getFileSigners(bytes32 _hash) constant returns(address[])
func (_ProxeusFS *ProxeusFSCallerSession) GetFileSigners(_hash [32]byte) ([]common.Address, error) {
	return _ProxeusFS.Contract.GetFileSigners(&_ProxeusFS.CallOpts, _hash)
}

// VerifyFile is a free data retrieval call binding the contract method 0x4b67d54b.
//
// Solidity: function verifyFile(bytes32 _hash) constant returns(bool exists, address issuer)
func (_ProxeusFS *ProxeusFSCaller) VerifyFile(opts *bind.CallOpts, _hash [32]byte) (struct {
	Exists bool
	Issuer common.Address
}, error) {
	ret := new(struct {
		Exists bool
		Issuer common.Address
	})
	out := ret
	err := _ProxeusFS.contract.Call(opts, out, "verifyFile", _hash)
	return *ret, err
}

// VerifyFile is a free data retrieval call binding the contract method 0x4b67d54b.
//
// Solidity: function verifyFile(bytes32 _hash) constant returns(bool exists, address issuer)
func (_ProxeusFS *ProxeusFSSession) VerifyFile(_hash [32]byte) (struct {
	Exists bool
	Issuer common.Address
}, error) {
	return _ProxeusFS.Contract.VerifyFile(&_ProxeusFS.CallOpts, _hash)
}

// VerifyFile is a free data retrieval call binding the contract method 0x4b67d54b.
//
// Solidity: function verifyFile(bytes32 _hash) constant returns(bool exists, address issuer)
func (_ProxeusFS *ProxeusFSCallerSession) VerifyFile(_hash [32]byte) (struct {
	Exists bool
	Issuer common.Address
}, error) {
	return _ProxeusFS.Contract.VerifyFile(&_ProxeusFS.CallOpts, _hash)
}

// RegisterFile is a paid mutator transaction binding the contract method 0x1aa9030c.
//
// Solidity: function registerFile(bytes32 _hash, bytes32 _data) returns()
func (_ProxeusFS *ProxeusFSTransactor) RegisterFile(opts *bind.TransactOpts, _hash [32]byte, _data [32]byte) (*types.Transaction, error) {
	return _ProxeusFS.contract.Transact(opts, "registerFile", _hash, _data)
}

// RegisterFile is a paid mutator transaction binding the contract method 0x1aa9030c.
//
// Solidity: function registerFile(bytes32 _hash, bytes32 _data) returns()
func (_ProxeusFS *ProxeusFSSession) RegisterFile(_hash [32]byte, _data [32]byte) (*types.Transaction, error) {
	return _ProxeusFS.Contract.RegisterFile(&_ProxeusFS.TransactOpts, _hash, _data)
}

// RegisterFile is a paid mutator transaction binding the contract method 0x1aa9030c.
//
// Solidity: function registerFile(bytes32 _hash, bytes32 _data) returns()
func (_ProxeusFS *ProxeusFSTransactorSession) RegisterFile(_hash [32]byte, _data [32]byte) (*types.Transaction, error) {
	return _ProxeusFS.Contract.RegisterFile(&_ProxeusFS.TransactOpts, _hash, _data)
}

// SignFile is a paid mutator transaction binding the contract method 0x57999119.
//
// Solidity: function signFile(bytes32 _hash) returns()
func (_ProxeusFS *ProxeusFSTransactor) SignFile(opts *bind.TransactOpts, _hash [32]byte) (*types.Transaction, error) {
	return _ProxeusFS.contract.Transact(opts, "signFile", _hash)
}

// SignFile is a paid mutator transaction binding the contract method 0x57999119.
//
// Solidity: function signFile(bytes32 _hash) returns()
func (_ProxeusFS *ProxeusFSSession) SignFile(_hash [32]byte) (*types.Transaction, error) {
	return _ProxeusFS.Contract.SignFile(&_ProxeusFS.TransactOpts, _hash)
}

// SignFile is a paid mutator transaction binding the contract method 0x57999119.
//
// Solidity: function signFile(bytes32 _hash) returns()
func (_ProxeusFS *ProxeusFSTransactorSession) SignFile(_hash [32]byte) (*types.Transaction, error) {
	return _ProxeusFS.Contract.SignFile(&_ProxeusFS.TransactOpts, _hash)
}

// ProxeusFSFileSignedEventIterator is returned from FilterFileSignedEvent and is used to iterate over the raw logs and unpacked data for FileSignedEvent events raised by the ProxeusFS contract.
type ProxeusFSFileSignedEventIterator struct {
	Event *ProxeusFSFileSignedEvent // Event containing the contract specifics and raw log

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
func (it *ProxeusFSFileSignedEventIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ProxeusFSFileSignedEvent)
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
		it.Event = new(ProxeusFSFileSignedEvent)
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
func (it *ProxeusFSFileSignedEventIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ProxeusFSFileSignedEventIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ProxeusFSFileSignedEvent represents a FileSignedEvent event raised by the ProxeusFS contract.
type ProxeusFSFileSignedEvent struct {
	Hash   [32]byte
	Signer common.Address
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterFileSignedEvent is a free log retrieval operation binding the contract event 0xe898b82efcc77a621bbc620d08e84d0b44e341fd7720cc544de745bdec8ebece.
//
// Solidity: event FileSignedEvent(bytes32 indexed hash, address indexed signer)
func (_ProxeusFS *ProxeusFSFilterer) FilterFileSignedEvent(opts *bind.FilterOpts, hash [][32]byte, signer []common.Address) (*ProxeusFSFileSignedEventIterator, error) {

	var hashRule []interface{}
	for _, hashItem := range hash {
		hashRule = append(hashRule, hashItem)
	}
	var signerRule []interface{}
	for _, signerItem := range signer {
		signerRule = append(signerRule, signerItem)
	}

	logs, sub, err := _ProxeusFS.contract.FilterLogs(opts, "FileSignedEvent", hashRule, signerRule)
	if err != nil {
		return nil, err
	}
	return &ProxeusFSFileSignedEventIterator{contract: _ProxeusFS.contract, event: "FileSignedEvent", logs: logs, sub: sub}, nil
}

// WatchFileSignedEvent is a free log subscription operation binding the contract event 0xe898b82efcc77a621bbc620d08e84d0b44e341fd7720cc544de745bdec8ebece.
//
// Solidity: event FileSignedEvent(bytes32 indexed hash, address indexed signer)
func (_ProxeusFS *ProxeusFSFilterer) WatchFileSignedEvent(opts *bind.WatchOpts, sink chan<- *ProxeusFSFileSignedEvent, hash [][32]byte, signer []common.Address) (event.Subscription, error) {

	var hashRule []interface{}
	for _, hashItem := range hash {
		hashRule = append(hashRule, hashItem)
	}
	var signerRule []interface{}
	for _, signerItem := range signer {
		signerRule = append(signerRule, signerItem)
	}

	logs, sub, err := _ProxeusFS.contract.WatchLogs(opts, "FileSignedEvent", hashRule, signerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ProxeusFSFileSignedEvent)
				if err := _ProxeusFS.contract.UnpackLog(event, "FileSignedEvent", log); err != nil {
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

// ParseFileSignedEvent is a log parse operation binding the contract event 0xe898b82efcc77a621bbc620d08e84d0b44e341fd7720cc544de745bdec8ebece.
//
// Solidity: event FileSignedEvent(bytes32 indexed hash, address indexed signer)
func (_ProxeusFS *ProxeusFSFilterer) ParseFileSignedEvent(log types.Log) (*ProxeusFSFileSignedEvent, error) {
	event := new(ProxeusFSFileSignedEvent)
	if err := _ProxeusFS.contract.UnpackLog(event, "FileSignedEvent", log); err != nil {
		return nil, err
	}
	return event, nil
}

// ProxeusFSUpdatedEventIterator is returned from FilterUpdatedEvent and is used to iterate over the raw logs and unpacked data for UpdatedEvent events raised by the ProxeusFS contract.
type ProxeusFSUpdatedEventIterator struct {
	Event *ProxeusFSUpdatedEvent // Event containing the contract specifics and raw log

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
func (it *ProxeusFSUpdatedEventIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ProxeusFSUpdatedEvent)
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
		it.Event = new(ProxeusFSUpdatedEvent)
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
func (it *ProxeusFSUpdatedEventIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ProxeusFSUpdatedEventIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ProxeusFSUpdatedEvent represents a UpdatedEvent event raised by the ProxeusFS contract.
type ProxeusFSUpdatedEvent struct {
	Hash [32]byte
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterUpdatedEvent is a free log retrieval operation binding the contract event 0xc2d5f8a998001c07f1f12d8cc85e7eb8f63960ba26ef8e12b0e7d0879ac9475d.
//
// Solidity: event UpdatedEvent(bytes32 indexed hash)
func (_ProxeusFS *ProxeusFSFilterer) FilterUpdatedEvent(opts *bind.FilterOpts, hash [][32]byte) (*ProxeusFSUpdatedEventIterator, error) {

	var hashRule []interface{}
	for _, hashItem := range hash {
		hashRule = append(hashRule, hashItem)
	}

	logs, sub, err := _ProxeusFS.contract.FilterLogs(opts, "UpdatedEvent", hashRule)
	if err != nil {
		return nil, err
	}
	return &ProxeusFSUpdatedEventIterator{contract: _ProxeusFS.contract, event: "UpdatedEvent", logs: logs, sub: sub}, nil
}

// WatchUpdatedEvent is a free log subscription operation binding the contract event 0xc2d5f8a998001c07f1f12d8cc85e7eb8f63960ba26ef8e12b0e7d0879ac9475d.
//
// Solidity: event UpdatedEvent(bytes32 indexed hash)
func (_ProxeusFS *ProxeusFSFilterer) WatchUpdatedEvent(opts *bind.WatchOpts, sink chan<- *ProxeusFSUpdatedEvent, hash [][32]byte) (event.Subscription, error) {

	var hashRule []interface{}
	for _, hashItem := range hash {
		hashRule = append(hashRule, hashItem)
	}

	logs, sub, err := _ProxeusFS.contract.WatchLogs(opts, "UpdatedEvent", hashRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ProxeusFSUpdatedEvent)
				if err := _ProxeusFS.contract.UnpackLog(event, "UpdatedEvent", log); err != nil {
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

// ParseUpdatedEvent is a log parse operation binding the contract event 0xc2d5f8a998001c07f1f12d8cc85e7eb8f63960ba26ef8e12b0e7d0879ac9475d.
//
// Solidity: event UpdatedEvent(bytes32 indexed hash)
func (_ProxeusFS *ProxeusFSFilterer) ParseUpdatedEvent(log types.Log) (*ProxeusFSUpdatedEvent, error) {
	event := new(ProxeusFSUpdatedEvent)
	if err := _ProxeusFS.contract.UnpackLog(event, "UpdatedEvent", log); err != nil {
		return nil, err
	}
	return event, nil
}
