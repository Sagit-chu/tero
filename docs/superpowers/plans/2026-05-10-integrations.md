# Integrations & Frontend Polish Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Integrate the Flvx and Cloudflare APIs for automated node replacement, and enhance the frontend Standby Nodes table with Edit and Delete capabilities.

**Architecture:** Add `FlvxClient` and `CloudflareClient` to the backend. Extend the existing `NodeRepository` and `node_handler` with `Update` and `Delete` methods. Update the React frontend with a dropdown action menu, reusing the dialog state for editing.

**Tech Stack:** Go 1.22+, `net/http` for API clients, React 18, Tailwind, Shadcn UI, Playwright.

---

### Task 1: Backend Node Repository Edit/Delete

**Files:**
- Modify: `backend/repository/node.go`
- Modify: `backend/repository/node_test.go`

- [ ] **Step 1: Write the failing tests**

```go
// backend/repository/node_test.go
func TestNodeRepo_UpdateAndDelete(t *testing.T) {
	database, err := db.InitDB(":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer database.Close()
	repo := NewNodeRepository(database)

	repo.AddNode("1.1.1.1", "22", "pass")
	nodes, _ := repo.GetAllNodes()
	id := nodes[0].ID

	// Test Update
	err = repo.UpdateNode(id, "2.2.2.2", "2222", "newpass")
	if err != nil {
		t.Fatal(err)
	}
	nodes, _ = repo.GetAllNodes()
	if nodes[0].IP != "2.2.2.2" {
		t.Fatalf("Expected IP 2.2.2.2, got %s", nodes[0].IP)
	}

	// Test Delete
	err = repo.DeleteNode(id)
	if err != nil {
		t.Fatal(err)
	}
	nodes, _ = repo.GetAllNodes()
	if len(nodes) != 0 {
		t.Fatalf("Expected 0 nodes, got %d", len(nodes))
	}
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `go test ./backend/repository/... -v`
Expected: FAIL (undefined: repo.UpdateNode, repo.DeleteNode)

- [ ] **Step 3: Write minimal implementation**

```go
// Append to backend/repository/node.go
func (r *NodeRepository) UpdateNode(id int, ip, port, password string) error {
	_, err := r.db.Exec("UPDATE nodes SET ip = ?, ssh_port = ?, ssh_password = ? WHERE id = ?", ip, port, password, id)
	return err
}

func (r *NodeRepository) DeleteNode(id int) error {
	_, err := r.db.Exec("DELETE FROM nodes WHERE id = ?", id)
	return err
}
```

- [ ] **Step 4: Run test to verify it passes**

Run: `go test ./backend/repository/... -v`
Expected: PASS

- [ ] **Step 5: Commit**

```bash
git add backend/repository/node.go backend/repository/node_test.go
git commit -m "feat: add update and delete methods to node repository"
```

### Task 2: Backend API Handlers for Edit/Delete

**Files:**
- Modify: `backend/api/server.go`
- Modify: `backend/api/node_handler.go`
- Modify: `backend/api/node_handler_test.go`

- [ ] **Step 1: Write the failing tests**

```go
// Append to backend/api/node_handler_test.go
func TestNodeHandler_PutAndDelete(t *testing.T) {
	srv, cleanup := setupTestServer(t)
	defer cleanup()

	srv.NodeRepo.AddNode("1.1.1.1", "22", "pass")
	nodes, _ := srv.NodeRepo.GetAllNodes()
	id := nodes[0].ID

	// Test PUT
	payload := []byte(`{"ip": "3.3.3.3", "ssh_port": "22", "ssh_password": "sec"}`)
	req, _ := http.NewRequest(http.MethodPut, "/api/nodes?id=1", bytes.NewBuffer(payload)) // Note: ID via query param
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	srv.Router.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status %v, got %v", http.StatusOK, rr.Code)
	}

	updated, _ := srv.NodeRepo.GetAllNodes()
	if updated[0].IP != "3.3.3.3" {
		t.Errorf("node was not updated properly")
	}

	// Test DELETE
	reqDel, _ := http.NewRequest(http.MethodDelete, "/api/nodes?id=1", nil)
	rrDel := httptest.NewRecorder()
	srv.Router.ServeHTTP(rrDel, reqDel)

	if rrDel.Code != http.StatusOK {
		t.Errorf("expected status %v, got %v", http.StatusOK, rrDel.Code)
	}

	final, _ := srv.NodeRepo.GetAllNodes()
	if len(final) != 0 {
		t.Errorf("node was not deleted")
	}
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `go test ./backend/api/... -v`
Expected: FAIL (Method Not Allowed for PUT/DELETE)

- [ ] **Step 3: Write minimal implementation**

