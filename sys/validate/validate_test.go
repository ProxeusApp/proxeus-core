package validate

import (
	"fmt"
	"testing"
)

func TestErrors_Translate(t *testing.T) {
	rules := Rules{"number": true}
	verrs := Field("bla", rules)
	if len(verrs) != 1 {
		t.Error(verrs)
		return
	}
	verrs.Translate(func(key string, args ...string) string {
		if key == "number invalid" {
			return "number is invalid"
		}
		return key
	})
	if verrs[0].Msg != "number is invalid" {
		t.Error("translation not working")
	}
}

func TestErrorMap_Translate(t *testing.T) {
	var verrs ErrorMap
	err := Struct(&struct {
		Number string `validate:"number=true"`
	}{Number: "bla"})
	if v, ok := err.(ErrorMap); ok {
		verrs = v
	}
	if len(verrs) != 1 {
		t.Error(verrs)
		return
	}
	verrs.Translate(func(key string, args ...string) string {
		if key == "number invalid" {
			return "number is invalid"
		}
		return key
	})
	if verrs["Number"][0].Msg != "number is invalid" {
		t.Error("translation not working")
	}
}

func TestMakeRulesFromString(t *testing.T) {
	var str = `min=3,max=2,children=[max=2,required=true],required=1,matches=^[a-zA-Z]*$`
	if len(fmt.Sprintf("%v", makeRules(str))) != len(`map[required:1 matches:^[a-zA-Z]*$ children:map[max:2 required:true] min:3 max:2]`) {
		t.Error("wrong rules")
	}
	str = `required=1,matches=^[a-zA-Z]*$,children=[max=2,required=true]`
	if rules := fmt.Sprintf("%v", makeRules(str)); len(rules) != len(`map[required:1 matches:^[a-zA-Z]*$ children:map[max:2 required:true]]`) {
		t.Error("wrong rules", rules)
	}
	str = `children=[max=2,required=true],required=1,matches=^[a-zA-Z]*$`
	if rules := fmt.Sprintf("%v", makeRules(str)); len(rules) != len(`map[required:1 matches:^[a-zA-Z]*$ children:map[max:2 required:true]]`) {
		t.Error("wrong rules", rules)
	}
	str = `children=[max=2,required=true]`
	if rules := fmt.Sprintf("%v", makeRules(str)); len(rules) != len(`map[children:map[max:2 required:true]]`) {
		t.Error("wrong rules", rules)
	}
	str = `max=2,required=true`
	if rules := fmt.Sprintf("%v", makeRules(str)); len(rules) != len(`map[max:2 required:true]`) {
		t.Error("wrong rules", rules)
	}
}
