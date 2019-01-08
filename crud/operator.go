package crud

import (
	"github.com/pavlo67/punctum/auth"
	"github.com/pavlo67/punctum/basis/viewshtml"
	"github.com/pavlo67/punctum/crud/selectors"
)

//type Result struct {
//	NumOk int64
//}

type Description struct {
	Title         string            `json:"title,omitempty"`
	Fields        []Field           `json:"fields,omitempty"`
	FieldsKey     []string          `json:"fields_key,omitempty"`
	SortByDefault []string          `json:"sort_by_default,omitempty"`
	Exemplar      interface{}       `json:"exemplar,omitempty"`
	View          []viewshtml.Field `json:"view,omitempty"`
	ViewList      []viewshtml.Field `json:"table_view,omitempty"`
}

func (descr Description) Field(key string) *Field {
	for _, field := range descr.Fields {
		if field.Key == key {
			return &field
		}
	}

	return nil
}

//type JoinTo struct {
//	Table    string             `json:"table,omitempty"`
//	Selector selectors.Selector `json:"selector,omitempty"`
//}

type ReadOptions struct {
	Selector selectors.Selector `json:"selector,omitempty"`
	SortBy   []string           `json:"sort_by,omitempty"`
	Limits   []uint64           `json:"limits,omitempty"`
	Exemplar interface{}        `json:"exemplar,omitempty"`

	//Values    []string           `json:"values,omitempty"`
	//JoinTo    []JoinTo           `json:"join_to,omitempty"`
	//GroupBy   []string           `json:"group_by,omitempty"`
	//ForAdmin  bool               `json:"for_admin,omitempty"`
	//ForExport bool               `json:"for_export,omitempty"`
}

// Operator is a common interface to manage create/read/update/delete operations
type Operator interface {
	Describe() (Description, error)

	StringMapToNative(data StringMap) (interface{}, error)

	NativeToStringMap(interface{}) (StringMap, error)

	IDFromNative(interface{}) (string, error)

	Create(userIS auth.ID, native interface{}) (id string, err error)

	// Read returns crud item (accordingly to requester's rights).
	Read(userIS auth.ID, id string) (interface{}, error)

	// ReadList returns crud items list (accordingly to requester's rights).
	ReadList(userIS auth.ID, options *ReadOptions) ([]interface{}, *uint64, error)

	// Update changes crud item (accordingly to requester's rights).
	Update(userIS auth.ID, native interface{}) error

	// Update deletes crud item (accordingly to requester's rights).
	Delete(userIS auth.ID, id string) error
}

type Cleaner func() error
