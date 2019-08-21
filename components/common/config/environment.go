package config

import (
	"fmt"
	"io/ioutil"

	"github.com/pkg/errors"
	"github.com/yosuke-furukawa/json5/encoding/json5"

	"github.com/pavlo67/constructor/components/common"
	"github.com/pavlo67/constructor/components/common/logger"
)

var ErrNoConfig = errors.New("no config")
var ErrNoLogger = errors.New("no logger")
var ErrNoValue = errors.New("no value")

// -----------------------------------------------------------------------------

func Get(path string, l logger.Operator) (*Config, error) {

	var data []byte
	var err error
	if len(path) < 1 {
		return nil, errors.New("empty config path")
	} else if path[len(path)-1] == '/' {
		data, err = ioutil.ReadFile(path + "cfg.json5")
	} else {
		data, err = ioutil.ReadFile(path)
	}
	if err != nil {
		return nil, errors.Wrapf(err, "no config file in %v", path)
	}

	return readConfig(data, l)
}

type Config struct {
	//Identity map[string]string
	//Credentials map[string]string

	Server ServerTLS

	SMTP ServerAccess
	POP3 ServerAccess

	MySQL    ServerAccess
	SQLite   ServerAccess
	Postgres ServerAccess

	//Facebook  []string
	//Twitter   map[string]string
	//Instagram map[string]string
	//Google    map[string]string

	Strings map[string]string
	Flags   map[string]bool

	Logger logger.Operator
}

func readConfig(data []byte, l logger.Operator) (*Config, error) {
	if l == nil {
		return nil, errors.New("no logger")
	}

	var configRaw map[string]json5.RawMessage
	err := json5.Unmarshal(data, &configRaw)
	if err != nil {
		return nil, errors.Wrapf(err, "error reading json to configRaw: %v", string(data))
	}

	var config map[string]interface{}

	err = json5.Unmarshal(data, &config)
	if err != nil {
		return nil, errors.Wrapf(err, "error reading json to config: %v", string(data))
	}

	conf := Config{
		//Identity: map[string]string{},
		Strings: map[string]string{},
		Flags:   map[string]bool{},
		Logger:  l,
	}

	for k, v0 := range config {
		switch k {
		//case "identity":
		//	err = json5.Unmarshal(configRaw[k], &conf.identity)

		case "mysql":
			err = json5.Unmarshal(configRaw[k], &conf.MySQL)
		case "sqlite":
			err = json5.Unmarshal(configRaw[k], &conf.SQLite)
		case "postgres":
			err = json5.Unmarshal(configRaw[k], &conf.Postgres)

		case "smtp":
			err = json5.Unmarshal(configRaw[k], &conf.SMTP)
		case "pop3":
			err = json5.Unmarshal(configRaw[k], &conf.POP3)
		//case "sender":
		//	err = json5.Unmarshal(configRaw[k], &conf.Sender)
		case "server":
			err = json5.Unmarshal(configRaw[k], &conf.Server)
		//case "fileslocal":
		//	err = json5.Unmarshal(configRaw[k], &conf.fileslocal)
		//case "paths":
		//	err = json5.Unmarshal(configRaw[k], &conf.paths)
		//case "credentials":
		//	err = json5.Unmarshal(configRaw[k], &conf.Credentials)

		//case "twitter":
		//	err = json5.Unmarshal(configRaw[k], &conf.twitter)
		//case "instagram":
		//	err = json5.Unmarshal(configRaw[k], &conf.instagram)
		//case "google":
		//	err = json5.Unmarshal(configRaw[k], &conf.google)
		//case "facebook":
		//	err = json5.Unmarshal(configRaw[k], &conf.facebook)
		default:
			switch v := v0.(type) {
			case string:
				conf.Strings[k] = v
			case []byte:
				conf.Strings[k] = string(v)
			case float64, float32:
				conf.Strings[k] = fmt.Sprintf("%.3f", v)
				// no integer values in JSON, only float
			case bool:
				conf.Flags[k] = v
			default:
				l.Errorf("bad config value for key `%s`: %#v\n", k, v)
			}
			continue

		}

		if err != nil {
			fmt.Printf("can't decode config value %v: %v\n", k, string(configRaw[k]))
		}

	}

	return &conf, nil
}

// String ...
func (c *Config) String(key string, errs common.Errors) (string, common.Errors) {
	if c == nil {
		return "", append(errs, ErrNoConfig)
	}
	if str, ok := c.Strings[key]; ok {
		return str, errs
	}
	return "", append(errs, errors.Wrapf(ErrNoValue, "no data for key '%s' in config.strings in %#v", key, c))
}

// Bool ...
func (c *Config) Bool(key string, errs common.Errors) (bool, common.Errors) {
	if c == nil {
		return false, append(errs, ErrNoConfig)
	}
	if flag, ok := c.Flags[key]; ok {
		return flag, errs
	}
	return false, append(errs, errors.Wrapf(ErrNoValue, "no data for key '%s' in config.flags in %#v", key, c))
}
