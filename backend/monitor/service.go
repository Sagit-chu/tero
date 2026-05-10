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
