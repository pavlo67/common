package data

import (
	"github.com/pavlo67/workshop/common/identity"
	"github.com/pavlo67/workshop/common/types"
)

const TypesKeyDataItems identity.Key = "data_items"

var TypeDataItems = types.Type{
	Key:      TypesKeyDataItems,
	Exemplar: []Item(nil),
}
