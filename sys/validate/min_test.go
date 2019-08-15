package validate

import "testing"

type minSecond struct {
	AnotherName string `validate:"min=4"`
}
type minStruct struct {
	Age        int               `validate:"min=4"`
	Weight     float64           `validate:"min=12.123123"`
	Name       string            `validate:"min=5"`
	Belongings []string          `validate:"min=2"`
	Children   map[string]string `validate:"min=2"`
	More       minSecond
}

func TestStruct_Min(t *testing.T) {
	if errs := Struct(&minStruct{
		More:       minSecond{AnotherName: "1234"},
		Age:        5,
		Weight:     12.123123,
		Name:       "Artan",
		Belongings: []string{"item1", "item2"},
		Children:   map[string]string{"Adam": "big boy", "Liam": "little boy"},
	}); errs != nil {
		t.Error(errs)
	}
	if errs := Struct(&minStruct{
		Age:        4,
		Weight:     12.123122,
		Name:       "Arta",
		Belongings: []string{"item1"},
		Children:   map[string]string{"Adam": "big boy"},
	}); len(errs.Error()) != len(`{"Weight":[{"msg":"min.lowly"}],"Name":[{"msg":"min.lowly"}],"Belongings":[{"msg":"min.lowly"}],"Children":[{"msg":"min.lowly"}],"More.AnotherName":[{"msg":"min.lowly"}]}`) {
		t.Error(errs)
	}
}

func TestField_Min(t *testing.T) {
	//min should be parsed as an int for checking the length of a string
	if errs := FieldByStrRules("Artan", "min=5.2"); errs != nil {
		t.Error("min validator error", errs)
	}
	//min should be parsed as an int for checking the length of a slice, array or map
	if errs := FieldByStrRules([]byte{1, 2, 3, 4}, "min=5.2"); errs == nil {
		t.Error("min validator error", errs)
	}
	//min should be parsed as an int for checking the length of a slice, array or map
	if errs := FieldByStrRules([]byte{1, 2, 3, 4, 5}, "min=5.2"); errs != nil {
		t.Error("min validator error", errs)
	}
	//min should be parsed as an int for checking the length of a slice, array or map
	if errs := FieldByStrRules(map[int]int{1: 1, 2: 2, 3: 3}, "min=5.2"); errs == nil {
		t.Error("min validator error", errs)
	}
	//min should be parsed as an int for checking the length of a slice, array or map
	if errs := FieldByStrRules(map[int]int{1: 1, 2: 2, 3: 3, 4: 4, 5: 5}, "min=5.2"); errs != nil {
		t.Error("min validator error", errs)
	}
	if errs := FieldByStrRules(6, "min=5.2"); errs != nil {
		t.Error("min validator error", errs)
	}
	if errs := FieldByStrRules(5.2, "min=5.2"); errs != nil {
		t.Error("min validator error", errs)
	}
	if errs := FieldByStrRules(5.1, "min=5.2"); errs == nil {
		t.Error("min validator error", errs)
	}
	if errs := FieldByStrRules(1.5, "min=asdf"); errs == nil {
		t.Error("expected: bad definition of min=asdf, expected int", errs)
	}
	if errs := FieldByStrRules("23424242343", "min=100000,number=t"); errs != nil {
		t.Error(errs)
	}
	if errs := FieldByStrRules("99999", "min=100000,number=t"); errs == nil {
		t.Error(errs)
	}
}
