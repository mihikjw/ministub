package config

// Service represents the cfg definition for a microservice
type Service struct {
	Hostname string `yaml:"hostname"`
	Port     int    `yaml:"port"`
	Probe    *Probe `yaml:"probe"`
}

// Probe represents a cfg definition for a probe for a microservice
type Probe struct {
	Endpoint string `yaml:"endpoint"`
	Result   int    `yaml:"result"`
	Timeout  int    `yaml:"timeout"`
}
