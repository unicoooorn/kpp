package docker

import (
	"context"
	"errors"
	"fmt"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/unicoooorn/docker-monitoring-tool/internal/model"
)

var (
	ErrNotFound = errors.New("container not found")
)

type DryRunClient struct {
	rawClient     *client.Client
	wrappedClient *Client
}

func NewDryRunClient() (*DryRunClient, error) {
	cli, err := client.NewClientWithOpts(
		client.FromEnv,
		client.WithAPIVersionNegotiation(),
	)
	if err != nil {
		return nil, err
	}

	wrappedClient, err := NewClient()
	if err != nil {
		return nil, err
	}

	return &DryRunClient{
		rawClient:     cli,
		wrappedClient: wrappedClient,
	}, nil
}

func (drc *DryRunClient) Close() error {
	return errors.Join(drc.rawClient.Close(), drc.wrappedClient.Close())
}

func (drc *DryRunClient) ContainersStats(ctx context.Context) (map[string]model.Stat, error) {
	return drc.wrappedClient.ContainersStats(ctx)
}

func (drc *DryRunClient) Kill(ctx context.Context, containerID string) error {
	ok, err := drc.isThereAlive(ctx, containerID)
	if err != nil {
		return fmt.Errorf("check if container is alive: %w", err)
	}

	if !ok {
		return ErrNotFound
	}

	fmt.Println("[dry-run] killing container:", containerID)

	return nil
}

func (drc *DryRunClient) Pause(ctx context.Context, containerID string) error {
	ok, err := drc.isThereAlive(ctx, containerID)
	if err != nil {
		return fmt.Errorf("check if container is alive: %w", err)
	}

	if !ok {
		return ErrNotFound
	}

	fmt.Println("[dry-run] pausing container:", containerID)

	return nil
}

func (drc *DryRunClient) Stop(ctx context.Context, containerID string) error {
	ok, err := drc.isThereAlive(ctx, containerID)
	if err != nil {
		return fmt.Errorf("check if container is alive: %w", err)
	}

	if !ok {
		return ErrNotFound
	}

	fmt.Println("[dry-run] stopping container:", containerID)

	return nil
}

func (drc *DryRunClient) Restart(ctx context.Context, containerID string) error {
	ok, err := drc.isThereAlive(ctx, containerID)
	if err != nil {
		return fmt.Errorf("check if container is alive: %w", err)
	}

	if !ok {
		return ErrNotFound
	}

	fmt.Println("[dry-run] restarting container:", containerID)

	return nil
}

func (drc *DryRunClient) isThereAlive(ctx context.Context, containerID string) (bool, error) {
	list, err := drc.rawClient.ContainerList(ctx, container.ListOptions{
		All: true,
	})
	if err != nil {
		return false, fmt.Errorf("list containers: %w", err)
	}

	found := false
	for _, c := range list {
		if c.ID == containerID {
			found = true
		}
	}
	if !found {
		return false, nil
	}
	return true, nil
}
