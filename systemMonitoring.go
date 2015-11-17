package monitoring
import (
	"encoding/json"
	"github.com/shirou/gopsutil/cpu"
	"time"
)

//MonitoringData represents a sample of monitoring data
type monitoringData struct{
	//Unix time of this sample
	Time int64
	//Tick for Time update (eg. CPUPercent)
	UpdateTime time.Duration
	//Computer Name
	ComputerName string
	//All CPU values
	CPUTime []cpu.CPUTimesStat
	GlobalCPUTime []cpu.CPUTimesStat
	CPUInfo []cpu.CPUInfoStat
	GlobalCPUPercent float64
	CPUPercent []float64
	CPUCounts int
	LogicalCPUCounts int
	//All disks values

}

//NewMonitoringData creates and fill a monitoringData
func NewMonitoringData(updateTime time.Duration) (monitorData monitoringData, err error){



	return nil, nil
}

//JSON() converts a monitoringData to JSON
func (self *monitoringData) JSON() (jsonString string, e error){
	b, err := json.Marshal(*self)

	if err != nil{
		return nil, err
	}

	return string(b), err
}

//String() converts a monitoringData to String
func (self *monitoringData) String() string{
	str, _ := self.JSON()
	return str
}