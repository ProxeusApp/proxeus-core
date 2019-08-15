package validate

import (
	"testing"
)

func TestErrorsAdd(t *testing.T) {
	errs := Errors{}
	errs.add(&Error{Msg: "first"})
	errs.add(&Error{Msg: "second"})
	errs.add(&Error{Msg: "third"})
	if len(errs) != 3 {
		t.Errorf("unexpected size %v", errs)
	}
	if errs[0].Msg != "first" {
		t.Errorf("unexpected msg")
	}
	if errs[1].Msg != "second" {
		t.Errorf("unexpected msg")
	}
	if errs[2].Msg != "third" {
		t.Errorf("unexpected msg")
	}
}
