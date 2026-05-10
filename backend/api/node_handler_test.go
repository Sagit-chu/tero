package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"github.com/sagit-chu/flvx-monitor/backend/repository"
)

func TestNodeHandler_Get(t *testing.T) {
	srv, cleanup := setupTestServer(t)
	defer cleanup()

	srv.NodeRepo.AddNode("1.1.1.1", "22", "pass")

	req, _ := http.NewRequest(http.MethodGet, "/api/nodes", nil)
	rr := httptest.NewRecorder()
	srv.Router.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status %v, got %v", http.StatusOK, rr.Code)
	}

	var resp []repository.Node
	json.NewDecoder(rr.Body).Decode(&resp)

	if len(resp) != 1 {
		t.Errorf("expected 1 node, got %d", len(resp))
	} else if resp[0].IP != "1.1.1.1" {
		t.Errorf("expected node IP 1.1.1.1, got %s", resp[0].IP)
	}
}

func TestNodeHandler_Post(t *testing.T) {
	srv, cleanup := setupTestServer(t)
	defer cleanup()

	payload := []byte(`{"ip": "2.2.2.2", "ssh_port": "22", "ssh_password": "sec"}`)
	req, _ := http.NewRequest(http.MethodPost, "/api/nodes", bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	srv.Router.ServeHTTP(rr, req)

	if rr.Code != http.StatusCreated {
		t.Errorf("expected status %v, got %v", http.StatusCreated, rr.Code)
	}

	nodes, _ := srv.NodeRepo.GetAllNodes()
	if len(nodes) != 1 || nodes[0].IP != "2.2.2.2" {
		t.Errorf("node was not added properly")
	}
}

func TestNodeHandler_InvalidMethod(t *testing.T) {
	srv, cleanup := setupTestServer(t)
	defer cleanup()

	req, _ := http.NewRequest(http.MethodPut, "/api/nodes", nil)
	rr := httptest.NewRecorder()
	srv.Router.ServeHTTP(rr, req)

	if rr.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected status %v, got %v", http.StatusMethodNotAllowed, rr.Code)
	}
}
