package sender

import (
	"github.com/pavlo67/workshop/common/crud"
	"github.com/pavlo67/workshop/common/identity"
	"github.com/pavlo67/workshop/common/joiner"

	"github.com/pavlo67/workshop/components/packs"
)

const InterfaceKey joiner.InterfaceKey = "sender"

const SentKey crud.ActionKey = "sent"
const DidntSendKey crud.ActionKey = "didn't send"

type Operator interface {
	Handle(pack *packs.Pack) (*packs.Pack, error)

	SendOne(pack *packs.Pack, to identity.Key, ignoreProblems bool) (response *packs.Pack, err error)
	History(packKey identity.Key) (trace []crud.Action, err error)
}
