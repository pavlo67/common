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
	"github.com/pavlo67/workshop/components/tags/tags_sqlite"

	"github.com/pavlo67/workshop/constructions/dataflow"
	"github.com/pavlo67/workshop/constructions/dataflow/flow_cleaner/flow_cleaner_sqlite"
	"github.com/pavlo67/workshop/constructions/dataflow/flow_server_http"
	"github.com/pavlo67/workshop/constructions/dataimporter/importer_tasks"
	"github.com/pavlo67/workshop/constructions/datastorage"
	"github.com/pavlo67/workshop/constructions/datastorage/storage_server_http"
	"github.com/pavlo67/workshop/constructions/taskscheduler"
	"github.com/pavlo67/workshop/constructions/taskscheduler/scheduler_timeout"

	"github.com/pavlo67/workshop/apps/workspace/ws_routes"
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

	// storage config

	configWorkspacePath := currentPath + "../../environments/workspace." + configEnv + ".yaml"
	cfgWorkspace, err := config.Get(configWorkspacePath, serializer.MarshalerYAML)
	if err != nil {
		l.Fatal(err)
	}

	// running starters

	const flowTable = dataflow.CollectionDefault

	label := "WORKSPACE REST BUILD"

	starters := []starter.Starter{
		{control.Starter(), nil},
		{auth_ecdsa.Starter(), nil},
		{server_http_jschmhr.Starter(), common.Map{"port": cfgEnvs["workspace_port"]}},

		{tags_sqlite.Starter(), nil},

		{data_sqlite.Starter(), common.Map{"interface_key": datastorage.InterfaceKey, "table": datastorage.CollectionDefault}},
		{data_tagged.Starter(), common.Map{"interface_key": datastorage.TaggedInterfaceKey, "data_key": datastorage.InterfaceKey}},
		{storage_server_http.Starter(), nil},

		{data_sqlite.Starter(), common.Map{"interface_key": dataflow.InterfaceKey, "table": flowTable}},
		{data_tagged.Starter(), common.Map{"interface_key": dataflow.TaggedInterfaceKey, "data_key": dataflow.InterfaceKey}},
		{flow_server_http.Starter(), nil},

		{ws_routes.Starter(), nil},

		{flow_cleaner_sqlite.Starter(), common.Map{"table": flowTable}},
		{scheduler_timeout.Starter(), nil},
		{importer_tasks.Starter(), nil},
	}

	joiner, err := starter.Run(starters, cfgCommon, cfgWorkspace, os.Args[1:], label)
	if err != nil {
		l.Fatal(err)
	}
	defer joiner.CloseAll()

	// scheduling importer task

	dataOp, ok := joiner.Interface(dataflow.InterfaceKey).(data.Operator)
	if !ok {
		l.Fatalf("no data.Operator with key %s", dataflow.InterfaceKey)
	}

	// TODO!!!
	url := "http://localhost:" + cfgEnvs["gatherer_port"] + "/gatherer/v1/export"

	task, err := importer_tasks.NewCopyTask(url, dataOp)
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

	err = schOp.Run(taskID, time.Minute, true)
	if err != nil {
		l.Fatalf("can't schOp.Run(%s, time.Hour, true): %s", taskID, err)
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
