package api

import (
	"encoding/json"
	"net/http"
	"github.com/sagit-chu/flvx-monitor/backend/repository"
)

func (s *Server) handleNodes(w http.ResponseWriter, r *http.Request) {
	if s.NodeRepo == nil {
		http.Error(w, "Repository not initialized", http.StatusInternalServerError)
		return
	}

	switch r.Method {
	case http.MethodGet:
		nodes, err := s.NodeRepo.GetAllNodes()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if nodes == nil {
			nodes = make([]repository.Node, 0)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(nodes)
	case http.MethodPost:
		var req struct {
			IP          string `json:"ip"`
			SSHPort     string `json:"ssh_port"`
			SSHPassword string `json:"ssh_password"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}
		if req.IP == "" || req.SSHPort == "" || req.SSHPassword == "" {
			http.Error(w, "Missing required fields", http.StatusBadRequest)
			return
		}
		if err := s.NodeRepo.AddNode(req.IP, req.SSHPort, req.SSHPassword); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}
