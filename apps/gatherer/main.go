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

	"github.com/pavlo67/workshop/components/data"
	"github.com/pavlo67/workshop/components/data/data_sqlite"
	"github.com/pavlo67/workshop/components/data/data_tagged"
	"github.com/pavlo67/workshop/constructions/dataflow"
	"github.com/pavlo67/workshop/constructions/dataflow/flow_cleaner/flow_cleaner_sqlite"
	"github.com/pavlo67/workshop/constructions/dataflow/flow_server_http_handler"
	"github.com/pavlo67/workshop/constructions/dataimporter/importer_tasks"
	"github.com/pavlo67/workshop/constructions/taskscheduler"
	"github.com/pavlo67/workshop/constructions/taskscheduler/scheduler_timeout"

	"github.com/pavlo67/workshop/apps/gatherer/gatherer_routes"
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

	var cfgSQLite config.Access
	err = cfgGatherer.Value("sqlite", &cfgSQLite)
	if err != nil {
		l.Fatal(err)
	}

	// running starters

	// TODO!!! vary it
	const flowTable = dataflow.CollectionDefault

	label := "GATHERER/SQLITE CLI BUILD"

	starters := []starter.Starter{
		{control.Starter(), nil},

		{data_sqlite.Starter(), common.Map{"table": flowTable, "interface_key": dataflow.DataInterfaceKey, "no_tagger": true}},
		{data_tagged.Starter(), common.Map{"data_key": dataflow.DataInterfaceKey, "interface_key": dataflow.InterfaceKey, "no_tagger": true}},
		{flow_server_http_handler.Starter(), nil},

		{auth_ecdsa.Starter(), nil},
		{server_http_jschmhr.Starter(), common.Map{"port": port}},
		{gatherer_routes.Starter(), nil},

		{flow_cleaner_sqlite.Starter(), common.Map{"table": flowTable}},
		{scheduler_timeout.Starter(), nil},
		{importer_tasks.Starter(), nil},
	}

	joiner, err := starter.Run(starters, cfgCommon, cfgGatherer, os.Args[1:], label)
	if err != nil {
		l.Fatal(err)
	}
	defer joiner.CloseAll()

	// scheduling importer task

	dataOp, ok := joiner.Interface(dataflow.DataInterfaceKey).(data.Operator)
	if !ok {
		l.Fatalf("no data.Operator with key %s", dataflow.DataInterfaceKey)
	}

	task, err := importer_tasks.NewLoader(dataOp)
	if err != nil {
		l.Fatal(err)
	}

	schOp, ok := joiner.Interface(taskscheduler.InterfaceKey).(taskscheduler.Operator)
	if !ok {
		l.Fatalf("no scheduler.Operator with key %s", taskscheduler.InterfaceKey)
	}

	taskID, err := schOp.Init(task)
	if err != nil {
		l.Fatalf("can't schOp.Init(%#v): %s", task, err)
	}

	err = schOp.Run(taskID, time.Hour, true)
	if err != nil {
		l.Fatalf("can't schOp.Run(%s, time.Hour, false): %s", taskID, err)
	}

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
