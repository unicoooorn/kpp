package docker

import (
	"context"
	"encoding/json"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
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

func NewClientWithOpts(ops ...client.Opt) (*Client, error) {
	cli, err := client.NewClientWithOpts(ops...)
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

func (dc *Client) ContainerStats(ctx context.Context, id string) (uint64, error) {
	response, err := dc.cli.ContainerStatsOneShot(ctx, id)
	if err != nil {
		return 0, err
	}
	defer response.Body.Close()
	var stats container.StatsResponse
	if err := json.NewDecoder(response.Body).Decode(&stats); err != nil {
		return 0, err
	}
	//TODO возвращать больше информации
	return stats.MemoryStats.Usage, nil //кол-во в байтах
}

// default - only active containers
func (dc *Client) Containers(ctx context.Context, options container.ListOptions) ([]types.Container, error) {
	containers, err := dc.cli.ContainerList(ctx, options)
	if err != nil {
		return nil, err
	}
	return containers, nil
}
