package flvx

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Client struct {
	BaseURL  string
	Account  string
	Password string
	HTTP     *http.Client
}

func NewClient(baseURL, account, password string) *Client {
	return &Client{
		BaseURL:  baseURL,
		Account:  account,
		Password: password,
		HTTP:     &http.Client{Timeout: 15 * time.Second},
	}
}

func (c *Client) ReplaceNode(ip, port, password string) error {
	payload := map[string]string{
		"ip":           ip,
		"ssh_port":     port,
		"ssh_password": password,
	}
	data, _ := json.Marshal(payload)

	req, err := http.NewRequest(http.MethodPost, c.BaseURL+"/nodes/replace", bytes.NewBuffer(data))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(c.Account, c.Password)

	resp, err := c.HTTP.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("flvx API returned status: %d", resp.StatusCode)
	}
	return nil
}
