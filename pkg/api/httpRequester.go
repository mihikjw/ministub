package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/MichaelWittgreffe/ministub/pkg/config"
)

// HTTPRequester sends HTTP/1.1 requests
type HTTPRequester struct {
	client *http.Client
	mutex  sync.Mutex
}

// NewHTTPRequester is a constructor for HTTPRequester
func NewHTTPRequester() *HTTPRequester {
	return &HTTPRequester{
		client: &http.Client{Timeout: time.Second * 5},
	}
}

// Request sends an HTTP/1.1 request defined in 'req' to the service defined in 'tgt'
func (h *HTTPRequester) Request(tgt *config.Service, req *config.Request) error {
	if tgt == nil || req == nil {
		return fmt.Errorf("Invalid Args")
	}

	request, err := h.setupRequest(tgt, req)
	if err != nil {
		return err
	}

	resp, err := h.makeRequest(request)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if err = h.validateResponse(resp, req); err != nil {
		return fmt.Errorf("Response Validation Error: %s", err.Error())
	}

	return nil
}

// setupRequest configured an http.Request object according to the specified config input
func (h *HTTPRequester) setupRequest(tgt *config.Service, req *config.Request) (*http.Request, error) {
	var body *bytes.Buffer
	url := fmt.Sprintf("%s://%s:%d%s", req.Protocol, tgt.Hostname, tgt.Port, req.URL)

	if req.Body != nil && len(req.Body) > 0 {
		if jsonBody, err := json.Marshal(req.Body); err == nil {
			body = bytes.NewBuffer(jsonBody)
		} else {
			return nil, fmt.Errorf("Error Marshalling Request Body: %s", err.Error())
		}
	}

	request, err := http.NewRequest(strings.ToUpper(req.Method), url, body)
	if err != nil {
		return nil, fmt.Errorf("http.NewRequest: %s", err.Error())
	}

	if req.Headers != nil && len(req.Headers) > 0 {
		for headerKey, headerVal := range req.Headers {
			request.Header.Add(headerKey, headerVal)
		}
	}

	return request, nil
}

// makeRequest makes the given http.Request and returns the response & error, is goroutine-safe
func (h *HTTPRequester) makeRequest(req *http.Request) (resp *http.Response, err error) {
	h.mutex.Lock()
	defer h.mutex.Unlock()
	resp, err = h.client.Do(req)
	return resp, err
}

// validateResponse ensures the givne http.Response is correct according to the specified request
func (h *HTTPRequester) validateResponse(resp *http.Response, req *config.Request) error {
	if resp.StatusCode != req.ExpectedResponse.StatusCode {
		return fmt.Errorf("Status Code MisMatch, Expected %d Got %d", resp.StatusCode, req.ExpectedResponse.StatusCode)
	}

	// validate response body if required
	if req.ExpectedResponse.Body != nil && len(req.ExpectedResponse.Body) > 0 {
		if resp.Body == nil {
			return fmt.Errorf("Response Body Expected, None Recieved")
		}

		body := make(map[string]interface{})
		if respBody, err := ioutil.ReadAll(resp.Body); err == nil {
			if err = json.Unmarshal(respBody, &body); err != nil {
				return fmt.Errorf("Unable To Unmarshal Response Body: %s", err.Error())
			}
		} else {
			return fmt.Errorf("Unable To Read Response Body: %s", err.Error())
		}

		if len(body) > 0 {
			for exName, exValue := range req.ExpectedResponse.Body {
				if err := AssertValidTypeFromPath(exName, exValue.(string), body); err != nil {
					return fmt.Errorf("Invalid Expected Body Field: %s", err.Error())
				}
			}
		}
	}

	// validate the headers
	for exHeadKey, exHeadVal := range req.Headers {
		if inHeadVal := resp.Header.Get(exHeadKey); exHeadVal != inHeadVal {
			return fmt.Errorf("Expected Header Field %s Does Not Equal Actual Response Value - Expected: %s; Got: %s", exHeadKey, exHeadVal, inHeadVal)
		}
	}

	return nil
}
