package main

import (
	"fmt"
	"log"
	"os"

	"github.com/pavlo67/workshop/basis/auth/auth_ecdsa"
	"github.com/pavlo67/workshop/basis/common/filelib"
	"github.com/pavlo67/workshop/basis/config"
	"github.com/pavlo67/workshop/basis/logger"
	"github.com/pavlo67/workshop/basis/server/server_http"
	"github.com/pavlo67/workshop/basis/server/server_http/server_http_jschmhr"
	"github.com/pavlo67/workshop/basis/starter"
)

func Starters() ([]starter.Starter, string) {
	//paramsServerStatic := basis.Info{
	//	"static_path": filelib.CurrentPath() + "../demo_server_http/static/",
	//}

	var starters []starter.Starter

	starters = append(starters, starter.Starter{auth_ecdsa.Starter(), nil})
	starters = append(starters, starter.Starter{server_http_jschmhr.Starter(), nil})
	// starters = append(starters, starter.Starter{rector_server.Starter(), nil})

	return starters, "DATA FLOW BUILD"
}

func main() {
	err := logger.Init(logger.Config{LogLevel: logger.DebugLevel})
	if err != nil {
		os.Stderr.WriteString(fmt.Sprintf("can't logger.Run(logger.Config{LogLevel: logger.DebugLevel}): %s", err))
		os.Exit(1)
	}

	l := logger.Get()

	cfgPath := filelib.CurrentPath() + "../../../environments/cfg.json5"
	conf, err := config.Get(cfgPath, l)
	if err != nil {
		l.Fatalf("can't config.Get(%s): %s", cfgPath, err)
	}

	// flag.Parse()

	//var ownerID basis.ID      // from flags
	//ownerPublKey := encrlib.Base58Decode([]byte(ownerID))
	//var managersJSON, signature []byte   // from some external channel
	//if !encrlib.ECDSAVerify(ownerPublKey, managersJSON, signature) {
	//
	//}

	starters, label := Starters()
	joiner, err := starter.Run(starters, conf, os.Args[1:], label)
	if err != nil {
		l.Fatal(err)

	}
	defer joiner.CloseAll()

	srvOp, ok := joiner.Interface(server_http.InterfaceKey).(server_http.Operator)
	if !ok {
		log.Fatalf("no server_http.Operator with key %s", server_http.InterfaceKey)
	}

	srvOp.Start()

	//go srvOp.Start()
	//
	//c := make(chan os.Signal, 1)
	//signal.Notify(c, os.Interrupt)
	//signal := <-c
	//l.Info("\nGot signal:", signal)

}
