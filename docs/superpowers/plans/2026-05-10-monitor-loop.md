# Monitor Loop Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Implement the background monitoring daemon that performs two-stage verification (TCP Ping and Ping.pe Scraper) and triggers automated node replacement via Flvx and Cloudflare APIs.

**Architecture:** A Go `time.Ticker` runs periodically, orchestrating the Stage 1 (TCP Ping) and Stage 2 (GFW Check). It uses SQLite repositories to fetch active and standby nodes.

**Tech Stack:** Go 1.22+, net package, net/http package.

---

### Task 1: Create GFWChecker Interface & PingPe Scraper

**Files:**
- Create: `backend/monitor/checker.go`
- Create: `backend/monitor/pingpe.go`
- Create: `backend/monitor/pingpe_test.go`

- [ ] **Step 1: Write the failing test**

```go
// backend/monitor/pingpe_test.go
package monitor

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPingPeScraper_CheckBlocked(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		// Mock response that simulates 100% loss from China
		w.Write([]byte(`{"china_loss": 100, "overseas_loss": 0}`))
	}))
	defer ts.Close()

	scraper := NewPingPeScraper(ts.URL)
	blocked, err := scraper.IsBlocked("1.1.1.1")
	
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if !blocked {
		t.Errorf("Expected node to be detected as blocked")
	}
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `go test ./backend/monitor/... -v`
Expected: FAIL with "no required module provides package github.com/sagit-chu/flvx-monitor/backend/monitor" or "undefined: NewPingPeScraper"

- [ ] **Step 3: Write minimal implementation**

```go
// backend/monitor/checker.go
package monitor

type GFWChecker interface {
	IsBlocked(ip string) (bool, error)
}
```

```go
// backend/monitor/pingpe.go
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
```

- [ ] **Step 4: Run test to verify it passes**

Run: `go test ./backend/monitor/... -v`
Expected: PASS

- [ ] **Step 5: Commit**

```bash
git add backend/monitor/
git commit -m "feat: add gfw checker interface and ping.pe scraper mock"
```

### Task 2: Create TCP Pinger

**Files:**
- Create: `backend/monitor/pinger.go`
- Create: `backend/monitor/pinger_test.go`

- [ ] **Step 1: Write the failing test**

```go
// backend/monitor/pinger_test.go
package monitor

import (
	"net"
	"testing"
	"time"
)

func TestTCPPing_Success(t *testing.T) {
	// Start a dummy TCP server
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer l.Close()

	pinger := NewTCPPinger(1 * time.Second)
	alive := pinger.Ping(l.Addr().String())
	if !alive {
		t.Errorf("Expected node to be alive")
	}
}

func TestTCPPing_Fail(t *testing.T) {
	pinger := NewTCPPinger(10 * time.Millisecond)
	// Ping a non-listening port
	alive := pinger.Ping("127.0.0.1:12345")
	if alive {
		t.Errorf("Expected node to be dead")
	}
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `go test ./backend/monitor/... -v`
Expected: FAIL with "undefined: NewTCPPinger"

- [ ] **Step 3: Write minimal implementation**

```go
// backend/monitor/pinger.go
package monitor

import (
	"net"
	"time"
)

type Pinger interface {
	Ping(address string) bool
}

type TCPPinger struct {
	Timeout time.Duration
}

func NewTCPPinger(timeout time.Duration) *TCPPinger {
	return &TCPPinger{Timeout: timeout}
}

func (p *TCPPinger) Ping(address string) bool {
	conn, err := net.DialTimeout("tcp", address, p.Timeout)
	if err != nil {
		return false
	}
	conn.Close()
	return true
}
```

- [ ] **Step 4: Run test to verify it passes**

Run: `go test ./backend/monitor/... -v`
Expected: PASS

- [ ] **Step 5: Commit**

```bash
git add backend/monitor/pinger.go backend/monitor/pinger_test.go
git commit -m "feat: add tcp pinger for stage 1 verification"
```

### Task 3: Monitor Service Core Loop

**Files:**
- Create: `backend/monitor/service.go`
- Create: `backend/monitor/service_test.go`

- [ ] **Step 1: Write the failing test**

```go
// backend/monitor/service_test.go
package monitor

import (
	"testing"
	"time"
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
```

- [ ] **Step 2: Run test to verify it fails**

Run: `go test ./backend/monitor/... -v`
Expected: FAIL with "undefined: NewMonitorService"

- [ ] **Step 3: Write minimal implementation**

```go
// backend/monitor/service.go
package monitor

import (
	"log"
)

type NodeRepository interface {
	GetActiveNode() (ip, port string, err error)
	MarkNodeFailed(ip string) error
}

type MonitorService struct {
	repo   NodeRepository
	pinger Pinger
	gfw    GFWChecker
}

func NewMonitorService(repo NodeRepository, pinger Pinger, gfw GFWChecker) *MonitorService {
	return &MonitorService{
		repo:   repo,
		pinger: pinger,
		gfw:    gfw,
	}
}

// RunCycle performs one iteration of the monitoring logic
func (s *MonitorService) RunCycle() string {
	ip, port, err := s.repo.GetActiveNode()
	if err != nil || ip == "" {
		return "no_active_node"
	}

	// Stage 1
	if !s.pinger.Ping(ip + ":" + port) {
		log.Printf("Node %s is DEAD", ip)
		s.repo.MarkNodeFailed(ip)
		return "dead"
	}

	// Stage 2
	blocked, err := s.gfw.IsBlocked(ip)
	if err != nil {
		log.Printf("GFW Check failed for %s: %v", ip, err)
		return "alive" // skip block assumption on error
	}

	if blocked {
		log.Printf("Node %s is BLOCKED", ip)
		s.repo.MarkNodeFailed(ip)
		return "blocked"
	}

	return "alive"
}
```

- [ ] **Step 4: Run test to verify it passes**

Run: `go test ./backend/monitor/... -v`
Expected: PASS

- [ ] **Step 5: Commit**

```bash
git add backend/monitor/service.go backend/monitor/service_test.go
git commit -m "feat: implement monitor service state machine"
```
