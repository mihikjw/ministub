package api

import (
	"encoding/json"
	"fmt"
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
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(err.StatusCode())
		if data, err := json.Marshal(err); err == nil {
			w.Write(data)
		}
		return
	}

	// evaluate query parameters
	if len(entry.Params.Query) > 0 {
		api.log.Info("Query Parameter Evaluation Not Written")
	}

	// evaluate body
	if entry.Recieves != nil {
		api.log.Info("Body Evaluation Not Written")
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
							switch {
							case pe.Type == "boolean":
								lowerIncomingBlock := strings.ToLower(incomingBlock)
								if lowerIncomingBlock != "true" && lowerIncomingBlock != "false" {
									return nil, &HTTPError{"Path Param Not Valid Boolean Value", http.StatusBadRequest}
								}
							case pe.Type == "integer":
								if _, err := strconv.Atoi(incomingBlock); err != nil {
									return nil, &HTTPError{"Path Param Not Valid Integer Value", http.StatusBadRequest}
								}
							case pe.Type == "string":
								if len(incomingBlock) == 0 {
									return nil, &HTTPError{"Path Param Not Valid String Value", http.StatusBadRequest}
								}
							case pe.Type == "float":
								if _, err := strconv.ParseFloat(incomingBlock, 64); err != nil {
									return nil, &HTTPError{"Path Param Not Valid Float Value", http.StatusBadRequest}
								}
							default:
								return nil, &HTTPError{"Path Param Type Not Supported", http.StatusBadRequest}
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
