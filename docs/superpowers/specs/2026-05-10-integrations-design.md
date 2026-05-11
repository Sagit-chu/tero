# Flvx Monitor - Integrations & Frontend Polish Design

## 1. Overview
This specification details the final phase of the Flvx Monitor implementation. It covers integrating the Flvx API for automated node replacement, the Cloudflare API for DNS updates, and polishing the frontend Standby Nodes table with Edit and Delete functionality via an Action Menu (Dropdown).

## 2. Frontend Node Management Polish
### 2.1 UI/UX Design
- **Action Menu:** Each row in the "Standby Nodes Pool" table will feature a "..." dropdown menu button at the end.
- **Dropdown Options:** 
  - **Edit:** Opens a modal pre-filled with the node's current IP, SSH Port, and SSH Password.
  - **Delete:** Triggers a confirmation dialog before deletion.

### 2.2 API Endpoints
The Go backend will be updated to support the new frontend actions:
- `PUT /api/nodes/{id}`: Updates an existing node's details.
- `DELETE /api/nodes/{id}`: Deletes a node from the pool.

### 2.3 Database Adjustments
- The `NodeRepository` will be expanded with `UpdateNode` and `DeleteNode` methods.

## 3. Flvx API Integration
### 3.1 Client Design
- **Component:** `backend/flvx/client.go` -> `FlvxClient`
- **Purpose:** Submit the new node's IP, SSH Port, and SSH Password to the flvx panel when a node is marked as dead/blocked.

### 3.2 Request Assumptions
- **Endpoint:** Derived from `flvx_api_url` in the database.
- **Authentication:** Standard Basic Auth using `flvx_account` and `flvx_password`.
- **Payload:** JSON payload containing `ip`, `ssh_port`, and `ssh_password`.

## 4. Cloudflare API Integration
### 4.1 Client Design
- **Component:** `backend/cloudflare/client.go` -> `CloudflareClient`
- **Purpose:** Update the domain's A record to point to the newly deployed standby node's IP.

### 4.2 Request Assumptions
- **Authentication:** Bearer token using `cf_token` from the database.
- **Workflow:** 
  1. `GET /client/v4/zones?name={domain}` -> Retrieve `zone_id`.
  2. `GET /client/v4/zones/{zone_id}/dns_records?type=A&name={domain}` -> Retrieve `record_id`.
  3. `PUT /client/v4/zones/{zone_id}/dns_records/{record_id}` -> Update the IP address.

## 5. Service Orchestration
The `MonitorService` will be updated to actually execute the replacement flow:
1. Fetch standby node.
2. Call `FlvxClient.ReplaceNode(...)`.
3. If successful, call `CloudflareClient.UpdateDNSRecord(...)`.
4. Update node statuses in the database (old -> failed, new -> active).

## 6. Testing Strategy
- **Backend:** Mock HTTP clients for Flvx and Cloudflare to verify request formatting and error handling without hitting real APIs.
- **Frontend:** Update Playwright E2E tests to cover opening the action menu, editing a node, and deleting a node.
