package confidence_scenarios

import (
	"encoding/json"
	"log"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/pavlo67/workshop/basis/auth"
	"github.com/pavlo67/workshop/basis/common"
)

func TestAuth(t *testing.T) {
	host := "http://localhost:3333"
	path := "/confidence/v1/auth"

	url := host + path

	body, err := json.Marshal([]auth.Creds{{Type: auth.CredsNickname, Value: "aaa"}, {Type: auth.CredsPassword, Value: "bbb"}})
	require.NoError(t, err)

	dataMap, err := common.RequestJSON("POST", url, body, nil)
	require.NoError(t, err)

	log.Printf("%#v", dataMap)
}
