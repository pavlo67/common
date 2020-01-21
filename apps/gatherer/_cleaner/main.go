package main

import (
	"flag"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/pavlo67/workshop/common"
	"github.com/pavlo67/workshop/common/config"
	"github.com/pavlo67/workshop/common/libraries/filelib"
	"github.com/pavlo67/workshop/common/logger"
	"github.com/pavlo67/workshop/common/serializer"
	"github.com/pavlo67/workshop/common/starter"
	"github.com/pavlo67/workshop/components/flow"
	"github.com/pavlo67/workshop/constructions/dataflow/flow_cleaner"
	"github.com/pavlo67/workshop/constructions/dataflow/flow_cleaner/flow_cleaner_sqlite"
)

var (
	BuildDate   = "unknown"
	BuildTag    = "unknown"
	BuildCommit = "unknown"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	var versionOnly bool
	flag.BoolVar(&versionOnly, "version", false, "show build vars only")
	flag.Parse()

	log.Printf("builded: %s, tag: %s, commit: %s\n", BuildDate, BuildTag, BuildCommit)

	if versionOnly {
		return
	}

	// logger

	l, err := logger.Init(logger.Config{})
	if err != nil {
		log.Fatal(err)
	}

	// getting config environments

	currentPath := filelib.CurrentPath()
	configEnv, ok := os.LookupEnv("ENV")
	if !ok {
		configEnv = "local"
	}

	// common config

	//configCommonPath := currentPath + "../../../environments/common." + configEnv + ".yaml"
	//cfgCommon, err := config.Get(configCommonPath, encodelib.MarshalerYAML)
	//if err != nil {
	//	l.Fatal(err)
	//}
	//var cfgEnvs map[string]string
	//err = cfgCommon.Value("envs", &cfgEnvs)
	//if err != nil {
	//	l.Fatal(err)
	//}

	// cleaner config

	serviceEnv, ok := os.LookupEnv("SERVICE")

	configGathererPath := currentPath + "../../../environments/" + serviceEnv + "." + configEnv + ".yaml"
	cfgGatherer, err := config.Get(configGathererPath, serializer.MarshalerYAML)
	if err != nil {
		l.Fatal(err)
	}

	var cfgSQLite config.Access
	err = cfgGatherer.Value("sqlite", &cfgSQLite)
	if err != nil {
		l.Fatal(err)
	}

	// running starters

	label := "FLOW CLEANER/SQLITE CLI BUILD"

	starters := []starter.Starter{
		{flow_cleaner_sqlite.Starter(), common.Map{"table": flow.CollectionDefault}},
	}

	joiner, err := starter.Run(starters, nil, cfgGatherer, os.Args[1:], label)
	if err != nil {
		l.Fatal(err)
	}
	defer joiner.CloseAll()

	fcOp, _ := joiner.Interface(flow_cleaner.InterfaceKey).(flow_cleaner.Operator)
	if fcOp == nil {
		l.Fatalf("no flow_cleaner.ActorKey with key %s", flow_cleaner.InterfaceKey)
	}

	err = fcOp.Clean(2850)
	if err != nil {
		l.Fatal(err)
	}
}
