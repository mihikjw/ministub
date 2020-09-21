package api

import (
	"github.com/MichaelWittgreffe/ministub/pkg/config"
)

// Requester represents an object able to make external requests
type Requester interface {
	// Request makes the given request to the given service
	Request(tgt *config.Service, req *config.Request) error
}

// NewRequester is a factory function for Requester objects
func NewRequester(reqMode string) Requester {
	switch {
	case reqMode == "http":
		return NewHTTPRequester()
	default:
		return nil
	}
}
