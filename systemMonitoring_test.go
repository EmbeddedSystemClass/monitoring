package monitoring

import (
	"testing"
	"time"
	"fmt"
)

func TestNewMonitoringData(t *testing.T) {
	_, err := NewMonitoringData(time.Second)

	if err != nil {
		t.Error("Error in NewMonitoringData should be nil : ", err)
	}
}

func TestMonitoringDataJSON(t *testing.T) {
	m, err := NewMonitoringData(time.Second)

	if err != nil {
		t.Error("Error in NewMonitoringData should be nil : ", err)
	}

	fmt.Println(m.String())
}

func TestMonitoringDataString(t *testing.T) {
	m, err := NewMonitoringData(time.Second)

	if err != nil {
		t.Error("Error in NewMonitoringData should be nil : ", err)
	}

	fmt.Println(m.JSON())
}
