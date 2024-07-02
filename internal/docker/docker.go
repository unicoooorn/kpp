package docker

import (
	"context"
	"fmt"

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

func (dc *Client) ContainersStats(ctx context.Context) (map[string]model.Stat, error) {
	stats := make(map[string]model.Stat)
	containers, err := dc.cli.ContainerList(ctx, container.ListOptions{
		All: true,
	})
	if err != nil {
		return nil, fmt.Errorf("list containers: %w", err)
	}

	for _, c := range containers {
		stats[c.ID] = model.Stat{DiskUsage: c.SizeRw + c.SizeRootFs}
	}

	return stats, nil
}

func (dc *Client) Kill(ctx context.Context, containerID string) error {
	return dc.ContainerKill(ctx, containerID, "SIGKILL")
}

func (dc *Client) Pause(ctx context.Context, containerID string) error {
	return dc.ContainerPause(ctx, containerID)
}

func (dc *Client) Stop(ctx, context.Context, containerID string) error {
	/*
	Timeout before SIGKILL & Signal can be specified
	*/
	options := container.StopOptions {nil, ""}
	return dc.ContainerStop(ctx, containerID, options)
}
