package sources

import (
	"time"

	"github.com/pavlo67/punctum/auth"
	"github.com/pavlo67/punctum/basis"
	"github.com/pavlo67/punctum/crud"
	"github.com/pavlo67/punctum/starter/joiner"
)

const InterfaceKey joiner.InterfaceKey = "sources"

const TypeFieldKey = "type"

type Item struct {
	ID    string
	URL   string
	Title string

	Type      joiner.InterfaceKey
	Params    basis.Info // for Create/Update methods for ex. tags list to set them on each imported item
	ParamsRaw string     // for Read/ReadList methods

	RView     auth.ID
	ROwner    auth.ID
	Managers  auth.Managers
	CreatedAt time.Time
	UpdatedAt *time.Time
}

type Operator interface {
	Create(userIS auth.ID, source Item) (string, error)
	Read(userIS auth.ID, id string) (*Item, error)
	ReadList(userIS auth.ID, options *content.ListOptions) ([]Item, uint64, error)
	Update(userIS auth.ID, source Item) (crud.Result, error)
	Delete(userIS auth.ID, id string) (crud.Result, error)
	Close() error
}
