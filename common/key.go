package common

import "strconv"

type Key string

func (key Key) Uint64() uint64 {
	keyUint64, _ := strconv.ParseUint(string(key), 10, 64)

	return keyUint64
}
