package validate

import "testing"

func TestStringDate(t *testing.T) {
	rules := Rules{"datePattern": "dd.MM.yyyy"}
	if errs := Field("04.05.2018", rules); len(errs) != 0 {
		t.Error(errs)
	}
	if errs := Field("44701.2016", rules); len(errs) == 0 {
		t.Error(errs)
	}
	if errs := Field("18.353.2212", rules); len(errs) == 0 {
		t.Error(errs)
	}
	if errs := Field("88.21.20125", rules); len(errs) == 0 {
		t.Error(errs)
	}
}
