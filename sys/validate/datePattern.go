package validate

import (
	r "reflect"
	s "strings"
	"time"
)

const msgDatePatternInvalid = "date invalid"

func (v *validator) datePattern(spec string) {
	if v.val.Kind() == r.String {
		_, err := goTime(spec, v.val.String())
		if err != nil {
			v.add(&Error{Msg: msgDatePatternInvalid})
		}
	}
}

// goTime can be extended to handle more formats
// reference time: Mon Jan 2 15:04:05 -0700 MST 2006
func goTime(javaPattern, value string) (time.Time, error) {
	goLayout := s.Replace(javaPattern, "MM", "1", 1)
	goLayout = s.Replace(goLayout, "dd", "2", 1)
	goLayout = s.Replace(goLayout, "yyyy", "2006", 1)
	goLayout = s.Replace(goLayout, "HH", "15", 1)
	goLayout = s.Replace(goLayout, "mm", "04", 1)
	goLayout = s.Replace(goLayout, "ss", "05", 1)
	return time.Parse(goLayout, value)
}
