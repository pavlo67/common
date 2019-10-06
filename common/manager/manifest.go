package manager

import (
	"io/ioutil"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

type Manifest struct {
	AppKey    string   `yaml:"app_key"`
	Requested []string `yaml:"requested"`
	Subpaths  []string `yaml:"subpaths"`
	Command   string   `yaml:"command"`
	Args      []string `yaml:"args"`
	Workdir   string   `yaml:"workdir"`
}

func ReadManifest(path string) (*Manifest, error) {
	data, err := ioutil.ReadFile(path + "/manifest.yaml")
	if err != nil {
		return nil, errors.Wrapf(err, "on ioutil.ReadFile('%s/manifest.yaml')", path)
	}

	var manifest Manifest
	err = yaml.Unmarshal(data, &manifest)
	if err != nil {
		return nil, errors.Wrapf(err, "on yaml.Unmarshal('%s'), data from '%s/manifest.yaml'", data, path)
	}

	return &manifest, nil
}
