package config

import "testing"

// TestValidateV1Config has no endpoints set on the incoming config, should fail
func TestValidateV1Config1(t *testing.T) {
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
	}

	if err := Validate(cfg); err == nil {
		t.Errorf("Validation Failed, No Error Raised")
	}
}

// TestValidateV1Config2 has an invalid Endpoint defined, should fail
func TestValidateV1Config2(t *testing.T) {
	endpoints := make(map[string]map[string]*Endpoint)
	endpoints["/test"] = map[string]*Endpoint{
		"get": {},
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

	if err := Validate(cfg); err == nil {
		t.Errorf("Validation Failed, No Error Raised")
	}
}

// TestValidateV1Config3 has an invalid Request defined, should fail
func TestValidateV1Config3(t *testing.T) {
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
			"testRequest": nil,
		},
		Endpoints: nil,
	}

	if err := Validate(cfg); err == nil {
		t.Errorf("Validation Failed, No Error Raised")
	}
}

// TestValidateV1Config4 has an invalid startup action defined, should fail
func TestValidateV1Config4(t *testing.T) {
	cfg := &Config{
		Version: 1.0,
		Services: map[string]*Service{
			"testService": {Hostname: "localhost", Port: 8080},
		},
		StartupActions: []map[string]interface{}{{"delay": nil}},
	}

	if err := Validate(cfg); err == nil {
		t.Errorf("Validation Failed, No Error Raised")
	}
}

// TestValidateV1Config5 has an invalid service field defined, should fail
func TestValidateV1Config5(t *testing.T) {
	cfg := &Config{
		Version: 1.0,
		Services: map[string]*Service{
			"testService": {Hostname: "", Port: 8080},
		},
	}

	if err := Validate(cfg); err == nil {
		t.Errorf("Validation Failed, No Error Raised")
	}
}

// TestValidateV1Config6 attempts to get the hostname of the OS but experiences an error getting it
func TestValidateV1Config6(t *testing.T) {
	osHostname = mockInvalidOsHostname

	cfg := &Config{
		Version: 1.0,
		Services: map[string]*Service{
			"testService": {Hostname: "$HOSTNAME", Port: 8080},
		},
	}

	if err := Validate(cfg); err == nil {
		t.Errorf("Validation Failed, No Error Raised")
	}
}
