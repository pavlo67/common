package storage_server_http

import (
	"net/http"

	"github.com/pavlo67/workshop/common/crud"

	"github.com/pkg/errors"

	"github.com/pavlo67/workshop/common/auth"
	"github.com/pavlo67/workshop/common/joiner"
	"github.com/pavlo67/workshop/common/server"
	"github.com/pavlo67/workshop/common/server/server_http"
)

const interfaceKeyParamName = "key"
const tagLabelParamName = "tag"

var listTagsEndpoint = server_http.Endpoint{Method: "GET", QueryParams: []string{interfaceKeyParamName}, WorkerHTTP: CountTags}

func CountTags(user *auth.User, _ server_http.Params, req *http.Request) (server.Response, error) {
	var interfaceKeyPtr *joiner.InterfaceKey
	if key := req.URL.Query().Get(interfaceKeyParamName); key != "" {
		interfaceKey := joiner.InterfaceKey(key)
		interfaceKeyPtr = &interfaceKey
	}

	counter, err := dataTaggedOp.CountTags(interfaceKeyPtr, &crud.GetOptions{ActorKey: user.KeyYet()})
	if err != nil {
		return server.ResponseRESTError(http.StatusInternalServerError, errors.Errorf("ERROR on GET storage/...CountTags (%#v): %s", req.URL.Query(), err))
	}

	return server.ResponseRESTOk(counter)
}

var listTaggedEndpoint = server_http.Endpoint{Method: "GET", QueryParams: []string{interfaceKeyParamName, tagLabelParamName}, WorkerHTTP: ListTagged}

func ListTagged(user *auth.User, _ server_http.Params, req *http.Request) (server.Response, error) {
	//var interfaceKeyPtr *joiner.InterfaceKey
	//if key := req.URL.Query().Get(interfaceKeyParamName); key != "" {
	//	interfaceKey := joiner.InterfaceKey(key)
	//	interfaceKeyPtr = &interfaceKey
	//}

	tagLabel := req.URL.Query().Get(tagLabelParamName)

	items, err := dataTaggedOp.ListTagged(tagLabel, nil, &crud.GetOptions{ActorKey: user.KeyYet()})

	if err != nil {
		return server.ResponseRESTError(http.StatusInternalServerError, errors.Errorf("ERROR on GET storage/...ListTagged (%#v): %s", req.URL.Query(), err))
	}

	return server.ResponseRESTOk(items)
}

var exportEndpoint = server_http.Endpoint{Method: "GET", QueryParams: []string{"after"}, WorkerHTTP: Export}

func Export(user *auth.User, _ server_http.Params, req *http.Request) (server.Response, error) {
	// TODO!!! use "after"

	crudData, err := exporterOp.Export(nil, "", &crud.GetOptions{ActorKey: user.KeyYet()})

	if err != nil {
		return server.ResponseRESTError(http.StatusInternalServerError, errors.Errorf("ERROR on GET storage/...Export: ", err))
	}

	return server.ResponseRESTOk(crudData)
}
