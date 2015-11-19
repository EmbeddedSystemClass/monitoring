package monitoring
import (
	"time"
	"testing"
)

func TestNewMonitoringServer(t *testing.T) {
	s := NewMonitoringServer(time.Second)

	if s == nil{
		t.Error("The monitoring data have to be not nil")
	}
}

func TestServer(t *testing.T) {
	s := NewMonitoringServer(time.Second)

	if s == nil{
		t.Error("The monitoring server have to be not nil")
	}

	s.Start()

	m := <-s.DataChannel

	if m == nil{
		t.Error("The monitoring data have to be not nil")
	}
	s.Stop()
}
