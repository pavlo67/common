package config

import (
	"io/ioutil"

	"github.com/pavlo67/common/common/serialization"

	"github.com/pavlo67/common/common"
	"github.com/pavlo67/common/common/errors"
)

type Envs struct {
	serviceName string
	data        map[string]interface{}
	marshaler   serialization.Marshaler
}

var errNoEnvs = errors.New("no envs")

func (c *Envs) ServiceName() string {
	if c == nil {
		return ""
	}

	return c.serviceName
}

func (c *Envs) Raw(key string) (interface{}, error) {
	if c == nil {
		return nil, errNoEnvs
	}

	valueRaw, ok := c.data[key]
	if !ok {
		return nil, errors.CommonError(common.NotFoundKey, common.Map{"reason": "no key in envs", "key": key})
	}

	return valueRaw, nil
}

func (c *Envs) Value(key string, target interface{}) error {
	if c == nil {
		return errNoEnvs
	}

	if value, ok := c.data[key]; ok {
		valueRaw, err := c.marshaler.Marshal(value)
		if err != nil {
			return errors.Wrapf(err, "can't marshal value (%s / %#v) to raw bytes", key, value)
		}

		return c.marshaler.Unmarshal(valueRaw, target)
	}

	return errors.CommonError(common.NotFoundKey, common.Map{"reason": "no key in envs", "key": key})
}

// -----------------------------------------------------------------------------

func Get(envsFile string, marshaler serialization.Marshaler) (*Envs, error) {

	if len(envsFile) < 1 {
		return nil, errors.New("empty envs path")
	}

	data, err := ioutil.ReadFile(envsFile)
	if err != nil {
		return nil, errors.Wrapf(err, "can't read envs file from '%s'", envsFile)
	}

	cfg := Envs{marshaler: marshaler}
	err = marshaler.Unmarshal(data, &cfg.data)
	if err != nil {
		return nil, errors.Wrapf(err, "can't .Unmarshal('%s') from envs '%s'", data, envsFile)
	}

	return &cfg, nil
}

//// Key ...
//func (c *Envs) Key(key string, errs common.multipleErrors) (string, common.multipleErrors) {
//	if c == nil {
//		return "", append(errs, ErrNoConfig)
//	}
//	if str, ok := c.Strings[key]; ok {
//		return str, errs
//	}
//	return "", append(errs, errors.Wrapf(ErrNoValue, "no data for key '%s' in config.strings in %#v", key, c))
//}
//
//// IsTrue ...
//func (c *Envs) IsTrue(key string, errs common.multipleErrors) (bool, common.multipleErrors) {
//	if c == nil {
//		return false, append(errs, ErrNoConfig)
//	}
//	if flag, ok := c.Flags[key]; ok {
//		return flag, errs
//	}
//	return false, append(errs, errors.Wrapf(ErrNoValue, "no data for key '%s' in config.flags in %#v", key, c))
//}