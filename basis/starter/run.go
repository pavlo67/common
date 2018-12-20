package starter

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/pkg/errors"

	"github.com/pavlo67/punctum/basis/config"
	"github.com/pavlo67/punctum/basis/logger"
	"github.com/pavlo67/punctum/basis/program"
)

func StartComponent(conf *config.PunctumConfig, c Starter, joiner program.Joiner) error {
	l := logger.Get()

	l.Info("  -------------   check component: ", c.Name(), "   ---------------")

	err := c.Prepare(conf, c.Params)
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

	err = c.Init(joiner)
	if err != nil {
		return fmt.Errorf("error calling Init() for component (%s): %s", c.Name(), err)
	}

	return nil
}

func Run(conf *config.PunctumConfig, starters []Starter, label string, runKeys []program.InterfaceKey) (program.Joiner, error) {
	l := logger.Get()

	if conf == nil {
		return nil, errors.New("no config data for starter.Run()")
	}

	joiner := program.NewJoiner()

	for _, c := range starters {
		err := StartComponent(conf, c, joiner)
		if err != nil {
			return joiner, err
		}
	}

	//for _, runKey := range runKeys {
	//	if runner, ok := joiner.Interface(runKey).(Runner); ok {
	//		err := runner.Run()
	//		if err != nil {
	//			return joiner, errors.Wrapf(err, "can't start .Runner for key %s", runKey)
	//		}
	//	} else {
	//		return joiner, errors.Errorf("no .Runner interface for key %s", runKey)
	//	}
	//}

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

	return joiner, nil
}
