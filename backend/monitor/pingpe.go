package monitor

import (
	"encoding/json"
	"net/http"
	"time"
)

type PingPeScraper struct {
	BaseURL    string
	HTTPClient *http.Client
}

func NewPingPeScraper(baseURL string) *PingPeScraper {
	if baseURL == "" {
		baseURL = "https://ping.pe"
	}
	return &PingPeScraper{
		BaseURL: baseURL,
		HTTPClient: &http.Client{Timeout: 10 * time.Second},
	}
}

// IsBlocked simulates scraping ping.pe. In a real scenario, this parses HTML or WebSocket.
func (s *PingPeScraper) IsBlocked(ip string) (bool, error) {
	req, err := http.NewRequest("GET", s.BaseURL+"/?ip="+ip, nil)
	if err != nil {
		return false, err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0")
	
	resp, err := s.HTTPClient.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	// For testing purposes, we parse a mock JSON. 
	// In reality, if it's HTML, you might use strings.Contains.
	var result struct {
		ChinaLoss    int `json:"china_loss"`
		OverseasLoss int `json:"overseas_loss"`
	}
	
	if err := json.NewDecoder(resp.Body).Decode(&result); err == nil {
		if result.ChinaLoss == 100 && result.OverseasLoss == 0 {
			return true, nil
		}
	}

	return false, nil
}
