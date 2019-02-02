package starter

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/pkg/errors"

	"github.com/pavlo67/punctum/starter/config"
)

func Setup(conf *config.PunctumConfig, starters []Starter, components []string, label string) error {
	if conf == nil {
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

		err := c.Prepare(conf, c.Options)
		if err != nil {
			return fmt.Errorf("error calling .Prepare() for component (%s): %s", c.Name(), err)
		}

		err = c.Setup()
		if err != nil {
			return err
		}
	}

	return nil
}
