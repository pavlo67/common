package main

import (
	"log"
	"reflect"
)

func main() {

	type aaa struct{ a int }

	exemplar := aaa{1}

	exemplar1 := reflect.New(reflect.TypeOf(exemplar)).Elem().Interface()

	log.Printf("%#v\n%#v", exemplar, exemplar1)
}
