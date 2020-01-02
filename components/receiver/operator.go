package receiver

import (
	"github.com/pavlo67/workshop/common/joiner"
	"github.com/pavlo67/workshop/common/types"

	"github.com/pavlo67/workshop/components/packs"
)

const InterfaceKey joiner.InterfaceKey = "receiver"
const ActionInterfaceKey joiner.InterfaceKey = "receiver_action"

type Operator interface {
	AddHandler(typeKey types.Key, handler packs.Handler) error
	RemoveHandler(types.Key)
}
