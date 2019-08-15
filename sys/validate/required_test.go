package validate

import (
	"testing"
)

type deepStruct struct {
	IAmTheLastOne string        `validate:"required=true"`
	SomeList      []interface{} `validate:"required=true"`
}

type anotherStruct struct {
	//json name needs to be taken as we use it for outgoing data
	Bla  string `json:"bla" validate:"required=true"`
	Blup int
	Omg  float64 `validate:"required=true"`
	Lol  bool
	Deep *deepStruct `validate:"required=true"`
}

type mySpecialStruct struct {
	Name     string `validate:"required=true"`
	Bla      bool
	Another  anotherStruct
	SomeList []anotherStruct `validate:"required=true"`
}

func TestRequiredValidationWithStruct(t *testing.T) {
	err := Struct(&mySpecialStruct{Name: "", SomeList: []anotherStruct{{Bla: ""}}})
	if len(err.Error()) != len(`{"Name":[{"msg":"required"}],"Another.bla":[{"msg":"required"}],"Another.Omg":[{"msg":"required"}],"Another.Deep":[{"msg":"required"}],"SomeList.0.bla":[{"msg":"required"}],"SomeList.0.Omg":[{"msg":"required"}],"SomeList.0.Deep":[{"msg":"required"}]}`) {
		t.Error(err.Error())
	}
	err = Struct(&mySpecialStruct{Name: "A", Another: anotherStruct{Bla: "Bla", Blup: 1, Omg: 2.4}, SomeList: []anotherStruct{{Bla: "B", Omg: 1.2, Deep: &deepStruct{}}}})
	if len(err.Error()) != len(`{"Another.Deep":[{"msg":"required"}],"SomeList.0.Deep.IAmTheLastOne":[{"msg":"required"}],"SomeList.0.Deep.SomeList":[{"msg":"required"}]}`) {
		t.Error(err)
	}
	err = Struct(&mySpecialStruct{Name: "A", Another: anotherStruct{Bla: "Bla", Blup: 1, Omg: 2.4, Deep: &deepStruct{IAmTheLastOne: "last", SomeList: []interface{}{1}}}, SomeList: []anotherStruct{{Bla: "B", Omg: 1.2, Deep: &deepStruct{IAmTheLastOne: "last 2", SomeList: []interface{}{2}}}}})
	if err != nil {
		t.Error(err)
	}
}

func TestRequiredValidation(t *testing.T) {
	reqTrue := "required=true"
	if err := FieldByStrRules("bla@abc", reqTrue); err != nil {
		t.Error("required validator error", err)
	}
	if err := FieldByStrRules("", reqTrue); err == nil {
		t.Error("required validator error", err)
	}
	if err := FieldByStrRules(nil, reqTrue); err == nil {
		t.Error("required validator error", err)
	}
	if err := FieldByStrRules([]string{}, reqTrue); err == nil {
		t.Error("required validator error", err)
	}
	if err := FieldByStrRules([]string{"a"}, reqTrue); err != nil {
		t.Error("required validator error", err)
	}
	//we do not support zero val on bool!
	if err := FieldByStrRules(false, reqTrue); err != nil {
		t.Error("required validator error", err)
	}
	if err := FieldByStrRules(true, reqTrue); err != nil {
		t.Error("required validator error", err)
	}
	var somePntr *mySpecialStruct
	if err := FieldByStrRules(somePntr, reqTrue); err == nil {
		t.Error("required validator error", err)
	}
	someZeroStruct := mySpecialStruct{}
	if err := FieldByStrRules(someZeroStruct, reqTrue); err == nil {
		t.Error("required validator error", err)
	}
	someNoneZeroStruct := mySpecialStruct{Name: "abc"}
	if err := FieldByStrRules(someNoneZeroStruct, reqTrue); err != nil {
		t.Error("required validator error", err)
	}
	if err := FieldByStrRules(0, reqTrue); err == nil {
		t.Error("required validator error", err)
	}
	if err := FieldByStrRules(1.1, reqTrue); err != nil {
		t.Error("required validator error", err)
	}
	if err := FieldByStrRules(1, reqTrue); err != nil {
		t.Error("required validator error", err)
	}
	if err := FieldByStrRules(-10, reqTrue); err != nil {
		t.Error("required validator error", err)
	}
	//empty map
	if err := FieldByStrRules(map[string]interface{}{}, reqTrue); err == nil {
		t.Error("required validator error", err)
	}
	if err := FieldByStrRules(map[int]interface{}{}, reqTrue); err == nil {
		t.Error("required validator error", err)
	}
	if err := FieldByStrRules(map[int]interface{}{1: "something"}, reqTrue); err != nil {
		t.Error("required validator error", err)
	}
	if err := FieldByStrRules("", "required=false"); err != nil {
		t.Error("required validator error", err)
	}
}
