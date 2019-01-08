package config

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/pkg/errors"
	"github.com/yosuke-furukawa/json5/encoding/json5"

	"github.com/pavlo67/punctum/basis"
	"github.com/pavlo67/punctum/starter/logger"
)

var ErrNoConfig = errors.New("no config")
var ErrNoValue = errors.New("no value")

var l logger.Operator

// -----------------------------------------------------------------------------

func Get(path string, loggerToSet logger.Operator) (*PunctumConfig, error) {
	l = loggerToSet

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

	return readConfig(data)
}

type PunctumConfig struct {
	identity map[string]map[string]string

	system     map[string]string
	server     map[string]ServerTLS
	fileslocal map[string]map[string]string

	sender      map[string]map[string]string
	credentials map[string]map[string]string
	facebook    map[string][]string

	smtp      map[string]ServerAccess
	pop3      map[string]ServerAccess
	mysql     map[string]ServerAccess
	bolt      map[string]ServerAccess
	twitter   map[string]map[string]string
	instagram map[string]map[string]string
	google    map[string]map[string]string

	strings map[string]string
	flags   map[string]bool
}

func readConfig(data []byte) (*PunctumConfig, error) {
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

	conf := PunctumConfig{
		identity:    map[string]map[string]string{},
		server:      map[string]ServerTLS{},
		mysql:       map[string]ServerAccess{},
		bolt:        map[string]ServerAccess{},
		smtp:        map[string]ServerAccess{},
		sender:      map[string]map[string]string{},
		credentials: map[string]map[string]string{},

		instagram: map[string]map[string]string{},
		twitter:   map[string]map[string]string{},
		google:    map[string]map[string]string{},
		facebook:  map[string][]string{},

		fileslocal: map[string]map[string]string{},
		system:     map[string]string{},
		strings:    map[string]string{},
		flags:      map[string]bool{},
	}

	for k, v0 := range config {
		if k == "identity" {
			err = json5.Unmarshal(configRaw[k], &conf.identity)
		} else if k == "mysql" {
			err = json5.Unmarshal(configRaw[k], &conf.mysql)
		} else if k == "bolt" {
			err = json5.Unmarshal(configRaw[k], &conf.bolt)
		} else if k == "smtp" {
			err = json5.Unmarshal(configRaw[k], &conf.smtp)
		} else if k == "pop3" {
			err = json5.Unmarshal(configRaw[k], &conf.pop3)
		} else if k == "sender" {
			err = json5.Unmarshal(configRaw[k], &conf.sender)
		} else if k == "server" {
			err = json5.Unmarshal(configRaw[k], &conf.server)
		} else if k == "fileslocal" {
			err = json5.Unmarshal(configRaw[k], &conf.fileslocal)
		} else if k == "system" {
			err = json5.Unmarshal(configRaw[k], &conf.system)
		} else if k == "credentials" {
			err = json5.Unmarshal(configRaw[k], &conf.credentials)

		} else if k == "twitter" {
			err = json5.Unmarshal(configRaw[k], &conf.twitter)
		} else if k == "instagram" {
			err = json5.Unmarshal(configRaw[k], &conf.instagram)
		} else if k == "google" {
			err = json5.Unmarshal(configRaw[k], &conf.google)
		} else if k == "facebook" {
			err = json5.Unmarshal(configRaw[k], &conf.facebook)
		} else {
			switch v := v0.(type) {
			case string:
				conf.strings[k] = v
			case []byte:
				conf.strings[k] = string(v)
			case float64, float32:
				conf.strings[k] = fmt.Sprintf("%.3f", v)
				// no integer values in JSON, only float
			case bool:
				conf.flags[k] = v
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

func (c *PunctumConfig) Identity(key string, errs basis.Errors) (map[string]string, basis.Errors) {
	if c == nil {
		return nil, append(errs, ErrNoConfig)
	}
	if identity, ok := c.identity[key]; ok {
		return identity, errs
	}
	return nil, append(errs, errors.Wrapf(ErrNoValue, "no data for key '%s' in config.identity in %#v", key, c))
}

func (c *PunctumConfig) Sender(key string, errs basis.Errors) (map[string]string, basis.Errors) {
	if c == nil {
		return nil, append(errs, ErrNoConfig)
	}
	if sender, ok := c.sender[key]; ok {
		return sender, errs
	}
	return nil, append(errs, errors.Wrapf(ErrNoValue, "no data for key '%s' in config.sender in %#v", key, c))
}

func (c *PunctumConfig) Server(key string, errs basis.Errors) (ServerTLS, basis.Errors) {
	if c == nil {
		return ServerTLS{}, append(errs, ErrNoConfig)
	}
	if srv, ok := c.server[key]; ok {
		return srv, errs
	}
	return ServerTLS{}, append(errs, errors.Wrapf(ErrNoValue, "no data for key '%s' in config.serverhttp_jschmhr in %#v", key, c))
}

func (c *PunctumConfig) MySQL(database string, errs basis.Errors) (ServerAccess, basis.Errors) {
	if c == nil {
		return ServerAccess{}, append(errs, ErrNoConfig)
	}
	if mc, ok := c.mysql[database]; ok {
		if env, ok := os.LookupEnv("ENV"); ok {
			mc.Path += "_" + env
			//log.Println("UserIS HERE?")
		}
		return mc, errs
	}
	return ServerAccess{}, append(errs, errors.Wrapf(ErrNoValue, "no config.mysql for key '%s'", database))
}

func (c *PunctumConfig) Bolt(database string, errs basis.Errors) (ServerAccess, basis.Errors) {
	if c == nil {
		return ServerAccess{}, append(errs, ErrNoConfig)
	}
	if mc, ok := c.bolt[database]; ok {
		if env, ok := os.LookupEnv("ENV"); ok {
			mc.Path += "_" + env
		}
		return mc, errs
	}
	return ServerAccess{}, append(errs, errors.Wrapf(ErrNoValue, "no config.bolt for key '%s'", database))
}

func (c *PunctumConfig) SMTP(server string, errs basis.Errors) (ServerAccess, basis.Errors) {
	if c == nil {
		return ServerAccess{}, append(errs, ErrNoConfig)
	}

	if sc, ok := c.smtp[server]; ok {
		return sc, errs
	}

	return ServerAccess{}, append(errs, errors.Wrapf(ErrNoValue, "no config.smtp for key '%s'", server))
}

func (c *PunctumConfig) POP3(server string, errs basis.Errors) (ServerAccess, basis.Errors) {
	if c == nil {
		return ServerAccess{}, append(errs, ErrNoConfig)
	}

	if sc, ok := c.pop3[server]; ok {
		return sc, errs
	}

	return ServerAccess{}, append(errs, errors.Wrapf(ErrNoValue, "no config.pop3 for key '%s'", server))
}

func (c *PunctumConfig) System(key string, errs basis.Errors) (string, basis.Errors) {
	if c == nil {
		return "", append(errs, ErrNoConfig)
	}
	if str, ok := c.system[key]; ok {
		return str, errs
	}
	return "", append(errs, errors.Wrapf(ErrNoValue, "no data for key '%s' in config.system in %#v", key, c))
}

func (c *PunctumConfig) Credentials(key string, errs basis.Errors) (map[string]string, basis.Errors) {
	if c == nil {
		return nil, append(errs, ErrNoConfig)
	}
	if str, ok := c.credentials[key]; ok {
		return str, errs
	}
	return nil, append(errs, errors.Wrapf(ErrNoValue, "no data for key '%s' in config.credentials in %#v", key, c))
}

func (c *PunctumConfig) Fileslocal(key string, errs basis.Errors) (map[string]string, basis.Errors) {
	if c == nil {
		return nil, append(errs, ErrNoConfig)
	}
	if str, ok := c.fileslocal[key]; ok {
		return str, errs
	}
	return nil, append(errs, errors.Wrapf(ErrNoValue, "no data for key '%s' in config.fileslocal in %#v", key, c))
}

func (c *PunctumConfig) Instagram(key string, errs basis.Errors) (map[string]string, basis.Errors) {
	if c == nil {
		return nil, append(errs, ErrNoConfig)
	}
	if str, ok := c.instagram[key]; ok {
		return str, errs
	}
	return nil, append(errs, errors.Wrapf(ErrNoValue, "no data for key '%s' in config.instagram in %#v", key, c))
}

// Twitter ...
func (c *PunctumConfig) Twitter(key string, errs basis.Errors) (map[string]string, basis.Errors) {
	if c == nil {
		return nil, append(errs, ErrNoConfig)
	}
	if str, ok := c.twitter[key]; ok {
		return str, errs
	}
	return nil, append(errs, errors.Wrapf(ErrNoValue, "no data for key '%s' in config.twitter in %#v", key, c))
}

// Facebook ...
func (c *PunctumConfig) Facebook(errs basis.Errors) (map[string][]string, basis.Errors) {
	if c == nil {
		return nil, append(errs, ErrNoConfig)
	}

	return c.facebook, nil

	//if str, ok := c.facebook[key]; ok {
	//	return str, errs
	//}
	//return nil, append(errs, nil.Wrapf(ErrNoValue, "no data for key '%s' in config.facebook in %#v", key, c))
}

// Google ...
func (c *PunctumConfig) Google(key string, errs basis.Errors) (map[string]string, basis.Errors) {
	if c == nil {
		return nil, append(errs, ErrNoConfig)
	}
	if str, ok := c.google[key]; ok {
		return str, errs
	}
	return nil, append(errs, errors.Wrapf(ErrNoValue, "no data for key '%s' in config.google in %#v", key, c))
}

// String ...
func (c *PunctumConfig) String(key string, errs basis.Errors) (string, basis.Errors) {
	if c == nil {
		return "", append(errs, ErrNoConfig)
	}
	if str, ok := c.strings[key]; ok {
		return str, errs
	}
	return "", append(errs, errors.Wrapf(ErrNoValue, "no data for key '%s' in config.strings in %#v", key, c))
}

// Bool ...
func (c *PunctumConfig) Bool(key string, errs basis.Errors) (bool, basis.Errors) {
	if c == nil {
		return false, append(errs, ErrNoConfig)
	}
	if flag, ok := c.flags[key]; ok {
		return flag, errs
	}
	return false, append(errs, errors.Wrapf(ErrNoValue, "no data for key '%s' in config.flags in %#v", key, c))
}

// -----------------------------------------------------------------------------

//// func DataStructParse(data interface{}) map[string]string {
//
//// 	values := make(map[string]string)
//
//// 	s := structs.New(data)
//// 	for _, f := range s.Fields() {
//// 		if f.IsExported() {
//// 			if reflect.TypeOf(f.Value()).Kind() == reflect.Struct {
//// 				v1 := DataStructParse(f.Value())
//// 				for key, val := range v1 {
//// 					if _, ok := v1[key]; ok {
//// 						values[key] = val
//// 					}
//// 				}
//// 			} else if reflect.TypeOf(f.Value()).Kind() == reflect.Map {
//// 				for key, val := range f.Value().(map[string]string) {
//// 					if _, ok := f.Value().(map[string]string)[key]; ok {
//// 						values[key] = val
//// 					}
//// 				}
//// 			} else {
//// 				values[f.Nick()] = f.Value().(string)
//// 			}
//// 		}
//// 	}
//// 	return values
//// }
