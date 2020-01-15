package crud

import (
	"github.com/pavlo67/workshop/common"
)

type TypeKey string

type Type struct {
	Key      TypeKey
	Exemplar interface{}
}

const KeyString TypeKey = "string"
const KeyHRefImage TypeKey = "href_image"
const KeyHRef TypeKey = "href"

type Data struct {
	TypeKey TypeKey `bson:",omitempty" json:",omitempty"`
	Content []byte  `bson:",omitempty" json:",omitempty"`
}

type Counter map[string]uint64

type Index map[string][]common.ID
