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

// TestValidateV1Actions1 ensures a correct set of actions does not raise any errors
func TestValidateV1Actions1(t *testing.T) {
	actions := []map[string]interface{}{
		{"delay": 10},
		{"request": map[interface{}]interface{}{"target": "testService", "id": "testRequest"}},
	}
	requests := map[string]*Request{"testRequest": nil}

	if err := validateV1Actions(actions, map[string]bool{"testService": true}, requests); err != nil {
		t.Errorf("validateV1Actions Error Raised: %s", err.Error())
	}
}

// TestValidateV1Actions2 ensures when a missing request is referenced, it is raised to the user
func TestValidateV1Actions2(t *testing.T) {
	actions := []map[string]interface{}{
		{"delay": 10},
		{"request": map[interface{}]interface{}{"target": "testService", "id": "testRequest"}},
	}

	if err := validateV1Actions(actions, map[string]bool{"testService": true}, make(map[string]*Request)); err == nil {
		t.Errorf("validateV1Actions No Error Raised")
	}
}

// TestValidateV1Actions3 ensures when a requestID is missing from a 'request' action, it is raised to the user
func TestValidateV1Actions3(t *testing.T) {
	actions := []map[string]interface{}{
		{"delay": 10},
		{"request": map[interface{}]interface{}{"target": "testService"}},
	}

	if err := validateV1Actions(actions, map[string]bool{"testService": true}, make(map[string]*Request)); err == nil {
		t.Errorf("validateV1Actions No Error Raised")
	}
}

// TestValidateV1Actions4 ensures when a request uses a target service that is not defined, it is raised to the user
func TestValidateV1Actions4(t *testing.T) {
	actions := []map[string]interface{}{
		{"delay": 10},
		{"request": map[interface{}]interface{}{"target": "testService"}},
	}

	if err := validateV1Actions(actions, make(map[string]bool), make(map[string]*Request)); err == nil {
		t.Errorf("validateV1Actions No Error Raised")
	}
}

// TestValidateV1Actions5 ensures when a request is missing a target service, it is raised to the user
func TestValidateV1Actions5(t *testing.T) {
	actions := []map[string]interface{}{
		{"delay": 10},
		{"request": make(map[interface{}]interface{})},
	}

	if err := validateV1Actions(actions, make(map[string]bool), make(map[string]*Request)); err == nil {
		t.Errorf("validateV1Actions No Error Raised")
	}
}

// TestValidateV1Actions6 ensures when a request is invalid formatted, it is raised to the user
func TestValidateV1Actions6(t *testing.T) {
	actions := []map[string]interface{}{
		{"delay": 10},
		{"request": "invalid"},
	}

	if err := validateV1Actions(actions, make(map[string]bool), make(map[string]*Request)); err == nil {
		t.Errorf("validateV1Actions No Error Raised")
	}
}

// TestValidateV1Actions7 ensures when an invalid action is requested, it is raised to the user
func TestValidateV1Actions7(t *testing.T) {
	actions := []map[string]interface{}{
		{"unsupported": "foo"},
	}

	if err := validateV1Actions(actions, make(map[string]bool), make(map[string]*Request)); err == nil {
		t.Errorf("validateV1Actions No Error Raised")
	}
}

// TestValidateV1Parameters1 ensures a given set of correctly formatted parameters is not raised as an error
func TestValidateV1Parameters1(t *testing.T) {
	params := map[string]*ParamEntry{
		"test": {Required: true, Type: "string"},
	}

	if err := validateV1Parameters(params); err != nil {
		t.Errorf("validateV1Parameters Error Raised: %s", err.Error())
	}
}

// TestValidateV1Parameters2 ensures an incorrect type is raised as an error
func TestValidateV1Parameters2(t *testing.T) {
	params := map[string]*ParamEntry{
		"test": {Required: true, Type: "unsupported"},
	}

	if err := validateV1Parameters(params); err == nil {
		t.Errorf("validateV1Parameters No Error Raised")
	}
}

