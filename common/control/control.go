package control

import (
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

var signalChan = make(chan os.Signal, 1000)

func (ws *controlStarter) Run(_ *config.Config, _ common.Map, joinerOp joiner.Operator, l_ logger.Operator) error {
	l = l_

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
