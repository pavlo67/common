package founts

import (
	"time"

	"github.com/pavlo67/constructor/processor"
	"github.com/pavlo67/constructor/starter/joiner"
)

const InterfaceKey joiner.ComponentKey = "founts"

type Item struct {
	URL     string              `json:"url,omitempty"`
	Log     []processor.LogItem `json:"log,omitempty"`
	SavedAt time.Time           `json:"saved_at"`
}

type Operator interface {
	Save(url string, logItems ...processor.LogItem) error
	Read(url string) (*Item, error)
	ReadList(content.ListOptions) ([]Item, *uint64, error)
	Delete(url string) error
	Close() error
}