// TestValidateV1Method1 ensures all supported http methods are correctly returned as supported
func TestValidateV1Method1(t *testing.T) {
	for _, method := range []string{"get", "post", "put", "delete"} {
		if !validateV1Method(method) {
			t.Errorf("Method '%s' Incorrectly Detected As Unsupported", method)
		}
	}
}

// TestValidateV1Method2 ensures an unsupported http method is correctly returned as unsupported
func TestValidateV1Method2(t *testing.T) {
	if validateV1Method("unsupported") {
		t.Errorf("Method 'unsupported' Incorrectly Detected As Supported")
	}
}

// TestValidateV1Protocol1 ensures all supported protocols are correctly returned as supported
func TestValidateV1Protocol1(t *testing.T) {
	for _, protocol := range []string{"http", "https"} {
		if !validateV1Protocol(protocol) {
			t.Errorf("Protocol '%s' Incorrectly Detected As Unsupported", protocol)
		}
	}
}

// TestValidateV1Protocol2 ensures an unsupported protocol is correctly raised as unsupported
func TestValidateV1Protocol2(t *testing.T) {
	if validateV1Protocol("unsupportted") {
		t.Errorf("Protocol 'unsupported' Incorrectly Detected As Supported")
	}
}

// TestValidateV1Endpoint1 ensures a correctly formatted endpoint object passes validation
func TestValidateV1Endpoint1(t *testing.T) {
	endpoint := &Endpoint{
		Params: &Parameters{
			Query: make(map[string]*ParamEntry),
			Path:  make(map[string]*ParamEntry),
		},
		Recieves: &Recieves{
			Headers: map[string]string{"Content-Type": "application/json"},
			Body:    map[string]string{"foo.bar": "integer"},
		},
		Responses: map[int]*Response{
			200: {
				StatusCode: 200,
				Body:       map[string]interface{}{"foo.bar": "string"},
				Headers:    nil,
				Weight:     100,
				Actions: []map[string]interface{}{
					{"delay": 10},
				},
			},
		},
		Actions: []map[string]interface{}{
			{"delay": 10},
		},
	}

	if err := validateV1Endpoint("http://test", "get", endpoint, map[string]bool{"testService": true}, map[string]*Request{"testRequest": nil}); err != nil {
		t.Errorf("Error Detected Validating Endpoint: %s", err.Error())
	}
}

// TestValidateV1Endpoint2 ensures incorrectly formatted endpoint actions are raised as an error
func TestValidateV1Endpoint2(t *testing.T) {
	endpoint := &Endpoint{
		Params: &Parameters{
			Query: make(map[string]*ParamEntry),
			Path:  make(map[string]*ParamEntry),
		},
		Recieves: &Recieves{
			Headers: map[string]string{"Content-Type": "application/json"},
			Body:    map[string]string{"foo.bar": "integer"},
		},
		Responses: map[int]*Response{
			200: {
				StatusCode: 200,
				Body:       map[string]interface{}{"foo.bar": "string"},
				Headers:    nil,
				Weight:     100,
				Actions: []map[string]interface{}{
					{"delay": 10},
				},
			},
		},
		Actions: []map[string]interface{}{
			{"delay": "foobar"},
		},
	}

	if err := validateV1Endpoint("http://test", "get", endpoint, map[string]bool{"testService": true}, map[string]*Request{"testRequest": nil}); err == nil {
		t.Errorf("No Error Detected Whilst Validating")
	}
}

// TestValidateV1Endpoint3 ensures incorrectly formatted response actions are raised as an error
func TestValidateV1Endpoint3(t *testing.T) {
	endpoint := &Endpoint{
		Params: &Parameters{
			Query: make(map[string]*ParamEntry),
			Path:  make(map[string]*ParamEntry),
		},
		Recieves: &Recieves{
			Headers: map[string]string{"Content-Type": "application/json"},
			Body:    map[string]string{"foo.bar": "integer"},
		},
		Responses: map[int]*Response{
			200: {
				StatusCode: 200,
				Body:       map[string]interface{}{"foo.bar": "string"},
				Headers:    nil,
				Weight:     100,
				Actions: []map[string]interface{}{
					{"delay": "foobar"},
				},
			},
		},
	}

	if err := validateV1Endpoint("http://test", "get", endpoint, map[string]bool{"testService": true}, map[string]*Request{"testRequest": nil}); err == nil {
		t.Errorf("No Error Detected Whilst Validating")
	}
}

