package gatherer_actions

import (
	"fmt"

	"github.com/pkg/errors"

	"github.com/pavlo67/workshop/common"
	"github.com/pavlo67/workshop/common/config"
	"github.com/pavlo67/workshop/common/joiner"
	"github.com/pavlo67/workshop/common/logger"
	"github.com/pavlo67/workshop/common/server/server_http"
	"github.com/pavlo67/workshop/common/starter"
)

const Name = "workspace_starter"

func Starter() starter.Operator {
	return &workspaceStarter{}
}

var l logger.Operator

var _ starter.Operator = &workspaceStarter{}

type workspaceStarter struct {
	// interfaceKey joiner.InterfaceKey
}

func (ss *workspaceStarter) Name() string {
	return logger.GetCallInfo().PackageName + "/" + Name
}

func (ss *workspaceStarter) Init(cfgCommon, cfg *config.Config, lCommon logger.Operator, options common.Map) ([]common.Map, error) {
	var errs common.Errors

	l = lCommon
	if l == nil {
		errs = append(errs, fmt.Errorf("no logger for %s:-(", Name))
	}

	// interfaceKey = joiner.InterfaceKey(options.StringDefault("interface_key", string(server_http.InterfaceKey)))

	return nil, errs.Err()
}

func (ss *workspaceStarter) Setup() error {
	return nil
}

func (ss *workspaceStarter) Run(joinerOp joiner.Operator) error {

	//// scheduling importer task
	//
	//dataOp, ok := joiner.Interface(flow.DataInterfaceKey).(data.Operator)
	//if !ok {
	//	l.Fatalf("no data.Operator with key %s", flow.DataInterfaceKey)
	//}
	//
	//task, err := flowimporter_task.NewLoader(dataOp)
	//if err != nil {
	//	l.Fatal(err)
	//}
	//
	//schOp, ok := joiner.Interface(taskscheduler.InterfaceKey).(taskscheduler.Operator)
	//if !ok {
	//	l.Fatalf("no scheduler.Operator with key %s", taskscheduler.InterfaceKey)
	//}
	//
	//taskID, err := schOp.Init(task)
	//if err != nil {
	//	l.Fatalf("can't schOp.Init(%#v): %s", task, err)
	//}
	//
	//err = schOp.Run(taskID, time.Hour, true)
	//if err != nil {
	//	l.Fatalf("can't schOp.Run(%s, time.Hour, false): %s", taskID, err)
	//}

	srvOp, ok := joinerOp.Interface(server_http.InterfaceKey).(server_http.Operator)
	if !ok {
		return errors.Errorf("no server_http.Operator with key %s", server_http.InterfaceKey)
	}

	srvPort, ok := joinerOp.Interface(server_http.PortInterfaceKey).(int)
	if !ok {
		return errors.Errorf("no server_http.Port with key %s", server_http.PortInterfaceKey)
	}

	for key, ep := range endpoints {
		ep.Handler, ok = joinerOp.Interface(ep.InterfaceKey).(*server_http.Endpoint)
		if !ok {
			return errors.Errorf("no server_http.Endpoint with key %s", ep.InterfaceKey)
		}
		endpoints[key] = ep
	}

	err := Init(srvOp, srvPort)

	err = srvOp.Start()
	if err != nil {
		l.Error(err)
	}

}
