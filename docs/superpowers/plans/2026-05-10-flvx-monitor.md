# Flvx Monitor Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Build a Go-based monitoring daemon with an embedded React frontend that checks node connectivity and auto-replaces failing nodes via the flvx API and updates Cloudflare DNS.

**Architecture:** A Go backend serving a REST API and static files, backed by SQLite. The backend runs a periodic ticker that pings the active node, calls a GFW check API, and triggers replacement logic if blocked.

**Tech Stack:** Go 1.22+, SQLite (mattn/go-sqlite3), React 18, Vite, Shadcn UI, TailwindCSS.

---

### Task 1: Go Project & Database Setup

**Files:**
- Create: `go.mod`
- Create: `backend/db/db.go`
- Test: `backend/db/db_test.go`

- [ ] **Step 1: Initialize Go module**
```bash
go mod init github.com/sagit-chu/flvx-monitor
go get github.com/mattn/go-sqlite3
```

- [ ] **Step 2: Write failing test for DB init**
```go
// backend/db/db_test.go
package db

import (
	"os"
	"testing"
)

func TestInitDB(t *testing.T) {
	os.Remove("test.db")
	db, err := InitDB("test.db")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	defer db.Close()

	var name string
	err = db.QueryRow("SELECT name FROM sqlite_master WHERE type='table' AND name='nodes'").Scan(&name)
	if err != nil {
		t.Fatalf("Table nodes not created")
	}
}
```

- [ ] **Step 3: Run test (Verify it fails)**
Run: `go test ./backend/db/...`
Expected: FAIL (InitDB not found)

- [ ] **Step 4: Write minimal implementation**
```go
// backend/db/db.go
package db

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

func InitDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dsn)
	if err != nil {
		return nil, err
	}
	schema := `
	CREATE TABLE IF NOT EXISTS config (key TEXT PRIMARY KEY, value TEXT);
	CREATE TABLE IF NOT EXISTS nodes (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		ip TEXT NOT NULL,
		ssh_port TEXT NOT NULL,
		ssh_password TEXT NOT NULL,
		status TEXT DEFAULT 'standby'
	);
	`
	_, err = db.Exec(schema)
	return db, err
}
```

- [ ] **Step 5: Run test (Verify it passes)**
Run: `go test ./backend/db/...`
Expected: PASS

- [ ] **Step 6: Commit**
```bash
git add go.mod go.sum backend/db/
git commit -m "feat: init go module and sqlite db schema"
```

### Task 2: Node Repository CRUD

**Files:**
- Create: `backend/repository/node.go`
- Test: `backend/repository/node_test.go`

- [ ] **Step 1: Write failing test for Add and Get standby node**
```go
// backend/repository/node_test.go
package repository

import (
	"testing"
	"github.com/sagit-chu/flvx-monitor/backend/db"
)

func TestNodeRepo(t *testing.T) {
	database, _ := db.InitDB(":memory:")
	defer database.Close()
	repo := NewNodeRepository(database)

	err := repo.AddNode("1.1.1.1", "22", "pass")
	if err != nil {
		t.Fatal(err)
	}

	node, err := repo.GetNextStandby()
	if err != nil {
		t.Fatal(err)
	}
	if node.IP != "1.1.1.1" {
		t.Fatalf("Expected IP 1.1.1.1, got %s", node.IP)
	}
}
```

- [ ] **Step 2: Run test (Verify it fails)**
Run: `go test ./backend/repository/...`
Expected: FAIL (undefined: NewNodeRepository)

- [ ] **Step 3: Write minimal implementation**
```go
// backend/repository/node.go
package repository

import "database/sql"

type Node struct {
	ID          int
	IP          string
	SSHPort     string
	SSHPassword string
	Status      string
}

type NodeRepository struct {
	db *sql.DB
}

func NewNodeRepository(db *sql.DB) *NodeRepository {
	return &NodeRepository{db: db}
}

func (r *NodeRepository) AddNode(ip, port, password string) error {
	_, err := r.db.Exec("INSERT INTO nodes (ip, ssh_port, ssh_password) VALUES (?, ?, ?)", ip, port, password)
	return err
}

func (r *NodeRepository) GetNextStandby() (*Node, error) {
	row := r.db.QueryRow("SELECT id, ip, ssh_port, ssh_password, status FROM nodes WHERE status = 'standby' LIMIT 1")
	var n Node
	err := row.Scan(&n.ID, &n.IP, &n.SSHPort, &n.SSHPassword, &n.Status)
	if err != nil {
		return nil, err
	}
	return &n, nil
}
```

- [ ] **Step 4: Run test (Verify it passes)**
Run: `go test ./backend/repository/...`
Expected: PASS

- [ ] **Step 5: Commit**
```bash
git add backend/repository/
git commit -m "feat: add node repository for sqlite"
```

### Task 3: API Server Setup

**Files:**
- Create: `backend/api/server.go`
- Test: `backend/api/server_test.go`

- [ ] **Step 1: Write failing test for /api/status endpoint**
```go
// backend/api/server_test.go
package api

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestStatusEndpoint(t *testing.T) {
	srv := NewServer()
	req, _ := http.NewRequest("GET", "/api/status", nil)
	rr := httptest.NewRecorder()
	srv.Router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
}
```

- [ ] **Step 2: Run test (Verify it fails)**
Run: `go test ./backend/api/...`
Expected: FAIL (undefined: NewServer)

- [ ] **Step 3: Write minimal implementation**
```go
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
```

- [ ] **Step 4: Run test (Verify it passes)**
Run: `go test ./backend/api/...`
Expected: PASS

- [ ] **Step 5: Commit**
```bash
git add backend/api/
git commit -m "feat: basic api server"
```

### Task 4: Frontend Scaffolding

- [ ] **Step 1: Scaffold Vite React app**
```bash
npm create vite@latest frontend -- --template react-ts
cd frontend
npm install
npm install -D tailwindcss postcss autoprefixer
npx tailwindcss init -p
```

- [ ] **Step 2: Commit**
```bash
git add frontend/
git commit -m "feat: scaffold react frontend"
```

### Task 5: Frontend Dashboard Component

**Files:**
- Create: `frontend/src/App.tsx`
- Test: `frontend/src/App.test.tsx`

- [ ] **Step 1: Write failing test**
```tsx
// frontend/src/App.test.tsx
import { render, screen } from '@testing-library/react';
import App from './App';
import '@testing-library/jest-dom';

test('renders dashboard header', () => {
  render(<App />);
  const linkElement = screen.getByText(/Flvx Monitor Dashboard/i);
  expect(linkElement).toBeInTheDocument();
});
```

- [ ] **Step 2: Implement Dashboard UI**
```tsx
// frontend/src/App.tsx
import React, { useEffect, useState } from 'react';

function App() {
  const [status, setStatus] = useState('Loading...');

  useEffect(() => {
    fetch('/api/status')
      .then(r => r.json())
      .then(d => setStatus(d.status))
      .catch(() => setStatus('Error'));
  }, []);

  return (
    <div className="p-8">
      <h1 className="text-2xl font-bold mb-4">Flvx Monitor Dashboard</h1>
      <div className="bg-gray-100 p-4 rounded shadow">
        System Status: <span className="font-semibold">{status}</span>
      </div>
    </div>
  );
}

export default App;
```

- [ ] **Step 3: Commit**
```bash
git add frontend/src/App.tsx frontend/src/App.test.tsx
git commit -m "feat: dashboard basic ui"
```
