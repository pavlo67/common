package gatherer_actions

import (
	"fmt"
	"strconv"
	"time"

	"github.com/pavlo67/workshop/components/runner"

	"github.com/pavlo67/workshop/common/scheduler"
	"github.com/pavlo67/workshop/components/flow"

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
	importerTaskKey     joiner.InterfaceKey
	importerImmediately bool
	cleanerTaskKey      joiner.InterfaceKey
	receiverHandlerKey  joiner.InterfaceKey
	schKey              joiner.InterfaceKey
}

func (gs *gathererStarter) Name() string {
	return logger.GetCallInfo().PackageName
}

// TODO!!! customize it
const importPeriod = time.Hour

func (gs *gathererStarter) Init(cfgCommon, cfg *config.Config, lCommon logger.Operator, options common.Map) ([]common.Map, error) {
	l = lCommon
	if l == nil {
		return nil, fmt.Errorf("no logger for %s:-(", gs.Name())
	}

	gs.importerTaskKey = joiner.InterfaceKey(options.StringDefault("importer_task_key", "")) // string(flow.ImporterTaskInterfaceKey)
	gs.importerImmediately = options.IsTrue("import_immediately")
	gs.receiverHandlerKey = joiner.InterfaceKey(options.StringDefault("receiver_handler_key", string(transport.HandlerInterfaceKey)))
	gs.cleanerTaskKey = joiner.InterfaceKey(options.StringDefault("cleaner_task_key", string(flow.CleanerTaskInterfaceKey)))
	gs.schKey = joiner.InterfaceKey(options.StringDefault("scheduler_key", string(scheduler.InterfaceKey)))

	return nil, nil
}

func (gs *gathererStarter) Setup() error {
	return nil
}

func (gs *gathererStarter) Run(joinerOp joiner.Operator) error {

	// scheduling importer task

	var ok bool
	var impTaskOp runner.Actor

	if gs.importerTaskKey != "" {
		impTaskOp, ok = joinerOp.Interface(gs.importerTaskKey).(runner.Actor)
		if !ok {
			l.Fatalf("no actor.ActorKey with key %s", gs.importerTaskKey)
		}

		schOp, ok := joinerOp.Interface(scheduler.InterfaceKey).(scheduler.Operator)
		if !ok {
			l.Fatalf("no scheduler.ActorKey with key %s", scheduler.InterfaceKey)
		}

		taskID, err := schOp.Init(impTaskOp)
		if err != nil {
			l.Fatalf("can't schOp.Init(%#v): %s", impTaskOp, err)
		}

		err = schOp.Run(taskID, importPeriod, gs.importerImmediately)
		if err != nil {
			l.Fatalf("can't schOp.Run(%s, %d, %t): %s", taskID, importPeriod, gs.importerImmediately, err)
		}
	}

	// scheduling cleaner task

	// TODO!!!

	// handling transport receiver

	srvOp, ok := joinerOp.Interface(server_http.InterfaceKey).(server_http.Operator)
	if !ok {
		return errors.Errorf("no server_http.ActorKey with key %s", server_http.InterfaceKey)
	}

	srvPort, ok := joinerOp.Interface(server_http.PortInterfaceKey).(int)
	if !ok {
		return errors.Errorf("no server_http.Port with key %s", server_http.PortInterfaceKey)
	}

	var endpoints = server_http.Endpoints{
		"transport": {Path: "/transport", Tags: []string{"transport"}, HandlerKey: gs.receiverHandlerKey},
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
		Prefix:    "",
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
	if err != nil {
		l.Error("on server_http.InitEndpointsWithSwaggerV2(): ", err)
	}

	WG.Add(1)

	go func() {
		defer WG.Done()
		err := srvOp.Start()
		if err != nil {
			l.Error("on srvOp.Start(): ", err)
		}
	}()
	return nil
}
