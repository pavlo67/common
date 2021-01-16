package data

import "github.com/pavlo67/workshop/common"

type Type struct {
	Key      TypeKey
	Exemplar interface{}
}

type Field struct {
	Key  string
	Type Type
}

type Description interface {
	Fields() []Field
	Required() []string
	Essential() []string
	IsEqualTo(Description) common.Error
	ValuesAreEqual(value1, value2 interface{}) common.Error
}

type Content interface {
	NewEmpty() Content
	Import(List, Description) common.Error
	Export() (List, Description, common.Error)
}
