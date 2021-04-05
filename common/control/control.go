package control

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/pavlo67/common/common"

	"github.com/pavlo67/common/common/config"
	"github.com/pavlo67/common/common/joiner"
	"github.com/pavlo67/common/common/logger"
	"github.com/pavlo67/common/common/starter"
)

// DEPRECATED
func Starter() starter.Operator {
	return &controlStarter{}
}

var l logger.Operator
var _ starter.Operator = &controlStarter{}

type controlStarter struct{}

func (ws *controlStarter) Name() string {
	return logger.GetCallInfo().PackageName
}

func (ws *controlStarter) Prepare(_ *config.Config, options common.Map) error {
	return nil
}

var signalChan = make(chan os.Signal, 1000)

func (ws *controlStarter) Run(joinerOp joiner.Operator) error {
	if l, _ = joinerOp.Interface(logger.InterfaceKey).(logger.Operator); l == nil {
		return fmt.Errorf("no logger.Operator with key %s", logger.InterfaceKey)
	}

	// signal.Notify(signalChan, os.Interrupt)
	// signal.Notify(signalChan, syscall.SIGTERM)
	// signal.Notify(signalChan, syscall.SIGQUIT)

	signal.Notify(signalChan, syscall.SIGPIPE)

	go processSignal()

	return nil
}

func processSignal() {
	for {
		sig := <-signalChan

		if sig == syscall.SIGPIPE {
			signal.Reset(sig)
			l.Warnf("got & ignored signal %s", sig)
		} else {
			l.Warnf("got signal %s", sig)
		}
	}
}
