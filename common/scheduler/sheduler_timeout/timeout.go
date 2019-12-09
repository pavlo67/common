package scheduler_timeout

import (
	"strconv"
	"time"

	"github.com/pkg/errors"

	"github.com/pavlo67/workshop/common"
	"github.com/pavlo67/workshop/common/scheduler"
)

func New() scheduler.Operator {
	return &schedulerTimeout{
		tasks: map[common.ID]scheduler.Task{},
	}
}

// implementation --------------------------------------------------------------------------------------

var _ scheduler.Operator = &schedulerTimeout{}

type schedulerTimeout struct {
	tasks map[common.ID]scheduler.Task
}

func (st *schedulerTimeout) Init(task scheduler.Task) (common.ID, error) {
	if st.tasks == nil {
		return "", errors.New("schedulerTimeout.tasks == nil")
	}

	id := common.ID(strconv.Itoa(len(st.tasks) + 1))
	st.tasks[id] = &task

	return id, nil
}

func (st *schedulerTimeout) Run(taskID common.ID, interval time.Duration, startImmediately bool) error {
	task := st.tasks[taskID]

	if task == nil {
		return errors.Errorf("schedulerTimeout.tasks[%s] == nil", taskID)
	}

	now := time.Now()

	if interval <= 0 {
		if !startImmediately {
			return errors.Errorf("schedulerTimeout: no action because interval = %d and startImmediately == false", interval)
		}

		err := task.Run(now)
		if err != nil {
			return errors.Errorf("on task(%s).Run(): %s", task.Name(), err)
		}

		return nil
	}

	delta := time.Duration(now.UnixNano() % int64(interval))

	var timeScheduled time.Time

	if startImmediately {
		timeScheduled = now.Add(-delta)
	} else {
		timeScheduled = now.Add(interval - delta)
	}

	for {
		rest := timeScheduled.Sub(time.Now())
		if rest > 0 {
			l.Infof("next task run scheduled on %s", timeScheduled.Format(time.RFC3339))
			time.Sleep(rest)
			continue
		}

		if rest > -interval {
			l.Infof("%s: task (%s) started...", timeScheduled.Format(time.RFC3339), task.Name())

			err := task.Run(timeScheduled)
			if err != nil {
				l.Errorf("on task(%s).Run(): %s", task.Name(), err)
			}

			l.Infof("%s: task (%s) finished", time.Now().Format(time.RFC3339), task.Name())
		}

		timeScheduled = timeScheduled.Add(interval)
	}
}
