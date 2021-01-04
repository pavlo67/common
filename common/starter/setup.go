package starter

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/pkg/errors"

	"github.com/pavlo67/workshop/common/config"
	"github.com/pavlo67/workshop/common/logger"
)

func Setup(cfg *config.Config, starters []Starter, components []string, label string) error {
	if cfg == nil {
		return errors.New("no config data for starter.Setup()")
	}

	toSetup := map[string]bool{}

	for _, c := range components {
		toSetup[c] = true
	}

	if len(toSetup) < 1 {
		log.Println("to setup: all")
	} else {
		log.Println("to setup: ", toSetup)
	}

	env, ok := os.LookupEnv("ENV")
	if !ok || strings.TrimSpace(env) == "" {
		env = "(default)"
	}

	log.Print(label + "; environment = " + env + ": ok?")
	bufio.NewReader(os.Stdin).ReadString('\n')

	for _, c := range starters {
		if len(toSetup) >= 1 && !toSetup[c.Name()] {
			continue
		}

		log.Println("  ---------- setup component: ", c.Name(), "   -----------")

		_, err := c.Init(cfg, logger.Get(), c.Options)
		if err != nil {
			return fmt.Errorf("error calling .Init() for component (%s): %s", c.Name(), err)
		}

		err = c.Setup()
		if err != nil {
			return err
		}
	}

	return nil
}
