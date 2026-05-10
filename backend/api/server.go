// backend/api/server.go
package api

import (
	"encoding/json"
	"net/http"
	"github.com/sagit-chu/flvx-monitor/backend/repository"
)

type Server struct {
	Router     *http.ServeMux
	NodeRepo   *repository.NodeRepository
	ConfigRepo *repository.ConfigRepository
}

func NewServer(nodeRepo *repository.NodeRepository, configRepo *repository.ConfigRepository) *Server {
	s := &Server{
		Router:     http.NewServeMux(),
		NodeRepo:   nodeRepo,
		ConfigRepo: configRepo,
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
	w.Header().Set("Content-Type", "application/json")
	body, err := json.Marshal(map[string]string{
		"status":      "ok",
		"node_status": "Alive", // TODO: Wire to actual monitoring state
	})
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	w.Write(body)
	w.Write([]byte("\n"))
}
