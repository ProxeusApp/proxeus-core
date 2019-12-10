package blockchain

import (
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

type (
	xesAdapter struct {
		XESContractAddress string
		XESABI             abi.ABI
	}
	Adapter interface {
		GetContractAddress() string
		EventFromLog(out interface{}, lg *types.Log, eventType string) error
	}
)

func NewAdapter(XESContractAddress string, XESABI abi.ABI) *xesAdapter {
	me := &xesAdapter{}
	me.XESContractAddress = XESContractAddress
	me.XESABI = XESABI
	return me
}

func (me *xesAdapter) EventFromLog(out interface{}, lg *types.Log, eventType string) error {
	pfsLogUnpacker := bind.NewBoundContract(common.HexToAddress(me.XESContractAddress), me.XESABI,
		nil, nil, nil)

	err := pfsLogUnpacker.UnpackLog(out, eventType, *lg)
	if err != nil {
		return err // not our event type
	}
	return nil
}
func (me *xesAdapter) GetContractAddress() string {
	return me.XESContractAddress
}
