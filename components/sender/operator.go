package sender

import (
	"github.com/pavlo67/workshop/common/crud"
	"github.com/pavlo67/workshop/common/identity"
	"github.com/pavlo67/workshop/common/joiner"

	"github.com/pavlo67/workshop/components/packs"
)

const InterfaceKey joiner.InterfaceKey = "sender"

type Operator interface {
	Send(pack *packs.Pack) (response *packs.Pack, err error)
	Trace(key identity.Key) (trace []crud.Action, err error)
}
