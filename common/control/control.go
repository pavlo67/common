package control

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/pavlo67/workshop/common/logger"
)

var (
	signalChan = make(chan os.Signal, 1000)
)

var l logger.Operator

func Init(lInit logger.Operator) {
	l = lInit

	signal.Notify(signalChan, os.Interrupt)
	signal.Notify(signalChan, syscall.SIGTERM)
	signal.Notify(signalChan, syscall.SIGQUIT)
	signal.Notify(signalChan, syscall.SIGPIPE)

	go processSignal()
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
