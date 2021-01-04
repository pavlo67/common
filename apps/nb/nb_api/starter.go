package nb_api

import (
	"fmt"
	"strconv"
	"time"

	"github.com/pavlo67/workshop/common/libraries/filelib"

	"github.com/pavlo67/workshop/components/storage"

	"github.com/pavlo67/workshop/common/auth"

	"github.com/pkg/errors"

	"github.com/pavlo67/workshop/common"
	"github.com/pavlo67/workshop/common/config"
	"github.com/pavlo67/workshop/common/joiner"
	"github.com/pavlo67/workshop/common/logger"
	"github.com/pavlo67/workshop/common/scheduler"
	"github.com/pavlo67/workshop/common/server/server_http"
	"github.com/pavlo67/workshop/common/starter"

	"github.com/pavlo67/workshop/components/flow"
	"github.com/pavlo67/workshop/components/runner"
	"github.com/pavlo67/workshop/components/transport"
)

func Starter() starter.Operator {
	return &workspaceStarter{}
}

var l logger.Operator

var _ starter.Operator = &workspaceStarter{}

type workspaceStarter struct {
	// baseDir string

	authHandlerKey      joiner.InterfaceKey
	setCredsHandlerKey  joiner.InterfaceKey
	getCredsHandlerKey  joiner.InterfaceKey
	copierTaskKey       joiner.InterfaceKey
	copyImmediately     bool
	cleanerTaskKey      joiner.InterfaceKey
	transportHandlerKey joiner.InterfaceKey
	schedulerKey        joiner.InterfaceKey
}

func (gs *workspaceStarter) Name() string {
	return logger.GetCallInfo().PackageName
}

// TODO!!! customize it
const copyPeriod = time.Minute

func (gs *workspaceStarter) Init(cfgCommon, cfg *config.Config, lCommon logger.Operator, options common.Map) ([]common.Map, error) {
	l = lCommon
	if l == nil {
		return nil, fmt.Errorf("no logger for %s:-(", gs.Name())
	}

	// gs.baseDir = options.StringDefault("base_dir", "")
	gs.authHandlerKey = joiner.InterfaceKey(options.StringDefault("auth_handler_key", string(auth.AuthorizeHandlerKey)))
	gs.setCredsHandlerKey = joiner.InterfaceKey(options.StringDefault("set_creds_handler_key", string(auth.SetCredsHandlerKey)))
	gs.getCredsHandlerKey = joiner.InterfaceKey(options.StringDefault("get_creds_handler_key", string(auth.SetCredsHandlerKey)))
	gs.copierTaskKey = joiner.InterfaceKey(options.StringDefault("copier_task_key", "")) // string(flow.copierTaskInterfaceKey)
	gs.copyImmediately = options.IsTrue("copy_immediately")
	gs.transportHandlerKey = joiner.InterfaceKey(options.StringDefault("transport_handler_key", string(transport.HandlerInterfaceKey)))
	gs.cleanerTaskKey = joiner.InterfaceKey(options.StringDefault("cleaner_task_key", string(flow.CleanerTaskInterfaceKey)))
	gs.schedulerKey = joiner.InterfaceKey(options.StringDefault("scheduler_key", string(scheduler.InterfaceKey)))

	return nil, nil
}

func (gs *workspaceStarter) Setup() error {
	return nil
}

func (gs *workspaceStarter) Run(joinerOp joiner.Operator) error {

	// scheduling copier task

	var ok bool
	var copyTaskOp runner.Actor

	if gs.copierTaskKey != "" {
		copyTaskOp, ok = joinerOp.Interface(gs.copierTaskKey).(runner.Actor)
		if !ok {
			l.Fatalf("no actor.ActorKey with key %s", gs.copierTaskKey)
		}

		schOp, ok := joinerOp.Interface(gs.schedulerKey).(scheduler.Operator)
		if !ok {
			l.Fatalf("no scheduler.ActorKey with key %s", gs.schedulerKey)
		}

		taskID, err := schOp.Init(copyTaskOp)
		if err != nil {
			l.Fatalf("can't schOp.Init(%#v): %s", copyTaskOp, err)
		}

		err = schOp.Run(taskID, copyPeriod, gs.copyImmediately)
		if err != nil {
			l.Fatalf("can't schOp.Run(%s, %d, %t): %s", taskID, copyPeriod, gs.copyImmediately, err)
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
		"set_creds": {Path: "/set_creds", Tags: []string{"auth"}, HandlerKey: gs.setCredsHandlerKey},
		"get_creds": {Path: "/get_creds", Tags: []string{"auth"}, HandlerKey: gs.getCredsHandlerKey},
		"authorize": {Path: "/authorize", Tags: []string{"auth"}, HandlerKey: gs.authHandlerKey},

		"transport": {Path: "/transport", Tags: []string{"transport"}, HandlerKey: gs.transportHandlerKey},

		"read":   {Path: "/notebook/read", Tags: []string{"notebook"}, HandlerKey: storage.ReadInterfaceKey},
		"save":   {Path: "/notebook/save", Tags: []string{"notebook"}, HandlerKey: storage.SaveInterfaceKey},
		"remove": {Path: "/notebook/remove", Tags: []string{"notebook"}, HandlerKey: storage.RemoveInterfaceKey},

		"export": {Path: "/notebook/export", Tags: []string{"notebook"}, HandlerKey: storage.ExportInterfaceKey},

		"recent": {Path: "/notebook/recent", Tags: []string{"notebook"}, HandlerKey: storage.RecentInterfaceKey},
		"tags":   {Path: "/notebook/tags", Tags: []string{"notebook"}, HandlerKey: storage.ListTagsInterfaceKey},
		"tagged": {Path: "/notebook/tagged", Tags: []string{"notebook"}, HandlerKey: storage.ListTaggedInterfaceKey},

		//"flow_read": {Path: "/flow/read", Tags: []string{"flow"}, HandlerKey: flow.ReadInterfaceKey},
		//"flow_list": {Path: "/flow/list", Tags: []string{"flow"}, HandlerKey: flow.RecentInterfaceKey},
	}

	for key, ep := range endpoints {
		ep.Handler, ok = joinerOp.Interface(ep.HandlerKey).(*server_http.Endpoint)
		if !ok {
			return errors.Errorf("no server_http.Endpoint with key %s", ep.HandlerKey)
		}
		endpoints[key] = ep
	}

	cfg := server_http.Config{
		Title:     "Pavlo's notebook REST API",
		Version:   "0.0.1",
		Prefix:    "",
		Endpoints: endpoints,
	}

	err := server_http.InitEndpointsWithSwaggerV2(
		cfg,
		":"+strconv.Itoa(srvPort),
		srvOp,
		filelib.CurrentPath()+"api-docs/",
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
