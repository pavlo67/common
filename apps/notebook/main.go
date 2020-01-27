package main

import (
	"flag"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/pavlo67/workshop/common"
	"github.com/pavlo67/workshop/common/auth"
	"github.com/pavlo67/workshop/common/auth/auth_http"
	"github.com/pavlo67/workshop/common/auth/auth_jwt"
	"github.com/pavlo67/workshop/common/auth/auth_users"
	"github.com/pavlo67/workshop/common/config"
	"github.com/pavlo67/workshop/common/control"
	"github.com/pavlo67/workshop/common/libraries/filelib"
	"github.com/pavlo67/workshop/common/logger"
	"github.com/pavlo67/workshop/common/serializer"
	"github.com/pavlo67/workshop/common/server/server_http/server_http_jschmhr"
	"github.com/pavlo67/workshop/common/starter"
	"github.com/pavlo67/workshop/common/users/users_stub"

	"github.com/pavlo67/workshop/components/data/data_pg"
	"github.com/pavlo67/workshop/components/datatagged"
	"github.com/pavlo67/workshop/components/packs/packs_pg"
	"github.com/pavlo67/workshop/components/runner_factory/runner_factory_goroutine"
	"github.com/pavlo67/workshop/components/storage"
	"github.com/pavlo67/workshop/components/storage/storage_server_http"
	"github.com/pavlo67/workshop/components/tagger/tagger_pg"
	"github.com/pavlo67/workshop/components/tasks/tasks_pg"
	"github.com/pavlo67/workshop/components/transport"
	"github.com/pavlo67/workshop/components/transport/transport_http"
	"github.com/pavlo67/workshop/components/transportrouter/transportrouter_stub"

	"github.com/pavlo67/workshop/apps/notebook/notebook_actions"
)

var (
	BuildDate   = "unknown"
	BuildTag    = "unknown"
	BuildCommit = "unknown"
)

const serviceNameDefault = "notebook"

func main() {
	rand.Seed(time.Now().UnixNano())

	var versionOnly, copyFlow bool
	var serviceName string
	flag.BoolVar(&versionOnly, "version_only", false, "show build vars only")
	flag.BoolVar(&copyFlow, "copy_flow", false, "copy flow data immediately")
	flag.StringVar(&serviceName, "service", serviceNameDefault, "service name")
	flag.Parse()

	log.Printf("builded: %s, tag: %s, commit: %s\n", BuildDate, BuildTag, BuildCommit)

	if versionOnly {
		return
	}

	// logger

	l, err := logger.Init(logger.Config{})
	if err != nil {
		log.Fatal(err)
	}

	// getting config environments

	currentPath := filelib.CurrentPath()
	configEnv, ok := os.LookupEnv("ENV")
	if !ok {
		configEnv = "local"
	}

	// routes config

	configCommonPath := currentPath + "../../environments/common." + configEnv + ".yaml"
	cfgCommon, err := config.Get(configCommonPath, serializer.MarshalerYAML)
	if err != nil {
		l.Fatal(err)
	}
	var routesCfg map[string]config.Access
	err = cfgCommon.Value("routes", &routesCfg)
	if err != nil {
		l.Fatal(err)
	}

	var port int
	if serviceAccess, ok := routesCfg[serviceName]; ok {
		port = serviceAccess.Port
	} else {
		l.Fatalf("no access config for key %s (%#v)", serviceName, routesCfg)
	}

	// notebook config

	configworkspacePath := currentPath + "../../environments/" + serviceName + "." + configEnv + ".yaml"
	cfgworkspace, err := config.Get(configworkspacePath, serializer.MarshalerYAML)
	if err != nil {
		l.Fatal(err)
	}

	// running starters

	label := "NOTEBOOK/PG CLI BUILD"

	starters := []starter.Starter{

		// general purposes components
		{control.Starter(), nil},

		// auth system
		{users_stub.Starter(), nil},
		{auth_users.Starter(), common.Map{"interface_key": auth_users.InterfaceKey}},
		{auth_jwt.Starter(), common.Map{"interface_key": auth_jwt.InterfaceKey}},
		{auth_http.Starter(), common.Map{
			"authorize_handler_key": auth.AuthorizeHandlerKey,
			"set_creds_handler_key": auth.SetCredsHandlerKey,
			"get_creds_handler_key": auth.GetCredsHandlerKey,
		}},

		// tasks system
		{tasks_pg.Starter(), nil},
		{runner_factory_goroutine.Starter(), nil},

		// action managers
		//{scheduler_timeout.Starter(), nil},
		{server_http_jschmhr.Starter(), common.Map{"port": port}},

		// transport system
		{packs_pg.Starter(), nil},
		{transportrouter_stub.Starter(), nil},
		{transport_http.Starter(), common.Map{"handler_key": transport.HandlerInterfaceKey, "domain": serviceName}},

		// database
		{tagger_pg.Starter(), nil},
		{data_pg.Starter(), common.Map{"table": storage.CollectionDefault, "interface_key": storage.DataInterfaceKey}},
		{datatagged.Starter(), common.Map{"data_key": storage.DataInterfaceKey, "interface_key": storage.InterfaceKey}},
		{storage_server_http.Starter(), common.Map{"data_key": storage.InterfaceKey}},

		// TODO: pass the interface_key of data_pg to front_end

		// flow actions
		//{data_pg.Starter(), common.Map{"table": flow.CollectionDefault, "interface_key": flow.DataInterfaceKey, "cleaner_key": flow.CleanerInterfaceKey}},
		//{datatagged.Starter(), common.Map{"data_key": flow.DataInterfaceKey, "interface_key": flow.InterfaceKey}},
		//{flow_server_http.Starter(), common.Map{"data_key": flow.DataInterfaceKey, "interface_key": flow.InterfaceKey}},
		// {sources_stub.Starter(), nil},
		// {flowcopier_task.Starter(), common.Map{"datatagged_key": flow.InterfaceKey}},
		// {flowcleaner_task.Starter(), common.Map{"cleaner_key": flow.CleanerInterfaceKey, "interface_key": flow.CleanerTaskInterfaceKey, "limit": 300000}},

		// actions starter (connecting specific actions to the corresponding action managers)
		{notebook_actions.Starter(), common.Map{
			"authorize_handler_key": auth.AuthorizeHandlerKey,
			"set_creds_handler_key": auth.SetCredsHandlerKey,
			"get_creds_handler_key": auth.GetCredsHandlerKey,

			"transport_handler_key": transport.HandlerInterfaceKey,

			// "copier_task_key":    flow.CopierTaskInterfaceKey,
			// "copy_immediately":   copyFlow,
			// "cleaner_task_key":   flow.CopierTaskInterfaceKey,
		}},
	}

	joinerOp, err := starter.Run(starters, cfgCommon, cfgworkspace, os.Args[1:], label)
	if err != nil {
		l.Fatal(err)
	}
	defer joinerOp.CloseAll()

	notebook_actions.WG.Wait()

}

//manifest, err := manager.ReadManifest(currentPath)
//if err != nil {
//	log.Fatal(err)
//}
//if manifest == nil {
//	log.Fatalf("can't load manifest, no data!")
//}
//for _, key := range manifest.Requested {
//	if os.Getenv(key) == "" {
//		log.Fatalf("no environment value for key '%s'", key)
//	}
//}
