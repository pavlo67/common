package main

import (
	"flag"
	"log"

	"github.com/pavlo67/common/apps/demo/demo_settings"
	"github.com/pavlo67/common/common/apps"
	"github.com/pavlo67/common/common/starter"
)

var BuildDate, BuildTag, BuildCommit string
var versionOnly bool

func main() {
	log.Printf("builded: %s, tag: %s, commit: %s\n", BuildDate, BuildTag, BuildCommit)
	flag.BoolVar(&versionOnly, "v", false, "show build vars only")
	flag.Parse()

	if versionOnly {
		return
	}

	envPath, cfgService, l := apps.Prepare("_environments/")
	label := "DEMO/REST BUILD"
	joinerOp, err := starter.Run(demo_settings.Components(envPath, true, false), cfgService, label, l)
	if err != nil {
		l.Fatal(err)
	}
	defer joinerOp.CloseAll()

	demo_settings.WG.Wait()
}
