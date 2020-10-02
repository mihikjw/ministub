package config

import (
	"fmt"
	"os"
)

// osHostname returns the hostname of the OS
var osHostname = os.Hostname

// osGetenv returns the requested environment variable
var osGetenv = os.Getenv

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

// validateJSON ensures the given input JSON is valid (all keys are string)
func validateJSON(input interface{}) interface{} {
	if json, valid := input.(map[string]interface{}); valid {
		correctOutput := make(map[string]interface{}, len(json))
		for k, v := range json {
			correctOutput[k] = validateJSON(v)
		}
		return correctOutput
	}
	if invalidJSON, valid := input.(map[interface{}]interface{}); valid {
		correctOutput := make(map[string]interface{}, len(invalidJSON))
		for k, v := range invalidJSON {
			if newK, valid := k.(string); valid {
				correctOutput[newK] = validateJSON(v)
			}
		}
		return correctOutput
	}
	return input
}

// getEnvValueForField checks if a given string value is actually trying to use an environment variable, and replaces it if so
func getEnvValueForField(field string) (string, error) {
	if len(field) > 0 && string(field[0]) == "$" {
		searchVar := field[1:]

		switch {
		case searchVar == "HOSTNAME":
			newField, err := osHostname()
			if err != nil {
				return "", fmt.Errorf("Unable To Get Requested Hostname: %s", err.Error())
			}
			return newField, nil
		default:
			if result := osGetenv(searchVar); len(result) != 0 {
				return result, nil
			}
			return "", fmt.Errorf("Env Var Not Found: %s", searchVar)
		}
	}

	return field, nil
}
