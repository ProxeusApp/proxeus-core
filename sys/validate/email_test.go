package validate

import (
	"log"
	"testing"
)

type emailStruct struct {
	//json name needs to be taken as we use it for outgoing data
	Email string  `json:"email" validate:"email=true"`
	Omg   float64 `validate:"email=true"` //nothing should happen here
}

type mySpecialEmailStruct struct {
	Email    string `validate:"email=true"`
	Another  *emailStruct
	SomeList []string `validate:"email=true"`
}

func TestEmailOnStruct(t *testing.T) {
	errs := Struct(&mySpecialEmailStruct{SomeList: []string{"hello", "hmmm??"}})
	log.Println(errs)
}

func TestEmailValidation(t *testing.T) {
	rules := Rules{"email": true}
	if err := Field("blaabc", rules); len(err) == 0 {
		t.Error("email validator error", err)
	}

	//anything other than string can not be validated for email
	//it should just be skipped
	if err := Field(1, rules); len(err) != 0 {
		t.Error("email validator error", err)
	}

	//we are validating fields only, no array, maps or structs
	errs := Field([]string{"bla@abc.com", "asdf", "", "abc@bla.com"}, rules)
	if errs != nil {
		t.Error(errs)
	}

	rules = Rules{"email": false}
	if err := Field("blaabc.com", rules); len(err) != 0 {
		t.Error("email validator error", err)
	}
	rules = Rules{"email": true}
	if err := Field("blaabc.com", rules); len(err) == 0 {
		t.Error("email validator error", err)
	}
	rules = Rules{"email": true}
	if err := Field("bla@abc.com", rules); len(err) != 0 {
		t.Error("email validator error", err)
	}
}
