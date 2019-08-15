package validate

import (
	"bytes"
	"fmt"
)

type (
	ErrorMap map[string]Errors
	Errors   []*Error

	Error struct {
		//I for index or key map, makes it easier to locate where msg belongs
		I   interface{} `json:"i,omitempty"`
		Msg string      `json:"msg"`
	}
)

func (me ErrorMap) Translate(trans func(key string, args ...string) string) {
	for _, errs := range me {
		if len(errs) > 0 {
			errs.Translate(trans)
		}
	}
}

func (me *Errors) Translate(trans func(key string, args ...string) string) {
	for _, er := range *me {
		er.Msg = trans(er.Msg)
	}
}

func (me ErrorMap) String() string {
	b := &bytes.Buffer{}
	b.WriteString("{")
	i := 0
	for k, v := range me {
		if i > 0 {
			b.WriteString(",")
		}
		i++
		b.WriteString(fmt.Sprintf(`"%v":%v`, k, v))
	}
	b.WriteString("}")
	return b.String()
}

func (me Error) String() string {
	if me.I == nil {
		return fmt.Sprintf(`{"msg":"%v"}`, me.Msg)
	}
	return fmt.Sprintf(`{"i":"%v","msg":"%v"}`, me.I, me.Msg)
}

func (me Error) Error() string {
	return me.String()
}

func (me Errors) Error() string {
	b := &bytes.Buffer{}
	b.WriteString("[")
	i := 0
	for _, v := range me {
		if i > 0 {
			b.WriteString(",")
		}
		i++
		b.WriteString(fmt.Sprintf(`%v`, v))
	}
	b.WriteString("]")
	return b.String()
}

func (me ErrorMap) Error() string {
	return me.String()
}

func (me *Errors) add(err *Error) {
	*me = append(*me, err)
}
