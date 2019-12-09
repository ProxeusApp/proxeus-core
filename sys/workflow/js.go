package workflow

import (
	"encoding/json"
	"fmt"

	"github.com/robertkrimen/otto"
)

type JS struct {
	vm *otto.Otto
}

func NewJSParser() *JS {
	js := &JS{}
	js.vm = otto.New()
	return js
}

func (js *JS) SetGlobal(data interface{}) error {
	if data != nil {
		if dataMap, ok := data.(map[string]interface{}); ok {
			for name, v := range dataMap {
				if v != nil {
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
		} else {
			b, err := json.Marshal(data)
			if err != nil {
				return err
			}
			_, err = js.vm.Eval(fmt.Sprintf("var global = %s;", string(b)))
			return err
		}

	}
	return nil
}

func (js *JS) Run(src interface{}) (otto.Value, error) {
	return js.vm.Run(src)
}

func (js *JS) Close() error {
	//TODO cannot close otto vm
	js.vm = nil
	return nil
}
