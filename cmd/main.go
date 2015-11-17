package main
import (
	"fmt"
	"github.com/shirou/gopsutil/disk"
)

func main(){
	count, _ := disk.DiskPartitions(true)

	fmt.Println(len(count))

	for d := range count{
		usage, err := disk.DiskUsage(count[d].Mountpoint)

		if err != nil{
			fmt.Println(err)
		}

		if usage != nil{
			fmt.Println("Disk usage : ", count[d].Mountpoint)
			fmt.Println("Disk ", usage.UsedPercent)
		}
	}
}
