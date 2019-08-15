package validate

import (
	"reflect"
	"regexp"
)

const msgEmailInvalid = "email invalid"

var emailRegexp = regexp.MustCompile("^[a-zA-Z0-9.!#$%&’*+/=?^_`{|}~-]+@[a-zA-Z0-9-]+(?:\\.[a-zA-Z0-9-]+)*$")

//email validates a string if the spec is { email:bool }
func (v *validator) email() {
	if v.val.Kind() == reflect.String {
		e := v.val.String()
		//"" is covered by required
		if e != "" && !emailRegexp.MatchString(e) {
			v.add(&Error{Msg: msgEmailInvalid})
		}
	}
}
