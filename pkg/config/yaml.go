package config

type YamlReader interface {
	ReadFromFile(string, interface{}) error
}

type Yaml struct{}

func (y *Yaml) ReadFromFile(path string, obj interface{}) error { return nil }
