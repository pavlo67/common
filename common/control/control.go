package control

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/pavlo67/workshop/common/logger"
)

var (
	signalChan = make(chan os.Signal, 1)
	exitChan   = make(chan int, 1)
)

var l logger.Operator

func Init(lInit logger.Operator) {
	l = lInit

	signal.Notify(signalChan, os.Interrupt)
	signal.Notify(signalChan, syscall.SIGTERM)
	signal.Notify(signalChan, syscall.SIGQUIT)
	signal.Notify(signalChan, syscall.SIGPIPE)
	//signal.Ignore(syscall.SIGPIPE) //this is totally unhandled ignoring of signal

	go processExit()
	go processSignal()
}

func processSignal() {
	sig := <-signalChan
	l.Errorf("got signal, value = \"%s\"\n", sig.String())
	if sig == syscall.SIGPIPE {
		signal.Reset(syscall.SIGPIPE)
		return
	}
	exitChan <- 1
}

func processExit() {
	select {
	case code := <-exitChan:
		//ohlc.Close()
		//sheduler.Stop()
		l.Infof("service exiting, code %d\n", code)

		os.Exit(code)
	}
}
