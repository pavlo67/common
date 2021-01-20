package types

import (
	"github.com/pavlo67/workshop/common"
)

type Type struct {
	Key      Key
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
	IsEqualTo(Description) error
	ValuesAreEqual(value1, value2 interface{}) error
}

type Content interface {
	NewEmpty() Content
	Import(common.List, Description) error
	Export() (common.List, Description, error)
}
