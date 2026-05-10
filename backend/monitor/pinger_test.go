package monitor

import (
	"net"
	"testing"
	"time"
)

func TestTCPPing_Success(t *testing.T) {
	// Start a dummy TCP server
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer l.Close()

	pinger := NewTCPPinger(1 * time.Second)
	alive := pinger.Ping(l.Addr().String())
	if !alive {
		t.Errorf("Expected node to be alive")
	}
}

func TestTCPPing_Fail(t *testing.T) {
	pinger := NewTCPPinger(10 * time.Millisecond)
	// Ping a non-listening port
	alive := pinger.Ping("127.0.0.1:12345")
	if alive {
		t.Errorf("Expected node to be dead")
	}
}
