package monitor

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPingPeScraper_CheckBlocked(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		// Mock response that simulates 100% loss from China
		w.Write([]byte(`{"china_loss": 100, "overseas_loss": 0}`))
	}))
	defer ts.Close()

	scraper := NewPingPeScraper(ts.URL)
	blocked, err := scraper.IsBlocked("1.1.1.1")
	
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if !blocked {
		t.Errorf("Expected node to be detected as blocked")
	}
}
