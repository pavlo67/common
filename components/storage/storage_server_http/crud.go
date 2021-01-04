package storage_server_http

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/pkg/errors"

	"github.com/pavlo67/workshop/common"
	"github.com/pavlo67/workshop/common/auth"
	"github.com/pavlo67/workshop/common/crud"
	"github.com/pavlo67/workshop/common/server"
	"github.com/pavlo67/workshop/common/server/server_http"

	"github.com/pavlo67/workshop/components/data"
)

// TODO!!! add parameters info the responses

// Save --------------------------------------------------------------------------------------

var saveEndpoint = server_http.Endpoint{Method: "POST", WorkerHTTP: Save}

func Save(user *auth.User, params server_http.Params, req *http.Request) (server.Response, error) {
	if user == nil {
		return server.ResponseRESTError(http.StatusUnauthorized, errors.New("ERROR on POST storage/...Save: no user"))
	}

	var item data.Item

	itemJSON, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return server.ResponseRESTError(http.StatusBadRequest, errors.Errorf("ERROR on POST storage/...Save: reading body: %s", err))
	}

	err = json.Unmarshal(itemJSON, &item)
	if err != nil {
		return server.ResponseRESTError(http.StatusBadRequest, errors.Errorf("ERROR on POST storage/...Save: can't json.Unmarshal(%s): %s", string(itemJSON), err))
	}

	id, err := dataOp.Save(item, &crud.SaveOptions{ActorKey: user.Key})
	if err != nil {
		return server.ResponseRESTError(http.StatusInternalServerError, errors.Errorf("ERROR on POST storage/...Save: %s", err))
	}
	if id == "" {
		return server.ResponseRESTError(http.StatusInternalServerError, errors.New("ERROR on POST storage/...Save: no id returned"))
	}

	return server.ResponseRESTOk(map[string]interface{}{"id": id})
}

// Read --------------------------------------------------------------------------------------

var readEndpoint = server_http.Endpoint{Method: "GET", PathParams: []string{"id"}, WorkerHTTP: Read}

func Read(user *auth.User, params server_http.Params, req *http.Request) (server.Response, error) {
	id := common.ID(params["id"])

	item, err := dataOp.Read(id, &crud.GetOptions{ActorKey: user.KeyYet()})
	if err == common.ErrNotFound {
		return server.ResponseRESTError(http.StatusNotFound, errors.Errorf("ERROR on GET storage/...Read: not found item with id = %s", id))
	} else if err != nil {
		return server.ResponseRESTError(http.StatusInternalServerError, errors.Errorf("ERROR on GET storage/...Read: ", err))
	}

	return server.ResponseRESTOk(item)
}

// ListFlow --------------------------------------------------------------------------------------

var recentEndpoint = server_http.Endpoint{Method: "GET", WorkerHTTP: Recent}

func Recent(user *auth.User, _ server_http.Params, req *http.Request) (server.Response, error) {
	items, err := dataOp.List(nil, &crud.GetOptions{OrderBy: []string{data.RecentOrder}, ActorKey: user.KeyYet()})

	if err != nil {
		return server.ResponseRESTError(http.StatusInternalServerError, errors.Errorf("ERROR on GET storage/...ListFlow: ", err))
	}

	return server.ResponseRESTOk(items)
}

// Read --------------------------------------------------------------------------------------

var removeEndpoint = server_http.Endpoint{Method: "DELETE", PathParams: []string{"id"}, WorkerHTTP: Remove}

func Remove(user *auth.User, params server_http.Params, req *http.Request) (server.Response, error) {
	id := common.ID(params["id"])

	err := dataOp.Remove(id, &crud.RemoveOptions{ActorKey: user.KeyYet()})
	if err != nil {
		return server.ResponseRESTError(http.StatusInternalServerError, errors.Errorf("ERROR on DELETE storage/...Remove: ", err))
	}

	return server.ResponseRESTOk(map[string]interface{}{"id": id})
}
