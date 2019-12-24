package taskqueue

import (
	"github.com/pavlo67/workshop/common"
	"github.com/pavlo67/workshop/common/crud"
	"github.com/pavlo67/workshop/common/selectors"
	"github.com/pavlo67/workshop/components/worker"
)

type Item struct {
	worker.Task `             bson:",inline"       json:",inline"`
	Status      `             bson:",inline"       json:",inline"`
	ID          common.ID    `bson:"_id,omitempty" json:",omitempty"`
	Results     []Result     `bson:",omitempty"    json:",omitempty"`
	History     crud.History `bson:",omitempty"    json:",omitempty"`
}

type Operator interface {
	Save(worker.Task, *crud.SaveOptions) (common.ID, error)
	Remove(common.ID, *crud.RemoveOptions) error
	Read(common.ID, *crud.GetOptions) (*Item, error)
	List(*selectors.Term, *crud.GetOptions) ([]Item, error)

	StartProcess() error
	StopProcess() error

	// SetStatus(common.ID, Status, *crud.SaveOptions) error
	// SetResult(common.ID, Result, *crud.SaveOptions) error

}
