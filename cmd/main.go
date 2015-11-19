package main

import (
	"fmt"

	"github.com/shirou/gopsutil/disk"
	"math"
	"log"
)

func main() {
	count, _ := disk.DiskPartitions(true)

	for d := range count {
		if count[d].Device != "none"{
			usage, err := disk.DiskUsage(count[d].Mountpoint)

			if err != nil {
				fmt.Println(err)
			}

			if usage != nil {
				if !math.IsNaN(usage.UsedPercent){
					log.Printf("%s : %d on %d bytes (%d %%)", count[d].Device, usage.Used, usage.Total, usage.UsedPercent)
				}
			}
		}
	}
}
