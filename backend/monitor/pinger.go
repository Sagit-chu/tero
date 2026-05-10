package monitor

import (
	"net"
	"time"
)

type Pinger interface {
	Ping(address string) bool
}

type TCPPinger struct {
	Timeout time.Duration
}

func NewTCPPinger(timeout time.Duration) *TCPPinger {
	return &TCPPinger{Timeout: timeout}
}

func (p *TCPPinger) Ping(address string) bool {
	conn, err := net.DialTimeout("tcp", address, p.Timeout)
	if err != nil {
		return false
	}
	conn.Close()
	return true
}
