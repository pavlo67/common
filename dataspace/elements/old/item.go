package old

import (
	"time"

	"github.com/pavlo67/punctum/auth"
	"github.com/pavlo67/punctum/confidenter/rights"
	"github.com/pavlo67/punctum/notebook/links"
)

const Private = "приватний запис"
const InGroup = "запис у групі "
const Public = "запис для загалу"

// this type need for rss importer & note(temporary)
type Text struct {
	Text     string `bson:"text,omitempty"     json:"text"`
	Language string `bson:"language,omitempty" json:"language"`
	// TODO: change to language.Tag
}

type Asked struct {
	ID    string `bson:"_id,omitempty" json:"id"`
	Genus string `bson:"genus"         json:"genus"`
	Name  string `bson:"name"          json:"name"`
}

type Item struct {
	// Asked
	ID    string `bson:"_id,omitempty"      json:"id"`
	Genus string `bson:"genus"              json:"genus"`
	Name  string `bson:"name"               json:"name"`

	Author      string       `bson:"author,omitempty"       json:"author,omitempty"`
	Visibility  string       `bson:"visibility,omitempty"   json:"visibility,omitempty"`
	Brief       string       `bson:"brief"                  json:"brief"`
	Content     string       `bson:"content,omitempty"      json:"content,omitempty"`
	Links       []links.Item `bson:"links,omitempty"        json:"links,omitempty"`
	Tags        string       `bson:"tags.comp,omitempty"         json:"tags.comp,omitempty"` // used for text search only
	CountLinked uint16       `bson:"count_linked,omitempty" json:"count_linked,omitempty"`

	RView    auth.ID         `bson:"r_view,omitempty"   json:"r_view,omitempty"`
	ROwner   auth.ID         `bson:"r_owner,omitempty"  json:"r_owner,omitempty"`
	Managers rights.Managers `bson:"managers,omitempty" json:"managers,omitempty"`

	CreatedAt time.Time  `bson:"created_at,omitempty" json:"created_at"`
	UpdatedAt *time.Time `bson:"updated_at,omitempty" json:"updated_at,omitempty"`

	GlobalIS string `bson:"global_is,omitempty" json:"global_is,omitempty"`
	History  string `bson:"history,omitempty"   json:"history,omitempty"`
	Status   string `bson:"status,omitempty"    json:"status,omitempty"`
}

type Content interface {
	ID() (string, error)
	Key() (string, error)
	FromObject(o *Item) error
	ObjectFrom(oDefault *Item) (*Item, error)
}
