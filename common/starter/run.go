package starter

import (
	"fmt"
	"os"
	"strings"

	"github.com/pavlo67/common/common/config"
	"github.com/pavlo67/common/common/joiner"
	"github.com/pavlo67/common/common/joiner/joiner_runtime"
	"github.com/pavlo67/common/common/logger"
)

func Run(starters []Component, serviceConfig *config.Config, label string, l logger.Operator) (joiner.Operator, error) {
	joinerOp := joiner_runtime.New(nil, l)

	for _, c := range starters {
		name := c.Name()
		if key, ok := c.Options.String("interface_key"); ok {
			name += " / " + key
		}

		l.Info("running component: ", name)

		starterConfig := c.Config
		if starterConfig == nil {
			starterConfig = serviceConfig
		}

		if err := c.Run(starterConfig, c.Options, joinerOp, l); err != nil {
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
