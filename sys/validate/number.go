package validate

import (
	"reflect"
	"strconv"
)

const msgNumberInvalid = "number invalid"

// number validates only strings for float if the spec is { number:bool }
func (v *validator) number() {
	if v.val.Kind() == reflect.String {
		n := v.val.String()
		if n == "" { //is covered by required
			return
		}
		_, err := strconv.ParseFloat(n, 64)
		if err != nil {
			v.add(&Error{Msg: msgNumberInvalid})
		} else {
			v.isNumber = true
		}
	}
}
