// backend/api/server.go
package api

import (
	"encoding/json"
	"net/http"

	"github.com/sagit-chu/flvx-monitor/backend/repository"
)

type StatusProvider interface {
	GetLastStatus() string
}

type Server struct {
	Router     *http.ServeMux
	NodeRepo   *repository.NodeRepository
	ConfigRepo *repository.ConfigRepository
	StatusProv StatusProvider
}

func NewServer(nodeRepo *repository.NodeRepository, configRepo *repository.ConfigRepository, statusProv StatusProvider) *Server {
	s := &Server{
		Router:     http.NewServeMux(),
		NodeRepo:   nodeRepo,
		ConfigRepo: configRepo,
		StatusProv: statusProv,
	}
	s.Router.HandleFunc("/api/status", s.handleStatus)
	s.Router.HandleFunc("/api/nodes", s.handleNodes)
	s.Router.HandleFunc("/api/config", s.handleConfig)
	return s
}

func (s *Server) handleStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	nodeStatus := "Unknown"
	if s.StatusProv != nil {
		nodeStatus = s.StatusProv.GetLastStatus()
	}

	w.Header().Set("Content-Type", "application/json")
	body, err := json.Marshal(struct {
		Status     string `json:"status"`
		NodeStatus string `json:"node_status"`
	}{
		Status:     "ok",
		NodeStatus: nodeStatus,
	})
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	w.Write(body)
}
