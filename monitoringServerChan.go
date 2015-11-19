package monitoring
import (
	"time"
)

type MonitoringServer struct{
	quitChannel chan struct{}
	DataChannel chan *monitoringData
	UpdateTime time.Duration
}

func NewMonitoringServer(updateTime time.Duration) *MonitoringServer{
	monitoringServer := new(MonitoringServer)

	monitoringServer.DataChannel = make(chan *monitoringData)
	monitoringServer.quitChannel = make(chan struct{}, 1)
	monitoringServer.UpdateTime = updateTime

	return monitoringServer
}

//Stop the monitoring server
func(self *MonitoringServer) Stop(){
	self.quitChannel <- struct{}{}
}

//Stop the monitoring server
func(self *MonitoringServer) Start(){
	go func(){
		for{
			select{
			case <-self.quitChannel:
				break
			default:
				m,err := NewMonitoringData(self.UpdateTime)

				if err == nil{
					self.DataChannel <- m
				}
			}
		}
	}()
}