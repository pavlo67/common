package packs

import "github.com/pavlo67/workshop/common/crud"

const HandleAction crud.ActionKey = "handle"

type Handler interface {
	Handle(pack *Pack) (*Pack, error)
}
