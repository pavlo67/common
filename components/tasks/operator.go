package tasks

import (
	"time"

	"github.com/pavlo67/workshop/common"
	"github.com/pavlo67/workshop/common/crud"
	"github.com/pavlo67/workshop/common/joiner"
	"github.com/pavlo67/workshop/common/selectors"
)

const InterfaceKey joiner.InterfaceKey = "tasks"
const CollectionDefault = "tasks"

type Timing struct {
	StartedAt  *time.Time
	FinishedAt *time.Time
}

func (timing Timing) UTC() Timing {
	if timing.StartedAt != nil {
		startedAt := timing.StartedAt.UTC()
		timing.StartedAt = &startedAt
	}

	if timing.FinishedAt != nil {
		finishedAt := timing.FinishedAt.UTC()
		timing.FinishedAt = &finishedAt
	}

	return timing
}

type Result struct {
	Timing    `              bson:",inline"    json:",inline"`
	ErrStr    string        `bson:",omitempty" json:",omitempty"`
	Info      common.Map    `bson:",omitempty" json:",omitempty"`
	Posterior []joiner.Link `bson:",omitempty" json:",omitempty"`
}

type Status struct {
	Timing `bson:",inline" json:",inline"`
}

type Item struct {
	crud.Data `             bson:",inline"       json:",inline"`
	Status    `             bson:",inline"       json:",inline"`
	ID        common.ID    `bson:"_id,omitempty" json:",omitempty"`
	Results   []Result     `bson:",omitempty"    json:",omitempty"`
	History   crud.History `bson:",omitempty"    json:",omitempty"`
}

type Operator interface {
	Save(crud.Data, *crud.SaveOptions) (common.ID, error)
	Remove(common.ID, *crud.RemoveOptions) error
	Read(common.ID, *crud.GetOptions) (*Item, error)
	List(*selectors.Term, *crud.GetOptions) ([]Item, error)

	Start(common.ID, *crud.SaveOptions) error
	Finish(common.ID, Result, *crud.SaveOptions) error
}
