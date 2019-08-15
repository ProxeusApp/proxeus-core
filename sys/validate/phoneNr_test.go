package validate

import (
	"testing"
)

func TestPhoneNr(t *testing.T) {
	if errs := FieldByStrRules("079 123 12 12", "phoneNr=true"); errs != nil {
		t.Error(errs)
	}
	if errs := FieldByStrRules("+41 79 123 12 12", "phoneNr=true"); errs != nil {
		t.Error(errs)
	}
	if errs := FieldByStrRules("+41791231212", "phoneNr=true"); errs != nil {
		t.Error(errs)
	}
	if errs := FieldByStrRules("41 79 123 12 12", "phoneNr=true"); errs != nil {
		t.Error(errs)
	}
	if errs := FieldByStrRules("41791231212", "phoneNr=true"); errs != nil {
		t.Error(errs)
	}
	if errs := FieldByStrRules("0791231212", "phoneNr=true"); errs != nil {
		t.Error(errs)
	}
	if errs := FieldByStrRules("079123ab12", "phoneNr=true"); errs == nil {
		t.Error(errs)
	}
}
