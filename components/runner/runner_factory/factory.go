package runner_factory

import (
	"time"

	"github.com/pkg/errors"

	"github.com/pavlo67/workshop/common"
	"github.com/pavlo67/workshop/common/crud"
	"github.com/pavlo67/workshop/common/joiner"

	"github.com/pavlo67/workshop/components/packs"
	"github.com/pavlo67/workshop/components/runner"
	"github.com/pavlo67/workshop/components/tasks"
	"github.com/pavlo67/workshop/components/transport"
)

func New(tasksOp tasks.Operator, joinerOp joiner.Operator) (runner.Factory, error) {
	if joinerOp == nil {
		return nil, errors.New("no joiner.Operator to create runner.Factory")
	}

	return &runnerFactory{
		joinerOp: joinerOp,
		tasksOp:  tasksOp,
	}, nil
}

// runner.Factory -------------------------------------------------------------------

var _ runner.Factory = &runnerFactory{}

type runnerFactory struct {
	joinerOp joiner.Operator
	tasksOp  tasks.Operator
}

func (rf runnerFactory) ItemRunner(item tasks.Item, saveOptions *crud.SaveOptions, transportOp transport.Operator, listener *transport.Listener) (runner.Operator, error) {
	if transportOp == nil {
		return nil, errors.Errorf("on runnerFactory.ItemRunner(): no transport.Operator for task(%#v)", item)
	}

	// TODO!!! check if listener is valid

	actor, ok := rf.joinerOp.Interface(item.ActorKey).(runner.Actor)
	if !ok {
		return nil, errors.Errorf("on runnerFactory.ItemRunner(): no runner.Actor with key %s to init new runner for task(%#v)", item.ActorKey, item)
	}

	return &runnerOp{
		tasksOp: rf.tasksOp,
		taskID:  item.ID,
		task:    item.Task,
		actor:   actor,

		transportOp: transportOp,
		listener:    listener,
	}, nil

}

func (rf runnerFactory) TaskRunner(task tasks.Task, saveOptions *crud.SaveOptions, transportOp transport.Operator, listener *transport.Listener) (runner.Operator, common.ID,
	error) {
	if transportOp == nil && listener != nil {
		return nil, "", errors.Errorf("on runnerFactory.TaskRunner(): no transport.Operator for task(%#v) with listener (%#v)", task, *listener)
	}

	// TODO!!! check if listener is valid

	actor, ok := rf.joinerOp.Interface(task.ActorKey).(runner.Actor)
	if !ok {
		return nil, "", errors.Errorf("on runnerFactory.TaskRunner(): no runner.Actor with key %s to init new runner for task(%#v)", task.ActorKey, task)
	}

	id, err := rf.tasksOp.Save(task, saveOptions)
	if err != nil {
		return nil, "", errors.Errorf("on runnerFactory.TaskRunner(): can'trf.tasksOp.Save(%#v, nil): %s", task, err)
	}
	return &runnerOp{
		tasksOp: rf.tasksOp,
		taskID:  id,
		task:    task,
		actor:   actor,

		transportOp: transportOp,
		listener:    listener,
	}, id, nil
}

// runner.Operator ------------------------------------------------------------------

var _ runner.Operator = &runnerOp{}

type runnerOp struct {
	tasksOp tasks.Operator
	taskID  common.ID
	task    tasks.Task
	actor   runner.Actor

	transportOp transport.Operator
	listener    *transport.Listener
}

const onRun = "on runnerOp.Run(): "

func (r runnerOp) Run() (estimate *runner.Estimate, err error) {
	if r.actor == nil {
		return nil, errors.New(onRun + "no runnerOp.actor")
	}

	estimate, err = r.actor.Init(r.task.Params)
	if err != nil {
		err1 := r.tasksOp.Finish(r.taskID, tasks.Result{ErrStr: err.Error()}, nil)
		if err1 != nil {
			err = errors.Wrap(err, err1.Error())
		}

	} else {
		err1 := r.tasksOp.Start(r.taskID, nil)
		if err1 != nil {
			err = err1
		} else {
			go r.runOnly()
		}
	}

	return estimate, err
}

func (r runnerOp) runOnly() {
	info, posterior, err := r.actor.Run()
	var errStr string
	if err != nil {
		errStr = err.Error()
	}

	if err1 := r.tasksOp.Finish(r.taskID, tasks.Result{Info: info, Posterior: posterior, ErrStr: errStr}, nil); err1 != nil {
		l.Error(err1) // TODO: wrap it
	}

	var task *tasks.Task

	response := info["response"]
	switch v := response.(type) {
	case tasks.Task:
		task = &v
	case *tasks.Task:
		task = v
	}

	if task != nil && r.transportOp != nil && r.listener != nil {
		_, _, err := r.transportOp.Send(&packs.Pack{
			From:    "", // TODO ???
			To:      r.listener.SenderKey,
			Options: nil,
			Task:    *task,
			History: crud.History{
				{Key: crud.ProducedAction, Related: &joiner.Link{InterfaceKey: packs.InterfaceKey, Key: r.listener.PackKey}},
				{Key: crud.ProducedAction, Related: &joiner.Link{InterfaceKey: tasks.InterfaceKey, ID: r.taskID}, DoneAt: time.Now()},
			},
		})
		if err != nil {
			l.Error(err) // TODO: wrap it
		}
	} else {
		// WTF???
	}
}

const onCheckResults = "on runnerOp.CheckResults(): "

func (r *runnerOp) CheckTask() (item *tasks.Item, err error) {
	return r.tasksOp.Read(r.taskID, nil)
}
