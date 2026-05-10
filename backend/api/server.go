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
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	body, err := json.Marshal(map[string]string{"status": "ok"})
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	w.Write(body)
	w.Write([]byte("\n"))
}
