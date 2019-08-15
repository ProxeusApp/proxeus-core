package validate

import (
	"fmt"
	r "reflect"
	"regexp"
)

const msgNotMatchingRegex = "doesn't match pattern"
const msgBadDefinitionForMatches = "bad definition for matches"

func (v *validator) matches(spec string) {
	k := v.val.Kind()
	if k != r.UnsafePointer &&
		k != r.Ptr &&
		k != r.Struct &&
		k != r.Array &&
		k != r.Slice &&
		k != r.Map &&
		k != r.Chan &&
		k != r.Func &&
		k != r.Invalid &&
		k != r.Uintptr {
		if spec == "" {
			v.add(&Error{Msg: msgBadDefinitionForMatches})
			return
		}
		if r.DeepEqual(v.val.Interface(), r.Zero(v.val.Type()).Interface()) {
			//covered by required
			return
		}
		reg, err := regexp.Compile(spec)
		if err != nil {
			v.add(&Error{Msg: msgBadDefinitionForMatches})
			return
		}
		val := fmt.Sprintf("%v", v.val.Interface())
		if !reg.MatchString(val) {
			v.add(&Error{Msg: msgNotMatchingRegex})
		}
	}
}
