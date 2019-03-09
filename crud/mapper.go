package crud

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/pavlo67/punctum/basis/encrlib"
)

type Mapper interface {
	Description() Description
	IDFromNative(interface{}) (string, error)
	StringMapToNative(StringMap) (interface{}, error)
	NativeToStringMap(interface{}) (StringMap, error)
}

type StringMap map[string]string

type Field struct {
	Key        string `json:"key"`
	Primary    bool   `json:"primary,omitempty"`
	MaxLength  int    `json:"max_length,omitempty"`
	Creatable  bool   `json:"creatable,omitempty"`
	Updatable  bool   `json:"editable,omitempty"`
	NotEmpty   bool   `json:"not_empty,omitempty"`
	Unique     bool   `json:"unique,omitempty"`
	AutoUnique bool   `json:"auto_unique,omitempty"`

	//Additable  bool      `json:"additable,omitempty"`
	//Type       FieldType `json:"type,omitempty"`
	//Format     string    `json:"format,omitempty"`
}

type Description struct {
	ExemplarNative interface{}
	Fields         []Field `json:"fields,omitempty"`
}

func (descr Description) PrimaryKeys() []string {
	var primaryKeys []string

	for _, field := range descr.Fields {
		if field.Primary {
			primaryKeys = append(primaryKeys, field.Key)
		}
	}

	return primaryKeys
}

//func (descr Description) Field(key string) *Field {
//	for _, field := range descr.Fields {
//		if field.Key == key {
//			return &field
//		}
//	}
//
//	return nil
//}

func MapperTest(t *testing.T, mapper Mapper, fields []Field) {

	// TODO: NotEmpty-fields
	// TODO: check all fields with reflect

	data0 := StringMap{}

	for _, f := range fields {
		if f.Creatable {
			maxLength := f.MaxLength
			if maxLength == 0 {
				maxLength = 1
			}
			data0[f.Key] = encrlib.RandomString(maxLength)
		}
	}
	mapped, err := mapper.StringMapToNative(data0)
	require.NoError(t, err)

	data1, err := mapper.NativeToStringMap(mapped)

	for k, v := range data0 {
		require.Equal(t, v, data1[k], fmt.Sprintf("??? for key: '%s'", k))
	}

}

//type FieldType string
//const TypeString FieldType = "str"
//const TypeBytes FieldType = "bytes"
//const TypeUInt FieldType = "uint"
//const TypeInt FieldType = "int"
//const TypeFloat FieldType = "float"
//const TypeBool FieldType = "bool"
//const TypeTime FieldType = "time"
