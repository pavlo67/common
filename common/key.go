package common

import (
	"strconv"
	"strings"
)

type ID string

func (key ID) Normalize() ID {
	return ID(strings.TrimSpace(string(key)))
}

func (key ID) Uint64() uint64 {
	keyUint64, _ := strconv.ParseUint(string(key), 10, 64)

	return keyUint64
}
