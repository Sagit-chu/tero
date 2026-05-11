package monitor

import (
	"log"
	"github.com/sagit-chu/flvx-monitor/backend/repository"
)

type NodeRepository interface {
	GetActiveNode() (ip, port string, err error)
	MarkNodeFailed(ip string) error
	GetNextStandby() (*repository.Node, error)
	MarkNodeActive(id int) error
}

type FlvxClient interface {
	ReplaceNode(ip, port, password string) error
}

type CloudflareClient interface {
	UpdateDNSRecord(domain, newIP string) error
}

type MonitorService struct {
	repo   NodeRepository
	pinger Pinger
	gfw    GFWChecker
	flvx   FlvxClient
	cf     CloudflareClient
	domain string
}

func NewMonitorService(repo NodeRepository, pinger Pinger, gfw GFWChecker, flvx FlvxClient, cf CloudflareClient, domain string) *MonitorService {
	return &MonitorService{
		repo:   repo,
		pinger: pinger,
		gfw:    gfw,
		flvx:   flvx,
		cf:     cf,
		domain: domain,
	}
}

// RunCycle performs one iteration of the monitoring logic
func (s *MonitorService) RunCycle() string {
	ip, port, err := s.repo.GetActiveNode()
	if err != nil || ip == "" {
		return "no_active_node"
	}

	needsReplacement := false
	statusStr := "alive"

	// Stage 1
	if !s.pinger.Ping(ip + ":" + port) {
		log.Printf("Node %s is DEAD", ip)
		s.repo.MarkNodeFailed(ip)
		needsReplacement = true
		statusStr = "dead"
	} else {
		// Stage 2
		blocked, err := s.gfw.IsBlocked(ip)
		if err != nil {
			log.Printf("GFW Check failed for %s: %v", ip, err)
			return "alive" // skip block assumption on error
		}

		if blocked {
			log.Printf("Node %s is BLOCKED", ip)
			s.repo.MarkNodeFailed(ip)
			needsReplacement = true
			statusStr = "blocked"
		}
	}

	if needsReplacement {
		s.performReplacement()
	}

	return statusStr
}

func (s *MonitorService) performReplacement() {
	for {
		node, err := s.repo.GetNextStandby()
		if err != nil || node == nil {
			log.Println("CRITICAL: No standby nodes available for replacement!")
			return
		}

		log.Printf("Attempting to replace with standby node: %s", node.IP)
		
		err = s.flvx.ReplaceNode(node.IP, node.SSHPort, node.SSHPassword)
		if err != nil {
			log.Printf("Flvx API replace failed for %s: %v", node.IP, err)
			s.repo.MarkNodeFailed(node.IP)
			continue // try next node
		}

		err = s.cf.UpdateDNSRecord(s.domain, node.IP)
		if err != nil {
			log.Printf("Cloudflare API update failed for %s: %v", node.IP, err)
			s.repo.MarkNodeFailed(node.IP)
			continue // try next node
		}

		s.repo.MarkNodeActive(node.ID)
		log.Printf("Successfully replaced with node: %s", node.IP)
		break
	}
}

