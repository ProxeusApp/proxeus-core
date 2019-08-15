package validate

import (
	"reflect"
	"regexp"
)

const msgPhoneNrInvalid = "phone number invalid"

var phoneNrRegexp = regexp.MustCompile(`^[+]?[0-9]{10,16}$`)
var phoneNrCleanRegexp = regexp.MustCompile(`[-\s/_]`)

func (v *validator) phoneNr() {
	if v.val.Kind() == reflect.String {
		pNr := v.val.String()
		if pNr == "" { //covered by required
			return
		}
		pNr = phoneNrCleanRegexp.ReplaceAllString(pNr, "")
		if !phoneNrRegexp.MatchString(pNr) {
			v.add(&Error{Msg: msgPhoneNrInvalid})
		}
	}
}
