package serialization

import (
	"bytes"
	"fmt"
	"os"
	"regexp"

	"github.com/pavlo67/common/common/errors"
	"github.com/pavlo67/common/common/filelib"
)

const onSave = "on serialization.Save()"

func Save(data interface{}, marshaler Marshaler, filename string) error {
	dataBytes, err := marshaler.Marshal(data)
	if err != nil {
		return fmt.Errorf("saving %#v got: %s / "+onSave, data, err)
	}

	if err := os.WriteFile(filename, dataBytes, 0644); err != nil {
		return fmt.Errorf("saving %s got: %s / "+onSave, filename, err)
	}

	return nil
}

const onRead = "on serialization.Read()"

func Read(filename string, marshaler Marshaler, data interface{}) error {
	dataBytes, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("reading %s got: %s / "+onRead, filename, err)
	}

	if err := marshaler.Unmarshal(dataBytes, data); err != nil {
		return fmt.Errorf("unmarshaling %s got: %s / "+onRead, dataBytes, err)
	}

	return nil
}

const onSavePart = "on serialization.SavePart()"

func SavePart(data interface{}, marshaler Marshaler, filename string) error {

	logData, err := marshaler.Marshal(data)
	if err != nil {
		return errors.Wrap(err, onSavePart)
	} else if err = filelib.AppendFile(filename, append(logData, '\n')); err != nil {
		return errors.Wrap(err, onSavePart)
	}

	return nil
}

const onReadPart = "on serialization.ReadPart()"

func ReadPart(filename string, n int, marshaler Marshaler, data interface{}) error {
	dataBytes, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("reading %s got: %s / "+onReadPart, filename, err)
	}
	lines := bytes.Split(dataBytes, []byte{'\n'})
	if n < 0 || n >= len(lines) {
		return fmt.Errorf("wrong n (%d) to get from %d lines / "+onReadPart, n, len(lines))
	}

	if err := marshaler.Unmarshal(lines[n], data); err != nil {
		return fmt.Errorf("unmarshaling %s got: %s / "+onReadPart, dataBytes, err)
	}

	return nil
}

var reEmptyLine = regexp.MustCompile(`^\s*$`)
