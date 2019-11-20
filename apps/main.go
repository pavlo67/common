package main

import (
	"flag"
	"log"
	"os"

	"github.com/pavlo67/workshop/common/config"
	"github.com/pavlo67/workshop/common/control"
	"github.com/pavlo67/workshop/common/manager"
	"github.com/pavlo67/workshop/libraries/filelib"
)

func main() {
	var path string
	flag.StringVar(&path, "path", "./", "application manifest path")
	flag.Parse()

	configPath := filelib.CurrentPath() + "../environments"
	configEnv, ok := os.LookupEnv("ENV")
	if !ok {
		configEnv = "local"
	}

	cfg, l, err := config.Get(configPath, configEnv)
	if err != nil {
		log.Fatalf("can't config.Get(%s): %s", configPath, err)
	}
	if cfg == nil {
		log.Fatal("can't load config, no data!")
	}
	if l == nil {
		log.Fatal("no logger!")
	}

	control.Init(l)

	app, err := manager.Init(path, cfg, l, nil)
	if err != nil {
		l.Fatal(err)
	}

	app.Start()
}
