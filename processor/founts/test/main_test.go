package test

import (
	"log"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	if err := os.Setenv("ENV", "test"); err != nil {
		log.Fatalln("No test environment!!!")
	}
	os.Exit(m.Run())
}
