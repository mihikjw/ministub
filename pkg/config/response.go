package config

// Response represents a response from request
type Response struct {
	StatusCode int                    `yaml:"statusCode"`
	Body       map[string]interface{} `yaml:"body"`
	Headers    map[string]string      `yaml:"headers"`
	Weight     int                    `yaml:"weight"`
}
