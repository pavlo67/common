package main

import (
	"flag"
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"

	"github.com/pavlo67/workshop/common/auth/auth_ecdsa"
	"github.com/pavlo67/workshop/common/config"
	"github.com/pavlo67/workshop/common/control"
	"github.com/pavlo67/workshop/common/manager"
	"github.com/pavlo67/workshop/common/server/server_http"
	"github.com/pavlo67/workshop/common/server/server_http/server_http_jschmhr"
	"github.com/pavlo67/workshop/common/starter"
	"github.com/pavlo67/workshop/components/data/data_sqlite"
	"github.com/pavlo67/workshop/libraries/filelib"

	"github.com/pavlo67/workshop/apps/gatherer/flow_routes/starter"
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

	manifest, err := manager.ReadManifest(currentPath)
	if err != nil {
		log.Fatal(err)
	}
	//if manifest == nil {
	//	log.Fatalf("can't load manifest, no data!")
	//}

	for _, key := range manifest.Requested {
		if os.Getenv(key) == "" {
			log.Fatalf("no environment value for key '%s'", key)
		}
	}

	configPath := currentPath + "../../environments"
	configEnv, ok := os.LookupEnv("ENV")
	if !ok {
		configEnv = "local"
	}

	cfg, l, err := config.Get(configPath, configEnv)
	if err != nil {
		log.Fatalf("can't config.Get(%s): %s", configPath, err)
	}
	//if cfg == nil {
	//	log.Fatalf("can't load config, no data!")
	//}

	control.Init(l)

	starters := []starter.Starter{
		{auth_ecdsa.Starter(), nil},
		{server_http_jschmhr.Starter(), nil},
		{data_sqlite.Starter(), nil},
		{flow_starter.Starter(), nil},
	}

	label := "DATA REST BUILD"

	joiner, err := starter.Run(starters, cfg, os.Args[1:], label)
	if err != nil {
		l.Fatal(err)
	}
	defer joiner.CloseAll()

	// TODO: synchronize with manifest.json
	portStr := os.Getenv("flow_port")
	port, err := strconv.Atoi(portStr)
	if err != nil {
		l.Fatalf("can't read port: '%s'", portStr)
	}

	srvOp, ok := joiner.Interface(server_http.InterfaceKey).(server_http.Operator)
	if !ok {
		log.Fatalf("no server_http.Operator with key %s", server_http.InterfaceKey)
	}

	srvOp.Start(port)

}

//srvOp.HandleFiles("/flow/api-docs/*filepath", filelib.CurrentPath()+"./api-docs/", nil)
//flow_v1.Init()
//for _, ep := range flow_routes.Endpoints {
//	srvOp.HandleEndpoint(ep)
//}
