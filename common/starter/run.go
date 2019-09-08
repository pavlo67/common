package starter

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/pkg/errors"

	"github.com/pavlo67/workshop/common"
	"github.com/pavlo67/workshop/common/config"
	"github.com/pavlo67/workshop/common/joiner"
	"github.com/pavlo67/workshop/common/logger"
)

func StartComponent(c Starter, conf *config.Config, args []string, joinerOp joiner.Operator) error {
	l := logger.Get()

	l.Info("checking component: ", c.Name())

	startOptions := c.CorrectedOptions(ReadOptions(args))

	info, err := c.Init(conf, startOptions)
	for _, i := range info {
		log.Println(i)
	}
	if err != nil {
		return fmt.Errorf("error calling .Init() for component (%s): %s", c.Name(), err)
	}

	err = c.Run(joinerOp)
	if err != nil {
		return fmt.Errorf("error calling .Prepare() for component (%s): %s", c.Name(), err)
	}

	return nil
}

func ReadOptions(args []string) common.Info {
	// TODO!!!

	return nil
}

func Run(starters []Starter, cfg *config.Config, args []string, label string) (joiner.Operator, error) {

	if cfg == nil {
		return nil, errors.New("no config data for starter.Prepare()")
	}

	config.SetActual(cfg)

	l := cfg.Logger

	joinerOp := joiner.New()
	for _, c := range starters {
		err := StartComponent(c, cfg, args, joinerOp)
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
