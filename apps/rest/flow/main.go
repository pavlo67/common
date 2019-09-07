package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/pavlo67/workshop/apps/rest/flow/flow_routes"
	"github.com/pavlo67/workshop/apps/rest/flow/flow_routes/v1"
	"github.com/pavlo67/workshop/basis/auth/auth_ecdsa"
	"github.com/pavlo67/workshop/basis/common/filelib"
	"github.com/pavlo67/workshop/basis/config"
	"github.com/pavlo67/workshop/basis/logger"
	"github.com/pavlo67/workshop/basis/server/server_http"
	"github.com/pavlo67/workshop/basis/server/server_http/server_http_jschmhr"
	"github.com/pavlo67/workshop/basis/starter"
	"github.com/pavlo67/workshop/components/data/data_sqlite"
)

var (
	BuildDate    = "unknown"
	BuildRelease = "unknown"
	BuildCommit  = "unknown"
)

func main() {
	start := time.Now()
	rand.Seed(start.UnixNano())

	var versionOnly bool
	flag.BoolVar(&versionOnly, "version", false, "show build vars only")
	flag.Parse()
	if versionOnly {
		fmt.Printf("builded: %s, revision: %s, commit: %s\n", BuildDate, BuildRelease, BuildCommit)
		return
	}

	configPath := filelib.CurrentPath() + "../../../environments"
	configEnv, ok := os.LookupEnv("ENV")
	if !ok {
		configEnv = "local"
	}

	err := logger.Init(logger.Config{LogLevel: logger.DebugLevel})
	if err != nil {
		fmt.Printf("can't logger.Init, error: %v\n", err)
		os.Exit(1)
	}

	l := logger.Get()
	if l == nil {
		fmt.Printf("no logger!")
		os.Exit(1)
	}

	cfg, err := config.Get(configPath, configEnv, l)
	if err != nil {
		l.Fatalf("can't config.Get(%s): %s", configPath, err)
	}

	if err != nil {
		fmt.Printf("can't load config, error: %v\n", err)
		os.Exit(1)
	}
	if cfg == nil {
		fmt.Printf("can't load config, no data!")
		os.Exit(1)
	}

	// flag.Parse()

	//var ownerID basis.ID      // from flags
	//ownerPublKey := encrlib.Base58Decode([]byte(ownerID))
	//var managersJSON, signature []byte   // from some external channel
	//if !encrlib.ECDSAVerify(ownerPublKey, managersJSON, signature) {
	//
	//}

	// !!! kostyl
	l.Info(ep_flow.ToInit)

	starters := []starter.Starter{
		{auth_ecdsa.Starter(), nil},
		{server_http_jschmhr.Starter(), nil},
		{data_sqlite.Starter(), nil},
		{flow_routes.Starter(), nil},
	}

	label := "DATA REST BUILD"

	joiner, err := starter.Run(starters, cfg, os.Args[1:], label)
	if err != nil {
		l.Fatal(err)

	}
	defer joiner.CloseAll()

	srvOp, ok := joiner.Interface(server_http.InterfaceKey).(server_http.Operator)
	if !ok {
		log.Fatalf("no server_http.Operator with key %s", server_http.InterfaceKey)
	}

	srvOp.HandleFiles("/flow/api-docs/*filepath", filelib.CurrentPath()+"../_api-docs/", nil)
	srvOp.HandleFiles("/flow/swagger/*filepath", filelib.CurrentPath()+"api-docs/", nil)

	srvOp.Start()

	//go srvOp.Start()
	//
	//c := make(chan os.Signal, 1)
	//signal.Notify(c, os.Interrupt)
	//signal := <-c
	//l.Info("\nGot signal:", signal)

}
