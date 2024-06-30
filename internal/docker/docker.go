package docker

import (
	"context"
	"encoding/json"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

type DockerClient struct {
	cli *client.Client
	ctx context.Context
}

func NewDockerClient() (*DockerClient, error) {
	cli, err := client.NewClientWithOpts(
		client.FromEnv,
		client.WithAPIVersionNegotiation(),
	)
	if err != nil {
		return nil, err
	}
	dc := &DockerClient{
		cli: cli,
		ctx: context.Background(), //хз что это но это нужно функциям апи
	}
	return dc, nil
}

func NewDockerClientWithOpts(ops ...client.Opt) (*DockerClient, error) {
	cli, err := client.NewClientWithOpts(ops...)
	if err != nil {
		return nil, err
	}
	dc := &DockerClient{
		cli: cli,
		ctx: context.Background(),
	}
	return dc, nil
}

func (dc *DockerClient) Close() error {
	dc.cli.Close()
	return nil
}

func (dc *DockerClient) GetContainerStats(id string) uint64 {
	response, err := dc.cli.ContainerStatsOneShot(dc.ctx, id)
	if err != nil {
		panic(err)
	}
	defer response.Body.Close()
	var stats container.StatsResponse
	if err := json.NewDecoder(response.Body).Decode(&stats); err != nil {
		panic(err)
	}

	//TODO возвращать больше информации
	return stats.MemoryStats.Usage //Bytes
}

func (dc *DockerClient) MonitorContainerStats(id string, deltaTime time.Duration, c chan uint64) error {
	defer close(c)
	response, err := dc.cli.ContainerStats(dc.ctx, id, true)
	if err != nil {
		panic(err)
	}
	var stats container.StatsResponse
	defer response.Body.Close()
	for {
		if err := json.NewDecoder(response.Body).Decode(&stats); err != nil {
			return err
		}
		c <- stats.MemoryStats.Usage
		time.Sleep(deltaTime)
	}
}

// default - only active containers
func (dc *DockerClient) GetContainers(options container.ListOptions) ([]types.Container, error) {
	containers, err := dc.cli.ContainerList(dc.ctx, options)
	if err != nil {
		return nil, err
	}
	return containers, nil
}

// func calculateCPUPercent(v *container.StatsResponse) float64 {
// 	var (
// 		cpuPercent  = 0.0
// 		cpuDelta    = float64(v.CPUStats.CPUUsage.TotalUsage) - float64(v.PreCPUStats.CPUUsage.TotalUsage)
// 		systemDelta = float64(v.CPUStats.SystemUsage) - float64(v.PreCPUStats.SystemUsage)
// 	)
// 	if systemDelta > 0.0 && cpuDelta > 0.0 {
// 		cpuPercent = (cpuDelta / systemDelta) * float64(runtime.NumCPU()) * 100.0
// 	}
// 	return cpuPercent
// }
// func calculateMemoryUsagePercent(v *container.StatsResponse) float64 {
// 	return float64(v.MemoryStats.Usage-v.MemoryStats.Stats["cache"]) / float64(v.MemoryStats.Limit) * 100.0
// }
