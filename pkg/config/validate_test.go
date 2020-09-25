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
