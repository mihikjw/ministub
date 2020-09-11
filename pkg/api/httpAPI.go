package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/MichaelWittgreffe/ministub/pkg/config"
	"github.com/MichaelWittgreffe/ministub/pkg/logger"
)

// HTTPAPI represents the HTTP API
type HTTPAPI struct {
	log logger.Logger
	cfg *config.Config
}

// NewHTTPAPI creates a new instance of HTTPAPI
func NewHTTPAPI(log logger.Logger, cfg *config.Config) *HTTPAPI {
	if log == nil || cfg == nil {
		return nil
	}
	api := &HTTPAPI{log: log, cfg: cfg}
	http.HandleFunc("/", api.requestHandler)
	return api
}

// ListenAndServe begins the API listening for requests
func (api *HTTPAPI) ListenAndServe(addressBind string, port int) error {
	return http.ListenAndServe(fmt.Sprintf("%s:%d", addressBind, port), nil)
}

// requestHandler is a handler for all incoming requests
func (api *HTTPAPI) requestHandler(w http.ResponseWriter, r *http.Request) {
	entry, err := api.getEndpointEntry(r)
	if err != nil {
		api.setupErrorResponse(err, w)
		return
	}

	// evaluate query parameters
	if len(entry.Params.Query) > 0 {
		if err = api.evaluateQueryParams(entry, r); err != nil {
			api.setupErrorResponse(err, w)
			return
		}
	}

	if entry.Recieves != nil {
		// evaluate headers
		if len(entry.Recieves.Headers) > 0 {
			if err := api.evaluateHeaders(entry.Recieves, r); err != nil {
				api.setupErrorResponse(err, w)
				return
			}
		}

		// evaluate body
		if len(entry.Recieves.Body) > 0 {
			if err := api.evaluateBody(entry.Recieves, r); err != nil {
				api.setupErrorResponse(err, w)
				return
			}
		}
	}

	// start actions goroutine
	if len(entry.Actions) > 0 {
		api.log.Info("Request Actions Evaluation Not Written")
	}

	// setup return value
	if entry.Response != 0 {
		w.WriteHeader(entry.Response)
		return
	}

	if len(entry.Responses) > 0 {
		api.log.Info("Responses Evaluation Not Written")
	}
}

// getEndpointEntry returns the Endpoint object for an incoming request, if it cannot be found immediatly we check all of them for parameter matching
func (api *HTTPAPI) getEndpointEntry(r *http.Request) (*config.Endpoint, *HTTPError) {
	urlEntry, found := api.cfg.Endpoints[r.URL.Path]
	if found {
		if entry, found := urlEntry[strings.ToLower(r.Method)]; found {
			return entry, nil
		}
		return nil, &HTTPError{"Method For URL Not Found", http.StatusMethodNotAllowed}
	}

	splitIncomingURL := strings.Split(r.URL.Path, "/")

	for url, data := range api.cfg.Endpoints {
		splitURL := strings.Split(url, "/")
		splitURLLen := len(splitURL)

		if len(splitIncomingURL) == splitURLLen {
			for i, incomingBlock := range splitIncomingURL {
				staticBlock := splitURL[i]

				// its the same and we're not the last item, skip to the next block
				if staticBlock == incomingBlock && (i+1) != splitURLLen {
					continue
				}

				// its a parameter point, get the value and ensure its the correct type
				if string(staticBlock[0]) == ":" {
					if endpoint, found := data[strings.ToLower(r.Method)]; found {
						if pe, found := endpoint.Params.Path[staticBlock[1:]]; found {
							if !api.assertValidType(interface{}(incomingBlock), pe.Type) {
								return nil, &HTTPError{fmt.Sprintf("Path Param Not Valid %s Value", pe.Type), http.StatusBadRequest}
							}
						}
					} else {
						return nil, &HTTPError{"Method For URL Not Found", http.StatusMethodNotAllowed}
					}
				}

				// we're on the last item and all fields have been correct
				if (i + 1) == splitURLLen {
					if entry, found := data[strings.ToLower(r.Method)]; found {
						return entry, nil
					}
					return nil, &HTTPError{"Method For URL Not Found", http.StatusMethodNotAllowed}
				}
			}
		}
	}

	return nil, &HTTPError{"URL Not Found", http.StatusNotFound}
}

