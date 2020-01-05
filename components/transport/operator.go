package transport

import (
	"github.com/pavlo67/workshop/common/crud"
	"github.com/pavlo67/workshop/common/identity"
	"github.com/pavlo67/workshop/common/joiner"

	"github.com/pavlo67/workshop/components/packs"
)

const InterfaceKey joiner.InterfaceKey = "transport"
const HandlerInterfaceKey joiner.InterfaceKey = "receiver_handler"

const SentKey crud.ActionKey = "sent"
const DidntSendKey crud.ActionKey = "didn't send"

type Operator interface {
	Send(pack *packs.Pack) (sentKey identity.Key, response *packs.Pack, err error)
	AddHandler(receiverKey, typeKey identity.Key, handler packs.Handler) error
	RemoveHandler(receiverKey, typeKey identity.Key)

	History(packKey identity.Key) (trace []crud.Action, err error)
}
