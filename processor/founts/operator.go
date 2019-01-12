package founts

import (
	"time"

	"github.com/pavlo67/punctum/crud"
	"github.com/pavlo67/punctum/processor"
	"github.com/pavlo67/punctum/starter/joiner"
)

const InterfaceKey joiner.InterfaceKey = "founts"

type Item struct {
	URL       string              `json:"url,omitempty"`
	Log       []processor.LogItem `json:"log,omitempty"`
	CreatedAt time.Time           `json:"created_at"`

	//CreatedAt    time.Time           `json:"-"`
	//CreatedAtStr string              `json:"created_at_str"`
}

type Operator interface {
	Save(url string, logItems ...processor.LogItem) error
	Read(url string) (*Item, error)
	ReadList(*crud.ReadOptions) ([]Item, *uint64, error)
	Delete(url string) error
	Close() error
}
