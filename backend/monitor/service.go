package monitor
import (
	"log"
	"strconv"
	"sync"
	"github.com/sagit-chu/flvx-monitor/backend/repository"
	"github.com/sagit-chu/flvx-monitor/backend/flvx"
)

type NodeRepository interface {
	GetNextStandby() (*repository.Node, error)
	MarkNodeFailed(ip string) error
	MarkNodeActive(id int) error
}

type ConfigRepository interface {
	GetConfig(key string) (string, error)
}

type FlvxClient interface {
	GetTunnelsRaw() ([]map[string]interface{}, error)
	UpdateTunnel(tunnelPayload map[string]interface{}) error
	GetNodes() ([]flvx.NodeItem, error)
}

type CloudflareClient interface {
	UpdateDNSRecord(domain, newIP string) error
}

type MonitorService struct {
	repo       NodeRepository
	configRepo ConfigRepository
	pinger     Pinger
	gfw        GFWChecker
	flvx       FlvxClient
	cf         CloudflareClient
	domain     string
	lastStatus string
	mu         sync.RWMutex
}

func NewMonitorService(repo NodeRepository, configRepo ConfigRepository, pinger Pinger, gfw GFWChecker, flvx FlvxClient, cf CloudflareClient, domain string) *MonitorService {
	return &MonitorService{
		repo:       repo,
		configRepo: configRepo,
		pinger:     pinger,
		gfw:        gfw,
		flvx:       flvx,
		cf:         cf,
		domain:     domain,
		lastStatus: "Unknown",
	}
}

func (s *MonitorService) GetLastStatus() string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.lastStatus
}

func (s *MonitorService) RunCycle() string {
	tunnelIDStr, err := s.configRepo.GetConfig("flvx_tunnel_id")
	if err != nil || tunnelIDStr == "" {
		s.mu.Lock()
		s.lastStatus = "Not Configured"
		s.mu.Unlock()
		return "not_configured"
	}

	tunnelID, _ := strconv.ParseInt(tunnelIDStr, 10, 64)
	if tunnelID <= 0 {
		s.mu.Lock()
		s.lastStatus = "Invalid Tunnel"
		s.mu.Unlock()
		return "invalid_tunnel"
	}

	tunnels, err := s.flvx.GetTunnelsRaw()
	if err != nil {
		log.Printf("Failed to fetch tunnels from Flvx: %v", err)
		s.mu.Lock()
		s.lastStatus = "Flvx API Error"
		s.mu.Unlock()
		return "api_error"
	}

	var targetTunnel map[string]interface{}
	for _, t := range tunnels {
		if idFloat, ok := t["id"].(float64); ok && int64(idFloat) == tunnelID {
			targetTunnel = t
			break
		}
	}

	if targetTunnel == nil {
		log.Printf("Target tunnel ID %d not found in Flvx", tunnelID)
		s.mu.Lock()
		s.lastStatus = "Tunnel Not Found"
		s.mu.Unlock()
		return "tunnel_not_found"
	}

	var currentNodeID int64
	if inNodes, ok := targetTunnel["inNodeId"].([]interface{}); ok && len(inNodes) > 0 {
		if firstNode, ok := inNodes[0].(map[string]interface{}); ok {
			if idFloat, ok := firstNode["nodeId"].(float64); ok {
				currentNodeID = int64(idFloat)
			}
		}
	}

	if currentNodeID <= 0 {
		log.Printf("Target tunnel %d has no entry node configured", tunnelID)
		s.mu.Lock()
		s.lastStatus = "No Entry Node"
		s.mu.Unlock()
		// Auto replace immediately if no entry node
		if s.performReplacement(targetTunnel) {
			return "alive"
		}
		return "no_entry_node"
	}

	nodes, err := s.flvx.GetNodes()
	if err != nil {
		log.Printf("Failed to fetch nodes from Flvx: %v", err)
		s.mu.Lock()
		s.lastStatus = "Flvx API Error"
		s.mu.Unlock()
		return "api_error"
	}

	var currentIP string
	for _, n := range nodes {
		if n.ID == currentNodeID {
			currentIP = n.ServerIP
			break
		}
	}

	if currentIP == "" {
		log.Printf("Node ID %d not found in Flvx nodes list", currentNodeID)
		s.mu.Lock()
		s.lastStatus = "Node Missing"
		s.mu.Unlock()
		// Try replacing
		if s.performReplacement(targetTunnel) {
			return "alive"
		}
		return "node_missing"
	}

	needsReplacement := false
	statusStr := "alive"
	displayStatus := "Alive"

	// We assume ssh_port is typically 22 or we ping just the IP. The pinger expects host:port.
	// Actually Ping just does a TCP connection. So port 22 is a reasonable default if not specified.
	pingTarget := currentIP + ":22"

	if !s.pinger.Ping(pingTarget) {
		log.Printf("Node %s (ID %d) is DEAD", currentIP, currentNodeID)
		needsReplacement = true
		statusStr = "dead"
		displayStatus = "Dead"
	} else {
		blocked, err := s.gfw.IsBlocked(currentIP)
		if err != nil {
			log.Printf("GFW Check failed for %s: %v", currentIP, err)
			s.mu.Lock()
			s.lastStatus = "Alive"
			s.mu.Unlock()
			return "alive" 
		}

		if blocked {
			log.Printf("Node %s (ID %d) is BLOCKED", currentIP, currentNodeID)
			needsReplacement = true
			statusStr = "blocked"
			displayStatus = "Blocked"
		}
	}

	if needsReplacement {
		s.mu.Lock()
		s.lastStatus = "Replacing"
		s.mu.Unlock()
		
		success := s.performReplacement(targetTunnel)
		s.mu.Lock()
		if success {
			s.lastStatus = "Alive" 
		} else {
			s.lastStatus = displayStatus 
		}
		s.mu.Unlock()
	} else {
		s.mu.Lock()
		s.lastStatus = displayStatus
		s.mu.Unlock()
	}

	return statusStr
}

