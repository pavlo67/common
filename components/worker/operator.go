package worker

import (
	"github.com/pavlo67/workshop/common"
	"github.com/pavlo67/workshop/common/joiner"
)

type Task struct {
	WorkerType        joiner.InterfaceKey `bson:",omitempty" json:",omitempty"`
	Params            interface{}         `bson:"-"          json:",omitempty"`
	ParamsRaw         []byte              `bson:",omitempty" json:",omitempty"` // shouldn't be used directly
	ContinueImmediate bool                `bson:",omitempty" json:",omitempty"`
}

type Operator interface {
	Name() string
	Run(task *Task, label string) (info common.Map, posterior []joiner.Link, err error)
}
