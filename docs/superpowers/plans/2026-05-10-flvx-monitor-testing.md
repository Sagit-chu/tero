# Flvx Monitor Testing Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Expand unit testing coverage for existing components and establish an End-to-End (E2E) testing framework.

**Architecture:** Unit tests in Go (`testing` package) and React (`vitest`, `@testing-library/react`), plus integration level E2E tests using Playwright.

**Tech Stack:** Go 1.22+, Vitest, React Testing Library, Playwright.

---

### Task 1: Backend Unit Tests - Config Repository CRUD

**Files:**
- Create: `backend/repository/config.go`
- Create: `backend/repository/config_test.go`

- [ ] **Step 1: Write the failing test**
```go
// backend/repository/config_test.go
package repository

import (
	"database/sql"
	"testing"
	"github.com/sagit-chu/flvx-monitor/backend/db"
)

func TestConfigRepo(t *testing.T) {
	database, err := db.InitDB(":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer database.Close()
	repo := NewConfigRepository(database)

	err = repo.SetConfig("interval", "5m")
	if err != nil {
		t.Fatal(err)
	}

	val, err := repo.GetConfig("interval")
	if err != nil {
		t.Fatal(err)
	}
	if val != "5m" {
		t.Fatalf("Expected 5m, got %s", val)
	}

	_, err = repo.GetConfig("missing")
	if err != sql.ErrNoRows {
		t.Fatalf("Expected ErrNoRows, got %v", err)
	}
}
```

- [ ] **Step 2: Run test to verify it fails**
Run: `go test ./backend/repository/...`
Expected: FAIL (undefined: NewConfigRepository)

- [ ] **Step 3: Write minimal implementation**
```go
// backend/repository/config.go
package repository

import "database/sql"

type ConfigRepository struct {
	db *sql.DB
}

func NewConfigRepository(db *sql.DB) *ConfigRepository {
	return &ConfigRepository{db: db}
}

func (r *ConfigRepository) SetConfig(key, value string) error {
	_, err := r.db.Exec("INSERT INTO config (key, value) VALUES (?, ?) ON CONFLICT(key) DO UPDATE SET value=excluded.value", key, value)
	return err
}

func (r *ConfigRepository) GetConfig(key string) (string, error) {
	row := r.db.QueryRow("SELECT value FROM config WHERE key = ?", key)
	var val string
	err := row.Scan(&val)
	return val, err
}
```

- [ ] **Step 4: Run test to verify it passes**
Run: `go test ./backend/repository/...`
Expected: PASS

- [ ] **Step 5: Commit**
```bash
git add backend/repository/config.go backend/repository/config_test.go
git commit -m "test: add config repository crud unit tests"
```

### Task 2: Backend Unit Tests - API Server Edge Cases

**Files:**
- Modify: `backend/api/server_test.go`
- Modify: `backend/api/server.go`

- [ ] **Step 1: Write the failing test**
```go
// backend/api/server_test.go
package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestStatusEndpoint(t *testing.T) {
	tests := []struct {
		name           string
		method         string
		expectedStatus int
		expectedBody   string
	}{
		{"Valid GET request", "GET", http.StatusOK, `{"status":"ok"}` + "\n"},
		{"Invalid POST request", "POST", http.StatusMethodNotAllowed, ""},
	}

	srv := NewServer()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest(tt.method, "/api/status", nil)
			rr := httptest.NewRecorder()
			srv.Router.ServeHTTP(rr, req)

			if rr.Code != tt.expectedStatus {
				t.Errorf("expected status %v, got %v", tt.expectedStatus, rr.Code)
			}
			if tt.expectedBody != "" && rr.Body.String() != tt.expectedBody {
				t.Errorf("expected body %v, got %v", tt.expectedBody, rr.Body.String())
			}
		})
	}
}
```

- [ ] **Step 2: Run test to verify it fails**
Run: `go test ./backend/api/...`
Expected: FAIL (Invalid POST request returns 200 instead of 405)

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
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
```

- [ ] **Step 4: Run test to verify it passes**
Run: `go test ./backend/api/...`
Expected: PASS

- [ ] **Step 5: Commit**
```bash
git add backend/api/server.go backend/api/server_test.go
git commit -m "test: expand api server unit tests with edge cases"
```

### Task 3: Frontend Unit Tests - React App Component States

**Files:**
- Modify: `frontend/src/App.test.tsx`

- [ ] **Step 1: Write the failing tests**
```tsx
// frontend/src/App.test.tsx
import { render, screen, waitFor } from '@testing-library/react';
import { vi } from 'vitest';
import App from './App';
import '@testing-library/jest-dom';

