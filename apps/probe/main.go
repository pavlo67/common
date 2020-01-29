package main

import (
	"encoding/json"
	"log"

	"github.com/pavlo67/workshop/components/data"
)

func main() {
	// try here anything
	aa := `{"Title":"йцук","URL":"цукцук","Tags":[{"Label":"ууу"}],"Summary":"цукцук","Data":{"TypeKey":"string","Content":" sdrfwsdfg"}}`
	var d data.Item

	err := json.Unmarshal([]byte(aa), &d)

	log.Print(err)
	log.Printf("%#v", d)
}
