// +build  linux

package filelib

import (
	"log"
	"os/exec"
	"strings"
	"testing"

	"github.com/pkg/errors"
)

func TestCurrentPath(t *testing.T) {
	path := CurrentPath()
	currentPath_, err := exec.Command("pwd").Output()
	//currentPath_, err := exec.Command("cmd", "/C", "cd",).Output()
	if err != nil {
		log.Fatal(err)
	}
	currentPath := strings.Trim(string(currentPath_), "\t\n\r") + "/"

	if currentPath != path {
		t.Error(errors.Errorf("bad current path: %v vs. %v", path, currentPath))
	}
}
