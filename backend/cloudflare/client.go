package cloudflare

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Client struct {
	BaseURL string
	Token   string
	HTTP    *http.Client
}

func NewClient(baseURL, token string) *Client {
	if baseURL == "" {
		baseURL = "https://api.cloudflare.com"
	}
	return &Client{
		BaseURL: baseURL,
		Token:   token,
		HTTP:    &http.Client{Timeout: 15 * time.Second},
	}
}

func (c *Client) doReq(method, path string, body []byte) (map[string]interface{}, error) {
	req, err := http.NewRequest(method, c.BaseURL+path, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+c.Token)
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.HTTP.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var res map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&res)
	if success, ok := res["success"].(bool); !ok || !success {
		return nil, fmt.Errorf("cloudflare API error: %v", res)
	}
	return res, nil
}

func (c *Client) UpdateDNSRecord(domain, newIP string) error {
	// 1. Get Zone
	zones, err := c.doReq(http.MethodGet, "/client/v4/zones?name="+domain, nil)
	if err != nil { return err }
	zoneResults := zones["result"].([]interface{})
	if len(zoneResults) == 0 { return fmt.Errorf("zone not found") }
	zoneID := zoneResults[0].(map[string]interface{})["id"].(string)

	// 2. Get Record
	records, err := c.doReq(http.MethodGet, "/client/v4/zones/"+zoneID+"/dns_records?type=A&name="+domain, nil)
	if err != nil { return err }
	recordResults := records["result"].([]interface{})
	if len(recordResults) == 0 { return fmt.Errorf("record not found") }
	recordID := recordResults[0].(map[string]interface{})["id"].(string)

	// 3. Update Record
	payload, _ := json.Marshal(map[string]string{"type": "A", "name": domain, "content": newIP, "proxied": "false"})
	_, err = c.doReq(http.MethodPut, "/client/v4/zones/"+zoneID+"/dns_records/"+recordID, payload)
	return err
}
