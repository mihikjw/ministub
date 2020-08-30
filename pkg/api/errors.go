package api

type error interface {
	Error() string
}

// HTTPError represents an error with an incoming API request
type HTTPError struct {
	Err        string `json:"error"`
	statusCode int
}

// Error returns the error string
func (me *HTTPError) Error() string { return me.Err }

// StatusCode returns the status code for this error
func (me *HTTPError) StatusCode() int { return me.statusCode }
