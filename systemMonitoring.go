package monitoring

import (
	"encoding/json"
	"sync"
	"time"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/net"
	"github.com/shirou/gopsutil/disk"
	"math"
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
	DiskStats map[string]disk.DiskUsageStat
}

//NewMonitoringData creates and fill a monitoringData
func NewMonitoringData(updateTime time.Duration) (monitorData monitoringData, err error) {
	returnMonitoringData := monitoringData{}

	//Time info
	returnMonitoringData.Time = time.Now().UnixNano()
	returnMonitoringData.UpdatePeriod = updateTime

	//Host info
	err = initHostInfo(&returnMonitoringData)

	//Memory info
	err = initMemoryInfo(&returnMonitoringData)

	//CPU info
	err = initCPUInfo(&returnMonitoringData)

	//Network info
	netIOCounter, err := net.NetIOCounters(true)

	if netIOCounter != nil {
		returnMonitoringData.NetworkIOCounter = netIOCounter
	}

	//Disk info
	err = initDiskInfo(&returnMonitoringData)

	return returnMonitoringData, nil
}

//JSON() converts a monitoringData to JSON
func (self *monitoringData) JSON() (jsonString string, e error) {
	b, err := json.Marshal(*self)
	return string(b), err
}

//String() converts a monitoringData to String
func (self *monitoringData) String() string {
	str, _ := self.JSON()
	return str
}

func initHostInfo(monitorData *monitoringData) error {
	hostInfo, err := host.HostInfo()

	if hostInfo != nil{
		monitorData.ComputerName = hostInfo.Hostname
		monitorData.OperatingSystem = hostInfo.OS
		monitorData.Platform = hostInfo.Platform
		monitorData.PlatformFamily = hostInfo.PlatformFamily
		monitorData.ProcessNumber = hostInfo.Procs
	}

	return err
}

func initMemoryInfo(monitorData *monitoringData) error {
	memoryInfo, err := mem.VirtualMemory()
	swapInfo, err := mem.SwapMemory()

	if memoryInfo != nil && swapInfo != nil{
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

	return err
}

func initCPUInfo(monitorData *monitoringData) error {
	info, err := cpu.CPUInfo()

	if info != nil{
		monitorData.CPUInfo = info

		monitorData.CPUModelName = monitorData.CPUInfo[0].ModelName

		t, err := cpu.CPUTimes(true)

		if t != nil {
			monitorData.CPUTime = t

			globalTime, err := cpu.CPUTimes(false)

			if globalTime != nil{
				monitorData.GlobalCPUTime = globalTime

				var percentWaitGroup sync.WaitGroup

				percentWaitGroup.Add(2)

				go func(monitorData *monitoringData, wg *sync.WaitGroup){
					defer wg.Done()
					cpuPercent, _ := cpu.CPUPercent(monitorData.UpdatePeriod, true)

					monitorData.CPUPercent = cpuPercent
				}(monitorData, &percentWaitGroup)

				go func(monitorData *monitoringData, wg *sync.WaitGroup){
					defer wg.Done()
					globalCpuPercent, _ := cpu.CPUPercent(monitorData.UpdatePeriod, false)

					monitorData.GlobalCPUPercent = globalCpuPercent[0]
				}(monitorData, &percentWaitGroup)

				percentWaitGroup.Wait()

				cpuCounts, err := cpu.CPUCounts(false)

				monitorData.CPUCounts = cpuCounts

				logicalCpuCounts, err := cpu.CPUCounts(false)

				monitorData.LogicalCPUCounts = logicalCpuCounts

				return err
			}

			return err
		}

		return err
	}

	return err
}

func initDiskInfo(monitorData *monitoringData) error{
	partitions, err := disk.DiskPartitions(true)
	monitorData.DiskStats = make(map[string]disk.DiskUsageStat)

	if err == nil{
		for d := range partitions {
			if partitions[d].Device != "none"{
				usage, _ := disk.DiskUsage(partitions[d].Mountpoint)

				if usage != nil {
					if !math.IsNaN(usage.UsedPercent){
						monitorData.DiskStats[partitions[d].Device] = *usage
					}
				}
			}
		}
	}

	return err
}