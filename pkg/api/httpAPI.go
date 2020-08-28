package api

import (
	"fmt"
	"net/http"

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
	api := &HTTPAPI{log: log}
	http.HandleFunc("/", api.requestHandler)
	return api
}

// ListenAndServe begins the API listening for requests
func (api *HTTPAPI) ListenAndServe(addressBind string, port int) error {
	return http.ListenAndServe(fmt.Sprintf("%s:%d", addressBind, port), nil)
}

func (api *HTTPAPI) requestHandler(w http.ResponseWriter, r *http.Request) {
	switch {
	case r.Method == "GET":
		api.getHandler(&w, r)
	case r.Method == "PUT":
		api.putHandler(&w, r)
	case r.Method == "POST":
		api.postHandler(&w, r)
	case r.Method == "DELETE":
		api.deleteHandler(&w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

// getHandler handles all GET requests
func (api *HTTPAPI) getHandler(w *http.ResponseWriter, r *http.Request) {
	api.log.Info("GET")
}

// putHandler handles all PUT requests
func (api *HTTPAPI) putHandler(w *http.ResponseWriter, r *http.Request) {
	api.log.Info("PUT")
}

// postHandler handles all POST requests
func (api *HTTPAPI) postHandler(w *http.ResponseWriter, r *http.Request) {
	api.log.Info("POST")
}

// deleteHandler handles all DELETE requests
func (api *HTTPAPI) deleteHandler(w *http.ResponseWriter, r *http.Request) {
	api.log.Info("DELETE")
}
