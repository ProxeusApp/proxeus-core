package validate

import (
	"fmt"
	r "reflect"
	"strconv"
)

const msgRequired = "required"

//Required validates fields if the spec is required=bool
//possible values: 0,1,t,T,TRUE,true,f,F,FALSE,false
//We do not support zero val on bool! Bool is always correct.
func (v *validator) required() {
	if !v.val.IsValid() {
		v.addRequired()
		return
	}
	switch v.val.Kind() {
	case r.Bool: //we do not support zero val on bool but it is important to have it for bool with unset val.
		//can only happen on the usage with the formSrc not struct
		return
	case r.Int, r.Int8, r.Int16, r.Int32, r.Int64, r.Float32, r.Float64, r.Complex64, r.Complex128:
		f, _ := strconv.ParseFloat(fmt.Sprintf("%v", v.val), 64)
		if f == 0 {
			v.addRequired()
		}
		return
	case r.String:
		if len(v.val.String()) == 0 {
			v.addRequired()
		}
		return
	case r.Struct:
		if r.DeepEqual(v.val.Interface(), r.Zero(v.val.Type()).Interface()) {
			//empty struct
			v.addRequired()
		}
		return
	case r.Slice, r.Array, r.Map:
		if v.val.IsNil() || v.val.Len() == 0 {
			v.addRequired()
		}
		return
	case r.Ptr:
		if v.val.IsNil() {
			v.addRequired()
		}
	}
}

func (v *validator) addRequired() {
	v.errs.add(&Error{Msg: msgRequired})
}
