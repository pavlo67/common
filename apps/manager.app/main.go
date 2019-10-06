package main

import (
	"flag"
	"fmt"
	"os"

	"log"

	"github.com/pavlo67/workshop/common/config"
	"github.com/pavlo67/workshop/common/control"
	"github.com/pavlo67/workshop/common/libs/filelib"
	"github.com/pavlo67/workshop/common/logger"
	"github.com/pavlo67/workshop/common/manager"
)

func main() {
	var path string
	flag.StringVar(&path, "path", "", "application manifest path")
	flag.Parse()

	configPath := filelib.CurrentPath() + "../../environments"
	configEnv, ok := os.LookupEnv("ENV")
	if !ok {
		configEnv = "local"
	}

	cfg, err := config.Get(configPath, configEnv)
	if err != nil {
		log.Fatalf("can't config.Get(%s): %s", configPath, err)
	}
	if cfg == nil {
		log.Fatalf("can't load config, no data!")
	}

	l, err := logger.Init(logger.Config{LogLevel: logger.DebugLevel}, cfg)
	if err != nil {
		fmt.Printf("can't logger.Init, error: %v\n", err)
		os.Exit(1)
	}
	if l == nil {
		fmt.Printf("no logger!")
		os.Exit(1)
	}

	control.Init(l)

	app, err := manager.Init(path, cfg, nil)
	if err != nil {
		l.Fatal(err)
	}

	app.Start()
}
