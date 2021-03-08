package main

import (
	"log"
	"reflect"
)

func main() {
	type Test struct {
		A string
	}

	a := []*Test{&Test{A: "111"}, nil, &Test{A: "211"}}

	v, ok := interface{}(a).([]interface{})

	log.Printf("%#v / %t", v, ok)

	if reflect.TypeOf(a).Kind() == reflect.Slice {
		v := reflect.ValueOf(a)
		lenA := v.Len()

		log.Print(lenA)

		var items []interface{}
		var numNotNil int
		for i := 0; i < lenA; i++ {
			itemI := v.Index(i).Interface()
			items = append(items, itemI)
			if !IsNil(itemI) {
				numNotNil++
			}
		}

		log.Printf("%#v / %d", items, numNotNil)

	}

}

func IsNil(i interface{}) bool {
	if i == nil {
		return true
	}
	switch reflect.TypeOf(i).Kind() {
	case reflect.Ptr, reflect.Map, reflect.Array, reflect.Chan, reflect.Slice:
		return reflect.ValueOf(i).IsNil()
	}
	return false
}
