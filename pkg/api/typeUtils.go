package api

import (
	"fmt"
	"strconv"
	"strings"
)

/*AssertValidType returns whether a given value is of the expected type
if the initial conversion fails, it will convert to string then convert to type where appropriate */
func AssertValidType(value interface{}, expectedType string) bool {
	switch {
	case expectedType == "boolean":
		if _, ok := value.(bool); !ok {
			if strValue, ok := value.(string); ok {
				if lowerInValue := strings.ToLower(strValue); lowerInValue != "true" && lowerInValue != "false" {
					return false
				}
			} else {
				return false
			}
		}
	case expectedType == "integer":
		/* when raw ints come in from some types they can only be extracted as float64 - convert to int from here is always possible so
		dont bother doing this test (cant even check its success, it always works) */
		if _, ok := value.(float64); !ok {
			if strValue, ok := value.(string); ok {
				if _, err := strconv.Atoi(strValue); err != nil {
					return false
				}
			} else {
				return false
			}
		}
	case expectedType == "string":
		if _, ok := value.(string); !ok {
			return false
		}
	case expectedType == "float":
		if _, ok := value.(float64); !ok {
			if _, err := strconv.ParseFloat(value.(string), 64); err != nil {
				return false
			}
		}
	case expectedType == "array":
		if _, ok := value.([]interface{}); !ok {
			return false
		}
	case expectedType == "object":
		if _, ok := value.(map[string]interface{}); !ok {
			return false
		}
	default:
		return false
	}

	return true
}

// AssertValidTypeFromPath goes down a given path for a given inputData JSON, and asserts the expected valid type when at the expected level for the given path
func AssertValidTypeFromPath(path, expectedType string, inputData interface{}) error {
	if splitPath := strings.Split(path, "."); len(splitPath) >= 1 {
		if len(splitPath) > 1 {
			if value, err := strconv.Atoi(splitPath[0]); err == nil {
				nextInputData := inputData.([]interface{})[value]
				return AssertValidTypeFromPath(strings.Join(splitPath[1:], "."), expectedType, nextInputData)
			} else if body, valid := inputData.(map[string]interface{}); valid {
				nextInputData := body[splitPath[0]]
				return AssertValidTypeFromPath(strings.Join(splitPath[1:], "."), expectedType, nextInputData)
			}
			return fmt.Errorf("Invalid Path Item %s", splitPath[0])
		}

		var valueToEvaluate interface{}
		if value, err := strconv.Atoi(splitPath[0]); err == nil {
			valueToEvaluate = inputData.([]interface{})[value]
		} else if body, valid := inputData.(map[string]interface{}); valid {
			valueToEvaluate = body[splitPath[0]]
		} else {
			valueToEvaluate = inputData
		}

		if !AssertValidType(valueToEvaluate, expectedType) {
			return fmt.Errorf("%s Is Invalid Expected Type %s", path, expectedType)
		}

		return nil
	}

	return fmt.Errorf("Empty Path Given For assertValidTypeFromPath")
}
