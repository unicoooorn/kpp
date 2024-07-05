package docker

import (
	"context"
	"fmt"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/unicoooorn/kpp/internal/model"
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
	return &Client{
		cli: cli,
	}, nil
}

func (dc *Client) Close() error {
	return dc.cli.Close()
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
		stat := model.Stat{}
		for _, mount := range c.Mounts {
			if mount.RW {
				if mount.Source == "" {
					continue
				}
				size, err := dirSize(mount.Source[9:]) // remove /host_mnt prefix
				if err != nil {
					return nil, err
				}
				stat.Volumes = append(stat.Volumes, mount.Source[9:])
				stat.DiskUsage += size
			}
			fmt.Println()
		}
		stats[c.ID] = stat
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
	options := container.StopOptions{Signal: "", Timeout: nil}
	return dc.cli.ContainerStop(ctx, containerID, options)
}

func (dc *Client) Restart(ctx context.Context, containerID string) error {
	options := container.StopOptions{Signal: "", Timeout: nil}
	return dc.cli.ContainerRestart(ctx, containerID, options)
}
