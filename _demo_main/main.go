package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/pavlo67/punctum/basis/filelib"
	"github.com/pavlo67/punctum/server/server_http"
	"github.com/pavlo67/punctum/starter"
	"github.com/pavlo67/punctum/starter/config"
	"github.com/pavlo67/punctum/starter/logger"

	"github.com/pavlo67/punctum/_demo_main/demo_starters"
)

func main() {
	err := logger.Init(logger.Config{LogLevel: logger.DebugLevel})
	if err != nil {
		os.Stderr.WriteString(fmt.Sprintf("can't logger.Init(logger.Config{LogLevel: logger.DebugLevel}): %s", err))
		os.Exit(1)
	}
	l := logger.Get()

	cfgPath := filelib.CurrentPath() + "../cfg.json5"
	conf, err := config.Get(cfgPath, l)
	if err != nil {
		l.Fatalf("can't config.Get(%s): %s", cfgPath, err)
	}

	// flag.Parse()

	starters, label := demo_starters.Starters()
	joiner, err := starter.Run(starters, conf, os.Args[1:], label)
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
