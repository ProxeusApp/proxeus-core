package validate

import (
	"fmt"
	r "reflect"
	"strconv"
)

const msgErrWhenParsingNumber = "error when parsing number"
const msgMaxExceeded = "max exceeded"
const msgBadDefinitionOfMax = "bad definition of max"

func (v *validator) max(spec string) {
	m := v.maxFloat(spec)
	if m == -1 {
		return
	}
	switch v.val.Kind() {
	case r.String:
		str := v.val.String()
		if v.isDefinedAs("number") {
			f, _ := strconv.ParseFloat(str, 64)
			//number as a string
			if f > m {
				v.addMaxErr()
			}
		} else {
			if len(str) > int(m) {
				v.addMaxErr()
			}
		}
	case r.Slice, r.Array, r.Map:
		if v.val.Len() > int(m) {
			v.addMaxErr()
		}
	case r.Int, r.Int8, r.Int16, r.Int32, r.Int64, r.Float32, r.Float64, r.Complex64, r.Complex128:
		num := fmt.Sprintf("%v", v.val.Interface())
		f, err := strconv.ParseFloat(num, 64)
		if err != nil {
			v.add(&Error{Msg: msgErrWhenParsingNumber})
			return
		}
		if f > m {
			v.addMaxErr()
		}
	}
}

func (v *validator) addMaxErr() {
	v.add(&Error{Msg: msgMaxExceeded})
}

func (v *validator) maxFloat(spec string) float64 {
	m, err := strconv.ParseFloat(spec, 64)
	if err != nil {
		v.add(&Error{Msg: msgBadDefinitionOfMax})
		return -1
	}
	return m
}
