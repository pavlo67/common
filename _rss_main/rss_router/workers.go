package rss_router

import (
	"github.com/pavlo67/punctum/server"
	"github.com/pavlo67/punctum/server/router"
)

var workers = map[string]router.WorkerFunc{
	"remove_old": removeOld,
	"load":       load,
	"read_list":  readList,
	"stat":       stat,
}

var endpoints = map[string]router.Endpoint{
	"remove_old": {Method: "GET", ServerPath: "/remove_old"},
	"load":       {Method: "POST", ServerPath: "/load"},
	"read_list":  {Method: "GET", ServerPath: "/read_list"},
	"stat":       {Method: "GET", ServerPath: "/stat"},
}

func removeOld(params router.Params, data []byte) (server.DataResponse, error) {
	responseData := server.DataResponse{
		Status: 0,
		Data:   nil,
	}

	return responseData, nil
}

func load(params router.Params, data []byte) (server.DataResponse, error) {
	responseData := server.DataResponse{
		Status: 0,
		Data:   nil,
	}

	return responseData, nil
}

func readList(params router.Params, data []byte) (server.DataResponse, error) {
	responseData := server.DataResponse{
		Status: 0,
		Data:   nil,
	}

	return responseData, nil
}

func stat(params router.Params, data []byte) (server.DataResponse, error) {
	responseData := server.DataResponse{
		Status: 0,
		Data:   nil,
	}

	return responseData, nil
}
