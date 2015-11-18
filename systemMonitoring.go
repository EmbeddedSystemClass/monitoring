package monitoring
import (
	"encoding/json"
	"github.com/shirou/gopsutil/cpu"
	"time"
	"github.com/shirou/gopsutil/host"
)

//MonitoringData represents a sample of monitoring data
type monitoringData struct{
	//Unix time of this sample
	Time             int64
	//Tick for Time update (eg. CPUPercent)
	UpdateTime       time.Duration
	//Computer Name
	ComputerName     string
	OperatingSystem  string
	Platform         string
	PlatformFamily   string
	ProcessNumber	 uint64
	//All CPU values
	CPUTime          []cpu.CPUTimesStat
	GlobalCPUTime    []cpu.CPUTimesStat
	CPUInfo          []cpu.CPUInfoStat
	GlobalCPUPercent float64
	CPUPercent       []float64
	CPUCounts        int
	LogicalCPUCounts int
	//All disks values

}

//NewMonitoringData creates and fill a monitoringData
func NewMonitoringData(updateTime time.Duration) (monitorData monitoringData, err error){
	returnMonitoringData := monitoringData{}

	//Time info
	returnMonitoringData.Time = time.Now().UnixNano()
	returnMonitoringData.UpdateTime = updateTime

	//Host info
	err = initHostInfo(&returnMonitoringData)

	if err != nil{
		return returnMonitoringData,err
	}


	return returnMonitoringData, nil
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

func initHostInfo(monitorData *monitoringData) error{
	hostInfo, err := host.HostInfo()

	if err != nil{
		return err
	}

	monitorData.ComputerName = hostInfo.Hostname
	monitorData.OperatingSystem = hostInfo.OS
	monitorData.Platform = hostInfo.Platform
	monitorData.PlatformFamily = hostInfo.PlatformFamily
	monitorData.ProcessNumber = hostInfo.Procs

	return nil
}
