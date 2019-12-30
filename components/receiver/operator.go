package receiver

import (
	"github.com/pavlo67/workshop/common"
	"github.com/pavlo67/workshop/common/joiner"
	"github.com/pavlo67/workshop/common/selectors"

	"github.com/pavlo67/workshop/components/packs"
)

const InterfaceKey joiner.InterfaceKey = "receiver"

type Handler interface {
	Handle(pack *packs.Pack) error
}

type Operator interface {
	AddHandler(id common.ID, handler Handler, selector *selectors.Term) error
	RemoveHandler(common.ID) error
}
