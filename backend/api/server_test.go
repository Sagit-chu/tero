// backend/api/server_test.go
package api

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestStatusEndpoint(t *testing.T) {
	tests := []struct {
		name           string
		method         string
		expectedStatus int
		expectedBody   string
	}{
		{"Valid GET request", "GET", http.StatusOK, `{"node_status":"Alive","status":"ok"}`},
		{"Invalid POST request", "POST", http.StatusMethodNotAllowed, "Method Not Allowed\n"},
	}

	srv := NewServer(nil, nil)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest(tt.method, "/api/status", nil)
			rr := httptest.NewRecorder()
			srv.Router.ServeHTTP(rr, req)

			if rr.Code != tt.expectedStatus {
				t.Errorf("expected status %v, got %v", tt.expectedStatus, rr.Code)
			}
			// For the valid GET request, the body includes a newline.
			// Compare using TrimSpace to avoid newline issues, or exact match if needed.
			actualBody := rr.Body.String()
			if tt.expectedStatus == http.StatusOK {
				actualBody = actualBody[:len(actualBody)-1] // remove newline
			}
			if actualBody != tt.expectedBody {
				t.Errorf("expected body %q, got %q", tt.expectedBody, actualBody)
			}
		})
	}
}
