package monitor

import (
	"testing"
)

type MockNodeRepo struct {
	ActiveIP string
}

func (m *MockNodeRepo) GetActiveNode() (string, string, error) {
	return m.ActiveIP, "22", nil
}
func (m *MockNodeRepo) MarkNodeFailed(ip string) error { return nil }

type MockPinger struct {
	Alive bool
}
func (m *MockPinger) Ping(addr string) bool { return m.Alive }

type MockGFW struct {
	Blocked bool
}
func (m *MockGFW) IsBlocked(ip string) (bool, error) { return m.Blocked, nil }

func TestMonitorService_RunCycle_DeadNode(t *testing.T) {
	repo := &MockNodeRepo{ActiveIP: "1.1.1.1"}
	pinger := &MockPinger{Alive: false}
	gfw := &MockGFW{Blocked: false}

	svc := NewMonitorService(repo, pinger, gfw)
	status := svc.RunCycle()
	
	if status != "dead" {
		t.Errorf("Expected dead, got %s", status)
	}
}

func TestMonitorService_RunCycle_BlockedNode(t *testing.T) {
	repo := &MockNodeRepo{ActiveIP: "2.2.2.2"}
	pinger := &MockPinger{Alive: true}
	gfw := &MockGFW{Blocked: true}

	svc := NewMonitorService(repo, pinger, gfw)
	status := svc.RunCycle()
	
	if status != "blocked" {
		t.Errorf("Expected blocked, got %s", status)
	}
}
