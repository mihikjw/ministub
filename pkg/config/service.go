package config

// Service represents the cfg definition for a microservice
type Service struct {
	Hostname string `yaml:"hostname"`
	Port     int    `yaml:"port"`
}
