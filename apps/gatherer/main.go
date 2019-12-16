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
	"github.com/pavlo67/workshop/common/scheduler"
	"github.com/pavlo67/workshop/common/scheduler/scheduler_timeout"
	"github.com/pavlo67/workshop/common/serializer"
	"github.com/pavlo67/workshop/common/server/server_http"
	"github.com/pavlo67/workshop/common/server/server_http/server_http_jschmhr"
	"github.com/pavlo67/workshop/common/starter"

	"github.com/pavlo67/workshop/components/data"
	"github.com/pavlo67/workshop/components/data/data_sqlite"
	"github.com/pavlo67/workshop/components/data/data_tagged"
	"github.com/pavlo67/workshop/components/flow"
	"github.com/pavlo67/workshop/components/flow/flow_server_http"

	"github.com/pavlo67/workshop/apps/gatherer/gatherer_routes"
	"github.com/pavlo67/workshop/components/flow/flow_cleaner/flow_cleaner_sqlite"
	"github.com/pavlo67/workshop/components/importer/importer_tasks"
)

var (
	BuildDate   = "unknown"
	BuildTag    = "unknown"
	BuildCommit = "unknown"
)

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

	// common config

	configCommonPath := currentPath + "../../environments/common." + configEnv + ".yaml"
	cfgCommon, err := config.Get(configCommonPath, serializer.MarshalerYAML)
	if err != nil {
		l.Fatal(err)
	}
	var cfgEnvs map[string]string
	err = cfgCommon.Value("envs", &cfgEnvs)
	if err != nil {
		l.Fatal(err)
	}

	// gatherer config

	configGathererPath := currentPath + "../../environments/gatherer." + configEnv + ".yaml"
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
	const flowTable = flow.CollectionDefault

	label := "GATHERER/SQLITE CLI BUILD"

	starters := []starter.Starter{
		{control.Starter(), nil},

		{data_sqlite.Starter(), common.Map{"table": flowTable, "interface_key": flow.InterfaceKey, "no_tagger": true}},
		{data_tagged.Starter(), common.Map{"data_key": flow.InterfaceKey, "interface_key": flow.TaggedInterfaceKey, "no_tagger": true}},
		{flow_server_http.Starter(), nil},

		{auth_ecdsa.Starter(), nil},
		{server_http_jschmhr.Starter(), common.Map{"port": cfgEnvs["gatherer_port"]}},
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

	dataOp, ok := joiner.Interface(flow.InterfaceKey).(data.Operator)
	if !ok {
		l.Fatalf("no data.Operator with key %s", flow.InterfaceKey)
	}

	task, err := importer_tasks.NewLoader(dataOp)
	if err != nil {
		l.Fatal(err)
	}

	schOp, ok := joiner.Interface(scheduler.InterfaceKey).(scheduler.Operator)
	if !ok {
		l.Fatalf("no scheduler.Operator with key %s", scheduler.InterfaceKey)
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
