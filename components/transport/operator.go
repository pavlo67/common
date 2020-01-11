package transport

import (
	"github.com/pavlo67/workshop/common"
	"github.com/pavlo67/workshop/common/crud"
	"github.com/pavlo67/workshop/common/identity"
	"github.com/pavlo67/workshop/common/joiner"
	"github.com/pavlo67/workshop/components/packs"
)

const InterfaceKey joiner.InterfaceKey = "transport"
const HandlerInterfaceKey joiner.InterfaceKey = "transport_handler"

const SentKey crud.ActionKey = "sent"
const DidntSendKey crud.ActionKey = "didn't send"

type Listener struct {
	SenderKey identity.Key
	PackKey   identity.Key
}

type Operator interface {
	Send(pack *packs.Pack) (sentKey identity.Key, targetTaskID common.ID, err error)
	History(packKey identity.Key) (history []crud.Action, err error)
}
