package config

import (
	"testing"
)

// TestValidate1 ensures a given valid V1 config with all components defined, is indeed valid. Failure cases beyond the initial func are on the individual funcs tests
func TestValidate1(t *testing.T) {
	endpoints := make(map[string]map[string]*Endpoint)
	endpoints["/test"] = map[string]*Endpoint{
		"get": {
			Params: &Parameters{
				Query: map[string]*ParamEntry{"test": {Type: "string", Required: true}},
				Path:  map[string]*ParamEntry{"test": {Type: "boolean", Required: true}},
			},
			Recieves: &Recieves{
				Headers: map[string]string{"foo": "bar"},
				Body:    map[string]string{"example_array.0.foo": "string"},
			},
			Responses: map[int]*Response{
				200: {
					Headers: map[string]string{"foo": "bar"},
					Body:    map[string]interface{}{"bar": "foo"},
					Weight:  100,
					Actions: []map[string]interface{}{
						{"delay": 10},
						{"request": map[interface{}]interface{}{"target": "testService", "id": "testRequest"}},
					},
				},
			},
			Actions: []map[string]interface{}{
				{"delay": 10},
				{"request": map[interface{}]interface{}{"target": "testService", "id": "testRequest"}},
			},
		},
	}

	cfg := &Config{
		Version: 1.0,
		Services: map[string]*Service{
			"testService": {Hostname: "localhost", Port: 8080},
		},
		StartupActions: []map[string]interface{}{
			{"delay": 10},
			{"request": map[interface{}]interface{}{"target": "testService", "id": "testRequest"}},
		},
		Requests: map[string]*Request{
			"testRequest": {
				URL:      "/test",
				Protocol: "http",
				Method:   "get",
				Headers:  map[string]string{"foo": "bar"},
				Body:     nil,
				ExpectedResponse: &Response{
					StatusCode: 200,
					Body:       map[string]interface{}{"foo.bar": "string"},
					Headers:    nil,
					Weight:     100,
					Actions: []map[string]interface{}{
						{"delay": 10},
						{"request": map[interface{}]interface{}{"target": "testService", "id": "testRequest"}},
					},
				},
			},
		},
		Endpoints: endpoints,
	}

	if err := Validate(cfg); err != nil {
		t.Errorf("Validation Failed: %s", err.Error())
	}
}

// TestValidate2 ensures an unsupported config version is returned as an error instead of processed
func TestValidate2(t *testing.T) {
	cfg := &Config{Version: 0}

	if err := Validate(cfg); err == nil {
		t.Errorf("No Error Raised For Invalid Version Param: %f", cfg.Version)
	}
}

// TestSupportedType1 ensures all the supported types are returned correctly as supported
func TestSupportedType1(t *testing.T) {
	for _, typeDef := range []string{"string", "integer", "float", "boolean", "array", "object"} {
		if !supportedType(typeDef) {
			t.Errorf("Supported Type '%s' Returned As Unsupported", typeDef)
		}
	}
}

// TestSupportedType2 ensures an unsupported type is not returned as supported
func TestSupportedType2(t *testing.T) {
	if supportedType("unsupported") {
		t.Errorf("Unsupported Type 'unsupported' Returned As Supported")
	}
}

// TestSupportedAction1 ensures all the supported action types are returned correctly as supported
func TestSupportedAction1(t *testing.T) {
	for _, action := range []string{"delay", "request"} {
		if !supportedAction(action) {
			t.Errorf("Supported Action '%s' Returned As Unsupported", action)
		}
	}
}

// TestSupportedAction2 ensures an unsupported action is not returned as supported
func TestSupportedAction2(t *testing.T) {
	if supportedAction("unsupported") {
		t.Errorf("Unsupported Action 'unsupported' Returned As Supported")
	}
}

// TestValidateJSON1 ensures a given piece of JSON without string-set keys is returned with string-set keys
func TestValidateJSON1(t *testing.T) {
	input := map[string]interface{}{
		"foo": map[interface{}]interface{}{
			"bar": 1234,
		},
	}

	output := validateJSON(input)

	if tmp1, ok := output.(map[string]interface{}); ok {
		if fooEntry, found := tmp1["foo"]; found {
			tmp2, ok := fooEntry.(map[string]interface{})
			if ok {
				if barEntry, found := tmp2["bar"]; found {
					if barEntry.(int) != 1234 {
						t.Errorf("barEntry Does Not Match Input")
					}
				} else {
					t.Errorf("Returned Validated JSON Is Invalid: Missing Key 'bar'")
				}
			} else {
				t.Errorf("Returned Validated JSON Is Invalid: tmp2")
			}
		} else {
			t.Errorf("Returned Validated JSON Is Invalid: Missing Key 'foo'")
		}
	} else {
		t.Errorf("Returned Validated JSON Is Invalid: tmp1")
	}
}

// TestValidGetEnvValueForField1 ensures a correct response for getting hostname, env vars and an unmodified value
func TestValidGetEnvValueForField(t *testing.T) {
	osHostname = mockValidOsHostname
	osGetenv = mockValidOsGetenv

	for _, dataToCheck := range []string{"$HOSTNAME", "$ENV_VAR", "test_value"} {
		checkedData, err := getEnvValueForField(dataToCheck)

		if err != nil {
			t.Errorf("Error Encountered Getting Env Value For Field: %s", err.Error())
		}

		if checkedData != "test_value" {
			t.Errorf("Value Does Not Match 'test_value': %s", checkedData)
		}
	}

}

// TestValidGetEnvValueForField2 ensures when a hostname cannot be found, an error is returned
func TestValidGetEnvValueForField2(t *testing.T) {
	osHostname = mockInvalidOsHostname

	if _, err := getEnvValueForField("$HOSTNAME"); err == nil {
		t.Errorf("Invalid Hostname Returned No Error")
	}
}

// TestValidGetEnvValueForField3 ensures when an env var cannot be found, an error is returned
func TestValidGetEnvValueForField3(t *testing.T) {
	osGetenv = mockInvalidOsGetenv

	if _, err := getEnvValueForField("$ENV_VAR"); err == nil {
		t.Errorf("Invalid Hostname Returned No Error")
	}
}
