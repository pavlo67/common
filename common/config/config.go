package config

import (
	"io/ioutil"

	"github.com/pavlo67/workshop/common/libraries/encodelib"
	"github.com/pkg/errors"
)

// -----------------------------------------------------------------------------

type Config struct {
	data      map[string]interface{}
	marshaler encodelib.Marshaler
}

type Access struct {
	Host    string
	Port    int
	User    string
	Pass    string
	Path    string
	Options string
}

func (c Config) Value(key string, target interface{}) error {
	if value, ok := c.data[key]; ok {
		valueRaw, err := c.marshaler.Marshal(value)
		if err != nil {
			return errors.Wrapf(err, "can't marshal value (%s / %#v) to raw bytes", key, value)
		}

		return c.marshaler.Unmarshal(valueRaw, target)
	}

	return nil
}

// -----------------------------------------------------------------------------

func Get(cfgFile string, marshaler encodelib.Marshaler) (*Config, error) {

	if len(cfgFile) < 1 {
		return nil, errors.New("empty config path")
	}

	data, err := ioutil.ReadFile(cfgFile)
	if err != nil {
		return nil, errors.Wrapf(err, "can't read config file from '%s'", cfgFile)
	}

	cfg := Config{marshaler: marshaler}
	err = marshaler.Unmarshal(data, &cfg.data)
	if err != nil {
		return nil, errors.Wrapf(err, "can't .Unmarshal('%s') from config '%s'", data, cfgFile)
	}

	return &cfg, nil
}

//// String ...
//func (c *Config) String(key string, errs common.Errors) (string, common.Errors) {
//	if c == nil {
//		return "", append(errs, ErrNoConfig)
//	}
//	if str, ok := c.Strings[key]; ok {
//		return str, errs
//	}
//	return "", append(errs, errors.Wrapf(ErrNoValue, "no data for key '%s' in config.strings in %#v", key, c))
//}
//
//// Bool ...
//func (c *Config) Bool(key string, errs common.Errors) (bool, common.Errors) {
//	if c == nil {
//		return false, append(errs, ErrNoConfig)
//	}
//	if flag, ok := c.Flags[key]; ok {
//		return flag, errs
//	}
//	return false, append(errs, errors.Wrapf(ErrNoValue, "no data for key '%s' in config.flags in %#v", key, c))
//}