package api

import (
	"encoding/json"
	"fmt"
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
			IP           string `json:"ip"`
			SSHPort      string `json:"ssh_port"`
			SSHPassword  string `json:"ssh_password"`
			FlvxNodeID   int    `json:"flvx_node_id"`
			FlvxNodeName string `json:"flvx_node_name"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}
		if req.IP == "" {
			http.Error(w, "Missing required fields", http.StatusBadRequest)
			return
		}
		if err := s.NodeRepo.AddNode(req.IP, req.SSHPort, req.SSHPassword, req.FlvxNodeID, req.FlvxNodeName); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
	case http.MethodPut:
		idStr := r.URL.Query().Get("id")
		if idStr == "" {
			http.Error(w, "Missing node id", http.StatusBadRequest)
			return
		}
		var id int
		fmt.Sscanf(idStr, "%d", &id)

		var req struct {
			IP           string `json:"ip"`
			SSHPort      string `json:"ssh_port"`
			SSHPassword  string `json:"ssh_password"`
			FlvxNodeID   int    `json:"flvx_node_id"`
			FlvxNodeName string `json:"flvx_node_name"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}
		if err := s.NodeRepo.UpdateNode(id, req.IP, req.SSHPort, req.SSHPassword, req.FlvxNodeID, req.FlvxNodeName); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)

	case http.MethodDelete:
		idStr := r.URL.Query().Get("id")
		if idStr == "" {
			http.Error(w, "Missing node id", http.StatusBadRequest)
			return
		}
		var id int
		fmt.Sscanf(idStr, "%d", &id)

		if err := s.NodeRepo.DeleteNode(id); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}