package storage_server_http

import (
	"net/http"

	"github.com/pkg/errors"

	"github.com/pavlo67/workshop/common/auth"
	"github.com/pavlo67/workshop/common/joiner"
	"github.com/pavlo67/workshop/common/server"
	"github.com/pavlo67/workshop/common/server/server_http"
)

const interfaceKeyParamName = "type_key"
const tagLabelParamName = "tag"

var CountTaggedEndpoint = server_http.Endpoint{Method: "GET", QueryParams: []string{interfaceKeyParamName}, WorkerHTTP: CountTagged}

func CountTagged(_ *auth.User, _ server_http.Params, req *http.Request) (server.Response, error) {
	var interfaceKeyPtr *joiner.InterfaceKey

	if key := req.URL.Query().Get(interfaceKeyParamName); key != "" {
		interfaceKey := joiner.InterfaceKey(key)
		interfaceKeyPtr = &interfaceKey
	}

	counter, err := dataTaggedOp.CountTagged(interfaceKeyPtr, nil)
	if err != nil {
		return server.ResponseRESTError(http.StatusInternalServerError, errors.Errorf("ERROR on GET storage/...CountTagged (%#v): %s", req.URL.Query(), err))
	}

	return server.ResponseRESTOk(counter)
}

var ListWithTagEndpoint = server_http.Endpoint{Method: "GET", QueryParams: []string{"tag"}, WorkerHTTP: ListWithTag}

func ListWithTag(user *auth.User, _ server_http.Params, req *http.Request) (server.Response, error) {

	tagLabel := req.URL.Query().Get(tagLabelParamName)
	items, err := dataTaggedOp.ListWithTag(nil, tagLabel, nil)

	if err != nil {
		return server.ResponseRESTError(http.StatusInternalServerError, errors.Errorf("ERROR on GET storage/...ListWithTag (%#v): %s", req.URL.Query(), err))
	}

	return server.ResponseRESTOk(items)
}
