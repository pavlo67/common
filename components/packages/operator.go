package packages

import (
	"github.com/pavlo67/workshop/common"
	"github.com/pavlo67/workshop/common/address"
	"github.com/pavlo67/workshop/common/crud"
	"github.com/pavlo67/workshop/common/types"
)

type Item struct {
	From    address.Item
	To      []address.Item
	Options common.Map

	TypeKey types.Key
	Content interface{}

	History crud.History
}

type Operator interface {
}
