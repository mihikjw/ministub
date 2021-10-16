package config

import "fmt"

func LoadCfgFile(path string, loader YamlReader) (*Config, error) {
	result := &Config{}
	if err := loader.ReadFromFile(path, result); err != nil {
		return nil, fmt.Errorf("error loading config: %s", err.Error())
	}
	return result, nil
}
