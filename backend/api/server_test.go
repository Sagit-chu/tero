// backend/api/server_test.go
package api

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestStatusEndpoint(t *testing.T) {
	srv := NewServer()
	req, _ := http.NewRequest("GET", "/api/status", nil)
	rr := httptest.NewRecorder()
	srv.Router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
}
