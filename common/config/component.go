package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"

	"github.com/pavlo67/workshop/common"
	"github.com/pkg/errors"
	"github.com/yosuke-furukawa/json5/encoding/json5"
)

type ComponentsIndex struct {
	Endpoints map[string]Endpoint
	MySQL     map[string]SQLTable
	SQLite    map[string]SQLTable
	Params    map[string]string
	ParamsArr map[string][]string
}

type componentsIndexRaw struct {
	Endpoints map[string]interface{}
	MySQL     map[string]SQLTable
	SQLite    map[string]SQLTable
	Params    map[string]interface{}
}

func ComponentIndex(indexPath string, errs common.Errors) (ComponentsIndex, common.Errors) {
	if indexPath == "" {
		return ComponentsIndex{}, errs
	}

	if indexPath[len(indexPath)-1] == '/' {
		indexPath += "index.json5"
	}

	if _, err := os.Stat(indexPath); err != nil {
		return ComponentsIndex{}, errs
	}

	data, err := ioutil.ReadFile(indexPath)
	if err != nil {
		return ComponentsIndex{}, append(errs, errors.Wrapf(err, "no index file in '%s'", indexPath))
	}

	var sciRaw componentsIndexRaw
	err = json5.Unmarshal(data, &sciRaw)
	if err != nil {
		return ComponentsIndex{}, append(errs, errors.Errorf("can't decode config data: '%s'\n", string(data)))
	}

	sci := &ComponentsIndex{
		Endpoints: map[string]Endpoint{},
		MySQL:     sciRaw.MySQL,
		Params:    map[string]string{},
		ParamsArr: map[string][]string{},
	}

	localPath := filepath.Dir(indexPath) + "/"
	for ke, e := range sciRaw.Endpoints {
		ep, err := readEndpoint(e, localPath)
		if err != nil {
			fmt.Printf("can't decode endpoint %s: %#v %s\n", ke, e, err)
		} else {
			sci.Endpoints[ke] = *ep
		}
	}

	for kp, p := range sciRaw.Params {
		par, err := readString(p)
		if err != nil {
			fmt.Printf("can't decode param %s: %#v %s\n", kp, p, err)
		} else if len(par) > 1 {
			sci.ParamsArr[kp] = par
		} else if len(par) == 1 {
			sci.Params[kp] = par[0]
		}
	}

	return *sci, errs
}

// -----------------------------------------------------------------------------

func stringifySlice(s0 []interface{}) ([]string, error) {
	var s1 []string
	for _, v0 := range s0 {
		if v, ok := v0.(string); ok {
			s1 = append(s1, v)
			//} else if v, ok := v0.(float64); ok {
			//	s1 = append(s1, strconv.FormatFloat(v, 10))
		} else {
			return nil, errors.Errorf("bad string value %v type %#v", v0, reflect.TypeOf(v0))
		}
	}
	return s1, nil
}

func readString(s0 interface{}) ([]string, error) {
	if s, ok := s0.(string); ok {
		return []string{s}, nil
	}

	var s1 []string
	if s, ok := s0.([]string); ok {
		s1 = s
	} else if s, ok := s0.([]interface{}); ok {
		sTmp, err := stringifySlice(s)
		if err != nil {
			return nil, errors.Wrapf(err, "bad string JSON: %#v", s0)
		}
		s1 = sTmp
		//} else if v, ok := s0.(float64); ok {
		//	s1 = append(s1, strconv.FormatFloat(v, 10))
	} else {
		return nil, errors.Errorf("bad string JSON type: %s", reflect.TypeOf(s0))
	}

	return s1, nil
}
