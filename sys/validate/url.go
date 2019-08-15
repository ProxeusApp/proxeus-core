package validate

import (
	"reflect"
	"regexp"
)

const msgUrlInvalid = "url invalid"

var urlRegexp = regexp.MustCompile("(\\b(https?|ftp)://)[-A-Za-z0-9+&@#/%?=~_|!:,.;]+[-A-Za-z0-9+&@#/%=~_|]")

// Url validates string if the spec is { url:bool }
func (v *validator) url() {
	if v.val.Kind() == reflect.String {
		strVal := v.val.String()
		if strVal == "" { //covered by required
			return
		}
		if !urlRegexp.MatchString(strVal) {
			v.add(&Error{Msg: msgUrlInvalid})
		}
	}
}
