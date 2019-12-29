package identity

import (
	"testing"

	"log"

	"github.com/stretchr/testify/require"
)

type TestCase struct {
	Key          Key
	PathExpected string
	KeyExpected  Key
}

func TestIdentity(t *testing.T) {
	testCases := []TestCase{
		{"", "", ""},
		{"abc", "", "abc"},
		{"abc/", "", "abc"},
		{"/abc", "abc", "/abc"},
		{"dumaj.org.ua/abc/123#11", "abc", "dumaj.org.ua/abc/123#11"},
		{"dumaj.org.ua//123", "", "dumaj.org.ua//123"},
		{"dumaj.org.ua/a/b/c/d/123", "a/b/c/d", "dumaj.org.ua/a/b/c/d/123"},
	}

	for _, tc := range testCases {
		item := tc.Key.Identity()

		log.Printf("%#v", item)
		require.Equal(t, tc.PathExpected, item.Path)
		require.Equal(t, tc.KeyExpected, item.Key())
	}
}
