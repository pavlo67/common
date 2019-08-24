package main

import (
	"os"

	"fmt"
	"log"

	"github.com/pavlo67/workshop/applications/flow"
	"github.com/pavlo67/workshop/applications/flow/flow_sqlite"
	"github.com/pavlo67/workshop/basis/config"
	"github.com/pavlo67/workshop/basis/common/filelib"
	"github.com/pavlo67/workshop/basis/logger"
	"github.com/pavlo67/workshop/basis/starter"
	"github.com/pavlo67/workshop/basis/instruments/importer/importer_rss"
	"github.com/pavlo67/workshop/basis/server/server_http"
)

func main() {
	err := logger.Init(logger.Config{LogLevel: logger.DebugLevel})
	if err != nil {
		os.Stderr.WriteString(fmt.Sprintf("can't logger.Run(logger.Config{LogLevel: logger.DebugLevel}): %s", err))
		os.Exit(1)
	}
	l := logger.Get()

	cfgPath := filelib.CurrentPath() + "../../../../environments/test.json5"
	cfg, err := config.Get(cfgPath, l)
	if err != nil {
		l.Fatalf("can't config.Get(%s): %s", cfgPath, err)
	}

	starters := []starter.Starter{
		{Operator: flow_sqlite.Starter()},
	}

	joiner, err := starter.Run(starters, cfg, os.Args[1:], "!!!")
	if err != nil {
		l.Fatal(err)

	}
	defer joiner.CloseAll()

	urls := []string{"https://rss.unian.net/site/news_ukr.rss"}

	impOp := &importer_rss.RSS{}

	adminOp, ok := joiner.Interface(flow.InterfaceKey).(flow.Administrator)
	if !ok {
		log.Fatalf("no server_http.Operator with key %s", server_http.InterfaceKey)

	}

	numAll, numProcessed, numNew, errs := flow.Load(urls, impOp, adminOp, l)

	l.Infof("numAll = %d, numProcessed = %d, numNew = %d, errs = %s", numAll, numProcessed, numNew, errs)
}
