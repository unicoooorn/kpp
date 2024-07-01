package tool

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/unicoooorn/docker-monitoring-tool/internal/config"
)

type Tool struct {
	cfg    config.Config
	logger slog.Logger
}

func NewTool(cfg config.Config, logger slog.Logger) *Tool {
	return &Tool{
		cfg:    cfg,
		logger: logger,
	}
}

func (t *Tool) Run(ctx context.Context) error {
	ticker := time.NewTicker(t.cfg.MonitoringPeriod)
	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("context canceled: %w", ctx.Err())
		case <-ticker.C:
			stats := t.check()
			for container, statusOk := range stats {
				if !statusOk {
					if err := t.kill(container); err != nil {
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

func (t *Tool) check() map[string]bool {
	return nil
}

func (t *Tool) kill(container string) error {
	return nil
}
