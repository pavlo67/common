package flow

import (
	"time"

	"github.com/pavlo67/workshop/common/joiner"
)

const InterfaceKey joiner.InterfaceKey = "flow"
const TaggedInterfaceKey joiner.InterfaceKey = "flow_tagged"
const CollectionDefault = "flow"

type Origin struct {
	Source string     `bson:",omitempty"    json:",omitempty"`
	Key    string     `bson:",omitempty"    json:",omitempty"`
	Time   *time.Time `bson:",omitempty"    json:",omitempty"`
	Data   string     `bson:",omitempty"    json:",omitempty"`
}
