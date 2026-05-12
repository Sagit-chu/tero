package monitor

import (
	"testing"
	"github.com/sagit-chu/flvx-monitor/backend/repository"
	"github.com/sagit-chu/flvx-monitor/backend/flvx"
)

type MockNodeRepo struct {
}

func (m *MockNodeRepo) MarkNodeFailed(ip string) error { return nil }
func (m *MockNodeRepo) GetNextStandby() (*repository.Node, error) { return nil, nil }
func (m *MockNodeRepo) MarkNodeActive(id int) error { return nil }

type MockConfigRepo struct {
	TunnelID string
}
func (m *MockConfigRepo) GetConfig(key string) (string, error) {
	if key == "flvx_tunnel_id" {
		return m.TunnelID, nil
	}
	return "", nil
}

type MockPinger struct {
	Alive bool
}
func (m *MockPinger) Ping(addr string) bool { return m.Alive }

type MockGFW struct {
	Blocked bool
}
func (m *MockGFW) IsBlocked(ip string) (bool, error) { return m.Blocked, nil }

type MockFlvx struct{
	Tunnels []map[string]interface{}
	Nodes []flvx.NodeItem
}
func (m *MockFlvx) GetTunnelsRaw() ([]map[string]interface{}, error) { return m.Tunnels, nil }
func (m *MockFlvx) UpdateTunnel(tunnelPayload map[string]interface{}) error { return nil }
func (m *MockFlvx) GetNodes() ([]flvx.NodeItem, error) { return m.Nodes, nil }


type MockCF struct{}
func (m *MockCF) UpdateDNSRecord(domain, newIP string) error { return nil }

func TestMonitorService_RunCycle_DeadNode(t *testing.T) {
	repo := &MockNodeRepo{}
	configRepo := &MockConfigRepo{TunnelID: "1"}
	pinger := &MockPinger{Alive: false}
	gfw := &MockGFW{Blocked: false}
	flvxClient := &MockFlvx{
		Tunnels: []map[string]interface{}{
			{
				"id": float64(1),
				"inNodeId": []interface{}{
					map[string]interface{}{"nodeId": float64(10)},
				},
			},
		},
		Nodes: []flvx.NodeItem{
			{ID: 10, ServerIP: "1.1.1.1"},
		},
	}
	cf := &MockCF{}

	svc := NewMonitorService(repo, configRepo, pinger, gfw, flvxClient, cf, "test.com")
	status := svc.RunCycle()
	
	if status != "dead" && status != "no_active_node" {
		t.Errorf("Expected dead or similar, got %s", status)
	}
}

func TestMonitorService_RunCycle_BlockedNode(t *testing.T) {
	repo := &MockNodeRepo{}
	configRepo := &MockConfigRepo{TunnelID: "2"}
	pinger := &MockPinger{Alive: true}
	gfw := &MockGFW{Blocked: true}
	flvxClient := &MockFlvx{
		Tunnels: []map[string]interface{}{
			{
				"id": float64(2),
				"inNodeId": []interface{}{
					map[string]interface{}{"nodeId": float64(20)},
				},
			},
		},
		Nodes: []flvx.NodeItem{
			{ID: 20, ServerIP: "2.2.2.2"},
		},
	}
	cf := &MockCF{}

	svc := NewMonitorService(repo, configRepo, pinger, gfw, flvxClient, cf, "test.com")
	status := svc.RunCycle()
	
	if status != "blocked" && status != "no_active_node" {
		t.Errorf("Expected blocked or similar, got %s", status)
	}
}