// TestValidateV1Endpoint4 ensures incorrectly weighted response totals are raised as an error
func TestValidateV1Endpoint4(t *testing.T) {
	endpoint := &Endpoint{
		Params: &Parameters{
			Query: make(map[string]*ParamEntry),
			Path:  make(map[string]*ParamEntry),
		},
		Recieves: &Recieves{
			Headers: map[string]string{"Content-Type": "application/json"},
			Body:    map[string]string{"foo.bar": "integer"},
		},
		Responses: map[int]*Response{
			200: {
				StatusCode: 200,
				Body:       map[string]interface{}{"foo.bar": "string"},
				Headers:    nil,
				Weight:     80,
				Actions:    nil,
			},
		},
	}

	if err := validateV1Endpoint("http://test", "get", endpoint, map[string]bool{"testService": true}, map[string]*Request{"testRequest": nil}); err == nil {
		t.Errorf("No Error Detected Whilst Validating")
	}
}

// TestValidateV1Endpoint5 ensures an incorrectly typed Recieves.Body element is raised as an error
func TestValidateV1Endpoint5(t *testing.T) {
	endpoint := &Endpoint{
		Params: &Parameters{
			Query: make(map[string]*ParamEntry),
			Path:  make(map[string]*ParamEntry),
		},
		Recieves: &Recieves{
			Headers: map[string]string{"Content-Type": "application/json"},
			Body:    map[string]string{"foo.bar": "unsupported type"},
		},
	}

	if err := validateV1Endpoint("http://test", "get", endpoint, map[string]bool{"testService": true}, map[string]*Request{"testRequest": nil}); err == nil {
		t.Errorf("No Error Detected Whilst Validating")
	}
}

// TestValidateV1Endpoint6 ensures when no 'Response' and no 'Responses' is defined, an error is raised
func TestValidateV1Endpoint6(t *testing.T) {
	endpoint := &Endpoint{
		Params: &Parameters{
			Query: make(map[string]*ParamEntry),
			Path:  make(map[string]*ParamEntry),
		},
		Recieves: &Recieves{
			Headers: map[string]string{"Content-Type": "application/json"},
			Body:    map[string]string{"foo.bar": "integer"},
		},
	}

	if err := validateV1Endpoint("http://test", "get", endpoint, map[string]bool{"testService": true}, map[string]*Request{"testRequest": nil}); err == nil {
		t.Errorf("No Error Detected Whilst Validating")
	}
}

// TestValidateV1Endpoint7 ensures when an invalid Path param is defined, it is raised as an error to the user
func TestValidateV1Endpoint7(t *testing.T) {
	endpoint := &Endpoint{
		Params: &Parameters{
			Query: nil,
			Path:  map[string]*ParamEntry{"test": {Required: true, Type: "unsupported"}},
		},
	}

	if err := validateV1Endpoint("http://test", "get", endpoint, map[string]bool{"testService": true}, map[string]*Request{"testRequest": nil}); err == nil {
		t.Errorf("No Error Detected Whilst Validating")
	}
}

// TestValidateV1Endpoint8 ensures when an invalid Query param is defined, it is raised as an error to the user
func TestValidateV1Endpoint8(t *testing.T) {
	endpoint := &Endpoint{
		Params: &Parameters{
			Query: map[string]*ParamEntry{"test": {Required: true, Type: "unsupported"}},
			Path:  nil,
		},
	}

	if err := validateV1Endpoint("http://test", "get", endpoint, map[string]bool{"testService": true}, map[string]*Request{"testRequest": nil}); err == nil {
		t.Errorf("No Error Detected Whilst Validating")
	}
}

