package starter

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/pkg/errors"

	"github.com/pavlo67/associatio/basis"
	"github.com/pavlo67/associatio/starter/config"
	"github.com/pavlo67/associatio/starter/joiner"
	"github.com/pavlo67/associatio/starter/logger"
)

func StartComponent(c Starter, conf *config.Config, runtimeOptions basis.Info, joinerOp joiner.Operator) error {
	l := logger.Get()

	l.Info("checking component: ", c.Name())

	info, err := c.Init(conf, c.Options, runtimeOptions)
	for _, i := range info {
		log.Println(i)
	}
	if err != nil {
		return fmt.Errorf("error calling .Init() for component (%s): %s", c.Name(), err)
	}

	err = c.Run(joinerOp)
	if err != nil {
		return fmt.Errorf("error calling zapInit() for component (%s): %s", c.Name(), err)
	}

	return nil
}

func ReadOptions(args []string) basis.Info {
	return nil
}

func Run(starters []Starter, conf *config.Config, args []string, label string) (joiner.Operator, error) {
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
