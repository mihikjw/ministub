package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"strings"

	"github.com/MichaelWittgreffe/ministub/pkg/config"
	"github.com/MichaelWittgreffe/ministub/pkg/logger"
)

// HTTPAPI represents the HTTP API
type HTTPAPI struct {
	log   logger.Logger
	cfg   *config.Config
	stats map[string]map[int]int // url -> statusCode: count
}

// NewHTTPAPI creates a new instance of HTTPAPI
func NewHTTPAPI(log logger.Logger, cfg *config.Config) *HTTPAPI {
	if log == nil || cfg == nil {
		return nil
	}
	api := &HTTPAPI{
		log:   log,
		cfg:   cfg,
		stats: make(map[string]map[int]int),
	}
	http.HandleFunc("/", api.requestHandler)
	return api
}

// ListenAndServe begins the API listening for requests
func (api *HTTPAPI) ListenAndServe(addressBind string, port int) error {
	api.log.Info(fmt.Sprintf("Beginning Listening For HTTP Requests On %s:%d", addressBind, port))
	return http.ListenAndServe(fmt.Sprintf("%s:%d", addressBind, port), nil)
}

// requestHandler is a handler for all incoming requests
func (api *HTTPAPI) requestHandler(w http.ResponseWriter, r *http.Request) {
	// check own endpoints first
	switch {
	case r.URL.Path == "/stats":
		api.statsHandler(w)
		api.log.Info(fmt.Sprintf("%s | %s | %d", r.Host, r.URL.Path, http.StatusOK))
		return
	case r.URL.Path == "/exit":
		api.exitHandler()
	}

	// get the entry for the incoming request
	url, entry, err := api.getEndpointEntry(r)
	if err != nil {
		api.setupErrorResponse(err, w)
		api.log.Error(fmt.Sprintf("%s | %s | %d - %s", r.Host, r.URL.Path, err.StatusCode(), err.Error()))
		return
	}

	// get stats entry before any processing
	stats, found := api.stats[url]
	if !found {
		api.addEndpointToStats(url, entry.Responses)
		if stats, found = api.stats[url]; !found {
			api.setupErrorResponse(&HTTPError{fmt.Sprintf("Unable To Init Stats For Endpoint: %s", url), http.StatusInternalServerError}, w)
			api.log.Error(fmt.Sprintf("%s | %s | %d - %s", r.Host, r.URL.Path, err.StatusCode(), err.Error()))
			return
		}
	}

	// evaluate query parameters
	if len(entry.Params.Query) > 0 {
		if err = api.evaluateQueryParams(entry, r); err != nil {
			api.setupErrorResponse(err, w)
			api.log.Error(fmt.Sprintf("%s | %s | %d - %s", r.Host, r.URL.Path, err.StatusCode(), err.Error()))
			return
		}
	}

	if entry.Recieves != nil {
		// evaluate headers
		if len(entry.Recieves.Headers) > 0 {
			if err := api.evaluateHeaders(entry.Recieves, r); err != nil {
				api.setupErrorResponse(err, w)
				api.log.Error(fmt.Sprintf("%s | %s | %d - %s", r.Host, r.URL.Path, err.StatusCode(), err.Error()))
				return
			}
		}

		// evaluate body
		if len(entry.Recieves.Body) > 0 {
			if err := api.evaluateBody(entry.Recieves, r); err != nil {
				api.setupErrorResponse(err, w)
				api.log.Error(fmt.Sprintf("%s | %s | %d - %s", r.Host, r.URL.Path, err.StatusCode(), err.Error()))
				return
			}
		}
	}

	// setup return value
	var statusCode int
	if entry.Response > 0 {
		statusCode = entry.Response
		w.WriteHeader(statusCode)
	} else if len(entry.Responses) > 0 {
		statusCode = api.setupResponse(url, entry.Responses, w)
	}

	// increment stats
	stats[statusCode]++
	api.stats[url] = stats

	// start actions
	if len(entry.Actions) > 0 {
		go ExecuteActions(entry.Actions, r.URL.Path, api.cfg, api.log)
	}
	if len(entry.Responses[statusCode].Actions) > 0 {
		go ExecuteActions(entry.Responses[statusCode].Actions, r.URL.Path, api.cfg, api.log)
	}

	api.log.Info(fmt.Sprintf("%s | %s | %d", r.Host, r.URL.Path, statusCode))
}

