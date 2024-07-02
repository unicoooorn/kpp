package tool

import (
	"context"
	"log/slog"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/unicoooorn/docker-monitoring-tool/internal/checker"
	"github.com/unicoooorn/docker-monitoring-tool/internal/config"
	"github.com/unicoooorn/docker-monitoring-tool/internal/model"
	"github.com/unicoooorn/docker-monitoring-tool/internal/tool/mocks"
)

func TestTool_check(t *testing.T) {
	mockDockerClient := new(mocks.ContainerManager)
	cfg := config.Config{
		MonitoringPeriod: time.Second * 1,
		DiskUsage: config.DiskUsageConfig{
			Max: 80,
		},
	}
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	tool := New(mockDockerClient, []Checker{checker.NewDiskUsageChecker(cfg.DiskUsage)}, cfg, *logger)

	containersStats := map[string]model.Stat{
		"container1": {DiskUsage: 50},
		"container2": {DiskUsage: 90},
	}
	mockDockerClient.On("ContainersStats", mock.Anything).Return(containersStats, nil)

	ctx := context.TODO()
	statuses, err := tool.check(ctx)
	require.NoError(t, err)

	expectedStatuses := map[string]bool{
		"container1": true,
		"container2": false,
	}
	assert.Equal(t, expectedStatuses, statuses)
}

func TestTool_Run(t *testing.T) {
	mockDockerClient := new(mocks.ContainerManager)
	cfg := config.Config{
		MonitoringPeriod: time.Millisecond * 10,
		DiskUsage: config.DiskUsageConfig{
			Max: 80,
		},
	}
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	tool := New(mockDockerClient, []Checker{checker.NewDiskUsageChecker(cfg.DiskUsage)}, cfg, *logger)

	containersStats := map[string]model.Stat{
		"container1": {DiskUsage: 50},
		"container2": {DiskUsage: 90},
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	mockDockerClient.On("ContainersStats", mock.Anything).Return(containersStats, nil)
	mockDockerClient.On("Kill", ctx, "container2").Return(nil)

	go func() {
		time.Sleep(time.Millisecond * 50)
		cancel()
	}()

	err := tool.Run(ctx)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "context canceled")
}
