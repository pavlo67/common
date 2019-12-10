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
	"github.com/pavlo67/workshop/common/libraries/encodelib"
	"github.com/pavlo67/workshop/common/libraries/filelib"
	"github.com/pavlo67/workshop/common/logger"
	"github.com/pavlo67/workshop/common/server/server_http"
	"github.com/pavlo67/workshop/common/server/server_http/server_http_jschmhr"
	"github.com/pavlo67/workshop/common/starter"

	"github.com/pavlo67/workshop/components/data/data_sqlite"
	"github.com/pavlo67/workshop/components/data/data_tagged"
	"github.com/pavlo67/workshop/components/data/data_tagged/data_tagged_server_http"
	"github.com/pavlo67/workshop/components/flow"
	"github.com/pavlo67/workshop/components/flow/flow_tagged/flow_tagged_server_http"
	"github.com/pavlo67/workshop/components/tagger/tagger_sqlite"

	"github.com/pavlo67/workshop/apps/workspace/ws_routes"
)

var (
	BuildDate    = "unknown"
	BuildRelease = "unknown"
	BuildCommit  = "unknown"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	var versionOnly bool
	flag.BoolVar(&versionOnly, "version", false, "show build vars only")
	flag.Parse()
	if versionOnly {
		log.Printf("builded: %s, revision: %s, commit: %s\n", BuildDate, BuildRelease, BuildCommit)
		return
	}

	currentPath := filelib.CurrentPath()

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

	l, err := logger.Init(logger.Config{})
	if err != nil {
		log.Fatal(err)
	}

	configEnv, ok := os.LookupEnv("ENV")
	if !ok {
		configEnv = "local"
	}
	configPath := currentPath + "../../environments/" + configEnv + ".yaml"

	cfg, err := config.Get(configPath, encodelib.MarshalerYAML)
	if err != nil {
		l.Fatal(err)
	}

	var cfgEnvs map[string]string
	err = cfg.Value("envs", &cfgEnvs)
	if err != nil {
		l.Fatal(err)
	}

	// TODO: synchronize with manifest.json
	//portStr := os.Getenv("workspace_port")
	//port, err := strconv.Atoi(portStr)
	//if err != nil {
	//	l.Fatalf("can't read port: '%s'", portStr)
	//}

	starters := []starter.Starter{
		{control.Starter(), nil},
		{auth_ecdsa.Starter(), nil},
		{server_http_jschmhr.Starter(), common.Map{"port": cfgEnvs["workspace_port"]}},

		{tagger_sqlite.Starter(), nil},

		{data_sqlite.Starter(), nil},
		{data_tagged.Starter(), nil},
		{data_tagged_server_http.Starter(), nil},

		{data_sqlite.Starter(), common.Map{"interface_key": flow.InterfaceKey, "table": flow.CollectionDefault}},
		{data_tagged.Starter(), common.Map{"interface_key": flow.TaggedInterfaceKey, "data_key": flow.InterfaceKey}},
		{flow_tagged_server_http.Starter(), nil},

		{ws_routes.Starter(), nil},
	}

	label := "WORKSPACE REST BUILD"

	joiner, err := starter.Run(starters, cfg, os.Args[1:], label)
	if err != nil {
		l.Fatal(err)
	}
	defer joiner.CloseAll()

	srvOp, ok := joiner.Interface(server_http.InterfaceKey).(server_http.Operator)
	if !ok {
		l.Fatalf("no server_http.Operator with key %s", server_http.InterfaceKey)
	}

	err = srvOp.Start()
	if err != nil {
		l.Error(err)
	}
}
