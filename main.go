package main

import (
	"context"
	"encoding/json"
	"fmt"
	"runtime"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

func calculateCPUPercent(v *container.StatsResponse) float64 {
	var (
		cpuPercent  = 0.0
		cpuDelta    = float64(v.CPUStats.CPUUsage.TotalUsage) - float64(v.PreCPUStats.CPUUsage.TotalUsage)
		systemDelta = float64(v.CPUStats.SystemUsage) - float64(v.PreCPUStats.SystemUsage)
	)

	if systemDelta > 0.0 && cpuDelta > 0.0 {
		cpuPercent = (cpuDelta / systemDelta) * float64(runtime.NumCPU()) * 100.0
	}
	return cpuPercent
}

func calculateMemoryUsagePercent(v *container.StatsResponse) float64 {
	return float64(v.MemoryStats.Usage-v.MemoryStats.Stats["cache"]) / float64(v.MemoryStats.Limit) * 100.0
}

// TODO определиться с выводом
func monitorStats(ctx context.Context, cli *client.Client, ctr *types.Container, d time.Duration) {
	response, _ := cli.ContainerStats(ctx, ctr.ID, true)
	var stats container.StatsResponse
	defer response.Body.Close()

	for {
		if err := json.NewDecoder(response.Body).Decode(&stats); err != nil {
			response.Body.Close()
			//panic(err)
			fmt.Print(err.Error())
			return
		}
		fmt.Printf("NAME\t%s\nCPU%%\t%.2f %%\nMEM%%\t %.2f %%\nMEM/LIM\t%v/%v Mb\nBLKIO\t%v/%v Mb\nRX/TX\t%v/%v Kb\n\n",
			ctr.Image,
			calculateCPUPercent(&stats),
			calculateMemoryUsagePercent(&stats),
			(stats.MemoryStats.Usage-stats.MemoryStats.Stats["cache"])/1024/1024,
			stats.MemoryStats.Limit/1024/1024,
			stats.BlkioStats.IoServiceBytesRecursive[0].Value/1024/1024,
			stats.BlkioStats.IoServiceBytesRecursive[1].Value/1024/1024,
			stats.Networks["eth0"].RxBytes/1024,
			stats.Networks["eth0"].TxBytes/1024)
		time.Sleep(d)
	}
}

func main() {
	cli, err := client.NewClientWithOpts(
		client.FromEnv,
		client.WithAPIVersionNegotiation(),
	)
	if err != nil {
		panic(err)
	}
	defer cli.Close()
	ctx := context.Background()
	containers, err := cli.ContainerList(ctx, container.ListOptions{})
	if err != nil {
		panic(err)
	}

	//TODO удалить или отредачить
	for _, ctr := range containers {
		go monitorStats(ctx, cli, &ctr, 1*time.Second)
	}
	var c string
	fmt.Scan(&c)
}
