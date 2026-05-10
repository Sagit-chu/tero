# Flvx Monitor - Monitor Loop & Verification Design

## 1. Overview
This specification details the design for the core background monitoring daemon of the Flvx Monitor project. It implements a two-stage verification process to detect node failures (downtime or GFW blocks) and automatically replaces them using the flvx panel API and Cloudflare DNS.

## 2. Architecture
The monitoring loop is built as a Go background service (`MonitorService`) that utilizes a standard `time.Ticker`. It queries the SQLite database for the configured monitoring interval (e.g., 5 minutes) and active node details.

## 3. Two-Stage Verification Process

### 3.1 Stage 1: Global Check (Downtime Detection)
- **Method:** TCP Ping to the node's SSH port.
- **Implementation:** Uses `net.DialTimeout("tcp", "ip:ssh_port", timeout)`.
- **Logic:** This avoids the need for root privileges (required for ICMP) and directly guarantees that the SSH daemon is responsive. If the connection fails or times out (e.g., 5 seconds), the node is immediately marked as "failed" (Dead) and triggers the replacement flow.

### 3.2 Stage 2: GFW Check (Block Detection)
- **Condition:** Only executes if Stage 1 is successful.
- **Method:** An unofficial Ping.pe scraper implemented via a `GFWChecker` interface.
- **Implementation:** 
  - The `PingPeScraper` initiates an HTTP request mimicking a browser session to Ping.pe.
  - It extracts the ping results for nodes located in China vs. Overseas.
- **Logic:** If China nodes show 100% packet loss while overseas nodes show 0% packet loss, the active node is confirmed to be blocked by the GFW. The node is marked as "failed" (Blocked) and triggers the replacement flow.
- **Resilience:** If the scraper fails to fetch results due to Ping.pe's anti-bot protections (e.g., Cloudflare turnstile), the service logs a warning and skips the block determination for the current cycle to prevent false positives.

## 4. Automated Replacement Mechanism

When a node is marked as failed, the service executes the following sequence:
1. **Fetch Standby Node:** Retrieve the next available node with `status = 'standby'` from the database. If the pool is empty, a critical alert is logged and the process halts.
2. **Flvx API Call:** Submit the new node's IP, SSH port, and SSH password to the flvx panel API.
3. **Cloudflare DNS Update:** Upon successful flvx replacement, call the Cloudflare API to update the domain's A record to the new IP.
4. **State Update:** The old node's status is updated to `failed`, and the new node is updated to `active`.
5. **Retry Logic:** If the flvx API or Cloudflare API fails, the standby node is marked as `failed` (as it couldn't be deployed), and the system immediately attempts the *next* standby node in the pool.

## 5. Security & Dependencies
- No new external dependencies are required for the ticker.
- The HTTP client for the Ping.pe scraper will include standard browser headers to minimize bot detection.
- API credentials for Flvx and Cloudflare are retrieved securely from the local SQLite `config` table.

## 6. Testing Strategy
- **Unit Tests:** Mock the `GFWChecker` interface, flvx API client, and Cloudflare API client to test the state machine and retry logic.
- **Integration Tests:** Verify the MonitorService correctly interacts with the SQLite database to fetch config and update node statuses.
