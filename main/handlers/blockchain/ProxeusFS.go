package blockchain

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// ProxeusFSABI is the input ABI used to generate the binding from.
const ProxeusFSABI = "[{\"constant\":false,\"inputs\":[{\"name\":\"_hash\",\"type\":\"bytes32\"},{\"name\":\"_data\",\"type\":\"bytes32\"}],\"name\":\"registerFile\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_hash\",\"type\":\"bytes32\"}],\"name\":\"getFileSigners\",\"outputs\":[{\"name\":\"\",\"type\":\"address[]\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_hash\",\"type\":\"bytes32\"}],\"name\":\"verifyFile\",\"outputs\":[{\"name\":\"exists\",\"type\":\"bool\"},{\"name\":\"issuer\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_hash\",\"type\":\"bytes32\"}],\"name\":\"signFile\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"files\",\"outputs\":[{\"name\":\"issuer\",\"type\":\"address\"},{\"name\":\"data\",\"type\":\"bytes32\"},{\"name\":\"exists\",\"type\":\"bool\"},{\"name\":\"signersCount\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"name\":\"_owner\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"hash\",\"type\":\"bytes32\"}],\"name\":\"UpdatedEvent\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"hash\",\"type\":\"bytes32\"},{\"indexed\":true,\"name\":\"signer\",\"type\":\"address\"}],\"name\":\"FileSignedEvent\",\"type\":\"event\"}]"

// ProxeusFSFileSignedEvent represents a FileSignedEvent event raised by the ProxeusFS contract.
type ProxeusFSFileSignedEvent struct {
	Hash   [32]byte
	Signer common.Address
	Raw    types.Log // Blockchain specific contextual infos
}
