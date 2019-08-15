package validate

import "testing"

type maxSecond struct {
	AnotherName string `validate:"max=4"`
}

type maxStruct struct {
	Age        int               `validate:"max=8"`
	Weight     float64           `validate:"max=12.123123"`
	Name       string            `validate:"max=5"`
	Belongings []string          `validate:"max=3,children=[min=1,max=10,required=true]"`
	Children   map[string]string `validate:"max=2,children=[min=1,max=10,required=true],min=1"`
	More       maxSecond
}

func TestStruct_Max(t *testing.T) {
	if errs := Struct(&maxStruct{
		More:       maxSecond{AnotherName: "1234"},
		Age:        5,
		Weight:     12.123123,
		Name:       "Artan",
		Belongings: []string{"item0", "item1", "item2"},
		Children:   map[string]string{"Adam": "big boy", "Liam": "little boy"},
	}); errs != nil {
		t.Error(errs)
	}
	if errs := Struct(&maxStruct{
		Age:        4,
		Weight:     12.123122,
		Name:       "Arta",
		Belongings: []string{"item", "item1"},
		Children:   map[string]string{"Adam": "big boy"},
	}); errs != nil {
		t.Error(errs)
	}
}

func TestField_Max(t *testing.T) {
	//min should be parsed as an int for checking the length of a string
	if errs := FieldByStrRules("Artan", "max=5.2"); errs != nil {
		t.Error("min validator error", errs)
	}
	if errs := FieldByStrRules("Artann", "max=5.2"); errs == nil {
		t.Error("min validator error", errs)
	}
	//min should be parsed as an int for checking the length of a slice, array or map
	if errs := FieldByStrRules([]byte{1, 2, 3, 4, 5, 6}, "max=5.2"); errs == nil {
		t.Error("min validator error", errs)
	}
	//min should be parsed as an int for checking the length of a slice, array or map
	if errs := FieldByStrRules([]byte{1, 2, 3, 4, 5}, "max=5.2"); errs != nil {
		t.Error("min validator error", errs)
	}
	//min should be parsed as an int for checking the length of a slice, array or map
	if errs := FieldByStrRules(map[int]int{1: 1, 2: 2, 3: 3, 4: 4, 5: 5, 6: 6}, "max=5.2"); errs == nil {
		t.Error("min validator error", errs)
	}
	//min should be parsed as an int for checking the length of a slice, array or map
	if errs := FieldByStrRules(map[int]int{1: 1, 2: 2, 3: 3, 4: 4, 5: 5}, "max=5.2"); errs != nil {
		t.Error("min validator error", errs)
	}
	if errs := FieldByStrRules(6, "max=5.2"); errs == nil {
		t.Error("min validator error", errs)
	}
	if errs := FieldByStrRules(5.3, "max=5.2"); errs == nil {
		t.Error("min validator error", errs)
	}
	if errs := FieldByStrRules(5.2, "max=5.2"); errs != nil {
		t.Error("min validator error", errs)
	}
	if errs := FieldByStrRules(1.5, "min=asdf"); errs == nil {
		t.Error("expected: bad definition of min=asdf, expected int", errs)
	}
}
