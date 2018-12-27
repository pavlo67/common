package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/pavlo67/punctum/basis/config"
	"github.com/pavlo67/punctum/basis/libs/filelib"
	"github.com/pavlo67/punctum/basis/logger"
	"github.com/pavlo67/punctum/basis/starter"
	"github.com/pavlo67/punctum/server_http"

	"github.com/pavlo67/punctum/_demo_main/demo_config"
)

func main() {
	err := logger.Init(logger.Config{LogLevel: logger.DebugLevel})
	if err != nil {
		os.Stderr.WriteString(fmt.Sprintf("can't logger.zapInit(logger.Config{LogLevel: zap.DebugLevel}): %s", err))
		os.Exit(1)
	}
	l := logger.Get()

	cfgPath := filelib.CurrentPath() + "../cfg.json5"
	conf, err := config.Get(cfgPath, l)
	if err != nil {
		l.Fatalf("can't config.zapGet(%s): %s", cfgPath, err)
	}

	// flag.Parse()

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
