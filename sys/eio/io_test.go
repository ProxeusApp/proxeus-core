package eio

import (
	"fmt"
	"log"
	"reflect"
	"runtime"
)

type IOType int

const (
	Input         IOType = 11
	Output        IOType = 22
	Somethingelse IOType = 33
)

type Foo struct {
	FirstName string `tag_name:"tag 1"`
	LastName  string `tag_name:"tag 2"`
	Age       int    `tag_name:"tag 3"`
}

func (f *Foo) Reflect() {

}

func (f *Foo) Reflect2(abc string, bla []byte) (interface{}, error) {
	return nil, nil
}

func (f *Foo) reflect() {
	log.Println(Input | Output | Somethingelse)

	val := reflect.ValueOf(f).Elem()
	val2 := reflect.ValueOf(f)
	nm := val2.NumMethod()
	fmt.Println("NumMethod", nm)
	for i := 0; i < nm; i++ {
		m := val2.Method(i)
		log.Println()
		log.Println("method: ", runtime.FuncForPC(reflect.Indirect(m).Pointer()).Name())

	}

	for i := 0; i < val.NumField(); i++ {
		valueField := val.Field(i)
		typeField := val.Type().Field(i)
		tag := typeField.Tag
		kin := val.Kind()
		fmt.Printf("Field Name: %s,\t Field Value: %v,\t Tag Value: %s, \t Kind: %v\n", typeField.Name, valueField.Interface(), tag.Get("tag_name"), kin)
	}
}
