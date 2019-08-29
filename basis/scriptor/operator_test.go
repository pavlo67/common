package scriptor

import (
	"log"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRead(t *testing.T) {
	var script = "aa bbb 0. .1 0.234-11 rwe )*()^&%(^%"

	item, err := Read(script)
	require.NoError(t, err)
	require.NotNil(t, item)

	log.Printf("%#v", *item)
}
