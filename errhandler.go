package errhandler

import (
	"encoding/json"
	"net/http"
)

// Wrap an http.HandlerFunc, allowing for errors to be handled in a more
// familiar way inside your handlers.
type Wrap func(w http.ResponseWriter, r *http.Request) error

// ServeHTTP is invoked when the HTTP handler is called, capturing and returning
// any errors encountered in the ErrHandlerFunc.
func (fn Wrap) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := fn(w, r); err != nil {
		if cerr, ok := err.(HTTPError); ok {
			http.Error(w, cerr.Error(), cerr.status)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

// SendJSON returns a JSON response to the caller, or fails with an error.
func SendJSON(w http.ResponseWriter, data interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		return err
	}
	return nil
}

// SendString returns a string response to the caller, or fails with an error.
func SendString(w http.ResponseWriter, message string) error {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte(message)); err != nil {
		return err
	}
	return nil
}

// ParseJSON can be used to parse a string from the body of a request.
func ParseJSON(r *http.Request, val any) error {
	return json.NewDecoder(r.Body).Decode(val)
}
