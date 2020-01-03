package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/pavlo67/workshop/apps/confidence/confidence_routes"
	"github.com/pavlo67/workshop/apps/confidence/confidence_routes/v1/auth"
	"github.com/pavlo67/workshop/common"
	"github.com/pavlo67/workshop/common/auth/auth_ecdsa"
	"github.com/pavlo67/workshop/common/auth/auth_jwt"
	"github.com/pavlo67/workshop/common/auth/auth_stub"
	"github.com/pavlo67/workshop/common/config"
	"github.com/pavlo67/workshop/common/control"
	"github.com/pavlo67/workshop/common/kv/kv_sqlite"
	"github.com/pavlo67/workshop/common/logger"
	"github.com/pavlo67/workshop/common/server/server_http"
	"github.com/pavlo67/workshop/common/server/server_http/server_http_jschmhr"
	"github.com/pavlo67/workshop/common/starter"
	"github.com/pavlo67/workshop/libraries/filelib"
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

	configPath := filelib.CurrentPath() + "../../environments"
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

	control.Init(l)

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

	starters := []starter.Starter{
		{auth_stub.Starter(), common.Map{"interface_key": string(auth_stub.InterfaceKey)}},
		// {auth_users_sqlite.Starter(), common.Map{"interface_key": string(auth_users_sqlite.HandlerKey)}},
		{auth_ecdsa.Starter(), common.Map{"interface_key": string(auth_ecdsa.InterfaceKey)}},
		{auth_jwt.Starter(), common.Map{"interface_key": string(auth_jwt.InterfaceKey)}},
		{kv_sqlite.Starter(), nil},
		{server_http_jschmhr.Starter(), nil},
		{confidence_routes.Starter(), nil},
	}

	label := "CONFIDENCE REST BUILD"

	joiner, err := starter.Run(starters, cfg, os.Args[1:], label)
	if err != nil {
		l.Fatal(err)

	}
	defer joiner.CloseAll()

	srvOp, ok := joiner.Interface(server_http.InterfaceKey).(server_http.Operator)
	if !ok {
		log.Fatalf("no server_http.Operator with key %s", server_http.InterfaceKey)
	}

	srvOp.HandleFiles("/confidence/api-docs/*filepath", filelib.CurrentPath()+"api-docs/", nil)

	v1_auth.Init()
	for _, ep := range confidence_routes.Endpoints {
		srvOp.HandleEndpoint(ep)
	}

	srvOp.Start()

	//go srvOp.Start()
	//
	//c := make(chan os.Signal, 1)
	//signal.Notify(c, os.Interrupt)
	//signal := <-c
	//l.Map("\nGot signal:", signal)

}
