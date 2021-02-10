package main

import (
	"log"

	"github.com/pavlo67/common/common/logger"
	"github.com/pavlo67/common/common/logger/logger_zap"
)

func main() {

	l, err := logger_zap.New(logger.Config{
		//OutputPaths:      []string{"stdout", "aaa.log"},
		//ErrorOutputPaths: []string{"stderr", "aaa.log"},
	})
	if err != nil {
		log.Fatal(err)
	}

	l.Info("!!!!")

	l.Errorf("!!!!%d", 111111)

	l.Error(5345)

}
