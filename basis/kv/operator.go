package kv

import (
	"time"

	"github.com/pavlo67/workshop/basis/joiner"
)

const InterfaceKey joiner.InterfaceKey = "kv"

type Item struct {
	Key      string
	Value    string
	StoredAt time.Time
}

type Operator interface {
	Set(key string, value string) error
	Get(key string) (*Item, error)
}
