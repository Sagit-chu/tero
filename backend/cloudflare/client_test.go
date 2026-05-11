package cloudflare

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCloudflareClient_UpdateDNSRecord(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if r.URL.Path == "/client/v4/zones" {
			w.Write([]byte(`{"success":true,"result":[{"id":"zone123"}]}`))
			return
		}
		if r.URL.Path == "/client/v4/zones/zone123/dns_records" && r.Method == http.MethodGet {
			w.Write([]byte(`{"success":true,"result":[{"id":"rec123"}]}`))
			return
		}
		if r.URL.Path == "/client/v4/zones/zone123/dns_records/rec123" && r.Method == http.MethodPut {
			w.Write([]byte(`{"success":true}`))
			return
		}
		w.WriteHeader(http.StatusNotFound)
	}))
	defer ts.Close()

	client := NewClient(ts.URL, "token123")
	err := client.UpdateDNSRecord("example.com", "1.1.1.1")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
}