func (s *MonitorService) performReplacement(tunnel map[string]interface{}) bool {
	for {
		node, err := s.repo.GetNextStandby()
		if err != nil || node == nil {
			log.Println("CRITICAL: No standby nodes available for replacement!")
			return false
		}

		if node.FlvxNodeID <= 0 {
			log.Printf("Local node record %d missing FlvxNodeID, marking failed", node.ID)
			s.repo.MarkNodeFailed(node.IP)
			continue
		}

		log.Printf("Attempting to replace tunnel entry with Flvx Node ID: %d (%s)", node.FlvxNodeID, node.FlvxNodeName)

		// Create a copy of the tunnel map to safely update it
		updatedTunnel := make(map[string]interface{})
		for k, v := range tunnel {
			updatedTunnel[k] = v
		}

		// Update inNodeId array
		updatedTunnel["inNodeId"] = []map[string]interface{}{
			{
				"nodeId": node.FlvxNodeID,
				"isExit": false,
			},
		}

		err = s.flvx.UpdateTunnel(updatedTunnel)
		if err != nil {
			log.Printf("Flvx API replace failed for node ID %d: %v", node.FlvxNodeID, err)
			s.repo.MarkNodeFailed(node.IP)
			continue 
		}

		// Try to fetch new IP for cloudflare update
		var newIP string
		nodes, _ := s.flvx.GetNodes()
		for _, n := range nodes {
			if n.ID == int64(node.FlvxNodeID) {
				newIP = n.ServerIP
				break
			}
		}

		if newIP != "" && s.domain != "" {
			err = s.cf.UpdateDNSRecord(s.domain, newIP)
			if err != nil {
				log.Printf("Cloudflare API update failed for %s: %v", newIP, err)
			}
		}

		s.repo.MarkNodeActive(node.ID)
		log.Printf("Successfully replaced with Flvx Node ID: %d", node.FlvxNodeID)
		return true
	}
}