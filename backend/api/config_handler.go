package api

import (
	"database/sql"
	"encoding/json"
	"net/http"
)

func (s *Server) handleConfig(w http.ResponseWriter, r *http.Request) {
	if s.ConfigRepo == nil {
		http.Error(w, "Config Repository not initialized", http.StatusInternalServerError)
		return
	}

	switch r.Method {
	case http.MethodGet:
		keys := []string{"flvx_api_key", "flvx_api_url", "cf_token", "domain_name", "check_api_url"}
		config := make(map[string]string)
		for _, key := range keys {
			val, err := s.ConfigRepo.GetConfig(key)
			if err != nil {
				if err == sql.ErrNoRows {
					config[key] = ""
				} else {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
			} else {
				config[key] = val
			}
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(config)

	case http.MethodPost:
		var req map[string]string
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}

		for key, val := range req {
			if err := s.ConfigRepo.SetConfig(key, val); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"status": "ok"})

	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}
