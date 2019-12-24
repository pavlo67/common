package taskqueue_postgres

import (
	"github.com/pavlo67/workshop/common/joiner"
	"github.com/pavlo67/workshop/constructions/taskqueue"
)

func Process(queueOp taskqueue.Operator, joinerOp joiner.Operator, tasks []taskqueue.Item) error {

	//for ; len(tasks) > 0; tasks = tasks[1:] {
	//	task := tasks[0]
	//	workerOp, ok := joinerOp.Interface(task.WorkerType).(worker.Operator)
	//	if !ok {
	//		return errors.Errorf("no worker.Operator for task (%#v)", task)
	//	}
	//
	//	var errs common.Errors
	//	var posterior []common.ID
	//
	//	info, tasks, err := workerOp.Run(task.Task, time.Now())
	//	if err != nil {
	//		errs = append(errs, err)
	//	}
	//
	//	for _, t := range tasks {
	//		id, err := op.Create(t, nil)
	//		if err != nil {
	//			errs = append(errs, err)
	//		}
	//		if id != "" {
	//			posterior = append(posterior, id)
	//		}
	//	}
	//
	//	// .SaveStatus
	//
	//	err = op.SetResult(task.ID, Result{
	//		Timing:    Timing{},
	//		Success:   len(errs) < 1,
	//		Info:      info,
	//		Posterior: posterior,
	//	}, nil)
	//
	//	// TODO!!!
	//	// if len(errs) > 0
	//	// if err != nil
	//
	//	// TODO: ???
	//	if task.ContinueImmediate {
	//		break
	//	}
	//}

	return nil
}
