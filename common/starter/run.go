package starter

import (
	"fmt"
	"os"
	"strings"

	"github.com/pkg/errors"

	"github.com/pavlo67/common/common/config"
	"github.com/pavlo67/common/common/joiner"
	"github.com/pavlo67/common/common/logger"
)

//func StartComponent(c Starter, cfg *config.Config, joinerOp joiner.Operator, l logger.Operator) error {
//	name := c.Name()
//
//	if key, ok := c.Options.String("interface_key"); ok {
//		name += " / " + key
//	}
//
//	l.Info("checking component: ", name)
//
//	startOptions := c.CorrectedOptions(nil)
//
//	if err := c.Prepare(cfg, startOptions); err != nil {
//		return fmt.Errorf("error calling .Prepare() for component (%s): %#v", name, err)
//	}
//
//	if err := c.Run(joinerOp); err != nil {
//		return fmt.Errorf("error calling .Run() for component (%s): %s", name, err)
//	}
//
//	return nil
//}

func Run(starters []Starter, cfg *config.Config, label string, l logger.Operator) (joiner.Operator, error) {
	for _, c := range starters {
		name := c.Name()
		if key, ok := c.Options.String("interface_key"); ok {
			name += " / " + key
		}

		l.Info("preparing component: ", name)
		startOptions := c.CorrectedOptions(nil)
		if err := c.Prepare(cfg, startOptions); err != nil {
			return nil, fmt.Errorf("error calling .Prepare() for component (%s): %#v", name, err)
		}
	}

	joinerOp := joiner.New(nil, l)
	if err := joinerOp.Join(l, logger.InterfaceKey); err != nil {
		return nil, errors.Errorf("can't join logger with key %s: %s", logger.InterfaceKey, err)
	}

	for _, c := range starters {
		name := c.Name()
		if key, ok := c.Options.String("interface_key"); ok {
			name += " / " + key
		}

		l.Info("running component: ", name)
		if err := c.Run(joinerOp); err != nil {
			return nil, fmt.Errorf("error calling .Run() for component (%s): %#v", name, err)
		}
	}

	env, ok := os.LookupEnv("ENV")
	if !ok || strings.TrimSpace(env) == "" {
		env = "(default)"
	}
	l.Info(label + "; environment = " + env)

	return joinerOp, nil
}