// TestValidateV1Request1 ensures a correctly formatted Request object is returned as valid
func TestValidateV1Request1(t *testing.T) {
	request := &Request{
		URL:      "/test",
		Method:   "post",
		Protocol: "http",
		Headers:  map[string]string{"foo": "bar"},
		Body:     map[string]interface{}{"bar": "foo"},
		ExpectedResponse: &Response{
			StatusCode: 201,
			Headers:    map[string]string{"foo": "bar"},
			Body:       map[string]interface{}{"foo.bar": "boolean"},
		},
	}

	if err := validateV1Request("testRequest", request); err != nil {
		t.Errorf("Valid Request Incorrectly Identified As Invalid: %s", err.Error())
	}
}

// TestValidateV1Request2 ensures an incorrectly configured ExpectedResponse.Body is raised as an error
func TestValidateV1Request2(t *testing.T) {
	request := &Request{
		URL:      "/test",
		Method:   "post",
		Protocol: "http",
		Headers:  map[string]string{"foo": "bar"},
		Body:     map[string]interface{}{"bar": "foo"},
		ExpectedResponse: &Response{
			StatusCode: 201,
			Headers:    map[string]string{"foo": "bar"},
			Body:       map[string]interface{}{"foo.bar": "unsupported"},
		},
	}

	if err := validateV1Request("testRequest", request); err == nil {
		t.Errorf("Invalid Request Incorrectly Identified As Valid")
	}
}

// TestValidateV1Request3 ensures an invalid ExpectedResponse.StatusCode is raised as an error
func TestValidateV1Request3(t *testing.T) {
	request := &Request{
		URL:              "/test",
		Method:           "post",
		Protocol:         "http",
		Headers:          map[string]string{"foo": "bar"},
		Body:             map[string]interface{}{"bar": "foo"},
		ExpectedResponse: &Response{},
	}

	if err := validateV1Request("testRequest", request); err == nil {
		t.Errorf("Invalid Request Incorrectly Identified As Valid")
	}
}

// TestValidateV1Request4 ensures when a protocol is missing, it defaults to 'http'
func TestValidateV1Request4(t *testing.T) {
	request := &Request{
		URL:     "/test",
		Method:  "post",
		Headers: map[string]string{"foo": "bar"},
		Body:    map[string]interface{}{"bar": "foo"},
		ExpectedResponse: &Response{
			StatusCode: 201,
			Headers:    map[string]string{"foo": "bar"},
			Body:       map[string]interface{}{"foo.bar": "boolean"},
		},
	}

	if err := validateV1Request("testRequest", request); err != nil {
		t.Errorf("Valid Request Incorrectly Identified As Invalid: %s", err.Error())
	}
}

// TestValidateV1Request5 ensures an unsupported protocol is raised as an error
func TestValidateV1Request5(t *testing.T) {
	request := &Request{
		URL:      "/test",
		Method:   "post",
		Protocol: "unsupported",
	}

	if err := validateV1Request("testRequest", request); err == nil {
		t.Errorf("Invalid Request Incorrectly Identified As Valid")
	}
}

// TestValidateV1Request6 ensures an unsupported method is raised as an error
func TestValidateV1Request6(t *testing.T) {
	request := &Request{
		URL:    "/test",
		Method: "unsupported",
	}

	if err := validateV1Request("testRequest", request); err == nil {
		t.Errorf("Invalid Request Incorrectly Identified As Valid")
	}
}

// TestValidateV1Request7 ensures a request with no URL defined is raised as an error
func TestValidateV1Request7(t *testing.T) {
	request := &Request{}

	if err := validateV1Request("testRequest", request); err == nil {
		t.Errorf("Invalid Request Incorrectly Identified As Valid")
	}
}

// TestValidateV1Request8 ensures an empty name is raised as an error
func TestValidateV1Request8(t *testing.T) {
	if err := validateV1Request("", &Request{}); err == nil {
		t.Errorf("Invalid Request Incorrectly Identified As Valid")
	}
}

// TestValidateV1Request9 ensures a nil-value request is raised as an error
func TestValidateV1Request9(t *testing.T) {
	if err := validateV1Request("", nil); err == nil {
		t.Errorf("Invalid Request Incorrectly Identified As Valid")
	}
}
