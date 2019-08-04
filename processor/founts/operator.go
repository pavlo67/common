package founts

import (
	"time"

	"github.com/pavlo67/constructor/processor"
	"github.com/pavlo67/constructor/starter/joiner"
	"github.com/pavlo67/constructor/structura"
)

const InterfaceKey joiner.InterfaceKey = "founts"

type Item struct {
	ID      string              `bson:"_id,omitempty" json:"id,omitempty"`
	URL     string              `bson:"url,omitempty" json:"url,omitempty"`
	Log     []processor.LogItem `bson:"log,omitempty" json:"log,omitempty"`
	SavedAt time.Time           `bson:"saved_at"      json:"saved_at"`
}

type Operator interface {
	Save(url string, logItems ...processor.LogItem) error
	Read(url string) (*Item, error)
	ReadList(structura.GetOptions) ([]Item, *uint64, error)
	Delete(url string) error
	Close() error
}
