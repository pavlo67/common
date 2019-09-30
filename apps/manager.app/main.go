package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/pavlo67/workshop/common/control"
	"github.com/pavlo67/workshop/common/logger"
	"github.com/pavlo67/workshop/common/manager"
)

func main() {
	var path string
	flag.StringVar(&path, "path", "", "application manifest path")
	flag.Parse()

	//configPath := filelib.CurrentPath() + "../../environments"
	//configEnv, ok := os.LookupEnv("ENV")
	//if !ok {
	//	configEnv = "local"
	//}

	err := logger.Init(logger.Config{LogLevel: logger.DebugLevel})
	if err != nil {
		fmt.Printf("can't logger.Init, error: %v\n", err)
		os.Exit(1)
	}

	l := logger.Get()
	if l == nil {
		fmt.Printf("no logger!")
		os.Exit(1)
	}

	control.Init(l)

	app, err := manager.Init(path, l)
	if err != nil {
		l.Fatal(err)
	}

	app.Start()
}
