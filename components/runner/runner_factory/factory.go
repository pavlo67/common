package runner_factory

import (
	"github.com/pkg/errors"

	"github.com/pavlo67/workshop/common"
	"github.com/pavlo67/workshop/common/joiner"

	"github.com/pavlo67/workshop/components/runner"
	"github.com/pavlo67/workshop/components/tasks"
)

func New(joinerOp joiner.Operator) (runner.Factory, error) {
	if joinerOp == nil {
		return nil, errors.New("no joiner.Operator to create runner.Factory")
	}

	return &runnerFactory{
		joinerOp: joinerOp,
	}, nil
}

// runner.Factory -------------------------------------------------------------------

var _ runner.Factory = &runnerFactory{}

type runnerFactory struct {
	joinerOp joiner.Operator
}

func (rf runnerFactory) NewRunner(item tasks.Item) (runner.Operator, error) {
	actor, ok := rf.joinerOp.Interface(item.ActorKey).(runner.Actor)
	if !ok {
		return nil, errors.Errorf("no runner.Actor with key %s to init new runner for task(%#v)", item.ActorKey, item)
	}

	return &runnerOp{
		item:  item,
		actor: actor,
	}, nil

}

// runner.Operator ------------------------------------------------------------------

var _ runner.Operator = &runnerOp{}

type runnerOp struct {
	item  tasks.Item
	actor runner.Actor
}

const onRun = "on runnerOp.Run(): "

func (r runnerOp) Run() (estimate *runner.Estimate, err error) {
	if r.actor == nil {
		return nil, errors.New(onRun + "no runnerOp.actor")
	}

	estimate, err = r.actor.Init(r.item.Params)
	if err != nil {
		go func() {
			_, _, _ = r.actor.Run()
		}()
	}

	return estimate, err
}

const onCheckResults = "on runnerOp.CheckResults(): "

func (r *runnerOp) CheckResults() (status *tasks.Status, posterior []joiner.Link, info common.Map, err error) {
	return nil, nil, nil, common.ErrNotImplemented
}
