package serialization

import (
	"bytes"
	"encoding/json"
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

const onReadAllPartsJSON = "on serialization.ReadAllPartsJSON()"

func ReadAllPartsJSON(filename string, data interface{}) error {
	dataBytesRaw, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("reading %s got: %s / "+onReadAllPartsJSON, filename, err)
	}
	var lines [][]byte
	for _, line := range bytes.Split(dataBytesRaw, []byte{'\n'}) {
		if !reEmptyLine.Match(line) {
			lines = append(lines, line)
		}
	}

	dataBytesLines := bytes.Join(lines, []byte{','})
	dataBytes := append([]byte{'['}, append(dataBytesLines, ']')...)

	if err := json.Unmarshal(dataBytes, data); err != nil {
		return fmt.Errorf("unmarshaling %s got: %s / "+onReadAllPartsJSON, dataBytes, err)
	}

	return nil
}

const onSaveAllPartsJSON = "on serialization.SaveAllPartsJSON()"

func SaveAllPartsJSON[T any](data []T, filename string) error {
	var dataBytes []byte

	for _, item := range data {
		itemBytes, err := json.Marshal(item)
		if err != nil {
			return fmt.Errorf("marshaling %v got: %s / "+onSaveAllPartsJSON, item, err)
		}

		dataBytes = append(dataBytes, itemBytes...)
		dataBytes = append(dataBytes, '\n')
	}

	err := os.WriteFile(filename, dataBytes, 0644)
	if err != nil {
		return errors.Wrap(err, onSaveAllPartsJSON)
	}

	return nil
}
