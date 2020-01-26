// package respond provides low touch API for sending HTTP responses.
package respond

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// With sends the headers with the provided status then writes the data on the
// provided http.ResponseWriter.
// If provided http.ResponseWriter errors while writing the response then a
// non-nil error is returned wrapping the original error.
func With(w http.ResponseWriter, status int, data []byte) error {
	w.WriteHeader(status)
	if _, err := w.Write(data); err != nil {
		return fmt.Errorf("write: %w", err)
	}
	return nil
}

// WithJSON also sends the headers with the provided status then writes the data
// on the provided http.ResponseWriter after marshalling the data into json.
// It also overwrites the the Content-Type header to application/json, if the
// provided data is non-nil.
// If marshalling of the provided data fails or the http.ResponseWriter errors
// while writing the response then a non-nil error is returned wrapping the
// original error.
func WithJSON(w http.ResponseWriter, status int, data interface{}) error {
	if data == nil {
		return With(w, status, nil)
	}
	b, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("marshal: %w", err)
	}
	w.Header().Add("Content-Type", "application/json; utf-8")
	return With(w, status, b)
}