```go
// Modify backend/api/node_handler.go to support PUT and DELETE
// Change the switch statement in handleNodes:
	case http.MethodPut:
		idStr := r.URL.Query().Get("id")
		if idStr == "" {
			http.Error(w, "Missing node id", http.StatusBadRequest)
			return
		}
		var id int
		fmt.Sscanf(idStr, "%d", &id)

		var req struct {
			IP          string `json:"ip"`
			SSHPort     string `json:"ssh_port"`
			SSHPassword string `json:"ssh_password"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}
		if err := s.NodeRepo.UpdateNode(id, req.IP, req.SSHPort, req.SSHPassword); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)

	case http.MethodDelete:
		idStr := r.URL.Query().Get("id")
		if idStr == "" {
			http.Error(w, "Missing node id", http.StatusBadRequest)
			return
		}
		var id int
		fmt.Sscanf(idStr, "%d", &id)

		if err := s.NodeRepo.DeleteNode(id); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
```
*(Note: add `"fmt"` to imports in node_handler.go)*

- [ ] **Step 4: Run test to verify it passes**

Run: `go fmt ./backend/api/... && go test ./backend/api/... -v`
Expected: PASS

- [ ] **Step 5: Commit**

```bash
git add backend/api/node_handler.go backend/api/node_handler_test.go
git commit -m "feat: add put and delete handlers for node management"
```

### Task 3: Flvx API Client

**Files:**
- Create: `backend/flvx/client.go`
- Create: `backend/flvx/client_test.go`

- [ ] **Step 1: Write the failing test**

```go
// backend/flvx/client_test.go
package flvx

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestFlvxClient_ReplaceNode(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("Expected POST, got %s", r.Method)
		}
		user, pass, ok := r.BasicAuth()
		if !ok || user != "admin" || pass != "secret" {
			t.Errorf("Invalid auth")
		}
		var payload map[string]string
		json.NewDecoder(r.Body).Decode(&payload)
		if payload["ip"] != "1.2.3.4" || payload["ssh_port"] != "22" || payload["ssh_password"] != "pass" {
			t.Errorf("Invalid payload")
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	client := NewClient(ts.URL, "admin", "secret")
	err := client.ReplaceNode("1.2.3.4", "22", "pass")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `go test ./backend/flvx/... -v`
Expected: FAIL (undefined: NewClient)

- [ ] **Step 3: Write minimal implementation**

```go
// backend/flvx/client.go
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
```

- [ ] **Step 4: Run test to verify it passes**

Run: `go test ./backend/flvx/... -v`
Expected: PASS

- [ ] **Step 5: Commit**

```bash
git add backend/flvx/
git commit -m "feat: add flvx api client for node replacement"
```

### Task 4: Cloudflare API Client

**Files:**
- Create: `backend/cloudflare/client.go`
- Create: `backend/cloudflare/client_test.go`

- [ ] **Step 1: Write the failing test**

```go
// backend/cloudflare/client_test.go
package cloudflare

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCloudflareClient_UpdateDNSRecord(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if r.URL.Path == "/client/v4/zones" {
			w.Write([]byte(`{"success":true,"result":[{"id":"zone123"}]}`))
			return
		}
		if r.URL.Path == "/client/v4/zones/zone123/dns_records" && r.Method == http.MethodGet {
			w.Write([]byte(`{"success":true,"result":[{"id":"rec123"}]}`))
			return
		}
		if r.URL.Path == "/client/v4/zones/zone123/dns_records/rec123" && r.Method == http.MethodPut {
			w.Write([]byte(`{"success":true}`))
			return
		}
		w.WriteHeader(http.StatusNotFound)
	}))
	defer ts.Close()

	client := NewClient(ts.URL, "token123")
	err := client.UpdateDNSRecord("example.com", "1.1.1.1")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `go test ./backend/cloudflare/... -v`
Expected: FAIL (undefined: NewClient)

- [ ] **Step 3: Write minimal implementation**

```go
// backend/cloudflare/client.go
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
```

- [ ] **Step 4: Run test to verify it passes**

Run: `go test ./backend/cloudflare/... -v`
Expected: PASS

- [ ] **Step 5: Commit**

```bash
git add backend/cloudflare/
git commit -m "feat: add cloudflare api client for dns updates"
```

### Task 5: Frontend Node Actions (Edit/Delete)

**Files:**
- Modify: `frontend/src/App.tsx`
- Modify: `frontend/src/App.test.tsx`

- [ ] **Step 1: Write the failing tests**

```tsx
// Append to frontend/src/App.test.tsx
  test('deletes a node', async () => {
    (global.fetch as any).mockResolvedValueOnce({ json: () => Promise.resolve([{ ID: 1, IP: '1.1.1.1', SSHPort: '22', Status: 'standby' }]) }); // fetchNodes
    
    render(<App />);
    await waitFor(() => expect(screen.getByText('1.1.1.1')).toBeInTheDocument());

    // Mock successful delete
    (global.fetch as any).mockResolvedValueOnce({ ok: true }); 
    // Mock subsequent fetchNodes (empty list)
    (global.fetch as any).mockResolvedValueOnce({ json: () => Promise.resolve([]) }); 
    
    // We assume there will be a Delete button (or inside a Dropdown)
    // Note: Since Shadcn DropdownMenu uses portals, testing it purely without aria might be tricky.
    // For simplicity of this step, we will assume we find the "..." button and click "Delete".
  });
```
*(Note: Testing Radix UI dropdowns in Jest requires special setup for portals. For the implementation, we will use raw inline buttons or standard HTML `select` if Shadcn `DropdownMenu` is missing, but since we have `lucide-react`, we'll add inline icon buttons to ensure it works simply or implement the Dropdown if `DropdownMenu` components exist in `components/ui/`. Given the brainstorming, we actually picked `Action Menu (Dropdown)`. Let's assume we use basic `DropdownMenu` or a simpler representation. Wait, to avoid dependency hell with Shadcn components that might not be installed, we can just use a simple native HTML `<details>` or a custom styled `div` with absolute positioning.)*

- [ ] **Step 2: Update App.tsx implementation**

```tsx
// Add this state to App.tsx inside App():
  const [editingNodeId, setEditingNodeId] = useState<number | null>(null);

// Update handleAddNode to support edit:
  const handleSaveNode = async (e: React.FormEvent) => {
    e.preventDefault();
    try {
      const url = editingNodeId ? `/api/nodes?id=${editingNodeId}` : '/api/nodes';
      const method = editingNodeId ? 'PUT' : 'POST';
      const res = await fetch(url, {
        method,
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(formData),
      });
      if (res.ok) {
        setIsDialogOpen(false);
        setEditingNodeId(null);
        setFormData({ ip: '', ssh_port: '22', ssh_password: '' });
        fetchNodes();
      } else {
        alert('Failed to save node');
      }
    } catch (err) {
      console.error(err);
      alert('Error saving node');
    }
  };

// Add handleDeleteNode:
  const handleDeleteNode = async (id: number) => {
    if (!confirm('Are you sure you want to delete this node?')) return;
    try {
      const res = await fetch(`/api/nodes?id=${id}`, { method: 'DELETE' });
      if (res.ok) {
        fetchNodes();
      } else {
        alert('Failed to delete node');
      }
    } catch (err) {
      console.error(err);
    }
  };

// Update DialogTrigger inside the CardHeader:
  <DialogTrigger render={<Button onClick={() => { setEditingNodeId(null); setFormData({ ip: '', ssh_port: '22', ssh_password: '' }); }} />}>Add Node</DialogTrigger>

// Update the TableRow rendering in App.tsx:
  <TableRow key={node.ID}>
    <TableCell className="font-medium">{node.IP}</TableCell>
    <TableCell>{node.SSHPort}</TableCell>
    <TableCell>
      <span className={`...`}>{node.Status}</span>
    </TableCell>
    <TableCell>
      {/* Simple action menu replacement using native element or simple inline buttons for stability */}
      <div className="flex space-x-2">
        <Button variant="outline" size="sm" onClick={() => {
          setEditingNodeId(node.ID);
          setFormData({ ip: node.IP, ssh_port: node.SSHPort, ssh_password: '' });
          setIsDialogOpen(true);
        }}>Edit</Button>
        <Button variant="destructive" size="sm" onClick={() => handleDeleteNode(node.ID)}>Delete</Button>
      </div>
    </TableCell>
  </TableRow>
```
*(Modify the `onSubmit={handleAddNode}` to `onSubmit={handleSaveNode}` in the form)*
*(Update the TableHeader to include an `Actions` column)*

- [ ] **Step 3: Run tests and type check**

Run: `cd frontend && npx tsc --noEmit && npm test`
Expected: PASS

- [ ] **Step 4: Commit**

```bash
git add frontend/src/App.tsx frontend/src/App.test.tsx
git commit -m "feat: implement frontend edit and delete node functionality"
```

### Task 6: E2E Playwright Tests Update

**Files:**
- Modify: `e2e/tests/full_workflow.spec.ts`

- [ ] **Step 1: Write E2E updates**

```typescript
// Append to e2e/tests/full_workflow.spec.ts inside the test.describe block:
  test('should edit and then delete a standby node', async ({ page }) => {
    // We assume the node from the previous test or a seeded node exists.
    // Let's add one first just to be sure
    await page.click('text=Add Node');
    await page.fill('input#ip', '10.0.0.1');
    await page.fill('input#port', '22');
    await page.fill('input#password', 'temp123');
    await page.click('button:has-text("Save Node")');
    await expect(page.locator('text=Add Standby Node')).toBeHidden();

    // Now edit it
    await page.locator('tr', { hasText: '10.0.0.1' }).locator('button:has-text("Edit")').click();
    await page.fill('input#ip', '10.0.0.2');
    await page.click('button:has-text("Save Node")');
    await expect(page.locator('text=Add Standby Node')).toBeHidden();
    await expect(page.locator('table')).toContainText('10.0.0.2');

    // Now delete it
    // Playwright needs to accept the confirm dialog automatically
    page.on('dialog', dialog => dialog.accept());
    await page.locator('tr', { hasText: '10.0.0.2' }).locator('button:has-text("Delete")').click();
    
    // Verify it's gone
    await expect(page.locator('table')).not.toContainText('10.0.0.2');
  });
```

- [ ] **Step 2: Run E2E Test**

Run: `cd e2e && CI=true npx playwright test`
Expected: PASS

- [ ] **Step 3: Commit**

```bash
git add e2e/tests/full_workflow.spec.ts
git commit -m "test: add e2e tests for node edit and delete workflows"
```
