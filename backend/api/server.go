// backend/api/server.go
package api

import (
	"encoding/json"
	"net/http"
)

type Server struct {
	Router *http.ServeMux
}

func NewServer() *Server {
	s := &Server{Router: http.NewServeMux()}
	s.Router.HandleFunc("/api/status", s.handleStatus)
	return s
}

func (s *Server) handleStatus(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}
