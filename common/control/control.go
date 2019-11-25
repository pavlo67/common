package control

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/pavlo67/workshop/common"
	"github.com/pavlo67/workshop/common/config"
	"github.com/pavlo67/workshop/common/joiner"
	"github.com/pavlo67/workshop/common/logger"
	"github.com/pavlo67/workshop/common/starter"
)

func Starter() starter.Operator {
	return &controlStarter{}
}

var l logger.Operator
var _ starter.Operator = &controlStarter{}

type controlStarter struct{}

func (ws *controlStarter) Name() string {
	return logger.GetCallInfo().PackageName
}

func (ws *controlStarter) Init(_ *config.Config, lCommon logger.Operator, options common.Options) ([]common.Options, error) {
	l = lCommon
	return nil, nil
}

func (ws *controlStarter) Setup() error {
	return nil
}

var signalChan = make(chan os.Signal, 1000)

func (ws *controlStarter) Run(joinerOp joiner.Operator) error {

	signal.Notify(signalChan, os.Interrupt)
	signal.Notify(signalChan, syscall.SIGTERM)
	signal.Notify(signalChan, syscall.SIGQUIT)
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
