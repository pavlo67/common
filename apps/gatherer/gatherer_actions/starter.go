package gatherer_actions

import (
	"fmt"
	"strconv"

	"github.com/pkg/errors"

	"github.com/pavlo67/workshop/common"
	"github.com/pavlo67/workshop/common/config"
	"github.com/pavlo67/workshop/common/joiner"
	"github.com/pavlo67/workshop/common/libraries/filelib"
	"github.com/pavlo67/workshop/common/logger"
	"github.com/pavlo67/workshop/common/server/server_http"
	"github.com/pavlo67/workshop/common/starter"
	"github.com/pavlo67/workshop/components/transport"
)

func Starter() starter.Operator {
	return &gathererStarter{}
}

var l logger.Operator

var _ starter.Operator = &gathererStarter{}

type gathererStarter struct {
	receiverHandlerKey joiner.InterfaceKey
}

func (ss *gathererStarter) Name() string {
	return logger.GetCallInfo().PackageName
}

func (ss *gathererStarter) Init(cfgCommon, cfg *config.Config, lCommon logger.Operator, options common.Map) ([]common.Map, error) {
	l = lCommon
	if l == nil {
		return nil, fmt.Errorf("no logger for %s:-(", ss.Name())
	}

	ss.receiverHandlerKey = joiner.InterfaceKey(options.StringDefault("receiver_handler_key", string(transport.HandlerInterfaceKey)))

	return nil, nil
}

func (ss *gathererStarter) Setup() error {
	return nil
}

func (ss *gathererStarter) Run(joinerOp joiner.Operator) error {

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
	//schOp, ok := joiner.Interface(taskscheduler.HandlerKey).(taskscheduler.Operator)
	//if !ok {
	//	l.Fatalf("no scheduler.Operator with key %s", taskscheduler.HandlerKey)
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

	var endpoints = server_http.Endpoints{
		"receive": {Path: "/v1/receive", Tags: []string{"transport"}, HandlerKey: ss.receiverHandlerKey},
	}

	for key, ep := range endpoints {
		ep.Handler, ok = joinerOp.Interface(ep.HandlerKey).(*server_http.Endpoint)
		if !ok {
			return errors.Errorf("no server_http.Endpoint with key %s", ep.HandlerKey)
		}
		endpoints[key] = ep
	}

	cfg := server_http.Config{
		Title:     "Pavlo's Gatherer REST API",
		Version:   "0.0.1",
		Prefix:    "/gatherer",
		Endpoints: endpoints,
	}

	err := server_http.InitEndpointsWithSwaggerV2(
		cfg,
		":"+strconv.Itoa(srvPort),
		srvOp,
		filelib.CurrentPath()+"api-docs/",
		"swagger.json",
		"api-docs",
		l,
	)

	go func() {
		srvOp.Start()
		if err != nil {
			l.Error("on srvOp.Start(): ", err)
		}
	}()
	return nil
}
