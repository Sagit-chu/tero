# Flvx Monitor - AI Agent Context

## Project Overview
Flvx Monitor is an automated node monitoring and replacement system. It actively monitors Linux server nodes for downtime or GFW blocks using a two-stage verification process. Upon detecting a failure, it uses the flvx panel API to replace the node and updates the Cloudflare DNS record.

## Architecture & Tech Stack
- **Backend:** Go 1.22+, `net/http`, `mattn/go-sqlite3`. Serves REST API and runs the background monitor daemon.
- **Frontend:** React 18, Vite, Tailwind CSS, Shadcn UI.
- **Database:** SQLite (Stores standby nodes pool and system configuration).
- **Testing:** Go `testing` (Backend), Vitest (Frontend), Playwright (E2E).

## Current Project State (as of May 2026)
### ✅ Completed Features
1. **Core Infrastructure:** Go module, SQLite initialization, API router structure.
2. **REST API & Repositories:** Endpoints for `/api/status`, `/api/nodes`, and `/api/config`. Repositories implemented for nodes and configuration.
3. **Frontend Dashboard:** Modern UI built with React. Features include System/Node status display, Settings dialog, and a Standby Nodes management table.
4. **Monitor Loop:** 
   - Uses `time.Ticker`.
   - **Stage 1 (Global Check):** TCP Ping to the SSH port.
   - **Stage 2 (GFW Check):** Unofficial Ping.pe scraper simulating an HTTP request to fetch China vs. Overseas packet loss.
5. **Testing Frameworks:** 
   - Backend unit tests achieving ~71% statement coverage.
   - Playwright E2E tests fully covering the dashboard, configuration saving, and node addition workflows.

### ⏳ Pending Features (TODO)
- **Flvx API Integration:** Call the flvx panel API to execute the actual replacement sequence when a node is confirmed dead/blocked.
- **Cloudflare API Integration:** Update the domain's A record pointing to the newly deployed standby node.
- **Frontend Node Management Polish:** Add editing and deleting capabilities to the Standby Nodes pool table in the dashboard.

## Agent Instructions
- Ensure all Go code follows idiomatic conventions and includes necessary tests.
- When writing tests, maintain the existing coverage standards (Backend `testing`, E2E `playwright`).
- The project aims for a single-binary deployment; keep dependencies minimal and standard library-focused where possible.
