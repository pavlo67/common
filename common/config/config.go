package config

import (
	"io/ioutil"

	"github.com/pavlo67/workshop/common/logger"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

var ErrNoConfig = errors.New("no config")
var ErrNoLogger = errors.New("no logger")
var ErrNoValue = errors.New("no value")

// -----------------------------------------------------------------------------

type Config struct {
	ServerHTTP Server `yaml:"server_http,omitempty"  json:"server_http,omitempty"`

	Postgres ServerAccess `yaml:"postgres,omitempty" json:"postgres,omitempty"`
	MySQL    ServerAccess `yaml:"mysql,omitempty"    json:"mysql,omitempty"`
	SQLite   ServerAccess `yaml:"sqlite,omitempty"   json:"sqlite,omitempty"`
	SMTP     ServerAccess `yaml:"smtp,omitempty"     json:"smtp,omitempty"`
	POP3     ServerAccess `yaml:"pop3,omitempty"     json:"pop3,omitempty"`

	Envs map[string]string `yaml:"envs,omitempty"  json:"envs,omitempty"`

	Logger logger.Operator `yaml:"-"  json:"-"`
}

// -----------------------------------------------------------------------------

func Get(path, environment string) (*Config, logger.Operator, error) {

	if len(path) < 1 {
		return nil, nil, errors.New("empty config path")
	}

	cfgFile := path + "/" + environment + ".yaml"
	data, err := ioutil.ReadFile(cfgFile)
	if err != nil {
		return nil, nil, errors.Wrapf(err, "can't read config file from '%s'", cfgFile)
	}

	cfg := &Config{}
	err = yaml.Unmarshal(data, cfg)
	if err != nil {
		return nil, nil, errors.Wrapf(err, "can't yaml.Unmarshal('%s') from config '%s'", data, cfgFile)
	}

	// TODO: use debug level from environment or config
	cfg.Logger, err = logger.Init(logger.Config{LogLevel: logger.DebugLevel})
	if err != nil {
		return nil, nil, errors.Wrap(err, "can't logger.Init()")
	}

	return cfg, cfg.Logger, nil
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
