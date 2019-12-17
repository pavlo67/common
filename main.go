package main

import (
	"log"

	"github.com/pkg/errors"
)

// executable cod for various probes only

func main() {

	err := errors.New("original error")

	log.Print(errors.Wrap(err, "wrapped text"))
}
