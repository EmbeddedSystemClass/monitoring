package monitoring

import (
	"encoding/json"
	"sync"
	"time"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/net"
)

//MonitoringData represents a sample of monitoring data
type monitoringData struct {
	//Unix nano time of this sample
	Time int64
	//Tick for Time update (eg. CPUPercent)
	UpdatePeriod time.Duration
	//Computer Name
	ComputerName    string
	OperatingSystem string
	Platform        string
	PlatformFamily  string
	ProcessNumber   uint64
	//All the memory values
	TotalMemory     uint64
	AvailableMemory uint64
	UsedMemory      uint64
	UsedPercent     float64
	FreeMemory      uint64
	SharedMemory    uint64
	TotalSwapMemory uint64
	UsedSwap        uint64
	FreeSwap        uint64
	SwapUsedPercent float64
	//All CPU values
	CPUTime          []cpu.CPUTimesStat
	GlobalCPUTime    []cpu.CPUTimesStat
	CPUInfo          []cpu.CPUInfoStat
	GlobalCPUPercent float64
	CPUPercent       []float64
	CPUCounts        int
	LogicalCPUCounts int
	CPUModelName     string
	//All network values
	NetworkIOCounter []net.NetIOCountersStat
	//All disks values
}

//NewMonitoringData creates and fill a monitoringData
func NewMonitoringData(updateTime time.Duration) (monitorData monitoringData, err error) {
	returnMonitoringData := monitoringData{}

	//Time info
	returnMonitoringData.Time = time.Now().UnixNano()
	returnMonitoringData.UpdatePeriod = updateTime

	//Host info
	err = initHostInfo(&returnMonitoringData)

	if err != nil {
		return returnMonitoringData, err
	}

	//Memory info
	err = initMemoryInfo(&returnMonitoringData)

	if err != nil {
		return returnMonitoringData, err
	}

	//CPU info
	err = initCPUInfo(&returnMonitoringData)

	if err != nil {
		return returnMonitoringData, err
	}

	//Network info
	netIOCounter, err := net.NetIOCounters(true)

	if err != nil {
		return returnMonitoringData, err
	}

	returnMonitoringData.NetworkIOCounter = netIOCounter

	return returnMonitoringData, nil
}

//JSON() converts a monitoringData to JSON
func (self *monitoringData) JSON() (jsonString string, e error) {
	b, err := json.Marshal(*self)

	if err != nil {
		return nil, err
	}

	return string(b), err
}

//String() converts a monitoringData to String
func (self *monitoringData) String() string {
	str, _ := self.JSON()
	return str
}

func initHostInfo(monitorData *monitoringData) error {
	hostInfo, err := host.HostInfo()

	if err != nil {
		return err
	}

	monitorData.ComputerName = hostInfo.Hostname
	monitorData.OperatingSystem = hostInfo.OS
	monitorData.Platform = hostInfo.Platform
	monitorData.PlatformFamily = hostInfo.PlatformFamily
	monitorData.ProcessNumber = hostInfo.Procs

	return nil
}

func initMemoryInfo(monitorData *monitoringData) error {
	memoryInfo, err := mem.VirtualMemory()
	if err != nil {
		return err
	}

	swapInfo, err := mem.SwapMemory()
	if err != nil {
		return err
	}

	monitorData.AvailableMemory = memoryInfo.Available
	monitorData.TotalMemory = memoryInfo.Total
	monitorData.UsedMemory = memoryInfo.Used
	monitorData.UsedPercent = memoryInfo.UsedPercent
	monitorData.FreeMemory = memoryInfo.Free
	monitorData.SharedMemory = memoryInfo.Shared
	monitorData.TotalSwapMemory = swapInfo.Total
	monitorData.UsedSwap = swapInfo.Used
	monitorData.FreeSwap = swapInfo.Free
	monitorData.SwapUsedPercent = swapInfo.UsedPercent
}

func initCPUInfo(monitorData *monitoringData) error {
	info, err := cpu.CPUInfo()

	if err != nil {
		return err
	}

	monitorData.CPUInfo = info

	monitorData.CPUModelName = monitorData.CPUInfo[0].ModelName

	t, err := cpu.CPUTimes(true)

	if err != nil {
		return err
	}

	monitorData.CPUTime = t

	globalTime, err := cpu.CPUTimes(false)

	if err != nil {
		return err
	}

	monitorData.GlobalCPUTime = globalTime[0]

	var percentWaitGroup sync.WaitGroup

	percentWaitGroup.Add(2)

	go func(monitorData *monitoringData, wg *sync.WaitGroup) {
		defer wg.Done()
		cpuPercent, err := cpu.CPUPercent(monitorData.UpdatePeriod, true)

		if err != nil {
			return err
		}

		monitorData.CPUPercent = cpuPercent
	}(monitorData)

	go func(monitorData *monitoringData, wg *sync.WaitGroup) {
		defer wg.Done()
		globalCpuPercent, err := cpu.CPUPercent(monitorData.UpdatePeriod, false)

		if err != nil {
			return err
		}

		monitorData.GlobalCPUPercent = globalCpuPercent[0]
	}(monitorData)

	percentWaitGroup.Wait()

	cpuCounts, err := cpu.CPUCounts(false)

	if err != nil {
		return err
	}

	monitorData.CPUCounts = cpuCounts

	logicalCpuCounts, err := cpu.CPUCounts(false)

	if err != nil {
		return err
	}

	monitorData.LogicalCPUCounts = logicalCpuCounts
}
