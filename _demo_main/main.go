package main

import (
	"flag"
	"log"
	"os"
	"os/signal"

	"go.uber.org/zap"

	"github.com/pavlo67/punctum/basis/filelib"
	"github.com/pavlo67/punctum/basis/program"
	"github.com/pavlo67/punctum/basis/starter"
	"github.com/pavlo67/punctum/server_http"

	"github.com/pavlo67/punctum/_demo_main/demo_config"
)

func main() {
	conf, l, err := program.Init(filelib.CurrentPath()+"../cfg.json5", zap.DebugLevel)
	if err != nil {
		log.Fatal(err)
	}

	flag.Parse()

	starters, label := demo_config.Starters()

	joiner, err := starter.Run(conf, starters, label, nil)
	if err != nil {
		l.Fatal(err)

	}

	defer joiner.CloseAll()

	srvOp, ok := joiner.Interface(server_http.InterfaceKey).(server_http.Operator)
	if !ok {
		log.Fatalf("no server_http.Operator with key %s", server_http.InterfaceKey)

	}
	go srvOp.Start()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal := <-c
	l.Info("\nGot signal:", signal)
}
