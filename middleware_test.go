package errhandler

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

type contextKey string

func TestMiddleware(t *testing.T) {
	tests := []struct {
		name           string
		handler        Wrap
		middleware     Middleware
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "single middleware",
			middleware: func(n Wrap) Wrap {
				return func(w http.ResponseWriter, r *http.Request) error {
					ctx := context.WithValue(r.Context(), contextKey("value"), 1)
					r = r.WithContext(ctx)

					return n(w, r)
				}
			},
			handler: Wrap(func(w http.ResponseWriter, r *http.Request) error {
				m := map[string]any{
					"value": r.Context().Value(contextKey("value")),
				}

				return SendJSON(w, m)
			}),
			expectedStatus: http.StatusOK,
			expectedBody:   "{\"value\":1}\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "http://example.com/foo", nil)
			w := httptest.NewRecorder()

			tt.middleware(tt.handler).ServeHTTP(w, req)

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

func TestChain(t *testing.T) {
	tests := []struct {
		name           string
		handler        Wrap
		middlewares    []Middleware
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "multiple middlewares",
			middlewares: []Middleware{
				func(n Wrap) Wrap {
					return func(w http.ResponseWriter, r *http.Request) error {
						ctx := context.WithValue(r.Context(), contextKey("value"), 1)
						r = r.WithContext(ctx)

						return n(w, r)
					}
				},
				func(n Wrap) Wrap {
					return func(w http.ResponseWriter, r *http.Request) error {
						value, ok := r.Context().Value(contextKey("value")).(int)
						if !ok {
							t.Fatalf("expected integer but didn't get one")
						}

						ctx := context.WithValue(r.Context(), contextKey("value"), value+2)
						r = r.WithContext(ctx)

						return n(w, r)
					}
				},
			},
			handler: Wrap(func(w http.ResponseWriter, r *http.Request) error {
				m := map[string]any{
					"value": r.Context().Value(contextKey("value")),
				}

				return SendJSON(w, m)
			}),
			expectedStatus: http.StatusOK,
			expectedBody:   "{\"value\":3}\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "http://example.com/foo", nil)
			w := httptest.NewRecorder()

			chain := Chain(tt.middlewares...)
			chain(tt.handler).ServeHTTP(w, req)

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
