package docker

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/unicoooorn/docker-monitoring-tool/internal/model"
)

type Client struct {
	cli *client.Client
}

func NewClient() (*Client, error) {
	cli, err := client.NewClientWithOpts(
		client.FromEnv,
		client.WithAPIVersionNegotiation(),
	)
	if err != nil {
		return nil, err
	}
	dc := &Client{
		cli: cli,
	}
	return dc, nil
}

func (dc *Client) Close() error {
	return dc.cli.Close()
}

func DirSize(path string) (int64, error) {
	var size int64
	err := filepath.Walk(path, func(_ string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			size += info.Size()
		}
		return err
	})
	return size, err
}

func (dc *Client) ContainersStats(ctx context.Context) (map[string]model.Stat, error) {
	stats := make(map[string]model.Stat)
	containers, err := dc.cli.ContainerList(ctx, container.ListOptions{
		All: true,
	})
	if err != nil {
		return nil, fmt.Errorf("list containers: %w", err)
	}
	for _, c := range containers {
		var totalDirSize int64 = 0
		for _, mount := range c.Mounts {
			if mount.RW {
				if mount.Source == "" {
					continue
				}
				size, err := DirSize(mount.Source[9:]) // remove /host_mnt prefix
				if err != nil {
					return nil, err
				}

				totalDirSize += size
			}
			fmt.Println()
		}
		stats[c.ID] = model.Stat{DiskUsage: totalDirSize}
	}

	return stats, nil
}

func (dc *Client) Kill(ctx context.Context, containerID string) error {
	return dc.cli.ContainerKill(ctx, containerID, "SIGKILL")
}

func (dc *Client) Pause(ctx context.Context, containerID string) error {
	return dc.cli.ContainerPause(ctx, containerID)
}

func (dc *Client) Stop(ctx context.Context, containerID string) error {
	/*
		Timeout before SIGKILL & Signal can be specified
	*/
	options := container.StopOptions{"", nil}
	return dc.cli.ContainerStop(ctx, containerID, options)
}

func (dc *Client) Restart(ctx context.Context, containerID string) error {
	options := container.StopOptions{"", nil}
	return dc.cli.ContainerRestart(ctx, containerID, options)
}

type ClientDry struct {
	cli *client.Client
}

func NewClientDry() (*ClientDry, error) {
	cli, err := client.NewClientWithOpts(
		client.FromEnv,
		client.WithAPIVersionNegotiation(),
	)
	if err != nil {
		return nil, err
	}
	dc := &ClientDry{
		cli: cli,
	}
	return dc, nil
}

func (dc *ClientDry) Close() error {
	return dc.cli.Close()
}

func (dc *ClientDry) ContainersStats(ctx context.Context) (map[string]model.Stat, error) {
	stats := make(map[string]model.Stat)
	containers, err := dc.cli.ContainerList(ctx, container.ListOptions{
		All: true,
	})
	if err != nil {
		return nil, fmt.Errorf("list containers: %w", err)
	}
	for _, c := range containers {
		var totalDirSize int64 = 0
		for _, mount := range c.Mounts {
			if mount.RW {
				if mount.Source == "" {
					continue
				}
				size, err := DirSize(mount.Source[9:]) // remove /host_mnt prefix
				if err != nil {
					return nil, err
				}

				totalDirSize += size
			}
			fmt.Println()
		}
		stats[c.ID] = model.Stat{DiskUsage: totalDirSize}
	}

	return stats, nil
}

func (dc *ClientDry) Kill(ctx context.Context, containerID string) error {
	return fmt.Errorf("[dry-run]: called method kill, container with id %s", containerID)
}

func (dc *ClientDry) Pause(ctx context.Context, containerID string) error {
	return fmt.Errorf("[dry-run]: called method pause, container with id %s", containerID)
}

func (dc *ClientDry) Stop(ctx context.Context, containerID string) error {
	return fmt.Errorf("[dry-run]: called method stop, container with id %s", containerID)
}

func (dc *ClientDry) Restart(ctx context.Context, containerID string) error {
	return fmt.Errorf("[dry-run]: called method restart, container with id %s", containerID)
}
