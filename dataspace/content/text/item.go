package text

import (
	"github.com/pavlo67/punctum/dataspace/content"
)

const Type content.Type = "text"

var _ content.Item = &Item{}

type Item struct {
	Text     string `bson:"text,omitempty"     json:"text"`
	Language string `bson:"language,omitempty" json:"language"`
	// TODO: change to language.Tag
}

func (text Item) Type() content.Type {
	return Type
}

func (text Item) Key() string {
	return ""
}

func (text Item) Set(interface{}) error {
	return nil
}

func (text Item) Refresh() error {
	return nil
}

func (text Item) String() string {
	return text.Text
}
