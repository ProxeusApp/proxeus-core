package validate

import "testing"

func TestNumber(t *testing.T) {
	if errs := FieldByStrRules("123", "number=true"); errs != nil {
		t.Error(errs)
	}
	if errs := FieldByStrRules(123, "number=true"); errs != nil {
		t.Error(errs)
	}
	if errs := FieldByStrRules("1.23", "number=true"); errs != nil {
		t.Error(errs)
	}
}
