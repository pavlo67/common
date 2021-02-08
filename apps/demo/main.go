package main

import (
	"log"

	"github.com/pavlo67/common/common/apps"
	"github.com/pavlo67/common/common/logger"
	"github.com/pavlo67/common/common/starter"

	"github.com/pavlo67/common/apps/demo/demo_api"
)

var (
	BuildDate   = "unknown"
	BuildTag    = "unknown"
	BuildCommit = "unknown"
)

const serviceName = "demo"

func main() {
	versionOnly, envPath, cfgService, l := apps.Prepare(BuildDate, BuildTag, BuildCommit, serviceName, apps.AppsSubpathDefault)
	if versionOnly {
		return
	}

	l, err := logger.Init(logger.Config{})
	if err != nil {
		log.Fatal(err)
	}

	// running starters

	label := "DEMO/PG/REST BUILD"
	joinerOp, err := starter.Run(demo_api.Components(envPath, true, false), cfgService, label, logger.Get())
	if err != nil {
		l.Fatal(err)
	}
	defer joinerOp.CloseAll()

	demo_api.WG.Wait()
}
