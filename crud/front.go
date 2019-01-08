package crud

import (
	"strconv"

	"github.com/pkg/errors"

	"github.com/pavlo67/punctum/basis"
	"github.com/pavlo67/punctum/crud/selectors"
)

type ReadData struct {
	Description Description
	Item        interface{}
}

type ReadListData struct {
	Description Description
	Items       []interface{}
	AllCnt      uint64
}

type ResultData struct {
	IDs  []string
	Info string
}

func ReadOptionsFromParams(paramsTree basis.Params, pageLengthDefault uint64, sortByDefault []string) (*ReadOptions, uint64, error) {
	var readOptions ReadOptions
	if sortByParams, ok := paramsTree["sort_by"]; ok {
		sortBy, ok := sortByParams.([]string)
		if !ok {
			return nil, 0, errors.Errorf("non-[]string value for sort_by parameter: %#v", sortByParams)
		}
		readOptions.SortBy = sortBy
	}
	if len(readOptions.SortBy) < 1 {
		readOptions.SortBy = sortByDefault
	}

	pageStr := paramsTree.StringKeyDefault("page", "1")
	page, err := strconv.ParseUint(pageStr, 10, 64)
	if err != nil {
		return nil, 0, errors.Errorf("bad query page value: '%s'", pageStr)
	}

	var pageLength uint64
	pageLengthStr, ok := paramsTree.StringKey("page_length")
	if ok {
		pageLength, err = strconv.ParseUint(pageLengthStr, 10, 64)
		if err != nil {
			return nil, 0, errors.Errorf("bad query page_length value: '%s'", pageLengthStr)
		}
	} else {
		pageLength = pageLengthDefault
	}
	if pageLength == 0 {
		// TODO: WARNING
		pageLength = 200
	}

	if page > 1 {
		readOptions.Limits = []uint64{(page - 1) * pageLength, pageLength}
	} else {
		readOptions.Limits = []uint64{pageLength}
	}

	if selectorParams, ok := paramsTree["selector"]; ok {
		params, ok := selectorParams.(basis.Params)
		if !ok {
			return nil, page, errors.Errorf("wrong value type for selector parameter: %#v", selectorParams)
		}
		selector, err := selectors.FromParams(params)
		if err != nil {
			return nil, page, errors.Wrapf(err, "wrong value for selector parameter: %#v", selectorParams)
		}
		readOptions.Selector = selector
	}

	return &readOptions, page, nil
}

//func ReadOptionsFromRequest(r *http.Request, pageLengthDefault uint64, sortByDefault []string) (*ReadOptions, *httplib.ReadOptionsHTTP, error) {
//	if r == nil {
//		return nil, nil, errors.New("null request")
//	}
//
//	paramsTree := basis.Params{}
//	query := r.URL.Query()
//	for k, v := range query {
//		paramsTree[k] = v
//	}
//
//	readOptions, page, err := ReadOptionsFromParams(paramsTree, pageLengthDefault, sortByDefault)
//	if err != nil {
//		return nil, nil, err
//	}
//
//	readOptionsHTTP := &httplib.ReadOptionsHTTP{
//		Path:    r.URL.Path,
//		PageNum: page,
//	}
//
//	return readOptions, readOptionsHTTP, nil
//}

//func GetReadOptions(r *http.Request, pageLength uint64, defaultSortBy ...string) *ReadOptions {
//	if pageLength <= 0 {
//		pageLength = defaultPageLength
//	}
//
//	opt := &ReadOptions{ReadOptions: crud.ReadOptions{Limits: []uint64{0, pageLength}}}
//	if r == nil {
//		return opt
//	}
//
//	query := r.URL.Query()
//
//	opt.Path = r.URL.Path
//	opt.SortBy = query["sort"]
//	if len(opt.SortBy) < 1 {
//		opt.SortBy = defaultSortBy
//	}
//
//	page := query.Get("page")
//	if page != "" {
//		pageNum, err := strconv.ParseUint(page, 10, 64)
//		if err != nil {
//			log.Println("bad query page value: ", page, err)
//		} else {
//			opt.Limits[0] = pageLength * pageNum
//			opt.PageNum = pageNum
//		}
//	}
//	return opt
//}
