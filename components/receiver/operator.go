package receiver

import (
	"github.com/pavlo67/workshop/common/joiner"
	"github.com/pavlo67/workshop/common/types"

	"github.com/pavlo67/workshop/components/packs"
)

const InterfaceKey joiner.InterfaceKey = "receiver"
const HandlerInterfaceKey joiner.InterfaceKey = "receiver_handler"

type Operator interface {
	AddHandler(typeKey types.Key, handler packs.Handler) error
	RemoveHandler(types.Key)
}
