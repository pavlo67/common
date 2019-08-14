package main

import (
	"os"

	"fmt"
	"log"

	"github.com/pavlo67/constructor/applications/flow"
	"github.com/pavlo67/constructor/applications/flow/flow_sqlite"
	"github.com/pavlo67/constructor/components/basis/config"
	"github.com/pavlo67/constructor/components/basis/filelib"
	"github.com/pavlo67/constructor/components/basis/logger"
	"github.com/pavlo67/constructor/components/basis/starter"
	"github.com/pavlo67/constructor/components/processor/importer/importer_rss"
	"github.com/pavlo67/constructor/components/server/server_http"
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

	numAll, numNewAll, errs := flow.Load(urls, impOp, adminOp, l)

	l.Infof("numAll = %d, numNewAll = %d, errs = %^#v", numAll, numNewAll, errs)
}
