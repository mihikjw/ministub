package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/MichaelWittgreffe/ministub/pkg/config"
)

// HTTPRequester sends HTTP/1.1 requests, is thread-safe
type HTTPRequester struct {
	client *http.Client
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

	url := fmt.Sprintf("http://%s:%d%s", tgt.Hostname, tgt.Port, req.URL)

	var body *bytes.Buffer
	if req.Body != nil && len(req.Body) > 0 {
		if jsonBody, err := json.Marshal(req.Body); err == nil {
			body = bytes.NewBuffer(jsonBody)
		} else {
			return fmt.Errorf("Error Marshalling Request Body: %s", err.Error())
		}
	}

	request, err := http.NewRequest(strings.ToUpper(req.Method), url, body)
	if err != nil {
		return fmt.Errorf("http.NewRequest: %s", err.Error())
	}

	if req.Headers != nil && len(req.Headers) > 0 {
		for headerKey, headerVal := range req.Headers {
			request.Header.Add(headerKey, headerVal)
		}
	}

	resp, err := h.client.Do(request)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != req.ExpectedResponse.StatusCode {
		return fmt.Errorf("Status Code MisMatch, Expected %d Got %d", resp.StatusCode, req.ExpectedResponse.StatusCode)
	}

	// validate response body
	if resp.Body != nil {
		// respBody, err := ioutil.ReadAll(resp.Body)
	}
	// validate the headers

	return nil
}
