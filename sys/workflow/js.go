package workflow

import (
	"encoding/json"
	"fmt"

	"github.com/ProxeusApp/proxeus-core/sys/model/compatability"

	"github.com/robertkrimen/otto"
)

type JS struct {
	vm *otto.Otto
}

func NewJSParser() *JS {
	return &JS{vm: otto.New()}
}

func (js *JS) SetGlobal(data interface{}) error {
	if data == nil {
		return nil
	}
	if dataMap, ok := compatability.ToMapStringIF(data); ok {
		for name, v := range dataMap {
			if v == nil {
				continue
			}
			b, err := json.Marshal(v)
			if err != nil {
				return err
			}
			_, err = js.vm.Eval(fmt.Sprintf("var %s = %s;", name, string(b)))
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (js *JS) Run(src interface{}) (otto.Value, error) {
	return js.vm.Run(src)
}
