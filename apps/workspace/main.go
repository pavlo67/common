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
	"github.com/pavlo67/workshop/common/server/server_http"
	"github.com/pavlo67/workshop/common/server/server_http/server_http_jschmhr"
	"github.com/pavlo67/workshop/common/starter"

	"github.com/pavlo67/workshop/components/data/data_sqlite"
	"github.com/pavlo67/workshop/components/datatagged"
	"github.com/pavlo67/workshop/components/flow"
	"github.com/pavlo67/workshop/components/flow/flow_server_http"
	"github.com/pavlo67/workshop/components/flowcleaner/flowcleaner_sqlite"
	"github.com/pavlo67/workshop/components/flowcopier"
	"github.com/pavlo67/workshop/components/packs/packs_pg"
	"github.com/pavlo67/workshop/components/storage"
	"github.com/pavlo67/workshop/components/storage/storage_server_http"
	"github.com/pavlo67/workshop/components/tagger/tagger_sqlite"
	"github.com/pavlo67/workshop/components/taskscheduler/scheduler_timeout"

	"github.com/pavlo67/workshop/apps/workspace/ws_routes"
)

var (
	BuildDate   = "unknown"
	BuildTag    = "unknown"
	BuildCommit = "unknown"
)

const serviceName = "workspace"

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

	// storage config

	configWorkspacePath := currentPath + "../../environments/" + serviceName + "." + configEnv + ".yaml"
	cfgWorkspace, err := config.Get(configWorkspacePath, serializer.MarshalerYAML)
	if err != nil {
		l.Fatal(err)
	}

	// running starters

	const storageTable = storage.CollectionDefault
	const flowTable = flow.CollectionDefault

	label := "WORKSPACE REST BUILD"

	starters := []starter.Starter{
		{control.Starter(), nil},
		{auth_ecdsa.Starter(), nil},

		{scheduler_timeout.Starter(), nil},
		{server_http_jschmhr.Starter(), common.Map{"port": port}},

		{packs_pg.Starter(), nil},
		{tagger_sqlite.Starter(), nil},

		{data_sqlite.Starter(), common.Map{"interface_key": storage.DataInterfaceKey, "table": storageTable}},
		{datatagged.Starter(), common.Map{"interface_key": storage.InterfaceKey, "data_key": storage.DataInterfaceKey}},
		{storage_server_http.Starter(), nil},

		{data_sqlite.Starter(), common.Map{"interface_key": flow.DataInterfaceKey, "table": flowTable}},
		{datatagged.Starter(), common.Map{"interface_key": flow.InterfaceKey, "data_key": flow.DataInterfaceKey}},
		{flowcopier.Starter(), common.Map{"client_http": true, "flow_key": flow.InterfaceKey}},
		{flowcleaner_sqlite.Starter(), common.Map{"limit": 3000, "flow_key": flow.CleanerInterfaceKey}},
		{flow_server_http.Starter(), nil},

		{ws_routes.Starter(), nil},
	}

	joiner, err := starter.Run(starters, cfgCommon, cfgWorkspace, os.Args[1:], label)
	if err != nil {
		l.Fatal(err)
	}
	defer joiner.CloseAll()

	// scheduling importer task

	//dataOp, ok := joiner.Interface(flow.DataInterfaceKey).(data.Operator)
	//if !ok {
	//	l.Fatalf("no data.Operator with key %s", flow.DataInterfaceKey)
	//}

	// TODO!!!
	//url := "http://localhost:" + cfgEnvs["gatherer_port"] + "/gatherer/v1/export"
	//
	//task, err := importer_tasks.NewCopyTask(url, dataOp)
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
	//err = schOp.Run(taskID, time.Minute, true)
	//if err != nil {
	//	l.Fatalf("can't schOp.Run(%s, time.Hour, true): %s", taskID, err)
	//}

	// http_server

	srvOp, ok := joiner.Interface(server_http.InterfaceKey).(server_http.Operator)
	if !ok {
		l.Fatalf("no server_http.Operator with key %s", server_http.InterfaceKey)
	}

	err = srvOp.Start()
	if err != nil {
		l.Error(err)
	}
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
