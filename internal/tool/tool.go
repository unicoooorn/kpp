package tool

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/unicoooorn/docker-monitoring-tool/internal/config"
)

type Tool struct {
	cfg              config.Config
	logger           slog.Logger
	containerManager ContainerManager
	checkers         []Checker
}

func New(statsRetriever ContainerManager, checkers []Checker, cfg config.Config, logger slog.Logger) *Tool {
	return &Tool{
		checkers:         checkers,
		containerManager: statsRetriever,
		cfg:              cfg,
		logger:           logger,
	}
}

func (t *Tool) Run(ctx context.Context) error {
	ticker := time.NewTicker(t.cfg.MonitoringPeriod)
	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("context canceled: %w", ctx.Err())
		case <-ticker.C:
			stats, err := t.check(ctx)
			if err != nil {
				slog.Error("check stats", slog.Any("err", err))
			}

			for container, statusOk := range stats {
				if !statusOk {
					if err := t.containerManager.Kill(ctx, container); err != nil {
						t.logger.Error("kill container",
							slog.String("container", container),
							slog.Any("err", err),
						)
					}
				}
			}
		}
	}
}

func (t *Tool) check(ctx context.Context) (map[string]bool, error) {
	statuses := make(map[string]bool)

	containersStats, err := t.containerManager.ContainersStats(ctx)
	if err != nil {
		return nil, err
	}

	for id, stat := range containersStats {
		ok := true
		for _, c := range t.checkers {
			if !c.Check(ctx, stat) {
				ok = false
				break
			}
		}
		statuses[id] = ok
	}

	return statuses, nil
}
