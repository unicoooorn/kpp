package checker

import (
	"context"

	"github.com/unicoooorn/docker-monitoring-tool/internal/config"
	"github.com/unicoooorn/docker-monitoring-tool/internal/model"
)

type DiskUsageChecker struct {
	cfg config.DiskUsageConfig
}

func NewDiskUsageChecker(cfg config.DiskUsageConfig) *DiskUsageChecker {
	return &DiskUsageChecker{cfg: cfg}
}

func (d *DiskUsageChecker) Check(_ context.Context, stat model.Stat) (ok bool) {
	return stat.DiskUsage <= d.cfg.Max
}
