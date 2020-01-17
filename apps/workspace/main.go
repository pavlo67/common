package main

import (
	"flag"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/pavlo67/workshop/common"
	"github.com/pavlo67/workshop/common/auth"
	"github.com/pavlo67/workshop/common/auth/auth_ecdsa"
	"github.com/pavlo67/workshop/common/auth/auth_http"
	"github.com/pavlo67/workshop/common/auth/auth_jwt"
	"github.com/pavlo67/workshop/common/config"
	"github.com/pavlo67/workshop/common/control"
	"github.com/pavlo67/workshop/common/libraries/filelib"
	"github.com/pavlo67/workshop/common/logger"
	"github.com/pavlo67/workshop/common/scheduler/scheduler_timeout"
	"github.com/pavlo67/workshop/common/serializer"
	"github.com/pavlo67/workshop/common/server/server_http/server_http_jschmhr"
	"github.com/pavlo67/workshop/common/starter"

	"github.com/pavlo67/workshop/components/data/data_pg"
	"github.com/pavlo67/workshop/components/datatagged"
	"github.com/pavlo67/workshop/components/flow"
	"github.com/pavlo67/workshop/components/packs/packs_pg"
	"github.com/pavlo67/workshop/components/runner_factory/runner_factory_goroutine"
	"github.com/pavlo67/workshop/components/sources/sources_stub"
	"github.com/pavlo67/workshop/components/tagger/tagger_pg"
	"github.com/pavlo67/workshop/components/tasks/tasks_pg"
	"github.com/pavlo67/workshop/components/transport"
	"github.com/pavlo67/workshop/components/transport/transport_http"
	"github.com/pavlo67/workshop/components/transportrouter/transportrouter_stub"

	"github.com/pavlo67/workshop/apps/workspace/workspace_actions"
)

var (
	BuildDate   = "unknown"
	BuildTag    = "unknown"
	BuildCommit = "unknown"
)

const serviceName = "workspace"

func main() {
	rand.Seed(time.Now().UnixNano())

	var versionOnly, copyImmediately bool
	flag.BoolVar(&versionOnly, "version_only", false, "show build vars only")
	flag.BoolVar(&copyImmediately, "copy_immediately", false, "immediately copy flow data")
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

	// workspace config

	configworkspacePath := currentPath + "../../environments/" + serviceName + "." + configEnv + ".yaml"
	cfgworkspace, err := config.Get(configworkspacePath, serializer.MarshalerYAML)
	if err != nil {
		l.Fatal(err)
	}

	// running starters

	label := "WORKSPACE/PG CLI BUILD"

	starters := []starter.Starter{

		// general purposes components
		{control.Starter(), nil},

		// auth system
		{auth_ecdsa.Starter(), common.Map{"interface_key": auth_ecdsa.InterfaceKey}},
		{auth_jwt.Starter(), common.Map{"interface_key": auth_jwt.InterfaceKey}},
		{auth_http.Starter(), common.Map{"auth_handler_key": auth.AuthorizeHandlerKey, "auth_init_handler_key": auth.AuthInitHandlerKey}},

		// tasks system
		{tasks_pg.Starter(), nil},
		{runner_factory_goroutine.Starter(), nil},

		// action managers
		{scheduler_timeout.Starter(), nil},
		{server_http_jschmhr.Starter(), common.Map{"port": port}},

		// transport system
		{packs_pg.Starter(), nil},
		{transportrouter_stub.Starter(), nil},
		{transport_http.Starter(), common.Map{"handler_key": transport.HandlerInterfaceKey, "domain": serviceName}},

		// database
		{tagger_pg.Starter(), nil},
		{data_pg.Starter(), common.Map{"table": flow.CollectionDefault, "interface_key": flow.DataInterfaceKey, "cleaner_key": flow.CleanerInterfaceKey, "no_tagger": true}},
		{datatagged.Starter(), common.Map{"data_key": flow.DataInterfaceKey, "interface_key": flow.InterfaceKey}},

		// flow actions
		{sources_stub.Starter(), nil},
		// {flowcopier_task.Starter(), common.Map{"datatagged_key": flow.InterfaceKey}},
		// {flowcleaner_task.Starter(), common.Map{"cleaner_key": flow.CleanerInterfaceKey, "interface_key": flow.CleanerTaskInterfaceKey, "limit": 300000}},

		// actions starter (connecting specific actions to the corresponding action managers)
		{workspace_actions.Starter(), common.Map{
			"auth_handler_key":      auth.AuthorizeHandlerKey,
			"auth_init_handler_key": auth.AuthInitHandlerKey,
			// "copier_task_key":       flow.CopierTaskInterfaceKey,
			"copy_immediately":      copyImmediately,
			"transport_handler_key": transport.HandlerInterfaceKey,

			// "cleaner_task_key":   flow.CopierTaskInterfaceKey,
		}},
	}

	joinerOp, err := starter.Run(starters, cfgCommon, cfgworkspace, os.Args[1:], label)
	if err != nil {
		l.Fatal(err)
	}
	defer joinerOp.CloseAll()

	workspace_actions.WG.Wait()

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
