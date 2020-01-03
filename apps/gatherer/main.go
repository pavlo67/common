package main

import (
	"flag"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/pavlo67/workshop/common"
	"github.com/pavlo67/workshop/common/auth/auth_ecdsa"
	"github.com/pavlo67/workshop/common/config"
	"github.com/pavlo67/workshop/common/control"
	"github.com/pavlo67/workshop/common/libraries/filelib"
	"github.com/pavlo67/workshop/common/logger"
	"github.com/pavlo67/workshop/common/serializer"
	"github.com/pavlo67/workshop/common/server/server_http/server_http_jschmhr"
	"github.com/pavlo67/workshop/common/starter"

	"github.com/pavlo67/workshop/components/dataimporter/flowimporter_task"
	"github.com/pavlo67/workshop/components/datatagged"
	"github.com/pavlo67/workshop/components/flow"
	"github.com/pavlo67/workshop/components/flow/flowcleaner_task"
	"github.com/pavlo67/workshop/components/flowcopier"
	"github.com/pavlo67/workshop/components/packs/packs_pg"
	"github.com/pavlo67/workshop/components/receiver"
	"github.com/pavlo67/workshop/components/receiver/receiver_server_http"
	"github.com/pavlo67/workshop/components/taskscheduler/scheduler_timeout"

	"github.com/pavlo67/workshop/apps/gatherer/gatherer_actions"
)

var (
	BuildDate   = "unknown"
	BuildTag    = "unknown"
	BuildCommit = "unknown"
)

const serviceName = "gatherer"

func main() {
	rand.Seed(time.Now().UnixNano())

	var versionOnly bool
	flag.BoolVar(&versionOnly, "version", false, "show build vars only")
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

	// gatherer config

	configGathererPath := currentPath + "../../environments/" + serviceName + "." + configEnv + ".yaml"
	cfgGatherer, err := config.Get(configGathererPath, serializer.MarshalerYAML)
	if err != nil {
		l.Fatal(err)
	}

	// running starters

	label := "GATHERER/SQLITE CLI BUILD"

	starters := []starter.Starter{

		// general purposes components
		{control.Starter(), nil},
		{auth_ecdsa.Starter(), nil},

		// action managers
		{scheduler_timeout.Starter(), nil},
		{server_http_jschmhr.Starter(), common.Map{"port": port}},

		// transport system
		{packs_pg.Starter(), nil},
		{receiver_server_http.Starter(), common.Map{"handler_key": receiver.HandlerInterfaceKey}},

		// database
		{data_pg.Starter(), common.Map{"table": flow.CollectionDefault, "interface_key": flow.DataInterfaceKey, "cleaner_key": flow.CleanerInterfaceKey, "no_tagger": true}},
		{datatagged.Starter(), common.Map{"data_key": flow.DataInterfaceKey, "interface_key": flow.InterfaceKey, "no_tagger": true}},

		// flow actions
		{flowimporter_task.Starter(), common.Map{"datatagged_key": flow.InterfaceKey, "interface_key": flow.ImporterTaskInterfaceKey}},
		{flowcleaner_task.Starter(), common.Map{"cleaner_key": flow.CleanerInterfaceKey, "interface_key": flow.CleanerTaskInterfaceKey, "limit": 300000}},
		{flowcopier.Starter(), common.Map{"datatagged_key": flow.InterfaceKey, "receiver_server_http": true}},

		// actions starter (connecting specific actions to the corresponding action managers)
		{gatherer_actions.Starter(), common.Map{
			"importer_task_key":    flow.ImporterTaskInterfaceKey,
			"cleaner_task_key":     flow.CopierTaskInterfaceKey,
			"receiver_handler_key": receiver.HandlerInterfaceKey,
		}},
	}

	joiner, err := starter.Run(starters, cfgCommon, cfgGatherer, os.Args[1:], label)
	if err != nil {
		l.Fatal(err)
	}
	defer joiner.CloseAll()

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
