package errhandler

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// TestWrap_ServeHTTP tests the ServeHTTP method of the Wrap type.
func TestWrap_ServeHTTP(t *testing.T) {
	tests := []struct {
		name           string
		body           string
		handler        Wrap
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "parse json without error",
			body: `{"a": 1}`,
			handler: Wrap(func(w http.ResponseWriter, r *http.Request) error {
				m := map[string]int{}

				return ParseJSON(r, &m)
			}),
			expectedStatus: http.StatusOK,
			expectedBody:   "",
		},
		{
			name: "parse json with error",
			body: "a",
			handler: Wrap(func(w http.ResponseWriter, r *http.Request) error {
				m := map[string]int{}

				return ParseJSON(r, &m)
			}),
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   "invalid character 'a' looking for beginning of value\n",
		},
		{
			name: "send string without error",
			handler: Wrap(func(w http.ResponseWriter, r *http.Request) error {
				return SendString(w, "test")
			}),
			expectedStatus: http.StatusOK,
			expectedBody:   "test",
		},
		{
			name: "send json without error",
			handler: Wrap(func(w http.ResponseWriter, r *http.Request) error {
				return SendJSON(w, map[string]int{"a": 1})
			}),
			expectedStatus: http.StatusOK,
			expectedBody:   "{\"a\":1}\n",
		},
		{
			name: "send error",
			handler: Wrap(func(w http.ResponseWriter, r *http.Request) error {
				return SendError(w, http.StatusUnprocessableEntity, fmt.Errorf("error doing stuff: %w", errors.New("database error")))
			}),
			expectedStatus: http.StatusUnprocessableEntity,
			expectedBody:   "error doing stuff: database error",
		},
		{
			name: "error",
			handler: Wrap(func(w http.ResponseWriter, r *http.Request) error {
				return fmt.Errorf("error doing stuff: %w", errors.New("database error"))
			}),
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   "error doing stuff: database error\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var body io.Reader
			if tt.body != "" {
				body = strings.NewReader(tt.body)
			}

			req := httptest.NewRequest(http.MethodGet, "http://example.com/foo", body)
			w := httptest.NewRecorder()

			tt.handler.ServeHTTP(w, req)

			res := w.Result()
			defer res.Body.Close()

			if res.StatusCode != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, res.StatusCode)
			}

			actualBody := w.Body.String()
			if actualBody != tt.expectedBody {
				t.Errorf("expected body %q, got %q", tt.expectedBody, actualBody)
			}
		})
	}
}