describe('App Component', () => {
  beforeEach(() => {
    global.fetch = vi.fn();
  });

  afterEach(() => {
    vi.resetAllMocks();
  });

  test('renders loading state initially', () => {
    // Mock fetch to return a promise that never resolves
    (global.fetch as any).mockReturnValue(new Promise(() => {}));
    render(<App />);
    expect(screen.getByText('Loading...')).toBeInTheDocument();
  });

  test('renders success state', async () => {
    (global.fetch as any).mockResolvedValue({
      json: () => Promise.resolve({ status: 'ok' }),
    });
    render(<App />);
    
    await waitFor(() => {
      expect(screen.getByText('ok')).toBeInTheDocument();
    });
    expect(screen.queryByText('Loading...')).not.toBeInTheDocument();
  });

  test('renders error state on fetch failure', async () => {
    (global.fetch as any).mockRejectedValue(new Error('Network Error'));
    render(<App />);
    
    await waitFor(() => {
      expect(screen.getByText('Error')).toBeInTheDocument();
    });
  });
});
```

- [ ] **Step 2: Run test to verify it passes**
Run: `cd frontend && npm test`
Note: No new implementation needed in `App.tsx` as it already handles these states.
Expected: PASS (if setup correctly, it should just pass)

- [ ] **Step 3: Commit**
```bash
git add frontend/src/App.test.tsx
git commit -m "test: add comprehensive frontend component state tests"
```

### Task 4: E2E Playwright Setup & Core Tests

**Files:**
- Create: `e2e/package.json`
- Create: `e2e/playwright.config.ts`
- Create: `e2e/tests/dashboard.spec.ts`

- [ ] **Step 1: Setup Playwright directory**
```bash
mkdir -p e2e/tests
cd e2e
npm init -y
npm install -D @playwright/test @types/node
npx playwright install --with-deps chromium
```

- [ ] **Step 2: Write Playwright configuration**
```typescript
// e2e/playwright.config.ts
import { defineConfig, devices } from '@playwright/test';

export default defineConfig({
  testDir: './tests',
  fullyParallel: true,
  retries: process.env.CI ? 2 : 0,
  workers: process.env.CI ? 1 : undefined,
  reporter: 'html',
  use: {
    baseURL: 'http://localhost:5173',
    trace: 'on-first-retry',
  },
  projects: [
    {
      name: 'chromium',
      use: { ...devices['Desktop Chrome'] },
    },
  ],
  webServer: [
    {
      command: 'cd ../ && go run backend/api/cmd/main.go', // We will need an entry point
      port: 8080,
      reuseExistingServer: !process.env.CI,
    },
    {
      command: 'cd ../frontend && npm run dev',
      port: 5173,
      reuseExistingServer: !process.env.CI,
    },
  ],
});
```

- [ ] **Step 3: Create Go Backend Entrypoint**
```go
// backend/api/cmd/main.go
package main

import (
	"log"
	"net/http"
	"github.com/sagit-chu/flvx-monitor/backend/api"
)

func main() {
	srv := api.NewServer()
	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", srv.Router); err != nil {
		log.Fatal(err)
	}
}
```

- [ ] **Step 4: Write failing E2E test**
```typescript
// e2e/tests/dashboard.spec.ts
import { test, expect } from '@playwright/test';

test('has title and displays status ok', async ({ page }) => {
  await page.goto('/');

  // Expect a title "to contain" a substring.
  await expect(page).toHaveTitle(/Vite \+ React/); // Default vite title, we need to update it
  
  // Wait for the loading state to disappear and status ok to appear
  await expect(page.locator('text=System Status:')).toBeVisible();
  await expect(page.locator('text=ok')).toBeVisible();
});
```

- [ ] **Step 5: Run test to verify it passes (with minimal fix)**
Run: `sed -i.bak 's/Vite + React/Flvx Monitor Dashboard/' ../frontend/index.html` (Fixing the title to make the test pass)
Run: `cd e2e && npx playwright test`
Expected: PASS

- [ ] **Step 6: Commit**
```bash
git add e2e/ backend/api/cmd/ frontend/index.html
git commit -m "test: setup playwright and core e2e dashboard tests"
```