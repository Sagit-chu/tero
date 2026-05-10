package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"github.com/sagit-chu/flvx-monitor/backend/db"
	"github.com/sagit-chu/flvx-monitor/backend/repository"
)

func setupTestServer(t *testing.T) (*Server, func()) {
	database, err := db.InitDB(":memory:")
	if err != nil {
		t.Fatalf("Failed to init db: %v", err)
	}

	nodeRepo := repository.NewNodeRepository(database)
	configRepo := repository.NewConfigRepository(database)
	srv := NewServer(nodeRepo, configRepo)

	return srv, func() {
		database.Close()
	}
}

func TestConfigHandler_Get(t *testing.T) {
	srv, cleanup := setupTestServer(t)
	defer cleanup()

	// Seed some data
	srv.ConfigRepo.SetConfig("flvx_account", "test_acc")

	req, _ := http.NewRequest(http.MethodGet, "/api/config", nil)
	rr := httptest.NewRecorder()
	srv.Router.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status %v, got %v", http.StatusOK, rr.Code)
	}

	var resp map[string]string
	json.NewDecoder(rr.Body).Decode(&resp)

	if resp["flvx_account"] != "test_acc" {
		t.Errorf("expected flvx_account to be test_acc, got %s", resp["flvx_account"])
	}
	if resp["domain_name"] != "" {
		t.Errorf("expected empty string for missing config, got %s", resp["domain_name"])
	}
}

func TestConfigHandler_Post(t *testing.T) {
	srv, cleanup := setupTestServer(t)
	defer cleanup()

	payload := []byte(`{"domain_name": "example.com", "cf_token": "token123"}`)
	req, _ := http.NewRequest(http.MethodPost, "/api/config", bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	srv.Router.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status %v, got %v", http.StatusOK, rr.Code)
	}

	val, _ := srv.ConfigRepo.GetConfig("domain_name")
	if val != "example.com" {
		t.Errorf("expected domain_name to be example.com, got %s", val)
	}
}

func TestConfigHandler_InvalidMethod(t *testing.T) {
	srv, cleanup := setupTestServer(t)
	defer cleanup()

	req, _ := http.NewRequest(http.MethodPut, "/api/config", nil)
	rr := httptest.NewRecorder()
	srv.Router.ServeHTTP(rr, req)

	if rr.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected status %v, got %v", http.StatusMethodNotAllowed, rr.Code)
	}
}
