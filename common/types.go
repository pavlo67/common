package common

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
	IsEqualTo(Description) Error
	ValuesAreEqual(value1, value2 interface{}) Error
}

type Content interface {
	NewEmpty() Content
	Import(List, Description) Error
	Export() (List, Description, Error)
}
