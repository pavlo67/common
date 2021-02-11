package main

import (
	"fmt"
	"log"

	"github.com/pavlo67/common/common/errata"
	"github.com/pkg/errors"
)

func main() {
	err := errors.Wrapf(errors.New("eeeeee"), "22222 %s", "111")
	log.Print(err)
	err1 := errata.CommonError(err, "can't init records.Operator")
	log.Print(fmt.Errorf("error calling .Run() for component (%s): %#v", "name", err1))
	log.Print(fmt.Errorf("error calling .Run() for component (%s): %s", "name", err1))

}
