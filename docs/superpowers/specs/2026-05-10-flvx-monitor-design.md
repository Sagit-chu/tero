# Flvx Node Monitor & Auto-Replacement System - Design Specification

## 1. Overview
A system designed to monitor a Linux machine's connectivity (checking if it is dead or blocked by the GFW). Upon detecting a failure, it automatically replaces the node using the flvx panel API and updates Cloudflare DNS. It features a React-based frontend to monitor status and manage a pool of standby nodes.

## 2. Architecture
- **Backend:** Go (Golang)
  - Runs a background ticker for monitoring.
  - Exposes a RESTful API for the frontend.
  - Serves the compiled frontend static files (Single-binary deployment).
- **Frontend:** React + Vite + Shadcn UI
  - A modern, responsive Single Page Application (SPA).
- **Database:** SQLite
  - Stores configuration, status logs, and the pool of standby nodes.

## 3. Components

### 3.1 Go Backend (Daemon & API)
- **Monitor Service:** Executes the two-stage verification on a scheduled interval (e.g., every 5 minutes).
- **Flvx Client:** Interfaces with the flvx panel API (`https://github.com/Sagit-chu/flvx`) to submit replacement requests using a standby node.
- **Cloudflare Client:** Updates the DNS A record via Cloudflare API.
- **API Server:** Provides endpoints (`/api/status`, `/api/nodes`, `/api/config`) for the frontend.

### 3.2 React Frontend
- **Dashboard:** Real-time visualization of current node status (Alive, Dead, Blocked, Replacing).
- **Standby Node Manager:** A table view to add, edit, or remove standby nodes. Each entry contains:
  - IP Address
  - SSH Port
  - SSH Password
- **Settings:** Form to configure flvx API details, Cloudflare tokens, Domain Name, and the 3rd-party GFW check API endpoint.

### 3.3 Database Schema (SQLite)
- **`config` table:** Key-value pairs for settings.
- **`nodes` table:** `id`, `ip`, `ssh_port`, `ssh_password`, `status` (standby, active, used, failed), `created_at`.
- **`logs` table:** Monitoring events and replacement history.

## 4. Data Flow

1. **Monitoring Loop:**
   - **Stage 1 (Global Check):** Go backend pings the active node. If it fails consecutively, mark as "Dead".
   - **Stage 2 (GFW Check):** If Stage 1 succeeds, Go backend calls the configured 3rd-party API to check connectivity from inside China. If it fails, mark as "Blocked".
2. **Replacement Process:**
   - If "Dead" or "Blocked", the backend queries SQLite for the next available node in the `nodes` table where `status = 'standby'`.
   - The backend calls the flvx API, sending the new node's IP, port, and password.
   - Upon successful replacement in flvx, the backend calls the Cloudflare API to update the domain's A record.
   - The new node is marked as `active`, and the old one is marked as `used`/`failed`.
3. **Frontend Interaction:**
   - The React app fetches data from the Go API to display the current state and allows the user to replenish the standby node pool.

## 5. Error Handling & Edge Cases
- **Empty Standby Pool:** If no standby nodes are available when a replacement is needed, the system will log a critical error and halt replacement attempts until the pool is replenished (optionally sending a notification).
- **Flvx API Failure:** If the flvx panel fails to replace the node, the backend will retry up to a configured limit before marking the standby node as failed and attempting the next one in the pool.
- **Cloudflare API Failure:** If the DNS update fails, it will retry exponentially, as the node is already replaced.

## 6. Testing Strategy
- **Backend:** Unit tests for the Monitor Service logic (mocking ping and HTTP requests), SQLite interactions, and flvx/Cloudflare API clients.
- **Frontend:** Component tests for the Shadcn UI components and integration tests for API calls.
