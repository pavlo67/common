package router_news

import (
	"strconv"
	"time"

	"github.com/pkg/errors"

	"github.com/pavlo67/punctum/basis"
	"github.com/pavlo67/punctum/crud"
	"github.com/pavlo67/punctum/server"
	"github.com/pavlo67/punctum/server/controller"
)

type URLs []string

var endpoints = map[string]controller.Endpoint{
	"clean":    {Method: "GET", Path: "clean", Worker: clean},
	"load":     {Method: "GET", Path: "load", Worker: load},
	"loadPost": {Method: "POST", Path: "load", Worker: load, DataItem: URLs{}},
	"list":     {Method: "GET", Path: "list", Worker: list},
	"listPost": {Method: "POST", Path: "list", Worker: list, DataItem: URLs{}},
}

const daysForCleanDefault = 7

const onClean = "on news_router.clean()"

func clean(endpoint controller.Endpoint, params basis.Params, _ basis.Options, _ interface{}) (*server.Response, error) {
	var err error

	daysStr := params.ByNum(0)
	days := daysForCleanDefault
	if daysStr != "" {
		days, err = strconv.Atoi(daysStr)
		if err != nil {
			return nil, errors.Wrap(err, onClean)
		}
	}

	err = newsOp.DeleteList(&crud.ReadOptions{
		Selector: basis.Lt(crud.TimeField, time.Now().Add(-time.Hour*24*time.Duration(days)).Format(time.RFC3339)),
	})

	return nil, err
}

func load(endpoint controller.Endpoint, params basis.Params, options basis.Options, data interface{}) (*server.Response, error) {
	var urls URLs
	if endpoint.Method == "POST" {
		var ok bool
		urls, ok = data.(URLs)
		if !ok {
			return nil, errors.New("wrong data type")
		}
	} else {
		urls = options.Strings("url")
	}

	num, numNew, errs := Load(urls, newsOp)
	responseData := server.Response{
		Data: map[string]int{"num": num, "num_new": numNew},
	}

	return &responseData, errs.Err()
}

func list(endpoint controller.Endpoint, params basis.Params, options basis.Options, data interface{}) (*server.Response, error) {
	var urls URLs
	if endpoint.Method == "POST" {
		var ok bool
		urls, ok = data.(URLs)
		if !ok {
			return nil, errors.New("wrong data type")
		}
	} else {
		urls = options.Strings("url")
	}

	days, _ := strconv.Atoi(options.StringDefault("days", "0"))
	now := time.Now().UTC()

	selector := basis.And(
		basis.InStr(string(crud.URLField), urls),
		basis.Unary(basis.Ge(crud.TimeField, time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC).Add(-time.Hour*time.Duration(days)).Format(time.RFC3339))),
	)
	news, _, err := newsOp.ReadList(&crud.ReadOptions{
		Selector: selector,
	})

	return &server.Response{Data: news}, err
}

//func stat(endpoint router.Endpoint, params basis.Options, options basis.Options, data interface{}) (*server.Response, error) {
//	var urls URLs
//	if endpoint.Method == "POST" {
//		var ok bool
//		urls, ok = data.(URLs)
//		if !ok {
//			return nil, errors.New("wrong data type")
//		}
//	} else {
//		urls = options.Strings("url")
//	}
//
//	responseData := server.Response{
//		Status: 0,
//		Data:   nil,
//	}
//
//	return &responseData, nil
//}
