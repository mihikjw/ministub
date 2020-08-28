package config

import (
	"fmt"
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"
)

/* save os/io/json methods to variables to allow mocking, normally these would be done through
a mock interface but its only for the LoadFromFile method so not worth the effort */
var osStat = os.Stat
var ioReadFile = ioutil.ReadFile
var yamlUnmarshal = yaml.Unmarshal

// Config holds all the data required to operate the application
type Config struct {
	Version        float32                         `yaml:"version"`
	Services       map[string]*Service             `yaml:"services"`
	StartupActions []map[string]interface{}        `yaml:"startupActions"`
	Requests       map[string]*Request             `yaml:"requests"`
	Endpoints      map[string]map[string]*Endpoint `yaml:"endpoints"`
}

// LoadFromFile creates a new Config object from the given filepath
func LoadFromFile(path string) (*Config, error) {
	if len(path) == 0 {
		return nil, fmt.Errorf("Invalid Path %s", path)
	}

	_, err := osStat(path)
	if err != nil {
		return nil, fmt.Errorf("File %s Not Found", path)
	}

	content, err := ioReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("Unable To Read File %s: %s", path, err.Error())
	}

	cfg := new(Config)
	if err = yamlUnmarshal(content, &cfg); err != nil {
		return nil, fmt.Errorf("Unable To Unmarshal Cfg: %s", err.Error())
	}

	return cfg, nil
}