// evaluateQueryParams ensures the incoming query params are compatible with the config definition for this endpoint
func (api *HTTPAPI) evaluateQueryParams(entry *config.Endpoint, r *http.Request) *HTTPError {
	inQueryValues := r.URL.Query()

	for expectedParamName, expectedParamEntry := range entry.Params.Query {
		inParamValue := inQueryValues.Get(expectedParamName)

		if len(inParamValue) == 0 && expectedParamEntry.Required {
			return &HTTPError{fmt.Sprintf("Missing Query Parameter: %s", expectedParamName), http.StatusBadRequest}
		}

		if len(inParamValue) > 0 && !api.assertValidType(inParamValue, expectedParamEntry.Type) {
			return &HTTPError{fmt.Sprintf("Query Param Not Valid %s Value", expectedParamEntry.Type), http.StatusBadRequest}
		}
	}

	return nil
}

// setupErrorResponse generates a standard error output ready for immediate return
func (api *HTTPAPI) setupErrorResponse(err *HTTPError, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(err.StatusCode())
	if data, err := json.Marshal(err); err == nil {
		w.Write(data)
	}
}

/* assertValidType returns whether a given value is of the expected type
if the initial conversion fails, it will convert to string then convert to type where appropriate */
func (api *HTTPAPI) assertValidType(value interface{}, expectedType string) bool {
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

// evaluateHeaders checks the request headers are valid
func (api *HTTPAPI) evaluateHeaders(in *config.Recieves, r *http.Request) *HTTPError {
	for exHeaderKey, exHeaderValue := range in.Headers {
		if inHeaderValue := r.Header.Get(exHeaderKey); inHeaderValue != exHeaderValue {
			return &HTTPError{fmt.Sprintf("Header Value %s Not Found", exHeaderKey), http.StatusBadRequest}
		}
	}
	return nil
}

// evaluateBody checks the request body is valid
func (api *HTTPAPI) evaluateBody(in *config.Recieves, r *http.Request) *HTTPError {
	var body map[string]interface{}
	rawBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return &HTTPError{"Error Reading Incoming Body", http.StatusInternalServerError}
	}
	if err := json.Unmarshal(rawBody, &body); err != nil {
		return &HTTPError{"Error Decoding Incoming Body", http.StatusInternalServerError}
	}

	for exName, exType := range in.Body {
		if err = api.assertValidTypeFromPath(exName, exType, body); err != nil {
			return &HTTPError{err.Error(), http.StatusBadRequest}
		}
	}

	return nil
}

// assertValidTypeFromPath goes down a given path for a given inputBody JSON, and asserts the expected valid type when at the expected level for the given path
func (api *HTTPAPI) assertValidTypeFromPath(path, expectedType string, inputData interface{}) error {
	if splitPath := strings.Split(path, "."); len(splitPath) >= 1 {
		if len(splitPath) > 1 {
			if value, err := strconv.Atoi(splitPath[0]); err == nil {
				nextInputData := inputData.([]interface{})[value]
				return api.assertValidTypeFromPath(strings.Join(splitPath[1:], "."), expectedType, nextInputData)
			} else if body, valid := inputData.(map[string]interface{}); valid {
				nextInputData := body[splitPath[0]]
				return api.assertValidTypeFromPath(strings.Join(splitPath[1:], "."), expectedType, nextInputData)
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

		if !api.assertValidType(valueToEvaluate, expectedType) {
			return fmt.Errorf("%s Is Invalid Expected Type %s", path, expectedType)
		}

		return nil
	}

	return fmt.Errorf("Empty Path Given For assertValidTypeFromPath")
}
