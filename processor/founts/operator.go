package founts

import (
	"time"

	"github.com/pavlo67/associatio/processor"
	"github.com/pavlo67/associatio/starter/joiner"
)

const InterfaceKey joiner.InterfaceKey = "founts"

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
