package main

import (
	"github.com/pavlo67/common/common/apps"
	"github.com/pavlo67/common/common/starter"

	"github.com/pavlo67/common/apps/demo/demo_settings"
)

var (
	BuildDate   = ""
	BuildTag    = ""
	BuildCommit = ""
)

func main() {
	versionOnly, envPath, cfgService, l := apps.Prepare(BuildDate, BuildTag, BuildCommit, "_environments/")
	if versionOnly {
		return
	}

	label := "DEMO/REST BUILD"
	joinerOp, err := starter.Run(demo_settings.Components(envPath, true, false), cfgService, label, l)
	if err != nil {
		l.Fatal(err)
	}
	defer joinerOp.CloseAll()

	demo_settings.WG.Wait()
}
