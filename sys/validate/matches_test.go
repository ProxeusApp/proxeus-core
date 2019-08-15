package validate

import (
	"testing"
)

func TestMatches(t *testing.T) {
	//covered by required
	if errs := FieldByStrRules("", "matches=^omfg$"); errs != nil {
		t.Error(errs)
	}
	if errs := FieldByStrRules("omfg", "matches=^omfg$"); errs != nil {
		t.Error(errs)
	}
	if errs := FieldByStrRules("099", "matches=^[0-9]{3}$"); errs != nil {
		t.Error(errs)
	}
	if errs := FieldByStrRules("099asasdf", "matches=^.{6}"); errs != nil {
		t.Error(errs)
	}
}
