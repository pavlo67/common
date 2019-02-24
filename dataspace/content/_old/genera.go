package _old

import (
	"time"

	"github.com/pavlo67/partes/crud"
	"github.com/pavlo67/punctum/auth"
	"github.com/pavlo67/punctum/basis"
	"github.com/pavlo67/punctum/basis/viewshtml"
	"github.com/pavlo67/punctum/confidenter/auth"
	"github.com/pavlo67/punctum/confidenter/rights"
	"github.com/pavlo67/punctum/starter/joiner"
)

const InterfaceKeyGenerium joiner.InterfaceKey = "genera"

type Set struct {
	Name   string `bson:"name,omitempty"   json:"name,omitempty"`
	RView  auth.ID
	ROwner auth.ID
	crud.ReadOptions
}

type Genus struct {
	ID  string `bson:"id,omitempty"  json:"id,omitempty"`
	Key string `bson:"key,omitempty" json:"key,omitempty"`

	Name       string         `bson:"name,omitempty"        json:"name,omitempty"`
	NamePlural string         `bson:"name_plural,omitempty" json:"name_plural,omitempty"`
	PxPreview  int            `bson:"px_preview,omitempty"  json:"px_preview,omitempty"`
	Sets       map[string]Set `bson:"sets,omitempty"        json:"sets,omitempty"`
	Translator

	RView    auth.ID         `bson:"r_view,omitempty"   json:"r_view,omitempty"`
	ROwner   auth.ID         `bson:"r_owner,omitempty"  json:"r_owner,omitempty"`
	Managers rights.Managers `bson:"managers,omitempty" json:"managers,omitempty"`

	CreatedAt time.Time  `bson:"created_at,omitempty" json:"created_at,omitempty"`
	UpdatedAt *time.Time `bson:"updated_at,omitempty" json:"updated_at,omitempty"`

	GlobalIS string `bson:"global_is,omitempty" json:"global_is,omitempty"`
	History  string `bson:"history,omitempty"   json:"history,omitempty"`
}

type OperatorGenerium interface {
	Create(userIS auth.ID, o *Genus) (string, error)

	Read(userIS auth.ID, id string) (*Genus, error)

	ReadByKey(userIS auth.ID, key string) (*Genus, error)

	Update(userIS auth.ID, o *Genus) (crud.Result, error)

	Delete(userIS auth.ID, id string) (crud.Result, error)

	ReadList(userIS auth.ID, options *crud.ReadOptions) ([]Genus, uint64, error)

	Close() error
}

//Unmarshal([]byte, interface{}) error
//Marshal(interface{}) ([]byte, error)

type Context struct {
	ID  string
	Tab string
}

type EditForm struct {
	InterfaceKey joiner.InterfaceKey               `bson:"interface_key,omitempty" json:"interface_key,omitempty"`
	Fields       []viewshtml.Field                 `bson:"fields,omitempty"        json:"fields,omitempty"`
	Options      map[string]viewshtml.SelectString `bson:"options,omitempty"       json:"options,omitempty"`
}

type DataRaw = basis.Options
type DataView map[string]string

type Translator interface {
	DataFromObject(userIS auth.ID, o *Item, dataDefault DataRaw) (dataRaw DataRaw, editForm EditForm, errs basis.Errors)

	ObjectFromData(userIS auth.ID, oOld *Item, dataRaw DataRaw, linksList []Item) (o *Item, index interface{}, errs basis.Errors)

	View(userIS auth.ID, o *Item, linkedObjects []Item, context *Context) DataView

	NewItem(user *auth.User, o *Item, context *Context) DataView

	Edit(user *auth.User, o *Item, context *Context) DataView
}
