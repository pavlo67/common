package main

import (
	"fmt"
	"log"
	"os"

	"github.com/pavlo67/workshop/common/config"
	"github.com/pavlo67/workshop/common/logger"
	"github.com/pavlo67/workshop/common/server/server_http"
	"github.com/pavlo67/workshop/common/starter"
	"github.com/pavlo67/workshop/libraries/filelib"

	"github.com/pavlo67/workshop/components/data"
	"github.com/pavlo67/workshop/components/data/data_importer"
	"github.com/pavlo67/workshop/components/data/data_sqlite"
	"github.com/pavlo67/workshop/components/instruments/importer/importer_rss"
)

func main() {
	err := logger.Init(logger.Config{LogLevel: logger.DebugLevel})
	if err != nil {
		os.Stderr.WriteString(fmt.Sprintf("can't logger.Prepare(logger.Config{LogLevel: logger.DebugLevel}): %s", err))
		os.Exit(1)
	}
	l := logger.Get()

	cfgPath := filelib.CurrentPath() + "../../../../environments/test.json5"
	cfg, err := config.Get(cfgPath, l)
	if err != nil {
		l.Fatalf("can't config.Get(%s): %s", cfgPath, err)
	}

	starters := []starter.Starter{
		{Operator: data_sqlite.Starter()},
	}

	joiner, err := starter.Run(starters, cfg, os.Args[1:], "!!!")
	if err != nil {
		l.Fatal(err)

	}
	defer joiner.CloseAll()

	urls := []string{"https://rss.unian.net/site/news_ukr.rss"}

	impOp := &importer_rss.RSS{}

	dataOp, ok := joiner.Interface(data.InterfaceKey).(data.Operator)
	if !ok {
		log.Fatalf("no server_http.Operator with key %s", server_http.InterfaceKey)

	}

	numAll, numProcessed, numNew, errs := data_importer.Load(urls, impOp, dataOp, l)

	l.Infof("numAll = %d, numProcessed = %d, numNew = %d, errs = %s", numAll, numProcessed, numNew, errs)
}
