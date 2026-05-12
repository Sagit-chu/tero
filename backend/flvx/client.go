package flvx

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Client struct {
	BaseURL  string
	Account  string
	Password string
	HTTP     *http.Client
	token    string
}

func NewClient(baseURL, account, password string) *Client {
	return &Client{
		BaseURL:  baseURL,
		Account:  account,
		Password: password,
		HTTP:     &http.Client{Timeout: 15 * time.Second},
	}
}

type FlvxResponse struct {
	Code int             `json:"code"`
	Msg  string          `json:"msg"`
	Data json.RawMessage `json:"data"`
}

func (c *Client) Login() error {
	payload := map[string]string{
		"user":     c.Account,
		"username": c.Account,
		"password": c.Password,
	}
	data, _ := json.Marshal(payload)

	req, err := http.NewRequest(http.MethodPost, c.BaseURL+"/api/v1/user/login", bytes.NewBuffer(data))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.HTTP.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var res FlvxResponse
	if err := json.Unmarshal(body, &res); err != nil {
		return fmt.Errorf("failed to decode login response: %v", err)
	}

	if res.Code != 0 {
		return fmt.Errorf("flvx login failed: %s", res.Msg)
	}

	var loginData struct {
		Token string `json:"token"`
	}
	if err := json.Unmarshal(res.Data, &loginData); err != nil {
		return fmt.Errorf("failed to decode login token: %v", err)
	}

	c.token = loginData.Token
	return nil
}

func (c *Client) request(method, path string, payload interface{}) ([]byte, error) {
	if c.token == "" {
		if err := c.Login(); err != nil {
			return nil, err
		}
	}

	var buf io.Reader
	if payload != nil {
		data, _ := json.Marshal(payload)
		buf = bytes.NewBuffer(data)
	}

	req, err := http.NewRequest(method, c.BaseURL+path, buf)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", c.token)

	resp, err := c.HTTP.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var res FlvxResponse
	if err := json.Unmarshal(body, &res); err != nil {
		return nil, fmt.Errorf("failed to decode response: %v", err)
	}

	if res.Code == 401 { // Token expired
		c.token = ""
		return c.request(method, path, payload)
	}

	if res.Code != 0 {
		return nil, fmt.Errorf("flvx api error: %s", res.Msg)
	}

	return res.Data, nil
}

type NodeItem struct {
	ID       int64  `json:"id"`
	Name     string `json:"name"`
	ServerIP string `json:"serverIp"`
	Status   int    `json:"status"`
}

type TunnelItem struct {
	ID       int64                    `json:"id"`
	Name     string                   `json:"name"`
	InIP     string                   `json:"inIp"`
	InNodeId []map[string]interface{} `json:"inNodeId"`
	Type     int                      `json:"type"`
	Status   int                      `json:"status"`
	// Keep other fields raw if we need to echo them
}

func (c *Client) GetNodes() ([]NodeItem, error) {
	data, err := c.request(http.MethodPost, "/api/v1/node/list", map[string]string{})
	if err != nil {
		return nil, err
	}
	var nodes []NodeItem
	if err := json.Unmarshal(data, &nodes); err != nil {
		return nil, err
	}
	return nodes, nil
}

func (c *Client) GetTunnelsRaw() ([]map[string]interface{}, error) {
	data, err := c.request(http.MethodPost, "/api/v1/tunnel/list", map[string]string{})
	if err != nil {
		return nil, err
	}
	var tunnels []map[string]interface{}
	if err := json.Unmarshal(data, &tunnels); err != nil {
		return nil, err
	}
	return tunnels, nil
}

func (c *Client) GetTunnels() ([]TunnelItem, error) {
	data, err := c.request(http.MethodPost, "/api/v1/tunnel/list", map[string]string{})
	if err != nil {
		return nil, err
	}
	var tunnels []TunnelItem
	if err := json.Unmarshal(data, &tunnels); err != nil {
		return nil, err
	}
	return tunnels, nil
}

func (c *Client) UpdateTunnel(tunnelPayload map[string]interface{}) error {
	_, err := c.request(http.MethodPost, "/api/v1/tunnel/update", tunnelPayload)
	return err
}

func (c *Client) ReplaceNode(ip, port, password string) error {
	// Deprecated: kept to satisfy existing interface until refactored completely.
	return fmt.Errorf("not implemented, use UpdateTunnel instead")
}
