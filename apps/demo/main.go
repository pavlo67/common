package main

import (
	"flag"
	"log"
	"os"

	"github.com/pavlo67/workshop/common/config"
	"github.com/pavlo67/workshop/common/logger"
	"github.com/pavlo67/workshop/common/serializer"
	"github.com/pavlo67/workshop/common/starter"

	"github.com/pavlo67/workshop/apps/demo/demo_api"
)

var (
	BuildDate   = "unknown"
	BuildTag    = "unknown"
	BuildCommit = "unknown"
)

const serviceNameDefault = "demo"
const appsSubpathDefault = "apps/"

func main() {
	//rand.Seed(time.Now().UnixNano())

	var versionOnly bool
	var serviceName, appsSubpath string
	flag.BoolVar(&versionOnly, "version_only", false, "show build vars only")
	flag.StringVar(&serviceName, "service", serviceNameDefault, "service name")
	flag.StringVar(&appsSubpath, "apps_subpath", appsSubpathDefault, "subpath to /apps directory")
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

	configEnv, ok := os.LookupEnv("ENV")
	if !ok {
		configEnv = "local"
	}

	cwd, err := os.Getwd()
	if err != nil {
		l.Fatal("can't os.Getwd(): ", err)
	}
	cwd += "/"
	l.Info("CWD: ", cwd)

	// get config

	envPath := cwd + appsSubpath + "_environments/"
	cfgServicePath := envPath + serviceName + "." + configEnv + ".yaml"
	cfgService, err := config.Get(cfgServicePath, serviceName, serializer.MarshalerYAML)
	if err != nil {
		l.Fatal(err)
	}

	// running starters

	label := "DEMO/PG/REST BUILD"

	// TODO: rename production ENV value

	joinerOp, err := starter.Run(demo_api.Components(envPath, true, false), cfgService, os.Args[1:], label)
	if err != nil {
		l.Fatal(err)
	}
	defer joinerOp.CloseAll()

	demo_api.WG.Wait()
}
