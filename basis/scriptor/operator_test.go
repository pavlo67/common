package scriptor

import (
	"testing"

	"log"

	"github.com/stretchr/testify/require"
)

func TestRead(t *testing.T) {
	var script = "1 + 2 + (3 - 5) * 2"

	elem, _, err := Read(script, "", nil)
	require.NoError(t, err)
	require.NotNil(t, elem)
	log.Printf("%#v", *elem)

}
