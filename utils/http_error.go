package utils

import (
	"fmt"
	"net/http"
)

// HTTPError is used to communicate http status code together with error message
// from the service layer to the web layer.
type HTTPError struct {
	// The http status code
	HTTPStatusCode int
	// The error message that should be returned in the body
	Error error
}

// Create a new HTTP error with the given status code and message and format parameters.
func NewHTTPErrorf(status int, msg string, args ...string) *HTTPError {
	return &HTTPError{
		HTTPStatusCode: status,
		Error: fmt.Errorf(msg, args),
	}
}

// Write the HTTPError to an http response.
func (e HTTPError) WriteTo(w http.ResponseWriter) {
	fmt.Fprintln(w, e.Error)
	w.WriteHeader(e.HTTPStatusCode)
}