package runner_factory_goroutine

import (
	"encoding/json"
	"time"

	"github.com/pavlo67/workshop/components/runner_factory"

	"github.com/pkg/errors"

	"github.com/pavlo67/workshop/common"
	"github.com/pavlo67/workshop/common/crud"
	"github.com/pavlo67/workshop/common/joiner"

	"github.com/pavlo67/workshop/components/packs"
	"github.com/pavlo67/workshop/components/runner"
	"github.com/pavlo67/workshop/components/tasks"
	"github.com/pavlo67/workshop/components/transport"
)

func New(tasksOp tasks.Operator, joinerOp joiner.Operator) (runner_factory.Factory, error) {
	if joinerOp == nil {
		return nil, errors.New("no joiner.Operator to create runner.Factory")
	}

	return &runnerFactory{
		joinerOp: joinerOp,
		tasksOp:  tasksOp,
	}, nil
}

// runner.Factory -------------------------------------------------------------------

var _ runner_factory.Factory = &runnerFactory{}

type runnerFactory struct {
	joinerOp joiner.Operator
	tasksOp  tasks.Operator
}

func (rf runnerFactory) ItemRunner(item tasks.Item, saveOptions *crud.SaveOptions, transpOp transport.Operator, listener *transport.Listener) (runner.Operator, error) {
	if transpOp == nil {
		return nil, errors.Errorf("on runnerFactory.ItemRunner(): no transport.Operator for data(%#v)", item)
	}

	// TODO!!! check if listener is valid

	actor, ok := rf.joinerOp.Interface(runner.DataInterfaceKey(item.TypeKey)).(runner.Actor)
	if !ok {
		return nil, errors.Errorf("on runnerFactory.ItemRunner(): no runner.Actor with key %s to init new runner for data(%#v)", item.TypeKey, item)
	}

	return &runnerOp{
		tasksOp: rf.tasksOp,
		taskID:  item.ID,
		data:    item.Data,
		actor:   actor,

		transportOp: transpOp,
		listener:    listener,
	}, nil

}

func (rf runnerFactory) TaskRunner(data crud.Data, saveOptions *crud.SaveOptions, transpOp transport.Operator, listener *transport.Listener) (runner.Operator, common.ID,
	error) {
	if transpOp == nil && listener != nil {
		return nil, "", errors.Errorf("on runnerFactory.TaskRunner(): no transport.Operator for data(%#v) with listener (%#v)", data, *listener)
	}

	// TODO!!! check if listener is valid

	actor, ok := rf.joinerOp.Interface(runner.DataInterfaceKey(data.TypeKey)).(runner.Actor)
	if !ok {
		return nil, "", errors.Errorf("on runnerFactory.TaskRunner(): no runner.Actor with key %s to init new runner for data(%#v)", data.TypeKey, data)
	}

	id, err := rf.tasksOp.Save(data, saveOptions)
	if err != nil {
		return nil, "", errors.Errorf("on runnerFactory.TaskRunner(): can'trf.tasksOp.Save(%#v, nil): %s", data, err)
	}
	return &runnerOp{
		tasksOp: rf.tasksOp,
		taskID:  id,
		data:    data,
		actor:   actor,

		transportOp: transpOp,
		listener:    listener,
	}, id, nil
}

// runner.Operator ------------------------------------------------------------------

var _ runner.Operator = &runnerOp{}

type runnerOp struct {
	tasksOp tasks.Operator
	taskID  common.ID
	data    crud.Data
	actor   runner.Actor

	transportOp transport.Operator
	listener    *transport.Listener
}

const onRun = "on runnerOp.Run(): "

func (r runnerOp) Run() (estimate *runner.Estimate, err error) {
	if r.actor == nil {
		return nil, errors.New(onRun + "no runnerOp.actor")
	}

	var params common.Map
	err = json.Unmarshal(r.data.Content, &params)
	if err != nil {
		err1 := r.tasksOp.Finish(r.taskID, tasks.Result{ErrStr: "on json.Unmarshal(r.data.Content, &params): " + err.Error()}, nil)
		if err1 != nil {
			err = errors.Wrap(err, err1.Error())
		}
	}

	l.Info("at runnerOp.Run(): %s --> %#v", r.data.Content, params)

	estimate, err = r.actor.Init(params)
	if err != nil {
		err1 := r.tasksOp.Finish(r.taskID, tasks.Result{ErrStr: "on r.actor.Init(params): " + err.Error()}, nil)
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

	var task *crud.Data

	response := info["response"]
	switch v := response.(type) {
	case crud.Data:
		task = &v
	case *crud.Data:
		task = v
	}

	if task != nil && r.transportOp != nil && r.listener != nil {
		_, _, err := r.transportOp.Send(&packs.Pack{
			From:    "", // TODO ???
			To:      r.listener.SenderKey,
			Options: nil,
			Data:    *task,
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
