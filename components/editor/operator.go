package editor

import "github.com/pavlo67/workshop/common"

type Type string

type Field struct {
	Type
	Options common.Map
	Value   interface{}
}

type Operator interface {
	PrepareToEdit() ([]Field, error)
	SaveEdited([]Field) error
}
