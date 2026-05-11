package flvx

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestFlvxClient_ReplaceNode(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("Expected POST, got %s", r.Method)
		}
		user, pass, ok := r.BasicAuth()
		if !ok || user != "admin" || pass != "secret" {
			t.Errorf("Invalid auth")
		}
		var payload map[string]string
		json.NewDecoder(r.Body).Decode(&payload)
		if payload["ip"] != "1.2.3.4" || payload["ssh_port"] != "22" || payload["ssh_password"] != "pass" {
			t.Errorf("Invalid payload")
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	client := NewClient(ts.URL, "admin", "secret")
	err := client.ReplaceNode("1.2.3.4", "22", "pass")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
}
