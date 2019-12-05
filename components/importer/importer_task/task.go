package importer_task

import (
	"errors"
	"time"

	"github.com/pavlo67/workshop/common/scheduler"
	"github.com/pavlo67/workshop/components/data"
)

func New(dataOp data.Operator) (scheduler.Task, error) {

	if dataOp == nil {
		return nil, errors.New("on importer_task.New(): data.Operator == nil")
	}

	return &task{dataOp}, nil
}

var _ scheduler.Task = &task{}

type task struct {
	dataOp data.Operator
}

func (it *task) Name() string {
	return "importer_actor"
}

func (it *task) Run(timeSheduled time.Time) error {
	if it == nil {
		return errors.New("on importer_task.Run(): task == nil")
	}

	return LoadAll(it.dataOp)
}
