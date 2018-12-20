package config

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

type EndpointPathWithParamsTestCase struct {
	Path                   string
	Params                 []string
	PathWithParamsExpected string
}

func TestEndpointPathWithParams(t *testing.T) {
	testCases := []EndpointPathWithParamsTestCase{
		{"/:aaaa/:bb", []string{"a", "b"}, "/a/b"},
		{"/:aaaa/:bb", []string{"a", "b", "c"}, "/a/b"},
		{"/:aaaa/:bb", []string{"a"}, "/a/:bb"},
	}

	for i, tc := range testCases {
		fmt.Println(i)

		ep := Endpoint{ServerPath: tc.Path}
		require.Equal(t, tc.PathWithParamsExpected, ep.Path(tc.Params...))
	}
}
