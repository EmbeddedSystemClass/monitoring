package monitoring
import (
	"testing"
	"time"
)

func TestNewMonitoringData(t *testing.T) {
	_ , err := NewMonitoringData(time.Second)

	if err != nil{
		t.Error("Error in NewMonitoringData should be nil : ", err)
	}
}

func TestMonitoringDataJSON(t *testing.T) {

}

func TestMonitoringDataString(t *testing.T) {

}

