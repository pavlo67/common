package starter

import (
	"fmt"
	"os"
	"strings"

	"github.com/pavlo67/common/common/joiner/joiner_runtime"

	"github.com/pkg/errors"

	"github.com/pavlo67/common/common/config"
	"github.com/pavlo67/common/common/joiner"
	"github.com/pavlo67/common/common/logger"
)

func Run(starters []Component, serviceConfig *config.Config, label string, l logger.Operator) (joiner.Operator, error) {
	for _, c := range starters {
		name := c.Name()
		if key, ok := c.Options.String("interface_key"); ok {
			name += " / " + key
		}

		l.Info("preparing component: ", name)
		starterOptions := c.CorrectedOptions(nil)
		starterConfig := c.Config
		if starterConfig == nil {
			starterConfig = serviceConfig
		}

		if err := c.Prepare(starterConfig, starterOptions); err != nil {
			return nil, fmt.Errorf("error calling .PrepareApp() for component (%s): %s", name, err)
		}
	}

	joinerOp := joiner_runtime.New(nil, l)
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
			return nil, fmt.Errorf("error calling .Run() for component (%s): %s", name, err)
		}
	}

	env, ok := os.LookupEnv("ENV")
	if !ok || strings.TrimSpace(env) == "" {
		env = "(default)"
	}
	l.Info(label + "; environment = " + env)

	return joinerOp, nil
}
