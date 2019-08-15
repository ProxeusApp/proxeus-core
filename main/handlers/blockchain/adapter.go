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
	adapter interface {
		getContractAddress() string
		eventFromLog(out interface{}, lg *types.Log, eventType string) error
	}
)

func NewAdapter(XESContractAddress string, XESABI abi.ABI) adapter {
	me := &xesAdapter{}
	me.XESContractAddress = XESContractAddress
	me.XESABI = XESABI
	return me
}

func (me *xesAdapter) eventFromLog(out interface{}, lg *types.Log, eventType string) error {
	pfsLogUnpacker := bind.NewBoundContract(common.HexToAddress(me.XESContractAddress), me.XESABI,
		nil, nil, nil)

	err := pfsLogUnpacker.UnpackLog(out, eventType, *lg)
	if err != nil {
		return err // not our event type
	}
	return nil
}
func (me *xesAdapter) getContractAddress() string {
	return me.XESContractAddress
}
