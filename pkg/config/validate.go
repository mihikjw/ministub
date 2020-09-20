package config

import "fmt"

// Validate ensures the loaded config meets the required standard
func Validate(cfg *Config) (err error) {
	switch cfg.Version {
	case 1.0:
		return validateV1Config(cfg)
	default:
		return fmt.Errorf("Unsupported Version: %f", cfg.Version)
	}
}

// supportedType checks if the given type string is supported
func supportedType(toCheck string) bool {
	switch {
	case toCheck == "string":
		return true
	case toCheck == "integer":
		return true
	case toCheck == "float":
		return true
	case toCheck == "boolean":
		return true
	case toCheck == "array":
		return true
	case toCheck == "object":
		return true
	default:
		return false
	}
}

// supportedAction ensures a given action string is supported by the application
func supportedAction(toCheck string) bool {
	switch {
	case toCheck == "delay":
		return true
	case toCheck == "request":
		return true
	default:
		return false
	}
}