// getEndpointEntry returns the Endpoint object for an incoming request, if it cannot be found immediatly we check all of them for parameter matching
func (api *HTTPAPI) getEndpointEntry(r *http.Request) (string, *config.Endpoint, *HTTPError) {
	urlEntry, found := api.cfg.Endpoints[r.URL.Path]
	if found {
		if entry, found := urlEntry[strings.ToLower(r.Method)]; found {
			return r.URL.Path, entry, nil
		}
		return "", nil, &HTTPError{"Method For URL Not Found", http.StatusMethodNotAllowed}
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
							if !AssertValidType(interface{}(incomingBlock), pe.Type) {
								return "", nil, &HTTPError{fmt.Sprintf("Path Param Not Valid %s Value", pe.Type), http.StatusBadRequest}
							}
						}
					} else {
						return "", nil, &HTTPError{"Method For URL Not Found", http.StatusMethodNotAllowed}
					}
				}

				// we're on the last item and all fields have been correct
				if (i + 1) == splitURLLen {
					if entry, found := data[strings.ToLower(r.Method)]; found {
						return url, entry, nil
					}
					return "", nil, &HTTPError{"Method For URL Not Found", http.StatusMethodNotAllowed}
				}
			}
		}
	}

	return "", nil, &HTTPError{"URL Not Found", http.StatusNotFound}
}

// evaluateQueryParams ensures the incoming query params are compatible with the config definition for this endpoint
func (api *HTTPAPI) evaluateQueryParams(entry *config.Endpoint, r *http.Request) *HTTPError {
	inQueryValues := r.URL.Query()

	for expectedParamName, expectedParamEntry := range entry.Params.Query {
		inParamValue := inQueryValues.Get(expectedParamName)

		if len(inParamValue) == 0 && expectedParamEntry.Required {
			return &HTTPError{fmt.Sprintf("Missing Query Parameter: %s", expectedParamName), http.StatusBadRequest}
		}

		if len(inParamValue) > 0 && !AssertValidType(inParamValue, expectedParamEntry.Type) {
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
		if err = AssertValidTypeFromPath(exName, exType, body); err != nil {
			return &HTTPError{err.Error(), http.StatusBadRequest}
		}
	}

	return nil
}

// setupResponse sets up the response to the users request based on the loaded cfg
func (api *HTTPAPI) setupResponse(url string, responses map[int]*config.Response, w http.ResponseWriter) int {
	var statusCode int
	var resp *config.Response

	if len(responses) > 1 {
		// setup based on specified weighting
		store := make([]int, 10)
		var addCount int
		processCount := 0
		addIndex := 0

		for statusCode, respEntry := range responses {
			addCount = respEntry.Weight / 10
			processCount = 0

			for processCount < addCount {
				store[addIndex] = statusCode
				processCount++
				addIndex++
			}
		}

		statusCode = store[rand.Intn(len(store))]
		resp = responses[statusCode]
	} else {
		// single response defined, just set it up
		for foundStatusCode, foundResp := range responses {
			statusCode = foundStatusCode
			resp = foundResp
		}
	}

	if len(resp.Headers) > 0 {
		for headerName, headerVal := range resp.Headers {
			w.Header().Set(headerName, headerVal)
		}
	}

	w.WriteHeader(statusCode)

	if resp.Body != nil {
		if data, err := json.Marshal(resp.Body); err == nil {
			w.Write(data)
		} else {
			api.setupErrorResponse(&HTTPError{
				fmt.Sprintf("Unable To Write Response Body For Endpoint %s: %s", url, err.Error()),
				http.StatusInternalServerError,
			}, w)
		}
	}

	return statusCode
}

// addEndpointToStats adds the given url to the statistics with zero-values for all status codes
func (api *HTTPAPI) addEndpointToStats(url string, responses map[int]*config.Response) {
	stats := make(map[int]int, len(responses))

	for statusCode := range responses {
		stats[statusCode] = 0
	}

	api.stats[url] = stats
}

// statsHandler returns the current application stats as JSON
func (api *HTTPAPI) statsHandler(w http.ResponseWriter) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if data, err := json.Marshal(api.stats); err == nil {
		w.Write(data)
	} else {
		api.setupErrorResponse(&HTTPError{
			fmt.Sprintf("Unable To Write Response Body For Endpoint /stats: %s", err.Error()),
			http.StatusInternalServerError,
		}, w)
	}
}

// exitHandler quits the application
func (api *HTTPAPI) exitHandler() {
	api.log.Info("Exit Requested, Shutting Down...")
	os.Exit(0)
}
