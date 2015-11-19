package monitoring

import (
	"testing"
	"time"
	"encoding/json"
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

	var obj monitoringData
	j, err := m.JSON()

	if err != nil{
		t.Error("Error in JSON() should be nil : ", err)
	}

	err = json.Unmarshal([]byte(j), &obj)

	if err != nil{
		t.Error("Malformated JSON", err)
	}
}

func TestMonitoringDataString(t *testing.T) {
	m, err := NewMonitoringData(time.Second)

	if err != nil {
		t.Error("Error in NewMonitoringData should be nil : ", err)
	}

	var obj monitoringData
	j, err := m.JSON()

	if err != nil{
		t.Error("Error in JSON() should be nil : ", err)
	}

	err = json.Unmarshal([]byte(j), &obj)

	if err != nil {
		t.Error("Malformated JSON", err)
	}
}
