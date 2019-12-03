package workspace_server_http

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/pkg/errors"

	"github.com/pavlo67/workshop/common"
	"github.com/pavlo67/workshop/common/auth"
	"github.com/pavlo67/workshop/common/server"
	"github.com/pavlo67/workshop/common/server/server_http"

	"github.com/pavlo67/workshop/components/data"
)

// TODO!!! add parameters info into responces

// Save --------------------------------------------------------------------------------------

var SaveEndpoint = server_http.Endpoint{Method: "POST", ParamNames: []string{"type"}, WorkerHTTP: Save}

func Save(user *auth.User, params server_http.Params, req *http.Request) (server.Response, error) {
	var item data.Item

	itemType := params["type"]
	switch itemType {
	case "test":
		item.Details = &data.Test{}
	default:
		return server.ResponseRESTError(http.StatusBadRequest, errors.Errorf("ERROR on POST workspace/...Save: wrong item type: %s", itemType))
	}

	var itemJSON []byte

	itemJSON, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return server.ResponseRESTError(http.StatusBadRequest, errors.Errorf("ERROR on POST workspace/...Save: reading body: %s", err))
	}

	err = json.Unmarshal(itemJSON, &item)
	if err != nil {
		return server.ResponseRESTError(http.StatusBadRequest, errors.Errorf("ERROR on POST workspace/...Save: can't json.Unmarshal(%s): %s", itemJSON, err))
	}

	ids, err := workspaceOp.Save([]data.Item{item}, nil)
	if err != nil {
		return server.ResponseRESTError(http.StatusInternalServerError, errors.Errorf("ERROR on POST workspace/...Save: %s", err))
	}

	if len(ids) != 1 {
		return server.ResponseRESTError(http.StatusInternalServerError, errors.Errorf("ERROR on POST workspace/...Save: returned wrong ids (%#v)", ids))
	}

	return server.ResponseRESTOk(map[string]interface{}{"id": ids[0]})
}

// Read --------------------------------------------------------------------------------------

var ReadEndpoint = server_http.Endpoint{Method: "GET", ParamNames: []string{"type"}, WorkerHTTP: Read}

func Read(user *auth.User, params server_http.Params, req *http.Request) (server.Response, error) {
	var details interface{}

	itemType := params["type"]
	switch itemType {
	case "test":
		details = &data.Test{}
	default:
		return server.ResponseRESTError(http.StatusBadRequest, errors.Errorf("ERROR on GET workspace/...Read: wrong item type: %s", itemType))
	}

	queryParams := req.URL.Query()
	id := common.ID(queryParams.Get("id"))

	item, err := workspaceOp.Read(id, nil)
	if err == common.ErrNotFound {
		return server.ResponseRESTError(http.StatusNotFound, errors.Errorf("ERROR on GET workspace/...Read: not found item with id = %s", id))
	} else if err != nil {
		return server.ResponseRESTError(http.StatusInternalServerError, errors.Errorf("ERROR on GET workspace/...Read: ", err))
	}

	item.Details = details
	err = workspaceOp.Details(item, item.Details)
	if err != nil {
		return server.ResponseRESTError(http.StatusInternalServerError, errors.Errorf("ERROR on GET workspace/...Read: ", err))
	}

	return server.ResponseRESTOk(item)
}

// l.Infof("item: %#v", item)
// l.Infof("details!!!: %#v", item.Details)

// List --------------------------------------------------------------------------------------

var ListEndpoint = server_http.Endpoint{Method: "GET", ParamNames: nil, WorkerHTTP: List}

func List(user *auth.User, _ server_http.Params, req *http.Request) (server.Response, error) {
	items, err := workspaceOp.List(nil, nil)

	l.Infof("%#v", items)

	if err != nil {
		return server.ResponseRESTError(http.StatusInternalServerError, errors.Errorf("ERROR on GET workspace/...List: ", err))
	}

	return server.ResponseRESTOk(items)
}

// Read --------------------------------------------------------------------------------------

var RemoveEndpoint = server_http.Endpoint{Method: "DELETE", ParamNames: nil, WorkerHTTP: Remove}

func Remove(user *auth.User, params server_http.Params, req *http.Request) (server.Response, error) {
	queryParams := req.URL.Query()
	id := common.ID(queryParams.Get("id"))

	err := workspaceOp.Remove(id, nil)
	if err != nil {
		return server.ResponseRESTError(http.StatusInternalServerError, errors.Errorf("ERROR on DELETE workspace/...Remove: ", err))
	}

	return server.ResponseRESTOk(map[string]interface{}{"id": id})
}
