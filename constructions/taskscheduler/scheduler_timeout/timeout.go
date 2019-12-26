package scheduler_timeout

import (
	"strconv"
	"sync"
	"time"

	"github.com/pkg/errors"

	"github.com/pavlo67/workshop/common"
	"github.com/pavlo67/workshop/common/actor"
	"github.com/pavlo67/workshop/constructions/taskscheduler"
)

func New() taskscheduler.Operator {
	return &schedulerTimeout{
		tasks: map[common.Key]*taskWithSignals{},
		mutex: &sync.RWMutex{},
	}
}

var MaxSleep = time.Second * 3

// implementation --------------------------------------------------------------------------------------

var _ taskscheduler.Operator = &schedulerTimeout{}

type taskWithSignals struct {
	actor.Operator
	isRunning            bool
	nextInterval         time.Duration
	nextStartImmediately bool
	mutex                *sync.Mutex
}

type schedulerTimeout struct {
	tasks map[common.Key]*taskWithSignals
	mutex *sync.RWMutex
}

func (st *schedulerTimeout) Init(task actor.Operator) (common.Key, error) {
	if st.tasks == nil {
		return "", errors.New("schedulerTimeout.tasks == nil")
	}

	id := common.Key(strconv.Itoa(len(st.tasks) + 1))

	st.mutex.Lock()
	st.tasks[id] = &taskWithSignals{
		Operator: task,
		mutex:    &sync.Mutex{},
	}
	st.mutex.Unlock()

	return id, nil
}

func (st *schedulerTimeout) Run(taskID common.Key, interval time.Duration, startImmediately bool) error {
	st.mutex.RLock()
	task := st.tasks[taskID]
	st.mutex.RUnlock()

	if task == nil {
		return errors.Errorf("schedulerTimeout.tasks[%s] == nil", taskID)
	}

	if interval <= 0 && !startImmediately {
		return errors.Errorf("schedulerTimeout: no action because interval = %d and startImmediately == false", interval)
	}

	task.mutex.Lock()
	defer task.mutex.Unlock()

	task.nextInterval, task.nextStartImmediately = interval, startImmediately
	if !task.isRunning {
		task.isRunning = true
		go st.run(task)
	}

	return nil
}

func (st *schedulerTimeout) run(task *taskWithSignals) {
	if task == nil {
		return
	}

	defer func() {
		task.mutex.Lock()
		task.isRunning = false
		task.mutex.Unlock()
	}()

	var interval, prevInterval time.Duration
	var timeScheduled time.Time
	var startImmediately, showSheduledTime bool

	for {

		// check settings

		prevInterval = interval

		task.mutex.Lock()
		interval = task.nextInterval
		startImmediately = task.nextStartImmediately
		task.nextStartImmediately = false
		task.mutex.Unlock()

		now := time.Now()

		// run immediately if it's necessary

		if startImmediately {
			_, _, err := task.Run(common.Map{"label": label(now)})
			if err != nil {
				l.Errorf("on task(%s).Run(): %s", task.Name(), err)
			}

			prevInterval = 0 // to prevent double run in the loop
		}

		// running loop

		if interval <= 0 {
			return
		}

		if interval != prevInterval {
			timeScheduled = now.Add(interval - time.Duration(now.UnixNano()%int64(interval)))
			showSheduledTime = true
		}

		rest := timeScheduled.Sub(time.Now())

		// slipping a while

		if rest > 0 {
			if rest > MaxSleep {
				rest = MaxSleep
			}

			if showSheduledTime {
				l.Infof("next task run scheduled on %s", timeScheduled.Format(time.RFC3339))
				showSheduledTime = false
			}
			time.Sleep(rest)
			continue
		}

		// running if the current interval isn't finished yet

		if rest > -interval {
			l.Infof("%s: task (%s) started...", timeScheduled.Format(time.RFC3339), task.Name())

			_, _, err := task.Run(common.Map{"label": label(timeScheduled)})
			if err != nil {
				l.Errorf("on task(%s).Run(): %s", task.Name(), err)
			}

			l.Infof("%s: task (%s) finished", time.Now().Format(time.RFC3339), task.Name())

			showSheduledTime = true
		}

		// moving the scheduled time

		timeScheduled = timeScheduled.Add(interval)
	}

}

func label(t time.Time) string {
	return t.Format(time.RFC3339)[:19]
}

func (st *schedulerTimeout) Stop(taskID common.Key) error {
	st.mutex.RLock()
	task := st.tasks[taskID]
	st.mutex.RUnlock()

	if task == nil {
		return errors.Errorf("schedulerTimeout.tasks[%s] == nil", taskID)
	}

	task.mutex.Lock()
	task.nextInterval = 0
	task.nextStartImmediately = false
	task.mutex.Unlock()

	return nil
}
