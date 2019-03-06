package starter

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/pkg/errors"

	"github.com/pavlo67/punctum/basis"
	"github.com/pavlo67/punctum/starter/config"
	"github.com/pavlo67/punctum/starter/joiner"
	"github.com/pavlo67/punctum/starter/logger"
)

func StartComponent(c Starter, conf *config.PunctumConfig, runtimeOptions basis.Options, joinerOp joiner.Operator) error {
	l := logger.Get()

	l.Info("checking component: ", c.Name())

	err := c.Prepare(conf, c.Options, runtimeOptions)
	if err != nil {
		return fmt.Errorf("error calling .Prepare() for component (%s): %s", c.Name(), err)
	}

	info, err := c.Check()
	if err != nil {
		for _, i := range info {
			log.Println(i)
		}
		return fmt.Errorf("error calling Check() for component (%s): %s", c.Name(), err)
	}

	err = c.Init(joinerOp)
	if err != nil {
		return fmt.Errorf("error calling zapInit() for component (%s): %s", c.Name(), err)
	}

	return nil
}

func ReadOptions(args []string) basis.Options {
	return nil
}

func Run(starters []Starter, conf *config.PunctumConfig, args []string, label string) (joiner.Operator, error) {
	l := logger.Get()

	if conf == nil {
		return nil, errors.New("no config data for starter.Run()")
	}

	runtimeOptions := ReadOptions(args)

	joinerOp := joiner.New()
	for _, c := range starters {
		err := StartComponent(c, conf, runtimeOptions, joinerOp)
		if err != nil {
			return joinerOp, err
		}
	}

	env, ok := os.LookupEnv("ENV")
	if !ok || strings.TrimSpace(env) == "" {
		env = "(default)"
	}
	l.Info(label + "; environment = " + env)

	// wait-runner
	//if waitForInterrupt {
	//	c := make(chan os.Signal, 1)
	//	signal.Notify(c, os.Interrupt)
	//	signal := <-c
	//	fmt.Println("\nGot signal:", signal)
	//}

	return joinerOp, nil
}
