// backend/api/server_test.go
package api

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestStatusEndpoint(t *testing.T) {
	tests := []struct {
		name           string
		method         string
		expectedStatus int
		expectedBody   string
	}{
		{"Valid GET request", "GET", http.StatusOK, `{"status":"ok","node_status":"Unknown"}`},
		{"Invalid POST request", "POST", http.StatusMethodNotAllowed, "Method Not Allowed\n"},
	}

	srv := NewServer(nil, nil, nil)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest(tt.method, "/api/status", nil)
			rr := httptest.NewRecorder()
			srv.Router.ServeHTTP(rr, req)

			if rr.Code != tt.expectedStatus {
				t.Errorf("expected status %v, got %v", tt.expectedStatus, rr.Code)
			}
			actualBody := strings.TrimSpace(rr.Body.String())
			expectedBody := strings.TrimSpace(tt.expectedBody)
			if actualBody != expectedBody {
				t.Errorf("expected body %q, got %q", expectedBody, actualBody)
			}
		})
	}
}
