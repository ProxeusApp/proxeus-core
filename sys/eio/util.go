package eio

import (
	"fmt"
	"log"
	"reflect"
)

func CallMethodByNameOn(inst interface{}, methodName string, args ...interface{}) (res []reflect.Value) {
	defer func() {
		if r := recover(); r != nil {
			res = []reflect.Value{reflect.ValueOf(nil), reflect.ValueOf(fmt.Errorf("%v", r))}
		}
	}()
	vargs := make([]reflect.Value, 0)
	for _, a := range args {
		vargs = append(vargs, reflect.ValueOf(a))
	}
	log.Println(methodName, len(vargs))
	res = reflect.ValueOf(inst).MethodByName(methodName).Call(vargs)
	return
}
