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
	IsNull       bool
}

func TestIdentity(t *testing.T) {
	testCases := []TestCase{
		{"", "", "", true},
		{"abc", "", "abc", false},
		{"abc/", "", "abc", false},
		{"/abc", "abc", "/abc", false},
		{"dumaj.org.ua/abc/123#11", "abc/123#11", "dumaj.org.ua/abc/123#11", false},
		{"dumaj.org.ua//123", "123", "dumaj.org.ua/123", false},
		{"dumaj.org.ua/a/b/c/d/123", "a/b/c/d/123", "dumaj.org.ua/a/b/c/d/123", false},
		{"dumaj.org.ua/a/b/c/d/123####abcd", "a/b/c/d/123", "dumaj.org.ua/a/b/c/d/123##abcd", false},
	}

	for _, tc := range testCases {
		item := tc.Key.Identity()

		if tc.IsNull {
			require.Nil(t, item)
		} else {
			require.NotNil(t, item)
			log.Printf("%s --> %#v", tc.Key, item)
			require.Equal(t, tc.PathExpected, item.Path)
			require.Equal(t, tc.KeyExpected, item.Key())
		}
	}
}
