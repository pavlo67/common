package main

import (
	"github.com/pavlo67/common/common/apps"
	"github.com/pavlo67/common/common/starter"

	"github.com/pavlo67/common/apps/demo/demo_api"
)

var (
	BuildDate   = ""
	BuildTag    = ""
	BuildCommit = ""
)

const serviceName = "demo"

func main() {
	versionOnly, envPath, cfgService, l := apps.Prepare(BuildDate, BuildTag, BuildCommit, serviceName, apps.AppsSubpathDefault)
	if versionOnly {
		return
	}

	label := "DEMO/REST BUILD"
	joinerOp, err := starter.Run(demo_api.Components(envPath, true, false), cfgService, label, l)
	if err != nil {
		l.Fatal(err)
	}
	defer joinerOp.CloseAll()

	demo_api.WG.Wait()
}
