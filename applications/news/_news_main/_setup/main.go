package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/pavlo67/constructor/components/basis/filelib"
	"github.com/pavlo67/constructor/components/basis/joiner"
	"github.com/pavlo67/constructor/components/basis/starter"

	"github.com/pavlo67/constructor/processor/_starter_process_rsss_rss/process_rss_config"
)

var setup = flag.Bool("setup", false, "recreate structures for the selected (or all if no) component")

func main() {
	_, conf, err := joiner.Init(filelib.CurrentPath()+"../../cfg.json5", false)
	if err != nil {
		log.Fatal(err)
	}

	flag.Parse()

	if !*setup {
		fmt.Println("no action selected...")
		return
	}

	starters, label := process_rss_config.Starters()
	err = starter.Setup(conf, starters, flag.Args(), label+" / to setup the component(s)")
	if err != nil {
		log.Print("ERROR: ", err)
	}

}
