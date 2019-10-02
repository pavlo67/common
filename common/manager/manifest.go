package manager

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type Manifest struct {
	AppKey    string   `yaml:"app_key"`
	Requested []string `yaml:"requested"`
	Subpaths  []string `yaml:"subpaths"`
	Command   string   `yaml:"command"`
	Args      []string `yaml:"args"`
	Workdir   string   `yaml:"workdir"`
	Logdir    string   `yaml:"logdir"`
}

func ReadManifest(path string) (*Manifest, error) {
	data, err := ioutil.ReadFile(path + "/manifest.yaml")
	if err != nil {
		return nil, err
	}

	var manifest Manifest
	err = yaml.Unmarshal(data, &manifest)
	if err != nil {
		return nil, err
	}

	return &manifest, nil
}
